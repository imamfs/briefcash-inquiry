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

type cimbClientRequest struct {
	cfg *entity.BankConfig
	req dto.InquiryRequest
}

type cimbClientResponse struct {
	cfg *entity.BankConfig
}

func NewCimbClientRequest(cfg *entity.BankConfig, req dto.InquiryRequest) *cimbClientRequest {
	return &cimbClientRequest{
		cfg: cfg, req: req,
	}
}

func NewCimbClientResponse(cfg *entity.BankConfig) *cimbClientResponse {
	return &cimbClientResponse{
		cfg: cfg,
	}
}

func (cimb *cimbClientRequest) BuildBodyRequest() []byte {
	if cimb.req.BankCode == "022" {
		payload := dto.CimbInternalInquiryRequest{
			PartnerReferenceNo:   cimb.req.PartnerReferenceNo,
			BeneficiaryAccountNo: cimb.req.BeneficiaryAccount,
			AdditionalInfo:       make(map[string]string),
		}
		return jsonhelper.WriteToJson(payload)
	} else {
		payload := dto.CimbExternalInquiryRequest{
			BeneficiaryBankCode:  cimb.req.BankCode,
			BeneficiaryAccountNo: cimb.req.BeneficiaryAccount,
			PartnerReferenceNo:   cimb.req.PartnerReferenceNo,
			AdditionalInfo: map[string]string{
				"trxType": func() string {
					if cimb.req.Type == "bifast" {
						return "02"
					} else {
						return "01"
					}
				}(),
				"proxyValue":     "01",
				"proxyType":      "01",
				"trxPurposeCode": "99",
			},
		}
		return jsonhelper.WriteToJson(payload)
	}
}

func (cimb *cimbClientRequest) GetUrl() string {
	var url string
	if cimb.req.BankCode == "022" {
		url = cimb.cfg.InternalInquiryURL
	} else {
		url = cimb.cfg.ExternalInquiryURL
	}
	return url
}

func (cimb *cimbClientRequest) GetHeaders(accessToken, externalId string, cfg *entity.BankConfig, payload []byte) map[string]string {
	hexPayload := authorization.HashSHA256Hex(payload)
	endpoint := func() string {
		if cfg.BankCode == "022" {
			return cfg.InternalInquiryURL
		} else {
			return cfg.ExternalInquiryURL
		}
	}()
	timestamp := timehelper.FormatTimeToISO7(time.Now())
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

func (cimb *cimbClientResponse) MapResponse(bankResponse []byte) BankResponseData {
	if cimb.cfg.BankCode == "022" {
		var respDto dto.CimbInternalInquiryResponse
		json.Unmarshal(bankResponse, &respDto)
		return BankResponseData{
			AccountName:     respDto.BeneficiaryAccountName,
			ResponseMessage: respDto.ResponseMessage,
		}
	} else {
		var respDto dto.CimbExternalInquiryResponse
		json.Unmarshal(bankResponse, &respDto)
		return BankResponseData{
			AccountName:     respDto.BeneficiaryAccountName,
			ResponseMessage: respDto.ResponseMessage,
		}
	}
}
