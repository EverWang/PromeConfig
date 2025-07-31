# PromeConfig Makefile
# 提供便捷的开发和部署命令

.PHONY: help dev prod build start stop restart logs clean test lint format install

# 默认目标
help: ## 显示帮助信息
	@echo "PromeConfig 项目管理命令"
	@echo ""
	@echo "可用命令:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# 开发环境
dev: ## 启动开发环境
	@echo "启动开发环境..."
	docker-compose up -d postgres redis
	@echo "数据库已启动，请在另一个终端运行:"
	@echo "  后端: cd backend && go run cmd/server/main.go"
	@echo "  前端: npm run dev"

dev-backend: ## 启动后端开发服务
	cd backend && go run cmd/server/main.go

dev-frontend: ## 启动前端开发服务
	npm run dev

# 生产环境
prod: ## 启动生产环境
	@echo "启动生产环境..."
	@if [ -f "backend/.env.production" ]; then \
		cp backend/.env.production backend/.env; \
		echo "已使用生产环境配置"; \
	fi
	docker-compose up -d --build
	@echo "生产环境已启动"
	@echo "访问地址: http://localhost"

build: ## 构建应用镜像
	@echo "构建应用镜像..."
	docker-compose build

start: ## 启动所有服务
	docker-compose up -d

stop: ## 停止所有服务
	@echo "停止所有服务..."
	docker-compose down

restart: stop start ## 重启所有服务

logs: ## 查看服务日志
	docker-compose logs -f

logs-app: ## 查看应用日志
	docker-compose logs -f app

logs-db: ## 查看数据库日志
	docker-compose logs -f postgres

# 清理
clean: ## 清理容器和镜像
	@echo "清理环境..."
	docker-compose down -v --rmi all
	docker system prune -f

clean-volumes: ## 清理数据卷
	@echo "清理数据卷..."
	docker-compose down -v

# 测试
test: ## 运行测试
	@echo "运行后端测试..."
	cd backend && go test ./...
	@echo "运行前端测试..."
	npm test

test-backend: ## 运行后端测试
	cd backend && go test ./...

test-frontend: ## 运行前端测试
	npm test

# 代码质量
lint: ## 代码检查
	@echo "检查后端代码..."
	cd backend && golangci-lint run
	@echo "检查前端代码..."
	npm run lint

format: ## 格式化代码
	@echo "格式化后端代码..."
	cd backend && go fmt ./...
	@echo "格式化前端代码..."
	npm run format

# 安装依赖
install: ## 安装依赖
	@echo "安装前端依赖..."
	npm install
	@echo "安装后端依赖..."
	cd backend && go mod download

install-frontend: ## 安装前端依赖
	npm install

install-backend: ## 安装后端依赖
	cd backend && go mod download

# 数据库操作
db-init: ## 初始化数据库
	@echo "初始化数据库..."
	cd backend && go run cmd/initdb/main.go

db-migrate: ## 运行数据库迁移
	@echo "运行数据库迁移..."
	cd backend && go run cmd/migrate/main.go

db-reset: ## 重置数据库
	@echo "重置数据库..."
	docker-compose down postgres
	docker volume rm promeconfig_postgres_data
	docker-compose up -d postgres
	sleep 5
	make db-init

# 监控
status: ## 查看服务状态
	docker-compose ps

health: ## 检查服务健康状态
	@echo "检查服务健康状态..."
	@curl -f http://localhost:8080/api/health || echo "后端服务不可用"
	@curl -f http://localhost/ || echo "前端服务不可用"

# 备份和恢复
backup: ## 备份数据库
	@echo "备份数据库..."
	mkdir -p backups
	docker-compose exec postgres pg_dump -U postgres promeconfig > backups/backup_$(shell date +%Y%m%d_%H%M%S).sql

restore: ## 恢复数据库 (使用: make restore FILE=backup.sql)
	@if [ -z "$(FILE)" ]; then \
		echo "请指定备份文件: make restore FILE=backup.sql"; \
		exit 1; \
	fi
	@echo "恢复数据库..."
	docker-compose exec -T postgres psql -U postgres promeconfig < $(FILE)

# 部署
deploy: ## 部署到生产环境
	@echo "部署到生产环境..."
	make build
	make prod
	@echo "部署完成"

# 开发工具
shell-backend: ## 进入后端容器shell
	docker-compose exec app sh

shell-db: ## 进入数据库容器shell
	docker-compose exec postgres psql -U postgres promeconfig