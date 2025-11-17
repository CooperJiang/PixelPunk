package ai

import (
	"encoding/json"
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/cache"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/utils"
	"time"

	"gorm.io/gorm"
)

// 获取打标状态列表
// status: 过滤状态
// page: 页码
// limit: 每页数量
// orderBy: 排序字段
// order: 排序方向
func GetTaggingStatus(status string, page, limit int, orderBy, order string) (map[string]interface{}, error) {
	db := database.GetDB()

	query := db.Model(&models.File{})

	if status != "" {
		query = query.Where("ai_tagging_status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询打标状态总数失败")
	}

	validOrderFields := map[string]bool{
		"created_at":          true,
		"updated_at":          true,
		"ai_tagging_tries":    true,
		"ai_tagging_status":   true,
		"ai_tagging_duration": true,
	}

	if orderBy == "" || !validOrderFields[orderBy] {
		orderBy = "updated_at"
	}

	if order != "asc" && order != "desc" {
		order = "desc"
	}

	var images []models.File
	offset := (page - 1) * limit
	if err := query.Order(orderBy + " " + order).
		Offset(offset).
		Limit(limit).
		Find(&images).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询打标状态列表失败")
	}

	// 为文件添加完整URL信息
	result := make([]map[string]interface{}, len(images))
	signer := utils.GetURLSigner()

	for i, file := range images {
		imgJSON, _ := json.Marshal(file)
		imgMap := make(map[string]interface{})
		json.Unmarshal(imgJSON, &imgMap)

		if file.ID != "" {
			signedURL := signer.SignFileURL(file.ID, utils.SIGNATURE_DURATION)
			signedThumbURL := signer.SignThumbURL(file.ID, utils.SIGNATURE_DURATION)

			imgMap["full_url"] = signedURL
			imgMap["full_thumb_url"] = signedThumbURL
		}

		if file.AITaggingDuration > 0 {
			imgMap["ai_tagging_duration_formatted"] = fmt.Sprintf("%.2f秒", float64(file.AITaggingDuration)/1000.0)
		}

		if file.AIHttpDuration > 0 {
			imgMap["ai_http_duration_formatted"] = fmt.Sprintf("%.2f秒", float64(file.AIHttpDuration)/1000.0)
		}

		result[i] = imgMap
	}

	return map[string]interface{}{
		"total": total,
		"page":  page,
		"limit": limit,
		"data":  result,
	}, nil
}

// 获取打标统计信息
func GetTaggingStats() (map[string]interface{}, error) {
	cacheKey := "dashboard:tagging_stats"

	// 先尝试从缓存读取（5s 缓存）
	if cached, err := cache.Get(cacheKey); err == nil && cached != "" {
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return result, nil
		}
	}

	db := database.GetDB()

	type StatusStats struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}

	var statusStats []StatusStats
	err := db.Model(&models.File{}).
		Select("ai_tagging_status as status, count(*) as count").
		Group("ai_tagging_status").
		Find(&statusStats).Error

	if err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询打标状态统计失败")
	}

	var recentLogs []models.FileTaggingLog
	err = db.Order("created_at DESC").
		Limit(10).
		Find(&recentLogs).Error

	if err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询打标日志失败")
	}

	result := map[string]interface{}{
		"status_stats": statusStats,
		"recent_logs":  recentLogs,
	}

	// 写入缓存（5秒）
	if data, err := json.Marshal(result); err == nil {
		cache.Set(cacheKey, string(data), 5*time.Second)
	}

	return result, nil
}


// 重置卡住的打标任务
func ResetStuckTaggingTasks(timeThresholdMinutes int, operatorID uint) (int, error) {
	if timeThresholdMinutes < 5 {
		timeThresholdMinutes = 5 // 最小5分钟
	}

	count, err := ResetStuckPendingFiles(timeThresholdMinutes)
	if err != nil {
		return 0, errors.Wrap(err, errors.CodeInternal, "重置卡住的打标任务失败")
	}

	db := database.GetDB()
	data, _ := json.Marshal(map[string]interface{}{
		"count":          count,
		"time_threshold": timeThresholdMinutes,
	})

	logEntry := models.FileTaggingLog{
		Status:     common.TaggingStatusReset,
		Action:     common.TaggingActionManual,
		Type:       "tagging.reset_stuck",
		Data:       string(data),
		OperatorID: operatorID,
	}
	db.Create(&logEntry)

	return count, nil
}

