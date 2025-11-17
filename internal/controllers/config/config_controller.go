package config

import (
	settingService "pixelpunk/internal/services/setting"
	"pixelpunk/internal/services/storage"
	"pixelpunk/pkg/config"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/imagex/formats"
	"pixelpunk/pkg/logger"

	"github.com/gin-gonic/gin"
)

func GetUploadConfig(c *gin.Context) {
	uploadConfig := config.GetUploadConfig()

	uploadSettingsResp, err := settingService.GetSettingsByGroupAsMap("upload")
	if err != nil {
		response := map[string]interface{}{
			"max_file_size":                 uploadConfig.MaxFileSize / (1024 * 1024), // 转换为MB
			"allowed_types":                 uploadConfig.AllowedTypes,
			"client_max_concurrent_uploads": 3, // 默认客户端最大并发数
			"chunked_upload": map[string]interface{}{
				"enabled":               true,                    // 默认启用
				"is_allow_chunk_upload": isLocalStorageDefault(), // 是否允许分片上传（仅local渠道允许）
				"threshold":             20,                      // 20MB
				"chunk_size":            2,                       // 2MB
				"max_chunk_size":        10,                      // 10MB
				"min_chunk_size":        1,                       // 1MB
				"max_concurrency":       5,                       // 默认并发数
				"session_timeout":       24,                      // 24小时
				"retry_count":           3,                       // 重试次数
				"cleanup_interval":      60,                      // 清理间隔（分钟）
			},
		}
		errors.ResponseSuccess(c, response, "获取上传配置成功")
		return
	}

	settingsMap := uploadSettingsResp.Settings

	response := map[string]interface{}{
		"max_file_size":                 getSettingValueOrDefault(settingsMap, "max_file_size", uploadConfig.MaxFileSize/(1024*1024)),
		"allowed_types":                 getSettingValueOrDefault(settingsMap, "allowed_file_formats", uploadConfig.AllowedTypes),
		"client_max_concurrent_uploads": getSettingValueOrDefault(settingsMap, "client_max_concurrent_uploads", 3),
		"chunked_upload": map[string]interface{}{
			"enabled":               getSettingValueOrDefault(settingsMap, "chunked_upload_enabled", true),
			"is_allow_chunk_upload": isLocalStorageDefault(),                                        // 是否允许分片上传（仅local渠道允许）
			"threshold":             getSettingValueOrDefault(settingsMap, "chunked_threshold", 20), // MB
			"chunk_size":            getSettingValueOrDefault(settingsMap, "chunk_size", 2),         // MB
			"max_chunk_size":        10,                                                             // 使用系统默认值 10MB
			"min_chunk_size":        1,                                                              // 使用系统默认值 1MB
			"max_concurrency":       getSettingValueOrDefault(settingsMap, "max_concurrency", 5),
			"session_timeout":       getSettingValueOrDefault(settingsMap, "session_timeout", 24), // 小时
			"retry_count":           getSettingValueOrDefault(settingsMap, "retry_count", 3),
			"cleanup_interval":      getSettingValueOrDefault(settingsMap, "cleanup_interval", 60), // 分钟
		},
	}

	errors.ResponseSuccess(c, response, "获取上传配置成功")
}

func GetUploadCapabilities(c *gin.Context) {
	caps := map[string]interface{}{
		"supported_extensions":          formats.SupportedExtensionsWithoutDot(),
		"supported_extensions_with_dot": formats.SupportedExtensionsWithDot(),
		"mime_map": map[string]string{
			"jpg":  formats.GetContentType("jpg"),
			"jpeg": formats.GetContentType("jpeg"),
			"png":  formats.GetContentType("png"),
			"gif":  formats.GetContentType("gif"),
			"webp": formats.GetContentType("webp"),
			"bmp":  formats.GetContentType("bmp"),
			"svg":  formats.GetContentType("svg"),
			"ico":  formats.GetContentType("ico"),
			"apng": formats.GetContentType("apng"),
			"jp2":  formats.GetContentType("jp2"),
			"tiff": formats.GetContentType("tiff"),
			"tif":  formats.GetContentType("tif"),
			"tga":  formats.GetContentType("tga"),
			"heic": formats.GetContentType("heic"),
			"heif": formats.GetContentType("heif"),
		},
		"thumbnail": map[string]interface{}{
			"rasterize_svg":        true,
			"webp_conversion":      true,
			"transparent_preserve": true,
		},
	}
	errors.ResponseSuccess(c, caps, "获取上传能力成功")
}

func getSettingValueOrDefault(settingsMap map[string]interface{}, key string, defaultValue interface{}) interface{} {
	if value, exists := settingsMap[key]; exists && value != nil {
		return value
	}
	return defaultValue
}

func isLocalStorageDefault() bool {
	defaultChannel, err := storage.GetDefaultChannel()
	if err != nil {
		logger.Error("获取默认存储渠道失败: %v", err)
		return false
	}

	return defaultChannel.Type == "local"
}
