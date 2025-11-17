package file

import (
	"encoding/json"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/activity"
	messageService "pixelpunk/internal/services/message"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/storage"

	"strings"

	"gorm.io/gorm"
)

// AdminFileSearchParams 管理员文件查询参数（语义化命名）
type AdminFileSearchParams struct {
	Page          int      // 页码
	Size          int      // 每页数量
	Keyword       string   // 关键字(匹配描述和标签)
	Tags          []string // 标签列表
	CategoryIDs   []string // 分类ID列表
	DominantColor []string // 主色调列表，支持多选
	Resolution    string   // 分辨率
	NSFWMinScore  float64  // NSFW最小评分
	NSFWMaxScore  float64  // NSFW最大评分
	IsNSFW        *bool    // 是否NSFW内容
	StorageType   string   // 存储类型
	Sort          string   // 排序方式
	MinWidth      int      // 最小宽度
	MaxWidth      int      // 最大宽度
	MinHeight     int      // 最小高度
	MaxHeight     int      // 最大高度
	UserID        uint     // 用户ID(可选)
	IsRecommended *bool    // 是否推荐内容(可选)
	FolderID      string   // 文件夹ID
	AccessLevel   string   // 访问级别
}

type AdminImageSearchParams = AdminFileSearchParams

// BatchOperationResult 批量操作结果
type BatchOperationResult struct {
	SuccessCount int                       `json:"success_count"`
	FailCount    int                       `json:"fail_count"`
	SuccessIDs   []string                  `json:"success_ids"`
	FailIDs      []string                  `json:"fail_ids"`
	Errors       map[string]string         `json:"errors,omitempty"`
	UpdatedFiles []AdminFileDetailResponse `json:"updated_files,omitempty"`
}

// BatchSetImageRecommendation 批量设置文件推荐状态
func BatchSetFileRecommendation(fileIDs []string, isRecommended bool) (*BatchOperationResult, error) {
	result := &BatchOperationResult{
		SuccessIDs:   make([]string, 0),
		FailIDs:      make([]string, 0),
		Errors:       make(map[string]string),
		UpdatedFiles: make([]AdminFileDetailResponse, 0),
	}

	var files []models.File
	if err := database.DB.Where("id IN ?", fileIDs).Find(&files).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件列表失败")
	}

	fileMap := make(map[string]*models.File)
	for i := range files {
		fileMap[files[i].ID] = &files[i]
	}

	aiInfoMap := batchGetFileAIInfo(fileIDs)

	var userIDs []uint
	for _, file := range files {
		if file.UserID > 0 {
			userIDs = append(userIDs, file.UserID)
		}
	}
	userMap := make(map[uint]string)
	if len(userIDs) > 0 {
		var users []models.User
		database.DB.Select("id, username").Where("id IN ?", userIDs).Find(&users)
		for _, user := range users {
			userMap[user.ID] = user.Username
		}
	}

	var stats []models.FileStats
	database.DB.Where("file_id IN ?", fileIDs).Find(&stats)
	statsMap := make(map[string]int64)
	for _, stat := range stats {
		statsMap[stat.FileID] = stat.Views
	}

	err := database.DB.Model(&models.File{}).
		Where("id IN ?", fileIDs).
		Update("is_recommended", isRecommended).Error

	if err != nil {
		return nil, errors.Wrap(err, errors.CodeDBUpdateFailed, "批量更新推荐状态失败")
	}

	for _, fileID := range fileIDs {
		file, exists := fileMap[fileID]
		if !exists {
			result.FailIDs = append(result.FailIDs, fileID)
			result.Errors[fileID] = "文件不存在"
			result.FailCount++
			continue
		}

		file.IsRecommended = isRecommended

		fullURL, fullThumbURL, _ := storage.GetFullURLs(*file)

		updatedFile := AdminFileDetailResponse{
			ID:                file.ID,
			URL:               file.URL,
			ThumbnailURL:      file.ThumbURL,
			FullURL:           fullURL,
			FullThumbURL:      fullThumbURL,
			OriginalName:      file.OriginalName,
			DisplayName:       file.DisplayName,
			Size:              file.Size,
			Width:             file.Width,
			Height:            file.Height,
			Format:            file.Format,
			AccessLevel:       file.AccessLevel,
			FolderID:          file.FolderID,
			CreatedAt:         file.CreatedAt,
			UpdatedAt:         file.UpdatedAt,
			Views:             statsMap[file.ID],
			IsRecommended:     file.IsRecommended,
			StorageProviderID: file.StorageProviderID,
			IsDuplicate:       file.IsDuplicate,
			MD5Hash:           file.MD5Hash,
			UserName:          userMap[file.UserID],
			AIInfo:            aiInfoMap[file.ID],
		}

		result.SuccessIDs = append(result.SuccessIDs, fileID)
		result.UpdatedFiles = append(result.UpdatedFiles, updatedFile)
		result.SuccessCount++
	}

	return result, nil
}

