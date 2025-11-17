package middleware

import (
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/bandwidth"
	"pixelpunk/pkg/assets"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"

	"github.com/gin-gonic/gin"
)

func BandwidthLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		userID := GetCurrentUserID(c)
		if userID == 0 {
			c.Next()
			return
		}

		fileObj, exists := c.Get("file_info")
		if !exists {
			logger.Error("[BANDWIDTH_LIMIT] 未找到文件信息，请确保 FileInfoExtractorMiddleware 在前")
			c.Next()
			return
		}

		file, ok := fileObj.(models.File)
		if !ok {
			logger.Error("[BANDWIDTH_LIMIT] 文件信息类型转换失败")
			c.Next()
			return
		}

		available, err := bandwidth.Service.CheckBandwidthAvailable(userID, file.Size)
		if err != nil {
			logger.Error("[BANDWIDTH_LIMIT] 检查带宽限制失败: %v", err)
			errors.HandleError(c, errors.Wrap(err, errors.CodeInternal, "检查带宽限制失败"))
			c.Abort()
			return
		}

		if !available {
			assets.ServeDefaultFile(c, assets.FileTypeBandwidthLimit)
			return
		}

		c.Next()
	}
}
