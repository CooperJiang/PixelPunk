package user

import (
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
)

func UpdateUserProfile(userID uint, data map[string]interface{}) (*models.User, error) {
	db := database.DB

	if username, ok := data["username"].(string); ok && username != "" {
		var count int64
		if err := db.Model(&models.User{}).Where("username = ? AND id != ?", username, userID).Count(&count).Error; err != nil {
			return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "检查用户名是否存在失败")
		}

		if count > 0 {
			return nil, errors.New(errors.CodeUserExists, "用户名已被使用")
		}
	}

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, errors.New(errors.CodeUserNotFound, "用户不存在")
	}

	if err := db.Model(&user).Updates(data).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBUpdateFailed, "更新用户资料失败")
	}

	if err := db.First(&user, userID).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取更新后的用户资料失败")
	}

	return &user, nil
}

func GetUserProfile(userID uint) (*models.User, error) {
	db := database.DB

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, errors.New(errors.CodeUserNotFound, "用户不存在")
	}

	return &user, nil
}
