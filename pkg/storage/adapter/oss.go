package adapter

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"pixelpunk/pkg/imagex/compress"
	"pixelpunk/pkg/imagex/convert"
	"pixelpunk/pkg/imagex/decode"
	"pixelpunk/pkg/imagex/formats"
	"pixelpunk/pkg/imagex/iox"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/storage/config"
	"pixelpunk/pkg/storage/pipeline"
	"pixelpunk/pkg/storage/tenant"
	"pixelpunk/pkg/storage/utils"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	osscred "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

// OSSAdapter 阿里云OSS存储适配器
type OSSAdapter struct {
	client        *oss.Client
	bucket        string
	region        string
	endpoint      string
	accessKey     string
	secretKey     string
	customDomain  string
	useHTTPS      bool
	accessControl string // 访问控制类型：public-read/private
	initialized   bool
}

func NewOSSAdapter() StorageAdapter {
	return &OSSAdapter{}
}

func (a *OSSAdapter) GetType() string {
	return "oss"
}

// Initialize 初始化适配器
func (a *OSSAdapter) Initialize(configData map[string]interface{}) error {
	cfg := config.NewMapConfig(configData)

	a.bucket = cfg.GetStringWithDefault("bucket", "")
	a.region = cfg.GetStringWithDefault("region", "")
	a.endpoint = cfg.GetStringWithDefault("endpoint", "")
	a.accessKey = cfg.GetStringWithDefault("access_key", "")
	a.secretKey = cfg.GetStringWithDefault("secret_key", "")
	a.useHTTPS = cfg.GetBoolWithDefault("use_https", true)
	a.customDomain, a.useHTTPS = normalizeDomainAndScheme(cfg.GetString("custom_domain"), a.useHTTPS)
	a.accessControl = cfg.GetString("access_control")

	if a.bucket == "" {
		return NewStorageError(ErrorTypeInternal, "bucket is required", nil)
	}
	if a.region == "" {
		return NewStorageError(ErrorTypeInternal, "region is required", nil)
	}
	if a.accessKey == "" {
		return NewStorageError(ErrorTypeInternal, "access_key is required", nil)
	}
	if a.secretKey == "" {
		return NewStorageError(ErrorTypeInternal, "secret_key is required", nil)
	}

	// 如果没有提供endpoint，根据region生成默认endpoint
	if a.endpoint == "" {
		a.endpoint = fmt.Sprintf("oss-%s.aliyuncs.com", a.region)
	}

	cfg_oss := oss.LoadDefaultConfig().
		WithCredentialsProvider(osscred.NewStaticCredentialsProvider(a.accessKey, a.secretKey, "")).
		WithRegion(a.region).
		WithEndpoint(a.endpoint)

	// OSS SDK v2不再需要设置HTTPS，直接通过endpoint控制

	client := oss.NewClient(cfg_oss)
	a.client = client
	a.initialized = true

	return nil
}

