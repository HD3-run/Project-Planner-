package controllers

import (
	"encoding/json"
	"net/http"

	"ecommitra-backend/config"
	"ecommitra-backend/models"
	"ecommitra-backend/middleware"
	"os"
)

// HandleGetArchitecture fetches all sections and their nested features
func HandleGetArchitecture(w http.ResponseWriter, r *http.Request) {
	var sections []models.Section
	
	// Preload automatically runs the necessary JOINs to get all features inside sections
	if err := config.DB.Preload("Features").Order("sort_order asc").Find(&sections).Error; err != nil {
		http.Error(w, "Failed to fetch architecture: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sections)
}

// HandleUpdateFeature creates or updates a feature in the roadmap
func HandleUpdateFeature(w http.ResponseWriter, r *http.Request) {
	// Security Check: Only Admin can edit
	role, _ := r.Context().Value(middleware.RoleContextKey).(string)
	adminRole := os.Getenv("ADMIN_ROLE")
	if adminRole == "" {
		adminRole = "admin"
	}
	
	if role != adminRole {
		http.Error(w, "Access Denied: Only the "+adminRole+" can modify the architecture.", http.StatusForbidden)
		return
	}

	var feature models.Feature
	if err := json.NewDecoder(r.Body).Decode(&feature); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON: " + err.Error()})
		return
	}

	if feature.ID == 0 {
		// New Feature
		if err := config.DB.Create(&feature).Error; err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "DB Create Error: " + err.Error()})
			return
		}
	} else {
		// Existing Feature - Use Save which updates all fields including those missing in JSON
		if err := config.DB.Save(&feature).Error; err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "DB Update Error: " + err.Error()})
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feature)
}

// HandleDeleteFeature removes a feature by ID
func HandleDeleteFeature(w http.ResponseWriter, r *http.Request) {
	// Security Check: Only Admin can delete
	role, _ := r.Context().Value(middleware.RoleContextKey).(string)
	adminRole := os.Getenv("ADMIN_ROLE")
	if adminRole == "" {
		adminRole = "admin"
	}

	if role != adminRole {
		http.Error(w, "Access Denied: Only the "+adminRole+" can modify the architecture.", http.StatusForbidden)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing Feature ID", http.StatusBadRequest)
		return
	}

	if err := config.DB.Delete(&models.Feature{}, id).Error; err != nil {
		http.Error(w, "Failed to delete feature", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
