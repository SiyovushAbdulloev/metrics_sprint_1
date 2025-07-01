// Package middleware предоставляет middleware-функции для HTTP-сервера, включая логгирование, хэширование и сжатие.
package middleware

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"slices"
	"strings"
)

var allowedExtensions = []string{
	"application/json",
	"text/html",
}

type gzipWriter struct {
	gin.ResponseWriter
	Writer *gzip.Writer
}

// Write записывает сжатые данные в поток ответа.
func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// WriteHeader отправляет HTTP-заголовок с указанным кодом состояния.
func (w gzipWriter) WriteHeader(code int) {
	//w.Header().Del("Content-Length") // Prevent incorrect Content-Length issues
	w.ResponseWriter.WriteHeader(code)
}

// Flush завершает запись буфера и передаёт данные клиенту.
func (w gzipWriter) Flush() {
	w.Writer.Flush()
	w.ResponseWriter.Flush()
}

// Compress возвращает middleware-функцию для Gin, которая реализует автоматическое GZIP-сжатие HTTP-ответов и декомпрессию входящих запросов.
// Работает только с типами, определёнными в allowedExtensions.
func Compress() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Content-Encoding") == "gzip" {
			gzipReader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				io.WriteString(c.Writer, err.Error())
				return
			}

			defer gzipReader.Close()

			c.Request.Body = io.NopCloser(gzipReader)
		}
		if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}
		accept := c.Request.Header.Get("Accept")
		if !slices.Contains(allowedExtensions, accept) {
			c.Next()
			return
		}

		gw, err := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			c.Abort()
			return
		}

		defer gw.Close()

		gzipW := gzipWriter{
			ResponseWriter: c.Writer,
			Writer:         gw,
		}

		c.Writer = gzipW
		c.Writer.Header().Set("Content-Encoding", "gzip")
		c.Writer.Header().Set("Content-Type", accept)

		c.Next()

		gzipW.Flush()
	}
}
