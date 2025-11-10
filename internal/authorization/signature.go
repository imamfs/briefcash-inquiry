package authorization

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func HashSHA256Hex(payload []byte) string {
	hash := sha256.New()
	hash.Write(payload)
	return hex.EncodeToString(hash.Sum(nil))
}

func HashSignature(httpMethod, relativeUrl, accessToken, bodyHash, timestamp, apiSecret string) string {
	stringToSign := fmt.Sprintf("%s:%s:%s:%s:%s:", httpMethod, relativeUrl, accessToken, bodyHash, timestamp)
	mac := hmac.New(sha256.New, []byte(apiSecret))
	mac.Write([]byte(stringToSign))
	return hex.EncodeToString(mac.Sum(nil))
}
