package migrations

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"pixelpunk/internal/controllers/setting/dto"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/setting"
	"pixelpunk/pkg/logger"

	"gorm.io/gorm"
)

// generateRandomJWTSecret 生成随机的JWT密钥
func generateRandomJWTSecret() string {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		logger.Errorf("无法生成随机JWT密钥: %v, 使用默认值", err)
		return "defaultSecretKey123456789012345678901234"
	}
	return base64.StdEncoding.EncodeToString(bytes)
}

// generateURLSafeToken 生成URL安全的随机令牌
func generateURLSafeToken(length int) string {
	bytes := make([]byte, length*2)
	_, err := rand.Read(bytes)
	if err != nil {
		logger.Errorf("生成随机令牌失败: %v", err)
		return "default_token_123456789012345678901234567890"
	}

	token := base64.RawURLEncoding.EncodeToString(bytes)

	// 确保长度正确，如果太长则截断
	if len(token) > length {
		token = token[:length]
	}

	return token
}

// AddSystemSettings 初始化系统设置
func AddSystemSettings(db *gorm.DB) error {
	migrationName := "add_system_settings_v2" // 更新版本号以包含所有功能

	// 检查迁移是否已执行
	if migrationRecord, err := IsMigrationApplied(db, migrationName); err != nil {
		return fmt.Errorf("检查迁移状态失败: %v", err)
	} else if migrationRecord {
		logger.Infof("系统设置迁移已应用，跳过")
		return nil
	}

	logger.Infof("正在执行系统设置初始化迁移...")

	// 生成URL安全的访问令牌
	internalServiceToken := generateURLSafeToken(32)

	// 定义所有系统设置
	var allSettings []dto.SettingCreateDTO

	// 网站后端功能设置
	websiteSettings := []dto.SettingCreateDTO{
		{
			Key:         "admin_email",
			Value:       DefaultSettings.Website.AdminEmail,
			Type:        "string",
			Group:       "website",
			Description: "管理员邮箱",
			IsSystem:    true,
		},
		{
			Key:         "site_base_url",
			Value:       DefaultSettings.Website.SiteBaseURL,
			Type:        "string",
			Group:       "website",
			Description: "网站基础地址",
			IsSystem:    true,
		},
	}
	allSettings = append(allSettings, websiteSettings...)

	// 网站信息配置设置 (前端显示用)
	websiteInfoSettings := []dto.SettingCreateDTO{
		{
			Key:         "site_name",
			Value:       DefaultSettings.WebsiteInfo.SiteName,
			Type:        "string",
			Group:       "website_info",
			Description: "网站名称",
			IsSystem:    true,
		},
		{
			Key:         "site_description",
			Value:       DefaultSettings.WebsiteInfo.SiteDescription,
			Type:        "string",
			Group:       "website_info",
			Description: "网站描述",
			IsSystem:    true,
		},
		{
			Key:         "site_keywords",
			Value:       DefaultSettings.WebsiteInfo.SiteKeywords,
			Type:        "string",
			Group:       "website_info",
			Description: "网站关键词",
			IsSystem:    true,
		},
		{
			Key:         "icp_number",
			Value:       DefaultSettings.WebsiteInfo.ICPNumber,
			Type:        "string",
			Group:       "website_info",
			Description: "ICP备案号",
			IsSystem:    true,
		},
		{
			Key:         "show_file_count",
			Value:       DefaultSettings.WebsiteInfo.ShowFileCount,
			Type:        "boolean",
			Group:       "website_info",
			Description: "显示文件数量统计",
			IsSystem:    true,
		},
		{
			Key:         "show_storage_usage",
			Value:       DefaultSettings.WebsiteInfo.ShowStorageUsage,
			Type:        "boolean",
			Group:       "website_info",
			Description: "显示存储使用统计",
			IsSystem:    true,
		},
		{
			Key:         "site_logo_url",
			Value:       DefaultSettings.WebsiteInfo.SiteLogoURL,
			Type:        "string",
			Group:       "website_info",
			Description: "网站Logo URL",
			IsSystem:    true,
		},
		{
			Key:         "favicon_url",
			Value:       DefaultSettings.WebsiteInfo.FaviconURL,
			Type:        "string",
			Group:       "website_info",
			Description: "Favicon URL",
			IsSystem:    true,
		},
		{
			Key:         "copyright_text",
			Value:       DefaultSettings.WebsiteInfo.CopyrightText,
			Type:        "string",
			Group:       "website_info",
			Description: "版权信息",
			IsSystem:    true,
		},
		{
			Key:         "contact_email",
			Value:       DefaultSettings.WebsiteInfo.ContactEmail,
			Type:        "string",
			Group:       "website_info",
			Description: "联系邮箱",
			IsSystem:    true,
		},
		{
			Key:         "footer_custom_text",
			Value:       DefaultSettings.WebsiteInfo.FooterCustomText,
			Type:        "string",
			Group:       "website_info",
			Description: "页脚自定义文字",
			IsSystem:    true,
		},
		{
			Key:         "site_hero_title",
			Value:       DefaultSettings.WebsiteInfo.SiteHeroTitle,
			Type:        "string",
			Group:       "website_info",
			Description: "网站主标语",
			IsSystem:    true,
		},
		{
			Key:         "site_features_text",
			Value:       DefaultSettings.WebsiteInfo.SiteFeaturesText,
			Type:        "string",
			Group:       "website_info",
			Description: "网站特色描述",
			IsSystem:    true,
		},
	}
	allSettings = append(allSettings, websiteInfoSettings...)

	// 注册设置
	registrationSettings := []dto.SettingCreateDTO{
		{
			Key:         "enable_registration",
			Value:       DefaultSettings.Registration.EnableRegistration,
			Type:        "boolean",
			Group:       "registration",
			Description: "开放注册",
			IsSystem:    true,
		},
		{
			Key:         "email_verification",
			Value:       DefaultSettings.Registration.EmailVerification,
			Type:        "boolean",
			Group:       "registration",
			Description: "邮箱验证",
			IsSystem:    true,
		},
		{
			Key:         "user_initial_storage",
			Value:       DefaultSettings.Registration.UserInitialStorage,
			Type:        "number",
			Group:       "registration",
			Description: "新用户默认存储空间(MB)",
			IsSystem:    true,
		},
		{
			Key:         "user_initial_bandwidth",
			Value:       DefaultSettings.Registration.UserInitialBandwidth,
			Type:        "number",
			Group:       "registration",
			Description: "新用户默认带宽流量(MB)",
			IsSystem:    true,
		},
	}
	allSettings = append(allSettings, registrationSettings...)

	// AI配置
	aiSettings := []dto.SettingCreateDTO{
		{
			Key:         "ai_enabled",
			Value:       DefaultSettings.AI.AIEnabled,
			Type:        "boolean",
			Group:       "ai",
			Description: "启用AI功能",
			IsSystem:    true,
		},
		{
			Key:         "ai_auto_processing_enabled",
			Value:       DefaultSettings.AI.AIAutoProcessingEnabled,
			Type:        "boolean",
			Group:       "ai",
			Description: "AI队列自动处理开关",
			IsSystem:    true,
		},
		{
			Key:         "ai_proxy",
			Value:       DefaultSettings.AI.AIProxy,
			Type:        "string",
			Group:       "ai",
			Description: "代理地址",
			IsSystem:    true,
		},
		{
			Key:         "ai_model",
			Value:       DefaultSettings.AI.AIModel,
			Type:        "string",
			Group:       "ai",
			Description: "模型选择",
			IsSystem:    true,
		},
		{
			Key:         "ai_api_key",
			Value:       DefaultSettings.AI.AIAPIKey,
			Type:        "string",
			Group:       "ai",
			Description: "API密钥",
			IsSystem:    true,
		},
		{
			Key:         "ai_temperature",
			Value:       DefaultSettings.AI.AITemperature,
			Type:        "number",
			Group:       "ai",
			Description: "温度值",
			IsSystem:    true,
		},
		{
			Key:         "ai_max_tokens",
			Value:       DefaultSettings.AI.AIMaxTokens,
			Type:        "number",
			Group:       "ai",
			Description: "最大生成令牌数",
			IsSystem:    true,
		},
		{
			Key:         "ai_concurrency",
			Value:       DefaultSettings.AI.AIConcurrency,
			Type:        "number",
			Group:       "ai",
			Description: "AI文件识别并发处理数量",
			IsSystem:    true,
		},
		{
			Key:         "nsfw_threshold",
			Value:       DefaultSettings.AI.NSFWThreshold,
			Type:        "number",
			Group:       "ai",
			Description: "NSFW内容检测阈值(0-1，值越高越严格)",
			IsSystem:    true,
		},
		{
			Key:         "pending_stuck_threshold_minutes",
			Value:       DefaultSettings.AI.PendingStuckThresholdMins,
			Type:        "number",
			Group:       "ai",
			Description: "卡住任务超时阈值(分钟)",
			IsSystem:    true,
		},
		{
			Key:         "ai_job_retention_days",
			Value:       DefaultSettings.AI.AIJobRetentionDays,
			Type:        "number",
			Group:       "ai",
			Description: "AI任务历史保留天数",
			IsSystem:    true,
		},
	}
	allSettings = append(allSettings, aiSettings...)

	// 邮件设置
	mailSettings := []dto.SettingCreateDTO{
		{
			Key:         "smtp_host",
			Value:       DefaultSettings.Mail.SMTPHost,
			Type:        "string",
			Group:       "mail",
			Description: "SMTP服务器",
			IsSystem:    true,
		},
		{
			Key:         "smtp_port",
			Value:       DefaultSettings.Mail.SMTPPort,
			Type:        "number",
			Group:       "mail",
			Description: "SMTP端口",
			IsSystem:    true,
		},
		{
			Key:         "smtp_encryption",
			Value:       DefaultSettings.Mail.SMTPEncryption,
			Type:        "string",
			Group:       "mail",
			Description: "加密类型",
			IsSystem:    true,
		},
		{
			Key:         "smtp_username",
			Value:       DefaultSettings.Mail.SMTPUsername,
			Type:        "string",
			Group:       "mail",
			Description: "SMTP用户名",
			IsSystem:    true,
		},
		{
			Key:         "smtp_password",
			Value:       DefaultSettings.Mail.SMTPPassword,
			Type:        "string",
			Group:       "mail",
			Description: "SMTP密码",
			IsSystem:    true,
		},
		{
			Key:         "smtp_from_address",
			Value:       DefaultSettings.Mail.SMTPFromAddr,
			Type:        "string",
			Group:       "mail",
			Description: "发件人地址",
			IsSystem:    true,
		},
		{
			Key:         "smtp_from_name",
			Value:       DefaultSettings.Mail.SMTPFromName,
			Type:        "string",
			Group:       "mail",
			Description: "发件人名称",
			IsSystem:    true,
		},
		{
			Key:         "smtp_reply_to",
			Value:       DefaultSettings.Mail.SMTPReplyTo,
			Type:        "string",
			Group:       "mail",
			Description: "回复地址",
			IsSystem:    true,
		},
	}
	allSettings = append(allSettings, mailSettings...)

	// 上传设置
	uploadSettings := []dto.SettingCreateDTO{
		{
			Key:         "allowed_file_formats",
			Value:       DefaultSettings.Upload.AllowedFileFormats,
			Type:        "array",
			Group:       "upload",
			Description: "允许上传的文件格式",
			IsSystem:    true,
		},
		{
			Key:         "max_file_size",
			Value:       DefaultSettings.Upload.MaxFileSize,
			Type:        "number",
			Group:       "upload",
			Description: "单个文件最大大小(MB)",
			IsSystem:    true,
		},
		{
			Key:         "max_batch_size",
			Value:       DefaultSettings.Upload.MaxBatchSize,
			Type:        "number",
			Group:       "upload",
			Description: "批量上传最大总大小(MB)",
			IsSystem:    true,
		},
		{
			Key:         "thumbnail_max_width",
			Value:       DefaultSettings.Upload.ThumbnailMaxWidth,
			Type:        "number",
			Group:       "upload",
			Description: "缩略图最大宽度",
			IsSystem:    true,
		},
		{
			Key:         "thumbnail_max_height",
			Value:       DefaultSettings.Upload.ThumbnailMaxHeight,
			Type:        "number",
			Group:       "upload",
			Description: "缩略图最大高度",
			IsSystem:    true,
		},
		{
			Key:         "thumbnail_quality",
			Value:       DefaultSettings.Upload.ThumbnailQuality,
			Type:        "number",
			Group:       "upload",
			Description: "缩略图质量设置(0-100)",
			IsSystem:    true,
		},
		{
			Key:         "preserve_exif",
			Value:       DefaultSettings.Upload.PreserveEXIF,
			Type:        "boolean",
			Group:       "upload",
			Description: "是否保留EXIF信息",
			IsSystem:    true,
		},
		{
			Key:         "daily_upload_limit",
			Value:       DefaultSettings.Upload.DailyUploadLimit,
			Type:        "number",
			Group:       "upload",
			Description: "用户每日上传数量限制",
			IsSystem:    true,
		},
		{
			Key:         "client_max_concurrent_uploads",
			Value:       DefaultSettings.Upload.ClientMaxConcurrentUploads,
			Type:        "number",
			Group:       "upload",
			Description: "客户端最大并发上传数",
			IsSystem:    true,
		},
		// 分片上传相关设置
		{
			Key:         "chunked_upload_enabled",
			Value:       DefaultSettings.Upload.ChunkedUploadEnabled,
			Type:        "boolean",
			Group:       "upload",
			Description: "启用分片上传",
			IsSystem:    true,
		},
		{
			Key:         "chunked_threshold",
			Value:       DefaultSettings.Upload.ChunkedThreshold,
			Type:        "number",
			Group:       "upload",
			Description: "分片上传阈值(MB)",
			IsSystem:    true,
		},
		{
			Key:         "chunk_size",
			Value:       DefaultSettings.Upload.ChunkSize,
			Type:        "number",
			Group:       "upload",
			Description: "分片大小(MB)",
			IsSystem:    true,
		},
		{
			Key:         "max_concurrency",
			Value:       DefaultSettings.Upload.MaxConcurrency,
			Type:        "number",
			Group:       "upload",
			Description: "分片最大并发数",
			IsSystem:    true,
		},
		{
			Key:         "retry_count",
			Value:       DefaultSettings.Upload.RetryCount,
			Type:        "number",
			Group:       "upload",
			Description: "分片重试次数",
			IsSystem:    true,
		},
		{
			Key:         "session_timeout",
			Value:       DefaultSettings.Upload.SessionTimeout,
			Type:        "number",
			Group:       "upload",
			Description: "分片上传会话超时(小时)",
			IsSystem:    true,
		},
		{
			Key:         "cleanup_interval",
			Value:       DefaultSettings.Upload.CleanupInterval,
			Type:        "number",
			Group:       "upload",
			Description: "分片会话清理间隔(分钟)",
			IsSystem:    true,
		},
		// 内容检测设置
		{
			Key:         "content_detection_enabled",
			Value:       DefaultSettings.Upload.ContentDetectionEnabled,
			Type:        "boolean",
			Group:       "upload",
			Description: "违规内容检测开关",
			IsSystem:    true,
		},
		{
			Key:         "sensitive_content_handling",
			Value:       DefaultSettings.Upload.SensitiveContentHandling,
			Type:        "string",
			Group:       "upload",
			Description: "敏感文件处理方式(auto_delete:自动删除, mark_only:仅标记, pending_review:等待审核)",
			IsSystem:    true,
		},
		{
			Key:         "ai_analysis_enabled",
			Value:       DefaultSettings.Upload.AIAnalysisEnabled,
			Type:        "boolean",
			Group:       "upload",
			Description: "AI智能分析开关(需配置AI模型)",
			IsSystem:    true,
		},
		// 存储时长设置
		{
			Key:         "user_allowed_storage_durations",
			Value:       DefaultSettings.Upload.UserAllowedStorageDurations,
			Type:        "array",
			Group:       "upload",
			Description: "已登录用户可选择的存储时长选项（permanent为内置选项）",
			IsSystem:    true,
		},
		{
			Key:         "user_default_storage_duration",
			Value:       DefaultSettings.Upload.UserDefaultStorageDuration,
			Type:        "string",
			Group:       "upload",
			Description: "已登录用户默认存储时长",
			IsSystem:    true,
		},
		// 秒传检测设置
		{
			Key:         "instant_upload_enabled",
			Value:       false,
			Type:        "boolean",
			Group:       "upload",
			Description: "检测上传图片是否重复实现秒传",
			IsSystem:    true,
		},
	}
	allSettings = append(allSettings, uploadSettings...)

	// 网站装修设置
	themeSettings := []dto.SettingCreateDTO{
		{
			Key:         "site_mode",
			Value:       DefaultSettings.Theme.SiteMode,
			Type:        "string",
			Group:       "theme",
			Description: "网站显示模式(website:传统网站模式, personal:个人工具模式, minimal:极简工具模式)",
			IsSystem:    true,
		},
	}
	allSettings = append(allSettings, themeSettings...)

	// 访客控制设置
	guestSettings := []dto.SettingCreateDTO{
		{
			Key:         "enable_guest_upload",
			Value:       DefaultSettings.Guest.EnableGuestUpload,
			Type:        "boolean",
			Group:       "guest",
			Description: "是否开放游客上传",
			IsSystem:    true,
		},
		{
			Key:         "guest_daily_limit",
			Value:       DefaultSettings.Guest.GuestDailyLimit,
			Type:        "number",
			Group:       "guest",
			Description: "游客每日上传次数限制",
			IsSystem:    true,
		},
		{
			Key:         "guest_default_access_level",
			Value:       DefaultSettings.Guest.GuestDefaultAccessLevel,
			Type:        "string",
			Group:       "guest",
			Description: "默认访问级别",
			IsSystem:    true,
		},
		{
			Key:         "guest_allowed_storage_durations",
			Value:       DefaultSettings.Guest.GuestAllowedStorageDurations,
			Type:        "array",
			Group:       "guest",
			Description: "允许的存储时长选项",
			IsSystem:    true,
		},
		{
			Key:         "guest_default_storage_duration",
			Value:       DefaultSettings.Guest.GuestDefaultStorageDuration,
			Type:        "string",
			Group:       "guest",
			Description: "默认存储时长",
			IsSystem:    true,
		},
		{
			Key:         "guest_ip_daily_limit",
			Value:       DefaultSettings.Guest.GuestIPDailyLimit,
			Type:        "number",
			Group:       "guest",
			Description: "IP每日限制，防止通过刷新指纹盗刷",
			IsSystem:    true,
		},
	}
	allSettings = append(allSettings, guestSettings...)

	// 安全设置
	securitySettings := []dto.SettingCreateDTO{
		{
			Key:         "max_login_attempts",
			Value:       DefaultSettings.Security.MaxLoginAttempts,
			Type:        "number",
			Group:       "security",
			Description: "最大登录尝试次数",
			IsSystem:    true,
		},
		{
			Key:         "account_lockout_minutes",
			Value:       DefaultSettings.Security.AccountLockoutMinutes,
			Type:        "number",
			Group:       "security",
			Description: "账户锁定分钟数",
			IsSystem:    true,
		},
		{
			Key:         "login_expire_hours",
			Value:       DefaultSettings.Security.LoginExpireHours,
			Type:        "number",
			Group:       "security",
			Description: "登录有效期(小时)",
			IsSystem:    true,
		},
		{
			Key:         "jwt_secret",
			Value:       generateRandomJWTSecret(),
			Type:        "string",
			Group:       "security",
			Description: "JWT签名密钥",
			IsSystem:    true,
		},
		{
			Key:         "hide_remote_url",
			Value:       DefaultSettings.Security.HideRemoteURL,
			Type:        "boolean",
			Group:       "security",
			Description: "是否隐藏远程存储地址",
			IsSystem:    true,
		},
		{
			Key:         "ip_whitelist",
			Value:       DefaultSettings.Security.IPWhitelist,
			Type:        "string",
			Group:       "security",
			Description: "IP白名单",
			IsSystem:    true,
		},
		{
			Key:         "ip_blacklist",
			Value:       DefaultSettings.Security.IPBlacklist,
			Type:        "string",
			Group:       "security",
			Description: "IP黑名单",
			IsSystem:    true,
		},
		{
			Key:         "domain_whitelist",
			Value:       DefaultSettings.Security.DomainWhitelist,
			Type:        "string",
			Group:       "security",
			Description: "域名白名单",
			IsSystem:    true,
		},
		{
			Key:         "domain_blacklist",
			Value:       DefaultSettings.Security.DomainBlacklist,
			Type:        "string",
			Group:       "security",
			Description: "域名黑名单",
			IsSystem:    true,
		},
	}
	allSettings = append(allSettings, securitySettings...)

	// 向量搜索设置
	vectorSettings := []dto.SettingCreateDTO{
		{
			Key:         "vector_enabled",
			Value:       DefaultSettings.Vector.VectorEnabled,
			Type:        "boolean",
			Group:       "vector",
			Description: "启用向量搜索功能",
			IsSystem:    true,
		},
		{
			Key:         "vector_auto_processing_enabled",
			Value:       DefaultSettings.Vector.VectorAutoProcessingEnabled,
			Type:        "boolean",
			Group:       "vector",
			Description: "向量队列自动处理开关",
			IsSystem:    true,
		},
		{
			Key:         "vector_provider",
			Value:       DefaultSettings.Vector.VectorProvider,
			Type:        "string",
			Group:       "vector",
			Description: "向量化提供者",
			IsSystem:    true,
		},
		{
			Key:         "vector_model",
			Value:       DefaultSettings.Vector.VectorModel,
			Type:        "string",
			Group:       "vector",
			Description: "向量化模型",
			IsSystem:    true,
		},
		{
			Key:         "vector_api_key",
			Value:       DefaultSettings.Vector.VectorAPIKey,
			Type:        "string",
			Group:       "vector",
			Description: "OpenAI API密钥",
			IsSystem:    true,
		},
		{
			Key:         "vector_base_url",
			Value:       DefaultSettings.Vector.VectorBaseURL,
			Type:        "string",
			Group:       "vector",
			Description: "OpenAI API代理地址",
			IsSystem:    true,
		},
		{
			Key:         "vector_timeout",
			Value:       DefaultSettings.Vector.VectorTimeout,
			Type:        "number",
			Group:       "vector",
			Description: "API调用超时时间(秒)",
			IsSystem:    true,
		},
		{
			Key:         "vector_similarity_threshold",
			Value:       DefaultSettings.Vector.VectorSimilarityThreshold,
			Type:        "number",
			Group:       "vector",
			Description: "相似图阈值(0-1)",
			IsSystem:    true,
		},
		{
			Key:         "vector_search_threshold",
			Value:       DefaultSettings.Vector.VectorSearchThreshold,
			Type:        "number",
			Group:       "vector",
			Description: "搜索阈值(0-1)",
			IsSystem:    true,
		},
		{
			Key:         "vector_max_results",
			Value:       DefaultSettings.Vector.VectorMaxResults,
			Type:        "number",
			Group:       "vector",
			Description: "最大搜索结果数",
			IsSystem:    true,
		},
		{
			Key:         "vector_concurrency",
			Value:       DefaultSettings.Vector.VectorConcurrency,
			Type:        "number",
			Group:       "vector",
			Description: "向量生成并发数量",
			IsSystem:    true,
		},
	}
	allSettings = append(allSettings, vectorSettings...)

	// 版本信息设置
	versionSettings := []dto.SettingCreateDTO{
		{
			Key:         "current_version",
			Value:       DefaultSettings.Version.CurrentVersion,
			Type:        "string",
			Group:       "version",
			Description: "当前系统版本",
			IsSystem:    true,
		},
		{
			Key:         "build_time",
			Value:       DefaultSettings.Version.BuildTime,
			Type:        "string",
			Group:       "version",
			Description: "构建时间",
			IsSystem:    true,
		},
		{
			Key:         "git_commit",
			Value:       DefaultSettings.Version.GitCommit,
			Type:        "string",
			Group:       "version",
			Description: "Git提交哈希",
			IsSystem:    true,
		},
		{
			Key:         "update_available",
			Value:       DefaultSettings.Version.UpdateAvailable,
			Type:        "boolean",
			Group:       "version",
			Description: "是否有可用更新",
			IsSystem:    true,
		},
		{
			Key:         "last_update_check",
			Value:       DefaultSettings.Version.LastUpdateCheck,
			Type:        "string",
			Group:       "version",
			Description: "最后更新检查时间",
			IsSystem:    true,
		},
		{
			Key:         "last_update_time",
			Value:       DefaultSettings.Version.LastUpdateTime,
			Type:        "string",
			Group:       "version",
			Description: "最后更新时间",
			IsSystem:    true,
		},
		{
			Key:         "update_logs",
			Value:       DefaultSettings.Version.UpdateLogs,
			Type:        "string",
			Group:       "version",
			Description: "更新日志",
			IsSystem:    true,
		},
	}
	allSettings = append(allSettings, versionSettings...)

	// 外观界面设置
	appearanceSettings := []dto.SettingCreateDTO{
		{
			Key:         "show_official_site",
			Value:       DefaultSettings.Appearance.ShowOfficialSite,
			Type:        "boolean",
			Group:       "appearance",
			Description: "是否显示官网链接",
			IsSystem:    true,
		},
		{
			Key:         "official_site_url",
			Value:       DefaultSettings.Appearance.OfficialSiteURL,
			Type:        "string",
			Group:       "appearance",
			Description: "官网地址",
			IsSystem:    true,
		},
		{
			Key:         "show_github_link",
			Value:       DefaultSettings.Appearance.ShowGitHubLink,
			Type:        "boolean",
			Group:       "appearance",
			Description: "是否显示GitHub链接",
			IsSystem:    true,
		},
		{
			Key:         "github_url",
			Value:       DefaultSettings.Appearance.GitHubURL,
			Type:        "string",
			Group:       "appearance",
			Description: "GitHub仓库地址",
			IsSystem:    true,
		},
		{
			Key:         "show_wechat_group",
			Value:       DefaultSettings.Appearance.ShowWeChatGroup,
			Type:        "boolean",
			Group:       "appearance",
			Description: "是否显示微信群",
			IsSystem:    true,
		},
		{
			Key:         "wechat_qr_image_url",
			Value:       DefaultSettings.Appearance.WeChatQRImageURL,
			Type:        "string",
			Group:       "appearance",
			Description: "微信群二维码文件地址",
			IsSystem:    true,
		},
		{
			Key:         "wechat_contact_account",
			Value:       DefaultSettings.Appearance.WeChatContactAcct,
			Type:        "string",
			Group:       "appearance",
			Description: "微信联系人账号",
			IsSystem:    true,
		},
		{
			Key:         "show_qq_group",
			Value:       DefaultSettings.Appearance.ShowQQGroup,
			Type:        "boolean",
			Group:       "appearance",
			Description: "是否显示QQ群",
			IsSystem:    true,
		},
		{
			Key:         "qq_qr_image_url",
			Value:       DefaultSettings.Appearance.QQQRImageURL,
			Type:        "string",
			Group:       "appearance",
			Description: "QQ群二维码文件地址",
			IsSystem:    true,
		},
		{
			Key:         "qq_group_number",
			Value:       DefaultSettings.Appearance.QQGroupNumber,
			Type:        "string",
			Group:       "appearance",
			Description: "QQ群号",
			IsSystem:    true,
		},
		{
			Key:         "enable_multi_layout",
			Value:       DefaultSettings.Appearance.EnableMultiLayout,
			Type:        "boolean",
			Group:       "appearance",
			Description: "是否开启多布局功能",
			IsSystem:    true,
		},
		{
			Key:         "default_layout",
			Value:       DefaultSettings.Appearance.DefaultLayout,
			Type:        "string",
			Group:       "appearance",
			Description: "默认布局 (top/sidebar)",
			IsSystem:    true,
		},

		{
			Key:         "enable_multi_language",
			Value:       DefaultSettings.Appearance.EnableMultiLanguage,
			Type:        "boolean",
			Group:       "appearance",
			Description: "是否开启多语言功能",
			IsSystem:    true,
		},
		{
			Key:         "default_language",
			Value:       DefaultSettings.Appearance.DefaultLanguage,
			Type:        "string",
			Group:       "appearance",
			Description: "默认语言 (zh-CN/en/ja等)",
			IsSystem:    true,
		},
	}
	allSettings = append(allSettings, appearanceSettings...)

	// 公告系统配置
	announcementSettings := []dto.SettingCreateDTO{
		{
			Key:         "announcement_enabled",
			Value:       DefaultSettings.Announcement.AnnouncementEnabled,
			Type:        "boolean",
			Group:       "announcement",
			Description: "公告系统总开关",
			IsSystem:    true,
		},
		{
			Key:         "announcement_drawer_position",
			Value:       DefaultSettings.Announcement.AnnouncementDrawerPos,
			Type:        "string",
			Group:       "announcement",
			Description: "抽屉位置(left/right)",
			IsSystem:    true,
		},
		{
			Key:         "announcement_display_limit",
			Value:       DefaultSettings.Announcement.AnnouncementDisplayLimit,
			Type:        "number",
			Group:       "announcement",
			Description: "显示最近N条公告",
			IsSystem:    true,
		},
		{
			Key:         "announcement_auto_show_delay",
			Value:       DefaultSettings.Announcement.AnnouncementAutoShowDelay,
			Type:        "number",
			Group:       "announcement",
			Description: "置顶公告自动弹窗延迟时间（秒）",
			IsSystem:    true,
		},
	}
	allSettings = append(allSettings, announcementSettings...)

	// 法律文档设置 - 使用预定义模板
	legalSettings := []dto.SettingCreateDTO{
		{
			Key:         "privacy_policy_content",
			Value:       DefaultPrivacyPolicyTemplate,
			Type:        "text",
			Group:       "legal",
			Description: "隐私政策内容",
			IsSystem:    true,
		},
		{
			Key:         "terms_of_service_content",
			Value:       DefaultTermsOfServiceTemplate,
			Type:        "text",
			Group:       "legal",
			Description: "服务条款内容",
			IsSystem:    true,
		},
	}
	allSettings = append(allSettings, legalSettings...)

	// 批量插入或更新设置
	settingsDTO := &dto.BatchUpsertSettingDTO{
		Settings: allSettings,
	}

	// 调用设置服务批量创建/更新设置
	result, err := setting.BatchUpsertSettings(settingsDTO)
	if err != nil {
		return fmt.Errorf("初始化系统设置失败: %v", err)
	}

	// 记录结果
	logger.Infof("设置初始化完成! 成功: %d, 失败: %d", len(result.Success), len(result.Failed))

	// 如果有失败的项，记录详细信息
	if len(result.Failed) > 0 {
		for _, failedItem := range result.Failed {
			logger.Errorf("设置 %s 创建失败: %s", failedItem.Key, failedItem.Message)
		}
	}

	// 记录此迁移已执行
	if err := RecordMigration(db, migrationName); err != nil {
		return fmt.Errorf("记录迁移状态失败: %v", err)
	}

	// 标记其他相关迁移为已完成，以避免重复执行
	migrationNames := []string{"add_image_access_tokens", "update_image_access_tokens"}
	for _, name := range migrationNames {
		if migrationRecord, _ := IsMigrationApplied(db, name); !migrationRecord {
			if err := RecordMigration(db, name); err != nil {
				logger.Errorf("标记迁移 %s 为已完成时失败: %v", name, err)
			}
		}
	}

	logger.Infof("已创建内部服务令牌: %s", internalServiceToken)

	// 初始化分类模板
	if err := initializeCategoryTemplates(db); err != nil {
		logger.Errorf("初始化分类模板失败: %v", err)
		// 不返回错误，因为分类模板失败不应该阻止系统初始化
	}

	// 初始化欢迎公告
	if err := initializeWelcomeAnnouncement(db); err != nil {
		logger.Errorf("初始化欢迎公告失败: %v", err)
		// 不返回错误，因为公告失败不应该阻止系统初始化
	}

	return nil
}

