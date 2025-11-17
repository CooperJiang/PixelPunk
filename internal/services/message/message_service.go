package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/cache"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/email"
	"pixelpunk/pkg/errors"
	"strconv"
	"strings"
	"text/template"
	"time"

	"gorm.io/gorm"
)

var messageService *MessageService

// MessageService 消息服务
type MessageService struct{}

// InitMessageService 初始化消息服务
func InitMessageService() {
	messageService = &MessageService{}
}

func GetMessageService() *MessageService {
	if messageService == nil {
		InitMessageService()
	}
	return messageService
}

// SendMessage 发送消息
func (s *MessageService) SendMessage(userID uint, msgType string, data map[string]interface{}, options *MessageOptions) error {

	db := database.GetDB()
	if db == nil {
		return errors.New(errors.CodeDBConnectionFailed, "数据库未初始化")
	}

	var dataJSON string
	if data != nil {
		if dataBytes, err := json.Marshal(data); err == nil {
			dataJSON = string(dataBytes)
		} else {
			dataJSON = "{}"
		}
	} else {
		dataJSON = "{}"
	}

	message := &models.Message{
		UserID: userID,
		Type:   msgType,
		Data:   dataJSON,
		Status: common.MessageStatusUnread,
	}

	if options != nil {
		message.Priority = options.Priority
		message.ExpiresAt = options.ExpiresAt
		message.RelatedType = options.RelatedType
		message.RelatedID = options.RelatedID
		message.IsActionable = options.IsActionable
		message.ActionType = options.ActionType
		message.ActionURL = options.ActionURL
		message.ActionText = options.ActionText
		message.ActionStyle = options.ActionStyle
		message.RequiresPermission = options.RequiresPermission

		// 序列化关联数据和元数据
		if options.RelatedData != nil {
			if data, err := json.Marshal(options.RelatedData); err == nil {
				message.RelatedData = string(data)
			} else {
			}
		}
		if options.MetaData != nil {
			if data, err := json.Marshal(options.MetaData); err == nil {
				message.MetaData = string(data)
			} else {
			}
		}
	} else {
	}

	result := db.Create(message)
	if result.Error != nil {
		return errors.Wrap(result.Error, errors.CodeDBCreateFailed, "创建消息失败")
	}

	s.clearUnreadCountCache(userID)

	return nil
}

// SendTemplateMessage 根据模板发送消息
func (s *MessageService) SendTemplateMessage(userID uint, templateType string, variables map[string]interface{}) error {

	// 尝试获取消息模板（可选）
	template, err := s.GetMessageTemplate(templateType)

	options := &MessageOptions{
		Priority: common.MessagePriorityNormal,
	}

	// 如果模板存在且启用，使用模板配置
	if err == nil && template.IsTemplateEnabled() {
		// 从模板配置构建交互选项
		if template.DefaultActionType != "" {
			options.IsActionable = true
			options.ActionType = template.DefaultActionType
			options.ActionText = template.DefaultActionText
			options.ActionStyle = template.DefaultActionStyle

			if template.ActionURLTemplate != "" {
				options.ActionURL = s.processTemplate(template.ActionURLTemplate, variables)
			}
		}

		if template.ShouldSendEmail() {
			title := s.processTemplate(template.Title, variables)
			content := s.processTemplate(template.Content, variables)
			go s.sendEmailNotification(userID, title, content)
		}
	} else {
		// 模板不存在或未启用，记录日志但继续发送消息
		if err != nil {
		} else {
		}
	}

	if relatedType, ok := variables["related_type"].(string); ok {
		options.RelatedType = relatedType
	}
	// RelatedID 可能是 uint 或 string，需要转换
	if relatedID, ok := variables["related_id"]; ok {
		switch v := relatedID.(type) {
		case uint:
			options.RelatedID = fmt.Sprintf("%d", v)
		case string:
			options.RelatedID = v
		default:
		}
	}

	options.RelatedData = variables
	options.MetaData = variables

	if err := s.SendMessage(userID, templateType, variables, options); err != nil {
		return err
	}

	return nil
}

func (s *MessageService) GetUserMessages(userID uint, page, pageSize, status int, msgType string) ([]models.Message, int64, error) {
	db := database.GetDB()

	query := db.Where("user_id = ? AND status != ?", userID, common.MessageStatusDeleted)

	if status > 0 {
		query = query.Where("status = ?", status)
	}

	if msgType != "" {
		if strings.Contains(msgType, ".") {
			query = query.Where("type = ?", msgType)
		} else {
			query = query.Where("type LIKE ?", msgType+".%")
		}
	}

	var total int64
	if err := query.Model(&models.Message{}).Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "查询消息总数失败")
	}

	// 分页查询，优先显示未读消息
	var messages []models.Message
	offset := (page - 1) * pageSize
	// 排序逻辑：先按status升序（1=未读在前，2=已读在后），再按创建时间倒序
	err := query.Order("status ASC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&messages).Error

	if err != nil {
		return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "查询消息列表失败")
	}

	return messages, total, nil
}

