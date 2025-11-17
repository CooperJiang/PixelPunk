package vector

import (
	"context"
	"fmt"
	"sync"
	"time"

	"pixelpunk/internal/controllers/websocket"
	metrics "pixelpunk/internal/metrics"
	"pixelpunk/internal/models"
	qqueue "pixelpunk/internal/queue"
	"pixelpunk/internal/services/setting"
	ws "pixelpunk/internal/websocket"
	"pixelpunk/pkg/cache"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"
	vector "pixelpunk/pkg/vector"

	"gorm.io/gorm"
)

/* VectorQueueService 向量队列服务（与打标对齐） */
type VectorQueueService struct {
	queue         qqueue.Queue
	paused        bool
	concurrent    int
	ctx           context.Context
	cancel        context.CancelFunc
	activeWorkers int
	reaperStop    chan struct{}
}

var globalVectorQueueService *VectorQueueService
var vectorHooksRegistered = false

// 防抖相关变量
var vectorScanTimer *time.Timer
var vectorScanMutex sync.Mutex

func GetGlobalVectorQueueService() *VectorQueueService { return globalVectorQueueService }

// RegisterVectorConfigHooks 注册向量配置变更钩子
// 采用统一处理策略：任何vector组配置变更都会触发初始化和扫描
func RegisterVectorConfigHooks() {
	if vectorHooksRegistered {
		return
	}
	vectorHooksRegistered = true

	handleVectorConfigChange := func() {
		vectorEnabled := setting.GetBoolDirectFromDB("vector", "vector_enabled", false)
		if !vectorEnabled {
			return
		}

		eng := vector.GetGlobalVectorEngine()
		if eng == nil {
			// 从数据库读取Qdrant配置并初始化引擎
			qdrantURL := setting.GetStringDirectFromDB("vector", "qdrant_url", "http://localhost:6333")
			qdrantTimeout := setting.GetIntDirectFromDB("vector", "qdrant_timeout", 30)

			if err := vector.InitQdrantVectorEngine(qdrantURL, qdrantTimeout); err != nil {
				logger.Error("[向量服务] 初始化向量引擎失败: %v", err)
				return
			}

			eng = vector.GetGlobalVectorEngine()
			if eng == nil {
				return
			}
			logger.Info("[向量服务] 向量引擎初始化成功")
		}

		if err := eng.ReloadConfig(); err != nil {
			logger.Error("[向量服务] 重新加载配置失败: %v", err)
			return
		}

		if globalVectorQueueService == nil {
			if err := InitGlobalVectorQueue(); err != nil {
				logger.Error("[向量服务] 初始化队列失败: %v", err)
				return
			}
			logger.Info("[向量服务] 向量队列初始化成功")
		}

		// 使用防抖机制延迟扫描
		if globalVectorQueueService != nil && eng.IsEnabled() && !globalVectorQueueService.IsPaused() {
			vectorScanMutex.Lock()
			if vectorScanTimer != nil {
				vectorScanTimer.Stop()
			}
			vectorScanTimer = time.AfterFunc(300*time.Millisecond, func() {
				if n, err := globalVectorQueueService.EnqueueAllPending(1000); err == nil && n > 0 {
				}
			})
			vectorScanMutex.Unlock()
		}
	}

	criticalKeys := []string{"vector_enabled", "vector_api_key", "vector_base_url", "vector_model", "qdrant_url"}
	for _, key := range criticalKeys {
		setting.RegisterSettingChangeHandler("vector", key, func(value string) {
			handleVectorConfigChange()
		})
	}

	setting.RegisterSettingChangeHandler("vector", "vector_concurrency", func(value string) {
		if globalVectorQueueService != nil {
			concurrency := setting.GetIntDirectFromDB("vector", "vector_concurrency", 3)
			if concurrency > 0 {
				globalVectorQueueService.UpdateConcurrency(concurrency)
			}
		}
	})

	setting.RegisterSettingChangeHandler("vector", "vector_auto_processing_enabled", func(value string) {
		enabled := setting.GetBoolDirectFromDB("vector", "vector_auto_processing_enabled", true)
		if globalVectorQueueService != nil {
			globalVectorQueueService.SetPaused(!enabled)
		} else if enabled {
			handleVectorConfigChange()
		}
	})

}

