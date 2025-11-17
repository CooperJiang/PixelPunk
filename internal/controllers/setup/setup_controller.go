package setup

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"pixelpunk/internal/controllers/setting/dto"
	"pixelpunk/internal/controllers/websocket"
	"pixelpunk/internal/cron"
	"pixelpunk/internal/models"
	ai "pixelpunk/internal/services/ai"
	"pixelpunk/internal/services/auth"
	"pixelpunk/internal/services/message"
	"pixelpunk/internal/services/setting"
	"pixelpunk/internal/services/storage"
	"pixelpunk/internal/services/user"
	vectorSvc "pixelpunk/internal/services/vector"
	"pixelpunk/migrations"
	"pixelpunk/pkg/cache"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/config"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/email"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/utils"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/gin-gonic/gin"
)

type SetupController struct{}

func (s *SetupController) GetStatus(c *gin.Context) {
	installManager := common.GetInstallManager()
	status := installManager.GetStatus()
	errors.ResponseSuccess(c, status, "获取安装状态成功")
}

func (s *SetupController) TestConnection(c *gin.Context) {
	var req struct {
		Type     string `json:"type" binding:"required"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Path     string `json:"path"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.New(errors.CodeValidationFailed, "参数错误: "+err.Error()))
		return
	}

	switch req.Type {
	case "mysql":
		if req.Host == "" || req.Username == "" || req.Name == "" {
			errors.HandleError(c, errors.New(errors.CodeValidationFailed, "MySQL数据库连接信息不完整"))
			return
		}
	case "sqlite":
		if req.Path == "" {
			errors.HandleError(c, errors.New(errors.CodeValidationFailed, "SQLite数据库文件路径不能为空"))
			return
		}
	default:
		errors.HandleError(c, errors.New(errors.CodeValidationFailed, "不支持的数据库类型"))
		return
	}

	if err := database.TestDatabaseConnection(req.Type, req.Host, req.Port, req.Username, req.Password, req.Name, req.Path); err != nil {
		errors.HandleError(c, errors.New(errors.CodeValidationFailed, "数据库连接失败: "+sanitizeDBError(err)))
		return
	}

	errors.ResponseSuccess(c, nil, "数据库连接测试成功")
}

