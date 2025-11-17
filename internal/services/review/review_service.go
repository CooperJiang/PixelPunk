package review

import (
	"fmt"
	"pixelpunk/internal/models"
	messageService "pixelpunk/internal/services/message"
	"pixelpunk/internal/services/setting"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"

	"gorm.io/gorm"
)

/* ApproveFileWithLog 批准文件并记录审核日志 */
func ApproveFileWithLog(fileID string, auditorID uint, reason string) error {
	db := database.GetDB()

	return db.Transaction(func(tx *gorm.DB) error {
		var file models.File
		if err := tx.Where("id = ? AND status = ?", fileID, "pending_review").First(&file).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("待审核文件不存在")
			}
			return fmt.Errorf("查询待审核文件失败: %v", err)
		}

		var nsfwScore, nsfwThreshold *float64
		var isNSFW *bool

		var aiInfo models.FileAIInfo
		if err := tx.Where("file_id = ?", fileID).First(&aiInfo).Error; err == nil {
			nsfwScore = &aiInfo.NSFWScore
			isNSFW = &aiInfo.IsNSFW
		}

		if threshold, err := getNSFWThreshold(); err == nil {
			nsfwThreshold = &threshold
		}

		reviewLog := &models.ReviewLog{
			FileID:        fileID,
			AuditorID:     auditorID,
			UploaderID:    file.UserID,
			Action:        "approve",
			Reason:        reason,
			NSFWScore:     nsfwScore,
			NSFWThreshold: nsfwThreshold,
			IsNSFW:        isNSFW,
		}

		if err := tx.Create(reviewLog).Error; err != nil {
			return fmt.Errorf("创建审核记录失败: %v", err)
		}

		if err := tx.Model(&models.File{}).
			Where("id = ? AND status = ?", fileID, "pending_review").
			Updates(map[string]interface{}{
				"status": "active",
				"nsfw":   false,
			}).Error; err != nil {
			return fmt.Errorf("批准文件失败: %v", err)
		}

		go sendFileReviewNotification(file.UserID, fileID, file.OriginalName, "approve", "", auditorID)

		return nil
	})
}

/* RejectFileWithLog 拒绝文件并记录审核日志（默认软删除） */
func RejectFileWithLog(fileID string, auditorID uint, reason string, hardDelete bool) error {
	db := database.GetDB()

	var fileToDelete models.File

	// 使用 GORM Transaction 方法替代手动事务管理，确保 SQLite 兼容性
	err := db.Transaction(func(tx *gorm.DB) error {
		var file models.File
		if err := tx.Where("id = ? AND status = ?", fileID, "pending_review").First(&file).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("待审核文件不存在，请确认文件ID是否正确且处于待审核状态")
			}
			return fmt.Errorf("查询待审核文件失败: %v", err)
		}

		// 保存文件信息用于后续删除
		fileToDelete = file

		var nsfwScore, nsfwThreshold *float64
		var isNSFW *bool

		var aiInfo models.FileAIInfo
		if err := tx.Where("file_id = ?", fileID).First(&aiInfo).Error; err == nil {
			nsfwScore = &aiInfo.NSFWScore
			isNSFW = &aiInfo.IsNSFW
		}

		if threshold, err := getNSFWThreshold(); err == nil {
			nsfwThreshold = &threshold
		}

		deleteType := "soft"
		if hardDelete {
			deleteType = "hard"
		}

		reviewLog := &models.ReviewLog{
			FileID:        fileID,
			AuditorID:     auditorID,
			UploaderID:    file.UserID,
			Action:        "reject",
			DeleteType:    deleteType,
			Reason:        reason,
			NSFWScore:     nsfwScore,
			NSFWThreshold: nsfwThreshold,
			IsNSFW:        isNSFW,
		}

		if err := tx.Create(reviewLog).Error; err != nil {
			return fmt.Errorf("创建审核记录失败: %v", err)
		}

		if hardDelete {
			// 硬删除：在事务中标记，然后在事务外执行删除
			// 注意：不能在事务内部调用 Commit()，应该让 Transaction() 自动提交
			if err := tx.Model(&models.File{}).
				Where("id = ?", fileID).
				Updates(map[string]interface{}{
					"status": "pending_deletion",
				}).Error; err != nil {
				return fmt.Errorf("标记文件待删除失败: %v", err)
			}
		} else {
			// 软删除：在事务内完成
			if err := tx.Model(&models.File{}).
				Where("id = ?", fileID).
				Updates(map[string]interface{}{
					"status": "deleted",
				}).Error; err != nil {
				return fmt.Errorf("软删除文件失败: %v", err)
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 在事务外执行硬删除操作（避免事务锁定）
	if hardDelete {
		// 使用 goroutine 异步执行硬删除，避免阻塞
		go func() {
			if err := executeFileHardDeletion(&fileToDelete); err != nil {
				logger.Error("硬删除文件失败: fileID=%s, error=%v", fileID, err)
			} else {
				go sendFileReviewNotification(fileToDelete.UserID, fileID, fileToDelete.OriginalName, "reject", reason, auditorID)
			}
		}()
	} else {
		go sendFileReviewNotification(fileToDelete.UserID, fileID, fileToDelete.OriginalName, "reject", reason, auditorID)
	}

	return nil
}

/* BatchReviewFilesWithLog 批量审核文件并记录审核日志 */
func BatchReviewFilesWithLog(fileIDs []string, action string, auditorID uint, reason string, hardDelete bool) (map[string]string, error) {
	results := make(map[string]string)

	for _, fileID := range fileIDs {
		var err error
		switch action {
		case "approve":
			err = ApproveFileWithLog(fileID, auditorID, reason)
		case "reject":
			err = RejectFileWithLog(fileID, auditorID, reason, hardDelete)
		default:
			err = fmt.Errorf("无效的审核操作: %s", action)
		}

		if err != nil {
			results[fileID] = err.Error()
		} else {
			results[fileID] = "success"
		}
	}

	return results, nil
}

/* CreateReviewLog 创建审核记录 */
func CreateReviewLog(log *models.ReviewLog) error {
	return database.GetDB().Create(log).Error
}

/* GetReviewLogsByFileID 根据文件ID获取审核记录 */
func GetReviewLogsByFileID(fileID string) ([]models.ReviewLog, error) {
	var logs []models.ReviewLog
	err := database.GetDB().Where("file_id = ?", fileID).
		Preload("Auditor").
		Preload("Uploader").
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

/* HardDeleteSoftDeletedFile 对已软删除的文件执行硬删除 */
func HardDeleteSoftDeletedFile(fileID string, operatorID uint) error {
	db := database.GetDB()

	var file models.File
	if err := db.Where("id = ? AND status = ?", fileID, "deleted").First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("找不到已软删除的文件")
		}
		return fmt.Errorf("查询文件失败: %v", err)
	}

	reviewLog := &models.ReviewLog{
		FileID:     fileID,
		AuditorID:  operatorID,
		UploaderID: file.UserID,
		Action:     "reject",
		DeleteType: "hard",
		Reason:     "管理员执行硬删除",
	}

	if err := db.Create(reviewLog).Error; err != nil {
		return fmt.Errorf("创建审核记录失败: %v", err)
	}

	if err := executeFileHardDeletion(&file); err != nil {
		return fmt.Errorf("物理删除失败: %v", err)
	}

	go sendHardDeleteNotification(file.UserID, fileID, file.OriginalName, "管理员执行硬删除")

	return nil
}

/* GetReviewLogsByAuditorID 根据审核员ID获取审核记录 */
func GetReviewLogsByAuditorID(auditorID uint, page, pageSize int) ([]models.ReviewLog, int64, error) {
	var logs []models.ReviewLog
	var total int64

	db := database.GetDB().Where("auditor_id = ?", auditorID)

	if err := db.Model(&models.ReviewLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := db.Preload("File").
		Preload("Uploader").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&logs).Error

	return logs, total, err
}

/* GetReviewStats 获取审核统计信息 */
func GetReviewStats() (map[string]interface{}, error) {
	var approveCount, rejectCount int64

	db := database.GetDB().Model(&models.ReviewLog{})

	if err := db.Where("action = ?", "approve").Count(&approveCount).Error; err != nil {
		return nil, err
	}

	if err := db.Where("action = ?", "reject").Count(&rejectCount).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"approve_count": approveCount,
		"reject_count":  rejectCount,
		"total_count":   approveCount + rejectCount,
	}, nil
}

