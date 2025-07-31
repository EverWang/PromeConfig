package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID           string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Email        string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	PasswordHash string    `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}