package model

import (
	"time"

	"github.com/google/uuid"
	
)

// PatientVisit represents a patient's visit record.
type PatientVisit struct {
	ID           uuid.UUID `json:"id"`
	PatientID    uuid.UUID `json:"patient_id"`
	DoctorID     uuid.UUID `json:"doctor_id"` // The user ID of the doctor who conducted the visit
	VisitDate    time.Time `json:"visit_date"`    // Timestamp of the visit
	Symptoms     *string   `json:"symptoms,omitempty"`
	Diagnosis    *string   `json:"diagnosis,omitempty"`
	Prescription *string   `json:"prescription,omitempty"`
	Notes        *string   `json:"notes,omitempty"` // Additional notes by doctor or about the visit
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// PatientVisitCreateRequest is used for recording a new patient visit.
// DoctorID is typically derived from the authenticated user's token, not sent in request body.
type PatientVisitCreateRequest struct {
	PatientID    uuid.UUID `json:"patient_id" validate:"required"`
	VisitDateStr string    `json:"visit_date,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"` // Expect RFC3339. If empty, handler might default to time.Now()
	Symptoms     *string   `json:"symptoms,omitempty" validate:"omitempty"`
	Diagnosis    *string   `json:"diagnosis,omitempty" validate:"omitempty"`
	Prescription *string   `json:"prescription,omitempty" validate:"omitempty"`
	Notes        *string   `json:"notes,omitempty" validate:"omitempty"`
}

// PatientVisitUpdateRequest is used for updating an existing patient visit.
// PatientID and DoctorID are generally not updatable for an existing visit.
type PatientVisitUpdateRequest struct {
	VisitDateStr *string `json:"visit_date,omitempty" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Symptoms     *string `json:"symptoms,omitempty" validate:"omitempty"`
	Diagnosis    *string `json:"diagnosis,omitempty" validate:"omitempty"`
	Prescription *string `json:"prescription,omitempty" validate:"omitempty"`
	Notes        *string `json:"notes,omitempty" validate:"omitempty"`
}

// ParsedPatientVisitRequest is an intermediate struct for services after parsing dates
// and potentially adding context like DoctorID.
type ParsedPatientVisitRequest struct {
	PatientID    uuid.UUID // From request for Create, fixed for Update
	DoctorID     uuid.UUID // From authenticated user context
	VisitDate    time.Time // Parsed from VisitDateStr
	Symptoms     *string
	Diagnosis    *string
	Prescription *string
	Notes        *string
}