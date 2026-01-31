// Package errors provides custom error types for the lottery system.
//
// It defines structured error types that can be used throughout the application
// for better error handling and user feedback.
package errors

import "fmt"

// AppError represents a structured application error
type AppError struct {
	Code    string                 // Error code for API responses
	Message string                 // Human-readable error message
	Details map[string]interface{} // Additional error details
	Err     error                  // Underlying error (if any)
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new application error
func NewAppError(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// NewAppErrorWithErr creates a new application error with an underlying error
func NewAppErrorWithErr(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string // Field name that failed validation
	Message string // Validation error message
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(message string) *ValidationError {
	return &ValidationError{
		Message: message,
	}
}

// NewValidationErrorWithField creates a new validation error with field name
func NewValidationErrorWithField(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// AuthenticationError represents an authentication failure
type AuthenticationError struct {
	Message string
}

// Error implements the error interface
func (e *AuthenticationError) Error() string {
	return e.Message
}

// NewAuthenticationError creates a new authentication error
func NewAuthenticationError(message string) *AuthenticationError {
	return &AuthenticationError{
		Message: message,
	}
}

// AuthorizationError represents an authorization failure
type AuthorizationError struct {
	Message string
}

// Error implements the error interface
func (e *AuthorizationError) Error() string {
	return e.Message
}

// NewAuthorizationError creates a new authorization error
func NewAuthorizationError(message string) *AuthorizationError {
	return &AuthorizationError{
		Message: message,
	}
}

// NotFoundError represents a resource not found error
type NotFoundError struct {
	Resource string // Resource type (e.g., "user", "admin", "company")
	ID       string // Resource identifier
}

// Error implements the error interface
func (e *NotFoundError) Error() string {
	if e.ID != "" {
		return fmt.Sprintf("%s not found: %s", e.Resource, e.ID)
	}
	return fmt.Sprintf("%s not found", e.Resource)
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(resource string) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
	}
}

// NewNotFoundErrorWithID creates a new not found error with ID
func NewNotFoundErrorWithID(resource, id string) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
	}
}

// BusinessLogicError represents a business logic violation
type BusinessLogicError struct {
	Message string
	Details map[string]interface{}
}

// Error implements the error interface
func (e *BusinessLogicError) Error() string {
	return e.Message
}

// NewBusinessLogicError creates a new business logic error
func NewBusinessLogicError(message string) *BusinessLogicError {
	return &BusinessLogicError{
		Message: message,
	}
}

// IsValidationError checks if an error is a ValidationError
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

// IsAuthenticationError checks if an error is an AuthenticationError
func IsAuthenticationError(err error) bool {
	_, ok := err.(*AuthenticationError)
	return ok
}

// IsAuthorizationError checks if an error is an AuthorizationError
func IsAuthorizationError(err error) bool {
	_, ok := err.(*AuthorizationError)
	return ok
}

// IsNotFoundError checks if an error is a NotFoundError
func IsNotFoundError(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}
