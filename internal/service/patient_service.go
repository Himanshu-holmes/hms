package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/himanshu-holmes/hms/internal/db"
	"github.com/himanshu-holmes/hms/internal/mapper"
	"github.com/himanshu-holmes/hms/internal/model"
	"github.com/himanshu-holmes/hms/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var ErrPatientNotFound = errors.New("patient not found")
var ErrPatientUpdateForbidden = errors.New("user not authorized to update this patient's medical history")
var ErrPatientConflict = errors.New("patient data conflicts with existing record")

type patientService struct {
	patientRepo repository.PatientRepository
}

func NewPatientService(patientRepo repository.PatientRepository) PatientService {
	return &patientService{patientRepo: patientRepo}
}

func (s *patientService) RegisterPatient(ctx context.Context, req model.ParsedPatientRequest, registeredByUserID uuid.UUID) (*model.Patient, error) {
	// patient := &model.Patient{
	// 	// ID will be generated by DB
	// 	FirstName:          req.FirstName,
	// 	LastName:           req.LastName,
	// 	DateOfBirth:        req.DateOfBirth,
	// 	Gender:             req.Gender,
	// 	ContactPhone:       req.ContactPhone,
	// 	ContactEmail:       req.ContactEmail,
	// 	Address:            req.Address,
	// 	MedicalHistory:     req.MedicalHistory, // Receptionist can set initial medical history
	// 	RegisteredByUserID: registeredByUserID,
	// }

	patientParams := db.CreatePatientParams{
		FirstName:          req.FirstName,
		LastName:           req.LastName,
		DateOfBirth: func() pgtype.Date {
			parsedDate, err := time.Parse("2006-01-02", req.DateOfBirth.Format("2006-01-02")) // Adjust format as per your date string
			if err != nil {
				log.Printf("PatientService: Invalid date format for DateOfBirth: %v", err)
				return pgtype.Date{Valid: false}
			}
			return pgtype.Date{Time: parsedDate, Valid: true}
		}(),
		Gender:             db.NullGenderEnum{GenderEnum: db.GenderEnum(*req.Gender), Valid: req.Gender != nil},
		ContactPhone:       pgtype.Text{String: *req.ContactPhone, Valid: req.ContactPhone != nil},
		ContactEmail:       pgtype.Text{String: *req.ContactEmail, Valid: req.ContactEmail != nil},
		Address:            pgtype.Text{String: *req.Address, Valid: req.Address != nil},
		MedicalHistory:     pgtype.Text{String: *req.MedicalHistory, Valid: req.MedicalHistory != nil}, // Receptionist can set initial medical history
		RegisteredByUserID: pgtype.UUID{Bytes: registeredByUserID, Valid: true},

	}

	patient, err := s.patientRepo.CreatePatient(ctx, patientParams)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique constraint") ||
			strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
			log.Printf("PatientService: Conflict creating patient: %v", err)
			return nil, ErrPatientConflict
		}
		log.Printf("PatientService: Failed to register patient in repo: %v", err)
		return nil, fmt.Errorf("failed to register patient: %w", err)
	}

 formattedPatient := mapper.ConvertDBPatientToModel(&patient)
	
	
	return &formattedPatient, nil
}

func (s *patientService) GetPatientDetails(ctx context.Context, patientID uuid.UUID) (*model.Patient, error) {
	convertedPatientID := pgtype.UUID{Bytes: patientID, Valid: true}
	patient, err := s.patientRepo.GetPatientByID(ctx, convertedPatientID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrPatientNotFound
		}
		log.Printf("PatientService: Failed to get patient details for ID %s: %v", patientID, err)
		return nil, fmt.Errorf("failed to get patient details: %w", err)
	}
    formattedPatient := mapper.ConvertDBPatientToModel(&patient)
	return &formattedPatient, nil
}

