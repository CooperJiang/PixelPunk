package middleware

import (
	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"

	"github.com/gin-gonic/gin"
)

func OptionalAuthForFileDownload() gin.HandlerFunc {
	return func(c *gin.Context) {
		fileID := c.Param("file_id")
		if fileID == "" {
			errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件ID不能为空"))
			c.Abort()
			return
		}

		var file models.File
		if err := database.DB.Where("id = ?", fileID).
			Where("status <> ?", "pending_deletion").
			First(&file).Error; err != nil {
			errors.HandleError(c, errors.New(errors.CodeFileNotFound, "文件不存在"))
			c.Abort()
			return
		}

		c.Set("file_info", file)

		switch file.AccessLevel {
		case "public", "private":
			c.Next()
			return

		case "protected":
			if authError, exists := c.Get(AuthErrorKey); exists {
				errors.HandleError(c, errors.New(errors.CodeUnauthorized, authError.(string)))
				c.Abort()
				return
			}

			currentUser := GetCurrentUser(c)
			if currentUser == nil {
				errors.HandleError(c, errors.New(errors.CodeUnauthorized, "需要登录才能下载此文件"))
				c.Abort()
				return
			}

			if !canUserDownloadFile(currentUser.UserID, file) {
				errors.HandleError(c, errors.New(errors.CodeFileAccessDenied, "无权限下载此文件"))
				c.Abort()
				return
			}

			c.Next()
			return

		default:
			errors.HandleError(c, errors.New(errors.CodeFileAccessDenied, "文件访问级别未知，拒绝访问"))
			c.Abort()
			return
		}
	}
}

func canUserDownloadFile(userID uint, file models.File) bool {
	switch file.AccessLevel {
	case "public", "private":
		return true
	case "protected":
		if userID == file.UserID {
			return true
		}
		var user models.User
		if err := database.DB.Where("id = ?", userID).First(&user).Error; err == nil {
			return user.Role == common.UserRoleAdmin || user.Role == common.UserRoleSuperAdmin
		}
		return false
	default:
		return false
	}
}
