package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode"

	"gorm.io/gorm/logger"
)

// Colors
const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

// Logger 自定义日志结构体
type Logger struct {
	*log.Logger
	config    *Config
	LogLevel  logger.LogLevel
	SlowQuery time.Duration
}

// Config 日志配置
type Config struct {
	SlowThreshold time.Duration // 慢查询阈值
	Colorful      bool          // 是否启用彩色输出
	LogLevel      logger.LogLevel
}

// 全局日志实例
var Log *Logger

// 全局方法
var (
	Infof  func(format string, args ...interface{})
	Warnf  func(format string, args ...interface{})
	Errorf func(format string, args ...interface{})
	Debugf func(format string, args ...interface{})
)

// 默认配置
var defaultConfig = &Config{
	SlowThreshold: 200 * time.Millisecond, // 慢查询阈值，超过这个时间的查询会被标记为慢查询
	Colorful:      true,                   // 默认开启彩色输出
	LogLevel:      logger.Info,            // 默认日志级别为 Info
}

// New 创建一个新的日志实例
func New(config *Config) *Logger {
	if config == nil {
		config = &Config{
			SlowThreshold: 200 * time.Millisecond,
			Colorful:      true,
			LogLevel:      logger.Info,
		}
	}

	l := &Logger{
		Logger:    log.New(os.Stdout, "", log.LstdFlags),
		config:    config,
		LogLevel:  config.LogLevel,
		SlowQuery: config.SlowThreshold,
	}

	return l
}

// LogMode 设置日志级别
func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info 打印信息日志
func (l *Logger) Info(ctx context.Context, format string, args ...interface{}) {
	if l.LogLevel >= logger.Info {
		safeFormat := sanitizeLogContent(format)
		safeArgs := sanitizeArgs(args)
		if l.config.Colorful {
			l.Printf(Green+"[INFO] "+safeFormat+Reset, safeArgs...)
		} else {
			l.Printf("[INFO] "+safeFormat, safeArgs...)
		}
	}
}

// Warn 打印警告日志
func (l *Logger) Warn(ctx context.Context, format string, args ...interface{}) {
	if l.LogLevel >= logger.Warn {
		safeFormat := sanitizeLogContent(format)
		safeArgs := sanitizeArgs(args)
		if l.config.Colorful {
			l.Printf(Yellow+"[WARN] "+safeFormat+Reset, safeArgs...)
		} else {
			l.Printf("[WARN] "+safeFormat, safeArgs...)
		}
	}
}

// Error 打印错误日志
func (l *Logger) Error(ctx context.Context, format string, args ...interface{}) {
	if l.LogLevel >= logger.Error {
		safeFormat := sanitizeLogContent(format)
		safeArgs := sanitizeArgs(args)
		if l.config.Colorful {
			l.Printf(Red+"[ERROR] "+safeFormat+Reset, safeArgs...)
		} else {
			l.Printf("[ERROR] "+safeFormat, safeArgs...)
		}
	}
}

// Trace SQL 追踪
func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		l.Error(ctx, "[%.3fms] [rows:%v] %s; %s", float64(elapsed.Nanoseconds())/1e6, rows, sql, err)
		return
	}

	if l.SlowQuery != 0 && elapsed > l.SlowQuery {
		l.Warn(ctx, "[%.3fms] [rows:%v] %s; %s", float64(elapsed.Nanoseconds())/1e6, rows, sql, "SLOW SQL")
		return
	}

	l.Info(ctx, "[%.3fms] [rows:%v] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
}

// InitLogger 初始化日志
func InitLogger(config *Config) {
	Log = New(config)

	Infof = func(format string, args ...interface{}) {
		safeFormat := sanitizeLogContent(format)
		safeArgs := sanitizeArgs(args)
		Log.Info(context.Background(), safeFormat, safeArgs...)
	}

	Warnf = func(format string, args ...interface{}) {
		safeFormat := sanitizeLogContent(format)
		safeArgs := sanitizeArgs(args)
		Log.Warn(context.Background(), safeFormat, safeArgs...)
	}

	Errorf = func(format string, args ...interface{}) {
		safeFormat := sanitizeLogContent(format)
		safeArgs := sanitizeArgs(args)
		Log.Error(context.Background(), safeFormat, safeArgs...)
	}

	Debugf = func(format string, args ...interface{}) {
		if os.Getenv("APP_DEBUG") == "true" || strings.EqualFold(os.Getenv("APP_DEBUG"), "1") {
			if Log.LogLevel >= logger.Info {
				safeFormat := sanitizeLogContent("[DEBUG] " + format)
				safeArgs := sanitizeArgs(args)
				Log.Info(context.Background(), safeFormat, safeArgs...)
			}
		}
	}
}