func batchGetFileAIInfo(fileIDs []string) map[string]*AIInfoResponse {
	result := make(map[string]*AIInfoResponse)

	if len(fileIDs) == 0 {
		return result
	}

	var aiInfos []models.FileAIInfo
	database.DB.Where("file_id IN ?", fileIDs).Find(&aiInfos)

	aiInfoMap := make(map[string]*models.FileAIInfo)
	for i := range aiInfos {
		aiInfoMap[aiInfos[i].FileID] = &aiInfos[i]
	}

	type TagRelation struct {
		FileID  string
		TagName string
	}
	var tagRelations []TagRelation
	database.DB.Table("file_global_tag_relation").
		Select("file_global_tag_relation.file_id, global_tag.name as tag_name").
		Joins("JOIN global_tag ON global_tag.id = file_global_tag_relation.tag_id").
		Where("file_global_tag_relation.file_id IN ?", fileIDs).
		Scan(&tagRelations)

	fileTagsMap := make(map[string][]string)
	for _, rel := range tagRelations {
		fileTagsMap[rel.FileID] = append(fileTagsMap[rel.FileID], rel.TagName)
	}

	for _, fileID := range fileIDs {
		aiInfo, exists := aiInfoMap[fileID]
		if !exists {
			continue
		}

		var tagNames []string
		if len(aiInfo.Tags) > 0 {
			var tags []string
			if err := json.Unmarshal(aiInfo.Tags, &tags); err == nil && len(tags) > 0 {
				tagNames = tags
			}
		}

		if len(tagNames) == 0 {
			tagNames = fileTagsMap[fileID]
		}

		var colorPalette []string
		if len(aiInfo.ColorPalette) > 0 {
			var colors []string
			if err := json.Unmarshal(aiInfo.ColorPalette, &colors); err == nil {
				colorPalette = colors
			}
		}

		var nsfwCategories map[string]float64
		if len(aiInfo.NSFWCategories) > 0 {
			var categories map[string]float64
			if err := json.Unmarshal(aiInfo.NSFWCategories, &categories); err == nil {
				nsfwCategories = categories
			}
		}

		result[fileID] = &AIInfoResponse{
			Description:    aiInfo.Description,
			Tags:           tagNames,
			DominantColor:  aiInfo.DominantColor,
			Resolution:     aiInfo.Resolution,
			IsNSFW:         aiInfo.IsNSFW,
			NSFWScore:      aiInfo.NSFWScore,
			NSFWEvaluation: aiInfo.NSFWEvaluation,
			ColorPalette:   colorPalette,
			AspectRatio:    aiInfo.AspectRatio,
			Composition:    aiInfo.Composition,
			ObjectsCount:   aiInfo.ObjectsCount,
			NSFWCategories: nsfwCategories,
		}
	}

	return result
}

