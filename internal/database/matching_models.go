package database

import (
	"time"

	"gorm.io/gorm"
)

// User ユーザーマスタモデル
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Gmail     string    `gorm:"uniqueIndex;not null;size:255" json:"gmail"`
	Name      string    `gorm:"not null;size:100" json:"name"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
	IconURL   *string   `gorm:"type:text" json:"icon_url,omitempty"`

	// リレーション
	Profile   *Profile   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"profile,omitempty"`
	Matchings []Matching `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"matchings,omitempty"`
	Bookmarks []Bookmark `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"bookmarks,omitempty"`
	
	// 逆参照: このユーザーをブックマークしているユーザー
	BookmarkedBy []Bookmark `gorm:"foreignKey:BookmarkedUserID;constraint:OnDelete:CASCADE" json:"bookmarked_by,omitempty"`
}

// Tag タグマスタモデル
type Tag struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"uniqueIndex;not null;size:50" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// リレーション
	Profiles []Profile `gorm:"foreignKey:TagID;constraint:OnDelete:SET NULL" json:"profiles,omitempty"`
}

// Profile ユーザープロフィールモデル
type Profile struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	Bio       *string   `gorm:"type:text" json:"bio,omitempty"`
	TagID     *uint     `json:"tag_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// リレーション
	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Tag  *Tag  `gorm:"foreignKey:TagID;constraint:OnDelete:SET NULL" json:"tag,omitempty"`
}

// Matching マッチング情報モデル
type Matching struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	NotifyID  *int      `gorm:"index" json:"notify_id,omitempty"`
	Content   *string   `gorm:"type:text" json:"content,omitempty"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// リレーション
	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

// Bookmark ブックマークモデル
type Bookmark struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           uint      `gorm:"index;not null" json:"user_id"`
	BookmarkedUserID uint      `gorm:"index;not null" json:"bookmarked_user_id"`
	CreatedAt        time.Time `gorm:"index" json:"created_at"`

	// リレーション
	User           *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	BookmarkedUser *User `gorm:"foreignKey:BookmarkedUserID;constraint:OnDelete:CASCADE" json:"bookmarked_user,omitempty"`
}

// BeforeCreate ユーザー作成前のフック
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	return nil
}

// BeforeCreate タグ作成前のフック
func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if t.CreatedAt.IsZero() {
		t.CreatedAt = now
	}
	if t.UpdatedAt.IsZero() {
		t.UpdatedAt = now
	}
	return nil
}

// BeforeUpdate タグ更新前のフック
func (t *Tag) BeforeUpdate(tx *gorm.DB) error {
	t.UpdatedAt = time.Now()
	return nil
}

// BeforeCreate プロフィール作成前のフック
func (p *Profile) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = now
	}
	return nil
}

// BeforeUpdate プロフィール更新前のフック
func (p *Profile) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}

// BeforeCreate マッチング作成前のフック
func (m *Matching) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}
	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = now
	}
	return nil
}

// BeforeUpdate マッチング更新前のフック
func (m *Matching) BeforeUpdate(tx *gorm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}

// BeforeCreate ブックマーク作成前のフック
func (b *Bookmark) BeforeCreate(tx *gorm.DB) error {
	if b.CreatedAt.IsZero() {
		b.CreatedAt = time.Now()
	}
	return nil
}

// TableName は、構造体に対応するテーブル名を返す
func (User) TableName() string {
	return "users"
}

func (Tag) TableName() string {
	return "tags"
}

func (Profile) TableName() string {
	return "profiles"
}

func (Matching) TableName() string {
	return "matchings"
}

func (Bookmark) TableName() string {
	return "bookmarks"
}

// マッチングアプリ用のモデルスライス（マイグレーション用）
var MatchingAppModels = []interface{}{
	&User{},
	&Tag{},
	&Profile{},
	&Matching{},
	&Bookmark{},
}