func InitGlobalVectorQueue() error {
	if globalVectorQueueService != nil {
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())

	concurrency := setting.GetIntDirectFromDB("vector", "vector_concurrency", 3)
	autoProcessingEnabled := setting.GetBoolDirectFromDB("vector", "vector_auto_processing_enabled", true)
	paused := !autoProcessingEnabled
	svc := &VectorQueueService{paused: paused, concurrent: concurrency, ctx: ctx, cancel: cancel, reaperStop: make(chan struct{})}

	if cache.IsRedisEnabled() {
		if rq := qqueue.NewRedisQueue(); rq != nil {
			svc.queue = rq.WithPrefix("vector")
			rq.StartReaper(1*time.Second, svc.reaperStop)
		}
	}
	if svc.queue == nil {
		svc.queue = qqueue.NewDBQueueVector()
	}

	for i := 0; i < svc.concurrent; i++ {
		go svc.worker(i + 1)
	}
	globalVectorQueueService = svc

	if !svc.IsPaused() {
		if n, err := svc.EnqueueAllPending(1000); err == nil && n > 0 {
		} else if err != nil {
			logger.Warn("向量队列启动自检入队失败: %v", err)
		}
	}

	return nil
}

func (s *VectorQueueService) worker(id int) {
	db := database.GetDB()
	for {
		if s.queue == nil {
			time.Sleep(1 * time.Second)
			continue
		}
		if s.paused || s.activeWorkers >= s.concurrent {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		s.activeWorkers++
		task, ack, nack, err := s.queue.Fetch(30 * time.Second)
		if err != nil || task == nil {
			s.activeWorkers--
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if s.paused {
			_ = nack(1*time.Second, false, "paused")
			s.activeWorkers--
			time.Sleep(100 * time.Millisecond)
			continue
		}

		s.processVectorTask(task, ack, nack, db)
		s.activeWorkers--
	}
}

func (s *VectorQueueService) processVectorTask(task *qqueue.TaggingTask, ack qqueue.AckFunc, nack qqueue.NackFunc, db *gorm.DB) {
	defer s.pushWS()

	var ai models.FileAIInfo
	if err := db.Where("file_id = ?", task.FileID).Take(&ai).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			_ = ack()
			return
		}
		_ = nack(1*time.Second, false, fmt.Sprintf("load ai info failed: %v", err))
		metrics.IncVectorNack()
		return
	}

	if ai.Description == "" {
		_ = ack()
		return
	}

	_ = db.Model(&models.FileVector{}).Where("file_id = ?", ai.FileID).Updates(map[string]interface{}{
		"status":        common.VectorStatusProcessing,
		"error_message": "",
	}).Error

	// 立即推送状态变更，让前端看到 processing 状态
	s.pushWS()

	engine := vector.GetGlobalVectorEngine()
	if engine == nil || !engine.IsEnabled() {
		_ = db.Model(&models.FileVector{}).Where("file_id = ?", ai.FileID).Updates(map[string]interface{}{
			"status":        common.VectorStatusFailed,
			"error_message": "vector engine not enabled",
		}).Error
		logger.Warn("向量引擎未启用，丢弃任务: file_id=%s", ai.FileID)
		_ = ack()
		return
	}

	errProc := engine.ProcessFile(ai.FileID, ai.Description)

	if errProc == nil {
		_ = db.Model(&models.FileVector{}).Where("file_id = ?", ai.FileID).Updates(map[string]interface{}{
			"status":        common.VectorStatusCompleted,
			"error_message": "",
			"retry_count":   0,
		}).Error
		_ = ack()
		metrics.IncVectorAck()

		go propagateVectorToDuplicates(ai.FileID)
	} else {
		var currentRetries int
		_ = db.Model(&models.FileVector{}).Where("file_id = ?", ai.FileID).Select("retry_count").Scan(&currentRetries).Error
		currentRetries++

		now := time.Now()
		_ = db.Model(&models.FileVector{}).Where("file_id = ?", ai.FileID).Updates(map[string]interface{}{
			"retry_count":   currentRetries,
			"last_retry_at": &now,
		}).Error

		maxRetries := 3
		if currentRetries >= maxRetries {
			_ = db.Model(&models.FileVector{}).Where("file_id = ?", ai.FileID).Updates(map[string]interface{}{
				"status":        common.VectorStatusFailed,
				"error_message": errProc.Error(),
				"retry_count":   maxRetries,
			}).Error
			logger.Error("向量生成失败（重试次数已达上限%d次）: file_id=%s, err=%v", maxRetries, ai.FileID, errProc)
			_ = ack()
			return
		}

		var delay time.Duration = 1 * time.Second
		if currentRetries >= 2 {
			delay = 5 * time.Second
		} else if currentRetries >= 1 {
			delay = 3 * time.Second
		}

		logger.Warn("向量生成重试: file_id=%s, 重试次数=%d/%d, 延迟=%v, err=%v",
			ai.FileID, currentRetries, maxRetries, delay, errProc)
		_ = nack(delay, false, errProc.Error())
		metrics.IncVectorNack()

		_ = db.Model(&models.FileVector{}).Where("file_id = ?", ai.FileID).Updates(map[string]interface{}{
			"status":        common.VectorStatusPending,
			"error_message": errProc.Error(),
		}).Error
	}
}

