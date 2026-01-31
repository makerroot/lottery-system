package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// InitLogger 初始化日志系统
func InitLogger() {
	Logger = logrus.New()

	// 设置输出为标准输出
	Logger.SetOutput(os.Stdout)

	// 设置JSON格式（生产环境）或文本格式（开发环境）
	if os.Getenv("GIN_MODE") == "release" {
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	}

	// 设置日志级别
	Logger.SetLevel(logrus.InfoLevel)
}

// WithFields 创建带有字段的日志
func WithFields(fields logrus.Fields) *logrus.Entry {
	if Logger == nil {
		InitLogger()
	}
	return Logger.WithFields(fields)
}

// Info 记录信息日志
func Info(args ...interface{}) {
	if Logger == nil {
		InitLogger()
	}
	Logger.Info(args...)
}

// Error 记录错误日志
func Error(args ...interface{}) {
	if Logger == nil {
		InitLogger()
	}
	Logger.Error(args...)
}

// Warn 记录警告日志
func Warn(args ...interface{}) {
	if Logger == nil {
		InitLogger()
	}
	Logger.Warn(args...)
}

// Debug 记录调试日志
func Debug(args ...interface{}) {
	if Logger == nil {
		InitLogger()
	}
	Logger.Debug(args...)
}
