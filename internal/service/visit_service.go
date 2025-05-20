package service

import (
	"context"
	"errors"
	"fmt"

	// "go/format"
	"log"
	"strings"

	// "time"
	// Not directly used here, but ParsedPatientVisitRequest has time.Time
	"github.com/google/uuid"
	"github.com/himanshu-holmes/hms/internal/db"
	"github.com/himanshu-holmes/hms/internal/mapper"
	"github.com/himanshu-holmes/hms/internal/model"
	"github.com/himanshu-holmes/hms/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var ErrVisitNotFound = errors.New("patient visit not found")
var ErrVisitUpdateForbidden = errors.New("not authorized to update this visit record")
var ErrPatientForVisitNotFound = errors.New("patient for visit not found")

type patientVisitService struct {
	visitRepo   repository.PatientVisitQuerier
	patientRepo repository.PatientRepository // To check if patient exists
}

func NewPatientVisitService(visitRepo repository.PatientVisitQuerier, patientRepo repository.PatientRepository) PatientVisitService {
	return &patientVisitService{visitRepo: visitRepo, patientRepo: patientRepo}
}

func derefString(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func (s *patientVisitService) RecordPatientVisit(ctx context.Context, req model.ParsedPatientVisitRequest) (*model.PatientVisit, error) {
// Check if patient exists
	_, err := s.patientRepo.GetPatientByID(ctx, pgtype.UUID{Bytes: [16]byte(req.PatientID), Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows)  {
			return nil, fmt.Errorf("%w: patient ID %s", ErrPatientForVisitNotFound, req.PatientID)
		}
		log.Printf("VisitService: Error checking patient %s for new visit: %v", req.PatientID, err)
	}
 
	visitParams  := &db.CreatePatientVisitParams{
		PatientID: pgtype.UUID{Bytes: [16]byte(req.PatientID), Valid: true},
		DoctorID:  pgtype.UUID{Bytes: [16]byte(req.DoctorID), Valid: true},
		VisitDate: pgtype.Timestamptz{Time: req.VisitDate, Valid: true},
		Symptoms:  pgtype.Text{String: derefString(req.Symptoms), Valid: req.Symptoms != nil},
		Diagnosis: pgtype.Text{String: derefString(req.Diagnosis), Valid: req.Diagnosis != nil},
		Prescription: pgtype.Text{String: derefString(req.Prescription), Valid: req.Prescription != nil},
		Notes: pgtype.Text{String: derefString(req.Notes), Valid: req.Notes != nil},

	}
	

	visit, err := s.visitRepo.CreatePatientVisit(ctx, *visitParams)
	if err != nil {
		// Handle potential foreign key constraint errors if patient_id or doctor_id is invalid at DB level
		if strings.Contains(strings.ToLower(err.Error()), "foreign key constraint") {
			log.Printf("VisitService: Foreign key constraint error creating visit for patient %s by doctor %s: %v", req.PatientID, req.DoctorID, err)
			return nil, fmt.Errorf("invalid patient or doctor reference for visit")
		}
		log.Printf("VisitService: Failed to record patient visit in repo: %v", err)
		return nil, fmt.Errorf("failed to record patient visit: %w", err)
	}
	

	formattedVisit,err := mapper.MapPatientVisit(&visit)
	
	return formattedVisit, nil
}

func (s *patientVisitService) GetPatientVisitDetails(ctx context.Context, visitID uuid.UUID) (*model.PatientVisit, error) {
	visit, err := s.visitRepo.GetPatientVisitByID(ctx, pgtype.UUID{Bytes: visitID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrVisitNotFound
		}
		log.Printf("VisitService: Failed to get visit details for ID %s: %v", visitID, err)
		return nil, fmt.Errorf("failed to get visit details: %w", err)
	}
	formmatedVisit,err := mapper.MapPatientVisit(&visit)
	return formmatedVisit, nil
}

func (s *patientVisitService) ListPatientVisits(ctx context.Context, patientID uuid.UUID, params model.PaginationParams) ([]model.PatientVisit, int64, error) {
	// Optional: Check if patient exists first to return a 404 if patient_id is invalid
	_, err := s.patientRepo.GetPatientByID(ctx, pgtype.UUID{Bytes: [16]byte(patientID), Valid: true})
	if err != nil  {
		if errors.Is(err, pgx.ErrNoRows)  {
			log.Printf("VisitService: Attempted to list visits for non-existent patient ID %s", patientID)
			return []model.PatientVisit{}, 0, fmt.Errorf("%w: patient ID %s", ErrPatientForVisitNotFound, patientID) // Return empty slice and 0 total
		}
		log.Printf("VisitService: Error checking patient %s before listing visits: %v", patientID, err)
		return nil, 0, fmt.Errorf("error verifying patient before listing visits: %w", err)
	}


	dbParams := db.ListPatientVisitsByPatientIDParams{
		PatientID: pgtype.UUID{Bytes: patientID, Valid: true},
		Limit:     int32(params.Limit),
		Offset:    int32(params.Offset),
	}
	visits, err := s.visitRepo.ListPatientVisitsByPatientID(ctx, dbParams)
	if err != nil {
		log.Printf("VisitService: Failed to list patient visits for patient %s: %v", patientID, err)
		return nil, 0, fmt.Errorf("failed to list patient visits: %w", err)
	}

	mappedVisits := make([]model.PatientVisit, len(visits))
	for i, visit := range visits {
		patientVisit := &db.PatientVisit{
			ID:          visit.ID,
			PatientID:   visit.PatientID,
			DoctorID:    visit.DoctorID,
			VisitDate:   visit.VisitDate,
			Symptoms:    visit.Symptoms,
			Diagnosis:   visit.Diagnosis,
			Prescription: visit.Prescription,
			Notes:       visit.Notes,
		}
		mappedVisit, err := mapper.MapPatientVisit(patientVisit)
		if err != nil {
			log.Printf("VisitService: Failed to map visit for patient %s: %v", patientID, err)
			return nil, 0, fmt.Errorf("failed to map patient visits: %w", err)
		}
		mappedVisits[i] = *mappedVisit
	}
	return mappedVisits, int64(len(mappedVisits)), nil
}

func (s *patientVisitService) UpdatePatientVisit(ctx context.Context, visitID uuid.UUID, req model.ParsedPatientVisitRequest) (*model.PatientVisit, error) {
	existingVisit, err := s.visitRepo.GetPatientVisitByID(ctx, pgtype.UUID{Bytes: visitID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrVisitNotFound
		}
		log.Printf("VisitService: Failed to fetch visit %s for update: %v", visitID, err)
		return nil, fmt.Errorf("failed to fetch visit for update: %w", err)
	}
	

	// Authorization: Ensure the authenticated doctor (from req.DoctorID) is the one who created this visit.
	// req.DoctorID is set by the handler/parser from the authenticated user's token.
	if existingVisit.DoctorID.Bytes != req.DoctorID {
		log.Printf("VisitService: Unauthorized attempt to update visit %s. Visit DoctorID: %s, Requester DoctorID: %s", visitID, existingVisit.DoctorID, req.DoctorID)
		return nil, ErrVisitUpdateForbidden
	}

	changed := false
	// Apply updates from req. ParsedPatientVisitRequest will have zero values for fields not in JSON.
	if !req.VisitDate.IsZero() && !existingVisit.VisitDate.Time.Equal(req.VisitDate) {
		existingVisit.VisitDate = pgtype.Timestamptz{Time: req.VisitDate, Valid: true}
		changed = true
	}
	if req.Symptoms != nil && (!existingVisit.Symptoms.Valid || existingVisit.Symptoms.String != *req.Symptoms) {
		existingVisit.Symptoms = pgtype.Text{String: derefString(req.Symptoms), Valid: req.Symptoms != nil}
		changed = true
	}
	if req.Diagnosis != nil && (!existingVisit.Diagnosis.Valid || existingVisit.Diagnosis.String != *req.Diagnosis) {
		existingVisit.Diagnosis = pgtype.Text{String: *req.Diagnosis, Valid: true}
		changed = true
	}
	if req.Prescription != nil && (!existingVisit.Prescription.Valid || existingVisit.Prescription.String != *req.Prescription) {
		existingVisit.Prescription = pgtype.Text{String: *req.Prescription, Valid: true}
		changed = true
	}
	if req.Notes != nil && (!existingVisit.Notes.Valid || existingVisit.Notes.String != *req.Notes) {
		existingVisit.Notes = pgtype.Text{String: *req.Notes, Valid: true}
		changed = true
	}
	// PatientID and DoctorID of the visit itself are not updatable.

	if !changed {
		log.Printf("VisitService: No changes detected for visit %s update. Returning existing.", visitID)
		mappedVisit, err := mapper.MapPatientVisit(&existingVisit)
		if err != nil {
			log.Printf("VisitService: Failed to map existing visit %s: %v", visitID, err)
			return nil, fmt.Errorf("failed to map existing visit: %w", err)
		}
		return mappedVisit, nil
	}

	updateParams := db.UpdatePatientVisitParams{
		ID:           existingVisit.ID,
		VisitDate:    existingVisit.VisitDate,
		Symptoms:     existingVisit.Symptoms,
		Diagnosis:    existingVisit.Diagnosis,
		Prescription: existingVisit.Prescription,
		Notes:        existingVisit.Notes,
	}
	patientVisit, err := s.visitRepo.UpdatePatientVisit(ctx, updateParams) // The repository should also check DoctorID for safety
	if err != nil {
		// The repo's UpdateVisit might return an error if rows affected is 0,
		// which could mean not found or not authorized at the DB level (if DoctorID is in WHERE clause).
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no rows updated") || errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrVisitNotFound // Or potentially ErrVisitUpdateForbidden if repo implies that
		}
		log.Printf("VisitService: Failed to update patient visit %s in repo: %v", visitID, err)
		return nil, fmt.Errorf("failed to update patient visit: %w", err)
	}

	updatedVisit, err := mapper.MapPatientVisit(&patientVisit)
	return updatedVisit, nil
}