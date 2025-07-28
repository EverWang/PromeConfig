import { supabase } from '../lib/supabase';
import type { AISettings } from '../lib/supabase';

export class AISettingsService {
  // 获取用户的AI设置
  static async getAISettings(): Promise<AISettings | null> {
    const { data, error } = await supabase
      .from('ai_settings')
      .select('*')
      .single();

    if (error) {
      if (error.code === 'PGRST116') {
        // 没有找到记录，返回null
        return null;
      }
      console.error('Error fetching AI settings:', error);
      throw error;
    }

    return data;
  }

  // 保存或更新AI设置
  static async saveAISettings(settings: Omit<AISettings, 'id' | 'user_id' | 'created_at' | 'updated_at'>): Promise<AISettings> {
    const { data: { user } } = await supabase.auth.getUser();
    
    if (!user) {
      throw new Error('User not authenticated');
    }

    // 尝试更新现有设置
    const { data: existingData, error: fetchError } = await supabase
      .from('ai_settings')
      .select('id')
      .eq('user_id', user.id)
      .single();

    if (existingData) {
      // 更新现有设置
      const { data, error } = await supabase
        .from('ai_settings')
        .update(settings)
        .eq('user_id', user.id)
        .select()
        .single();

      if (error) {
        console.error('Error updating AI settings:', error);
        throw error;
      }

      return data;
    } else {
      // 创建新设置
      const { data, error } = await supabase
        .from('ai_settings')
        .insert([{
          ...settings,
          user_id: user.id
        }])
        .select()
        .single();

      if (error) {
        console.error('Error creating AI settings:', error);
        throw error;
      }

      return data;
    }
  }

  // 删除AI设置
  static async deleteAISettings(): Promise<void> {
    const { error } = await supabase
      .from('ai_settings')
      .delete()
      .eq('user_id', (await supabase.auth.getUser()).data.user?.id);

    if (error) {
      console.error('Error deleting AI settings:', error);
      throw error;
    }
  }
}