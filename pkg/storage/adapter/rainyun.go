package adapter

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"pixelpunk/pkg/imagex/formats"
	"pixelpunk/pkg/imagex/iox"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/storage/config"
	"pixelpunk/pkg/storage/tenant"
	"pixelpunk/pkg/storage/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// RainyunAdapter 雨云对象存储适配器
type RainyunAdapter struct {
	client        *s3.Client
	presignClient *s3.PresignClient
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

func NewRainyunAdapter() StorageAdapter {
	return &RainyunAdapter{}
}

func (a *RainyunAdapter) GetType() string {
	return "rainyun"
}

// Initialize 初始化适配器
func (a *RainyunAdapter) Initialize(configData map[string]interface{}) error {
	cfg := config.NewMapConfig(configData)

	a.bucket = cfg.GetStringWithDefault("bucket", "")
	a.region = cfg.GetStringWithDefault("region", "rainyun")
	a.endpoint = cfg.GetStringWithDefault("endpoint", "")
	a.accessKey = cfg.GetStringWithDefault("access_key", "")
	a.secretKey = cfg.GetStringWithDefault("secret_key", "")
	a.useHTTPS = cfg.GetBoolWithDefault("use_https", true)
	a.customDomain, a.useHTTPS = normalizeDomainAndScheme(cfg.GetString("custom_domain"), a.useHTTPS)
	a.accessControl = cfg.GetString("access_control")
	// RainyUN不支持存储类型配置，统一使用标准存储

	if a.bucket == "" {
		return NewStorageError(ErrorTypeInternal, "bucket is required", nil)
	}
	if a.endpoint == "" {
		return NewStorageError(ErrorTypeInternal, "endpoint is required", nil)
	}
	if a.accessKey == "" {
		return NewStorageError(ErrorTypeInternal, "access_key is required", nil)
	}
	if a.secretKey == "" {
		return NewStorageError(ErrorTypeInternal, "secret_key is required", nil)
	}

	awsConfig := aws.Config{
		Region:      a.region,
		Credentials: credentials.NewStaticCredentialsProvider(a.accessKey, a.secretKey, ""),
		EndpointResolverWithOptions: aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:               a.endpoint,
				SigningRegion:     a.region,
				HostnameImmutable: true,
			}, nil
		}),
	}

	a.client = s3.NewFromConfig(awsConfig, func(o *s3.Options) {
		o.UsePathStyle = true // 使用路径样式
	})

	a.presignClient = s3.NewPresignClient(a.client)

	a.initialized = true
	return nil
}

// Upload 上传文件
func (a *RainyunAdapter) Upload(ctx context.Context, req *UploadRequest) (*UploadResult, error) {
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

	processedBytes, width, height, format := processUploadData(data, req)

	originalFileName := req.FileName
	objectPath, err := tenant.BuildObjectKey(req.UserID, req.FolderPath, originalFileName)
	if err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to build object key", err)
	}
	logicalPath := utils.BuildLogicalPath(req.FolderPath, originalFileName)

	uploadResult, err := a.uploadToRainyun(processedBytes, objectPath, req.ContentType)
	if err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to upload to RainYun", err)
	}

	var thumbnailPath string
	var thumbnailURL string

	if req.Options != nil && req.Options.GenerateThumb {
		thumbPath, _, _, err := a.generateThumbnail(bytes.NewReader(data), req, objectPath)
		if err != nil {
		} else {
			thumbnailPath = thumbPath
			thumbFormat := "jpg"
			if ext := filepath.Ext(thumbPath); ext != "" {
				thumbFormat = strings.TrimPrefix(strings.ToLower(ext), ".")
			}
			thumbLogicalName := utils.MakeThumbName(originalFileName, thumbFormat)
			thumbnailURL = utils.BuildLogicalPath(req.FolderPath, thumbLogicalName)
		}
	}

	hash := fmt.Sprintf("%x", md5.Sum(processedBytes))

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
func (a *RainyunAdapter) Delete(ctx context.Context, path string) error {
	if !a.initialized {
		return NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}

	_, err := a.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(path),
	})
	return err
}

func (a *RainyunAdapter) GetURL(path string, options *URLOptions) (string, error) {
	if !a.initialized {
		return "", NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}

	// 如果是私有访问，生成预签名URL
	if a.accessControl == "private" {
		return a.generatePresignedURL(path, options)
	}

	// 雨云使用虚拟主机样式 (bucket.endpoint/path)
	scheme := "https"
	if !a.useHTTPS {
		scheme = "http"
	}
	host := s3ExtractHost(a.endpoint)
	return s3BuildURL(scheme, a.bucket, host, path, false, a.customDomain), nil
}

//

// HealthCheck 健康检查
func (a *RainyunAdapter) HealthCheck(ctx context.Context) error {
	if !a.initialized {
		return NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}

	// 尝试列出bucket以验证连接
	_, err := a.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:  aws.String(a.bucket),
		MaxKeys: aws.Int32(1),
	})
	return err
}

