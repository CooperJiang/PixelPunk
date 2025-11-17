package models

import (
	"pixelpunk/pkg/common"

	"gorm.io/gorm"
)

/* UploadChunk 分片上传分片信息 */
type UploadChunk struct {
	ID        uint            `gorm:"primarykey" json:"id"`
	CreatedAt common.JSONTime `json:"created_at"`
	UpdatedAt common.JSONTime `json:"updated_at"`

	SessionID   string `gorm:"size:32;not null;index:idx_upload_chunk_session_id;uniqueIndex:idx_session_chunk" json:"session_id"`
	ChunkNumber int    `gorm:"not null;uniqueIndex:idx_session_chunk" json:"chunk_number"` // 分片序号 (0-based)
	ChunkSize   int64  `gorm:"not null" json:"chunk_size"`                                 // 实际分片大小
	ChunkMD5    string `gorm:"size:32" json:"chunk_md5"`                                   // 分片MD5校验
	Status      string `gorm:"size:20;default:pending" json:"status"`                      // pending/uploaded/verified/failed
	StoragePath string `gorm:"size:255" json:"storage_path"`                               // 临时存储路径
	RetryCount  int    `gorm:"default:0" json:"retry_count"`                               // 重试次数
	ErrorMsg    string `gorm:"type:text" json:"error_msg"`                                 // 错误信息

}

func (UploadChunk) TableName() string {
	return "upload_chunk"
}

func (uc *UploadChunk) BeforeCreate(tx *gorm.DB) error {
	if uc.Status == "" {
		uc.Status = "pending"
	}

	return nil
}

/* IsUploaded 检查分片是否已上传 */
func (uc *UploadChunk) IsUploaded() bool {
	return uc.Status == "uploaded" || uc.Status == "verified"
}

/* IsPending 检查分片是否待上传 */
func (uc *UploadChunk) IsPending() bool {
	return uc.Status == "pending"
}

/* IsFailed 检查分片是否上传失败 */
func (uc *UploadChunk) IsFailed() bool {
	return uc.Status == "failed"
}

/* CanRetry 检查分片是否可以重试 */
func (uc *UploadChunk) CanRetry() bool {
	return uc.RetryCount < 3 // 最多重试3次
}
