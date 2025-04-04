package models

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	UID       string    `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex;not null"`
	UserUID   string    `gorm:"type:uuid;not null;references:users.uid"`
	User      User      `gorm:"foreignKey:UserUID;references:UID"`
	ExpiresAt time.Time `gorm:"not null"`
	TokenHash string    `gorm:"not null;uniqueIndex"`
	UserAgent string    `gorm:"not null"`
	IPAddress string    `gorm:"not null"`
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (s *Session) SetToken(token string) {
	hash := sha256.Sum256([]byte(token))
	s.TokenHash = hex.EncodeToString(hash[:])
}

func (s *Session) ValidateToken(token string) bool {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:]) == s.TokenHash
}

func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.UID == "" {
		s.UID = uuid.New().String()
	}
	if s.ExpiresAt.IsZero() {
		s.ExpiresAt = time.Now().Add(30 * 24 * time.Hour)
	}
	return nil
}

func (Session) TableName() string {
	return "auth.sessions"
}
