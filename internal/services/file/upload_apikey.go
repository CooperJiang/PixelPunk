package file

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"pixelpunk/internal/models"
	"pixelpunk/internal/services/apikey"
	"pixelpunk/internal/services/folder"
	"pixelpunk/internal/services/stats"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
)

/* APIKeyUploadResult API密钥上传结果 */
type APIKeyUploadResult struct {
	Uploaded         []*ExternalAPIFileResponse `json:"uploaded,omitempty"`          // 上传成功的文件（多文件时使用）
	UploadedSingle   *ExternalAPIFileResponse   `json:"uploaded_single,omitempty"`   // 单文件上传成功时使用
	OversizedFiles   []string                   `json:"oversized_files,omitempty"`   // 超出大小限制的文件
	UnsupportedFiles []string                   `json:"unsupported_files,omitempty"` // 格式不支持的文件
	InvalidFiles     []string                   `json:"invalid_files,omitempty"`     // 其他无效文件（空文件、无效文件名等）
	SizeLimit        string                     `json:"size_limit,omitempty"`        // 大小限制（可读格式）
	UploadErrors     []string                   `json:"upload_errors,omitempty"`     // 上传过程中的错误
	Message          string                     `json:"message,omitempty"`           // 处理结果消息
	IsSingleUpload   bool                       `json:"-"`                           // 内部字段，标识是否为单文件上传
}

/* FileValidationResult 文件验证结果 */
type FileValidationResult struct {
	ValidFiles       []*multipart.FileHeader `json:"valid_files,omitempty"`
	OversizedFiles   []string                `json:"oversized_files,omitempty"`   // 大小超限的文件
	UnsupportedFiles []string                `json:"unsupported_files,omitempty"` // 格式不支持的文件
	InvalidFiles     []string                `json:"invalid_files,omitempty"`     // 其他无效文件
}

/* UploadFileWithAPIKey 使用API密钥上传文件 */
func UploadFileWithAPIKey(c *gin.Context, key *models.APIKey, folderID, filePath, accessLevel string, optimize bool, files []*multipart.FileHeader, singleFile *multipart.FileHeader) (*APIKeyUploadResult, error) {
	targetFolderID, err := determineTargetFolder(key, folderID, filePath)
	if err != nil {
		return &APIKeyUploadResult{Message: err.Error()}, err
	}

	if len(files) > 0 {
		return processMultipleFilesUpload(c, key, targetFolderID, accessLevel, optimize, files)
	}

	if singleFile != nil {
		return processSingleFileUpload(c, key, targetFolderID, accessLevel, optimize, singleFile)
	}

	return &APIKeyUploadResult{Message: "未检测到上传文件"}, errors.New(errors.CodeInvalidParameter, "未检测到上传文件")
}

func determineTargetFolder(key *models.APIKey, folderID, filePath string) (string, error) {
	if filePath != "" {
		return folder.CreateFolderByPath(key.UserID, filePath)
	}
	if folderID != "" && folderID != "null" {
		return folderID, nil
	}
	return key.FolderID, nil
}

func processMultipleFilesUpload(c *gin.Context, key *models.APIKey, folderID, accessLevel string, optimize bool, files []*multipart.FileHeader) (*APIKeyUploadResult, error) {
	result := &APIKeyUploadResult{}
	validationResult := validateFiles(key, files)

	if len(validationResult.ValidFiles) == 0 {
		result.OversizedFiles = validationResult.OversizedFiles
		result.UnsupportedFiles = validationResult.UnsupportedFiles
		result.InvalidFiles = validationResult.InvalidFiles

		if len(validationResult.OversizedFiles) > 0 {
			result.SizeLimit = fmt.Sprintf("%.1fMB", float64(key.SingleFileLimit)/1024/1024)
		}

		var errorMessages []string
		if len(validationResult.OversizedFiles) > 0 {
			errorMessages = append(errorMessages, fmt.Sprintf("%d个文件超过大小限制", len(validationResult.OversizedFiles)))
		}
		if len(validationResult.UnsupportedFiles) > 0 {
			errorMessages = append(errorMessages, fmt.Sprintf("%d个文件格式不支持", len(validationResult.UnsupportedFiles)))
		}
		if len(validationResult.InvalidFiles) > 0 {
			errorMessages = append(errorMessages, fmt.Sprintf("%d个文件无效", len(validationResult.InvalidFiles)))
		}

		result.Message = "所有文件均无法上传: " + strings.Join(errorMessages, ", ")
		return result, errors.New(errors.CodeFileTooLarge, result.Message)
	}

	if err := validateStorageAndUploadLimits(key, validationResult.ValidFiles); err != nil {
		return result, err
	}

	return uploadValidFiles(c, key, folderID, accessLevel, optimize, validationResult)
}

func validateFiles(key *models.APIKey, files []*multipart.FileHeader) *FileValidationResult {
	result := &FileValidationResult{
		ValidFiles:       []*multipart.FileHeader{},
		OversizedFiles:   []string{},
		UnsupportedFiles: []string{},
		InvalidFiles:     []string{},
	}

	for _, file := range files {
		if key != nil && key.SingleFileLimit > 0 && file.Size > key.SingleFileLimit {
			result.OversizedFiles = append(result.OversizedFiles, file.Filename)
			continue
		}

		fileExt := strings.ToLower(filepath.Ext(file.Filename))
		if !isValidFileType(fileExt) {
			result.UnsupportedFiles = append(result.UnsupportedFiles, file.Filename)
			continue
		}

		if file.Size == 0 {
			result.InvalidFiles = append(result.InvalidFiles, file.Filename)
			continue
		}

		if strings.TrimSpace(file.Filename) == "" {
			result.InvalidFiles = append(result.InvalidFiles, "无效文件名")
			continue
		}

		result.ValidFiles = append(result.ValidFiles, file)
	}

	return result
}

