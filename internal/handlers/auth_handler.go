package handlers

import (
	"aibo/internal/database"
	"aibo/internal/types"
	"aibo/internal/utilitaries"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthService handles authentication-related requests.
type AuthService struct {
	DB             *gorm.DB
	AiboRepository *database.AiboRepository
}

// NewAuthService returns a new AuthService instance.
//
// The AuthService instance is configured with the provided db instance.
func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{DB: db, AiboRepository: database.NewAiboRepository(db)}
}

// Register creates a new aibo and returns a 201 status with a JSON response containing a message "aibo created successfully".
//
// The request body should contain an "email", a "password", and a "daily_budget" field.
//
// If the request body is invalid, it returns a 400 error with a JSON response containing the error message.
//
// If the aibo already exists, it returns a 409 error with a JSON response containing the error message.
//
// If there is an error creating the aibo, it returns a 500 error with a JSON response containing the error message.
// @Summary Register a new aibo
// @Description Create a new aibo account
// @Tags auth
// @Accept json
// @Produce json
// @Param aibo body types.RegisterRequest true "Aibo registration details"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
func (h *AuthService) Register(c *gin.Context) {
	var req types.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Failed to bind JSON", "error", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utilitaries.HashPassword(req.Password)
	if err != nil {
		slog.Error("Failed to hash password", "error", err)
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	req.Password = string(hashedPassword)

	var aibo types.Aibo = types.Aibo{
		Email:        req.Email,
		Password:     req.Password,
		CurrentDelta: 0,
	}

	if err := h.AiboRepository.CreateAibo(&aibo); err != nil {
		slog.Error("Failed to create aibo", "error", err)
		c.JSON(500, gin.H{"error": "Failed to create aibo"})
		return
	}

	c.JSON(201, gin.H{"message": "aibo created successfully"})
}

// Login authenticates an aibo and returns a JWT token if the credentials are valid.
//
// The request body should contain an "email" and a "password" field.
//
// If the credentials are invalid, it returns a 401 error with a message "Invalid credentials".
//
// If the credentials are valid, it returns a 200 status with a JSON response containing the JWT token.
// @Summary Login
// @Description Authenticate an aibo and receive a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body object{email=string,password=string} true "Login credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /login [post]
func (h *AuthService) Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		slog.Error("Failed to bind JSON", "error", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	aibo, err := h.AiboRepository.GetAiboByEmail(loginData.Email)
	if err != nil {
		slog.Error("Failed to get aibo", "error", err)
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	if !utilitaries.CheckPasswordHash(loginData.Password, aibo.Password) {
		slog.Error("Failed to check password hash", "error", err)
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utilitaries.GenerateJWT(fmt.Sprintf("%d", aibo.ID))
	if err != nil {
		slog.Error("Failed to generate token", "error", err)
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

// GetProfile returns the profile of the aibo that made the request.
//
// It reads the aibo ID from the JWT token and uses it to query the database.
// If the aibo is not found, it returns a 404 error.
// @Summary Get aibo profile
// @Description Get the profile of the authenticated aibo
// @Tags profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} types.Aibo
// @Failure 404 {object} map[string]string
// @Router /profile [get]
func (h *AuthService) GetProfile(c *gin.Context) {
	// Get aibo ID from JWT token
	aiboID := c.GetString("aibo_id")

	aibo, err := h.AiboRepository.GetAiboByID(aiboID)
	if err != nil {
		slog.Error("Failed to get aibo", "error", err)
		c.JSON(404, gin.H{"error": "aibo not found"})
		return
	}

	c.JSON(200, gin.H{"aibo": aibo})
}

// UpdateProfile updates the profile of the aibo that made the request.
//
// It reads the aibo ID from the JWT token and uses it to query the database.
// If the aibo is not found, it returns a 404 error.
// If the request body is invalid, it returns a 400 error.
// @Summary Update aibo profile
// @Description Update the profile of the authenticated aibo
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param aibo body types.UpdateProfileRequest true "Aibo update details"
// @Success 200 {object} types.Aibo
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /profile [put]
func (h *AuthService) UpdateProfile(c *gin.Context) {
	// Get aibo ID from JWT token
	aiboID := c.GetString("aibo_id")

	var req types.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Failed to bind JSON", "error", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	aibo, err := h.AiboRepository.GetAiboByID(aiboID)
	if err != nil {
		slog.Error("Failed to get aibo", "error", err)
		c.JSON(404, gin.H{"error": "aibo not found"})
		return
	}

	// Update fields only if they are provided in the request
	if req.FirstName != "" {
		aibo.FirstName = req.FirstName
	}
	if req.LastName != "" {
		aibo.LastName = req.LastName
	}
	if req.BirthDate != "" {
		// Parse the date string
		birthDate, err := time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			slog.Error("Failed to parse date string", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
		aibo.BirthDate = birthDate
	}

	err = h.AiboRepository.UpdateAibo(aibo)
	if err != nil {
		slog.Error("Failed to update aibo", "error", err)
		c.JSON(500, gin.H{"error": "Failed to update aibo"})
		return
	}

	c.JSON(200, gin.H{"aibo": aibo})
}

// UpdatePassword updates the password of the aibo that made the request.
//
// It reads the aibo ID from the JWT token and uses it to query the database.
// If the aibo is not found, it returns a 404 error.
// If the request body is invalid, it returns a 400 error.
// @Summary Update aibo password
// @Description Update the password of the authenticated aibo
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param aibo body types.UpdatePasswordRequest true "Aibo password update details"
// @Success 200 {object} types.Aibo
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /profile/password [put]
func (h *AuthService) UpdatePassword(c *gin.Context) {
	// Get aibo ID from JWT token
	aiboID := c.GetString("aibo_id")

	var req types.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Failed to bind JSON", "error", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	aibo, err := h.AiboRepository.GetAiboByID(aiboID)
	if err != nil {
		slog.Error("Failed to get aibo", "error", err)
		c.JSON(404, gin.H{"error": "aibo not found"})
		return
	}

	// Update password only if it is provided in the request
	if req.NewPassword != "" {
		hashedPassword, err := utilitaries.HashPassword(req.NewPassword)
		if err != nil {
			slog.Error("Failed to hash password", "error", err)
			c.JSON(500, gin.H{"error": "Failed to update password"})
			return
		}
		aibo.Password = hashedPassword
	}

	err = h.AiboRepository.UpdateAibo(aibo)
	if err != nil {
		slog.Error("Failed to update aibo", "error", err)
		c.JSON(500, gin.H{"error": "Failed to update aibo"})
		return
	}

	c.JSON(200, gin.H{"aibo": aibo})
}

// Logout logs out the aibo that made the request.
//
// It reads the aibo ID from the JWT token and uses it to query the database.
// If the aibo is not found, it returns a 404 error.
// @Summary Logout
// @Description Log out the authenticated aibo
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /logout [get]
func (h *AuthService) Logout(c *gin.Context) {
	c.Set("aibo_id", "")
	c.JSON(200, gin.H{"message": "aibo logged out successfully"})
}
