# 数据库初始化说明

## 概述

本项目现在支持智能数据库初始化功能，可以自动检查数据库是否已经初始化，避免重复初始化。

## 使用方法

### 1. 完整 Docker 开发环境启动（推荐）

```bash
make docker-dev-full
```

这个命令会：
- 启动所有 Docker 服务（MySQL、Redis、后端、前端、管理工具）
- 自动检查数据库是否已初始化
- 如果未初始化，则自动执行初始化
- 如果已初始化，则跳过初始化步骤

### 2. 手动数据库初始化

#### 智能初始化（推荐）
```bash
make init-db-smart
```
- 自动检查数据库状态
- 只在需要时执行初始化

#### 强制重新初始化
```bash
make init-db
```
- 强制执行数据库初始化
- 会重新创建所有表和数据

### 3. 生产环境初始化
```bash
make init-db-prod
```

## 服务端口

启动后可以通过以下地址访问各个服务：

- **后端 API**: http://localhost:8081
- **前端应用**: http://localhost:8082  
- **API 文档**: http://localhost:8081/swagger/index.html
- **数据库管理**: http://localhost:8083 (Adminer)
- **Redis 管理**: http://localhost:8084 (Redis Commander)

## 默认账号

初始化完成后，系统会创建以下默认账号：

### 管理员账号
- 用户名: `admin`
- 密码: `123456`
- 权限: 系统管理员（拥有所有 API 权限）

### 普通用户账号
- 用户名: `member`
- 密码: `123456`
- 权限: 普通用户

⚠️ **安全提醒**: 请在生产环境中立即修改默认密码！

## 工作原理

智能初始化功能通过以下步骤工作：

1. 连接到 MySQL 数据库
2. 检查 `users` 表中是否存在 `admin` 用户
3. 如果不存在，执行完整的数据库初始化
4. 如果已存在，跳过初始化步骤

这确保了：
- 首次启动时自动初始化数据库
- 后续启动时不会重复初始化
- 避免数据丢失和重复操作

## 故障排除

### 端口冲突
如果遇到端口冲突，可以通过环境变量修改端口：

```bash
export BACKEND_PORT=8091
export FRONTEND_PORT=8092
make docker-dev-full
```

### 数据库连接失败
确保 MySQL 服务已经启动并且健康检查通过：

```bash
docker-compose -f docker-compose.dev.yml ps
```

### 重置数据库
如果需要完全重置数据库：

```bash
make docker-dev-down
docker volume rm mcprapi_mysql-dev-data
make docker-dev-full
```