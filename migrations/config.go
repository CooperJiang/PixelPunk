package migrations

// DefaultSettings 默认配置值统一管理
var DefaultSettings = struct {
	Website      WebsiteSettings
	WebsiteInfo  WebsiteInfoSettings
	Registration RegistrationSettings
	AI           AISettings
	Mail         MailSettings
	Upload       UploadSettings
	Theme        ThemeSettings
	Guest        GuestSettings
	Security     SecuritySettings
	Vector       VectorSettings
	Version      VersionSettings
	Appearance   AppearanceSettings
	Announcement AnnouncementSettings
}{
	Website: WebsiteSettings{
		AdminEmail:  "",
		SiteBaseURL: "",
	},

	WebsiteInfo: WebsiteInfoSettings{
		SiteName:         "Pixel Punk",
		SiteDescription:  "高效的文件管理与分享平台",
		SiteKeywords:     "图床,文件分享,存储,PixelPunk,AI图床,智能图床,自动化图床",
		ICPNumber:        "",
		ShowFileCount:    true,
		ShowStorageUsage: true,
		SiteLogoURL:      "",
		FaviconURL:       "",
		CopyrightText:    "© 2024 PixelPunk. All rights reserved.",
		ContactEmail:     "",
		FooterCustomText: "",
		SiteHeroTitle:    "让图片管理从繁琐到简单，让文件分享从等待到极速",
		SiteFeaturesText: "AI自动识别 · 智能分类整理 · 秒传极速分享 · 开源社区驱动",
	},

	Registration: RegistrationSettings{
		EnableRegistration:   true,
		EmailVerification:    true,
		UserInitialStorage:   1024,  // MB
		UserInitialBandwidth: 10240, // MB
	},

	AI: AISettings{
		AIEnabled:                 false,
		AIAutoProcessingEnabled:   false,
		AIProxy:                   "https://api.openai.com/v1",
		AIModel:                   "gpt-5-mini",
		AIAPIKey:                  "sk-xxxxxxxxxxxxxxx",
		AITemperature:             0.1,
		AIMaxTokens:               16000,
		AIConcurrency:             5,
		NSFWThreshold:             0.6,
		PendingStuckThresholdMins: 30,
		AIJobRetentionDays:        14,
	},

	Mail: MailSettings{
		SMTPHost:       "smtp.example.com",
		SMTPPort:       587,
		SMTPEncryption: "none",
		SMTPUsername:   "noreply@example.com",
		SMTPPassword:   "",
		SMTPFromAddr:   "noreply@example.com",
		SMTPFromName:   "PixelPunk",
		SMTPReplyTo:    "noreply@example.com",
	},

	Upload: UploadSettings{
		AllowedFileFormats: []string{
			"jpg", "jpeg", "png", "gif", "webp", "bmp", "svg", "ico",
			"apng", "jp2", "tiff", "tif", "tga", "heic", "heif",
		},
		MaxFileSize:                 20,
		MaxBatchSize:                100,
		ThumbnailMaxWidth:           1000,
		ThumbnailMaxHeight:          800,
		ThumbnailQuality:            80,
		PreserveEXIF:                true,
		DailyUploadLimit:            1000,
		ClientMaxConcurrentUploads:  5,
		ChunkedUploadEnabled:        true,
		ChunkedThreshold:            10,
		ChunkSize:                   2,
		MaxConcurrency:              3,
		RetryCount:                  3,
		SessionTimeout:              24,
		CleanupInterval:             60,
		ContentDetectionEnabled:     true,
		SensitiveContentHandling:    "mark_only",
		AIAnalysisEnabled:           true,
		UserAllowedStorageDurations: []string{"1h", "3d", "7d", "30d", "permanent"},
		UserDefaultStorageDuration:  "permanent",
	},

	Theme: ThemeSettings{
		SiteMode: "website", // website/personal/minimal
	},

	Guest: GuestSettings{
		EnableGuestUpload:            true,
		GuestDailyLimit:              10,
		GuestDefaultAccessLevel:      "public",
		GuestAllowedStorageDurations: []string{"3d", "7d", "30d"},
		GuestDefaultStorageDuration:  "7d",
		GuestIPDailyLimit:            50,
	},

	Security: SecuritySettings{
		MaxLoginAttempts:      5,
		AccountLockoutMinutes: 30,
		LoginExpireHours:      60,
		HideRemoteURL:         true,
		IPWhitelist:           "",
		IPBlacklist:           "",
		DomainWhitelist:       "",
		DomainBlacklist:       "",
	},

	Vector: VectorSettings{
		VectorEnabled:               true,
		VectorAutoProcessingEnabled: true,
		VectorProvider:              "openai",
		VectorModel:                 "text-embedding-3-small",
		VectorAPIKey:                "",
		VectorBaseURL:               "https://api.openai.com/v1",
		VectorTimeout:               30,
		VectorSimilarityThreshold:   0.7,
		VectorSearchThreshold:       0.36,
		VectorMaxResults:            100,
		VectorConcurrency:           3,
	},

	Version: VersionSettings{
		CurrentVersion:  "1.0.0",
		BuildTime:       "",
		GitCommit:       "",
		UpdateAvailable: false,
		LastUpdateCheck: "",
		LastUpdateTime:  "",
		UpdateLogs:      "",
	},

	Appearance: AppearanceSettings{
		ShowOfficialSite:    true,
		OfficialSiteURL:     "https://pixelpunk.cc/",
		ShowGitHubLink:      true,
		GitHubURL:           "https://github.com/CooperJiang/PixelPunk",
		ShowWeChatGroup:     true,
		WeChatQRImageURL:    "",
		WeChatContactAcct:   "",
		ShowQQGroup:         true,
		QQQRImageURL:        "",
		QQGroupNumber:       "",
		EnableMultiLayout:   true,
		DefaultLayout:       "top",
		EnableMultiLanguage: false,
		DefaultLanguage:     "auto",
	},

	Announcement: AnnouncementSettings{
		AnnouncementEnabled:       true,
		AnnouncementDrawerPos:     "right",
		AnnouncementDisplayLimit:  10,
		AnnouncementAutoShowDelay: 2, // 秒
	},
}

