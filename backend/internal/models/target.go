package models

import (
	"time"

	"gorm.io/datatypes"
)

// RelabelConfig 重标签配置
type RelabelConfig struct {
	SourceLabels []string `json:"source_labels,omitempty"`
	Separator    string   `json:"separator,omitempty"`
	TargetLabel  string   `json:"target_label,omitempty"`
	Regex        string   `json:"regex,omitempty"`
	Modulus      int      `json:"modulus,omitempty"`
	Replacement  string   `json:"replacement,omitempty"`
	Action       string   `json:"action,omitempty"`
}

// Target 监控目标模型
type Target struct {
	ID                     string           `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID                 string           `gorm:"type:uuid;not null" json:"user_id"`
	JobName                string           `gorm:"not null" json:"job_name"`
	Targets                datatypes.JSON   `gorm:"type:jsonb;not null;default:'[]'" json:"targets"`
	ScrapeInterval         string           `gorm:"not null;default:'15s'" json:"scrape_interval"`
	MetricsPath            string           `gorm:"not null;default:'/metrics'" json:"metrics_path"`
	RelabelConfigs         datatypes.JSON   `gorm:"type:jsonb" json:"relabel_configs,omitempty"`
	MetricRelabelConfigs   datatypes.JSON   `gorm:"type:jsonb" json:"metric_relabel_configs,omitempty"`
	CreatedAt              time.Time        `gorm:"default:now()" json:"created_at"`
	UpdatedAt              time.Time        `gorm:"default:now()" json:"updated_at"`
}