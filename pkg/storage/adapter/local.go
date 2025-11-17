package adapter

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"pixelpunk/pkg/imagex/compress"
	"pixelpunk/pkg/imagex/convert"
	"pixelpunk/pkg/imagex/decode"
	"pixelpunk/pkg/imagex/formats"
	imgHash "pixelpunk/pkg/imagex/hash"
	"pixelpunk/pkg/imagex/iox"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/storage/config"
	"pixelpunk/pkg/storage/middleware"
	"pixelpunk/pkg/storage/pipeline"
	"pixelpunk/pkg/storage/tenant"
	storageutils "pixelpunk/pkg/storage/utils"
	"strings"
)

// LocalAdapter 本地存储适配器
type LocalAdapter struct {
	basePath      string // 基础存储路径
	thumbnailPath string // 缩略图存储路径
	customDomain  string // 自定义域名
	cdnDomain     string // CDN域名
	initialized   bool   // 是否已初始化
}

func NewLocalAdapter() StorageAdapter {
	return &LocalAdapter{}
}

func (a *LocalAdapter) GetType() string {
	return "local"
}

// Initialize 初始化适配器
func (a *LocalAdapter) Initialize(configData map[string]interface{}) error {
	cfg := config.NewMapConfig(configData)

	a.basePath = cfg.GetStringWithDefault("base_path", "uploads/files")
	a.thumbnailPath = cfg.GetStringWithDefault("thumbnail_path", "uploads/thumbnails")
	a.customDomain = cfg.GetString("custom_domain")
	a.cdnDomain = cfg.GetString("cdn_domain")

	if err := a.ensureDirectories(); err != nil {
		return NewStorageError(
			ErrorTypeInternal,
			"failed to create storage directories",
			err,
		)
	}

	a.initialized = true
	return nil
}

// Upload 上传文件
func (a *LocalAdapter) Upload(ctx context.Context, req *UploadRequest) (*UploadResult, error) {
	if !a.initialized {
		return nil, NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}

	if err := a.validateFile(req); err != nil {
		return nil, NewStorageError(ErrorTypeInvalidFormat, "file validation failed", err)
	}

	var src io.Reader
	var isPreProcessed bool = false

	if len(req.ProcessedData) > 0 {
		// 使用预处理后的数据（如水印处理后的数据）
		src = bytes.NewReader(req.ProcessedData)
		isPreProcessed = true
	} else {
		file, err := req.File.Open()
		if err != nil {
			return nil, NewStorageError(ErrorTypeInternal, "failed to open source file", err)
		}
		defer file.Close()
		src = file
	}

	var processedData io.Reader
	var format string
	var width, height int
	var err error

	if isPreProcessed {
		// 对于预处理数据，只检测格式和尺寸，不进行进一步处理
		processedData, format, width, height, err = a.processPreProcessedImage(src, req)
	} else {
		// 对于原始文件，进行正常的图像处理
		processedData, format, width, height, err = a.processImage(src, req)
	}
	if err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "image processing failed", err)
	}

	originalFileName := req.FileName
	objectKey, err := tenant.BuildObjectKey(req.UserID, req.FolderPath, originalFileName)
	if err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to build object key", err)
	}
	// objectKey: files/{shard}/alias/(folder)/file
	// 本地物理路径: basePath/(after "files/")
	rel := strings.TrimPrefix(objectKey, "files/")
	fullPath := filepath.Join(a.basePath, rel)
	logicalRelativePath := storageutils.BuildLogicalPath(req.FolderPath, originalFileName)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to create directory", err)
	}

	fileSize, err := a.saveFile(processedData, fullPath)
	if err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to save file", err)
	}

	var thumbnailPath string
	var thumbnailURL string
	var thumbnailErr error
	if req.Options != nil && req.Options.GenerateThumb {
		// 生成缩略图文件路径：使用对象键映射
		_, thumbnailPath, thumbnailURL, thumbnailErr = a.generateThumbnailBeforeWebP(fullPath, req, rel, logicalRelativePath)
		if thumbnailErr != nil {
			// 缩略图生成失败不影响主文件上传
			logger.Warn("Local storage: 缩略图生成失败: %v", thumbnailErr)
		}
	}

	hash, _ := imgHash.FromFile(fullPath)

	result := &UploadResult{
		OriginalPath:              fullPath,
		ThumbnailPath:             thumbnailPath,
		URL:                       logicalRelativePath,
		ThumbnailURL:              thumbnailURL,
		Size:                      fileSize,
		Width:                     width,
		Height:                    height,
		Hash:                      hash,
		ContentType:               a.getContentType(format),
		Format:                    format,
		ThumbnailGenerationFailed: thumbnailErr != nil,
		ThumbnailFailureReason: func() string {
			if thumbnailErr != nil {
				return thumbnailErr.Error()
			}
			return ""
		}(),
	}
	result.FullURL = a.buildFullURL(logicalRelativePath, false)
	if thumbnailURL != "" {
		result.FullThumbURL = a.buildFullURL(thumbnailURL, true)
	}
	result.RemoteURL = ""
	result.RemoteThumbURL = ""
	return result, nil
}

