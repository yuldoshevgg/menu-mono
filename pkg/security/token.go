package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
)

func GenerateToken(data map[string]interface{}, secret []byte) string {
	payload, _ := json.Marshal(data)
	encoded := base64.StdEncoding.EncodeToString(payload)

	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(encoded))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return encoded + "." + signature
}

func VerifyToken(token string, secret []byte) (map[string]interface{}, bool) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, false
	}
	payloadEncoded, signature := parts[0], parts[1]

	// Recalculate signature
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(payloadEncoded))
	expectedSig := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expectedSig)) {
		return nil, false
	}

	payloadJson, _ := base64.StdEncoding.DecodeString(payloadEncoded)
	var data map[string]interface{}
	if err := json.Unmarshal(payloadJson, &data); err != nil {
		return nil, false
	}

	return data, true
}
