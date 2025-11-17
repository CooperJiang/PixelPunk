package middleware

import (
	"net/http"
	"strings"

	"pixelpunk/internal/services/apikey"
	"pixelpunk/pkg/errors"

	"github.com/gin-gonic/gin"
)

func APIKeyAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		keyValue := c.GetHeader("x-pixelpunk-key")
		if keyValue == "" {
			keyValue = c.GetHeader("x-api-key")
		}

		if keyValue == "" {
			authz := c.GetHeader("Authorization")
			if authz != "" {
				parts := strings.SplitN(authz, " ", 2)
				if len(parts) == 2 {
					prefix := strings.ToLower(strings.TrimSpace(parts[0]))
					token := strings.TrimSpace(parts[1])
					if prefix == "bearer" || prefix == "apikey" {
						keyValue = token
					}
				}
			}
		}

		if keyValue == "" {
			keyValue = c.Query("key")
			if keyValue == "" {
				keyValue = c.Query("api_key")
			}
		}

		if keyValue == "" {
			errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未提供API密钥"))
			c.Abort()
			return
		}

		key, err := apikey.ValidateAPIKey(keyValue)
		if err != nil {
			errors.HandleError(c, err)
			c.Abort()
			return
		}

		if c.Request.Method == "POST" {
			if !key.CheckUploadCountLimit() {
				errors.HandleError(c, errors.New(errors.CodeForbidden, "已达到上传次数限制"))
				c.Abort()
				return
			}

		}

		c.Set("api_key", key)
		c.Set("api_key_id", key.ID)
		c.Set("user_id", key.UserID)

		c.Next()
	}
}
