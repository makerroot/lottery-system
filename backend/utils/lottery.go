package utils

import (
	"fmt"
	"math/rand"

	"lottery-system/models"

	"gorm.io/gorm"
)

// RandomInt 生成0到max-1之间的随机整数
func RandomInt(max int) int {
	return rand.Intn(max)
}

// RandomIndices 生成count个不重复的随机索引（0到max-1）
func RandomIndices(max, count int) []int {
	if count >= max {
		count = max
	}

	// 创建索引池
	indices := make([]int, max)
	for i := 0; i < max; i++ {
		indices[i] = i
	}

	// Fisher-Yates洗牌算法
	for i := max - 1; i > max-count-1; i-- {
		j := rand.Intn(i + 1)
		indices[i], indices[j] = indices[j], indices[i]
	}

	// 返回前count个
	return indices[max-count:]
}

// RandomFloat 生成0到1之间的随机浮点数
func RandomFloat() float64 {
	return rand.Float64()
}

// DrawLotteryWithLevel 执行抽奖逻辑（指定奖项等级）
func DrawLotteryWithLevel(db *gorm.DB, user *models.User, ip string, companyID int, levelID int) (*models.DrawRecord, error) {
	// 检查用户是否已经抽过奖
	if user.HasDrawn {
		return nil, fmt.Errorf("user has already drawn")
	}

	// 开启事务
	tx := db.Begin()

	// 获取指定的奖项等级
	var level models.PrizeLevel
	if err := tx.Where("id = ? AND company_id = ?", levelID, companyID).First(&level).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("prize level not found")
	}

	// 获取该奖项下的所有奖品
	var prizes []models.Prize
	if err := tx.Where("level_id = ?", level.ID).Find(&prizes).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 检查是否有奖品
	if len(prizes) == 0 {
		tx.Rollback()
		return nil, fmt.Errorf("no prizes found for this level")
	}

	// 筛选出有库存的奖品
	var availablePrizes []models.Prize
	for _, prize := range prizes {
		if prize.UsedStock < prize.TotalStock {
			availablePrizes = append(availablePrizes, prize)
		}
	}

	// 检查是否有可用奖品
	if len(availablePrizes) == 0 {
		tx.Rollback()
		return nil, fmt.Errorf("all prizes are out of stock")
	}

	// 随机选择一个有库存的奖品
	selectedPrize := availablePrizes[rand.Intn(len(availablePrizes))]

	// 创建抽奖记录
	record := &models.DrawRecord{
		CompanyID: companyID,
		UserID:    user.ID,
		LevelID:   level.ID,
		PrizeID:   selectedPrize.ID,
		IP:        ip,
	}

	if err := tx.Create(record).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 更新奖品库存
	if err := tx.Model(&selectedPrize).Update("used_stock", selectedPrize.UsedStock+1).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 更新用户已抽奖状态
	if err := tx.Model(user).Update("has_drawn", true).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 加载关联数据
	db.Preload("Level").Preload("Prize").First(record, record.ID)

	return record, nil
}

// 执行抽奖逻辑（完全随机，不使用概率）
func DrawLottery(db *gorm.DB, user *models.User, ip string, companyID int) (*models.DrawRecord, error) {
	// 检查用户是否已经抽过奖
	if user.HasDrawn {
		return nil, fmt.Errorf("user has already drawn")
	}

	// 开启事务
	tx := db.Begin()

	// 获取该公司的所有有库存的奖品（直接查询奖品表）
	var availablePrizes []models.Prize
	if err := tx.Raw(`
		SELECT p.* FROM prizes p
		INNER JOIN prize_levels l ON p.level_id = l.id
		WHERE l.company_id = ? AND l.is_active = ? AND p.used_stock < p.total_stock
		ORDER BY l.sort_order ASC, p.id ASC
	`, companyID, true).Find(&availablePrizes).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 如果没有可用奖品，返回未中奖
	if len(availablePrizes) == 0 {
		tx.Rollback()
		return nil, fmt.Errorf("no prizes available")
	}

	// 完全随机选择一个奖品
	selectedPrize := availablePrizes[rand.Intn(len(availablePrizes))]

	// 获取该奖品的奖项等级
	var level models.PrizeLevel
	if err := tx.First(&level, selectedPrize.LevelID).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("prize level not found")
	}

	// 创建抽奖记录
	record := &models.DrawRecord{
		CompanyID: companyID,
		UserID:    user.ID,
		LevelID:   level.ID,
		PrizeID:   selectedPrize.ID,
		IP:        ip,
	}

	if err := tx.Create(record).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 更新奖品库存
	if err := tx.Model(&selectedPrize).Update("used_stock", selectedPrize.UsedStock+1).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 更新用户已抽奖状态
	if err := tx.Model(user).Update("has_drawn", true).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 加载关联数据
	db.Preload("Level").Preload("Prize").First(record, record.ID)

	return record, nil
}
