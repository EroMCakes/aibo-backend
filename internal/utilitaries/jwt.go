package utilitaries

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWTClaim represents the claims in the JWT
type JWTClaim struct {
	AiboID string `json:"aibo_id"`
	jwt.StandardClaims
}

// GenerateJWT generates a new JWT token for a given user ID
func GenerateJWT(aiboID string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &JWTClaim{
		AiboID: aiboID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Get the JWT secret key from the environment variable
	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT validates the JWT token
func ValidateJWT(tokenString string) (*JWTClaim, error) {
	// Get the JWT secret key from the environment variable
	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return nil, fmt.Errorf("couldn't parse claims")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}