func (s *SetupController) Install(c *gin.Context) {
	var req struct {
		Database struct {
			Type     string `json:"type" binding:"required"`
			Host     string `json:"host"`
			Port     int    `json:"port"`
			Username string `json:"username"`
			Password string `json:"password"`
			Name     string `json:"name"`
			Path     string `json:"path"`
		} `json:"database" binding:"required"`
		Redis struct {
			Host     string `json:"host"`
			Port     int    `json:"port"`
			Password string `json:"password"`
			DB       int    `json:"db"`
		} `json:"redis"`
		Vector struct {
			QdrantURL     string `json:"qdrant_url"`
			QdrantTimeout int    `json:"qdrant_timeout"`
			UseBuiltin    bool   `json:"use_builtin"`
			HTTPPort      int    `json:"http_port"`
			GRPCPort      int    `json:"grpc_port"`
		} `json:"vector"`
		AdminUsername string `json:"admin_username" binding:"required"`
		AdminPassword string `json:"admin_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.New(errors.CodeValidationFailed, "参数错误: "+err.Error()))
		return
	}

	installManager := common.GetInstallManager()
	if err := installManager.StartInstall(); err != nil {
		errors.HandleError(c, errors.New(errors.CodeValidationFailed, err.Error()))
		return
	}

	var installSuccess bool
	defer func() {
		installManager.FinishInstall(installSuccess)
	}()

	// 检查是否为Docker/Compose模式且配置已预设
	deployMode := common.GetDeployMode()
	configPreset := common.IsConfigPreset()

	// 在Docker/Compose模式下且配置已预设时，跳过数据库测试和配置文件写入
	if !configPreset || deployMode == "standalone" {
		if err := database.TestDatabaseConnection(req.Database.Type, req.Database.Host, req.Database.Port, req.Database.Username, req.Database.Password, req.Database.Name, req.Database.Path); err != nil {
			errors.HandleError(c, errors.New(errors.CodeValidationFailed, "数据库连接失败: "+sanitizeDBError(err)))
			return
		}

		if err := writeConfigFile(req); err != nil {
			errors.HandleError(c, errors.New(errors.CodeInternal, "写入配置文件失败: "+err.Error()))
			return
		}
	}

	config.ReloadConfig()

	if err := database.ReconnectDatabase(); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, "重新连接数据库失败: "+sanitizeDBError(err)))
		return
	}

	adminUser, userMessage, err := createAdminUser(req.AdminUsername, req.AdminPassword)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, "创建管理员用户失败: "+err.Error()))
		return
	}

	if err := initializeSystemServices(req.Vector); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, "系统初始化失败: "+err.Error()))
		return
	}

	installManager.SetInstallMode(false)
	installManager.SetSystemInstalled(true)

	jwtSecret := "pixelpunk_default_secret_key"
	expiresHours := 168

	if securitySettings, err := setting.GetSettingsByGroupAsMap("security"); err == nil {
		if val, ok := securitySettings.Settings["jwt_secret"]; ok {
			if secretStr, ok := val.(string); ok && secretStr != "" {
				jwtSecret = secretStr
			}
		}
		if val, ok := securitySettings.Settings["login_expire_hours"]; ok {
			if hours, ok := val.(float64); ok && hours > 0 {
				expiresHours = int(hours)
			}
		}
	}

	token, err := auth.GenerateToken(adminUser.ID, adminUser.Username, int(adminUser.Role), jwtSecret, expiresHours)
	if err != nil {
		installSuccess = true
		errors.ResponseSuccess(c, gin.H{"message": userMessage}, userMessage)
		return
	}

	installSuccess = true
	errors.ResponseSuccess(c, gin.H{
		"message": userMessage,
		"token":   token,
		"user": gin.H{
			"id":       adminUser.ID,
			"username": adminUser.Username,
			"role":     adminUser.Role,
		},
	}, userMessage)
}

func writeConfigFile(req struct {
	Database struct {
		Type     string `json:"type" binding:"required"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Path     string `json:"path"`
	} `json:"database" binding:"required"`
	Redis struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	} `json:"redis"`
	Vector struct {
		QdrantURL     string `json:"qdrant_url"`
		QdrantTimeout int    `json:"qdrant_timeout"`
		UseBuiltin    bool   `json:"use_builtin"`
		HTTPPort      int    `json:"http_port"`
		GRPCPort      int    `json:"grpc_port"`
	} `json:"vector"`
	AdminUsername string `json:"admin_username" binding:"required"`
	AdminPassword string `json:"admin_password" binding:"required"`
}) error {
	redisHost := req.Redis.Host
	redisPort := req.Redis.Port
	if redisHost == "" {
		redisHost = "localhost"
	}
	if redisPort == 0 {
		redisPort = 6379
	}

	configPaths := []string{
		"configs/config.yaml",
		"config.yaml",
	}

	var configPath string
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			configPath = path
			break
		}
	}

	if configPath != "" {
		return updateExistingConfigFile(configPath, req, redisHost, redisPort)
	}

	configPath = "configs/config.yaml"
	if err := os.MkdirAll("configs", 0755); err != nil {
		configPath = "config.yaml"
	}

	return createNewConfigFile(configPath, req, redisHost, redisPort)
}

