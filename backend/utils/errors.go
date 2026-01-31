package utils

import (
	"lottery-system/errors"
)

// NewValidationError creates a new validation error (convenience function)
func NewValidationError(message string) *errors.ValidationError {
	return errors.NewValidationError(message)
}

// NewValidationErrorWithField creates a new validation error with field name
func NewValidationErrorWithField(field, message string) *errors.ValidationError {
	return errors.NewValidationErrorWithField(field, message)
}

// NewAppError creates a new application error (convenience function)
func NewAppError(code, message string) *errors.AppError {
	return errors.NewAppError(code, message)
}

// NewAuthenticationError creates a new authentication error (convenience function)
func NewAuthenticationError(message string) *errors.AuthenticationError {
	return errors.NewAuthenticationError(message)
}

// NewAuthorizationError creates a new authorization error (convenience function)
func NewAuthorizationError(message string) *errors.AuthorizationError {
	return errors.NewAuthorizationError(message)
}

// NewNotFoundError creates a new not found error (convenience function)
func NewNotFoundError(resource string) *errors.NotFoundError {
	return errors.NewNotFoundError(resource)
}

// NewBusinessLogicError creates a new business logic error (convenience function)
func NewBusinessLogicError(message string) *errors.BusinessLogicError {
	return errors.NewBusinessLogicError(message)
}
