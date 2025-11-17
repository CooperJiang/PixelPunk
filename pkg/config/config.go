package config

import (
	"os"
	"pixelpunk/pkg/logger"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

const (
	// 环境变量前缀
	envPrefix = "APP_"
)

// Config 应用配置结构
type Config struct {
	App      AppConfig      `yaml:"app" env:"APP"`
	Database DatabaseConfig `yaml:"database" env:"DB"`
	Redis    RedisConfig    `yaml:"redis" env:"REDIS"`
	Upload   UploadConfig   `yaml:"upload" env:"UPLOAD"`
	Vector   VectorConfig   `yaml:"vector" env:"VECTOR"`
}

// 更新服务配置已移除

// AppConfig 应用基础配置
type AppConfig struct {
	Port      int    `yaml:"port" env:"PORT"`
	Mode      string `yaml:"mode" env:"MODE"`
	Namespace string `yaml:"ns" env:"NS"` // 命名空间，用于缓存隔离，默认: pixelpunk
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type     string `yaml:"type" env:"TYPE"` // 数据库类型: mysql/sqlite
	Host     string `yaml:"host" env:"HOST"`
	Port     int    `yaml:"port" env:"PORT"`
	Username string `yaml:"username" env:"USERNAME"`
	Password string `yaml:"password" env:"PASSWORD"`
	Name     string `yaml:"name" env:"NAME"`
	Path     string `yaml:"path" env:"PATH"` // SQLite数据库文件路径
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `yaml:"host" env:"HOST"`
	Port     int    `yaml:"port" env:"PORT"`
	Password string `yaml:"password" env:"PASSWORD"`
	DB       int    `yaml:"db" env:"DB"`
}

// UploadConfig 上传配置
type UploadConfig struct {
	MaxFileSize  int64    `yaml:"max_file_size" env:"MAX_FILE_SIZE"` // 最大文件大小（字节）
	AllowedTypes []string `yaml:"allowed_types" env:"ALLOWED_TYPES"` // 允许的文件类型
}

// VectorConfig 向量数据库配置
type VectorConfig struct {
	Enabled       bool   `yaml:"enabled" env:"ENABLED"`                 // 是否启用向量搜索功能
	QdrantURL     string `yaml:"qdrant_url" env:"QDRANT_URL"`           // Qdrant数据库地址
	Timeout       int    `yaml:"timeout" env:"TIMEOUT"`                 // 请求超时时间（秒）
	OpenAIAPIKey  string `yaml:"openai_api_key" env:"OPENAI_API_KEY"`   // OpenAI API密钥
	OpenAIBaseURL string `yaml:"openai_base_url" env:"OPENAI_BASE_URL"` // OpenAI API地址
	OpenAIModel   string `yaml:"openai_model" env:"OPENAI_MODEL"`       // 向量化模型
}

var (
	config Config
	once   sync.Once
)

// setDefaultConfig 设置默认配置值
func setDefaultConfig(cfg *Config) {
	cfg.App.Port = 9520
	cfg.App.Mode = "debug"
	cfg.App.Namespace = "pixelpunk"

	// 数据库默认配置（用于安装模式）
	cfg.Database.Port = 3306

	cfg.Redis.Host = "localhost"
	cfg.Redis.Port = 6379
	cfg.Redis.DB = 0
}

// InitConfig 初始化配置
func InitConfig() {
	once.Do(func() {
		setDefaultConfig(&config)

		// 先从配置文件读取默认配置
		loadConfigFromFile(&config)

		// 然后从环境变量中覆盖配置
		loadConfigFromEnv(&config)

		// 打印加载配置信息（带颜色）
	})
}

// ReloadConfig 强制重新加载配置（用于安装完成后）
func ReloadConfig() {
	setDefaultConfig(&config)

	// 先从配置文件读取默认配置
	loadConfigFromFile(&config)

	// 然后从环境变量中覆盖配置
	loadConfigFromEnv(&config)

	// 打印加载配置信息（带颜色）
}

// loadConfigFromFile 从配置文件加载配置
func loadConfigFromFile(cfg *Config) {
	// 按优先级尝试读取配置文件
	configPaths := []string{
		"configs/config.yaml", // 新路径
		"config.yaml",         // 旧路径（向后兼容）
	}

	var data []byte
	var err error

	for _, path := range configPaths {
		data, err = os.ReadFile(path)
		if err == nil {
			break
		}
	}

	if err != nil {
		logger.Warn("无法读取配置文件: %v，将使用默认配置和环境变量", err)
		return
	}


	if err := yaml.Unmarshal(data, cfg); err != nil {
		logger.Warn("无法解析配置文件: %v，将使用默认配置和环境变量", err)
	}
}

// loadConfigFromEnv 从环境变量加载配置
func loadConfigFromEnv(cfg *Config) {
	// 处理App配置的环境变量
	loadEnvToStruct(envPrefix+"APP_", &cfg.App)

	// 处理Database配置的环境变量
	loadEnvToStruct(envPrefix+"DB_", &cfg.Database)

	// 处理Redis配置的环境变量
	loadEnvToStruct(envPrefix+"REDIS_", &cfg.Redis)

	// 处理Upload配置的环境变量
	loadEnvToStruct(envPrefix+"UPLOAD_", &cfg.Upload)

	// 处理Vector配置的环境变量
	loadEnvToStruct(envPrefix+"VECTOR_", &cfg.Vector)

}

// loadEnvToStruct 加载环境变量到结构体
func loadEnvToStruct(prefix string, obj interface{}) {
	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, prefix) {
			continue
		}

		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimPrefix(parts[0], prefix)
		value := parts[1]

		// 将环境变量的值设置到结构体中
		setStructField(obj, key, value)
	}
}

// setStructField 设置结构体字段的值
func setStructField(obj interface{}, key string, value string) {
	v := reflect.ValueOf(obj).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// 根据env标签找到对应的环境变量字段
		envKey := field.Tag.Get("env")
		if envKey == key {
			// 找到匹配的字段，根据字段类型设置值
			setFieldValue(v.Field(i), value)
			return
		}
	}
}

// setFieldValue 根据字段类型设置值
func setFieldValue(field reflect.Value, value string) {
	if !field.CanSet() {
		return
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			field.SetInt(intValue)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if uintValue, err := strconv.ParseUint(value, 10, 64); err == nil {
			field.SetUint(uintValue)
		}
	case reflect.Float32, reflect.Float64:
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			field.SetFloat(floatValue)
		}
	case reflect.Bool:
		if boolValue, err := strconv.ParseBool(value); err == nil {
			field.SetBool(boolValue)
		}
	}
}

func GetConfig() *Config {
	if reflect.DeepEqual(config, Config{}) {
		InitConfig()
	}
	return &config
}

// GetEnvString 从环境变量获取字符串值
func GetEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// GetEnvInt 从环境变量获取整数值
func GetEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetEnvBool 从环境变量获取布尔值
func GetEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func GetUploadConfig() *UploadConfig {
	config := GetConfig()

	if config.Upload.MaxFileSize == 0 {
		config.Upload.MaxFileSize = 100 * 1024 * 1024 // 100MB
	}

	if len(config.Upload.AllowedTypes) == 0 {
		config.Upload.AllowedTypes = []string{
			"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp", "image/bmp",
		}
	}

	return &config.Upload
}
