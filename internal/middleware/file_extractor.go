package middleware

import (
	"pixelpunk/internal/models"
	"pixelpunk/pkg/assets"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"
	"strings"

	"github.com/gin-gonic/gin"
)

/* FileInfoExtractorMiddleware 提取文件信息到上下文（file_info, isThumb） */
func FileInfoExtractorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		var file models.File
		var err error
		var isThumb bool

		if strings.HasPrefix(path, "/f/") {
			fileID := c.Param("fileID")
			if fileID == "" {
				logger.Error("[IMAGE_EXTRACTOR] 文件ID参数为空")
				assets.ServeDefaultFile(c, assets.FileTypeNotFound)
				return
			}

			err = database.DB.Where("id = ?", fileID).First(&file).Error
			isThumb = false

		} else if strings.HasPrefix(path, "/t/") || strings.HasPrefix(path, "/thumb/") {
			fileID := c.Param("fileID")
			if fileID == "" {
				logger.Error("[IMAGE_EXTRACTOR] 文件ID参数为空")
				assets.ServeDefaultFile(c, assets.FileTypeNotFound)
				return
			}

			err = database.DB.Where("id = ?", fileID).First(&file).Error
			isThumb = true

		} else if strings.HasPrefix(path, "/s/") {
			shortURL := c.Param("shortURL")
			if shortURL == "" {
				logger.Error("[IMAGE_EXTRACTOR] shortURL参数为空")
				assets.ServeDefaultFile(c, assets.FileTypeNotFound)
				return
			}

			err = database.DB.Where("short_url = ?", shortURL).First(&file).Error
			isThumb = false

		} else {
			logger.Error("[IMAGE_EXTRACTOR] 不支持的路径格式: %s", path)
			assets.ServeDefaultFile(c, assets.FileTypeNotFound)
			return
		}

		if err != nil {
			logger.Error("[IMAGE_EXTRACTOR] 查找文件失败: %v", err)
			assets.ServeDefaultFile(c, assets.FileTypeNotFound)
			return
		}

		if file.Status == "pending_deletion" {
			assets.ServeDefaultFile(c, assets.FileTypeNotFound)
			return
		}

		c.Set("file_info", file)
		c.Set("isThumb", isThumb)

		c.Next()
	}
}
