package ai

import (
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"time"

	"gorm.io/gorm"
)

func GetAIQueueStats() (map[string]interface{}, error) {
	db := GetDBFromContext()
	if db == nil {
		return nil, fmt.Errorf("无法获取数据库连接")
	}

	type StatusCount struct {
		Status string
		Count  int64
	}
	var statusCounts []StatusCount
	if err := db.Table("file").Select("ai_tagging_status as status, count(*) as count").Group("ai_tagging_status").Find(&statusCounts).Error; err != nil {
		return nil, fmt.Errorf("查询文件状态统计失败: %v", err)
	}

	stats := map[string]int64{
		"none_count": 0, "pending_count": 0, "done_count": 0, "failed_count": 0, "skipped_count": 0, "ignored_count": 0, "total_count": 0,
	}
	for _, sc := range statusCounts {
		switch sc.Status {
		case common.AITaggingStatusNone:
			stats["none_count"] = sc.Count
		case common.AITaggingStatusPending:
			stats["pending_count"] = sc.Count
		case common.AITaggingStatusDone:
			stats["done_count"] = sc.Count
		case common.AITaggingStatusFailed:
			stats["failed_count"] = sc.Count
		case common.AITaggingStatusSkipped:
			stats["skipped_count"] = sc.Count
		case common.AITaggingStatusIgnored:
			stats["ignored_count"] = sc.Count
		}
		stats["total_count"] += sc.Count
	}

	var configuredConcurrency, activeWorkers, queueLength, inFlight int
	gs := GetGlobalTaggingService()
	if gs == nil {
		_ = InitGlobalTaggingQueue()
		gs = GetGlobalTaggingService()
	}
	if gs != nil {
		// 运行态来自服务本身
		gsStats := gs.GetQueueStats()
		if v, ok := gsStats["active_workers"].(int); ok {
			activeWorkers = v
		}
		if v, ok := gsStats["max_workers"].(int); ok {
			configuredConcurrency = v
		}
	}
	if gs != nil && gs.taskQueue != nil {
		if m, err := gs.taskQueue.Metrics(); err == nil && m != nil {
			queueLength = m.QueueLength
			inFlight = m.InFlight
		}
	}

	// 为了与前端 /admin/tagging WebSocket 适配，提供 *_count 命名和 total_count 字段
	result := map[string]interface{}{
		"queue_stats": map[string]int64{
			"total_count":   stats["total_count"],
			"none_count":    stats["none_count"],
			"pending_count": stats["pending_count"],
			"done_count":    stats["done_count"],
			"failed_count":  stats["failed_count"],
			"ignored_count": stats["ignored_count"],
			// 兼容保留（前端目前未使用）
			"skipped_count": stats["skipped_count"],
		},
		// 前端 /admin/tagging 期望的扩展字段命名：queued / processing
		"queue_stats_ext": map[string]int{
			"queued": queueLength,
			// 处理中数量以队列侧 InFlight 为准（更接近真实处理中的任务数）
			"processing": inFlight,
			"delayed":    0,
			"dlq":        0,
		},
		"config": map[string]interface{}{
			"current_concurrency":    configuredConcurrency,
			"max_concurrency":        configuredConcurrency,
			"active_workers":         activeWorkers,
			"configured_concurrency": configuredConcurrency,
			"paused":                 gs != nil && gs.IsPaused(),
			"recent_failures": func() int {
				if gs != nil {
					return gs.RecentFailures()
				}
				return 0
			}(),
			"queue_length":   queueLength,
			"queue_capacity": 0,
		},
		"runtime": map[string]interface{}{
			"active_workers":         activeWorkers,
			"configured_concurrency": configuredConcurrency,
			"paused":                 gs != nil && gs.IsPaused(),
		},
		"performance": map[string]interface{}{"processing_rate": 0, "average_duration": 0},
	}
	return result, nil
}

