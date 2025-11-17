package file

/* Mutation helpers (no behavior change). */

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"pixelpunk/internal/models"
	storageChannelService "pixelpunk/internal/services/storage"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

/* UpdateFile 更新文件信息（名称/文件夹/访问级别） */
func UpdateFile(userID uint, fileID, name, folderID, accessLevel string) (*FileDetailResponse, error) {
	var file models.File
	if err := database.DB.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.CodeFileNotFound, "文件不存在")
		}
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件失败")
	}
	if folderID == "null" {
		folderID = ""
	}
	file.FolderID = folderID

	if name != "" {
		file.DisplayName = name
	}
	if accessLevel != "" {
		file.AccessLevel = accessLevel
	}
	if err := database.DB.Save(&file).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBUpdateFailed, "保存文件信息失败")
	}

	var stats models.FileStats
	if err := database.DB.Where("file_id = ?", fileID).First(&stats).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			stats = models.FileStats{FileID: fileID}
		}
	}
	aiInfo, _ := GetFileAIInfo(file.ID)
	resp := BuildFileDetailResponse(file, stats.Views, aiInfo)
	return &resp, nil
}

/* ToggleFileAccessLevel 切换文件访问级别（新语义） */
func ToggleFileAccessLevel(userID uint, fileID string) (*FileDetailResponse, error) {
	var file models.File
	if err := database.DB.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.CodeFileNotFound, "文件不存在")
		}
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件失败")
	}
	switch file.AccessLevel {
	case AccessPublic:
		file.AccessLevel = AccessPrivate
	case AccessPrivate:
		file.AccessLevel = AccessProtected
	case AccessProtected:
		file.AccessLevel = AccessPublic
	default:
		file.AccessLevel = AccessPublic
	}
	if err := database.DB.Save(&file).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBUpdateFailed, "更新文件失败")
	}
	var stats models.FileStats
	if err := database.DB.Where("file_id = ?", fileID).First(&stats).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			stats = models.FileStats{FileID: fileID}
		}
	}
	aiInfo, _ := GetFileAIInfo(file.ID)
	resp2 := BuildFileDetailResponse(file, stats.Views, aiInfo)
	return &resp2, nil
}

/* DeleteFile 删除文件（异步标记+后台删除） */
func DeleteFile(userID uint, fileID string) error {
	var file models.File
	if err := database.DB.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.CodeFileNotFound, "文件不存在")
		}
		return errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件失败")
	}
	if err := database.DB.Model(&models.File{}).Where("id = ? AND user_id = ?", fileID, userID).Updates(map[string]interface{}{"status": StatusPendingDeletion}).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBUpdateFailed, "标记文件为待删除失败")
	}
	imgCopy := file
	go func() {
		if err := deleteFileWithCascade(&imgCopy, userID); err != nil {
			logger.Warn("后台删除文件失败，将由定时任务兜底处理，file=%s err=%v", fileID, err)
		}
	}()
	return nil
}

/* GetTotalFileCount 获取文件总数（新语义） */
func GetTotalFileCount() (int64, error) {
	var count int64
	err := database.DB.Model(&models.File{}).Where("status IS NULL OR status <> ?", StatusPendingDeletion).Count(&count).Error
	if err != nil {
		return 0, errors.Wrap(err, errors.CodeDBQueryFailed, "获取文件总数失败")
	}
	return count, nil
}

/* GetTotalStorageUsage 获取文件总存储使用量（通过渠道模块统计） */
func GetTotalStorageUsage() (int64, string, error) {
	totalStorage, formattedStorage, err := storageChannelService.GetTotalStorageUsageByChannels()
	if err != nil {
		return 0, "", errors.Wrap(err, errors.CodeServiceUnavailable, "获取存储使用量失败")
	}

	return totalStorage, formattedStorage, nil
}

