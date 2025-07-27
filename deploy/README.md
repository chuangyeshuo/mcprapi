# Deploy 目录说明

本目录包含 MCPRAPI 项目的所有部署相关文件和配置。

## 📁 目录结构

```
deploy/
├── .env.production          # 生产环境配置文件
├── docker-compose.ghcr.yml  # GitHub Container Registry 部署配置
├── deploy-ghcr.sh           # 快速部署脚本
├── login-ghcr.sh            # GitHub Container Registry 登录助手
├── check-token-permissions.sh # GitHub Token 权限检查工具
├── build-and-push.sh        # 镜像构建和推送脚本
├── BUILD_STATUS.md          # 构建状态和镜像信息
├── GITHUB_CONTAINER_REGISTRY_GUIDE.md  # GHCR 详细指南
└── README.md               # 本说明文件
```

## 🔐 安全配置

### 1. 环境变量配置

**`.env.production`** - 生产环境配置文件
- 包含所有生产环境所需的密码和密钥
- **重要**: 请修改所有默认密码和密钥
- 不要将此文件提交到版本控制系统

### 2. 必须修改的配置项

```bash
# GitHub认证
GITHUB_TOKEN=your_github_personal_access_token_here

# 数据库密码
MYSQL_ROOT_PASSWORD=your_very_secure_mysql_root_password_here
MYSQL_PASSWORD=your_very_secure_mysql_user_password_here

# Redis密码
REDIS_PASSWORD=your_very_secure_redis_password_here

# JWT密钥 (至少32个字符)
JWT_SECRET=your_very_secure_jwt_secret_key_at_least_32_characters_long

# 加密密钥
ENCRYPTION_KEY=your_very_secure_encryption_key_32_chars

# 域名配置
CORS_ORIGINS=https://yourdomain.com,https://api.yourdomain.com
```

## 🚀 快速部署

### 方式一：使用登录助手（推荐）

```bash
# 1. 运行登录助手
./login-ghcr.sh

# 2. 按照提示完成登录

# 3. 运行部署脚本
./deploy-ghcr.sh
```

### 方式二：手动登录

```bash
# 1. 设置GitHub Token环境变量
export GITHUB_TOKEN=your_github_token_here

# 2. 登录到GitHub Container Registry
echo $GITHUB_TOKEN | docker login ghcr.io -u chuangyeshuo --password-stdin

# 3. 运行部署脚本
./deploy-ghcr.sh
```

### 方式三：配置环境变量文件

```bash
# 1. 复制环境变量模板
cp .env.production .env

# 2. 编辑 .env 文件，设置 GITHUB_TOKEN
# GITHUB_TOKEN=your_github_token_here

# 3. 加载环境变量
source .env

# 4. 运行部署脚本
./deploy-ghcr.sh
```

## 🛠️ 开发者工具

### 可用镜像

项目已成功构建并推送到 GitHub Container Registry：

**后端镜像**:
- `ghcr.io/chuangyeshuo/mcprapi-backend:latest`
- 镜像摘要: `sha256:fdc3ce25d58fe4ce1da2db004d4b4c38b3c188d35de8e58087727d1cb11727f0`

**前端镜像**:
- `ghcr.io/chuangyeshuo/mcprapi-frontend:latest`
- 镜像摘要: `sha256:e1024a94bb9f0491b47324c4561783c67773119f4a1205585fdc80205fbbd4ec`

### 镜像使用

```bash
# 拉取镜像
docker pull ghcr.io/chuangyeshuo/mcprapi-backend:latest
docker pull ghcr.io/chuangyeshuo/mcprapi-frontend:latest

# 使用 docker-compose 部署
BACKEND_VERSION=latest FRONTEND_VERSION=latest docker-compose -f docker-compose.ghcr.yml up -d

# 或直接使用部署脚本
./deploy-ghcr.sh
```

### 构建和推送镜像

```bash
# 构建最新版本
./build-and-push.sh

# 构建特定版本
./build-and-push.sh 1.2.0

# 不使用缓存构建
./build-and-push.sh 1.2.0 --no-cache
```

### 部署管理

```bash
# 部署最新版本
./deploy-ghcr.sh

# 部署特定版本
./deploy-ghcr.sh 1.2.0

# 查看部署状态
docker-compose -f docker-compose.ghcr.yml ps

# 查看日志
docker-compose -f docker-compose.ghcr.yml logs -f
```

## 🔒 安全最佳实践

### 1. 文件权限

