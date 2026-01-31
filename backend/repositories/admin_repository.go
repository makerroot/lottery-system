package repositories

import (
	"lottery-system/config"
	"lottery-system/models"
)

// AdminRepository handles admin data operations
type AdminRepository struct{}

// NewAdminRepository creates a new admin repository
func NewAdminRepository() *AdminRepository {
	return &AdminRepository{}
}

// FindByID finds an admin by ID
func (r *AdminRepository) FindByID(id int) (*models.Admin, error) {
	var admin models.Admin
	err := config.DB.First(&admin, id).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// FindByIDWithPreload finds an admin by ID with preloaded associations
func (r *AdminRepository) FindByIDWithPreload(id int) (*models.Admin, error) {
	var admin models.Admin
	err := config.DB.Preload("Company").First(&admin, id).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// FindByUsername finds an admin by username
func (r *AdminRepository) FindByUsername(username string) (*models.Admin, error) {
	var admin models.Admin
	err := config.DB.Where("username = ?", username).Preload("Company").First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// Create creates a new admin
func (r *AdminRepository) Create(admin *models.Admin) error {
	return config.DB.Create(admin).Error
}

// Update updates an admin
func (r *AdminRepository) Update(admin *models.Admin) error {
	return config.DB.Save(admin).Error
}

// Delete deletes an admin by ID
func (r *AdminRepository) Delete(id int) error {
	return config.DB.Delete(&models.Admin{}, id).Error
}

// FindAll finds all admins with optional company filter
func (r *AdminRepository) FindAll(companyID *int) ([]models.Admin, error) {
	var admins []models.Admin
	query := config.DB.Model(&models.Admin{})

	if companyID != nil {
		query = query.Where("company_id = ?", *companyID)
	}

	err := query.Preload("Company").Order("id ASC").Find(&admins).Error
	return admins, err
}

// ExistsByUsername checks if an admin exists by username
func (r *AdminRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := config.DB.Model(&models.Admin{}).
		Where("username = ?", username).
		Count(&count).Error
	return count > 0, err
}

// ExistsByUsernameAndID checks if an admin exists by username excluding a specific ID
func (r *AdminRepository) ExistsByUsernameAndID(username string, excludeID int) (bool, error) {
	var count int64
	err := config.DB.Model(&models.Admin{}).
		Where("username = ? AND id != ?", username, excludeID).
		Count(&count).Error
	return count > 0, err
}
