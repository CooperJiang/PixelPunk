package models

import (
	"pixelpunk/pkg/common"

	"gorm.io/gorm"
)

/* StorageChannel 存储渠道模型 */
type StorageChannel struct {
	ID        string          `gorm:"primarykey;size:36" json:"id"`
	CreatedAt common.JSONTime `json:"created_at"`
	UpdatedAt common.JSONTime `json:"updated_at"`

	Name         string           `gorm:"size:50;not null;index:idx_storage_channel_name" json:"name"`
	Type         string           `gorm:"size:20;not null" json:"type"`
	Status       int8             `gorm:"not null;default:1" json:"status"`
	IsDefault    bool             `gorm:"default:false" json:"is_default"`
	IsLocal      bool             `gorm:"default:false" json:"is_local"`
	Remark       string           `gorm:"size:255" json:"remark"`
	FileCount    int64            `gorm:"-" json:"file_count"`
	CustomDomain string           `gorm:"-" json:"custom_domain"`
	Bucket       string           `gorm:"-" json:"bucket"`
	LastUploadAt *common.JSONTime `gorm:"-" json:"last_upload_at"`
}

func (s *StorageChannel) BeforeCreate(tx *gorm.DB) (err error) {
	if s.Type == "local" || s.IsLocal {
		s.IsLocal = true
		s.Status = 1
		s.IsDefault = true
	}

	if s.IsDefault {
		if err := tx.Model(&StorageChannel{}).
			Where("is_default = ?", true).
			Update("is_default", false).Error; err != nil {
			return err
		}
	}
	return nil
}

/* StorageConfigItem 存储配置项模型 */
type StorageConfigItem struct {
	ID        string          `gorm:"primarykey;size:36" json:"id"`
	CreatedAt common.JSONTime `json:"created_at"`
	UpdatedAt common.JSONTime `json:"updated_at"`

	ChannelID   string `gorm:"size:36;not null;index:idx_storage_config_channel_id" json:"channel_id"`
	Name        string `gorm:"size:50;not null" json:"name"`
	KeyName     string `gorm:"size:50;not null;index:idx_storage_config_key_name" json:"key_name"`
	Value       string `gorm:"type:text" json:"value,omitempty"`
	Type        string `gorm:"size:20;not null;default:string" json:"type"`
	IsSecret    bool   `gorm:"default:false" json:"is_secret"`
	Required    bool   `gorm:"default:true" json:"required"`
	Description string `gorm:"type:text" json:"description"`
}

type ConfigTemplate struct {
	Name        string   `json:"name"`
	KeyName     string   `json:"key_name"`
	Type        string   `json:"type"`
	IsSecret    bool     `json:"is_secret"`
	Required    bool     `json:"required"`
	Description string   `json:"description"`
	Options     []string `json:"options,omitempty"`
	Default     string   `json:"default,omitempty"`
}

func (s *StorageConfigItem) GetDecryptedValue() (string, error) {
	return s.Value, nil
}

/* GetChannelConfigWithDefaults 获取渠道配置并应用默认值 */
func GetChannelConfigWithDefaults(channelID string, configMap map[string]interface{}, channelType string) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range configMap {
		result[k] = v
	}

	if channelType != "local" {
		if _, exists := result["use_https"]; !exists {
			result["use_https"] = true
		}
		if _, exists := result["hide_remote_url"]; !exists {
			result["hide_remote_url"] = false
		}
		if _, exists := result["access_control"]; !exists {
			result["access_control"] = "public-read"
		}
	}

	return result
}