func propagateVectorToDuplicates(originalID string) {
	db := database.GetDB()
	if db == nil {
		return
	}
	var dups []models.File
	if err := db.Where("original_file_id = ?", originalID).Where("status <> ?", "pending_deletion").Find(&dups).Error; err != nil {
		return
	}
	if len(dups) == 0 {
		return
	}
	eng := vector.GetGlobalVectorEngine()
	if eng == nil || !eng.IsEnabled() {
		return
	}
	for _, dup := range dups {
		var iv models.FileVector
		if err := db.Where("file_id = ?", dup.ID).First(&iv).Error; err == nil && iv.Status == common.VectorStatusCompleted {
			continue
		}
		desc := ""
		var ai models.FileAIInfo
		if err := db.Where("file_id = ?", dup.ID).First(&ai).Error; err == nil {
			if ai.SearchContent != "" {
				desc = ai.SearchContent
			} else {
				desc = ai.Description
			}
		}
		if err := eng.CloneVectorFrom(originalID, dup.ID, desc); err != nil {
			logger.Warn("克隆向量到重复文件失败: orig=%s dup=%s err=%v", originalID, dup.ID, err)
		}
	}
}

func (s *VectorQueueService) EnqueueVector(fileID string) error {
	if s == nil || s.queue == nil {
		return fmt.Errorf("vector queue not ready")
	}
	return s.queue.EnqueueUnique(fileID, 0)
}

func (s *VectorQueueService) EnqueueAllPending(batch int) (int, error) {
	if s.paused {
		return 0, nil
	}
	if batch <= 0 {
		batch = 1000
	}
	db := database.GetDB()
	var ids []string
	if err := db.Table("file_vector").Where("status IN (?)", []string{"pending", "reset", "failed"}).Order("created_at asc").Limit(batch).Pluck("file_id", &ids).Error; err != nil {
		return 0, err
	}
	if len(ids) == 0 {
		return 0, nil
	}

	// 清理vector_job表中的旧记录，以便重新入队
	_ = db.Where("file_id IN ? AND status IN ?", ids, []string{"done", "failed"}).Delete(&models.VectorJob{}).Error

	enq := 0
	for _, id := range ids {
		if s.queue.EnqueueUnique(id, 0) == nil {
			enq++
		}
	}
	if enq > 0 {
		s.pushWS()
	}
	return enq, nil
}

