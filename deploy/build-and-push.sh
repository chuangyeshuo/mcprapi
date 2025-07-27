#!/bin/bash

# MCPRAPI GitHub Container Registry 构建和推送脚本
# 使用方法: ./build-and-push.sh [version] [--no-cache]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目配置
REGISTRY="ghcr.io"
NAMESPACE="chuangyeshuo"
PROJECT="mcprapi"
DEFAULT_VERSION="latest"

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

# 显示帮助信息
show_help() {
    echo "MCPRAPI Docker 镜像构建和推送脚本"
    echo
    echo "使用方法:"
    echo "  $0 [version] [options]"
    echo
    echo "参数:"
    echo "  version     镜像版本标签 (默认: latest)"
    echo
    echo "选项:"
    echo "  --no-cache  不使用Docker缓存构建"
    echo "  --help      显示此帮助信息"
    echo
    echo "示例:"
    echo "  $0                    # 构建 latest 版本"
    echo "  $0 1.2.0             # 构建 1.2.0 版本"
    echo "  $0 1.2.0 --no-cache  # 不使用缓存构建 1.2.0 版本"
}

# 检查Docker是否安装和运行
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

# 检查是否已登录GHCR
check_login() {
    log_info "检查GitHub Container Registry登录状态..."
    
    if ! docker system info | grep -q "ghcr.io"; then
        log_warning "未登录到GitHub Container Registry"
        log_info "请先登录: echo \$GITHUB_TOKEN | docker login ghcr.io -u $NAMESPACE --password-stdin"
        
        if [ -n "$GITHUB_TOKEN" ]; then
            log_info "检测到GITHUB_TOKEN环境变量，尝试自动登录..."
            echo $GITHUB_TOKEN | docker login ghcr.io -u $NAMESPACE --password-stdin
            if [ $? -eq 0 ]; then
                log_success "自动登录成功"
            else
                log_error "自动登录失败，请检查GITHUB_TOKEN"
                exit 1
            fi
        else
            log_error "请设置GITHUB_TOKEN环境变量或手动登录"
            exit 1
        fi
    else
        log_success "已登录到GitHub Container Registry"
    fi
}

# 构建镜像
build_image() {
    local service=$1
    local version=$2
    local no_cache=$3
    
    local image_name="$REGISTRY/$NAMESPACE/$PROJECT-$service:$version"
    local cache_option=""
    
    if [ "$no_cache" = "true" ]; then
        cache_option="--no-cache"
        log_info "构建 $service 镜像 (不使用缓存): $image_name"
    else
        log_info "构建 $service 镜像: $image_name"
    fi
    
    # 检查Dockerfile是否存在
    if [ ! -f "../$service/Dockerfile" ]; then
        log_error "未找到 ../$service/Dockerfile"
        return 1
    fi
    
    # 构建镜像
    docker build $cache_option -t $image_name ../$service
    
    if [ $? -eq 0 ]; then
        log_success "$service 镜像构建成功"
        
        # 如果不是latest版本，同时打latest标签
        if [ "$version" != "latest" ]; then
            local latest_image="$REGISTRY/$NAMESPACE/$PROJECT-$service:latest"
            docker tag $image_name $latest_image
            log_info "已为 $service 添加 latest 标签"
        fi
        
        return 0
    else
        log_error "$service 镜像构建失败"
        return 1
    fi
}

# 推送镜像
push_image() {
    local service=$1
    local version=$2
    
    local image_name="$REGISTRY/$NAMESPACE/$PROJECT-$service:$version"
    
    log_info "推送 $service 镜像: $image_name"
    docker push $image_name
    
    if [ $? -eq 0 ]; then
        log_success "$service 镜像推送成功"
        
        # 如果不是latest版本，也推送latest标签
        if [ "$version" != "latest" ]; then
            local latest_image="$REGISTRY/$NAMESPACE/$PROJECT-$service:latest"
            log_info "推送 $service latest 镜像: $latest_image"
            docker push $latest_image
            if [ $? -eq 0 ]; then
                log_success "$service latest 镜像推送成功"
            else
                log_error "$service latest 镜像推送失败"
                return 1
            fi
        fi
        
        return 0
    else
        log_error "$service 镜像推送失败"
        return 1
    fi
}

# 显示镜像信息
show_image_info() {
    local version=$1
    
    echo
    log_success "构建完成！镜像信息："
    echo
    
    # 获取镜像摘要
    local backend_digest=$(docker inspect --format='{{index .RepoDigests 0}}' "$REGISTRY/$NAMESPACE/$PROJECT-backend:$version" 2>/dev/null | cut -d'@' -f2)
    local frontend_digest=$(docker inspect --format='{{index .RepoDigests 0}}' "$REGISTRY/$NAMESPACE/$PROJECT-frontend:$version" 2>/dev/null | cut -d'@' -f2)
    
    echo "后端镜像:"
    echo "  - 地址: $REGISTRY/$NAMESPACE/$PROJECT-backend:$version"
    if [ -n "$backend_digest" ]; then
        echo "  - 摘要: $backend_digest"
    fi
    if [ "$version" != "latest" ]; then
        echo "  - 别名: $REGISTRY/$NAMESPACE/$PROJECT-backend:latest"
    fi
    echo
    echo "前端镜像:"
    echo "  - 地址: $REGISTRY/$NAMESPACE/$PROJECT-frontend:$version"
    if [ -n "$frontend_digest" ]; then
        echo "  - 摘要: $frontend_digest"
    fi
    if [ "$version" != "latest" ]; then
        echo "  - 别名: $REGISTRY/$NAMESPACE/$PROJECT-frontend:latest"
    fi
    echo
    echo "使用方法:"
    echo "  # 拉取镜像"
    echo "  docker pull $REGISTRY/$NAMESPACE/$PROJECT-backend:$version"
    echo "  docker pull $REGISTRY/$NAMESPACE/$PROJECT-frontend:$version"
    echo
    echo "  # 使用部署脚本"
    echo "  ./deploy-ghcr.sh"
    echo
    echo "  # 或使用docker-compose"
    echo "  BACKEND_VERSION=$version FRONTEND_VERSION=$version docker-compose -f docker-compose.ghcr.yml up -d"
    echo
    echo "✅ 镜像已成功推送到 GitHub Container Registry"
    echo
}

# 主函数
main() {
    local version=$DEFAULT_VERSION
    local no_cache=false
    
    # 解析参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            --help)
                show_help
                exit 0
                ;;
            --no-cache)
                no_cache=true
                shift
                ;;
            -*)
                log_error "未知选项: $1"
                show_help
                exit 1
                ;;
            *)
                version=$1
                shift
                ;;
        esac
    done
    
    echo "========================================"
    echo "  MCPRAPI Docker 镜像构建和推送"
    echo "========================================"
    echo "版本: $version"
    echo "项目: $REGISTRY/$NAMESPACE/$PROJECT"
    echo "不使用缓存: $no_cache"
    echo "========================================"
    echo
    
    # 检查环境
    check_docker
    check_login
    
    # 构建后端镜像
    if ! build_image "backend" "$version" "$no_cache"; then
        exit 1
    fi
    
    # 构建前端镜像
    if ! build_image "frontend" "$version" "$no_cache"; then
        exit 1
    fi
    
    # 推送后端镜像
    if ! push_image "backend" "$version"; then
        exit 1
    fi
    
    # 推送前端镜像
    if ! push_image "frontend" "$version"; then
        exit 1
    fi
    
    # 显示完成信息
    show_image_info "$version"
}

# 执行主函数
main "$@"