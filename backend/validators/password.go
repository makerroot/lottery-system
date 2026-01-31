package validators

import (
	"errors"
	"strings"
	"unicode"

	"lottery-system/constants"
)

// PasswordStrength represents the strength level of a password
type PasswordStrength int

const (
	WeakPassword   PasswordStrength = 0
	MediumPassword PasswordStrength = 1
	StrongPassword PasswordStrength = 2
)

// PasswordResult represents the result of password validation
type PasswordResult struct {
	Valid    bool
	Strength PasswordStrength
	Errors   []string
}

// ValidatePasswordEnhanced validates a password with enhanced security checks
func ValidatePasswordEnhanced(password string) *PasswordResult {
	result := &PasswordResult{
		Valid:  true,
		Errors: make([]string, 0),
	}

	// Check length
	if len(password) < 8 {
		result.Valid = false
		result.Errors = append(result.Errors, "密码长度至少8位")
	}

	if len(password) > 128 {
		result.Valid = false
		result.Errors = append(result.Errors, "密码长度不能超过128位")
	}

	// Check for common weak passwords
	if isCommonPassword(password) {
		result.Valid = false
		result.Errors = append(result.Errors, "密码过于简单，请使用更复杂的密码")
	}

	// Check character variety
	hasLower := hasLowercase(password)
	hasUpper := hasUppercase(password)
	hasDigit := hasDigit(password)
	hasSpecial := hasSpecial(password)

	varietyCount := 0
	if hasLower {
		varietyCount++
	}
	if hasUpper {
		varietyCount++
	}
	if hasDigit {
		varietyCount++
	}
	if hasSpecial {
		varietyCount++
	}

	// Require at least 3 types of characters
	if varietyCount < 3 {
		result.Valid = false
		result.Errors = append(result.Errors, "密码必须包含大写字母、小写字母、数字中的至少三种")
	}

	// Calculate strength
	if len(result.Errors) == 0 {
		if varietyCount >= 4 && len(password) >= 12 {
			result.Strength = StrongPassword
		} else if varietyCount >= 3 && len(password) >= 8 {
			result.Strength = MediumPassword
		} else {
			result.Strength = WeakPassword
		}
	} else {
		result.Strength = WeakPassword
	}

	return result
}

// ValidatePasswordForUser validates password for regular users (simpler requirements)
func ValidatePasswordForUser(password string) error {
	if len(password) < 6 {
		return errors.New(constants.ErrInvalidPasswordLength)
	}

	if len(password) > 100 {
		return errors.New("密码长度不能超过100位")
	}

	// Check for obviously weak passwords
	if isCommonPassword(password) {
		return errors.New("密码过于简单")
	}

	return nil
}

// ValidatePasswordForAdmin validates password for admins (stronger requirements)
func ValidatePasswordForAdmin(password string) error {
	if len(password) < 8 {
		return errors.New(constants.ErrInvalidPasswordLength8)
	}

	if len(password) > 100 {
		return errors.New("密码长度不能超过100位")
	}

	// Check for common weak passwords
	if isCommonPassword(password) {
		return errors.New("密码过于简单，请使用更复杂的密码")
	}

	// Require at least 2 types of characters
	hasLower := hasLowercase(password)
	hasUpper := hasUppercase(password)
	hasDigit := hasDigit(password)
	hasSpecial := hasSpecial(password)

	varietyCount := 0
	if hasLower {
		varietyCount++
	}
	if hasUpper {
		varietyCount++
	}
	if hasDigit {
		varietyCount++
	}
	if hasSpecial {
		varietyCount++
	}

	if varietyCount < 2 {
		return errors.New("密码必须包含大写字母、小写字母、数字中的至少两种")
	}

	return nil
}

// hasLowercase checks if password has lowercase letters
func hasLowercase(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}

// hasUppercase checks if password has uppercase letters
func hasUppercase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

// hasDigit checks if password has digits
func hasDigit(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

// hasSpecial checks if password has special characters
func hasSpecial(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

// isCommonPassword checks if password is in the common weak passwords list
func isCommonPassword(password string) bool {
	// Convert to lowercase for comparison
	passwordLower := strings.ToLower(password)

	// List of common weak passwords
	commonPasswords := []string{
		"password", "12345678", "123456789", "qwerty", "abc123",
		"monkey", "1234567", "letmein", "trustno1", "dragon",
		"baseball", "11111111", "iloveyou", "master", "sunshine",
		"ashley", "bailey", "shadow", "123123", "654321",
		"superman", "qazwsx", "michael", "123456", "password1",
		"admin", "welcome", "login", "princess", "solo",
		"azerty", "password123", "123qwe", "1q2w3e4r", "qwertyuiop",
	}

	for _, common := range commonPasswords {
		if passwordLower == common {
			return true
		}
	}

	// Check for all same characters
	if len(password) > 0 {
		firstChar := password[0]
		allSame := true
		for _, c := range password {
			if c != rune(firstChar) {
				allSame = false
				break
			}
		}
		if allSame {
			return true
		}
	}

	// Check for sequential characters
	if isSequential(password) {
		return true
	}

	return false
}

// isSequential checks if password contains sequential characters
func isSequential(password string) bool {
	if len(password) < 4 {
		return false
	}

	// Check for forward sequences (like "1234", "abcd")
	for i := 0; i <= len(password)-4; i++ {
		if isSequentialAt(password, i, 1) {
			return true
		}
	}

	// Check for reverse sequences (like "4321", "dcba")
	for i := 0; i <= len(password)-4; i++ {
		if isSequentialAt(password, i, -1) {
			return true
		}
	}

	return false
}

// isSequentialAt checks if sequence starts at index with given direction
func isSequentialAt(s string, start, direction int) bool {
	for i := 0; i < 3; i++ {
		current := s[start+i]
		next := s[start+i+1]

		// Check if characters are sequential
		if next != current+byte(direction) {
			return false
		}
	}
	return true
}
