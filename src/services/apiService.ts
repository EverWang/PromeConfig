// 统一的API服务，支持Supabase和Golang后端切换
import { supabase } from '../lib/supabase';
import { apiClient, apiConfig } from '../lib/api';
import type { Target, AlertRule, AISettings } from '../lib/supabase';

// 用户认证服务
export class AuthService {
  static async signUp(email: string, password: string) {
    if (apiConfig.type === 'golang' && apiClient) {
      const response = await apiClient.signUp(email, password);
      return { user: response.user, error: null };
    } else {
      if (!supabase) throw new Error('Supabase not initialized');
      const { data, error } = await supabase.auth.signUp({ email, password });
      return { user: data.user, error };
    }
  }

  static async signIn(email: string, password: string) {
    if (apiConfig.type === 'golang' && apiClient) {
      const response = await apiClient.signIn(email, password);
      return { user: response.user, error: null };
    } else {
      if (!supabase) throw new Error('Supabase not initialized');
      const { data, error } = await supabase.auth.signInWithPassword({ email, password });
      return { user: data.user, error };
    }
  }

  static async signOut() {
    if (apiConfig.type === 'golang' && apiClient) {
      await apiClient.signOut();
      return { error: null };
    } else {
      if (!supabase) throw new Error('Supabase not initialized');
      const { error } = await supabase.auth.signOut();
      return { error };
    }
  }

  static async getUser() {
    if (apiConfig.type === 'golang' && apiClient) {
      try {
        const response = await apiClient.getUser();
        return { user: response.user, error: null };
      } catch (error) {
        return { user: null, error };
      }
    } else {
      if (!supabase) throw new Error('Supabase not initialized');
      const { data, error } = await supabase.auth.getUser();
      return { user: data.user, error };
    }
  }

  static onAuthStateChange(callback: (user: any) => void) {
    if (apiConfig.type === 'golang') {
      // 对于Golang后端，我们需要手动检查token状态
      const token = localStorage.getItem('access_token');
      if (token) {
        // 验证token并获取用户信息
        this.getUser().then(({ user }) => {
          callback(user);
        }).catch(() => {
          callback(null);
        });
      } else {
        callback(null);
      }
      
      // 返回一个模拟的subscription对象
      return {
        data: {
          subscription: {
            unsubscribe: () => {}
          }
        }
      };
    } else {
      if (!supabase) {
        callback(null);
        return {
          data: {
            subscription: {
              unsubscribe: () => {}
            }
          }
        };
      }
      
      return supabase.auth.onAuthStateChange((event, session) => {
        callback(session?.user ?? null);
      });
    }
  }
}

// Target服务
export class TargetService {
  static async getTargets(): Promise<Target[]> {
    if (apiConfig.type === 'golang' && apiClient) {
      return await apiClient.getTargets();
    } else {
      if (!supabase) throw new Error('Supabase not initialized');
      const { data, error } = await supabase
        .from('targets')
        .select('*')
        .order('created_at', { ascending: false });

      if (error) throw error;
      return data || [];
    }
  }

  static async createTarget(targetData: Omit<Target, 'id' | 'user_id' | 'created_at' | 'updated_at'>): Promise<Target> {
    if (apiConfig.type === 'golang' && apiClient) {
      return await apiClient.createTarget(targetData);
    } else {
      if (!supabase) throw new Error('Supabase not initialized');
      const { data: { user } } = await supabase.auth.getUser();
      if (!user) throw new Error('User not authenticated');

      const { data, error } = await supabase
        .from('targets')
        .insert([{ ...targetData, user_id: user.id }])
        .select()
        .single();

      if (error) throw error;
      return data;
    }
  }

  static async updateTarget(id: string, targetData: Partial<Omit<Target, 'id' | 'user_id' | 'created_at' | 'updated_at'>>): Promise<Target> {
    if (apiConfig.type === 'golang' && apiClient) {
      return await apiClient.updateTarget(id, targetData);
    } else {
      if (!supabase) throw new Error('Supabase not initialized');
      const { data, error } = await supabase
        .from('targets')
        .update(targetData)
        .eq('id', id)
        .select()
        .single();

      if (error) throw error;
      return data;
    }
  }

