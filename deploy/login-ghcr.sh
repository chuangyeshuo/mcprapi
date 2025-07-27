#!/bin/bash

# GitHub Container Registry 登录帮助脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

echo "========================================"
echo "  GitHub Container Registry 登录助手"
echo "========================================"
echo

log_info "此脚本将帮助您登录到GitHub Container Registry (GHCR)"
echo

# 检查是否已经登录
log_info "检查当前登录状态..."
if docker pull ghcr.io/chuangyeshuo/mcprapi-backend:latest &> /dev/null; then
    log_success "您已经登录到GHCR，无需重新登录"
    exit 0
fi

log_warning "您尚未登录到GHCR"
echo

# 提供登录选项
echo "请选择登录方式："
echo "1. 使用环境变量 GITHUB_TOKEN"
echo "2. 手动输入 GitHub Personal Access Token"
echo "3. 查看如何创建 GitHub Token"
echo

read -p "请选择 (1-3): " choice

case $choice in
    1)
        if [ -z "$GITHUB_TOKEN" ]; then
            log_error "环境变量 GITHUB_TOKEN 未设置"
            log_info "请先设置环境变量: export GITHUB_TOKEN=your_token_here"
            exit 1
        fi
        
        log_info "使用环境变量中的 GITHUB_TOKEN 登录..."
        echo $GITHUB_TOKEN | docker login ghcr.io -u chuangyeshuo --password-stdin
        
        if [ $? -eq 0 ]; then
            log_success "登录成功！"
        else
            log_error "登录失败，请检查token权限"
            exit 1
        fi
        ;;
        
    2)
        log_info "请输入您的GitHub Personal Access Token:"
        log_warning "注意：输入时不会显示字符（这是正常的安全行为）"
        read -s token
        echo
        
        if [ -z "$token" ]; then
            log_error "Token不能为空"
            exit 1
        fi
        
        log_info "正在登录..."
        echo $token | docker login ghcr.io -u chuangyeshuo --password-stdin
        
        if [ $? -eq 0 ]; then
            log_success "登录成功！"
            log_info "建议将token保存到环境变量: export GITHUB_TOKEN=$token"
        else
            log_error "登录失败，请检查token权限"
            exit 1
        fi
        ;;
        
    3)
        echo
        log_info "如何创建GitHub Personal Access Token："
        echo
        echo "1. 访问 GitHub Settings: https://github.com/settings/tokens"
        echo "2. 点击 'Generate new token' -> 'Generate new token (classic)'"
        echo "3. 设置token名称，例如: 'MCPRAPI GHCR Access'"
        echo "4. 选择过期时间（建议选择较长时间或无过期）"
        echo "5. 勾选以下权限："
        echo "   ✓ read:packages   (读取包权限 - 拉取镜像必需)"
        echo "   ✓ write:packages  (写入包权限 - 推送镜像必需)"
        echo "   ✓ repo           (仓库权限 - 如果是私有仓库)"
        echo "6. 点击 'Generate token'"
        echo "7. 复制生成的token（只会显示一次）"
        echo
        log_warning "重要提示："
        echo "- 如果您只需要拉取镜像，read:packages 权限就足够了"
        echo "- 如果您需要推送镜像（构建新版本），必须包含 write:packages 权限"
        echo "- 请妥善保存您的token，GitHub不会再次显示"
        echo
        log_info "创建token后，可以重新运行此脚本进行登录"
        ;;
        
    *)
        log_error "无效选择"
        exit 1
        ;;
esac

echo
log_success "现在您可以运行部署脚本了："
echo "  ./deploy-ghcr.sh"
echo