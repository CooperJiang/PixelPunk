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
	"pixelpunk/pkg/storage"
	"pixelpunk/pkg/vector"
	"time"

	"github.com/gin-gonic/gin"
)

func UserVectorSearch(c *gin.Context) {
	startTime := time.Now()

	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "用户未认证"))
		return
	}

	req, err := common.ValidateRequest[dto.UserVectorSearchRequest](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}

	threshold := getSearchThreshold()
	_, maxResults, err := getVectorConfig()
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, fmt.Sprintf("获取向量配置失败: %v", err)))
		return
	}

	size := maxResults
	if size > 100 {
		size = 100
	}
	if size <= 0 {
		size = 20
	}

	engine := vector.GetEngine()
	if engine == nil || !engine.IsEnabled() {
		errors.HandleError(c, errors.New(errors.CodeServiceUnavailable, "向量搜索服务不可用"))
		return
	}

	offset := (page - 1) * size
	searchLimit := size * 2

	searchResults, err := engine.SearchFiles(req.Query, searchLimit, userID, threshold)
	if err != nil {
		logger.Error("用户向量搜索失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, fmt.Sprintf("搜索失败: %v", err)))
		return
	}

	db := database.DB

	var results []map[string]interface{}
	totalResults := len(searchResults)

	for i, result := range searchResults {
		if i < offset {
			continue
		}
		if len(results) >= size {
			break
		}

		var file models.File
		if err := db.Where("id = ? AND user_id = ?", result.FileID, userID).
			Where("status <> ?", "pending_deletion").
			First(&file).Error; err != nil {
			continue
		}

		fullURL, fullThumbURL, _ := storage.GetFullURLs(file)

		fileMap := map[string]interface{}{
			"id":                  file.ID,
			"created_at":          file.CreatedAt,
			"updated_at":          file.UpdatedAt,
			"user_id":             file.UserID,
			"folder_id":           file.FolderID,
			"original_name":       file.OriginalName,
			"display_name":        file.DisplayName,
			"file_name":           file.FileName,
			"file_path":           file.FilePath,
			"full_path":           file.FullPath,
			"local_file_path":     file.LocalFilePath,
			"local_thumb_path":    file.LocalThumbPath,
			"url":                 file.URL,
			"thumb_url":           file.ThumbURL,
			"full_url":            fullURL,
			"full_thumb_url":      fullThumbURL,
			"remote_url":          file.RemoteURL,
			"remote_thumb_url":    file.RemoteThumbURL,
			"md5_hash":            file.MD5Hash,
			"size":                file.Size,
			"size_formatted":      file.SizeFormatted,
			"width":               file.Width,
			"height":              file.Height,
			"ratio":               file.Ratio,
			"format":              file.Format,
			"mime":                file.Mime,
			"resolution":          file.Resolution,
			"description":         file.Description,
			"nsfw":                file.NSFW,
			"access_level":        file.AccessLevel,
			"is_duplicate":        file.IsDuplicate,
			"is_recommended":      file.IsRecommended,
			"api_key_id":          file.APIKeyID,
			"storage_provider_id": file.StorageProviderID,
			"storage_type":        file.StorageType,
			"ai_tagging_status":   file.AITaggingStatus,
			"ai_tagging_tries":    file.AITaggingTries,
			"ai_tagging_duration": file.AITaggingDuration,
			"ai_http_duration":    file.AIHttpDuration,
			"sort_order":          file.SortOrder,
			"similarity":          result.Similarity,
		}

		results = append(results, fileMap)
	}

	processTime := time.Since(startTime)

	data := gin.H{
		"items": results,
		"pagination": gin.H{
			"total":        int64(totalResults),
			"size":         size,
			"current_page": page,
			"last_page":    (int64(totalResults) + int64(size) - 1) / int64(size),
		},
		"search_info": gin.H{
			"query":        req.Query,
			"threshold":    threshold,
			"process_time": processTime.String(),
		},
	}

	errors.ResponseSuccess(c, data, "用户文件向量搜索成功")
}

