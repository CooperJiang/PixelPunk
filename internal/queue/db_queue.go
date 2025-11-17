package queue

import (
	"errors"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"time"

	"gorm.io/gorm"
)

// DBQueue 基于数据库的简易可靠队列（支持SQLite/MySQL）
type DBQueue struct{ db *gorm.DB }

func NewDBQueue() *DBQueue {
	db := database.GetDB()
	return &DBQueue{db: db}
}

// EnqueueUnique: UPSERT 进入 ai_job，避免重复（改进版：检查File表心跳）
func (q *DBQueue) EnqueueUnique(fileID string, priority int) error {
	if q.db == nil {
		return errors.New("db not initialized")
	}

	// 先查询File表，检查是否正在被处理（通过心跳时间判断）
	var file struct {
		AITaggingStatus      string     `gorm:"column:ai_tagging_status"`
		AILastHeartbeatAt    *time.Time `gorm:"column:ai_last_heartbeat_at"`
		AIProcessingWorkerID string     `gorm:"column:ai_processing_worker_id"`
	}
	if err := q.db.Table("file").
		Where("id = ?", fileID).
		Take(&file).Error; err != nil {
		return nil
	}

	// 如果文件正在处理中（pending + 心跳时间在2分钟内），跳过入队
	if file.AITaggingStatus == "pending" &&
		file.AILastHeartbeatAt != nil &&
		time.Since(*file.AILastHeartbeatAt) < 2*time.Minute {
		return nil
	}

	// 为简化复用一张表（AIJob）。向量队列将独立使用 VectorJob 表；由调用方选择不同实现。
	job := models.AIJob{FileID: fileID, Status: "queued", Priority: priority}
	// SQLite/MySQL 兼容的幂等：先查再插/更
	var existing models.AIJob
	if err := q.db.Where("file_id = ?", fileID).Take(&existing).Error; err == nil {
		if existing.Status == "done" || existing.Status == "ignored" || existing.Status == "skipped" {
			return nil
		}
		// 如果是processing状态，检查租约是否过期
		if existing.Status == "processing" {
			if existing.LeaseUntil != nil && time.Now().Before(*existing.LeaseUntil) {
				return nil
			}
			return q.db.Model(&models.AIJob{}).Where("id = ?", existing.ID).Updates(map[string]interface{}{
				"status":      "queued",
				"lease_until": gorm.Expr("NULL"),
				"lease_by":    "",
			}).Error
		}
		// 非终态且非processing，保持queued
		if existing.Status != "queued" {
			return q.db.Model(&models.AIJob{}).Where("id = ?", existing.ID).Updates(map[string]interface{}{
				"status":      "queued",
				"lease_until": gorm.Expr("NULL"),
				"lease_by":    "",
			}).Error
		}
		return nil
	}
	if err := q.db.Create(&job).Error; err != nil {
		return nil
	}
	return nil
}

// Fetch 使用乐观锁方式抢占任务（无事务，高并发）
func (q *DBQueue) Fetch(lease time.Duration) (*TaggingTask, AckFunc, NackFunc, error) {
	if q.db == nil {
		return nil, nil, nil, errors.New("db not initialized")
	}

	now := time.Now()
	leaseUntil := now.Add(lease)

	var candidate models.AIJob
	if err := q.db.Where("(status = ? OR (status = ? AND lease_until < ?))", "queued", "processing", now).
		Order("priority DESC, created_at ASC").
		Take(&candidate).Error; err != nil {
		return nil, nil, nil, err
	}

	// 只有当任务状态未被其他 Worker 改变时才会更新成功
	result := q.db.Model(&models.AIJob{}).
		Where("id = ? AND (status = ? OR (status = ? AND lease_until < ?))",
			candidate.ID, "queued", "processing", now).
		Updates(map[string]interface{}{
			"status":      "processing",
			"lease_until": leaseUntil,
			"lease_by":    "tagger",
		})

	if result.Error != nil {
		return nil, nil, nil, result.Error
	}

	// 如果 RowsAffected == 0，说明被其他 Worker 抢走了，返回空
	if result.RowsAffected == 0 {
		return nil, nil, nil, gorm.ErrRecordNotFound
	}

	task := &TaggingTask{FileID: candidate.FileID}
	picked := candidate

	ack := func() error {
		return q.db.Model(&models.AIJob{}).Where("id = ?", picked.ID).Updates(map[string]interface{}{
			"status":      "done",
			"lease_until": gorm.Expr("NULL"),
			"lease_by":    "",
		}).Error
	}

	nack := func(delay time.Duration, toDLQ bool, lastError string) error {
		if toDLQ {
			return q.db.Model(&models.AIJob{}).Where("id = ?", picked.ID).Updates(map[string]interface{}{
				"status":      "failed",
				"last_error":  lastError,
				"lease_until": gorm.Expr("NULL"),
				"lease_by":    "",
			}).Error
		}
		// 指数退避由上层计算；这里直接将租约推后（使其在 delay 后可再取）
		return q.db.Model(&models.AIJob{}).Where("id = ?", picked.ID).Updates(map[string]interface{}{
			"status":      "queued",
			"tries":       gorm.Expr("tries + 1"),
			"last_error":  lastError,
			"lease_until": time.Now().Add(delay),
			"lease_by":    "",
		}).Error
	}

	return task, ack, nack, nil
}

func (q *DBQueue) Metrics() (*Metrics, error) {
	if q.db == nil {
		return nil, errors.New("db not initialized")
	}
	var queued, processing, delayed, dlq int64
	now := time.Now()
	if err := q.db.Model(&models.AIJob{}).Where("status = ?", "queued").Count(&queued).Error; err != nil {
		return nil, err
	}
	if err := q.db.Model(&models.AIJob{}).Where("status = ?", "processing").Count(&processing).Error; err != nil {
		return nil, err
	}
	if err := q.db.Model(&models.AIJob{}).Where("status = ? AND lease_until > ?", "queued", now).Count(&delayed).Error; err != nil {
		return nil, err
	}
	if err := q.db.Model(&models.AIJob{}).Where("status = ?", "failed").Count(&dlq).Error; err != nil {
		return nil, err
	}
	return &Metrics{QueueLength: int(queued), InFlight: int(processing), DelayedCount: int(delayed), DLQCount: int(dlq)}, nil
}

func (q *DBQueue) Close() error { return nil }
