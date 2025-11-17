package message

import (
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"

	"gorm.io/gorm"
)

/* TemplateService 消息模板服务 */
type TemplateService struct{}

/* GetTemplateService 获取模板服务实例 */
func GetTemplateService() *TemplateService {
	return &TemplateService{}
}

/* GetAllTemplates 获取所有消息模板 */
func (s *TemplateService) GetAllTemplates() ([]models.MessageTemplate, error) {
	db := database.GetDB()

	var templates []models.MessageTemplate
	err := db.Order("type ASC").Find(&templates).Error
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询消息模板失败")
	}

	return templates, nil
}

/* GetTemplateByType 根据类型获取模板 */
func (s *TemplateService) GetTemplateByType(templateType string) (*models.MessageTemplate, error) {
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

/* CreateTemplate 创建消息模板 */
func (s *TemplateService) CreateTemplate(template *models.MessageTemplate) error {
	db := database.GetDB()

	var existing models.MessageTemplate
	if err := db.Where("type = ?", template.Type).First(&existing).Error; err == nil {
		return errors.New(errors.CodeUserExists, "消息模板类型已存在")
	}

	if err := db.Create(template).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBCreateFailed, "创建消息模板失败")
	}

	return nil
}

/* UpdateTemplate 更新消息模板 */
func (s *TemplateService) UpdateTemplate(id uint, updates map[string]interface{}) error {
	db := database.GetDB()

	var template models.MessageTemplate
	if err := db.First(&template, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.CodeNotFound, "消息模板不存在")
		}
		return errors.Wrap(err, errors.CodeDBQueryFailed, "查询消息模板失败")
	}

	if err := db.Model(&template).Updates(updates).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBUpdateFailed, "更新消息模板失败")
	}

	return nil
}

/* DeleteTemplate 删除消息模板 */
func (s *TemplateService) DeleteTemplate(id uint) error {
	db := database.GetDB()

	var template models.MessageTemplate
	if err := db.First(&template, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.CodeNotFound, "消息模板不存在")
		}
		return errors.Wrap(err, errors.CodeDBQueryFailed, "查询消息模板失败")
	}

	if err := db.Delete(&template).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBDeleteFailed, "删除消息模板失败")
	}

	return nil
}

/* ToggleTemplate 切换模板启用状态 */
func (s *TemplateService) ToggleTemplate(id uint, enabled bool) error {
	db := database.GetDB()

	var template models.MessageTemplate
	if err := db.First(&template, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.CodeNotFound, "消息模板不存在")
		}
		return errors.Wrap(err, errors.CodeDBQueryFailed, "查询消息模板失败")
	}

	if err := db.Model(&template).Update("is_enabled", enabled).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBUpdateFailed, "更新模板状态失败")
	}

	return nil
}