// WebsiteSettings 网站后端功能设置
type WebsiteSettings struct {
	AdminEmail  string
	SiteBaseURL string
}

// WebsiteInfoSettings 网站信息配置（前端显示）
type WebsiteInfoSettings struct {
	SiteName         string
	SiteDescription  string
	SiteKeywords     string
	ICPNumber        string
	ShowFileCount    bool
	ShowStorageUsage bool
	SiteLogoURL      string
	FaviconURL       string
	CopyrightText    string
	ContactEmail     string
	FooterCustomText string
	SiteHeroTitle    string
	SiteFeaturesText string
}

// RegistrationSettings 注册设置
type RegistrationSettings struct {
	EnableRegistration   bool
	EmailVerification    bool
	UserInitialStorage   int // MB
	UserInitialBandwidth int // MB
}

// AISettings AI配置
type AISettings struct {
	AIEnabled                 bool
	AIAutoProcessingEnabled   bool
	AIProxy                   string
	AIModel                   string
	AIAPIKey                  string
	AITemperature             float64
	AIMaxTokens               int
	AIConcurrency             int
	NSFWThreshold             float64
	PendingStuckThresholdMins int
	AIJobRetentionDays        int
}

// MailSettings 邮件设置
type MailSettings struct {
	SMTPHost       string
	SMTPPort       int
	SMTPEncryption string
	SMTPUsername   string
	SMTPPassword   string
	SMTPFromAddr   string
	SMTPFromName   string
	SMTPReplyTo    string
}

// UploadSettings 上传设置
type UploadSettings struct {
	AllowedFileFormats          []string
	MaxFileSize                 int
	MaxBatchSize                int
	ThumbnailMaxWidth           int
	ThumbnailMaxHeight          int
	ThumbnailQuality            int
	PreserveEXIF                bool
	DailyUploadLimit            int
	ClientMaxConcurrentUploads  int
	ChunkedUploadEnabled        bool
	ChunkedThreshold            int
	ChunkSize                   int
	MaxConcurrency              int
	RetryCount                  int
	SessionTimeout              int
	CleanupInterval             int
	ContentDetectionEnabled     bool
	SensitiveContentHandling    string
	AIAnalysisEnabled           bool
	UserAllowedStorageDurations []string
	UserDefaultStorageDuration  string
}

