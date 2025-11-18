package database

import (
	"fmt"
	"os"
	"path/filepath"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/config"
	log "pixelpunk/pkg/logger"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	_ "modernc.org/sqlite"
)

var DB *gorm.DB

func GetDB() *gorm.DB {
	return DB
}

func getGormConfig() *gorm.Config {
	return &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NowFunc: func() time.Time {
			return time.Now().In(time.Local)
		},
		SkipDefaultTransaction: false, // SQLite需要事务保护
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}
}

func getDialector(dbType, host, username, password, name, path string, port int) (gorm.Dialector, error) {
	switch dbType {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FShanghai&sql_mode=%%27STRICT_TRANS_TABLES,NO_ZERO_DATE,NO_ZERO_IN_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION%%27",
			username, password, host, port, name)
		return mysql.Open(dsn), nil
	case "sqlite":
		dbPath := path
		if !filepath.IsAbs(dbPath) {
			dbPath = filepath.Join(".", dbPath)
		}
		// 为SQLite添加并发优化参数
		// _busy_timeout: 等待锁的超时时间（毫秒），10秒
		// _journal_mode=WAL: 启用WAL模式，支持并发读取
		// _foreign_keys=on: 启用外键约束
		// _cache_size: 设置缓存大小（页面数）
		// _synchronous=NORMAL: 同步模式，平衡性能和安全性
		dsn := fmt.Sprintf("%s?_busy_timeout=10000&_journal_mode=WAL&_foreign_keys=on&_cache_size=1000&_synchronous=NORMAL&_temp_store=MEMORY", dbPath)

		// 使用纯Go版本的SQLite驱动，不需要CGO
		return sqlite.Dialector{
			DriverName: "sqlite",
			DSN:        dsn,
		}, nil
	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", dbType)
	}
}

// checkConfigFileExists 检查配置文件是否存在（支持多路径）
func checkConfigFileExists() bool {
	configPaths := []string{
		"configs/config.yaml", // 新路径
		"config.yaml",         // 旧路径（向后兼容）
	}

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}
	return false
}

