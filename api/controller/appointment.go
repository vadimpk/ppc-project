package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/vadimpk/ppc-project/controller/middleware"
	"github.com/vadimpk/ppc-project/controller/response"
	"github.com/vadimpk/ppc-project/entity"
	"github.com/vadimpk/ppc-project/services"
)

type AppointmentHandler struct {
	appointmentService services.AppointmentService
}

func NewAppointmentHandler(service services.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentService: service,
	}
}

type CreateAppointmentRequest struct {
	ClientID     int       `json:"client_id"`
	EmployeeID   int       `json:"employee_id"`
	ServiceID    int       `json:"service_id"`
	StartTime    time.Time `json:"start_time"`
	ReminderTime *int      `json:"reminder_time,omitempty"`
}

type UpdateAppointmentRequest struct {
	StartTime    time.Time `json:"start_time"`
	ReminderTime *int      `json:"reminder_time,omitempty"`
}

type GetAvailableSlotsQuery struct {
	EmployeeID int       `json:"employee_id"`
	ServiceID  int       `json:"service_id"`
	Date       time.Time `json:"date"`
}

func (h *AppointmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Get business ID from URL parameters
	businessID, err := strconv.Atoi(chi.URLParam(r, "businessID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid business ID")
		return
	}

	var req CreateAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// If client ID is not provided, use the authenticated user's ID
	if req.ClientID == 0 {
		userID, ok := middleware.GetUserID(r.Context())
		if !ok {
			response.Error(w, http.StatusBadRequest, "client ID is required")
			return
		}
		req.ClientID = userID
	} else {
		// If client ID is provided, only admins can create appointments for other users
		userRole, _ := middleware.GetRole(r.Context())
		if userRole != entity.RoleAdmin {
			response.Error(w, http.StatusForbidden, "unauthorized to create appointments for other users")
			return
		}
	}

	appointment := &entity.Appointment{
		BusinessID:   businessID,
		ClientID:     req.ClientID,
		EmployeeID:   req.EmployeeID,
		ServiceID:    req.ServiceID,
		StartTime:    req.StartTime,
		ReminderTime: req.ReminderTime,
	}

	if err := h.appointmentService.Create(r.Context(), appointment); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, appointment)
}

func (h *AppointmentHandler) Get(w http.ResponseWriter, r *http.Request) {
	appointmentID, err := strconv.Atoi(chi.URLParam(r, "appointmentID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid appointment ID")
		return
	}

	appointment, err := h.appointmentService.Get(r.Context(), appointmentID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get appointment")
		return
	}

	// Verify access rights
	userID, _ := middleware.GetUserID(r.Context())
	userRole, _ := middleware.GetRole(r.Context())
	businessID, _ := middleware.GetBusinessID(r.Context())

	if userRole != entity.RoleAdmin &&
		appointment.BusinessID != businessID &&
		appointment.ClientID != userID {
		response.Error(w, http.StatusForbidden, "unauthorized")
		return
	}

	response.JSON(w, http.StatusOK, appointment)
}

func (h *AppointmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	appointmentID, err := strconv.Atoi(chi.URLParam(r, "appointmentID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid appointment ID")
		return
	}

	var req UpdateAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Get existing appointment
	existing, err := h.appointmentService.Get(r.Context(), appointmentID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get appointment")
		return
	}

	// Verify access rights
	userRole, _ := middleware.GetRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	if userRole != entity.RoleAdmin && existing.ClientID != userID {
		response.Error(w, http.StatusForbidden, "unauthorized")
		return
	}

	// Update only allowed fields
	existing.StartTime = req.StartTime
	existing.ReminderTime = req.ReminderTime

	if err := h.appointmentService.Update(r.Context(), existing); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, existing)
}

func (h *AppointmentHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	appointmentID, err := strconv.Atoi(chi.URLParam(r, "appointmentID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid appointment ID")
		return
	}

	// Get existing appointment to check permissions
	existing, err := h.appointmentService.Get(r.Context(), appointmentID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get appointment")
		return
	}

	// Verify access rights
	userRole, _ := middleware.GetRole(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	if userRole != entity.RoleAdmin && existing.ClientID != userID {
		response.Error(w, http.StatusForbidden, "unauthorized")
		return
	}

	if err := h.appointmentService.Cancel(r.Context(), appointmentID); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"status": "cancelled"})
}

func (h *AppointmentHandler) ListByBusiness(w http.ResponseWriter, r *http.Request) {
	businessID, err := strconv.Atoi(chi.URLParam(r, "businessID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid business ID")
		return
	}

	// Verify business context
	currentBusinessID, _ := middleware.GetBusinessID(r.Context())
	if businessID != currentBusinessID {
		response.Error(w, http.StatusForbidden, "unauthorized")
		return
	}

	startTime, endTime, err := parseDateRangeQuery(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	appointments, err := h.appointmentService.ListByBusiness(r.Context(), businessID, startTime, endTime)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list appointments")
		return
	}

	response.JSON(w, http.StatusOK, appointments)
}

func (h *AppointmentHandler) ListByClient(w http.ResponseWriter, r *http.Request) {
	clientID, err := strconv.Atoi(chi.URLParam(r, "userID")) // from /users/{userID}/appointments
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	// Verify access rights
	userID, _ := middleware.GetUserID(r.Context())
	userRole, _ := middleware.GetRole(r.Context())

	if userRole != entity.RoleAdmin && clientID != userID {
		response.Error(w, http.StatusForbidden, "unauthorized")
		return
	}

	startTime, endTime, err := parseDateRangeQuery(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	appointments, err := h.appointmentService.ListByClient(r.Context(), clientID, startTime, endTime)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list appointments")
		return
	}

	response.JSON(w, http.StatusOK, appointments)
}

func (h *AppointmentHandler) GetAvailableSlots(w http.ResponseWriter, r *http.Request) {
	employeeID, err := strconv.Atoi(r.URL.Query().Get("employee_id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid employee ID")
		return
	}

	serviceID, err := strconv.Atoi(r.URL.Query().Get("service_id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid service ID")
		return
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		response.Error(w, http.StatusBadRequest, "date is required")
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid date format")
		return
	}

	slots, err := h.appointmentService.GetAvailableSlots(r.Context(), employeeID, serviceID, date)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, slots)
}

// Helper function to parse date range from query parameters
func parseDateRangeQuery(r *http.Request) (time.Time, time.Time, error) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if startDate == "" || endDate == "" {
		return time.Time{}, time.Time{}, fmt.Errorf("start_date and end_date are required")
	}

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start_date format")
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end_date format")
	}

	if end.Before(start) {
		return time.Time{}, time.Time{}, fmt.Errorf("end_date must be after start_date")
	}

	return start, end.Add(24*time.Hour - time.Second), nil
}
