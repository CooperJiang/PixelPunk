package file

import (
	"context"
	"fmt"
	"mime/multipart"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/activity"
	"pixelpunk/internal/services/ai"
	messageService "pixelpunk/internal/services/message"
	"pixelpunk/internal/services/stats"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	pathutil "pixelpunk/pkg/storage/path"
	"pixelpunk/pkg/utils"
	"pixelpunk/pkg/vector"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/* UploadFile 上传单张文件 */
func UploadFile(c *gin.Context, userID uint, file *multipart.FileHeader, folderID, accessLevel string, optimize bool) (*FileDetailResponse, error) {
	return UploadFileWithDuration(c, userID, file, folderID, accessLevel, optimize, "")
}

/* UploadFileWithDuration 上传单张文件（支持存储时长） */
func UploadFileWithDuration(c *gin.Context, userID uint, file *multipart.FileHeader, folderID, accessLevel string, optimize bool, storageDuration string) (*FileDetailResponse, error) {
	available, err := stats.CheckUserStorageAvailable(userID, file.Size)
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

	ctx := CreateUploadContextWithDuration(c, userID, file, folderID, accessLevel, optimize, storageDuration)

	if err := validateUploadRequest(ctx); err != nil {
		return nil, err
	}

	if err := processFileAndUpload(ctx); err != nil {
		return nil, err
	}

	return completeFileUpload(ctx)
}

func completeFileUpload(ctx *UploadContext) (*FileDetailResponse, error) {
	if err := saveFileRecordAndStats(ctx); err != nil {
		return nil, err
	}

	response := createFileResponse(ctx)

	if response != nil {
		activity.LogImageUploadByID(ctx.FileID, ctx.FolderID)
	}

	return response, nil
}

func validateUploadRequest(ctx *UploadContext) error {
	if err := validateUploadInput(ctx); err != nil {
		return err
	}
	return prepareUploadEnvironment(ctx)
}

func processFileAndUpload(ctx *UploadContext) error {
	if err := processFile(ctx); err != nil {
		return err
	}
	return executeUpload(ctx)
}

func saveFileRecordAndStats(ctx *UploadContext) error {
	if err := saveFileData(ctx); err != nil {
		return err
	}
	updateStatisticsAsync(ctx)
	return nil
}

func executeUpload(ctx *UploadContext) error {
	if !ctx.ReuseExistingFile {
		return uploadNewFile(ctx)
	}
	return reuseExistingFile(ctx)
}

func uploadNewFile(ctx *UploadContext) error {
	storageService, err := GetStorageServiceInstance()
	if err != nil {
		logger.Error("获取存储服务失败: %v", err)
		return errors.Wrap(err, errors.CodeInternal, "存储服务初始化失败")
	}

	uploadReq := convertToNewStorageRequest(ctx)

	uploadResult, err := storageService.Upload(context.Background(), uploadReq)
	if err != nil {
		logger.Error("新存储服务上传失败: %v", err)
		return errors.Wrap(err, errors.CodeFileUploadFailed, "上传文件失败")
	}

	ctx.Result = convertFromNewStorageResult(uploadResult)

	prevHash := ctx.FileHash
	if uploadResult.Hash != "" && len(uploadResult.Hash) == 32 {
		ctx.FileHash = uploadResult.Hash
	} else {
		_ = prevHash // 无有效哈希则保留 reader MD5
	}
	ctx.FileSize = uploadResult.Size
	ctx.FileFormat = uploadResult.Format
	ctx.ActualChannelID = uploadResult.ChannelID

	return handleAccessLevel(ctx)
}

func reuseExistingFile(ctx *UploadContext) error {
	existingImage := ctx.ExistingFile
	ctx.Result = &UploadResult{
		URL:                       existingImage.URL,                       // 访问文件的URL
		LocalUrlPath:              existingImage.LocalFilePath,             // 本地存储路径
		ThumbUrl:                  existingImage.ThumbURL,                  // 缩略图URL
		LocalThumbPath:            existingImage.LocalThumbPath,            // 本地缩略图路径
		RemoteUrl:                 existingImage.RemoteURL,                 // 远程URL
		RemoteThumbUrl:            existingImage.RemoteThumbURL,            // 远程缩略图URL
		Width:                     existingImage.Width,                     // 文件宽度
		Height:                    existingImage.Height,                    // 文件高度
		ThumbnailGenerationFailed: existingImage.ThumbnailGenerationFailed, // 复用原文件的缩略图失败状态
		ThumbnailFailureReason:    existingImage.ThumbnailFailureReason,    // 复用原文件的缩略图失败原因
	}
	return handleAccessLevel(ctx)
}

func saveFileData(ctx *UploadContext) error {
	file := createFileModel(ctx)

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		ctx.Tx = tx

		if err := saveFileRecord(tx, file); err != nil {
			return err
		}

		if err := updateUserStats(tx, ctx); err != nil {
			return err
		}

		if ctx.EXIFData != nil {
			ctx.EXIFData.FileID = file.ID
			if err := tx.Create(ctx.EXIFData).Error; err != nil {
				logger.Warn("保存 EXIF 数据失败: %v", err)
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	ctx.SavedFile = file
	ctx.FileModel = file

	// 检查缩略图生成是否失败，如果失败则发送通知
	if file.ThumbnailGenerationFailed {
		userID := ctx.UserID
		go func() {
			msgService := messageService.GetMessageService()
			variables := map[string]interface{}{
				"file_id":      file.ID,
				"file_name":    file.OriginalName,
				"reason":       file.ThumbnailFailureReason,
				"related_type": "file",
				"related_id":   file.ID,
			}
			if err := msgService.SendTemplateMessage(userID, common.MessageTypeFileThumbnailFailed, variables); err != nil {
				logger.Warn("发送缩略图生成失败通知失败: userID=%d, fileID=%s, error=%v", userID, file.ID, err)
			}
		}()
	}

	if ctx.OriginalFileID != "" {
		go func(origID, newID string) {
			defer func() {
				if r := recover(); r != nil {
					logger.Error("重复文件向量复用 panic: %v", r)
				}
			}()

			asyncCtx := &UploadContext{
				OriginalFileID: origID,
				FileID:         newID,
			}

			if err := reuseAnalysisAndVectorForDuplicate(asyncCtx); err != nil {
				logger.Warn("重复文件复用AI/向量失败: %v", err)
			}
		}(ctx.OriginalFileID, ctx.FileID)
	}

	// 异步执行所有后处理操作，避免阻塞上传接口返回
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("[上传后处理] panic: %v, 文件ID: %s", r, file.ID)
			}
		}()

		if utils.GetAiAnalysisEnabled() {
			if err := captureThumbnailBase64(ctx); err != nil {
				logger.Warn("[上传后处理] 捕获缩略图base64数据失败: %v, file_id=%s", err, file.ID)
			}

			if err := ai.AddFileToQueue(*file); err != nil {
				logger.Error("[上传后处理] 将文件加入AI处理队列失败，文件ID: %s, 错误: %v", file.ID, err)
			}
		}

		if vector.IsVectorEnabled() && file.Description != "" {
			vector.AddFileToVectorQueue(*file)
		}

	}()

	return nil
}

