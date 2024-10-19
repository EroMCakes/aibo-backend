package database

import (
	"aibo/internal/types"
	"fmt"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CatBudRepository struct {
	db *gorm.DB
}

// NewCatBudRepository creates a new CatBudRepository instance.
//
// The CatBudRepository instance is configured with the provided db instance.
func NewCatBudRepository(db *gorm.DB) *CatBudRepository {
	return &CatBudRepository{db: db}
}

// CreateCatBud creates a new CatBud entry in the database.
//
// The CatBud is created using the provided CatBud instance. If the CatBud is created successfully,
// a nil error is returned. If there is an error during creation, a gorm error is returned.
func (r *CatBudRepository) CreateCatBud(catBud *types.CatBud) error {
	return r.db.Create(catBud).Error
}

// UpdateCatBud updates an existing CatBud entry in the database.
//
// The CatBud is updated using the provided CatBud instance. If the CatBud is updated successfully,
// a nil error is returned. If there is an error during updating, a gorm error is returned.
func (r *CatBudRepository) UpdateCatBud(catBud *types.CatBud) error {
	return r.db.Save(catBud).Error
}

// GetCatBudByID retrieves a CatBud entry by its ID from the database.
//
// The CatBud is queried by its ID and returned if found. If the CatBud is not found,
// a gorm.NotFound error is returned.
func (r *CatBudRepository) GetCatBudByID(id snowflake.ID) (*types.CatBud, error) {
	var catBud types.CatBud
	err := r.db.First(&catBud, "id = ?", id).Error
	return &catBud, err
}

// GetAllCatBudsByAiboID retrieves all CatBud entries associated with a specific AiboID from the database.
//
// The function queries the database for CatBud entries that match the provided AiboID. If entries are found,
// they are returned as a slice. If no entries are found, a custom error is returned indicating that no CatBuds
// were found for the given AiboID.
func (r *CatBudRepository) GetAllCatBudsByAiboID(aiboID uuid.UUID) ([]types.CatBud, error) {
	var catBuds []types.CatBud
	result := r.db.Where("aibo_id = ?", aiboID).Find(&catBuds)
	if result.Error != nil {
		return nil, result.Error
	}

	// If no records are found, Find() will return an empty slice without an error
	// We can check if the slice is empty and return a custom error if needed
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("no CatBuds found for AiboID: %v", aiboID)
	}

	return catBuds, nil
}

// DeleteCatBudByID deletes a CatBud entry by its ID from the database.
//
// The CatBud is deleted using the provided ID. If the CatBud is deleted successfully,
// a nil error is returned. If there is an error during deletion, a gorm error is returned.
func (r *CatBudRepository) DeleteCatBudByID(id snowflake.ID) error {
	return r.db.Delete(&types.CatBud{}, "id = ?", id).Error
}
