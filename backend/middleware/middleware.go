package middleware

import (
	"context"
	"net/http"
	"strings"
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
			http.Error(w, "Unauthorized: Missing or malformed token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		
		// Validate using our secure JWT utility
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Inject email and role into context
		ctx := r.Context()
		if email, ok := claims["email"].(string); ok {
			ctx = context.WithValue(ctx, UserContextKey, email)
		}
		if role, ok := claims["role"].(string); ok {
			ctx = context.WithValue(ctx, RoleContextKey, role)
		}
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
