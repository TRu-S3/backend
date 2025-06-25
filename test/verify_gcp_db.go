package main

import (
	"fmt"
	"log"
	"os"

	"github.com/TRu-S3/backend/internal/config"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables from .env.gcp.backup
	if err := godotenv.Load(".env.gcp.backup"); err != nil {
		log.Printf("Warning: Could not load .env.gcp.backup file: %v", err)
	}

	// Load configuration
	cfg := config.Load()

	// Create direct GORM connection to GCP database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	log.Printf("Attempting to connect to GCP database: %s:%s/%s", cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("❌ Failed to connect to GCP database: %v", err)
		os.Exit(1)
	}

	log.Println("✅ Successfully connected to GCP database")

	// Test 1: Check all tables
	var tables []string
	result := db.Raw(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		ORDER BY table_name
	`).Scan(&tables)

	if result.Error != nil {
		log.Printf("❌ Error listing tables: %v", result.Error)
		os.Exit(1)
	}

	log.Printf("📋 Found %d tables in database:", len(tables))
	for i, table := range tables {
		log.Printf("  %d. %s", i+1, table)
	}

	// Test 2: Check required tables
	requiredTables := []string{"users", "tags", "profiles", "matchings", "bookmarks", "contests", "file_metadata", "hackathons", "hackathon_participants"}
	var foundCount int64

	for _, table := range requiredTables {
		var exists bool
		result := db.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = ? AND table_schema = 'public')", table).Scan(&exists)
		if result.Error != nil {
			log.Printf("❌ Error checking table %s: %v", table, result.Error)
			continue
		}
		if exists {
			foundCount++
			log.Printf("✅ %s - EXISTS", table)
		} else {
			log.Printf("❌ %s - MISSING", table)
		}
	}

	log.Printf("📊 Required tables found: %d/9", foundCount)

	// Test 3: Check hackathon data
	var hackathonCount int64
	result = db.Raw("SELECT COUNT(*) FROM hackathons").Scan(&hackathonCount)
	if result.Error != nil {
		log.Printf("❌ Error counting hackathons: %v", result.Error)
	} else {
		log.Printf("🏆 Hackathons in database: %d", hackathonCount)
	}

	// Test 4: Sample hackathon data
	type HackathonInfo struct {
		Name      string `json:"name"`
		Organizer string `json:"organizer"`
		Status    string `json:"status"`
	}
	
	var hackathons []HackathonInfo
	result = db.Raw("SELECT name, organizer, status FROM hackathons ORDER BY created_at").Scan(&hackathons)
	if result.Error != nil {
		log.Printf("❌ Error fetching hackathon details: %v", result.Error)
	} else {
		log.Println("🎯 Hackathon details:")
		for i, h := range hackathons {
			log.Printf("  %d. %s (by %s) - %s", i+1, h.Name, h.Organizer, h.Status)
		}
	}

	// Final summary
	if foundCount == 9 {
		log.Println("🎉 SUCCESS: All required tables are present in GCP database!")
	} else {
		log.Printf("⚠️  WARNING: Only %d out of 9 required tables found", foundCount)
	}

	// Test row counts for each existing table
	log.Println("📈 Row counts per table:")
	for _, table := range requiredTables {
		var count int64
		result := db.Raw(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
		if result.Error != nil {
			log.Printf("  %s: ERROR - %v", table, result.Error)
		} else {
			log.Printf("  %s: %d rows", table, count)
		}
	}

	os.Exit(0)
}