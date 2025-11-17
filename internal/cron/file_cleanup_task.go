package cron

import (
	filesvc "pixelpunk/internal/services/file"
	"pixelpunk/internal/services/message"
	"pixelpunk/pkg/logger"
	"time"
)

/* ImageCleanupJob 文件清理任务 */
type ImageCleanupJob struct{}

/* NewImageCleanupJob 创建文件清理任务 */
func NewImageCleanupJob() *ImageCleanupJob {
	return &ImageCleanupJob{}
}

/* GetSchedule 获取任务调度配置 */
func (j *ImageCleanupJob) GetSchedule() string {
	return "0 * * * * *"
}

/* Execute 执行清理任务 */
func (j *ImageCleanupJob) Execute() error {
	startTime := time.Now()

	var totalExpired, deletedSuccess, deletedFailed, notifiedSuccess, notifiedFailed int

	if n, err := filesvc.CleanupPendingDeletionFiles(200); err != nil {
		logger.Error("❌ 清理待删除文件失败: %v", err)
	} else if n > 0 {
		logger.Info("🧹 清理待删除文件：%d", n)
	} else {
	}

	if success, failed, err := filesvc.CleanupExpiredFiles(); err != nil {
		logger.Error("❌ 清理过期文件失败: %v", err)
		return err
	} else {
		deletedSuccess = success
		deletedFailed = failed
		totalExpired = success + failed
	}

	if success, failed, err := message.GetMessageService().SendExpiryNotifications(); err != nil {
		logger.Error("❌ 发送过期提醒失败: %v", err)
		return err
	} else {
		notifiedSuccess = success
		notifiedFailed = failed
	}

	if err := j.cleanupGuestUploadLogs(); err != nil {
		logger.Error("❌ 清理游客上传日志失败: %v", err)
		return err
	}

	duration := time.Since(startTime)

	if totalExpired > 0 || notifiedSuccess > 0 {
		logger.Info("🎉 文件清理任务完成，耗时: %v, 统计: 过期=%d, 删除成功=%d, 删除失败=%d, 提醒成功=%d, 提醒失败=%d",
			duration, totalExpired, deletedSuccess, deletedFailed,
			notifiedSuccess, notifiedFailed)
	} else {
	}

	return nil
}

func (j *ImageCleanupJob) cleanupGuestUploadLogs() error {
	logService := filesvc.GetGuestUploadLogService()
	err := logService.DeleteExpiredLogs()
	if err != nil {
		logger.Error("❌ 清理游客上传日志失败: %v", err)
		return err
	}
	return nil
}
