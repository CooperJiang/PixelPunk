package bandwidth

import (
	"encoding/json"
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/cache"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"
	"time"

	"gorm.io/gorm"
)

/* BandwidthService 带宽服务 */
type BandwidthService struct{}

/* GetUserSettings 获取用户设置（支持缓存） */
func (s *BandwidthService) GetUserSettings(userID uint) (models.UserSettings, error) {
	cacheKey := fmt.Sprintf("user_settings:%d", userID)

	if cachedData, err := cache.Get(cacheKey); err == nil {
		var settings models.UserSettings
		if err := json.Unmarshal([]byte(cachedData), &settings); err == nil {
			return settings, nil
		}
	}

	db := database.GetDB()
	var settings models.UserSettings
	err := db.Where("user_id = ?", userID).First(&settings).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			settings = models.UserSettings{
				UserID:         userID,
				StorageLimit:   models.DefaultStorageLimit,
				BandwidthLimit: models.DefaultBandwidthLimit,
			}
		} else {
			logger.Error("[BANDWIDTH_SERVICE] 查询用户设置失败: %v", err)
			return settings, err
		}
	}

	if settingsJSON, err := json.Marshal(settings); err == nil {
		cache.Set(cacheKey, string(settingsJSON), time.Duration(common.UserSettingsCacheExpire)*time.Second)
	}

	return settings, nil
}

/* GetCurrentYearMonth 获取当前年月 */
func (s *BandwidthService) GetCurrentYearMonth() (int, int, error) {
	// 直接使用Go的time包获取当前年月，避免数据库兼容性问题
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	return year, month, nil
}

/* GetMonthlyUsage 获取用户月度带宽使用情况（支持缓存） */
func (s *BandwidthService) GetMonthlyUsage(userID uint, year, month int) (models.UserBandwidthUsage, error) {
	cacheKey := fmt.Sprintf("bandwidth_usage:%d:%d:%d", userID, year, month)

	if cachedData, err := cache.Get(cacheKey); err == nil {
		var usage models.UserBandwidthUsage
		if err := json.Unmarshal([]byte(cachedData), &usage); err == nil {
			return usage, nil
		}
	}

	db := database.GetDB()
	var usage models.UserBandwidthUsage
	err := db.Where("user_id = ? AND year = ? AND month = ?",
		userID, year, month).First(&usage).Error

	if err == gorm.ErrRecordNotFound {
		usage = models.UserBandwidthUsage{
			UserID:    userID,
			Year:      year,
			Month:     month,
			UsedBytes: 0,
		}
	} else if err != nil {
		logger.Error("[BANDWIDTH_SERVICE] 查询带宽使用记录失败: %v", err)
		return usage, err
	}

	if usageJSON, err := json.Marshal(usage); err == nil {
		cache.Set(cacheKey, string(usageJSON), time.Duration(common.BandwidthUsageCacheExpire)*time.Second)
	}

	return usage, nil
}

/* CheckBandwidthAvailable 检查用户带宽是否可用 */
func (s *BandwidthService) CheckBandwidthAvailable(userID uint, estimatedBytes int64) (bool, error) {
	settings, err := s.GetUserSettings(userID)
	if err != nil {
		return false, err
	}

	year, month, err := s.GetCurrentYearMonth()
	if err != nil {
		return false, err
	}

	usage, err := s.GetMonthlyUsage(userID, year, month)
	if err != nil {
		return false, err
	}

	available := usage.UsedBytes+estimatedBytes <= settings.BandwidthLimit

	return available, nil
}

/* GetUserBandwidthStats 获取用户带宽统计（兼容旧接口） */
func (s *BandwidthService) GetUserBandwidthStats(userID uint) (models.UserBandwidthUsage, error) {
	year, month, err := s.GetCurrentYearMonth()
	if err != nil {
		return models.UserBandwidthUsage{}, err
	}

	return s.GetMonthlyUsage(userID, year, month)
}

/* RecordBandwidthTransfer 记录带宽传输（仅更新月度统计） */
func (s *BandwidthService) RecordBandwidthTransfer(userID uint, bytesCount int64) error {
	db := database.GetDB()

	year, month, err := s.GetCurrentYearMonth()
	if err != nil {
		return err
	}

	// 使用事务更新月度统计
	return db.Transaction(func(tx *gorm.DB) error {
		err := s.updateMonthlyUsage(tx, userID, year, month, bytesCount)
		if err != nil {
			logger.Error("[BANDWIDTH_SERVICE] 更新月度带宽使用失败: %v", err)
			return err
		}

		s.clearBandwidthCache(userID, year, month)
		return nil
	})
}

func (s *BandwidthService) updateMonthlyUsage(tx *gorm.DB, userID uint, year, month int, bytes int64) error {
	var usage models.UserBandwidthUsage
	err := tx.Where("user_id = ? AND year = ? AND month = ?",
		userID, year, month).First(&usage).Error

	if err == gorm.ErrRecordNotFound {
		usage = models.UserBandwidthUsage{
			UserID:    userID,
			Year:      year,
			Month:     month,
			UsedBytes: bytes,
		}
		return tx.Create(&usage).Error
	} else if err != nil {
		return err
	} else {
		return tx.Model(&usage).Update("used_bytes", usage.UsedBytes+bytes).Error
	}
}

func (s *BandwidthService) clearBandwidthCache(userID uint, year, month int) {
	cacheKey := fmt.Sprintf("bandwidth_usage:%d:%d:%d", userID, year, month)
	cache.Del(cacheKey)
}

var Service = &BandwidthService{}
