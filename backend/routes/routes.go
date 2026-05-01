package routes

import (
	"encoding/json"
	"net/http"

	"ecommitra-backend/controllers"
	"ecommitra-backend/middleware"
)

// SetupRouter initializes all application routes and injects middleware
func SetupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Public Routes
	mux.HandleFunc("/api/signup", controllers.HandleSignup)
	mux.HandleFunc("/api/login", controllers.HandleLogin)
	mux.HandleFunc("/api/refresh", controllers.HandleRefresh)
	mux.HandleFunc("/api/logout", controllers.HandleLogout)
	mux.HandleFunc("/api/architecture", controllers.HandleGetArchitecture)

	// Protected Routes (Require JWT)
	// Consolidating all Feature operations to a single clean endpoint
	mux.HandleFunc("/api/architecture/feature", middleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controllers.HandleUpdateFeature(w, r)
		case http.MethodDelete:
			controllers.HandleDeleteFeature(w, r)
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		}
	}))

	return mux
}