// Delete 删除文件
func (a *LocalAdapter) Delete(ctx context.Context, path string) error {
	if !a.initialized {
		return NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}

	p := strings.TrimSpace(path)
	if p == "" {
		return nil
	}

	if strings.HasPrefix(p, "thumbnails/") || strings.HasPrefix(p, "/thumbnails/") {
		tclean := strings.TrimPrefix(p, "/")
		tclean = strings.TrimPrefix(tclean, "thumbnails/")
		_ = os.Remove(filepath.Join(a.thumbnailPath, tclean))
		return nil
	}

	if filepath.IsAbs(p) {
		if strings.HasPrefix(p, a.basePath) || strings.HasPrefix(p, a.thumbnailPath) {
			_ = os.Remove(p)
			return nil
		}
		// 非本适配器目录，忽略
		return nil
	}

	clean := strings.TrimPrefix(p, "/")
	// 如果是对象键 files/
	if strings.HasPrefix(clean, "files/") {
		clean = strings.TrimPrefix(clean, "files/")
		_ = os.Remove(filepath.Join(a.basePath, clean))
		return nil
	}
	// 如果是缩略图物理相对路径（带有 thumbnailPath 前缀）
	thumbRelBase := strings.TrimPrefix(filepath.Clean(a.thumbnailPath), "/")
	if strings.HasPrefix(clean, thumbRelBase+"/") {
		rel := strings.TrimPrefix(clean, thumbRelBase+"/")
		_ = os.Remove(filepath.Join(a.thumbnailPath, rel))
		return nil
	}
	// 如果是原图物理相对路径（带有 basePath 前缀）
	baseRel := strings.TrimPrefix(filepath.Clean(a.basePath), "/")
	if strings.HasPrefix(clean, baseRel+"/") {
		rel := strings.TrimPrefix(clean, baseRel+"/")
		_ = os.Remove(filepath.Join(a.basePath, rel))
		return nil
	}
	// 否则作为相对路径直接拼到 basePath
	_ = os.Remove(filepath.Join(a.basePath, clean))
	return nil
}

// Exists 检查文件是否存在
func (a *LocalAdapter) Exists(ctx context.Context, path string) (bool, error) {
	if !a.initialized {
		return false, NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}

	fullPath := filepath.Join(a.basePath, path)
	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, NewStorageError(
			ErrorTypeInternal,
			"failed to check file existence",
			err,
		)
	}
	return true, nil
}

func (a *LocalAdapter) SetObjectACL(ctx context.Context, path string, acl string) error {
	// 本地存储不支持ACL设置，直接返回成功
	// 本地存储不支持ACL设置
	return nil
}

func (a *LocalAdapter) GetURL(path string, options *URLOptions) (string, error) {
	if !a.initialized {
		return "", NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}

	return path, nil
}

// ReadFile 读取文件
func (a *LocalAdapter) ReadFile(ctx context.Context, path string) (io.ReadCloser, error) {
	if !a.initialized {
		return nil, NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}

	// 支持对象键前缀映射：files/ -> basePath, thumbnails/ -> thumbnailPath
	clean := strings.TrimPrefix(path, "/")
	var fullPath string
	if strings.HasPrefix(clean, "thumbnails/") {
		rel := strings.TrimPrefix(clean, "thumbnails/")
		fullPath = filepath.Join(a.thumbnailPath, rel)
	} else if strings.HasPrefix(clean, "files/") {
		rel := strings.TrimPrefix(clean, "files/")
		fullPath = filepath.Join(a.basePath, rel)
	} else {
		fullPath = filepath.Join(a.basePath, clean)
	}
	file, err := os.Open(fullPath)
	if os.IsNotExist(err) {
		return nil, NewStorageError(ErrorTypeNotFound, "file not found", err)
	}
	if err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to open file", err)
	}

	return file, nil
}

// GetBase64 获取文件的Base64编码
// GetBase64 / GetThumbnailBase64 已统一到 Manager 层实现

// HealthCheck 健康检查
func (a *LocalAdapter) HealthCheck(ctx context.Context) error {
	if !a.initialized {
		return NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}

	testFile := filepath.Join(a.basePath, ".health_check")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return NewStorageError(
			ErrorTypePermission,
			"storage directory not writable",
			err,
		)
	}

	os.Remove(testFile)

	return nil
}

