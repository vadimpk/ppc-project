package controller

import (
	"net/http"

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

func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {

}

func (h *EmployeeHandler) Get(w http.ResponseWriter, r *http.Request) {

}

func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {

}

func (h *EmployeeHandler) List(w http.ResponseWriter, r *http.Request) {

}

func (h *EmployeeHandler) ListServices(w http.ResponseWriter, r *http.Request) {

}

func (h *EmployeeHandler) AssignServices(w http.ResponseWriter, r *http.Request) {

}

func (h *EmployeeHandler) RemoveServices(w http.ResponseWriter, r *http.Request) {

}
