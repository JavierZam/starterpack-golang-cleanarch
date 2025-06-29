package employee

import "starterpack-golang-cleanarch/internal/utils"

// CreateEmployeeRequest is the DTO for creating a new employee.
type CreateEmployeeRequest struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	// Add other fields relevant for employee creation
}

// EmployeeResponse is the DTO for responding with employee details.
type EmployeeResponse struct {
	ID          int64  `json:"id"` // Changed to int64 to match DB SERIAL PRIMARY KEY
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	CreatedAt   string `json:"created_at"` // Formatted as ISO8601 string
	UpdatedAt   string `json:"updated_at"`
	// Add other fields that should be exposed in the API response
	PasswordHash string `json:"-"` // FIX: Don't expose password_hash in API response
}

// GetEmployeesRequest is the DTO for querying employees with pagination and filters.
type GetEmployeesRequest struct {
	utils.PaginationRequest
	Status string `query:"status"` // Example custom filter for employees
	// Add other specific filter fields
}

// GetEmployeesResponse is the DTO for responding with a paginated list of employees.
type GetEmployeesResponse = utils.PaginationResponse[EmployeeResponse]