// ScheduledTaggingTask 定时打标任务
func ScheduledTaggingTask() {
	// 新逻辑：需要 AI 分析开关 与 队列开关（二者均为 true）
	if !utils.GetAiAnalysisEnabled() {
		return
	}

	db := database.GetDB()

	startTime := time.Now()

	resetCount, resetErr := ResetStuckPendingFiles(30) // 30分钟视为卡住

	// 扫描遗漏的文件并添加到全局队列中（暂停时不扫描）
	globalService := GetGlobalTaggingService()
	var triggerErr error

	if globalService != nil {
		if globalService.IsPaused() {
		} else {
			// 扫描所有（分批）
			if _, e := EnqueueAllPending(1000); e != nil {
				triggerErr = e
			}
		}
	} else {
		triggerErr = fmt.Errorf("全局TaggingService未初始化")
	}

	// 每15分钟进行一次对账式补入队（幂等）
	if time.Now().Unix()%900 < 60 { // 近似每15分钟窗口内执行一次
		if err := ReconcileAIJobs(); err != nil {
			logger.Warn("对账补入队失败: %v", err)
		}
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	logData := map[string]interface{}{
		"duration":    duration.Seconds(),
		"reset_count": resetCount,
	}

	if resetErr != nil || triggerErr != nil {
		logData["errors"] = fmt.Sprintf("%v %v", resetErr, triggerErr)
	}

	data, _ := json.Marshal(logData)
	logEntry := models.FileTaggingLog{
		Status: common.TaggingStatusScheduled,
		Action: common.TaggingActionAuto,
		Type:   "tagging.scheduled_scan",
		Data:   string(data),
	}
	db.Create(&logEntry)
}

// EnqueueAllPending 扫描所有需要处理的文件并幂等入队（按批，限制最大扫描次数）
func EnqueueAllPending(batch int) (int, error) {
	if batch <= 0 {
		batch = 1000
	}
	db := database.GetDB()
	svc := GetGlobalTaggingService()
	if db == nil || svc == nil || svc.taskQueue == nil {
		return 0, nil
	}
	if svc.IsPaused() {
		return 0, nil
	}
	if !utils.GetAiAnalysisEnabled() {
		return 0, nil
	}

	total := 0
	maxRounds := 10 // 最多扫描10轮,避免无限循环
	for round := 0; round < maxRounds; round++ {
		var ids []string
		if err := db.Table("file").
			Where("ai_tagging_status IN ? AND ai_tagging_tries < 3",
				[]string{common.AITaggingStatusNone, common.AITaggingStatusFailed}).
			Order("created_at ASC").
			Limit(batch).
			Pluck("id", &ids).Error; err != nil {
			return total, err
		}
		if len(ids) == 0 {
			break
		}
		if len(ids) > 0 {
			_ = db.Where("file_id IN ? AND status IN ?", ids, []string{"done", "ignored", "skipped"}).
				Delete(&models.AIJob{}).Error
		}
		// 然后更新file状态为pending
		if len(ids) > 0 {
			_ = db.Model(&models.File{}).Where("id IN ?", ids).
				Update("ai_tagging_status", common.AITaggingStatusPending).Error
		}
		for _, id := range ids {
			_ = svc.taskQueue.EnqueueUnique(id, 0)
		}
		total += len(ids)
		if len(ids) < batch {
			break
		}
		if svc.IsPaused() {
			break
		}
	}
	return total, nil
}

// ReconcileAIJobs: 将可能遗漏的文件幂等入队（none/failed，但排除已失败3次和已跳过的）
func ReconcileAIJobs() error {
	db := database.GetDB()
	svc := GetGlobalTaggingService()
	if db == nil || svc == nil || svc.taskQueue == nil {
		return nil
	}

	var ids []string
	if err := db.Table("file").
		Where("ai_tagging_status IN ? AND ai_tagging_tries < 3", []string{common.AITaggingStatusNone, common.AITaggingStatusFailed}).
		Limit(200).
		Pluck("id", &ids).Error; err != nil {
		return err
	}
	for _, id := range ids {
		_ = svc.taskQueue.EnqueueUnique(id, 0)
	}
	return nil
}

// IgnoreTagging 批量将文件标记为已忽略
func IgnoreTagging(imageIDs []string, reason string, operatorID uint) (int, []string, error) {
	if len(imageIDs) == 0 {
		return 0, nil, errors.New(errors.CodeInvalidParameter, "文件ID列表不能为空")
	}
	db := database.GetDB()
	var updatedIDs []string

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.File{}).
			Where("id IN ?", imageIDs).
			Where("ai_tagging_status <> ?", common.AITaggingStatusIgnored).
			Update("ai_tagging_status", common.AITaggingStatusIgnored).Error; err != nil {
			return err
		}
		var images []models.File
		if err := tx.Where("id IN ?", imageIDs).Find(&images).Error; err != nil {
			return err
		}
		for _, file := range images {
			if file.AITaggingStatus == common.AITaggingStatusIgnored {
				updatedIDs = append(updatedIDs, file.ID)
			}
		}
		for _, id := range updatedIDs {
			logData := map[string]interface{}{}
			if reason != "" {
				logData["reason"] = reason
			}
			data, _ := json.Marshal(logData)

			logEntry := models.FileTaggingLog{
				FileID:     id,
				Status:     common.TaggingStatusIgnore,
				Action:     common.TaggingActionManual,
				Type:       "tagging.ignore",
				Data:       string(data),
				OperatorID: operatorID,
			}
			if err := tx.Create(&logEntry).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return 0, nil, errors.Wrap(err, errors.CodeDBUpdateFailed, "设置忽略状态失败")
	}

	// 触发队列状态更新通知（因为忽略状态变化会影响统计数据）
	if globalService := GetGlobalTaggingService(); globalService != nil && len(updatedIDs) > 0 {
		globalService.NotifyQueueStatsChange()
	}

	return len(updatedIDs), updatedIDs, nil
}

