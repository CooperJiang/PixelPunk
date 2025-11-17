package vector

import (
	"encoding/json"
	"fmt"
	websocket "pixelpunk/internal/controllers/websocket"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/setting"
	ws "pixelpunk/internal/websocket"
	"pixelpunk/pkg/cache"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/storage"
	"pixelpunk/pkg/vector"
	"time"

	"gorm.io/gorm"
)

/* GetVectorList 获取向量列表 */
func GetVectorList(page, pageSize int, status, model, keyword string) ([]models.FileVector, int64, error) {
	db := database.GetDB()
	if db == nil {
		return nil, 0, fmt.Errorf("数据库连接不可用")
	}

	query := db.Model(&models.FileVector{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if model != "" {
		query = query.Where("model = ?", model)
	}

	if keyword != "" {
		query = query.Where("file_id LIKE ? OR description LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询向量总数失败: %v", err)
	}

	var vectors []models.FileVector
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&vectors).Error; err != nil {
		return nil, 0, fmt.Errorf("查询向量列表失败: %v", err)
	}

	return vectors, total, nil
}

/* VectorWithFileInfo 带文件信息的向量数据结构 */
type VectorWithFileInfo struct {
	models.FileVector
	FileName     string `json:"file_name"`
	FileURL      string `json:"file_url"`
	ThumbnailURL string `json:"thumbnail_url"`
}

/* GetVectorWithFileInfo 获取带文件信息的向量列表 */
func GetVectorWithFileInfo(page, pageSize int, status, model, keyword string) ([]VectorWithFileInfo, int64, error) {
	db := database.GetDB()
	if db == nil {
		return nil, 0, fmt.Errorf("数据库连接不可用")
	}

	query := db.Model(&models.FileVector{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if model != "" {
		query = query.Where("model = ?", model)
	}

	if keyword != "" {
		query = query.Where("file_id LIKE ? OR description LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询向量总数失败: %v", err)
	}

	var vectors []models.FileVector
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&vectors).Error; err != nil {
		return nil, 0, fmt.Errorf("查询向量列表失败: %v", err)
	}

	var results []VectorWithFileInfo
	for _, vector := range vectors {
		result := VectorWithFileInfo{
			FileVector: vector,
		}

		var file models.File
		if err := db.Where("id = ?", vector.FileID).First(&file).Error; err == nil {
			result.FileName = file.OriginalName
			fullURL, fullThumbURL, _ := storage.GetFullURLs(file)
			result.FileURL = fullURL
			result.ThumbnailURL = fullThumbURL
		}

		results = append(results, result)
	}

	return results, total, nil
}

/* GetVectorStats 获取向量统计信息 */
func GetVectorStats() (map[string]interface{}, error) {
	cacheKey := "dashboard:vector_stats"

	if cached, err := cache.Get(cacheKey); err == nil && cached != "" {
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return result, nil
		}
	}

	svc := GetGlobalVectorQueueService()
	if svc == nil {
		return nil, fmt.Errorf("向量队列不可用")
	}
	db := database.GetDB()
	var stats struct {
		Pending   int64
		Running   int64
		Completed int64
		Failed    int64
		Reset     int64
		Total     int64
	}
	db.Model(&models.FileVector{}).
		Select("SUM(CASE WHEN status='pending' THEN 1 ELSE 0 END) AS pending, " +
			"SUM(CASE WHEN status='running' THEN 1 ELSE 0 END) AS running, " +
			"SUM(CASE WHEN status='completed' THEN 1 ELSE 0 END) AS completed, " +
			"SUM(CASE WHEN status='failed' THEN 1 ELSE 0 END) AS failed, " +
			"SUM(CASE WHEN status='reset' THEN 1 ELSE 0 END) AS reset, " +
			"COUNT(*) AS total").Scan(&stats)
	result := map[string]interface{}{
		"pending_count":   stats.Pending,
		"running_count":   stats.Running,
		"completed_count": stats.Completed,
		"failed_count":    stats.Failed,
		"reset_count":     stats.Reset,
		"active_workers":  0,
		"max_workers":     0,
		"total_tasks":     stats.Total,
	}

	if data, err := json.Marshal(result); err == nil {
		cache.Set(cacheKey, string(data), 5*time.Second)
	}

	return result, nil
}

/* BatchVectorAction 批量向量操作 */
func BatchVectorAction(fileIDs []string, action string) error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("数据库连接不可用")
	}

	switch action {
	case "reset":
		return batchResetVectors(db, fileIDs)
	case "retry":
		return batchRetryVectors(db, fileIDs)
	case "delete":
		return batchDeleteVectors(db, fileIDs)
	default:
		return fmt.Errorf("不支持的操作类型: %s", action)
	}
}

