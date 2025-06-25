package hackathon

import (
	"time"

	"gorm.io/gorm"
)

// Hackathon represents a hackathon event
type Hackathon struct {
	ID                   uint                   `gorm:"primarykey" json:"id"`
	Name                 string                 `gorm:"not null;type:varchar(255)" json:"name"`
	Description          string                 `gorm:"type:text" json:"description"`
	StartDate            time.Time              `gorm:"not null" json:"start_date"`
	EndDate              time.Time              `gorm:"not null" json:"end_date"`
	RegistrationStart    time.Time              `gorm:"not null" json:"registration_start"`
	RegistrationDeadline time.Time              `gorm:"not null" json:"registration_deadline"`
	MaxParticipants      int                    `gorm:"default:0" json:"max_participants"`
	Location             string                 `gorm:"type:varchar(255)" json:"location"`
	Organizer            string                 `gorm:"not null;type:varchar(255)" json:"organizer"`
	ContactEmail         string                 `gorm:"type:varchar(255)" json:"contact_email"`
	PrizeInfo            string                 `gorm:"type:text" json:"prize_info"`
	Rules                string                 `gorm:"type:text" json:"rules"`
	TechStack            string                 `gorm:"type:text" json:"tech_stack"` // JSON string
	Status               string                 `gorm:"default:upcoming;type:varchar(50)" json:"status"`
	IsPublic             bool                   `gorm:"default:true" json:"is_public"`
	BannerURL            string                 `gorm:"type:text" json:"banner_url"`
	WebsiteURL           string                 `gorm:"type:text" json:"website_url"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
	Participants         []HackathonParticipant `gorm:"foreignKey:HackathonID" json:"participants,omitempty"`
}

// HackathonParticipant represents a participant in a hackathon
type HackathonParticipant struct {
	ID               uint      `gorm:"primarykey" json:"id"`
	HackathonID      uint      `gorm:"not null" json:"hackathon_id"`
	UserID           uint      `gorm:"not null" json:"user_id"`
	TeamName         string    `gorm:"type:varchar(255)" json:"team_name"`
	Role             string    `gorm:"type:varchar(100)" json:"role"`
	RegistrationDate time.Time `gorm:"autoCreateTime" json:"registration_date"`
	Status           string    `gorm:"default:registered;type:varchar(50)" json:"status"`
	Notes            string    `gorm:"type:text" json:"notes"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// Foreign key constraints
	Hackathon Hackathon `gorm:"foreignKey:HackathonID;constraint:OnDelete:CASCADE" json:"hackathon,omitempty"`
}

// Models returns all hackathon-related models for migration
var Models = []interface{}{
	&Hackathon{},
	&HackathonParticipant{},
}

// AutoMigrate performs auto-migration for hackathon models
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(Models...)
}