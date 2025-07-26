# Makefile for MCP RAPI项目

# 变量定义
BACKEND_DIR = ./backend
FRONTEND_DIR = ./frontend
BIN_NAME = mcprapi
CONFIG_FILE = configs/dev.yaml

# Docker相关变量
COMPOSE_FILE = docker-compose.yml
COMPOSE_DEV_FILE = docker-compose.dev.yml
BACKEND_PORT = 8081
FRONTEND_PORT = 8082

# 默认目标
.PHONY: all
all: build

# 构建后端
.PHONY: build-backend
build-backend:
	cd $(BACKEND_DIR) && go build -o $(BIN_NAME) ./cmd/main.go

# 运行后端
.PHONY: run-backend
run-backend:
	cd $(BACKEND_DIR) && go run ./cmd/main.go --config $(CONFIG_FILE)

# 安装前端依赖
.PHONY: install-frontend
install-frontend:
	cd $(FRONTEND_DIR) && npm install

# 运行前端开发服务器
.PHONY: run-frontend
run-frontend:
	cd $(FRONTEND_DIR) && npm run serve

# 构建前端
.PHONY: build-frontend
build-frontend:
	cd $(FRONTEND_DIR) && npm run build

# 构建整个项目
.PHONY: build
build: build-backend build-frontend

# 运行测试
.PHONY: test
test:
	cd $(BACKEND_DIR) && go test -v ./...

# 清理构建产物
.PHONY: clean
clean:
	rm -f $(BACKEND_DIR)/$(BIN_NAME)
	rm -rf $(FRONTEND_DIR)/dist
	docker system prune -f

# ==================== 数据库初始化命令 ====================

# 检查并初始化开发环境数据库（智能检查）
.PHONY: init-db-smart
init-db-smart:
	@echo "正在检查开发环境数据库状态..."
	@if docker exec mcprapi-mysql-dev mysql -u mcprapi -pdevpassword api_auth_dev -e "SELECT COUNT(*) FROM users WHERE username='admin';" 2>/dev/null | grep -q "1"; then \
		echo "✅ 数据库已经初始化，跳过初始化步骤"; \
	else \
		echo "🔄 数据库未初始化，开始执行初始化..."; \
		cd $(BACKEND_DIR) && go run scripts/init_admin.go -config configs/dev.yaml; \
		echo "✅ 数据库初始化完成！"; \
		echo "管理员账号: admin / 123456"; \
		echo "普通用户账号: member / 123456"; \
	fi

# 初始化开发环境数据库（强制重新初始化）
.PHONY: init-db
init-db:
	@echo "正在初始化开发环境数据库..."
	cd $(BACKEND_DIR) && go run scripts/init_admin.go -config configs/dev.yaml
	@echo "数据库初始化完成！"
	@echo "管理员账号: admin / 123456"
	@echo "普通用户账号: member / 123456"

# 初始化生产环境数据库
.PHONY: init-db-prod
init-db-prod:
	@echo "正在初始化生产环境数据库..."
	cd $(BACKEND_DIR) && go run scripts/init_admin.go -config configs/prod.yaml
	@echo "数据库初始化完成！"
	@echo "⚠️  请立即修改默认密码！"

# 安装后端依赖
.PHONY: install-backend
install-backend:
	cd $(BACKEND_DIR) && go mod tidy

# 安装所有依赖
.PHONY: install
install: install-backend install-frontend
	@echo "依赖安装完成！"

# 开发环境快速启动
.PHONY: dev-setup
dev-setup: install
	@echo "开发环境设置完成！"
	@echo "现在可以运行 'make run-backend' 和 'make run-frontend' 启动服务"

# Docker开发环境完整启动（包含数据库初始化）
.PHONY: docker-dev-full
docker-dev-full: docker-dev-up
	@echo "等待数据库服务启动..."
	@sleep 10
	@$(MAKE) init-db-smart
	@echo "🎉 Docker开发环境启动完成！"
	@echo "后端服务: http://localhost:8081"
	@echo "前端服务: http://localhost:8082"
	@echo "API文档: http://localhost:8081/swagger/index.html"
	@echo "数据库管理: http://localhost:8083"
	@echo "Redis管理: http://localhost:8084"

# ==================== Docker 生产环境命令 ====================