// Upload 上传文件
func (a *OSSAdapter) Upload(ctx context.Context, req *UploadRequest) (*UploadResult, error) {
	if !a.initialized {
		return nil, NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}

	src, err := req.File.Open()
	if err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to open file", err)
	}
	defer src.Close()

	data, err := iox.ReadAllWithLimit(src, iox.DefaultMaxReadBytes)
	if err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to read file data", err)
	}

	var processedData io.Reader = bytes.NewReader(data)
	var width, height int
	var format string = a.getFileFormat(req.FileName)

	if w, h, f, err := decode.DetectFormat(bytes.NewReader(data)); err == nil {
		width, height = w, h
		if f != "" {
			format = f
		}
	}

	if req.Options != nil && req.Options.Compress {
		compressResult, err := compress.CompressToTargetSize(bytes.NewReader(data), 5.0, &compress.Options{
			MaxWidth:  req.Options.MaxWidth,
			MaxHeight: req.Options.MaxHeight,
			Quality:   req.Options.Quality,
		})
		if err != nil {
			logger.Warn("压缩失败，使用原图: %v", err)
			// 压缩失败，使用原始数据
			processedData = bytes.NewReader(data)
		} else {
			processedData = compressResult.Reader
			width = compressResult.Width
			height = compressResult.Height
			format = compressResult.Format
		}
	}

	if req.Options != nil && req.Options.WebPEnabled {
		buf, _ := io.ReadAll(processedData)
		if webpResult, err := convert.ToWebP(buf, convert.WebPOptions{Quality: req.Options.Quality}); err == nil && webpResult.Converted {
			processedBytes, _ := io.ReadAll(webpResult.Reader)
			processedData = bytes.NewReader(processedBytes)
			format = "webp"
		} else if err != nil {
			logger.Warn("WebP转换失败，使用原格式: %v", err)
		} else {
			processedData = bytes.NewReader(buf)
		}
	}

	originalFileName := req.FileName
	objectPath, err := tenant.BuildObjectKey(req.UserID, req.FolderPath, originalFileName)
	if err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to build object key", err)
	}
	logicalPath := utils.BuildLogicalPath(req.FolderPath, originalFileName)

	savedBytes, err := io.ReadAll(processedData)
	if err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to buffer processed data", err)
	}
	uploadResult, err := a.uploadToOSS(savedBytes, objectPath, req.ContentType)
	if err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to upload to OSS", err)
	}

	var thumbnailPath string
	var thumbnailURL string

	if req.Options != nil && req.Options.GenerateThumb {
		thumbPath, _, _, err := a.generateThumbnail(bytes.NewReader(data), req, objectPath)
		if err != nil {
			logger.Warn("缩略图生成失败: %v", err)
		} else {
			thumbnailPath = thumbPath
			// 缩略图使用逻辑路径，格式与原图一致，添加_thumb后缀
			// 获取缩略图的实际格式（从生成的缩略图路径推断）
			thumbFormat := "jpg" // 默认格式
			if ext := filepath.Ext(thumbPath); ext != "" {
				thumbFormat = strings.TrimPrefix(strings.ToLower(ext), ".")
			}
			thumbLogicalName := utils.MakeThumbName(originalFileName, thumbFormat)
			thumbnailURL = utils.BuildLogicalPath(req.FolderPath, thumbLogicalName)
		}
	}

	hash := fmt.Sprintf("%x", md5.Sum(savedBytes))

	remoteDirectURL, _ := a.GetURL(objectPath, nil)
	thumbDirectURL := ""
	if thumbnailPath != "" {
		thumbDirectURL, _ = a.GetURL(thumbnailPath, nil)
	}

	result := &UploadResult{
		OriginalPath:   objectPath,
		ThumbnailPath:  thumbnailPath,
		URL:            logicalPath,
		ThumbnailURL:   thumbnailURL,
		FullURL:        remoteDirectURL,
		FullThumbURL:   thumbDirectURL,
		RemoteURL:      objectPath,
		RemoteThumbURL: thumbnailPath,
		Size:           uploadResult.Size,
		Width:          width,
		Height:         height,
		Hash:           hash,
		ContentType:    a.getContentType(format),
		Format:         format,
	}

	return result, nil
}

// Delete 删除文件
func (a *OSSAdapter) Delete(ctx context.Context, path string) error {
	if !a.initialized {
		return NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}

	_, err := a.client.DeleteObject(ctx, &oss.DeleteObjectRequest{
		Bucket: oss.Ptr(a.bucket),
		Key:    oss.Ptr(path),
	})
	return err
}

func (a *OSSAdapter) GetURL(path string, options *URLOptions) (string, error) {
	if !a.initialized {
		return "", NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}

	// 如果设置了自定义域名，使用自定义域名
	if a.customDomain != "" {
		scheme := "https"
		if !a.useHTTPS {
			scheme = "http"
		}
		return fmt.Sprintf("%s://%s/%s", scheme, a.customDomain, encodePathSegments(path)), nil
	}

	scheme := "https"
	if !a.useHTTPS {
		scheme = "http"
	}

	return fmt.Sprintf("%s://%s.%s/%s", scheme, a.bucket, a.endpoint, encodePathSegments(path)), nil
}

//

// HealthCheck 健康检查
func (a *OSSAdapter) HealthCheck(ctx context.Context) error {
	if !a.initialized {
		return NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}

	// 尝试列出bucket以验证连接
	_, err := a.client.ListObjects(ctx, &oss.ListObjectsRequest{
		Bucket:  oss.Ptr(a.bucket),
		MaxKeys: int32(1),
	})
	return err
}

func (a *OSSAdapter) GetCapabilities() Capabilities {
	return Capabilities{
		SupportsSignedURL: true,
		SupportsCDN:       false,
		SupportsResize:    false,
		SupportsWebP:      true,
		MaxFileSize:       5 * 1024 * 1024 * 1024, // 5GB
		SupportedFormats:  []string{"jpg", "jpeg", "png", "gif", "webp", "bmp", "svg", "ico", "apng", "jp2", "tiff", "tif", "tga"},
	}
}