/* HardDeleteFile 强制硬删除文件 */
func HardDeleteFile(fileID string, auditorID uint, reason string) error {
	return RejectFileWithLog(fileID, auditorID, reason, true)
}

func getNSFWThreshold() (float64, error) {
	// 直接从数据库读取配置（绕过缓存）
	threshold := setting.GetFloatDirectFromDB("ai", "nsfw_threshold", 0.6)
	return threshold, nil
}

func executeFileHardDeletion(file *models.File) error {
	fileID := file.ID
	db := database.GetDB()

	db.Where("file_id = ?", fileID).Delete(&models.FileAIInfo{})

	db.Where("file_id = ?", fileID).Delete(&models.FileGlobalTagRelation{})

	db.Where("file_id = ?", fileID).Delete(&models.FileVector{})

	db.Where("item_type = ? AND item_id = ?", "file", fileID).Delete(&models.ShareItem{})

	db.Where("file_id = ?", fileID).Delete(&models.UploadSession{})

	if err := deletePhysicalFiles(file); err != nil {
		logger.Warn("删除物理文件失败，但继续删除数据库记录: %v", err)
	}

	if err := db.Delete(&models.File{}, "id = ?", fileID).Error; err != nil {
		return fmt.Errorf("删除文件记录失败: %v", err)
	}

	return nil
}

func deletePhysicalFiles(file *models.File) error {
	return nil
}

func sendFileReviewNotification(userID uint, fileID, fileName, action, reason string, auditorID uint) {
	var messageType string
	variables := map[string]interface{}{
		"file_id":      fileID,
		"file_name":    fileName,
		"related_type": "file",
		"related_id":   fileID,
	}

	switch action {
	case "approve":
		messageType = common.MessageTypeContentReviewApproved
	case "reject":
		messageType = common.MessageTypeContentReviewRejected
		variables["reason"] = reason
		variables["review_id"] = fmt.Sprintf("%s_%d", fileID, auditorID)
	default:
		logger.Warn("未知的审核操作类型: %s, 跳过发送消息", action)
		return
	}

	msgService := messageService.GetMessageService()
	if err := msgService.SendTemplateMessage(userID, messageType, variables); err != nil {
		logger.Warn("发送文件审核消息失败: userID=%d, fileID=%s, action=%s, error=%v", userID, fileID, action, err)
	} else {
	}
}

func sendHardDeleteNotification(userID uint, fileID, fileName, reason string) {
	variables := map[string]interface{}{
		"file_id":      fileID,
		"file_name":    fileName,
		"reason":       reason,
		"related_type": "file",
		"related_id":   fileID,
	}

	msgService := messageService.GetMessageService()
	if err := msgService.SendTemplateMessage(userID, common.MessageTypeFileHardDeletedByAdmin, variables); err != nil {
		logger.Warn("发送文件硬删除消息失败: userID=%d, fileID=%s, error=%v", userID, fileID, err)
	} else {
	}
}