// ThemeSettings 网站装修设置
type ThemeSettings struct {
	SiteMode string // website/personal/minimal
}

// GuestSettings 访客控制设置
type GuestSettings struct {
	EnableGuestUpload            bool
	GuestDailyLimit              int
	GuestDefaultAccessLevel      string
	GuestAllowedStorageDurations []string
	GuestDefaultStorageDuration  string
	GuestIPDailyLimit            int
}

// SecuritySettings 安全设置
type SecuritySettings struct {
	MaxLoginAttempts      int
	AccountLockoutMinutes int
	LoginExpireHours      int
	HideRemoteURL         bool
	IPWhitelist           string
	IPBlacklist           string
	DomainWhitelist       string
	DomainBlacklist       string
}

// VectorSettings 向量搜索设置
type VectorSettings struct {
	VectorEnabled               bool
	VectorAutoProcessingEnabled bool
	VectorProvider              string
	VectorModel                 string
	VectorAPIKey                string
	VectorBaseURL               string
	VectorTimeout               int
	VectorSimilarityThreshold   float64
	VectorSearchThreshold       float64
	VectorMaxResults            int
	VectorConcurrency           int
}

// VersionSettings 版本信息设置
type VersionSettings struct {
	CurrentVersion  string
	BuildTime       string
	GitCommit       string
	UpdateAvailable bool
	LastUpdateCheck string
	LastUpdateTime  string
	UpdateLogs      string
}

// AppearanceSettings 外观界面设置
type AppearanceSettings struct {
	ShowOfficialSite    bool
	OfficialSiteURL     string
	ShowGitHubLink      bool
	GitHubURL           string
	ShowWeChatGroup     bool
	WeChatQRImageURL    string
	WeChatContactAcct   string
	ShowQQGroup         bool
	QQQRImageURL        string
	QQGroupNumber       string
	EnableMultiLayout   bool
	DefaultLayout       string
	EnableMultiLanguage bool
	DefaultLanguage     string
}

// AnnouncementSettings 公告系统配置
type AnnouncementSettings struct {
	AnnouncementEnabled       bool
	AnnouncementDrawerPos     string // left/right
	AnnouncementDisplayLimit  int
	AnnouncementAutoShowDelay int // 秒
}

// CategoryTemplateConfig 分类模板配置
type CategoryTemplateConfig struct {
	Name        string
	Description string
	Icon        string
	IsPopular   bool
	SortOrder   int
}

// DefaultCategoryTemplates 默认分类模板（通用生活分类）
var DefaultCategoryTemplates = []CategoryTemplateConfig{
	{
		Name:        "人物肖像",
		Description: "人物照片、头像、肖像等",
		Icon:        "fas fa-user",
		IsPopular:   true,
		SortOrder:   10,
	},
	{
		Name:        "自然风景",
		Description: "山川、湖泊、天空、海洋等自然景观",
		Icon:        "fas fa-mountain",
		IsPopular:   true,
		SortOrder:   20,
	},
	{
		Name:        "城市建筑",
		Description: "建筑物、街道、城市景观等",
		Icon:        "fas fa-building",
		IsPopular:   true,
		SortOrder:   30,
	},
	{
		Name:        "动物宠物",
		Description: "各种动物、宠物照片",
		Icon:        "fas fa-paw",
		IsPopular:   true,
		SortOrder:   40,
	},
	{
		Name:        "美食料理",
		Description: "食物、饮品、烹饪相关",
		Icon:        "fas fa-utensils",
		IsPopular:   true,
		SortOrder:   50,
	},
	{
		Name:        "植物花卉",
		Description: "花朵、树木、植物等",
		Icon:        "fas fa-leaf",
		IsPopular:   false,
		SortOrder:   60,
	},
	{
		Name:        "交通工具",
		Description: "汽车、飞机、船只等交通工具",
		Icon:        "fas fa-car",
		IsPopular:   false,
		SortOrder:   70,
	},
	{
		Name:        "生活物品",
		Description: "日常用品、家具、电子设备等",
		Icon:        "fas fa-home",
		IsPopular:   false,
		SortOrder:   80,
	},
	{
		Name:        "艺术创作",
		Description: "绘画、雕塑、设计作品等",
		Icon:        "fas fa-palette",
		IsPopular:   false,
		SortOrder:   90,
	},
	{
		Name:        "其他杂项",
		Description: "不属于其他分类的文件",
		Icon:        "fas fa-folder",
		IsPopular:   false,
		SortOrder:   100,
	},
}

