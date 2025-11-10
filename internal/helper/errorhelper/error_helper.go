package errorhelper

import (
	"briefcash-inquiry/internal/dto"
	"fmt"
)

type ErrorDetail struct {
	Code       string
	Message    string
	LogMessage string
	Source     string
}

const (
	SourceBank     = "bank"
	SourceInternal = "internal"
	SourceClient   = "client"
)

var ErrorMap = map[int]ErrorDetail{
	400: {Code: "INVALID_BODY", Message: "Invalid payload request", LogMessage: "Invalid body verified by bank", Source: SourceBank},
	401: {Code: "UNAUTHORIZED", Message: "Access unauthorized", LogMessage: "Bank return unauthorized access", Source: SourceBank},
	403: {Code: "FORBIDDEN_FEATURE", Message: "Service not allowed", LogMessage: "Feature forbidden by bank", Source: SourceBank},
	404: {Code: "ACCOUNT_NOT_FOUND", Message: "Account number not found", LogMessage: "Account number not found in bank system", Source: SourceBank},
	409: {Code: "DUPLICATE_REFERENCE", Message: "Duplicate external id in same day", LogMessage: "Duplicate external id request", Source: SourceBank},
	500: {Code: "BANK_INTERNAL_ERROR", Message: "Bank internal error, please use check status service", LogMessage: "Bank returned internal error", Source: SourceBank},
	504: {Code: "BANK_TIMEOUT", Message: "Bank timeout, please use check status service", LogMessage: "Bank timeout while processing request", Source: SourceBank},
}

var DefaultBankError = ErrorDetail{
	Code:       "BANK_UNKOWN_ERROR",
	Message:    "Unexpecter error from bank",
	LogMessage: "Unkown bank error occured",
	Source:     SourceBank,
}

func BuildErrorResponse(errDetail ErrorDetail, message string, err error) (*dto.InquiryResponse, error) {
	return &dto.InquiryResponse{
		Status:  false,
		Message: errDetail.Message,
		Code:    errDetail.Code,
		Source:  errDetail.Source,
		Data:    dto.InquiryData{},
	}, fmt.Errorf("%s: %w", errDetail.LogMessage, err)
}
