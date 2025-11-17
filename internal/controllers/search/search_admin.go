package search

import (
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/setting"
	vectorSvc "pixelpunk/internal/services/vector"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"

	"github.com/gin-gonic/gin"
)

func RegenerateVectors(c *gin.Context) {

	db := database.GetDB()

	// 直接从数据库读取配置（绕过缓存）
	currentModel := setting.GetStringDirectFromDB("vector", "vector_model", "text-embedding-3-small")

	if singleID := c.Query("file_id"); singleID != "" {
		updated := db.Model(&models.FileVector{}).Where("file_id = ?", singleID).
			Updates(map[string]interface{}{"status": common.VectorStatusReset, "model": currentModel})
		if updated.Error != nil {
			logger.Error("单文件更新向量状态失败: %v", updated.Error)
			errors.HandleError(c, errors.New(errors.CodeDBUpdateFailed, "更新向量状态失败"))
			return
		}
		created := 0
		if updated.RowsAffected == 0 {
			var aiInfo models.FileAIInfo
			if err := db.Where("file_id = ?", singleID).First(&aiInfo).Error; err == nil && aiInfo.Description != "" {
				fv := &models.FileVector{FileID: singleID, Description: aiInfo.Description, Model: currentModel, Status: common.VectorStatusReset}
				if err := db.Create(fv).Error; err == nil {
					created = 1
				}
			}
		}
		if svc := vectorSvc.GetGlobalVectorQueueService(); svc != nil {
			_ = svc.EnqueueVector(singleID)
		}
		errors.ResponseSuccess(c, map[string]interface{}{
			"total":         1,
			"updated_count": 1 - created,
			"created_count": created,
			"model":         currentModel,
			"message":       "已重置该文件向量，队列将自动处理",
		}, "单文件向量重置成功")
		return
	}

	var totalCount int64
	if err := db.Table("file i").
		Joins("JOIN file_ai_info ai ON i.id = ai.file_id").
		Where("ai.description IS NOT NULL AND ai.description != ''").
		Count(&totalCount).Error; err != nil {
		logger.Error("统计文件总数失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeDBQueryFailed, "统计文件失败"))
		return
	}

	if totalCount == 0 {
		logger.Warn("没有找到需要生成向量的文件")
		errors.ResponseSuccess(c, map[string]interface{}{
			"total":   0,
			"message": "没有找到需要生成向量的文件",
		}, "重新生成向量完成")
		return
	}

	result := db.Model(&models.FileVector{}).
		Where("file_id IN (SELECT i.id FROM file i JOIN file_ai_info ai ON i.id = ai.file_id WHERE ai.description IS NOT NULL AND ai.description != '')").
		Updates(map[string]interface{}{
			"status": common.VectorStatusReset,
			"model":  currentModel,
		})

	if result.Error != nil {
		logger.Error("更新现有向量状态失败: %v", result.Error)
		errors.HandleError(c, errors.New(errors.CodeDBUpdateFailed, "更新向量状态失败"))
		return
	}

	updatedCount := result.RowsAffected

	var missingFiles []models.File
	if err := db.Table("file i").
		Joins("JOIN file_ai_info ai ON i.id = ai.file_id").
		Where("ai.description IS NOT NULL AND ai.description != ''").
		Where("i.id NOT IN (SELECT file_id FROM file_vector)").
		Find(&missingFiles).Error; err != nil {
		logger.Error("查询缺失向量记录的文件失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeDBQueryFailed, "查询文件失败"))
		return
	}

	createdCount := 0
	for _, file := range missingFiles {
		var aiInfo models.FileAIInfo
		if err := db.Where("file_id = ?", file.ID).First(&aiInfo).Error; err != nil {
			continue
		}

		resetVector := &models.FileVector{
			FileID:      file.ID,
			Description: aiInfo.Description,
			Model:       currentModel,
			Status:      common.VectorStatusReset,
		}

		if err := db.Create(resetVector).Error; err != nil {
			logger.Warn("为文件 %s 创建重置向量记录失败: %v", file.ID, err)
		} else {
			createdCount++
		}
	}

	queueStats := vectorSvc.GetGlobalVectorQueueService().GetQueueStats()

	response := map[string]interface{}{
		"total":         totalCount,
		"updated_count": updatedCount,
		"created_count": createdCount,
		"model":         currentModel,
		"queue_stats":   queueStats,
		"message":       fmt.Sprintf("已将 %d 个文件标记为待重新生成向量状态 (模型: %s)，队列将自动处理", totalCount, currentModel),
	}

	errors.ResponseSuccess(c, response, "全量向量重新生成任务已启动")
}

func GetVectorQueueStatus(c *gin.Context) {
	stats := vectorSvc.GetGlobalVectorQueueService().GetQueueStats()
	errors.ResponseSuccess(c, stats, "获取队列状态成功")
}

func GetVectorTaskStatus(c *gin.Context) {
	taskID := c.Param("task_id")
	if taskID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "任务ID不能为空"))
		return
	}

	errors.HandleError(c, errors.New(errors.CodeNotFound, "任务不存在"))
}

