package file

import (
	"fmt"
	"net/http"
	"strings"

	"pixelpunk/internal/controllers/file/dto"
	"pixelpunk/internal/middleware"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/apikey"
	filesvc "pixelpunk/internal/services/file"
	folder "pixelpunk/internal/services/folder"
	"pixelpunk/internal/services/setting"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"

	"mime/multipart"

	"github.com/gin-gonic/gin"
)

// Upload 上传单张文件
func Upload(c *gin.Context) {
	var userID uint
	var apiKeyID string

	// 优先从API密钥中获取用户ID
	if apiKey, exists := c.Get("api_key"); exists {
		if key, ok := apiKey.(*models.APIKey); ok {
			userID = key.UserID
			apiKeyID = key.ID
		}
	}

	// 如果没有API密钥，则从认证用户中获取
	if userID == 0 {
		userID = middleware.GetCurrentUserID(c)
	}

	req, err := common.ValidateRequest[dto.UploadFileDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	// 使用StorageConfig验证存储时长（已登录用户）
	if req.StorageDuration != "" {
		storageConfig, err := setting.CreateStorageConfig()
		if err != nil {
			logger.Error("Upload: 创建存储配置失败: %v", err)
			errors.HandleError(c, errors.Wrap(err, errors.CodeDBQueryFailed, "创建存储配置失败"))
			return
		}
		// 验证存储时长（false表示已登录用户模式）
		if err := storageConfig.ValidateStorageDuration(req.StorageDuration, false); err != nil {
			logger.Error("Upload: 存储时长验证失败: %v", err)
			errors.HandleError(c, errors.New(errors.CodeInvalidParameter, err.Error()))
			return
		}
	}

	file, err := c.FormFile("file")
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件上传失败: "+err.Error()))
		return
	}

	// 如果使用API密钥，且指定了文件夹ID，则使用API密钥的文件夹ID
	if apiKeyID != "" && c.Request.FormValue("folder_id") == "" {
		// 从数据库获取API密钥
		var apiKey models.APIKey
		if err := database.DB.Where("id = ?", apiKeyID).First(&apiKey).Error; err == nil {
			if apiKey.FolderID != "" {
				req.FolderID = apiKey.FolderID
			}
		}
	}

	// 新参数优先：统一的watermark（JSON字符串）
	wmEnabled := false
	wmConfig := ""
	if strings.TrimSpace(req.Watermark) != "" {
		wmEnabled = true
		wmConfig = req.Watermark
	} else if req.WatermarkEnabled && strings.TrimSpace(req.WatermarkConfig) != "" {
		wmEnabled = true
		wmConfig = req.WatermarkConfig
	}

	fileInfo, err := filesvc.UploadFileWithWatermark(c, userID, file, req.FolderID, req.AccessLevel, req.Optimize, req.StorageDuration, wmEnabled, wmConfig)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	// 如果是通过API密钥上传的，记录API密钥信息和更新使用情况
	if apiKeyID != "" {
		if err := database.DB.Model(&models.File{}).Where("id = ?", fileInfo.ID).Update("api_key_id", apiKeyID).Error; err != nil {
			// 只记录错误，不影响上传结果
		}
		// 异步更新API密钥使用情况
		go func(id string, size int64) {
			_ = apikey.UpdateAPIKeyUsage(id, size)
		}(apiKeyID, file.Size)
	}

	errors.ResponseSuccess(c, fileInfo, "上传成功")
}

// BatchUpload 批量上传文件
func BatchUpload(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	req, err := common.ValidateRequest[dto.UploadFileDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件上传失败: "+err.Error()))
		return
	}

	files := form.File["files[]"]
	if len(files) == 0 {
		// 兼容仅 file/file[]/upload[]，不再支持 images*/image*
		if alt := form.File["files"]; len(alt) > 0 {
			files = alt
		}
		if len(files) == 0 {
			if alt := form.File["file[]"]; len(alt) > 0 {
				files = alt
			}
		}
		if len(files) == 0 {
			if alt := form.File["upload[]"]; len(alt) > 0 {
				files = alt
			}
		}
	}
	if len(files) == 0 {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "未检测到上传文件"))
		return
	}

	wmEnabled := false
	wmConfig := ""
	if strings.TrimSpace(req.Watermark) != "" {
		wmEnabled = true
		wmConfig = req.Watermark
	} else if req.WatermarkEnabled && strings.TrimSpace(req.WatermarkConfig) != "" {
		wmEnabled = true
		wmConfig = req.WatermarkConfig
	}

	result, err := filesvc.UploadFileBatchWithWatermark(c, userID, files, req.FolderID, req.AccessLevel, req.Optimize, req.StorageDuration, wmEnabled, wmConfig)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	// 将结果转换为DTO格式
	resultDTO := dto.BatchUploadResultDTO{
		TotalFiles:   len(files),
		SuccessCount: len(result),
		FailureCount: len(files) - len(result),
		SuccessFiles: make([]interface{}, len(result)),
		Failures:     make([]dto.BatchUploadFailure, 0),
		Message:      fmt.Sprintf("批量上传完成: 成功 %d 个，失败 %d 个", len(result), len(files)-len(result)),
	}
	for i, file := range result {
		resultDTO.SuccessFiles[i] = file
	}
	errors.ResponseSuccess(c, resultDTO, resultDTO.Message)
}

