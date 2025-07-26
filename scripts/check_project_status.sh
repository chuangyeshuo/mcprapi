#!/bin/bash

# MCP RAPI 项目状态检查脚本
# 用于快速检查项目的当前状态和配置

echo "🔍 MCP RAPI 项目状态检查"
echo "=========================="
echo ""

# 检查基本环境
echo "📋 环境检查"
echo "----------"

# 检查Docker
if command -v docker >/dev/null 2>&1; then
    echo "✅ Docker: $(docker --version | cut -d' ' -f3 | cut -d',' -f1)"
    
    # 检查Docker Compose
    if command -v docker-compose >/dev/null 2>&1; then
        echo "✅ Docker Compose: $(docker-compose --version | cut -d' ' -f3 | cut -d',' -f1)"
    else
        echo "❌ Docker Compose: 未安装"
    fi
else
    echo "❌ Docker: 未安装"
fi

# 检查Go
if command -v go >/dev/null 2>&1; then
    echo "✅ Go: $(go version | cut -d' ' -f3)"
else
    echo "❌ Go: 未安装"
fi

# 检查Node.js
if command -v node >/dev/null 2>&1; then
    echo "✅ Node.js: $(node --version)"
    
    # 检查npm
    if command -v npm >/dev/null 2>&1; then
        echo "✅ npm: $(npm --version)"
    else
        echo "❌ npm: 未安装"
    fi
else
    echo "❌ Node.js: 未安装"
fi

echo ""

# 检查项目文件结构
echo "📁 项目文件检查"
echo "-------------"

# 检查主要目录
for dir in "backend" "frontend" "configs" "scripts" "docs"; do
    if [ -d "$dir" ]; then
        echo "✅ $dir/ 目录存在"
    else
        echo "❌ $dir/ 目录缺失"
    fi
done

# 检查重要文件
important_files=(
    "README.md"
    "QUICK_START.md"
    "DEPLOYMENT_GUIDE.md"
    "DOCS_OVERVIEW.md"
    "docker-compose.yml"
    "docker-compose.dev.yml"
    "Makefile"
    ".env.example"
)

for file in "${important_files[@]}"; do
    if [ -f "$file" ]; then
        echo "✅ $file 存在"
    else
        echo "❌ $file 缺失"
    fi
done

echo ""

# 检查Docker容器状态
echo "🐳 Docker 容器状态"
echo "----------------"

if command -v docker >/dev/null 2>&1; then
    # 检查是否有运行的容器
    running_containers=$(docker ps --format "table {{.Names}}\t{{.Status}}" | grep mcprapi)
    
    if [ -n "$running_containers" ]; then
        echo "运行中的容器:"
        echo "$running_containers"
    else
        echo "❌ 没有运行中的 MCP RAPI 容器"
    fi
    
    echo ""
    
    # 检查所有相关容器
    all_containers=$(docker ps -a --format "table {{.Names}}\t{{.Status}}" | grep mcprapi)
    
    if [ -n "$all_containers" ]; then
        echo "所有相关容器:"
        echo "$all_containers"
    else
        echo "ℹ️  没有找到 MCP RAPI 相关容器"
    fi
else
    echo "❌ Docker 未安装，无法检查容器状态"
fi

echo ""

# 检查端口占用
echo "🌐 端口状态检查"
echo "-------------"

ports=(8081 8082 3306 6379 8083 8084)
port_names=("后端API" "前端Web" "MySQL" "Redis" "Adminer" "Redis Commander")

for i in "${!ports[@]}"; do
    port=${ports[$i]}
    name=${port_names[$i]}
    
    if command -v lsof >/dev/null 2>&1; then
        if lsof -i :$port >/dev/null 2>&1; then
            echo "✅ 端口 $port ($name) 正在使用"
        else
            echo "⚪ 端口 $port ($name) 空闲"
        fi
    elif command -v netstat >/dev/null 2>&1; then
        if netstat -an | grep ":$port " >/dev/null 2>&1; then
            echo "✅ 端口 $port ($name) 正在使用"
        else
            echo "⚪ 端口 $port ($name) 空闲"
        fi
    else
        echo "❓ 无法检查端口 $port ($name) 状态"
    fi
done

echo ""

# 检查配置文件
echo "⚙️  配置文件检查"
echo "-------------"

config_files=(
    "backend/configs/dev.yaml"
    "backend/configs/prod.yaml"
    "frontend/.env.development"
    "frontend/.env.production"
)

for config in "${config_files[@]}"; do
    if [ -f "$config" ]; then
        echo "✅ $config 存在"
    else
        echo "❌ $config 缺失"
    fi
done

echo ""

# 检查依赖
echo "📦 依赖检查"
echo "----------"

# 检查后端依赖
if [ -f "backend/go.mod" ]; then
    echo "✅ 后端 go.mod 存在"
    if [ -f "backend/go.sum" ]; then
        echo "✅ 后端 go.sum 存在"
    else
        echo "⚠️  后端 go.sum 缺失，可能需要运行 go mod tidy"
    fi
else
    echo "❌ 后端 go.mod 缺失"
fi

# 检查前端依赖
if [ -f "frontend/package.json" ]; then
    echo "✅ 前端 package.json 存在"
    if [ -d "frontend/node_modules" ]; then
        echo "✅ 前端 node_modules 存在"
    else
        echo "⚠️  前端 node_modules 缺失，需要运行 npm install"
    fi
else
    echo "❌ 前端 package.json 缺失"
fi

echo ""

# 提供建议
echo "💡 建议操作"
echo "----------"

if ! command -v docker >/dev/null 2>&1; then
    echo "🔧 请安装 Docker 和 Docker Compose"
fi

if [ ! -f ".env" ]; then
    echo "🔧 运行 'make env-setup' 创建环境配置文件"
fi

if [ ! -d "frontend/node_modules" ]; then
    echo "🔧 运行 'make install-frontend' 安装前端依赖"
fi

# 检查是否有运行的容器
if command -v docker >/dev/null 2>&1; then
    if ! docker ps | grep mcprapi >/dev/null 2>&1; then
        echo "🚀 运行 'make docker-dev-full' 启动完整开发环境"
        echo "🚀 或运行 'make docker-up' 启动生产环境"
    fi
fi

echo ""
echo "📚 更多信息请查看:"
echo "  - README.md - 项目概览"
echo "  - QUICK_START.md - 快速开始指南"
echo "  - DOCS_OVERVIEW.md - 完整文档导航"
echo "  - 运行 'make help' 查看所有可用命令"
echo ""
echo "✨ 检查完成！"