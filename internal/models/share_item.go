package models

import (
	"pixelpunk/pkg/common"
)

/* ShareItem 分享项目表 */
type ShareItem struct {
	ID      string `gorm:"primarykey;size:32" json:"id"`
	ShareID string `gorm:"size:32;not null;index" json:"share_id"`

	ItemType string `gorm:"size:10;not null"` // folder或image
	ItemID   string `gorm:"size:32;not null"` // 文件夹ID或文件ID

	SortOrder int `gorm:"default:0" json:"sort_order"` // 排序顺序

	CreatedAt common.JSONTime `json:"created_at"` // 创建时间
}

func (ShareItem) TableName() string {
	return "share_item"
}
