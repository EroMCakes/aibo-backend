package types

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

// CatBud represents a category-budget pair for an Aibo
// @Description Category and Budget pair model
type CatBud struct {
	// Unique identifier for the CatBud
	// @exemple 1234567890123456
	ID snowflake.ID `gorm:"primaryKey;type:bigint" json:"id"`
	// ID of the Aibo this CatBud belongs to
	AiboID uuid.UUID `gorm:"type:char(36);not null;" json:"aibo_id" swaggertype:"string" format:"uuid"`
	// Reference to the Aibo
	Aibo Aibo `gorm:"foreignKey:AiboID" json:"-"`
	// Name of the category
	Category string `gorm:"type:varchar(255);not null;" json:"category"`
	// Budget amount for the category (can be null)
	Budget *float64 `gorm:"type:decimal(10,2);default:null" json:"budget" swaggertype:"number"`
	// Timestamp of when the CatBud was created
	CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	// Timestamp of when the CatBud was last updated
	UpdatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}