func (a *RainyunAdapter) GetCapabilities() Capabilities {
	return Capabilities{
		SupportsSignedURL: true,
		SupportsCDN:       false,
		SupportsResize:    false,
		SupportsWebP:      true,
		MaxFileSize:       5 * 1024 * 1024 * 1024, // 5GB
		SupportedFormats:  []string{"jpg", "jpeg", "png", "gif", "webp"},
	}
}

// generatePresignedURL 生成预签名URL用于私有访问
func (a *RainyunAdapter) generatePresignedURL(path string, options *URLOptions) (string, error) {
	if a.presignClient == nil {
		return "", NewStorageError(ErrorTypeInternal, "presign client not initialized", nil)
	}

	// 设置过期时间，默认1小时
	expiration := time.Hour
	if options != nil && options.Expires > 0 {
		expiration = time.Duration(options.Expires) * time.Second
	}

	request, err := a.presignClient.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(path),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})

	if err != nil {
		logger.Error("生成预签名URL失败: %v", err)
		return "", NewStorageError(ErrorTypeInternal, "failed to generate presigned URL", err)
	}

	return request.URL, nil
}

// ReadFile 读取文件
func (a *RainyunAdapter) ReadFile(ctx context.Context, path string) (io.ReadCloser, error) {
	if !a.initialized {
		return nil, NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}

	if strings.TrimSpace(path) == "" {
		logger.Error("RainYun ReadFile收到空path，拒绝调用S3 GetObject")
		return nil, NewStorageError(ErrorTypeInvalidFormat, "empty object key", nil)
	}

	resp, err := a.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		logger.Error("RainYun ReadFile失败: %v", err)
		return nil, err
	}

	return resp.Body, nil
}

// GetBase64 / GetThumbnailBase64 已统一到 Manager 层实现

// Exists 检查文件是否存在
func (a *RainyunAdapter) Exists(ctx context.Context, path string) (bool, error) {
	if !a.initialized {
		return false, NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}

	_, err := a.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (a *RainyunAdapter) SetObjectACL(ctx context.Context, path string, acl string) error {
	if !a.initialized {
		return NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}

	var cannedACL types.ObjectCannedACL
	switch acl {
	case "public-read":
		cannedACL = types.ObjectCannedACLPublicRead
	case "private":
		cannedACL = types.ObjectCannedACLPrivate
	default:
		return NewStorageError(ErrorTypeInternal, "unsupported ACL type: "+acl, nil)
	}

	_, err := a.client.PutObjectAcl(ctx, &s3.PutObjectAclInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(path),
		ACL:    cannedACL,
	})

	if err != nil {
		logger.Error("RainyUN设置对象ACL失败: %v", err)
		return NewStorageError(ErrorTypeInternal, "failed to set object ACL", err)
	}

	return nil
}

// 私有辅助方法

// uploadToRainyun 上传数据到RainyUN
func (a *RainyunAdapter) uploadToRainyun(dataBytes []byte, objectPath, contentType string) (*UploadResult, error) {
	reader := bytes.NewReader(dataBytes)

	input := &s3.PutObjectInput{
		Bucket:      aws.String(a.bucket),
		Key:         aws.String(objectPath),
		Body:        reader,
		ContentType: aws.String(contentType),
	}

	if a.accessControl != "" {
		switch a.accessControl {
		case "public-read":
			input.ACL = types.ObjectCannedACLPublicRead
		case "private":
			input.ACL = types.ObjectCannedACLPrivate
		}
	}

	// RainyUN统一使用标准存储
	input.StorageClass = types.StorageClassStandard

	_, err := a.client.PutObject(context.Background(), input)
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

// generateThumbnail 生成缩略图 (与COS/OSS保持一致)
func (a *RainyunAdapter) generateThumbnail(src io.Reader, req *UploadRequest, originalPath string) (string, string, string, error) {
	// 重新打开源文件进行缩略图处理
	srcFile, err := req.File.Open()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to reopen source file: %w", err)
	}
	defer srcFile.Close()

	data, err := iox.ReadAllWithLimit(srcFile, iox.DefaultMaxReadBytes)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to read source data: %w", err)
	}

	thumbBytes, thumbFormat := buildThumbnailBytes(data, req)

	// 保存缩略图到 Rainyun
	thumbFileName := utils.MakeThumbName(filepath.Base(originalPath), thumbFormat)
	thumbObjectPath, _ := tenant.BuildThumbObjectKey(req.UserID, req.FolderPath, thumbFileName)

	thumbUploadResult, err := a.uploadToRainyun(thumbBytes, thumbObjectPath, formats.GetContentType(thumbFormat))
	if err != nil {
		return "", "", "", fmt.Errorf("failed to upload thumbnail: %w", err)
	}

	return thumbObjectPath, thumbUploadResult.URL, thumbUploadResult.URL, nil
}

// calculateDataHash 计算数据哈希值
// removed calculateDataHash: simplified to md5.Sum on []byte

// getContentType 根据格式获取Content-Type
func (a *RainyunAdapter) getContentType(format string) string {
	return formats.GetContentType(format)
}