func batchResetVectors(db *gorm.DB, fileIDs []string) error {
	result := db.Model(&models.FileVector{}).
		Where("file_id IN ?", fileIDs).
		Updates(map[string]interface{}{
			"status":        common.VectorStatusReset,
			"error_message": "",
			"retry_count":   0,
		})

	if result.Error != nil {
		return fmt.Errorf("批量重置向量失败: %v", result.Error)
	}

	for _, fileID := range fileIDs {
		var vector models.FileVector
		model := ""
		if db.Where("file_id = ?", fileID).First(&vector).Error == nil {
			model = vector.Model
		}

		logVectorProcessing(fileID, models.VectorLogActionReset, "vector.reset",
			map[string]interface{}{}, model, 0, "", "")
	}

	notifyVectorStatsChange()

	return nil
}

func batchRetryVectors(db *gorm.DB, fileIDs []string) error {
	result := db.Model(&models.FileVector{}).
		Where("file_id IN ? AND status = ?", fileIDs, common.VectorStatusFailed).
		Updates(map[string]interface{}{
			"status":        common.VectorStatusPending,
			"error_message": "",
		})

	if result.Error != nil {
		return fmt.Errorf("批量重试向量失败: %v", result.Error)
	}

	var vectors []models.FileVector
	if err := db.Where("file_id IN ? AND status = ?", fileIDs, common.VectorStatusPending).
		Find(&vectors).Error; err != nil {
		logger.Error("获取待重试向量记录失败: %v", err)
		return err
	}

	pushedCount := 0
	for _, vectorRecord := range vectors {
		var aiInfo models.FileAIInfo
		if err := db.Where("file_id = ?", vectorRecord.FileID).First(&aiInfo).Error; err != nil {
			logger.Warn("获取文件AI描述失败，跳过重试: %s, 错误: %v", vectorRecord.FileID, err)
			continue
		}

		if aiInfo.Description == "" {
			logger.Warn("文件描述为空，跳过重试: %s", vectorRecord.FileID)
			continue
		}

		if svc := GetGlobalVectorQueueService(); svc != nil && svc.EnqueueVector(vectorRecord.FileID) == nil {
			pushedCount++
		} else {
			logger.Warn("重试任务推送失败: %s", vectorRecord.FileID)
		}

		logVectorProcessing(vectorRecord.FileID, models.VectorLogActionRetry, "vector.retry",
			map[string]interface{}{"retry_count": vectorRecord.RetryCount + 1}, vectorRecord.Model, 0, "", "")
	}

	notifyVectorStatsChange()

	return nil
}

func batchDeleteVectors(db *gorm.DB, fileIDs []string) error {
	var vectorModels []struct {
		FileID string
		Model  string
	}
	db.Model(&models.FileVector{}).
		Select("file_id, model").
		Where("file_id IN ?", fileIDs).
		Scan(&vectorModels)

	result := db.Where("file_id IN ?", fileIDs).Delete(&models.FileVector{})

	if result.Error != nil {
		return fmt.Errorf("批量删除向量失败: %v", result.Error)
	}

	for _, vm := range vectorModels {
		logVectorProcessing(vm.FileID, "delete", "vector.delete",
			map[string]interface{}{}, vm.Model, 0, "", "")
	}

	notifyVectorStatsChange()

	return nil
}

/* RetryVector 重试单个向量 */
func RetryVector(fileID string) error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("数据库连接不可用")
	}

	var vector models.FileVector
	if err := db.Where("file_id = ?", fileID).First(&vector).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("向量记录不存在: %s", fileID)
		}
		return fmt.Errorf("查询向量记录失败: %v", err)
	}

	if !vector.IsFailed() {
		return fmt.Errorf("只能重试失败状态的向量")
	}

	if err := db.Model(&vector).Updates(map[string]interface{}{
		"status":        common.VectorStatusPending,
		"error_message": "",
	}).Error; err != nil {
		return fmt.Errorf("更新向量状态失败: %v", err)
	}

	logVectorProcessing(fileID, models.VectorLogActionRetry, "vector.retry",
		map[string]interface{}{}, vector.Model, 0, "", "")

	notifyVectorStatsChange()

	return nil
}

