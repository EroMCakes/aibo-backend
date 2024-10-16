package types

import (
	"time"

	"github.com/google/uuid"
)

// Aibo represents a user in the system
// @Description Aibo user model
type Aibo struct {
	// Unique identifier for the Aibo
	ID uuid.UUID `gorm:"type:char(36);primary_key;" json:"id" swaggertype:"string" format:"uuid"`
	// Timestamp of when the Aibo was created
	CreatedAt time.Time `json:"created_at"`
	// Timestamp of when the Aibo was last updated
	UpdatedAt time.Time `json:"updated_at"`
	// Email address of the Aibo
	Email string `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	// Hashed password of the Aibo
	Password string `gorm:"not null" json:"-"` // "-" means this field will be omitted in JSON responses
	// First name of the Aibo
	FirstName string `gorm:"type:varchar(255)" json:"first_name"`
	// Last name of the Aibo
	LastName string `gorm:"type:varchar(255)" json:"last_name"`
	// Birth date of the Aibo
	BirthDate time.Time `gorm:"type:date;" json:"birth_date"`
	// Whether the Aibo has a premium account
	IsPremium bool `gorm:"default:false" json:"is_premium"`
	// Daily budget set by the Aibo
	DailyBudget float64 `gorm:"default:0" json:"daily_budget"`
	// Current delta (difference) from the daily budget
	CurrentDelta float64 `gorm:"default:0" json:"current_delta"`
	// List of category-budget pairs associated with this Aibo
	CatBuds []CatBud `gorm:"foreignKey:AiboID" json:"cat_buds"`
}
