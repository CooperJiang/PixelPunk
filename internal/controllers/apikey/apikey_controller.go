package apikey

import (
	"pixelpunk/internal/controllers/apikey/dto"
	"pixelpunk/internal/middleware"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/activity"
	"pixelpunk/internal/services/apikey"
	messageService "pixelpunk/internal/services/message"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"

	"github.com/gin-gonic/gin"
)

func CreateAPIKey(c *gin.Context) {
	req, err := common.ValidateRequest[dto.CreateAPIKeyDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	userID := middleware.GetCurrentUserID(c)

	apiKeyModel, keyValue, err := apikey.CreateAPIKey(
		userID,
		req.Name,
		req.StorageLimit,
		req.SingleFileLimit,
		req.UploadCountLimit,
		req.AllowedTypes,
		req.FolderID,
		req.ExpiresInDays,
	)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	response := gin.H{
		"id":                 apiKeyModel.ID,
		"name":               apiKeyModel.Name,
		"key":                keyValue,
		"status":             apiKeyModel.Status,
		"storage_limit":      apiKeyModel.StorageLimit,
		"single_file_limit":  apiKeyModel.SingleFileLimit,
		"upload_count_limit": apiKeyModel.UploadCountLimit,
		"allowed_types":      apikey.ParseAllowedTypes(apiKeyModel.AllowedTypes),
		"folder_id":          apiKeyModel.FolderID,
		"expires_at":         apiKeyModel.ExpiresAt,
		"created_at":         apiKeyModel.CreatedAt,
	}

	// 记录API密钥创建活动日志
	activity.LogAPIKeyCreate(userID, apiKeyModel.Name, apiKeyModel.ID)

	go func() {
		msgService := messageService.GetMessageService()
		variables := map[string]interface{}{
			"key_id":       apiKeyModel.ID,
			"key_name":     apiKeyModel.Name,
			"related_type": "apikey",
			"related_id":   apiKeyModel.ID,
		}
		if err := msgService.SendTemplateMessage(userID, common.MessageTypeAPIKeyCreated, variables); err != nil {
			logger.Warn("发送API密钥创建通知失败: userID=%d, keyID=%s, error=%v", userID, apiKeyModel.ID, err)
		}
	}()

	errors.ResponseSuccess(c, response, "创建API密钥成功")
}

func GetAPIKeyList(c *gin.Context) {
	req, err := common.ValidateRequest[dto.APIKeyQueryDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	userID := middleware.GetCurrentUserID(c)

	page := req.Page
	if page <= 0 {
		page = 1
	}
	size := req.Size
	if size <= 0 {
		size = common.DefaultPageSize
	}

	apikeys, total, err := apikey.GetAPIKeyList(userID, page, size, req.Status, req.Search)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	items := make([]gin.H, 0, len(apikeys))
	for _, key := range apikeys {
		folderPath := apikey.GetFolderFullPath(userID, key.FolderID)

		items = append(items, gin.H{
			"id":                 key.ID,
			"key":                key.KeyValue, // 添加API密钥本身
			"name":               key.Name,
			"status":             key.Status,
			"status_text":        getStatusText(key.Status),
			"is_active":          key.IsActive(),
			"storage_limit":      key.StorageLimit,
			"storage_used":       key.StorageUsed,
			"upload_count_limit": key.UploadCountLimit,
			"upload_count_used":  key.UploadCountUsed,
			"single_file_limit":  key.SingleFileLimit,
			"folder_id":          key.FolderID,
			"folder_path":        folderPath,
			"allowed_types":      apikey.ParseAllowedTypes(key.AllowedTypes),
			"is_expired":         key.IsExpired(),
			"expires_at":         key.ExpiresAt,
			"last_used_at":       key.LastUsedAt,
			"created_at":         key.CreatedAt,
		})
	}

	response := gin.H{
		"items": items,
		"pagination": gin.H{
			"total":       total,
			"page":        page,
			"size":        size,
			"total_pages": (total + int64(size) - 1) / int64(size),
		},
	}

	errors.ResponseSuccess(c, response, "获取API密钥列表成功")
}

func GetAPIKeyDetail(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	keyID := c.Param("key_id")
	if keyID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "API密钥ID不能为空"))
		return
	}

	key, err := apikey.GetAPIKey(userID, keyID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	folderPath := apikey.GetFolderFullPath(userID, key.FolderID)

	response := gin.H{
		"id":                 key.ID,
		"name":               key.Name,
		"status":             key.Status,
		"status_text":        getStatusText(key.Status),
		"is_active":          key.IsActive(),
		"storage_limit":      key.StorageLimit,
		"storage_used":       key.StorageUsed,
		"upload_count_limit": key.UploadCountLimit,
		"upload_count_used":  key.UploadCountUsed,
		"single_file_limit":  key.SingleFileLimit,
		"folder_id":          key.FolderID,
		"folder_path":        folderPath,
		"allowed_types":      apikey.ParseAllowedTypes(key.AllowedTypes),
		"is_expired":         key.IsExpired(),
		"expires_at":         key.ExpiresAt,
		"last_used_at":       key.LastUsedAt,
		"created_at":         key.CreatedAt,
		"updated_at":         key.UpdatedAt,
	}

	errors.ResponseSuccess(c, response, "获取API密钥详情成功")
}

