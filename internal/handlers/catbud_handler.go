package handlers

import (
	"aibo/internal/database"
	"aibo/internal/types"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CatBudService struct {
	DB               *gorm.DB
	CatBudRepository *database.CatBudRepository
}

// NewCatBudService creates a new CatBudService instance.
//
// The CatBudService instance is configured with the provided db instance.
func NewCatBudService(db *gorm.DB) *CatBudService {
	return &CatBudService{
		DB:               db,
		CatBudRepository: database.NewCatBudRepository(db),
	}
}

// GetCatBuds retrieves all CatBud entries from the database.
//
// The function queries the database for all CatBud entries and returns them as a slice.
//
// If no entries are found, a custom error is returned indicating that no CatBuds were found.
func (s *CatBudService) GetCatBuds(c *gin.Context) {
	aiboID := c.Param("aiboID")

	catBuds, err := s.CatBudRepository.GetAllCatBudsByAiboID(uuid.MustParse(aiboID))
	if err != nil {
		slog.Error("Failed to get cat buds", "error", err)
		c.JSON(404, gin.H{"error": err.Error()})
	}

	var resp types.GetCatBudsResponse
	resp.CatBuds = catBuds

	c.JSON(200, resp)
}

// CreateCatBuds creates a new CatBud entry in the database.
//
// The function reads the request body and creates a new CatBud entry using the provided data.
//
// If the request body is invalid, it returns a 400 error.
//
// If the CatBud is created successfully, it returns a 201 status with a JSON response containing the created CatBud.
func (s *CatBudService) CreateCatBuds(c *gin.Context) {
	var req types.CreateCatBudsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Failed to bind JSON", "error", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	authRepo := database.NewAiboRepository(s.DB)

	_, err := authRepo.GetAiboByID(req.AiboID.String())

	if err != nil {
		slog.Error("Failed to get aibo", "error", err)
		c.JSON(404, gin.H{"messge": "aibo not found", "error": err.Error()})
		return
	}

	for _, catBud := range req.CatBuds {
		if err := s.CatBudRepository.CreateCatBud(&catBud); err != nil {
			slog.Error("Failed to create cat bud", "error", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(201, gin.H{"message": "Cat bud created successfully"})
}

// UpdateCatBud updates an existing CatBud entry in the database.
//
// The function reads the request body and updates the CatBud entry using the provided data.
//
// If the request body is invalid, it returns a 400 error.
//
// If the CatBud is updated successfully, it returns a 200 status with a JSON response containing the updated CatBud.
func (s *CatBudService) UpdateCatBud(c *gin.Context) {
	var req types.UpdateCatBudRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Failed to bind JSON", "error", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	cb, err := s.CatBudRepository.GetCatBudByID(req.ID)

	if err != nil {
		slog.Error("Failed to get cat bud", "error", err)
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	if req.Category != "" {
		cb.Category = req.Category
	}

	if req.Budget != nil {
		cb.Budget = req.Budget
	}

	err = s.CatBudRepository.UpdateCatBud(cb)

	if err != nil {
		slog.Error("Failed to update cat bud", "error", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Cat bud updated successfully"})
}

// DeleteCatBud deletes an existing CatBud entry in the database.
//
// The function reads the ID of the CatBud to be deleted from the request body.
//
// If the CatBud is deleted successfully, it returns a 200 status with a JSON response indicating success.
func (s *CatBudService) DeleteCatBud(c *gin.Context) {
	var req types.DeleteCatBudRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Failed to bind JSON", "error", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := s.CatBudRepository.DeleteCatBudByID(req.ID); err != nil {
		slog.Error("Failed to delete cat bud", "error", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Cat bud deleted successfully"})
}