func (a *LocalAdapter) GetCapabilities() Capabilities {
	return Capabilities{
		SupportsSignedURL: false,
		SupportsCDN:       a.cdnDomain != "",
		SupportsResize:    false,
		SupportsWebP:      true,
		MaxFileSize:       100 * 1024 * 1024, // 100MB
		SupportedFormats:  []string{"jpg", "jpeg", "png", "gif", "webp", "bmp"},
	}
}

// 私有辅助方法

// ensureDirectories 确保目录存在
func (a *LocalAdapter) ensureDirectories() error {
	if err := os.MkdirAll(a.basePath, 0755); err != nil {
		return fmt.Errorf("failed to create base directory: %w", err)
	}
	if err := os.MkdirAll(a.thumbnailPath, 0755); err != nil {
		return fmt.Errorf("failed to create thumbnail directory: %w", err)
	}
	return nil
}

// generateThumbnailBeforeWebP 基于已保存原图路径生成缩略图
func (a *LocalAdapter) generateThumbnailBeforeWebP(originalFullPath string, req *UploadRequest, physicalRelativePath, logicalRelativePath string) (io.Reader, string, string, error) {
	data, err := os.ReadFile(originalFullPath)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to read source data: %w", err)
	}

	// 归一化（无日志版）
	data = NormalizePossiblyTextualBytes(data, "Local thumb")

	if req.Options.Compress && (req.Options.MaxWidth > 0 || req.Options.MaxHeight > 0) {
		compressOptions := &compress.Options{
			MaxWidth:  req.Options.MaxWidth,
			MaxHeight: req.Options.MaxHeight,
			Quality:   req.Options.Quality,
			Preserve:  true,
		}
		if cr, err := compress.CompressFile(bytes.NewReader(data), compressOptions); err == nil {
			if buf, e := io.ReadAll(cr.Reader); e == nil {
				data = buf
			}
		}
	}

	// 统一生成缩略图（带回退）
	q := req.Options.ThumbQuality
	if q <= 0 {
		q = 85
	}
	w := req.Options.ThumbWidth
	h := req.Options.ThumbHeight
	thumbBytes, thumbFormat, _ := pipeline.GenerateOrFallback(data, pipeline.Options{
		Width: w, Height: h, Quality: q, EnableWebP: true, FallbackOnError: true,
	})
	thumbData := bytes.NewReader(thumbBytes)
	if thumbFormat == "" {
		thumbFormat = "jpg"
	}

	// 保存缩略图（按 alias 分片对象键映射到本地路径）
	thumbFileName := storageutils.MakeThumbName(req.FileName, thumbFormat)
	thumbKey, _ := tenant.BuildThumbObjectKey(req.UserID, req.FolderPath, thumbFileName)
	thumbRel := strings.TrimPrefix(thumbKey, "thumbnails/")
	thumbFullPath := filepath.Join(a.thumbnailPath, thumbRel)
	if err := os.MkdirAll(filepath.Dir(thumbFullPath), 0755); err != nil {
		return nil, "", "", fmt.Errorf("failed to create thumbnail directory: %w", err)
	}
	if _, err := a.saveFile(thumbData, thumbFullPath); err != nil {
		return nil, "", "", fmt.Errorf("failed to save thumbnail: %w", err)
	}
	thumbLogicalPath := filepath.Join(req.FolderPath, thumbFileName)
	thumbURL := thumbLogicalPath
	return thumbData, thumbFullPath, thumbURL, nil
}

// 本地适配器文件哈希已统一使用 imagex/hash

// getFileFormat 获取文件格式
func (a *LocalAdapter) getFileFormat(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))
	return strings.TrimPrefix(ext, ".")
}

func (a *LocalAdapter) buildFullURL(relativePath string, isThumbnail bool) string {
	return ""
}

// processPreProcessedImage 处理预处理过的图像数据（如水印处理后的数据）
// 只检测格式和尺寸，不进行进一步的图像处理以保持预处理效果
func (a *LocalAdapter) processPreProcessedImage(src io.Reader, req *UploadRequest) (io.Reader, string, int, int, error) {
	data, err := iox.ReadAllWithLimit(src, iox.DefaultMaxReadBytes)
	if err != nil {
		return nil, "", 0, 0, err
	}

	width, height, format, err := decode.DetectFormat(bytes.NewReader(data))
	if err != nil {
		// 如果无法检测格式，使用文件扩展名
		format = a.getFileFormat(req.FileName)
		width, height = 0, 0 // 设为0表示无法获取尺寸
	}

	return bytes.NewReader(data), format, width, height, nil
}

// validateFile 验证文件
func (a *LocalAdapter) validateFile(req *UploadRequest) error {
	// 如果使用预处理数据，跳过文件验证（假设预处理数据已经是有效的）
	if len(req.ProcessedData) > 0 {
		if len(req.ProcessedData) > 20*1024*1024 { // 20MB
			return fmt.Errorf("processed data size %d exceeds maximum limit", len(req.ProcessedData))
		}
		return nil
	}

	options := &middleware.ValidationOptions{
		MaxFileSize:     20 * 1024 * 1024,
		AllowedFormats:  formats.SupportedExtensionsWithDot(),
		CheckFileHeader: true,
		CheckMimeType:   true,
	}
	return middleware.ValidateSingleFile(req.File, options)
}