func reuseAnalysisAndVectorForDuplicate(ctx *UploadContext) error {
	db := database.DB
	newID := ctx.FileID
	origID := ctx.OriginalFileID

	var origAI models.FileAIInfo
	if err := db.Where("file_id = ?", origID).First(&origAI).Error; err == nil {
		clone := origAI
		clone.ID = 0
		clone.FileID = newID
		_ = db.Where("file_id = ?", newID).Delete(&models.FileAIInfo{})
		if err := db.Create(&clone).Error; err != nil {
			logger.Warn("复制AI信息失败: %v", err)
		}
		if clone.Description != "" {
			_ = db.Model(&models.File{}).Where("id = ?", newID).Update("description", clone.Description).Error
		}
	}

	var rels []models.FileGlobalTagRelation
	if err := db.Where("file_id = ?", origID).Find(&rels).Error; err == nil {
		for _, r := range rels {
			nr := models.FileGlobalTagRelation{
				FileID:      newID,
				TagID:       r.TagID,
				UserID:      r.UserID,
				AccessLevel: r.AccessLevel,
				Source:      r.Source,
				Confidence:  r.Confidence,
			}
			_ = db.Create(&nr).Error
		}
	}

	var origImg models.File
	if err := db.Where("id = ?", origID).First(&origImg).Error; err == nil {
		updates := map[string]interface{}{}
		if origImg.CategoryID != nil {
			updates["category_id"] = *origImg.CategoryID
			updates["category_source"] = origImg.CategorySource
		}
		if len(updates) > 0 {
			_ = db.Model(&models.File{}).Where("id = ?", newID).Updates(updates).Error
		}
	}

	if vector.IsVectorEnabled() {
		desc := ""
		var newAI models.FileAIInfo
		if err := db.Where("file_id = ?", newID).First(&newAI).Error; err == nil {
			if newAI.SearchContent != "" {
				desc = newAI.SearchContent
			} else {
				desc = newAI.Description
			}
		}
		if eng := vector.GetGlobalVectorEngine(); eng != nil && eng.IsEnabled() {
			if err := eng.CloneVectorFrom(origID, newID, desc); err != nil {
				msg := err.Error()
				if strings.Contains(msg, "向量不存在") || strings.Contains(msg, "Not found") || strings.Contains(msg, "获取向量失败") {
				} else {
					logger.Warn("复制向量失败（将由后续传播兜底）: %v", err)
				}
			}
		}
	}

	return nil
}

