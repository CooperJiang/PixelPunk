package user

import (
	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/config"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// getGreatestExpr 根据数据库类型返回兼容的 GREATEST 表达式
// SQLite 不支持 GREATEST()，使用 CASE WHEN 替代
func getGreatestExpr(column string, delta interface{}) clause.Expr {
	cfg := config.GetConfig()
	if cfg.Database.Type == "sqlite" {
		// SQLite 使用 CASE WHEN 表达式：确保结果不小于 0
		// CASE WHEN (total_bandwidth + ?) < 0 THEN 0 ELSE (total_bandwidth + ?) END
		return gorm.Expr("CASE WHEN ("+column+" + ?) < 0 THEN 0 ELSE ("+column+" + ?) END", delta, delta)
	}
	// MySQL 使用 GREATEST 函数
	return gorm.Expr("GREATEST("+column+" + ?, 0)", delta)
}

func InitUserUsageStats(userID uint) error {
	var count int64
	if err := database.DB.Model(&models.UserUsageStats{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBQueryFailed, "检查用户统计记录失败")
	}

	if count > 0 {
		return nil
	}

	stats := models.UserUsageStats{
		UserID:         userID,
		TotalImages:    0,
		TotalSize:      0,
		TotalBandwidth: 0,
		TotalViews:     0,
		UpdatedAt:      common.JSONTimeNow(),
	}

	if err := database.DB.Create(&stats).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBCreateFailed, "创建用户统计记录失败")
	}

	return nil
}

func GetUserUsageStats(userID uint) (*models.UserUsageStats, error) {
	var stats models.UserUsageStats

	err := database.DB.Where("user_id = ?", userID).First(&stats).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := InitUserUsageStats(userID); err != nil {
				return nil, err
			}
			err = database.DB.Where("user_id = ?", userID).First(&stats).Error
			if err != nil {
				return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取用户统计失败")
			}
		} else {
			return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取用户统计失败")
		}
	}

	return &stats, nil
}

func UpdateStorageUsage(userID uint, sizeDelta int64) error {
	result := database.DB.Model(&models.UserUsageStats{}).
		Where("user_id = ?", userID).
		UpdateColumn("total_size", getGreatestExpr("total_size", sizeDelta))

	if result.Error != nil {
		return errors.Wrap(result.Error, errors.CodeDBUpdateFailed, "更新存储使用量失败")
	}

	if result.RowsAffected == 0 {
		if err := InitUserUsageStats(userID); err != nil {
			return err
		}
		result = database.DB.Model(&models.UserUsageStats{}).
			Where("user_id = ?", userID).
			UpdateColumn("total_size", getGreatestExpr("total_size", sizeDelta))

		if result.Error != nil {
			return errors.Wrap(result.Error, errors.CodeDBUpdateFailed, "更新存储使用量失败")
		}
	}

	return nil
}

func UpdateBandwidthUsage(userID uint, bandwidthDelta int64) error {
	result := database.DB.Model(&models.UserUsageStats{}).
		Where("user_id = ?", userID).
		UpdateColumn("total_bandwidth", getGreatestExpr("total_bandwidth", bandwidthDelta))

	if result.Error != nil {
		return errors.Wrap(result.Error, errors.CodeDBUpdateFailed, "更新带宽使用量失败")
	}

	if result.RowsAffected == 0 {
		if err := InitUserUsageStats(userID); err != nil {
			return err
		}
		result = database.DB.Model(&models.UserUsageStats{}).
			Where("user_id = ?", userID).
			UpdateColumn("total_bandwidth", getGreatestExpr("total_bandwidth", bandwidthDelta))

		if result.Error != nil {
			return errors.Wrap(result.Error, errors.CodeDBUpdateFailed, "更新带宽使用量失败")
		}
	}

	return nil
}

func UpdateViewsUsage(userID uint, viewsDelta int64) error {
	result := database.DB.Model(&models.UserUsageStats{}).
		Where("user_id = ?", userID).
		UpdateColumn("total_views", getGreatestExpr("total_views", viewsDelta))

	if result.Error != nil {
		return errors.Wrap(result.Error, errors.CodeDBUpdateFailed, "更新文件浏览次数失败")
	}

	if result.RowsAffected == 0 {
		if err := InitUserUsageStats(userID); err != nil {
			return err
		}
		result = database.DB.Model(&models.UserUsageStats{}).
			Where("user_id = ?", userID).
			UpdateColumn("total_views", getGreatestExpr("total_views", viewsDelta))

		if result.Error != nil {
			return errors.Wrap(result.Error, errors.CodeDBUpdateFailed, "更新文件浏览次数失败")
		}
	}

	return nil
}

func UpdateFileUploadStats(tx *gorm.DB, userID uint, fileSize int64) error {
	var userStats models.UserUsageStats
	result := tx.Where("user_id = ?", userID).First(&userStats)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			userStats = models.UserUsageStats{
				UserID:      userID,
				TotalImages: 1,
				TotalSize:   fileSize,
				UpdatedAt:   common.JSONTimeNow(),
			}
			if err := tx.Create(&userStats).Error; err != nil {
				return errors.Wrap(err, errors.CodeDBCreateFailed, "创建用户统计记录失败")
			}
		} else {
			return errors.Wrap(result.Error, errors.CodeDBQueryFailed, "查询用户统计记录失败")
		}
	} else {
		userStats.TotalImages++
		userStats.TotalSize += fileSize
		userStats.UpdatedAt = common.JSONTimeNow()

		if err := tx.Save(&userStats).Error; err != nil {
			return errors.Wrap(err, errors.CodeDBUpdateFailed, "更新用户统计记录失败")
		}
	}

	return nil
}
