package file

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pixelpunk/internal/controllers/file/dto"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/setting"
	"pixelpunk/pkg/config"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChunkedUploadService struct {
	tempStoragePath string
}

func NewChunkedUploadService() *ChunkedUploadService {
	tempPath := filepath.Join("temp", "chunks")
	return &ChunkedUploadService{
		tempStoragePath: tempPath,
	}
}

/* InitChunkedUpload 初始化分片上传 */
func InitChunkedUpload(userID uint, req *dto.InitChunkedUploadDTO) (*models.UploadSession, error) {
	chunkedUploadEnabled, err := setting.GetBoolValue("chunked_upload_enabled", false)
	if err != nil {
		logger.Error("获取分片上传开关配置失败: %v", err)
		return nil, errors.New(errors.CodeInternal, "获取分片上传配置失败")
	}

	if !chunkedUploadEnabled {
		return nil, errors.New(errors.CodeInvalidParameter, "分片上传功能已禁用")
	}

	if !isValidFileTypeForChunked(req.MimeType) {
		return nil, errors.New(errors.CodeFileTypeNotSupported, "不支持的文件格式")
	}

	const (
		minChunkSize = 1024 * 1024
		maxChunkSize = 10 * 1024 * 1024
	)

	if req.ChunkSize < minChunkSize {
		return nil, errors.New(errors.CodeInvalidParameter, fmt.Sprintf("分片大小不能小于 %d 字节", minChunkSize))
	}
	if req.ChunkSize > maxChunkSize {
		return nil, errors.New(errors.CodeInvalidParameter, fmt.Sprintf("分片大小不能大于 %d 字节", maxChunkSize))
	}

	maxFileSizeMB, err := setting.GetNumberValue("max_file_size", 100.0)
	if err != nil {
		logger.Error("获取最大文件大小配置失败: %v", err)
		maxFileSizeMB = 100.0
	}
	maxFileSize := int64(maxFileSizeMB * 1024 * 1024)

	if req.FileSize > maxFileSize {
		return nil, errors.New(errors.CodeInvalidParameter, fmt.Sprintf("文件大小不能超过 %d 字节", maxFileSize))
	}

	totalChunks := int(math.Ceil(float64(req.FileSize) / float64(req.ChunkSize)))
	sessionID := strings.ReplaceAll(uuid.New().String(), "-", "")

	sessionTimeoutHours, err := setting.GetNumberValue("session_timeout", 24.0)
	if err != nil {
		logger.Error("获取会话超时时间配置失败: %v", err)
		sessionTimeoutHours = 24.0
	}

	var watermarkConfigJSON string
	if req.WatermarkConfig != nil {
		watermarkBytes, err := json.Marshal(req.WatermarkConfig)
		if err != nil {
			logger.Error("序列化水印配置失败: %v", err)
		} else {
			watermarkConfigJSON = string(watermarkBytes)
		}
	}

	session := &models.UploadSession{
		SessionID:       sessionID,
		UserID:          userID,
		FileName:        req.FileName,
		FileSize:        req.FileSize,
		FileMD5:         req.FileMD5,
		MimeType:        req.MimeType,
		ChunkSize:       req.ChunkSize,
		TotalChunks:     totalChunks,
		Status:          "pending",
		Progress:        0,
		FolderID:        req.FolderID,
		AccessLevel:     req.AccessLevel,
		Optimize:        req.Optimize,
		WatermarkConfig: watermarkConfigJSON,
		ExpiresAt:       time.Now().Add(time.Duration(sessionTimeoutHours) * time.Hour),
	}

	if err := database.DB.Create(session).Error; err != nil {
		logger.Error("创建分片上传会话失败: %v", err)
		return nil, errors.Wrap(err, errors.CodeInternal, "创建上传会话失败")
	}

	sessionDir := filepath.Join("temp", "chunks", sessionID)
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		logger.Error("创建临时存储目录失败: %v", err)
		return nil, errors.Wrap(err, errors.CodeInternal, "创建临时目录失败")
	}
	for i := 0; i < totalChunks; i++ {
		chunk := &models.UploadChunk{
			SessionID:   sessionID,
			ChunkNumber: i,
			ChunkSize:   req.ChunkSize,
			Status:      "pending",
			StoragePath: filepath.Join(sessionDir, fmt.Sprintf("chunk_%d", i)),
		}

		if i == totalChunks-1 {
			chunk.ChunkSize = req.FileSize - int64(i)*req.ChunkSize
		}

		if err := database.DB.Create(chunk).Error; err != nil {
			logger.Error("创建分片记录失败: %v", err)
			return nil, errors.Wrap(err, errors.CodeInternal, "创建分片记录失败")
		}
	}

	return session, nil
}