func CleanOldAIJobs(retentionDays int) (int64, error) {
	if retentionDays <= 0 {
		return 0, nil
	}
	db := database.GetDB()
	if db == nil {
		return 0, fmt.Errorf("无法获取数据库连接")
	}
	cutoff := time.Now().Add(-time.Duration(retentionDays) * 24 * time.Hour)
	res := db.Where("status IN (?) AND updated_at < ?", []string{"done", "failed", "ignored", "skipped"}, cutoff).Delete(&models.AIJob{})
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

// TriggerFullTaggingProcess 触发全量文件标记处理
func TriggerFullTaggingProcess(maxImages int) (int, error) {
	db := GetDBFromContext()
	if db == nil {
		return 0, fmt.Errorf("无法获取数据库连接")
	}
	globalService := GetGlobalTaggingService()
	if globalService == nil {
		if err := InitGlobalTaggingQueue(); err != nil {
			return 0, fmt.Errorf("初始化AI队列失败: %v", err)
		}
		globalService = GetGlobalTaggingService()
		if globalService == nil {
			return 0, fmt.Errorf("全局TaggingService未初始化，无法处理任务")
		}
	}

	// 解冻 pending → none 当队列空闲（无事务）
	if qs := globalService.GetQueueStats(); true {
		qlen, _ := qs["queue_length"].(int)
		active, _ := qs["active_workers"].(int)
		if qlen == 0 && active == 0 {
			var pendingCount int64
			if err := db.Table("file").Where("ai_tagging_status = ?", common.AITaggingStatusPending).Count(&pendingCount).Error; err == nil && pendingCount > 0 {
				_ = db.Model(&models.File{}).Where("ai_tagging_status = ?", common.AITaggingStatusPending).Updates(map[string]interface{}{
					"ai_tagging_status": common.AITaggingStatusNone, "ai_tagging_tries": 0,
				}).Error
			}
		}
	}

	// 解冻超上限的 tries>=3 任务和已跳过的任务，置 pending 并入队（无事务）
	// 手动触发时包含 Skipped 状态，允许用户修正配置后重新处理
	rescueLimit := 500
	if maxImages > 0 && maxImages < rescueLimit {
		rescueLimit = maxImages
	}
	var rescueIDs []string
	if err := db.Table("file").Where("ai_tagging_status IN ?", []string{common.AITaggingStatusNone, common.AITaggingStatusFailed, common.AITaggingStatusSkipped}).
		Where("ai_tagging_tries >= ?", 3).Order("updated_at ASC").Limit(rescueLimit).Pluck("id", &rescueIDs).Error; err == nil && len(rescueIDs) > 0 {
		_ = db.Model(&models.File{}).Where("id IN ?", rescueIDs).Updates(map[string]interface{}{
			"ai_tagging_status": common.AITaggingStatusPending, "ai_tagging_tries": 0,
		}).Error
		var rescueImages []models.File
		if err := db.Where("id IN ?", rescueIDs).Find(&rescueImages).Error; err == nil {
			if s := GetGlobalTaggingService(); s != nil {
				s.BatchProcessFiles(rescueImages)
			}
		}
	}

	// 常规 tries<3 任务入队（无事务，手动触发时包含已跳过的文件）
	var images []models.File
	query := db.Table("file").Where("ai_tagging_status IN ? AND ai_tagging_tries < 3", []string{common.AITaggingStatusNone, common.AITaggingStatusFailed, common.AITaggingStatusSkipped})
	if maxImages > 0 {
		query = query.Limit(maxImages)
	}
	query = query.Order("created_at ASC")
	var imageIDs []string
	if err := query.Pluck("id", &imageIDs).Error; err != nil {
		return 0, fmt.Errorf("查询待处理文件失败: %v", err)
	}
	if len(imageIDs) > 0 {
		_ = db.Where("file_id IN ? AND status IN ?", imageIDs, []string{"done", "ignored", "skipped"}).
			Delete(&models.AIJob{}).Error
		_ = db.Model(&models.File{}).Where("id IN ? AND ai_tagging_status IN ?", imageIDs, []string{common.AITaggingStatusNone, common.AITaggingStatusFailed, common.AITaggingStatusSkipped}).
			Update("ai_tagging_status", common.AITaggingStatusPending).Error
		if err := db.Where("id IN ?", imageIDs).Find(&images).Error; err != nil {
			return 0, fmt.Errorf("查询文件详情失败: %v", err)
		}
	}

	count := len(images)
	if count == 0 {
		return 0, nil
	}
	globalService.BatchProcessFiles(images)
	return count, nil
}

// ResetStuckPendingFiles 重置卡在pending状态的文件（改进版：使用心跳时间判断，无事务）
func ResetStuckPendingFiles(timeThresholdMinutes int) (int, error) {
	if timeThresholdMinutes <= 0 {
		timeThresholdMinutes = 30
	}
	db := GetDBFromContext()
	if db == nil {
		return 0, fmt.Errorf("无法获取数据库连接")
	}
	thresholdTime := time.Now().Add(-time.Duration(timeThresholdMinutes) * time.Minute)

	// 改进的卡住检测逻辑（无事务）：
	result := db.Model(&models.File{}).
		Where("ai_tagging_status = ?", common.AITaggingStatusPending).
		Where("(ai_last_heartbeat_at < ? OR (ai_last_heartbeat_at IS NULL AND updated_at < ?))",
			thresholdTime, thresholdTime).
		Updates(map[string]interface{}{
			"ai_tagging_status":    common.AITaggingStatusFailed,
			"ai_tagging_tries":     gorm.Expr("ai_tagging_tries + 1"),
			"ai_last_heartbeat_at": nil,
		})
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}
