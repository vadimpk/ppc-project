package controller

import (
	"net/http"

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

func (h *BusinessHandler) Create(w http.ResponseWriter, r *http.Request) {

}

func (h *BusinessHandler) Get(w http.ResponseWriter, r *http.Request) {

}

func (h *BusinessHandler) Update(w http.ResponseWriter, r *http.Request) {

}

func (h *BusinessHandler) UpdateAppearance(w http.ResponseWriter, r *http.Request) {

}
