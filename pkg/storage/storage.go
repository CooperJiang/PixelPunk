package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"pixelpunk/pkg/storage/adapter"
	"pixelpunk/pkg/storage/factory"
	"pixelpunk/pkg/storage/manager"

	"github.com/google/uuid"
)

// Storage 存储服务统一入口
// 提供高级的存储操作接口，封装底层的适配器和管理器复杂性
type Storage struct {
	manager *manager.StorageManager
}

// New 创建存储服务
func New(channelRepo manager.ChannelRepository) *Storage {
	globalFactory := factory.GetGlobalFactory()
	mgr := manager.NewStorageManagerWithFactory(channelRepo, globalFactory)

	return &Storage{manager: mgr}
}

// NewWithManager 使用指定管理器创建存储服务
func NewWithManager(mgr *manager.StorageManager) *Storage {
	return &Storage{manager: mgr}
}

// 注意：这个函数用于兼容性，实际使用中应该注入正确的渠道仓库
func NewGlobalStorage() *Storage {
	channelRepo := &CompatChannelRepository{}
	globalFactory := factory.GetGlobalFactory()
	mgr := manager.NewStorageManagerWithFactory(channelRepo, globalFactory)

	return &Storage{manager: mgr}
}

// UploadRequest 上传请求
type UploadRequest struct {
	File          *multipart.FileHeader // 上传的文件
	ProcessedData []byte                // 预处理后的数据（如水印处理后的数据），优先级高于File
	ChannelID     string                // 存储渠道ID（可选，为空时使用最佳渠道）
	UserID        uint                  // 用户ID
	FolderPath    string                // 文件夹路径
	FileName      string                // 文件名（可选，为空时自动生成）
	ContentType   string                // 内容类型
	Quality       int                   // 压缩质量 (1-100)
	MaxWidth      int                   // 最大宽度
	MaxHeight     int                   // 最大高度
	GenerateThumb bool                  // 是否生成缩略图
	ThumbWidth    int                   // 缩略图最大宽度
	ThumbHeight   int                   // 缩略图最大高度
	ThumbQuality  int                   // 缩略图质量 (1-100)
	Compress      bool                  // 是否压缩
	WebPEnabled   bool                  // 是否启用WebP转换
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
	ChannelID                 string // 实际使用的存储渠道ID
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

// Upload 上传文件
func (s *Storage) Upload(ctx context.Context, req *UploadRequest) (*UploadResult, error) {
	// 生成文件名（如果未提供）
	fileName := req.FileName
	if fileName == "" {
		fileName = generateUniqueFileName(req.File.Filename)
	}

	adapterReq := &adapter.UploadRequest{
		File:          req.File,
		ProcessedData: req.ProcessedData, // 传递预处理后的数据
		UserID:        req.UserID,
		FolderPath:    req.FolderPath,
		FileName:      fileName,
		ContentType:   req.ContentType,
		Options: &adapter.UploadOptions{
			Quality:       req.Quality,
			MaxWidth:      req.MaxWidth,
			MaxHeight:     req.MaxHeight,
			GenerateThumb: req.GenerateThumb,
			ThumbWidth:    req.ThumbWidth,
			ThumbHeight:   req.ThumbHeight,
			ThumbQuality:  req.ThumbQuality,
			Compress:      req.Compress,
			WebPEnabled:   req.WebPEnabled,
		},
	}

	var result *adapter.UploadResult
	var err error
	var channelID string

	if req.ChannelID != "" {
		result, err = s.manager.Upload(ctx, req.ChannelID, adapterReq)
		channelID = req.ChannelID
	} else {
		result, err = s.manager.UploadWithBest(ctx, adapterReq)
		channelID, _ = s.manager.GetDefaultChannelID()
	}

	if err != nil {
		return nil, err
	}

	return &UploadResult{
		OriginalPath:   result.OriginalPath,
		ThumbnailPath:  result.ThumbnailPath,
		URL:            result.URL,
		ThumbnailURL:   result.ThumbnailURL,
		FullURL:        result.FullURL,
		FullThumbURL:   result.FullThumbURL,
		RemoteURL:      result.RemoteURL,
		RemoteThumbURL: result.RemoteThumbURL,
		Size:           result.Size,
		Width:          result.Width,
		Height:         result.Height,
		Hash:           result.Hash,
		ContentType:    result.ContentType,
		Format:         result.Format,
		ChannelID:      channelID,
	}, nil
}

// UploadWithDefault 使用默认渠道上传
func (s *Storage) UploadWithDefault(ctx context.Context, req *UploadRequest) (*UploadResult, error) {
	req.ChannelID = "" // 确保使用默认渠道
	return s.Upload(ctx, req)
}

// Delete 删除文件
func (s *Storage) Delete(ctx context.Context, channelID, path string) error {
	return s.manager.Delete(ctx, channelID, path)
}

func (s *Storage) GetURL(channelID, path string, options *URLOptions) (string, error) {
	adapterOptions := &adapter.URLOptions{
		IsThumbnail:  options.IsThumbnail,
		CustomDomain: options.CustomDomain,
		UseCDN:       options.UseCDN,
		Quality:      options.Quality,
		Width:        options.Width,
		Height:       options.Height,
		Expires:      options.Expires,
		ForceHTTPS:   options.ForceHTTPS,
	}

	return s.manager.GetURL(channelID, path, adapterOptions)
}

// HealthCheck 健康检查
func (s *Storage) HealthCheck(ctx context.Context, channelID string) error {
	return s.manager.HealthCheck(ctx, channelID)
}

// HealthCheckAll 检查所有渠道健康状态
func (s *Storage) HealthCheckAll(ctx context.Context) map[string]error {
	return s.manager.HealthCheckAll(ctx)
}

func (s *Storage) GetCapabilities(channelID string) (adapter.Capabilities, error) {
	return s.manager.GetAdapterCapabilities(channelID)
}

// RefreshChannel 刷新指定渠道
func (s *Storage) RefreshChannel(channelID string) error {
	return s.manager.RefreshAdapter(channelID)
}

func (s *Storage) GetManager() *manager.StorageManager {
	return s.manager
}

// BatchUpload 批量上传文件
func (s *Storage) BatchUpload(ctx context.Context, requests []*UploadRequest) ([]*UploadResult, []error) {
	var results []*UploadResult
	var errors []error

	for _, req := range requests {
		result, err := s.Upload(ctx, req)
		results = append(results, result)
		errors = append(errors, err)
	}

	return results, errors
}

// 辅助函数

// generateUniqueFileName 生成唯一文件名
func generateUniqueFileName(originalName string) string {
	ext := filepath.Ext(originalName)

	// 生成36位长度的唯一名称：UUID(32位) + 时间戳(4位)
	uuidStr := strings.ReplaceAll(uuid.New().String(), "-", "")
	timestamp := fmt.Sprintf("%04d", time.Now().Unix()%10000) // 取时间戳后4位
	uniqueName := uuidStr + timestamp

	return uniqueName + ext
}

func GetSupportedTypes() []string {
	return factory.GetGlobalSupportedTypes()
}

// RegisterStorageAdapter 注册存储适配器
func RegisterStorageAdapter(storageType string, adapterFactory adapter.AdapterFactory) {
	factory.RegisterGlobalAdapter(storageType, adapterFactory)
}

// GetBase64 获取文件的Base64编码
func (s *Storage) GetBase64(ctx context.Context, channelID, path string) (string, error) {
	return s.manager.GetBase64(ctx, channelID, path)
}

// GetThumbnailBase64 获取缩略图的Base64编码
func (s *Storage) GetThumbnailBase64(ctx context.Context, channelID, path string) (string, error) {
	return s.manager.GetThumbnailBase64(ctx, channelID, path)
}

// ReadFile 读取文件内容
func (s *Storage) ReadFile(ctx context.Context, channelID, path string) (io.ReadCloser, error) {
	return s.manager.ReadFile(ctx, channelID, path)
}
