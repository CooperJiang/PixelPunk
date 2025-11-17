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
	"pixelpunk/pkg/storage/config"
	"pixelpunk/pkg/storage/tenant"
	"pixelpunk/pkg/storage/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// R2Adapter Cloudflare R2 存储适配器（S3 兼容）
type R2Adapter struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucket        string
	region        string
	endpoint      string
	accessKey     string
	secretKey     string
	customDomain  string
	useHTTPS      bool
	accessControl string // public-read/private（R2 多由桶策略控制，仅用于URL生成是否使用签名）
	initialized   bool
}

// NewR2Adapter 创建 R2 适配器
func NewR2Adapter() StorageAdapter { return &R2Adapter{} }

func (a *R2Adapter) GetType() string { return "r2" }

// Initialize 初始化适配器
func (a *R2Adapter) Initialize(configData map[string]interface{}) error {
	cfg := config.NewMapConfig(configData)

	a.bucket = cfg.GetStringWithDefault("bucket", "")
	a.region = cfg.GetStringWithDefault("region", "auto")
	a.endpoint = cfg.GetStringWithDefault("endpoint", "")
	a.accessKey = cfg.GetStringWithDefault("access_key", "")
	a.secretKey = cfg.GetStringWithDefault("secret_key", "")
	a.useHTTPS = cfg.GetBoolWithDefault("use_https", true)
	a.customDomain, a.useHTTPS = normalizeDomainAndScheme(cfg.GetString("custom_domain"), a.useHTTPS)
	a.accessControl = cfg.GetString("access_control")

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

	// AWS S3 SDK v2 配置，Region 用 auto，固定 Endpoint 到 R2；PathStyle 必须 true
	awsCfg := aws.Config{
		Region:      a.region,
		Credentials: credentials.NewStaticCredentialsProvider(a.accessKey, a.secretKey, ""),
		EndpointResolverWithOptions: aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{URL: a.endpoint, SigningRegion: a.region, HostnameImmutable: true}, nil
		}),
	}

	a.client = s3.NewFromConfig(awsCfg, func(o *s3.Options) { o.UsePathStyle = true })
	a.presignClient = s3.NewPresignClient(a.client)
	a.initialized = true
	return nil
}

// Upload 上传文件
func (a *R2Adapter) Upload(ctx context.Context, req *UploadRequest) (*UploadResult, error) {
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

	// 上传到 R2（不设置 ACL，R2 多由桶策略控制）
	if _, err := a.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(a.bucket),
		Key:         aws.String(objectPath),
		Body:        bytes.NewReader(processedBytes),
		ContentType: aws.String(a.getContentType(format)),
		// StorageClass 留空；ACL 不设置
	}); err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to upload to R2", err)
	}

	var thumbnailPath, thumbnailURL, thumbRemoteDirect string
	if req.Options != nil && req.Options.GenerateThumb {
		// 使用原始数据生成缩略图
		if thumbPath, _, remote, err := a.generateThumbnail(bytes.NewReader(data), req, objectPath); err == nil {
			thumbnailPath = thumbPath
			thumbFormat := "jpg"
			if ext := filepath.Ext(thumbPath); ext != "" {
				thumbFormat = strings.TrimPrefix(strings.ToLower(ext), ".")
			}
			thumbLogical := utils.MakeThumbName(originalFileName, thumbFormat)
			thumbnailURL = utils.BuildLogicalPath(req.FolderPath, thumbLogical)
			thumbRemoteDirect = remote
		}
	}

	hash := fmt.Sprintf("%x", md5.Sum(processedBytes))

	url, _ := a.GetURL(objectPath, nil)
	return &UploadResult{
		OriginalPath:   objectPath,
		ThumbnailPath:  thumbnailPath,
		URL:            logicalPath,
		ThumbnailURL:   thumbnailURL,
		FullURL:        url,
		FullThumbURL:   thumbRemoteDirect,
		RemoteURL:      objectPath,
		RemoteThumbURL: thumbnailPath,
		Size:           int64(len(processedBytes)),
		Width:          width,
		Height:         height,
		Hash:           hash,
		ContentType:    a.getContentType(format),
		Format:         format,
	}, nil
}

// Delete 删除
func (a *R2Adapter) Delete(ctx context.Context, path string) error {
	if !a.initialized {
		return NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}
	_, err := a.client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: aws.String(a.bucket), Key: aws.String(path)})
	return err
}

