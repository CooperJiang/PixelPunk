package common

import (
	"pixelpunk/internal/middleware"
	"pixelpunk/pkg/errors"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	authParams := middleware.GetAuthParams(c)
	res := gin.H{
		"ip":      authParams.IP,
		"domain":  authParams.Domain,
		"referer": authParams.Referer,
		"message": "您的访问权限正常",
	}
	errors.ResponseSuccess(c, res, "健康检查成功")
}
