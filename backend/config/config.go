package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret      string
	DatabaseURL    string
	ServerPort     string
	JWTExpiration  int64  // hours
	AllowedOrigins string // 逗号分隔的允许来源

	// Redis 配置
	RedisAddr     string // Redis 地址，如 localhost:6379
	RedisPassword string // Redis 密码
	RedisDB       int    // Redis 数据库编号
	EnableRedis   bool   // 是否启用 Redis（用于限流）

	// 限流配置
	RateLimitRPS int // 每秒请求数（requests per second）
	RateLimitBurst int // 突发请求数（burst）

	// 默认管理员配置
	DefaultAdminUsername string // 默认管理员用户名
	DefaultAdminPassword string // 默认管理员密码
}

var AppConfig *Config

func LoadConfig() {
	// 尝试从多个位置加载.env文件
	envFiles := []string{
		".env",                     // 当前目录
		"backend/.env",             // backend子目录
		"/etc/lottery-system/.env", // 系统配置目录
	}

	loaded := false
	for _, envFile := range envFiles {
		if err := godotenv.Load(envFile); err == nil {
			log.Printf("✅ Loaded environment variables from: %s", envFile)
			loaded = true
			break // 找到一个就停止
		}
	}

	if !loaded {
		log.Println("⚠️  No .env file found, using system environment variables or defaults")
	}

	// 强制要求JWT_SECRET在非开发环境必须设置
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" || jwtSecret == "your-secret-key-change-in-production" || jwtSecret == "dev-secret-key-for-testing" {
		log.Fatal("❌ JWT_SECRET environment variable is required and must be set to a strong random value!")
	}

	AppConfig = &Config{
		JWTSecret:      jwtSecret,
		DatabaseURL:    getEnv("DATABASE_URL", "lottery.db"),
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		JWTExpiration:  24 * 7, // 7 days
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:3001"),
		RedisAddr:      getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
		RedisDB:        getEnvInt("REDIS_DB", 0),
		EnableRedis:    getEnvBool("ENABLE_REDIS", false),
		// 限流配置
		RateLimitRPS:   getEnvInt("RATE_LIMIT_RPS", 10),  // 默认每秒10个请求
		RateLimitBurst: getEnvInt("RATE_LIMIT_BURST", 20), // 默认突发20个请求
		// 默认管理员配置
		DefaultAdminUsername: getEnv("DEFAULT_ADMIN_USERNAME", "makerroot"),
		DefaultAdminPassword: getEnv("DEFAULT_ADMIN_PASSWORD", "123456"),
	}

	log.Println("✅ Configuration loaded successfully")
	log.Printf("📋 Server will listen on port: %s", AppConfig.ServerPort)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	// 简单的字符串转 int
	var value int
	fmt.Sscanf(valueStr, "%d", &value)
	return value
}

func getEnvBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	return valueStr == "true" || valueStr == "1"
}