/* GetVectorLogs 获取向量处理日志 */
func GetVectorLogs(page, pageSize int, fileID, action, model string) ([]models.VectorProcessingLog, int64, error) {
	db := database.GetDB()
	if db == nil {
		return nil, 0, fmt.Errorf("数据库连接不可用")
	}

	query := db.Model(&models.VectorProcessingLog{})

	if fileID != "" {
		query = query.Where("file_id = ?", fileID)
	}

	if action != "" {
		query = query.Where("action = ?", action)
	}

	if model != "" {
		query = query.Where("model = ?", model)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询日志总数失败: %v", err)
	}

	var logs []models.VectorProcessingLog
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("查询日志列表失败: %v", err)
	}

	return logs, total, nil
}

/* RegenerateAllVectors 重新生成所有向量 */
func RegenerateAllVectors(force bool) error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("数据库连接不可用")
	}

	var aiInfos []models.FileAIInfo
	query := db.Model(&models.FileAIInfo{}).Where("description != '' AND description IS NOT NULL")

	if !force {
		query = query.Where("NOT EXISTS (SELECT 1 FROM file_vector iv WHERE iv.file_id = file_ai_info.file_id AND iv.status = ?)",
			common.VectorStatusCompleted)
	}

	if err := query.Find(&aiInfos).Error; err != nil {
		return fmt.Errorf("查询AI信息失败: %v", err)
	}

	if len(aiInfos) == 0 {
		return nil
	}

	successCount := 0
	for _, aiInfo := range aiInfos {
		var existingVector models.FileVector
		err := db.Where("file_id = ?", aiInfo.FileID).First(&existingVector).Error

		if err == gorm.ErrRecordNotFound {
			newVector := &models.FileVector{
				FileID:      aiInfo.FileID,
				Description: aiInfo.Description,
				Status:      common.VectorStatusPending,
				Model:       "text-embedding-3-small",
				Dimension:   1536,
			}

			if err := db.Create(newVector).Error; err != nil {
				logger.Error("创建向量记录失败: %s, 错误: %v", aiInfo.FileID, err)
				continue
			}
			successCount++
		} else if err != nil {
			logger.Error("查询向量记录失败: %s, 错误: %v", aiInfo.FileID, err)
			continue
		} else {
			updateData := map[string]interface{}{
				"status":        common.VectorStatusReset,
				"description":   aiInfo.Description,
				"error_message": "",
				"retry_count":   0,
			}

			if err := db.Model(&existingVector).Updates(updateData).Error; err != nil {
				logger.Error("更新向量记录失败: %s, 错误: %v", aiInfo.FileID, err)
				continue
			}
			successCount++
		}
	}

	if successCount > 0 {
		pushedCount := 0
		for _, aiInfo := range aiInfos {
			if svc := GetGlobalVectorQueueService(); svc != nil && svc.EnqueueVector(aiInfo.FileID) == nil {
				pushedCount++
			}
		}

		notifyVectorStatsChange()
	}

	return nil
}

/* RetryAllFailedVectors 重试所有失败的向量任务 */
func RetryAllFailedVectors() error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("数据库连接不可用")
	}

	var failedVectors []models.FileVector
	if err := db.Where("status = ?", common.VectorStatusFailed).Find(&failedVectors).Error; err != nil {
		return fmt.Errorf("查询失败向量记录失败: %v", err)
	}

	if len(failedVectors) == 0 {
		return nil
	}

	successCount := 0
	taskID := fmt.Sprintf("retry-all-%d", time.Now().Unix())

	for _, vector := range failedVectors {
		updateData := map[string]interface{}{
			"status":        common.VectorStatusPending,
			"error_message": nil,
			"retry_count":   vector.RetryCount + 1,
		}

		if err := db.Model(&vector).Updates(updateData).Error; err != nil {
			logger.Error("更新失败向量记录失败: %s, 错误: %v", vector.FileID, err)
			continue
		}

		data, _ := json.Marshal(map[string]interface{}{
			"retry_count": vector.RetryCount + 1,
		})
		log := &models.VectorProcessingLog{
			FileID: vector.FileID,
			Action: "retry",
			Type:   "vector.retry_all",
			Data:   string(data),
			Model:  vector.Model,
			TaskID: taskID,
		}

		if err := db.Create(log).Error; err != nil {
			logger.Error("记录重试日志失败: %s, 错误: %v", vector.FileID, err)
		}

		successCount++
	}

	if successCount > 0 {
		notifyVectorStatsChange()
	}

	return nil
}

