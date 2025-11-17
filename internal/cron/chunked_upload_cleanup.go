package cron

import (
	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"
	"time"

	"gorm.io/gorm"
)

/* ChunkedUploadCleanupJob 分片上传清理任务 */
type ChunkedUploadCleanupJob struct {
	db *gorm.DB
}

/* NewChunkedUploadCleanupJob 创建分片上传清理任务 */
func NewChunkedUploadCleanupJob() *ChunkedUploadCleanupJob {
	return &ChunkedUploadCleanupJob{
		db: database.GetDB(),
	}
}

/* Execute 执行清理任务 */
func (j *ChunkedUploadCleanupJob) Execute() error {
	expiredTime := time.Now().Add(-24 * time.Hour)

	var expiredSessions []models.UploadSession
	err := j.db.Where("updated_at < ? AND status NOT IN (?)", expiredTime, []string{"completed", "failed"}).
		Find(&expiredSessions).Error
	if err != nil {
		logger.Error("查询过期上传会话失败: %v", err)
		return err
	}

	if len(expiredSessions) == 0 {
		return nil
	}

	batchSize := 50
	for i := 0; i < len(expiredSessions); i += batchSize {
		end := i + batchSize
		if end > len(expiredSessions) {
			end = len(expiredSessions)
		}

		batch := expiredSessions[i:end]
		if err := j.cleanupSessionBatch(batch); err != nil {
			logger.Error("清理会话批次失败: %v", err)
			continue
		}
	}

	return nil
}

func (j *ChunkedUploadCleanupJob) cleanupSessionBatch(sessions []models.UploadSession) error {
	return j.db.Transaction(func(tx *gorm.DB) error {
		for _, session := range sessions {
			if err := tx.Where("session_id = ?", session.SessionID).Delete(&models.UploadChunk{}).Error; err != nil {
				logger.Error("删除会话 %s 的分片记录失败: %v", session.SessionID, err)
				return err
			}

			session.Status = "cleaned"
			session.UpdatedAt = common.JSONTime(time.Now())
			if err := tx.Save(&session).Error; err != nil {
				logger.Error("更新会话 %s 状态失败: %v", session.SessionID, err)
				return err
			}

		}

		return nil
	})
}

/* GetName 获取任务名称 */
func (j *ChunkedUploadCleanupJob) GetName() string {
	return "chunked_upload_cleanup"
}

/* GetSchedule 获取调度配置（每小时执行一次） */
func (j *ChunkedUploadCleanupJob) GetSchedule() string {
	return "0 0 * * * *" // 每小时的0分执行
}