// InitDB 初始化数据库连接
func InitDB() {
	cfg := config.GetConfig().Database
	installManager := common.GetInstallManager()

	// 首先检查配置文件是否存在（支持多路径）
	configExists := checkConfigFileExists()
	if !configExists {
		log.Info("配置文件不存在，进入安装模式")
		installManager.SetInstallMode(true)
		installManager.SetSystemInstalled(false)
		return
	}

	if cfg.Type == "" {
		// 兼容老版本配置：如果没有type字段，但有MySQL配置，默认为mysql
		if cfg.Host != "" && cfg.Username != "" && cfg.Name != "" {
			log.Info("检测到老版本配置，自动设置数据库类型为MySQL")
			cfg.Type = "mysql"
		} else {
			// 如果配置文件存在但数据库配置不完整，进入安装模式
			log.Info("数据库类型未配置，进入安装模式")
			installManager.SetInstallMode(true)
			installManager.SetSystemInstalled(false)
			return
		}
	}

	// 根据数据库类型检查必要配置
	switch cfg.Type {
	case "mysql":
		if cfg.Host == "" || cfg.Username == "" || cfg.Name == "" {
			if configExists {
				log.Warn("配置文件存在但MySQL配置不完整，请检查配置文件")
				log.Warn("如需重新配置，请删除 configs/config.yaml 后重启应用")
				return
			} else {
				log.Info("MySQL数据库配置不完整，进入安装模式")
				installManager.SetInstallMode(true)
				installManager.SetSystemInstalled(false)
				return
			}
		}
	case "sqlite":
		if cfg.Path == "" {
			if configExists {
				log.Warn("配置文件存在但SQLite路径未配置，请检查配置文件")
				log.Warn("如需重新配置，请删除 configs/config.yaml 后重启应用")
				return
			} else {
				log.Info("SQLite数据库路径未配置，进入安装模式")
				installManager.SetInstallMode(true)
				installManager.SetSystemInstalled(false)
				return
			}
		}

		dbPath := cfg.Path
		if !filepath.IsAbs(dbPath) {
			dbPath = filepath.Join(".", dbPath)
		}

		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			if configExists {
				// 配置文件存在但数据库文件不存在，尝试创建数据库文件而不是进入安装模式
				log.Info("SQLite数据库文件不存在，尝试创建: %s", dbPath)
				if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
					log.Warn("创建数据库目录失败: %v", err)
					return
				}
			} else {
				log.Info("SQLite数据库文件不存在，进入安装模式")
				installManager.SetInstallMode(true)
				installManager.SetSystemInstalled(false)
				return
			}
		}
	default:
		log.Warn("不支持的数据库类型: %s", cfg.Type)
		if configExists {
			log.Warn("如需重新配置，请删除 configs/config.yaml 后重启应用")
			return
		} else {
			log.Info("进入安装模式")
			installManager.SetInstallMode(true)
			installManager.SetSystemInstalled(false)
			return
		}
	}

	gormConfig := getGormConfig()
	dialector, err := getDialector(cfg.Type, cfg.Host, cfg.Username, cfg.Password, cfg.Name, cfg.Path, cfg.Port)
	if err != nil {
		log.Warn("数据库配置错误，进入安装模式: %v", err)
		installManager.SetInstallMode(true)
		installManager.SetSystemInstalled(false)
		return
	}

	DB, err = gorm.Open(dialector, gormConfig)
	if err != nil {
		if configExists {
			log.Warn("数据库连接失败: %v", err)
			log.Warn("请检查数据库服务是否运行，或配置是否正确")
			log.Warn("如需重新配置，请删除 configs/config.yaml 后重启应用")
			return
		} else {
			log.Warn("数据库连接失败，进入安装模式: %v", err)
			installManager.SetInstallMode(true)
			installManager.SetSystemInstalled(false)
			return
		}
	}

	// 为SQLite配置连接池参数，避免并发锁定
	if cfg.Type == "sqlite" {
		sqlDB, err := DB.DB()
		if err == nil {
			// SQLite建议的连接池配置，避免并发写入冲突
			sqlDB.SetMaxOpenConns(1)            // SQLite只允许一个写连接
			sqlDB.SetMaxIdleConns(1)            // 保持一个空闲连接
			sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生存时间
		}
	}

	if cfg.Type == "mysql" {
		sqlDB, err := DB.DB()
		if err == nil {
			sqlDB.SetMaxOpenConns(150)
			sqlDB.SetMaxIdleConns(50)
			sqlDB.SetConnMaxLifetime(time.Hour)
			sqlDB.SetConnMaxIdleTime(10 * time.Minute)
		}
	}

	if err := autoMigrate(); err != nil {
		if configExists {
			log.Warn("数据库迁移失败: %v", err)
			log.Warn("请检查数据库权限或表结构，如需重新配置，请删除 configs/config.yaml 后重启应用")
			if status, checkErr := checkDatabaseStatus(); checkErr == nil && status.HasUsers {
				log.Warn("检测到数据库中已有用户数据，系统将继续运行，但建议修复迁移问题")
				installManager.SetInstallMode(false)
				installManager.SetSystemInstalled(true)
				return
			}
			return
		} else {
			log.Warn("数据库迁移失败，进入安装模式: %v", err)
			installManager.SetInstallMode(true)
			installManager.SetSystemInstalled(false)
			return
		}
	}

	if err := checkAdminUserExists(); err != nil {
		if configExists {
			log.Warn("系统初始化检查: %v", err)
			log.Warn("请访问Web界面完成初始化，如需重新配置，请删除 configs/config.yaml 后重启应用")
			installManager.SetInstallMode(true)
			installManager.SetSystemInstalled(false)
			return
		} else {
			log.Warn("系统需要初始化: %v", err)
			installManager.SetInstallMode(true)
			installManager.SetSystemInstalled(false)
			return
		}
	}

	installManager.SetInstallMode(false)
	installManager.SetSystemInstalled(true)
}

// TestDatabaseConnection 测试数据库连接
func TestDatabaseConnection(dbType, host string, port int, username, password, name, path string) error {
	gormConfig := getGormConfig()
	dialector, err := getDialector(dbType, host, username, password, name, path, port)
	if err != nil {
		return err
	}

	testDB, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	sqlDB, err := testDB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库ping失败: %v", err)
	}

	sqlDB.Close()
	return nil
}

// ReconnectDatabase 重新连接数据库（用于安装完成后）
func ReconnectDatabase() error {
	cfg := config.GetConfig().Database
	installManager := common.GetInstallManager()
	gormConfig := getGormConfig()

	dialector, err := getDialector(cfg.Type, cfg.Host, cfg.Username, cfg.Password, cfg.Name, cfg.Path, cfg.Port)
	if err != nil {
		return err
	}

	DB, err = gorm.Open(dialector, gormConfig)
	if err != nil {
		return fmt.Errorf("重新连接数据库失败: %v", err)
	}

	// 为SQLite配置连接池参数，避免并发锁定
	if cfg.Type == "sqlite" {
		sqlDB, err := DB.DB()
		if err == nil {
			// SQLite建议的连接池配置，避免并发写入冲突
			sqlDB.SetMaxOpenConns(1)            // SQLite只允许一个写连接
			sqlDB.SetMaxIdleConns(1)            // 保持一个空闲连接
			sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生存时间
		}
	}

	if cfg.Type == "mysql" {
		sqlDB, err := DB.DB()
		if err == nil {
			sqlDB.SetMaxOpenConns(150)
			sqlDB.SetMaxIdleConns(50)
			sqlDB.SetConnMaxLifetime(time.Hour)
			sqlDB.SetConnMaxIdleTime(10 * time.Minute)
		}
	}

	if err := autoMigrate(); err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	// 在安装过程中不检查 root 用户，因为会在安装控制器中单独创建
	if !installManager.IsInstalling() {
		// 只有在非安装过程中才检查并创建 root 用户
		if err := checkAdminUserExists(); err != nil {
			return fmt.Errorf("系统初始化检查失败: %v", err)
		}

		installManager.SetInstallMode(false)
		installManager.SetSystemInstalled(true)
	}

	log.Info("数据库重新连接成功")

	return nil
}

