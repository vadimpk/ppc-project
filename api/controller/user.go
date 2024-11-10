package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/vadimpk/ppc-project/controller/response"
	"github.com/vadimpk/ppc-project/entity"
	"github.com/vadimpk/ppc-project/pkg/auth"
	"github.com/vadimpk/ppc-project/services"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userService  services.UserService
	tokenManager *auth.TokenManager
}

func NewUserHandler(service services.UserService, manager *auth.TokenManager) *UserHandler {
	return &UserHandler{
		userService:  service,
		tokenManager: manager,
	}
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {

}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {

}

type RegisterRequest struct {
	BusinessName string `json:"business_name,omitempty"` // only for business registration
	Email        string `json:"email,omitempty"`
	Phone        string `json:"phone,omitempty"`
	FullName     string `json:"full_name"`
	Password     string `json:"password"`
}

type LoginRequest struct {
	BusinessID int    `json:"business_id"`
	Email      string `json:"email,omitempty"`
	Phone      string `json:"phone,omitempty"`
	Password   string `json:"password"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  entity.User `json:"user"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate request
	if req.FullName == "" || req.Password == "" {
		response.Error(w, http.StatusBadRequest, "full name and password are required")
		return
	}
	if req.Email == "" && req.Phone == "" {
		response.Error(w, http.StatusBadRequest, "either email or phone is required")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	var user *entity.User
	if req.BusinessName != "" {
		// Register business admin
		user, err = h.userService.CreateBusinessAdmin(r.Context(), req.BusinessName, &entity.User{
			Email:        &req.Email,
			Phone:        &req.Phone,
			FullName:     req.FullName,
			PasswordHash: string(hashedPassword),
			Role:         entity.RoleAdmin,
		})
	} else {
		// Register regular user
		user, err = h.userService.Create(r.Context(), &entity.User{
			Email:        &req.Email,
			Phone:        &req.Phone,
			FullName:     req.FullName,
			PasswordHash: string(hashedPassword),
			Role:         entity.RoleClient,
		})
	}

	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			response.Error(w, http.StatusConflict, "user already exists")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	// Generate token
	token, err := h.tokenManager.GenerateToken(user.ID, user.BusinessID, user.Role)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	response.JSON(w, http.StatusCreated, AuthResponse{
		Token: token,
		User:  *user,
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate request
	if req.BusinessID == 0 || req.Password == "" {
		response.Error(w, http.StatusBadRequest, "business_id and password are required")
		return
	}
	if req.Email == "" && req.Phone == "" {
		response.Error(w, http.StatusBadRequest, "either email or phone is required")
		return
	}

	// Get user
	var user *entity.User
	var err error
	if req.Email != "" {
		user, err = h.userService.GetByEmail(r.Context(), req.BusinessID, req.Email)
	} else {
		user, err = h.userService.GetByPhone(r.Context(), req.BusinessID, req.Phone)
	}

	if err != nil {
		response.Error(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		response.Error(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	// Generate token
	token, err := h.tokenManager.GenerateToken(user.ID, user.BusinessID, user.Role)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	response.JSON(w, http.StatusOK, AuthResponse{
		Token: token,
		User:  *user,
	})
}
