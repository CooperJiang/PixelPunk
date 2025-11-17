package message

import (
	"pixelpunk/internal/controllers/message/dto"
	"pixelpunk/internal/middleware"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/message"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetUserMessages(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	req, err := common.ValidateRequest[dto.GetMessagesDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	messageService := message.GetMessageService()
	messages, total, err := messageService.GetUserMessages(userID, req.Page, req.PageSize, req.Status, req.Type)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	// 转换消息列表，解析 JSON 字符串字段为对象
	messageResponses := make([]gin.H, len(messages))
	for i, msg := range messages {
		messageResponses[i] = dto.MessageToResponse(&msg)
	}

	data := gin.H{
		"items": messageResponses,
		"pagination": gin.H{
			"total":        total,
			"size":         req.PageSize,
			"current_page": req.Page,
			"last_page":    (total + int64(req.PageSize) - 1) / int64(req.PageSize),
		},
	}

	errors.ResponseSuccess(c, data, "获取成功")
}

func MarkMessageRead(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	messageID := c.Param("id")
	if messageID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "消息ID不能为空"))
		return
	}

	messageService := message.GetMessageService()
	if err := messageService.MarkMessageRead(userID, messageID); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "标记成功")
}

func MarkAllMessagesRead(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	messageService := message.GetMessageService()
	if err := messageService.MarkAllMessagesRead(userID); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "标记成功")
}

func DeleteMessage(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	messageID := c.Param("id")
	if messageID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "消息ID不能为空"))
		return
	}

	messageService := message.GetMessageService()
	if err := messageService.DeleteMessage(userID, messageID); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "删除成功")
}

func GetUnreadCount(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	messageService := message.GetMessageService()
	count, err := messageService.GetUnreadCount(userID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, gin.H{
		"count": count,
	}, "获取成功")
}

func GetAllTemplates(c *gin.Context) {
	templateService := message.GetTemplateService()
	templates, err := templateService.GetAllTemplates()
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, gin.H{
		"templates": templates,
	}, "获取成功")
}

func GetTemplateByID(c *gin.Context) {
	req, err := common.ValidateRequest[dto.GetTemplateDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	var template models.MessageTemplate
	db := database.GetDB()
	if err := db.First(&template, req.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			errors.HandleError(c, errors.New(errors.CodeNotFound, "消息模板不存在"))
		} else {
			errors.HandleError(c, errors.Wrap(err, errors.CodeDBQueryFailed, "查询消息模板失败"))
		}
		return
	}

	errors.ResponseSuccess(c, gin.H{
		"template": template,
	}, "获取成功")
}

func CreateTemplate(c *gin.Context) {
	req, err := common.ValidateRequest[dto.CreateTemplateDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	template := &models.MessageTemplate{
		Type:               req.Type,
		Title:              req.Title,
		Content:            req.Content,
		Description:        req.Description,
		IsEnabled:          req.IsEnabled,
		SendEmail:          req.SendEmail,
		ShowToast:          req.ShowToast,
		ToastType:          req.ToastType,
		DefaultActionType:  req.DefaultActionType,
		DefaultActionText:  req.DefaultActionText,
		DefaultActionStyle: req.DefaultActionStyle,
		ActionURLTemplate:  req.ActionURLTemplate,
	}

	templateService := message.GetTemplateService()
	if err := templateService.CreateTemplate(template); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, gin.H{
		"template": template,
	}, "创建成功")
}

func UpdateTemplate(c *gin.Context) {
	req, err := common.ValidateRequest[dto.UpdateTemplateDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.IsEnabled != nil {
		updates["is_enabled"] = *req.IsEnabled
	}
	if req.SendEmail != nil {
		updates["send_email"] = *req.SendEmail
	}
	if req.ShowToast != nil {
		updates["show_toast"] = *req.ShowToast
	}
	if req.ToastType != "" {
		updates["toast_type"] = req.ToastType
	}
	if req.DefaultActionType != "" {
		updates["default_action_type"] = req.DefaultActionType
	}
	if req.DefaultActionText != "" {
		updates["default_action_text"] = req.DefaultActionText
	}
	if req.DefaultActionStyle != "" {
		updates["default_action_style"] = req.DefaultActionStyle
	}
	if req.ActionURLTemplate != "" {
		updates["action_url_template"] = req.ActionURLTemplate
	}

	templateService := message.GetTemplateService()
	if err := templateService.UpdateTemplate(req.ID, updates); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "更新成功")
}

func DeleteTemplate(c *gin.Context) {
	req, err := common.ValidateRequest[dto.DeleteTemplateDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	templateService := message.GetTemplateService()
	if err := templateService.DeleteTemplate(req.ID); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "删除成功")
}

func ToggleTemplate(c *gin.Context) {
	req, err := common.ValidateRequest[dto.ToggleTemplateDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	templateService := message.GetTemplateService()
	if err := templateService.ToggleTemplate(req.ID, req.Enabled); err != nil {
		errors.HandleError(c, err)
		return
	}

	status := "禁用"
	if req.Enabled {
		status = "启用"
	}

	errors.ResponseSuccess(c, nil, "模板已"+status)
}

func SendMessage(c *gin.Context) {
	req, err := common.ValidateRequest[dto.SendMessageDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	// 构建消息数据（管理员手动发送时，将 title 和 content 放入 data）
	data := map[string]interface{}{
		"title":   req.Title,
		"content": req.Content,
	}

	// 如果有额外的 RelatedData，合并进去
	if req.RelatedData != nil {
		for k, v := range req.RelatedData {
			data[k] = v
		}
	}

	options := &message.MessageOptions{
		Priority:           req.Priority,
		ExpiresAt:          req.ExpiresAt,
		RelatedType:        req.RelatedType,
		RelatedID:          req.RelatedID,
		RelatedData:        req.RelatedData,
		IsActionable:       req.IsActionable,
		ActionType:         req.ActionType,
		ActionURL:          req.ActionURL,
		ActionText:         req.ActionText,
		ActionStyle:        req.ActionStyle,
		RequiresPermission: req.RequiresPermission,
		MetaData:           req.MetaData,
	}

	messageService := message.GetMessageService()
	if err := messageService.SendMessage(req.UserID, req.Type, data, options); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "消息发送成功")
}

func GetMessageStats(c *gin.Context) {
	errors.ResponseSuccess(c, gin.H{
		"total_messages":  "待实现",
		"unread_messages": "待实现",
		"daily_messages":  "待实现",
		"template_count":  "待实现",
	}, "获取成功")
}
