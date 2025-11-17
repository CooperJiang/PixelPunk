package models

import (
	"fmt"
	"path/filepath"
	"pixelpunk/pkg/common"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/* File 文件模型 */
type File struct {
	ID        string          `gorm:"primarykey;size:32" json:"id"`
	CreatedAt common.JSONTime `gorm:"index:idx_file_created_at" json:"created_at"`
	UpdatedAt common.JSONTime `json:"updated_at"`

	UserID       uint   `gorm:"not null;index:idx_file_user_id" json:"user_id"`
	FolderID     string `gorm:"size:32;index:idx_file_folder_id" json:"folder_id"`
	OriginalName string `gorm:"size:255;not null" json:"original_name"`
	DisplayName  string `gorm:"size:255" json:"display_name"` // 显示名称，可能包含时间戳等

	FileName       string `gorm:"size:255;not null" json:"file_name"`
	FilePath       string `gorm:"size:255;not null" json:"file_path"`
	FullPath       string `gorm:"size:255" json:"full_path"`
	LocalFilePath  string `gorm:"size:255" json:"local_file_path"`
	LocalThumbPath string `gorm:"size:255" json:"local_thumb_path"`
	URL            string `gorm:"size:255" json:"url"`
	ThumbURL       string `gorm:"size:255" json:"thumb_url"`
	RemoteURL      string `gorm:"size:255" json:"remote_url"`
	RemoteThumbURL string `gorm:"size:255" json:"remote_thumb_url"`

	ShortURL string `gorm:"size:32;index:idx_file_short_url" json:"short_url"`

	MD5Hash       string  `gorm:"size:32;index:idx_file_md5_hash" json:"md5_hash"`
	Size          int64   `gorm:"not null" json:"size"`
	SizeFormatted string  `gorm:"size:20" json:"size_formatted"`
	Width         int     `json:"width"`  // 文件/视频专用
	Height        int     `json:"height"` // 文件/视频专用
	Ratio         float64 `json:"ratio"`
	Format        string  `gorm:"size:10;not null" json:"format"`
	Mime          string  `gorm:"size:50" json:"mime"`
	Resolution    string  `gorm:"size:20" json:"resolution"`

	FileType string `gorm:"size:20;not null;default:'image';index:idx_file_type" json:"file_type"` // image,video,document,archive,audio,other
	MimeType string `gorm:"size:100" json:"mime_type"`

	Description       string `gorm:"type:text" json:"description"`
	NSFW              bool   `gorm:"default:false" json:"nsfw"`
	Status            string `gorm:"size:20;not null;default:active" json:"status"`
	AccessLevel       string `gorm:"size:20;not null;default:private" json:"access_level"`
	AccessKey         string `gorm:"size:32" json:"access_key,omitempty"`
	IsDuplicate       bool   `gorm:"default:false" json:"is_duplicate"`
	OriginalFileID    string `gorm:"size:32" json:"-"`
	IsRecommended     bool   `gorm:"default:false" json:"is_recommended"`
	APIKeyID          string `gorm:"size:32" json:"api_key_id"`
	StorageProviderID string `gorm:"size:36" json:"storage_provider_id"`
	StorageType       string `gorm:"size:20;not null;default:local" json:"storage_type"`

	AITaggingStatus      string     `gorm:"size:20;not null;default:none" json:"ai_tagging_status"`
	AITaggingTries       int        `gorm:"default:0" json:"ai_tagging_tries"`
	AITaggingDuration    int64      `gorm:"default:0" json:"ai_tagging_duration"`      // 总耗时（毫秒）
	AIHttpDuration       int64      `gorm:"default:0" json:"ai_http_duration"`         // HTTP调用耗时（毫秒）
	AILastHeartbeatAt    *time.Time `gorm:"index:idx_file_ai_heartbeat" json:"ai_last_heartbeat_at"`
	AIProcessingWorkerID string     `gorm:"size:64" json:"ai_processing_worker_id"`

	CategoryID     *uint  `gorm:"null;index:idx_file_category_id" json:"category_id"`
	CategorySource string `gorm:"size:20;default:manual" json:"category_source"`

	ThumbnailGenerationFailed bool   `gorm:"default:false" json:"thumbnail_generation_failed"`
	ThumbnailFailureReason    string `gorm:"size:255" json:"thumbnail_failure_reason,omitempty"`

	ExpiresAt              *time.Time `gorm:"column:expires_at;index" json:"expires_at"`
	StorageDuration        string     `gorm:"column:storage_duration;size:20;default:'permanent'" json:"storage_duration"`
	IsGuestUpload          bool       `gorm:"column:is_guest_upload;default:false" json:"is_guest_upload"`
	GuestFingerprint       string     `gorm:"column:guest_fingerprint;size:64;index" json:"guest_fingerprint"`
	GuestIP                string     `gorm:"column:guest_ip;size:45;index" json:"guest_ip"`
	ExpiryNotificationSent bool       `gorm:"column:expiry_notification_sent;default:false" json:"expiry_notification_sent"`

	SortOrder int `gorm:"default:0" json:"sort_order"`

	User     *User         `gorm:"foreignKey:UserID;references:ID" json:"user"`
	AIInfo   *FileAIInfo   `gorm:"foreignKey:FileID;references:ID" json:"ai_info"`
	Category *FileCategory `gorm:"foreignKey:CategoryID;references:ID" json:"category"`
}

const (
	FileTypeImage    = "image"
	FileTypeVideo    = "video"
	FileTypeDocument = "document"
	FileTypeArchive  = "archive"
	FileTypeAudio    = "audio"
	FileTypeOther    = "other"
)

func (File) TableName() string {
	return "file"
}

func (f *File) IsPublic() bool {
	return f.AccessLevel == "public"
}

func (f *File) IsPrivate() bool {
	return f.AccessLevel == "private"
}

func (f *File) IsProtected() bool {
	return f.AccessLevel == "protected"
}

func (f *File) HasLinkAccess() bool {
	return false
}

func (f *File) IsImage() bool {
	return f.FileType == FileTypeImage
}

func (f *File) IsVideo() bool {
	return f.FileType == FileTypeVideo
}

func (f *File) IsDocument() bool {
	return f.FileType == FileTypeDocument
}

func (f *File) IsArchive() bool {
	return f.FileType == FileTypeArchive
}

func (f *File) IsAudio() bool {
	return f.FileType == FileTypeAudio
}

func (f *File) HasDimensions() bool {
	return f.IsImage() || f.IsVideo()
}

func (f *File) BeforeCreate(tx *gorm.DB) error {
	if f.AccessLevel == "" {
		f.AccessLevel = "private"
	}

	if f.Status == "" {
		f.Status = "active"
	}

	if f.DisplayName == "" {
		f.DisplayName = f.OriginalName
	}

	if f.AITaggingStatus == "" {
		f.AITaggingStatus = common.AITaggingStatusNone
	}

	if f.StorageDuration == "" {
		f.StorageDuration = "permanent"
	}

	if f.FileType == "" {
		f.FileType = f.DetectFileType()
	}

	if f.SortOrder == 0 {
		var maxOrder int
		query := tx.Model(&File{}).Where("user_id = ?", f.UserID)
		if f.FolderID != "" {
			query = query.Where("folder_id = ?", f.FolderID)
		} else {
			query = query.Where("folder_id = '' OR folder_id IS NULL")
		}
		query.Select("COALESCE(MAX(sort_order), 0)").Scan(&maxOrder)
		f.SortOrder = maxOrder + 1
	}

	if f.ShortURL == "" {
		f.ShortURL = f.GenerateUniqueShortURL(tx)
	}

	return nil
}

func (f *File) BeforeSave(tx *gorm.DB) error {
	if f.StorageDuration != "permanent" && f.StorageDuration != "" {
		if f.ExpiresAt == nil {
			expiresAt := f.CalculateExpiryTime()
			f.ExpiresAt = &expiresAt
		}
	} else if f.StorageDuration == "permanent" {
		f.ExpiresAt = nil
	}
	return nil
}

func (f *File) DetectFileType() string {
	if f.OriginalName == "" && f.FileName == "" {
		return FileTypeOther
	}

	filename := f.OriginalName
	if filename == "" {
		filename = f.FileName
	}

	return DetectFileTypeFromFilename(filename)
}

func DetectFileTypeFromFilename(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	ext = strings.TrimPrefix(ext, ".")

	switch ext {
	case "jpg", "jpeg", "png", "gif", "bmp", "svg", "webp", "ico", "tiff", "tif":
		return FileTypeImage
	case "mp4", "avi", "mov", "wmv", "flv", "webm", "mkv", "m4v", "3gp", "ogv":
		return FileTypeVideo
	case "pdf", "doc", "docx", "xls", "xlsx", "ppt", "pptx", "txt", "rtf", "odt", "ods", "odp":
		return FileTypeDocument
	case "zip", "rar", "7z", "tar", "gz", "bz2", "xz", "cab", "iso":
		return FileTypeArchive
	case "mp3", "wav", "flac", "aac", "ogg", "wma", "m4a", "opus":
		return FileTypeAudio
	default:
		return FileTypeOther
	}
}

func (f *File) CalculateExpiryTime() time.Time {
	return common.CalculateExpiryTime(f.StorageDuration)
}

func (f *File) ParseStorageDuration() time.Duration {
	return common.ParseStorageDuration(f.StorageDuration)
}

func (f *File) IsExpired() bool {
	if f.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*f.ExpiresAt)
}

