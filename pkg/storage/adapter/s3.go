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
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3Adapter 通用 S3 兼容存储适配器（AWS S3 / MinIO / 其他兼容端点）
type S3Adapter struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucket        string
	region        string
	endpoint      string // 为空表示使用 AWS 官方端点
	accessKey     string
	secretKey     string
	customDomain  string
	useHTTPS      bool
	usePathStyle  bool
	accessControl string // public-read/private（也可留空，使用桶策略）
	initialized   bool
}

func NewS3Adapter() StorageAdapter   { return &S3Adapter{} }
func (a *S3Adapter) GetType() string { return "s3" }

// Initialize 初始化
func (a *S3Adapter) Initialize(configData map[string]interface{}) error {
	cfg := config.NewMapConfig(configData)
	a.bucket = cfg.GetStringWithDefault("bucket", "")
	a.region = cfg.GetStringWithDefault("region", "us-east-1")
	a.endpoint = cfg.GetStringWithDefault("endpoint", "")
	a.accessKey = cfg.GetStringWithDefault("access_key", "")
	a.secretKey = cfg.GetStringWithDefault("secret_key", "")
	a.useHTTPS = cfg.GetBoolWithDefault("use_https", true)
	a.usePathStyle = cfg.GetBoolWithDefault("use_path_style", false)
	a.accessControl = cfg.GetString("access_control")
	a.customDomain, a.useHTTPS = normalizeDomainAndScheme(cfg.GetString("custom_domain"), a.useHTTPS)

	if a.bucket == "" {
		return NewStorageError(ErrorTypeInternal, "bucket is required", nil)
	}
	if a.accessKey == "" {
		return NewStorageError(ErrorTypeInternal, "access_key is required", nil)
	}
	if a.secretKey == "" {
		return NewStorageError(ErrorTypeInternal, "secret_key is required", nil)
	}

	awsCfg := aws.Config{
		Region:      a.region,
		Credentials: credentials.NewStaticCredentialsProvider(a.accessKey, a.secretKey, ""),
	}
	if strings.TrimSpace(a.endpoint) != "" {
		ep := a.endpoint
		awsCfg.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{URL: ep, SigningRegion: a.region, HostnameImmutable: true}, nil
		})
	}

	if a.endpoint != "" {
		a.client = s3.NewFromConfig(awsCfg, func(o *s3.Options) { o.UsePathStyle = a.usePathStyle })
	} else {
		// 官方 S3 端点：默认使用虚拟主机样式，更符合规范
		a.client = s3.NewFromConfig(awsCfg, func(o *s3.Options) { o.UsePathStyle = a.usePathStyle })
	}
	a.presignClient = s3.NewPresignClient(a.client)
	a.initialized = true
	return nil
}

// Upload 上传
func (a *S3Adapter) Upload(ctx context.Context, req *UploadRequest) (*UploadResult, error) {
	if !a.initialized {
		return nil, NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}

	var data []byte
	if len(req.ProcessedData) > 0 {
		data = req.ProcessedData
	} else {
		src, err := req.File.Open()
		if err != nil {
			return nil, NewStorageError(ErrorTypeInternal, "failed to open file", err)
		}
		defer src.Close()
		buf, err := iox.ReadAllWithLimit(src, iox.DefaultMaxReadBytes)
		if err != nil {
			return nil, NewStorageError(ErrorTypeInternal, "failed to read file data", err)
		}
		data = buf
	}

	processed, width, height, format := processUploadData(data, req)

	originalFileName := req.FileName
	objectPath, err := tenant.BuildObjectKey(req.UserID, req.FolderPath, originalFileName)
	if err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to build object key", err)
	}
	logicalPath := utils.BuildLogicalPath(req.FolderPath, originalFileName)

	put := &s3.PutObjectInput{
		Bucket:      aws.String(a.bucket),
		Key:         aws.String(objectPath),
		Body:        bytes.NewReader(processed),
		ContentType: aws.String(a.getContentType(format)),
	}
	if acl, ok := s3MapACL(a.accessControl); ok {
		put.ACL = acl
	}
	if _, err := a.client.PutObject(ctx, put); err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to upload to S3", err)
	}

	var thumbnailPath, thumbnailURL string
	if req.Options != nil && req.Options.GenerateThumb {
		if tpath, _, _, err := a.generateThumbnail(bytes.NewReader(data), req, objectPath); err == nil {
			thumbnailPath = tpath
			tf := "jpg"
			if ext := filepath.Ext(tpath); ext != "" {
				tf = strings.TrimPrefix(strings.ToLower(ext), ".")
			}
			thumbnailURL = utils.BuildLogicalPath(req.FolderPath, utils.MakeThumbName(originalFileName, tf))
		}
	}

	hash := fmt.Sprintf("%x", md5.Sum(processed))
	direct, _ := a.GetURL(objectPath, nil)
	thumbDirectURL := ""
	if thumbnailPath != "" {
		thumbDirectURL, _ = a.GetURL(thumbnailPath, nil)
	}
	return &UploadResult{
		OriginalPath:   objectPath,
		ThumbnailPath:  thumbnailPath,
		URL:            logicalPath,
		ThumbnailURL:   thumbnailURL,
		FullURL:        direct,
		FullThumbURL:   thumbDirectURL,
		RemoteURL:      objectPath,
		RemoteThumbURL: thumbnailPath,
		Size:           int64(len(processed)),
		Width:          width,
		Height:         height,
		Hash:           hash,
		ContentType:    a.getContentType(format),
		Format:         format,
	}, nil
}

