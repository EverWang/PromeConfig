package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func Initialize(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func Migrate(db *sql.DB) error {
	migrations := []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,
		
		// Users表
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);`,

		// Targets表
		`CREATE TABLE IF NOT EXISTS targets (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			job_name TEXT NOT NULL,
			targets JSONB NOT NULL DEFAULT '[]'::jsonb,
			scrape_interval TEXT NOT NULL DEFAULT '15s',
			metrics_path TEXT NOT NULL DEFAULT '/metrics',
			relabel_configs JSONB,
			metric_relabel_configs JSONB,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);`,

		// Alert Rules表
		`CREATE TABLE IF NOT EXISTS alert_rules (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			alert_name TEXT NOT NULL,
			expr TEXT NOT NULL,
			for_duration TEXT NOT NULL DEFAULT '5m',
			labels JSONB NOT NULL DEFAULT '{}'::jsonb,
			annotations JSONB NOT NULL DEFAULT '{}'::jsonb,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);`,

		// AI Settings表
		`CREATE TABLE IF NOT EXISTS ai_settings (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			provider TEXT NOT NULL DEFAULT 'openai',
			api_key TEXT,
			base_url TEXT,
			model TEXT NOT NULL DEFAULT 'gpt-3.5-turbo',
			temperature DECIMAL(3,2) NOT NULL DEFAULT 0.3,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);`,

		// 创建索引
		`CREATE INDEX IF NOT EXISTS idx_targets_user_id ON targets(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_targets_job_name ON targets(job_name);`,
		`CREATE INDEX IF NOT EXISTS idx_alert_rules_user_id ON alert_rules(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_alert_rules_alert_name ON alert_rules(alert_name);`,
		`CREATE INDEX IF NOT EXISTS idx_ai_settings_user_id ON ai_settings(user_id);`,

		// 创建更新时间触发器函数
		`CREATE OR REPLACE FUNCTION update_updated_at_column()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = NOW();
			RETURN NEW;
		END;
		$$ language 'plpgsql';`,

		// 为各表创建更新时间触发器
		`DROP TRIGGER IF EXISTS update_users_updated_at ON users;
		CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();`,

		`DROP TRIGGER IF EXISTS update_targets_updated_at ON targets;
		CREATE TRIGGER update_targets_updated_at BEFORE UPDATE ON targets FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();`,

		`DROP TRIGGER IF EXISTS update_alert_rules_updated_at ON alert_rules;
		CREATE TRIGGER update_alert_rules_updated_at BEFORE UPDATE ON alert_rules FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();`,

		`DROP TRIGGER IF EXISTS update_ai_settings_updated_at ON ai_settings;
		CREATE TRIGGER update_ai_settings_updated_at BEFORE UPDATE ON ai_settings FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("failed to execute migration: %w", err)
		}
	}

	return nil
}