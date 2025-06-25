package contest

import (
	"time"

	"gorm.io/gorm"
)

// Contest represents a programming contest
type Contest struct {
	ID                  uint      `gorm:"primarykey" json:"id"`
	BackendQuota        int       `gorm:"not null;default:0" json:"backend_quota"`
	FrontendQuota       int       `gorm:"not null;default:0" json:"frontend_quota"`
	AIQuota             int       `gorm:"not null;default:0" json:"ai_quota"`
	ApplicationDeadline time.Time `gorm:"not null" json:"application_deadline"`
	Purpose             string    `gorm:"not null;type:text" json:"purpose"`
	Message             string    `gorm:"not null;type:text" json:"message"`
	AuthorID            uint      `gorm:"not null" json:"author_id"`
	StartTime           time.Time `json:"start_time"`
	EndTime             time.Time `json:"end_time"`
	Title               string    `json:"title"`
	Description         string    `json:"description"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// Models returns all contest-related models for migration
var Models = []interface{}{
	&Contest{},
}

// AutoMigrate performs auto-migration for contest models
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(Models...)
}