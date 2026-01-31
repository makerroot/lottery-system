package models

import (
	"time"
)

// OperationLog 操作日志模型
type OperationLog struct {
	ID         uint      `gorm:"type:integer;primarykey" json:"id"`
	AdminID    uint      `gorm:"type:integer;not null;index" json:"admin_id"`    // 操作管理员ID
	AdminName  string    `gorm:"type:varchar(100)" json:"admin_name"`            // 操作管理员名称（冗余字段，便于查询）
	CompanyID  *uint     `gorm:"type:integer;index" json:"company_id,omitempty"` // 公司ID（null表示超级管理员操作）
	Action     string    `gorm:"type:varchar(50);not null" json:"action"`        // 操作类型：create, update, delete, login, logout等
	Resource   string    `gorm:"type:varchar(100)" json:"resource"`              // 操作资源：company, admin, user, prize_level等
	ResourceID *uint     `gorm:"type:integer" json:"resource_id,omitempty"`      // 资源ID
	Details    string    `gorm:"type:text" json:"details"`                       // 操作详情（JSON格式）
	IPAddress  string    `gorm:"type:varchar(50)" json:"ip_address"`             // IP地址
	UserAgent  string    `gorm:"type:varchar(500)" json:"user_agent"`            // User-Agent
	CreatedAt  time.Time `json:"created_at"`
}

// TableName 指定表名
func (OperationLog) TableName() string {
	return "operation_logs"
}
