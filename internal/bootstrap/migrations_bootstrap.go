package bootstrap

import (
	"pixelpunk/migrations"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"
)

func RunMigrations() {
	db := database.GetDB()
	if db == nil {
		logger.Warn("无法获取数据库连接，跳过迁移")
		return
	}

	if sqlDB, err := db.DB(); err != nil {
		logger.Warn("获取底层数据库连接失败，跳过迁移: %v", err)
		return
	} else if err := sqlDB.Ping(); err != nil {
		logger.Warn("数据库不可用（Ping失败），跳过迁移: %v", err)
		return
	}

	if err := migrations.RegisterAllMigrations(db); err != nil {
		logger.Warn("部分迁移可能执行失败: %v", err)
		return
	}
}