func autoMigrate() error {
	models := []interface{}{
		&models.User{},
		&models.File{},
		&models.FileStats{},
		&models.FileDownloadLog{},
		&models.Folder{},
		&models.UserUsageStats{},
		&models.UserSettings{},
		&models.GlobalStats{},
		&models.APIKey{},
		&models.RandomImageAPI{},
		&models.StorageChannel{},
		&models.StorageConfigItem{},
		&models.Setting{},
		&models.FileAIInfo{},
		&models.FileTaggingLog{},
		&models.UserAccessControl{},
		&models.Share{},
		&models.ShareItem{},
		&models.ShareAccessLog{},
		&models.ShareVisitorInfo{},
		&models.ShareAccessToken{},
		&models.UploadSession{},
		&models.UploadChunk{},
		&models.FileVector{},
		&models.VectorProcessingLog{},
		&models.VectorVerificationTask{},
		&models.ReviewLog{},
		&models.Message{},
		&models.MessageTemplate{},
		&models.ActivityLog{},
		&models.GuestUploadLimit{},
		&models.GuestUploadLog{},
		&models.UserBandwidthUsage{},
		&models.GlobalTag{},
		&models.UserTagReference{},
		&models.TagCategoryRelation{},
		&models.FileGlobalTagRelation{},
		&models.GlobalTagOperationLog{},
		&models.GlobalTagStatsCache{},
		&models.FileCategory{},
		&models.CategoryTemplate{},
		&models.FileCategoryRelation{},
		// 队列表模型（改为自动迁移）
		&models.AIJob{},
		&models.VectorJob{},
		&models.Announcement{},
	}

	silentDB := DB.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)})

	for _, model := range models {
		if err := silentDB.AutoMigrate(model); err != nil {
			if isIndexError(err) {
				continue
			}
			return err
		}
	}

	return nil
}

func isIndexError(err error) bool {
	errorMsg := err.Error()
	indexErrors := []string{
		"Can't DROP", "BLOB/TEXT column", "key specification",
		"Duplicate key name", "already exists",
	}

	for _, errType := range indexErrors {
		if strings.Contains(errorMsg, errType) {
			return true
		}
	}
	return false
}

// DatabaseStatus 数据库状态
type DatabaseStatus struct {
	Exists       bool // 数据库文件/连接是否存在
	HasTables    bool // 是否有表结构
	HasUsers     bool // 是否有用户数据
	HasAdminUser bool // 是否有管理员用户
}

// checkDatabaseStatus 检查数据库状态
func checkDatabaseStatus() (*DatabaseStatus, error) {
	status := &DatabaseStatus{
		Exists:       false,
		HasTables:    false,
		HasUsers:     false,
		HasAdminUser: false,
	}

	if DB == nil {
		return status, nil
	}

	status.Exists = true

	if DB.Migrator().HasTable(&models.User{}) {
		status.HasTables = true

		var userCount int64
		if err := DB.Model(&models.User{}).Count(&userCount).Error; err != nil {
			return status, err
		}

		if userCount > 0 {
			status.HasUsers = true

			var adminCount int64
			if err := DB.Model(&models.User{}).Where("role IN ?", []int{1, 2}).Count(&adminCount).Error; err != nil {
				return status, err
			}

			status.HasAdminUser = adminCount > 0
		}
	}

	return status, nil
}

// checkAdminUserExists 检查管理员用户是否存在
func checkAdminUserExists() error {
	status, err := checkDatabaseStatus()
	if err != nil {
		return fmt.Errorf("检查数据库状态失败: %v", err)
	}

	if !status.Exists || !status.HasTables {
		return fmt.Errorf("数据库表结构未初始化，请通过安装向导完成设置")
	}

	// 如果没有任何用户，说明需要创建管理员账户
	if !status.HasUsers {
		return fmt.Errorf("系统未配置管理员账户，请通过Web界面完成初始化")
	}

	// 如果有用户但没有管理员用户，说明数据不完整
	if !status.HasAdminUser {
		return fmt.Errorf("缺少管理员用户，请通过安装向导重新配置")
	}

	return nil
}

// Close 关闭数据库连接
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
