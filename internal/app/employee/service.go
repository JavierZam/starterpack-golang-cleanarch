package employee

import (
	"context"
	"fmt"
	"strconv"

	"starterpack-golang-cleanarch/internal/domain"
	"starterpack-golang-cleanarch/internal/utils"
	"starterpack-golang-cleanarch/internal/utils/errors"
)

type EmployeeService struct {
	employeeRepo domain.EmployeeRepository
}

// NewEmployeeService creates a new instance of EmployeeService.
func NewEmployeeService(repo domain.EmployeeRepository) *EmployeeService {
	return &EmployeeService{repo}
}

// CreateEmployee handles the business logic for creating a new employee.
func (s *EmployeeService) CreateEmployee(ctx context.Context, tenantID string, req CreateEmployeeRequest) (*EmployeeResponse, error) {
	// TenantID di service sekarang bertipe string.
	// Tidak perlu parsing UUID di sini jika di domain/repo sudah VARCHAR.
	// Jika tenantID dari token selalu UUID, mungkin perlu validasi format di sini.

	// 1. Business Validation: Check if email already exists for this tenant
	existingEmployee, err := s.employeeRepo.FindByEmail(ctx, tenantID, req.Email) // tenantID langsung string
	if err != nil {
		return nil, errors.NewInternalServerError(fmt.Errorf("failed to check existing employee: %w", err), "Internal error during employee creation check.")
	}
	if existingEmployee != nil {
		return nil, ErrEmployeeAlreadyExists
	}

	// 2. Map Request DTO to Domain Model
	employee := &domain.Employee{
		TenantID:     tenantID, // tenantID langsung string
		Name:         req.Name,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		PasswordHash: "dummy-hashed-password", // FIX: Provide a dummy password hash
	}
	employee.GenerateID() // Ini akan mengisi TenantID dan Created/Updated timestamps

	// 3. Call Repository to persist data
	if err := s.employeeRepo.Save(ctx, employee); err != nil {
		return nil, errors.NewInternalServerError(fmt.Errorf("failed to save employee to database: %w", err), "Internal error saving employee.")
	}

	// 4. Map Domain Model to Response DTO
	resp := &EmployeeResponse{
		ID:          employee.ID, // ID sekarang int64
		Name:        employee.Name,
		Email:       employee.Email,
		PhoneNumber: employee.PhoneNumber,
		CreatedAt:   employee.CreatedAt.Format(utils.ISO8601TimeFormat),
		UpdatedAt:   employee.UpdatedAt.Format(utils.ISO8601TimeFormat),
		// PasswordHash field tidak perlu diisi di sini karena di EmployeeResponse sudah ada `json:"-"`
	}

	return resp, nil
}

// GetEmployees retrieves a list of employees with pagination.
func (s *EmployeeService) GetEmployees(ctx context.Context, tenantID string, req GetEmployeesRequest) (*GetEmployeesResponse, error) {
	// TenantID di service sekarang bertipe string.

	// 1. Call repository for total count and paginated data
	total, employees, err := s.employeeRepo.FindAll(ctx, tenantID, req.Page, req.Limit, req.Query) // tenantID langsung string
	if err != nil {
		return nil, errors.NewInternalServerError(fmt.Errorf("failed to fetch employees from repository: %w", err), "Internal error fetching employees.")
	}

	// 2. Map domain models to response DTOs
	employeeResponses := make([]EmployeeResponse, len(employees))
	for i, emp := range employees {
		employeeResponses[i] = EmployeeResponse{
			ID:          emp.ID, // ID sekarang int64
			Name:        emp.Name,
			Email:       emp.Email,
			PhoneNumber: emp.PhoneNumber,
			CreatedAt:   emp.CreatedAt.Format(utils.ISO8601TimeFormat),
			UpdatedAt:   emp.UpdatedAt.Format(utils.ISO8601TimeFormat),
		}
	}

	// 3. Build PaginationResponse using the helper from utils
	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))
	var nextPage, prevPage *int
	if req.Page < totalPages {
		np := req.Page + 1
		nextPage = &np
	}
	if req.Page > 1 {
		pp := req.Page - 1
		prevPage = &pp
	}

	return &GetEmployeesResponse{
		Data:       employeeResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
		NextPage:   nextPage,
		PrevPage:   prevPage,
	}, nil
}

// GetEmployeeByID retrieves a single employee by ID.
func (s *EmployeeService) GetEmployeeByID(ctx context.Context, tenantID string, employeeID string) (*EmployeeResponse, error) {
	// employeeID sekarang string, tidak perlu parsing UUID di sini jika di domain/repo sudah SERIAL.
	// Lakukan konversi ke int64 untuk FindByID repository
	parsedEmployeeID, err := strconv.ParseInt(employeeID, 10, 64)
	if err != nil {
		return nil, errors.NewBadRequest("Invalid employee ID format (must be integer)", nil)
	}

	employee, err := s.employeeRepo.FindByID(ctx, tenantID, parsedEmployeeID) // tenantID string, employeeID int64
	if err != nil {
		return nil, errors.NewInternalServerError(fmt.Errorf("failed to get employee from repository: %w", err), "Internal error fetching employee by ID.")
	}
	if employee == nil {
		return nil, ErrEmployeeNotFound
	}

	resp := &EmployeeResponse{
		ID:          employee.ID, // ID sekarang int64
		Name:        employee.Name,
		Email:       employee.Email,
		PhoneNumber: employee.PhoneNumber,
		CreatedAt:   employee.CreatedAt.Format(utils.ISO8601TimeFormat),
		UpdatedAt:   employee.UpdatedAt.Format(utils.ISO8601TimeFormat),
	}
	return resp, nil
}
