@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

REM PromeConfig 一键部署脚本 (Windows版本)
REM 使用方法: deploy.bat [dev|prod|stop|logs|clean]

set "RED=[91m"
set "GREEN=[92m"
set "YELLOW=[93m"
set "BLUE=[94m"
set "NC=[0m"

REM 打印带颜色的消息
:print_message
echo %~2[%date% %time%] %~1%NC%
goto :eof

:print_info
call :print_message "%~1" "%BLUE%"
goto :eof

:print_success
call :print_message "%~1" "%GREEN%"
goto :eof

:print_warning
call :print_message "%~1" "%YELLOW%"
goto :eof

:print_error
call :print_message "%~1" "%RED%"
goto :eof

REM 检查Docker和Docker Compose
:check_dependencies
call :print_info "检查依赖..."

docker --version >nul 2>&1
if errorlevel 1 (
    call :print_error "Docker 未安装，请先安装 Docker Desktop"
    exit /b 1
)

docker-compose --version >nul 2>&1
if errorlevel 1 (
    docker compose version >nul 2>&1
    if errorlevel 1 (
        call :print_error "Docker Compose 未安装，请先安装 Docker Compose"
        exit /b 1
    )
)

call :print_success "依赖检查完成"
goto :eof

REM 创建必要的目录和文件
:setup_environment
call :print_info "设置环境..."

REM 创建日志目录
if not exist "logs" mkdir logs

REM 检查环境变量文件
if not exist "backend\.env" (
    call :print_warning "backend\.env 文件不存在，从示例文件复制..."
    copy "backend\.env.example" "backend\.env" >nul
)

call :print_success "环境设置完成"
goto :eof

REM 开发环境启动
:start_dev
call :print_info "启动开发环境..."

REM 启动数据库
docker-compose up -d postgres redis

call :print_success "开发环境数据库已启动"
call :print_info "请在另一个终端运行以下命令启动服务:"
call :print_info "后端: cd backend && go run cmd/server/main.go"
call :print_info "前端: npm run dev"
goto :eof

REM 生产环境启动
:start_prod
call :print_info "启动生产环境..."

REM 复制生产环境配置
if exist "backend\.env.production" (
    copy "backend\.env.production" "backend\.env" >nul
    call :print_info "已使用生产环境配置"
)

REM 构建并启动所有服务
docker-compose up -d --build

call :print_success "生产环境已启动"
call :print_info "应用访问地址: http://localhost"
call :print_info "API访问地址: http://localhost:8080"
call :print_info "数据库访问地址: localhost:5432"
goto :eof

REM 停止服务
:stop_services
call :print_info "停止所有服务..."
docker-compose down
call :print_success "所有服务已停止"
goto :eof

REM 查看日志
:show_logs
call :print_info "显示服务日志..."
docker-compose logs -f
goto :eof

REM 清理环境
:clean_environment
call :print_warning "这将删除所有容器、镜像和数据卷，确定要继续吗？(y/N)"
set /p response="请输入选择: "
if /i "!response!"=="y" (
    call :print_info "清理环境..."
    docker-compose down -v --rmi all
    docker system prune -f
    call :print_success "环境清理完成"
) else (
    call :print_info "取消清理操作"
)
goto :eof

REM 显示帮助信息
:show_help
echo PromeConfig 部署脚本 (Windows版本)
echo.
echo 使用方法: %~nx0 [命令]
echo.
echo 命令:
echo   dev     启动开发环境（仅数据库）
echo   prod    启动生产环境（完整服务）
echo   stop    停止所有服务
echo   logs    查看服务日志
echo   clean   清理所有容器和数据
echo   help    显示此帮助信息
echo.
echo 示例:
echo   %~nx0 prod    # 启动生产环境
echo   %~nx0 dev     # 启动开发环境
echo   %~nx0 logs    # 查看日志
goto :eof

REM 主函数
set "command=%~1"
if "%command%"=="" set "command=help"

if "%command%"=="dev" (
    call :check_dependencies
    call :setup_environment
    call :start_dev
) else if "%command%"=="prod" (
    call :check_dependencies
    call :setup_environment
    call :start_prod
) else if "%command%"=="stop" (
    call :stop_services
) else if "%command%"=="logs" (
    call :show_logs
) else if "%command%"=="clean" (
    call :clean_environment
) else if "%command%"=="help" (
    call :show_help
) else (
    call :print_error "未知命令: %command%"
    call :show_help
    exit /b 1
)

endlocal