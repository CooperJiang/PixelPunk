package user

import (
	"pixelpunk/internal/middleware"
	userService "pixelpunk/internal/services/user"
	"pixelpunk/pkg/errors"

	"github.com/gin-gonic/gin"
)

func GetWorkspaceStats(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "获取用户信息失败"))
		return
	}

	data, err := userService.GetWorkspaceStats(userID)
	if err != nil {
		errors.HandleError(c, errors.Wrap(err, errors.CodeInternal, "获取工作台数据失败"))
		return
	}

	errors.ResponseSuccess(c, data, "获取工作台数据成功")
}
