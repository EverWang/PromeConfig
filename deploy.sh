#!/bin/bash

# PromeConfig 一键部署脚本
# 使用方法: ./deploy.sh [dev|prod|stop|logs|clean]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_message() {
    echo -e "${2}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

print_info() {
    print_message "$1" "$BLUE"
}

print_success() {
    print_message "$1" "$GREEN"
}

print_warning() {
    print_message "$1" "$YELLOW"
}

print_error() {
    print_message "$1" "$RED"
}

# 检查Docker和Docker Compose
check_dependencies() {
    print_info "检查依赖..."
    
    if ! command -v docker &> /dev/null; then
        print_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        print_error "Docker Compose 未安装，请先安装 Docker Compose"
        exit 1
    fi
    
    print_success "依赖检查完成"
}

# 创建必要的目录和文件
setup_environment() {
    print_info "设置环境..."
    
    # 创建日志目录
    mkdir -p logs
    
    # 检查环境变量文件
    if [ ! -f "backend/.env" ]; then
        print_warning "backend/.env 文件不存在，从示例文件复制..."
        cp backend/.env.example backend/.env
    fi
    
    print_success "环境设置完成"
}

# 开发环境启动
start_dev() {
    print_info "启动开发环境..."
    
    # 启动数据库
    docker-compose up -d postgres redis
    
    print_success "开发环境数据库已启动"
    print_info "请在另一个终端运行以下命令启动服务:"
    print_info "后端: cd backend && go run cmd/server/main.go"
    print_info "前端: npm run dev"
}

# 生产环境启动
start_prod() {
    print_info "启动生产环境..."
    
    # 复制生产环境配置
    if [ -f "backend/.env.production" ]; then
        cp backend/.env.production backend/.env
        print_info "已使用生产环境配置"
    fi
    
    # 构建并启动所有服务
    docker-compose up -d --build
    
    print_success "生产环境已启动"
    print_info "应用访问地址: http://localhost"
    print_info "API访问地址: http://localhost:8080"
    print_info "数据库访问地址: localhost:5432"
}

# 停止服务
stop_services() {
    print_info "停止所有服务..."
    docker-compose down
    print_success "所有服务已停止"
}

# 查看日志
show_logs() {
    print_info "显示服务日志..."
    docker-compose logs -f
}

# 清理环境
clean_environment() {
    print_warning "这将删除所有容器、镜像和数据卷，确定要继续吗？(y/N)"
    read -r response
    if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
        print_info "清理环境..."
        docker-compose down -v --rmi all
        docker system prune -f
        print_success "环境清理完成"
    else
        print_info "取消清理操作"
    fi
}

# 显示帮助信息
show_help() {
    echo "PromeConfig 部署脚本"
    echo ""
    echo "使用方法: $0 [命令]"
    echo ""
    echo "命令:"
    echo "  dev     启动开发环境（仅数据库）"
    echo "  prod    启动生产环境（完整服务）"
    echo "  stop    停止所有服务"
    echo "  logs    查看服务日志"
    echo "  clean   清理所有容器和数据"
    echo "  help    显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 prod    # 启动生产环境"
    echo "  $0 dev     # 启动开发环境"
    echo "  $0 logs    # 查看日志"
}

# 主函数
main() {
    case "${1:-help}" in
        "dev")
            check_dependencies
            setup_environment
            start_dev
            ;;
        "prod")
            check_dependencies
            setup_environment
            start_prod
            ;;
        "stop")
            stop_services
            ;;
        "logs")
            show_logs
            ;;
        "clean")
            clean_environment
            ;;
        "help")
            show_help
            ;;
        *)
            print_error "未知命令: $1"
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"