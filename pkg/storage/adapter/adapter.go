package adapter

import (
	"context"
	"io"
	"mime/multipart"
)

// StorageAdapter 存储适配器接口
// 定义了所有存储类型必须实现的核心功能
type StorageAdapter interface {
	Upload(ctx context.Context, req *UploadRequest) (*UploadResult, error)
	Delete(ctx context.Context, path string) error
	Exists(ctx context.Context, path string) (bool, error)

	// URL 生成（仅保留相对/直链的基本能力；完整URL由上层策略生成）
	GetURL(path string, options *URLOptions) (string, error)

	ReadFile(ctx context.Context, path string) (io.ReadCloser, error)

	// Base64 数据获取已统一在 Manager 层实现（通过 ReadFile 实现）

	SetObjectACL(ctx context.Context, path string, acl string) error

	Initialize(config map[string]interface{}) error
	HealthCheck(ctx context.Context) error
	GetType() string
	GetCapabilities() Capabilities
}

// UploadRequest 上传请求
type UploadRequest struct {
	File          *multipart.FileHeader // 上传的文件
	ProcessedData []byte                // 预处理后的数据（如水印处理后的数据），优先级高于File
	UserID        uint                  // 用户ID
	FolderPath    string                // 文件夹路径
	FileName      string                // 文件名
	ContentType   string                // 内容类型
	Options       *UploadOptions        // 上传选项
}

// UploadOptions 上传选项
type UploadOptions struct {
	Quality       int  // 压缩质量 (1-100)
	MaxWidth      int  // 最大宽度
	MaxHeight     int  // 最大高度
	GenerateThumb bool // 是否生成缩略图
	ThumbWidth    int  // 缩略图最大宽度
	ThumbHeight   int  // 缩略图最大高度
	ThumbQuality  int  // 缩略图质量 (1-100)
	Compress      bool // 是否压缩
	WebPEnabled   bool // 是否启用WebP转换
}

// UploadResult 上传结果
type UploadResult struct {
	OriginalPath              string // 原始文件路径
	ThumbnailPath             string // 缩略图路径
	URL                       string // 相对URL
	ThumbnailURL              string // 缩略图相对URL
	FullURL                   string // 完整URL
	FullThumbURL              string // 完整缩略图URL
	RemoteURL                 string // 远程URL（用于云存储）
	RemoteThumbURL            string // 远程缩略图URL
	Size                      int64  // 文件大小
	Width                     int    // 文件宽度
	Height                    int    // 文件高度
	Hash                      string // 文件哈希
	ContentType               string // 内容类型
	Format                    string // 文件格式
	ThumbnailGenerationFailed bool   // 缩略图生成是否失败
	ThumbnailFailureReason    string // 缩略图失败原因
}

// URLOptions URL选项
type URLOptions struct {
	IsThumbnail  bool   // 是否缩略图
	CustomDomain string // 自定义域名
	UseCDN       bool   // 是否使用CDN
	Quality      int    // 文件质量
	Width        int    // 宽度
	Height       int    // 高度
	Expires      int64  // 过期时间(签名URL)
	ForceHTTPS   bool   // 强制HTTPS
}

// Capabilities 存储能力
type Capabilities struct {
	SupportsSignedURL bool     // 支持签名URL
	SupportsCDN       bool     // 支持CDN
	SupportsResize    bool     // 支持在线缩放
	SupportsWebP      bool     // 支持WebP转换
	MaxFileSize       int64    // 最大文件大小
	SupportedFormats  []string // 支持的格式
}

// Config 配置接口
type Config interface {
	Get(key string) interface{}
	GetString(key string) string
	GetInt(key string) int
	GetInt64(key string) int64
	GetBool(key string) bool
	GetFloat64(key string) float64
	Set(key string, value interface{})
	Has(key string) bool
}

// AdapterFactory 适配器工厂函数
type AdapterFactory func() StorageAdapter

// ErrorType 错误类型
type ErrorType string

const (
	ErrorTypeNotFound      ErrorType = "not_found"
	ErrorTypePermission    ErrorType = "permission"
	ErrorTypeQuotaExceeded ErrorType = "quota_exceeded"
	ErrorTypeInvalidFormat ErrorType = "invalid_format"
	ErrorTypeNetwork       ErrorType = "network"
	ErrorTypeInternal      ErrorType = "internal"
)

// StorageError 存储错误
type StorageError struct {
	Type    ErrorType
	Message string
	Cause   error
}

func (e *StorageError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e *StorageError) Unwrap() error {
	return e.Cause
}

func NewStorageError(errType ErrorType, message string, cause error) *StorageError {
	return &StorageError{
		Type:    errType,
		Message: message,
		Cause:   cause,
	}
}

// IsNotFoundError 检查是否为文件不存在错误
func IsNotFoundError(err error) bool {
	if storageErr, ok := err.(*StorageError); ok {
		return storageErr.Type == ErrorTypeNotFound
	}
	return false
}

// IsPermissionError 检查是否为权限错误
func IsPermissionError(err error) bool {
	if storageErr, ok := err.(*StorageError); ok {
		return storageErr.Type == ErrorTypePermission
	}
	return false
}

// IsQuotaExceededError 检查是否为配额超限错误
func IsQuotaExceededError(err error) bool {
	if storageErr, ok := err.(*StorageError); ok {
		return storageErr.Type == ErrorTypeQuotaExceeded
	}
	return false
}
