package controller

import (
	"net/http"

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

func (h *AppointmentHandler) Create(w http.ResponseWriter, r *http.Request) {

}

func (h *AppointmentHandler) Get(w http.ResponseWriter, r *http.Request) {

}

func (h *AppointmentHandler) Update(w http.ResponseWriter, r *http.Request) {

}

func (h *AppointmentHandler) Cancel(w http.ResponseWriter, r *http.Request) {

}

func (h *AppointmentHandler) ListByBusiness(w http.ResponseWriter, r *http.Request) {

}

func (h *AppointmentHandler) ListByClient(w http.ResponseWriter, r *http.Request) {

}

func (h *AppointmentHandler) GetAvailableSlots(w http.ResponseWriter, r *http.Request) {

}
