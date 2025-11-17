package adapter

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"net/url"
	"path/filepath"
	"strings"

	"pixelpunk/pkg/imagex/formats"
	"pixelpunk/pkg/imagex/iox"
	"pixelpunk/pkg/storage/config"
	"pixelpunk/pkg/storage/tenant"
	"pixelpunk/pkg/storage/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// QiniuAdapter 七牛云 Kodo（S3 兼容）适配器
type QiniuAdapter struct {
	client                *s3.Client
	presignClient         *s3.PresignClient
	bucket                string
	region                string
	endpoint              string
	accessKey             string
	secretKey             string
	customDomain          string
	customDomainScheme    string
	customDomainHasScheme bool
	useHTTPS              bool
	usePathStyle          bool
	initialized           bool
}

func NewQiniuAdapter() StorageAdapter { return &QiniuAdapter{} }

func (a *QiniuAdapter) GetType() string { return "qiniu" }

// Initialize 初始化
func (a *QiniuAdapter) Initialize(configData map[string]interface{}) error {
	cfg := config.NewMapConfig(configData)
	a.bucket = cfg.GetStringWithDefault("bucket", "")
	a.region = cfg.GetString("region")
	a.endpoint = cfg.GetStringWithDefault("endpoint", "")
	a.accessKey = cfg.GetStringWithDefault("access_key", "")
	a.secretKey = cfg.GetStringWithDefault("secret_key", "")
	// Normalize custom domain (may include scheme)
	rawDomain := strings.TrimSpace(cfg.GetString("custom_domain"))
	if rawDomain != "" {
		if strings.HasPrefix(rawDomain, "http://") || strings.HasPrefix(rawDomain, "https://") {
			if u, err := url.Parse(rawDomain); err == nil {
				a.customDomain = strings.TrimSuffix(u.Host, "/")
				a.customDomainScheme = u.Scheme
				a.customDomainHasScheme = true
			} else {
				a.customDomain = strings.TrimPrefix(strings.TrimPrefix(rawDomain, "https://"), "http://")
			}
		} else {
			a.customDomain = strings.TrimSuffix(rawDomain, "/")
		}
	}
	a.useHTTPS = cfg.GetBoolWithDefault("use_https", true)
	// 默认：若配置了自定义域名，则采用路径样式 domain/bucket/key；否则使用虚拟主机样式 bucket.endpoint/key
	a.usePathStyle = cfg.GetBoolWithDefault("use_path_style", a.customDomain != "")

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

	// Normalize endpoint: strip scheme if present
	ep := strings.TrimSpace(a.endpoint)
	if strings.HasPrefix(ep, "http://") || strings.HasPrefix(ep, "https://") {
		if u, err := url.Parse(ep); err == nil {
			a.endpoint = u.Host
		} else {
			a.endpoint = strings.TrimPrefix(strings.TrimPrefix(ep, "https://"), "http://")
		}
	}

	awsCfg := aws.Config{
		Region: func() string {
			if a.region != "" {
				return a.region
			}
			return "auto"
		}(),
		Credentials: credentials.NewStaticCredentialsProvider(a.accessKey, a.secretKey, ""),
		EndpointResolverWithOptions: aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{URL: "https://" + a.endpoint, SigningRegion: a.region, HostnameImmutable: true}, nil
		}),
	}
	// 默认使用虚拟主机样式（bucket.endpoint）
	a.client = s3.NewFromConfig(awsCfg)
	a.presignClient = s3.NewPresignClient(a.client)
	a.initialized = true
	return nil
}

// Upload 上传
func (a *QiniuAdapter) Upload(ctx context.Context, req *UploadRequest) (*UploadResult, error) {
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

	// 处理上传数据（压缩等）
	processedBytes, width, height, format := processUploadData(data, req)

	originalFileName := req.FileName
	objectPath, err := tenant.BuildObjectKey(req.UserID, req.FolderPath, originalFileName)
	if err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to build object key", err)
	}
	logicalPath := utils.BuildLogicalPath(req.FolderPath, originalFileName)

	if _, err := a.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(a.bucket),
		Key:         aws.String(objectPath),
		Body:        bytes.NewReader(processedBytes),
		ContentType: aws.String(a.getContentType(format)),
	}); err != nil {
		return nil, NewStorageError(ErrorTypeInternal, "failed to upload to Qiniu", err)
	}

	// thumbnail (optional)
	var thumbnailPath, thumbnailURL, thumbRemoteDirect string
	if req.Options != nil && req.Options.GenerateThumb {
		if tpath, _, remote, err := a.generateThumbnail(bytes.NewReader(data), req, objectPath); err == nil {
			thumbnailPath = tpath
			tf := "jpg"
			if ext := filepath.Ext(tpath); ext != "" {
				tf = strings.TrimPrefix(strings.ToLower(ext), ".")
			}
			thumbnailURL = utils.BuildLogicalPath(req.FolderPath, utils.MakeThumbName(originalFileName, tf))
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
func (a *QiniuAdapter) Delete(ctx context.Context, path string) error {
	if !a.initialized {
		return NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}
	_, err := a.client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: aws.String(a.bucket), Key: aws.String(path)})
	return err
}

