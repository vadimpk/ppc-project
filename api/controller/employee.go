package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/vadimpk/ppc-project/controller/middleware"
	"github.com/vadimpk/ppc-project/controller/response"
	"github.com/vadimpk/ppc-project/entity"
	"github.com/vadimpk/ppc-project/services"
)

type EmployeeHandler struct {
	employeeService services.EmployeeService
}

func NewEmployeeHandler(service services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		employeeService: service,
	}
}

type CreateEmployeeRequest struct {
	UserID         int     `json:"user_id"`
	Specialization *string `json:"specialization,omitempty"`
}

type UpdateEmployeeRequest struct {
	Specialization *string `json:"specialization,omitempty"`
	IsActive       bool    `json:"is_active"`
}

type AssignServicesRequest struct {
	ServiceIDs []int `json:"service_ids"`
}

func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Get business ID from URL parameters
	businessID, err := strconv.Atoi(chi.URLParam(r, "businessID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid business ID")
		return
	}

	// Verify user has permission to create employees
	userRole, _ := middleware.GetRole(r.Context())
	if userRole != entity.RoleAdmin {
		response.Error(w, http.StatusForbidden, "unauthorized")
		return
	}

	var req CreateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	employee := &entity.Employee{
		BusinessID:     businessID,
		UserID:         req.UserID,
		Specialization: req.Specialization,
		IsActive:       true,
	}

	if err := h.employeeService.Create(r.Context(), employee); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create employee")
		return
	}

	response.JSON(w, http.StatusCreated, employee)
}

func (h *EmployeeHandler) Get(w http.ResponseWriter, r *http.Request) {
	// Get employee ID from URL parameters
	employeeID, err := strconv.Atoi(chi.URLParam(r, "employeeID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid employee ID")
		return
	}

	employee, err := h.employeeService.Get(r.Context(), employeeID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get employee")
		return
	}

	// Verify user has permission to view the employee
	businessID, _ := middleware.GetBusinessID(r.Context())
	if employee.BusinessID != businessID {
		response.Error(w, http.StatusForbidden, "unauthorized")
		return
	}

	response.JSON(w, http.StatusOK, employee)
}

func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Get employee ID from URL parameters
	employeeID, err := strconv.Atoi(chi.URLParam(r, "employeeID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid employee ID")
		return
	}

	// Verify user has permission to update employees
	userRole, _ := middleware.GetRole(r.Context())
	if userRole != entity.RoleAdmin {
		response.Error(w, http.StatusForbidden, "unauthorized")
		return
	}

	var req UpdateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	employee := &entity.Employee{
		ID:             employeeID,
		Specialization: req.Specialization,
		IsActive:       req.IsActive,
	}

	if err := h.employeeService.Update(r.Context(), employee); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to update employee")
		return
	}

	response.JSON(w, http.StatusOK, employee)
}

func (h *EmployeeHandler) List(w http.ResponseWriter, r *http.Request) {
	// Get business ID from URL parameters
	businessID, err := strconv.Atoi(chi.URLParam(r, "businessID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid business ID")
		return
	}

	// Verify user has permission to list employees
	currentBusinessID, _ := middleware.GetBusinessID(r.Context())
	if businessID != currentBusinessID {
		response.Error(w, http.StatusForbidden, "unauthorized")
		return
	}

	employees, err := h.employeeService.List(r.Context(), businessID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list employees")
		return
	}

	response.JSON(w, http.StatusOK, employees)
}

func (h *EmployeeHandler) ListServices(w http.ResponseWriter, r *http.Request) {
	// Get employee ID from URL parameters
	employeeID, err := strconv.Atoi(chi.URLParam(r, "employeeID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid employee ID")
		return
	}

	services, err := h.employeeService.GetServices(r.Context(), employeeID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list employee services")
		return
	}

	response.JSON(w, http.StatusOK, services)
}

func (h *EmployeeHandler) AssignServices(w http.ResponseWriter, r *http.Request) {
	// Get employee ID from URL parameters
	employeeID, err := strconv.Atoi(chi.URLParam(r, "employeeID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid employee ID")
		return
	}

	// Verify user has permission to assign services
	userRole, _ := middleware.GetRole(r.Context())
	if userRole != entity.RoleAdmin {
		response.Error(w, http.StatusForbidden, "unauthorized")
		return
	}

	var req AssignServicesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.employeeService.AssignServices(r.Context(), employeeID, req.ServiceIDs); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to assign services")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"status": "services assigned"})
}

func (h *EmployeeHandler) RemoveServices(w http.ResponseWriter, r *http.Request) {
	// Get employee ID from URL parameters
	employeeID, err := strconv.Atoi(chi.URLParam(r, "employeeID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid employee ID")
		return
	}

	// Verify user has permission to remove services
	userRole, _ := middleware.GetRole(r.Context())
	if userRole != entity.RoleAdmin {
		response.Error(w, http.StatusForbidden, "unauthorized")
		return
	}

	var req AssignServicesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.employeeService.RemoveServices(r.Context(), employeeID, req.ServiceIDs); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to remove services")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"status": "services removed"})
}
