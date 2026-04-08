package hdmigration

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SchemaMigration 版本追踪表，每个连接只保留一条记录
type SchemaMigration struct {
	ID          uint      `gorm:"primarykey"`
	Name        string    `gorm:"column:name;type:varchar(64);not null;uniqueIndex;comment:连接名称"`
	Version     int       `gorm:"column:version;not null;comment:版本号"`
	Description string    `gorm:"column:description;type:varchar(255);comment:变更说明"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime;autoUpdateTime;comment:更新时间"`
}

func (SchemaMigration) TableName() string {
	return "schema_migrations"
}

// CheckVersion 确保版本表存在并检查是否需要迁移
// name: 连接名称，用于区分不同数据库连接的迁移版本
// 返回: (needsMigration, appliedVersion, error)
func CheckVersion(db *gorm.DB, name string, targetVersion int) (bool, int, error) {
	if err := db.AutoMigrate(&SchemaMigration{}); err != nil {
		return false, 0, fmt.Errorf("migrate schema_migrations table: %w", err)
	}

	var record SchemaMigration
	err := db.Where("name = ?", name).First(&record).Error
	if err != nil {
		// 没有记录，说明首次运行，需要迁移
		return true, 0, nil
	}

	if record.Version >= targetVersion {
		return false, record.Version, nil
	}

	return true, record.Version, nil
}

// RecordVersion 记录迁移成功的版本号（upsert，每个 name 只保留一条）
// desc: 可选的变更说明，如 "新增 order 表 xxx 字段"
func RecordVersion(db *gorm.DB, name string, version int, desc string) error {
	record := SchemaMigration{
		Name:        name,
		Version:     version,
		Description: desc,
		UpdatedAt:   time.Now(),
	}
	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"version", "description", "updated_at"}),
	}).Create(&record).Error; err != nil {
		return fmt.Errorf("record schema version: %w", err)
	}
	return nil
}
