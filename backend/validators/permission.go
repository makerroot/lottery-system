package validators

import (
	"lottery-system/constants"
	"lottery-system/errors"
)

// CanAdminModifyCompany checks if an admin can modify a company
func CanAdminModifyCompany(isSuperAdmin bool, adminCompanyID *int, targetCompanyID int) bool {
	// Super admins can modify any company
	if isSuperAdmin {
		return true
	}

	// Regular admins can only modify their own company
	if adminCompanyID != nil && *adminCompanyID == targetCompanyID {
		return true
	}

	return false
}

// CanAdminAccessUser checks if an admin can access a user
func CanAdminAccessUser(isSuperAdmin bool, adminCompanyID *int, userCompanyID int) bool {
	// Super admins can access any user
	if isSuperAdmin {
		return true
	}

	// Regular admins can only access users in their company
	if adminCompanyID != nil && *adminCompanyID == userCompanyID {
		return true
	}

	return false
}

// CanAdminDeleteAdmin checks if an admin can delete another admin
func CanAdminDeleteAdmin(isSuperAdmin bool, adminID int, targetAdminID int) error {
	// Cannot delete yourself
	if adminID == targetAdminID {
		return errors.NewBusinessLogicError(constants.ErrCannotDeleteSelf)
	}

	// Only super admins can delete admins
	if !isSuperAdmin {
		return errors.NewAuthorizationError(constants.ErrPermissionDenied)
	}

	return nil
}

// CanAdminUpdateAdmin checks if an admin can update another admin
func CanAdminUpdateAdmin(isSuperAdmin bool, adminID int, targetAdminID int) error {
	// Super admins can update any admin
	if isSuperAdmin {
		return nil
	}

	// Regular admins can only update themselves
	if adminID != targetAdminID {
		return errors.NewAuthorizationError(constants.ErrPermissionDenied)
	}

	return nil
}

// CanAdminCreateAdmin checks if an admin can create another admin
func CanAdminCreateAdmin(isSuperAdmin bool) error {
	// Only super admins can create admins
	if !isSuperAdmin {
		return errors.NewAuthorizationError(constants.ErrPermissionDenied)
	}

	return nil
}

// CanAdminManagePrizeLevels checks if an admin can manage prize levels for a company
func CanAdminManagePrizeLevels(isSuperAdmin bool, adminCompanyID *int, targetCompanyID int) error {
	if !CanAdminModifyCompany(isSuperAdmin, adminCompanyID, targetCompanyID) {
		return errors.NewAuthorizationError(constants.ErrPermissionDenied)
	}
	return nil
}

// CanAdminViewCompanyStats checks if an admin can view company statistics
func CanAdminViewCompanyStats(isSuperAdmin bool, adminCompanyID *int, targetCompanyID int) error {
	if !CanAdminModifyCompany(isSuperAdmin, adminCompanyID, targetCompanyID) {
		return errors.NewAuthorizationError(constants.ErrPermissionDenied)
	}
	return nil
}

// GetCompanyFilter returns a database filter for company-based queries
func GetCompanyFilter(isSuperAdmin bool, companyID *int) map[string]interface{} {
	if isSuperAdmin {
		// Super admins can see all companies (no filter)
		return nil
	}

	// Regular admins are filtered to their company
	if companyID != nil {
		return map[string]interface{}{
			"company_id": *companyID,
		}
	}

	// No company assigned (shouldn't happen for regular admins)
	return map[string]interface{}{
		"company_id": 0, // This will match nothing
	}
}

// ValidateAdminCompany validates admin company assignment
func ValidateAdminCompany(isSuperAdmin bool, companyID *int) error {
	if isSuperAdmin && companyID != nil {
		return errors.NewValidationError(constants.ErrSuperAdminNoCompany)
	}

	if !isSuperAdmin && companyID == nil {
		return errors.NewValidationError(constants.ErrAdminMustHaveCompany)
	}

	return nil
}
