package models

import (
	"time"
)

// AISettings AI设置模型
type AISettings struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID      string    `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	Provider    string    `gorm:"type:text;not null;default:'openai'" json:"provider"`
	APIKey      string    `gorm:"type:text" json:"api_key,omitempty"`
	BaseURL     string    `gorm:"type:text" json:"base_url,omitempty"`
	Model       string    `gorm:"type:text;not null;default:'gpt-3.5-turbo'" json:"model"`
	Temperature float64   `gorm:"default:0.3" json:"temperature"`
	CreatedAt   time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:now()" json:"updated_at"`

}