// Package middleware предоставляет middleware-функции для HTTP-сервера, включая логгирование, хэширование и сжатие.
package middleware

import (
	"github.com/SiyovushAbdulloev/metriks_sprint_1/pkg/logger"
	"github.com/gin-gonic/gin"
	"time"
)

type responseData struct {
	data   int
	status int
}

type loggerResponseWriter struct {
	gin.ResponseWriter
	responseData *responseData
}

// Write записывает тело ответа и подсчитывает количество байт.
func (w *loggerResponseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.responseData.data += size
	return size, err
}

// WriteHeader сохраняет статус-код и передаёт его в оригинальный ResponseWriter.
func (w *loggerResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.responseData.status = statusCode
}

// Logger возвращает middleware, логирующее метод, URI, статус ответа, размер ответа и время обработки запроса.
// Использует переданный интерфейс логгера.
func Logger(l logger.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		responseWriter := &loggerResponseWriter{
			ResponseWriter: c.Writer,
			responseData: &responseData{
				data:   0,
				status: 0,
			},
		}
		t := time.Now()
		c.Writer = responseWriter

		c.Next()

		latency := time.Since(t)
		method := c.Request.Method
		uri := c.Request.RequestURI

		l.Info(
			"Data",
			"latency", latency,
			"method", method,
			"uri", uri,
			"status", responseWriter.responseData.status,
			"data", responseWriter.responseData.data,
		)
	}
}
