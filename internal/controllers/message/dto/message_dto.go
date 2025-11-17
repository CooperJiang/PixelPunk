package dto

import (
	"encoding/json"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"

	"github.com/gin-gonic/gin"
)

type GetMessagesDTO struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"pageSize,default=20" binding:"min=1,max=100"`
	Status   int    `form:"status" binding:"omitempty,oneof=1 2"` // 1:未读 2:已读
	Type     string `form:"type" binding:"omitempty,max=100"`     // 消息类型筛选
}

func (d *GetMessagesDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Page.min":     "页码不能小于1",
		"PageSize.min": "每页数量不能小于1",
		"PageSize.max": "每页数量不能大于100",
		"Status.oneof": "状态只能是1(未读)或2(已读)",
		"Type.max":     "消息类型长度不能超过100",
	}
}

type MarkMessageReadDTO struct {
	MessageID string `uri:"id" binding:"required"`
}

func (d *MarkMessageReadDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"MessageID.required": "消息ID不能为空",
	}
}

type DeleteMessageDTO struct {
	MessageID string `uri:"id" binding:"required"`
}

func (d *DeleteMessageDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"MessageID.required": "消息ID不能为空",
	}
}

type CreateTemplateDTO struct {
	Type               string `json:"type" binding:"required,max=100"`
	Title              string `json:"title" binding:"required,max=200"`
	Content            string `json:"content" binding:"required"`
	Description        string `json:"description" binding:"omitempty,max=500"`
	IsEnabled          bool   `json:"is_enabled"`
	SendEmail          bool   `json:"send_email"`
	ShowToast          bool   `json:"show_toast"`
	ToastType          string `json:"toast_type" binding:"omitempty,oneof=success error info warning"`
	DefaultActionType  string `json:"default_action_type" binding:"omitempty,max=50"`
	DefaultActionText  string `json:"default_action_text" binding:"omitempty,max=100"`
	DefaultActionStyle string `json:"default_action_style" binding:"omitempty,oneof=primary secondary danger warning"`
	ActionURLTemplate  string `json:"action_url_template" binding:"omitempty,max=500"`
}

func (d *CreateTemplateDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Type.required":            "模板类型不能为空",
		"Type.max":                 "模板类型长度不能超过100",
		"Title.required":           "模板标题不能为空",
		"Title.max":                "模板标题长度不能超过200",
		"Content.required":         "模板内容不能为空",
		"Description.max":          "模板描述长度不能超过500",
		"ToastType.oneof":          "Toast类型只能是success、error、info、warning之一",
		"DefaultActionType.max":    "操作类型长度不能超过50",
		"DefaultActionText.max":    "操作文本长度不能超过100",
		"DefaultActionStyle.oneof": "操作样式只能是primary、secondary、danger、warning之一",
		"ActionURLTemplate.max":    "URL模板长度不能超过500",
	}
}

type UpdateTemplateDTO struct {
	ID                 uint   `uri:"id" binding:"required"`
	Title              string `json:"title" binding:"omitempty,max=200"`
	Content            string `json:"content" binding:"omitempty"`
	Description        string `json:"description" binding:"omitempty,max=500"`
	IsEnabled          *bool  `json:"is_enabled"`
	SendEmail          *bool  `json:"send_email"`
	ShowToast          *bool  `json:"show_toast"`
	ToastType          string `json:"toast_type" binding:"omitempty,oneof=success error info warning"`
	DefaultActionType  string `json:"default_action_type" binding:"omitempty,max=50"`
	DefaultActionText  string `json:"default_action_text" binding:"omitempty,max=100"`
	DefaultActionStyle string `json:"default_action_style" binding:"omitempty,oneof=primary secondary danger warning"`
	ActionURLTemplate  string `json:"action_url_template" binding:"omitempty,max=500"`
}

func (d *UpdateTemplateDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"ID.required":              "模板ID不能为空",
		"Title.max":                "模板标题长度不能超过200",
		"Description.max":          "模板描述长度不能超过500",
		"ToastType.oneof":          "Toast类型只能是success、error、info、warning之一",
		"DefaultActionType.max":    "操作类型长度不能超过50",
		"DefaultActionText.max":    "操作文本长度不能超过100",
		"DefaultActionStyle.oneof": "操作样式只能是primary、secondary、danger、warning之一",
		"ActionURLTemplate.max":    "URL模板长度不能超过500",
	}
}

type DeleteTemplateDTO struct {
	ID uint `uri:"id" binding:"required"`
}

func (d *DeleteTemplateDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"ID.required": "模板ID不能为空",
	}
}

type ToggleTemplateDTO struct {
	ID      uint `uri:"id" binding:"required"`
	Enabled bool `json:"enabled"`
}

func (d *ToggleTemplateDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"ID.required": "模板ID不能为空",
	}
}

type SendMessageDTO struct {
	UserID             uint                   `json:"user_id" binding:"required"`
	Type               string                 `json:"type" binding:"required,max=100"`
	Title              string                 `json:"title" binding:"required,max=200"`
	Content            string                 `json:"content" binding:"required"`
	Priority           int                    `json:"priority" binding:"omitempty,oneof=1 2 3"`
	ExpiresAt          *common.JSONTime       `json:"expires_at"`
	RelatedType        string                 `json:"related_type" binding:"omitempty,max=50"`
	RelatedID          string                 `json:"related_id" binding:"omitempty,max=100"`
	RelatedData        map[string]interface{} `json:"related_data"`
	IsActionable       bool                   `json:"is_actionable"`
	ActionType         string                 `json:"action_type" binding:"omitempty,max=50"`
	ActionURL          string                 `json:"action_url" binding:"omitempty,max=500"`
	ActionText         string                 `json:"action_text" binding:"omitempty,max=100"`
	ActionStyle        string                 `json:"action_style" binding:"omitempty,oneof=primary secondary danger warning"`
	RequiresPermission string                 `json:"requires_permission" binding:"omitempty,max=100"`
	MetaData           map[string]interface{} `json:"meta_data"`
}

func (d *SendMessageDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"UserID.required":        "用户ID不能为空",
		"Type.required":          "消息类型不能为空",
		"Type.max":               "消息类型长度不能超过100",
		"Title.required":         "消息标题不能为空",
		"Title.max":              "消息标题长度不能超过200",
		"Content.required":       "消息内容不能为空",
		"Priority.oneof":         "优先级只能是1(高)、2(中)、3(低)",
		"RelatedType.max":        "关联类型长度不能超过50",
		"RelatedID.max":          "关联ID长度不能超过100",
		"ActionType.max":         "操作类型长度不能超过50",
		"ActionURL.max":          "操作URL长度不能超过500",
		"ActionText.max":         "操作文本长度不能超过100",
		"ActionStyle.oneof":      "操作样式只能是primary、secondary、danger、warning之一",
		"RequiresPermission.max": "权限要求长度不能超过100",
	}
}

type GetTemplateDTO struct {
	ID uint `uri:"id" binding:"required"`
}

func (d *GetTemplateDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"ID.required": "模板ID不能为空",
	}
}

func MessageToResponse(msg *models.Message) gin.H {
	response := gin.H{
		"id":                  msg.ID,
		"created_at":          msg.CreatedAt,
		"updated_at":          msg.UpdatedAt,
		"user_id":             msg.UserID,
		"type":                msg.Type,
		"status":              msg.Status,
		"priority":            msg.Priority,
		"read_at":             msg.ReadAt,
		"expires_at":          msg.ExpiresAt,
		"related_type":        msg.RelatedType,
		"related_id":          msg.RelatedID,
		"is_actionable":       msg.IsActionable,
		"action_type":         msg.ActionType,
		"action_url":          msg.ActionURL,
		"action_text":         msg.ActionText,
		"action_style":        msg.ActionStyle,
		"requires_permission": msg.RequiresPermission,
	}

	var data map[string]interface{}
	if msg.Data != "" {
		if err := json.Unmarshal([]byte(msg.Data), &data); err != nil {
			data = make(map[string]interface{})
		}
	} else {
		data = make(map[string]interface{})
	}
	response["data"] = data

	var relatedData map[string]interface{}
	if msg.RelatedData != "" {
		if err := json.Unmarshal([]byte(msg.RelatedData), &relatedData); err != nil {
			relatedData = make(map[string]interface{})
		}
	} else {
		relatedData = make(map[string]interface{})
	}
	response["related_data"] = relatedData

	var metaData map[string]interface{}
	if msg.MetaData != "" {
		if err := json.Unmarshal([]byte(msg.MetaData), &metaData); err != nil {
			metaData = make(map[string]interface{})
		}
	} else {
		metaData = make(map[string]interface{})
	}
	response["meta_data"] = metaData

	return response
}