func updateExistingConfigFile(configPath string, req struct {
	Database struct {
		Type     string `json:"type" binding:"required"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Path     string `json:"path"`
	} `json:"database" binding:"required"`
	Redis struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	} `json:"redis"`
	Vector struct {
		QdrantURL     string `json:"qdrant_url"`
		QdrantTimeout int    `json:"qdrant_timeout"`
		UseBuiltin    bool   `json:"use_builtin"`
		HTTPPort      int    `json:"http_port"`
		GRPCPort      int    `json:"grpc_port"`
	} `json:"vector"`
	AdminUsername string `json:"admin_username" binding:"required"`
	AdminPassword string `json:"admin_password" binding:"required"`
}, redisHost string, redisPort int) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	var existingConfig map[string]interface{}
	if err := yaml.Unmarshal(data, &existingConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	if existingConfig["database"] == nil {
		existingConfig["database"] = make(map[string]interface{})
	}
	dbConfig := existingConfig["database"].(map[string]interface{})

	dbConfig["type"] = req.Database.Type
	if req.Database.Type == "mysql" {
		dbConfig["host"] = req.Database.Host
		dbConfig["port"] = req.Database.Port
		dbConfig["username"] = req.Database.Username
		dbConfig["password"] = req.Database.Password
		dbConfig["name"] = req.Database.Name
		delete(dbConfig, "path")
	} else if req.Database.Type == "sqlite" {
		dbConfig["path"] = req.Database.Path
		delete(dbConfig, "host")
		delete(dbConfig, "port")
		delete(dbConfig, "username")
		delete(dbConfig, "password")
		delete(dbConfig, "name")
	}

	if existingConfig["redis"] == nil {
		existingConfig["redis"] = make(map[string]interface{})
	}
	redisConfig := existingConfig["redis"].(map[string]interface{})
	redisConfig["host"] = redisHost
	redisConfig["port"] = redisPort
	redisConfig["password"] = req.Redis.Password
	redisConfig["db"] = req.Redis.DB

	updatedData, err := yaml.Marshal(existingConfig)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	return os.WriteFile(configPath, updatedData, 0600)
}

func createNewConfigFile(configPath string, req struct {
	Database struct {
		Type     string `json:"type" binding:"required"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Path     string `json:"path"`
	} `json:"database" binding:"required"`
	Redis struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	} `json:"redis"`
	Vector struct {
		QdrantURL     string `json:"qdrant_url"`
		QdrantTimeout int    `json:"qdrant_timeout"`
		UseBuiltin    bool   `json:"use_builtin"`
		HTTPPort      int    `json:"http_port"`
		GRPCPort      int    `json:"grpc_port"`
	} `json:"vector"`
	AdminUsername string `json:"admin_username" binding:"required"`
	AdminPassword string `json:"admin_password" binding:"required"`
}, redisHost string, redisPort int) error {
	var configContent string

	switch req.Database.Type {
	case "mysql":
		configContent = fmt.Sprintf(`# 应用基本配置
app:
  port: 9520
  mode: "debug"

# 数据库配置
database:
  type: "mysql"
  host: "%s"
  port: %d
  username: "%s"
  password: "%s"
  name: "%s"

# Redis配置
redis:
  host: "%s"
  port: %d
  password: "%s"
  db: %d

# 跨域(CORS)配置
cors:
  enabled: true
  allowed_origins:
    - "*"
`,
			req.Database.Host,
			req.Database.Port,
			req.Database.Username,
			req.Database.Password,
			req.Database.Name,
			redisHost,
			redisPort,
			req.Redis.Password,
			req.Redis.DB,
		)
	case "sqlite":
		configContent = fmt.Sprintf(`# 应用基本配置
app:
  port: 9520
  mode: "debug"

# 数据库配置
database:
  type: "sqlite"
  path: "%s"

# Redis配置
redis:
  host: "%s"
  port: %d
  password: "%s"
  db: %d

# 跨域(CORS)配置
cors:
  enabled: true
  allowed_origins:
    - "*"
`,
			req.Database.Path,
			redisHost,
			redisPort,
			req.Redis.Password,
			req.Redis.DB,
		)
	default:
		return fmt.Errorf("不支持的数据库类型: %s", req.Database.Type)
	}

	return os.WriteFile(configPath, []byte(configContent), 0600)
}

func createAdminUser(username, password string) (*models.User, string, error) {
	db := database.GetDB()
	if db == nil {
		return nil, "", fmt.Errorf("数据库连接未初始化")
	}

	var existingUser models.User
	err := db.Model(&models.User{}).Where("username = ?", username).First(&existingUser).Error
	if err == nil {
		return &existingUser, fmt.Sprintf("数据库迁移完成，检测到已存在管理员账户 %s，已自动登录", username), nil
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", fmt.Errorf("密码加密失败: %v", err)
	}

	adminUser := models.User{
		Username: username,
		Password: hashedPassword,
		Email:    fmt.Sprintf("%s@pixelpunk.local", username),
		Role:     common.UserRoleSuperAdmin,
		Status:   common.UserStatusNormal,
	}

	if err := db.Create(&adminUser).Error; err != nil {
		return nil, "", fmt.Errorf("创建管理员用户失败: %v", err)
	}

	return &adminUser, fmt.Sprintf("系统安装完成，已创建管理员账户 %s，已自动登录", username), nil
}

func sanitizeDBError(err error) string {
	if err == nil {
		return ""
	}

	errStr := err.Error()

	if strings.Contains(errStr, "@") && strings.Contains(errStr, "using password") {
		errStr = strings.ReplaceAll(errStr, "using password: YES", "using password: ***")
		errStr = strings.ReplaceAll(errStr, "using password: NO", "using password: ***")

		parts := strings.Split(errStr, "'")
		if len(parts) >= 3 {
			for i := 1; i < len(parts); i += 2 {
				if strings.Contains(parts[i], "@") {
					parts[i] = "***@***"
				}
			}
			errStr = strings.Join(parts, "'")
		}
	}

	return errStr
}

