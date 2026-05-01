package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"ecommitra-backend/config"
	"ecommitra-backend/models"
	"ecommitra-backend/utils"
)

// SignupRequest structure for incoming JSON
type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// HandleSignup processes new user registration
func HandleSignup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Hash password using military-grade Argon2id
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to secure password", http.StatusInternalServerError)
		return
	}

	// Save to DB
	user := models.User{Email: req.Email, PasswordHash: hash}
	if err := config.DB.Create(&user).Error; err != nil {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	// Generate Session Token
	token, _ := utils.GenerateJWT(fmt.Sprint(user.ID), user.Email)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token, "email": user.Email})
}

// HandleLogin processes user authentication
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Lookup user
	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Verify constant-time Argon2 hash
	valid, err := utils.VerifyPassword(req.Password, user.PasswordHash)
	if err != nil || !valid {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate Session Token
	token, _ := utils.GenerateJWT(fmt.Sprint(user.ID), user.Email)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token, "email": user.Email})
}
