// Package services provides business logic layer for the lottery system.
//
// It handles all business operations and orchestrates interactions between
// handlers and repositories.
package services

import (
	"lottery-system/constants"
	"lottery-system/models"
	"lottery-system/repositories"
	"lottery-system/utils"
	"lottery-system/validators"
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo     *repositories.UserRepository
	adminRepo    *repositories.AdminRepository
	jwtSecret    string
	jwtExpiration int64
}

// NewAuthService creates a new auth service
func NewAuthService(jwtSecret string, jwtExpiration int64) *AuthService {
	return &AuthService{
		userRepo:     repositories.NewUserRepository(),
		adminRepo:    repositories.NewAdminRepository(),
		jwtSecret:    jwtSecret,
		jwtExpiration: jwtExpiration,
	}
}

// AuthenticateUser authenticates a user with username and password
func (s *AuthService) AuthenticateUser(username, password string, companyID int) (*models.User, error) {
	// Validate input
	if err := validators.ValidateUsername(username); err != nil {
		return nil, err
	}
	if err := validators.ValidatePassword(password); err != nil {
		return nil, err
	}

	// Find user
	user, err := s.userRepo.FindByUsername(username, companyID)
	if err != nil {
		return nil, utils.NewNotFoundError("用户")
	}

	// Check password
	if !utils.CheckPassword(password, user.Password) {
		return nil, utils.NewAuthenticationError(constants.ErrInvalidCredentials)
	}

	// Check role
	if user.Role != models.RoleUser {
		return nil, utils.NewAuthenticationError(constants.ErrInvalidCredentials)
	}

	return user, nil
}

// AuthenticateAdmin authenticates an admin with username and password
func (s *AuthService) AuthenticateAdmin(username, password string) (*models.Admin, error) {
	// Validate input
	if err := validators.ValidateUsername(username); err != nil {
		return nil, err
	}
	if err := validators.ValidatePassword(password); err != nil {
		return nil, err
	}

	// Find admin
	admin, err := s.adminRepo.FindByUsername(username)
	if err != nil {
		return nil, utils.NewNotFoundError("管理员")
	}

	// Check password
	if !utils.CheckPassword(password, admin.Password) {
		return nil, utils.NewAuthenticationError(constants.ErrInvalidCredentials)
	}

	return admin, nil
}

// GenerateUserToken generates a JWT token for a user
func (s *AuthService) GenerateUserToken(user *models.User) (string, error) {
	return utils.GenerateUserToken(user.ID, user.Username, s.jwtSecret, s.jwtExpiration)
}

// GenerateAdminToken generates a JWT token for an admin
func (s *AuthService) GenerateAdminToken(admin *models.Admin) (string, error) {
	return utils.GenerateToken(admin.ID, admin.Username, s.jwtSecret, s.jwtExpiration)
}

// ValidateUserPassword validates a user's password
func (s *AuthService) ValidateUserPassword(user *models.User, password string) error {
	if !utils.CheckPassword(password, user.Password) {
		return utils.NewAuthenticationError(constants.ErrInvalidCredentials)
	}
	return nil
}

// ValidateAdminPassword validates an admin's password
func (s *AuthService) ValidateAdminPassword(admin *models.Admin, password string) error {
	if !utils.CheckPassword(password, admin.Password) {
		return utils.NewAuthenticationError(constants.ErrInvalidCredentials)
	}
	return nil
}

// HashPassword hashes a plaintext password
func (s *AuthService) HashPassword(password string) (string, error) {
	return utils.HashPassword(password)
}

// ChangeUserPassword changes a user's password
func (s *AuthService) ChangeUserPassword(user *models.User, oldPassword, newPassword string) error {
	// Validate old password
	if err := s.ValidateUserPassword(user, oldPassword); err != nil {
		return err
	}

	// Validate new password
	if err := validators.ValidatePassword(newPassword); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := s.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password
	user.Password = hashedPassword
	return s.userRepo.Update(user)
}

// ChangeAdminPassword changes an admin's password
func (s *AuthService) ChangeAdminPassword(admin *models.Admin, oldPassword, newPassword string) error {
	// Validate old password
	if err := s.ValidateAdminPassword(admin, oldPassword); err != nil {
		return err
	}

	// Validate new password
	if err := validators.ValidatePassword(newPassword); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := s.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password
	admin.Password = hashedPassword
	return s.adminRepo.Update(admin)
}
