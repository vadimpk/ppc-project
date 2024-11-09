package controller

import (
	"net/http"

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

func (h *ScheduleHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) {

}

func (h *ScheduleHandler) UpdateTemplate(w http.ResponseWriter, r *http.Request) {

}

func (h *ScheduleHandler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {

}

func (h *ScheduleHandler) ListTemplates(w http.ResponseWriter, r *http.Request) {

}

func (h *ScheduleHandler) CreateOverride(w http.ResponseWriter, r *http.Request) {

}

func (h *ScheduleHandler) UpdateOverride(w http.ResponseWriter, r *http.Request) {

}

func (h *ScheduleHandler) DeleteOverride(w http.ResponseWriter, r *http.Request) {

}

func (h *ScheduleHandler) ListOverrides(w http.ResponseWriter, r *http.Request) {

}
