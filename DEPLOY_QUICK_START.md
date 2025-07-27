# 🚀 部署快速入门

## 📦 使用 GitHub Container Registry (推荐)

最快的部署方式，使用预构建的Docker镜像：

### ✅ 镜像状态

项目镜像已成功构建并推送到 GitHub Container Registry：

- **后端镜像**: `ghcr.io/chuangyeshuo/mcprapi-backend:latest` ✅
- **前端镜像**: `ghcr.io/chuangyeshuo/mcprapi-frontend:latest` ✅

### 🚀 快速部署

```bash
# 1. 克隆项目
git clone https://github.com/chuangyeshuo/mcprapi.git
cd mcprapi

# 2. 进入部署目录
cd deploy

# 3. 登录到GitHub Container Registry（使用登录助手）
./login-ghcr.sh

# 4. 配置环境变量（可选，使用默认配置也可以）
cp .env.production .env
# 编辑 .env 文件，至少修改以下配置：
# - MYSQL_ROOT_PASSWORD=your_secure_password
# - JWT_SECRET=your_jwt_secret_key
# - ENCRYPTION_KEY=your_32_char_encryption_key

# 5. 启动服务
./deploy-ghcr.sh

# 6. 检查服务状态
docker-compose -f docker-compose.ghcr.yml ps
```

## 🔑 GitHub Token 获取

如果您没有GitHub Personal Access Token：

1. 运行登录助手：`./login-ghcr.sh`
2. 选择选项 3 查看详细的Token创建指南
3. 或直接访问：https://github.com/settings/tokens

**需要的权限**：
- ✅ `read:packages` - 读取包权限（拉取镜像必需）
- ✅ `write:packages` - 写入包权限（推送镜像必需，仅开发者需要）
- ✅ `repo` - 仓库权限（如果是私有仓库）

**权限检查**：
```bash
# 检查您的Token权限是否正确
./check-token-permissions.sh
```

## 🌐 访问应用

- **前端界面**: http://localhost:8082
- **后端API**: http://localhost:8081
- **API文档**: http://localhost:8081/swagger/index.html

## 🔑 默认登录

```
用户名: admin
密码: admin
```

> ⚠️ **安全提醒**: 首次登录后请立即修改默认密码！

## 📚 详细文档

- [完整部署指南](deploy/README.md)
- [GitHub Container Registry 指南](deploy/GITHUB_CONTAINER_REGISTRY_GUIDE.md)
- [环境变量配置说明](deploy/.env.production)

## 🛠️ 开发者工具

```bash
# 构建并推送自定义镜像
cd deploy
./build-and-push.sh

# 快速部署脚本
./deploy-ghcr.sh
```

## ❓ 遇到问题？

1. 检查 [故障排除指南](deploy/README.md#故障排除)
2. 查看服务日志: `docker-compose -f docker-compose.ghcr.yml logs`
3. 提交 [Issue](https://github.com/chuangyeshuo/mcprapi/issues)