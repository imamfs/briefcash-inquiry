package routinghelper

import (
	"briefcash-inquiry/internal/dto"
	"briefcash-inquiry/internal/entity"
	"briefcash-inquiry/internal/mapper"
)

type BankRouteRequest interface {
	BuildBodyRequest() []byte
	GetUrl() string
	GetHeaders(accessToken, externalId string, cfg *entity.BankConfig, payload []byte) map[string]string
}

func NewBankRouteRequest(req dto.InquiryRequest, cfg *entity.BankConfig, partnerRefNo string) BankRouteRequest {
	switch req.BankCode {
	case "002":
		return mapper.NewBriClientRequest(cfg, req)
	case "013":
		return mapper.NewPermataClientRequest(cfg, req)
	case "014":
		return mapper.NewBcaClientRequest(cfg, req)
	case "022":
		return mapper.NewCimbClientRequest(cfg, req)
	default:
		return mapper.NewBcaClientRequest(cfg, req)
	}
}

type BankRouteResponse interface {
	MapResponse(bankResponse []byte) mapper.BankResponseData
}

func NewBankRouteResponse(cfg *entity.BankConfig, httpStatus int) BankRouteResponse {
	switch cfg.BankCode {
	case "002":
		return mapper.NewBriClientResponse(cfg, httpStatus)
	case "013":
		return mapper.NewPermataClientResponse(cfg)
	case "014":
		return mapper.NewBcaClientResponse(cfg)
	case "022":
		return mapper.NewCimbClientResponse(cfg)
	default:
		return mapper.NewBcaClientResponse(cfg)
	}
}
