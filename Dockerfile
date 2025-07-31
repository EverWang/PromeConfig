# 前端构建阶段
FROM node:18-alpine AS frontend-builder

WORKDIR /app

# 复制前端依赖文件
COPY package*.json ./
RUN npm ci --only=production

# 复制前端源代码
COPY . .

# 构建前端应用
RUN npm run build

# 后端构建阶段
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app

# 复制后端依赖文件
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# 复制后端源代码
COPY backend/ .

# 构建后端应用
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# 生产环境镜像
FROM nginx:alpine

# 安装必要的工具
RUN apk --no-cache add ca-certificates supervisor

# 创建应用目录
WORKDIR /app

# 复制后端二进制文件
COPY --from=backend-builder /app/server /app/

# 复制前端构建文件到nginx目录
COPY --from=frontend-builder /app/dist /usr/share/nginx/html

# 复制nginx配置
COPY deploy/nginx.conf /etc/nginx/nginx.conf

# 复制supervisor配置
COPY deploy/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# 复制环境变量文件
COPY backend/.env* /app/

# 暴露端口
EXPOSE 80 8080

# 使用supervisor启动nginx和后端服务
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]