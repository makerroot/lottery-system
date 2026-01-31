package services

import (
	"lottery-system/constants"
	"lottery-system/models"
	"lottery-system/repositories"
	"lottery-system/utils"
	"lottery-system/validators"
)

// CreateAdminRequest represents a request to create an admin
type CreateAdminRequest struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	IsSuperAdmin bool   `json:"is_super_admin"`
	CompanyID    *int   `json:"company_id"`
}

// UpdateAdminRequest represents a request to update an admin
type UpdateAdminRequest struct {
	Username     *string `json:"username"`
	Password     string  `json:"password"`
	IsSuperAdmin *bool   `json:"is_super_admin"`
	CompanyID    *int    `json:"company_id"`
}

// AdminService handles admin business logic
type AdminService struct {
	adminRepo   *repositories.AdminRepository
	companyRepo *repositories.CompanyRepository
	authService *AuthService
}

// NewAdminService creates a new admin service
func NewAdminService() *AdminService {
	return &AdminService{
		adminRepo:   repositories.NewAdminRepository(),
		companyRepo: repositories.NewCompanyRepository(),
		authService: NewAuthService(),
	}
}

// CreateAdmin creates a new admin
func (s *AdminService) CreateAdmin(req *CreateAdminRequest) (*models.Admin, error) {
	// Validate username
	if err := validators.ValidateUsername(req.Username); err != nil {
		return nil, err
	}

	// Validate password
	if err := validators.ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	// Validate admin company assignment
	if err := validators.ValidateAdminCompany(req.IsSuperAdmin, req.CompanyID); err != nil {
		return nil, err
	}

	// Check if username already exists
	exists, err := s.adminRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, utils.NewBusinessLogicError(constants.ErrAdminExists)
	}

	// Validate company if specified
	if req.CompanyID != nil {
		_, err := s.companyRepo.FindByID(*req.CompanyID)
		if err != nil {
			return nil, utils.NewNotFoundError("公司")
		}
	}

	// Hash password
	hashedPassword, err := s.authService.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create admin
	admin := &models.Admin{
		Username:     req.Username,
		Password:     hashedPassword,
		CompanyID:    req.CompanyID,
		IsSuperAdmin: req.IsSuperAdmin,
	}

	// Set role based on super admin status
	if req.IsSuperAdmin {
		admin.Role = models.RoleSuperAdmin
	} else {
		admin.Role = models.RoleAdmin
	}

	if err := s.adminRepo.Create(admin); err != nil {
		return nil, err
	}

	// Reload with company info
	return s.adminRepo.FindByIDWithPreload(admin.ID)
}

// UpdateAdmin updates an admin
func (s *AdminService) UpdateAdmin(adminID int, req *UpdateAdminRequest, updaterIsSuperAdmin bool) (*models.Admin, error) {
	// Find admin
	admin, err := s.adminRepo.FindByID(adminID)
	if err != nil {
		return nil, utils.NewNotFoundError("管理员")
	}

	// Check update permission
	if err := validators.CanAdminUpdateAdmin(updaterIsSuperAdmin, 0, adminID); err != nil {
		return nil, err
	}

	// Update username (only super admin)
	if updaterIsSuperAdmin && req.Username != nil {
		// Check if new username already exists
		exists, err := s.adminRepo.ExistsByUsernameAndID(*req.Username, adminID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, utils.NewBusinessLogicError(constants.ErrAdminExists)
		}
		admin.Username = *req.Username
	}

	// Update password
	if req.Password != "" {
		hashedPassword, err := s.authService.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		admin.Password = hashedPassword
	}

	// Update admin type and company (only super admin)
	if updaterIsSuperAdmin {
		if req.IsSuperAdmin != nil {
			if *req.IsSuperAdmin {
				// Set as super admin
				admin.CompanyID = nil
				admin.IsSuperAdmin = true
				admin.Role = models.RoleSuperAdmin
			} else {
				// Set as regular admin
				if req.CompanyID == nil {
					return nil, utils.NewValidationError(constants.ErrAdminMustHaveCompany)
				}
				// Validate company exists
				_, err := s.companyRepo.FindByID(*req.CompanyID)
				if err != nil {
					return nil, utils.NewNotFoundError("公司")
				}
				admin.CompanyID = req.CompanyID
				admin.IsSuperAdmin = false
				admin.Role = models.RoleAdmin
			}
		} else if req.CompanyID != nil {
			// Only updating company, keep as regular admin
			_, err := s.companyRepo.FindByID(*req.CompanyID)
			if err != nil {
				return nil, utils.NewNotFoundError("公司")
			}
			admin.CompanyID = req.CompanyID
			admin.IsSuperAdmin = false
			admin.Role = models.RoleAdmin
		}
	}

	if err := s.adminRepo.Update(admin); err != nil {
		return nil, err
	}

	return s.adminRepo.FindByIDWithPreload(admin.ID)
}

// DeleteAdmin deletes an admin
func (s *AdminService) DeleteAdmin(adminID int, requesterID int, requesterIsSuperAdmin bool) error {
	// Check permission
	if err := validators.CanAdminDeleteAdmin(requesterIsSuperAdmin, requesterID, adminID); err != nil {
		return err
	}

	// Check if admin exists
	_, err := s.adminRepo.FindByID(adminID)
	if err != nil {
		return utils.NewNotFoundError("管理员")
	}

	return s.adminRepo.Delete(adminID)
}

// GetAdmins gets admins with optional company filter
func (s *AdminService) GetAdmins(isSuperAdmin bool, companyID *int) ([]models.Admin, error) {
	if isSuperAdmin {
		// Super admin sees all admins
		return s.adminRepo.FindAll(nil)
	}
	// Regular admin sees admins in their company
	return s.adminRepo.FindAll(companyID)
}

// GetAdminByID gets an admin by ID
func (s *AdminService) GetAdminByID(adminID int) (*models.Admin, error) {
	return s.adminRepo.FindByIDWithPreload(adminID)
}

// ChangePassword changes an admin's password
func (s *AdminService) ChangePassword(adminID int, oldPassword, newPassword string) error {
	admin, err := s.adminRepo.FindByID(adminID)
	if err != nil {
		return utils.NewNotFoundError("管理员")
	}

	return s.authService.ChangeAdminPassword(admin, oldPassword, newPassword)
}
