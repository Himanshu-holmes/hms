package repository

import (
	"context"

	"github.com/himanshu-holmes/hms/internal/db"
	
	"github.com/jackc/pgx/v5/pgtype"
)

// UserRepository defines the interface for user data persistence.
type UserRepository interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetUserByID(ctx context.Context, id pgtype.UUID) (db.User, error)
	GetUserByUsername(ctx context.Context, username string) (db.User, error)
	ListUsers(ctx context.Context, arg db.ListUsersParams) ([]db.User, error)
	UpdateUser(ctx context.Context, arg db.UpdateUserParams) (db.User, error)
	SetUserActiveStatus(ctx context.Context, arg db.SetUserActiveStatusParams) (db.User, error)
	DeleteUser(ctx context.Context, id pgtype.UUID) error
}

// PatientRepository defines the interface for patient data persistence.
type PatientRepository interface {
	CreatePatient(ctx context.Context, arg db.CreatePatientParams) (db.Patient, error)
	GetPatientByID(ctx context.Context, id pgtype.UUID) (db.Patient, error)
	ListPatients(ctx context.Context,arg db.ListPatientsParams) ([]db.Patient, error)
	UpdatePatient(ctx context.Context, arg db.UpdatePatientParams) (db.Patient, error)
	UpdatePatientMedicalInfo(ctx context.Context, arg db.UpdatePatientMedicalInfoParams) (db.Patient, error)
	SoftDeletePatient(ctx context.Context, id pgtype.UUID) (db.Patient, error)
	HardDeletePatient(ctx context.Context, id pgtype.UUID) error
	CountPatients(ctx context.Context) (int64, error)
}

// PatientVisitRepository defines the interface for patient visit data persistence.
type PatientVisitQuerier interface {
	CreatePatientVisit(ctx context.Context, arg db.CreatePatientVisitParams) (db.PatientVisit, error)
	DeletePatientVisit(ctx context.Context, id pgtype.UUID) error
	GetPatientVisitByID(ctx context.Context, id pgtype.UUID) (db.PatientVisit, error)
	ListPatientVisitsByDoctorID(ctx context.Context, arg db.ListPatientVisitsByDoctorIDParams) ([]db.ListPatientVisitsByDoctorIDRow, error)
	ListPatientVisitsByPatientID(ctx context.Context, arg db.ListPatientVisitsByPatientIDParams) ([]db.ListPatientVisitsByPatientIDRow, error)
	UpdatePatientVisit(ctx context.Context, arg db.UpdatePatientVisitParams) (db.PatientVisit, error)
}