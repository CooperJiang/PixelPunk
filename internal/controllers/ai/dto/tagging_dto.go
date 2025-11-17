package dto

import (
	"encoding/json"
	"pixelpunk/internal/models"

	"github.com/gin-gonic/gin"
)

type TaggingQueryDTO struct {
	Status  string `form:"status" binding:"omitempty"`
	Page    int    `form:"page" binding:"required,min=1"`
	Limit   int    `form:"limit" binding:"required,min=1,max=100"`
	OrderBy string `form:"order_by"`
	Order   string `form:"order" default:"desc"`
}

func (d *TaggingQueryDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Page.required":  "页码不能为空",
		"Page.min":       "页码必须大于等于1",
		"Limit.required": "每页数量不能为空",
		"Limit.min":      "每页数量必须大于等于1",
		"Limit.max":      "每页数量不能超过100",
	}
}

type ResetPendingDTO struct {
	TimeThresholdMinutes int `json:"time_threshold_minutes" binding:"required,min=5"`
}

func (d *ResetPendingDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"TimeThresholdMinutes.required": "时间阈值不能为空",
		"TimeThresholdMinutes.min":      "时间阈值必须大于等于5分钟",
	}
}

type IgnoreTaggingDTO struct {
	FileIDs []string `json:"file_ids" binding:"required,min=1,max=500"`
	Reason  string   `json:"reason"`
}

func (d *IgnoreTaggingDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"FileIDs.required": "必须提供file_ids",
		"FileIDs.min":      "至少选择1个文件",
		"FileIDs.max":      "一次最多忽略500个文件",
	}
}

type UnignoreTaggingDTO struct {
	FileIDs []string `json:"file_ids" binding:"required,min=1,max=500"`
}

func (d *UnignoreTaggingDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"FileIDs.required": "必须提供file_ids",
		"FileIDs.min":      "至少选择1个文件",
		"FileIDs.max":      "一次最多操作500个文件",
	}
}

type RetryFailedAllDTO struct {
	Limit int `json:"limit" binding:"omitempty,min=1,max=1000"`
}

func (d *RetryFailedAllDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Limit.min": "重试数量至少为1",
		"Limit.max": "单次最多重试1000条",
	}
}

type RetryTaggingDTO struct {
	FileIDs []string `json:"file_ids" binding:"required,min=1,max=500"`
}

func (d *RetryTaggingDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"FileIDs.required": "必须提供file_ids",
		"FileIDs.min":      "至少选择1个文件",
		"FileIDs.max":      "一次最多操作500个文件",
	}
}

type TriggerTaggingDTO struct {
	MaxFiles int `json:"max_files" binding:"omitempty,min=1,max=5000"`
}

func (d *TriggerTaggingDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"MaxFiles.min": "数量至少为1",
		"MaxFiles.max": "单次最多触发5000个文件",
	}
}

func TaggingLogToResponse(log *models.FileTaggingLog) gin.H {
	response := gin.H{
		"id":          log.ID,
		"file_id":     log.FileID,
		"file_type":   log.FileType,
		"status":      log.Status,
		"action":      log.Action,
		"type":        log.Type,
		"operator_id": log.OperatorID,
		"duration":    log.Duration,
		"created_at":  log.CreatedAt,
		"updated_at":  log.UpdatedAt,
	}

	var data map[string]interface{}
	if log.Data != "" {
		if err := json.Unmarshal([]byte(log.Data), &data); err != nil {
			data = make(map[string]interface{})
		}
	} else {
		data = make(map[string]interface{})
	}
	response["data"] = data

	return response
}
