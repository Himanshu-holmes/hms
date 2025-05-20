package model

import (
	"time"

	"github.com/google/uuid"
	
)

type UserRole string

const (
	RoleReceptionist UserRole = "receptionist"
	RoleDoctor       UserRole = "doctor"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Exclude from JSON responses, still needed internally for auth
	Role         UserRole  `json:"role"`
	FirstName    *string   `json:"first_name,omitempty"`
	LastName     *string   `json:"last_name,omitempty"`
	Email        *string   `json:"email,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserCreateRequest is used for creating users via the API.
type UserCreateRequest struct {
	Username  string   `json:"username" validate:"required,min=3,max=100"`
	Password  string   `json:"password" validate:"required,min=6"`
	Role      UserRole `json:"role" validate:"required,oneof=receptionist doctor"`
	FirstName *string  `json:"first_name,omitempty" validate:"omitempty,max=100"`
	LastName  *string  `json:"last_name,omitempty" validate:"omitempty,max=100"`
	Email     *string  `json:"email,omitempty" validate:"omitempty,email,max=255"`
}

// LoginRequest is used for user login via the API.
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse is the API response for a successful login.
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User  User   `json:"user"` // This User struct is from our model package
}