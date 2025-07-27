#!/bin/bash

# GitHub Token 权限检查工具
# 用于验证 GitHub Personal Access Token 的权限设置

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
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
    cat << EOF
GitHub Token 权限检查工具

用法:
    $0 [选项]

选项:
    --token TOKEN    指定要检查的 GitHub Token
    --help          显示此帮助信息

环境变量:
    GITHUB_TOKEN    GitHub Personal Access Token

示例:
    # 使用环境变量
    export GITHUB_TOKEN=your_token_here
    $0

    # 直接指定 token
    $0 --token your_token_here

权限要求:
    - read:packages   (拉取镜像必需)
    - write:packages  (推送镜像必需)
    - repo           (私有仓库访问)

EOF
}

# 检查 GitHub Token 权限
check_token_permissions() {
    local token=$1
    
    if [ -z "$token" ]; then
        log_error "未提供 GitHub Token"
        echo
        echo "请通过以下方式之一提供 Token："
        echo "1. 设置环境变量: export GITHUB_TOKEN=your_token_here"
        echo "2. 使用参数: $0 --token your_token_here"
        echo
        return 1
    fi
    
    log_info "检查 GitHub Token 权限..."
    echo
    
    # 检查 token 基本信息
    log_info "获取 Token 基本信息..."
    local response=$(curl -s -H "Authorization: token $token" \
                          -H "Accept: application/vnd.github.v3+json" \
                          https://api.github.com/user)
    
    if echo "$response" | grep -q "Bad credentials"; then
        log_error "Token 无效或已过期"
        echo
        echo "请检查您的 GitHub Token 是否正确，或创建新的 Token："
        echo "https://github.com/settings/tokens"
        return 1
    fi
    
    local username=$(echo "$response" | grep -o '"login":"[^"]*"' | cut -d'"' -f4)
    if [ -n "$username" ]; then
        log_success "Token 有效，用户: $username"
    else
        log_error "无法获取用户信息"
        return 1
    fi
    
    # 检查 packages 权限
    log_info "检查 packages 权限..."
    
    # 尝试访问 packages API
    local packages_response=$(curl -s -w "%{http_code}" -o /dev/null \
                                   -H "Authorization: token $token" \
                                   -H "Accept: application/vnd.github.v3+json" \
                                   "https://api.github.com/user/packages?package_type=container")
    
    if [ "$packages_response" = "200" ]; then
        log_success "✅ read:packages 权限正常"
    elif [ "$packages_response" = "403" ]; then
        log_error "❌ 缺少 read:packages 权限"
        echo "   需要此权限来拉取镜像"
    else
        log_warning "⚠️  无法确定 read:packages 权限状态 (HTTP: $packages_response)"
    fi
    
    # 检查 repo 权限（通过尝试访问私有仓库信息）
    log_info "检查 repo 权限..."
    local repo_response=$(curl -s -w "%{http_code}" -o /dev/null \
                               -H "Authorization: token $token" \
                               -H "Accept: application/vnd.github.v3+json" \
                               "https://api.github.com/repos/chuangyeshuo/mcprapi")
    
    if [ "$repo_response" = "200" ]; then
        log_success "✅ repo 权限正常"
    elif [ "$repo_response" = "404" ]; then
        log_warning "⚠️  可能缺少 repo 权限或仓库不存在"
        echo "   如果是私有仓库，需要 repo 权限"
    else
        log_warning "⚠️  无法确定 repo 权限状态 (HTTP: $repo_response)"
    fi
    
    echo
    log_info "权限检查完成"
    echo
    
    # 提供建议
    echo "📋 权限建议："
    echo "  ✅ read:packages  - 拉取镜像必需"
    echo "  ✅ write:packages - 推送镜像必需（开发者）"
    echo "  ✅ repo          - 私有仓库访问必需"
    echo
    echo "🔗 创建新 Token: https://github.com/settings/tokens"
    echo "📚 详细指南: ./GITHUB_CONTAINER_REGISTRY_GUIDE.md"
    echo
}

# 主函数
main() {
    local token=""
    
    # 解析参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            --token)
                token="$2"
                shift 2
                ;;
            --help)
                show_help
                exit 0
                ;;
            *)
                log_error "未知参数: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # 如果没有通过参数指定 token，尝试从环境变量获取
    if [ -z "$token" ]; then
        token="$GITHUB_TOKEN"
    fi
    
    echo "========================================"
    echo "  GitHub Token 权限检查工具"
    echo "========================================"
    echo
    
    check_token_permissions "$token"
}

# 执行主函数
main "$@"