// MarkMessageRead 标记消息已读
func (s *MessageService) MarkMessageRead(userID uint, messageID string) error {
	db := database.GetDB()

	var message models.Message
	if err := db.Where("id = ? AND user_id = ?", messageID, userID).First(&message).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.CodeNotFound, "消息不存在")
		}
		return errors.Wrap(err, errors.CodeDBQueryFailed, "查询消息失败")
	}

	// 如果已经是已读状态，直接返回
	if message.IsRead() {
		return nil
	}

	now := common.JSONTime(time.Now())
	err := db.Model(&message).Updates(map[string]interface{}{
		"status":  common.MessageStatusRead,
		"read_at": &now,
	}).Error

	if err != nil {
		return errors.Wrap(err, errors.CodeDBUpdateFailed, "更新消息状态失败")
	}

	s.clearUnreadCountCache(userID)

	return nil
}

// MarkAllMessagesRead 批量标记用户所有未读消息为已读
func (s *MessageService) MarkAllMessagesRead(userID uint) error {
	db := database.GetDB()

	now := common.JSONTime(time.Now())
	err := db.Model(&models.Message{}).
		Where("user_id = ? AND status = ?", userID, common.MessageStatusUnread).
		Updates(map[string]interface{}{
			"status":  common.MessageStatusRead,
			"read_at": &now,
		}).Error

	if err != nil {
		return errors.Wrap(err, errors.CodeDBUpdateFailed, "批量标记已读失败")
	}

	s.clearUnreadCountCache(userID)

	return nil
}

func (s *MessageService) DeleteMessage(userID uint, messageID string) error {
	db := database.GetDB()

	var message models.Message
	if err := db.Where("id = ? AND user_id = ?", messageID, userID).First(&message).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.CodeNotFound, "消息不存在")
		}
		return errors.Wrap(err, errors.CodeDBQueryFailed, "查询消息失败")
	}

	err := db.Model(&message).Update("status", common.MessageStatusDeleted).Error
	if err != nil {
		return errors.Wrap(err, errors.CodeDBUpdateFailed, "删除消息失败")
	}

	// 如果删除的是未读消息，清除缓存
	if message.IsUnread() {
		s.clearUnreadCountCache(userID)
	}

	return nil
}

func (s *MessageService) GetUnreadCount(userID uint) (int64, error) {
	cacheKey := fmt.Sprintf("user:unread_messages:%d", userID)

	if countStr, err := cache.GetCache().Get(cacheKey); err == nil && countStr != "" {
		if count, err := strconv.ParseInt(countStr, 10, 64); err == nil {
			return count, nil
		}
	}

	// 缓存未命中，查询数据库
	count, err := s.countUnreadFromDB(userID)
	if err != nil {
		return 0, err
	}

	// 写入缓存，5分钟过期
	_ = cache.GetCache().Set(cacheKey, fmt.Sprintf("%d", count), 5*time.Minute)

	return count, nil
}

func (s *MessageService) GetMessageTemplate(templateType string) (*models.MessageTemplate, error) {
	db := database.GetDB()

	var template models.MessageTemplate
	err := db.Where("type = ?", templateType).First(&template).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.CodeNotFound, fmt.Sprintf("消息模板不存在: %s", templateType))
		}
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询消息模板失败")
	}

	return &template, nil
}

// processTemplate 处理模板变量替换
func (s *MessageService) processTemplate(templateStr string, variables map[string]interface{}) string {
	// 使用Go模板引擎处理条件语句和变量替换
	tmpl, err := template.New("message").Parse(templateStr)
	if err != nil {
		// 如果模板解析失败，回退到简单替换
		return s.processTemplateSimple(templateStr, variables)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, variables); err != nil {
		// 如果模板执行失败，回退到简单替换
		return s.processTemplateSimple(templateStr, variables)
	}

	return buf.String()
}

