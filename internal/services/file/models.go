package file

import (
	imodels "pixelpunk/internal/models"
	"pixelpunk/pkg/common"
)

/* AIInfoResponse AI信息响应 */
type AIInfoResponse struct {
	Description    string             `json:"description"`
	Tags           []string           `json:"tags"`
	DominantColor  string             `json:"dominant_color"`
	Resolution     string             `json:"resolution"`
	IsNSFW         bool               `json:"is_nsfw"`
	NSFWScore      float64            `json:"nsfw_score"`
	NSFWEvaluation string             `json:"nsfw_evaluation"`
	ColorPalette   []string           `json:"color_palette,omitempty"`   // 颜色调色板
	AspectRatio    float64            `json:"aspect_ratio,omitempty"`    // 宽高比
	Composition    string             `json:"composition,omitempty"`     // 构图类型
	ObjectsCount   int                `json:"objects_count,omitempty"`   // 识别到的物体数量
	NSFWCategories map[string]float64 `json:"nsfw_categories,omitempty"` // NSFW分类评分
}

type FileDetailResponse struct {
	ID                string            `json:"id"`
	FullURL           string            `json:"full_url"`       // 完整的文件URL（包含域名）
	FullThumbURL      string            `json:"full_thumb_url"` // 完整的缩略图URL（包含域名）
	OriginalName      string            `json:"original_name"`
	DisplayName       string            `json:"display_name"` // 显示名称
	Size              int64             `json:"size"`
	Width             int               `json:"width"`
	Height            int               `json:"height"`
	Format            string            `json:"format"`
	AccessLevel       string            `json:"access_level"`
	FolderID          string            `json:"folder_id,omitempty"`
	CreatedAt         common.JSONTime   `json:"created_at"`
	UpdatedAt         common.JSONTime   `json:"updated_at"`
	Views             int64             `json:"views,omitempty"`
	IsDuplicate       bool              `json:"is_duplicate,omitempty"`        // 是否是重复文件
	MD5Hash           string            `json:"md5_hash,omitempty"`            // MD5哈希值
	IsRecommended     bool              `json:"is_recommended"`                // 是否推荐
	StorageProviderID string            `json:"storage_provider_id,omitempty"` // 存储提供者ID
	AIInfo            *AIInfoResponse   `json:"ai_info,omitempty"`
	EXIFInfo          *imodels.FileEXIF `json:"exif_info,omitempty"` // EXIF 元数据
}

/* ExternalAPIFileResponse 外部API文件响应结构（简化版） */
type ExternalAPIFileResponse struct {
	ID           string          `json:"id"`
	URL          string          `json:"url"`       // 完整的文件URL（包含域名）
	ThumbURL     string          `json:"thumb_url"` // 完整的缩略图URL（包含域名）
	OriginalName string          `json:"original_name"`
	Size         int64           `json:"size"`
	Width        int             `json:"width"`
	Height       int             `json:"height"`
	Format       string          `json:"format"`
	AccessLevel  string          `json:"access_level"`
	CreatedAt    common.JSONTime `json:"created_at"`
}

/* AdminFileDetailResponse 管理员文件详情响应 */
type AdminFileDetailResponse struct {
	ID                string          `json:"id"`
	URL               string          `json:"url"`
	ThumbnailURL      string          `json:"thumb_url"`
	FullURL           string          `json:"full_url"`       // 完整的文件URL（包含域名）
	FullThumbURL      string          `json:"full_thumb_url"` // 完整的缩略图URL（包含域名）
	ShortURL          string          `json:"short_url"`      // 短链URL
	OriginalName      string          `json:"original_name"`
	DisplayName       string          `json:"display_name"` // 显示名称
	Size              int64           `json:"size"`
	Width             int             `json:"width"`
	Height            int             `json:"height"`
	Format            string          `json:"format"`
	AccessLevel       string          `json:"access_level"`
	FolderID          string          `json:"folder_id,omitempty"`
	CreatedAt         common.JSONTime `json:"created_at"`
	UpdatedAt         common.JSONTime `json:"updated_at"`
	Views             int64           `json:"views,omitempty"`
	IsDuplicate       bool            `json:"is_duplicate"`
	MD5Hash           string          `json:"md5_hash,omitempty"`
	IsRecommended     bool            `json:"is_recommended"`
	StorageProviderID string          `json:"storage_provider_id,omitempty"`

	UserID          uint             `json:"user_id,omitempty"`
	UserName        string           `json:"user_name,omitempty"`
	UserInfo        interface{}      `json:"user_info,omitempty"` // 用户详细信息
	AIInfo          *AIInfoResponse  `json:"ai_info,omitempty"`
	StorageDuration string           `json:"storage_duration"` // 存储时长：3d/7d/30d/permanent
	ExpiresAt       *common.JSONTime `json:"expires_at"`       // 过期时间
	IsTimeLimited   bool             `json:"is_time_limited"`  // 是否为限时存储
}

/* TagWithCount 带有计数的标签结构 */
type TagWithCount struct {
	imodels.GlobalTag
	Count int64 `json:"count"`
}

/* UserDetailInfo 用户详细信息（用于创作者卡片） */
type UserDetailInfo struct {
	ID          uint            `json:"id"`
	Username    string          `json:"username"`
	Avatar      string          `json:"avatar,omitempty"`
	CreatedAt   common.JSONTime `json:"created_at"`
	TotalImages int64           `json:"total_files"` // 公开文件总数
	TotalViews  int64           `json:"total_views"` // 所有文件总浏览量
	DaysJoined  int             `json:"days_joined"` // 注册天数
	OtherImages []FileThumbnail `json:"other_files"` // 其他作品（最多10张）
	TopTags     []string        `json:"top_tags"`    // 热门标签（最多5个）
}

/* FileThumbnail 文件缩略图信息 */
type FileThumbnail struct {
	ID           string `json:"id"`
	DisplayName  string `json:"display_name"`
	FullURL      string `json:"full_url"`
	FullThumbURL string `json:"full_thumb_url"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	IsNSFW       bool   `json:"is_nsfw,omitempty"`
}
