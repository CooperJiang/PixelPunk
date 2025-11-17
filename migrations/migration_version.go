package migrations

import (
	"time"

	"gorm.io/gorm"
)

// MigrationRecord 表示一条迁移记录
type MigrationRecord struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	AppliedAt time.Time `gorm:"not null"`
}

// TableName 设置迁移记录表名
func (MigrationRecord) TableName() string {
	return "migration_versions"
}

// EnsureMigrationTable 确保迁移版本表存在
func EnsureMigrationTable(db *gorm.DB) error {
	if !db.Migrator().HasTable(&MigrationRecord{}) {
		if err := db.Migrator().CreateTable(&MigrationRecord{}); err != nil {
			return err
		}
	}
	return nil
}

// IsMigrationApplied 检查指定名称的迁移是否已应用
func IsMigrationApplied(db *gorm.DB, name string) (bool, error) {
	if err := EnsureMigrationTable(db); err != nil {
		return false, err
	}

	var count int64
	err := db.Model(&MigrationRecord{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// RecordMigration 记录一次成功的迁移
func RecordMigration(db *gorm.DB, name string) error {
	if err := EnsureMigrationTable(db); err != nil {
		return err
	}

	// 创建迁移记录
	record := MigrationRecord{
		Name:      name,
		AppliedAt: time.Now(),
	}

	return db.Create(&record).Error
}
