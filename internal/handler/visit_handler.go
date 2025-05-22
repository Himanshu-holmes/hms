package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/himanshu-holmes/hms/internal/authorization"
	"github.com/himanshu-holmes/hms/internal/model"
	"github.com/himanshu-holmes/hms/internal/service"
	util "github.com/himanshu-holmes/hms/internal/utils"
)

type PatientVisitHandler struct {
	visitService service.PatientVisitService
}

func NewPatientVisitHandler(visitService service.PatientVisitService) *PatientVisitHandler {
	return &PatientVisitHandler{visitService: visitService}
}

// parseVisitRequest handles parsing for both create and update visit requests.
// It converts string dates to time.Time and sets the DoctorID from context.
func parseVisitRequest(c *gin.Context, req interface{}) (*model.ParsedPatientVisitRequest, error) {
	var visitDate time.Time
	var err error

	parsedReq := &model.ParsedPatientVisitRequest{}

	user,ok := c.Get("info")
	if !ok {
		// This should not happen if AuthMiddleware is working correctly
		return nil, fmt.Errorf("user not found in context (critical error)")
	}
	 info,ok := user.(authorization.Info)
	if !ok {
		// This should not happen if AuthMiddleware is working correctly
		return nil, fmt.Errorf("user not found in context (critical error)")
	}


	// DoctorID for the visit is always the authenticated user performing the action
	doctorID, err := uuid.Parse(info.ID.String())
	if err != nil {
		return nil, fmt.Errorf("invalid doctor_id: %w", err)
	}
	parsedReq.DoctorID = doctorID

	switch r := req.(type) {
	case model.PatientVisitCreateRequest:
		if r.VisitDateStr != "" {
			visitDate, err = time.Parse(time.RFC3339, r.VisitDateStr)
			if err != nil {
				return nil, fmt.Errorf("invalid visit_date format: %w. Expected RFC3339 (e.g., 2023-10-27T10:00:00Z)", err)
			}
		} else {
			// Default to now if not provided, or handle as error based on requirements
			visitDate = time.Now()
		}
		parsedReq.PatientID = r.PatientID // PatientID must be provided for creation
		parsedReq.VisitDate = visitDate
		parsedReq.Symptoms = r.Symptoms
		parsedReq.Diagnosis = r.Diagnosis
		parsedReq.Prescription = r.Prescription
		parsedReq.Notes = r.Notes
	case model.PatientVisitUpdateRequest:
		// PatientID is not updatable for an existing visit through this request.
		// It's tied to the visit's identity.
		if r.VisitDateStr != nil && *r.VisitDateStr != "" {
			visitDate, err = time.Parse(time.RFC3339, *r.VisitDateStr)
			if err != nil {
				return nil, fmt.Errorf("invalid visit_date format: %w. Expected RFC3339", err)
			}
			parsedReq.VisitDate = visitDate // Only update if provided and valid
		}
		// Other fields are optional for update
		parsedReq.Symptoms = r.Symptoms
		parsedReq.Diagnosis = r.Diagnosis
		parsedReq.Prescription = r.Prescription
		parsedReq.Notes = r.Notes
	default:
		return nil, fmt.Errorf("unsupported request type for patient visit parsing")
	}
	return parsedReq, nil
}

// RecordPatientVisit godoc
// @Summary Record a visit for a patient
// @Description Doctors can record a visit for a patient.
// @Tags Visits
// @Security BearerAuth
// @Produce json
// @Param request body model.PatientVisitCreateRequest true "Visit details"
// @Success 201 {object} model.PatientVisit
// @Failure 400 {object} model.APIError "Invalid request body"
// @Failure 401 {object} model.APIError "Unauthorized"
// @Failure 404 {object} model.APIError "Patient not found"
// @Failure 500 {object} model.APIError "Failed to record patient visit"
// @Router /visits/create [post]

func (h *PatientVisitHandler) RecordPatientVisit(c *gin.Context) {
	var req model.PatientVisitCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{Message: "Invalid request body", Details: err.Error()})
		return
	}

	if err := util.ValidateStruct(req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{Message: "Validation failed", Details: util.FormatValidationErrors(err)})
		return
	}

	parsedReq, err := parseVisitRequest(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{Message: err.Error()})
		return
	}
	// Ensure PatientID was part of the create request
	if parsedReq.PatientID == uuid.Nil {
		c.JSON(http.StatusBadRequest, model.APIError{Message: "patient_id is required for creating a visit"})
		return
	}


	visit, err := h.visitService.RecordPatientVisit(c.Request.Context(), *parsedReq)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "patient with id") && strings.Contains(strings.ToLower(err.Error()), "not found") {
			c.JSON(http.StatusNotFound, model.APIError{Message: fmt.Sprintf("Patient with ID %s not found", parsedReq.PatientID)})
		} else if strings.Contains(strings.ToLower(err.Error()), "doctor id not found") { // from parseVisitRequest
		    log.Printf("CRITICAL: DoctorID not found in context for RecordPatientVisit")
			c.JSON(http.StatusInternalServerError, model.APIError{Message: "User context error"})
		} else {
			log.Printf("Record patient visit error for patient %s by doctor %s: %v", parsedReq.PatientID, parsedReq.DoctorID, err)
			c.JSON(http.StatusInternalServerError, model.APIError{Message: "Failed to record patient visit"})
		}
		return
	}

	c.JSON(http.StatusCreated, visit)
}

