import { createClient } from '@supabase/supabase-js';

const apiType = import.meta.env.VITE_API_TYPE || 'supabase';
const supabaseUrl = import.meta.env.VITE_SUPABASE_URL;
const supabaseAnonKey = import.meta.env.VITE_SUPABASE_ANON_KEY;

// Validate Supabase URL format
const isValidUrl = (url: string) => {
  try {
    new URL(url);
    return true;
  } catch {
    return false;
  }
};

// Only require and validate Supabase environment variables when using Supabase backend
if (apiType === 'supabase') {
  if (!supabaseUrl || !supabaseAnonKey) {
    console.warn('Missing Supabase environment variables');
  } else if (!isValidUrl(supabaseUrl)) {
    console.warn('Invalid Supabase URL format:', supabaseUrl);
  }
}

// Create Supabase client only when using Supabase backend
export const supabase = apiType === 'supabase' && supabaseUrl && supabaseAnonKey && isValidUrl(supabaseUrl)
  ? createClient(supabaseUrl, supabaseAnonKey)
  : null;

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