```bash
# 设置配置文件权限
chmod 600 .env.production
chmod 600 .env

# 设置脚本执行权限
chmod +x deploy-ghcr.sh
chmod +x build-and-push.sh
```

### 2. 密码安全

- 使用强密码（至少16个字符，包含大小写字母、数字和特殊字符）
- JWT_SECRET 至少32个字符
- 定期轮换密码和密钥
- 不要在日志中记录敏感信息

### 3. 网络安全

- 配置防火墙规则
- 使用HTTPS
- 设置适当的CORS策略
- 启用访问控制

### 4. 监控和备份

- 定期备份数据库
- 监控系统资源使用情况
- 设置日志轮转
- 配置告警机制

## 📋 部署检查清单

部署前请确认以下项目：

- [ ] 已修改所有默认密码
- [ ] JWT_SECRET 至少32个字符
- [ ] 配置了正确的域名和CORS
- [ ] 设置了适当的文件权限
- [ ] 配置了HTTPS（生产环境）
- [ ] 设置了防火墙规则
- [ ] 配置了备份策略
- [ ] 测试了所有服务功能

## 🆘 故障排除

### 常见问题

1. **GitHub Container Registry 登录问题**
   ```bash
   # 问题：未登录到GHCR
   [WARNING] 未登录到GitHub Container Registry
   [ERROR] 请设置GITHUB_TOKEN环境变量或手动登录
   
   # 解决方案1：使用登录助手
   ./login-ghcr.sh
   
   # 解决方案2：手动登录
   export GITHUB_TOKEN=your_token_here
   echo $GITHUB_TOKEN | docker login ghcr.io -u chuangyeshuo --password-stdin
   
   # 解决方案3：检查Token权限
   # 确保Token具有以下权限：
   # - read:packages (拉取镜像必需)
   # - repo (如果是私有仓库)
   ```

2. **推送权限错误**
   ```bash
   # 问题：推送镜像时权限被拒绝
   denied: permission_denied: The token provided does not match expected scopes.
   
   # 原因：GitHub Token缺少写入权限
   # 解决方案：重新创建Token并包含以下权限：
   # - read:packages (读取包权限)
   # - write:packages (写入包权限 - 推送镜像必需)
   # - repo (仓库权限)
   
   # 步骤：
   # 1. 访问 https://github.com/settings/tokens
   # 2. 创建新Token，确保勾选 write:packages
   # 3. 重新登录
   export GITHUB_TOKEN=new_token_with_write_permission
   echo $GITHUB_TOKEN | docker login ghcr.io -u chuangyeshuo --password-stdin
   ```

3. **权限错误**
   ```bash
   chmod 600 .env
   chmod +x *.sh
   ```

4. **镜像拉取失败**
   ```bash
   # 检查网络连接
   docker pull hello-world
   
   # 检查GitHub Token权限
   echo $GITHUB_TOKEN | docker login ghcr.io -u chuangyeshuo --password-stdin
   
   # 手动拉取测试
   docker pull ghcr.io/chuangyeshuo/mcprapi-backend:latest
   ```

5. **服务启动失败**
   ```bash
   # 查看详细日志
   docker-compose -f docker-compose.ghcr.yml logs
   
   # 检查端口占用
   netstat -tulpn | grep :8081
   netstat -tulpn | grep :8082
   
   # 重启服务
   docker-compose -f docker-compose.ghcr.yml restart
   ```

### GitHub Token 创建指南

如果您没有GitHub Personal Access Token，请按以下步骤创建：

1. 访问 [GitHub Settings](https://github.com/settings/tokens)
2. 点击 "Generate new token" → "Generate new token (classic)"
3. 设置Token名称：`MCPRAPI GHCR Access`
4. 选择过期时间（建议选择较长时间）
5. 勾选权限：
   - ✅ `read:packages` - 读取包权限（拉取镜像必需）
   - ✅ `write:packages` - 写入包权限（推送镜像必需）
   - ✅ `repo` - 仓库权限（如果是私有仓库）
6. 点击 "Generate token"
7. **重要**：立即复制Token（只会显示一次）

### 获取帮助

- 检查Token权限: `./check-token-permissions.sh`
- 使用登录助手: `./login-ghcr.sh`
- 查看主项目文档: `../GITHUB_CONTAINER_REGISTRY_GUIDE.md`
- 查看快速开始: `../DEPLOY_QUICK_START.md`
- 提交Issue: https://github.com/chuangyeshuo/mcprapi/issues

---

**⚠️ 重要提醒**: 
- 生产环境部署前请仔细检查所有配置
- 定期更新镜像版本和安全补丁
- 保持配置文件的安全性