  static async deleteTarget(id: string): Promise<void> {
    if (apiConfig.type === 'golang' && apiClient) {
      await apiClient.deleteTarget(id);
    } else {
      if (!supabase) throw new Error('Supabase not initialized');
      const { error } = await supabase
        .from('targets')
        .delete()
        .eq('id', id);

      if (error) throw error;
    }
  }
}

// AlertRule服务
export class AlertRuleService {
  static async getAlertRules(): Promise<AlertRule[]> {
    if (apiConfig.type === 'golang' && apiClient) {
      return await apiClient.getAlertRules();
    } else {
      if (!supabase) throw new Error('Supabase not initialized');
      const { data, error } = await supabase
        .from('alert_rules')
        .select('*')
        .order('created_at', { ascending: false });

      if (error) throw error;
      return data || [];
    }
  }

  static async createAlertRule(ruleData: Omit<AlertRule, 'id' | 'user_id' | 'created_at' | 'updated_at'>): Promise<AlertRule> {
    if (apiConfig.type === 'golang' && apiClient) {
      return await apiClient.createAlertRule(ruleData);
    } else {
      if (!supabase) throw new Error('Supabase not initialized');
      const { data: { user } } = await supabase.auth.getUser();
      if (!user) throw new Error('User not authenticated');

      const { data, error } = await supabase
        .from('alert_rules')
        .insert([{ ...ruleData, user_id: user.id }])
        .select()
        .single();

      if (error) throw error;
      return data;
    }
  }

  static async updateAlertRule(id: string, ruleData: Partial<Omit<AlertRule, 'id' | 'user_id' | 'created_at' | 'updated_at'>>): Promise<AlertRule> {
    if (apiConfig.type === 'golang' && apiClient) {
      return await apiClient.updateAlertRule(id, ruleData);
    } else {
      if (!supabase) throw new Error('Supabase not initialized');
      const { data, error } = await supabase
        .from('alert_rules')
        .update(ruleData)
        .eq('id', id)
        .select()
        .single();

      if (error) throw error;
      return data;
    }
  }

  static async deleteAlertRule(id: string): Promise<void> {
    if (apiConfig.type === 'golang' && apiClient) {
      await apiClient.deleteAlertRule(id);
    } else {
      if (!supabase) throw new Error('Supabase not initialized');
      const { error } = await supabase
        .from('alert_rules')
        .delete()
        .eq('id', id);

      if (error) throw error;
    }
  }
}

// AISettings服务
export class AISettingsService {
  static async getAISettings(): Promise<AISettings | null> {
    if (apiConfig.type === 'golang' && apiClient) {
      try {
        return await apiClient.getAISettings();
      } catch (error) {
        return null;
      }
    } else {
      if (!supabase) throw new Error('Supabase not initialized');
      const { data, error } = await supabase
        .from('ai_settings')
        .select('*')
        .maybeSingle();

      if (error) throw error;
      return data;
    }
  }

  static async saveAISettings(settings: Omit<AISettings, 'id' | 'user_id' | 'created_at' | 'updated_at'>): Promise<AISettings> {
    if (apiConfig.type === 'golang' && apiClient) {
      return await apiClient.saveAISettings(settings);
    } else {
      if (!supabase) throw new Error('Supabase not initialized');
      const { data: { user } } = await supabase.auth.getUser();
      if (!user) throw new Error('User not authenticated');

      // 尝试更新现有设置
      const { data: existingData } = await supabase
        .from('ai_settings')
        .select('id')
        .eq('user_id', user.id)
        .maybeSingle();

      if (existingData) {
        const { data, error } = await supabase
          .from('ai_settings')
          .update(settings)
          .eq('user_id', user.id)
          .select()
          .single();

        if (error) throw error;
        return data;
      } else {
        if (!supabase) throw new Error('Supabase not initialized');
        const { data, error } = await supabase
          .from('ai_settings')
          .insert([{ ...settings, user_id: user.id }])
          .select()
          .single();

        if (error) throw error;
        return data;
      }
    }
  }

  static async deleteAISettings(): Promise<void> {
    if (apiConfig.type === 'golang' && apiClient) {
      await apiClient.deleteAISettings();
    } else {
      if (!supabase) throw new Error('Supabase not initialized');
      const { data: { user } } = await supabase.auth.getUser();
      if (!user) throw new Error('User not authenticated');

      const { error } = await supabase
        .from('ai_settings')
        .delete()
        .eq('user_id', user.id);

      if (error) throw error;
    }
  }
}