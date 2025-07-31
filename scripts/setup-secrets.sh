#!/bin/bash

# PromeConfig 密钥设置脚本
# 用于生成和管理生产环境的密钥

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 生成随机密码
generate_password() {
    local length=${1:-32}
    openssl rand -base64 $length | tr -d "=+/" | cut -c1-$length
}

# 生成JWT密钥
generate_jwt_secret() {
    openssl rand -hex 64
}

# 创建密钥目录
setup_secrets_directory() {
    print_info "创建密钥目录..."
    
    local secrets_dir="./secrets"
    
    if [ ! -d "$secrets_dir" ]; then
        mkdir -p "$secrets_dir"
        chmod 700 "$secrets_dir"
        print_success "密钥目录已创建: $secrets_dir"
    else
        print_warning "密钥目录已存在: $secrets_dir"
    fi
}

# 生成数据库密码
setup_db_password() {
    local secrets_file="./secrets/db_password.txt"
    
    if [ ! -f "$secrets_file" ]; then
        print_info "生成数据库密码..."
        generate_password 32 > "$secrets_file"
        chmod 600 "$secrets_file"
        print_success "数据库密码已生成: $secrets_file"
    else
        print_warning "数据库密码文件已存在: $secrets_file"
    fi
}

# 生成JWT密钥
setup_jwt_secret() {
    local secrets_file="./secrets/jwt_secret.txt"
    
    if [ ! -f "$secrets_file" ]; then
        print_info "生成JWT密钥..."
        generate_jwt_secret > "$secrets_file"
        chmod 600 "$secrets_file"
        print_success "JWT密钥已生成: $secrets_file"
    else
        print_warning "JWT密钥文件已存在: $secrets_file"
    fi
}

# 生成Grafana密码
setup_grafana_password() {
    local secrets_file="./secrets/grafana_password.txt"
    
    if [ ! -f "$secrets_file" ]; then
        print_info "生成Grafana密码..."
        generate_password 24 > "$secrets_file"
        chmod 600 "$secrets_file"
        print_success "Grafana密码已生成: $secrets_file"
    else
        print_warning "Grafana密码文件已存在: $secrets_file"
    fi
}

# 生成Redis密码
setup_redis_password() {
    local secrets_file="./secrets/redis_password.txt"
    
    if [ ! -f "$secrets_file" ]; then
        print_info "生成Redis密码..."
        generate_password 32 > "$secrets_file"
        chmod 600 "$secrets_file"
        print_success "Redis密码已生成: $secrets_file"
    else
        print_warning "Redis密码文件已存在: $secrets_file"
    fi
}

# 更新环境变量文件
update_env_file() {
    local env_file="./backend/.env.production"
    
    if [ -f "$env_file" ]; then
        print_info "更新生产环境配置文件..."
        
        # 读取生成的密码
        local db_password=$(cat ./secrets/db_password.txt)
        local jwt_secret=$(cat ./secrets/jwt_secret.txt)
        
        # 更新配置文件
        sed -i "s/DB_PASSWORD=.*/DB_PASSWORD=$db_password/" "$env_file"
        sed -i "s/JWT_SECRET=.*/JWT_SECRET=$jwt_secret/" "$env_file"
        
        print_success "生产环境配置文件已更新"
    else
        print_warning "生产环境配置文件不存在: $env_file"
    fi
}

# 生成SSL证书（自签名，仅用于测试）
generate_ssl_cert() {
    local ssl_dir="./deploy/ssl"
    local cert_file="$ssl_dir/cert.pem"
    local key_file="$ssl_dir/key.pem"
    
    if [ ! -d "$ssl_dir" ]; then
        mkdir -p "$ssl_dir"
    fi
    
    if [ ! -f "$cert_file" ] || [ ! -f "$key_file" ]; then
        print_info "生成自签名SSL证书（仅用于测试）..."
        
        openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
            -keyout "$key_file" \
            -out "$cert_file" \
            -subj "/C=CN/ST=Beijing/L=Beijing/O=PromeConfig/OU=IT/CN=localhost"
        
        chmod 600 "$key_file"
        chmod 644 "$cert_file"
        
        print_success "SSL证书已生成: $ssl_dir"
        print_warning "注意: 这是自签名证书，仅用于测试。生产环境请使用有效的SSL证书。"
    else
        print_warning "SSL证书已存在: $ssl_dir"
    fi
}

# 显示密钥信息
show_secrets_info() {
    print_info "密钥信息摘要:"
    echo ""
    
    if [ -f "./secrets/db_password.txt" ]; then
        echo "数据库密码: $(head -c 8 ./secrets/db_password.txt)..."
    fi
    
    if [ -f "./secrets/jwt_secret.txt" ]; then
        echo "JWT密钥: $(head -c 16 ./secrets/jwt_secret.txt)..."
    fi
    
    if [ -f "./secrets/grafana_password.txt" ]; then
        echo "Grafana密码: $(head -c 8 ./secrets/grafana_password.txt)..."
    fi
    
    if [ -f "./secrets/redis_password.txt" ]; then
        echo "Redis密码: $(head -c 8 ./secrets/redis_password.txt)..."
    fi
    
    echo ""
    print_warning "完整密钥存储在 ./secrets/ 目录中，请妥善保管！"
}

# 备份密钥
backup_secrets() {
    local backup_dir="./backups/secrets_$(date +%Y%m%d_%H%M%S)"
    
    if [ -d "./secrets" ]; then
        print_info "备份密钥到: $backup_dir"
        mkdir -p "$backup_dir"
        cp -r ./secrets/* "$backup_dir/"
        chmod -R 600 "$backup_dir"
        print_success "密钥备份完成"
    else
        print_error "密钥目录不存在，无法备份"
    fi
}

# 清理密钥
clean_secrets() {
    print_warning "这将删除所有生成的密钥文件，确定要继续吗？(y/N)"
    read -r response
    
    if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
        print_info "清理密钥文件..."
        rm -rf ./secrets
        rm -rf ./deploy/ssl
        print_success "密钥文件已清理"
    else
        print_info "取消清理操作"
    fi
}

# 显示帮助信息
show_help() {
    echo "PromeConfig 密钥管理脚本"
    echo ""
    echo "使用方法: $0 [命令]"
    echo ""
    echo "命令:"
    echo "  setup     生成所有密钥和证书"
    echo "  backup    备份现有密钥"
    echo "  clean     清理所有密钥"
    echo "  info      显示密钥信息"
    echo "  ssl       仅生成SSL证书"
    echo "  help      显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 setup    # 生成所有密钥"
    echo "  $0 backup   # 备份密钥"
    echo "  $0 info     # 查看密钥信息"
}

# 主函数
main() {
    case "${1:-setup}" in
        "setup")
            print_info "开始设置密钥..."
            setup_secrets_directory
            setup_db_password
            setup_jwt_secret
            setup_grafana_password
            setup_redis_password
            generate_ssl_cert
            update_env_file
            show_secrets_info
            print_success "密钥设置完成！"
            ;;
        "backup")
            backup_secrets
            ;;
        "clean")
            clean_secrets
            ;;
        "info")
            show_secrets_info
            ;;
        "ssl")
            generate_ssl_cert
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

# 检查依赖
if ! command -v openssl &> /dev/null; then
    print_error "openssl 未安装，请先安装 openssl"
    exit 1
fi

# 执行主函数
main "$@"