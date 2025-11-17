package user

import (
	"pixelpunk/internal/controllers/user/dto"
	"pixelpunk/internal/middleware"
	userService "pixelpunk/internal/services/user"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/errors"

	"github.com/gin-gonic/gin"
)

func GetUserFolders(c *gin.Context) {
	requestUserID := middleware.GetCurrentUserID(c)
	if requestUserID == 0 {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "获取用户信息失败"))
		return
	}

	query, err := common.ValidateRequest[dto.FolderQueryDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	data, err := userService.GetUserFolders(requestUserID, query)
	if err != nil {
		errors.HandleError(c, errors.Wrap(err, errors.CodeInternal, "获取文件夹列表失败"))
		return
	}

	errors.ResponseSuccess(c, data, "获取文件夹列表成功")
}
