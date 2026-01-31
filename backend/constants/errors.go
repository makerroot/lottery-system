// Package constants defines all constant values used throughout the application.
package constants

// Error message constants
const (
	// Authentication errors
	ErrInvalidCredentials = "用户名或密码错误"
	ErrUserNotFound       = "用户不存在"
	ErrAdminNotFound      = "管理员不存在"
	ErrCompanyNotFound    = "公司不存在"
	ErrInvalidToken       = "无效的令牌"
	ErrTokenExpired       = "令牌已过期"
	ErrUnauthorized       = "未授权访问"
	ErrInvalidInput       = "输入数据无效"

	// Validation errors
	ErrInvalidUsername        = "用户名格式错误"
	ErrInvalidPassword        = "密码格式错误"
	ErrInvalidPasswordLength  = "密码长度至少6位"
	ErrInvalidPasswordLength8 = "密码长度至少8位"
	ErrInvalidPhone           = "手机号格式错误"
	ErrInvalidName            = "姓名格式错误"
	ErrInvalidEmail           = "邮箱格式错误"
	ErrRequiredField          = "必填字段不能为空"

	// User errors
	ErrUsernameExists   = "该用户名已存在"
	ErrPhoneExists      = "该手机号已存在"
	ErrUserAlreadyDrawn = "用户已经抽过奖"
	ErrNoUsersAvailable = "没有可抽奖的用户"
	ErrCannotDeleteSelf = "不能删除自己"

	// Admin errors
	ErrAdminExists          = "该管理员已存在"
	ErrAdminMustHaveCompany = "普通管理员必须指定所属公司"
	ErrSuperAdminNoCompany  = "超级管理员不应该指定所属公司"
	ErrPermissionDenied     = "权限不足"
	ErrCannotModifyCompany  = "不能修改公司关联"

	// Company errors
	ErrCompanyCodeExists  = "公司代码已存在"
	ErrCompanyInactive    = "公司未激活"
	ErrInvalidCompanyCode = "无效的公司代码"

	// Prize errors
	ErrPrizeNotFound      = "奖品不存在"
	ErrPrizeLevelNotFound = "奖项等级不存在"
	ErrPrizeOutOfStock    = "该奖品已抽完"
	ErrNoPrizesAvailable  = "没有可用的奖品"
	ErrInvalidDrawCount   = "抽奖数量无效"

	// Request errors
	ErrInvalidRequestFormat = "请求参数格式错误"
	ErrInvalidJSON          = "JSON格式错误"
	ErrMissingParameter     = "缺少必要参数"

	// Operation errors
	ErrOperationFailed    = "操作失败，请稍后重试"
	ErrDatabaseError      = "数据库错误"
	ErrInternalError      = "内部错误"
	ErrServiceUnavailable = "服务暂时不可用"
)

// Error codes for API responses
const (
	CodeBadRequest         = "BAD_REQUEST"
	CodeUnauthorized       = "UNAUTHORIZED"
	CodeForbidden          = "FORBIDDEN"
	CodeNotFound           = "NOT_FOUND"
	CodeConflict           = "CONFLICT"
	CodeInternalError      = "INTERNAL_ERROR"
	CodeInvalidCredentials = "INVALID_CREDENTIALS"
	CodeInvalidToken       = "INVALID_TOKEN"
	CodeUserNotFound       = "USER_NOT_FOUND"
	CodeAdminNotFound      = "ADMIN_NOT_FOUND"
	CodeCompanyNotFound    = "COMPANY_NOT_FOUND"
	CodeUsernameExists     = "USERNAME_EXISTS"
	CodePermissionDenied   = "PERMISSION_DENIED"
	CodePrizeOutOfStock    = "PRIZE_OUT_OF_STOCK"
	CodeInvalidPassword    = "INVALID_PASSWORD"
	CodeInvalidInput       = "INVALID_INPUT"
)
