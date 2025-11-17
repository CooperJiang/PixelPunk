package file

import (
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/setting"
	"pixelpunk/pkg/config"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"

	"gorm.io/gorm"
)

/* GuestUploadLimitService 游客上传限制服务 */
type GuestUploadLimitService struct{}

/* NewGuestUploadLimitService 创建游客上传限制服务实例 */
func NewGuestUploadLimitService() *GuestUploadLimitService {
	return &GuestUploadLimitService{}
}

/* CheckUploadLimit 检查游客上传限制 */
func (s *GuestUploadLimitService) CheckUploadLimit(fingerprint string) (bool, int, error) {

	if s == nil {
		logger.Error("CheckUploadLimit: service is nil")
		return false, 0, errors.New(errors.CodeInternal, "游客上传限制服务未初始化")
	}

	if database.DB == nil {
		logger.Error("CheckUploadLimit: database.DB is nil")
		return false, 0, errors.New(errors.CodeInternal, "数据库连接未初始化")
	}

	guestSettings, err := setting.GetSettingsByGroupAsMap("guest")
	if err != nil {
		logger.Error("CheckUploadLimit: 获取游客设置失败, err=%v", err)
		return false, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "获取游客设置失败")
	}

	enableGuestUpload, ok := guestSettings.Settings["enable_guest_upload"].(bool)
	if !ok || !enableGuestUpload {
		return false, 0, errors.New(errors.CodeForbidden, "游客上传功能已禁用")
	}

	dailyLimit := 10 // 默认值
	if limit, ok := guestSettings.Settings["guest_daily_limit"].(float64); ok {
		dailyLimit = int(limit)
	}

	limit, err := s.getOrCreateLimit(fingerprint, dailyLimit)
	if err != nil {
		return false, 0, err
	}

	limit.ResetIfNeeded()

	if limit.DailyCount >= limit.DailyLimit {
		return false, 0, errors.New(errors.CodeUploadLimitExceeded, "每日上传次数已达上限")
	}

	remaining := dailyLimit - limit.DailyCount

	return true, remaining, nil
}

/* RecordUpload 记录游客上传 */
func (s *GuestUploadLimitService) RecordUpload(fingerprint string, fileSize int64) error {
	var limit models.GuestUploadLimit
	if err := database.DB.Where("fingerprint = ?", fingerprint).First(&limit).Error; err != nil {
		return err
	}

	limit.IncrementUploadCount()

	return database.DB.Save(&limit).Error
}

/* CheckIPUploadLimit 检查IP级别的上传限制 */
func (s *GuestUploadLimitService) CheckIPUploadLimit(ip string) (bool, error) {
	guestSettings, err := setting.GetSettingsByGroupAsMap("guest")
	if err != nil {
		logger.Error("CheckIPUploadLimit: 获取游客设置失败, err=%v", err)
		return false, errors.Wrap(err, errors.CodeDBQueryFailed, "获取游客设置失败")
	}

	ipDailyLimit := 50 // 默认值
	if limit, ok := guestSettings.Settings["guest_ip_daily_limit"].(float64); ok {
		ipDailyLimit = int(limit)
	}

	var count int64
	cfg := config.GetConfig()
	var dateCondition string
	if cfg.Database.Type == "sqlite" {
		dateCondition = "DATE(created_at) = DATE('now')"
	} else {
		dateCondition = "DATE(created_at) = CURDATE()"
	}
	err = database.DB.Model(&models.GuestUploadLog{}).
		Where("ip = ? AND "+dateCondition, ip).
		Count(&count).Error

	if err != nil {
		return false, errors.Wrap(err, errors.CodeDBQueryFailed, "查询IP上传次数失败")
	}

	if int(count) >= ipDailyLimit {
		return false, errors.New(errors.CodeUploadLimitExceeded, "IP每日上传次数已达上限")
	}

	return true, nil
}

func (s *GuestUploadLimitService) getOrCreateLimit(fingerprint string, dailyLimit int) (*models.GuestUploadLimit, error) {
	var limit models.GuestUploadLimit

	err := database.DB.Where("fingerprint = ?", fingerprint).First(&limit).Error
	if err == nil {
		return &limit, nil
	}

	if err == gorm.ErrRecordNotFound {
		limit = models.GuestUploadLimit{
			Fingerprint: fingerprint,
			DailyLimit:  dailyLimit,
		}
		if err := database.DB.Create(&limit).Error; err != nil {
			return nil, err
		}
		return &limit, nil
	}

	return nil, err
}

/* GetGuestUploadStats 获取游客上传统计信息 */
func (s *GuestUploadLimitService) GetGuestUploadStats(fingerprint string) (map[string]interface{}, error) {
	guestSettings, err := setting.GetSettingsByGroupAsMap("guest")
	if err != nil {
		logger.Error("GetGuestUploadStats: 获取游客设置失败, err=%v", err)
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取游客设置失败")
	}

	dailyLimit := 10 // 默认值
	if limit, ok := guestSettings.Settings["guest_daily_limit"].(float64); ok {
		dailyLimit = int(limit)
	}

	var limit models.GuestUploadLimit
	err = database.DB.Where("fingerprint = ?", fingerprint).First(&limit).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return map[string]interface{}{
				"daily_limit": dailyLimit,
				"daily_count": 0,
			}, nil
		}
		return nil, err
	}

	return map[string]interface{}{
		"daily_limit": limit.DailyLimit,
		"daily_count": limit.DailyCount,
	}, nil
}

var guestUploadLimitService *GuestUploadLimitService

/* GetGuestUploadLimitService 获取游客上传限制服务实例（延迟初始化） */
func GetGuestUploadLimitService() *GuestUploadLimitService {
	if guestUploadLimitService == nil {
		guestUploadLimitService = NewGuestUploadLimitService()
	}
	return guestUploadLimitService
}

func getGuestUploadLimitByFingerprint(fingerprint string) (*models.GuestUploadLimit, error) {
	var limit models.GuestUploadLimit
	if err := database.DB.Where("fingerprint = ?", fingerprint).First(&limit).Error; err != nil {
		return nil, err
	}
	return &limit, nil
}

/* CheckUploadLimitAndGetLimit 检查游客上传限制并返回limit记录 */
func (s *GuestUploadLimitService) CheckUploadLimitAndGetLimit(fingerprint string) (*models.GuestUploadLimit, bool, int, error) {
	guestSettings, err := setting.GetSettingsByGroupAsMap("guest")
	if err != nil {
		logger.Error("CheckUploadLimitAndGetLimit: 获取游客设置失败, err=%v", err)
		return nil, false, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "获取游客设置失败")
	}

	enableGuestUpload, ok := guestSettings.Settings["enable_guest_upload"].(bool)
	if !ok || !enableGuestUpload {
		return nil, false, 0, errors.New(errors.CodeForbidden, "游客上传功能已禁用")
	}

	dailyLimit := 10 // 默认值
	if limit, ok := guestSettings.Settings["guest_daily_limit"].(float64); ok {
		dailyLimit = int(limit)
	}

	limit, err := s.getOrCreateLimit(fingerprint, dailyLimit)
	if err != nil {
		return nil, false, 0, err
	}

	limit.ResetIfNeeded()

	if limit.DailyCount >= limit.DailyLimit {
		return limit, false, 0, errors.New(errors.CodeUploadLimitExceeded, "每日上传次数已达上限")
	}

	remaining := dailyLimit - limit.DailyCount

	return limit, true, remaining, nil
}
