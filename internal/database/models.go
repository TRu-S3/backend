package database

import (
	"log"
	"time"

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

// Contest represents contest information stored in database
type Contest struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	BackendQuota         int       `gorm:"not null;default:0" json:"backend_quota"`
	FrontendQuota        int       `gorm:"not null;default:0" json:"frontend_quota"`
	AIQuota              int       `gorm:"not null;default:0" json:"ai_quota"`
	ApplicationDeadline  time.Time `gorm:"not null" json:"application_deadline"`
	Purpose              string    `gorm:"not null;type:text" json:"purpose"`
	Message              string    `gorm:"not null;type:text" json:"message"`
	AuthorID             uint      `gorm:"not null" json:"author_id"`
	CreatedAt            time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName sets the table name for Contest
func (Contest) TableName() string {
	return "contests"
}

// Migrate runs database migrations
func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	if err := db.AutoMigrate(&FileMetadata{}, &Contest{}); err != nil {
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

	// Add contest indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_contests_author_id ON contests(author_id)").Error; err != nil {
		log.Printf("Warning: Could not create index on contests author_id: %v", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_contests_application_deadline ON contests(application_deadline)").Error; err != nil {
		log.Printf("Warning: Could not create index on contests application_deadline: %v", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_contests_created_at ON contests(created_at)").Error; err != nil {
		log.Printf("Warning: Could not create index on contests created_at: %v", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}
