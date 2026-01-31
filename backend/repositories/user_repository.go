// Package repositories provides data access layer for the lottery system.
//
// It abstracts database operations and provides a clean interface for services
// to interact with the database.
package repositories

import (
	"lottery-system/config"
	"lottery-system/models"
)

// UserRepository handles user data operations
type UserRepository struct{}

// NewUserRepository creates a new user repository
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id int) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByIDWithPreload finds a user by ID with preloaded associations
func (r *UserRepository) FindByIDWithPreload(id int) (*models.User, error) {
	var user models.User
	err := config.DB.Preload("Company").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername finds a user by username and company ID
func (r *UserRepository) FindByUsername(username string, companyID int) (*models.User, error) {
	var user models.User
	err := config.DB.Where("username = ? AND company_id = ?", username, companyID).
		Preload("Company").
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByPhone finds a user by phone number and company ID
func (r *UserRepository) FindByPhone(phone string, companyID int) (*models.User, error) {
	var user models.User
	err := config.DB.Where("phone = ? AND company_id = ?", phone, companyID).
		Preload("Company").
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	return config.DB.Create(user).Error
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	return config.DB.Save(user).Error
}

// UpdateFields updates specific fields of a user
func (r *UserRepository) UpdateFields(id int, fields map[string]interface{}) error {
	return config.DB.Model(&models.User{}).Where("id = ?", id).Updates(fields).Error
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(id int) error {
	return config.DB.Delete(&models.User{}, id).Error
}

// FindAll finds all users with optional filters
func (r *UserRepository) FindAll(filters map[string]interface{}) ([]models.User, error) {
	var users []models.User
	query := config.DB.Model(&models.User{})

	// Apply filters
	if companyID, ok := filters["company_id"]; ok {
		query = query.Where("company_id = ?", companyID)
	}
	if hasDrawn, ok := filters["has_drawn"]; ok {
		query = query.Where("has_drawn = ?", hasDrawn)
	}

	err := query.Order("id ASC").Find(&users).Error
	return users, err
}

// FindAllWithPreload finds all users with preloaded associations
func (r *UserRepository) FindAllWithPreload(filters map[string]interface{}) ([]models.User, error) {
	var users []models.User
	query := config.DB.Model(&models.User{})

	// Apply filters
	if companyID, ok := filters["company_id"]; ok {
		query = query.Where("company_id = ?", companyID)
	}
	if hasDrawn, ok := filters["has_drawn"]; ok {
		query = query.Where("has_drawn = ?", hasDrawn)
	}

	err := query.Preload("Company").Order("id ASC").Find(&users).Error
	return users, err
}

// CountByCompany counts users by company ID
func (r *UserRepository) CountByCompany(companyID int) (int64, error) {
	var count int64
	err := config.DB.Model(&models.User{}).Where("company_id = ?", companyID).Count(&count).Error
	return count, err
}

// CountByCompanyAndStatus counts users by company ID and drawn status
func (r *UserRepository) CountByCompanyAndStatus(companyID int, hasDrawn bool) (int64, error) {
	var count int64
	err := config.DB.Model(&models.User{}).
		Where("company_id = ? AND has_drawn = ?", companyID, hasDrawn).
		Count(&count).Error
	return count, err
}

// ExistsByUsername checks if a user exists by username in a company
func (r *UserRepository) ExistsByUsername(username string, companyID int) (bool, error) {
	var count int64
	err := config.DB.Model(&models.User{}).
		Where("username = ? AND company_id = ?", username, companyID).
		Count(&count).Error
	return count > 0, err
}

// ExistsByPhone checks if a user exists by phone number in a company
func (r *UserRepository) ExistsByPhone(phone string, companyID int) (bool, error) {
	var count int64
	err := config.DB.Model(&models.User{}).
		Where("phone = ? AND company_id = ?", phone, companyID).
		Count(&count).Error
	return count > 0, err
}

// FindAvailableUsers finds users who haven't drawn yet in a company
func (r *UserRepository) FindAvailableUsers(companyID int) ([]models.User, error) {
	var users []models.User
	err := config.DB.Where("company_id = ? AND has_drawn = ?", companyID, false).
		Order("id ASC").
		Find(&users).Error
	return users, err
}
