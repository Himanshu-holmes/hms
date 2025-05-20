package repository

import (
	"context"

	"github.com/himanshu-holmes/hms/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type patientVisitQuerierRepo struct {
	queries *db.Queries
}

func NewPatientVisitRepo(queries *db.Queries) *patientVisitQuerierRepo {
	return &patientVisitQuerierRepo{queries: queries}
}

func (r *patientVisitQuerierRepo) CreatePatientVisit(ctx context.Context, arg db.CreatePatientVisitParams) (db.PatientVisit, error) {
	return r.queries.CreatePatientVisit(ctx, arg)
}

func (r *patientVisitQuerierRepo) GetPatientVisitByID(ctx context.Context, id pgtype.UUID) (db.PatientVisit, error) {
	return r.queries.GetPatientVisitByID(ctx, id)
}

func (r *patientVisitQuerierRepo) ListPatientVisitsByPatientID(ctx context.Context, arg db.ListPatientVisitsByPatientIDParams) ([]db.ListPatientVisitsByPatientIDRow, error) {
	return r.queries.ListPatientVisitsByPatientID(ctx, arg)
}

func (r *patientVisitQuerierRepo) ListPatientVisitsByDoctorID(ctx context.Context, arg db.ListPatientVisitsByDoctorIDParams) ([]db.ListPatientVisitsByDoctorIDRow, error) {
	return r.queries.ListPatientVisitsByDoctorID(ctx, arg)
}

func (r *patientVisitQuerierRepo) UpdatePatientVisit(ctx context.Context, arg db.UpdatePatientVisitParams) (db.PatientVisit, error) {
	return r.queries.UpdatePatientVisit(ctx, arg)
}
func (r *patientVisitQuerierRepo) DeletePatientVisit(ctx context.Context, id pgtype.UUID) error {
	return r.queries.DeletePatientVisit(ctx, id)
}