// ReadFile 读取文件
func (a *OSSAdapter) ReadFile(ctx context.Context, path string) (io.ReadCloser, error) {
	if !a.initialized {
		return nil, NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}

	resp, err := a.client.GetObject(ctx, &oss.GetObjectRequest{
		Bucket: oss.Ptr(a.bucket),
		Key:    oss.Ptr(path),
	})
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

// GetBase64 获取文件的Base64编码
// GetBase64 / GetThumbnailBase64 已统一到 Manager 层实现

// Exists 检查文件是否存在
func (a *OSSAdapter) Exists(ctx context.Context, path string) (bool, error) {
	if !a.initialized {
		return false, NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}

	_, err := a.client.HeadObject(ctx, &oss.HeadObjectRequest{
		Bucket: oss.Ptr(a.bucket),
		Key:    oss.Ptr(path),
	})
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (a *OSSAdapter) SetObjectACL(ctx context.Context, path string, acl string) error {
	if !a.initialized {
		return NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}

	var ossACL oss.ObjectACLType
	switch acl {
	case "public-read":
		ossACL = oss.ObjectACLPublicRead
	case "private":
		ossACL = oss.ObjectACLPrivate
	default:
		return NewStorageError(ErrorTypeInternal, "unsupported ACL type: "+acl, nil)
	}

	_, err := a.client.PutObjectAcl(ctx, &oss.PutObjectAclRequest{
		Bucket: oss.Ptr(a.bucket),
		Key:    oss.Ptr(path),
		Acl:    ossACL,
	})

	if err != nil {
		logger.Error("OSS设置对象ACL失败: %v", err)
		return NewStorageError(ErrorTypeInternal, "failed to set object ACL", err)
	}

	return nil
}

// 私有辅助方法

// uploadToOSS 上传数据到OSS
func (a *OSSAdapter) uploadToOSS(dataBytes []byte, objectPath, contentType string) (*UploadResult, error) {
	reader := bytes.NewReader(dataBytes)

	var acl oss.ObjectACLType
	if a.accessControl != "" {
		acl = oss.ObjectACLType(a.accessControl)
	}

	req := &oss.PutObjectRequest{
		Bucket:      oss.Ptr(a.bucket),
		Key:         oss.Ptr(objectPath),
		Body:        reader,
		ContentType: oss.Ptr(contentType),
	}
	if a.accessControl != "" {
		req.Acl = acl
	}

	_, err := a.client.PutObject(context.Background(), req)

	if err != nil {
		return nil, err
	}

	url, err := a.GetURL(objectPath, nil)
	if err != nil {
		return nil, err
	}

	return &UploadResult{
		URL:  url,
		Size: int64(len(dataBytes)),
	}, nil
}

// generateThumbnail 生成缩略图 (与COS保持一致)
func (a *OSSAdapter) generateThumbnail(src io.Reader, req *UploadRequest, originalPath string) (string, string, string, error) {
	// 重新打开源文件进行缩略图处理
	srcFile, err := req.File.Open()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to reopen source file: %w", err)
	}
	defer srcFile.Close()

	data, err := io.ReadAll(srcFile)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to read source data: %w", err)
	}

	// SVG 特判：直接拷贝为缩略图
	if strings.EqualFold(strings.TrimPrefix(strings.ToLower(filepath.Ext(req.FileName)), "."), "svg") {
		thumbFileName := utils.MakeThumbName(req.FileName, "svg")
		thumbObjectPath, _ := tenant.BuildThumbObjectKey(req.UserID, req.FolderPath, thumbFileName)
		uploadResult, err := a.uploadToOSS(data, thumbObjectPath, "image/svg+xml")
		if err != nil {
			return "", "", "", fmt.Errorf("failed to upload thumbnail: %w", err)
		}
		thumbURL := thumbObjectPath
		return thumbObjectPath, thumbURL, uploadResult.URL, nil
	}

	// 统一生成缩略图（带回退）
	q := 85
	if req.Options.ThumbQuality > 0 {
		q = req.Options.ThumbQuality
	}
	w := max(1, coalesceInt(req.Options.ThumbWidth, 1200))
	h := max(1, coalesceInt(req.Options.ThumbHeight, 900))
	thumbBytes, thumbFormat, _ := pipeline.GenerateOrFallback(data, pipeline.Options{
		Width: w, Height: h, Quality: q, EnableWebP: true, FallbackOnError: true,
	})

	thumbFileName := utils.MakeThumbName(req.FileName, thumbFormat)
	thumbObjectPath, _ := tenant.BuildThumbObjectKey(req.UserID, req.FolderPath, thumbFileName)

	uploadResult, err := a.uploadToOSS(thumbBytes, thumbObjectPath, formats.GetContentType(thumbFormat))
	if err != nil {
		return "", "", "", fmt.Errorf("failed to upload thumbnail: %w", err)
	}
	thumbURL := thumbObjectPath
	return thumbObjectPath, thumbURL, uploadResult.URL, nil
}

// removed mustReadAll: no longer needed

func coalesceInt(v int, def int) int {
	if v > 0 {
		return v
	}
	return def
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// getFileFormat 根据文件名获取格式（返回去点小写扩展名；未知返回 unknown）
func (a *OSSAdapter) getFileFormat(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))
	if len(ext) > 1 {
		return strings.TrimPrefix(ext, ".")
	}
	return "unknown"
}

// getContentType 根据格式获取Content-Type
func (a *OSSAdapter) getContentType(format string) string {
	return formats.GetContentType(format)
}
