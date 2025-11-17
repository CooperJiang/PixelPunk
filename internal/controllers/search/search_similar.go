package search

import (
	"fmt"
	"pixelpunk/internal/controllers/search/dto"
	"pixelpunk/internal/middleware"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/vector"
	"time"

	"github.com/gin-gonic/gin"
)

func SearchSimilarFiles(c *gin.Context) {
	startTime := time.Now()

	fileID := c.Param("fileId")
	if fileID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件ID不能为空"))
		return
	}

	limit := 15

	threshold, _, err := getVectorConfig()
	if err != nil {
		logger.Warn("获取向量配置失败，使用默认阈值: %v", err)
		threshold = 0.7
	}

	engine := vector.GetGlobalVectorEngine()
	if engine == nil || !engine.IsEnabled() {
		errors.HandleError(c, errors.New(errors.CodeServiceUnavailable, "向量搜索功能不可用"))
		return
	}

	userIDUint := uint(0)
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(uint); ok {
			userIDUint = uid
		}
	}

	searchResults, err := engine.SearchSimilarByFileID(fileID, limit+1, userIDUint, threshold)
	if err != nil {
		logger.Error("相似文件搜索失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, fmt.Sprintf("搜索失败: %v", err)))
		return
	}

	db := database.GetDB()

	results := make([]dto.VectorSearchResult, 0, limit)
	for _, result := range searchResults {
		if result.FileID == fileID {
			continue
		}

		var file models.File
		if err := db.Where("id = ?", result.FileID).
			Where("status <> ?", "pending_deletion").
			First(&file).Error; err != nil {
			logger.Warn("文件信息查询失败 [%s]: %v", result.FileID, err)
			continue
		}

		fileInfo := dto.ConvertFileToInfo(&file)
		results = append(results, dto.VectorSearchResult{
			FileID:      result.FileID,
			Similarity:  result.Similarity,
			Description: result.Description,
			FileInfo:    fileInfo,
		})

		if len(results) >= limit {
			break
		}
	}

	processTime := time.Since(startTime)

	items := make([]map[string]interface{}, 0, len(results))
	for _, result := range results {
		if result.FileInfo != nil {
			item := map[string]interface{}{
				"id":             result.FileID,
				"original_name":  result.FileInfo.OriginalName,
				"display_name":   result.FileInfo.DisplayName,
				"size":           result.FileInfo.Size,
				"width":          result.FileInfo.Width,
				"height":         result.FileInfo.Height,
				"format":         result.FileInfo.Format,
				"url":            result.FileInfo.URL,
				"thumb_url":      result.FileInfo.ThumbURL,
				"full_url":       result.FileInfo.FullURL,
				"full_thumb_url": result.FileInfo.FullThumbURL,
				"created_at":     result.FileInfo.CreatedAt,
				"similarity":     result.Similarity,
				"size_formatted": result.FileInfo.SizeFormatted,
				"resolution":     result.FileInfo.Resolution,
			}
			items = append(items, item)
		}
	}

	response := gin.H{
		"items": items,
		"pagination": gin.H{
			"total":        int64(len(results)),
			"size":         limit,
			"current_page": 1,
			"last_page":    1,
		},
		"search_info": gin.H{
			"query":        fmt.Sprintf("similar_to:%s", fileID),
			"threshold":    threshold,
			"process_time": processTime.String(),
			"used_cache":   false,
		},
	}

	errors.ResponseSuccess(c, response, "相似文件搜索成功")
}

