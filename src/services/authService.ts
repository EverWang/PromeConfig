const API_BASE_URL = '/api';

export interface User {
  id: string;
  email: string;
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
}

class AuthService {
  private token: string | null = null;
  private user: User | null = null;

  constructor() {
    // 从localStorage恢复token和用户信息
    this.token = localStorage.getItem('auth_token');
    const savedUser = localStorage.getItem('auth_user');
    if (savedUser) {
      try {
        this.user = JSON.parse(savedUser);
      } catch (e) {
        localStorage.removeItem('auth_user');
      }
    }
  }

  async login(credentials: LoginRequest): Promise<AuthResponse> {
    console.log('🔐 开始登录请求:', {
      url: `${API_BASE_URL}/auth/login`,
      email: credentials.email,
      passwordLength: credentials.password.length
    });

    try {
      const response = await fetch(`${API_BASE_URL}/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(credentials),
      });

      console.log('📡 收到响应:', {
        status: response.status,
        statusText: response.statusText,
        headers: Object.fromEntries(response.headers.entries())
      });

      if (!response.ok) {
        let errorMessage = 'Login failed';
        try {
          const error = await response.json();
          errorMessage = error.error || errorMessage;
          console.error('❌ 登录失败 - 服务器错误:', error);
        } catch (e) {
          console.error('❌ 登录失败 - 无法解析错误响应:', e);
        }
        throw new Error(errorMessage);
      }

      const data: AuthResponse = await response.json();
      console.log('✅ 登录成功:', {
        hasToken: !!data.token,
        tokenLength: data.token?.length,
        userEmail: data.user?.email
      });
      
      // 保存token和用户信息
      this.token = data.token;
      this.user = data.user;
      localStorage.setItem('auth_token', data.token);
      localStorage.setItem('auth_user', JSON.stringify(data.user));

      return data;
    } catch (error) {
      console.error('🚨 登录请求异常:', error);
      throw error;
    }
  }

  async register(credentials: RegisterRequest): Promise<AuthResponse> {
    const response = await fetch(`${API_BASE_URL}/auth/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(credentials),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Registration failed');
    }

    const data: AuthResponse = await response.json();
    
    // 保存token和用户信息
    this.token = data.token;
    this.user = data.user;
    localStorage.setItem('auth_token', data.token);
    localStorage.setItem('auth_user', JSON.stringify(data.user));

    return data;
  }

  async logout(): Promise<void> {
    try {
      if (this.token) {
        await fetch(`${API_BASE_URL}/auth/logout`, {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${this.token}`,
          },
        });
      }
    } catch (e) {
      // 忽略登出API错误
    } finally {
      // 清除本地存储
      this.token = null;
      this.user = null;
      localStorage.removeItem('auth_token');
      localStorage.removeItem('auth_user');
    }
  }

  async getCurrentUser(): Promise<User | null> {
    if (!this.token) {
      return null;
    }

    try {
      const response = await fetch(`${API_BASE_URL}/auth/user`, {
        headers: {
          'Authorization': `Bearer ${this.token}`,
        },
      });

      if (!response.ok) {
        // Token可能已过期，清除本地存储
        this.logout();
        return null;
      }

      const user: User = await response.json();
      this.user = user;
      localStorage.setItem('auth_user', JSON.stringify(user));
      return user;
    } catch (e) {
      this.logout();
      return null;
    }
  }

  getToken(): string | null {
    return this.token;
  }

  getUser(): User | null {
    return this.user;
  }

  isAuthenticated(): boolean {
    return this.token !== null && this.user !== null;
  }

  // 获取带认证头的fetch选项
  getAuthHeaders(): HeadersInit {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }

    return headers;
  }
}

export const authService = new AuthService();
export default authService;