package middleware

import (
	"bytes"
	hash2 "github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/hash"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type responseWriter struct {
	gin.ResponseWriter
	body    *bytes.Buffer
	hashKey string
}

func (rw responseWriter) Write(b []byte) (int, error) {
	if len(rw.hashKey) > 0 {
		hash := hash2.CalculateHashSHA256(b, rw.hashKey)
		rw.Header().Set("HashSHA256", hash)
	}

	rw.body.Write(b) // сохраняем тело
	return rw.ResponseWriter.Write(b)
}

func Hash(hashKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(hashKey) > 0 {
			hash := c.Request.Header.Get("HashSHA256")
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{})
				return
			}

			if !hash2.ValidateHash(body, hash, hashKey) {
				c.JSON(http.StatusBadRequest, gin.H{})
				return
			}

			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		buf := &bytes.Buffer{}
		rw := &responseWriter{
			body:           buf,
			ResponseWriter: c.Writer,
			hashKey:        hashKey,
		}
		c.Writer = rw

		c.Next()
	}
}