/* UploadChunk 上传分片 */
func UploadChunk(req *dto.UploadChunkDTO, file *multipart.FileHeader) error {
	var session models.UploadSession
	if err := database.DB.Where("session_id = ?", req.SessionID).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.CodeNotFound, "上传会话不存在")
		}
		return errors.Wrap(err, errors.CodeInternal, "查询上传会话失败")
	}

	if session.IsExpired() {
		return errors.New(errors.CodeUploadSessionExpired, "上传会话已过期")
	}

	if !session.IsActive() {
		return errors.New(errors.CodeInvalidParameter, "上传会话状态无效")
	}

	var chunk models.UploadChunk
	if err := database.DB.Where("session_id = ? AND chunk_number = ?", req.SessionID, req.ChunkNumber).First(&chunk).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.CodeNotFound, "分片记录不存在")
		}
		return errors.Wrap(err, errors.CodeInternal, "查询分片记录失败")
	}

	if chunk.IsUploaded() {
		if err := updateSessionProgress(req.SessionID); err != nil {
			logger.Error("更新会话进度失败: %v", err)
		}
		return nil
	}

	src, err := file.Open()
	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "打开上传文件失败")
	}
	defer src.Close()

	if file.Size != chunk.ChunkSize {
		return errors.New(errors.CodeInvalidParameter, "分片大小不匹配")
	}

	hasher := md5.New()
	src.Seek(0, 0)
	if _, err := io.Copy(hasher, src); err != nil {
		return errors.Wrap(err, errors.CodeInternal, "计算分片MD5失败")
	}
	calculatedMD5 := fmt.Sprintf("%x", hasher.Sum(nil))

	if calculatedMD5 != req.ChunkMD5 {
		return errors.New(errors.CodeInvalidParameter, "分片MD5校验失败")
	}

	src.Seek(0, 0)
	dst, err := os.Create(chunk.StoragePath)
	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "创建分片文件失败")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return errors.Wrap(err, errors.CodeInternal, "保存分片文件失败")
	}

	chunk.Status = "uploaded"
	chunk.ChunkMD5 = req.ChunkMD5
	if err := database.DB.Save(&chunk).Error; err != nil {
		return errors.Wrap(err, errors.CodeInternal, "更新分片状态失败")
	}

	if err := updateSessionProgress(req.SessionID); err != nil {
		logger.Error("更新会话进度失败: %v", err)
	}

	return nil
}

/* CompleteChunkedUpload 完成分片上传 */
func CompleteChunkedUpload(sessionID string) (*FileDetailResponse, error) {
	var session models.UploadSession
	if err := database.DB.Where("session_id = ?", sessionID).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.CodeNotFound, "上传会话不存在")
		}
		return nil, errors.Wrap(err, errors.CodeInternal, "查询上传会话失败")
	}

	if session.IsCompleted() {
		if session.FileID != "" {
			var file models.File
			if err := database.DB.Where("id = ?", session.FileID).First(&file).Error; err == nil {
				resp := BuildFileDetailResponse(file, 0, nil)
				return &resp, nil
			}
		}
		return nil, errors.New(errors.CodeInvalidParameter, "上传会话已完成")
	}

	var uploadedCount int64
	if err := database.DB.Model(&models.UploadChunk{}).
		Where("session_id = ? AND status = ?", sessionID, "uploaded").
		Count(&uploadedCount).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "查询已上传分片数量失败")
	}

	if int(uploadedCount) != session.TotalChunks {
		return nil, errors.New(errors.CodeInvalidParameter, "分片上传未完成")
	}

	mergedFilePath, err := mergeChunks(sessionID, session.FileName)
	if err != nil {
		return nil, err
	}

	if err := validateMergedFile(mergedFilePath, session.FileMD5); err != nil {
		return nil, err
	}

	imageResponse, err := processUploadedFile(session.UserID, mergedFilePath, &session)
	if err != nil {
		return nil, err
	}

	session.Status = "completed"
	session.Progress = 100
	session.FileID = imageResponse.ID
	if err := database.DB.Save(&session).Error; err != nil {
		logger.Error("更新会话状态失败: %v", err)
	}

	go cleanupTempFiles(sessionID)

	var file models.File
	if err := database.DB.Where("id = ?", imageResponse.ID).First(&file).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "查询文件失败")
	}
	resp := BuildFileDetailResponse(file, 0, nil)
	return &resp, nil
}

