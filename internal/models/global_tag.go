package models

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"pixelpunk/pkg/common"
	"strings"

	"gorm.io/gorm"
)

/* GlobalTag 全局标签模型 - 新的标签引用架构核心表 */
type GlobalTag struct {
	ID          uint            `gorm:"primarykey" json:"id"`
	Name        string          `gorm:"size:50;not null;uniqueIndex:idx_global_tag_name" json:"name"`
	Slug        string          `gorm:"size:50;not null;uniqueIndex:idx_global_tag_slug" json:"slug"`
	Description string          `gorm:"type:text" json:"description"`
	IsSystem    bool            `gorm:"default:false;index:idx_global_tag_system" json:"is_system"`
	CreatorID   uint            `gorm:"not null;index:idx_global_tag_creator" json:"creator_id"`
	UsageCount  int             `gorm:"default:0;index:idx_global_tag_usage" json:"usage_count"` // 全局使用次数统计
	SortOrder   int             `gorm:"default:0" json:"sort_order"`
	CreatedAt   common.JSONTime `json:"created_at"`
	UpdatedAt   common.JSONTime `json:"updated_at"`

	Creator User `gorm:"foreignKey:CreatorID;references:ID" json:"creator,omitempty"`
}

func (GlobalTag) TableName() string {
	return "global_tag"
}

func (t *GlobalTag) BeforeCreate(tx *gorm.DB) error {
	if t.Slug == "" {
		t.Slug = generateSlug(t.Name)
	}

	var count int64
	for i := 0; i < 100; i++ { // 最多尝试100次
		slug := t.Slug
		if i > 0 {
			slug = fmt.Sprintf("%s-%d", t.Slug, i)
		}

		err := tx.Model(&GlobalTag{}).Where("slug = ?", slug).Count(&count).Error
		if err != nil {
			return err
		}

		if count == 0 {
			t.Slug = slug
			break
		}
	}

	return nil
}

/* UserTagReference 用户标签引用表 - 用户与全局标签的多对多关系 */
type UserTagReference struct {
	ID        uint            `gorm:"primarykey" json:"id"`
	UserID    uint            `gorm:"not null;index:idx_user_tag_ref_user;uniqueIndex:idx_user_tag_unique" json:"user_id"`
	TagID     uint            `gorm:"not null;index:idx_user_tag_ref_tag;uniqueIndex:idx_user_tag_unique" json:"tag_id"`
	Source    string          `gorm:"size:20;not null;default:manual" json:"source"` // manual(手动创建), ai(AI生成), system(系统导入)
	CreatedAt common.JSONTime `json:"created_at"`

	User User      `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Tag  GlobalTag `gorm:"foreignKey:TagID;references:ID" json:"tag,omitempty"`
}

func (UserTagReference) TableName() string {
	return "user_tag_reference"
}

func (r *UserTagReference) BeforeCreate(tx *gorm.DB) error {
	if r.Source == "" {
		r.Source = "manual"
	}

	return nil
}

/* TagCategoryRelation 标签分类关联表 - 标签与分类的多对多关系 */
type TagCategoryRelation struct {
	ID         uint            `gorm:"primarykey" json:"id"`
	TagID      uint            `gorm:"not null;index:idx_tag_category_tag;uniqueIndex:idx_tag_category_unique" json:"tag_id"`
	CategoryID uint            `gorm:"not null;index:idx_tag_category_category;uniqueIndex:idx_tag_category_unique" json:"category_id"`
	Weight     int             `gorm:"default:1" json:"weight"`      // 权重，用于排序
	IsCore     bool            `gorm:"default:false" json:"is_core"` // 是否为核心标签
	CreatedAt  common.JSONTime `json:"created_at"`

	Tag      GlobalTag    `gorm:"foreignKey:TagID;references:ID" json:"tag,omitempty"`
	Category FileCategory `gorm:"foreignKey:CategoryID;references:ID" json:"category,omitempty"`
}

func (TagCategoryRelation) TableName() string {
	return "tag_category_relation"
}

/* FileGlobalTagRelation 文件全局标签关联表 - 替换原有的ImageTagRelation */
type FileGlobalTagRelation struct {
	ID          uint            `gorm:"primarykey" json:"id"`
	FileID      string          `gorm:"size:32;not null;index:idx_img_global_tag_file" json:"file_id"`
	TagID       uint            `gorm:"not null;index:idx_img_global_tag_tag" json:"tag_id"`
	UserID      uint            `gorm:"not null;index:idx_img_global_tag_user" json:"user_id"`                // 冗余字段，便于查询
	AccessLevel string          `gorm:"size:20;not null;index:idx_img_global_tag_access" json:"access_level"` // 冗余字段，便于查询
	Source      string          `gorm:"size:20;not null;default:manual" json:"source"`                        // manual(手动), ai(AI生成), import(导入)
	Confidence  float64         `gorm:"default:1.0000" json:"confidence"`                                     // AI标注的置信度
	CreatedAt   common.JSONTime `json:"created_at"`

	Tag  GlobalTag `gorm:"foreignKey:TagID;references:ID" json:"tag,omitempty"`
	File File      `gorm:"foreignKey:FileID;references:ID" json:"file,omitempty"`
	User User      `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}