func getDescriptionFromContext(ctx *UploadContext) string {
	if ctx.Context == nil {
		return ""
	}

	return ctx.Context.PostForm("description")
}

/* BatchUploadFailure 批量上传失败信息 */
type BatchUploadFailure struct {
	Filename string
	Error    string
	Index    int
}

/* BatchUploadResult 批量上传结果 */
type BatchUploadResult struct {
	TotalFiles   int
	SuccessCount int
	FailureCount int
	SuccessFiles []*FileDetailResponse
	Failures     []BatchUploadFailure
	Message      string
}

/* UploadFileBatch 批量上传文件（支持部分成功） */
func UploadFileBatch(c *gin.Context, userID uint, files []*multipart.FileHeader, folderID, accessLevel string, optimize bool) (*BatchUploadResult, error) {
	if err := validateBatchUploadFiles(files); err != nil {
		return nil, err
	}

	fileCount := len(files)
	if exceeded, err := checkDailyUploadLimit(userID, fileCount); err != nil {
		logger.Warn("检查每日上传限制失败: %v", err)
	} else if exceeded {
		return nil, errors.New(errors.CodeUploadLimitExceeded, "上传数量超过每日限制")
	}

	if folderID == "null" {
		folderID = ""
	}

	result := &BatchUploadResult{
		TotalFiles:   len(files),
		SuccessFiles: make([]*FileDetailResponse, 0),
		Failures:     make([]BatchUploadFailure, 0),
	}

	for index, file := range files {
		resp, err := UploadFile(c, userID, file, folderID, accessLevel, optimize)
		if err != nil {
			result.Failures = append(result.Failures, BatchUploadFailure{
				Filename: file.Filename,
				Error:    err.Error(),
				Index:    index,
			})
			result.FailureCount++
		} else {
			result.SuccessFiles = append(result.SuccessFiles, resp)
			result.SuccessCount++
		}
	}

	if result.FailureCount == 0 {
		result.Message = fmt.Sprintf("全部%d个文件上传成功", result.SuccessCount)
	} else if result.SuccessCount == 0 {
		result.Message = fmt.Sprintf("全部%d个文件上传失败", result.FailureCount)
	} else {
		result.Message = fmt.Sprintf("成功上传%d个文件，%d个文件失败", result.SuccessCount, result.FailureCount)
	}

	return result, nil
}

