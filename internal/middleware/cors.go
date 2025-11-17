package middleware

import (
	"github.com/gin-gonic/gin"
)

/* CORSMiddleware 写死放开所有跨域，始终允许任何来源和常见请求头 */
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Vary", "Origin")
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")

		requestedHeaders := c.Request.Header.Get("Access-Control-Request-Headers")
		baseAllowedHeaders := "Authorization, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Accept, X-Requested-With, X-API-Key, x-pixelpunk-key"
		if requestedHeaders != "" {
			c.Writer.Header().Set("Access-Control-Allow-Headers", baseAllowedHeaders+", "+requestedHeaders)
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Headers", baseAllowedHeaders)
		}

		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition, Content-Type, X-Request-Id, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
