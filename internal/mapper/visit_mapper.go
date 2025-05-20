package mapper

import (
	"github.com/google/uuid"
	"github.com/himanshu-holmes/hms/internal/db"
	"github.com/himanshu-holmes/hms/internal/model"
)

// MapPatientVisit maps a *PatientVisit (pgtype) to model.PatientVisit (standard types)
func MapPatientVisit(pv *db.PatientVisit) (*model.PatientVisit, error) {
    // Handle NULL timestamps and UUIDs
    visitDate := pv.VisitDate.Time
    createdAt := pv.CreatedAt.Time
    updatedAt := pv.UpdatedAt.Time

	if pv.VisitDate.Valid {
		visitDate = pv.VisitDate.Time
	}
	if pv.CreatedAt.Valid {
		createdAt = pv.CreatedAt.Time
	}
	if pv.UpdatedAt.Valid {
		updatedAt = pv.UpdatedAt.Time
	}

    // Optional text fields
    var symptoms, diagnosis, prescription, notes *string
    if pv.Symptoms.Valid {
        symptoms = &pv.Symptoms.String
    }
    if pv.Diagnosis.Valid {
        diagnosis = &pv.Diagnosis.String
    }
    if pv.Prescription.Valid {
        prescription = &pv.Prescription.String
    }
    if pv.Notes.Valid {
        notes = &pv.Notes.String
    }

    return &model.PatientVisit{
        ID:           uuid.UUID(pv.ID.Bytes),
        PatientID:    uuid.UUID(pv.PatientID.Bytes),
        DoctorID:     uuid.UUID(pv.DoctorID.Bytes),
        VisitDate:    visitDate,
        Symptoms:     symptoms,
        Diagnosis:    diagnosis,
        Prescription: prescription,
        Notes:        notes,
        CreatedAt:    createdAt,
        UpdatedAt:    updatedAt,
    }, nil
}