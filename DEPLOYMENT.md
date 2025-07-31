# PromeConfig 部署指南

本文档介绍如何部署 PromeConfig 应用到生产环境。

## 目录结构

```
PromeConfig/
├── deploy/                 # 部署配置文件
│   ├── nginx.conf         # Nginx配置
│   └── supervisord.conf   # Supervisor配置
├── backend/
│   ├── .env.production    # 生产环境配置
│   └── ...
├── docker-compose.yml     # Docker Compose配置
├── Dockerfile            # 应用镜像构建文件
├── deploy.sh            # Linux/Mac部署脚本
├── deploy.bat           # Windows部署脚本
├── Makefile            # Make命令集合
└── DEPLOYMENT.md       # 本文档
```

## 快速开始

### 1. 环境要求

- Docker 20.10+
- Docker Compose 2.0+
- 至少 2GB 可用内存
- 至少 5GB 可用磁盘空间

### 2. 一键部署

#### Linux/Mac

```bash
# 给脚本执行权限
chmod +x deploy.sh

# 启动生产环境
./deploy.sh prod
```

#### Windows

```cmd
# 启动生产环境
deploy.bat prod
```

#### 使用 Make (推荐)

```bash
# 查看所有可用命令
make help

# 启动生产环境
make prod
```

### 3. 访问应用

部署完成后，可以通过以下地址访问：

- **前端应用**: http://localhost
- **API接口**: http://localhost:8080
- **数据库**: localhost:5432

## 详细部署说明

### 环境配置

#### 生产环境变量

编辑 `backend/.env.production` 文件：

```env
# 生产环境配置
ENV=production
PORT=8080

# 数据库配置
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-secure-password  # 请修改为安全密码
DB_NAME=promeconfig
DB_SSL_MODE=disable

# JWT配置 - 请使用强密钥
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRATION=24h

# Prometheus配置
PROMETHEUS_URL=http://your-prometheus:9090
PROMETHEUS_USERNAME=your-username
PROMETHEUS_PASSWORD=your-password
```

#### 安全配置建议

1. **修改默认密码**: 更改数据库密码和JWT密钥
2. **使用HTTPS**: 在生产环境中配置SSL证书
3. **限制访问**: 配置防火墙规则，只开放必要端口
4. **定期备份**: 设置数据库自动备份

### 服务架构

生产环境包含以下服务：

- **app**: 前端+后端应用容器
- **postgres**: PostgreSQL数据库
- **redis**: Redis缓存（可选）

### 部署命令详解

#### 开发环境

```bash
# 启动开发环境（仅数据库）
./deploy.sh dev
# 或
make dev

# 然后在另一个终端启动服务
cd backend && go run cmd/server/main.go  # 后端
npm run dev                              # 前端
```

#### 生产环境

```bash
# 启动生产环境
./deploy.sh prod
# 或
make prod
```

#### 服务管理

```bash
# 停止服务
./deploy.sh stop
make stop

# 查看日志
./deploy.sh logs
make logs

# 查看服务状态
make status

# 重启服务
make restart
```

#### 清理环境

```bash
# 清理所有容器和数据
./deploy.sh clean
make clean

# 仅清理数据卷
make clean-volumes
```

### 数据库管理

#### 初始化数据库

```bash
make db-init
```

#### 备份数据库

```bash
# 创建备份
make backup

# 恢复备份
make restore FILE=backups/backup_20240101_120000.sql
```

#### 重置数据库

```bash
make db-reset
```

### 监控和维护

#### 健康检查

```bash
# 检查服务健康状态
make health

# 查看服务状态
make status
```

#### 日志管理

```bash
# 查看所有日志
make logs

# 查看应用日志
make logs-app

# 查看数据库日志
make logs-db
```

### 性能优化

#### 资源限制

在 `docker-compose.yml` 中添加资源限制：

```yaml
services:
  app:
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 1G
        reservations:
          cpus: '1.0'
          memory: 512M
```

#### 缓存配置

Redis 缓存已包含在配置中，可以用于：
- 会话存储
- API响应缓存
- 临时数据存储

### 故障排除

#### 常见问题

1. **端口冲突**
   ```bash
   # 检查端口占用
   netstat -tulpn | grep :80
   netstat -tulpn | grep :8080
   ```

2. **容器启动失败**
   ```bash
   # 查看详细日志
   docker-compose logs app
   ```

3. **数据库连接失败**
   ```bash
   # 检查数据库状态
   docker-compose exec postgres pg_isready -U postgres
   ```

4. **权限问题**
   ```bash
   # 给脚本执行权限
   chmod +x deploy.sh
   ```

#### 调试模式

```bash
# 以调试模式启动
docker-compose up --build

# 进入容器调试
make shell-backend
make shell-db
```

### 更新部署

#### 应用更新

```bash
# 拉取最新代码
git pull

# 重新构建并部署
make deploy
```

#### 滚动更新

```bash
# 构建新镜像
make build

# 重启服务
make restart
```

### 扩展配置

#### 负载均衡

可以通过修改 `docker-compose.yml` 添加多个应用实例：

```yaml
services:
  app:
    deploy:
      replicas: 3
```

#### 外部数据库

如果使用外部数据库，修改环境变量：

```env
DB_HOST=your-external-db-host
DB_PORT=5432
DB_SSL_MODE=require
```

### 安全最佳实践

1. **使用强密码**: 数据库和JWT密钥
2. **启用SSL**: 配置HTTPS证书
3. **网络隔离**: 使用Docker网络隔离
4. **定期更新**: 保持镜像和依赖最新
5. **监控日志**: 设置日志监控和告警
6. **备份策略**: 定期备份数据库

### 支持

如果遇到问题，请：

1. 查看日志: `make logs`
2. 检查服务状态: `make status`
3. 运行健康检查: `make health`
4. 查看本文档的故障排除部分

---

**注意**: 在生产环境中，请确保修改所有默认密码和密钥，并配置适当的安全措施。