func GalleryVectorSearch(c *gin.Context) {
	startTime := time.Now()

	req, err := common.ValidateRequest[dto.GalleryVectorSearchRequest](c)
	if err != nil {
		logger.Error("Gallery向量搜索参数验证失败: %v", err)
		errors.HandleError(c, err)
		return
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}

	threshold := getSearchThreshold()
	_, maxResults, err := getVectorConfig()
	if err != nil {
		logger.Error("Gallery向量搜索获取配置失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, fmt.Sprintf("获取向量配置失败: %v", err)))
		return
	}

	size := maxResults
	if size > 100 {
		size = 100
	}
	if size <= 0 {
		size = 20
	}

	engine := vector.GetEngine()
	if engine == nil {
		logger.Error("Gallery向量搜索: 向量引擎为空")
		errors.HandleError(c, errors.New(errors.CodeServiceUnavailable, "向量搜索服务不可用"))
		return
	}
	if !engine.IsEnabled() {
		logger.Error("Gallery向量搜索: 向量引擎未启用")
		errors.HandleError(c, errors.New(errors.CodeServiceUnavailable, "向量搜索服务不可用"))
		return
	}

	offset := (page - 1) * size
	searchLimit := size * 3

	searchResults, err := engine.SearchFiles(req.Query, searchLimit, 0, threshold)
	if err != nil {
		logger.Error("Gallery向量搜索失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, fmt.Sprintf("搜索失败: %v", err)))
		return
	}

	db := database.DB

	var results []map[string]interface{}
	var filteredResults []string

	publicCount := 0
	for _, result := range searchResults {

		var file models.File
		err := db.Where("id = ? AND access_level = ? AND is_recommended = ?",
			result.FileID, "public", true).First(&file).Error
		if err == nil {
			filteredResults = append(filteredResults, result.FileID)
			publicCount++
		} else {
		}
	}

	paginatedCount := 0

	for _, result := range searchResults {
		var file models.File
		err := db.Where("id = ? AND access_level = ? AND is_recommended = ?",
			result.FileID, "public", true).First(&file).Error
		if err != nil {
			continue
		}

		currentIndex := 0
		for _, filteredID := range filteredResults {
			if filteredID == result.FileID {
				break
			}
			currentIndex++
		}

		if currentIndex < offset {
			continue
		}
		if len(results) >= size {
			break
		}

		paginatedCount++

		fullURL, fullThumbURL, _ := storage.GetFullURLs(file)

		fileMap := map[string]interface{}{
			"id":                  file.ID,
			"created_at":          file.CreatedAt,
			"updated_at":          file.UpdatedAt,
			"user_id":             file.UserID,
			"folder_id":           file.FolderID,
			"original_name":       file.OriginalName,
			"display_name":        file.DisplayName,
			"file_name":           file.FileName,
			"file_path":           file.FilePath,
			"full_path":           file.FullPath,
			"local_file_path":     file.LocalFilePath,
			"local_thumb_path":    file.LocalThumbPath,
			"url":                 file.URL,
			"thumb_url":           file.ThumbURL,
			"full_url":            fullURL,
			"full_thumb_url":      fullThumbURL,
			"remote_url":          file.RemoteURL,
			"remote_thumb_url":    file.RemoteThumbURL,
			"md5_hash":            file.MD5Hash,
			"size":                file.Size,
			"size_formatted":      file.SizeFormatted,
			"width":               file.Width,
			"height":              file.Height,
			"ratio":               file.Ratio,
			"format":              file.Format,
			"mime":                file.Mime,
			"resolution":          file.Resolution,
			"description":         file.Description,
			"nsfw":                file.NSFW,
			"access_level":        file.AccessLevel,
			"is_duplicate":        file.IsDuplicate,
			"is_recommended":      file.IsRecommended,
			"api_key_id":          file.APIKeyID,
			"storage_provider_id": file.StorageProviderID,
			"storage_type":        file.StorageType,
			"ai_tagging_status":   file.AITaggingStatus,
			"ai_tagging_tries":    file.AITaggingTries,
			"ai_tagging_duration": file.AITaggingDuration,
			"ai_http_duration":    file.AIHttpDuration,
			"sort_order":          file.SortOrder,
			"similarity":          result.Similarity,
		}

		results = append(results, fileMap)
	}

	processTime := time.Since(startTime)
	totalResults := len(filteredResults)

	data := gin.H{
		"items": results,
		"pagination": gin.H{
			"total":        int64(totalResults),
			"size":         size,
			"current_page": page,
			"last_page":    (int64(totalResults) + int64(size) - 1) / int64(size),
		},
		"search_info": gin.H{
			"query":        req.Query,
			"threshold":    threshold,
			"process_time": processTime.String(),
		},
	}

	errors.ResponseSuccess(c, data, "Gallery向量搜索成功")
}
