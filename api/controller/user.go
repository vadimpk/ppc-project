package controller

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/vadimpk/ppc-project/controller/middleware"
	"net/http"
	"strconv"
	"strings"

	"github.com/vadimpk/ppc-project/controller/response"
	"github.com/vadimpk/ppc-project/entity"
	"github.com/vadimpk/ppc-project/pkg/auth"
	"github.com/vadimpk/ppc-project/services"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userService     services.UserService
	employeeService services.EmployeeService
	tokenManager    *auth.TokenManager
}

func NewUserHandler(service services.UserService, employeeService services.EmployeeService, manager *auth.TokenManager) *UserHandler {
	return &UserHandler{
		userService:     service,
		employeeService: employeeService,
		tokenManager:    manager,
	}
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {

}

type UpdateUserRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	FullName string `json:"full_name"`
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID")) // from /users/{userID}/appointments
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Phone == "" || req.FullName == "" {
		response.Error(w, http.StatusBadRequest, "email, phone, and full name are required")
		return
	}

	actualUserID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, http.StatusInternalServerError, "failed to get user ID")
		return
	}

	if actualUserID != userID {
		response.Error(w, http.StatusForbidden, "forbidden")
		return
	}

	user, err := h.userService.Get(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get user")
		return
	}

	user.Email = &req.Email
	user.Phone = &req.Phone
	user.FullName = req.FullName

	user, err = h.userService.Update(r.Context(), user)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to update user")
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

type RegisterRequest struct {
	BusinessName string `json:"business_name,omitempty"` // only for business registration
	BusinessID   int    `json:"business_id,omitempty"`   // only for employee registration
	Email        string `json:"email,omitempty"`
	Phone        string `json:"phone,omitempty"`
	FullName     string `json:"full_name"`
	Password     string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  entity.User `json:"user"`
}

func (h *UserHandler) RegisterBusiness(w http.ResponseWriter, r *http.Request) {
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
	} else if req.BusinessID != 0 {
		// Register regular user
		user, err = h.userService.Create(r.Context(), &entity.User{
			BusinessID:   req.BusinessID,
			Email:        &req.Email,
			Phone:        &req.Phone,
			FullName:     req.FullName,
			PasswordHash: string(hashedPassword),
			Role:         entity.RoleEmployee,
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

	if user.Role == entity.RoleEmployee {
		employeeID, err := h.employeeService.GetIDByUserID(r.Context(), user.ID)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "failed to get employee")
			return
		}
		user.EmployeeID = &employeeID
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
	if req.Password == "" {
		response.Error(w, http.StatusBadRequest, "password is required")
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
		user, err = h.userService.GetByEmail(r.Context(), req.Email)
	} else {
		user, err = h.userService.GetByPhone(r.Context(), req.Phone)
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

	if user.Role == entity.RoleEmployee {
		employeeID, err := h.employeeService.GetIDByUserID(r.Context(), user.ID)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "failed to get employee")
			return
		}
		user.EmployeeID = &employeeID
	}

	response.JSON(w, http.StatusOK, AuthResponse{
		Token: token,
		User:  *user,
	})
}
