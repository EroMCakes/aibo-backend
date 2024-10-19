package types

import (
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

// CreateCatBudsRequest represents the request to create multiple CatBuds
// @Description Create multiple CatBuds request structure
type CreateCatBudsRequest struct {
	// ID of the Aibo associated with these CatBuds
	// @example 123e4567-e89b-12d3-a456-426614174000
	AiboID uuid.UUID `json:"aibo_id"`
	// List of CatBuds to create
	CatBuds []CatBud `json:"cat_buds"`
}

// UpdateCatBudRequest represents the request to update a CatBud
// @Description Update CatBud request structure
type UpdateCatBudRequest struct {
	// ID of the CatBud to update
	// @example 1234567890123456
	ID snowflake.ID `json:"id"`
	// Updated CatBud information
	Category string   `json:"category"`
	Budget   *float64 `json:"budget"`
}

// DeleteCatBudRequest represents the request to delete a CatBud
// @Description Delete CatBud request structure
type DeleteCatBudRequest struct {
	// ID of the CatBud to delete
	// @example 1234567890123456
	ID snowflake.ID `json:"id"`
}

// GetCatBudsRequest represents the request to get all CatBuds for an Aibo
// @Description Get CatBuds request structure
type GetCatBudsRequest struct {
	// ID of the Aibo to get CatBuds for
	// @example 123e4567-e89b-12d3-a456-426614174000
	AiboID uuid.UUID `json:"aibo_id"`
}

// GetCatBudsResponse represents the response containing multiple CatBuds
// @Description Get CatBuds response structure
type GetCatBudsResponse struct {
	// List of CatBuds
	CatBuds []CatBud `json:"cat_buds"`
}

// GetCatBudResponse represents the response containing a single CatBud
// @Description Get single CatBud response structure
type GetCatBudResponse struct {
	// The requested CatBud
	CatBud CatBud `json:"cat_bud"`
}

// UpdateCatBudResponse represents the response after updating a CatBud
// @Description Update CatBud response structure
type UpdateCatBudResponse struct {
	// The updated CatBud
	CatBud CatBud `json:"cat_bud"`
}

// DeleteCatBudResponse represents the response after deleting a CatBud
// @Description Delete CatBud response structure
type DeleteCatBudResponse struct {
	// Success message
	// @example CatBud successfully deleted
	Message string `json:"message"`
}

// CreateCatBudsResponse represents the response after creating CatBuds
// @Description Create CatBuds response structure
type CreateCatBudsResponse struct {
	// Success message
	// @example CatBuds successfully created
	Message string `json:"message"`
}
