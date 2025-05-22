package handler

import (
	"fmt" // For error formatting
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/himanshu-holmes/hms/internal/authorization"
	"github.com/himanshu-holmes/hms/internal/middleware"
	"github.com/himanshu-holmes/hms/internal/model"
	"github.com/himanshu-holmes/hms/internal/service"
	util "github.com/himanshu-holmes/hms/internal/utils"
)

type PatientHandler struct {
	patientService service.PatientService
}

func NewPatientHandler(patientService service.PatientService) *PatientHandler {
	return &PatientHandler{patientService: patientService}
}
// 

// parsePatientRequest handles parsing for both create and update patient requests.
// It converts string dates to time.Time.
func parsePatientRequest(c *gin.Context, req interface{}) (*model.ParsedPatientRequest, error) {
	var dob time.Time
	var err error

	parsedReq := &model.ParsedPatientRequest{}

	switch r := req.(type) {
	case model.PatientCreateRequest:
		if r.DateOfBirth != "" {
			dob, err = time.Parse("2006-01-02", r.DateOfBirth)
			if err != nil {
				return nil, fmt.Errorf("invalid date_of_birth format: %w. Expected YYYY-MM-DD", err)
			}
		}
		parsedReq.FirstName = r.FirstName
		parsedReq.LastName = r.LastName
		parsedReq.DateOfBirth = dob
		parsedReq.Gender = r.Gender
		parsedReq.ContactPhone = r.ContactPhone
		parsedReq.ContactEmail = r.ContactEmail
		parsedReq.Address = r.Address
		parsedReq.MedicalHistory = r.MedicalHistory
	case model.PatientUpdateRequest:
		if r.DateOfBirthStr != nil && *r.DateOfBirthStr != "" {
			dob, err = time.Parse("2006-01-02", *r.DateOfBirthStr)
			if err != nil {
				return nil, fmt.Errorf("invalid date_of_birth format: %w. Expected YYYY-MM-DD", err)
			}
		}
		// Only set fields if they are provided in the update request
		if r.FirstName != nil {
			parsedReq.FirstName = *r.FirstName
		}
		if r.LastName != nil {
			parsedReq.LastName = *r.LastName
		}
		if !dob.IsZero() { // Only update DOB if it was parsed successfully
			parsedReq.DateOfBirth = dob
		}
		parsedReq.Gender = r.Gender // If nil, it means no update for gender
		parsedReq.ContactPhone = r.ContactPhone
		parsedReq.ContactEmail = r.ContactEmail
		parsedReq.Address = r.Address
		parsedReq.MedicalHistory = r.MedicalHistory
	default:
		return nil, fmt.Errorf("unsupported request type for patient parsing")
	}
	return parsedReq, nil
}

// RegisterPatient godoc
// @Summary Register a new patient
// @Description Doctors can register a new patient.
// @Tags Patients
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param patientRequest body model.PatientCreateRequest true "Patient Registration Data"
// @Success 201 {object} model.Patient
// @Failure 400 {object} model.APIError "Validation error or invalid input"
// @Failure 401 {object} model.APIError "Unauthorized"
// @Failure 500 {object} model.APIError "Internal server error"
// @Router /patients/create [post]
func (h *PatientHandler) RegisterPatient(c *gin.Context) {
	var req model.PatientCreateRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{Message: "Invalid request body", Details: err.Error()})
		return
	}

	if err := util.ValidateStruct(req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{Message: "Validation failed", Details: util.FormatValidationErrors(err)})
		return
	}

	parsedReq, err := parsePatientRequest(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{Message: err.Error()}) // Error from parsePatientRequest
		return
	}

	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		// This should ideally not happen if AuthMiddleware is working
		log.Printf("CRITICAL: UserID not found in context for an authenticated route in RegisterPatient")
		c.JSON(http.StatusInternalServerError, model.APIError{Message: "User context error"})
		return
	}

	patient, err := h.patientService.RegisterPatient(c.Request.Context(), *parsedReq, userID)
	if err != nil {
		log.Printf("Register patient error: %v by user %s", err, userID)
		// Check for specific errors, e.g., duplicate contact info if your DB has unique constraints
		if strings.Contains(strings.ToLower(err.Error()), "unique constraint") || strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
			c.JSON(http.StatusConflict, model.APIError{Message: "Failed to register patient due to conflicting data (e.g., phone or email already exists)."})
		} else {
			c.JSON(http.StatusInternalServerError, model.APIError{Message: "Failed to register patient"})
		}
		return
	}

	c.JSON(http.StatusCreated, patient)
}

// GetPatient godoc
// @Summary Get details of a specific patient
// @Description Receptionists and Doctors can get details of a specific patient.
// @Tags Patients
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Patient ID"
// @Success 200 {object} model.Patient
// @Failure 400 {object} model.APIError "Invalid patient ID format"
// @Failure 404 {object} model.APIError "Patient not found"
// @Failure 401 {object} model.APIError "Unauthorized"
// @Failure 500 {object} model.APIError "Internal server error"
// @Router /patients/{id} [get]
func (h *PatientHandler) GetPatient(c *gin.Context) {
	patientIDStr := c.Param("id")
	patientID,err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{Message: "Invalid patient ID format"})
		return
	}

	patientUUID, err := patientID.Value()
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{Message: "Invalid patient ID format"})
		return
	}

	patient, err := h.patientService.GetPatientDetails(c.Request.Context(), patientUUID.(uuid.UUID))
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			c.JSON(http.StatusNotFound, model.APIError{Message: "Patient not found"})
		} else {
			log.Printf("Get patient error for ID %s: %v", patientIDStr, err)
			c.JSON(http.StatusInternalServerError, model.APIError{Message: "Failed to get patient details"})
		}
		return
	}

	c.JSON(http.StatusOK, patient)
}

