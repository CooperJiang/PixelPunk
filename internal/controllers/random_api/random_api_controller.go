package random_api

import (
	"io"
	"net/http"
	"strconv"

	"pixelpunk/internal/controllers/random_api/dto"
	"pixelpunk/internal/middleware"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/activity"
	filesvc "pixelpunk/internal/services/file"
	messageService "pixelpunk/internal/services/message"
	"pixelpunk/internal/services/random_api"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/storage"

	"github.com/gin-gonic/gin"
)

func parseIDParam(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return 0, errors.New(errors.CodeInvalidParameter, "无效的ID")
	}
	return uint(id), nil
}

func CreateRandomAPI(c *gin.Context) {
	req, err := common.ValidateRequest[dto.CreateRandomAPIDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	userID := middleware.GetCurrentUserID(c)
	api, err := random_api.CreateRandomAPI(userID, req.Name, req.FolderID, req.ReturnType)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	activity.LogRandomAPICreate(userID, api.Name, api.ID)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("❌ [RandomAPI] 发送创建消息时panic: %v", r)
			}
		}()

		msgService := messageService.GetMessageService()
		variables := map[string]interface{}{
			"api_id":       api.ID,
			"api_name":     api.Name,
			"related_type": "random_api",
			"related_id":   api.ID,
		}
		if err := msgService.SendTemplateMessage(userID, common.MessageTypeRandomAPICreated, variables); err != nil {
			logger.Error("❌ [RandomAPI] 发送随机API创建通知失败: userID=%d, apiID=%d, error=%v", userID, api.ID, err)
		} else {
		}
	}()

	errors.ResponseSuccess(c, gin.H{
		"id":          api.ID,
		"name":        api.Name,
		"api_key":     api.APIKey,
		"folder_id":   api.FolderID,
		"return_type": api.ReturnType,
		"status":      api.Status,
		"call_count":  api.CallCount,
		"created_at":  api.CreatedAt,
	}, "创建成功")
}

func GetRandomAPIList(c *gin.Context) {
	req, err := common.ValidateRequest[dto.RandomAPIQueryDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	page, size := req.Page, req.Size
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = common.DefaultPageSize
	}

	apis, total, err := random_api.GetRandomAPIList(middleware.GetCurrentUserID(c), page, size, req.Status, req.Search)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	items := make([]gin.H, 0, len(apis))
	for _, api := range apis {
		items = append(items, gin.H{
			"id":             api.ID,
			"name":           api.Name,
			"api_key":        api.APIKey,
			"folder_id":      api.FolderID,
			"folder_name":    random_api.GetFolderNameByID(api.FolderID),
			"return_type":    api.ReturnType,
			"status":         api.Status,
			"is_active":      api.IsActive(),
			"call_count":     api.CallCount,
			"last_called_at": api.LastCalledAt,
			"created_at":     api.CreatedAt,
			"updated_at":     api.UpdatedAt,
		})
	}

	errors.ResponseSuccess(c, gin.H{
		"items": items,
		"total": total,
		"page":  page,
		"size":  size,
	}, "获取成功")
}

func UpdateRandomAPIStatus(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	req, err := common.ValidateRequest[dto.UpdateRandomAPIStatusDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	userID := middleware.GetCurrentUserID(c)

	api, err := random_api.GetRandomAPIByID(id, userID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if err := random_api.UpdateRandomAPIStatus(id, userID, req.Status); err != nil {
		errors.HandleError(c, err)
		return
	}

	activity.LogRandomAPIToggleStatus(userID, api.Name, api.ID, req.Status)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("❌ [RandomAPI] 发送状态切换消息时panic: %v", r)
			}
		}()

		var messageType string
		if req.Status == models.RandomImageAPIStatusActive {
			messageType = common.MessageTypeRandomAPIEnabled
		} else {
			messageType = common.MessageTypeRandomAPIDisabled
		}

		msgService := messageService.GetMessageService()
		variables := map[string]interface{}{
			"api_id":       api.ID,
			"api_name":     api.Name,
			"related_type": "random_api",
			"related_id":   api.ID,
		}

		if err := msgService.SendTemplateMessage(userID, messageType, variables); err != nil {
			logger.Error("❌ [RandomAPI] 发送随机API状态切换通知失败: userID=%d, apiID=%d, type=%s, error=%v", userID, api.ID, messageType, err)
		} else {
		}
	}()

	statusText := "已激活"
	if req.Status == models.RandomImageAPIStatusDisabled {
		statusText = "已禁用"
	}
	errors.ResponseSuccess(c, nil, statusText)
}

func UpdateRandomAPIConfig(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	req, err := common.ValidateRequest[dto.UpdateRandomAPIConfigDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if err := random_api.UpdateRandomAPIConfig(id, middleware.GetCurrentUserID(c), req.FolderID, req.ReturnType); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "配置更新成功")
}

func DeleteRandomAPI(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	userID := middleware.GetCurrentUserID(c)

	api, err := random_api.GetRandomAPIByID(id, userID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if err := random_api.DeleteRandomAPI(id, userID); err != nil {
		errors.HandleError(c, err)
		return
	}

	activity.LogRandomAPIDelete(userID, api.Name, api.ID)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("❌ [RandomAPI] 发送删除消息时panic: %v", r)
			}
		}()

		msgService := messageService.GetMessageService()
		variables := map[string]interface{}{
			"api_id":       api.ID,
			"api_name":     api.Name,
			"related_type": "random_api",
			"related_id":   api.ID,
		}
		if err := msgService.SendTemplateMessage(userID, common.MessageTypeRandomAPIDeleted, variables); err != nil {
			logger.Error("❌ [RandomAPI] 发送随机API删除通知失败: userID=%d, apiID=%d, error=%v", userID, api.ID, err)
		} else {
		}
	}()

	errors.ResponseSuccess(c, nil, "删除成功")
}

func GetRandomImage(c *gin.Context) {
	apiKey := c.Param("api_key")
	if apiKey == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "缺少API密钥"))
		return
	}

	file, api, err := random_api.GetRandomImageByAPIKeyWithConfig(apiKey)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	if api.ReturnType == models.RandomImageAPIReturnTypeDirect {
		provider, err := storage.GetStorageProviderByChannelID(file.StorageProviderID)
		if err != nil {
			errors.HandleError(c, err)
			return
		}

		if provider.IsDirectAccess() {
			localPath, err := filesvc.GetFileLocalPath(*file, false)
			if err != nil {
				errors.HandleError(c, err)
				return
			}
			c.Header("Content-Type", file.MimeType)
			c.Header("Content-Disposition", "inline; filename=\""+file.OriginalName+"\"")
			c.File(localPath)
		} else {
			content, contentType, err := provider.GetRemoteContent(file.RemoteURL, false, file.UserID)
			if err != nil {
				errors.HandleError(c, err)
				return
			}
			defer content.Close()

			c.Header("Content-Type", contentType)
			c.Header("Content-Length", strconv.FormatInt(file.Size, 10))
			c.Header("Content-Disposition", "inline; filename=\""+file.OriginalName+"\"")
			c.Status(http.StatusOK)
			io.Copy(c.Writer, content)
		}
	} else {
		fullURL, _, _ := random_api.GetFileFullURL(*file)
		c.Redirect(302, fullURL)
	}
}
