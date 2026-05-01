package models

import (
	"time"
)

// Session represents an active refresh token session in the database
type Session struct {
	ID        string    `gorm:"primaryKey;type:varchar(64)" json:"id"` // The unique session ID inside the refresh token
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	User      User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