/* BatchDeleteUserFiles 用户批量删除自己的文件 */
func BatchDeleteUserFiles(userID uint, fileIDs []string) ([]string, []string, error) {
	var successIds []string
	var failIds []string
	for _, fileID := range fileIDs {
		if err := DeleteFile(userID, fileID); err != nil {
			logger.Error("用户批量删除文件失败 - 用户ID: %d, 文件ID: %s, 错误: %v", userID, fileID, err)
			failIds = append(failIds, fileID)
		} else {
			successIds = append(successIds, fileID)
		}
	}
	return successIds, failIds, nil
}

/* CleanupPendingDeletionFiles 查找并删除标记为待删除的文件 */
func CleanupPendingDeletionFiles(maxImages int) (int, error) {
	var imageIDs []string
	query := database.DB.Model(&models.File{}).Where("status = ?", "pending_deletion").Select("id")
	if maxImages > 0 {
		query = query.Limit(maxImages)
	}
	if err := query.Find(&imageIDs).Error; err != nil {
		return 0, errors.Wrap(err, errors.CodeDBQueryFailed, "查询待删除文件失败")
	}
	deleted := 0
	for _, id := range imageIDs {
		if err := DeleteNSFWFile(id); err == nil {
			deleted++
		}
	}
	return deleted, nil
}

/* UploadAvatar 上传用户头像（系统目录） */
func UploadAvatar(file *multipart.FileHeader) (string, string, error) {
	if file.Size > 2*1024*1024 {
		return "", "", errors.New(errors.CodeFileTooLarge, "头像大小不能超过2MB")
	}
	fileExt := strings.ToLower(filepath.Ext(file.Filename))
	if fileExt != ".jpg" && fileExt != ".jpeg" && fileExt != ".png" && fileExt != ".gif" && fileExt != ".webp" {
		return "", "", errors.New(errors.CodeFileTypeNotSupported, "只支持jpg、jpeg、png、webp和gif格式的文件")
	}
	src, err := file.Open()
	if err != nil {
		return "", "", errors.Wrap(err, errors.CodeFileUploadFailed, "打开上传文件失败")
	}
	defer src.Close()
	if err := os.MkdirAll(AvatarUploadDir, 0755); err != nil {
		return "", "", errors.Wrap(err, errors.CodeFileUploadFailed, "创建头像目录失败")
	}
	timestamp := time.Now().UnixNano() / 1000000
	randomStr := utils.GenerateRandomString(8)
	newFileName := fmt.Sprintf("avatar_%d_%s%s", timestamp, randomStr, fileExt)
	filePath := filepath.Join(AvatarUploadDir, newFileName)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", "", errors.Wrap(err, errors.CodeFileUploadFailed, "创建目标文件失败")
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return "", "", errors.Wrap(err, errors.CodeFileUploadFailed, "保存文件失败")
	}
	publicURL := fmt.Sprintf("/file/avatar/%s", newFileName)
	fullURL := utils.GetSystemFileURL(publicURL)
	return publicURL, fullURL, nil
}

/* UploadAdminFile 管理员文件上传（系统目录） */
func UploadAdminFile(file *multipart.FileHeader) (string, string, error) {
	if file.Size > 30*1024*1024 {
		return "", "", errors.New(errors.CodeFileTooLarge, "文件大小不能超过30MB")
	}
	src, err := file.Open()
	if err != nil {
		return "", "", errors.Wrap(err, errors.CodeFileUploadFailed, "打开上传文件失败")
	}
	defer src.Close()
	if err := os.MkdirAll(FileUploadDir, 0755); err != nil {
		return "", "", errors.Wrap(err, errors.CodeFileUploadFailed, "创建文件目录失败")
	}
	timestamp := time.Now().UnixNano() / 1000000
	randomStr := utils.GenerateRandomString(8)
	fileExt := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("file_%d_%s%s", timestamp, randomStr, fileExt)
	filePath := filepath.Join(FileUploadDir, newFileName)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", "", errors.Wrap(err, errors.CodeFileUploadFailed, "创建目标文件失败")
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		os.Remove(filePath)
		return "", "", errors.Wrap(err, errors.CodeFileUploadFailed, "文件保存失败")
	}
	relativePath := fmt.Sprintf("/file/admin/%s", newFileName)
	fullURL := utils.GetSystemFileURL(relativePath)
	return relativePath, fullURL, nil
}
