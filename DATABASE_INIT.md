# 🗃️ 数据库初始化指南

## 📋 概述

本项目支持智能数据库初始化功能，可以自动检查数据库状态并执行相应的初始化操作，避免重复初始化和数据丢失。

## 🚀 快速开始

### 🐳 方式一：Docker 完整环境（推荐）

```bash
# 启动完整开发环境（包含自动数据库初始化）
make docker-dev-full
```

这个命令会：
- 🔄 启动所有 Docker 服务（MySQL、Redis、后端、前端、管理工具）
- 🔍 自动检查数据库是否已初始化
- ⚡ 如果未初始化，则自动执行初始化
- ✅ 如果已初始化，则跳过初始化步骤

### 🛠️ 方式二：智能初始化

```bash
# 智能检查并初始化（推荐）
make init-db-smart

# 强制重新初始化（谨慎使用）
make init-db

# 生产环境初始化
make init-db-prod
```

### 🔧 方式三：手动脚本初始化

```bash
# 进入后端目录
cd backend

# 运行初始化脚本（使用默认配置）
go run scripts/init_admin.go

# 或指定配置文件
go run scripts/init_admin.go -config configs/prod.yaml
```

### 🎯 方式四：服务启动时自动初始化

当启动后端服务时，系统会自动检查数据库是否为空，如果为空则自动初始化基础数据。

```bash
# 启动后端服务
cd backend
go run cmd/main.go
```

## 🌐 服务端口

启动后可以通过以下地址访问各个服务：

- **🔗 后端 API**: http://localhost:8081
- **🎨 前端应用**: http://localhost:8082  
- **📚 API 文档**: http://localhost:8081/swagger/index.html
- **🗃️ 数据库管理**: http://localhost:8083 (Adminer)
- **🔴 Redis 管理**: http://localhost:8084 (Redis Commander)

## 🔑 默认账号

初始化完成后，系统会创建以下默认账号：

### 👑 管理员账号
- **用户名**: `admin`
- **密码**: `123456`
- **权限**: 系统管理员（拥有所有 API 权限）

### 👤 普通用户账号
- **用户名**: `member`
- **密码**: `123456`
- **权限**: 普通用户

> ⚠️ **安全提醒**: 请在生产环境中立即修改默认密码！

## ⚙️ 智能初始化工作原理

智能初始化功能通过以下步骤工作：

1. 🔌 连接到 MySQL 数据库
2. 🔍 检查 `users` 表中是否存在 `admin` 用户
3. ⚡ 如果不存在，执行完整的数据库初始化
4. ✅ 如果已存在，跳过初始化步骤

这确保了：
- 🚀 首次启动时自动初始化数据库
- 🔄 后续启动时不会重复初始化
- 🛡️ 避免数据丢失和重复操作

## 📦 初始化内容

数据库初始化会创建以下内容：

### 1. 默认部门
- **名称**: 默认部门
- **编码**: default
- **级别**: 1
- **状态**: 启用

### 2. 默认角色
- **管理员角色**
  - 名称: 管理员
  - 编码: admin
  - 描述: 系统管理员
  
- **普通用户角色**
  - 名称: 普通用户
  - 编码: user
  - 描述: 普通用户

### 3. 默认用户
- **管理员用户**
  - 用户名: `admin`
  - 密码: `123456`
  - 邮箱: admin@example.com
  - 角色: 管理员
  
- **普通用户**
  - 用户名: `member`
  - 密码: `123456`
  - 邮箱: member@example.com
  - 角色: 普通用户

### 4. 默认业务线
- **名称**: 默认业务线
- **编码**: default
- **负责人**: 管理员
- **邮箱**: admin@example.com

### 5. 默认API分类
- **名称**: 系统管理
- **编码**: system

### 6. 示例API
- **名称**: 用户登录
- **路径**: /api/v1/auth/login
- **方法**: POST
- **描述**: 用户登录接口

## 密码加密方式

系统使用SHA256哈希算法对密码进行加密：

```go
func GenerateHash(password string) string {
    hash := sha256.Sum256([]byte(password))
    return base64.StdEncoding.EncodeToString(hash[:])
}
```

默认密码`123456`的哈希值为：`jGl25bVBBBW96Qi9Te4V37Fnqchz/Eu4qB2JKbpbGKw=`

## 数据库配置

### 开发环境配置 (configs/dev.yaml)
```yaml
mysql:
  host: localhost
  port: 3306
  username: root
  password: lidi10
  database: api_auth_dev
  charset: utf8mb4
```

### 生产环境配置 (configs/prod.yaml)
```yaml
mysql:
  host: localhost
  port: 3306
  username: root
  password: your_production_password
  database: api_auth
  charset: utf8mb4
```

## 注意事项

1. **安全性**: 生产环境中请务必修改默认密码
2. **数据检查**: 初始化脚本会检查数据库中是否已有用户数据，如果有则跳过初始化
3. **事务处理**: 所有初始化操作都在数据库事务中执行，确保数据一致性
4. **配置文件**: 确保数据库配置正确，特别是数据库名称和连接信息

## 🆘 故障排除

### 🔌 端口冲突
如果遇到端口冲突，可以通过环境变量修改端口：

```bash
export BACKEND_PORT=8091
export FRONTEND_PORT=8092
make docker-dev-full
```

### 🗃️ 数据库连接失败
确保 MySQL 服务已经启动并且健康检查通过：

```bash
# 检查服务状态
docker-compose -f docker-compose.dev.yml ps

# 查看数据库日志
docker-compose -f docker-compose.dev.yml logs mysql-dev
```

**常见解决方案**：
- 检查数据库服务是否启动
- 验证配置文件中的连接信息
- 确认数据库用户权限

### 🔄 重置数据库
如果需要完全重置数据库：

```bash
# 停止所有服务
make docker-dev-down

# 删除数据卷
docker volume rm mcprapi_mysql-dev-data

# 重新启动
make docker-dev-full
```

### ⚡ 初始化失败
- 检查数据库用户是否有创建表的权限
- 查看错误日志获取详细信息：
  ```bash
  docker-compose -f docker-compose.dev.yml logs backend-dev
  ```
- 确认数据库字符集设置

### 🔧 手动重新初始化
如果需要手动重新初始化，请先清空数据库：
```sql
DROP DATABASE api_auth_dev;
CREATE DATABASE api_auth_dev CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

然后重新运行初始化脚本。

### 📊 检查初始化状态
```bash
# 检查数据库中的表
docker exec -it mcprapi_mysql-dev_1 mysql -u root -p -e "USE api_auth_dev; SHOW TABLES;"

# 检查用户数据
docker exec -it mcprapi_mysql-dev_1 mysql -u root -p -e "USE api_auth_dev; SELECT username FROM users;"
```