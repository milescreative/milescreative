package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model         // This includes ID, CreatedAt, UpdatedAt, DeletedAt
	UID         string `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex;not null"`
	Email       string `gorm:"uniqueIndex;not null"`
	Name        string
	Provider    string    `gorm:"not null"`
	ProviderID  string    `gorm:"not null"`
	LastLoginAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	AvatarURL   string
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.UID == "" {
		u.UID = uuid.New().String()
	}
	return nil
}

func (User) TableName() string {
	return "auth.users"
}
