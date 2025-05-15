package model

import (
	"github.com/google/uuid"
)

type UserRole string

const (
	RoleReceptionist UserRole = "receptionist"
	RoleDoctor       UserRole = "doctor"
)

type User struct {
	// ID           uuid.UUID `json:"id" db:"id"`
	// Username     string    `json:"username" db:"username"`
	// PasswordHash string    `json:"-" db:"password_hash"` // Omit from JSON responses
	// Role         UserRole  `json:"role" db:"role"`
	// CreatedAt    time.Time `json:"created_at" db:"created_at"`
	// UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
    ID uuid.UUID `json:"id"`
	Username string `json:"username"`
	
	
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}