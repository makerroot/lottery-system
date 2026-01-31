package repositories

import (
	"lottery-system/config"
	"lottery-system/models"
)

// PrizeRepository handles prize and prize level data operations
type PrizeRepository struct{}

// NewPrizeRepository creates a new prize repository
func NewPrizeRepository() *PrizeRepository {
	return &PrizeRepository{}
}

// FindActiveLevelsByCompany finds all active prize levels for a company
func (r *PrizeRepository) FindActiveLevelsByCompany(companyID int) ([]models.PrizeLevel, error) {
	var levels []models.PrizeLevel
	err := config.DB.Where("company_id = ? AND is_active = ?", companyID, true).
		Order("sort_order ASC").
		Find(&levels).Error
	return levels, err
}

// FindAllLevelsByCompany finds all prize levels for a company
func (r *PrizeRepository) FindAllLevelsByCompany(companyID int) ([]models.PrizeLevel, error) {
	var levels []models.PrizeLevel
	err := config.DB.Where("company_id = ?", companyID).
		Order("sort_order ASC").
		Find(&levels).Error
	return levels, err
}

// FindLevelByID finds a prize level by ID and company ID
func (r *PrizeRepository) FindLevelByID(id, companyID int) (*models.PrizeLevel, error) {
	var level models.PrizeLevel
	err := config.DB.Where("id = ? AND company_id = ?", id, companyID).First(&level).Error
	if err != nil {
		return nil, err
	}
	return &level, nil
}

// FindActiveLevelByID finds an active prize level by ID and company ID
func (r *PrizeRepository) FindActiveLevelByID(id, companyID int) (*models.PrizeLevel, error) {
	var level models.PrizeLevel
	err := config.DB.Where("id = ? AND company_id = ? AND is_active = ?", id, companyID, true).
		First(&level).Error
	if err != nil {
		return nil, err
	}
	return &level, nil
}

// CreateLevel creates a new prize level
func (r *PrizeRepository) CreateLevel(level *models.PrizeLevel) error {
	return config.DB.Create(level).Error
}

// UpdateLevel updates a prize level
func (r *PrizeRepository) UpdateLevel(level *models.PrizeLevel) error {
	return config.DB.Save(level).Error
}

// DeleteLevel deletes a prize level by ID
func (r *PrizeRepository) DeleteLevel(id int) error {
	return config.DB.Delete(&models.PrizeLevel{}, id).Error
}

// UpdateStock updates the used stock of a prize level
func (r *PrizeRepository) UpdateStock(levelID int, usedStock int) error {
	return config.DB.Model(&models.PrizeLevel{}).
		Where("id = ?", levelID).
		Update("used_stock", usedStock).Error
}

// FindPrizesByLevel finds all prizes for a prize level
func (r *PrizeRepository) FindPrizesByLevel(levelID int) ([]models.Prize, error) {
	var prizes []models.Prize
	err := config.DB.Where("level_id = ?", levelID).Find(&prizes).Error
	return prizes, err
}

// FindAnyPrizeByLevel finds any one prize for a prize level
func (r *PrizeRepository) FindAnyPrizeByLevel(levelID int) (*models.Prize, error) {
	var prize models.Prize
	err := config.DB.Where("level_id = ?", levelID).First(&prize).Error
	if err != nil {
		return nil, err
	}
	return &prize, nil
}

// CreatePrize creates a new prize
func (r *PrizeRepository) CreatePrize(prize *models.Prize) error {
	return config.DB.Create(prize).Error
}

// UpdatePrize updates a prize
func (r *PrizeRepository) UpdatePrize(prize *models.Prize) error {
	return config.DB.Save(prize).Error
}

// DeletePrize deletes a prize by ID
func (r *PrizeRepository) DeletePrize(id int) error {
	return config.DB.Delete(&models.Prize{}, id).Error
}

// FindPrizeByID finds a prize by ID
func (r *PrizeRepository) FindPrizeByID(id int) (*models.Prize, error) {
	var prize models.Prize
	err := config.DB.First(&prize, id).Error
	if err != nil {
		return nil, err
	}
	return &prize, nil
}
