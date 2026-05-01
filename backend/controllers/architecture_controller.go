package controllers

import (
	"encoding/json"
	"net/http"

	"ecommitra-backend/config"
	"ecommitra-backend/models"
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
	var feature models.Feature
	if err := json.NewDecoder(r.Body).Decode(&feature); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if feature.ID == 0 {
		// New Feature
		if err := config.DB.Create(&feature).Error; err != nil {
			http.Error(w, "Failed to create feature", http.StatusInternalServerError)
			return
		}
	} else {
		// Existing Feature
		if err := config.DB.Save(&feature).Error; err != nil {
			http.Error(w, "Failed to update feature", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feature)
}

// HandleDeleteFeature removes a feature by ID
func HandleDeleteFeature(w http.ResponseWriter, r *http.Request) {
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
