package dto

type SettingQueryDTO struct {
	Group string `form:"group" json:"group"` // 设置分组，为空则查询所有
	Key   string `form:"key" json:"key"`     // 设置键名，为空则查询所有
}

type SettingCreateDTO struct {
	Key         string      `json:"key" binding:"required,max=100"`  // 设置键名
	Value       interface{} `json:"value" binding:"required"`        // 设置值
	Type        string      `json:"type" binding:"required,max=20"`  // 值类型
	Group       string      `json:"group" binding:"required,max=50"` // 设置分组
	Description string      `json:"description" binding:"max=500"`   // 设置描述
	IsSystem    bool        `json:"is_system"`                       // 是否为系统设置
}

func (d *SettingCreateDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Key.required":    "设置键名不能为空",
		"Key.max":         "设置键名长度不能超过100个字符",
		"Value.required":  "设置值不能为空",
		"Type.required":   "值类型不能为空",
		"Type.max":        "值类型长度不能超过20个字符",
		"Group.required":  "设置分组不能为空",
		"Group.max":       "设置分组长度不能超过50个字符",
		"Description.max": "设置描述长度不能超过500个字符",
	}
}

type SettingUpdateDTO struct {
	Key         string      `json:"key" binding:"required,max=100"`  // 设置键名
	Value       interface{} `json:"value" binding:"required"`        // 设置值
	Type        string      `json:"type" binding:"required,max=20"`  // 值类型
	Group       string      `json:"group" binding:"required,max=50"` // 设置分组
	Description string      `json:"description" binding:"max=500"`   // 设置描述
}

func (d *SettingUpdateDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Key.required":    "设置键名不能为空",
		"Key.max":         "设置键名长度不能超过100个字符",
		"Value.required":  "设置值不能为空",
		"Type.required":   "值类型不能为空",
		"Type.max":        "值类型长度不能超过20个字符",
		"Group.required":  "设置分组不能为空",
		"Group.max":       "设置分组长度不能超过50个字符",
		"Description.max": "设置描述长度不能超过500个字符",
	}
}

