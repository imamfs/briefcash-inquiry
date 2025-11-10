package service

import (
	"briefcash-inquiry/internal/authorization"
	"briefcash-inquiry/internal/dto"
	"briefcash-inquiry/internal/entity"
	"briefcash-inquiry/internal/helper/errorhelper"
	"briefcash-inquiry/internal/helper/httphelper"
	"briefcash-inquiry/internal/helper/loghelper"
	"briefcash-inquiry/internal/helper/routinghelper"
	"briefcash-inquiry/internal/mapper"
	"briefcash-inquiry/internal/repository"
	"fmt"
	"net/http"
	"sync"
	"time"

	"context"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type InquiryService interface {
	InquiryAccount(ctx context.Context, dto dto.InquiryRequest, partnerRefNo string) (*dto.InquiryResponse, error)
}

type inquiryService struct {
	repo     repository.InquiryRepository
	tokenSvc TokenService
	bankRepo BankPartner
	db       *gorm.DB
}

type inquiryContext struct {
	Request      dto.InquiryRequest
	BankConfig   *entity.BankConfig
	PartnerRefNo string
	Context      context.Context
}

func NewInquiryService(repo repository.InquiryRepository, tokenSvc TokenService, bankRepo BankPartner, db *gorm.DB) InquiryService {
	return &inquiryService{repo, tokenSvc, bankRepo, db}
}

func (is *inquiryService) InquiryAccount(ctx context.Context, req dto.InquiryRequest, externalId string) (*dto.InquiryResponse, error) {
	log := loghelper.Logger.WithFields(logrus.Fields{
		"service":     "inquiry_service",
		"operation":   "inquiry_account",
		"bank_code":   req.BankCode,
		"external_id": externalId,
	})

	log.WithField("step", "get_bank_route").Info("Check available bank routes")
	bankConfig := is.bankRepo.GetBankConfig(req.BankCode)
	bankRoute := routinghelper.NewBankRouteRequest(req, &bankConfig, externalId)
	log.Infof("Bank available, will send request from bank %s", bankConfig.BankName)

	data := inquiryContext{Request: req, BankConfig: &bankConfig, PartnerRefNo: externalId, Context: ctx}
	var accessToken string

	log.WithField("step", "check_active_access_token").Info("Checking active access token in redis and database")
	accessToken, err := is.tokenSvc.GetActiveAccessToken(ctx, bankConfig.BankName)

	if err != nil {
		log.WithField("step", "get_new_access_token").Info("Get new access token from bank")
		respToken, err := authorization.GetAccessToken(&bankConfig, log)
		if err != nil {
			log.WithField("step", "get_new_access_token").WithError(err).Error("Failed to get new access token from bank")
			return is.handleInquiryResponse(&data, nil, 0, err, log)
		}

		token := &entity.AccessToken{
			AccessToken: respToken.AccessToken,
			TokenType:   "Bearer",
			ExpiresIn:   respToken.ExpiresIn,
			ExpiresDate: time.Now().Add(time.Duration(respToken.ExpiresIn-30) * time.Second),
		}

		log.WithField("step", "get_new_access_token").Info("Access token retrieved, saving data to database and redis, running on goroutine")
		errorChn := make(chan error, 2)
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			err := is.tokenSvc.SaveAccessTokenDB(ctx, token)
			if err != nil {
				errorChn <- fmt.Errorf("failed to save access token to database: %v", err)
			}
		}()

		go func() {
			defer wg.Done()
			err := is.tokenSvc.SaveAccessTokenRedis(ctx, bankConfig.BankName, token)
			if err != nil {
				errorChn <- fmt.Errorf("failed to save access token to redis: %v", err)
			}
		}()

		wg.Wait()
		close(errorChn)

		for er := range errorChn {
			log.WithField("step", "get_new_access_token").WithError(er).Warn("Failed to save new access token (Redis/DB)")
		}

		log.WithField("step", "get_new_access_token").Info("Saving data is done")
		accessToken = respToken.AccessToken

		if accessToken == "" {
			return is.handleInquiryResponse(&data, nil, 0, fmt.Errorf("missing access token after refresh"), log)
		}
	}

	log.WithField("step", "set_param_request").Info("Setting up url, payload, and http header parameters")
	url := bankRoute.GetUrl()
	payload := bankRoute.BuildBodyRequest()
	headers := bankRoute.GetHeaders(accessToken, externalId, &bankConfig, payload)

	log.WithField("step", "send_request").Info("Send request inquiry to destination bank")
	client := httphelper.NewHttpClientHelper(10 * time.Second)
	resp, httpStatus, err := client.SendRequest("POST", url, payload, headers)
	return is.handleInquiryResponse(&data, resp, httpStatus, err, log)
}

