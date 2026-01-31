package main

import (
	"context"
	"log"
	"lottery-system/config"
	"lottery-system/middleware"
	"lottery-system/router"
	"lottery-system/utils"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	// 初始化随机数种子（只调用一次）
	rand.Seed(time.Now().UnixNano())

	// 初始化日志系统
	utils.InitLogger()
	utils.Info("🚀 抽奖系统启动中...")

	// 加载配置
	config.LoadConfig()

	// 初始化数据库
	config.InitDB()

	// 初始化 Redis（如果启用）
	var redisClient *redis.Client
	if config.AppConfig.EnableRedis {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     config.AppConfig.RedisAddr,
			Password: config.AppConfig.RedisPassword,
			DB:       config.AppConfig.RedisDB,
		})

		// 测试 Redis 连接
		_, err := redisClient.Ping(context.Background()).Result()
		if err != nil {
			log.Printf("⚠️  Redis 连接失败: %v，将降级到内存限流", err)
			redisClient = nil
		} else {
			log.Println("✅ Redis 连接成功")
			// 初始化 Redis 限流器（使用配置文件中的参数）
			middleware.InitRedisRateLimiter(redisClient, config.AppConfig.RateLimitRPS, config.AppConfig.RateLimitBurst)
			log.Printf("✅ Redis 限流器已启用（%d req/sec, %d burst）", config.AppConfig.RateLimitRPS, config.AppConfig.RateLimitBurst)
		}
	} else {
		log.Println("ℹ️  Redis 未启用，使用内存限流")
		// 初始化内存限流器（使用配置文件中的参数）
		middleware.InitRateLimiter(config.AppConfig.RateLimitRPS, config.AppConfig.RateLimitBurst)
		log.Printf("✅ 内存限流器已初始化（%d req/sec, %d burst）", config.AppConfig.RateLimitRPS, config.AppConfig.RateLimitBurst)
	}

	// 设置路由（自动应用中间件和限流）
	r := router.SetupRouter()

	// 启动服务器
	r.Run(":" + config.AppConfig.ServerPort)
}