// GetPatientVisitDetails godoc
// @Summary Get details of a specific patient visit
// @Description Doctors and Receptionists can get details of a specific patient visit.
// @Tags Visits
// @Security BearerAuth
// @Produce json
// @Param id path string true "Visit ID (UUID) for which to get details" Format(uuid)
// @Success 200 {object} model.PatientVisit
// @Failure 400 {object} model.APIError "Invalid visit ID format"
// @Failure 404 {object} model.APIError "Visit not found"
// @Failure 401 {object} model.APIError "Unauthorized"
// @Router /visits/{id} [get]

func (h *PatientVisitHandler) GetPatientVisitDetails(c *gin.Context) {
	visitIDStr := c.Param("id")
	visitID, err := util.GetUserIDFromString(visitIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{Message: "Invalid visit ID format"})
		return
	}

	visit, err := h.visitService.GetPatientVisitDetails(c.Request.Context(), visitID.Bytes)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			c.JSON(http.StatusNotFound, model.APIError{Message: "Visit not found"})
		} else {
			log.Printf("Get visit details error for ID %s: %v", visitIDStr, err)
			c.JSON(http.StatusInternalServerError, model.APIError{Message: "Failed to get visit details"})
		}
		return
	}

	c.JSON(http.StatusOK, visit)
}

// ListPatientVisits godoc
// @Summary List patient visits for a specific patient
// @Description Doctors and Receptionists can list patient visits for a specific patient.
// @Tags Visits
// @Security BearerAuth
// @Produce json
// @Param id path string true "Patient ID (UUID) for which to list visits" Format(uuid)
// @Param pagination query model.PaginationParams true "Pagination parameters"
// @Success 200 {object} model.PaginatedResponse
// @Failure 400 {object} model.APIError "Invalid patient ID format"
// @Failure 401 {object} model.APIError "Unauthorized"
// @Router /visits/{id}/list [get]
func (h *PatientVisitHandler) ListPatientVisits(c *gin.Context) {
	patientIDStr := c.Param("id")
	// patientID, err := util.GetUserIDFromString(patientIDStr)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, model.APIError{Message: "Invalid patient ID format"})
	// 	return
	// }

	var params model.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{Message: "Invalid pagination parameters", Details: err.Error()})
		return
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}
	if params.Offset < 0 {
		params.Offset = 0
	}

	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{Message: "Invalid patient ID format"})
		return
	}
	visits, total, err := h.visitService.ListPatientVisits(c.Request.Context(), patientID, params)
	if err != nil {
		log.Printf("List patient visits error for patient %s: %v", patientIDStr, err)
		c.JSON(http.StatusInternalServerError, model.APIError{Message: "Failed to list patient visits"})
		return
	}

	response := model.PaginatedResponse{
		Data:   visits,
		Total:  total,
		Limit:  params.Limit,
		Offset: params.Offset,
	}
	c.JSON(http.StatusOK, response)
}

// UpdatePatientVisit godoc
// @Summary Update a specific patient visit
// @Description Doctors can update patient visit details they recorded. Doctor ID is taken from authenticated user.
// @Tags Visits
// @Security BearerAuth
// @Produce json
// @Param id path string true "Visit ID (UUID) for which to update details" Format(uuid)
// @Param visit body model.PatientVisitUpdateRequest true "Patient visit details to update"
// @Success 200 {object} model.PatientVisit
// @Failure 400 {object} model.APIError "Invalid visit ID format"
// @Failure 400 {object} model.APIError "Invalid request body"
// @Failure 400 {object} model.APIError "Validation failed"
// @Failure 404 {object} model.APIError "Visit not found"
// @Failure 401 {object} model.APIError "Unauthorized"
// @Router /visits/{id} [patch]
func (h *PatientVisitHandler) UpdatePatientVisit(c *gin.Context) {
	visitIDStr := c.Param("id")
	visitID, err := uuid.Parse(visitIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{Message: "Invalid visit ID format"})
		return
	}

	var req model.PatientVisitUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{Message: "Invalid request body", Details: err.Error()})
		return
	}

	if err := util.ValidateStruct(req); err != nil { // Validate the update request struct
		c.JSON(http.StatusBadRequest, model.APIError{Message: "Validation failed", Details: util.FormatValidationErrors(err)})
		return
	}

	parsedReq, err := parseVisitRequest(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{Message: err.Error()})
		return
	}
	// PatientID is not set from request for update, it's tied to the visitID
	// DoctorID is set from context in parseVisitRequest

	updatedVisit, err := h.visitService.UpdatePatientVisit(c.Request.Context(), visitID, *parsedReq)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			c.JSON(http.StatusNotFound, model.APIError{Message: "Visit not found"})
		} else if strings.Contains(strings.ToLower(err.Error()), "not authorized to update") {
			c.JSON(http.StatusForbidden, model.APIError{Message: "Not authorized to update this visit record"})
		} else if strings.Contains(strings.ToLower(err.Error()), "doctor id not found") { // from parseVisitRequest
		    log.Printf("CRITICAL: DoctorID not found in context for UpdatePatientVisit")
			c.JSON(http.StatusInternalServerError, model.APIError{Message: "User context error"})
		} else {
			log.Printf("Update patient visit error for visit %s by doctor %s: %v", visitIDStr, parsedReq.DoctorID, err)
			c.JSON(http.StatusInternalServerError, model.APIError{Message: "Failed to update patient visit"})
		}
		return
	}

	c.JSON(http.StatusOK, updatedVisit)
}