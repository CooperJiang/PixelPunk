package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"pixelpunk/internal/controllers/websocket"
	"pixelpunk/internal/cron"
	middlewareInternal "pixelpunk/internal/middleware"
	"pixelpunk/internal/routes"
	"pixelpunk/internal/services/storage"
	"pixelpunk/pkg/cache"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/config"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/email"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/vector"

	"github.com/gin-gonic/gin"
	gormLogger "gorm.io/gorm/logger"
)

type App struct {
	Version string
	Engine  *gin.Engine
	Server  *http.Server
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewApp(version string) *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		Version: version,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (app *App) Initialize() error {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return fmt.Errorf("设置时区失败: %v", err)
	}
	time.Local = loc

	logger.InitWithConfig(&logger.Config{LogLevel: gormLogger.Info, Colorful: true})
	config.InitConfig()
	database.InitDB()

	installManager := common.GetInstallManager()
	if installManager.IsInstallMode() {
		if err := app.initializeHTTPServer(); err != nil {
			return fmt.Errorf("HTTP服务器初始化失败: %v", err)
		}
		return nil
	}

	cache.InitCache()
	RunMigrations()
	storage.CheckAndInitDefaultChannel()
	email.Init()
	websocket.InitWebSocketManager()
	InitAllServices(app.Version)
	cron.InitCronManager()

	if err := app.initializeHTTPServer(); err != nil {
		return fmt.Errorf("HTTP服务器初始化失败: %v", err)
	}

	return nil
}

func (app *App) initializeHTTPServer() error {
	gin.SetMode(config.GetConfig().App.Mode)
	app.Engine = gin.New()
	app.configureMiddleware()
	routes.RegisterRoutes(app.Engine)
	return nil
}

func (app *App) configureMiddleware() {
	app.Engine.Use(middlewareInternal.CORSMiddleware())
	app.Engine.Use(gin.Recovery())
	app.Engine.Use(errors.ErrorHandler())
	app.Engine.SetTrustedProxies([]string{"127.0.0.1", "localhost"})
}

func (app *App) Start() error {
	appCfg := config.GetConfig().App
	app.Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", appCfg.Port),
		Handler: app.Engine,
	}

	logger.Info("🚀 启动HTTP服务器，地址: %s", app.Server.Addr)
	if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("HTTP服务器启动失败: %v", err)
	}

	return nil
}

func (app *App) Shutdown(ctx context.Context) error {
	if app.Server != nil {
		if err := app.Server.Shutdown(ctx); err != nil {
			logger.Error("HTTP服务器关闭失败: %v", err)
		}
	}

	app.cancel()
	cron.Stop()

	if vectorEngine := vector.GetGlobalVectorEngine(); vectorEngine != nil {
		if err := vectorEngine.Close(); err != nil {
			logger.Error("关闭向量引擎失败: %v", err)
		}
	}

	if err := database.Close(); err != nil {
		logger.Error("关闭数据库连接失败: %v", err)
	}

	if err := cache.Close(); err != nil {
		logger.Error("关闭缓存连接失败: %v", err)
	}

	return nil
}