func (s *patientService) ListPatients(ctx context.Context, params model.PaginationParams) ([]model.Patient, int64, error) {

	patients, err := s.patientRepo.ListPatients(ctx, db.ListPatientsParams{
		Limit:  int32(params.Limit),
		Offset: int32(params.Offset),
	})
	if err != nil {
		log.Printf("PatientService: Failed to list patients: %v", err)
		return nil, 0, fmt.Errorf("failed to list patients: %w", err)
	}
    var formattedPatients []model.Patient
	for _, patient := range patients {
		formattedPatient := mapper.ConvertDBPatientToModel(&patient)
		formattedPatients = append(formattedPatients, formattedPatient)
	}

	total, err := s.patientRepo.CountPatients(ctx)
	if err != nil {
		log.Printf("PatientService: Failed to count patients: %v", err)
		return nil, 0, fmt.Errorf("failed to count patients: %w", err)
	}
	
	return formattedPatients, total, nil
}

func (s *patientService) UpdatePatientDetails(
	ctx context.Context,
	patientID uuid.UUID,
	req model.ParsedPatientRequest,
	updaterRole model.UserRole,
	updaterID uuid.UUID, // ID of the user performing the update
) (*model.Patient, error) {
	convertedPatientID := pgtype.UUID{Bytes: patientID, Valid: true}
	existingPatient, err := s.patientRepo.GetPatientByID(ctx, convertedPatientID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows)  {
			return nil, ErrPatientNotFound
		}
		log.Printf("PatientService: Failed to fetch patient %s for update: %v", patientID, err)
		return nil, fmt.Errorf("failed to fetch patient for update: %w", err)
	}
	

	// Apply updates based on request fields.
	// ParsedPatientRequest will have zero values for fields not in the JSON update request.
	// We only update if the field in 'req' is non-zero (for pointers: not nil).

	// Receptionist can update demographic info.
	// Doctor can update demographic info AND medical history.
	// This logic assumes ParsedPatientRequest fields are pointers for optional updates,
	// or uses a sentinel like non-empty string / non-zero time for required fields that are being updated.

	changed := false // Flag to track if any change was made

	if req.FirstName != "" && existingPatient.FirstName != req.FirstName {
		existingPatient.FirstName = req.FirstName
		changed = true
	}
	if req.LastName != "" && existingPatient.LastName != req.LastName {
		existingPatient.LastName = req.LastName
		changed = true
	}
	if !req.DateOfBirth.IsZero() && existingPatient.DateOfBirth.Time != req.DateOfBirth {

		existingPatient.DateOfBirth = pgtype.Date{Time: req.DateOfBirth, Valid: true}
		changed = true
	}
	if req.Gender != nil && (!existingPatient.Gender.Valid || existingPatient.Gender.GenderEnum != db.GenderEnum(*req.Gender)) {
		existingPatient.Gender = db.NullGenderEnum{
			GenderEnum: db.GenderEnum(*req.Gender),
			Valid:      req.Gender != nil,
		}
		changed = true
	}
	if req.ContactPhone != nil && (!existingPatient.ContactPhone.Valid || existingPatient.ContactPhone.String != *req.ContactPhone) {
		existingPatient.ContactPhone = pgtype.Text{String: *req.ContactPhone, Valid: req.ContactPhone != nil}
		changed = true
	}
	if req.ContactEmail != nil && (!existingPatient.ContactEmail.Valid || existingPatient.ContactEmail.String != *req.ContactEmail) {
		existingPatient.ContactEmail = pgtype.Text{String: *req.ContactEmail, Valid: req.ContactEmail != nil}
		changed = true
	}
	if req.Address != nil && (!existingPatient.Address.Valid || existingPatient.Address.String != *req.Address) {
		existingPatient.Address = pgtype.Text{String: *req.Address, Valid: req.Address != nil}
		changed = true
	}

	// Medical History: Only Doctors can update it IF IT ALREADY EXISTS AND IS DIFFERENT.
	// Receptionists can set it if it's currently NULL (during initial creation or if it was cleared).
	// This allows receptionists to fill in initial medical history but not change an existing one set by a doctor.
	if req.MedicalHistory != nil {
		isNewMH := !existingPatient.MedicalHistory.Valid
		mhChanged := !existingPatient.MedicalHistory.Valid || existingPatient.MedicalHistory.String != *req.MedicalHistory

		if mhChanged {
			if updaterRole == model.RoleDoctor {
				existingPatient.MedicalHistory = pgtype.Text{String: *req.MedicalHistory, Valid: req.MedicalHistory != nil}
				changed = true
			} else if updaterRole == model.RoleReceptionist && isNewMH { // Receptionist can set if currently null
				existingPatient.MedicalHistory = pgtype.Text{String: *req.MedicalHistory, Valid: req.MedicalHistory != nil}
				changed = true
			} else if updaterRole == model.RoleReceptionist && !isNewMH && *req.MedicalHistory != existingPatient.MedicalHistory.String {
				// Receptionist trying to change existing non-null medical history
				log.Printf("PatientService: Receptionist %s attempted to update existing medical history for patient %s.", updaterID, patientID)
				return nil, ErrPatientUpdateForbidden
			}
		}
	}

   
	if !changed {
		log.Printf("PatientService: No changes detected for patient %s update. Returning existing.", patientID)
		formattedPatient := mapper.ConvertDBPatientToModel(&existingPatient)
		return &formattedPatient, nil // No actual update needed
	}

	updateParams := db.UpdatePatientParams{
		ID:               existingPatient.ID,
		FirstName:        pgtype.Text{String: existingPatient.FirstName, Valid: existingPatient.FirstName != ""},
		LastName:         pgtype.Text{String: existingPatient.LastName, Valid: existingPatient.LastName != ""},
		DateOfBirth:      existingPatient.DateOfBirth,
		Gender:           existingPatient.Gender,
		ContactPhone:     existingPatient.ContactPhone,
		ContactEmail:     existingPatient.ContactEmail,
		Address:          existingPatient.Address,
		MedicalHistory:   existingPatient.MedicalHistory,
		 // Assuming updaterID is the user performing the update
	}
	_, err = s.patientRepo.UpdatePatient(ctx, updateParams)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique constraint") ||
			strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
			log.Printf("PatientService: Conflict updating patient %s: %v", patientID, err)
			return nil, ErrPatientConflict
		}
		log.Printf("PatientService: Failed to update patient %s in repo: %v", patientID, err)
		return nil, fmt.Errorf("failed to update patient: %w", err)
	}

	// Refetch to get updated timestamps and ensure consistency
	convertedPatientID = pgtype.UUID{Bytes: patientID, Valid: true}
	updatedPatient, err := s.patientRepo.GetPatientByID(ctx, convertedPatientID)
	if err != nil {
		log.Printf("PatientService: Error fetching updated patient %s: %v. Returning locally modified patient.", patientID, err)
		formattedPatient := mapper.ConvertDBPatientToModel(&updatedPatient)
		return &formattedPatient, nil // Return what we have if refetch fails
	}
	formattedPatient := mapper.ConvertDBPatientToModel(&updatedPatient)
	return &formattedPatient, nil
}

