// API配置和客户端
export interface ApiConfig {
  baseUrl: string;
  type: 'supabase' | 'golang';
}

// 从环境变量获取API配置
const getApiConfig = (): ApiConfig => {
  const apiType = import.meta.env.VITE_API_TYPE || 'supabase';
  const baseUrl = apiType === 'golang' 
    ? import.meta.env.VITE_GOLANG_API_URL || 'http://localhost:8080/api'
    : import.meta.env.VITE_SUPABASE_URL;

  return {
    baseUrl,
    type: apiType as 'supabase' | 'golang'
  };
};

export const apiConfig = getApiConfig();

// API客户端类
export class ApiClient {
  private baseUrl: string;
  private token: string | null = null;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
    this.loadToken();
  }

  private loadToken() {
    this.token = localStorage.getItem('access_token');
  }

  private saveToken(token: string) {
    this.token = token;
    localStorage.setItem('access_token', token);
  }

  private clearToken() {
    this.token = null;
    localStorage.removeItem('access_token');
  }

  private async request<T>(
    endpoint: string, 
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    if (this.token) {
      headers.Authorization = `Bearer ${this.token}`;
    }

    const response = await fetch(url, {
      ...options,
      headers,
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.error || `HTTP ${response.status}: ${response.statusText}`);
    }

    return response.json();
  }

  // 认证相关
  async signUp(email: string, password: string) {
    const response = await this.request<{ user: any; access_token: string }>('/auth/signup', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
    
    this.saveToken(response.access_token);
    return response;
  }

  async signIn(email: string, password: string) {
    const response = await this.request<{ user: any; access_token: string }>('/auth/signin', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
    
    this.saveToken(response.access_token);
    return response;
  }

  async signOut() {
    try {
      await this.request('/auth/signout', { method: 'POST' });
    } finally {
      this.clearToken();
    }
  }

  async getUser() {
    return this.request<{ user: any }>('/user');
  }

  // Targets相关
  async getTargets() {
    return this.request<any[]>('/targets');
  }

  async createTarget(data: any) {
    return this.request<any>('/targets', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async updateTarget(id: string, data: any) {
    return this.request<any>(`/targets/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async deleteTarget(id: string) {
    return this.request(`/targets/${id}`, { method: 'DELETE' });
  }

  // Alert Rules相关
  async getAlertRules() {
    return this.request<any[]>('/alert-rules');
  }

  async createAlertRule(data: any) {
    return this.request<any>('/alert-rules', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async updateAlertRule(id: string, data: any) {
    return this.request<any>(`/alert-rules/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async deleteAlertRule(id: string) {
    return this.request(`/alert-rules/${id}`, { method: 'DELETE' });
  }

  // AI Settings相关
  async getAISettings() {
    return this.request<any>('/ai-settings');
  }

  async saveAISettings(data: any) {
    return this.request<any>('/ai-settings', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async deleteAISettings() {
    return this.request('/ai-settings', { method: 'DELETE' });
  }
}

// 创建API客户端实例
export const apiClient = apiConfig.type === 'golang' 
  ? new ApiClient(apiConfig.baseUrl)
  : null;