/* GetAvailableModels 获取可用的向量模型列表 */
func GetAvailableModels() []string {
	db := database.GetDB()
	if db == nil {
		return []string{"text-embedding-3-small", "text-embedding-3-large"}
	}

	var usedModels []string
	db.Model(&models.FileVector{}).
		Distinct("model").
		Where("model != ''").
		Pluck("model", &usedModels)

	defaultModels := []string{"text-embedding-3-small", "text-embedding-3-large"}
	for _, defaultModel := range defaultModels {
		found := false
		for _, model := range usedModels {
			if model == defaultModel {
				found = true
				break
			}
		}
		if !found {
			usedModels = append(usedModels, defaultModel)
		}
	}

	return usedModels
}

func logVectorProcessing(fileID, action, logType string, data map[string]interface{}, model string, duration int, errorCode, taskID string) {
	db := database.GetDB()
	if db == nil {
		logger.Error("数据库连接不可用，无法记录向量处理日志")
		return
	}

	dataJSON, _ := json.Marshal(data)
	log := &models.VectorProcessingLog{
		FileID:    fileID,
		Action:    action,
		Type:      logType,
		Data:      string(dataJSON),
		Model:     model,
		Duration:  duration,
		ErrorCode: errorCode,
		TaskID:    taskID,
	}

	if err := db.Create(log).Error; err != nil {
		logger.Error("创建向量处理日志失败: %v", err)
	}
}

func notifyVectorStatsChange() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("向量WebSocket通知发生panic: %v", r)
			}
		}()

		if stats, err := GetVectorStats(); err == nil {
			websocket.BroadcastToAdmins(ws.MessageTypeVectorStats, stats)
		} else {
			logger.Error("获取向量统计信息失败，无法发送WebSocket通知: %v", err)
		}
	}()
}

/* RecoverStuckRunningTasks 恢复卡住的running状态任务 */
func RecoverStuckRunningTasks() (int64, error) {
	db := database.GetDB()
	if db == nil {
		return 0, fmt.Errorf("数据库连接不可用")
	}

	result := db.Model(&models.FileVector{}).
		Where("status = ?", "running").
		Updates(map[string]interface{}{
			"status":        common.VectorStatusPending,
			"error_message": "检测到卡住的任务，已手动恢复",
		})

	if result.Error != nil {
		logger.Error("恢复卡住的向量任务失败: %v", result.Error)
		return 0, result.Error
	}

	if result.RowsAffected > 0 {

		logVectorProcessing("", "recover_stuck", "vector.recover_stuck",
			map[string]interface{}{"count": result.RowsAffected},
			"", 0, "MANUAL_STUCK_RECOVERY", "")

		notifyVectorStatsChange()
	}

	return result.RowsAffected, nil
}

/* QdrantRealStatsResponse Qdrant 实际统计信息响应 */
type QdrantRealStatsResponse struct {
	QdrantVectorCount   int64   `json:"qdrant_vector_count"`   // Qdrant 中实际的向量数量
	QdrantIndexedCount  int64   `json:"qdrant_indexed_count"`  // Qdrant 中已索引的向量数量
	MySQLTotalCount     int64   `json:"mysql_total_count"`     // MySQL 中的总记录数
	MySQLCompletedCount int64   `json:"mysql_completed_count"` // MySQL 中已完成的记录数
	SyncRatio           float64 `json:"sync_ratio"`            // 同步比例 (Qdrant/MySQL completed)
	IsHealthy           bool    `json:"is_healthy"`            // 是否健康（同步率 > 95%）
	CollectionStatus    string  `json:"collection_status"`     // 集合状态
	LastChecked         string  `json:"last_checked"`          // 最后检查时间
}

