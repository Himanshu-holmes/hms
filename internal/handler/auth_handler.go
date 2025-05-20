package handler

import (
	"log"
	"net/http"
	"strings" // Added for error checking

	"github.com/gin-gonic/gin"
	"github.com/himanshu-holmes/hms/internal/model"
	"github.com/himanshu-holmes/hms/internal/service"
	util "github.com/himanshu-holmes/hms/internal/utils"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login godoc
// @Summary User login
// @Description Authenticate a user and return a JWT token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginRequest body model.LoginRequest true "Login Credentials"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} model.APIError "Validation error or invalid input"
// @Failure 401 {object} model.APIError "Invalid credentials"
// @Failure 500 {object} model.APIError "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{Message: "Invalid request body", Details: err.Error()})
		return
	}

	if err := util.ValidateStruct(req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{Message: "Validation failed", Details: util.FormatValidationErrors(err)})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(),req)
	if err != nil {
		// Distinguish between bad credentials and server errors
		if err.Error() == "invalid username or password" { // Specific error check
			c.JSON(http.StatusUnauthorized, model.APIError{Message: err.Error()})
		} else {
			log.Printf("Login error: %v", err) // Log internal errors
			c.JSON(http.StatusInternalServerError, model.APIError{Message: "Login failed due to an internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreateUser godoc
// @Summary Create a new user (for seeding or admin purposes)
// @Description Allows creating a new user. In a real app, this might be restricted to admins.
// @Tags Users
// @Accept json
// @Produce json
// @Param userRequest body model.UserCreateRequest true "User Creation Data"
// @Success 201 {object} model.User
// @Failure 400 {object} model.APIError "Validation error or invalid input"
// @Failure 409 {object} model.APIError "User already exists"
// @Failure 500 {object} model.APIError "Internal server error"
// @Router /users [post]
func (h *AuthHandler) CreateUser(c *gin.Context) {
	var req model.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{Message: "Invalid request body", Details: err.Error()})
		return
	}

	if err := util.ValidateStruct(req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIError{Message: "Validation failed", Details: util.FormatValidationErrors(err)})
		return
	}

	user, err := h.authService.CreateUser(c.Request.Context(), req)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "already exists") ||
			strings.Contains(strings.ToLower(err.Error()), "unique constraint") {
			c.JSON(http.StatusConflict, model.APIError{Message: err.Error()})
		} else {
			log.Printf("Create user error: %v", err)
			c.JSON(http.StatusInternalServerError, model.APIError{Message: "Failed to create user"})
		}
		return
	}

	c.JSON(http.StatusCreated, user)
}