/*
  # Prometheus配置管理数据库架构

  1. 新建表
    - `targets` - 存储Prometheus监控目标配置
      - `id` (uuid, 主键)
      - `user_id` (uuid, 外键到auth.users)
      - `job_name` (text, 任务名称)
      - `targets` (jsonb, 目标列表)
      - `scrape_interval` (text, 抓取间隔)
      - `metrics_path` (text, 指标路径)
      - `relabel_configs` (jsonb, 重标记配置)
      - `metric_relabel_configs` (jsonb, 指标重标记配置)
      - `created_at` (timestamptz)
      - `updated_at` (timestamptz)
    
    - `alert_rules` - 存储告警规则配置
      - `id` (uuid, 主键)
      - `user_id` (uuid, 外键到auth.users)
      - `alert_name` (text, 告警名称)
      - `expr` (text, PromQL表达式)
      - `for_duration` (text, 持续时间)
      - `labels` (jsonb, 标签)
      - `annotations` (jsonb, 注释)
      - `created_at` (timestamptz)
      - `updated_at` (timestamptz)

    - `ai_settings` - 存储AI配置
      - `id` (uuid, 主键)
      - `user_id` (uuid, 外键到auth.users)
      - `provider` (text, API提供商)
      - `api_key` (text, 加密存储)
      - `base_url` (text, API地址)
      - `model` (text, 模型名称)
      - `temperature` (numeric, 温度参数)
      - `created_at` (timestamptz)
      - `updated_at` (timestamptz)

  2. 安全策略
    - 启用所有表的RLS
    - 用户只能访问自己的数据
    - 支持完整的CRUD操作
*/

-- 创建targets表
CREATE TABLE IF NOT EXISTS targets (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid REFERENCES auth.users(id) ON DELETE CASCADE NOT NULL,
  job_name text NOT NULL,
  targets jsonb NOT NULL DEFAULT '[]',
  scrape_interval text NOT NULL DEFAULT '15s',
  metrics_path text NOT NULL DEFAULT '/metrics',
  relabel_configs jsonb DEFAULT NULL,
  metric_relabel_configs jsonb DEFAULT NULL,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

-- 创建alert_rules表
CREATE TABLE IF NOT EXISTS alert_rules (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid REFERENCES auth.users(id) ON DELETE CASCADE NOT NULL,
  alert_name text NOT NULL,
  expr text NOT NULL,
  for_duration text NOT NULL DEFAULT '5m',
  labels jsonb NOT NULL DEFAULT '{}',
  annotations jsonb NOT NULL DEFAULT '{}',
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

-- 创建ai_settings表
CREATE TABLE IF NOT EXISTS ai_settings (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid REFERENCES auth.users(id) ON DELETE CASCADE NOT NULL,
  provider text NOT NULL DEFAULT 'openai',
  api_key text DEFAULT NULL,
  base_url text DEFAULT NULL,
  model text NOT NULL DEFAULT 'gpt-3.5-turbo',
  temperature numeric(3,2) NOT NULL DEFAULT 0.3,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  UNIQUE(user_id)
);

-- 启用行级安全
ALTER TABLE targets ENABLE ROW LEVEL SECURITY;
ALTER TABLE alert_rules ENABLE ROW LEVEL SECURITY;
ALTER TABLE ai_settings ENABLE ROW LEVEL SECURITY;

-- 创建targets表的RLS策略
CREATE POLICY "用户可以查看自己的targets"
  ON targets
  FOR SELECT
  TO authenticated
  USING (auth.uid() = user_id);

CREATE POLICY "用户可以插入自己的targets"
  ON targets
  FOR INSERT
  TO authenticated
  WITH CHECK (auth.uid() = user_id);

CREATE POLICY "用户可以更新自己的targets"
  ON targets
  FOR UPDATE
  TO authenticated
  USING (auth.uid() = user_id)
  WITH CHECK (auth.uid() = user_id);

CREATE POLICY "用户可以删除自己的targets"
  ON targets
  FOR DELETE
  TO authenticated
  USING (auth.uid() = user_id);

-- 创建alert_rules表的RLS策略
CREATE POLICY "用户可以查看自己的alert_rules"
  ON alert_rules
  FOR SELECT
  TO authenticated
  USING (auth.uid() = user_id);

CREATE POLICY "用户可以插入自己的alert_rules"
  ON alert_rules
  FOR INSERT
  TO authenticated
  WITH CHECK (auth.uid() = user_id);

CREATE POLICY "用户可以更新自己的alert_rules"
  ON alert_rules
  FOR UPDATE
  TO authenticated
  USING (auth.uid() = user_id)
  WITH CHECK (auth.uid() = user_id);

CREATE POLICY "用户可以删除自己的alert_rules"
  ON alert_rules
  FOR DELETE
  TO authenticated
  USING (auth.uid() = user_id);

-- 创建ai_settings表的RLS策略
CREATE POLICY "用户可以查看自己的ai_settings"
  ON ai_settings
  FOR SELECT
  TO authenticated
  USING (auth.uid() = user_id);

CREATE POLICY "用户可以插入自己的ai_settings"
  ON ai_settings
  FOR INSERT
  TO authenticated
  WITH CHECK (auth.uid() = user_id);

CREATE POLICY "用户可以更新自己的ai_settings"
  ON ai_settings
  FOR UPDATE
  TO authenticated
  USING (auth.uid() = user_id)
  WITH CHECK (auth.uid() = user_id);

CREATE POLICY "用户可以删除自己的ai_settings"
  ON ai_settings
  FOR DELETE
  TO authenticated
  USING (auth.uid() = user_id);

-- 创建更新时间触发器函数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ language 'plpgsql';

-- 为所有表添加更新时间触发器
CREATE TRIGGER update_targets_updated_at
  BEFORE UPDATE ON targets
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_alert_rules_updated_at
  BEFORE UPDATE ON alert_rules
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_ai_settings_updated_at
  BEFORE UPDATE ON ai_settings
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

-- 创建索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_targets_user_id ON targets(user_id);
CREATE INDEX IF NOT EXISTS idx_targets_job_name ON targets(job_name);
CREATE INDEX IF NOT EXISTS idx_alert_rules_user_id ON alert_rules(user_id);
CREATE INDEX IF NOT EXISTS idx_alert_rules_alert_name ON alert_rules(alert_name);
CREATE INDEX IF NOT EXISTS idx_ai_settings_user_id ON ai_settings(user_id);