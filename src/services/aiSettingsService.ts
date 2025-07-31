import apiClient from './apiClient';
import type { AISettings, CreateAISettingsRequest, UpdateAISettingsRequest } from '../types';

export const aiSettingsService = {
  async getAISettings(): Promise<AISettings | null> {
    try {
      return await apiClient.get<AISettings>('/aisettings');
    } catch (error: any) {
      if (error.message.includes('404') || error.message.includes('not found')) {
        return null;
      }
      throw error;
    }
  },

  async createAISettings(settings: CreateAISettingsRequest): Promise<AISettings> {
    return apiClient.post<AISettings>('/aisettings', settings);
  },

  async updateAISettings(id: string, updates: UpdateAISettingsRequest): Promise<AISettings> {
    return apiClient.put<AISettings>(`/aisettings/${id}`, updates);
  },

  async deleteAISettings(id: string): Promise<void> {
    return apiClient.delete<void>(`/aisettings/${id}`);
  },
};