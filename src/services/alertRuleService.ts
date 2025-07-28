import { supabase } from '../lib/supabase';
import type { AlertRule } from '../lib/supabase';

export class AlertRuleService {
  // 获取用户的所有告警规则
  static async getAlertRules(): Promise<AlertRule[]> {
    const { data, error } = await supabase
      .from('alert_rules')
      .select('*')
      .order('created_at', { ascending: false });

    if (error) {
      console.error('Error fetching alert rules:', error);
      throw error;
    }

    return data || [];
  }

  // 创建新的告警规则
  static async createAlertRule(ruleData: Omit<AlertRule, 'id' | 'user_id' | 'created_at' | 'updated_at'>): Promise<AlertRule> {
    const { data: { user } } = await supabase.auth.getUser();
    
    if (!user) {
      throw new Error('User not authenticated');
    }

    const { data, error } = await supabase
      .from('alert_rules')
      .insert([{
        ...ruleData,
        user_id: user.id
      }])
      .select()
      .single();

    if (error) {
      console.error('Error creating alert rule:', error);
      throw error;
    }

    return data;
  }

  // 更新告警规则
  static async updateAlertRule(id: string, ruleData: Partial<Omit<AlertRule, 'id' | 'user_id' | 'created_at' | 'updated_at'>>): Promise<AlertRule> {
    const { data, error } = await supabase
      .from('alert_rules')
      .update(ruleData)
      .eq('id', id)
      .select()
      .single();

    if (error) {
      console.error('Error updating alert rule:', error);
      throw error;
    }

    return data;
  }

  // 删除告警规则
  static async deleteAlertRule(id: string): Promise<void> {
    const { error } = await supabase
      .from('alert_rules')
      .delete()
      .eq('id', id);

    if (error) {
      console.error('Error deleting alert rule:', error);
      throw error;
    }
  }
}