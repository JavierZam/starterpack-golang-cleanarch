package employee

import (
	"net/http"

	"starterpack-golang-cleanarch/internal/utils/errors"
)

// Module-specific custom errors for the Employee domain.
var (
	ErrEmployeeNotFound      = errors.New("EMPLOYEE_NOT_FOUND", "Employee with given ID not found", http.StatusNotFound, nil, nil)
	ErrEmployeeAlreadyExists = errors.New("EMPLOYEE_ALREADY_EXISTS", "Employee with the given email or phone number already exists", http.StatusConflict, nil, nil)
	// Add other specific errors related to employee operations
)
