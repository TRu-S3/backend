package file

import (
	"time"

	"gorm.io/gorm"
)

// FileMetadata represents file metadata stored in database
type FileMetadata struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"unique;not null" json:"name"`
	Path        string    `gorm:"not null" json:"path"`
	Size        int64     `gorm:"not null" json:"size"`
	ContentType string    `gorm:"not null" json:"content_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Models returns all file-related models for migration
var Models = []interface{}{
	&FileMetadata{},
}

// AutoMigrate performs auto-migration for file models
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(Models...)
}