func GallerySimilarFiles(c *gin.Context) {
	startTime := time.Now()

	fileID := c.Param("fileId")
	if fileID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件ID不能为空"))
		return
	}

	limit := 15
	threshold, _, err := getVectorConfig()
	if err != nil {
		logger.Warn("获取向量配置失败，使用默认阈值: %v", err)
		threshold = 0.7
	}

	engine := vector.GetGlobalVectorEngine()
	if engine == nil || !engine.IsEnabled() {
		errors.HandleError(c, errors.New(errors.CodeServiceUnavailable, "向量搜索功能不可用"))
		return
	}

	searchResults, err := engine.SearchSimilarByFileID(fileID, limit*3, 0, threshold)
	if err != nil {
		logger.Error("Gallery相似文件搜索失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, fmt.Sprintf("搜索失败: %v", err)))
		return
	}

	db := database.GetDB()
	results := make([]dto.VectorSearchResult, 0, limit)

	for _, result := range searchResults {
		if result.FileID == fileID {
			continue
		}

		var file models.File
		if err := db.Where("id = ? AND access_level = ? AND is_recommended = ?",
			result.FileID, "public", true).
			Where("status <> ?", "pending_deletion").
			First(&file).Error; err != nil {
			continue
		}

		fileInfo := dto.ConvertFileToInfo(&file)
		results = append(results, dto.VectorSearchResult{
			FileID:      result.FileID,
			Similarity:  result.Similarity,
			Description: result.Description,
			FileInfo:    fileInfo,
		})

		if len(results) >= limit {
			break
		}
	}

	processTime := time.Since(startTime)

	items := make([]map[string]interface{}, 0, len(results))
	for _, result := range results {
		if result.FileInfo != nil {
			item := map[string]interface{}{
				"id":             result.FileID,
				"original_name":  result.FileInfo.OriginalName,
				"display_name":   result.FileInfo.DisplayName,
				"size":           result.FileInfo.Size,
				"width":          result.FileInfo.Width,
				"height":         result.FileInfo.Height,
				"format":         result.FileInfo.Format,
				"url":            result.FileInfo.URL,
				"thumb_url":      result.FileInfo.ThumbURL,
				"full_url":       result.FileInfo.FullURL,
				"full_thumb_url": result.FileInfo.FullThumbURL,
				"created_at":     result.FileInfo.CreatedAt,
				"similarity":     result.Similarity,
				"size_formatted": result.FileInfo.SizeFormatted,
				"resolution":     result.FileInfo.Resolution,
			}
			items = append(items, item)
		}
	}

	response := gin.H{
		"items": items,
		"pagination": gin.H{
			"total":        int64(len(results)),
			"size":         limit,
			"current_page": 1,
			"last_page":    1,
		},
		"search_info": gin.H{
			"query":        fmt.Sprintf("similar_to:%s", fileID),
			"threshold":    threshold,
			"process_time": processTime.String(),
			"used_cache":   false,
			"scope":        "gallery",
		},
	}

	errors.ResponseSuccess(c, response, "Gallery相似文件搜索成功")
}

func UserSimilarFiles(c *gin.Context) {
	startTime := time.Now()

	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "请先登录"))
		return
	}

	fileID := c.Param("fileId")
	if fileID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件ID不能为空"))
		return
	}

	limit := 15
	threshold, _, err := getVectorConfig()
	if err != nil {
		logger.Warn("获取向量配置失败，使用默认阈值: %v", err)
		threshold = 0.7
	}

	engine := vector.GetGlobalVectorEngine()
	if engine == nil || !engine.IsEnabled() {
		errors.HandleError(c, errors.New(errors.CodeServiceUnavailable, "向量搜索功能不可用"))
		return
	}

	searchResults, err := engine.SearchSimilarByFileID(fileID, limit+1, userID, threshold)
	if err != nil {
		logger.Error("用户相似文件搜索失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, fmt.Sprintf("搜索失败: %v", err)))
		return
	}

	db := database.GetDB()
	results := make([]dto.VectorSearchResult, 0, limit)

	for _, result := range searchResults {
		if result.FileID == fileID {
			continue
		}

		var file models.File
		if err := db.Where("id = ? AND user_id = ?", result.FileID, userID).
			Where("status <> ?", "pending_deletion").
			First(&file).Error; err != nil {
			continue
		}

		fileInfo := dto.ConvertFileToInfo(&file)
		results = append(results, dto.VectorSearchResult{
			FileID:      result.FileID,
			Similarity:  result.Similarity,
			Description: result.Description,
			FileInfo:    fileInfo,
		})

		if len(results) >= limit {
			break
		}
	}

	processTime := time.Since(startTime)

	items := make([]map[string]interface{}, 0, len(results))
	for _, result := range results {
		if result.FileInfo != nil {
			item := map[string]interface{}{
				"id":             result.FileID,
				"original_name":  result.FileInfo.OriginalName,
				"display_name":   result.FileInfo.DisplayName,
				"size":           result.FileInfo.Size,
				"width":          result.FileInfo.Width,
				"height":         result.FileInfo.Height,
				"format":         result.FileInfo.Format,
				"url":            result.FileInfo.URL,
				"thumb_url":      result.FileInfo.ThumbURL,
				"full_url":       result.FileInfo.FullURL,
				"full_thumb_url": result.FileInfo.FullThumbURL,
				"created_at":     result.FileInfo.CreatedAt,
				"similarity":     result.Similarity,
				"size_formatted": result.FileInfo.SizeFormatted,
				"resolution":     result.FileInfo.Resolution,
			}
			items = append(items, item)
		}
	}

	response := gin.H{
		"items": items,
		"pagination": gin.H{
			"total":        int64(len(results)),
			"size":         limit,
			"current_page": 1,
			"last_page":    1,
		},
		"search_info": gin.H{
			"query":        fmt.Sprintf("similar_to:%s", fileID),
			"threshold":    threshold,
			"process_time": processTime.String(),
			"used_cache":   false,
			"scope":        "user",
		},
	}

	errors.ResponseSuccess(c, response, "用户相似文件搜索成功")
}

