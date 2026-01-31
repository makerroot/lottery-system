package migrations

import (
	"log"

	"gorm.io/gorm"
)

// Migration20260125AddPrizeStock 为奖品表添加库存字段
type Migration20260125AddPrizeStock struct{}

// Name 返回迁移名称
func (m *Migration20260125AddPrizeStock) Name() string {
	return "20260125_add_prize_stock"
}

// Up 执行迁移
func (m *Migration20260125AddPrizeStock) Up(tx *gorm.DB) error {
	log.Println("  → 检查 prizes 表的库存字段...")

	// 检查列是否已存在
	var columnExists int64
	checkSQL := `
		SELECT COUNT(*)
		FROM information_schema.columns
		WHERE table_schema = DATABASE()
		AND table_name = 'prizes'
		AND column_name = 'total_stock'
	`

	if err := tx.Raw(checkSQL).Count(&columnExists).Error; err != nil {
		return err
	}

	if columnExists == 0 {
		log.Println("  → 添加 total_stock 字段...")
		if err := tx.Exec("ALTER TABLE prizes ADD COLUMN total_stock INT NOT NULL DEFAULT 0").Error; err != nil {
			return err
		}
		log.Println("  ✓ 添加 total_stock 字段成功")
	} else {
		log.Println("  ℹ️  total_stock 字段已存在")
	}

	// 检查 used_stock 列
	var usedStockColumnExists int64
	checkUsedStockSQL := `
		SELECT COUNT(*)
		FROM information_schema.columns
		WHERE table_schema = DATABASE()
		AND table_name = 'prizes'
		AND column_name = 'used_stock'
	`

	if err := tx.Raw(checkUsedStockSQL).Count(&usedStockColumnExists).Error; err != nil {
		return err
	}

	if usedStockColumnExists == 0 {
		log.Println("  → 添加 used_stock 字段...")
		if err := tx.Exec("ALTER TABLE prizes ADD COLUMN used_stock INT NOT NULL DEFAULT 0").Error; err != nil {
			return err
		}
		log.Println("  ✓ 添加 used_stock 字段成功")
	} else {
		log.Println("  ℹ️  used_stock 字段已存在")
	}

	return nil
}

// Down 回滚迁移
func (m *Migration20260125AddPrizeStock) Down(tx *gorm.DB) error {
	log.Println("  → 删除库存字段...")
	tx.Exec("ALTER TABLE prizes DROP COLUMN total_stock")
	tx.Exec("ALTER TABLE prizes DROP COLUMN used_stock")
	return nil
}
