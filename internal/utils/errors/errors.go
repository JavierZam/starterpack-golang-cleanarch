package errors

import (
	"net/http"
)

// AppError defines a generic interface for application-specific errors.
// All custom errors in the application should implement this interface.
type AppError interface {
	error
	Code() string                    // Unique error code (e.g., "NOT_FOUND", "INVALID_INPUT")
	Message() string                 // Human-readable message, suitable for user display
	Status() int                     // HTTP Status Code associated with the error
	Details() map[string]interface{} // Optional: for more detailed error information (e.g., validation errors)
	Unwrap() error                   // For error wrapping (Go 1.13+)
}

// baseError implements the AppError interface for common error types.
type baseError struct {
	err     error                  // Wrapped original error
	code    string                 // Unique code for the error type
	message string                 // User-friendly message
	status  int                    // HTTP status code
	details map[string]interface{} // Additional error details
}

// Error implements the standard Go error interface.
func (e *baseError) Error() string {
	if e.message != "" {
		return e.message
	}
	if e.err != nil {
		return e.err.Error()
	}
	return "An unknown error occurred"
}

// Code returns the unique error code.
func (e *baseError) Code() string { return e.code }

// Message returns the human-readable message.
func (e *baseError) Message() string { return e.message }

// Status returns the HTTP status code.
func (e *baseError) Status() int { return e.status }

// Details returns optional detailed error information.
func (e *baseError) Details() map[string]interface{} { return e.details }

// Unwrap returns the wrapped error, allowing errors.Is and errors.As to work.
func (e *baseError) Unwrap() error { return e.err }

// New creates a new baseError. This is the constructor for custom application errors.
func New(code, message string, status int, originalErr error, details map[string]interface{}) AppError {
	return &baseError{
		err:     originalErr,
		code:    code,
		message: message,
		status:  status,
		details: details,
	}
}

// --- Predefined Common Application Errors for consistent usage ---
var (
	ErrBadRequest         = New("BAD_REQUEST", "Invalid request payload or parameters", http.StatusBadRequest, nil, nil)
	ErrUnauthorized       = New("UNAUTHORIZED", "Authentication required", http.StatusUnauthorized, nil, nil)
	ErrForbidden          = New("FORBIDDEN", "Access denied for this resource", http.StatusForbidden, nil, nil)
	ErrNotFound           = New("NOT_FOUND", "Resource not found", http.StatusNotFound, nil, nil)
	ErrConflict           = New("CONFLICT", "Resource conflict or already exists", http.StatusConflict, nil, nil)
	ErrInternalServer     = New("INTERNAL_SERVER_ERROR", "An unexpected internal server error occurred", http.StatusInternalServerError, nil, nil)
	ErrServiceUnavailable = New("SERVICE_UNAVAILABLE", "Service is temporarily unavailable, please try again later", http.StatusServiceUnavailable, nil, nil)
)

// NewBadRequest creates a new bad request error with optional details.
// Useful for validation errors.
func NewBadRequest(message string, details map[string]interface{}) AppError {
	return New(ErrBadRequest.Code(), message, ErrBadRequest.Status(), nil, details)
}

// NewInternalServerError creates a new internal server error, wrapping the original error.
func NewInternalServerError(originalErr error, message string) AppError {
	if message == "" {
		message = ErrInternalServer.Message()
	}
	return New(ErrInternalServer.Code(), message, ErrInternalServer.Status(), originalErr, nil)
}

// Wrap wraps an existing error into an AppError with specified code, message, and status.
func Wrap(err error, code, message string, status int) AppError {
	if err == nil {
		return nil
	}
	return New(code, message, status, err, nil)
}
