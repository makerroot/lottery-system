package migrations

import (
	"log"

	"gorm.io/gorm"
)

// Migration20260131AllowDuplicateUsername 允许同一公司内有重名用户
type Migration20260131AllowDuplicateUsername struct{}

// Name 返回迁移名称
func (m *Migration20260131AllowDuplicateUsername) Name() string {
	return "20260131_allow_duplicate_username"
}

// Up 执行迁移
func (m *Migration20260131AllowDuplicateUsername) Up(tx *gorm.DB) error {
	log.Println("  → 检查旧的唯一索引...")

	// 步骤1: 检查旧索引是否存在 (company_id, username)
	var indexExists int64
	checkSQL := `
		SELECT COUNT(*)
		FROM information_schema.statistics
		WHERE table_schema = DATABASE()
		AND table_name = 'users'
		AND index_name = 'idx_username_company'
	`

	if err := tx.Raw(checkSQL).Count(&indexExists).Error; err != nil {
		return err
	}

	// 步骤2: 删除旧的唯一索引（如果存在）
	if indexExists > 0 {
		log.Println("  → 删除 (company_id, username) 唯一索引，允许重名用户...")
		dropIndexSQL := `ALTER TABLE users DROP INDEX idx_username_company`
		if err := tx.Exec(dropIndexSQL).Error; err != nil {
			log.Printf("  ⚠️  删除旧索引失败: %v", err)
			return err
		}
		log.Println("  ✓ 删除旧索引成功，现在允许重名用户")
	} else {
		log.Println("  ℹ️  未发现旧的唯一索引")
	}

	// 步骤3: 确保 username 有普通索引（用于查询优化）
	log.Println("  → 确保 username 索引存在...")

	// 检查普通索引是否存在
	var usernameIndexExists int64
	checkUsernameIndexSQL := `
		SELECT COUNT(*)
		FROM information_schema.statistics
		WHERE table_schema = DATABASE()
		AND table_name = 'users'
		AND index_name = 'idx_users_username'
	`

	if err := tx.Raw(checkUsernameIndexSQL).Count(&usernameIndexExists).Error; err != nil {
		return err
	}

	if usernameIndexExists == 0 {
		createIndexSQL := `
			CREATE INDEX idx_users_username ON users(username)
		`
		if err := tx.Exec(createIndexSQL).Error; err != nil {
			log.Printf("  ⚠️  创建 username 索引失败: %v", err)
			// 不返回错误，因为这不是致命错误
		} else {
			log.Println("  ✓ 创建 username 索引成功")
		}
	} else {
		log.Println("  ℹ️  username 索引已存在")
	}

	log.Println("  ✓ 迁移完成：现在允许同一公司内有重名用户")
	log.Println("  💡 用户通过手机号(phone)或ID进行区分")

	return nil
}

// Down 回滚迁移
func (m *Migration20260131AllowDuplicateUsername) Down(tx *gorm.DB) error {
	log.Println("  → 回滚：恢复 (company_id, username) 唯一约束...")

	// 删除普通索引
	tx.Exec(`ALTER TABLE users DROP INDEX idx_users_username`)

	// 恢复复合唯一索引
	createIndexSQL := `
		CREATE UNIQUE INDEX idx_username_company
		ON users(company_id, username)
	`

	if err := tx.Exec(createIndexSQL).Error; err != nil {
		return err
	}

	log.Println("  ✓ 回滚完成")
	return nil
}