// processImage 处理图像（先压缩，后WebP转换以避免缩略图问题）
func (a *LocalAdapter) processImage(src io.Reader, req *UploadRequest) (io.Reader, string, int, int, error) {
	if req.Options == nil {
		// 没有处理选项，直接读取原始数据
		data, err := iox.ReadAllWithLimit(src, iox.DefaultMaxReadBytes)
		if err != nil {
			return nil, "", 0, 0, err
		}
		// 多轮归一化 - 修复前端传来的错误格式数据（无条件调用，内部判断）
		data = NormalizePossiblyTextualBytes(data, "Local process")

		// 检查并转换 HEIC/HEIF 为 JPEG
		if convert.IsHEICFormat(data) {
			heicResult, err := convert.ToJPEGFromHEIC(data, convert.HEICToJPEGOptions{Quality: 95})
			if err == nil && heicResult.Converted {
				if buf, e := io.ReadAll(heicResult.Reader); e == nil {
					data = buf
					req.FileName = replaceHEICExtension(req.FileName)
				}
			}
		}

		width, height, format, err := decode.DetectFormat(bytes.NewReader(data))
		if err != nil {
			// 如果无法获取尺寸，可能不是图像文件
			return bytes.NewReader(data), a.getFileFormat(req.FileName), 0, 0, nil
		}
		return bytes.NewReader(data), format, width, height, nil
	}

	data, err := iox.ReadAllWithLimit(src, iox.DefaultMaxReadBytes)
	if err != nil {
		return nil, "", 0, 0, err
	}
	// 多轮归一化 - 修复前端传来的错误格式数据（无条件调用，内部判断）
	data = NormalizePossiblyTextualBytes(data, "Local process")

	// 检查并转换 HEIC/HEIF 为 JPEG
	if convert.IsHEICFormat(data) {
		heicResult, err := convert.ToJPEGFromHEIC(data, convert.HEICToJPEGOptions{Quality: 95})
		if err == nil && heicResult.Converted {
			if buf, e := io.ReadAll(heicResult.Reader); e == nil {
				data = buf
				req.FileName = replaceHEICExtension(req.FileName)
			}
		}
	}

	originalFormat := a.getFileFormat(req.FileName)
	currentData := bytes.NewReader(data)
	currentFormat := originalFormat
	var width, height int

	// 获取原始图像尺寸（容错处理）
	if w, h, f, derr := decode.DetectFormat(bytes.NewReader(data)); derr == nil {
		width, height = w, h
		if f != "" {
			currentFormat = f
		}
	} else {
		// 格式检测失败时使用文件扩展名作为格式，保持兼容性
		width, height = 0, 0 // 设置为0表示尺寸未知
	}

	if req.Options.Compress && (req.Options.MaxWidth > 0 || req.Options.MaxHeight > 0) {
		compressOptions := &compress.Options{
			MaxWidth:  req.Options.MaxWidth,
			MaxHeight: req.Options.MaxHeight,
			Quality:   req.Options.Quality,
			Preserve:  true,
		}

		compressResult, err := compress.CompressFile(currentData, compressOptions)
		if err == nil {
			// 读取压缩结果到bytes.Reader
			compressData, readErr := io.ReadAll(compressResult.Reader)
			if readErr == nil {
				currentData = bytes.NewReader(compressData)
				width = compressResult.Width
				height = compressResult.Height
			} else {
			}
		} else {
		}
	}

	if req.Options.WebPEnabled {
		buf, _ := io.ReadAll(currentData)
		if webpResult, err := convert.ToWebP(buf, convert.WebPOptions{Quality: req.Options.Quality}); err == nil && webpResult.Converted {
			webpData, _ := io.ReadAll(webpResult.Reader)
			currentData = bytes.NewReader(webpData)
			currentFormat = "webp"
		} else {
			currentData = bytes.NewReader(buf)
		}
	}

	return currentData, currentFormat, width, height, nil
}

// saveFile 保存文件数据
func (a *LocalAdapter) saveFile(data io.Reader, filePath string) (int64, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	size, err := io.Copy(file, data)
	if err != nil {
		os.Remove(filePath) // 清理失败的文件
		return 0, err
	}

	return size, nil
}

// getContentType 根据格式获取内容类型
func (a *LocalAdapter) getContentType(format string) string {
	return formats.GetContentType(format)
}

// 删除了 useDefaultFailThumbnail 方法，因为已由 pipeline.GenerateOrFallback 统一回退逻辑替代
