package config

import (
	"log"
	"os"
	
	"ecommitra-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global database connection instance
var DB *gorm.DB

// ConnectDatabase initializes the PostgreSQL database connection and runs migrations
func ConnectDatabase() {
	var err error
	
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set. Please add your Supabase Postgres URI to the .env file.")
	}

	// Open connection to Supabase Postgres
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established successfully.")

	// Run migrations in background so server starts immediately
	go func() {
		log.Println("🚀 Running Auto Migrations in background...")
		err := DB.AutoMigrate(&models.User{}, &models.Section{}, &models.Feature{}, &models.Session{})
		if err != nil {
			log.Printf("⚠️ Background migration error: %v", err)
			return
		}
		seedInitialData()
		log.Println("✅ Database Migrations & Seeding Complete.")
	}()
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