// WelcomeAnnouncementContent 欢迎公告内容模板
const WelcomeAnnouncementContent = `# 🎉 欢迎来到 PixelPunk

您已成功安装 **PixelPunk 开源图床系统**！

这是一个功能极其全面且强大、界面精美的现代化AI智能图片管理平台，专为个人和团队打造。

---

## 📚 重要资源

| 类型        | 链接                                                              | 说明                         |
| ----------- | ----------------------------------------------------------------- | ---------------------------- |
| 📖 官方文档 | [pixelpunk.cc](https://pixelpunk.cc)                              | 完整使用教程与API文档        |
| 💻 GitHub   | [CooperJiang/PixelPunk](https://github.com/CooperJiang/PixelPunk) | 开源代码，欢迎 Star & Fork   |
| 💬 QQ交流群 | **826708512**                                                     | 问题反馈、功能建议、技术交流 |

---

## ✨ 核心特性

### 🤖 AI 智能功能

- **AI自动审核** - NSFW内容审核，智能识别违规敏感图片
- **AI自动标签** - 智能识别图片内容，自动打标签分类
- **AI自动分类** - 智能识别图片场景，自动归档整理
- **向量搜索** - 基于AI的智能搜索引擎，自然语言即可搜索
- **以图搜图** - 上传图片查找相似内容，搜索近似图片
- **内容审核** - 多级审核机制，保护平台内容安全

### 🎨 主题与界面

- **10+精美主题** - 赛博朋克、简约、炫彩等风格任选
- **双风格切换** - 常规风格与赛博朋克风格无缝切换
- **多语言支持** - 中文、英文等多种语言界面
- **暗黑/明亮模式** - 护眼舒适，随心切换
- **响应式设计** - 完美适配桌面、平板、手机

### 📁 文件管理

- **文件夹组织** - 多级目录结构，井井有条
- **批量操作** - 上传、下载、移动、删除一键完成
- **标签系统** - 灵活标记，快速筛选定位
- **分类管理** - 自定义分类，智能归档
- **作者主页** - 展示个人作品，构建创作者档案
- **秒传优化** - 智能去重，节省空间与带宽

### 🔗 分享与协作

- **共享资源** - 创建公共资源池，团队协作更高效
- **分享链接** - 一键生成，自定义有效期与访问权限
- **访客统计** - 查看访问记录、来源、地域分析
- **访客模式** - 无需登录即可浏览分享内容
- **权限控制** - 公开/私密/受保护三级权限管理
- **水印功能** - 自定义水印保护版权，防止盗用
- **防盗链** - 完善的防盗链机制，保护资源安全

### 🔐 安全与认证

- **多模式登录** - 账号密码、邮箱、手机多种方式
- **三方登录** - 支持GitHub、Google等第三方平台登录
- **安全设置** - 两步验证、登录日志、设备管理
- **用户认证** - 完善的权限管理系统
- **带宽流量控制** - 灵活配置用户上传下载限制
- **API密钥** - 安全的API访问控制机制

### 🌐 API与集成

- **RESTful API** - 完整的API文档，轻松集成
- **随机API** - 随机获取图片，适配各类场景
- **三方API支持** - 兼容PicGo、Typora等主流工具
- **开放API管理** - 后台可视化配置API密钥
- **Webhook支持** - 事件通知，实时同步

### 📊 管理与监控

- **可视化仪表盘** - 实时数据统计，一目了然
- **存储监控** - 磁盘使用情况实时掌握
- **上传队列** - 批量上传进度管理
- **操作日志** - 完整的操作记录追溯
- **消息中心** - 系统通知、任务提醒集中管理
- **WebSocket实时通信** - 消息即时推送，状态实时更新

### 🛠️ 系统功能

- **公告系统** - Markdown编辑，多版本管理，实时发布
- **后台配置中心** - 数十项可视化配置，无需改代码
- **Docker部署** - 一键启动，简单快捷
- **数据库支持** - PostgreSQL/MySQL/SQLite多选
- **图片优化** - 自动压缩，智能缩略图生成
- **备份恢复** - 完整的数据备份与恢复方案

---

## 🚀 快速开始

### 第一步：创建文件夹

点击左侧菜单 **文件夹管理**，创建您的第一个分类文件夹

### 第二步：上传图片

支持拖拽上传、粘贴上传、批量上传多种方式

### 第三步：体验AI功能

上传后自动识别标签，尝试向量搜索和以图搜图

### 第四步：探索主题

访问 **个人设置 → 外观** 切换您喜欢的主题风格

### 第五步：生成分享链接

选择图片，一键生成分享链接，发送给好友

### 第六步：查看消息中心

点击右上角消息图标，查看系统通知与任务进度

---

## 💡 实用技巧

> **快捷上传：** 在任何页面按 Ctrl + V 粘贴剪贴板图片快速上传

> **批量选择：** 按住 Shift 点击图片可批量选择

> **搜索技巧：** 使用标签组合搜索，如 #风景 #日落 查找特定图片

> **API集成：** 查看文档中的API密钥管理，将PixelPunk集成到您的博客或应用

> **随机图API：** 使用 /api/v1/random 接口，为您的网站添加随机图片

> **防盗链：** 在后台配置 Referer 白名单，保护您的图片资源

> **访客模式：** 分享链接支持无需登录访问，方便快捷

---

## 🎯 下一步计划

我们正在开发更多令人兴奋的功能：

- 🖥️ **桌面端应用** - Windows/macOS/Linux 原生客户端开发中
- 🎥 **视频全面支持** - 视频上传、预览、管理一站式体验
- 🤖 **更多AI模型** - 支持更多AI服务商，识别更精准
- 🔌 **插件系统升级** - 开放插件市场，社区扩展生态
- 📸 **照片墙功能** - 瀑布流展示，打造个性化相册
- 🌍 **更多语言支持** - 日语、韩语、法语等国际化语言
- 🎨 **更多主题加入** - 北欧简约、未来科技、温馨暖色等风格
- 🚀 **CDN加速** - 全球访问更快
- 📱 **移动端APP** - iOS/Android原生应用

---

## 🤝 参与贡献

PixelPunk 是一个开源项目，我们欢迎：

- 🐛 **Bug反馈** - 在GitHub提Issue或加入QQ群反馈
- 💡 **功能建议** - 告诉我们您想要的功能
- 🌐 **翻译贡献** - 帮助我们翻译更多语言
- 💻 **代码贡献** - Fork项目，提交PR
- 🎨 **主题设计** - 设计并分享您的自定义主题
- 🔌 **插件开发** - 为社区开发实用插件
- ⭐ **Star支持** - 给我们一个Star，是最大的鼓励

---

## 📞 获取帮助

遇到问题？我们随时为您服务：

1. 📖 查看 [官方文档](https://pixelpunk.cc) 寻找答案
2. 💬 加入 QQ群 **826708512** 与社区交流
3. 🐛 在 [GitHub Issues](https://github.com/CooperJiang/PixelPunk/issues) 提交问题
4. 📧 通过邮件联系开发团队

---

## 🏆 特色亮点

### 为什么选择 PixelPunk？

✅ **全能型图床** - 从上传到分享，从管理到统计，一站式解决方案

✅ **AI加持** - 智能标签、向量搜索、以图搜图，体验未来图片管理

✅ **开源免费** - 完全开源，永久免费，无任何隐藏费用

✅ **持续更新** - 活跃的开发团队，快速响应，定期更新

✅ **易于部署** - Docker一键部署，5分钟即可上线

✅ **界面精美** - 赛博朋克等多套主题，视觉体验极佳

✅ **功能丰富** - 70+核心功能，满足各种使用场景

✅ **高性能** - Go语言开发，并发处理能力强，响应速度快

✅ **安全可靠** - 多重安全机制，保护您的数据安全

✅ **社区活跃** - QQ群随时交流，问题快速解决

---

<div align="center">

**感谢您选择 PixelPunk！**

开源 • 免费 • 持续更新 🚀

_让图片管理变得简单而美好_

🌟 如果喜欢，请给我们一个 Star 🌟

</div>`
