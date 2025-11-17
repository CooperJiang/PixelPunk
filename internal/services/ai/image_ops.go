package ai

import (
	"context"
	"fmt"
	"pixelpunk/internal/models"
	messageService "pixelpunk/internal/services/message"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/storage"

	"gorm.io/gorm"
)

func markFileForDeletion(tx *gorm.DB, fileID string) error {
	var file models.File
	if err := tx.Where("id = ?", fileID).First(&file).Error; err != nil {
		logger.Error("查询待删除违规文件失败: %v", err)
		return err
	}

	// 在事务中标记文件状态，然后在事务外执行删除
	// 注意：不能在事务内部调用 Commit()，应该让 Transaction() 自动提交
	if err := tx.Model(&models.File{}).
		Where("id = ?", fileID).
		Updates(map[string]interface{}{
			"status": "pending_deletion",
			"nsfw":   true,
		}).Error; err != nil {
		logger.Error("标记文件待删除失败: %v", err)
		return err
	}

	// 在事务外执行删除操作（避免事务锁定）
	// 使用 goroutine 异步执行，避免阻塞当前事务
	go func() {
		if err := executeFileDeletion(&file); err != nil {
			logger.Error("立即删除违规文件失败，文件已标记为待删除: %v", err)
			// 如果立即删除失败，文件已标记为 pending_deletion，由定时任务处理
		} else {
		}
	}()

	return nil
}

func executeFileDeletion(file *models.File) error {
	fileID := file.ID

	database.DB.Where("file_id = ?", fileID).Delete(&models.FileAIInfo{})

	database.DB.Where("file_id = ?", fileID).Delete(&models.FileGlobalTagRelation{})

	database.DB.Where("file_id = ?", fileID).Delete(&models.FileVector{})

	database.DB.Where("item_type = ? AND item_id = ?", "file", fileID).Delete(&models.ShareItem{})

	database.DB.Where("file_id = ?", fileID).Delete(&models.UploadSession{})

	storageService := storage.NewGlobalStorage()
	ctx := context.Background()
	if file.LocalFilePath != "" {
		if err := storageService.Delete(ctx, file.StorageProviderID, file.LocalFilePath); err != nil {
			logger.Warn("删除原图文件失败: %v", err)
		}
	}
	if file.LocalThumbPath != "" {
		if err := storageService.Delete(ctx, file.StorageProviderID, file.LocalThumbPath); err != nil {
			logger.Warn("删除缩略图文件失败: %v", err)
		}
	}

	if err := database.DB.Delete(&models.File{}, "id = ?", fileID).Error; err != nil {
		return fmt.Errorf("删除文件记录失败: %v", err)
	}

	return nil
}

func markFileForReview(tx *gorm.DB, fileID string, nsfwReason string) error {
	var file models.File
	if err := tx.Where("id = ?", fileID).First(&file).Error; err != nil {
		logger.Error("获取文件信息失败: %v", err)
		return err
	}

	err := tx.Model(&models.File{}).
		Where("id = ?", fileID).
		Updates(map[string]interface{}{
			"status": "pending_review",
			"nsfw":   true,
		}).Error
	if err != nil {
		return err
	}

	go sendFilePendingReviewNotification(file.UserID, fileID, file.OriginalName, nsfwReason)
	return nil
}

func sendFilePendingReviewNotification(userID uint, fileID, fileName, nsfwReason string) {
	variables := map[string]interface{}{
		"file_id":      fileID,
		"file_name":    fileName,
		"nsfw_reason":  nsfwReason,
		"related_type": "file",
		"related_id":   fileID,
	}
	msgService := messageService.GetMessageService()
	if err := msgService.SendTemplateMessage(userID, common.MessageTypeContentReviewPending, variables); err != nil {
		logger.Warn("发送文件待审核消息失败: userID=%d, fileID=%s, error=%v", userID, fileID, err)
	}
}

