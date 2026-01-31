package config

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"lottery-system/migrations"
	"lottery-system/models"
	"lottery-system/utils"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// 如果是MySQL，先创建数据库
	if strings.Contains(AppConfig.DatabaseURL, "@tcp(") || strings.Contains(AppConfig.DatabaseURL, "mysql:") {
		fmt.Println("📦 Connecting to MySQL database...")
		if err := createMySQLDatabaseIfNotExists(); err != nil {
			log.Printf("⚠️  Failed to create database: %v", err)
			log.Println("ℹ️  Trying to connect anyway...")
		}
	}

	// 带重试的数据库连接
	db, err := openDatabaseWithRetry()
	if err != nil {
		log.Fatal("Failed to connect to database after retries:", err)
	}

	// 自动迁移
	err = models.AutoMigrate(db)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	DB = db
	fmt.Println("✅ Database connected and migrated successfully")

	// 执行数据库迁移（如果有新的迁移）
	// 注意：这必须在 InitializeData 之前执行
	if err := runMigrations(); err != nil {
		log.Printf("⚠️  迁移执行失败: %v", err)
		// 迁移失败不终止应用，让用户决定是否继续
	}

	// 初始化数据（创建默认管理员等）
	if err := InitializeData(db); err != nil {
		log.Printf("⚠️  Failed to initialize data: %v", err)
	}
}

// createMySQLDatabaseIfNotExists 创建MySQL数据库（如果不存在）
func createMySQLDatabaseIfNotExists() error {
	// 从DSN中提取数据库名、主机、端口等信息
	dsn := AppConfig.DatabaseURL

	// 解析DSN，获取数据库名
	var dbName, dsnWithoutDB string
	if strings.Contains(dsn, "/") {
		parts := strings.Split(dsn, "/")
		dbName = parts[len(parts)-1]
		// 移除参数部分
		if idx := strings.Index(dbName, "?"); idx > 0 {
			dbName = dbName[:idx]
		}
		dsnWithoutDB = strings.Join(parts[:len(parts)-1], "/") + "/?"
	} else {
		return fmt.Errorf("invalid DSN format")
	}

	// 连接到MySQL服务器（不指定数据库）
	db, err := sql.Open("mysql", dsnWithoutDB)
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL server: %w", err)
	}
	defer db.Close()

	// 检查数据库是否存在
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?", dbName).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check database existence: %w", err)
	}

	if count > 0 {
		fmt.Printf("ℹ️  Database '%s' already exists\n", dbName)
		return nil
	}

	// 创建数据库
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName))
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	fmt.Printf("✅ Database '%s' created successfully\n", dbName)
	return nil
}

func openDatabase() (*gorm.DB, error) {
	dsn := AppConfig.DatabaseURL

	// 检测数据库类型
	if strings.Contains(dsn, "@tcp(") || strings.Contains(dsn, "mysql:") {
		// MySQL连接
		fmt.Println("📦 Connecting to MySQL database...")
		return gorm.Open(mysql.Open(dsn), &gorm.Config{})
	} else {
		// SQLite连接（默认）
		fmt.Println("📦 Connecting to SQLite database...")
		return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	}
}

// openDatabaseWithRetry 带重试的数据库连接
func openDatabaseWithRetry() (*gorm.DB, error) {
	const maxRetries = 10
	const retryInterval = 3 // 秒

	var lastErr error

	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			fmt.Printf("⏳ Retry %d/%d in %d seconds...\n", i, maxRetries, retryInterval)
			time.Sleep(time.Duration(retryInterval) * time.Second)
		}

		db, err := openDatabase()
		if err == nil {
			// 测试连接
			sqlDB, err := db.DB()
			if err != nil {
				lastErr = err
				continue
			}

			if err := sqlDB.Ping(); err != nil {
				fmt.Printf("⚠️  Database ping failed: %v\n", err)
				lastErr = err
				continue
			}

			fmt.Printf("✅ Database connection successful (attempt %d)\n", i+1)
			return db, nil
		}

		fmt.Printf("⚠️  Database connection failed: %v\n", err)
		lastErr = err
	}

	return nil, fmt.Errorf("failed to connect after %d retries: %w", maxRetries, lastErr)
}