// UnignoreTagging 取消忽略（回到 none，以便重新打标）
func UnignoreTagging(imageIDs []string, operatorID uint) (int, []string, int, int, error) {
	if len(imageIDs) == 0 {
		return 0, nil, 0, 0, errors.New(errors.CodeInvalidParameter, "文件ID列表不能为空")
	}
	db := database.GetDB()
	var updatedIDs []string

	// 查询符合条件的被忽略的文件
	var ignoredImages []models.File
	err := db.Where("id IN ? AND ai_tagging_status = ?", imageIDs, common.AITaggingStatusIgnored).Find(&ignoredImages).Error
	if err != nil {
		return 0, nil, 0, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "查询被忽略的文件失败")
	}

	if len(ignoredImages) == 0 {
		return 0, nil, 0, 0, nil // 没有需要取消忽略的文件
	}

	globalService := GetGlobalTaggingService()
	if globalService == nil {
		if err := InitGlobalTaggingQueue(); err != nil {
			return 0, nil, 0, 0, errors.New(errors.CodeInternal, "队列服务不可用")
		}
		globalService = GetGlobalTaggingService()
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		// 将状态从 ignored 改为 none，同时重置重试次数为 0
		if err := tx.Model(&models.File{}).
			Where("id IN ?", imageIDs).
			Where("ai_tagging_status = ?", common.AITaggingStatusIgnored).
			Updates(map[string]interface{}{
				"ai_tagging_status": common.AITaggingStatusNone,
				"ai_tagging_tries":  0,
			}).Error; err != nil {
			return err
		}

		for _, file := range ignoredImages {
			updatedIDs = append(updatedIDs, file.ID)
			data, _ := json.Marshal(map[string]interface{}{})
			logEntry := models.FileTaggingLog{
				FileID:     file.ID,
				Status:     common.TaggingStatusUnignore,
				Action:     common.TaggingActionManual,
				Type:       "tagging.unignore",
				Data:       string(data),
				OperatorID: operatorID,
			}
			if err := tx.Create(&logEntry).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return 0, nil, 0, 0, errors.Wrap(err, errors.CodeDBUpdateFailed, "取消忽略失败")
	}

	// 添加到队列进行处理（同步获取统计）
	enqueued, skipped := globalService.BatchProcessFilesWithResult(ignoredImages)

	return len(updatedIDs), updatedIDs, enqueued, skipped, nil
}

