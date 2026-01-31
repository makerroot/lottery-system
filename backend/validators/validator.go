// Package validators provides input validation functions for the lottery system.
//
// It centralizes all validation logic to ensure consistency across the application
// and provide clear, user-friendly error messages.
package validators

import (
	"regexp"
	"strings"

	"lottery-system/constants"
	"lottery-system/utils"
)

var (
	phoneRegex       = regexp.MustCompile(constants.PhonePattern)
	usernameRegex    = regexp.MustCompile(constants.UsernamePattern)
	emailRegex       = regexp.MustCompile(constants.EmailPattern)
	companyCodeRegex = regexp.MustCompile(constants.CompanyCodePattern)
)

// ValidateUsername validates a username string
func ValidateUsername(username string) error {
	if username == "" {
		return utils.NewValidationError(constants.ErrRequiredField)
	}

	if len(username) < constants.MinUsernameLength || len(username) > constants.MaxUsernameLength {
		return utils.NewValidationError(constants.ErrInvalidUsername)
	}

	if !usernameRegex.MatchString(username) {
		return utils.NewValidationError(constants.ErrInvalidUsername)
	}

	return nil
}

// ValidatePassword validates a password string
func ValidatePassword(password string) error {
	if password == "" {
		return utils.NewValidationError(constants.ErrRequiredField)
	}

	if len(password) < constants.MinPasswordLength {
		return utils.NewValidationError(constants.ErrInvalidPasswordLength)
	}

	return nil
}

// ValidatePasswordStrong validates a password with stronger requirements
func ValidatePasswordStrong(password string) error {
	if password == "" {
		return utils.NewValidationError(constants.ErrRequiredField)
	}

	if len(password) < 8 {
		return utils.NewValidationError(constants.ErrInvalidPasswordLength8)
	}

	return nil
}

// ValidatePhone validates a phone number string
func ValidatePhone(phone string) error {
	if phone == "" {
		return nil // Phone is optional
	}

	if !phoneRegex.MatchString(phone) {
		return utils.NewValidationError(constants.ErrInvalidPhone)
	}

	return nil
}

// ValidateName validates a name string
func ValidateName(name string) error {
	if name == "" {
		return nil // Name is optional in some contexts
	}

	name = strings.TrimSpace(name)
	if len(name) < constants.MinNameLength || len(name) > constants.MaxNameLength {
		return utils.NewValidationError(constants.ErrInvalidName)
	}

	return nil
}

// ValidateEmail validates an email string
func ValidateEmail(email string) error {
	if email == "" {
		return nil // Email is optional
	}

	if !emailRegex.MatchString(email) {
		return utils.NewValidationError(constants.ErrInvalidEmail)
	}

	return nil
}

// ValidateCompanyCode validates a company code string
func ValidateCompanyCode(code string) error {
	if code == "" {
		return utils.NewValidationError(constants.ErrRequiredField)
	}

	code = strings.TrimSpace(code)
	if len(code) > constants.MaxCompanyCodeLength {
		return utils.NewValidationError(constants.ErrInvalidCompanyCode)
	}

	if !companyCodeRegex.MatchString(code) {
		return utils.NewValidationError(constants.ErrInvalidCompanyCode)
	}

	return nil
}

// ValidateRequired validates that a string field is not empty
func ValidateRequired(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return utils.NewValidationError(fieldName + "不能为空")
	}
	return nil
}

// ValidateLength validates string length
func ValidateLength(value string, min, max int) error {
	length := len(strings.TrimSpace(value))
	if length < min || length > max {
		return utils.NewValidationError(constants.ErrInvalidInput)
	}
	return nil
}

// ValidateIntRange validates an integer is within range
func ValidateIntRange(value, min, max int) error {
	if value < min || value > max {
		return utils.NewValidationError(constants.ErrInvalidInput)
	}
	return nil
}
