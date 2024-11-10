package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/vadimpk/ppc-project/controller/middleware"
	"github.com/vadimpk/ppc-project/controller/response"
	"github.com/vadimpk/ppc-project/entity"
	"github.com/vadimpk/ppc-project/services"
)

type ScheduleHandler struct {
	scheduleService services.ScheduleService
}

func NewScheduleHandler(service services.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleService: service,
	}
}

type CreateTemplateRequest struct {
	DayOfWeek int       `json:"day_of_week"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	IsBreak   bool      `json:"is_break"`
}

type UpdateTemplateRequest struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	IsBreak   bool      `json:"is_break"`
}

type CreateOverrideRequest struct {
	OverrideDate time.Time  `json:"override_date"`
	StartTime    *time.Time `json:"start_time,omitempty"`
	EndTime      *time.Time `json:"end_time,omitempty"`
	IsWorkingDay bool       `json:"is_working_day"`
	IsBreak      bool       `json:"is_break"`
}

type UpdateOverrideRequest struct {
	StartTime    *time.Time `json:"start_time,omitempty"`
	EndTime      *time.Time `json:"end_time,omitempty"`
	IsWorkingDay bool       `json:"is_working_day"`
	IsBreak      bool       `json:"is_break"`
}

func (h *ScheduleHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	employeeID, err := strconv.Atoi(chi.URLParam(r, "employeeID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid employee ID")
		return
	}

	// Verify permissions
	userRole, _ := middleware.GetRole(r.Context())
	if userRole != entity.RoleAdmin {
		response.Error(w, http.StatusForbidden, "unauthorized")
		return
	}

	var req CreateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	template := &entity.ScheduleTemplate{
		EmployeeID: employeeID,
		DayOfWeek:  req.DayOfWeek,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		IsBreak:    req.IsBreak,
	}

	if err := h.scheduleService.CreateTemplate(r.Context(), template); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create template")
		return
	}

	response.JSON(w, http.StatusCreated, template)
}

func (h *ScheduleHandler) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	templateID, err := strconv.Atoi(chi.URLParam(r, "templateID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid template ID")
		return
	}

	employeeID, err := strconv.Atoi(chi.URLParam(r, "employeeID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid employee ID")
		return
	}

	// Verify permissions
	userRole, _ := middleware.GetRole(r.Context())
	if userRole != entity.RoleAdmin {
		response.Error(w, http.StatusForbidden, "unauthorized")
		return
	}

	var req UpdateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	template := &entity.ScheduleTemplate{
		ID:         templateID,
		EmployeeID: employeeID,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		IsBreak:    req.IsBreak,
	}

	if err := h.scheduleService.UpdateTemplate(r.Context(), template); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to update template")
		return
	}

	response.JSON(w, http.StatusOK, template)
}

func (h *ScheduleHandler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	templateID, err := strconv.Atoi(chi.URLParam(r, "templateID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid template ID")
		return
	}

	// Verify permissions
	userRole, _ := middleware.GetRole(r.Context())
	if userRole != entity.RoleAdmin {
		response.Error(w, http.StatusForbidden, "unauthorized")
		return
	}

	if err := h.scheduleService.DeleteTemplate(r.Context(), templateID); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to delete template")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *ScheduleHandler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	employeeID, err := strconv.Atoi(chi.URLParam(r, "employeeID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid employee ID")
		return
	}

	templates, err := h.scheduleService.ListTemplates(r.Context(), employeeID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list templates")
		return
	}

	response.JSON(w, http.StatusOK, templates)
}

func (h *ScheduleHandler) CreateOverride(w http.ResponseWriter, r *http.Request) {
	employeeID, err := strconv.Atoi(chi.URLParam(r, "employeeID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid employee ID")
		return
	}

	// Verify permissions
	userRole, _ := middleware.GetRole(r.Context())
	if userRole != entity.RoleAdmin {
		response.Error(w, http.StatusForbidden, "unauthorized")
		return
	}

	var req CreateOverrideRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	override := &entity.ScheduleOverride{
		EmployeeID:   employeeID,
		OverrideDate: req.OverrideDate,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		IsWorkingDay: req.IsWorkingDay,
		IsBreak:      req.IsBreak,
	}

	if err := h.scheduleService.CreateOverride(r.Context(), override); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create override")
		return
	}

	response.JSON(w, http.StatusCreated, override)
}

func (h *ScheduleHandler) ListOverrides(w http.ResponseWriter, r *http.Request) {
	employeeID, err := strconv.Atoi(chi.URLParam(r, "employeeID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid employee ID")
		return
	}

	// Parse date range from query parameters
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if startDate == "" || endDate == "" {
		response.Error(w, http.StatusBadRequest, "start_date and end_date are required")
		return
	}

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid start_date format")
		return
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid end_date format")
		return
	}

	if end.Before(start) {
		response.Error(w, http.StatusBadRequest, "end_date must be after start_date")
		return
	}

	overrides, err := h.scheduleService.ListOverrides(r.Context(), employeeID, start, end)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list overrides")
		return
	}

	response.JSON(w, http.StatusOK, overrides)
}

func (h *ScheduleHandler) UpdateOverride(w http.ResponseWriter, r *http.Request) {
	overrideID, err := strconv.Atoi(chi.URLParam(r, "overrideID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid override ID")
		return
	}

	employeeID, err := strconv.Atoi(chi.URLParam(r, "employeeID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid employee ID")
		return
	}

	// Verify permissions
	userRole, _ := middleware.GetRole(r.Context())
	if userRole != entity.RoleAdmin {
		response.Error(w, http.StatusForbidden, "unauthorized")
		return
	}

	var req UpdateOverrideRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	override := &entity.ScheduleOverride{
		ID:           overrideID,
		EmployeeID:   employeeID,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		IsWorkingDay: req.IsWorkingDay,
		IsBreak:      req.IsBreak,
	}

	if err := h.scheduleService.UpdateOverride(r.Context(), override); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to update override")
		return
	}

	response.JSON(w, http.StatusOK, override)
}

func (h *ScheduleHandler) DeleteOverride(w http.ResponseWriter, r *http.Request) {
	overrideID, err := strconv.Atoi(chi.URLParam(r, "overrideID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid override ID")
		return
	}

	// Verify permissions
	userRole, _ := middleware.GetRole(r.Context())
	if userRole != entity.RoleAdmin {
		response.Error(w, http.StatusForbidden, "unauthorized")
		return
	}

	if err := h.scheduleService.DeleteOverride(r.Context(), overrideID); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to delete override")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// Additional handler for getting the combined schedule
//func (h *ScheduleHandler) GetSchedule(w http.ResponseWriter, r *http.Request) {
//	employeeID, err := strconv.Atoi(chi.URLParam(r, "employeeID"))
//	if err != nil {
//		response.Error(w, http.StatusBadRequest, "invalid employee ID")
//		return
//	}
//
//	// Parse date range from query parameters
//	startDate := r.URL.Query().Get("start_date")
//	endDate := r.URL.Query().Get("end_date")
//
//	if startDate == "" || endDate == "" {
//		response.Error(w, http.StatusBadRequest, "start_date and end_date are required")
//		return
//	}
//
//	start, err := time.Parse("2006-01-02", startDate)
//	if err != nil {
//		response.Error(w, http.StatusBadRequest, "invalid start_date format")
//		return
//	}
//
//	end, err := time.Parse("2006-01-02", endDate)
//	if err != nil {
//		response.Error(w, http.StatusBadRequest, "invalid end_date format")
//		return
//	}
//
//	if end.Before(start) {
//		response.Error(w, http.StatusBadRequest, "end_date must be after start_date")
//		return
//	}
//
//	// Limit the date range to prevent excessive data retrieval
//	maxDays := 31
//	if end.Sub(start).Hours()/24 > float64(maxDays) {
//		response.Error(w, http.StatusBadRequest, fmt.Sprintf("date range cannot exceed %d days", maxDays))
//		return
//	}
//
//	schedule, err := h.scheduleService.GetSchedule(r.Context(), employeeID, start, end)
//	if err != nil {
//		response.Error(w, http.StatusInternalServerError, "failed to get schedule")
//		return
//	}
//
//	response.JSON(w, http.StatusOK, schedule)
//}
