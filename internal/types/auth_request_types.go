package types

import (
	"time"
)

// RegisterRequest represents the structure of the registration request
// @Description Registration request structure
type RegisterRequest struct {
	// User's email address
	// @example user@example.com
	Email string `json:"email" binding:"required,email"`
	// User's password (minimum 8 characters)
	// @example password123
	Password string `json:"password" binding:"required,min=8"`
	// User's first name
	// @example John
	FirstName string `json:"first_name" binding:"required"`
	// User's last name
	// @example Doe
	LastName string `json:"last_name" binding:"required"`
}

// LoginRequest represents the structure of the login request
// @Description Login request structure
type LoginRequest struct {
	// User's email address
	// @example user@example.com
	Email string `json:"email" binding:"required,email"`
	// User's password
	// @example password123
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest represents the structure of the update profile request
// @Description Update profile request structure
type UpdateProfileRequest struct {
	// User's new first name
	// @example John
	FirstName string `json:"first_name"`
	// User's new last name
	// @example Doe
	LastName string `json:"last_name"`
	// User's new birth date (format: YYYY-MM-DD)
	// @example 1990-01-01
	BirthDate string `json:"birth_date"`
}

// UpdatePasswordRequest represents the structure of the update password request
// @Description Update password request structure
type UpdatePasswordRequest struct {
	// User's current password
	// @example oldpassword123
	OldPassword string `json:"old_password" binding:"required"`
	// User's new password
	// @example newpassword123
	NewPassword string `json:"new_password" binding:"required"`
}

// UpdateProfileResponse represents the structure of the update profile response
// @Description Update profile response structure
type UpdateProfileResponse struct {
	// User's updated first name
	// @example John
	FirstName string `json:"first_name"`
	// User's updated last name
	// @example Doe
	LastName string `json:"last_name"`
	// User's updated birth date
	// @example 1990-01-01
	BirthDate string `json:"birth_date"`
}

// UpdatePasswordResponse represents the structure of the update password response
// @Description Update password response structure
type UpdatePasswordResponse struct {
	// Success message
	// @example Password updated successfully
	Message string `json:"message"`
}

// AuthResponse represents the structure of the authentication response
// @Description Authentication response structure
type AuthResponse struct {
	// JWT token for authentication
	// @example eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	Token string `json:"token"`
}

// UserResponse represents the structure of the user data in responses
// @Description User response structure
type UserResponse struct {
	// User's unique identifier
	// @example 1
	ID uint `json:"id"`
	// User's email address
	// @example user@example.com
	Email string `json:"email"`
	// User's first name
	// @example John
	FirstName string `json:"first_name"`
	// User's last name
	// @example Doe
	LastName string `json:"last_name"`
	// Whether the user has a premium account
	// @example false
	IsPremium bool `json:"is_premium"`
	// User's daily budget
	// @example 100.00
	DailyBudget float64 `json:"daily_budget"`
	// Current delta from the daily budget
	// @example -20.50
	CurrentDelta float64 `json:"current_delta"`
	// Timestamp of when the user was created
	// @example 2023-01-01T00:00:00Z
	CreatedAt time.Time `json:"created_at"`
	// Timestamp of when the user was last updated
	// @example 2023-01-02T00:00:00Z
	UpdatedAt time.Time `json:"updated_at"`
}

// ErrorResponse represents the structure of error responses
// @Description Error response structure
type ErrorResponse struct {
	// Error message
	// @example Invalid credentials
	Error string `json:"error"`
}
