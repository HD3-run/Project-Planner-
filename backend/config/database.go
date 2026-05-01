package config

import (
	"log"
	
	"ecommitra-backend/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// DB is the global database connection instance
var DB *gorm.DB

// ConnectDatabase initializes the pure-Go SQLite database and runs migrations
func ConnectDatabase() {
	var err error
	
	// Open connection to the SQLite file
	DB, err = gorm.Open(sqlite.Open("planner.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established successfully.")

	// AutoMigrate applies schema updates automatically based on our struct definitions
	log.Println("Running Auto Migrations...")
	err = DB.AutoMigrate(&models.User{}, &models.Section{}, &models.Feature{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	seedInitialData()
}

// seedInitialData ensures the database isn't completely empty on first launch
func seedInitialData() {
	var count int64
	DB.Model(&models.Section{}).Count(&count)
	
	if count == 0 {
		log.Println("Database is empty. Seeding initial Core Section data...")
		coreSection := models.Section{
			ID:        "core",
			Title:     "Core System",
			Color:     "#3b82f6",
			SortOrder: 1,
		}
		DB.Create(&coreSection)

		DB.Create(&models.Feature{
			SectionID: "core",
			Title:     "Initial Architecture",
			Status:    "live",
		})
	}
}
