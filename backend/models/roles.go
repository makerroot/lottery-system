package models

// 角色常量定义
const (
	// 用户角色
	RoleUser = "user" // 普通用户

	// 管理员角色
	RoleAdmin      = "admin"       // 普通管理员
	RoleSuperAdmin = "super_admin" // 超级管理员
)

// RoleIsValid 检查角色是否有效
func RoleIsValid(role string) bool {
	switch role {
	case RoleUser, RoleAdmin, RoleSuperAdmin:
		return true
	default:
		return false
	}
}

// RoleIsAdmin 检查是否是管理员角色
func RoleIsAdmin(role string) bool {
	return role == RoleAdmin || role == RoleSuperAdmin
}

// RoleIsSuperAdmin 检查是否是超级管理员角色
func RoleIsSuperAdmin(role string) bool {
	return role == RoleSuperAdmin
}

// RoleIsUser 检查是否是用户角色
func RoleIsUser(role string) bool {
	return role == RoleUser
}
