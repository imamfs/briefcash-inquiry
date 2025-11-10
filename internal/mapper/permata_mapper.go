package mapper

import (
	"briefcash-inquiry/internal/dto"
	"briefcash-inquiry/internal/entity"
	"briefcash-inquiry/internal/helper/jsonhelper"
	"briefcash-inquiry/internal/helper/timehelper"
	"encoding/json"
	"time"
)

type permataClientRequest struct {
	cfg *entity.BankConfig
	req dto.InquiryRequest
}

type permataClientResponse struct {
	cfg *entity.BankConfig
}

func NewPermataClientRequest(cfg *entity.BankConfig, req dto.InquiryRequest) *permataClientRequest {
	return &permataClientRequest{
		cfg: cfg, req: req,
	}
}

func NewPermataClientResponse(cfg *entity.BankConfig) *permataClientResponse {
	return &permataClientResponse{
		cfg: cfg,
	}
}

func (permata *permataClientRequest) BuildBodyRequest() []byte {
	var wrapper map[string]any
	headerMsg := dto.PermataInquiryHeaderRequest{
		RequestTimeStamp: timehelper.FormatTimeToISO7(time.Now()),
		CustReffID:       permata.req.CompanyId,
	}

	bodyMsg := dto.PermataInternalInquiryBodyRequest{
		AccountNumber: permata.req.BeneficiaryAccount,
	}

	payload := dto.PermataInternalInquiryRequest{
		MessageHeader: headerMsg,
		MessageBody:   bodyMsg,
	}

	wrapper = map[string]any{
		"AcctInqRq": payload,
	}

	return jsonhelper.WriteToJson(wrapper)
}

func (permata *permataClientRequest) GetUrl() string {
	return permata.cfg.InternalInquiryURL
}

func (permata *permataClientRequest) GetHeaders(accessToken, externalId string, cfg *entity.BankConfig, payload []byte) map[string]string {
	return map[string]string{
		"Content-Type":     "application/json",
		"OrganizationName": permata.req.CompanyId,
	}
}

func (permata *permataClientResponse) MapResponse(bankResponse []byte) BankResponseData {
	if permata.cfg.BankCode == "013" {
		var wrapper map[string]interface{}
		json.Unmarshal(bankResponse, &wrapper)
		data := wrapper["AcctInqRs"].(dto.PermataInternalInquiryResponse)
		accountName := data.MessageBody.AccountName
		responseMessage := data.MessageHeader.StatusDesc
		return BankResponseData{
			AccountName:     accountName,
			ResponseMessage: responseMessage,
		}
	} else {
		var wrapper map[string]interface{}
		json.Unmarshal(bankResponse, &wrapper)
		data := wrapper["OlXferInqRs"].(dto.PermataExternalInquiryResponse)
		accountName := data.MessageBody.ToAccountFullName
		responseMessage := data.MessageHeader.StatusDesc
		return BankResponseData{
			AccountName:     accountName,
			ResponseMessage: responseMessage,
		}
	}
}
