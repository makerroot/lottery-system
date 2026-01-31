package migrations

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// Migration 迁移接口
type Migration interface {
	// Name 返回迁移名称
	Name() string

	// Up 执行迁移
	Up(*gorm.DB) error

	// Down 回滚迁移（可选）
	Down(*gorm.DB) error
}

// migrationRegistry 注册的迁移列表
var migrationRegistry []Migration

// RegisterMigration 注册迁移
func RegisterMigration(migration Migration) {
	migrationRegistry = append(migrationRegistry, migration)
}

// RunMigrations 执行所有未运行的迁移
func RunMigrations(db *gorm.DB) error {
	log.Println("🔄 检查数据库迁移...")

	// 确保迁移记录表存在
	if err := createMigrationTableIfNotExists(db); err != nil {
		return fmt.Errorf("创建迁移表失败: %w", err)
	}

	// 获取已执行的迁移列表
	executedMigrations, err := getExecutedMigrations(db)
	if err != nil {
		return fmt.Errorf("获取已执行迁移失败: %w", err)
	}

	// 执行未运行的迁移
	for _, migration := range migrationRegistry {
		migrationName := migration.Name()

		if isExecuted(executedMigrations, migrationName) {
			log.Printf("✓ 迁移已执行: %s", migrationName)
			continue
		}

		log.Printf("📌 执行迁移: %s", migrationName)

		// 在事务中执行迁移
		tx := db.Begin()
		if err := migration.Up(tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("迁移 %s 失败: %w", migrationName, err)
		}

		// 记录迁移
		if err := recordMigration(tx, migrationName); err != nil {
			tx.Rollback()
			return fmt.Errorf("记录迁移 %s 失败: %w", migrationName, err)
		}

		// 提交事务
		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("提交迁移 %s 失败: %w", migrationName, err)
		}

		log.Printf("✅ 迁移成功: %s", migrationName)
	}

	log.Println("✅ 所有迁移检查完成")
	return nil
}

// createMigrationTableIfNotExists 创建迁移记录表（如果不存在）
func createMigrationTableIfNotExists(db *gorm.DB) error {
	// 使用原生SQL创建表以避免GORM的自动迁移干扰
	sql := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			executed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_name (name)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`

	if err := db.Exec(sql).Error; err != nil {
		return err
	}

	return nil
}

// getExecutedMigrations 获取已执行的迁移列表
func getExecutedMigrations(db *gorm.DB) ([]string, error) {
	var migrations []string

	err := db.Model(&SchemaMigration{}).
		Select("name").
		Order("executed_at ASC").
		Pluck("name", &migrations).
		Error

	return migrations, err
}

// isExecuted 检查迁移是否已执行
func isExecuted(executedMigrations []string, migrationName string) bool {
	for _, name := range executedMigrations {
		if name == migrationName {
			return true
		}
	}
	return false
}

// recordMigration 记录已执行的迁移
func recordMigration(db *gorm.DB, migrationName string) error {
	migration := &SchemaMigration{
		Name: migrationName,
	}

	return db.Create(migration).Error
}

// SchemaMigration 迁移记录模型
type SchemaMigration struct {
	ID         int    `gorm:"type:integer;primarykey"`
	Name       string `gorm:"type:varchar(255);not null;unique;index:idx_name"`
	ExecutedAt string `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
}

// TableName 指定表名
func (SchemaMigration) TableName() string {
	return "schema_migrations"
}