/* GetChunkedUploadStatus 获取分片上传状态 */
func GetChunkedUploadStatus(sessionID string) (*dto.ChunkedUploadResponse, error) {
	var session models.UploadSession
	if err := database.DB.Where("session_id = ?", sessionID).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.CodeNotFound, "上传会话不存在")
		}
		return nil, errors.Wrap(err, errors.CodeInternal, "查询上传会话失败")
	}

	var uploadedCount int64
	if err := database.DB.Model(&models.UploadChunk{}).
		Where("session_id = ? AND status = ?", sessionID, "uploaded").
		Count(&uploadedCount).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "查询已上传分片数量失败")
	}

	response := &dto.ChunkedUploadResponse{
		SessionID:      session.SessionID,
		Status:         session.Status,
		Progress:       session.Progress,
		TotalChunks:    session.TotalChunks,
		UploadedChunks: int(uploadedCount),
		Message:        fmt.Sprintf("已上传 %d/%d 个分片", uploadedCount, session.TotalChunks),
	}

	return response, nil
}

/* CancelChunkedUpload 取消分片上传 */
func CancelChunkedUpload(sessionID string) error {
	var session models.UploadSession
	if err := database.DB.Where("session_id = ?", sessionID).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.CodeNotFound, "上传会话不存在")
		}
		return errors.Wrap(err, errors.CodeInternal, "查询上传会话失败")
	}

	session.Status = "failed"
	if err := database.DB.Save(&session).Error; err != nil {
		return errors.Wrap(err, errors.CodeInternal, "更新会话状态失败")
	}

	go cleanupTempFiles(sessionID)

	return nil
}

func isValidFileTypeForChunked(mimeType string) bool {
	allowedFormatsInterface, err := setting.GetSettingValue("allowed_file_formats", []interface{}{})
	if err != nil {
		logger.Error("获取允许文件格式配置失败: %v", err)
		uploadConfig := config.GetUploadConfig()
		for _, allowedType := range uploadConfig.AllowedTypes {
			if allowedType == mimeType {
				return true
			}
		}
		return false
	}

	if allowedFormatsArray, ok := allowedFormatsInterface.([]interface{}); ok {
		for _, formatInterface := range allowedFormatsArray {
			if format, ok := formatInterface.(string); ok && format == mimeType {
				return true
			}
		}
	}

	uploadConfig := config.GetUploadConfig()
	for _, allowedType := range uploadConfig.AllowedTypes {
		if allowedType == mimeType {
			return true
		}
	}
	return false
}

func updateSessionProgress(sessionID string) error {
	var uploadedCount int64
	if err := database.DB.Model(&models.UploadChunk{}).
		Where("session_id = ? AND status = ?", sessionID, "uploaded").
		Count(&uploadedCount).Error; err != nil {
		return err
	}

	var session models.UploadSession
	if err := database.DB.Where("session_id = ?", sessionID).First(&session).Error; err != nil {
		return err
	}

	progress := int(float64(uploadedCount) / float64(session.TotalChunks) * 100)

	updates := map[string]interface{}{
		"progress": progress,
	}

	if uploadedCount > 0 && session.Status == "pending" {
		updates["status"] = "uploading"
	}

	return database.DB.Model(&session).Where("session_id = ?", sessionID).Updates(updates).Error
}

