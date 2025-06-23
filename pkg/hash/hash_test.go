package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateHashSHA256(t *testing.T) {
	value := []byte("test-body")
	key := "secret"

	result := CalculateHashSHA256(value, key)

	// Пример: длина HMAC SHA256 в hex — 64 символа
	assert.Len(t, result, 64)
	assert.NotEmpty(t, result)
}

func TestValidateHash_Valid(t *testing.T) {
	value := []byte("payload")
	key := "key123"
	expected := CalculateHashSHA256(value, key)

	valid := ValidateHash(value, expected, key)
	assert.True(t, valid)
}

func TestValidateHash_Invalid(t *testing.T) {
	value := []byte("payload")
	key := "key123"
	invalidHash := "deadbeef"

	valid := ValidateHash(value, invalidHash, key)
	assert.False(t, valid)
}
