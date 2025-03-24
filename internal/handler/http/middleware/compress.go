package middleware

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"slices"
	"strings"
)

var allowedExtensions = []string{
	"application/json",
	"text/html",
}

type gzipWriter struct {
	gin.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

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
		if !slices.Contains(allowedExtensions, c.Request.Header.Get("Content-Type")) {
			c.Next()
			return
		}

		gw, err := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed)
		if err != nil {
			io.WriteString(c.Writer, err.Error())
			return
		}

		defer gw.Close()

		gzipW := gzipWriter{
			ResponseWriter: c.Writer,
			Writer:         gw,
		}

		gzipW.ResponseWriter.Header().Set("Content-Encoding", "gzip")
		c.Writer = gzipW
	}
}
