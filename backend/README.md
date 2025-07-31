# PromeConfig 后端实现

这是PromeConfig的Golang后端实现，用于替代原有的Supabase实现。

## 技术栈

- Golang
- Gin Web框架
- GORM ORM框架
- PostgreSQL数据库
- JWT认证

## 项目结构

```
backend/
├── cmd/                  # 应用程序入口
│   └── server/           # 服务器入口
│       └── main.go       # 主程序
├── config/               # 配置文件
│   └── config.go         # 配置加载
├── internal/             # 内部包
│   ├── api/              # API处理器
│   │   ├── auth.go       # 认证相关API
│   │   ├── targets.go    # 监控目标API
│   │   ├── alertrules.go # 告警规则API
│   │   ├── aisettings.go # AI设置API
│   │   └── prometheus.go # Prometheus API
│   ├── middleware/       # 中间件
│   │   ├── auth.go       # 认证中间件
│   │   └── cors.go       # CORS中间件
│   ├── models/           # 数据模型
│   │   ├── user.go       # 用户模型
│   │   ├── target.go     # 监控目标模型
│   │   ├── alertrule.go  # 告警规则模型
│   │   └── aisetting.go  # AI设置模型
│   ├── repository/       # 数据访问层
│   │   ├── user.go       # 用户仓库
│   │   ├── target.go     # 监控目标仓库
│   │   ├── alertrule.go  # 告警规则仓库
│   │   └── aisetting.go  # AI设置仓库
│   └── service/          # 业务逻辑层
│       ├── auth.go       # 认证服务
│       ├── target.go     # 监控目标服务
│       ├── alertrule.go  # 告警规则服务
│       ├── aisetting.go  # AI设置服务
│       └── prometheus.go # Prometheus服务
├── pkg/                  # 公共包
│   ├── jwt/              # JWT工具
│   ├── password/         # 密码工具
│   └── prometheus/       # Prometheus客户端
├── migrations/           # 数据库迁移
└── README.md             # 项目说明
```

## API设计

### 认证API

- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/logout` - 用户登出
- `GET /api/auth/user` - 获取当前用户信息

### 监控目标API

- `GET /api/targets` - 获取所有监控目标
- `POST /api/targets` - 创建监控目标
- `PUT /api/targets/:id` - 更新监控目标
- `DELETE /api/targets/:id` - 删除监控目标

### 告警规则API

- `GET /api/alertrules` - 获取所有告警规则
- `POST /api/alertrules` - 创建告警规则
- `PUT /api/alertrules/:id` - 更新告警规则
- `DELETE /api/alertrules/:id` - 删除告警规则

### AI设置API

- `GET /api/aisettings` - 获取AI设置
- `POST /api/aisettings` - 创建或更新AI设置
- `DELETE /api/aisettings` - 删除AI设置

### Prometheus API

- `POST /api/prometheus/reload` - 重新加载Prometheus配置
- `GET /api/prometheus/config` - 获取生成的Prometheus配置
- `POST /api/prometheus/config` - 上传Prometheus配置

## 数据库设计

数据库设计与原Supabase实现保持一致，包括以下表：

- `users` - 用户表
- `targets` - 监控目标表
- `alert_rules` - 告警规则表
- `ai_settings` - AI设置表

## 开发计划

1. 设置项目基础结构
2. 实现数据库模型和迁移
3. 实现认证功能
4. 实现监控目标管理
5. 实现告警规则管理
6. 实现AI设置管理
7. 实现Prometheus API集成
8. 前后端集成测试
9. 部署和文档