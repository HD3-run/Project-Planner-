package models

import "gorm.io/gorm"

// User represents a registered administrator in the system
type User struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"default:'user'"`
	gorm.Model
}