/* InitDefaultTemplates 初始化默认消息模板 */
func (s *TemplateService) InitDefaultTemplates() error {
	db := database.GetDB()

	templates := []models.MessageTemplate{
		{
			Type:               common.MessageTypeAccountRegister,
			Title:              "欢迎来到 {{.site_name}}",
			Content:            "注册成功！你已获得 {{.storage_mb}}MB 存储空间和 {{.bandwidth_mb}}MB 月流量，快来上传你的第一个文件吧！",
			Description:        "用户注册成功通知",
			IsEnabled:          true,
			SendEmail:          false,
			ShowToast:          true,
			ToastType:          "success",
			DefaultActionType:  common.ActionTypeView,
			DefaultActionText:  "开始上传",
			DefaultActionStyle: "primary",
			ActionURLTemplate:  "/upload",
		},
		{
			Type:               common.MessageTypeContentReviewApproved,
			Title:              "文件审核通过",
			Content:            "您的文件 \"{{.file_name}}\" 已通过审核，现在可以正常访问和分享。",
			Description:        "文件审核通过通知",
			IsEnabled:          true,
			SendEmail:          false,
			ShowToast:          true,
			ToastType:          "success",
			DefaultActionType:  common.ActionTypeView,
			DefaultActionText:  "查看文件",
			DefaultActionStyle: "primary",
			ActionURLTemplate:  "/files/{{.file_id}}",
		},
		{
			Type:               common.MessageTypeContentReviewRejected,
			Title:              "文件审核未通过",
			Content:            "很抱歉，您的文件 \"{{.file_name}}\" 未通过审核并已被删除。原因：{{.reason}}",
			Description:        "文件审核拒绝通知",
			IsEnabled:          true,
			SendEmail:          false,
			ShowToast:          true,
			ToastType:          "error",
			DefaultActionType:  common.ActionTypeView,
			DefaultActionText:  "查看详情",
			DefaultActionStyle: "secondary",
			ActionURLTemplate:  "/review-details/{{.review_id}}",
		},
		{
			Type:               common.MessageTypeContentReviewPending,
			Title:              "文件进入审核",
			Content:            "您的文件 \"{{.file_name}}\" 疑似违规，进入审核状态中。{{if .nsfw_reason}}检测原因：{{.nsfw_reason}}。{{end}}审核期间暂不可访问，我们会尽快处理。",
			Description:        "文件进入审核通知",
			IsEnabled:          true,
			SendEmail:          false,
			ShowToast:          true,
			ToastType:          "warning",
			DefaultActionType:  common.ActionTypeView,
			DefaultActionText:  "查看进度",
			DefaultActionStyle: "secondary",
			ActionURLTemplate:  "/my-files?filter=pending",
		},
		{
			Type:               common.MessageTypeAccountStorageGranted,
			Title:              "存储空间赠送",
			Content:            "管理员为您赠送了 {{.granted_mb}}MB 存储空间，当前总容量：{{.total_mb}}MB。感谢您的支持！",
			Description:        "存储空间赠送通知",
			IsEnabled:          true,
			SendEmail:          false,
			ShowToast:          true,
			ToastType:          "success",
			DefaultActionType:  common.ActionTypeManage,
			DefaultActionText:  "管理存储",
			DefaultActionStyle: "primary",
			ActionURLTemplate:  "/storage/overview",
		},
		{
			Type:               common.MessageTypeAccountBandwidthGranted,
			Title:              "带宽流量赠送",
			Content:            "管理员为您赠送了 {{.granted_mb}}MB 月流量，当前总流量：{{.total_mb}}MB。",
			Description:        "带宽流量赠送通知",
			IsEnabled:          true,
			SendEmail:          false,
			ShowToast:          true,
			ToastType:          "success",
			DefaultActionType:  common.ActionTypeView,
			DefaultActionText:  "查看统计",
			DefaultActionStyle: "primary",
			ActionURLTemplate:  "/stats/bandwidth",
		},
		{
			Type:               common.MessageTypeSystemMaintenance,
			Title:              "系统维护通知",
			Content:            "系统将于 {{.start_time}} 至 {{.end_time}} 进行维护升级，期间服务可能中断，请提前做好准备。",
			Description:        "系统维护通知",
			IsEnabled:          true,
			SendEmail:          false,
			ShowToast:          true,
			ToastType:          "warning",
			DefaultActionType:  common.ActionTypeView,
			DefaultActionText:  "查看公告",
			DefaultActionStyle: "secondary",
			ActionURLTemplate:  "/notices/{{.notice_id}}",
		},
		{
			Type:               common.MessageTypeStorageQuotaWarning,
			Title:              "存储空间不足",
			Content:            "您的存储空间已使用 {{.used_percent}}%（{{.used_mb}}MB/{{.total_mb}}MB），请及时清理或申请扩容。",
			Description:        "存储空间不足警告",
			IsEnabled:          true,
			SendEmail:          false,
			ShowToast:          true,
			ToastType:          "warning",
			DefaultActionType:  common.ActionTypeManage,
			DefaultActionText:  "管理存储",
			DefaultActionStyle: "warning",
			ActionURLTemplate:  "/storage/manage",
		},
		{
			Type:               common.MessageTypeStorageQuotaIncreased,
			Title:              "存储空间增加",
			Content:            "管理员为您增加了存储空间：{{.old_size}} → {{.new_size}}（+{{.increased_size}}）。感谢您的支持！",
			Description:        "存储空间增加通知",
			IsEnabled:          true,
			SendEmail:          false,
			ShowToast:          true,
			ToastType:          "success",
			DefaultActionType:  common.ActionTypeView,
			DefaultActionText:  "查看详情",
			DefaultActionStyle: "primary",
			ActionURLTemplate:  "/storage/overview",
		},
		{
			Type:               common.MessageTypeStorageQuotaDecreased,
			Title:              "存储空间调整",
			Content:            "您的存储空间已调整：{{.old_size}} → {{.new_size}}（-{{.decreased_size}}）。调整原因：{{.reason}}",
			Description:        "存储空间减少通知",
			IsEnabled:          true,
			SendEmail:          false,
			ShowToast:          true,
			ToastType:          "warning",
			DefaultActionType:  common.ActionTypeManage,
			DefaultActionText:  "管理存储",
			DefaultActionStyle: "warning",
			ActionURLTemplate:  "/storage/manage",
		},
		{
			Type:               common.MessageTypeFileDeletedByAdmin,
			Title:              "文件被删除",
			Content:            "管理员删除了您的文件 \"{{.file_name}}\"。删除原因：{{.reason}}",
			Description:        "管理员删除文件通知",
			IsEnabled:          true,
			SendEmail:          false,
			ShowToast:          true,
			ToastType:          "warning",
			DefaultActionType:  common.ActionTypeView,
			DefaultActionText:  "查看我的文件",
			DefaultActionStyle: "secondary",
			ActionURLTemplate:  "/my-files",
		},
		{
			Type:               common.MessageTypeFileBatchDeletedByAdmin,
			Title:              "文件批量删除",
			Content:            "管理员批量删除了您的 {{.count}} 张文件：{{.file_names}}。删除原因：{{.reason}}",
			Description:        "管理员批量删除文件通知",
			IsEnabled:          true,
			SendEmail:          false,
			ShowToast:          true,
			ToastType:          "warning",
			DefaultActionType:  common.ActionTypeView,
			DefaultActionText:  "查看我的文件",
			DefaultActionStyle: "secondary",
			ActionURLTemplate:  "/my-files",
		},
		{
			Type:               common.MessageTypeFileHardDeletedByAdmin,
			Title:              "文件永久删除",
			Content:            "管理员永久删除了您的文件 \"{{.file_name}}\"，该文件无法恢复。删除原因：{{.reason}}",
			Description:        "管理员硬删除文件通知",
			IsEnabled:          true,
			SendEmail:          false,
			ShowToast:          true,
			ToastType:          "error",
			DefaultActionType:  common.ActionTypeView,
			DefaultActionText:  "查看我的文件",
			DefaultActionStyle: "secondary",
			ActionURLTemplate:  "/my-files",
		},
		{
			Type:               common.MessageTypeFileThumbnailFailed,
			Title:              "缩略图生成失败",
			Content:            "文件「{{.file_name}}」上传成功，但缩略图生成失败。原因：{{.reason}}。您仍然可以正常查看和使用原图。",
			Description:        "缩略图生成失败通知",
			IsEnabled:          true,
			SendEmail:          false,
			ShowToast:          true,
			ToastType:          "warning",
			DefaultActionType:  common.ActionTypeView,
			DefaultActionText:  "查看文件",
			DefaultActionStyle: "secondary",
			ActionURLTemplate:  "/f/{{.file_id}}",
		},
		{
			Type:               common.MessageTypeAPIKeyCreated,
			Title:              "API 密钥创建成功",
			Content:            "您已成功创建 API 密钥「{{.key_name}}」，请妥善保管您的密钥。密钥仅在创建时显示一次，丢失后需要重新生成。",
			Description:        "API密钥创建成功通知",
			IsEnabled:          true,
			SendEmail:          false,
			ShowToast:          true,
			ToastType:          "success",
			DefaultActionType:  common.ActionTypeView,
			DefaultActionText:  "查看密钥",
			DefaultActionStyle: "primary",
			ActionURLTemplate:  "/admin/open-api",
		},
		{
			Type:               common.MessageTypeAPIKeyDeleted,
			Title:              "API 密钥已删除",
			Content:            "您已删除 API 密钥「{{.key_name}}」，使用该密钥的所有请求将立即失效。",
			Description:        "API密钥删除通知",
			IsEnabled:          true,
			SendEmail:          false,
			ShowToast:          true,
			ToastType:          "warning",
			DefaultActionType:  common.ActionTypeView,
			DefaultActionText:  "查看密钥列表",
			DefaultActionStyle: "secondary",
			ActionURLTemplate:  "/admin/open-api",
		},
		{
			Type:               common.MessageTypeAPIKeyRegenerated,
			Title:              "API 密钥已重新生成",
			Content:            "您已重新生成 API 密钥「{{.key_name}}」，旧密钥已失效。请使用新密钥更新您的应用程序配置。",
			Description:        "API密钥重新生成通知",
			IsEnabled:          true,
			SendEmail:          true,
			ShowToast:          true,
			ToastType:          "warning",
			DefaultActionType:  common.ActionTypeView,
			DefaultActionText:  "查看新密钥",
			DefaultActionStyle: "primary",
			ActionURLTemplate:  "/admin/open-api",
		},
		{
			Type:               common.MessageTypeAPIKeyDisabled,
			Title:              "API 密钥已停用",
			Content:            "您已停用 API 密钥「{{.key_name}}」，使用该密钥的所有请求将暂时失效。您可以随时重新启用。",
			Description:        "API密钥停用通知",
			IsEnabled:          true,
			SendEmail:          false,
			ShowToast:          true,
			ToastType:          "info",
			DefaultActionType:  common.ActionTypeView,
			DefaultActionText:  "查看密钥",
			DefaultActionStyle: "secondary",
			ActionURLTemplate:  "/admin/open-api",
		},
		{
			Type:               common.MessageTypeAPIKeyEnabled,
			Title:              "API 密钥已启用",
			Content:            "您已启用 API 密钥「{{.key_name}}」，该密钥现在可以正常使用。",
			Description:        "API密钥启用通知",
			IsEnabled:          true,
			SendEmail:          false,
			ShowToast:          false,
			ToastType:          "success",
			DefaultActionType:  common.ActionTypeView,
			DefaultActionText:  "查看密钥",
			DefaultActionStyle: "secondary",
			ActionURLTemplate:  "/admin/open-api",
		},
		{
			Type:               common.MessageTypeRandomAPICreated,
			Title:              "随机图片 API 创建成功",
			Content:            "您已成功创建随机图片 API「{{.api_name}}」，现在可以通过 API 链接随机获取图片。",
			Description:        "随机图片API创建成功通知",
			IsEnabled:          true,
			SendEmail:          false,
			ShowToast:          true,
			ToastType:          "success",
			DefaultActionType:  common.ActionTypeView,
			DefaultActionText:  "查看 API",
			DefaultActionStyle: "primary",
			ActionURLTemplate:  "/admin/open-api",
		},
		{
			Type:               common.MessageTypeRandomAPIDeleted,
			Title:              "随机图片 API 已删除",
			Content:            "您已删除随机图片 API「{{.api_name}}」，使用该 API 的所有请求将立即失效。",
			Description:        "随机图片API删除通知",
			IsEnabled:          true,
			SendEmail:          false,
			ShowToast:          true,
			ToastType:          "warning",
			DefaultActionType:  common.ActionTypeView,
			DefaultActionText:  "查看 API 列表",
			DefaultActionStyle: "secondary",
			ActionURLTemplate:  "/admin/open-api",
		},
		{
			Type:               common.MessageTypeRandomAPIDisabled,
			Title:              "随机图片 API 已停用",
			Content:            "您已停用随机图片 API「{{.api_name}}」，使用该 API 的所有请求将暂时失效。您可以随时重新启用。",
			Description:        "随机图片API停用通知",
			IsEnabled:          true,
			SendEmail:          false,
			ShowToast:          true,
			ToastType:          "info",
			DefaultActionType:  common.ActionTypeView,
			DefaultActionText:  "查看 API",
			DefaultActionStyle: "secondary",
			ActionURLTemplate:  "/admin/open-api",
		},
		{
			Type:               common.MessageTypeRandomAPIEnabled,
			Title:              "随机图片 API 已启用",
			Content:            "您已启用随机图片 API「{{.api_name}}」，该 API 现在可以正常使用。",
			Description:        "随机图片API启用通知",
			IsEnabled:          true,
			SendEmail:          false,
			ShowToast:          false,
			ToastType:          "success",
			DefaultActionType:  common.ActionTypeView,
			DefaultActionText:  "查看 API",
			DefaultActionStyle: "secondary",
			ActionURLTemplate:  "/admin/open-api",
		},
		{
			Type:               common.MessageTypeShareExpiryWarning,
			Title:              "分享即将过期",
			Content:            "您的分享「{{.share_name}}」将于 {{.expires_in}} 后过期。过期后，访问链接将失效，请及时续期或备份。",
			Description:        "分享过期提醒通知",
			IsEnabled:          true,
			SendEmail:          false,
			ShowToast:          true,
			ToastType:          "warning",
			DefaultActionType:  common.ActionTypeView,
			DefaultActionText:  "查看分享",
			DefaultActionStyle: "secondary",
			ActionURLTemplate:  "/admin/shares",
		},
	}

	for _, template := range templates {
		var existing models.MessageTemplate
		if err := db.Where("type = ?", template.Type).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&template).Error; err != nil {
				} else {
				}
			} else {
			}
		} else {
			existing.Content = template.Content
			existing.SendEmail = template.SendEmail
			if err := db.Save(&existing).Error; err != nil {
			}
			// 删除成功日志，避免启动时输出过多日志
		}
	}

	return nil
}

/* GetTemplatesByPrimaryType 根据主类型获取模板列表 */
func (s *TemplateService) GetTemplatesByPrimaryType(primaryType string) ([]models.MessageTemplate, error) {
	db := database.GetDB()

	var templates []models.MessageTemplate
	err := db.Where("type LIKE ?", primaryType+".%").
		Order("type ASC").
		Find(&templates).Error

	if err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询消息模板失败")
	}

	return templates, nil
}

/* GetEnabledTemplates 获取所有启用的模板 */
func (s *TemplateService) GetEnabledTemplates() ([]models.MessageTemplate, error) {
	db := database.GetDB()

	var templates []models.MessageTemplate
	err := db.Where("is_enabled = ?", true).
		Order("type ASC").
		Find(&templates).Error

	if err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询启用模板失败")
	}

	return templates, nil
}