func DebugVectorStatus(c *gin.Context) {

	db := database.GetDB()

	var stats = make(map[string]interface{})

	var totalCount int64
	db.Model(&models.FileVector{}).Count(&totalCount)
	stats["total_vectors"] = totalCount

	var completedCount, pendingCount, failedCount, resetCount int64
	db.Model(&models.FileVector{}).Where("status = ?", common.VectorStatusCompleted).Count(&completedCount)
	db.Model(&models.FileVector{}).Where("status = ?", common.VectorStatusPending).Count(&pendingCount)
	db.Model(&models.FileVector{}).Where("status = ?", common.VectorStatusFailed).Count(&failedCount)
	db.Model(&models.FileVector{}).Where("status = ?", common.VectorStatusReset).Count(&resetCount)

	stats["completed_count"] = completedCount
	stats["pending_count"] = pendingCount
	stats["failed_count"] = failedCount
	stats["reset_count"] = resetCount

	var modelStats []map[string]interface{}
	// 使用 GORM ORM 方法替代 Raw SQL，确保 SQLite 兼容性
	var results []struct {
		Model  string
		Status string
		Count  int64
	}

	if err := db.Model(&models.FileVector{}).
		Select("model, status, COUNT(*) as count").
		Group("model, status").
		Order("model, status").
		Scan(&results).Error; err == nil {
		for _, r := range results {
			modelStats = append(modelStats, map[string]interface{}{
				"model":  r.Model,
				"status": r.Status,
				"count":  r.Count,
			})
		}
	}
	stats["model_distribution"] = modelStats

	var publicVectorCount int64
	// 使用 GORM ORM 方法替代 Raw SQL，确保 SQLite 兼容性
	db.Model(&models.FileVector{}).
		Joins("JOIN file ON file.id = file_vector.file_id").
		Where("file_vector.status = ? AND file.access_level = ? AND file.is_recommended = ?",
			common.VectorStatusCompleted, "public", true).
		Count(&publicVectorCount)
	stats["public_recommended_vectors"] = publicVectorCount

	var recentVectors []models.FileVector
	db.Order("updated_at DESC").Limit(5).Find(&recentVectors)

	var recentSamples []map[string]interface{}
	for _, v := range recentVectors {
		recentSamples = append(recentSamples, map[string]interface{}{
			"file_id":     v.FileID,
			"model":       v.Model,
			"status":      v.Status,
			"dimension":   v.Dimension,
			"updated_at":  v.UpdatedAt,
			"retry_count": v.RetryCount,
		})
	}
	stats["recent_samples"] = recentSamples

	// 直接从数据库读取向量配置（绕过缓存）
	vectorConfig, err := setting.GetMultipleSettingsDirectFromDB("vector", []string{
		"vector_enabled", "vector_api_key", "vector_base_url", "vector_model",
		"vector_concurrency", "vector_auto_processing_enabled",
	})
	if err == nil {
		stats["current_config"] = vectorConfig
	}

	errors.ResponseSuccess(c, stats, "向量状态调试完成")
}
