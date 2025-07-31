# PromeConfig 后端设置指南

## 前置要求

1. 安装 Docker 和 Docker Compose
2. 确保 Docker 服务正在运行
3. 安装 Go 1.21 或更高版本

## 启动步骤

### 1. 启动 PostgreSQL 数据库

```bash
# 在 backend 目录下运行
docker-compose up -d postgres
```

如果遇到权限问题，请：
- 确保 Docker Desktop 正在运行
- 或者以管理员身份运行 PowerShell/命令提示符

### 2. 验证数据库启动

```bash
# 检查容器状态
docker ps

# 查看数据库日志
docker-compose logs postgres
```

### 3. 启动后端服务

```bash
# 安装依赖
go mod tidy

# 运行服务
go run cmd/server/main.go
```

### 4. 验证服务运行

服务启动后，你应该看到类似以下的输出：
```
Successfully connected to PostgreSQL database
Server starting on port 8080
```

## 数据库信息

- **主机**: localhost
- **端口**: 5432
- **数据库名**: promeconfig
- **用户名**: postgres
- **密码**: postgres123

## 测试用户

系统会自动创建一个测试用户：
- **邮箱**: admin@example.com
- **密码**: password123

## 故障排除

### Docker 相关问题

1. **Docker 服务未启动**
   - 启动 Docker Desktop
   - 或运行 `net start docker`（需要管理员权限）

2. **权限问题**
   - 以管理员身份运行终端
   - 或将用户添加到 docker-users 组

3. **端口冲突**
   - 检查端口 5432 是否被占用
   - 修改 docker-compose.yml 中的端口映射

### Go 相关问题

1. **依赖下载失败**
   ```bash
   go clean -modcache
   go mod tidy
   ```

2. **编译错误**
   - 确保 Go 版本 >= 1.21
   - 检查 GOPATH 和 GOROOT 环境变量

## API 端点

服务启动后，可以访问以下端点：

- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `GET /api/targets` - 获取监控目标
- `GET /api/alert-rules` - 获取告警规则
- `GET /api/prometheus/config` - 获取 Prometheus 配置

## 停止服务

```bash
# 停止后端服务
Ctrl+C

# 停止数据库
docker-compose down
```