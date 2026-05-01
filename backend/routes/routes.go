package routes

import (
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
	// We wrap the controller functions in the RequireAuth middleware
	mux.HandleFunc("/api/features/update", middleware.RequireAuth(controllers.HandleUpdateFeature))
	mux.HandleFunc("/api/features/delete", middleware.RequireAuth(controllers.HandleDeleteFeature))

	return mux
}
