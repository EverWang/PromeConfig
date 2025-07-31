package models

import (
	"time"

	"gorm.io/datatypes"
)

// AlertRule 告警规则模型
type AlertRule struct {
	ID          string          `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      string          `json:"user_id" gorm:"type:uuid;not null"`
	AlertName   string          `json:"alert_name" gorm:"column:alert_name;not null"`
	Expr        string          `json:"expr" gorm:"column:expr;not null"`
	ForDuration string          `json:"for_duration" gorm:"column:for_duration;default:'5m'"`
	Severity    string          `json:"severity" gorm:"column:severity;default:'warning'"`
	Labels      datatypes.JSON  `json:"labels" gorm:"type:jsonb;default:'{}'"`
	Annotations datatypes.JSON  `json:"annotations" gorm:"type:jsonb;default:'{}'"`
	CreatedAt   time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}