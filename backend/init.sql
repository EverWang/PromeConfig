-- 启用UUID扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 删除旧表（如果存在）
DROP TABLE IF EXISTS ai_settings;
DROP TABLE IF EXISTS alert_rules;
DROP TABLE IF EXISTS targets;
DROP TABLE IF EXISTS monitoring_targets;
DROP TABLE IF EXISTS users;

-- 创建用户表
CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- 创建监控目标表
CREATE TABLE targets (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    job_name VARCHAR(255) NOT NULL,
    targets jsonb NOT NULL DEFAULT '[]',
    scrape_interval VARCHAR(50) NOT NULL DEFAULT '15s',
    metrics_path VARCHAR(255) NOT NULL DEFAULT '/metrics',
    relabel_configs jsonb,
    metric_relabel_configs jsonb,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- 创建告警规则表
CREATE TABLE alert_rules (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    expression TEXT NOT NULL,
    duration VARCHAR(50) DEFAULT '5m',
    severity VARCHAR(50) DEFAULT 'warning',
    labels jsonb,
    annotations jsonb,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- 创建AI设置表
CREATE TABLE ai_settings (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    api_key VARCHAR(255) NOT NULL,
    model VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    UNIQUE(user_id)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_targets_user_id ON targets(user_id);
CREATE INDEX IF NOT EXISTS idx_alert_rules_user_id ON alert_rules(user_id);
CREATE INDEX IF NOT EXISTS idx_ai_settings_user_id ON ai_settings(user_id);

-- 插入测试用户（密码是 'password123' 的哈希值）
INSERT INTO users (email, password_hash) VALUES 
('admin@example.com', '$2a$10$gvmN4CVPoWg8L2K.VaIPE.4.ez9/fdbbf3Yy3lylLsbWlwj2ESDWm')
ON CONFLICT (email) DO NOTHING;