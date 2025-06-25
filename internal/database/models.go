package database

import (
	"log"

	contestDB "github.com/TRu-S3/backend/internal/database/contest"
	fileDB "github.com/TRu-S3/backend/internal/database/file"
	hackathonDB "github.com/TRu-S3/backend/internal/database/hackathon"
	userDB "github.com/TRu-S3/backend/internal/database/user"
	"gorm.io/gorm"
)

// Legacy type aliases for backward compatibility
type FileMetadata = fileDB.FileMetadata
type Contest = contestDB.Contest
type Hackathon = hackathonDB.Hackathon
type HackathonParticipant = hackathonDB.HackathonParticipant
type User = userDB.User
type Tag = userDB.Tag
type Profile = userDB.Profile
type Matching = userDB.Matching
type Bookmark = userDB.Bookmark

// Migrate runs database migrations for all domains
func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Migrate all domain models
	if err := fileDB.AutoMigrate(db); err != nil {
		return err
	}

	if err := contestDB.AutoMigrate(db); err != nil {
		return err
	}

	if err := hackathonDB.AutoMigrate(db); err != nil {
		return err
	}

	if err := userDB.AutoMigrate(db); err != nil {
		return err
	}

	// Add custom indexes
	if err := createIndexes(db); err != nil {
		log.Printf("Warning: Some indexes could not be created: %v", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// createIndexes creates additional database indexes for performance
func createIndexes(db *gorm.DB) error {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_file_metadata_name ON file_metadata(name)",
		"CREATE INDEX IF NOT EXISTS idx_file_metadata_path ON file_metadata(path)",
		"CREATE INDEX IF NOT EXISTS idx_file_metadata_created_at ON file_metadata(created_at)",
		"CREATE INDEX IF NOT EXISTS idx_contests_application_deadline ON contests(application_deadline)",
		"CREATE INDEX IF NOT EXISTS idx_contests_created_at ON contests(created_at)",
		"CREATE INDEX IF NOT EXISTS idx_hackathons_start_date ON hackathons(start_date)",
		"CREATE INDEX IF NOT EXISTS idx_hackathons_status ON hackathons(status)",
		"CREATE INDEX IF NOT EXISTS idx_hackathon_participants_hackathon_id ON hackathon_participants(hackathon_id)",
		"CREATE INDEX IF NOT EXISTS idx_users_gmail ON users(gmail)",
		"CREATE INDEX IF NOT EXISTS idx_profiles_user_id ON profiles(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_bookmarks_user_id ON bookmarks(user_id)",
	}

	for _, index := range indexes {
		if err := db.Exec(index).Error; err != nil {
			log.Printf("Warning: Could not create index: %v", err)
		}
	}

	return nil
}