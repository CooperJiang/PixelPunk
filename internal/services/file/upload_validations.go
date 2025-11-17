package file

/* Validation helpers split from upload_service.go (no behavior change). */

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/setting"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"strings"
	"time"

	"gorm.io/gorm"
)

func validateUploadInput(ctx *UploadContext) error {
	settingsMap, err := setting.GetSettingsByGroupAsMap("upload")
	maxFileSize := int64(20 * 1024 * 1024) // 默认20MB

	if err == nil {
		if maxSizeVal, ok := settingsMap.Settings["max_file_size"]; ok {
			if maxSizeMB, ok := maxSizeVal.(float64); ok {
				maxFileSize = int64(maxSizeMB * 1024 * 1024) // 转换MB为字节
			}
		}
	}

	if maxFileSize > 0 && ctx.File.Size > maxFileSize {
		maxSizeMB := maxFileSize / (1024 * 1024)
		return errors.New(errors.CodeFileTooLarge, fmt.Sprintf("文件大小不能超过%dMB", maxSizeMB))
	}

	fileExt := strings.ToLower(filepath.Ext(ctx.File.Filename))
	ctx.FileExt = fileExt

	if !isValidFileType(fileExt) {
		return errors.New(errors.CodeFileTypeNotSupported, "当前格式不被支持、请联系管理员解除限制！")
	}

	if ctx.FolderID == "null" {
		ctx.FolderID = ""
	}
	if ctx.FolderID != "" {
		return validateFolder(ctx)
	}

	if ctx.StorageDuration != "" {
		storageConfig, err := setting.CreateStorageConfig()
		if err != nil {
			logger.Warn("获取存储配置失败，使用默认配置: %v", err)
			storageConfig = common.CreateDefaultStorageConfig()
		}
		if err := storageConfig.ValidateStorageDuration(ctx.StorageDuration, ctx.IsGuestUpload); err != nil {
			return errors.Wrap(err, errors.CodeInvalidParameter, err.Error())
		}
	}
	return nil
}

func isValidFileType(ext string) bool {
	settingsMap, err := setting.GetSettingsByGroupAsMap("upload")
	if err != nil {
		logger.Warn("获取文件格式设置失败，使用默认配置: %v", err)
		validTypes := map[string]bool{
			".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
			".bmp": true, ".apng": true, ".svg": true, ".ico": true, ".jp2": true,
			".tiff": true, ".tif": true, ".tga": true, ".heic": true, ".heif": true,
		}
		return validTypes[ext]
	}
	if formatsInterface, ok := settingsMap.Settings["allowed_file_formats"]; ok {
		if formats, ok := formatsInterface.([]any); ok {
			extWithoutDot := strings.TrimPrefix(ext, ".")
			for _, format := range formats {
				if formatStr, ok := format.(string); ok && formatStr == extWithoutDot {
					return true
				}
			}
			return false
		}
	}
	validTypes := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
		".bmp": true, ".apng": true, ".svg": true, ".ico": true, ".jp2": true,
		".tiff": true, ".tif": true, ".tga": true, ".heic": true, ".heif": true,
	}
	return validTypes[ext]
}

func validateFolder(ctx *UploadContext) error {
	if ctx.FolderID == "" {
		return nil
	}
	var folder models.Folder
	if err := database.DB.Where("id = ? AND user_id = ?", ctx.FolderID, ctx.UserID).First(&folder).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.CodeFolderNotFound, "文件夹不存在")
		}
		return errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件夹失败")
	}
	return nil
}

func validateBatchUploadFiles(files []*multipart.FileHeader) error {
	settingsMap, err := setting.GetSettingsByGroupAsMap("upload")
	maxFileSize := int64(20 * 1024 * 1024)   // 默认20MB单文件限制
	maxBatchSize := int64(100 * 1024 * 1024) // 默认100MB批量限制
	if err == nil {
		if maxSizeVal, ok := settingsMap.Settings["max_file_size"]; ok {
			if maxSizeMB, ok := maxSizeVal.(float64); ok {
				maxFileSize = int64(maxSizeMB * 1024 * 1024)
			}
		}
		if maxBatchSizeVal, ok := settingsMap.Settings["max_batch_size"]; ok {
			if maxBatchSizeMB, ok := maxBatchSizeVal.(float64); ok {
				maxBatchSize = int64(maxBatchSizeMB * 1024 * 1024)
			}
		}
	}
	var totalSize int64
	for _, file := range files {
		if maxFileSize > 0 && file.Size > maxFileSize {
			maxSizeMB := maxFileSize / (1024 * 1024)
			return errors.New(errors.CodeFileTooLarge, fmt.Sprintf("文件%s大小超过单文件限制%dMB", file.Filename, maxSizeMB))
		}
		totalSize += file.Size
	}
	if maxBatchSize > 0 && totalSize > maxBatchSize {
		maxSizeMB := maxBatchSize / (1024 * 1024)
		return errors.New(errors.CodeFileTooLarge, fmt.Sprintf("批量上传总大小不能超过%dMB", maxSizeMB))
	}
	return nil
}

func checkDailyUploadLimit(userID uint, uploadCount int) (bool, error) {
	settingsMap, err := setting.GetSettingsByGroupAsMap("upload")
	if err != nil {
		return false, err
	}
	var dailyLimit int = 50 // 默认值
	if limitVal, ok := settingsMap.Settings["daily_upload_limit"]; ok {
		if limit, ok := limitVal.(float64); ok {
			dailyLimit = int(limit)
		}
	}
	if dailyLimit == -1 {
		return false, nil
	}
	db := database.DB
	var todayCount int64
	startOfDay := time.Now().Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Second)
	err = db.Model(&models.File{}).Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, startOfDay, endOfDay).Count(&todayCount).Error
	if err != nil {
		return false, err
	}
	return int(todayCount)+uploadCount > dailyLimit, nil
}
