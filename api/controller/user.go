package controller

import (
	"net/http"

	"github.com/vadimpk/ppc-project/services"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{
		userService: service,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {

}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {

}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {

}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {

}
