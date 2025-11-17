package stats

import (
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"

	"gorm.io/gorm"
)

type StatsAdapter struct {
	statsService *GlobalStatsService
}

var adapter *StatsAdapter

func GetStatsAdapter() *StatsAdapter {
	if adapter == nil {
		adapter = &StatsAdapter{
			statsService: NewGlobalStatsService(GetDB()),
		}
	}
	return adapter
}

func GetDB() *gorm.DB {
	db := database.GetDB()
	return db
}

func (a *StatsAdapter) RecordFileCreated(size int64) {
	go func() {
		if err := a.statsService.IncrementFileStats(size); err != nil {
			logger.Warn("记录文件创建统计失败: %v", err)
		}
	}()
}

func (a *StatsAdapter) RecordFolderCreated() {
	go func() {
		if err := a.statsService.IncrementFolderStats(); err != nil {
			logger.Warn("记录文件夹创建统计失败: %v", err)
		}
	}()
}

func (a *StatsAdapter) RecordUserCreated() {
	go func() {
		if err := a.statsService.IncrementUserStats(); err != nil {
			logger.Warn("记录用户创建统计失败: %v", err)
		}
	}()
}

func (a *StatsAdapter) RecordFileViewed(bandwidth int64) {
	go func() {
		if err := a.statsService.IncrementViewStats(bandwidth); err != nil {
			logger.Warn("记录文件访问统计失败: %v", err)
		}
	}()
}
