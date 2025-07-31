import apiClient from './apiClient';
import type { Target, CreateTargetRequest, UpdateTargetRequest } from '../types';

export const targetService = {
  async getTargets(): Promise<Target[]> {
    return apiClient.get<Target[]>('/targets');
  },

  async createTarget(target: CreateTargetRequest): Promise<Target> {
    return apiClient.post<Target>('/targets', target);
  },

  async updateTarget(id: string, updates: UpdateTargetRequest): Promise<Target> {
    return apiClient.put<Target>(`/targets/${id}`, updates);
  },

  async deleteTarget(id: string): Promise<void> {
    return apiClient.delete<void>(`/targets/${id}`);
  },
};