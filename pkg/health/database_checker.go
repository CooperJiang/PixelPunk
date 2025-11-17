package health

import (
	"pixelpunk/pkg/database"
	"time"
)

// DatabaseChecker 数据库健康检查器
type DatabaseChecker struct{}

// Name 返回检查项名称
func (c *DatabaseChecker) Name() string {
	return "database"
}

// Check 执行健康检查
func (c *DatabaseChecker) Check() (Status, map[string]interface{}) {
	db := database.GetDB()
	if db == nil {
		return StatusDown, map[string]interface{}{
			"error": "数据库连接未初始化",
		}
	}

	sqlDB, err := db.DB()
	if err != nil {
		return StatusDown, map[string]interface{}{
			"error": err.Error(),
		}
	}

	start := time.Now()
	if err := sqlDB.Ping(); err != nil {
		return StatusDown, map[string]interface{}{
			"error": err.Error(),
		}
	}
	pingTime := time.Since(start)

	stats := sqlDB.Stats()

	details := map[string]interface{}{
		"ping_time_ms":        pingTime.Milliseconds(),
		"max_open_conns":      stats.MaxOpenConnections,
		"open_conns":          stats.OpenConnections,
		"in_use":              stats.InUse,
		"idle":                stats.Idle,
		"wait_count":          stats.WaitCount,
		"wait_duration_ms":    stats.WaitDuration.Milliseconds(),
		"max_idle_closed":     stats.MaxIdleClosed,
		"max_lifetime_closed": stats.MaxLifetimeClosed,
	}

	// 如果响应时间超过100ms, 视为性能下降
	if pingTime > 100*time.Millisecond {
		return StatusDegraded, details
	}

	return StatusUp, details
}

// Type 返回检查类型
func (c *DatabaseChecker) Type() CheckType {
	return CheckTypeBasic
}

// init 注册数据库健康检查器
func init() {
	RegisterChecker(&DatabaseChecker{})
}
