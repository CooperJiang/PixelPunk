package file

import (
	"mime/multipart"
	"net"
	"path/filepath"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/setting"
	"pixelpunk/internal/services/stats"
	"pixelpunk/pkg/common"
	pkgStorage "pixelpunk/pkg/storage"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/* UploadResult 上传结果结构体，包含文件上传的所有信息 */
type UploadResult struct {
	URL                       string // 访问文件的URL
	LocalUrlPath              string // 本地存储路径
	ThumbUrl                  string // 缩略图URL
	LocalThumbPath            string // 本地缩略图路径
	RemoteUrl                 string // 远程URL
	RemoteThumbUrl            string // 远程缩略图URL
	Width                     int    // 文件宽度
	Height                    int    // 文件高度
	ThumbnailGenerationFailed bool   // 缩略图生成是否失败
	ThumbnailFailureReason    string // 缩略图失败原因
}

/* FileUploadResponse 文件上传响应结构体 */
type FileUploadResponse struct {
	ID           string `json:"id"`
	URL          string `json:"url"`
	FullURL      string `json:"full_url"`       // 完整URL
	FullThumbURL string `json:"full_thumb_url"` // 完整缩略图URL

	ShortURL     string `json:"short_url"`      // 短链标识符
	FileIDURL    string `json:"file_id_url"`    // 基于imageID的访问URL
	ThumbIDURL   string `json:"thumb_id_url"`   // 基于imageID的缩略图URL
	ShortLinkURL string `json:"short_link_url"` // 短链URL

	Width             int                    `json:"width"`
	Height            int                    `json:"height"`
	Size              int64                  `json:"size"`
	SizeFormatted     string                 `json:"size_formatted"`
	Ratio             float64                `json:"ratio"`
	Format            string                 `json:"format"`
	OriginalName      string                 `json:"original_name"`
	DisplayName       string                 `json:"display_name"`
	Description       string                 `json:"description"`
	AccessLevel       string                 `json:"access_level"`
	AccessKey         string                 `json:"access_key,omitempty"`
	IsDuplicate       bool                   `json:"is_duplicate"`
	OriginalFileID    string                 `json:"original_file_id,omitempty"`
	StorageProvider   string                 `json:"storage_provider"`
	StorageProviderID string                 `json:"storage_provider_id"` // 存储渠道ID
	StorageChannel    *models.StorageChannel `json:"storage_channel"`
	MD5Hash           string                 `json:"md5_hash"`       // MD5哈希
	IsRecommended     bool                   `json:"is_recommended"` // 是否推荐
	CreatedAt         *time.Time             `json:"created_at"`     // 创建时间
	UpdatedAt         *time.Time             `json:"updated_at"`     // 更新时间

	StorageDuration string     `json:"storage_duration"`
	ExpiresAt       *time.Time `json:"expires_at"`
	IsGuestUpload   bool       `json:"is_guest_upload"`

	ThumbnailGenerationFailed bool   `json:"thumbnail_generation_failed"`
	ThumbnailFailureReason    string `json:"thumbnail_failure_reason,omitempty"`

	WatermarkApplied       bool   `json:"watermark_applied"`
	WatermarkFailureReason string `json:"watermark_failure_reason,omitempty"`
}

/* CompressOptions 图像压缩选项，用于替代file包中的版本 */
type CompressOptions struct {
	MaxWidth  int  // 最大宽度
	MaxHeight int  // 最大高度
	Quality   int  // 压缩质量 (1-100)，仅对JPEG有效
	Preserve  bool // 是否保持原始宽高比
}

/* UploadContext 上传上下文，包含上传过程中需要共享的状态 */
type UploadContext struct {
	UserID   uint                  // 用户ID，0表示游客上传
	File     *multipart.FileHeader // 上传的文件
	Context  *gin.Context          // Gin上下文
	Optimize bool                  // 是否优化文件

	FolderID   string // 文件夹ID
	FolderPath string // 文件夹相对路径

	FileExt           string // 文件扩展名
	FileHash          string // 文件MD5哈希
	IsDuplicate       bool   // 是否重复文件
	OriginalFileID    string // 原始文件ID（重复文件时有值）
	ReuseExistingFile bool   // 是否复用现有文件

	StorageChannel *models.StorageChannel // 存储渠道

	FileID           string           // 生成的文件ID
	FileSize         int64            // 文件大小（新）
	FileFormat       string           // 文件格式（新）
	ActualChannelID  string           // 实际使用的渠道ID（新）
	OriginalName     string           // 原始文件名
	SafeOriginalName string           // 安全处理后的原始文件名
	DisplayName      string           // 显示名称
	AccessLevel      string           // 访问级别
	AccessKey        string           // 访问密钥（当access_level=protected时）
	CompressOptions  *CompressOptions // 压缩选项

	Tx *gorm.DB // 数据库事务

	Result       *UploadResult // 上传结果
	ExistingFile *models.File  // 现有文件（重复文件时）
	SavedFile    *models.File  // 保存到数据库的完整文件记录（包含创建时间等）

	ThumbnailBase64 string // 缩略图的base64数据，用于AI分析

	StorageDuration string     // 存储时长：3d/7d/30d/permanent
	ExpiresAt       *time.Time // 过期时间（自动计算）

	IsGuestUpload    bool   // 是否为游客上传
	GuestFingerprint string // 游客指纹
	GuestIP          string // 游客IP地址
	GuestUserAgent   string // 游客User-Agent

	WatermarkEnabled       bool        // 是否启用水印
	WatermarkConfig        string      // 水印配置JSON字符串
	WatermarkWrapper       interface{} // 水印处理后的文件包装器（内部使用）
	WatermarkApplied       bool        // 水印是否成功应用
	WatermarkFailureReason string      // 水印失败原因
	OriginalFileData       []byte      // 原始文件数据（一次性读取，供多次使用）

	EXIFData  *models.FileEXIF // 提取的 EXIF 元数据
	FileModel *models.File     // 文件模型（用于后续操作）
}

/* CreateUploadContext 创建一个新的上传上下文 */
func CreateUploadContext(c *gin.Context, userID uint, file *multipart.FileHeader, folderID, accessLevel string, optimize bool) *UploadContext {
	return CreateUploadContextWithDuration(c, userID, file, folderID, accessLevel, optimize, "")
}

/* CreateUploadContextWithDuration 创建带存储时长的上传上下文 */
func CreateUploadContextWithDuration(c *gin.Context, userID uint, file *multipart.FileHeader, folderID, accessLevel string, optimize bool, storageDuration string) *UploadContext {
	ctx := &UploadContext{
		UserID:          userID,
		File:            file,
		Context:         c,
		FolderID:        folderID,
		AccessLevel:     accessLevel,
		Optimize:        optimize,
		IsDuplicate:     false,
		StorageDuration: storageDuration,
		IsGuestUpload:   userID == 0, // 用户ID为0表示游客
	}

	if ctx.IsGuestUpload {
		ctx.GuestIP = getClientIP(c)
		ctx.GuestUserAgent = c.GetHeader("User-Agent")
	}

	if storageDuration != "" && storageDuration != "permanent" {
		expiresAt := common.CalculateExpiryTime(storageDuration)
		ctx.ExpiresAt = &expiresAt
	}

	if optimize {
		ctx.CompressOptions = createCompressOptions()
	}

	return ctx
}

func createCompressOptions() *CompressOptions {
	settingsMap, err := setting.GetSettingsByGroupAsMap("upload")
	maxWidth, maxHeight, quality := 600, 600, 85

	if err == nil {
		if widthVal, ok := settingsMap.Settings["thumbnail_max_width"]; ok {
			if width, ok := widthVal.(float64); ok {
				maxWidth = int(width)
			}
		}
		if heightVal, ok := settingsMap.Settings["thumbnail_max_height"]; ok {
			if height, ok := heightVal.(float64); ok {
				maxHeight = int(height)
			}
		}
		if qualityVal, ok := settingsMap.Settings["thumbnail_quality"]; ok {
			if q, ok := qualityVal.(float64); ok {
				quality = int(q)
			}
		}
	}

	return &CompressOptions{
		MaxWidth:  maxWidth,
		MaxHeight: maxHeight,
		Quality:   quality,
		Preserve:  true,
	}
}

/* CreateInstantUploadContext 为秒传创建虚拟的上传上下文 */
func CreateInstantUploadContext(c *gin.Context, userID uint, originalFile *models.File, fileName string, fileSize int64, folderID, accessLevel string, optimize bool) *UploadContext {
	virtualFile := &multipart.FileHeader{
		Filename: fileName,
		Size:     fileSize,
		Header:   make(map[string][]string),
	}

	if originalFile.Mime != "" {
		virtualFile.Header.Set("Content-Type", originalFile.Mime)
	}

	ctx := &UploadContext{
		UserID:            userID,
		File:              virtualFile,
		Context:           c,
		FolderID:          folderID,
		AccessLevel:       accessLevel,
		Optimize:          optimize,
		IsDuplicate:       true,                 // 秒传必然是重复文件
		OriginalFileID:    originalFile.ID,      // 设置原始文件ID
		ReuseExistingFile: true,                 // 标记为复用现有文件
		ExistingFile:      originalFile,         // 保存原文件信息
		FileHash:          originalFile.MD5Hash, // 使用原文件的MD5
		FileSize:          fileSize,
		FileFormat:        originalFile.Format,
		ActualChannelID:   originalFile.StorageProviderID,
	}

	if originalFile.StorageProviderID != "" {
		ctx.StorageChannel = &models.StorageChannel{
			ID:   originalFile.StorageProviderID,
			Type: originalFile.StorageType,
		}
	}

	ctx.Result = &UploadResult{
		URL:                       originalFile.URL,
		LocalUrlPath:              originalFile.LocalFilePath,
		ThumbUrl:                  originalFile.ThumbURL,
		LocalThumbPath:            originalFile.LocalThumbPath,
		RemoteUrl:                 originalFile.RemoteURL,
		RemoteThumbUrl:            originalFile.RemoteThumbURL,
		Width:                     originalFile.Width,
		Height:                    originalFile.Height,
		ThumbnailGenerationFailed: originalFile.ThumbnailGenerationFailed,
		ThumbnailFailureReason:    originalFile.ThumbnailFailureReason,
	}

	ctx.OriginalName = strings.TrimSuffix(fileName, filepath.Ext(fileName))
	ctx.SafeOriginalName = pkgStorage.SanitizeFileName(ctx.OriginalName)
	ctx.DisplayName = ctx.SafeOriginalName
	ctx.FileExt = filepath.Ext(fileName)

	ctx.FileID = generateFileID()

	if optimize {
		ctx.CompressOptions = createCompressOptions()
	}

	return ctx
}

/* CreateInstantUploadContextWithDuration 为秒传创建虚拟的上传上下文（支持存储时长） */
func CreateInstantUploadContextWithDuration(c *gin.Context, userID uint, originalFile *models.File, fileName string, fileSize int64, folderID, accessLevel string, optimize bool, storageDuration string) *UploadContext {
	ctx := CreateInstantUploadContext(c, userID, originalFile, fileName, fileSize, folderID, accessLevel, optimize)

	ctx.StorageDuration = storageDuration

	if storageDuration != "" && storageDuration != "permanent" {
		expiresAt := common.CalculateExpiryTime(storageDuration)
		ctx.ExpiresAt = &expiresAt
	}

	return ctx
}

/* StatsEvent 统计事件结构体，用于异步更新统计 */
type StatsEvent struct {
	Type   string // 事件类型
	UserID uint   // 用户ID
	FileID string // 文件ID
	Size   int64  // 文件大小
}

var statsChannel = make(chan StatsEvent, 100)

func init() {
	go handleStatsEvents()
}

func handleStatsEvents() {
	for event := range statsChannel {
		switch event.Type {
		case "file_created":
			stats.GetStatsAdapter().RecordFileCreated(event.Size)
		}
	}
}

func getClientIP(c *gin.Context) string {
	if ip := c.GetHeader("X-Forwarded-For"); ip != "" {
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}

	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		return ip
	}

	if ip := c.GetHeader("CF-Connecting-IP"); ip != "" {
		return ip
	}

	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}

	return ip
}