func initializeSystemServices(vectorConfig struct {
	QdrantURL     string `json:"qdrant_url"`
	QdrantTimeout int    `json:"qdrant_timeout"`
	UseBuiltin    bool   `json:"use_builtin"`
	HTTPPort      int    `json:"http_port"`
	GRPCPort      int    `json:"grpc_port"`
}) error {
	cache.InitCache()
	runMigrations()

	var qdrantURL string
	var qdrantTimeout int

	configPreset := common.IsConfigPreset()

	if configPreset {
		cfg := config.GetConfig()
		qdrantURL = cfg.Vector.QdrantURL
		qdrantTimeout = cfg.Vector.Timeout
		if qdrantURL == "" {
			qdrantURL = "http://localhost:6333"
		}
		if qdrantTimeout <= 0 {
			qdrantTimeout = 30
		}
	} else {
		if vectorConfig.UseBuiltin {
			actualHTTPPort := vectorConfig.HTTPPort
			if actualHTTPPort <= 0 {
				actualHTTPPort = 6333
			}
			qdrantURL = fmt.Sprintf("http://localhost:%d", actualHTTPPort)
		} else {
			qdrantURL = vectorConfig.QdrantURL
		}
		qdrantTimeout = vectorConfig.QdrantTimeout
	}

	writeVectorConfigToDatabase(qdrantURL, qdrantTimeout)

	if vectorConfig.UseBuiltin {
		go func() {
			startBuiltinQdrant(vectorConfig.HTTPPort, vectorConfig.GRPCPort)
		}()
		time.Sleep(2 * time.Second)
	}

	storage.CheckAndInitDefaultChannel()
	email.Init()
	initAllServices()
	websocket.InitWebSocketManager()
	cron.InitCronManager()

	return nil
}

func runMigrations() {
	db := database.GetDB()
	if db == nil {
		return
	}
	migrations.RegisterAllMigrations(db)
}

func initAllServices() {
	user.InitUserService()
	setting.InitSettingService()

	message.InitMessageService()
	templateService := message.GetTemplateService()
	templateService.InitDefaultTemplates()

	ai.RegisterAISettingHooks()
	vectorSvc.RegisterVectorConfigHooks()
}

func writeVectorConfigToDatabase(qdrantURL string, qdrantTimeout int) error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("数据库连接未初始化")
	}

	if qdrantURL == "" {
		qdrantURL = "http://localhost:6333"
	}
	if qdrantTimeout <= 0 {
		qdrantTimeout = 30
	}

	upsertDTO := &dto.BatchUpsertSettingDTO{
		Settings: []dto.SettingCreateDTO{
			{
				Key:         "vector_enabled",
				Value:       true,
				Type:        "boolean",
				Group:       "vector",
				Description: "向量功能开关",
				IsSystem:    true,
			},
			{
				Key:         "qdrant_url",
				Value:       qdrantURL,
				Type:        "string",
				Group:       "vector",
				Description: "Qdrant向量数据库地址",
				IsSystem:    true,
			},
			{
				Key:         "qdrant_timeout",
				Value:       qdrantTimeout,
				Type:        "number",
				Group:       "vector",
				Description: "Qdrant连接超时时间(秒)",
				IsSystem:    true,
			},
		},
	}

	result, err := setting.BatchUpsertSettings(upsertDTO)
	if err != nil {
		return fmt.Errorf("写入向量配置失败: %v", err)
	}

	if len(result.Failed) > 0 {
		return fmt.Errorf("部分向量配置写入失败: %v", result.Failed)
	}

	return nil
}