func (is *inquiryService) handleInquiryResponse(data *inquiryContext, respData []byte, httpStatus int, er error, log *logrus.Entry) (*dto.InquiryResponse, error) {
	log.WithField("step", "handle_transport_error").Info("Check transport data from bank")
	if errResp, e := is.handleTransportError(er, respData, log); errResp != nil {
		log.WithField("step", "handle_transport_error").WithError(e).Error("Failed to send request to bank due connection issue")
		return errResp, e
	}

	log.WithField("step", "parse_response").Info("Parsing and validating response data from bank")
	mapData, err := is.parseBankResponse(respData, data.BankConfig, httpStatus)
	if err != nil {
		log.WithField("step", "parse_response").WithError(err).Error("Failed to parsing bank response, please check response format")
		return errorhelper.BuildErrorResponse(errorhelper.ErrorDetail{
			Code:       "BANK_FORMAT_ERROR",
			Message:    "Invalid response format from bank",
			LogMessage: err.Error(),
			Source:     errorhelper.SourceInternal,
		}, "", err)
	}

	log.WithField("step", "handle_bank_error").Info("Evaluating HTTP response status from bank")
	if httpStatus != http.StatusOK {
		return is.handleBankError(httpStatus, mapData, log)
	}

	log.WithField("step", "persist_data").Info("Save inquiry response to database")
	if err := is.persistInquiry(data, mapData); err != nil {
		log.WithField("step", "persist_data").WithError(err).Error("Failed to save inquiry to database")
		return errorhelper.BuildErrorResponse(errorhelper.ErrorDetail{
			Code:       "INTERNAL_SERVER_ERROR",
			Message:    "Internal server error occured",
			LogMessage: err.Error(),
			Source:     errorhelper.SourceInternal,
		}, "", err)
	}

	log.WithField("step", "build_response").Info("Map inquiry response to client")
	return is.buildSuccessResponse(data, mapData)
}

func (is *inquiryService) handleTransportError(err error, respData []byte, log *logrus.Entry) (*dto.InquiryResponse, error) {
	if err != nil {
		detail := errorhelper.ErrorDetail{
			Code:       "INTERNAL_CONNECTION_ERROR",
			Message:    "Fail to send request to bank",
			LogMessage: "Failed to send request, please check connection",
			Source:     errorhelper.SourceInternal,
		}
		log.WithField("step", "handle_transport_error").WithError(err).Error("Return error from the bank")
		return errorhelper.BuildErrorResponse(detail, "", err)
	}

	if respData == nil {
		detail := errorhelper.ErrorDetail{
			Code:       "BANK_NO_RESPONSE",
			Message:    "No response from bank",
			LogMessage: "Empty or nil response from bank",
			Source:     errorhelper.SourceBank,
		}
		log.WithField("step", "handle_transport_error").Error("No response from the bank")
		return errorhelper.BuildErrorResponse(detail, "", nil)
	}
	return nil, nil
}

func (is *inquiryService) parseBankResponse(bankResp []byte, bankCfg *entity.BankConfig, httpStatus int) (mapper.BankResponseData, error) {
	routeResp := routinghelper.NewBankRouteResponse(bankCfg, httpStatus)
	mapData := routeResp.MapResponse(bankResp)
	return mapData, nil
}

func (is *inquiryService) handleBankError(httpStatus int, mapData mapper.BankResponseData, log *logrus.Entry) (*dto.InquiryResponse, error) {
	detail, ok := errorhelper.ErrorMap[httpStatus]
	if !ok {
		detail = errorhelper.DefaultBankError
	}
	log.WithField("step", "handle_bank_error").Errorf("bank error status: %d, with message: %s", httpStatus, mapData.ResponseMessage)
	return errorhelper.BuildErrorResponse(detail, mapData.ResponseMessage, fmt.Errorf("bank error status: %d", httpStatus))
}

func (is *inquiryService) saveInquiry(ctx context.Context, inq *entity.Inquiry) error {
	return is.db.Transaction(func(trx *gorm.DB) error {
		repo := is.repo.WithTransaction(trx)
		return repo.SaveInquiry(ctx, inq)
	})
}

func (is *inquiryService) persistInquiry(data *inquiryContext, mapData mapper.BankResponseData) error {
	inquiry := &entity.Inquiry{
		MerchantCode:           data.Request.CompanyId,
		PartnerReferenceNo:     data.PartnerRefNo,
		BeneficiaryAccount:     data.Request.BeneficiaryAccount,
		BeneficiaryBankCode:    data.Request.BankCode,
		BeneficiaryAccountName: mapData.AccountName,
		InquiryDate:            time.Now(),
		Status:                 "SUCCESS",
	}
	return is.saveInquiry(data.Context, inquiry)
}

func (is *inquiryService) buildSuccessResponse(data *inquiryContext, mapData mapper.BankResponseData) (*dto.InquiryResponse, error) {
	return &dto.InquiryResponse{
		Status:  true,
		Code:    "SUCCESS",
		Message: "Inquiry successful",
		Source:  errorhelper.SourceBank,
		Data: dto.InquiryData{
			BeneficiaryAccount: data.Request.BeneficiaryAccount,
			BankCode:           data.Request.BankCode,
			BeneficiaryName:    mapData.AccountName,
		},
	}, nil
}