func validateStorageAndUploadLimits(key *models.APIKey, files []*multipart.FileHeader) error {
	var totalSize int64
	for _, file := range files {
		totalSize += file.Size
	}

	if key.StorageLimit > 0 && key.StorageUsed+totalSize > key.StorageLimit {
		return errors.New(errors.CodeStorageLimitExceeded, "API密钥存储容量已用尽")
	}

	if key.UploadCountLimit > 0 && key.UploadCountUsed+len(files) > key.UploadCountLimit {
		return errors.New(errors.CodeUploadLimitExceeded, "API密钥上传次数不足")
	}

	return nil
}

func uploadValidFiles(c *gin.Context, key *models.APIKey, folderID, accessLevel string, optimize bool, validationResult *FileValidationResult) (*APIKeyUploadResult, error) {
	result := &APIKeyUploadResult{
		OversizedFiles:   validationResult.OversizedFiles,
		UnsupportedFiles: validationResult.UnsupportedFiles,
		InvalidFiles:     validationResult.InvalidFiles,
	}

	var responses []*ExternalAPIFileResponse
	var uploadErrors []string

	for _, file := range validationResult.ValidFiles {
		imgInfo, err := UploadFileForAPI(c, key.UserID, file, folderID, accessLevel, optimize)
		if err != nil {
			uploadErrors = append(uploadErrors, fmt.Sprintf("%s: %s", file.Filename, err.Error()))
			continue
		}

		if err := associateFileWithAPIKey(imgInfo.ID, key.ID); err != nil {
			logger.Error("更新文件API密钥关联失败", "fileID", imgInfo.ID, "error", err)
		}

		responses = append(responses, imgInfo)
		go updateAPIKeyUsageAsync(key.ID, file.Size)
	}

	result.Uploaded = responses
	if len(validationResult.OversizedFiles) > 0 {
		result.SizeLimit = fmt.Sprintf("%.1fMB", float64(key.SingleFileLimit)/1024/1024)
	}
	if len(uploadErrors) > 0 {
		result.UploadErrors = uploadErrors
	}

	result.Message = determineUploadMessage(len(responses), len(validationResult.ValidFiles))
	return result, nil
}

func associateFileWithAPIKey(fileID, apiKeyID string) error {
	return database.DB.Model(&models.File{}).Where("id = ?", fileID).Update("api_key_id", apiKeyID).Error
}

func updateAPIKeyUsageAsync(apiKeyID string, fileSize int64) {
	if err := apikey.UpdateAPIKeyUsage(apiKeyID, fileSize); err != nil {
		logger.Error("更新API密钥使用情况失败", "apiKeyID", apiKeyID, "error", err)
	}
}

func processSingleFileUpload(c *gin.Context, key *models.APIKey, folderID, accessLevel string, optimize bool, file *multipart.FileHeader) (*APIKeyUploadResult, error) {
	result := &APIKeyUploadResult{
		IsSingleUpload: true,
	}

	if err := validateSingleFileLimits(key, file); err != nil {
		return result, err
	}

	imgInfo, err := UploadFileForAPI(c, key.UserID, file, folderID, accessLevel, optimize)
	if err != nil {
		result.UploadErrors = []string{fmt.Sprintf("%s: %s", file.Filename, err.Error())}
		result.Message = "上传失败"
		return result, err
	}

	if err := associateFileWithAPIKey(imgInfo.ID, key.ID); err != nil {
		logger.Error("更新文件API密钥关联失败", "fileID", imgInfo.ID, "error", err)
	}

	go updateAPIKeyUsageAsync(key.ID, file.Size)

	result.UploadedSingle = imgInfo
	result.Message = "上传成功"
	return result, nil
}

func validateSingleFileLimits(key *models.APIKey, file *multipart.FileHeader) error {
	if key.SingleFileLimit > 0 && file.Size > key.SingleFileLimit {
		return errors.New(errors.CodeFileTooLarge, fmt.Sprintf("文件大小超过API密钥限制(%.1fMB)", float64(key.SingleFileLimit)/1024/1024))
	}

	if key.StorageLimit > 0 && key.StorageUsed+file.Size > key.StorageLimit {
		return errors.New(errors.CodeStorageLimitExceeded, "API密钥存储容量已用尽")
	}

	if key.UploadCountLimit > 0 && key.UploadCountUsed >= key.UploadCountLimit {
		return errors.New(errors.CodeUploadLimitExceeded, "API密钥上传次数已用尽")
	}

	return nil
}

/* UploadFileForAPI API专用的文件上传，返回简化响应 */
func UploadFileForAPI(c *gin.Context, userID uint, file *multipart.FileHeader, folderID, accessLevel string, optimize bool) (*ExternalAPIFileResponse, error) {
	available, err := stats.CheckUserStorageAvailable(userID, file.Size)
	if err != nil {
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

	ctx := CreateUploadContext(c, userID, file, folderID, accessLevel, optimize)

	if err := validateUploadRequest(ctx); err != nil {
		return nil, err
	}

	if err := processFileAndUpload(ctx); err != nil {
		return nil, err
	}

	if err := saveFileRecordAndStats(ctx); err != nil {
		return nil, err
	}

	return createExternalAPIResponse(ctx), nil
}
