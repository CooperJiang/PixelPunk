package models

import (
	"pixelpunk/pkg/common"
	"time"

	"gorm.io/gorm"
)

/* APIKey API密钥模型 */
type APIKey struct {
	ID        string          `gorm:"primarykey;size:32" json:"id"`
	CreatedAt common.JSONTime `json:"created_at"`
	UpdatedAt common.JSONTime `json:"updated_at"`

	UserID   uint   `gorm:"not null;index" json:"user_id"`
	Name     string `gorm:"size:100;not null" json:"name"`                                // 密钥名称/备注
	KeyValue string `gorm:"size:128;not null;uniqueIndex:idx_api_key_key_value" json:"-"` // 密钥值，不对外暴露
	Status   int    `gorm:"default:1;index" json:"status"`                                // 1:正常 2:禁用

	StorageLimit     int64 `gorm:"default:0" json:"storage_limit"`      // 存储容量限制(bytes)，0表示不限制
	StorageUsed      int64 `gorm:"default:0" json:"storage_used"`       // 已使用的存储容量(bytes)
	UploadCountLimit int   `gorm:"default:0" json:"upload_count_limit"` // 上传次数限制，0表示不限制
	UploadCountUsed  int   `gorm:"default:0" json:"upload_count_used"`  // 已使用的上传次数
	SingleFileLimit  int64 `gorm:"default:0" json:"single_file_limit"`  // 单文件大小限制(bytes)，0表示不限制

	AllowedTypes string `gorm:"size:255" json:"allowed_types"` // 允许的文件类型，如: "jpg,jpeg,png,gif"
	FolderID     string `gorm:"size:32" json:"folder_id"`      // 指定上传目录

	ExpiresAt  *common.JSONTime `json:"expires_at"`   // 过期时间，nil表示永不过期
	LastUsedAt *common.JSONTime `json:"last_used_at"` // 最后使用时间
}

/* APIKeyStatus API密钥状态常量 */
const (
	APIKeyStatusActive   = 1 // 正常状态
	APIKeyStatusDisabled = 2 // 禁用状态
)

func (APIKey) TableName() string {
	return "api_key"
}

func (k *APIKey) BeforeCreate(tx *gorm.DB) error {
	if k.Status == 0 {
		k.Status = APIKeyStatusActive
	}
	return nil
}

func (k *APIKey) IsActive() bool {
	return k.Status == APIKeyStatusActive
}

func (k *APIKey) IsDisabled() bool {
	return k.Status == APIKeyStatusDisabled
}

func (k *APIKey) IsExpired() bool {
	if k.ExpiresAt == nil {
		return false // 没有过期时间，表示永不过期
	}

	now := time.Now()
	expiryTime := time.Time(*k.ExpiresAt)
	return now.After(expiryTime)
}

func (k *APIKey) CheckStorageLimit(fileSize int64) bool {
	if k.StorageLimit <= 0 {
		return true // 不限制
	}

	return k.StorageUsed+fileSize <= k.StorageLimit
}

func (k *APIKey) CheckUploadCountLimit() bool {
	if k.UploadCountLimit <= 0 {
		return true // 不限制
	}

	return k.UploadCountUsed < k.UploadCountLimit
}

func (k *APIKey) CheckSingleFileLimit(fileSize int64) bool {
	if k.SingleFileLimit <= 0 {
		return true // 不限制
	}

	return fileSize <= k.SingleFileLimit
}

func (k *APIKey) CanUpload(fileSize int64) (bool, string) {
	if !k.IsActive() {
		return false, "API密钥已禁用"
	}

	if k.IsExpired() {
		return false, "API密钥已过期"
	}

	if !k.CheckUploadCountLimit() {
		return false, "已达到上传次数限制"
	}

	if !k.CheckStorageLimit(fileSize) {
		return false, "已达到存储容量限制"
	}

	if !k.CheckSingleFileLimit(fileSize) {
		return false, "文件大小超过单文件限制"
	}

	return true, ""
}
