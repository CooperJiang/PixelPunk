package cron

import (
	"pixelpunk/internal/services/ai"
	"pixelpunk/internal/services/stats"
	"pixelpunk/internal/services/tag"
	vectorSvc "pixelpunk/internal/services/vector"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

var cronManager *cron.Cron
var db *gorm.DB
var taggingService *ai.TaggingService

/* InitCronManager 初始化定时任务管理器 */
func InitCronManager() {
	db = database.GetDB()
	if db == nil {
		logger.Error("定时任务管理器初始化失败：无法获取数据库连接")
		return
	}

	cronManager = cron.New(cron.WithSeconds())

	taggingService = ai.NewTaggingServiceWithConfig(db) // 使用配置的并发数
	if taggingService == nil {
		logger.Error("AI文件标记服务创建失败")
		return
	}

	ai.SetGlobalTaggingService(taggingService)

	registerTasks()

	cronManager.Start()
}

func registerTasks() {

	registerStatsTask()

	registerImageTaggingTask()

	registerShareTask()

	registerChunkedUploadCleanupTask()

	registerVectorVerificationTask()

	registerImageCleanupTask()

	registerVectorQueueTask()

	registerVectorReconcileTasks()

	registerTagUsageCountCalibrationTask()

}

func registerStatsTask() {
	statsService := stats.NewGlobalStatsService(db)

	_, err := cronManager.AddFunc("0 0 1 * * *", func() {
		if err := statsService.ReconcileAllStats(); err != nil {
			logger.Warn("全量统计数据校准失败: %v", err)
		}
	})
	if err != nil {
		logger.Warn("注册全量统计校准任务失败: %v", err)
	}

	_, err = cronManager.AddFunc("0 0 * * * *", func() {
		if err := statsService.ReconcileTodayStats(); err != nil {
			logger.Warn("今日统计数据校准失败: %v", err)
		}
	})
	if err != nil {
		logger.Warn("注册今日统计校准任务失败: %v", err)
	}
}

func registerImageTaggingTask() {
	_, err := cronManager.AddFunc("0 */1 * * * *", func() {
		ai.ScheduledTaggingTask()
	})

	if err != nil {
		logger.Error("注册文件标记定时任务失败: %v", err)
	}

	_, err = cronManager.AddFunc("0 0 * * * *", func() {
		count, err := ai.ResetStuckPendingFiles(60) // 60分钟未更新的视为卡住
		if err != nil {
			logger.Error("重置卡住的文件处理任务失败: %v", err)
		} else if count > 0 {
		}
	})
	if err != nil {
		logger.Error("注册卡住任务检查任务失败: %v", err)
	}
}

func registerVectorQueueTask() {
	_, err := cronManager.AddFunc("0 0/15 * * * *", func() {
		if svc := vectorSvc.GetGlobalVectorQueueService(); svc != nil && !svc.IsPaused() {
			if n, err := svc.EnqueueAllPending(1000); err == nil && n > 0 {
				logger.Info("向量队列对账入队: %d", n)
			}
		}
	})
	if err != nil {
		logger.Warn("注册向量队列定时任务失败: %v", err)
	}
}

func registerVectorReconcileTasks() {
	_, err := cronManager.AddFunc("0 0/15 * * * *", func() {
		if svc := vectorSvc.GetGlobalVectorQueueService(); svc != nil && !svc.IsPaused() {
			if total, enq, err := svc.ReconcileMissing(1000, false); err == nil {
				if total > 0 {
					logger.Info("向量补齐缺失：发现 %d，入队 %d", total, enq)
				}
			} else {
				logger.Warn("向量补齐缺失失败: %v", err)
			}
		}
	})
	_ = err // 忽略重复注册导致的警告

	_, err = cronManager.AddFunc("0 30 3 * * *", func() {
		if svc := vectorSvc.GetGlobalVectorQueueService(); svc != nil {
			if total, removed, err := svc.CleanOrphans(2000, false); err == nil {
				if total > 0 {
					logger.Info("向量清理孤儿：发现 %d，删除 %d", total, removed)
				}
			} else {
				logger.Warn("向量清理孤儿失败: %v", err)
			}
		}
	})
	_ = err // 忽略重复注册导致的警告
}

func registerChunkedUploadCleanupTask() {
	cleanupJob := NewChunkedUploadCleanupJob()

	_, err := cronManager.AddFunc(cleanupJob.GetSchedule(), func() {
		if err := cleanupJob.Execute(); err != nil {
			logger.Error("分片上传清理任务执行失败: %v", err)
		}
	})
	if err != nil {
		logger.Error("注册分片上传清理任务失败: %v", err)
	}
}

func registerVectorVerificationTask() {
	verificationJob := NewVectorVerificationJob()

	_, err := cronManager.AddFunc(verificationJob.GetSchedule(), func() {
		if err := verificationJob.Execute(); err != nil {
			logger.Error("向量验证任务执行失败: %v", err)
		}
	})
	if err != nil {
		logger.Error("注册向量验证任务失败: %v", err)
	}
}

func registerImageCleanupTask() {
	cleanupJob := NewImageCleanupJob()

	_, err := cronManager.AddFunc(cleanupJob.GetSchedule(), func() {
		if err := cleanupJob.Execute(); err != nil {
			logger.Error("文件清理任务执行失败: %v", err)
		}
	})
	if err != nil {
		logger.Error("注册文件清理任务失败: %v", err)
	}
}

func registerTagUsageCountCalibrationTask() {
	tagService := tag.NewFileGlobalTagService()

	// 每天凌晨2点执行标签使用次数校准任务
	_, err := cronManager.AddFunc("0 0 2 * * *", func() {
		if err := tagService.CalibrateAllTagUsageCount(); err != nil {
		} else {
		}
	})
	if err != nil {
		logger.Error("注册标签使用次数校准任务失败: %v", err)
	}
}

/* Stop 停止所有定时任务 */
func Stop() {
	if cronManager != nil {
		cronManager.Stop()
	}

	if taggingService != nil {
		taggingService.Stop()
	}
}
