package handlers

import (
	"aibo/internal/database"
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
// func (s *CatBudService) GetCatBuds(c *gin.Context) {
// 	aiboID := c.Param("aiboID")

// 	catBuds, err := s.CatBudRepository.GetAllCatBudsByAiboID(uuid.MustParse(aiboID))
// 	if err != nil {
// 		slog.Error("Failed to get cat buds", "error", err)
// 		c.JSON(404, gin.H{"error": err.Error()})
// 	}
// 	return s.CatBudRepository.GetAllCatBudsByAiboID()
// }
