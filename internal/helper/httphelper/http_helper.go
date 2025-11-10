package httphelper

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HttpClientHelper struct {
	client *http.Client
}

func NewHttpClientHelper(timeout time.Duration) *HttpClientHelper {
	return &HttpClientHelper{client: &http.Client{
		Timeout: timeout,
	}}
}

func (c *HttpClientHelper) SendRequest(method, url string, payload []byte, headers map[string]string) ([]byte, int, error) {
	var body io.Reader

	if payload != nil {
		body = bytes.NewBuffer(payload)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %v", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response: %v", err)
	}
	return respBody, resp.StatusCode, nil
}
