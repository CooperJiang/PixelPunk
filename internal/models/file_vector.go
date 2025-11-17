package models

import (
	"pixelpunk/pkg/common"
	"time"

	"gorm.io/gorm"
)

/* FileVector 文件向量元数据模型（向量数据存储在Qdrant中） */
type FileVector struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	FileID      string `gorm:"size:32;not null;uniqueIndex" json:"file_id"`
	Dimension   int    `gorm:"not null;default:1536" json:"dimension"`                          // 向量维度
	Model       string `gorm:"size:100;not null;default:'text-embedding-3-small'" json:"model"` // 向量化模型名称
	Description string `gorm:"type:text" json:"description"`                                    // 原始描述文本
	Status      string `gorm:"size:20;not null;default:pending" json:"status"`                  // 向量状态：pending/completed/failed/reset/stale

	EmbeddingSchemaVersion int    `gorm:"not null;default:1" json:"embedding_schema_version"` // 嵌入Schema版本（配合系统当前版本判断是否过期）
	DescChecksum           string `gorm:"size:64;index" json:"desc_checksum"`                 // 描述哈希（变更即置stale）

	DeletedAt *common.JSONTime `gorm:"type:timestamp;index" json:"deleted_at,omitempty"`

	RetryCount         int              `gorm:"default:0" json:"retry_count"`                  // 重试次数
	ErrorMessage       string           `gorm:"type:text" json:"error_message,omitempty"`      // 错误信息
	LastRetryAt        *common.JSONTime `gorm:"type:timestamp" json:"last_retry_at,omitempty"` // 最后重试时间
	ProcessingDuration int              `gorm:"default:0" json:"processing_duration"`          // 处理耗时(秒)

	ActualStatus      string           `gorm:"type:varchar(20);default:unknown" json:"actual_status"` // 实际向量存在状态
	LastVerified      *common.JSONTime `gorm:"type:timestamp" json:"last_verified,omitempty"`         // 最后验证时间
	VerificationError string           `gorm:"type:text" json:"verification_error,omitempty"`         // 验证错误信息

	CreatedAt common.JSONTime `json:"created_at"`
	UpdatedAt common.JSONTime `json:"updated_at"`

	File File `gorm:"foreignKey:FileID;references:ID" json:"-"`
}

func (FileVector) TableName() string {
	return "file_vector"
}

func (fv *FileVector) SetVectorMetadata(dimension int, model string) {
	fv.Dimension = dimension
	if model != "" {
		fv.Model = model
	}
	fv.Status = common.VectorStatusCompleted
}

func (fv *FileVector) IsCompleted() bool {
	return fv.Status == common.VectorStatusCompleted
}

func (fv *FileVector) IsPending() bool {
	return fv.Status == common.VectorStatusPending
}

func (fv *FileVector) IsFailed() bool {
	return fv.Status == common.VectorStatusFailed
}

func (fv *FileVector) IsReset() bool {
	return fv.Status == common.VectorStatusReset
}

func (fv *FileVector) IsProcessing() bool {
	return fv.Status == "processing"
}

func (fv *FileVector) CanRetry() bool {
	return fv.IsFailed() && fv.RetryCount < 3 // 最多重试3次
}

func (fv *FileVector) IncrementRetry() {
	fv.RetryCount++
	now := common.JSONTime(time.Now())
	fv.LastRetryAt = &now
}

func (fv *FileVector) SetError(errorMsg string) {
	fv.Status = common.VectorStatusFailed
	fv.ErrorMessage = errorMsg
}

const (
	FileVectorActualStatusUnknown  = "unknown"
	FileVectorActualStatusVerified = "verified"
	FileVectorActualStatusMissing  = "missing"
)

/* IsActualVerified 检查向量是否已验证存在 */
func (fv *FileVector) IsActualVerified() bool {
	return fv.ActualStatus == FileVectorActualStatusVerified
}

/* IsActualMissing 检查向量是否确实缺失 */
func (fv *FileVector) IsActualMissing() bool {
	return fv.ActualStatus == FileVectorActualStatusMissing
}

/* IsActualUnknown 检查向量状态是否未知 */
func (fv *FileVector) IsActualUnknown() bool {
	return fv.ActualStatus == FileVectorActualStatusUnknown
}

func (fv *FileVector) SetVerificationStatus(status string, errorMsg string) {
	fv.ActualStatus = status
	now := common.JSONTime(time.Now())
	fv.LastVerified = &now

	if errorMsg != "" {
		fv.VerificationError = errorMsg
	} else {
		fv.VerificationError = ""
	}
}

/* NeedsVerification 检查是否需要验证 */
func (fv *FileVector) NeedsVerification() bool {
	if !fv.IsCompleted() {
		return false
	}

	if fv.IsActualUnknown() {
		return true
	}

	if fv.LastVerified != nil {
		duration := time.Since(time.Time(*fv.LastVerified))
		return duration > 24*time.Hour
	}

	return true
}

func (fv *FileVector) BeforeCreate(tx *gorm.DB) error {
	if fv.Status == "" {
		fv.Status = common.VectorStatusPending
	}

	if fv.Dimension == 0 {
		fv.Dimension = 1536 // text-embedding-3-small默认维度
	}

	return nil
}
