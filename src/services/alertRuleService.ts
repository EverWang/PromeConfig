import apiClient from './apiClient';
import type { AlertRule, CreateAlertRuleRequest, UpdateAlertRuleRequest } from '../types';

export const alertRuleService = {
  async getAlertRules(): Promise<AlertRule[]> {
    return apiClient.get<AlertRule[]>('/alertrules');
  },

  async createAlertRule(alertRule: CreateAlertRuleRequest): Promise<AlertRule> {
    return apiClient.post<AlertRule>('/alertrules', alertRule);
  },

  async updateAlertRule(id: string, updates: UpdateAlertRuleRequest): Promise<AlertRule> {
    return apiClient.put<AlertRule>(`/alertrules/${id}`, updates);
  },

  async deleteAlertRule(id: string): Promise<void> {
    return apiClient.delete<void>(`/alertrules/${id}`);
  },
};