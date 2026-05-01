package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"ecommitra-backend/config"
	"ecommitra-backend/middleware"
	"ecommitra-backend/routes"
)

func main() {
	// 1. Initialize Database & GORM Migrations
	config.ConnectDatabase()

	// 2. Setup Routing
	router := routes.SetupRouter()

	// 3. Apply Global Middleware (CORS)
	handler := middleware.CORS(router)

	// 4. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Secure Custom Backend running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