func mergeChunks(sessionID, fileName string) (string, error) {
	tempDir := filepath.Join("temp", "chunks", sessionID)
	mergedFilePath := filepath.Join(tempDir, "merged_"+fileName)

	mergedFile, err := os.Create(mergedFilePath)
	if err != nil {
		return "", errors.Wrap(err, errors.CodeInternal, "创建合并文件失败")
	}
	defer mergedFile.Close()

	var chunks []models.UploadChunk
	if err := database.DB.Where("session_id = ?", sessionID).
		Order("chunk_number ASC").Find(&chunks).Error; err != nil {
		return "", errors.Wrap(err, errors.CodeInternal, "查询分片记录失败")
	}

	for _, chunk := range chunks {
		chunkFile, err := os.Open(chunk.StoragePath)
		if err != nil {
			return "", errors.Wrap(err, errors.CodeInternal, fmt.Sprintf("打开分片文件失败: %d", chunk.ChunkNumber))
		}

		if _, err := io.Copy(mergedFile, chunkFile); err != nil {
			chunkFile.Close()
			return "", errors.Wrap(err, errors.CodeInternal, fmt.Sprintf("写入分片数据失败: %d", chunk.ChunkNumber))
		}

		chunkFile.Close()
	}

	return mergedFilePath, nil
}

func validateMergedFile(filePath, expectedMD5 string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "打开合并文件失败")
	}
	defer file.Close()

	hasher := md5.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return errors.Wrap(err, errors.CodeInternal, "计算文件MD5失败")
	}

	calculatedMD5 := fmt.Sprintf("%x", hasher.Sum(nil))
	if calculatedMD5 != expectedMD5 {
		return errors.New(errors.CodeInvalidParameter, "文件MD5校验失败")
	}

	return nil
}

func processUploadedFile(userID uint, filePath string, session *models.UploadSession) (*FileDetailResponse, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "打开合并文件失败")
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "获取文件信息失败")
	}

	header := &multipart.FileHeader{
		Filename: session.FileName,
		Size:     fileInfo.Size(),
		Header:   make(map[string][]string),
	}
	header.Header.Set("Content-Type", session.MimeType)

	ctx := CreateUploadContext(nil, userID, header, session.FolderID, session.AccessLevel, session.Optimize)

	// 设置水印配置（如果有）
	if session.WatermarkConfig != "" {
		ctx.WatermarkEnabled = true
		ctx.WatermarkConfig = session.WatermarkConfig
	}

	ctx.FileExt = filepath.Ext(session.FileName)
	ctx.FileHash = session.FileMD5
	ctx.FileID = utils.GenerateFileID()

	if err := validateFolder(ctx); err != nil {
		return nil, err
	}

	if err := processFolderPath(ctx); err != nil {
		return nil, err
	}

	if err := prepareUploadEnvironment(ctx); err != nil {
		return nil, err
	}

	ctx.FileExt = filepath.Ext(session.FileName)
	ctx.FileHash = session.FileMD5

	if err := uploadMergedFileDirectly(ctx, filePath); err != nil {
		return nil, err
	}

	return completeFileUpload(ctx)
}

func uploadMergedFileDirectly(ctx *UploadContext, mergedFilePath string) error {
	file, err := os.Open(mergedFilePath)
	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "打开合并文件失败")
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "获取文件信息失败")
	}

	fileName := filepath.Base(mergedFilePath)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "创建表单文件失败")
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "复制文件内容失败")
	}

	err = writer.Close()
	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "关闭multipart writer失败")
	}

	reader := multipart.NewReader(&body, writer.Boundary())
	form, err := reader.ReadForm(fileInfo.Size() + 1024)
	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "解析multipart form失败")
	}
	defer form.RemoveAll()

	if files := form.File["file"]; len(files) > 0 {
		originalFile := ctx.File
		ctx.File = files[0]

		// 如果启用了水印，调用带水印的上传流程
		if ctx.WatermarkEnabled && ctx.WatermarkConfig != "" {
			err = processFileAndUploadWithWatermark(ctx)
		} else {
			err = uploadNewFile(ctx)
		}

		ctx.File = originalFile

		return err
	}

	return errors.New(errors.CodeInternal, "multipart form中没有找到文件")
}

func cleanupTempFiles(sessionID string) {
	tempDir := filepath.Join("temp", "chunks", sessionID)
	if err := os.RemoveAll(tempDir); err != nil {
		logger.Error("清理临时文件失败: session_id=%s, error=%v", sessionID, err)
	}
}
