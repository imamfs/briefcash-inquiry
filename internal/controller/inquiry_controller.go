package controller

import (
	"briefcash-inquiry/internal/dto"
	"briefcash-inquiry/internal/helper/errorhelper"
	"briefcash-inquiry/internal/helper/loghelper"
	"briefcash-inquiry/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type inquiryController struct {
	svc service.InquiryService
}

func NewInquiryController(svc service.InquiryService) *inquiryController {
	return &inquiryController{svc}
}

func (ctr *inquiryController) InquiryAccountNumber(c *gin.Context) {
	start := time.Now()

	var req dto.InquiryRequest
	partnerRefNo := c.GetHeader("X-PARTNER-REFERENCE")

	if partnerRefNo == "" {
		c.JSON(http.StatusBadRequest, dto.InquiryResponse{
			Status:  false,
			Code:    "CLIENT_MISSING_HEADER",
			Message: "Missing X-PARTNER-REFERENCE header",
			Source:  errorhelper.SourceClient,
			Data:    dto.InquiryData{},
		})
		return
	}

	log := loghelper.Logger.WithFields(logrus.Fields{
		"service":  "inquiry_controller",
		"trace_id": partnerRefNo,
	})

	defer func() {
		log.WithFields(logrus.Fields{
			"step":            "return_success_response",
			"processing_time": time.Since(start).Milliseconds(),
		}).Info("Inquiry account number successfuly requested")
	}()

	log.WithField("step", "payload_validation").Info("Validating payload request")
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.InquiryResponse{
			Status:  false,
			Code:    "CLIENT_ERROR_REQUEST",
			Message: "Invalid body request",
			Source:  errorhelper.SourceClient,
			Data:    dto.InquiryData{},
		})
		return
	}

	log.WithField("step", "send_inquiry_request").Info("Sending inquiry account request")
	response, err := ctr.svc.InquiryAccount(c.Request.Context(), req, partnerRefNo)

	log.WithField("step", "error_validation").Info("Validating error response from bank")
	if err != nil {
		log.WithFields(logrus.Fields{
			"step":            "return_failed_response",
			"processing_time": time.Since(start).Milliseconds(),
		}).Error(response.Message)

		statusMap := map[string]int{
			"INTERNAL_CONNECTION_ERROR": http.StatusGatewayTimeout,
			"BANK_NO_RESPONSE":          http.StatusGatewayTimeout,
			"BANK_FORMAT_ERROR":         http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR":     http.StatusInternalServerError,
			"INVALID_BODY":              http.StatusInternalServerError,
			"UNAUTHORIZED":              http.StatusInternalServerError,
			"FORBIDDEN_FEATURE":         http.StatusInternalServerError,
			"ACCOUNT_NOT_FOUND":         http.StatusNotFound,
			"DUPLICATE_REFERENCE":       http.StatusConflict,
			"BANK_INTERNAL_ERROR":       http.StatusBadGateway,
			"BANK_TIMEOUT":              http.StatusGatewayTimeout,
		}

		status := statusMap[response.Code]
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, response)
}
