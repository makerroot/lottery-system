package migrations

import (
	"log"

	"gorm.io/gorm"
)

// Migration20260125ModifyUserUnique 修改用户表的唯一约束
type Migration20260125ModifyUserUnique struct{}

// Name 返回迁移名称
func (m *Migration20260125ModifyUserUnique) Name() string {
	return "20260125_modify_user_unique_constraint"
}

// Up 执行迁移
func (m *Migration20260125ModifyUserUnique) Up(tx *gorm.DB) error {
	log.Println("  → 检查旧的全局唯一索引...")

	// 步骤1: 检查旧索引是否存在
	var indexExists int64
	checkSQL := `
		SELECT COUNT(*)
		FROM information_schema.statistics
		WHERE table_schema = DATABASE()
		AND table_name = 'users'
		AND index_name = 'uq_users_username'
	`

	if err := tx.Raw(checkSQL).Count(&indexExists).Error; err != nil {
		return err
	}

	// 步骤2: 删除旧索引（如果存在）
	if indexExists > 0 {
		log.Println("  → 删除旧的全局唯一索引...")
		dropIndexSQL := `ALTER TABLE users DROP INDEX uq_users_username`
		if err := tx.Exec(dropIndexSQL).Error; err != nil {
			log.Printf("  ⚠️  删除旧索引失败: %v", err)
		} else {
			log.Println("  ✓ 删除旧索引成功")
		}
	} else {
		log.Println("  ℹ️  未发现旧索引")
	}

	// 步骤3: 创建新的复合唯一索引
	log.Println("  → 创建新的复合唯一索引 (company_id, username)...")

	// 检查新索引是否已存在
	var newIndexExists int64
	checkNewIndexSQL := `
		SELECT COUNT(*)
		FROM information_schema.statistics
		WHERE table_schema = DATABASE()
		AND table_name = 'users'
		AND index_name = 'idx_username_company'
	`

	if err := tx.Raw(checkNewIndexSQL).Count(&newIndexExists).Error; err != nil {
		return err
	}

	if newIndexExists == 0 {
		createIndexSQL := `
			CREATE UNIQUE INDEX idx_username_company
			ON users(company_id, username)
		`

		if err := tx.Exec(createIndexSQL).Error; err != nil {
			return err
		}
		log.Println("  ✓ 创建复合唯一索引成功")
	} else {
		log.Println("  ℹ️  复合唯一索引已存在")
	}

	return nil
}

// Down 回滚迁移（可选）
func (m *Migration20260125ModifyUserUnique) Down(tx *gorm.DB) error {
	// 回滚：删除复合唯一索引，恢复全局唯一索引
	log.Println("  → 删除复合唯一索引...")
	tx.Exec(`ALTER TABLE users DROP INDEX idx_username_company`)

	log.Println("  → 恢复全局唯一索引...")
	tx.Exec(`CREATE UNIQUE INDEX uq_users_username ON users(username)`)

	return nil
}
