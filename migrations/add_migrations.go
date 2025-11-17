package migrations

import (
	"pixelpunk/pkg/logger"

	"gorm.io/gorm"
)

// 注册迁移任务
type migrationTask struct {
	name string
	fn   func(*gorm.DB) error
}

// 注册的迁移列表
var registeredMigrations = []migrationTask{
	{"add_system_settings", AddSystemSettings},
}

// RegisterAllMigrations 注册所有迁移函数
func RegisterAllMigrations(db *gorm.DB) error {
	if err := EnsureMigrationTable(db); err != nil {
		logger.Error("确保迁移表存在失败: %v", err)
		return err
	}

	for _, task := range registeredMigrations {
		applied, err := IsMigrationApplied(db, task.name)
		if err != nil {
			logger.Error("检查迁移状态失败: %v", err)
			continue
		}

		if applied {
			continue
		}

		if err := task.fn(db); err != nil {
			logger.Error("迁移 %s 执行失败: %v", task.name, err)
		} else {
			if err := RecordMigration(db, task.name); err != nil {
				logger.Error("记录迁移状态失败: %s, 错误: %v", task.name, err)
			}
		}
	}
	return nil
}

func GetAllMigrationNames() []string {
	names := make([]string, len(registeredMigrations))
	for i, task := range registeredMigrations {
		names[i] = task.name
	}
	return names
}