func (s *VectorQueueService) SetPaused(v bool) {
	wasPaused := s.paused
	s.paused = v
	s.pushWS()

	// 如果从暂停变为恢复，自动扫描入队待处理任务
	if wasPaused && !v {
		go func() {
			if n, err := s.EnqueueAllPending(1000); err == nil && n > 0 {
			}
		}()
	}
}
func (s *VectorQueueService) UpdateConcurrency(n int) {
	if n <= 0 {
		return
	}
	if n > s.concurrent {
		add := n - s.concurrent
		for i := 0; i < add; i++ {
			go s.worker(s.concurrent + i + 1)
		}
	}
	s.concurrent = n
	s.pushWS()
}
func (s *VectorQueueService) pushWS() {
	go func() { stats := s.GetQueueStats(); websocket.BroadcastToAdmins(ws.MessageTypeVectorStats, stats) }()
}

/* GetQueueStats 供控制器/WS使用 */
func (s *VectorQueueService) GetQueueStats() map[string]interface{} {
	db := database.GetDB()
	var counts struct{ Pending, Processing, Completed, Failed, Reset, Stale, Total int64 }
	db.Model(&models.FileVector{}).
		Select("SUM(CASE WHEN status='pending' THEN 1 ELSE 0 END) AS pending, " +
			"SUM(CASE WHEN status='processing' THEN 1 ELSE 0 END) AS processing, " +
			"SUM(CASE WHEN status='completed' THEN 1 ELSE 0 END) AS completed, " +
			"SUM(CASE WHEN status='failed' THEN 1 ELSE 0 END) AS failed, " +
			"SUM(CASE WHEN status='reset' THEN 1 ELSE 0 END) AS reset, " +
			"SUM(CASE WHEN status='stale' THEN 1 ELSE 0 END) AS stale, " +
			"COUNT(*) AS total").Scan(&counts)

	ext := map[string]int{"queued": 0, "processing": 0, "delayed": 0, "dlq": 0}
	if s.queue != nil {
		if m, err := s.queue.Metrics(); err == nil && m != nil {
			ext["queued"], ext["processing"], ext["delayed"], ext["dlq"] = m.QueueLength, m.InFlight, m.DelayedCount, m.DLQCount
		}
	}
	runtime := map[string]interface{}{
		"active_workers":         s.activeWorkers,
		"configured_concurrency": s.concurrent,
		"paused":                 s.paused,
	}

	var totalWithDesc int64
	db.Table("file_ai_info").Where("description IS NOT NULL AND description != ''").Count(&totalWithDesc)
	withVectors := counts.Completed
	withoutVectors := totalWithDesc - withVectors

	divergence := map[string]interface{}{"qdrant_total": 0, "db_completed": withVectors, "delta": 0}
	if eng := vector.GetGlobalVectorEngine(); eng != nil {
		if st, err := eng.GetStorageStats(); err == nil && st != nil {
			delta := int64(st.TotalVectors) - withVectors
			if delta < 0 {
				delta = -delta
			}
			divergence = map[string]interface{}{
				"qdrant_total": st.TotalVectors,
				"db_completed": withVectors,
				"delta":        delta,
			}
		}
	}

	return map[string]interface{}{
		"queue_stats": map[string]interface{}{
			"pending_count":    counts.Pending,
			"processing_count": counts.Processing,
			"completed_count":  counts.Completed,
			"failed_count":     counts.Failed,
			"reset_count":      counts.Reset,
			"stale_count":      counts.Stale,
			"total_count":      counts.Total,
		},
		"queue_stats_ext": ext,
		"runtime":         runtime,
		"coverage": map[string]interface{}{
			"total_files_with_desc": totalWithDesc,
			"with_vectors":          withVectors,
			"without_vectors":       withoutVectors,
		},
		"divergence": divergence,
	}
}

/* IsPaused 返回暂停状态 */
func (s *VectorQueueService) IsPaused() bool { return s.paused }

