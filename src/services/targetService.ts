import { supabase } from '../lib/supabase';
import type { Target } from '../lib/supabase';

export class TargetService {
  // 获取用户的所有targets
  static async getTargets(): Promise<Target[]> {
    const { data, error } = await supabase
      .from('targets')
      .select('*')
      .order('created_at', { ascending: false });

    if (error) {
      console.error('Error fetching targets:', error);
      throw error;
    }

    return data || [];
  }

  // 创建新的target
  static async createTarget(targetData: Omit<Target, 'id' | 'user_id' | 'created_at' | 'updated_at'>): Promise<Target> {
    const { data: { user } } = await supabase.auth.getUser();
    
    if (!user) {
      throw new Error('User not authenticated');
    }

    const { data, error } = await supabase
      .from('targets')
      .insert([{
        ...targetData,
        user_id: user.id,
        targets: targetData.targets
      }])
      .select()
      .single();

    if (error) {
      console.error('Error creating target:', error);
      throw error;
    }

    return data;
  }

  // 更新target
  static async updateTarget(id: string, targetData: Partial<Omit<Target, 'id' | 'user_id' | 'created_at' | 'updated_at'>>): Promise<Target> {
    const { data, error } = await supabase
      .from('targets')
      .update(targetData)
      .eq('id', id)
      .select()
      .single();

    if (error) {
      console.error('Error updating target:', error);
      throw error;
    }

    return data;
  }

  // 删除target
  static async deleteTarget(id: string): Promise<void> {
    const { error } = await supabase
      .from('targets')
      .delete()
      .eq('id', id);

    if (error) {
      console.error('Error deleting target:', error);
      throw error;
    }
  }
}