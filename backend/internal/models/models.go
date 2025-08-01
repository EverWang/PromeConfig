package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Target struct {
	ID                    uuid.UUID       `json:"id" db:"id"`
	UserID               uuid.UUID       `json:"user_id" db:"user_id"`
	JobName              string          `json:"job_name" db:"job_name"`
	Targets              json.RawMessage `json:"targets" db:"targets"`
	ScrapeInterval       string          `json:"scrape_interval" db:"scrape_interval"`
	MetricsPath          string          `json:"metrics_path" db:"metrics_path"`
	RelabelConfigs       json.RawMessage `json:"relabel_configs,omitempty" db:"relabel_configs"`
	MetricRelabelConfigs json.RawMessage `json:"metric_relabel_configs,omitempty" db:"metric_relabel_configs"`
	CreatedAt            time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at" db:"updated_at"`
}

type AlertRule struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	UserID      uuid.UUID       `json:"user_id" db:"user_id"`
	AlertName   string          `json:"alert_name" db:"alert_name"`
	Expr        string          `json:"expr" db:"expr"`
	ForDuration string          `json:"for_duration" db:"for_duration"`
	Labels      json.RawMessage `json:"labels" db:"labels"`
	Annotations json.RawMessage `json:"annotations" db:"annotations"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

type AISettings struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Provider    string    `json:"provider" db:"provider"`
	APIKey      *string   `json:"api_key,omitempty" db:"api_key"`
	BaseURL     *string   `json:"base_url,omitempty" db:"base_url"`
	Model       string    `json:"model" db:"model"`
	Temperature float64   `json:"temperature" db:"temperature"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// 请求/响应结构体
type SignUpRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type SignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User        User   `json:"user"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type CreateTargetRequest struct {
	JobName              string          `json:"job_name" binding:"required"`
	Targets              json.RawMessage `json:"targets" binding:"required"`
	ScrapeInterval       string          `json:"scrape_interval"`
	MetricsPath          string          `json:"metrics_path"`
	RelabelConfigs       json.RawMessage `json:"relabel_configs,omitempty"`
	MetricRelabelConfigs json.RawMessage `json:"metric_relabel_configs,omitempty"`
}

type CreateAlertRuleRequest struct {
	AlertName   string          `json:"alert_name" binding:"required"`
	Expr        string          `json:"expr" binding:"required"`
	ForDuration string          `json:"for_duration"`
	Labels      json.RawMessage `json:"labels"`
	Annotations json.RawMessage `json:"annotations"`
}

type SaveAISettingsRequest struct {
	Provider    string  `json:"provider" binding:"required"`
	APIKey      *string `json:"api_key,omitempty"`
	BaseURL     *string `json:"base_url,omitempty"`
	Model       string  `json:"model" binding:"required"`
	Temperature float64 `json:"temperature"`
}