func isPortInUse(port int) bool {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func updateQdrantConfig(configPath string, httpPort, grpcPort int) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取Qdrant配置文件失败: %v", err)
	}

	content := string(data)

	httpPortPattern := regexp.MustCompile(`(?m)^(\s*http_port:\s*)(\d+)`)
	if httpPortPattern.MatchString(content) {
		content = httpPortPattern.ReplaceAllString(content, fmt.Sprintf("${1}%d", httpPort))
	} else {
		servicePattern := regexp.MustCompile(`(?m)^(service:)`)
		if servicePattern.MatchString(content) {
			content = servicePattern.ReplaceAllString(content, fmt.Sprintf("${1}\n  http_port: %d", httpPort))
		}
	}

	grpcPortPattern := regexp.MustCompile(`(?m)^(\s*grpc_port:\s*)(\d+)`)
	if grpcPortPattern.MatchString(content) {
		content = grpcPortPattern.ReplaceAllString(content, fmt.Sprintf("${1}%d", grpcPort))
	} else {
		servicePattern := regexp.MustCompile(`(?m)^(service:)`)
		if servicePattern.MatchString(content) {
			content = servicePattern.ReplaceAllString(content, fmt.Sprintf("${1}\n  grpc_port: %d", grpcPort))
		}
	}

	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入Qdrant配置文件失败: %v", err)
	}

	return nil
}

func startBuiltinQdrant(httpPort, grpcPort int) error {
	if httpPort <= 0 {
		httpPort = 6333
	}
	if grpcPort <= 0 {
		grpcPort = 6334
	}

	if isPortInUse(httpPort) {
		return fmt.Errorf("HTTP端口 %d 已被占用", httpPort)
	}

	qdrantBinPath := findQdrantBinary()
	if qdrantBinPath == "" {
		return fmt.Errorf("未找到Qdrant二进制文件")
	}

	if _, err := os.Stat(qdrantBinPath); err != nil {
		return fmt.Errorf("无法访问二进制文件: %v", err)
	}

	configPath, err := prepareQdrantConfig(httpPort, grpcPort)
	if err != nil {
		return fmt.Errorf("准备配置文件失败: %v", err)
	}

	return launchQdrantProcess(qdrantBinPath, configPath)
}

func findQdrantBinary() string {
	paths := []string{
		"qdrant/bin/qdrant",
		"qdrant/qdrant",
		"./qdrant/bin/qdrant",
		"qdrant/bin/qdrant.exe",
	}
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

func prepareQdrantConfig(httpPort, grpcPort int) (string, error) {
	configPath := "qdrant/config/config.yaml"

	if _, err := os.Stat(configPath); err == nil {
		return configPath, updateQdrantConfig(configPath, httpPort, grpcPort)
	}

	if err := os.MkdirAll("qdrant/config", 0755); err != nil {
		return "", err
	}

	config := fmt.Sprintf(`service:
  http_port: %d
  grpc_port: %d
storage:
  storage_path: ./storage
`, httpPort, grpcPort)

	return configPath, os.WriteFile(configPath, []byte(config), 0644)
}

func launchQdrantProcess(binPath, configPath string) error {
	absBinPath, err := filepath.Abs(binPath)
	if err != nil {
		return fmt.Errorf("无法获取二进制文件绝对路径: %v", err)
	}

	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		return fmt.Errorf("无法获取配置文件绝对路径: %v", err)
	}

	if err := os.Chmod(absBinPath, 0755); err != nil {
		return fmt.Errorf("无法设置执行权限: %v", err)
	}

	testCmd := exec.Command(absBinPath, "--version")
	testOutput, testErr := testCmd.CombinedOutput()
	if testErr != nil {
		if strings.Contains(string(testOutput), "GLIBC") || strings.Contains(testErr.Error(), "GLIBC") {
			return fmt.Errorf("系统GLIBC版本不兼容，需要GLIBC 2.31+")
		}
		return fmt.Errorf("Qdrant二进制文件无法运行: %v", testErr)
	}

	workDir, _ := os.Getwd()
	cmd := exec.Command(absBinPath, "--config-path", absConfigPath)
	cmd.Dir = workDir

	if err := os.MkdirAll("logs", 0755); err != nil {
		return fmt.Errorf("无法创建日志目录: %v", err)
	}

	logFile, err := os.OpenFile("logs/qdrant.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("无法打开日志文件: %v", err)
	}
	defer func() {
		if logFile != nil {
			logFile.Close()
		}
	}()

	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动进程失败: %v", err)
	}

	pid := cmd.Process.Pid
	os.MkdirAll("qdrant", 0755)
	pidFile := "qdrant/qdrant.pid"
	os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", pid)), 0644)

	time.Sleep(1 * time.Second)

	if err := cmd.Process.Signal(os.Signal(nil)); err != nil {
		return fmt.Errorf("进程启动后立即退出")
	}

	return nil
}
