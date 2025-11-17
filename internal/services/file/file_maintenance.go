package file

// Maintenance queries and cleanup (no behavior change).

import (
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/activity"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"time"
)

func GetExpiredFiles() ([]models.File, error) {
	var expiredFiles []models.File
	now := time.Now()
	if err := database.DB.Where("expires_at IS NOT NULL AND expires_at < ?", now).Find(&expiredFiles).Error; err != nil {
		logger.Error("查询过期文件失败: %v", err)
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询过期文件失败")
	}
	return expiredFiles, nil
}

// CleanupExpiredImages 批量清理过期文件
func CleanupExpiredFiles() (int, int, error) {
	expiredFiles, err := GetExpiredFiles()
	if err != nil {
		return 0, 0, err
	}
	if len(expiredFiles) == 0 {
		return 0, 0, nil
	}

	userImages := make(map[uint][]models.File)
	guestImages := make([]models.File, 0)
	for _, img := range expiredFiles {
		if img.UserID != 0 {
			userImages[img.UserID] = append(userImages[img.UserID], img)
		} else {
			guestImages = append(guestImages, img)
		}
	}
	successCount, failedCount := 0, 0
	for userID, images := range userImages {
		deletedFiles := make([]models.File, 0)
		for _, img := range images {
			if deleteErr := DeleteFile(img.UserID, img.ID); deleteErr != nil {
				logger.Error("❌ 删除过期文件失败: ID=%s, 文件名=%s, 过期时间=%s, 错误=%v", img.ID, img.OriginalName, img.ExpiresAt.Format("2006-01-02 15:04:05"), deleteErr)
				failedCount++
			} else {
				successCount++
				deletedFiles = append(deletedFiles, img)
			}
		}
		if len(deletedFiles) > 0 {
			activity.LogImageExpired(userID, deletedFiles)
		}
	}
	if len(guestImages) > 0 {
		deletedGuestImages := make([]models.File, 0)
		for _, img := range guestImages {
			if deleteErr := DeleteNSFWFile(img.ID); deleteErr != nil {
				logger.Error("❌ 删除游客过期文件失败: ID=%s, 文件名=%s, 过期时间=%s, 错误=%v", img.ID, img.OriginalName, img.ExpiresAt.Format("2006-01-02 15:04:05"), deleteErr)
				failedCount++
			} else {
				successCount++
				deletedGuestImages = append(deletedGuestImages, img)
			}
		}
		if len(deletedGuestImages) > 0 {
			activity.LogGuestImageExpired(deletedGuestImages)
		}
	}
	return successCount, failedCount, nil
}