type SettingResponseDTO struct {
	ID          uint        `json:"id"`
	Key         string      `json:"key"`
	Value       interface{} `json:"value"`
	Type        string      `json:"type"`
	Group       string      `json:"group"`
	Description string      `json:"description"`
	IsSystem    bool        `json:"is_system"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
}

type SettingListResponseDTO struct {
	Settings []SettingResponseDTO `json:"settings"`
}

type SettingMapResponseDTO struct {
	Group     string                 `json:"group"`      // 设置分组名称
	Settings  map[string]interface{} `json:"settings"`   // 键值对形式的设置
	UpdatedAt string                 `json:"updated_at"` // 最后更新时间
}

type GlobalSettingsGroupDTO struct {
	Group     string                 `json:"group"`      // 设置分组名称
	Settings  map[string]interface{} `json:"settings"`   // 键值对形式的设置
	UpdatedAt string                 `json:"updated_at"` // 最后更新时间
}

type OAuthProvidersDTO struct {
	GithubEnabled  bool `json:"github_enabled"`
	GoogleEnabled  bool `json:"google_enabled"`
	LinuxdoEnabled bool `json:"linuxdo_enabled"`
}

type GlobalSettingsResponseDTO struct {
	Website        map[string]interface{} `json:"website,omitempty"`      // 网站后端功能设置
	WebsiteInfo    map[string]interface{} `json:"website_info,omitempty"` // 网站前端显示设置
	Upload         map[string]interface{} `json:"upload,omitempty"`       // 上传相关设置
	Theme          map[string]interface{} `json:"theme,omitempty"`        // 网站装修设置
	Registration   map[string]interface{} `json:"registration,omitempty"` // 注册相关设置
	Version        map[string]interface{} `json:"version,omitempty"`      // 版本信息设置
	AI             map[string]interface{} `json:"ai,omitempty"`           // AI相关设置
	Vector         map[string]interface{} `json:"vector,omitempty"`       // 向量搜索相关设置
	Guest          map[string]interface{} `json:"guest,omitempty"`        // 游客上传相关设置
	Appearance     map[string]interface{} `json:"appearance,omitempty"`   // 外观设置
	Analytics      map[string]interface{} `json:"analytics,omitempty"`    // 埋点统计设置
	OAuthProviders OAuthProvidersDTO      `json:"oauth_providers"`        // OAuth 提供商状态
	DeployMode     string                 `json:"deploy_mode"`            // 部署模式：standalone, docker, compose
}

type LegalDocumentsResponseDTO struct {
	PrivacyPolicy  string `json:"privacy_policy"`   // 隐私政策内容
	TermsOfService string `json:"terms_of_service"` // 服务条款内容
}

type BatchSettingCreateDTO struct {
	Settings []SettingCreateDTO `json:"settings" binding:"required,min=1,dive"`
}

func (d *BatchSettingCreateDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Settings.required": "设置列表不能为空",
		"Settings.min":      "设置列表至少需要一项",
	}
}

type BatchSettingUpdateDTO struct {
	Settings []SettingUpdateDTO `json:"settings" binding:"required,min=1,dive"`
}

func (d *BatchSettingUpdateDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Settings.required": "设置列表不能为空",
		"Settings.min":      "设置列表至少需要一项",
	}
}

type BatchSettingResponseDTO struct {
	Success []SettingResponseDTO `json:"success"` // 成功的设置
	Failed  []BatchFailedItem    `json:"failed"`  // 失败的设置
}

type BatchFailedItem struct {
	Key     string `json:"key"`     // 设置键名
	Message string `json:"message"` // 错误信息
}

type BatchUpsertSettingDTO struct {
	Settings []SettingCreateDTO `json:"settings" binding:"required,min=1,dive"`
}

func (d *BatchUpsertSettingDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Settings.required": "设置列表不能为空",
		"Settings.min":      "设置列表至少需要一项",
	}
}

type EmailTestDTO struct {
	Host        string `json:"host" binding:"required"`               // SMTP服务器地址
	Port        int    `json:"port" binding:"required,min=1"`         // SMTP服务器端口
	Username    string `json:"username" binding:"required"`           // SMTP用户名
	Password    string `json:"password" binding:"required"`           // SMTP密码
	Encryption  string `json:"encryption"`                            // 加密类型(ssl, tls, starttls)
	FromAddress string `json:"from_address" binding:"required,email"` // 发件人地址
	FromName    string `json:"from_name" binding:"required"`          // 发件人名称
	TestEmail   string `json:"test_email" binding:"required,email"`   // 测试邮件接收地址
}

func (d *EmailTestDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Host.required":        "SMTP服务器地址不能为空",
		"Port.required":        "SMTP服务器端口不能为空",
		"Port.min":             "SMTP服务器端口必须大于0",
		"Username.required":    "SMTP用户名不能为空",
		"Password.required":    "SMTP密码不能为空",
		"FromAddress.required": "发件人地址不能为空",
		"FromAddress.email":    "发件人地址必须是有效的邮箱格式",
		"FromName.required":    "发件人名称不能为空",
		"TestEmail.required":   "测试邮件接收地址不能为空",
		"TestEmail.email":      "测试邮件接收地址必须是有效的邮箱格式",
	}
}

type TestMailDTO struct {
	Email           string `json:"email" binding:"required,email"`             // 测试接收邮箱
	SmtpHost        string `json:"smtp_host" binding:"required"`               // SMTP服务器地址
	SmtpPort        int    `json:"smtp_port" binding:"required,min=1"`         // SMTP服务器端口
	SmtpEncryption  string `json:"smtp_encryption"`                            // 加密类型(ssl, tls, starttls)
	SmtpUsername    string `json:"smtp_username" binding:"required"`           // SMTP用户名
	SmtpPassword    string `json:"smtp_password"`                              // SMTP密码
	SmtpFromAddress string `json:"smtp_from_address" binding:"required,email"` // 发件人地址
	SmtpFromName    string `json:"smtp_from_name" binding:"required"`          // 发件人名称
}

func (d *TestMailDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Email.required":           "测试邮箱不能为空",
		"Email.email":              "测试邮箱格式不正确",
		"SmtpHost.required":        "SMTP服务器地址不能为空",
		"SmtpPort.required":        "SMTP服务器端口不能为空",
		"SmtpPort.min":             "SMTP服务器端口必须大于0",
		"SmtpUsername.required":    "SMTP用户名不能为空",
		"SmtpFromAddress.required": "发件人地址不能为空",
		"SmtpFromAddress.email":    "发件人地址必须是有效的邮箱格式",
		"SmtpFromName.required":    "发件人名称不能为空",
	}
}

type VectorTestDTO struct {
	Provider string `json:"provider" binding:"required"`      // 向量化提供者
	Model    string `json:"model" binding:"required"`         // 向量化模型
	APIKey   string `json:"api_key" binding:"required"`       // API密钥
	BaseURL  string `json:"base_url"`                         // API代理地址
	Timeout  int    `json:"timeout" binding:"required,min=1"` // 超时时间(秒)
}

func (d *VectorTestDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Provider.required": "向量化提供者不能为空",
		"Model.required":    "向量化模型不能为空",
		"APIKey.required":   "API密钥不能为空",
		"Timeout.required":  "超时时间不能为空",
		"Timeout.min":       "超时时间必须大于0",
	}
}

type QdrantTestDTO struct {
	QdrantURL     string `json:"qdrant_url" binding:"required"`           // Qdrant服务器地址
	QdrantTimeout int    `json:"qdrant_timeout" binding:"required,min=1"` // 连接超时时间(秒)
}

func (d *QdrantTestDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"QdrantURL.required":     "Qdrant服务器地址不能为空",
		"QdrantTimeout.required": "连接超时时间不能为空",
		"QdrantTimeout.min":      "连接超时时间必须大于0",
	}
}

type GithubOAuthConfig struct {
	Enabled       bool   `json:"enabled"`        // 是否启用
	ClientID      string `json:"client_id"`      // Client ID
	ClientSecret  string `json:"client_secret"`  // Client Secret
	RedirectURI   string `json:"redirect_uri"`   // 回调地址
	Scope         string `json:"scope"`          // 授权范围
	ProxyEnabled  bool   `json:"proxy_enabled"`  // 是否启用代理
	ProxyDynamic  bool   `json:"proxy_dynamic"`  // 是否使用动态代理
	ProxyAPIURL   string `json:"proxy_api_url"`  // 动态代理API地址
	ProxyType     string `json:"proxy_type"`     // 代理类型
	ProxyHost     string `json:"proxy_host"`     // 代理地址（静态代理）
	ProxyPort     string `json:"proxy_port"`     // 代理端口（静态代理）
	ProxyUsername string `json:"proxy_username"` // 代理用户名
	ProxyPassword string `json:"proxy_password"` // 代理密码
}

type GoogleOAuthConfig struct {
	Enabled       bool   `json:"enabled"`
	ClientID      string `json:"client_id"`
	ClientSecret  string `json:"client_secret"`
	RedirectURI   string `json:"redirect_uri"`
	Scope         string `json:"scope"`
	ProxyEnabled  bool   `json:"proxy_enabled"`
	ProxyDynamic  bool   `json:"proxy_dynamic"`
	ProxyAPIURL   string `json:"proxy_api_url"`
	ProxyType     string `json:"proxy_type"`
	ProxyHost     string `json:"proxy_host"`
	ProxyPort     string `json:"proxy_port"`
	ProxyUsername string `json:"proxy_username"`
	ProxyPassword string `json:"proxy_password"`
}

type LinuxdoOAuthConfig struct {
	Enabled       bool   `json:"enabled"`
	ClientID      string `json:"client_id"`
	ClientSecret  string `json:"client_secret"`
	RedirectURI   string `json:"redirect_uri"`
	Scope         string `json:"scope"`
	ProxyEnabled  bool   `json:"proxy_enabled"`
	ProxyDynamic  bool   `json:"proxy_dynamic"`
	ProxyAPIURL   string `json:"proxy_api_url"`
	ProxyType     string `json:"proxy_type"`
	ProxyHost     string `json:"proxy_host"`
	ProxyPort     string `json:"proxy_port"`
	ProxyUsername string `json:"proxy_username"`
	ProxyPassword string `json:"proxy_password"`
}

type OAuthConfigResponseDTO struct {
	Github  GithubOAuthConfig  `json:"github"`
	Google  GoogleOAuthConfig  `json:"google"`
	Linuxdo LinuxdoOAuthConfig `json:"linuxdo"`
}

type GithubOAuthLoginDTO struct {
	Code string `json:"code" binding:"required"` // GitHub 返回的授权码
}

func (d *GithubOAuthLoginDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Code.required": "授权码不能为空",
	}
}

type GoogleOAuthLoginDTO struct {
	Code string `json:"code" binding:"required"`
}

func (d *GoogleOAuthLoginDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Code.required": "授权码不能为空",
	}
}

type LinuxdoOAuthLoginDTO struct {
	Code string `json:"code" binding:"required"`
}

func (d *LinuxdoOAuthLoginDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Code.required": "授权码不能为空",
	}
}

type TestProxyDTO struct {
	ProxyDynamic  bool   `json:"proxy_dynamic"`  // 是否使用动态代理
	ProxyAPIURL   string `json:"proxy_api_url"`  // 动态代理API地址
	ProxyType     string `json:"proxy_type"`     // 代理类型: http, https, socks5
	ProxyHost     string `json:"proxy_host"`     // 代理地址
	ProxyPort     string `json:"proxy_port"`     // 代理端口
	ProxyUsername string `json:"proxy_username"` // 代理用户名(可选)
	ProxyPassword string `json:"proxy_password"` // 代理密码(可选)
}

func (d *TestProxyDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"ProxyType": "代理类型格式错误",
		"ProxyHost": "代理地址格式错误",
		"ProxyPort": "代理端口格式错误",
	}
}

type TestProxyResponseDTO struct {
	Success bool   `json:"success"` // 是否成功
	Message string `json:"message"` // 消息
	Latency int64  `json:"latency"` // 延迟(毫秒)
}
