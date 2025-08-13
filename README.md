# PromeConfig

一个现代化的 Prometheus 配置管理 Web 应用程序，提供直观的界面来管理监控目标、告警规则和配置文件。

## ✨ 功能特性

### 🎯 核心功能
- **监控目标管理** - 添加、编辑和删除 Prometheus 抓取目标
- **告警规则管理** - 创建和管理 Prometheus 告警规则
- **配置预览** - 实时预览生成的 prometheus.yml 和 alerts.yml 配置文件
- **API 管理** - 连接到 Prometheus 服务器并管理配置重载
- **用户认证** - 安全的用户注册和登录系统

### 🤖 AI 增强功能
- **AI 告警生成** - 使用 AI 根据自然语言描述生成告警规则
- **多 AI 提供商支持** - 支持 OpenAI、Azure OpenAI、Anthropic Claude 等
- **智能配置建议** - AI 辅助的配置优化建议

### 🎨 用户体验
- **现代化 UI** - 基于 Tailwind CSS 的响应式设计
- **实时验证** - 配置文件语法验证和错误提示
- **导出功能** - 一键下载配置文件
- **暗色主题** - 专业的暗色界面设计

## 🚀 在线演示

访问在线演示：[https://graceful-figolla-1821dc.netlify.app](https://graceful-figolla-1821dc.netlify.app)

## 🛠️ 技术栈

### 前端
- **React 18** - 现代化的用户界面框架
- **TypeScript** - 类型安全的 JavaScript
- **Tailwind CSS** - 实用优先的 CSS 框架
- **Vite** - 快速的构建工具
- **Lucide React** - 美观的图标库

### 后端选项
- **Supabase** - 现代化的 BaaS 平台（推荐）
- **Go + Gin** - 高性能的 REST API 服务器（可选）

### 数据库
- **PostgreSQL** - 通过 Supabase 或独立部署

## 📦 快速开始

### 方式一：使用 Supabase（推荐）

1. **克隆项目**
   ```bash
   git clone <repository-url>
   cd promeconfig
   ```

2. **安装依赖**
   ```bash
   npm install
   ```

3. **配置环境变量**
   ```bash
   cp .env.example .env
   ```
   
   编辑 `.env` 文件：
   ```env
   VITE_API_TYPE=supabase
   VITE_SUPABASE_URL=your_supabase_project_url
   VITE_SUPABASE_ANON_KEY=your_supabase_anon_key
   ```

4. **启动开发服务器**
   ```bash
   npm run dev
   ```

5. **访问应用**
   打开浏览器访问 `http://localhost:5173`

### 方式二：使用 Go 后端

1. **启动 PostgreSQL 数据库**
   ```bash
   docker run --name promeconfig-postgres \
     -e POSTGRES_DB=promeconfig \
     -e POSTGRES_USER=user \
     -e POSTGRES_PASSWORD=password \
     -p 5432:5432 \
     -d postgres:15
   ```

2. **配置后端环境变量**
   ```bash
   cd backend
   cp .env.example .env
   ```

3. **启动后端服务**
   ```bash
   cd backend
   go mod tidy
   go run main.go
   ```

4. **配置前端环境变量**
   ```env
   VITE_API_TYPE=golang
   VITE_GOLANG_API_URL=http://localhost:8080/api
   ```

5. **启动前端**
   ```bash
   npm run dev
   ```

## 🔧 配置说明

### Supabase 配置

1. 在 [Supabase](https://supabase.com) 创建新项目
2. 在项目设置中获取 URL 和 anon key
3. 数据库表会自动创建（通过 RLS 策略）

### AI 功能配置

在应用中点击 "AI设置" 按钮配置：

- **OpenAI**: 需要 API Key
- **Azure OpenAI**: 需要 API Key 和 Base URL
- **Anthropic Claude**: 需要 API Key
- **自定义 API**: 配置自定义端点

### Prometheus 连接配置

在 "API Management" 页面配置：

- **Prometheus URL**: 如 `https://prometheus.example.com:9090`
- **用户名/密码**: 基础认证凭据
- **连接测试**: 验证连接状态

## 📚 使用指南

### 1. 管理监控目标

- 点击 "Targets" 进入目标管理页面
- 添加新的抓取目标，配置：
  - Job 名称
  - 目标地址列表
  - 抓取间隔
  - 指标路径
  - 重标签配置（可选）

### 2. 创建告警规则

- 点击 "Alert Rules" 进入告警管理页面
- 手动创建或使用 AI 生成告警规则
- 配置告警表达式、持续时间、标签和注释

### 3. 预览和导出配置

- 点击 "Config Preview" 查看生成的配置文件
- 验证配置语法
- 下载 prometheus.yml 和 alerts.yml 文件

### 4. 同步到 Prometheus

- 在 "API Management" 页面连接到 Prometheus 服务器
- 同步配置文件到服务器
- 重载 Prometheus 配置

## 🏗️ 项目结构

```
promeconfig/
├── src/
│   ├── components/          # React 组件
│   │   ├── AuthWrapper.tsx
│   │   ├── Dashboard.tsx
│   │   ├── TargetManagement.tsx
│   │   ├── AlertRuleManagement.tsx
│   │   ├── ConfigPreview.tsx
│   │   ├── PrometheusAPI.tsx
│   │   └── Sidebar.tsx
│   ├── lib/                 # 工具库
│   │   ├── supabase.ts
│   │   └── api.ts
│   ├── services/            # API 服务
│   │   └── apiService.ts
│   └── App.tsx
├── backend/                 # Go 后端（可选）
│   ├── internal/
│   │   ├── config/
│   │   ├── database/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── models/
│   └── main.go
└── supabase/
    └── migrations/          # 数据库迁移文件
```

## 🔒 安全特性

- **行级安全 (RLS)** - Supabase 数据库安全策略
- **JWT 认证** - 安全的用户会话管理
- **API 密钥加密** - AI 设置中的敏感信息保护
- **CORS 配置** - 跨域请求安全控制

## 🚀 部署

### Netlify 部署

1. **构建项目**
   ```bash
   npm run build
   ```

2. **部署到 Netlify**
   - 连接 GitHub 仓库
   - 设置构建命令：`npm run build`
   - 设置发布目录：`dist`
   - 配置环境变量

### Docker 部署

```dockerfile
# 前端 Dockerfile
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build
EXPOSE 3000
CMD ["npm", "run", "preview"]
```

### Go 后端部署

```dockerfile
# 后端 Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🆘 支持

如果您遇到问题或有疑问：

1. 查看 [Issues](../../issues) 页面
2. 创建新的 Issue 描述问题
3. 参考文档和示例代码

## 🎯 路线图

- [ ] 支持更多 Prometheus 配置选项
- [ ] 添加配置模板功能
- [ ] 实现配置版本控制
- [ ] 支持多环境配置管理
- [ ] 添加配置导入/导出功能
- [ ] 集成更多监控系统

## 🙏 致谢

感谢以下开源项目：

- [React](https://reactjs.org/)
- [Tailwind CSS](https://tailwindcss.com/)
- [Supabase](https://supabase.com/)
- [Lucide Icons](https://lucide.dev/)
- [Vite](https://vitejs.dev/)

---

⭐ 如果这个项目对您有帮助，请给它一个星标！