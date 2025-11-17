package middleware

import (
	"net"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func getClientIP(c *gin.Context) string {
	if xForwardedFor := c.GetHeader("X-Forwarded-For"); xForwardedFor != "" {
		if ip, _, err := net.SplitHostPort(xForwardedFor); err == nil {
			return ip
		}
		return xForwardedFor
	}

	if xRealIP := c.GetHeader("X-Real-IP"); xRealIP != "" {
		return xRealIP
	}

	if ip, _, err := net.SplitHostPort(c.Request.RemoteAddr); err == nil {
		return ip
	}

	return c.Request.RemoteAddr
}

func updateUserActivity(userID uint, clientIP string) {
	go func() {
		db := database.GetDB()
		if db == nil {
			logger.Error("无法获取数据库连接，跳过用户活动更新")
			return
		}

		now := common.JSONTime(time.Now())
		result := db.Model(&models.User{}).
			Where("id = ?", userID).
			Updates(map[string]interface{}{
				"last_activity_at": &now,
				"last_activity_ip": clientIP,
				"updated_at":       &now,
			})

		if result.Error != nil {
			logger.Error("更新用户 %d 活动信息失败: %v", userID, result.Error)
		}
	}()
}

func TrackUserActivity() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if _, hasAuthError := c.Get(AuthErrorKey); !hasAuthError {
			if claims := GetCurrentUser(c); claims != nil {
				clientIP := getClientIP(c)

				updateUserActivity(claims.UserID, clientIP)
			}
		}
	}
}