// processTemplateSimple 简单的变量替换（向后兼容）
func (s *MessageService) processTemplateSimple(template string, variables map[string]interface{}) string {
	result := template
	for key, value := range variables {
		placeholder := fmt.Sprintf("{{.%s}}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}
	return result
}

// countUnreadFromDB 从数据库查询未读消息数量
func (s *MessageService) countUnreadFromDB(userID uint) (int64, error) {
	db := database.GetDB()

	var count int64
	err := db.Model(&models.Message{}).
		Where("user_id = ? AND status = ?", userID, common.MessageStatusUnread).
		Count(&count).Error

	if err != nil {
		return 0, errors.Wrap(err, errors.CodeDBQueryFailed, "查询未读消息数量失败")
	}

	return count, nil
}

// clearUnreadCountCache 清除用户未读消息数量缓存
func (s *MessageService) clearUnreadCountCache(userID uint) {
	cacheKey := fmt.Sprintf("user:unread_messages:%d", userID)
	_ = cache.GetCache().Del(cacheKey)
}

// sendEmailNotification 发送邮件通知
func (s *MessageService) sendEmailNotification(userID uint, title, content string) {
	db := database.GetDB()
	var user models.User
	if err := db.Select("email").Where("id = ?", userID).First(&user).Error; err != nil {
		return
	}

	if !email.IsMailEnabled() {
		return
	}

	if err := email.SendMail(user.Email, title, content); err != nil {
	} else {
	}
}

// MessageOptions 消息选项
type MessageOptions struct {
	Priority           int                    `json:"priority"`
	ExpiresAt          *common.JSONTime       `json:"expires_at"`
	RelatedType        string                 `json:"related_type"`
	RelatedID          string                 `json:"related_id"`
	RelatedData        map[string]interface{} `json:"related_data"`
	IsActionable       bool                   `json:"is_actionable"`
	ActionType         string                 `json:"action_type"`
	ActionURL          string                 `json:"action_url"`
	ActionText         string                 `json:"action_text"`
	ActionStyle        string                 `json:"action_style"`
	RequiresPermission string                 `json:"requires_permission"`
	MetaData           map[string]interface{} `json:"meta_data"`
}

// SendExpiryNotifications 发送过期提醒
func (s *MessageService) SendExpiryNotifications() (int, int, error) {
	db := database.GetDB()

	// 查询即将过期（3天内）且未发送提醒的文件
	threeDaysLater := time.Now().Add(72 * time.Hour)

	// 使用join查询获取用户信息，避免表名硬编码问题
	var results []struct {
		FileID       string    `gorm:"column:file_id"`
		OriginalName string    `gorm:"column:original_name"`
		ExpiresAt    time.Time `gorm:"column:expires_at"`
		UserID       uint      `gorm:"column:user_id"`
		Username     string    `gorm:"column:username"`
		Email        string    `gorm:"column:email"`
	}

	// 使用子查询的方式避免直接join可能的表名问题
	err := db.Table("file i").
		Select("i.id as file_id, i.original_name, i.expires_at, i.user_id, u.username, u.email").
		Joins("LEFT JOIN user u ON i.user_id = u.id").
		Where("i.expires_at IS NOT NULL AND i.expires_at < ? AND i.expiry_notification_sent = ? AND i.user_id != 0",
			threeDaysLater, false).
		Find(&results).Error

	if err != nil {
		return 0, 1, errors.Wrap(err, errors.CodeDBQueryFailed, "查询即将过期文件失败")
	}

	if len(results) == 0 {
		return 0, 0, nil
	}

	userImages := make(map[uint][]struct {
		FileID       string
		OriginalName string
		ExpiresAt    time.Time
	})

	for _, result := range results {
		userImages[result.UserID] = append(userImages[result.UserID], struct {
			FileID       string
			OriginalName string
			ExpiresAt    time.Time
		}{
			FileID:       result.FileID,
			OriginalName: result.OriginalName,
			ExpiresAt:    result.ExpiresAt,
		})
	}

	successCount := 0
	failedCount := 0

	// 为每个用户发送提醒
	for userID, images := range userImages {

		if err := s.sendUserExpiryNotification(userID, images); err != nil {
			failedCount++
		} else {
			successCount++

			imageIDs := make([]string, len(images))
			for i, img := range images {
				imageIDs[i] = img.FileID
			}

			if updateErr := db.Model(&models.File{}).
				Where("id IN ?", imageIDs).
				Update("expiry_notification_sent", true).Error; updateErr != nil {
			}
		}
	}

	return successCount, failedCount, nil
}

// sendUserExpiryNotification 发送用户过期提醒
func (s *MessageService) sendUserExpiryNotification(userID uint, images []struct {
	FileID       string
	OriginalName string
	ExpiresAt    time.Time
}) error {
	fileList := make([]map[string]interface{}, len(images))
	for i, img := range images {
		fileList[i] = map[string]interface{}{
			"file_id":    img.FileID,
			"file_name":  img.OriginalName,
			"expires_at": img.ExpiresAt.Format("2006-01-02 15:04:05"),
		}
	}

	data := map[string]interface{}{
		"file_count": len(images),
		"files":      fileList,
	}

	err := s.SendMessage(userID, common.MessageTypeFileExpiryWarning, data, nil)
	if err != nil {
		return err
	}

	return nil
}
