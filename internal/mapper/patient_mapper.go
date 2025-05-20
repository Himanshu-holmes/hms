package mapper

import (
	"github.com/himanshu-holmes/hms/internal/db"
	"github.com/himanshu-holmes/hms/internal/model"
	"github.com/google/uuid"
)

// ConvertDBPatientToModel maps db.CreatePatientParams to model.Patient
func ConvertDBPatientToModel(p *db.Patient) model.Patient {
	// Handle pgtype.Text and other nullable fields safely
	var contactPhone *string
	if p.ContactPhone.Valid {
		contactPhone = &p.ContactPhone.String
	}

	var contactEmail *string
	if p.ContactEmail.Valid {
		contactEmail = &p.ContactEmail.String
	}

	var address *string
	if p.Address.Valid {
		address = &p.Address.String
	}

	var medicalHistory *string
	if p.MedicalHistory.Valid {
		medicalHistory = &p.MedicalHistory.String
	}

	var gender *model.Gender
	if p.Gender.Valid {
		g := model.Gender(p.Gender.GenderEnum)
		gender = &g
	}

	return model.Patient{
		ID:                 p.ID.Bytes,
		FirstName:          p.FirstName,
		LastName:           p.LastName,
		DateOfBirth:        p.DateOfBirth.Time,
		Gender:             gender,
		ContactPhone:       contactPhone,
		ContactEmail:       contactEmail,
		Address:            address,
		MedicalHistory:     medicalHistory,
		RegisteredByUserID: uuid.UUID(p.RegisteredByUserID.Bytes),
		CreatedAt:          p.CreatedAt.Time,
		UpdatedAt:          p.UpdatedAt.Time,
		DeletedAt:          nil,
	}
}
