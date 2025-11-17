package file

/* File deletion and cleanup helpers (no behavior change). */

import (
	"context"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/storage"
	"pixelpunk/pkg/vector"

	"gorm.io/gorm"
)

func cleanupFileResources(fileID string, file models.File, totalReferences int64) {
	cleanupFileLogs(fileID)
	cleanupFileShares(fileID)
	cleanupFileUploadSessions(fileID)
	cleanupFileVectors(fileID)
	if totalReferences == 0 {
		cleanupPhysicalFiles(file)
	}
}

func cleanupFileLogs(fileID string) {
	db := database.DB
	if err := db.Unscoped().Where("file_id = ?", fileID).Delete(&models.VectorProcessingLog{}).Error; err != nil {
		logger.Error("删除向量处理日志失败 [%s]: %v", fileID, err)
	}
	if err := db.Unscoped().Where("file_id = ?", fileID).Delete(&models.FileDownloadLog{}).Error; err != nil {
		logger.Error("删除文件下载日志失败 [%s]: %v", fileID, err)
	}
	if err := db.Unscoped().Where("file_id = ?", fileID).Delete(&models.FileTaggingLog{}).Error; err != nil {
		logger.Error("删除文件打标日志失败 [%s]: %v", fileID, err)
	}
}

func cleanupFileShares(fileID string) {
	db := database.DB
	if err := db.Unscoped().Where("item_type = ? AND item_id = ?", "file", fileID).Delete(&models.ShareItem{}).Error; err != nil {
		logger.Error("删除分享项目失败 [%s]: %v", fileID, err)
	}
}

func cleanupFileUploadSessions(fileID string) {
	db := database.DB
	if err := db.Model(&models.UploadSession{}).Where("file_id = ?", fileID).Update("file_id", "").Error; err != nil {
		logger.Error("清理上传会话文件引用失败 [%s]: %v", fileID, err)
	}
}

func cleanupFileVectors(fileID string) {
	ve := vector.GetGlobalVectorEngine()
	if ve == nil || !ve.IsEnabled() {
		logger.Warn("向量引擎不可用，跳过向量数据清理 [%s]", fileID)
		return
	}
	if err := ve.DeleteVector(fileID); err != nil {
		logger.Error("从向量数据库删除向量失败 [%s]: %v", fileID, err)
	} else {
	}
}

func cleanupPhysicalFiles(file models.File) {
	if file.StorageProviderID == "" {
		logger.Error("文件 %s 缺少存储提供者ID，无法删除", file.ID)
		return
	}
	ctx := context.Background()
	st := storage.NewGlobalStorage()
	if file.URL != "" {
		if err := st.Delete(ctx, file.StorageProviderID, file.URL); err != nil {
			logger.Error("删除存储文件失败 %s: %v", file.URL, err)
		} else {
		}
	}
	if file.ThumbURL != "" {
		if err := st.Delete(ctx, file.StorageProviderID, file.ThumbURL); err != nil {
			logger.Error("删除存储缩略图失败 %s: %v", file.ThumbURL, err)
		} else {
		}
	}
}

/* DeleteNSFWFile 自动删除违规文件（被AI标记为NSFW） */
func DeleteNSFWFile(fileID string) error {
	var file models.File
	if err := database.DB.Where("id = ?", fileID).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.CodeFileNotFound, "文件不存在")
		}
		return errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件失败")
	}
	logger.Warn("违规文件自动删除：用户ID=%d, 文件ID=%s, 文件名=%s", file.UserID, fileID, file.OriginalName)
	return deleteFileWithCascade(&file, file.UserID)
}
