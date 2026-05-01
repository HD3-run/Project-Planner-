package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"ecommitra-backend/config"
	"ecommitra-backend/models"
	"ecommitra-backend/utils"
)

// Define a custom type for context keys to avoid collisions
type ContextKey string
const UserContextKey ContextKey = "userEmail"
const RoleContextKey ContextKey = "userRole"

// RequireAuth is a middleware that intercepts incoming requests, verifies the JWT, and injects user context
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		
		// Ensure the header exists and starts with "Bearer "
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized: Missing or malformed token"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		
		// Validate using our secure JWT utility
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized: " + err.Error()})
			return
		}

		// Inject email and role into context
		ctx := r.Context()
		email, _ := claims["email"].(string)
		role, _ := claims["role"].(string)
		sessionID, _ := claims["session_id"].(string)

		// STATEFUL CHECK: Verify the session still exists in the DB
		var session models.Session
		if err := config.DB.Where("id = ? AND expires_at > ?", sessionID, time.Now()).First(&session).Error; err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized: Session has been revoked or is invalid. Please login again."})
			return
		}

		ctx = context.WithValue(ctx, UserContextKey, email)
		ctx = context.WithValue(ctx, RoleContextKey, role)
		r = r.WithContext(ctx)

		// Pass execution to the next handler
		next.ServeHTTP(w, r)
	}
}

// CORS Middleware to allow the frontend to communicate with the backend
func CORS(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	}
}