func AdminSimilarFiles(c *gin.Context) {
	startTime := time.Now()

	userRole := middleware.GetCurrentUserRole(c)
	if userRole != common.UserRoleAdmin && userRole != common.UserRoleSuperAdmin {
		errors.HandleError(c, errors.New(errors.CodeForbidden, "需要管理员权限"))
		return
	}

	fileID := c.Param("fileId")
	if fileID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件ID不能为空"))
		return
	}

	limit := 15
	threshold, _, err := getVectorConfig()
	if err != nil {
		logger.Warn("获取向量配置失败，使用默认阈值: %v", err)
		threshold = 0.7
	}

	engine := vector.GetGlobalVectorEngine()
	if engine == nil || !engine.IsEnabled() {
		errors.HandleError(c, errors.New(errors.CodeServiceUnavailable, "向量搜索功能不可用"))
		return
	}

	searchResults, err := engine.SearchSimilarByFileID(fileID, limit+1, 0, threshold)
	if err != nil {
		logger.Error("管理员相似文件搜索失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, fmt.Sprintf("搜索失败: %v", err)))
		return
	}

	db := database.GetDB()

	type AdminSearchResult struct {
		dto.VectorSearchResult
		File *models.File `json:"-"`
	}

	results := make([]AdminSearchResult, 0, limit)

	for _, result := range searchResults {
		if result.FileID == fileID {
			continue
		}

		var file models.File
		if err := db.Where("id = ?", result.FileID).
			Where("status <> ?", "pending_deletion").
			First(&file).Error; err != nil {
			logger.Warn("文件信息查询失败 [%s]: %v", result.FileID, err)
			continue
		}

		fileInfo := dto.ConvertFileToInfo(&file)
		results = append(results, AdminSearchResult{
			VectorSearchResult: dto.VectorSearchResult{
				FileID:      result.FileID,
				Similarity:  result.Similarity,
				Description: result.Description,
				FileInfo:    fileInfo,
			},
			File: &file,
		})

		if len(results) >= limit {
			break
		}
	}

	processTime := time.Since(startTime)

	items := make([]map[string]interface{}, 0, len(results))
	for _, result := range results {
		if result.FileInfo != nil && result.File != nil {
			item := map[string]interface{}{
				"id":             result.FileID,
				"original_name":  result.FileInfo.OriginalName,
				"display_name":   result.FileInfo.DisplayName,
				"size":           result.FileInfo.Size,
				"width":          result.FileInfo.Width,
				"height":         result.FileInfo.Height,
				"format":         result.FileInfo.Format,
				"url":            result.FileInfo.URL,
				"thumb_url":      result.FileInfo.ThumbURL,
				"full_url":       result.FileInfo.FullURL,
				"full_thumb_url": result.FileInfo.FullThumbURL,
				"created_at":     result.FileInfo.CreatedAt,
				"similarity":     result.Similarity,
				"size_formatted": result.FileInfo.SizeFormatted,
				"resolution":     result.FileInfo.Resolution,
				"user_id":        result.File.UserID,
				"access_level":   result.File.AccessLevel,
				"is_recommended": result.File.IsRecommended,
				"folder_id":      result.File.FolderID,
			}
			items = append(items, item)
		}
	}

	response := gin.H{
		"items": items,
		"pagination": gin.H{
			"total":        int64(len(results)),
			"size":         limit,
			"current_page": 1,
			"last_page":    1,
		},
		"search_info": gin.H{
			"query":        fmt.Sprintf("similar_to:%s", fileID),
			"threshold":    threshold,
			"process_time": processTime.String(),
			"used_cache":   false,
			"scope":        "admin",
		},
	}

	errors.ResponseSuccess(c, response, "管理员相似文件搜索成功")
}
