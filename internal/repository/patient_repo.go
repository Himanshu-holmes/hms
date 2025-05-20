package repository

import (
	"context"

	"github.com/himanshu-holmes/hms/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type patientRepo struct {
	queries *db.Queries
}

func NewPatientRepo(queries *db.Queries) PatientRepository {
	return &patientRepo{queries: queries}
}

func (r *patientRepo) CreatePatient(ctx context.Context, arg db.CreatePatientParams) (db.Patient, error) {
	return r.queries.CreatePatient(ctx, arg)
}

func (r *patientRepo) GetPatientByID(ctx context.Context, id pgtype.UUID) (db.Patient, error) {
	return r.queries.GetPatientByID(ctx, id)
}

func (r *patientRepo) ListPatients(ctx context.Context, arg db.ListPatientsParams) ([]db.Patient, error) {
	return r.queries.ListPatients(ctx, arg)
}

func (r *patientRepo) UpdatePatient(ctx context.Context, arg db.UpdatePatientParams) (db.Patient, error) {
	return r.queries.UpdatePatient(ctx, arg)
}

func (r *patientRepo) UpdatePatientMedicalInfo(ctx context.Context, arg db.UpdatePatientMedicalInfoParams) (db.Patient, error) {
	return r.queries.UpdatePatientMedicalInfo(ctx, arg)
}

func (r *patientRepo) SoftDeletePatient(ctx context.Context, id pgtype.UUID) (db.Patient, error) {
	return r.queries.SoftDeletePatient(ctx, id)
}

func (r *patientRepo) HardDeletePatient(ctx context.Context, id pgtype.UUID) error {
	return r.queries.HardDeletePatient(ctx, id)
}

func (r *patientRepo) CountPatients(ctx context.Context) (int64, error) {
	return r.queries.CountPatients(ctx)
}

