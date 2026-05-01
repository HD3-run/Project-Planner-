package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// getJWTSecret fetches the secret from environment, or uses a fallback for local dev
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return []byte("super_secret_fallback_key_for_dev_only")
	}
	return []byte(secret)
}

// GenerateJWT creates a signed JSON Web Token valid for 24 hours
func GenerateJWT(userID string, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), 
		"iat":   time.Now().Unix(),
	})
	return token.SignedString(getJWTSecret())
}

// ValidateJWT verifies the signature and expiration of a token
func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC to prevent algo downgrade attacks
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid or expired token")
}
