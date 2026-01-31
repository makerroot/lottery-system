package utils

import (
	"time"
)

// SecurityLogger handles security-related logging
type SecurityLogger struct{}

// NewSecurityLogger creates a new security logger
func NewSecurityLogger() *SecurityLogger {
	return &SecurityLogger{}
}

// LogFailedLogin logs a failed login attempt
func (l *SecurityLogger) LogFailedLogin(username, ip, userAgent string) {
	WithFields(map[string]interface{}{
		"event":      "failed_login",
		"username":   username,
		"ip":         ip,
		"user_agent": userAgent,
		"timestamp":  time.Now().Unix(),
	}).Warn("Failed login attempt")
}

// LogSuccessfulLogin logs a successful login
func (l *SecurityLogger) LogSuccessfulLogin(userID int, username, ip string) {
	WithFields(map[string]interface{}{
		"event":     "successful_login",
		"user_id":   userID,
		"username":  username,
		"ip":        ip,
		"timestamp": time.Now().Unix(),
	}).Info("User logged in successfully")
}

// LogPasswordChange logs a password change
func (l *SecurityLogger) LogPasswordChange(userID int, username string) {
	WithFields(map[string]interface{}{
		"event":     "password_change",
		"user_id":   userID,
		"username":  username,
		"timestamp": time.Now().Unix(),
	}).Info("Password changed successfully")
}

// LogFailedPasswordChange logs a failed password change attempt
func (l *SecurityLogger) LogFailedPasswordChange(userID int, username, reason string) {
	WithFields(map[string]interface{}{
		"event":     "failed_password_change",
		"user_id":   userID,
		"username":  username,
		"reason":    reason,
		"timestamp": time.Now().Unix(),
	}).Warn("Failed password change attempt")
}

// LogAdminCreation logs admin creation
func (l *SecurityLogger) LogAdminCreation(createdBy int, adminUsername string, isSuperAdmin bool) {
	WithFields(map[string]interface{}{
		"event":          "admin_created",
		"created_by":     createdBy,
		"admin_username": adminUsername,
		"is_super_admin": isSuperAdmin,
		"timestamp":      time.Now().Unix(),
	}).Info("New admin created")
}

// LogSuspiciousActivity logs suspicious activity
func (l *SecurityLogger) LogSuspiciousActivity(activity, description, ip string) {
	WithFields(map[string]interface{}{
		"event":       "suspicious_activity",
		"activity":    activity,
		"description": description,
		"ip":          ip,
		"timestamp":   time.Now().Unix(),
	}).Warn("Suspicious activity detected")
}

// LogRateLimitExceeded logs when rate limit is exceeded
func (l *SecurityLogger) LogRateLimitExceeded(operation, key string) {
	WithFields(map[string]interface{}{
		"event":     "rate_limit_exceeded",
		"operation": operation,
		"key":       key,
		"timestamp": time.Now().Unix(),
	}).Warn("Rate limit exceeded")
}

// LogXSSAttempt logs potential XSS attack
func (l *SecurityLogger) LogXSSAttempt(input, ip string) {
	WithFields(map[string]interface{}{
		"event":     "xss_attempt",
		"input":     input[:100], // Truncate for safety
		"ip":        ip,
		"timestamp": time.Now().Unix(),
	}).Warn("Potential XSS attempt detected")
}

// LogSQLInjectionAttempt logs potential SQL injection
func (l *SecurityLogger) LogSQLInjectionAttempt(input, ip string) {
	WithFields(map[string]interface{}{
		"event":     "sql_injection_attempt",
		"input":     input[:100], // Truncate for safety
		"ip":        ip,
		"timestamp": time.Now().Unix(),
	}).Warn("Potential SQL injection attempt detected")
}

// LogSecurityEvent logs a general security event
func (l *SecurityLogger) LogSecurityEvent(event, message string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["event"] = event
	fields["timestamp"] = time.Now().Unix()

	WithFields(fields).Info(message)
}

// Security audit log structure
type SecurityAuditLog struct {
	Event     string                 `json:"event"`
	UserID    *int                   `json:"user_id,omitempty"`
	Username  string                 `json:"username,omitempty"`
	IP        string                 `json:"ip"`
	UserAgent string                 `json:"user_agent,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Timestamp int64                  `json:"timestamp"`
	Success   bool                   `json:"success"`
}
