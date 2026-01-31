package models

import (
	"crypto/md5"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Admin 管理员模型
type Admin struct {
	ID           int       `gorm:"type:integer;primarykey" json:"id"`
	Username     string    `gorm:"type:varchar(100);unique;not null" json:"username"`
	Password     string    `gorm:"type:varchar(255);not null" json:"-"`
	Role         string    `gorm:"type:varchar(50);not null;default:'admin';index" json:"role"` // 角色: admin, super_admin
	CompanyID    *int      `gorm:"type:integer" json:"company_id,omitempty"`                    // null表示超级管理员
	Company      *Company  `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	IsSuperAdmin bool      `gorm:"default:false" json:"is_super_admin"` // 是否超级管理员（保留用于兼容）
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// User 用户模型
type User struct {
	ID        int       `gorm:"type:integer;primarykey" json:"id"`
	CompanyID int       `gorm:"type:integer;not null;index" json:"company_id"` // 所属公司
	Company   Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Username  string    `gorm:"type:varchar(100);not null;index" json:"username"` // 允许重名
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	Role      string    `gorm:"type:varchar(50);not null;default:'user';index" json:"role"` // 角色: user
	Name      string    `gorm:"type:varchar(100)" json:"name"`
	Phone     string    `gorm:"type:varchar(20);index" json:"phone"` // 手机号（可选，用于区分重名用户）
	HasDrawn  bool      `gorm:"default:false" json:"has_drawn"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PrizeLevel 奖项等级（一等奖、二等奖等）
type PrizeLevel struct {
	ID          int       `gorm:"type:integer;primarykey" json:"id"`
	CompanyID   int       `gorm:"type:integer;not null;index" json:"company_id"` // 所属公司
	Company     Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Name        string    `gorm:"type:varchar(50);not null" json:"name"`
	Description string    `gorm:"type:varchar(200)" json:"description"`
	Probability float64   `gorm:"type:real;not null" json:"probability"`
	TotalStock  int       `gorm:"type:integer;not null" json:"total_stock"`
	UsedStock   int       `gorm:"type:integer;default:0" json:"used_stock"`
	SortOrder   int       `gorm:"type:integer;default:0" json:"sort_order"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Prize 具体奖品
type Prize struct {
	ID         int       `gorm:"type:integer;primarykey" json:"id"`
	LevelID    int       `gorm:"type:integer;not null" json:"level_id"`
	Name       string    `gorm:"type:varchar(100);not null" json:"name"`
	Image      string    `gorm:"type:varchar(255)" json:"image"`
	TotalStock int       `gorm:"type:integer;not null;default:0" json:"total_stock"` // 奖品总库存
	UsedStock  int       `gorm:"type:integer;default:0" json:"used_stock"`           // 已使用库存
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// DrawRecord 抽奖记录
type DrawRecord struct {
	ID        int        `gorm:"type:integer;primarykey" json:"id"`
	CompanyID int        `gorm:"type:integer;not null;index" json:"company_id"` // 所属公司
	Company   Company    `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	UserID    int        `gorm:"type:integer;not null" json:"user_id"`
	User      User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	LevelID   int        `json:"level_id"`
	Level     PrizeLevel `gorm:"foreignKey:LevelID" json:"level,omitempty"`
	PrizeID   int        `json:"prize_id"`
	Prize     Prize      `gorm:"foreignKey:PrizeID" json:"prize,omitempty"`
	IP        string     `gorm:"type:varchar(50)" json:"ip"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// AutoMigrate 自动迁移数据库表（仅在表结构变化时执行）
func AutoMigrate(db *gorm.DB) error {
	// 1. 确保表结构版本表存在
	if err := createSchemaVersionTableIfNotExists(db); err != nil {
		return fmt.Errorf("failed to create schema version table: %w", err)
	}

	// 2. 获取当前模型的哈希值
	currentHash := computeModelsHash()

	// 3. 检查是否需要迁移
	needsMigration, err := needsMigration(db, currentHash)
	if err != nil {
		return fmt.Errorf("failed to check migration status: %w", err)
	}

	if !needsMigration {
		log.Println("✅ 表结构未变化，跳过迁移")
		return nil
	}

	log.Println("🔄 检测到表结构变化，开始迁移...")

	// 4. 执行迁移
	if err := db.AutoMigrate(
		&Company{},
		&Admin{},
		&User{},
		&PrizeLevel{},
		&Prize{},
		&DrawRecord{},
		&OperationLog{},
	); err != nil {
		return fmt.Errorf("failed to run auto migration: %w", err)
	}

	// 5. 更新版本记录
	if err := updateSchemaVersion(db, currentHash); err != nil {
		return fmt.Errorf("failed to update schema version: %w", err)
	}

	log.Println("✅ 表结构迁移完成")
	return nil
}

// createSchemaVersionTableIfNotExists 创建表结构版本表
func createSchemaVersionTableIfNotExists(db *gorm.DB) error {
	sql := `
		CREATE TABLE IF NOT EXISTS schema_versions (
			id INT AUTO_INCREMENT PRIMARY KEY,
			model_hash VARCHAR(64) NOT NULL UNIQUE,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_hash (model_hash)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`
	return db.Exec(sql).Error
}

// computeModelsHash 计算模型的哈希值
func computeModelsHash() string {
	// 使用模型的结构信息生成哈希
	// 这里我们用模型名称列表作为简化版本
	models := []string{
		"Company", "Admin", "User", "PrizeLevel", "Prize", "DrawRecord", "OperationLog",
	}

	// TODO: 未来可以使用反射获取实际的结构信息
	// 当前使用简单的字符串拼接作为哈希依据
	hashStr := strings.Join(models, ",")
	return fmt.Sprintf("%x", md5.Sum([]byte(hashStr)))
}

// needsMigration 检查是否需要迁移
func needsMigration(db *gorm.DB, currentHash string) (bool, error) {
	var count int64
	err := db.Table("schema_versions").Where("model_hash = ?", currentHash).Count(&count).Error
	if err != nil {
		return false, err
	}

	// 如果找不到当前哈希记录，说明表结构可能变化了
	return count == 0, nil
}

// updateSchemaVersion 更新表结构版本
func updateSchemaVersion(db *gorm.DB, hash string) error {
	// 删除旧记录
	db.Exec("DELETE FROM schema_versions")

	// 插入新记录
	return db.Table("schema_versions").Create(map[string]interface{}{
		"model_hash": hash,
	}).Error
}
