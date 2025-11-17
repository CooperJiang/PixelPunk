package share

import (
	"pixelpunk/internal/controllers/share/dto"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/utils"
	"time"

	"gorm.io/gorm"
)

func CreateShare(userID uint, req *dto.CreateShareDTO) (models.Share, error) {
	shareKey := utils.GenerateRandomString(16)

	for {
		var count int64
		if err := database.DB.Model(&models.Share{}).Where("share_key = ?", shareKey).Count(&count).Error; err != nil {
			return models.Share{}, err
		}
		if count == 0 {
			break
		}
		shareKey = utils.GenerateRandomString(16)
	}

	share := models.Share{
		ID:                   generateID(),
		UserID:               userID,
		ShareKey:             shareKey,
		Name:                 req.Name,
		Description:          req.Description,
		Password:             req.Password,
		ExpiredDays:          req.ExpiredDays,
		MaxViews:             req.MaxViews,
		Status:               common.ShareStatusNormal,
		CollectVisitorInfo:   req.CollectVisitorInfo,
		NotificationOnAccess: req.NotificationOnAccess,
	}

	if req.NotificationOnAccess && req.NotificationThreshold > 0 {
		share.NotificationThreshold = req.NotificationThreshold
	} else if req.NotificationOnAccess {
		share.NotificationThreshold = 100
	}

	if req.ExpiredDays > 0 {
		expiredAt := time.Now().AddDate(0, 0, req.ExpiredDays)
		jsonTime := common.JSONTime(expiredAt)
		share.ExpiredAt = &jsonTime
	}

	// 使用 GORM Transaction 方法替代手动事务管理，确保 SQLite 兼容性
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&share).Error; err != nil {
			return err
		}

		for i, item := range req.Items {
			shareItem := models.ShareItem{
				ID:        generateID(),
				ShareID:   share.ID,
				ItemType:  item.ItemType,
				ItemID:    item.ItemID,
				SortOrder: i,
			}

			if err := tx.Create(&shareItem).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return models.Share{}, err
	}

	return share, nil
}
