package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"ecommitra-backend/config"
	"ecommitra-backend/models"
	"ecommitra-backend/utils"
	"time"
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

	// Generate Session Tokens
	access, refresh, sessionID, err := utils.GenerateTokens(fmt.Sprint(user.ID), user.Email)
	if err != nil {
		http.Error(w, "Error generating tokens", http.StatusInternalServerError)
		return
	}

	// Save Session to DB for revocation control
	session := models.Session{
		ID:        sessionID,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}
	config.DB.Create(&session)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token":         access, 
		"refresh_token": refresh,
		"email":         user.Email,
	})
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

	// Generate Session Tokens
	access, refresh, sessionID, err := utils.GenerateTokens(fmt.Sprint(user.ID), user.Email)
	if err != nil {
		http.Error(w, "Error generating tokens", http.StatusInternalServerError)
		return
	}

	// Save Session to DB
	session := models.Session{
		ID:        sessionID,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}
	config.DB.Create(&session)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token":         access, 
		"refresh_token": refresh,
		"email":         user.Email,
	})
}

// RefreshRequest structure
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// HandleRefresh processes a refresh token to issue a new access token
func HandleRefresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	claims, err := utils.ValidateJWT(req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	// Verify it's actually a refresh token
	if tokenType, ok := claims["type"].(string); !ok || tokenType != "refresh" {
		http.Error(w, "Invalid token type", http.StatusUnauthorized)
		return
	}

	sessionID := claims["session_id"].(string)

	// Verify session exists in DB
	var session models.Session
	if err := config.DB.Where("id = ?", sessionID).First(&session).Error; err != nil {
		http.Error(w, "Session revoked or not found", http.StatusUnauthorized)
		return
	}

	// Delete old session (rolling session)
	config.DB.Delete(&session)

	userID := claims["sub"].(string)
	email := claims["email"].(string)

	// Issue new tokens
	access, newRefresh, newSessionID, err := utils.GenerateTokens(userID, email)
	if err != nil {
		http.Error(w, "Error generating new tokens", http.StatusInternalServerError)
		return
	}

	// Save new Session to DB
	newSession := models.Session{
		ID:        newSessionID,
		UserID:    session.UserID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}
	config.DB.Create(&newSession)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token":         access,
		"refresh_token": newRefresh,
	})
}

// HandleLogout revokes a session
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest // We reuse RefreshRequest because it has the token
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// If no body, we might be able to get it from header but refresh token is usually in body
		http.Error(w, "Refresh token required to logout", http.StatusBadRequest)
		return
	}

	claims, err := utils.ValidateJWT(req.RefreshToken)
	if err == nil {
		if sessionID, ok := claims["session_id"].(string); ok {
			config.DB.Where("id = ?", sessionID).Delete(&models.Session{})
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}
