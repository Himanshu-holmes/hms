package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/himanshu-holmes/hms/internal/model"
	
)

type AuthService interface {
	Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error)
	CreateUser(ctx context.Context, req model.UserCreateRequest) (*model.User, error)
}

type PatientService interface {
	RegisterPatient(ctx context.Context, req model.ParsedPatientRequest, registeredByUserID uuid.UUID) (*model.Patient, error)
	GetPatientDetails(ctx context.Context, patientID uuid.UUID) (*model.Patient, error)
	ListPatients(ctx context.Context, params model.PaginationParams) ([]model.Patient, int64, error)
	UpdatePatientDetails(context.Context, uuid.UUID, model.ParsedPatientRequest, model.UserRole, uuid.UUID) (*model.Patient, error)
	DeletePatientRecord(ctx context.Context, patientID uuid.UUID, deletedByUserID uuid.UUID) error
}

type PatientVisitService interface {
	RecordPatientVisit(ctx context.Context, req model.ParsedPatientVisitRequest) (*model.PatientVisit, error)
	GetPatientVisitDetails(ctx context.Context, visitID uuid.UUID) (*model.PatientVisit, error)
	ListPatientVisits(ctx context.Context, patientID uuid.UUID, params model.PaginationParams) ([]model.PatientVisit, int64, error)
	UpdatePatientVisit(ctx context.Context, visitID uuid.UUID, req model.ParsedPatientVisitRequest) (*model.PatientVisit, error) // DoctorID is in ParsedPatientVisitRequest
}