// Delete 删除
func (a *S3Adapter) Delete(ctx context.Context, path string) error {
	if !a.initialized {
		return NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}
	_, err := a.client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: aws.String(a.bucket), Key: aws.String(path)})
	return err
}

// GetURL 生成访问 URL
func (a *S3Adapter) GetURL(path string, options *URLOptions) (string, error) {
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
	host := func() string {
		if strings.TrimSpace(a.endpoint) != "" {
			return s3ExtractHost(a.endpoint)
		}
		return s3DefaultHostForRegion(a.region)
	}()
	return s3BuildURL(scheme, a.bucket, host, path, a.usePathStyle || strings.TrimSpace(a.endpoint) != "", a.customDomain), nil
}

//

func (a *S3Adapter) generatePresignedURL(path string, options *URLOptions) (string, error) {
	if a.presignClient == nil {
		return "", NewStorageError(ErrorTypeInternal, "presign client not initialized", nil)
	}
	exp := time.Hour
	if options != nil && options.Expires > 0 {
		exp = time.Duration(options.Expires) * time.Second
	}
	req, err := a.presignClient.PresignGetObject(context.Background(), &s3.GetObjectInput{Bucket: aws.String(a.bucket), Key: aws.String(path)}, func(o *s3.PresignOptions) { o.Expires = exp })
	if err != nil {
		return "", NewStorageError(ErrorTypeInternal, "failed to generate presigned URL", err)
	}
	return req.URL, nil
}

func (a *S3Adapter) SetObjectACL(ctx context.Context, path string, acl string) error {
	if !a.initialized {
		return NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}
	var canned types.ObjectCannedACL
	switch acl {
	case "public-read":
		canned = types.ObjectCannedACLPublicRead
	case "private":
		canned = types.ObjectCannedACLPrivate
	default:
		return nil
	}
	_, err := a.client.PutObjectAcl(ctx, &s3.PutObjectAclInput{Bucket: aws.String(a.bucket), Key: aws.String(path), ACL: canned})
	return err
}

// HealthCheck 简单列举
func (a *S3Adapter) HealthCheck(ctx context.Context) error {
	if !a.initialized {
		return NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}
	_, err := a.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{Bucket: aws.String(a.bucket), MaxKeys: aws.Int32(1)})
	return err
}

// ReadFile 读取对象
func (a *S3Adapter) ReadFile(ctx context.Context, path string) (io.ReadCloser, error) {
	if !a.initialized {
		return nil, NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}
	resp, err := a.client.GetObject(ctx, &s3.GetObjectInput{Bucket: aws.String(a.bucket), Key: aws.String(path)})
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// Exists 检查对象是否存在
func (a *S3Adapter) Exists(ctx context.Context, path string) (bool, error) {
	if !a.initialized {
		return false, NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}
	_, err := a.client.HeadObject(ctx, &s3.HeadObjectInput{Bucket: aws.String(a.bucket), Key: aws.String(path)})
	if err != nil {
		return false, nil
	}
	return true, nil
}

// generateThumbnail 生成缩略图
func (a *S3Adapter) generateThumbnail(src io.Reader, req *UploadRequest, originalPath string) (string, string, string, error) {
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
	thumbFileName := utils.MakeThumbName(filepath.Base(originalPath), thumbFormat)
	thumbObjectPath, _ := tenant.BuildThumbObjectKey(req.UserID, req.FolderPath, thumbFileName)
	put := &s3.PutObjectInput{Bucket: aws.String(a.bucket), Key: aws.String(thumbObjectPath), Body: bytes.NewReader(thumbBytes), ContentType: aws.String(formats.GetContentType(thumbFormat))}
	if acl, ok := s3MapACL(a.accessControl); ok {
		put.ACL = acl
	}
	if _, err := a.client.PutObject(context.Background(), put); err != nil {
		return "", "", "", fmt.Errorf("failed to upload thumbnail: %w", err)
	}
	url, _ := a.GetURL(thumbObjectPath, nil)
	return thumbObjectPath, url, url, nil
}

func (a *S3Adapter) GetCapabilities() Capabilities {
	return Capabilities{SupportsSignedURL: true, SupportsCDN: false, SupportsResize: false, SupportsWebP: true, MaxFileSize: 5 * 1024 * 1024 * 1024, SupportedFormats: []string{"jpg", "jpeg", "png", "gif", "webp"}}
}

func (a *S3Adapter) getContentType(format string) string { return formats.GetContentType(format) }