/* IsExpiringSoon 检查文件是否即将过期（默认24小时内） */
func (f *File) IsExpiringSoon() bool {
	if f.ExpiresAt == nil {
		return false
	}
	return time.Until(*f.ExpiresAt) < 24*time.Hour && time.Until(*f.ExpiresAt) > 0
}

func (f *File) GetRemainingTime() time.Duration {
	if f.ExpiresAt == nil {
		return -1 // 返回-1表示永久存储
	}
	remaining := time.Until(*f.ExpiresAt)
	if remaining < 0 {
		return 0 // 已过期
	}
	return remaining
}

func (f *File) GetRemainingDays() int {
	remaining := f.GetRemainingTime()
	if remaining < 0 {
		return -1 // 永久存储
	}
	return int(remaining.Hours() / 24)
}

func (f *File) IsTimeLimitedStorage() bool {
	return f.StorageDuration != "permanent" && f.StorageDuration != ""
}

func (f *File) CanExtendStorage() bool {
	return f.IsTimeLimitedStorage() && !f.IsExpired()
}

func (f *File) ExtendStorage(newDuration string) error {
	if !f.CanExtendStorage() {
		return fmt.Errorf("无法延长存储时间")
	}

	f.StorageDuration = newDuration
	expiresAt := f.CalculateExpiryTime()
	f.ExpiresAt = &expiresAt

	return nil
}

func (f *File) GenerateShortURL() string {
	if f.ShortURL != "" {
		return f.ShortURL
	}
	return common.GenerateBase62ShortURL()
}

/* GenerateUniqueShortURL 生成唯一的短链标识符，包含冲突检测 */
func (f *File) GenerateUniqueShortURL(tx *gorm.DB) string {
	if f.ShortURL != "" {
		return f.ShortURL
	}

	const maxRetries = 5
	for attempt := 0; attempt < maxRetries; attempt++ {
		shortURL := common.GenerateBase62ShortURL()

		var existingCount int64
		if err := tx.Model(&File{}).Where("short_url = ?", shortURL).Count(&existingCount).Error; err != nil {
			continue
		}

		if existingCount == 0 {
			return shortURL
		}

	}

	uuidStr := strings.ReplaceAll(uuid.New().String(), "-", "")
	return uuidStr[:7] // 取UUID的前7位作为短链
}

func (f *File) GetShortLinkURL() string {
	shortURL := f.GenerateShortURL()
	return fmt.Sprintf("/s/%s", shortURL)
}
