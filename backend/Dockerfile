FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装依赖
COPY go.mod go.sum* ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# 使用轻量级的alpine镜像
FROM alpine:latest

WORKDIR /app

# 安装CA证书
RUN apk --no-cache add ca-certificates

# 从构建阶段复制二进制文件
COPY --from=builder /app/server /app/

# 创建配置目录
RUN mkdir -p /app/configs

# 复制环境变量文件
COPY .env* /app/

# 设置时区
RUN apk --no-cache add tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["/app/server"]