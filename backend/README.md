# PromeConfig Backend

Go + Gin + PostgreSQL 后端服务

## 快速开始

### 1. 安装依赖

```bash
cd backend
go mod tidy
```

### 2. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，配置数据库连接等信息
```

### 3. 启动PostgreSQL数据库

```bash
# 使用Docker启动PostgreSQL
docker run --name promeconfig-postgres \
  -e POSTGRES_DB=promeconfig \
  -e POSTGRES_USER=user \
  -e POSTGRES_PASSWORD=password \
  -p 5432:5432 \
  -d postgres:15
```

### 4. 运行服务

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动

## API文档

### 认证相关

- `POST /api/auth/signup` - 用户注册
- `POST /api/auth/signin` - 用户登录
- `POST /api/auth/signout` - 用户登出
- `GET /api/user` - 获取当前用户信息

### Targets管理

- `GET /api/targets` - 获取所有targets
- `POST /api/targets` - 创建target
- `PUT /api/targets/:id` - 更新target
- `DELETE /api/targets/:id` - 删除target

### Alert Rules管理

- `GET /api/alert-rules` - 获取所有告警规则
- `POST /api/alert-rules` - 创建告警规则
- `PUT /api/alert-rules/:id` - 更新告警规则
- `DELETE /api/alert-rules/:id` - 删除告警规则

### AI Settings管理

- `GET /api/ai-settings` - 获取AI设置
- `POST /api/ai-settings` - 保存AI设置
- `DELETE /api/ai-settings` - 删除AI设置

### Prometheus配置管理

- `POST /api/prometheus/sync` - 同步配置到Prometheus
- `POST /api/prometheus/reload` - 重载Prometheus配置
- `GET /api/prometheus/status` - 获取Prometheus状态

## 数据库结构

数据库会自动创建以下表：

- `users` - 用户表
- `targets` - 监控目标表
- `alert_rules` - 告警规则表
- `ai_settings` - AI设置表

## 部署

### 使用Docker

```bash
# 构建镜像
docker build -t promeconfig-backend .

# 运行容器
docker run -p 8080:8080 --env-file .env promeconfig-backend
```

### 直接部署

```bash
# 构建二进制文件
go build -o promeconfig-backend main.go

# 运行
./promeconfig-backend
```