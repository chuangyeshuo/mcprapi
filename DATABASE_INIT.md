# 数据库初始化指南

## 概述

本项目包含完整的数据库初始化和迁移代码，可以快速设置admin用户和基础数据。

## 数据库初始化方式

### 方式一：自动初始化（推荐）

当启动后端服务时，系统会自动检查数据库是否为空，如果为空则自动初始化基础数据。

```bash
# 启动后端服务
cd backend
go run cmd/main.go
```

### 方式二：使用初始化脚本

如果需要单独初始化数据库，可以使用专门的初始化脚本：

```bash
# 进入后端目录
cd backend

# 运行初始化脚本（使用默认配置）
go run scripts/init_admin.go

# 或指定配置文件
go run scripts/init_admin.go -config configs/prod.yaml
```

## 初始化内容

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

## 故障排除

### 1. 数据库连接失败
- 检查数据库服务是否启动
- 验证配置文件中的连接信息
- 确认数据库用户权限

### 2. 初始化失败
- 检查数据库用户是否有创建表的权限
- 查看错误日志获取详细信息
- 确认数据库字符集设置

### 3. 重新初始化
如果需要重新初始化，请先清空数据库：
```sql
DROP DATABASE api_auth_dev;
CREATE DATABASE api_auth_dev CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

然后重新运行初始化脚本。