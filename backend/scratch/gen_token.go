package main

import (
	"fmt"
	"time"
	"github.com/golang-jwt/jwt/v4"
)

func main() {
	secret := []byte("ecommitra_ultra_secure_jwt_secret_2026_!!")
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":        "admin-user-id-123",
		"email":      "admin@ecommitra.com",
		"role":       "BABA",
		"session_id": "manual-gen-session",
		"type":       "access",
		"exp":        time.Now().Add(time.Hour * 24 * 365).Unix(), // 1 year expiry for testing
		"iat":        time.Now().Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println(tokenString)
}
