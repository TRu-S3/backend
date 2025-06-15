package database

import (
	"log"

	"gorm.io/gorm"
)

// FileMetadata represents file metadata stored in database
type FileMetadata struct {
	ID          string `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Name        string `gorm:"not null;type:varchar(255)" json:"name"`
	Path        string `gorm:"not null;type:varchar(500)" json:"path"`
	Size        int64  `gorm:"not null" json:"size"`
	ContentType string `gorm:"type:varchar(100)" json:"content_type"`
	Checksum    string `gorm:"type:varchar(64)" json:"checksum,omitempty"`
	Tags        string `gorm:"type:text" json:"tags,omitempty"` // JSON string
	CreatedAt   int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   int64  `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName sets the table name for FileMetadata
func (FileMetadata) TableName() string {
	return "file_metadata"
}

// Migrate runs database migrations
func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	if err := db.AutoMigrate(&FileMetadata{}); err != nil {
		return err
	}

	// Add indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_file_metadata_name ON file_metadata(name)").Error; err != nil {
		log.Printf("Warning: Could not create index on name: %v", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_file_metadata_path ON file_metadata(path)").Error; err != nil {
		log.Printf("Warning: Could not create index on path: %v", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_file_metadata_created_at ON file_metadata(created_at)").Error; err != nil {
		log.Printf("Warning: Could not create index on created_at: %v", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}