# 启动生产环境Docker容器
.PHONY: docker-up
docker-up:
	docker-compose -f $(COMPOSE_FILE) up -d

# 停止生产环境Docker容器
.PHONY: docker-down
docker-down:
	docker-compose -f $(COMPOSE_FILE) down

# 构建生产环境Docker镜像
.PHONY: docker-build
docker-build:
	docker-compose -f $(COMPOSE_FILE) build --no-cache

# 重启生产环境Docker容器
.PHONY: docker-restart
docker-restart: docker-down docker-up

# 显示生产环境Docker容器日志
.PHONY: docker-logs
docker-logs:
	docker-compose -f $(COMPOSE_FILE) logs -f

# ==================== Docker 开发环境命令 ====================

# 启动开发环境Docker容器
.PHONY: docker-dev-up
docker-dev-up:
	docker-compose -f $(COMPOSE_DEV_FILE) up -d

# 停止开发环境Docker容器
.PHONY: docker-dev-down
docker-dev-down:
	docker-compose -f $(COMPOSE_DEV_FILE) down

# 构建开发环境Docker镜像
.PHONY: docker-dev-build
docker-dev-build:
	docker-compose -f $(COMPOSE_DEV_FILE) build --no-cache

# 重启开发环境Docker容器
.PHONY: docker-dev-restart
docker-dev-restart: docker-dev-down docker-dev-up

# 显示开发环境Docker容器日志
.PHONY: docker-dev-logs
docker-dev-logs:
	docker-compose -f $(COMPOSE_DEV_FILE) logs -f

# ==================== Docker 管理命令 ====================

# 查看Docker容器状态
.PHONY: docker-ps
docker-ps:
	docker ps -a

# 清理Docker资源
.PHONY: docker-clean
docker-clean:
	docker-compose -f $(COMPOSE_FILE) down -v
	docker-compose -f $(COMPOSE_DEV_FILE) down -v
	docker system prune -af
	docker volume prune -f

# 进入后端容器
.PHONY: docker-exec-backend
docker-exec-backend:
	docker exec -it mcprapi-backend /bin/sh

# 进入前端容器
.PHONY: docker-exec-frontend
docker-exec-frontend:
	docker exec -it mcprapi-frontend /bin/sh

# 进入MySQL容器
.PHONY: docker-exec-mysql
docker-exec-mysql:
	docker exec -it mcprapi-mysql mysql -u root -p

# 进入Redis容器
.PHONY: docker-exec-redis
docker-exec-redis:
	docker exec -it mcprapi-redis redis-cli

# ==================== 文档命令 ====================

# 查看文档总览
.PHONY: docs
docs:
	@echo "📚 MCP RAPI 项目文档总览"
	@echo ""
	@echo "🚀 快速开始:"
	@echo "  - README.md           - 项目概览和快速开始"
	@echo "  - QUICK_START.md      - 一键启动指南"
	@echo "  - DOCS_OVERVIEW.md    - 完整文档导航"
	@echo ""
	@echo "🏗️ 部署文档:"
	@echo "  - DEPLOYMENT_GUIDE.md - 完整部署指南"
	@echo "  - DOCKER_DEPLOYMENT.md - Docker部署指南"
	@echo ""
	@echo "🏛️ 架构文档:"
	@echo "  - API多租户授权管理系统架构文档.md - 系统架构"
	@echo "  - DATABASE_INIT.md    - 数据库初始化"
	@echo ""
	@echo "🔧 功能文档:"
	@echo "  - 新增部门流程文档.md - 部门管理流程"
	@echo "  - mcp-example/        - MCP集成示例"
	@echo ""
	@echo "使用 'make docs-serve' 启动文档服务器"

# 启动文档服务器（如果有的话）
.PHONY: docs-serve
docs-serve:
	@if command -v mdbook >/dev/null 2>&1; then \
		echo "启动 mdbook 文档服务器..."; \
		mdbook serve; \
	elif command -v docsify >/dev/null 2>&1; then \
		echo "启动 docsify 文档服务器..."; \
		docsify serve .; \
	elif command -v python3 >/dev/null 2>&1; then \
		echo "启动简单HTTP服务器查看文档..."; \
		python3 -m http.server 3000; \
	else \
		echo "未找到文档服务器工具，请手动查看 Markdown 文件"; \
	fi