// OpenDatabase 打开数据库连接（不执行自动迁移）
func OpenDatabase() (*gorm.DB, error) {
	return openDatabase()
}

func GetDB() *gorm.DB {
	return DB
}

// InitializeData 初始化数据（创建默认管理员、公司、奖品等）
func InitializeData(db *gorm.DB) error {
	fmt.Println("🔧 Initializing data...")

	// 检查是否已有数据
	var adminCount int64
	if err := db.Model(&models.Admin{}).Count(&adminCount).Error; err != nil {
		return fmt.Errorf("failed to check admin count: %w", err)
	}

	// 如果已有数据，只检查并创建缺失的奖品
	if adminCount > 0 {
		fmt.Println("ℹ️  Admin data exists, checking for missing prizes...")

		// 检查是否有奖品数据
		var prizeCount int64
		if err := db.Model(&models.Prize{}).Count(&prizeCount).Error; err != nil {
			return fmt.Errorf("failed to check prize count: %w", err)
		}

		// 如果没有奖品，创建默认奖品
		if prizeCount == 0 {
			fmt.Println("⚠️  No prizes found, creating default prizes...")

			// 获取默认公司
			var defaultCompany models.Company
			if err := db.Where("code = ?", "DEFAULT").First(&defaultCompany).Error; err != nil {
				return fmt.Errorf("failed to find default company: %w", err)
			}

			// 获取奖品等级
			var prizeLevels []models.PrizeLevel
			if err := db.Where("company_id = ?", defaultCompany.ID).Find(&prizeLevels).Error; err != nil {
				return fmt.Errorf("failed to find prize levels: %w", err)
			}

			// 为每个等级创建默认奖品（带库存）
			for _, level := range prizeLevels {
				var prizeCount int64
				db.Model(&models.Prize{}).Where("level_id = ?", level.ID).Count(&prizeCount)

				if prizeCount == 0 {
					prize := models.Prize{
						LevelID:    int(level.ID),
						Name:       level.Description,
						TotalStock: 0, // 稍后会在管理后台设置
						UsedStock:  0,
						Image:      "",
					}
					if err := db.Create(&prize).Error; err != nil {
						return fmt.Errorf("failed to create prize for level %s: %w", level.Name, err)
					}
					fmt.Printf("   ✅ Created prize: %s (请在管理后台设置库存)\n", prize.Name)
				}
			}

			fmt.Println("✅ Default prizes created successfully")
		} else {
			fmt.Printf("ℹ️  Prizes already exist (%d found)\n", prizeCount)
		}

		return nil
	}

	// 1. 创建默认管理员（全新安装）
	defaultUsername := "makerroot"
	if AppConfig.DefaultAdminUsername != "" {
		defaultUsername = AppConfig.DefaultAdminUsername
	}

	defaultPassword := "123456"
	if AppConfig.DefaultAdminPassword != "" {
		defaultPassword = AppConfig.DefaultAdminPassword
	}

	hashedPassword, err := utils.HashPassword(defaultPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	admin := models.Admin{
		Username:     defaultUsername,
		Password:     hashedPassword,
		IsSuperAdmin: true,
		Role:         "super_admin",
		CompanyID:    nil,
	}

	if err := db.Create(&admin).Error; err != nil {
		return fmt.Errorf("failed to create default admin: %w", err)
	}

	fmt.Println("✅ Default admin created successfully")
	fmt.Printf("   Username: %s\n", defaultUsername)
	fmt.Printf("   Password: %s\n", defaultPassword)
	fmt.Println("   ⚠️  Please change the password after first login!")

	// 2. 创建默认公司
	defaultCompany := models.Company{
		Name:       "默认公司",
		Code:       "DEFAULT",
		ThemeColor: "#8b5cf6",
		IsActive:   true,
	}

	if err := db.Create(&defaultCompany).Error; err != nil {
		return fmt.Errorf("failed to create default company: %w", err)
	}

	fmt.Println("✅ Default company created successfully")
	fmt.Printf("   Name: %s\n", defaultCompany.Name)
	fmt.Printf("   Code: %s\n", defaultCompany.Code)

	// 3. 创建默认奖品等级（移除概率，库存由奖品管理）
	prizeLevels := []models.PrizeLevel{
		{
			CompanyID:   int(defaultCompany.ID),
			Name:        "一等奖",
			Description: "iPhone 15 Pro",
			TotalStock:  0, // 库存由奖品管理，固定为0
			UsedStock:   0,
			SortOrder:   1,
			IsActive:    true,
		},
		{
			CompanyID:   int(defaultCompany.ID),
			Name:        "二等奖",
			Description: "iPad Pro",
			TotalStock:  0,
			UsedStock:   0,
			SortOrder:   2,
			IsActive:    true,
		},
		{
			CompanyID:   int(defaultCompany.ID),
			Name:        "三等奖",
			Description: "AirPods Pro",
			TotalStock:  0,
			UsedStock:   0,
			SortOrder:   3,
			IsActive:    true,
		},
		{
			CompanyID:   int(defaultCompany.ID),
			Name:        "四等奖",
			Description: "小米充电宝",
			TotalStock:  0,
			UsedStock:   0,
			SortOrder:   4,
			IsActive:    true,
		},
		{
			CompanyID:   int(defaultCompany.ID),
			Name:        "参与奖",
			Description: "定制纪念品",
			TotalStock:  0,
			UsedStock:   0,
			SortOrder:   5,
			IsActive:    true,
		},
	}

	for _, level := range prizeLevels {
		if err := db.Create(&level).Error; err != nil {
			return fmt.Errorf("failed to create prize level %s: %w", level.Name, err)
		}
		fmt.Printf("   ✅ Created: %s\n", level.Name)
	}

	fmt.Println("✅ Default prize levels created successfully")
	fmt.Println("   Total: 5 prize levels")
	fmt.Println("   ℹ️  Stock managed by prizes (add prizes in admin panel)")

	// 4. 为每个等级创建默认奖品（带库存）
	prizes := []models.Prize{
		{LevelID: int(prizeLevels[0].ID), Name: "iPhone 15 Pro 256GB", TotalStock: 3, UsedStock: 0, Image: ""},
		{LevelID: int(prizeLevels[1].ID), Name: "iPad Pro 11英寸 256GB", TotalStock: 10, UsedStock: 0, Image: ""},
		{LevelID: int(prizeLevels[2].ID), Name: "AirPods Pro (第2代)", TotalStock: 30, UsedStock: 0, Image: ""},
		{LevelID: int(prizeLevels[3].ID), Name: "小米极充套装 120W", TotalStock: 100, UsedStock: 0, Image: ""},
		{LevelID: int(prizeLevels[4].ID), Name: "定制U盘 64GB", TotalStock: 500, UsedStock: 0, Image: ""},
	}

	for _, prize := range prizes {
		if err := db.Create(&prize).Error; err != nil {
			return fmt.Errorf("failed to create prize: %w", err)
		}
		fmt.Printf("   ✅ Created: %s (库存: %d)\n", prize.Name, prize.TotalStock)
	}

	fmt.Println("✅ Default prizes created successfully")
	fmt.Println("   Total: 5 prizes")
	fmt.Println("   💡 You can modify stock quantities in the admin panel")

	fmt.Println("\n🎉 All data initialized successfully!")
	fmt.Println("   You can now:")
	fmt.Println("   1. Login with the default admin account")
	fmt.Println("   2. View and manage prize levels in the admin panel")
	fmt.Println("   3. Add prizes and set stock quantities")
	fmt.Println("   4. Start the lottery draw")

	return nil
}

// runMigrations 执行数据库迁移
func runMigrations() error {
	// 注册所有迁移
	// 注意：迁移的顺序很重要，新迁移添加到末尾
	migrations.RegisterMigration(&migrations.Migration20260125ModifyUserUnique{})
	migrations.RegisterMigration(&migrations.Migration20260125AddPrizeStock{})
	migrations.RegisterMigration(&migrations.Migration20260131AllowDuplicateUsername{})

	// 执行迁移
	return migrations.RunMigrations(DB)
}