// GetURL 生成访问 URL
func (a *QiniuAdapter) GetURL(path string, options *URLOptions) (string, error) {
	if !a.initialized {
		return "", NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}
	scheme := "https"
	if !a.useHTTPS {
		scheme = "http"
	}
	if a.customDomain != "" {
		// Choose scheme: respect domain's own scheme if provided
		if a.customDomainHasScheme {
			scheme = a.customDomainScheme
		}
		domain := strings.TrimSuffix(a.customDomain, "/")
		if a.usePathStyle {
			return fmt.Sprintf("%s://%s/%s/%s", scheme, domain, a.bucket, encodePathSegments(path)), nil
		}
		return fmt.Sprintf("%s://%s/%s", scheme, domain, encodePathSegments(path)), nil
	}
	// Endpoint 直连
	host := a.endpoint
	if a.usePathStyle {
		return fmt.Sprintf("%s://%s/%s/%s", scheme, host, a.bucket, encodePathSegments(path)), nil
	}
	return fmt.Sprintf("%s://%s.%s/%s", scheme, a.bucket, host, encodePathSegments(path)), nil
}

//

// SetObjectACL：Kodo 使用桶策略控制，忽略对象 ACL
func (a *QiniuAdapter) SetObjectACL(ctx context.Context, path string, acl string) error { return nil }

// HealthCheck：List 1 object
func (a *QiniuAdapter) HealthCheck(ctx context.Context) error {
	if !a.initialized {
		return NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}
	_, err := a.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{Bucket: aws.String(a.bucket), MaxKeys: aws.Int32(1)})
	return err
}

// ReadFile 读取对象
func (a *QiniuAdapter) ReadFile(ctx context.Context, path string) (io.ReadCloser, error) {
	if !a.initialized {
		return nil, NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}
	resp, err := a.client.GetObject(ctx, &s3.GetObjectInput{Bucket: aws.String(a.bucket), Key: aws.String(path)})
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// generateThumbnail 与 R2 类似
func (a *QiniuAdapter) generateThumbnail(src io.Reader, req *UploadRequest, originalPath string) (string, string, string, error) {
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
	if _, err := a.client.PutObject(context.Background(), &s3.PutObjectInput{Bucket: aws.String(a.bucket), Key: aws.String(thumbObjectPath), Body: bytes.NewReader(thumbBytes), ContentType: aws.String(formats.GetContentType(thumbFormat))}); err != nil {
		return "", "", "", fmt.Errorf("failed to upload thumbnail: %w", err)
	}
	url, _ := a.GetURL(thumbObjectPath, nil)
	return thumbObjectPath, url, url, nil
}

// Exists 检查对象存在
func (a *QiniuAdapter) Exists(ctx context.Context, path string) (bool, error) {
	if !a.initialized {
		return false, NewStorageError(ErrorTypeInternal, "adapter not initialized", nil)
	}
	_, err := a.client.HeadObject(ctx, &s3.HeadObjectInput{Bucket: aws.String(a.bucket), Key: aws.String(path)})
	if err != nil {
		return false, nil
	}
	return true, nil
}

// GetCapabilities 返回能力
func (a *QiniuAdapter) GetCapabilities() Capabilities {
	return Capabilities{
		SupportsSignedURL: true,
		SupportsCDN:       false,
		SupportsResize:    false,
		SupportsWebP:      true,
		MaxFileSize:       5 * 1024 * 1024 * 1024, // 5GB
		SupportedFormats:  []string{"jpg", "jpeg", "png", "gif", "webp"},
	}
}

func (a *QiniuAdapter) getContentType(format string) string { return formats.GetContentType(format) }