// GetURL 生成访问 URL（私有返回预签名）
func (a *R2Adapter) GetURL(path string, options *URLOptions) (string, error) {
	if !a.initialized {
		return "", NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}
	if a.accessControl == "private" {
		return a.generatePresignedURL(path, options)
	}

	scheme := "https"
	if !a.useHTTPS {
		scheme = "http"
	}
	host := s3ExtractHost(a.endpoint)
	// 使用 path-style
	return s3BuildURL(scheme, a.bucket, host, path, true, a.customDomain), nil
}

//

// SetObjectACL R2 不支持对象级 ACL，这里直接返回 nil
func (a *R2Adapter) SetObjectACL(ctx context.Context, path string, acl string) error { return nil }

// HealthCheck：最小化检查（List 1 object）
func (a *R2Adapter) HealthCheck(ctx context.Context) error {
	if !a.initialized {
		return NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}
	_, err := a.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{Bucket: aws.String(a.bucket), MaxKeys: aws.Int32(1)})
	return err
}

// ReadFile 读取对象
func (a *R2Adapter) ReadFile(ctx context.Context, path string) (io.ReadCloser, error) {
	if !a.initialized {
		return nil, NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}
	resp, err := a.client.GetObject(ctx, &s3.GetObjectInput{Bucket: aws.String(a.bucket), Key: aws.String(path)})
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// generatePresignedURL 生成私有访问签名 URL
func (a *R2Adapter) generatePresignedURL(path string, options *URLOptions) (string, error) {
	if a.presignClient == nil {
		return "", NewStorageError(ErrorTypeInternal, "presign client not initialized", nil)
	}
	expiration := time.Hour
	if options != nil && options.Expires > 0 {
		expiration = time.Duration(options.Expires) * time.Second
	}
	req, err := a.presignClient.PresignGetObject(context.Background(), &s3.GetObjectInput{Bucket: aws.String(a.bucket), Key: aws.String(path)}, func(o *s3.PresignOptions) { o.Expires = expiration })
	if err != nil {
		return "", NewStorageError(ErrorTypeInternal, "failed to generate presigned URL", err)
	}
	return req.URL, nil
}

// generateThumbnail 生成缩略图（复用与 Rainyun 相同逻辑）
func (a *R2Adapter) generateThumbnail(src io.Reader, req *UploadRequest, originalPath string) (string, string, string, error) {
	srcFile, err := req.File.Open()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to reopen source file: %w", err)
	}
	defer srcFile.Close()

	data, err := iox.ReadAllWithLimit(srcFile, iox.DefaultMaxReadBytes)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to read source data: %w", err)
	}

	// thumbnail bytes from original
	thumbBytes, thumbFormat := buildThumbnailBytes(data, req)

	thumbFileName := utils.MakeThumbName(filepath.Base(originalPath), thumbFormat)
	thumbObjectPath, _ := tenant.BuildThumbObjectKey(req.UserID, req.FolderPath, thumbFileName)

	if _, err := a.client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(a.bucket),
		Key:         aws.String(thumbObjectPath),
		Body:        bytes.NewReader(thumbBytes),
		ContentType: aws.String(formats.GetContentType(thumbFormat)),
	}); err != nil {
		return "", "", "", fmt.Errorf("failed to upload thumbnail: %w", err)
	}

	url, _ := a.GetURL(thumbObjectPath, nil)
	return thumbObjectPath, url, url, nil
}

// Exists 检查是否存在
func (a *R2Adapter) Exists(ctx context.Context, path string) (bool, error) {
	if !a.initialized {
		return false, NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}
	_, err := a.client.HeadObject(ctx, &s3.HeadObjectInput{Bucket: aws.String(a.bucket), Key: aws.String(path)})
	if err != nil {
		return false, nil
	}
	return true, nil
}

// GetCapabilities 返回 R2 能力
func (a *R2Adapter) GetCapabilities() Capabilities {
	return Capabilities{
		SupportsSignedURL: true,
		SupportsCDN:       false,
		SupportsResize:    false,
		SupportsWebP:      true,
		MaxFileSize:       5 * 1024 * 1024 * 1024, // 5GB
		SupportedFormats:  []string{"jpg", "jpeg", "png", "gif", "webp"},
	}
}

func (a *R2Adapter) getContentType(format string) string { return formats.GetContentType(format) }
