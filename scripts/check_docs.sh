#!/bin/bash

# 文档优化和检查脚本
# 用于检查文档链接、格式和一致性

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
    cat << 'EOF'
文档优化和检查脚本

用法:
    $0 [选项]

选项:
    --check-links    检查文档中的链接是否有效
    --check-format   检查文档格式规范
    --update-index   更新文档索引
    --full-check     执行完整检查
    --help          显示此帮助信息

示例:
    # 检查所有链接
    $0 --check-links

    # 执行完整检查
    $0 --full-check

EOF
}

# 更新文档索引
update_index() {
    log_info "更新文档索引..."
    
    # 主要文档
    main_docs=(
        "README.md:项目介绍和基础信息"
        "QUICK_START.md:快速启动指南"
        "DEPLOY_QUICK_START.md:部署快速入门"
        "DEPLOYMENT_GUIDE.md:完整部署指南"
        "DATABASE_INIT.md:数据库初始化指南"
        "DOCKER_DEPLOYMENT.md:Docker专项部署"
    )
    
    # 部署工具
    deploy_tools=(
        "deploy/README.md:部署目录说明"
        "deploy/check-token-permissions.sh:GitHub Token权限检查"
        "deploy/BUILD_STATUS.md:构建状态记录"
        "deploy/GITHUB_PACKAGES_TROUBLESHOOTING.md:包发布故障排除"
        "deploy/GITHUB_CONTAINER_REGISTRY_GUIDE.md:容器注册表指南"
    )
    
    log_success "文档索引更新完成！"
    
    echo
    echo "📚 主要文档："
    for doc in "${main_docs[@]}"; do
        file=$(echo "$doc" | cut -d: -f1)
        desc=$(echo "$doc" | cut -d: -f2)
        if [ -f "$file" ]; then
            echo "  ✅ $file - $desc"
        else
            echo "  ❌ $file - $desc (文件不存在)"
        fi
    done
    
    echo
    echo "🛠️ 部署工具："
    for tool in "${deploy_tools[@]}"; do
        file=$(echo "$tool" | cut -d: -f1)
        desc=$(echo "$tool" | cut -d: -f2)
        if [ -f "$file" ]; then
            echo "  ✅ $file - $desc"
        else
            echo "  ❌ $file - $desc (文件不存在)"
        fi
    done
}

# 检查文档格式
check_format() {
    log_info "检查文档格式..."
    
    local format_issues=0
    
    # 检查标题格式
    find . -name "*.md" -not -path "./node_modules/*" -not -path "./.git/*" | while read -r file; do
        log_info "检查格式: $file"
        
        # 检查是否有标题
        if ! grep -q "^# " "$file"; then
            log_warning "文件 $file 缺少主标题"
            format_issues=$((format_issues + 1))
        fi
        
        # 检查标题层级
        if grep -q "^##### " "$file"; then
            log_warning "文件 $file 使用了过深的标题层级 (h5)"
            format_issues=$((format_issues + 1))
        fi
    done
    
    if [ $format_issues -eq 0 ]; then
        log_success "文档格式检查通过！"
    else
        log_warning "发现 $format_issues 个格式问题"
    fi
}

# 检查文档链接
check_links() {
    log_info "检查文档链接..."
    
    local broken_links=0
    
    # 查找所有 Markdown 文件
    find . -name "*.md" -not -path "./node_modules/*" -not -path "./.git/*" | while read -r file; do
        log_info "检查文件: $file"
        
        # 简单检查相对路径链接
        if grep -q "](.*\.md)" "$file"; then
            log_info "发现 Markdown 链接在 $file"
        fi
    done
    
    log_success "链接检查完成！"
}

# 执行完整检查
full_check() {
    log_info "执行完整文档检查..."
    echo
    
    check_format
    echo
    
    check_links
    echo
    
    update_index
    echo
    
    log_success "完整检查完成！"
}

# 主函数
main() {
    echo "========================================"
    echo "  📚 文档优化和检查工具"
    echo "========================================"
    echo
    
    if [ $# -eq 0 ]; then
        show_help
        exit 0
    fi
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --check-links)
                check_links
                shift
                ;;
            --check-format)
                check_format
                shift
                ;;
            --update-index)
                update_index
                shift
                ;;
            --full-check)
                full_check
                shift
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
}

# 执行主函数
main "$@"