// ListPatients godoc
// @Summary List all registered patients
// @Description Receptionists and Doctors can list all registered patients.
// @Tags Patients
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default: 10)" minimum(1) maximum(100)
// @Param offset query int false "Offset (default: 0)"
// @Success 200 {object} model.PaginatedResponse{data=[]model.Patient}
// @Failure 400 {object} model.APIError "Invalid pagination parameters"
// @Failure 401 {object} model.APIError "Unauthorized"
// @Failure 500 {object} model.APIError "Internal server error"
// @Router /patients [get]
func (h *PatientHandler) ListPatients(c *gin.Context) {
	var params model.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{Message: "Invalid pagination parameters", Details: err.Error()})
		return
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 { // Max limit
		params.Limit = 100
	}
	if params.Offset < 0 {
		params.Offset = 0
	}

	patients, total, err := h.patientService.ListPatients(c.Request.Context(), params)
	if err != nil {
		log.Printf("List patients error: %v", err)
		c.JSON(http.StatusInternalServerError, model.APIError{Message: "Failed to list patients"})
		return
	}

	response := model.PaginatedResponse{
		Data:   patients,
		Total:  total,
		Limit:  params.Limit,
		Offset: params.Offset,
	}

	c.JSON(http.StatusOK, response)
}

// UpdatePatient godoc
// @Summary Update patient details
// @Description Receptionists can update most patient details. Doctors can update patient details, especially medical history.
// @Tags Patients
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Patient ID (UUID)" Format(uuid)
// @Param patientRequest body model.PatientUpdateRequest true "Patient Update Data (fields to update)"
// @Success 200 {object} model.Patient
// @Failure 400 {object} model.APIError "Validation error, invalid input, or invalid patient ID"
// @Failure 401 {object} model.APIError "Unauthorized"
// @Failure 403 {object} model.APIError "Forbidden (e.g., if trying to update restricted fields)"
// @Failure 404 {object} model.APIError "Patient not found"
// @Failure 500 {object} model.APIError "Internal server error"
// @Router /patients/{id} [patch]
func (h *PatientHandler) UpdatePatient(c *gin.Context) {
	patientIDStr := c.Param("id")
	patientID := uuid.MustParse(patientIDStr)
	

	var req model.PatientUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{Message: "Invalid request body", Details: err.Error()})
		return
	}

	if err := util.ValidateStruct(req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{Message: "Validation failed", Details: util.FormatValidationErrors(err)})
		return
	}

	parsedReq, err := parsePatientRequest(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{Message: err.Error()})
		return
	}

	user, ok := c.Get("info")
	if !ok {
		log.Printf("CRITICAL: User context error")
		c.JSON(http.StatusInternalServerError, model.APIError{Message: "User context error"})
		return
	}
	_,ok = user.(authorization.Info)

	if !ok {
		log.Printf("CRITICAL: User context error")
		c.JSON(http.StatusInternalServerError, model.APIError{Message: "User context error"})
		return
	}

	role := user.(authorization.Info).Role
	userRole := model.UserRole(role)


	patient, err := h.patientService.UpdatePatientDetails(c.Request.Context(), patientID, *parsedReq, userRole, user.(authorization.Info).ID.Bytes)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			c.JSON(http.StatusNotFound, model.APIError{Message: "Patient not found"})
		} else if strings.Contains(strings.ToLower(err.Error()), "restricted") || strings.Contains(strings.ToLower(err.Error()), "authorized to update") {
			c.JSON(http.StatusForbidden, model.APIError{Message: err.Error()})
		} else if strings.Contains(strings.ToLower(err.Error()), "unique constraint") || strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
			c.JSON(http.StatusConflict, model.APIError{Message: "Update failed due to conflicting data (e.g., phone or email already exists for another patient)."})
		} else {
			log.Printf("Update patient error for ID %s: %v", patientIDStr, err)
			c.JSON(http.StatusInternalServerError, model.APIError{Message: "Failed to update patient details"})
		}
		return
	}

	c.JSON(http.StatusOK, patient)
}

// DeletePatient godoc
// @Summary Delete a patient record
// @Description Only doctors can delete patient records
// @Tags Patients
// @Security BearerAuth
// @Param id path string true "Patient ID (UUID)" Format(uuid)
// @Success 204 "Patient record deleted successfully"
// @Failure 400 {object} model.APIError "Validation error, invalid input, or invalid patient ID"
// @Failure 401 {object} model.APIError "Unauthorized"
// @Failure 403 {object} model.APIError "Forbidden (e.g., if trying to update restricted fields)"
// @Failure 404 {object} model.APIError "Patient not found"
// @Failure 500 {object} model.APIError "Internal server error"
// @Router /patients/{id} [delete]
func (h *PatientHandler) DeletePatient(c *gin.Context) {
	patientIDStr := c.Param("id")
	patientID := uuid.MustParse(patientIDStr)

	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		log.Printf("CRITICAL: UserID not found in context for an authenticated route in DeletePatient")
		c.JSON(http.StatusInternalServerError, model.APIError{Message: "User context error"})
		return
	}

	err := h.patientService.DeletePatientRecord(c.Request.Context(), patientID, userID) // userID for audit
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") || strings.Contains(strings.ToLower(err.Error()), "already deleted") {
			c.JSON(http.StatusNotFound, model.APIError{Message: "Patient not found or already deleted"})
		} else {
			log.Printf("Delete patient error for ID %s: %v", patientIDStr, err)
			c.JSON(http.StatusInternalServerError, model.APIError{Message: "Failed to delete patient"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