func (s *patientService) DeletePatientRecord(ctx context.Context, patientID uuid.UUID, deletedByUserID uuid.UUID) error {
	// First, check if patient exists to provide a better error message
	convertedPatientID := pgtype.UUID{Bytes: [16]byte(patientID), Valid: true}
	_, err := s.patientRepo.GetPatientByID(ctx, convertedPatientID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrPatientNotFound
		}
		log.Printf("PatientService: Error checking patient %s before delete: %v", patientID, err)
		return fmt.Errorf("error preparing to delete patient: %w", err)
	}

	convertedPatientID = pgtype.UUID{Bytes: patientID, Valid: true}
	_, err = s.patientRepo.SoftDeletePatient(ctx, convertedPatientID)
	if err != nil {
		// The repository DeletePatient might return specific errors for "already deleted" or "not found"
		// which could be sql.ErrNoRows or a custom error.
		if errors.Is(err, pgx.ErrNoRows) || strings.Contains(err.Error(), "already deleted") { // Adapt to repo's specific error
			return ErrPatientNotFound // Treat as "not found" if already deleted for idempotency
		}
		log.Printf("PatientService: Failed to delete patient %s by user %s: %v", patientID, deletedByUserID, err)
		return fmt.Errorf("failed to delete patient: %w", err)
	}
	return nil
}