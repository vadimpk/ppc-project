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

type BusinessHandler struct {
	businessService services.BusinessService
}

func NewBusinessHandler(service services.BusinessService) *BusinessHandler {
	return &BusinessHandler{
		businessService: service,
	}
}

type CreateBusinessRequest struct {
	Name string `json:"name"`
}

type UpdateBusinessRequest struct {
	Name string `json:"name"`
}

type UpdateBusinessAppearanceRequest struct {
	LogoURL     string                 `json:"logo_url"`
	ColorScheme map[string]interface{} `json:"color_scheme"`
}

func (h *BusinessHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateBusinessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	business := &entity.Business{
		Name: req.Name,
	}

	if err := h.businessService.Create(r.Context(), business); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create business")
		return
	}

	response.JSON(w, http.StatusCreated, business)
}

func (h *BusinessHandler) Get(w http.ResponseWriter, r *http.Request) {
	businessID, err := strconv.Atoi(chi.URLParam(r, "businessID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid business ID")
		return
	}

	business, err := h.businessService.Get(r.Context(), businessID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get business")
		return
	}

	response.JSON(w, http.StatusOK, business)
}

func (h *BusinessHandler) Update(w http.ResponseWriter, r *http.Request) {
	businessID, err := strconv.Atoi(chi.URLParam(r, "businessID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid business ID")
		return
	}

	// Verify user has permission to update the business
	userRole, _ := middleware.GetRole(r.Context())
	if userRole != entity.RoleAdmin {
		response.Error(w, http.StatusForbidden, "unauthorized")
		return
	}

	var req UpdateBusinessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	business := &entity.Business{
		ID:   businessID,
		Name: req.Name,
	}

	if err := h.businessService.Update(r.Context(), business); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to update business")
		return
	}

	response.JSON(w, http.StatusOK, business)
}

func (h *BusinessHandler) UpdateAppearance(w http.ResponseWriter, r *http.Request) {
	businessID, err := strconv.Atoi(chi.URLParam(r, "businessID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid business ID")
		return
	}

	// Verify user has permission to update the business appearance
	userRole, _ := middleware.GetRole(r.Context())
	if userRole != entity.RoleAdmin {
		response.Error(w, http.StatusForbidden, "unauthorized")
		return
	}

	var req UpdateBusinessAppearanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.businessService.UpdateAppearance(r.Context(), businessID, req.LogoURL, req.ColorScheme); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to update business appearance")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

type SearchBusinessResponse struct {
	Businesses []entity.Business        `json:"businesses"`
	Services   []entity.BusinessService `json:"services"`
}

func (h *BusinessHandler) SearchBusinessAndServices(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	if search == "" {
		response.Error(w, http.StatusBadRequest, "search query is required")
		return
	}

	businesses, err := h.businessService.ListBySearch(r.Context(), search)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to search businesses")
		return
	}

	businessServices, err := h.businessService.ListServicesBySearch(r.Context(), search)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to search services")
		return
	}

	response.JSON(w, http.StatusOK, SearchBusinessResponse{
		Businesses: businesses,
		Services:   businessServices,
	})
}