func BatchDeleteFiles(fileIDs []string) (*BatchOperationResult, error) {
	result := &BatchOperationResult{
		SuccessIDs: make([]string, 0),
		FailIDs:    make([]string, 0),
		Errors:     make(map[string]string),
	}

	if len(fileIDs) == 0 {
		return result, nil
	}

	var files []models.File
	if err := database.DB.Where("id IN ?", fileIDs).Find(&files).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件列表失败")
	}

	fileMap := make(map[string]*models.File)
	for i := range files {
		fileMap[files[i].ID] = &files[i]
	}

	for _, fileID := range fileIDs {
		if _, exists := fileMap[fileID]; !exists {
			result.FailIDs = append(result.FailIDs, fileID)
			result.Errors[fileID] = "文件不存在"
			result.FailCount++
		} else {
			result.SuccessIDs = append(result.SuccessIDs, fileID)
		}
	}

	if len(files) == 0 {
		return result, nil
	}

	validFileIDs := make([]string, 0, len(files))
	for _, f := range files {
		validFileIDs = append(validFileIDs, f.ID)
	}

	database.DB.Where("file_id IN ?", validFileIDs).Delete(&models.FileGlobalTagRelation{})
	database.DB.Where("file_id IN ?", validFileIDs).Delete(&models.FileAIInfo{})
	database.DB.Unscoped().Where("file_id IN ?", validFileIDs).Delete(&models.FileStats{})
	database.DB.Unscoped().Where("file_id IN ?", validFileIDs).Delete(&models.FileVector{})

	userFileMap := make(map[uint][]models.File)
	categoryFiles := make(map[uint]int)

	for _, file := range files {
		userFileMap[file.UserID] = append(userFileMap[file.UserID], file)
		if file.CategoryID != nil && *file.CategoryID > 0 {
			categoryFiles[*file.CategoryID]++
		}
	}

	for userID, userFiles := range userFileMap {
		var totalSize int64
		var totalCount int
		for _, f := range userFiles {
			totalSize += f.Size
			totalCount++
		}

		database.DB.Model(&models.UserUsageStats{}).
			Where("user_id = ?", userID).
			Updates(map[string]interface{}{
				"total_images": gorm.Expr("CASE WHEN total_images >= ? THEN total_images - ? ELSE 0 END", totalCount, totalCount),
				"total_size":   gorm.Expr("CASE WHEN total_size >= ? THEN total_size - ? ELSE 0 END", totalSize, totalSize),
			})
	}

	for categoryID, count := range categoryFiles {
		database.DB.Model(&models.FileCategory{}).
			Where("id = ?", categoryID).
			UpdateColumn("file_count", gorm.Expr("CASE WHEN file_count >= ? THEN file_count - ? ELSE 0 END", count, count))
	}

	if err := database.DB.Unscoped().Where("id IN ?", validFileIDs).Delete(&models.File{}).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBDeleteFailed, "批量删除文件记录失败")
	}

	result.SuccessCount = len(validFileIDs)

	go func() {
		for _, file := range files {
			var duplicateCount int64
			if file.MD5Hash != "" {
				database.DB.Model(&models.File{}).
					Where("md5_hash = ? AND id != ?", file.MD5Hash, file.ID).
					Count(&duplicateCount)
			}
			var referenceCount int64
			database.DB.Model(&models.File{}).
				Where("original_file_id = ?", file.ID).
				Count(&referenceCount)
			cleanupFileResources(file.ID, file, duplicateCount+referenceCount)
		}
	}()

	for userID, deletedFiles := range userFileMap {
		go sendBatchDeleteNotification(userID, deletedFiles)
		go func(uid uint, images []models.File) {
			for _, img := range images {
				activity.LogAdminDelete(uid, img.ID, 0)
			}
		}(userID, deletedFiles)
	}

	return result, nil
}

// sendBatchDeleteNotification 发送批量删除文件通知
func sendBatchDeleteNotification(userID uint, deletedFiles []models.File) {
	if len(deletedFiles) == 0 {
		return
	}

	fileNames := make([]string, len(deletedFiles))
	for i, file := range deletedFiles {
		fileNames[i] = file.OriginalName
	}
	fileNamesStr := strings.Join(fileNames, "、")

	// 限制显示的文件名称长度，避免消息过长
	if len(fileNamesStr) > 200 {
		fileNamesStr = fileNamesStr[:200] + "..."
	}

	variables := map[string]interface{}{
		"count":        len(deletedFiles),
		"image_names":  fileNamesStr,
		"reason":       "管理员删除",
		"related_type": "file",
		"related_id":   "batch_delete",
	}

	var messageType string
	if len(deletedFiles) == 1 {
		messageType = common.MessageTypeFileDeletedByAdmin
		variables["image_name"] = deletedFiles[0].OriginalName
		variables["related_id"] = deletedFiles[0].ID
	} else {
		messageType = common.MessageTypeFileBatchDeletedByAdmin
	}

	msgService := messageService.GetMessageService()
	if err := msgService.SendTemplateMessage(userID, messageType, variables); err != nil {
		logger.Warn("发送文件删除消息失败: userID=%d, count=%d, error=%v", userID, len(deletedFiles), err)
	} else {
	}
}