// RetryFailedAll 批量重试失败与超过尝试上限的任务
// 规则：
// 1) 选择 ai_tagging_status in ('failed','none') 且 ai_tagging_tries >= 3 的文件作为优先重试目标
// 2) 可选 limit 限制单次重试数量，默认最多500
// 3) 将这些文件的 tries 重置为0，状态改为 none，并入队
func RetryFailedAll(limit int, operatorID uint) (int, int, int, error) {
	if limit <= 0 || limit > 1000 {
		limit = 500
	}

	db := database.GetDB()

	// 选择候选文件ID（包含 Skipped 和失败多次的任务，允许用户修正配置后重试）
	var imageIDs []string
	query := db.Table("file").
		Where("(ai_tagging_status IN ? AND ai_tagging_tries >= ?) OR ai_tagging_status = ?",
			[]string{common.AITaggingStatusFailed, common.AITaggingStatusNone}, 3, common.AITaggingStatusSkipped).
		Order("updated_at ASC").
		Limit(limit)
	if err := query.Pluck("id", &imageIDs).Error; err != nil {
		return 0, 0, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "查询需重试的任务失败")
	}
	if len(imageIDs) == 0 {
		return 0, 0, 0, nil
	}

	var images []models.File
	if err := db.Where("id IN ?", imageIDs).Find(&images).Error; err != nil {
		return 0, 0, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件详情失败")
	}

	// 兜底：全局队列
	svc := GetGlobalTaggingService()
	if svc == nil {
		if err := InitGlobalTaggingQueue(); err != nil {
			return 0, 0, 0, errors.New(errors.CodeInternal, "队列服务不可用")
		}
		svc = GetGlobalTaggingService()
	}

	// 重置 tries 与状态 + 记录日志
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.File{}).
			Where("id IN ?", imageIDs).
			Updates(map[string]interface{}{
				"ai_tagging_tries":  0,
				"ai_tagging_status": common.AITaggingStatusNone,
			}).Error; err != nil {
			return err
		}
		for _, id := range imageIDs {
			data, _ := json.Marshal(map[string]interface{}{})
			logEntry := models.FileTaggingLog{
				FileID:     id,
				Status:     common.TaggingStatusRetry,
				Action:     common.TaggingActionManual,
				Type:       "tagging.retry_failed_all",
				Data:       string(data),
				OperatorID: operatorID,
			}
			if err := tx.Create(&logEntry).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 0, 0, 0, errors.Wrap(err, errors.CodeDBUpdateFailed, "重置任务失败")
	}

	enq, skip := svc.BatchProcessFilesWithResult(images)
	return len(imageIDs), enq, skip, nil
}

// RetryTagging 指定文件批量重试打标：将状态重置为 none、清零重试次数并入队
// 返回：requested(请求总数)、enqueued(入队数)、skipped(跳过数)
func RetryTagging(imageIDs []string, operatorID uint) (int, int, int, error) {
	if len(imageIDs) == 0 {
		return 0, 0, 0, errors.New(errors.CodeInvalidParameter, "文件ID列表不能为空")
	}

	db := database.GetDB()

	var images []models.File
	if err := db.Where("id IN ?", imageIDs).Find(&images).Error; err != nil {
		return 0, 0, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件失败")
	}
	if len(images) == 0 {
		return 0, 0, 0, nil
	}

	svc := GetGlobalTaggingService()
	if svc == nil {
		if err := InitGlobalTaggingQueue(); err != nil {
			return 0, 0, 0, errors.New(errors.CodeInternal, "队列服务不可用")
		}
		svc = GetGlobalTaggingService()
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.File{}).
			Where("id IN ?", imageIDs).
			Updates(map[string]interface{}{
				"ai_tagging_status": common.AITaggingStatusNone,
				"ai_tagging_tries":  0,
			}).Error; err != nil {
			return err
		}
		for _, id := range imageIDs {
			data, _ := json.Marshal(map[string]interface{}{})
			entry := models.FileTaggingLog{
				FileID:     id,
				Status:     common.TaggingStatusRetry,
				Action:     common.TaggingActionManual,
				Type:       "tagging.retry",
				Data:       string(data),
				OperatorID: operatorID,
			}
			if err := tx.Create(&entry).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return 0, 0, 0, errors.Wrap(err, errors.CodeDBUpdateFailed, "重置打标状态失败")
	}

	// 入队（统计 enqueued/skipped）
	enqueued, skipped := svc.BatchProcessFilesWithResult(images)
	return len(imageIDs), enqueued, skipped, nil
}

// fileID: 可选,指定文件ID过滤
// page: 页码,从1开始
// limit: 每页数量
func GetTaggingLogs(fileID string, page, limit int) (map[string]interface{}, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	db := database.GetDB()
	query := db.Model(&models.FileTaggingLog{})

	// 如果指定了文件ID,添加过滤
	if fileID != "" {
		query = query.Where("file_id = ?", fileID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询日志总数失败")
	}

	var logs []models.FileTaggingLog
	offset := (page - 1) * limit
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error

	if err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询打标日志失败")
	}

	result := map[string]interface{}{
		"logs": logs,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}

	return result, nil
}
