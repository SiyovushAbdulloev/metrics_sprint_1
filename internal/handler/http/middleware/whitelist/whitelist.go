package whitelist

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

func IPWhitelist(trustedCIDR string) gin.HandlerFunc {
	if trustedCIDR == "" {
		return func(c *gin.Context) { c.Next() }
	}

	_, subnet, err := net.ParseCIDR(trustedCIDR)
	if err != nil {
		return func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid trusted_subnet"})
		}
	}

	return func(c *gin.Context) {
		ip := c.GetHeader("X-Real-IP")
		parsedIP := net.ParseIP(ip)
		if parsedIP == nil || !subnet.Contains(parsedIP) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}
