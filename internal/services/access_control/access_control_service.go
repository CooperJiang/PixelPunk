package access_control

import (
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"

	"gorm.io/gorm"
)

/* GetUserAccessControl 获取用户的访问控制配置 */
func GetUserAccessControl(userID uint) (*models.UserAccessControl, error) {
	var config models.UserAccessControl
	result := database.DB.Where("user_id = ?", userID).First(&config)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			config = *models.CreateDefaultConfig(userID)
			if err := database.DB.Create(&config).Error; err != nil {
				logger.Error("创建用户访问控制默认配置失败: %v", err)
				return nil, errors.Wrap(err, errors.CodeDBCreateFailed, "创建用户访问控制配置失败")
			}
			return &config, nil
		}
		logger.Error("获取用户访问控制配置失败: %v", result.Error)
		return nil, errors.Wrap(result.Error, errors.CodeDBQueryFailed, "获取用户访问控制配置失败")
	}
	return &config, nil
}

/* UpdateUserAccessControl 更新或创建用户的访问控制配置 */
func UpdateUserAccessControl(config *models.UserAccessControl) error {
	var existingConfig models.UserAccessControl
	result := database.DB.Where("user_id = ?", config.UserID).First(&existingConfig)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			if err := database.DB.Create(config).Error; err != nil {
				logger.Error("创建用户访问控制配置失败: %v", err)
				return errors.Wrap(err, errors.CodeDBCreateFailed, "创建用户访问控制配置失败")
			}
			return nil
		}
		logger.Error("查询用户访问控制配置失败: %v", result.Error)
		return errors.Wrap(result.Error, errors.CodeDBQueryFailed, "查询用户访问控制配置失败")
	}

	config.ID = existingConfig.ID
	if err := database.DB.Save(config).Error; err != nil {
		logger.Error("更新用户访问控制配置失败: %v", err)
		return errors.Wrap(err, errors.CodeDBUpdateFailed, "更新用户访问控制配置失败")
	}

	return nil
}

/* DeleteUserAccessControl 删除用户的访问控制配置 */
func DeleteUserAccessControl(userID uint) error {
	result := database.DB.Where("user_id = ?", userID).Delete(&models.UserAccessControl{})
	if result.Error != nil {
		logger.Error("删除用户访问控制配置失败: %v", result.Error)
		return errors.Wrap(result.Error, errors.CodeDBDeleteFailed, "删除用户访问控制配置失败")
	}
	return nil
}

/* ApplyControlPreset 应用预设的访问控制配置 */
func ApplyControlPreset(userID uint, preset string) error {
	config, err := GetUserAccessControl(userID)
	if err != nil {
		return err
	}

	config.ApplyPreset(preset)

	return UpdateUserAccessControl(config)
}

/* CheckUserIP 检查IP是否在用户允许的范围内 */
func CheckUserIP(config *models.UserAccessControl, ip string, isIPInListFunc func(string, string) bool) bool {
	if !config.EnableIPCheck {
		return true // 未启用IP检查，直接通过
	}

	if config.IPWhitelist != "" {
		return isIPInListFunc(ip, config.IPWhitelist)
	}

	if config.IPBlacklist != "" {
		return !isIPInListFunc(ip, config.IPBlacklist)
	}

	return true // 默认通过
}

/* CheckUserDomain 检查域名是否在用户允许的范围内 */
func CheckUserDomain(config *models.UserAccessControl, domain string, referer string, isDomainInListFunc func(string, string) bool) bool {
	if !config.EnableRefererCheck {
		return true // 未启用防盗链检查，直接通过
	}

	if referer == "" {
		return config.AllowEmptyReferer
	}

	if domain == "" {
		return config.AllowEmptyReferer // 按照空Referer处理
	}

	if config.DomainWhitelist != "" {
		return isDomainInListFunc(domain, config.DomainWhitelist)
	}

	if config.DomainBlacklist != "" {
		return !isDomainInListFunc(domain, config.DomainBlacklist)
	}

	return true // 默认通过
}
