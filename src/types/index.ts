// 数据库类型定义
export interface Target {
  id: string;
  user_id: string;
  job_name: string;
  targets: string[];
  scrape_interval: string;
  metrics_path: string;
  relabel_configs?: Array<{
    source_labels?: string[];
    separator?: string;
    target_label?: string;
    regex?: string;
    modulus?: number;
    replacement?: string;
    action?: 'replace' | 'keep' | 'drop' | 'hashmod' | 'labelmap' | 'labeldrop' | 'labelkeep';
  }>;
  metric_relabel_configs?: Array<{
    source_labels?: string[];
    separator?: string;
    target_label?: string;
    regex?: string;
    modulus?: number;
    replacement?: string;
    action?: 'replace' | 'keep' | 'drop' | 'hashmod' | 'labelmap' | 'labeldrop' | 'labelkeep';
  }>;
  created_at: string;
  updated_at: string;
}

export interface AlertRule {
  id: string;
  user_id: string;
  alert_name: string;
  expr: string;
  for_duration: string;
  labels: Record<string, string>;
  annotations: Record<string, string>;
  created_at: string;
  updated_at: string;
}

export interface AISettings {
  id: string;
  user_id: string;
  provider: string;
  api_key?: string;
  base_url?: string;
  model: string;
  temperature: number;
  created_at: string;
  updated_at: string;
}

// API请求/响应类型
export interface CreateTargetRequest {
  job_name: string;
  targets: string[];
  scrape_interval: string;
  metrics_path: string;
  relabel_configs?: Target['relabel_configs'];
  metric_relabel_configs?: Target['metric_relabel_configs'];
}

export interface UpdateTargetRequest extends Partial<CreateTargetRequest> {}

export interface CreateAlertRuleRequest {
  alert_name: string;
  expr: string;
  for_duration: string;
  labels?: Record<string, string>;
  annotations?: Record<string, string>;
}

export interface UpdateAlertRuleRequest extends Partial<CreateAlertRuleRequest> {}

export interface CreateAISettingsRequest {
  provider: string;
  api_key?: string;
  base_url?: string;
  model: string;
  temperature: number;
}

export interface UpdateAISettingsRequest extends Partial<CreateAISettingsRequest> {}