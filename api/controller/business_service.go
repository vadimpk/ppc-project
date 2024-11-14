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

type BusinessServiceHandler struct {
	serviceService services.BusinessServiceService
}

func NewBusinessServiceHandler(service services.BusinessServiceService) *BusinessServiceHandler {
	return &BusinessServiceHandler{
		serviceService: service,
	}
}

type CreateServiceRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Duration    int     `json:"duration"` // in minutes
	Price       int     `json:"price"`    // in cents
}

type UpdateServiceRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Duration    int     `json:"duration"` // in minutes
	Price       int     `json:"price"`    // in cents
	IsActive    bool    `json:"is_active"`
}

func (h *BusinessServiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Get business ID from URL parameters
	businessID, err := strconv.Atoi(chi.URLParam(r, "businessID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid business ID")
		return
	}

	// Verify user has permission to create services
	userRole, _ := middleware.GetRole(r.Context())
	if userRole != entity.RoleAdmin {
		response.Error(w, http.StatusForbidden, "unauthorized")
		return
	}

	var req CreateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	service := &entity.BusinessService{
		BusinessID:  businessID,
		Name:        req.Name,
		Description: req.Description,
		Duration:    req.Duration,
		Price:       req.Price,
		IsActive:    true,
	}

	if err := h.serviceService.Create(r.Context(), service); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create service")
		return
	}

	response.JSON(w, http.StatusCreated, service)
}

func (h *BusinessServiceHandler) Get(w http.ResponseWriter, r *http.Request) {
	// Get service ID from URL parameters
	serviceID, err := strconv.Atoi(chi.URLParam(r, "serviceID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid service ID")
		return
	}

	service, err := h.serviceService.Get(r.Context(), serviceID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get service")
		return
	}

	response.JSON(w, http.StatusOK, service)
}

func (h *BusinessServiceHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Get service ID from URL parameters
	serviceID, err := strconv.Atoi(chi.URLParam(r, "serviceID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid service ID")
		return
	}

	businessID, _ := middleware.GetBusinessID(r.Context())

	// Verify user has permission to update services
	userRole, _ := middleware.GetRole(r.Context())
	if userRole != entity.RoleAdmin {
		response.Error(w, http.StatusForbidden, "unauthorized")
		return
	}

	var req UpdateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	service := &entity.BusinessService{
		ID:          serviceID,
		BusinessID:  businessID,
		Name:        req.Name,
		Description: req.Description,
		Duration:    req.Duration,
		Price:       req.Price,
		IsActive:    req.IsActive,
	}

	if err := h.serviceService.Update(r.Context(), service); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to update service")
		return
	}

	response.JSON(w, http.StatusOK, service)
}

func (h *BusinessServiceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Get service ID from URL parameters
	serviceID, err := strconv.Atoi(chi.URLParam(r, "serviceID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid service ID")
		return
	}

	// Verify user has permission to delete services
	userRole, _ := middleware.GetRole(r.Context())
	if userRole != entity.RoleAdmin {
		response.Error(w, http.StatusForbidden, "unauthorized")
		return
	}

	if err := h.serviceService.Delete(r.Context(), serviceID); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to delete service")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *BusinessServiceHandler) List(w http.ResponseWriter, r *http.Request) {
	// Get business ID from URL parameters
	businessID, err := strconv.Atoi(chi.URLParam(r, "businessID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid business ID")
		return
	}

	// Verify user has permission to list services for this business
	currentBusinessID, _ := middleware.GetBusinessID(r.Context())
	if businessID != currentBusinessID {
		response.Error(w, http.StatusForbidden, "unauthorized")
		return
	}

	// Check query parameter for active-only filter
	activeOnly := r.URL.Query().Get("active") == "true"

	var services []entity.BusinessService
	var listErr error

	if activeOnly {
		services, listErr = h.serviceService.ListActive(r.Context(), businessID)
	} else {
		services, listErr = h.serviceService.List(r.Context(), businessID)
	}

	if listErr != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list services")
		return
	}

	response.JSON(w, http.StatusOK, services)
}

func (h *BusinessServiceHandler) ListEmployees(w http.ResponseWriter, r *http.Request) {
	serviceID, err := strconv.Atoi(chi.URLParam(r, "serviceID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid service ID")
		return
	}

	employees, err := h.serviceService.ListEmployee(r.Context(), serviceID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list employees")
		return
	}

	response.JSON(w, http.StatusOK, employees)
}