// GuestUpload 游客上传文件
func GuestUpload(c *gin.Context) {
	req, err := common.ValidateRequest[dto.GuestUploadDTO](c)
	if err != nil {
		logger.Error("GuestUpload: 请求验证失败: %v", err)
		errors.HandleError(c, err)
		return
	}

	storageConfig, err := setting.CreateStorageConfig()
	if err != nil {
		logger.Error("GuestUpload: 创建存储配置失败: %v", err)
		errors.HandleError(c, errors.Wrap(err, errors.CodeDBQueryFailed, "创建存储配置失败"))
		return
	}
	if err := storageConfig.ValidateStorageDuration(req.StorageDuration, true); err != nil {
		logger.Error("GuestUpload: 存储时长验证失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, err.Error()))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		logger.Error("GuestUpload: 获取文件失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件上传失败: "+err.Error()))
		return
	}

	fingerprint := common.GenerateFingerprint(c.Request)

	guestFolderID, err := folder.GetOrCreateGuestFolder()
	if err != nil {
		logger.Error("GuestUpload: 获取游客文件夹失败: %v", err)
		errors.HandleError(c, err)
		return
	}

	wmEnabled := false
	wmConfig := ""
	if strings.TrimSpace(req.Watermark) != "" {
		wmEnabled = true
		wmConfig = req.Watermark
	} else if req.WatermarkEnabled && strings.TrimSpace(req.WatermarkConfig) != "" {
		wmEnabled = true
		wmConfig = req.WatermarkConfig
	}

	fileInfo, remainingCount, err := filesvc.GuestUploadWithWatermark(c, file, guestFolderID, req.AccessLevel, req.Optimize, req.StorageDuration, fingerprint, wmEnabled, wmConfig)
	if err != nil {
		logger.Error("GuestUpload: 服务层处理失败: %v", err)
		errors.HandleError(c, err)
		return
	}

	response := map[string]interface{}{"file_info": fileInfo, "remaining_count": remainingCount}
	errors.ResponseSuccess(c, response, "上传成功")
}

// UploadForApiKey 使用 API Key 上传
func UploadForApiKey(c *gin.Context) {
	apiKeyObj, _ := c.Get("api_key")
	key := apiKeyObj.(*models.APIKey)

	folderID := c.PostForm("folderId")
	filePath := c.PostForm("filePath")
	accessLevel := c.PostForm("access_level")
	optimizeStr := c.PostForm("optimize")
	optimize := optimizeStr == "true" || optimizeStr == "1"

	form, err := c.MultipartForm()
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "解析表单数据失败: "+err.Error()))
		return
	}

	var files []*multipart.FileHeader
	for _, fieldName := range []string{"files[]", "files", "file[]", "upload[]"} {
		if formFiles := form.File[fieldName]; len(formFiles) > 0 {
			files = formFiles
			break
		}
	}
	var singleFile *multipart.FileHeader
	if len(files) == 0 {
		singleFile, err = c.FormFile("file")
		if err != nil && err != http.ErrMissingFile {
			errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件上传失败: "+err.Error()))
			return
		}
	}

	result, err := filesvc.UploadFileWithAPIKey(c, key, folderID, filePath, accessLevel, optimize, files, singleFile)
	if err != nil {
		if result != nil && (len(result.Uploaded) > 0 || result.UploadedSingle != nil) {
			partialResponse := gin.H{}
			if result.IsSingleUpload && result.UploadedSingle != nil {
				partialResponse["uploaded"] = result.UploadedSingle
			} else if len(result.Uploaded) > 0 {
				partialResponse["uploaded"] = result.Uploaded
			}
			partialResponse["oversized_files"] = result.OversizedFiles
			partialResponse["size_limit"] = result.SizeLimit
			partialResponse["upload_errors"] = result.UploadErrors
			partialResponse["message"] = result.Message
			errors.ResponseSuccess(c, partialResponse, result.Message)
			return
		}
		errors.HandleError(c, err)
		return
	}

	response := gin.H{}
	if result.IsSingleUpload && result.UploadedSingle != nil {
		response["uploaded"] = result.UploadedSingle
	} else if len(result.Uploaded) > 0 {
		response["uploaded"] = result.Uploaded
	}
	if len(result.OversizedFiles) > 0 {
		response["oversized_files"] = result.OversizedFiles
		response["size_limit"] = result.SizeLimit
	}
	if len(result.UploadErrors) > 0 {
		response["upload_errors"] = result.UploadErrors
	}
	errors.ResponseSuccess(c, response, result.Message)
}

// CheckDuplicate MD5预检查重复文件
func CheckDuplicate(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	req, err := common.ValidateRequest[dto.CheckDuplicateDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}
	result, err := filesvc.CheckDuplicateByMD5(userID, req.MD5, req.FileName, req.FileSize)
	if err != nil {
		errors.HandleError(c, err)
		return
	}
	errors.ResponseSuccess(c, result, "检查完成")
}

// InstantUpload 秒传上传
func InstantUpload(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	req, err := common.ValidateRequest[dto.InstantUploadDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}
	if req.StorageDuration != "" {
		storageConfig, err := setting.CreateStorageConfig()
		if err != nil {
			logger.Error("InstantUpload: 创建存储配置失败: %v", err)
			errors.HandleError(c, errors.Wrap(err, errors.CodeDBQueryFailed, "创建存储配置失败"))
			return
		}
		isGuest := userID == 0
		if err := storageConfig.ValidateStorageDuration(req.StorageDuration, isGuest); err != nil {
			logger.Error("InstantUpload: 存储时长验证失败: %v", err)
			errors.HandleError(c, errors.New(errors.CodeInvalidParameter, err.Error()))
			return
		}
	}
	result, err := filesvc.InstantUploadWithDuration(c, userID, req.MD5, req.FileName, req.FileSize, req.FolderID, req.AccessLevel, req.Optimize, req.StorageDuration)
	if err != nil {
		errors.HandleError(c, err)
		return
	}
	errors.ResponseSuccess(c, result, "秒传成功")
}
