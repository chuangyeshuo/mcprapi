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

# 初始化开发环境数据库
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
dev-setup: install init-db
	@echo "开发环境设置完成！"
	@echo "现在可以运行 'make run-backend' 和 'make run-frontend' 启动服务"

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

# ==================== 环境配置命令 ====================

# 创建环境变量文件
.PHONY: env-setup
env-setup:
	cp .env.example .env
	@echo "请编辑 .env 文件配置您的环境变量"

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
	@echo "  make env-setup         - 创建环境变量文件"
	@echo "  make env-check         - 检查环境依赖"
	@echo "  make init              - 初始化项目"
	@echo "  make help              - 显示帮助信息"