func UpdateAPIKey(c *gin.Context) {
	req, err := common.ValidateRequest[dto.UpdateAPIKeyDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	userID := middleware.GetCurrentUserID(c)

	keyID := c.Param("key_id")
	if keyID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "API密钥ID不能为空"))
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Status != 0 {
		updates["status"] = req.Status
	}
	if req.StorageLimit >= 0 {
		updates["storage_limit"] = req.StorageLimit
	}
	if req.SingleFileLimit >= 0 {
		updates["single_file_limit"] = req.SingleFileLimit
	}
	if req.UploadCountLimit >= 0 {
		updates["upload_count_limit"] = req.UploadCountLimit
	}
	if req.AllowedTypes != nil {
		updates["allowed_types"] = req.AllowedTypes
	}
	if req.FolderID != "" {
		updates["folder_id"] = req.FolderID
	}
	if c.Request.Method == "PUT" || c.PostForm("expires_in_days") != "" || c.Request.Header.Get("Content-Type") == "application/json" {
		updates["expires_in_days"] = req.ExpiresInDays
	}

	updatedKey, err := apikey.UpdateAPIKey(userID, keyID, updates)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	response := gin.H{
		"id":                 updatedKey.ID,
		"name":               updatedKey.Name,
		"status":             updatedKey.Status,
		"status_text":        getStatusText(updatedKey.Status),
		"is_active":          updatedKey.IsActive(),
		"storage_limit":      updatedKey.StorageLimit,
		"storage_used":       updatedKey.StorageUsed,
		"upload_count_limit": updatedKey.UploadCountLimit,
		"upload_count_used":  updatedKey.UploadCountUsed,
		"single_file_limit":  updatedKey.SingleFileLimit,
		"folder_id":          updatedKey.FolderID,
		"allowed_types":      apikey.ParseAllowedTypes(updatedKey.AllowedTypes),
		"is_expired":         updatedKey.IsExpired(),
		"expires_at":         updatedKey.ExpiresAt,
		"updated_at":         updatedKey.UpdatedAt,
	}

	errors.ResponseSuccess(c, response, "更新API密钥成功")
}

func DeleteAPIKey(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	keyID := c.Param("key_id")
	if keyID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "API密钥ID不能为空"))
		return
	}

	// 先获取API密钥信息以记录日志
	keyInfo, err := apikey.GetAPIKey(userID, keyID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	err = apikey.DeleteAPIKey(userID, keyID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	// 记录API密钥删除活动日志
	activity.LogAPIKeyDelete(userID, keyInfo.Name, keyID)

	go func() {
		msgService := messageService.GetMessageService()
		variables := map[string]interface{}{
			"key_id":       keyID,
			"key_name":     keyInfo.Name,
			"related_type": "apikey",
			"related_id":   keyID,
		}
		if err := msgService.SendTemplateMessage(userID, common.MessageTypeAPIKeyDeleted, variables); err != nil {
			logger.Warn("发送API密钥删除通知失败: userID=%d, keyID=%s, error=%v", userID, keyID, err)
		}
	}()

	errors.ResponseSuccess(c, gin.H{"id": keyID}, "删除API密钥成功")
}

func ToggleAPIKeyStatus(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	keyID := c.Param("key_id")
	if keyID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "API密钥ID不能为空"))
		return
	}

	updatedKey, err := apikey.ToggleAPIKeyStatus(userID, keyID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	// 记录API密钥状态切换活动日志
	activity.LogAPIKeyToggleStatus(userID, updatedKey.Name, updatedKey.ID, updatedKey.Status)

	go func() {
		msgService := messageService.GetMessageService()
		variables := map[string]interface{}{
			"key_id":       updatedKey.ID,
			"key_name":     updatedKey.Name,
			"related_type": "apikey",
			"related_id":   updatedKey.ID,
		}

		var messageType string
		if updatedKey.Status == models.APIKeyStatusActive {
			messageType = common.MessageTypeAPIKeyEnabled
		} else {
			messageType = common.MessageTypeAPIKeyDisabled
		}

		if err := msgService.SendTemplateMessage(userID, messageType, variables); err != nil {
			logger.Warn("发送API密钥状态切换通知失败: userID=%d, keyID=%s, error=%v", userID, updatedKey.ID, err)
		}
	}()

	response := gin.H{
		"id":          updatedKey.ID,
		"status":      updatedKey.Status,
		"status_text": getStatusText(updatedKey.Status),
		"is_active":   updatedKey.IsActive(),
	}

	errors.ResponseSuccess(c, response, "切换API密钥状态成功")
}

func GetAPIKeyStats(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	keyID := c.Param("key_id")
	if keyID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "API密钥ID不能为空"))
		return
	}

	stats, err := apikey.GetAPIKeyStats(userID, keyID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, stats, "获取API密钥统计数据成功")
}

func RegenerateAPIKey(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	keyID := c.Param("key_id")
	if keyID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "API密钥ID不能为空"))
		return
	}

	// 先获取API密钥信息以记录日志
	keyInfo, err := apikey.GetAPIKey(userID, keyID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	newKeyValue, err := apikey.RegenerateAPIKey(userID, keyID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	// 记录API密钥重新生成活动日志
	activity.LogAPIKeyRegenerate(userID, keyInfo.Name, keyID)

	// 发送消息通知（重新生成是安全相关操作，重要性高）
	go func() {
		msgService := messageService.GetMessageService()
		variables := map[string]interface{}{
			"key_id":       keyID,
			"key_name":     keyInfo.Name,
			"related_type": "apikey",
			"related_id":   keyID,
		}
		if err := msgService.SendTemplateMessage(userID, common.MessageTypeAPIKeyRegenerated, variables); err != nil {
			logger.Warn("发送API密钥重新生成通知失败: userID=%d, keyID=%s, error=%v", userID, keyID, err)
		}
	}()

	response := gin.H{
		"id":  keyID,
		"key": newKeyValue,
	}

	errors.ResponseSuccess(c, response, "重新生成API密钥成功")
}

func getStatusText(status int) string {
	switch status {
	case models.APIKeyStatusActive:
		return "正常"
	case models.APIKeyStatusDisabled:
		return "已禁用"
	default:
		return "未知状态"
	}
}
