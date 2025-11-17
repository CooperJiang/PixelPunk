package setting

import (
	"fmt"
	"pixelpunk/internal/controllers/setting/dto"
	"pixelpunk/internal/middleware"
	"pixelpunk/internal/services/activity"
	"pixelpunk/internal/services/setting"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/email"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/vector"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetSettings(c *gin.Context) {
	req, err := common.ValidateRequest[dto.SettingQueryDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	result, err := setting.GetSettings(req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, result, "获取设置列表成功")
}

func GetSetting(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "设置键名不能为空"))
		return
	}

	result, err := setting.GetSetting(key)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, result, "获取设置成功")
}

func CreateSetting(c *gin.Context) {
	req, err := common.ValidateRequest[dto.SettingCreateDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	result, err := setting.CreateSetting(req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, result, "创建设置成功")
}

func UpdateSetting(c *gin.Context) {
	req, err := common.ValidateRequest[dto.SettingUpdateDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	result, err := setting.UpdateSetting(req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, result, "更新设置成功")
}

func DeleteSetting(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "设置键名不能为空"))
		return
	}

	err := setting.DeleteSetting(key)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "删除设置成功")
}

func BatchCreateSettings(c *gin.Context) {
	req, err := common.ValidateRequest[dto.BatchSettingCreateDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	result, err := setting.BatchCreateSettings(req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if len(result.Failed) > 0 {
		message := fmt.Sprintf("部分设置创建成功，%d项失败", len(result.Failed))
		errors.ResponseSuccess(c, result, message)
		return
	}

	errors.ResponseSuccess(c, result, "批量创建设置成功")
}

func BatchUpdateSettings(c *gin.Context) {
	req, err := common.ValidateRequest[dto.BatchSettingUpdateDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	result, err := setting.BatchUpdateSettings(req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if len(result.Failed) > 0 {
		message := fmt.Sprintf("部分设置更新成功，%d项失败", len(result.Failed))
		errors.ResponseSuccess(c, result, message)
		return
	}

	errors.ResponseSuccess(c, result, "批量更新设置成功")
}

func BatchUpsertSettings(c *gin.Context) {
	req, err := common.ValidateRequest[dto.BatchUpsertSettingDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	containsSecuritySettings := false
	for _, item := range req.Settings {
		if item.Group == "security" {
			containsSecuritySettings = true
			break
		}
	}

	if containsSecuritySettings {
		websiteSettings, err := setting.GetSettingsByGroupAsMap("website")
		if err != nil {
			errors.HandleError(c, errors.Wrap(err, errors.CodeDBQueryFailed, "获取网站设置失败"))
			return
		}

		baseUrlValue, exists := websiteSettings.Settings["site_base_url"]
		if !exists || baseUrlValue == "" {
			errors.HandleError(c, errors.New(errors.CodeValidationFailed, "请先在网站信息中设置网站基础地址，再进行安全设置配置"))
			return
		}

		baseUrl, ok := baseUrlValue.(string)
		if !ok || baseUrl == "" {
			errors.HandleError(c, errors.New(errors.CodeValidationFailed, "网站基础地址配置无效，请先设置有效的网站基础地址"))
			return
		}
	}

	result, err := setting.BatchUpsertSettings(req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if containsSecuritySettings && len(result.Failed) == 0 {
		userID := middleware.GetCurrentUserID(c)

		settingsData := make(map[string]any)
		for _, item := range req.Settings {
			if item.Group == "security" {
				settingsData[item.Key] = map[string]any{
					"value":       item.Value,
					"type":        item.Type,
					"description": item.Description,
				}
			}
		}

		activity.LogHotlinkProtectionChange(userID, settingsData)
	}

	if len(result.Failed) > 0 {
		message := fmt.Sprintf("部分设置操作成功，%d项失败", len(result.Failed))
		errors.ResponseSuccess(c, result, message)
		return
	}

	errors.ResponseSuccess(c, result, "批量设置操作成功")
}

func GetSettingsByGroupAsMap(c *gin.Context) {
	group := c.Param("group")
	if group == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "设置分组不能为空"))
		return
	}

	result, err := setting.GetSettingsByGroupAsMap(group)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, result, "获取分组设置成功")
}

func TestEmailSettings(c *gin.Context) {
	req, err := common.ValidateRequest[dto.EmailTestDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	err = email.TestMailSettings(
		req.Host,
		req.Port,
		req.Username,
		req.Password,
		req.Encryption,
		req.FromAddress,
		req.FromName,
		req.TestEmail,
	)

	if err != nil {
		errors.HandleError(c, errors.Wrap(err, errors.CodeEmailServiceError, "邮件设置测试失败"))
		return
	}

	errors.ResponseSuccess(c, nil, "邮件设置测试成功，请检查收件箱")
}

func RefreshEmailSettings(c *gin.Context) {
	err := email.RefreshMailer()
	if err != nil {
		errors.HandleError(c, errors.Wrap(err, errors.CodeEmailServiceError, "刷新邮件设置失败"))
		return
	}

	errors.ResponseSuccess(c, nil, "邮件设置已刷新")
}

func TestMail(c *gin.Context) {
	req, err := common.ValidateRequest[dto.TestMailDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	err = email.TestMailSettings(
		req.SmtpHost,
		req.SmtpPort,
		req.SmtpUsername,
		req.SmtpPassword,
		req.SmtpEncryption,
		req.SmtpFromAddress,
		req.SmtpFromName,
		req.Email,
	)

	if err != nil {
		errors.HandleError(c, errors.Wrap(err, errors.CodeEmailServiceError, "邮件发送测试失败"))
		return
	}

	errors.ResponseSuccess(c, nil, "邮件发送测试成功，请检查收件箱")
}

func GetGlobalSettings(c *gin.Context) {
	result, err := setting.GetGlobalSettingsGroups()
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, result, "获取全局设置成功")
}

func GetLegalDocuments(c *gin.Context) {
	result, err := setting.GetLegalDocuments()
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, result, "获取法律文档成功")
}

func TestVectorSettings(c *gin.Context) {
	req, err := common.ValidateRequest[dto.VectorTestDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if req.APIKey == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "API密钥不能为空"))
		return
	}

	if req.Provider != "openai" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "目前仅支持OpenAI提供者"))
		return
	}

	client := vector.NewOpenAIEmbeddingClientWithConfig(
		req.APIKey,
		req.BaseURL,
		req.Model,
		req.Timeout,
	)

	if client == nil {
		logger.Error("临时OpenAI客户端创建失败")
		errors.HandleError(c, errors.New(errors.CodeInternal, "向量化客户端初始化失败，请检查配置参数"))
		return
	}

	testText := "这是一个测试文本，用于验证向量化配置是否正确"
	_, err = client.GenerateEmbedding(testText)
	if err != nil {
		logger.Error("向量化测试失败: %v", err)

		var errorMessage string
		var details string

		if errorMsg := err.Error(); errorMsg != "" {
			if strings.Contains(errorMsg, "invalid_api_key") || strings.Contains(errorMsg, "401") {
				errorMessage = "API密钥无效或已过期"
				details = "请检查API密钥是否正确配置"
			} else if strings.Contains(errorMsg, "timeout") || strings.Contains(errorMsg, "deadline exceeded") {
				errorMessage = "请求超时"
				details = fmt.Sprintf("API调用超时（%d秒），请检查网络连接或增加超时时间", req.Timeout)
			} else if strings.Contains(errorMsg, "network") || strings.Contains(errorMsg, "connection") {
				errorMessage = "网络连接失败"
				details = "请检查网络连接或API代理地址配置"
			} else {
				errorMessage = "向量化API调用失败"
				details = fmt.Sprintf("错误详情: %s", errorMsg)
			}
		} else {
			errorMessage = "向量化配置测试失败"
			details = "未知错误，请检查所有配置参数"
		}

		result := map[string]interface{}{
			"success": false,
			"message": errorMessage,
			"details": details,
			"config": map[string]interface{}{
				"provider": req.Provider,
				"model":    req.Model,
				"baseUrl":  req.BaseURL,
				"timeout":  req.Timeout,
			},
		}

		errors.ResponseSuccess(c, result, "向量配置测试完成")
		return
	}

	dimension := client.GetDimension()
	model := client.GetModel()

	result := map[string]interface{}{
		"success": true,
		"message": "向量配置测试成功！API连接正常，模型可用。",
		"details": fmt.Sprintf("模型: %s\nAPI地址: %s\n向量维度: %d\n超时时间: %d秒",
			model, req.BaseURL, dimension, req.Timeout),
		"config": map[string]interface{}{
			"provider":  req.Provider,
			"model":     model,
			"baseUrl":   req.BaseURL,
			"timeout":   req.Timeout,
			"dimension": dimension,
		},
	}

	errors.ResponseSuccess(c, result, "向量配置测试成功")
}

