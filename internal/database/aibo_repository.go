package database

import (
	"aibo/internal/types"

	"gorm.io/gorm"
)

type AiboRepository struct {
	db *gorm.DB
}

// NewAiboRepository creates a new AiboRepository instance.
//
// The AiboRepository instance is configured with the provided db instance.
func NewAiboRepository(db *gorm.DB) *AiboRepository {
	return &AiboRepository{db: db}
}

// CreateAibo creates a new Aibo in the database.
//
// The Aibo is created using the provided Aibo instance. If the Aibo is created successfully, a
// nil error is returned. If there is an error creating the Aibo, a gorm.error is returned.
func (r *AiboRepository) CreateAibo(Aibo *types.Aibo) error {
	return r.db.Create(Aibo).Error
}

// GetAiboByEmail returns an Aibo by its email.
//
// The Aibo is queried by its email and returned if found. If the Aibo is not found, a
// gorm.NotFound error is returned.
func (r *AiboRepository) GetAiboByEmail(email string) (*types.Aibo, error) {
	var Aibo types.Aibo
	err := r.db.Where("email = ?", email).First(&Aibo).Error
	return &Aibo, err
}

// GetAiboByID returns an Aibo by its ID.
//
// The Aibo is queried by its ID and returned if found. If the Aibo is not found, a
// gorm.NotFound error is returned.
func (r *AiboRepository) GetAiboByID(id string) (*types.Aibo, error) {
	var Aibo types.Aibo
	err := r.db.First(&Aibo, "id = ?", id).Error
	return &Aibo, err
}

// UpdateAibo updates an existing Aibo in the database.
//
// The Aibo is updated using the provided Aibo instance. If the Aibo is updated successfully, a
// nil error is returned. If there is an error updating the Aibo, a gorm.error is returned.
func (r *AiboRepository) UpdateAibo(Aibo *types.Aibo) error {
	return r.db.Save(Aibo).Error
}
