package entity

import "time"

type User struct {
	ID           int       `json:"id" db:"id"`
	BusinessID   int       `json:"business_id" db:"business_id"`
	Email        *string   `json:"email" db:"email"`
	Phone        *string   `json:"phone" db:"phone"`
	FullName     string    `json:"full_name" db:"full_name"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Role         string    `json:"role" db:"role"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

const (
	RoleAdmin    = "admin"
	RoleEmployee = "employee"
	RoleClient   = "client"
)
