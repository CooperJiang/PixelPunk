package queue

import (
	"errors"
	"time"

	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"

	"gorm.io/gorm"
)

// DBQueueVector 使用 vector_job 表的队列实现
type DBQueueVector struct{ db *gorm.DB }

func NewDBQueueVector() *DBQueueVector { return &DBQueueVector{db: database.GetDB()} }

func (q *DBQueueVector) EnqueueUnique(fileID string, priority int) error {
	if q.db == nil {
		return errors.New("db not initialized")
	}
	var existing models.VectorJob
	if err := q.db.Where("file_id = ?", fileID).Take(&existing).Error; err == nil {
		if existing.Status == "done" {
			return nil
		}
		if existing.Status != "queued" {
			return q.db.Model(&models.VectorJob{}).Where("id = ?", existing.ID).Updates(map[string]interface{}{
				"status":      "queued",
				"lease_until": gorm.Expr("NULL"),
				"lease_by":    "",
			}).Error
		}
		return nil
	}
	job := models.VectorJob{FileID: fileID, Status: "queued", Priority: priority}
	return q.db.Create(&job).Error
}

func (q *DBQueueVector) Fetch(lease time.Duration) (*TaggingTask, AckFunc, NackFunc, error) {
	if q.db == nil {
		return nil, nil, nil, errors.New("db not initialized")
	}

	var candidate models.VectorJob
	now := time.Now()

	if err := q.db.Where("(status = ? OR (status = ? AND lease_until < ?))", "queued", "processing", now).
		Order("priority DESC, created_at ASC").Take(&candidate).Error; err != nil {
		return nil, nil, nil, err
	}

	leaseUntil := now.Add(lease)
	result := q.db.Model(&models.VectorJob{}).
		Where("id = ? AND (status = ? OR (status = ? AND lease_until < ?))", candidate.ID, "queued", "processing", now).
		Updates(map[string]interface{}{
			"status":      "processing",
			"lease_until": leaseUntil,
			"lease_by":    "vector",
		})

	if result.Error != nil {
		return nil, nil, nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil, nil, errors.New("task taken by another worker")
	}

	picked := candidate

	task := &TaggingTask{FileID: picked.FileID}
	ack := func() error {
		return q.db.Model(&models.VectorJob{}).Where("id = ?", picked.ID).Updates(map[string]interface{}{
			"status": "done", "lease_until": gorm.Expr("NULL"), "lease_by": "",
		}).Error
	}
	nack := func(delay time.Duration, toDLQ bool, lastError string) error {
		if toDLQ {
			return q.db.Model(&models.VectorJob{}).Where("id = ?", picked.ID).Updates(map[string]interface{}{
				"status": "failed", "last_error": lastError, "lease_until": gorm.Expr("NULL"), "lease_by": "",
			}).Error
		}
		return q.db.Model(&models.VectorJob{}).Where("id = ?", picked.ID).Updates(map[string]interface{}{
			"status": "queued", "attempt": gorm.Expr("attempt + 1"), "last_error": lastError, "lease_until": time.Now().Add(delay), "lease_by": "",
		}).Error
	}
	return task, ack, nack, nil
}

func (q *DBQueueVector) Metrics() (*Metrics, error) {
	if q.db == nil {
		return nil, errors.New("db not initialized")
	}
	var queued, processing, delayed, dlq int64
	now := time.Now()
	if err := q.db.Model(&models.VectorJob{}).Where("status = ?", "queued").Count(&queued).Error; err != nil {
		return nil, err
	}
	if err := q.db.Model(&models.VectorJob{}).Where("status = ?", "processing").Count(&processing).Error; err != nil {
		return nil, err
	}
	if err := q.db.Model(&models.VectorJob{}).Where("status = ? AND lease_until > ?", "queued", now).Count(&delayed).Error; err != nil {
		return nil, err
	}
	if err := q.db.Model(&models.VectorJob{}).Where("status = ?", "failed").Count(&dlq).Error; err != nil {
		return nil, err
	}
	return &Metrics{QueueLength: int(queued), InFlight: int(processing), DelayedCount: int(delayed), DLQCount: int(dlq)}, nil
}

func (q *DBQueueVector) Close() error { return nil }
