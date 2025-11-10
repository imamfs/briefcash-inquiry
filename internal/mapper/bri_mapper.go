package mapper

import (
	"briefcash-inquiry/internal/authorization"
	"briefcash-inquiry/internal/dto"
	"briefcash-inquiry/internal/entity"
	"briefcash-inquiry/internal/helper/jsonhelper"
	"briefcash-inquiry/internal/helper/timehelper"
	"encoding/json"
	"net/http"
	"time"
)

type briClientRequest struct {
	cfg *entity.BankConfig
	req dto.InquiryRequest
}

type briClientResponse struct {
	cfg        *entity.BankConfig
	httpStatus int
}

func NewBriClientRequest(cfg *entity.BankConfig, req dto.InquiryRequest) *briClientRequest {
	return &briClientRequest{
		cfg: cfg, req: req,
	}
}

func NewBriClientResponse(cfg *entity.BankConfig, httpStatus int) *briClientResponse {
	return &briClientResponse{
		cfg:        cfg,
		httpStatus: httpStatus,
	}
}

func (bri *briClientRequest) BuildBodyRequest() []byte {
	if bri.req.BankCode == "002" {
		payload := dto.BRIInternalInquiryRequest{
			BeneficiaryAccountNo: bri.req.BeneficiaryAccount,
			AdditionalInfo: map[string]string{
				"channel":  "",
				"deviceId": "",
			},
		}
		return jsonhelper.WriteToJson(payload)
	} else {
		payload := dto.BRIExternalInquiryRequest{
			BeneficiaryBankCode:  bri.req.BankCode,
			BeneficiaryAccountNo: bri.req.BeneficiaryAccount,
			AdditionalInfo: map[string]string{
				"serviceCode": func() string {
					if bri.req.Type == "bifast" {
						return "81"
					} else {
						return "16"
					}
				}(),
				"deviceId": "",
				"channel":  "",
			},
		}
		return jsonhelper.WriteToJson(payload)
	}
}

func (bri *briClientRequest) GetUrl() string {
	return bri.cfg.InternalInquiryURL
}

func (bri *briClientRequest) GetHeaders(accessToken, externalId string, cfg *entity.BankConfig, payload []byte) map[string]string {
	hexPaylod := authorization.HashSHA256Hex(payload)
	endpoint := func() string {
		if cfg.BankCode == "002" {
			return cfg.InternalInquiryURL
		} else {
			return cfg.ExternalInquiryURL
		}
	}()
	timestamp := timehelper.FormatTimeToISO7(time.Now())
	signature := authorization.HashSignature("POST", endpoint, accessToken, hexPaylod, timestamp, cfg.ClientSecret)
	return map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + accessToken,
		"X-TIMESTAMP":   timestamp,
		"X-SIGNATURE":   signature,
		"X-PARTNER-ID":  cfg.PartnerId,
		"X-EXTERNAL-ID": externalId,
		"CHANNEL-ID":    cfg.ChannelId,
	}
}

func (bri *briClientResponse) MapResponse(bankResponse []byte) BankResponseData {

	if bri.httpStatus != http.StatusOK {
		var resDto dto.BRIErrorResponse
		json.Unmarshal(bankResponse, &resDto)
		return BankResponseData{
			ResponseMessage: resDto.ResponseMessage,
			AccountName:     "",
		}
	}

	if bri.cfg.BankCode == "002" {
		var resDto dto.BRIInternalInquiryResponse
		json.Unmarshal(bankResponse, &resDto)
		return BankResponseData{
			AccountName:     resDto.BeneficiaryAccountName,
			ResponseMessage: resDto.ResponseMessage,
		}
	} else {
		var resDto dto.BRIExternalInquiryResponse
		json.Unmarshal(bankResponse, &resDto)
		return BankResponseData{
			AccountName:     resDto.BeneficiaryAccountName,
			ResponseMessage: resDto.ResponseMessage,
		}
	}
}