# 检查文档链接
.PHONY: docs-check
docs-check:
	@echo "检查文档文件是否存在..."
	@for doc in README.md QUICK_START.md DEPLOYMENT_GUIDE.md DOCS_OVERVIEW.md DOCKER_DEPLOYMENT.md DATABASE_INIT.md; do \
		if [ -f "$$doc" ]; then \
			echo "✅ $$doc"; \
		else \
			echo "❌ $$doc (缺失)"; \
		fi; \
	done

# ==================== 环境配置命令 ====================

# 创建环境变量文件
.PHONY: env-setup
env-setup:
	cp .env.example .env
	@echo "请编辑 .env 文件配置您的环境变量"

# 检查项目状态
.PHONY: status
status:
	@./scripts/check_project_status.sh

# 检查环境
.PHONY: env-check
env-check:
	@echo "检查Docker环境..."
	docker --version
	docker-compose --version
	@echo "检查Go环境..."
	go version
	@echo "检查Node.js环境..."
	node --version
	npm --version

# 初始化项目
.PHONY: init
init: env-check env-setup install-frontend
	cd $(BACKEND_DIR) && go mod download
	@echo "项目初始化完成！"
	@echo "生产环境端口: 后端 $(BACKEND_PORT), 前端 $(FRONTEND_PORT)"
	@echo "使用 'make docker-up' 启动生产环境"
	@echo "使用 'make docker-dev-up' 启动开发环境"

# 帮助信息
.PHONY: help
help:
	@echo "MCP RAPI 项目 Makefile 命令:"
	@echo ""
	@echo "构建命令:"
	@echo "  make build-backend     - 构建后端"
	@echo "  make run-backend       - 运行后端"
	@echo "  make install-frontend  - 安装前端依赖"
	@echo "  make run-frontend      - 运行前端开发服务器"
	@echo "  make build-frontend    - 构建前端"
	@echo "  make build             - 构建整个项目"
	@echo "  make test              - 运行测试"
	@echo "  make clean             - 清理构建产物"
	@echo ""
	@echo "数据库初始化:"
	@echo "  make init-db           - 初始化开发环境数据库 (admin/123456)"
	@echo "  make init-db-prod      - 初始化生产环境数据库"
	@echo "  make install-backend   - 安装后端依赖"
	@echo "  make install           - 安装所有依赖"
	@echo "  make dev-setup         - 开发环境快速设置"
	@echo ""
	@echo "Docker 生产环境:"
	@echo "  make docker-up         - 启动生产环境 (后端:$(BACKEND_PORT), 前端:$(FRONTEND_PORT))"
	@echo "  make docker-down       - 停止生产环境"
	@echo "  make docker-build      - 构建生产环境镜像"
	@echo "  make docker-restart    - 重启生产环境"
	@echo "  make docker-logs       - 查看生产环境日志"
	@echo ""
	@echo "Docker 开发环境:"
	@echo "  make docker-dev-up     - 启动开发环境"
	@echo "  make docker-dev-down   - 停止开发环境"
	@echo "  make docker-dev-build  - 构建开发环境镜像"
	@echo "  make docker-dev-restart - 重启开发环境"
	@echo "  make docker-dev-logs   - 查看开发环境日志"
	@echo ""
	@echo "Docker 管理:"
	@echo "  make docker-ps         - 查看容器状态"
	@echo "  make docker-clean      - 清理Docker资源"
	@echo "  make docker-exec-backend - 进入后端容器"
	@echo "  make docker-exec-frontend - 进入前端容器"
	@echo "  make docker-exec-mysql - 进入MySQL容器"
	@echo "  make docker-exec-redis - 进入Redis容器"
	@echo ""
	@echo "环境配置:"
	@echo "  make status            - 检查项目状态"
	@echo "  make env-setup         - 创建环境变量文件"
	@echo "  make env-check         - 检查环境依赖"
	@echo "  make init              - 初始化项目"
	@echo ""
	@echo "文档命令:"
	@echo "  make docs              - 查看文档总览"
	@echo "  make docs-serve        - 启动文档服务器"
	@echo "  make docs-check        - 检查文档文件"
	@echo ""
	@echo "其他:"
	@echo "  make help              - 显示帮助信息"