/* ReconcileMissing 补齐缺失：对有描述但未完成的文件入队；支持dry-run与limit */
func (s *VectorQueueService) ReconcileMissing(limit int, dryRun bool) (int, int, error) {
	db := database.GetDB()
	if db == nil {
		return 0, 0, fmt.Errorf("数据库不可用")
	}
	if limit <= 0 {
		limit = 1000
	}

	var ids []string
	qry := db.Table("file_ai_info ai").
		Select("ai.file_id").
		Joins("LEFT JOIN file_vector iv ON iv.file_id = ai.file_id").
		Where("ai.description IS NOT NULL AND ai.description != ''").
		Where("iv.file_id IS NULL OR iv.status IN (?)", []string{"pending", "reset", "failed", "stale"}).
		Limit(limit)
	if err := qry.Pluck("ai.file_id", &ids).Error; err != nil {
		return 0, 0, err
	}
	if dryRun {
		return len(ids), 0, nil
	}

	for _, id := range ids {
		var cnt int64
		if err := db.Table("file_vector").Where("file_id = ?", id).Count(&cnt).Error; err == nil && cnt == 0 {
			var ai models.FileAIInfo
			_ = db.Where("file_id = ?", id).First(&ai).Error
			_ = db.Create(&models.FileVector{FileID: id, Description: ai.Description, Status: common.VectorStatusPending}).Error
		}
	}

	// 清理vector_job表中的旧记录
	if len(ids) > 0 {
		_ = db.Where("file_id IN ? AND status IN ?", ids, []string{"done", "failed"}).Delete(&models.VectorJob{}).Error
	}

	enq := 0
	for _, id := range ids {
		if s.EnqueueVector(id) == nil {
			enq++
		}
	}
	if enq > 0 {
		s.pushWS()
	}
	return len(ids), enq, nil
}

/* CleanOrphans 清理孤儿：DB 无对应文件或已软删除(DeletedAt)的记录，删除Qdrant向量；支持dry-run与limit */
func (s *VectorQueueService) CleanOrphans(limit int, dryRun bool) (int, int, error) {
	db := database.GetDB()
	if db == nil {
		return 0, 0, fmt.Errorf("数据库不可用")
	}
	if limit <= 0 {
		limit = 1000
	}

	eng := vector.GetGlobalVectorEngine()
	if eng == nil || !eng.IsEnabled() {
		return 0, 0, fmt.Errorf("向量引擎不可用")
	}

	qIDs, err := eng.GetAllFileIDs(limit)
	if err != nil {
		return 0, 0, err
	}
	if len(qIDs) == 0 {
		return 0, 0, nil
	}

	var existing []string
	if err := db.Table("file").Where("id IN (?)", qIDs).Pluck("id", &existing).Error; err != nil {
		return 0, 0, err
	}
	existSet := make(map[string]struct{}, len(existing))
	for _, id := range existing {
		existSet[id] = struct{}{}
	}

	var orphans []string
	for _, id := range qIDs {
		if _, ok := existSet[id]; !ok {
			orphans = append(orphans, id)
		}
	}
	if dryRun {
		return len(orphans), 0, nil
	}

	removed := 0
	for _, id := range orphans {
		if err := eng.DeleteVector(id); err == nil {
			removed++
		} else {
			logger.Warn("删除孤儿向量失败: %s, err=%v", id, err)
		}
	}
	if removed > 0 {
		s.pushWS()
	}
	return len(orphans), removed, nil
}

/* RebuildStale 重建过期：将 stale 记录入队；支持limit */
func (s *VectorQueueService) RebuildStale(limit int) (int, error) {
	db := database.GetDB()
	if db == nil {
		return 0, fmt.Errorf("数据库不可用")
	}
	if limit <= 0 {
		limit = 1000
	}
	var ids []string
	if err := db.Table("file_vector").Where("status = ?", common.VectorStatusStale).Order("updated_at asc").Limit(limit).Pluck("file_id", &ids).Error; err != nil {
		return 0, err
	}

	// 清理vector_job表中的旧记录
	if len(ids) > 0 {
		_ = db.Where("file_id IN ? AND status IN ?", ids, []string{"done", "failed"}).Delete(&models.VectorJob{}).Error
	}

	enq := 0
	for _, id := range ids {
		if s.EnqueueVector(id) == nil {
			enq++
		}
	}
	if enq > 0 {
		s.pushWS()
	}
	return enq, nil
}
