package repositories

import (
	"lottery-system/config"
	"lottery-system/models"
)

// DrawRepository handles draw record data operations
type DrawRepository struct{}

// NewDrawRepository creates a new draw repository
func NewDrawRepository() *DrawRepository {
	return &DrawRepository{}
}

// Create creates a new draw record
func (r *DrawRepository) Create(record *models.DrawRecord) error {
	return config.DB.Create(record).Error
}

// FindByID finds a draw record by ID
func (r *DrawRepository) FindByID(id int) (*models.DrawRecord, error) {
	var record models.DrawRecord
	err := config.DB.First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// FindByIDWithPreload finds a draw record by ID with preloaded associations
func (r *DrawRepository) FindByIDWithPreload(id int) (*models.DrawRecord, error) {
	var record models.DrawRecord
	err := config.DB.Preload("User").
		Preload("Level").
		Preload("Prize").
		Preload("Company").
		First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// FindByUser finds all draw records for a user
func (r *DrawRepository) FindByUser(userID int) ([]models.DrawRecord, error) {
	var records []models.DrawRecord
	err := config.DB.Where("user_id = ?", userID).
		Preload("Level").
		Preload("Prize").
		Order("created_at DESC").
		Find(&records).Error
	return records, err
}

// FindByCompany finds draw records for a company with optional limit
func (r *DrawRepository) FindByCompany(companyID int, limit int) ([]models.DrawRecord, error) {
	var records []models.DrawRecord
	query := config.DB.Where("company_id = ?", companyID).
		Preload("User").
		Preload("Level").
		Preload("Prize").
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&records).Error
	return records, err
}

// FindByCompanyAndUser finds a draw record by company and user
func (r *DrawRepository) FindByCompanyAndUser(companyID, userID int) (*models.DrawRecord, error) {
	var record models.DrawRecord
	err := config.DB.Where("company_id = ? AND user_id = ?", companyID, userID).
		Preload("Level").
		Preload("Prize").
		First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// CountByCompany counts draw records by company
func (r *DrawRepository) CountByCompany(companyID int) (int64, error) {
	var count int64
	err := config.DB.Model(&models.DrawRecord{}).
		Where("company_id = ?", companyID).
		Count(&count).Error
	return count, err
}

// CountByLevel counts draw records by prize level
func (r *DrawRepository) CountByLevel(levelID int) (int64, error) {
	var count int64
	err := config.DB.Model(&models.DrawRecord{}).
		Where("level_id = ?", levelID).
		Count(&count).Error
	return count, err
}

// Delete deletes a draw record by ID
func (r *DrawRepository) Delete(id int) error {
	return config.DB.Delete(&models.DrawRecord{}, id).Error
}
