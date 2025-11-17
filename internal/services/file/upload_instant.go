package file

import (
	"github.com/gin-gonic/gin"

	"pixelpunk/internal/models"
	"pixelpunk/internal/services/stats"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"

	"gorm.io/gorm"
)

/* CheckDuplicateResponse 重复检查响应 */
type CheckDuplicateResponse struct {
	Exists       bool                    `json:"exists"`
	OriginalFile *CheckDuplicateFileInfo `json:"original_file,omitempty"`
}

/* CheckDuplicateFileInfo 检查重复时返回的最小安全信息 */
type CheckDuplicateFileInfo struct {
	ID           string `json:"id"`
	OriginalName string `json:"original_name"`
	Size         int64  `json:"size"`
	Format       string `json:"format"`
}

func CheckDuplicateByMD5(userID uint, md5Hash, fileName string, fileSize int64) (*CheckDuplicateResponse, error) {
	var existingImage models.File

	err := database.DB.Where("user_id = ? AND md5_hash = ?", userID, md5Hash).
		Where("status <> ?", "pending_deletion").
		First(&existingImage).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &CheckDuplicateResponse{Exists: false}, nil
		}
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询重复文件失败")
	}

	return &CheckDuplicateResponse{
		Exists: true,
		OriginalFile: &CheckDuplicateFileInfo{
			ID:           existingImage.ID,
			OriginalName: existingImage.OriginalName,
			Size:         existingImage.Size,
			Format:       existingImage.Format,
		},
	}, nil
}

/* InstantUploadResponse 秒传响应 */
type InstantUploadResponse struct {
	ImageInfo *FileDetailResponse `json:"file_info"`
	IsInstant bool                `json:"is_instant"`
	Message   string              `json:"message"`
}

func validateInstantUploadRequest(ctx *UploadContext) error {
	if err := validateFolder(ctx); err != nil {
		return err
	}
	return processFolderPath(ctx)
}

/* InstantUpload 秒传上传（复用已有文件，执行完整AI处理） */
func InstantUpload(c *gin.Context, userID uint, md5Hash, fileName string, fileSize int64, folderID, accessLevel string, optimize bool) (*InstantUploadResponse, error) {
	checkResult, err := CheckDuplicateByMD5(userID, md5Hash, fileName, fileSize)
	if err != nil {
		return nil, err
	}

	if !checkResult.Exists {
		return nil, errors.New(errors.CodeFileNotFound, "未找到可复用的文件，请使用正常上传")
	}

	var originalImage models.File
	if err := database.DB.Where("user_id = ? AND md5_hash = ? AND (original_file_id = '' OR original_file_id IS NULL)",
		userID, md5Hash).
		Where("status <> ?", "pending_deletion").
		First(&originalImage).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询原始文件失败")
	}

	available, err := stats.CheckUserStorageAvailable(userID, fileSize)
	if err != nil {
		logger.Error("检查用户存储空间失败: %v", err)
		return nil, errors.Wrap(err, errors.CodeInternal, "检查用户存储空间失败")
	}
	if !available {
		return nil, errors.New(errors.CodeStorageLimitExceeded, "存储空间不足，无法上传文件")
	}

	if exceeded, err := checkDailyUploadLimit(userID, 1); err != nil {
		logger.Warn("检查每日上传限制失败: %v", err)
	} else if exceeded {
		return nil, errors.New(errors.CodeUploadLimitExceeded, "已达到每日上传限制")
	}

	if folderID == "null" {
		folderID = ""
	}

	ctx := CreateInstantUploadContext(c, userID, &originalImage, fileName, fileSize, folderID, accessLevel, optimize)

	if err := validateInstantUploadRequest(ctx); err != nil {
		return nil, err
	}

	response, err := completeFileUpload(ctx)
	if err != nil {
		return nil, err
	}

	return &InstantUploadResponse{
		ImageInfo: response,
		IsInstant: true,
		Message:   "秒传成功",
	}, nil
}

/* InstantUploadWithDuration 秒传文件（支持存储时长） */
func InstantUploadWithDuration(c *gin.Context, userID uint, md5Hash, fileName string, fileSize int64, folderID, accessLevel string, optimize bool, storageDuration string) (*InstantUploadResponse, error) {
	checkResult, err := CheckDuplicateByMD5(userID, md5Hash, fileName, fileSize)
	if err != nil {
		return nil, err
	}

	if !checkResult.Exists {
		return nil, errors.New(errors.CodeFileNotFound, "未找到可复用的文件，请使用正常上传")
	}

	var originalImage models.File
	if err := database.DB.Where("user_id = ? AND md5_hash = ? AND (original_file_id = '' OR original_file_id IS NULL)",
		userID, md5Hash).
		Where("status <> ?", "pending_deletion").
		First(&originalImage).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询原始文件失败")
	}

	available, err := stats.CheckUserStorageAvailable(userID, fileSize)
	if err != nil {
		logger.Error("检查用户存储空间失败: %v", err)
		return nil, errors.Wrap(err, errors.CodeInternal, "检查用户存储空间失败")
	}
	if !available {
		return nil, errors.New(errors.CodeStorageLimitExceeded, "存储空间不足，无法上传文件")
	}

	if exceeded, err := checkDailyUploadLimit(userID, 1); err != nil {
		logger.Warn("检查每日上传限制失败: %v", err)
	} else if exceeded {
		return nil, errors.New(errors.CodeUploadLimitExceeded, "已达到每日上传限制")
	}

	if folderID == "null" {
		folderID = ""
	}

	ctx := CreateInstantUploadContextWithDuration(c, userID, &originalImage, fileName, fileSize, folderID, accessLevel, optimize, storageDuration)

	if err := validateInstantUploadRequest(ctx); err != nil {
		return nil, err
	}

	response, err := completeFileUpload(ctx)
	if err != nil {
		return nil, err
	}

	return &InstantUploadResponse{
		ImageInfo: response,
		IsInstant: true,
		Message:   "秒传成功",
	}, nil
}
