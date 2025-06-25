package main

import (
	"log"
	"os"

	"github.com/TRu-S3/backend/internal/config"
	"github.com/TRu-S3/backend/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Load configuration
	cfg := config.Load()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	log.Printf("Connecting to GCP database: %s:%s/%s", cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	log.Println("Database connection established successfully")

	// Run database migrations
	log.Println("Running GORM auto-migrations...")
	if err := database.Migrate(database.GetDB()); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Migration completed successfully!")

	// Test query to confirm tables exist
	db := database.GetDB()
	
	var tableCount int64
	result := db.Raw(`
		SELECT COUNT(*) 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_name IN ('users', 'tags', 'profiles', 'matchings', 'bookmarks', 'contests', 'file_metadata', 'hackathons', 'hackathon_participants')
	`).Scan(&tableCount)
	
	if result.Error != nil {
		log.Printf("Error checking tables: %v", result.Error)
	} else {
		log.Printf("Required tables found: %d/9", tableCount)
		if tableCount == 9 {
			log.Println("✅ All required tables are present!")
		} else {
			log.Printf("⚠️  Missing %d tables", 9-tableCount)
		}
	}

	// List all tables
	var tables []string
	result = db.Raw(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		ORDER BY table_name
	`).Scan(&tables)
	
	if result.Error != nil {
		log.Printf("Error listing tables: %v", result.Error)
	} else {
		log.Println("Current tables in database:")
		for _, table := range tables {
			log.Printf("  - %s", table)
		}
	}

	// Test hackathon data
	var hackathonCount int64
	result = db.Raw("SELECT COUNT(*) FROM hackathons").Scan(&hackathonCount)
	if result.Error != nil {
		log.Printf("Error checking hackathons: %v", result.Error)
	} else {
		log.Printf("Hackathons in database: %d", hackathonCount)
	}

	os.Exit(0)
}