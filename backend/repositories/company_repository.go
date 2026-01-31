package repositories

import (
	"lottery-system/config"
	"lottery-system/models"
)

// CompanyRepository handles company data operations
type CompanyRepository struct{}

// NewCompanyRepository creates a new company repository
func NewCompanyRepository() *CompanyRepository {
	return &CompanyRepository{}
}

// FindByID finds a company by ID
func (r *CompanyRepository) FindByID(id int) (*models.Company, error) {
	var company models.Company
	err := config.DB.First(&company, id).Error
	if err != nil {
		return nil, err
	}
	return &company, nil
}

// FindByCode finds a company by code
func (r *CompanyRepository) FindByCode(code string) (*models.Company, error) {
	var company models.Company
	err := config.DB.Where("code = ?", code).First(&company).Error
	if err != nil {
		return nil, err
	}
	return &company, nil
}

// FindByCodeWithActive checks if a company exists and is active
func (r *CompanyRepository) FindByCodeWithActive(code string) (*models.Company, error) {
	var company models.Company
	err := config.DB.Where("code = ? AND is_active = ?", code, true).First(&company).Error
	if err != nil {
		return nil, err
	}
	return &company, nil
}

// Create creates a new company
func (r *CompanyRepository) Create(company *models.Company) error {
	return config.DB.Create(company).Error
}

// Update updates a company
func (r *CompanyRepository) Update(company *models.Company) error {
	return config.DB.Save(company).Error
}

// Delete deletes a company by ID
func (r *CompanyRepository) Delete(id int) error {
	return config.DB.Delete(&models.Company{}, id).Error
}

// FindAll finds all companies
func (r *CompanyRepository) FindAll() ([]models.Company, error) {
	var companies []models.Company
	err := config.DB.Order("id ASC").Find(&companies).Error
	return companies, err
}

// ExistsByCode checks if a company exists by code
func (r *CompanyRepository) ExistsByCode(code string) (bool, error) {
	var count int64
	err := config.DB.Model(&models.Company{}).
		Where("code = ?", code).
		Count(&count).Error
	return count > 0, err
}

// ExistsByCodeAndID checks if a company exists by code excluding a specific ID
func (r *CompanyRepository) ExistsByCodeAndID(code string, excludeID int) (bool, error) {
	var count int64
	err := config.DB.Model(&models.Company{}).
		Where("code = ? AND id != ?", code, excludeID).
		Count(&count).Error
	return count > 0, err
}

// Count counts total companies
func (r *CompanyRepository) Count() (int64, error) {
	var count int64
	err := config.DB.Model(&models.Company{}).Count(&count).Error
	return count, err
}

// CountActive counts active companies
func (r *CompanyRepository) CountActive() (int64, error) {
	var count int64
	err := config.DB.Model(&models.Company{}).Where("is_active = ?", true).Count(&count).Error
	return count, err
}
