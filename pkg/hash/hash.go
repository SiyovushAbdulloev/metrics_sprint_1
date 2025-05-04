package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func CalculateHashSHA256(value []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(value)
	return hex.EncodeToString(h.Sum(nil))
}

func ValidateHash(body []byte, receivedHash, key string) bool {
	expected := CalculateHashSHA256(body, key)
	return receivedHash == expected
}
