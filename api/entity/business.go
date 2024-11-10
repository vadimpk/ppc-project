package entity

import "time"

type Business struct {
	ID          int                    `json:"id" db:"id"`
	Name        string                 `json:"name" db:"name"`
	LogoURL     *string                `json:"logo_url" db:"logo_url"`
	ColorScheme map[string]interface{} `json:"color_scheme" db:"color_scheme"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
}

type BusinessService struct {
	ID          int       `json:"id" db:"id"`
	BusinessID  int       `json:"business_id" db:"business_id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	Duration    int       `json:"duration" db:"duration"`
	Price       int       `json:"price" db:"price"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Employee struct {
	ID             int       `json:"id" db:"id"`
	BusinessID     int       `json:"business_id" db:"business_id"`
	UserID         int       `json:"user_id" db:"user_id"`
	Specialization *string   `json:"specialization" db:"specialization"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	User           *User     `json:"user" db:"user"`
}

type EmployeeService struct {
	EmployeeID int `json:"employee_id" db:"employee_id"`
	ServiceID  int `json:"service_id" db:"service_id"`
}
