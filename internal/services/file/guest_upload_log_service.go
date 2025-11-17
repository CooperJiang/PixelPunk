package file

import (
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"
	"time"
)

/* GuestUploadLogService 游客上传日志服务 */
type GuestUploadLogService struct{}

/* NewGuestUploadLogService 创建游客上传日志服务实例 */
func NewGuestUploadLogService() *GuestUploadLogService {
	return &GuestUploadLogService{}
}

/* RecordGuestUpload 记录游客上传日志 */
func (s *GuestUploadLogService) RecordGuestUpload(log *models.GuestUploadLog) error {
	if err := database.DB.Create(log).Error; err != nil {
		logger.Error("记录游客上传日志失败: %v", err)
		return err
	}

	return nil
}

/* GetGuestUploadLogs 获取游客上传日志列表 */
func (s *GuestUploadLogService) GetGuestUploadLogs(fingerprint string, page, pageSize int) ([]*models.GuestUploadLog, int64, error) {
	var logs []*models.GuestUploadLog
	var total int64

	offset := (page - 1) * pageSize

	query := database.DB.Model(&models.GuestUploadLog{})
	if fingerprint != "" {
		query = query.Where("fingerprint = ?", fingerprint)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Order("created_at DESC")
	if fingerprint != "" {
		query = query.Where("fingerprint = ?", fingerprint)
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

/* GetGuestUploadStatsByFingerprint 根据指纹获取游客上传统计 */
func (s *GuestUploadLogService) GetGuestUploadStatsByFingerprint(fingerprint string) (map[string]interface{}, error) {
	var stats struct {
		TotalUploads int64      `json:"total_uploads"`
		TotalSize    int64      `json:"total_size"`
		LastUploadAt *time.Time `json:"last_upload_at"`
	}

	err := database.DB.Model(&models.GuestUploadLog{}).
		Select("COUNT(*) as total_uploads, SUM(file_size) as total_size, MAX(upload_time) as last_upload_at").
		Where("fingerprint = ?", fingerprint).
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	var durationStats []struct {
		StorageDuration string `json:"storage_duration"`
		Count           int64  `json:"count"`
	}

	if err := database.DB.Model(&models.GuestUploadLog{}).
		Select("storage_duration, COUNT(*) as count").
		Where("fingerprint = ?", fingerprint).
		Group("storage_duration").
		Find(&durationStats).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_uploads":  stats.TotalUploads,
		"total_size":     stats.TotalSize,
		"last_upload_at": stats.LastUploadAt,
		"duration_stats": durationStats,
	}, nil
}

/* DeleteExpiredLogs 清理过期的上传日志（保留最近30天） */
func (s *GuestUploadLogService) DeleteExpiredLogs() error {
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	result := database.DB.Where("created_at < ?", thirtyDaysAgo).Delete(&models.GuestUploadLog{})
	if result.Error != nil {
		logger.Error("清理过期上传日志失败: %v", result.Error)
		return result.Error
	}

	if result.RowsAffected > 0 {
	}

	return nil
}

var guestUploadLogService *GuestUploadLogService

/* GetGuestUploadLogService 获取游客上传日志服务实例（延迟初始化） */
func GetGuestUploadLogService() *GuestUploadLogService {
	if guestUploadLogService == nil {
		guestUploadLogService = NewGuestUploadLogService()
	}
	return guestUploadLogService
}
