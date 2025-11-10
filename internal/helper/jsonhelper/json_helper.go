package jsonhelper

import (
	dto "briefcash-inquiry/internal/dto"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func WriteHttpJson(w http.ResponseWriter, httpStatus int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(data)
}

func WriteJsonError(w http.ResponseWriter, httpStatus int, message string) {
	WriteHttpJson(w, httpStatus, dto.InquiryResponse{
		Status:  false,
		Message: message,
	})
}

func ParseJsonBody(r *http.Request, destination interface{}) ([]byte, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body request: %w", err)
	}
	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	if err := json.Unmarshal(bodyBytes, destination); err != nil {
		return nil, fmt.Errorf("invalid json format: %w", err)
	}
	return bodyBytes, nil
}

func WriteToJson(payload any) []byte {
	body, _ := json.Marshal(payload)
	return body
}
