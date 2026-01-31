package models

import (
	"time"

	"gorm.io/gorm"
)

// Company 公司模型
type Company struct {
	ID         int    `gorm:"type:integer;primarykey" json:"id"`
	Code       string `gorm:"type:varchar(50);unique;not null" json:"code"`          // 公司代码，用于URL识别
	Name       string `gorm:"type:varchar(100);not null" json:"name"`                // 公司名称
	Logo       string `gorm:"type:varchar(255)" json:"logo"`                         // Logo URL
	ThemeColor string `gorm:"type:varchar(20);default:'#00fff5'" json:"theme_color"` // 主题颜色
	BgColor    string `gorm:"type:varchar(20);default:'#0a0f14'" json:"bg_color"`      // 背景颜色

	// 文案配置
	Title          string `gorm:"type:varchar(100);default:'幸运抽奖'" json:"title"`           // 系统标题
	Subtitle       string `gorm:"type:varchar(200)" json:"subtitle"`                       // 副标题
	WelcomeText    string `gorm:"type:text" json:"welcome_text"`                           // 欢迎语
	RulesText      string `gorm:"type:text" json:"rules_text"`                             // 规则说明
	DrawButtonText string `gorm:"type:varchar(50);default:'点击抽奖'" json:"draw_button_text"` // 抽奖按钮文字
	SuccessText    string `gorm:"type:varchar(200);default:'恭喜中奖！'" json:"success_text"`   // 中奖恭喜语

	// 联系信息
	ContactName  string `gorm:"type:varchar(100)" json:"contact_name"`  // 联系人
	ContactPhone string `gorm:"type:varchar(20)" json:"contact_phone"`  // 联系电话
	ContactEmail string `gorm:"type:varchar(100)" json:"contact_email"` // 联系邮箱

	IsActive  bool      `gorm:"default:true" json:"is_active"` // 是否启用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AutoMigrate 自动迁移数据库表（包含Company）
func AutoMigrateWithCompany(db *gorm.DB) error {
	return db.AutoMigrate(
		&Company{},
		&Admin{},
		&User{},
		&PrizeLevel{},
		&Prize{},
		&DrawRecord{},
	)
}