func TestQdrantConnection(c *gin.Context) {
	req, err := common.ValidateRequest[dto.QdrantTestDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if req.QdrantURL == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "Qdrant服务器地址不能为空"))
		return
	}

	timeout := req.QdrantTimeout
	if timeout <= 0 {
		timeout = 30
	}

	startTime := time.Now()
	client := vector.NewQdrantClient(req.QdrantURL, timeout)

	err = client.InitCollection()
	responseTime := time.Since(startTime).Milliseconds()

	if err != nil {
		logger.Error("Qdrant连接测试失败: %v", err)

		errorMessage := "Qdrant连接失败"
		var details string

		errorMsg := err.Error()
		if strings.Contains(errorMsg, "connection refused") || strings.Contains(errorMsg, "dial tcp") {
			details = "无法连接到Qdrant服务，请确保服务已启动且地址正确"
		} else if strings.Contains(errorMsg, "timeout") || strings.Contains(errorMsg, "deadline exceeded") {
			details = fmt.Sprintf("连接超时（%d秒），请检查网络或增加超时时间", timeout)
		} else if strings.Contains(errorMsg, "404") {
			details = "Qdrant服务未找到，请检查URL路径是否正确"
		} else {
			details = fmt.Sprintf("错误详情: %s", errorMsg)
		}

		result := map[string]interface{}{
			"success": false,
			"message": errorMessage,
			"error":   details,
			"details": map[string]interface{}{
				"url":           req.QdrantURL,
				"response_time": responseTime,
			},
		}

		errors.ResponseSuccess(c, result, "Qdrant连接测试完成")
		return
	}

	vectorEngine := vector.GetGlobalVectorEngine()
	var totalVectors int64 = 0
	var qdrantVersion string = "未知"

	if vectorEngine != nil && vectorEngine.IsEnabled() {
		if stats, err := vectorEngine.GetStorageStats(); err == nil && stats != nil {
			totalVectors = stats.TotalVectors
		}
		// 尝试获取 Qdrant 版本信息（可选）
		qdrantVersion = "1.x" // Qdrant 客户端暂不提供版本API，使用占位符
	}

	result := map[string]interface{}{
		"success": true,
		"message": "✅ Qdrant连接成功！向量数据库工作正常。",
		"details": map[string]interface{}{
			"url":           req.QdrantURL,
			"version":       qdrantVersion,
			"total_vectors": totalVectors,
			"response_time": responseTime,
		},
	}

	errors.ResponseSuccess(c, result, "Qdrant连接测试成功")
}

func GetOAuthConfig(c *gin.Context) {
	result, err := setting.GetOAuthConfig()
	if err != nil {
		errors.HandleError(c, err)
		return
	}
	errors.ResponseSuccess(c, result, "获取 OAuth 配置成功")
}

func TestProxy(c *gin.Context) {
	req, err := common.ValidateRequest[dto.TestProxyDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	result, err := setting.TestProxyConnection(req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, result, "代理测试完成")
}
