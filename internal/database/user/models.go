package user

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Gmail     string    `gorm:"size:100;uniqueIndex;not null" json:"gmail"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Profile   *Profile   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"profile,omitempty"`
	Bookmarks []Bookmark `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"bookmarks,omitempty"`
}

// Tag represents a tag for categorization
type Tag struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"size:50;uniqueIndex;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Profile represents a user's profile information
type Profile struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	TagID     uint      `json:"tag_id"`
	Bio       string    `gorm:"type:text" json:"bio"`
	Age       int       `json:"age"`
	Location  string    `gorm:"size:100" json:"location"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Foreign key relationships
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Tag  Tag  `gorm:"foreignKey:TagID;constraint:OnDelete:SET NULL" json:"tag,omitempty"`
}

// Matching represents a match between two users
type Matching struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	User1ID   uint      `gorm:"not null" json:"user1_id"`
	User2ID   uint      `gorm:"not null" json:"user2_id"`
	Status    string    `gorm:"default:'pending'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Foreign key relationships
	User1 User `gorm:"foreignKey:User1ID;constraint:OnDelete:CASCADE" json:"-"`
	User2 User `gorm:"foreignKey:User2ID;constraint:OnDelete:CASCADE" json:"-"`
}

// Bookmark represents a user's bookmark of another user
type Bookmark struct {
	ID               uint      `gorm:"primarykey" json:"id"`
	UserID           uint      `gorm:"not null" json:"user_id"`
	BookmarkedUserID uint      `gorm:"not null" json:"bookmarked_user_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// Foreign key relationships
	User           User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	BookmarkedUser User `gorm:"foreignKey:BookmarkedUserID;constraint:OnDelete:CASCADE" json:"bookmarked_user,omitempty"`
}

// Models returns all user-related models for migration
var Models = []interface{}{
	&User{},
	&Tag{},
	&Profile{},
	&Matching{},
	&Bookmark{},
}

// AutoMigrate performs auto-migration for user models
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(Models...)
}