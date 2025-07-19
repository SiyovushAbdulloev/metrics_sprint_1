package middleware

import (
	"bytes"
	"crypto/rsa"
	"fmt"
	"github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/crypto"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

// DecryptBody расшифровывает тело запроса, используя RSA приватный ключ
func DecryptBody(privKey *rsa.PrivateKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Прочитать тело
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
			c.Abort()
			return
		}

		// Расшифровать
		plaintext, err := crypto.DecryptWithPrivateKey(body, privKey)
		fmt.Println("Plaintext: ", string(plaintext))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to decrypt request body"})
			c.Abort()
			return
		}

		// Заменить тело на расшифрованное
		c.Request.Body = io.NopCloser(bytes.NewReader(plaintext))
		c.Next()
	}
}
