package employee

import (
	"encoding/json"
	"net/http"
	"strconv"

	"starterpack-golang-cleanarch/internal/platform/http/middleware"
	"starterpack-golang-cleanarch/internal/utils"
	"starterpack-golang-cleanarch/internal/utils/errors"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type EmployeeHandler struct {
	service   *EmployeeService
	validator *validator.Validate
}

// NewEmployeeHandler creates a new instance of EmployeeHandler.
func NewEmployeeHandler(s *EmployeeService, v *validator.Validate) *EmployeeHandler {
	return &EmployeeHandler{service: s, validator: v}
}

// RegisterRoutes registers Employee-related API routes to the router.
func (h *EmployeeHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/employees", h.CreateEmployee).Methods("POST")
	router.HandleFunc("/employees", h.GetEmployees).Methods("GET")
	router.HandleFunc("/employees/{id}", h.GetEmployeeByID).Methods("GET") // Path ID is string for mux
}

// CreateEmployee handles the request to create a new employee.
func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var req CreateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleHTTPError(w, errors.NewBadRequest("Invalid request payload", nil), r)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.HandleHTTPError(w, errors.NewBadRequest(err.Error(), nil), r)
		return
	}

	tenantID, ok := r.Context().Value(middleware.ContextKeyTenantID).(string)
	if !ok || tenantID == "" {
		utils.HandleHTTPError(w, errors.ErrUnauthorized, r)
		return
	}

	employee, err := h.service.CreateEmployee(r.Context(), tenantID, req)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusCreated, employee)
}

// GetEmployees handles the request to retrieve a list of employees with pagination.
func (h *EmployeeHandler) GetEmployees(w http.ResponseWriter, r *http.Request) {
	var req GetEmployeesRequest
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		req.Page, _ = strconv.Atoi(pageStr)
	} else {
		req.Page = 1
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		req.Limit, _ = strconv.Atoi(limitStr)
	} else {
		req.Limit = 10
	}
	req.Query = r.URL.Query().Get("query")
	req.Status = r.URL.Query().Get("status")

	if err := h.validator.Struct(req); err != nil {
		utils.HandleHTTPError(w, errors.NewBadRequest(err.Error(), nil), r)
		return
	}

	tenantID, ok := r.Context().Value(middleware.ContextKeyTenantID).(string)
	if !ok || tenantID == "" {
		utils.HandleHTTPError(w, errors.ErrUnauthorized, r)
		return
	}

	response, err := h.service.GetEmployees(r.Context(), tenantID, req)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, response)
}

// GetEmployeeByID handles the request to retrieve an employee by ID.
func (h *EmployeeHandler) GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeIDStr := vars["id"] // ID from path is always string

	// No need to parse UUID here, as DB ID is SERIAL (int).
	// The service layer will handle string to int64 conversion.
	// Basic validation for numeric ID can be added here if needed.

	tenantID, ok := r.Context().Value(middleware.ContextKeyTenantID).(string)
	if !ok || tenantID == "" {
		utils.HandleHTTPError(w, errors.ErrUnauthorized, r)
		return
	}

	employee, err := h.service.GetEmployeeByID(r.Context(), tenantID, employeeIDStr) // Pass string ID to service
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, employee)
}
