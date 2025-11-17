package automation

import (
	"fmt"
	"pixelpunk/internal/controllers/automation/dto"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/storage"
	"time"

	"gorm.io/gorm"
)

/* 用户自动任务服务 */

func GetUserAutomationOverview(userID uint) (*dto.AutomationOverviewResponse, error) {
	db := database.GetDB()

	taggingStats, err := getUserTaggingStats(db, userID)
	if err != nil {
		return nil, fmt.Errorf("获取打标统计失败: %w", err)
	}

	vectorStats, err := getUserVectorStats(db, userID)
	if err != nil {
		return nil, fmt.Errorf("获取向量统计失败: %w", err)
	}

	systemStatus := getSystemStatus(db)

	return &dto.AutomationOverviewResponse{
		Tagging:      *taggingStats,
		Vector:       *vectorStats,
		SystemStatus: *systemStatus,
	}, nil
}

// getUserTaggingStats 获取用户打标任务统计
func getUserTaggingStats(db *gorm.DB, userID uint) (*dto.TaggingTaskStats, error) {
	stats := &dto.TaggingTaskStats{}

	var counts []struct {
		Status string
		Count  int64
	}

	err := db.Model(&models.File{}).
		Select("ai_tagging_status as status, COUNT(*) as count").
		Where("user_id = ? AND file_type = ?", userID, models.FileTypeImage).
		Group("ai_tagging_status").
		Find(&counts).Error

	if err != nil {
		return nil, err
	}

	for _, c := range counts {
		switch c.Status {
		case common.AITaggingStatusNone:
			stats.NoneCount = int(c.Count)
		case common.AITaggingStatusPending:
			stats.PendingCount = int(c.Count)
		case "processing":
			stats.ProcessingCount = int(c.Count)
		case common.AITaggingStatusDone:
			stats.DoneCount = int(c.Count)
		case common.AITaggingStatusFailed:
			stats.FailedCount = int(c.Count)
		case common.AITaggingStatusIgnored:
			stats.IgnoredCount = int(c.Count)
		}
		stats.TotalCount += int(c.Count)
	}

	return stats, nil
}

// getUserVectorStats 获取用户向量任务统计
func getUserVectorStats(db *gorm.DB, userID uint) (*dto.VectorTaskStats, error) {
	stats := &dto.VectorTaskStats{}

	var counts []struct {
		Status string
		Count  int64
	}

	err := db.Model(&models.FileVector{}).
		Select("file_vector.status as status, COUNT(*) as count").
		Joins("JOIN file ON file.id = file_vector.file_id").
		Where("file.user_id = ?", userID).
		Group("file_vector.status").
		Find(&counts).Error

	if err != nil {
		return nil, err
	}

	for _, c := range counts {
		switch c.Status {
		case common.VectorStatusPending:
			stats.PendingCount = int(c.Count)
		case "processing":
			stats.ProcessingCount = int(c.Count)
		case common.VectorStatusCompleted:
			stats.CompletedCount = int(c.Count)
		case common.VectorStatusFailed:
			stats.FailedCount = int(c.Count)
		case common.VectorStatusReset:
			stats.ResetCount = int(c.Count)
		}
		stats.TotalCount += int(c.Count)
	}

	return stats, nil
}

// getSystemStatus 获取系统队列状态
func getSystemStatus(db *gorm.DB) *dto.SystemStatus {
	status := &dto.SystemStatus{
		TaggingEnabled:     true, // 从配置或AI服务获取
		VectorEnabled:      true, // 从配置或向量服务获取
		TaggingConcurrency: 5,    // 从配置获取
		VectorConcurrency:  3,    // 从配置获取
	}

	// 统计总队列长度（pending + processing 状态的任务）
	var taggingQueueCount int64
	db.Model(&models.File{}).
		Where("ai_tagging_status IN ?", []string{common.AITaggingStatusPending, "processing"}).
		Count(&taggingQueueCount)

	var vectorQueueCount int64
	db.Model(&models.FileVector{}).
		Where("status IN ?", []string{common.VectorStatusPending, "processing"}).
		Count(&vectorQueueCount)

	status.TotalQueueLength = int(taggingQueueCount + vectorQueueCount)

	return status
}

func GetUserTaggingTasks(userID uint, query dto.TaskListQuery) (*dto.TaskListResponse, error) {
	db := database.GetDB()

	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 10
	}

	offset := (query.Page - 1) * query.Limit

	dbQuery := db.Model(&models.File{}).
		Where("user_id = ? AND file_type = ?", userID, models.FileTypeImage)

	if query.Status != "" {
		dbQuery = dbQuery.Where("ai_tagging_status = ?", query.Status)
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	var files []models.File
	err := dbQuery.
		Order("updated_at DESC").
		Offset(offset).
		Limit(query.Limit).
		Find(&files).Error

	if err != nil {
		return nil, err
	}

	items := make([]dto.TaggingTaskItem, 0, len(files))
	for _, file := range files {
		_, fullThumbURL, _ := storage.GetFullURLs(file)

		displayName := file.DisplayName
		if displayName == "" {
			displayName = file.OriginalName
		}

		items = append(items, dto.TaggingTaskItem{
			ID:            file.ID, // file.ID是string类型，直接使用
			FileID:        file.ID,
			FileName:      displayName,
			ThumbnailURL:  fullThumbURL,
			Status:        file.AITaggingStatus,
			Tries:         file.AITaggingTries,
			Format:        file.Format,
			SizeFormatted: file.SizeFormatted,
			Resolution:    file.Resolution,
			UpdatedAt:     time.Time(file.UpdatedAt).Format("2006-01-02 15:04:05"),
			CreatedAt:     time.Time(file.CreatedAt).Format("2006-01-02 15:04:05"),
		})
	}

	return &dto.TaskListResponse{
		Data:  items,
		Total: total,
		Page:  query.Page,
		Limit: query.Limit,
	}, nil
}

func GetUserVectorTasks(userID uint, query dto.TaskListQuery) (*dto.TaskListResponse, error) {
	db := database.GetDB()

	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 10
	}

	offset := (query.Page - 1) * query.Limit

	dbQuery := db.Model(&models.FileVector{}).
		Joins("JOIN file ON file.id = file_vector.file_id").
		Where("file.user_id = ?", userID)

	if query.Status != "" {
		dbQuery = dbQuery.Where("file_vector.status = ?", query.Status)
	}

	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	var vectors []models.FileVector
	err := dbQuery.
		Preload("File").
		Order("file_vector.updated_at DESC").
		Offset(offset).
		Limit(query.Limit).
		Find(&vectors).Error

	if err != nil {
		return nil, err
	}

	items := make([]dto.VectorTaskItem, 0, len(vectors))
	for _, vec := range vectors {
		item := dto.VectorTaskItem{
			ID:           int(vec.ID),
			FileID:       vec.FileID,
			Status:       vec.Status,
			Model:        vec.Model,
			RetryCount:   vec.RetryCount,
			ErrorMessage: vec.ErrorMessage,
			UpdatedAt:    time.Time(vec.UpdatedAt).Format("2006-01-02 15:04:05"),
			CreatedAt:    time.Time(vec.CreatedAt).Format("2006-01-02 15:04:05"),
		}

		if vec.File.ID != "" {
			_, fullThumbURL, _ := storage.GetFullURLs(vec.File)
			item.FileName = vec.File.OriginalName
			item.ThumbnailURL = fullThumbURL
		}

		items = append(items, item)
	}

	return &dto.TaskListResponse{
		Data:  items,
		Total: total,
		Page:  query.Page,
		Limit: query.Limit,
	}, nil
}
