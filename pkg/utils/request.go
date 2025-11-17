package utils

import (
	"net"

	"github.com/gin-gonic/gin"
)

func GetClientIP(c *gin.Context) string {
	// 尝试从 X-Forwarded-For 头获取
	if xForwardedFor := c.GetHeader("X-Forwarded-For"); xForwardedFor != "" {
		// X-Forwarded-For 可能包含多个IP，取第一个
		if ip, _, err := net.SplitHostPort(xForwardedFor); err == nil {
			return ip
		}
		return xForwardedFor
	}

	// 尝试从 X-Real-IP 头获取
	if xRealIP := c.GetHeader("X-Real-IP"); xRealIP != "" {
		return xRealIP
	}

	// 从 RemoteAddr 获取
	if ip, _, err := net.SplitHostPort(c.Request.RemoteAddr); err == nil {
		return ip
	}

	return c.Request.RemoteAddr
}