/* GetQdrantRealStats 获取 Qdrant 实际存储统计信息 */
func GetQdrantRealStats() (*QdrantRealStatsResponse, error) {
	db := database.GetDB()
	if db == nil {
		return nil, fmt.Errorf("数据库连接不可用")
	}

	var mysqlTotal int64
	if err := db.Model(&models.FileVector{}).Count(&mysqlTotal).Error; err != nil {
		return nil, fmt.Errorf("查询 MySQL 总记录数失败: %v", err)
	}

	var mysqlCompleted int64
	if err := db.Model(&models.FileVector{}).
		Where("status = ?", "completed").
		Count(&mysqlCompleted).Error; err != nil {
		return nil, fmt.Errorf("查询 MySQL 已完成记录数失败: %v", err)
	}

	qdrantURL := setting.GetStringDirectFromDB("vector", "qdrant_url", "")
	if qdrantURL == "" {
		return &QdrantRealStatsResponse{
			QdrantVectorCount:   0,
			QdrantIndexedCount:  0,
			MySQLTotalCount:     mysqlTotal,
			MySQLCompletedCount: mysqlCompleted,
			SyncRatio:           0,
			IsHealthy:           false,
			CollectionStatus:    "engine_not_initialized",
			LastChecked:         time.Now().Format("2006-01-02 15:04:05"),
		}, nil
	}

	qdrantTimeout := setting.GetIntDirectFromDB("vector", "qdrant_timeout", 30)
	if qdrantTimeout <= 0 {
		qdrantTimeout = 30
	}

	client := vector.NewQdrantClient(qdrantURL, qdrantTimeout)

	if err := client.HealthCheck(); err != nil {
		return &QdrantRealStatsResponse{
			QdrantVectorCount:   0,
			QdrantIndexedCount:  0,
			MySQLTotalCount:     mysqlTotal,
			MySQLCompletedCount: mysqlCompleted,
			SyncRatio:           0,
			IsHealthy:           false,
			CollectionStatus:    "qdrant_connection_failed",
			LastChecked:         time.Now().Format("2006-01-02 15:04:05"),
		}, nil
	}

	storageStats, err := client.GetStorageStats()
	if err != nil {
		return &QdrantRealStatsResponse{
			QdrantVectorCount:   0,
			QdrantIndexedCount:  0,
			MySQLTotalCount:     mysqlTotal,
			MySQLCompletedCount: mysqlCompleted,
			SyncRatio:           0,
			IsHealthy:           false,
			CollectionStatus:    "qdrant_connection_failed",
			LastChecked:         time.Now().Format("2006-01-02 15:04:05"),
		}, nil
	}

	var syncRatio float64
	if mysqlCompleted > 0 {
		actualSynced := storageStats.TotalVectors
		if actualSynced > mysqlCompleted {
			actualSynced = mysqlCompleted
		}
		syncRatio = float64(actualSynced) / float64(mysqlCompleted) * 100
	} else {
		syncRatio = 100.0
	}

	delta := storageStats.TotalVectors - mysqlCompleted
	if delta < 0 {
		delta = -delta
	}

	tolerance := int64(10)
	if mysqlCompleted > 200 {
		tolerance = mysqlCompleted / 20
	}
	isHealthy := syncRatio > 95.0 && delta <= tolerance

	collectionStatus := "unknown"
	if storageStats.TotalVectors == 0 {
		if mysqlCompleted > 0 {
			collectionStatus = "empty_but_should_have_data"
		} else {
			collectionStatus = "empty_normal"
		}
	} else {
		if isHealthy {
			collectionStatus = "healthy"
		} else if syncRatio >= 80 {
			collectionStatus = "partially_synced"
		} else {
			collectionStatus = "needs_resync"
		}
	}

	return &QdrantRealStatsResponse{
		QdrantVectorCount:   storageStats.TotalVectors,
		QdrantIndexedCount:  storageStats.CompletedCount,
		MySQLTotalCount:     mysqlTotal,
		MySQLCompletedCount: mysqlCompleted,
		SyncRatio:           syncRatio,
		IsHealthy:           isHealthy,
		CollectionStatus:    collectionStatus,
		LastChecked:         time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}
