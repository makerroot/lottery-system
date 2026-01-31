package main

import (
	"fmt"
	"log"

	"lottery-system/config"
	"lottery-system/models"
)

func main() {
	// 初始化数据库
	config.InitDB()
	defer func() {
		sqlDB, _ := config.DB.DB()
		sqlDB.Close()
	}()

	// 获取默认公司
	var company models.Company
	if err := config.DB.Where("code = ?", "DEFAULT").First(&company).Error; err != nil {
		log.Fatal("默认公司不存在，请先创建公司")
	}

	// 检查是否已有奖项等级
	var existingCount int64
	config.DB.Model(&models.PrizeLevel{}).Where("company_id = ?", company.ID).Count(&existingCount)

	if existingCount > 0 {
		log.Printf("⚠️  已存在 %d 个奖项等级，跳过初始化", existingCount)
		log.Println("💡 如需重新初始化，请先在数据库中删除现有奖项等级")
		return
	}

	log.Println("📋 开始初始化奖项等级...")

	// 创建默认奖项等级
	prizeLevels := []models.PrizeLevel{
		{
			CompanyID:   int(company.ID),
			Name:        "一等奖",
			Description: "iPhone 15 Pro",
			Probability: 0.01, // 1%
			TotalStock:  3,
			UsedStock:   0,
			SortOrder:   1,
			IsActive:    true,
		},
		{
			CompanyID:   int(company.ID),
			Name:        "二等奖",
			Description: "iPad Air",
			Probability: 0.05, // 5%
			TotalStock:  10,
			UsedStock:   0,
			SortOrder:   2,
			IsActive:    true,
		},
		{
			CompanyID:   int(company.ID),
			Name:        "三等奖",
			Description: "AirPods Pro",
			Probability: 0.15, // 15%
			TotalStock:  20,
			UsedStock:   0,
			SortOrder:   3,
			IsActive:    true,
		},
		{
			CompanyID:   int(company.ID),
			Name:        "幸运奖",
			Description: "精美礼品",
			Probability: 0.30, // 30%
			TotalStock:  50,
			UsedStock:   0,
			SortOrder:   4,
			IsActive:    true,
		},
		{
			CompanyID:   int(company.ID),
			Name:        "参与奖",
			Description: "纪念品",
			Probability: 0.49, // 49%
			TotalStock:  100,
			UsedStock:   0,
			SortOrder:   5,
			IsActive:    true,
		},
	}

	// 批量创建（带重复检查）
	createdCount := 0
	skippedCount := 0

	for i := range prizeLevels {
		// 检查是否已存在同名奖项
		var existingLevel models.PrizeLevel
		err := config.DB.Where(
			"company_id = ? AND name = ?",
			company.ID,
			prizeLevels[i].Name,
		).First(&existingLevel).Error

		if err == nil {
			// 奖项已存在，跳过
			log.Printf("⏭️  跳过已存在的奖项: %s", prizeLevels[i].Name)
			skippedCount++
			continue
		}

		// 创建新奖项
		if err := config.DB.Create(&prizeLevels[i]).Error; err != nil {
			log.Printf("❌ 创建奖项失败: %s - %v", prizeLevels[i].Name, err)
		} else {
			log.Printf("✅ 创建奖项成功: %s", prizeLevels[i].Name)
			createdCount++
		}
	}

	fmt.Println("\n===========================================")
	if createdCount > 0 {
		fmt.Println("✅ 奖项等级初始化完成！")
	} else {
		fmt.Println("ℹ️  所有奖项等级已存在，无需创建")
	}
	fmt.Println("===========================================")
	fmt.Printf("公司: %s (%s)\n", company.Name, company.Code)
	fmt.Printf("创建: %d 个\n", createdCount)
	fmt.Printf("跳过: %d 个\n", skippedCount)
	fmt.Println("\n奖项等级列表:")
	for _, level := range prizeLevels {
		fmt.Printf("  - %s: %s (概率: %.1f%%, 库存: %d)\n",
			level.Name, level.Description, level.Probability*100, level.TotalStock)
	}
	fmt.Println("\n总概率: 100%")
	fmt.Println("===========================================")
}