func (FileGlobalTagRelation) TableName() string {
	return "file_global_tag_relation"
}

func (r *FileGlobalTagRelation) BeforeCreate(tx *gorm.DB) error {
	if r.UserID == 0 || r.AccessLevel == "" {
		var file File
		err := tx.Where("id = ?", r.FileID).First(&file).Error
		if err != nil {
			return fmt.Errorf("获取文件信息失败: %v", err)
		}

		if r.UserID == 0 {
			r.UserID = file.UserID
		}
		if r.AccessLevel == "" {
			r.AccessLevel = file.AccessLevel
		}
	}

	if r.Source == "" {
		r.Source = "manual"
	}

	if r.Confidence == 0 {
		r.Confidence = 1.0
	}

	return nil
}

/* GlobalTagOperationLog 全局标签操作日志 */
type GlobalTagOperationLog struct {
	ID            uint            `gorm:"primarykey" json:"id"`
	OperationType string          `gorm:"size:50;not null;index:idx_global_tag_log_type" json:"operation_type"`
	TagID         *uint           `gorm:"index:idx_global_tag_log_tag" json:"tag_id"`
	FileID        *string         `gorm:"size:32;index:idx_global_tag_log_file" json:"file_id"`
	UserID        *uint           `gorm:"index:idx_global_tag_log_user" json:"user_id"`
	OperatorID    *uint           `gorm:"index:idx_global_tag_log_operator" json:"operator_id"`
	OldValue      *string         `gorm:"type:json" json:"old_value"`
	NewValue      *string         `gorm:"type:json" json:"new_value"`
	Source        string          `gorm:"size:20;not null" json:"source"`
	Status        string          `gorm:"size:20;not null;default:success" json:"status"`
	ErrorMessage  *string         `gorm:"type:text" json:"error_message"`
	CreatedAt     common.JSONTime `gorm:"index:idx_global_tag_log_created" json:"created_at"`
}

func (GlobalTagOperationLog) TableName() string {
	return "global_tag_operation_log"
}

func (l *GlobalTagOperationLog) BeforeCreate(tx *gorm.DB) error {
	if l.Status == "" {
		l.Status = "success"
	}

	if l.Source == "" {
		l.Source = "api"
	}

	return nil
}

/* GlobalTagStatsCache 全局标签统计缓存表 */
type GlobalTagStatsCache struct {
	CacheKey    string          `gorm:"primarykey;size:100" json:"cache_key"`
	TagID       *uint           `gorm:"index:idx_global_cache_tag" json:"tag_id"`
	UserID      *uint           `gorm:"index:idx_global_cache_user" json:"user_id"`
	CategoryID  *uint           `gorm:"index:idx_global_cache_category" json:"category_id"`
	AccessLevel *string         `gorm:"size:20" json:"access_level"`
	CountValue  int64           `gorm:"not null;default:0" json:"count_value"`
	CacheType   string          `gorm:"size:50;not null;index:idx_global_cache_type" json:"cache_type"`
	UpdatedAt   common.JSONTime `json:"updated_at"`
	ExpiresAt   common.JSONTime `gorm:"index:idx_global_cache_expires" json:"expires_at"`
}

func (GlobalTagStatsCache) TableName() string {
	return "global_tag_stats_cache"
}

func (c *GlobalTagStatsCache) BeforeCreate(tx *gorm.DB) error {
	if c.CacheKey == "" {
		c.CacheKey = c.GenerateCacheKey()
	}

	return nil
}

func (c *GlobalTagStatsCache) BeforeUpdate(tx *gorm.DB) error {
	c.CacheKey = c.GenerateCacheKey()
	return nil
}

/* GenerateCacheKey 生成缓存键 */
func (c *GlobalTagStatsCache) GenerateCacheKey() string {
	data := fmt.Sprintf("%s:%v:%v:%v:%v", c.CacheType, c.TagID, c.UserID, c.CategoryID, c.AccessLevel)
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%s:%s", c.CacheType, hex.EncodeToString(hash[:8]))
}

func generateSlug(name string) string {
	slug := strings.ToLower(name)

	replacements := map[string]string{
		" ":  "-",
		"_":  "-",
		"/":  "-",
		"\\": "-",
		".":  "-",
		",":  "-",
		"(":  "",
		")":  "",
		"[":  "",
		"]":  "",
		"{":  "",
		"}":  "",
		"<":  "",
		">":  "",
		"&":  "and",
		"+":  "plus",
		"=":  "eq",
		"?":  "",
		"!":  "",
		"@":  "at",
		"#":  "",
		"$":  "",
		"%":  "",
		"^":  "",
		"*":  "",
		"|":  "",
		"\"": "",
		"'":  "",
		"`":  "",
		"~":  "",
	}

	for old, new := range replacements {
		slug = strings.ReplaceAll(slug, old, new)
	}

	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	slug = strings.Trim(slug, "-")

	if slug == "" {
		slug = "unnamed"
	}

	if len(slug) > 50 {
		slug = slug[:50]
		slug = strings.Trim(slug, "-")
	}
	return slug
}