func GetLogger() *Logger {
	if Log == nil {
		InitLogger(nil)
	}
	return Log
}

// sanitizeLogContent 清理日志内容，移除不可打印字符和过长的内容
func sanitizeLogContent(content string) string {
	maxLength := 1000
	if len(content) > maxLength {
		content = content[:maxLength] + "...[截断]"
	}

	content = regexp.MustCompile(`[\x00-\x1f\x7f-\x9f]`).ReplaceAllString(content, "")

	// 替换可能导致问题的字符
	content = strings.ReplaceAll(content, "\u0000", "")

	// 确保所有字符都是可打印的
	var result strings.Builder
	for _, r := range content {
		if unicode.IsPrint(r) || unicode.IsSpace(r) {
			result.WriteRune(r)
		} else {
			result.WriteString(fmt.Sprintf("\\x%02x", r))
		}
	}

	return result.String()
}

// sanitizeArgs 清理日志参数
func sanitizeArgs(args []interface{}) []interface{} {
	sanitized := make([]interface{}, len(args))
	for i, arg := range args {
		switch v := arg.(type) {
		case string:
			// 只对字符串类型进行清理
			sanitized[i] = sanitizeLogContent(v)
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
			// 对于基本数据类型，保持原样
			sanitized[i] = v
		case error:
			// 对于错误类型，清理错误信息
			sanitized[i] = sanitizeLogContent(v.Error())
		default:
			// 对于其他类型，转换为字符串然后清理
			strVal := fmt.Sprintf("%v", v)
			sanitized[i] = sanitizeLogContent(strVal)
		}
	}
	return sanitized
}

// 便捷方法
func Info(format string, args ...interface{}) {
	safeFormat := sanitizeLogContent(format)
	safeArgs := sanitizeArgs(args)
	GetLogger().Info(context.Background(), safeFormat, safeArgs...)
}

func Warn(format string, args ...interface{}) {
	safeFormat := sanitizeLogContent(format)
	safeArgs := sanitizeArgs(args)
	GetLogger().Warn(context.Background(), safeFormat, safeArgs...)
}

func Error(format string, args ...interface{}) {
	safeFormat := sanitizeLogContent(format)
	safeArgs := sanitizeArgs(args)
	GetLogger().Error(context.Background(), safeFormat, safeArgs...)
}

func Debug(format string, args ...interface{}) {
	if os.Getenv("APP_DEBUG") == "true" || strings.EqualFold(os.Getenv("APP_DEBUG"), "1") {
		if GetLogger().LogLevel >= logger.Info {
			safeFormat := sanitizeLogContent("[DEBUG] " + format)
			safeArgs := sanitizeArgs(args)
			GetLogger().Info(context.Background(), safeFormat, safeArgs...)
		}
	}
}

// Init 使用默认配置初始化日志
func Init() {
	InitLogger(defaultConfig)
}

// InitWithConfig 使用自定义配置初始化日志
func InitWithConfig(config *Config) {
	InitLogger(config)
}

// Fatal 打印错误日志并退出程序
func Fatal(format string, args ...interface{}) {
	safeFormat := sanitizeLogContent(format)
	safeArgs := sanitizeArgs(args)
	GetLogger().Error(context.Background(), safeFormat, safeArgs...)
	os.Exit(1)
}

// ErrorExit 如果有错误就打印并退出
func ErrorExit(err error, msg string) {
	if err != nil {
		Fatal("%s: %v", msg, err)
	}
}

// ErrorReturn 如果有错误就打印并返回错误
func ErrorReturn(err error, msg string) error {
	if err != nil {
		Error("%s: %v", msg, err)
		return err
	}
	return nil
}

func PrintSeparator() {
	log.Println("")
}

func DefaultLogger(msg string) {
	log.Println(msg)
}
