package authorization

import (
	"briefcash-inquiry/internal/dto"
	"briefcash-inquiry/internal/entity"
	"briefcash-inquiry/internal/helper/httphelper"
	"briefcash-inquiry/internal/helper/timehelper"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func readPEMFile(locFilePath string) ([]byte, error) {
	key, err := os.ReadFile(locFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}
	return key, nil
}

func loadPrivateKey(pemData []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block with private key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		key2, err2 := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err2 != nil {
			return nil, err
		}
		return key2, nil
	}

	pk, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not RSA private key")
	}

	return pk, nil
}

func signWithRSA(privateKey *rsa.PrivateKey, data string) (string, error) {
	hashed := sha256.Sum256([]byte(data))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

func GetAccessToken(cfg *entity.BankConfig, log *logrus.Entry) (dto.SNAPAccessToken, error) {
	var tokenResponse dto.SNAPAccessToken
	endpoint := cfg.BaseURL + cfg.AccessTokenURL
	timestamp := timehelper.FormatTimeToISO7(time.Now())
	stringToSign := fmt.Sprintf("%s|%s", cfg.ClientKey, timestamp)

	log.WithField("step", "read_pem").Info("Reading PEM file")
	pemFile, err := readPEMFile("./resource/private_key.pem")
	if err != nil {
		return dto.SNAPAccessToken{}, err
	}

	log.WithField("step", "get_private_key").Info("Extract private key from PEM file")
	pk, err := loadPrivateKey(pemFile)
	if err != nil {
		return dto.SNAPAccessToken{}, err
	}

	log.WithField("step", "sign_rsa").Info("Signing data with RSA")
	signature, err := signWithRSA(pk, stringToSign)
	if err != nil {
		return dto.SNAPAccessToken{}, err
	}

	log.WithField("step", "handle_payload").Info("Preparing payload and parse to JSON")
	payload := map[string]string{
		"grant_type": "client_credentials",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return dto.SNAPAccessToken{}, err
	}

	log.WithField("step", "handle_headers").Info("Define header paramaters")
	headers := map[string]string{
		"Content-Type": "application/json",
		"X-TIMESTAMP":  timestamp,
		"X-CLIENT-KEY": cfg.ClientKey,
		"X-SIGNATURE":  signature,
	}

	log.WithField("step", "send_request").Info("Send request access token to bank")
	client := httphelper.NewHttpClientHelper(10 * time.Second)
	resp, httpStatus, err := client.SendRequest("POST", endpoint, payloadBytes, headers)

	log.WithField("step", "handle_error").Info("Checking error return from bank")
	if err != nil {
		return dto.SNAPAccessToken{}, err
	}

	log.WithField("step", "handle_httpstatus").Info("Checking http status return")
	if httpStatus != http.StatusOK {
		return dto.SNAPAccessToken{}, fmt.Errorf("access unauthorized: http code %d", httpStatus)
	}

	log.WithField("step", "parse_response").Info("Parsing body response from JSON to struct")
	if err := json.Unmarshal(resp, &tokenResponse); err != nil {
		return dto.SNAPAccessToken{}, err
	}

	if tokenResponse.AccessToken == "" {
		return dto.SNAPAccessToken{}, fmt.Errorf("access token empty, response: %+v", tokenResponse)
	}

	log.WithField("step", "finalise_access_token").Info("Access token successfully retrieved from bank")
	return tokenResponse, nil
}
