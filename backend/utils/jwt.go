package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// getJWTSecret fetches the secret from environment, or uses a fallback for local dev
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return []byte("super_secret_fallback_key_for_dev_only")
	}
	return []byte(secret)
}

// GenerateSessionID creates a fast, secure random hex string for the session ID
func GenerateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// GenerateTokens creates both an Access Token (15m) and Refresh Token (7d)
func GenerateTokens(userID string, email string, role string) (accessTokenString string, refreshTokenString string, sessionID string, err error) {
	sessionID = GenerateSessionID()

	// Access Token: 15 minutes
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":        userID,
		"email":      email,
		"role":       role,
		"session_id": sessionID, // Added for stateful revocation!
		"type":       "access",
		"exp":        time.Now().Add(time.Minute * 15).Unix(),
		"iat":        time.Now().Unix(),
	})
	accessString, err := accessToken.SignedString(getJWTSecret())
	if err != nil {
		return "", "", "", err
	}

	// Refresh Token: 7 days
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":        userID,
		"email":      email,
		"role":       role,
		"type":       "refresh",
		"session_id": sessionID, // Embed the DB Session ID so we can revoke it!
		"exp":        time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat":        time.Now().Unix(),
	})
	refreshString, err := refreshToken.SignedString(getJWTSecret())
	if err != nil {
		return "", "", "", err
	}

	return accessString, refreshString, sessionID, nil
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
