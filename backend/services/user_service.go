package services

import (
	"fmt"
	"strings"

	"lottery-system/constants"
	"lottery-system/models"
	"lottery-system/repositories"
	"lottery-system/utils"
	"lottery-system/validators"
)

// CreateUserRequest represents a request to create a user
type CreateUserRequest struct {
	CompanyID int    `json:"company_id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
}

// BatchCreateUserRequest represents a request to batch create users
type BatchCreateUserRequest struct {
	CompanyID int      `json:"company_id"`
	Users     []string `json:"users"` // Format: "username,password,name"
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	HasDrawn *bool  `json:"has_drawn"`
}

// UserService handles user business logic
type UserService struct {
	userRepo    *repositories.UserRepository
	companyRepo *repositories.CompanyRepository
	authService *AuthService
}

// NewUserService creates a new user service
func NewUserService() *UserService {
	return &UserService{
		userRepo:    repositories.NewUserRepository(),
		companyRepo: repositories.NewCompanyRepository(),
		authService: NewAuthService(),
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(req *CreateUserRequest) (*models.User, error) {
	// Validate username
	if err := validators.ValidateUsername(req.Username); err != nil {
		return nil, err
	}

	// Validate password
	if err := validators.ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	// Validate name (if provided)
	if req.Name != "" {
		if err := validators.ValidateName(req.Name); err != nil {
			return nil, err
		}
	}

	// Validate phone (if provided)
	if req.Phone != "" {
		if err := validators.ValidatePhone(req.Phone); err != nil {
			return nil, err
		}
	}

	// Check if company exists
	company, err := s.companyRepo.FindByID(req.CompanyID)
	if err != nil {
		return nil, utils.NewNotFoundError("公司")
	}
	_ = company // Used for validation

	// Check if username already exists
	exists, err := s.userRepo.ExistsByUsername(req.Username, req.CompanyID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, utils.NewBusinessLogicError(constants.ErrUsernameExists)
	}

	// Check if phone already exists (if provided)
	if req.Phone != "" {
		exists, err := s.userRepo.ExistsByPhone(req.Phone, req.CompanyID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, utils.NewBusinessLogicError(constants.ErrPhoneExists)
		}
	}

	// Hash password
	hashedPassword, err := s.authService.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		CompanyID: req.CompanyID,
		Username:  req.Username,
		Password:  hashedPassword,
		Role:      models.RoleUser,
		Name:      req.Name,
		Phone:     req.Phone,
		HasDrawn:  false,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Reload with company info
	return s.userRepo.FindByIDWithPreload(user.ID)
}

// BatchCreateUsers creates multiple users in batch
func (s *UserService) BatchCreateUsers(req *BatchCreateUserRequest) ([]models.User, []string, error) {
	// Check if company exists
	company, err := s.companyRepo.FindByID(req.CompanyID)
	if err != nil {
		return nil, nil, utils.NewNotFoundError("公司")
	}
	_ = company // Used for validation

	var createdUsers []models.User
	var failedUsers []string

	for _, userStr := range req.Users {
		username, password, name, err := s.parseUserString(userStr)
		if err != nil {
			failedUsers = append(failedUsers, userStr+" (格式错误)")
			continue
		}

		// Validate
		if err := validators.ValidateUsername(username); err != nil {
			failedUsers = append(failedUsers, username+" (用户名格式错误)")
			continue
		}

		if err := validators.ValidatePassword(password); err != nil {
			failedUsers = append(failedUsers, username+" (密码格式错误)")
			continue
		}

		// Check if exists
		exists, err := s.userRepo.ExistsByUsername(username, req.CompanyID)
		if err != nil {
			failedUsers = append(failedUsers, username+" (查询失败)")
			continue
		}
		if exists {
			failedUsers = append(failedUsers, username+" (已存在)")
			continue
		}

		// Hash password
		hashedPassword, err := s.authService.HashPassword(password)
		if err != nil {
			failedUsers = append(failedUsers, username+" (密码加密失败)")
			continue
		}

		// Create user
		user := &models.User{
			CompanyID: req.CompanyID,
			Username:  username,
			Password:  hashedPassword,
			Role:      models.RoleUser,
			Name:      name,
			HasDrawn:  false,
		}

		if err := s.userRepo.Create(user); err != nil {
			failedUsers = append(failedUsers, username+" (创建失败)")
			continue
		}

		createdUsers = append(createdUsers, *user)
	}

	return createdUsers, failedUsers, nil
}

// parseUserString parses a user string in format "username,password,name"
func (s *UserService) parseUserString(userStr string) (username, password, name string, err error) {
	if userStr == "" {
		return "", "", "", fmt.Errorf("empty string")
	}

	// Trim and split by comma
	parts := splitAndTrim(userStr, ",")

	switch len(parts) {
	case 1:
		// Only username: use username as password and name
		username = parts[0]
		password = parts[0]
		name = parts[0]
	case 2:
		// username,password: use username as name
		username = parts[0]
		password = parts[1]
		name = parts[0]
	case 3:
		// username,password,name
		username = parts[0]
		password = parts[1]
		name = parts[2]
	default:
		return "", "", "", fmt.Errorf("invalid format")
	}

	return username, password, name, nil
}

// splitAndTrim splits a string and trims each part
func splitAndTrim(s, sep string) []string {
	parts := make([]string, 0)
	for _, part := range splitString(s, sep) {
		trimmed := trimSpace(part)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

// splitString splits a string by separator
func splitString(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, sep)
}

// trimSpace trims whitespace from a string
func trimSpace(s string) string {
	return strings.TrimSpace(s)
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(userID int, req *UpdateUserRequest) (*models.User, error) {
	// Check if user exists
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, utils.NewNotFoundError("用户")
	}

	// Build updates map
	updates := make(map[string]interface{})

	if req.Name != "" {
		if err := validators.ValidateName(req.Name); err != nil {
			return nil, err
		}
		updates["name"] = req.Name
	}

	if req.Phone != "" {
		if err := validators.ValidatePhone(req.Phone); err != nil {
			return nil, err
		}
		updates["phone"] = req.Phone
	}

	if req.HasDrawn != nil {
		updates["has_drawn"] = *req.HasDrawn
	}

	if len(updates) == 0 {
		return nil, utils.NewValidationError(constants.ErrInvalidInput)
	}

	// Update
	if err := s.userRepo.UpdateFields(userID, updates); err != nil {
		return nil, err
	}

	// Reload
	return s.userRepo.FindByIDWithPreload(userID)
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(userID int) error {
	// Check if user exists
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return utils.NewNotFoundError("用户")
	}

	return s.userRepo.Delete(userID)
}

// GetUsers gets users with optional filters
func (s *UserService) GetUsers(filters map[string]interface{}) ([]models.User, error) {
	return s.userRepo.FindAllWithPreload(filters)
}

// GetUserStats gets user statistics for a company
func (s *UserService) GetUserStats(companyID int) (map[string]interface{}, error) {
	// Count total users
	total, err := s.userRepo.CountByCompany(companyID)
	if err != nil {
		return nil, err
	}

	// Count undrawn users
	undrawn, err := s.userRepo.CountByCompanyAndStatus(companyID, false)
	if err != nil {
		return nil, err
	}

	// Count drawn users
	drawn, err := s.userRepo.CountByCompanyAndStatus(companyID, true)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_users":   total,
		"undrawn_users": undrawn,
		"drawn_users":   drawn,
	}, nil
}

// GetAvailableUsers gets users who haven't drawn yet
func (s *UserService) GetAvailableUsers(companyID int) ([]models.User, error) {
	return s.userRepo.FindAvailableUsers(companyID)
}
