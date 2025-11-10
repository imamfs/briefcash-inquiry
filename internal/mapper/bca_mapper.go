package mapper

import (
	"briefcash-inquiry/internal/authorization"
	"briefcash-inquiry/internal/dto"
	"briefcash-inquiry/internal/entity"
	"briefcash-inquiry/internal/helper/jsonhelper"
	"briefcash-inquiry/internal/helper/timehelper"
	"encoding/json"
	"time"
)

type bcaClientRequest struct {
	cfg *entity.BankConfig
	req dto.InquiryRequest
}

type bcaClientResponse struct {
	cfg *entity.BankConfig
}

func NewBcaClientRequest(cfg *entity.BankConfig, req dto.InquiryRequest) *bcaClientRequest {
	return &bcaClientRequest{
		cfg: cfg, req: req,
	}
}

func NewBcaClientResponse(cfg *entity.BankConfig) *bcaClientResponse {
	return &bcaClientResponse{cfg: cfg}
}

func (bca *bcaClientRequest) BuildBodyRequest() []byte {
	if bca.req.BankCode == "014" {
		request := dto.BCAInternalInquiryRequest{
			PartnerReferenceNo:   bca.req.PartnerReferenceNo,
			BeneficiaryAccountNo: bca.req.BeneficiaryAccount,
		}
		return jsonhelper.WriteToJson(request)
	} else {
		additionalInfo := dto.BCAAdditionalInfo{
			InquiryService: func() string {
				if bca.req.Type == "bifast" {
					return "2"
				} else {
					return "1"
				}
			}(),
		}
		request := dto.BCAExternalInquiryRequest{
			BeneficiaryBankCode:  bca.req.BankCode,
			BeneficiaryAccountNo: bca.req.BeneficiaryAccount,
			PartnerReferenceNo:   bca.req.PartnerReferenceNo,
			AdditionalInfo:       additionalInfo,
		}
		return jsonhelper.WriteToJson(request)
	}
}

func (bca *bcaClientRequest) GetUrl() string {
	var url string
	if bca.req.BankCode == "014" {
		url = bca.cfg.InternalInquiryURL
	} else {
		url = bca.cfg.ExternalInquiryURL
	}
	return url
}

func (bca *bcaClientRequest) GetHeaders(accessToken, externalId string, cfg *entity.BankConfig, payload []byte) map[string]string {
	hexPayload := authorization.HashSHA256Hex(payload)
	timestamp := timehelper.FormatTimeToISO7(time.Now())
	endpoint := func() string {
		if cfg.BankCode == "014" {
			return cfg.InternalInquiryURL
		} else {
			return cfg.ExternalInquiryURL
		}
	}()
	signature := authorization.HashSignature("POST", endpoint, accessToken, hexPayload, timestamp, cfg.ClientSecret)
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

func (bca *bcaClientResponse) MapResponse(bankResponse []byte) BankResponseData {
	var inquiryResponse dto.BCAInquiryResponse
	json.Unmarshal(bankResponse, &inquiryResponse)
	return BankResponseData{
		AccountName:     inquiryResponse.BeneficiaryAccountName,
		ResponseMessage: inquiryResponse.ResponseMessage,
	}
}
