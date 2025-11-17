package bootstrap

import (
	"time"

	ai "pixelpunk/internal/services/ai"
	"pixelpunk/internal/services/message"
	"pixelpunk/internal/services/setting"
	"pixelpunk/internal/services/user"
	vectorSvc "pixelpunk/internal/services/vector"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/vector"
)

func InitAllServices(appVersion string) {
	user.InitUserService()
	setting.InitSettingService()
	initMessageService()
	initVectorEngine()
	ai.RegisterAISettingHooks()
	vectorSvc.RegisterVectorConfigHooks()
	if err := ai.InitGlobalTaggingQueue(); err != nil {
		logger.Warn("AI打标队列初始化警告: %v", err)
	}
}

func initVectorEngine() {
	vectorEnabled := setting.GetBoolDirectFromDB("vector", "vector_enabled", false)
	if !vectorEnabled {
		logger.Info("向量功能未启用，跳过初始化")
		return
	}

	qdrantURL := setting.GetStringDirectFromDB("vector", "qdrant_url", "")
	if qdrantURL == "" {
		logger.Warn("向量功能已启用，但未配置 qdrant_url，跳过初始化")
		return
	}

	qdrantTimeout := setting.GetIntDirectFromDB("vector", "qdrant_timeout", 30)

	if err := vector.InitQdrantVectorEngine(qdrantURL, qdrantTimeout); err != nil {
		logger.Error("向量引擎初始化失败: %v", err)
		return
	}

	engine := vector.GetGlobalVectorEngine()
	if engine == nil {
		logger.Error("向量引擎初始化后仍为nil")
		return
	}


	if err := vectorSvc.InitGlobalVectorQueue(); err != nil {
		logger.Error("向量队列初始化失败: %v", err)
		return
	}


	if queueSvc := vectorSvc.GetGlobalVectorQueueService(); queueSvc != nil {
		go func() {
			time.Sleep(3 * time.Second)
			if n, err := queueSvc.EnqueueAllPending(1000); err == nil && n > 0 {
				logger.Info("启动时自动扫描入队 %d 个待处理向量任务", n)
			}
		}()
	}
}

func initMessageService() {
	message.InitMessageService()

	templateService := message.GetTemplateService()
	if err := templateService.InitDefaultTemplates(); err != nil {
		logger.Warn("初始化默认消息模板失败: %v", err)
	}
}