// initializeCategoryTemplates 初始化默认分类模板
func initializeCategoryTemplates(db *gorm.DB) error {
	logger.Infof("正在初始化默认分类模板...")

	// 检查是否已有分类模板
	var count int64
	if err := db.Model(&models.CategoryTemplate{}).Count(&count).Error; err != nil {
		return fmt.Errorf("检查分类模板数量失败: %v", err)
	}

	if count > 0 {
		logger.Infof("分类模板已存在 (%d 个)，跳过初始化", count)
		return nil
	}

	// 从配置文件加载默认分类模板
	defaultTemplates := make([]models.CategoryTemplate, 0, len(DefaultCategoryTemplates))
	for _, tmpl := range DefaultCategoryTemplates {
		defaultTemplates = append(defaultTemplates, models.CategoryTemplate{
			Name:        tmpl.Name,
			Description: tmpl.Description,
			Icon:        tmpl.Icon,
			IsPopular:   tmpl.IsPopular,
			SortOrder:   tmpl.SortOrder,
			UsageCount:  0,
		})
	}

	// 批量插入分类模板
	if err := db.Create(&defaultTemplates).Error; err != nil {
		return fmt.Errorf("创建默认分类模板失败: %v", err)
	}

	logger.Infof("成功创建 %d 个默认分类模板", len(defaultTemplates))
	return nil
}

// initializeWelcomeAnnouncement 初始化欢迎公告
func initializeWelcomeAnnouncement(db *gorm.DB) error {
	logger.Infof("正在初始化欢迎公告...")

	// 检查是否已有公告
	var count int64
	if err := db.Model(&models.Announcement{}).Count(&count).Error; err != nil {
		return fmt.Errorf("检查公告数量失败: %v", err)
	}

	if count > 0 {
		logger.Infof("公告已存在 (%d 条)，跳过初始化", count)
		return nil
	}

	// 创建欢迎公告（内容来自配置文件）
	announcement := models.Announcement{
		Title:     "PixelPunk 安装完成",
		Content:   WelcomeAnnouncementContent,
		Summary:   "🎉 欢迎使用 PixelPunk！这是一个功能强大的AI智能图床系统，快来探索70+核心功能吧！",
		IsPinned:  true,
		Status:    "published",
		ViewCount: 0,
		CreatedBy: 1, // 管理员ID
	}

	// 插入欢迎公告
	if err := db.Create(&announcement).Error; err != nil {
		return fmt.Errorf("创建欢迎公告失败: %v", err)
	}

	logger.Infof("成功创建欢迎公告 (ID: %d)", announcement.ID)
	return nil
}
