#!/bin/bash

# GitHub Container Registry 快速部署脚本
# 使用方法: ./deploy-ghcr.sh [version]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 默认配置
DEFAULT_VERSION="latest"
REGISTRY="ghcr.io"
NAMESPACE="chuangyeshuo"
PROJECT="mcprapi"

# 函数定义
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查Docker是否安装
check_docker() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker未安装，请先安装Docker"
        exit 1
    fi
    
    if ! docker info &> /dev/null; then
        log_error "Docker服务未运行，请启动Docker服务"
        exit 1
    fi
    
    log_success "Docker检查通过"
}

# 检查docker-compose是否安装
check_docker_compose() {
    if ! command -v docker-compose &> /dev/null; then
        log_error "docker-compose未安装，请先安装docker-compose"
        exit 1
    fi
    
    log_success "docker-compose检查通过"
}

# 检查是否已登录到GitHub Container Registry
check_ghcr_login() {
    log_info "检查GitHub Container Registry登录状态..."
    
    # 尝试拉取一个小的测试镜像来检查登录状态
    if docker pull $REGISTRY/$NAMESPACE/$PROJECT-backend:latest &> /dev/null; then
        log_success "已登录到GitHub Container Registry"
        return 0
    else
        log_warning "未登录到GitHub Container Registry"
        log_info "请先登录: echo \$GITHUB_TOKEN | docker login ghcr.io -u $NAMESPACE --password-stdin"
        log_error "请设置GITHUB_TOKEN环境变量或手动登录"
        return 1
    fi
}

# 登录到GitHub Container Registry
login_ghcr() {
    log_info "登录到GitHub Container Registry..."
    
    if [ -z "$GITHUB_TOKEN" ]; then
        log_warning "未设置GITHUB_TOKEN环境变量"
        log_info "请输入GitHub Personal Access Token (需要packages:read权限):"
        read -s token
        export GITHUB_TOKEN=$token
    fi
    
    echo $GITHUB_TOKEN | docker login $REGISTRY -u $NAMESPACE --password-stdin
    
    if [ $? -eq 0 ]; then
        log_success "登录成功"
    else
        log_error "登录失败，请检查token权限"
        exit 1
    fi
}

# 拉取最新镜像
pull_images() {
    local version=${1:-$DEFAULT_VERSION}
    
    log_info "拉取镜像版本: $version"
    
    local backend_image="$REGISTRY/$NAMESPACE/$PROJECT-backend:$version"
    local frontend_image="$REGISTRY/$NAMESPACE/$PROJECT-frontend:$version"
    
    log_info "拉取后端镜像: $backend_image"
    docker pull $backend_image
    
    log_info "拉取前端镜像: $frontend_image"
    docker pull $frontend_image
    
    log_success "镜像拉取完成"
}

# 停止现有服务
stop_services() {
    log_info "停止现有服务..."
    
    if [ -f "docker-compose.ghcr.yml" ]; then
        docker-compose -f docker-compose.ghcr.yml down
    elif [ -f "../docker-compose.yml" ]; then
        docker-compose -f ../docker-compose.yml down
    fi
    
    log_success "服务已停止"
}

# 启动服务
start_services() {
    local version=${1:-$DEFAULT_VERSION}
    
    log_info "启动服务..."
    
    # 设置环境变量
    export BACKEND_VERSION=$version
    export FRONTEND_VERSION=$version
    
    # 使用GitHub Container Registry配置文件
    if [ -f "docker-compose.ghcr.yml" ]; then
        docker-compose -f docker-compose.ghcr.yml up -d
    else
        log_error "未找到docker-compose.ghcr.yml文件"
        exit 1
    fi
    
    log_success "服务启动完成"
}

# 检查服务状态
check_services() {
    log_info "检查服务状态..."
    
    sleep 10  # 等待服务启动
    
    # 检查容器状态
    docker-compose -f docker-compose.ghcr.yml ps
    
    # 检查健康状态
    log_info "等待服务健康检查..."
    
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if docker-compose -f docker-compose.ghcr.yml ps | grep -q "healthy"; then
            log_success "服务健康检查通过"
            break
        fi
        
        log_info "等待服务启动... ($attempt/$max_attempts)"
        sleep 5
        ((attempt++))
    done
    
    if [ $attempt -gt $max_attempts ]; then
        log_warning "服务启动超时，请检查日志"
        docker-compose -f docker-compose.ghcr.yml logs
    fi
}

# 显示访问信息
show_access_info() {
    log_success "部署完成！"
    echo
    echo "访问信息:"
    echo "- 前端地址: http://localhost:8082"
    echo "- 后端API: http://localhost:8081"
    echo "- MySQL: localhost:3306"
    echo "- Redis: localhost:6379"
    echo
    echo "管理命令:"
    echo "- 查看日志: docker-compose -f docker-compose.ghcr.yml logs -f"
    echo "- 停止服务: docker-compose -f docker-compose.ghcr.yml down"
    echo "- 重启服务: docker-compose -f docker-compose.ghcr.yml restart"
    echo
}

# 清理旧镜像
cleanup_images() {
    log_info "清理未使用的镜像..."
    docker image prune -f
    log_success "镜像清理完成"
}

# 主函数
main() {
    local version=${1:-$DEFAULT_VERSION}
    
    echo "========================================"
    echo "  GitHub Container Registry 部署脚本"
    echo "========================================"
    echo
    
    log_info "部署版本: $version"
    echo
    
    # 检查环境
    check_docker
    check_docker_compose
    
    # 如果需要拉取镜像（非本地构建）
    if [ "$version" != "local" ]; then
        # 检查是否已登录到GitHub Container Registry
        if ! check_ghcr_login; then
            log_info "尝试自动登录..."
            login_ghcr
        fi
        
        # 拉取镜像
        pull_images $version
    fi
    
    # 停止现有服务
    stop_services
    
    # 启动新服务
    start_services $version
    
    # 检查服务状态
    check_services
    
    # 显示访问信息
    show_access_info
    
    # 清理旧镜像
    cleanup_images
}

# 帮助信息
show_help() {
    echo "使用方法: $0 [OPTIONS] [VERSION]"
    echo
    echo "选项:"
    echo "  -h, --help     显示帮助信息"
    echo "  -c, --cleanup  仅清理镜像"
    echo "  -s, --stop     仅停止服务"
    echo "  -l, --logs     查看日志"
    echo
    echo "示例:"
    echo "  $0                    # 部署latest版本"
    echo "  $0 1.2.0             # 部署指定版本"
    echo "  $0 --stop            # 停止服务"
    echo "  $0 --cleanup         # 清理镜像"
    echo "  $0 --logs            # 查看日志"
    echo
    echo "环境变量:"
    echo "  GITHUB_TOKEN         GitHub Personal Access Token"
    echo "  GITHUB_USERNAME      GitHub用户名"
    echo
}

# 参数处理
case "${1:-}" in
    -h|--help)
        show_help
        exit 0
        ;;
    -c|--cleanup)
        cleanup_images
        exit 0
        ;;
    -s|--stop)
        stop_services
        exit 0
        ;;
    -l|--logs)
        docker-compose -f docker-compose.ghcr.yml logs -f
        exit 0
        ;;
    *)
        main "$@"
        ;;
esac