/* UploadFileBatchWithDuration 批量上传文件（支持存储时长） */
func UploadFileBatchWithDuration(c *gin.Context, userID uint, files []*multipart.FileHeader, folderID, accessLevel string, optimize bool, storageDuration string) (*BatchUploadResult, error) {
	if err := validateBatchUploadFiles(files); err != nil {
		return nil, err
	}

	fileCount := len(files)
	if exceeded, err := checkDailyUploadLimit(userID, fileCount); err != nil {
		logger.Warn("检查每日上传限制失败: %v", err)
	} else if exceeded {
		return nil, errors.New(errors.CodeUploadLimitExceeded, "上传数量超过每日限制")
	}

	if folderID == "null" {
		folderID = ""
	}

	result := &BatchUploadResult{
		TotalFiles:   len(files),
		SuccessCount: 0,
		FailureCount: 0,
		SuccessFiles: make([]*FileDetailResponse, 0),
		Failures:     make([]BatchUploadFailure, 0),
	}

	validationResult := validateFiles(nil, files)
	if len(validationResult.ValidFiles) == 0 {
		for _, filename := range validationResult.OversizedFiles {
			result.Failures = append(result.Failures, BatchUploadFailure{Filename: filename, Error: "文件大小超过限制"})
		}
		for _, filename := range validationResult.UnsupportedFiles {
			result.Failures = append(result.Failures, BatchUploadFailure{Filename: filename, Error: "文件格式不支持"})
		}
		for _, filename := range validationResult.InvalidFiles {
			result.Failures = append(result.Failures, BatchUploadFailure{Filename: filename, Error: "文件无效"})
		}
		result.FailureCount = len(result.Failures)
		result.Message = "没有有效的文件可上传"
		return result, nil
	}

	return uploadValidFilesWithDuration(c, userID, folderID, accessLevel, optimize, storageDuration, validationResult)
}

func uploadValidFilesWithDuration(c *gin.Context, userID uint, folderID, accessLevel string, optimize bool, storageDuration string, validationResult *FileValidationResult) (*BatchUploadResult, error) {
	result := &BatchUploadResult{
		TotalFiles:   len(validationResult.ValidFiles) + len(validationResult.OversizedFiles) + len(validationResult.UnsupportedFiles) + len(validationResult.InvalidFiles),
		SuccessCount: 0,
		FailureCount: 0,
		SuccessFiles: make([]*FileDetailResponse, 0),
		Failures:     make([]BatchUploadFailure, 0),
	}

	responses := make([]*FileDetailResponse, 0, len(validationResult.ValidFiles))

	for _, file := range validationResult.ValidFiles {
		imgInfo, err := UploadFileWithDuration(c, userID, file, folderID, accessLevel, optimize, storageDuration)
		if err != nil {
			result.Failures = append(result.Failures, BatchUploadFailure{
				Filename: file.Filename,
				Error:    err.Error(),
			})
			continue
		}
		responses = append(responses, imgInfo)
	}

	result.SuccessFiles = responses
	result.SuccessCount = len(responses)

	result.FailureCount = len(result.Failures)

	for _, filename := range validationResult.OversizedFiles {
		result.Failures = append(result.Failures, BatchUploadFailure{Filename: filename, Error: "文件大小超过限制"})
	}
	for _, filename := range validationResult.UnsupportedFiles {
		result.Failures = append(result.Failures, BatchUploadFailure{Filename: filename, Error: "文件格式不支持"})
	}
	for _, filename := range validationResult.InvalidFiles {
		result.Failures = append(result.Failures, BatchUploadFailure{Filename: filename, Error: "文件无效"})
	}

	result.Message = determineUploadMessage(len(responses), len(validationResult.ValidFiles))
	return result, nil
}

/* UploadImageWithAPIKey 逻辑已迁移至 upload_apikey.go */

func determineUploadMessage(successCount, totalCount int) string {
	if successCount == totalCount {
		return "全部文件上传成功"
	} else if successCount > 0 {
		return "部分文件上传成功"
	}
	return "所有文件上传失败"
}

/* CheckDuplicateResponse / CheckDuplicateFileInfo 已迁移至 upload_instant.go */

/* InstantUploadResponse 已迁移至 upload_instant.go */

/* InstantUpload 已迁移至 upload_instant.go */

/* InstantUploadWithDuration 已迁移至 upload_instant.go */

/* GuestUpload 已迁移至 upload_guest.go */

/* UploadImageWithWatermark 上传单张文件（支持水印） */
/* UploadImageWithWatermark 已迁移至 upload_watermark.go */

/* UploadFileBatchWithWatermark 批量上传文件（支持水印） */
/* UploadFileBatchWithWatermark 已迁移至 upload_watermark.go */

/* GuestUploadWithWatermark 已迁移至 upload_guest.go */

func ensureFullObjectPath(ctx *UploadContext, p string) string {
	return pathutil.EnsureObjectKey(ctx.UserID, p, false)
}
