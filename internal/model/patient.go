package model

import (
	"time"

	"github.com/google/uuid"
)

// Gender defines the possible gender options for a patient.
type Gender string

const (
	GenderMale           Gender = "male"
	GenderFemale         Gender = "female"
	GenderOther          Gender = "other"
)

// Patient represents patient data for API responses and internal logic.
type Patient struct {
	ID                  uuid.UUID  `json:"id"`
	FirstName           string     `json:"first_name"`
	LastName            string     `json:"last_name"`
	DateOfBirth         time.Time  `json:"date_of_birth"` // Recommended to store as full timestamp, format for display
	Gender              *Gender    `json:"gender,omitempty"`
	ContactPhone        *string    `json:"contact_phone,omitempty"`
	ContactEmail        *string    `json:"contact_email,omitempty"`
	Address             *string    `json:"address,omitempty"`
	MedicalHistory      *string    `json:"medical_history,omitempty"`
	RegisteredByUserID  uuid.UUID  `json:"registered_by_user_id"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"-"` // For soft delete, typically excluded from normal JSON responses.
}

type PatientCreateRequest struct {
	FirstName      string  `json:"first_name" validate:"required,max=100"`
	LastName       string  `json:"last_name" validate:"required,max=100"`
	DateOfBirth string  `json:"date_of_birth" validate:"required,datetime=2006-01-02"` 
	Gender         *Gender `json:"gender,omitempty" validate:"omitempty,oneof=male female other"`
	ContactPhone   *string `json:"contact_phone,omitempty" validate:"omitempty,e164"`
	ContactEmail   *string `json:"contact_email,omitempty" validate:"omitempty,email,max=255"`
	Address        *string `json:"address,omitempty" validate:"omitempty,max=500"`
	MedicalHistory *string `json:"medical_history,omitempty"` 
}

// PatientUpdateRequest is used for updating an existing patient's details.

type PatientUpdateRequest struct {
	FirstName      *string `json:"first_name,omitempty" validate:"omitempty,max=100"`
	LastName       *string `json:"last_name,omitempty" validate:"omitempty,max=100"`
	DateOfBirthStr *string `json:"date_of_birth,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Gender         *Gender `json:"gender,omitempty" validate:"omitempty,oneof=male female other prefer_not_to_say"`
	ContactPhone   *string `json:"contact_phone,omitempty" validate:"omitempty,e164"`
	ContactEmail   *string `json:"contact_email,omitempty" validate:"omitempty,email,max=255"`
	Address        *string `json:"address,omitempty" validate:"omitempty,max=500"`
	MedicalHistory *string `json:"medical_history,omitempty"` // Doctors can update this
}

// ParsedPatientRequest is an intermediate struct used by services after parsing string dates.
// This avoids direct use of time.Time in JSON request structs which can be tricky.
type ParsedPatientRequest struct {
	FirstName      string
	LastName       string
	DateOfBirth    time.Time // Parsed from DateOfBirthStr
	Gender         *Gender
	ContactPhone   *string
	ContactEmail   *string
	Address        *string
	MedicalHistory *string
}