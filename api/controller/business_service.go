package controller

import (
	"net/http"

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

func (h *BusinessServiceHandler) Create(w http.ResponseWriter, r *http.Request) {

}

func (h *BusinessServiceHandler) Get(w http.ResponseWriter, r *http.Request) {

}

func (h *BusinessServiceHandler) Update(w http.ResponseWriter, r *http.Request) {

}

func (h *BusinessServiceHandler) Delete(w http.ResponseWriter, r *http.Request) {

}

func (h *BusinessServiceHandler) List(w http.ResponseWriter, r *http.Request) {

}
