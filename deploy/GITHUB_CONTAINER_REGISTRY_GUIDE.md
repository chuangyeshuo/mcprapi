# GitHub Container Registry (GHCR) 使用指南

本指南将帮助您将 mcprapi 项目的Docker镜像发布到GitHub Container Registry，并在生产环境中使用。

## 项目信息

- **GitHub仓库**: https://github.com/chuangyeshuo/mcprapi
- **容器注册表**: ghcr.io/chuangyeshuo
- **后端镜像**: ghcr.io/chuangyeshuo/mcprapi-backend
- **前端镜像**: ghcr.io/chuangyeshuo/mcprapi-frontend

### 当前可用镜像

**后端镜像 (latest)**:
- 镜像地址: `ghcr.io/chuangyeshuo/mcprapi-backend:latest`
- 镜像摘要: `sha256:fdc3ce25d58fe4ce1da2db004d4b4c38b3c188d35de8e58087727d1cb11727f0`
- 状态: ✅ 已构建并推送成功

**前端镜像 (latest)**:
- 镜像地址: `ghcr.io/chuangyeshuo/mcprapi-frontend:latest`
- 镜像摘要: `sha256:e1024a94bb9f0491b47324c4561783c67773119f4a1205585fdc80205fbbd4ec`
- 状态: ✅ 已构建并推送成功

### 快速使用

```bash
# 拉取镜像
docker pull ghcr.io/chuangyeshuo/mcprapi-backend:latest
docker pull ghcr.io/chuangyeshuo/mcprapi-frontend:latest

# 使用部署脚本
cd deploy
./deploy-ghcr.sh
```

## 概述

GitHub Container Registry (ghcr.io) 是GitHub提供的容器镜像托管服务，与您的GitHub仓库紧密集成。

### 相比Docker Hub的优势

- **免费且无限制**: 公共镜像完全免费，私有镜像有慷慨的免费额度
- **与GitHub集成**: 自动继承仓库权限，支持GitHub Actions
- **更好的性能**: 全球CDN分发，下载速度更快
- **安全扫描**: 内置漏洞扫描功能
- **版本管理**: 与Git标签自动关联

## 配置步骤

### 1. 创建Personal Access Token (PAT)

1. 访问 GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. 点击 "Generate new token (classic)"
3. 设置token名称，如 "MCPRAPI GHCR Access"
4. 选择权限：
   - `read:packages` - 读取包权限（拉取镜像必需）
   - `write:packages` - 写入包权限（推送镜像必需）
   - `delete:packages` - 删除包权限（可选）
   - `repo` - 仓库访问权限（用于GitHub Actions）

   **重要提示**: 如果您需要推送镜像到GHCR，必须包含 `write:packages` 权限。仅有 `read:packages` 权限只能拉取镜像，无法推送。
5. 点击 "Generate token" 并保存token

### 2. 登录到GitHub Container Registry

```bash
# 使用Personal Access Token登录
echo $GITHUB_TOKEN | docker login ghcr.io -u chuangyeshuo --password-stdin

# 或者交互式登录
docker login ghcr.io -u chuangyeshuo
```

## 手动构建和推送镜像

### 构建后端镜像

```bash
# 构建后端镜像
docker build -t ghcr.io/chuangyeshuo/mcprapi-backend:latest ./backend

# 推送到GHCR
docker push ghcr.io/chuangyeshuo/mcprapi-backend:latest

# 构建并推送特定版本
docker build -t ghcr.io/chuangyeshuo/mcprapi-backend:1.2.0 ./backend
docker push ghcr.io/chuangyeshuo/mcprapi-backend:1.2.0
```

### 构建前端镜像

```bash
# 构建前端镜像
docker build -t ghcr.io/chuangyeshuo/mcprapi-frontend:latest ./frontend

# 推送到GHCR
docker push ghcr.io/chuangyeshuo/mcprapi-frontend:latest

# 构建并推送特定版本
docker build -t ghcr.io/chuangyeshuo/mcprapi-frontend:1.2.0 ./frontend
docker push ghcr.io/chuangyeshuo/mcprapi-frontend:1.2.0
```

### 批量构建脚本

创建 `build-and-push.sh` 脚本：

```bash
#!/bin/bash
VERSION=${1:-latest}

echo "构建版本: $VERSION"

# 构建后端
echo "构建后端镜像..."
docker build -t ghcr.io/chuangyeshuo/mcprapi-backend:$VERSION ./backend
docker push ghcr.io/chuangyeshuo/mcprapi-backend:$VERSION

# 构建前端
echo "构建前端镜像..."
docker build -t ghcr.io/chuangyeshuo/mcprapi-frontend:$VERSION ./frontend
docker push ghcr.io/chuangyeshuo/mcprapi-frontend:$VERSION

echo "构建完成!"
```

使用方法：
```bash
chmod +x build-and-push.sh
./build-and-push.sh 1.2.0  # 构建特定版本
./build-and-push.sh        # 构建latest版本
```

## GitHub Actions自动化

### 创建自动构建和推送的工作流

创建 `.github/workflows/docker-publish.yml`:

```yaml
name: Docker Build and Push

on:
  push:
    branches: [ main ]
    tags: [ 'v*' ]
  pull_request:
    branches: [ main ]

env:
  REGISTRY: ghcr.io
  IMAGE_NAME_BACKEND: ${{ github.repository }}-backend
  IMAGE_NAME_FRONTEND: ${{ github.repository }}-frontend

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write

    steps:
    - name: Checkout repository
      uses: actions/checkout@v4

    - name: Log in to Container Registry
      uses: docker/login-action@v3
      with:
        registry: ${{ env.REGISTRY }}
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}

    - name: Extract metadata for backend
      id: meta-backend
      uses: docker/metadata-action@v5
      with:
        images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME_BACKEND }}
        tags: |
          type=ref,event=branch
          type=ref,event=pr
          type=semver,pattern={{version}}
          type=semver,pattern={{major}}.{{minor}}

    - name: Extract metadata for frontend
      id: meta-frontend
      uses: docker/metadata-action@v5
      with:
        images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME_FRONTEND }}
        tags: |
          type=ref,event=branch
          type=ref,event=pr
          type=semver,pattern={{version}}
          type=semver,pattern={{major}}.{{minor}}

    - name: Build and push backend image
      uses: docker/build-push-action@v5
      with:
        context: ./backend
        push: ${{ github.event_name != 'pull_request' }}
        tags: ${{ steps.meta-backend.outputs.tags }}
        labels: ${{ steps.meta-backend.outputs.labels }}

    - name: Build and push frontend image
      uses: docker/build-push-action@v5
      with:
        context: ./frontend
        push: ${{ github.event_name != 'pull_request' }}
        tags: ${{ steps.meta-frontend.outputs.tags }}
        labels: ${{ steps.meta-frontend.outputs.labels }}
```

## 使用GitHub Container Registry镜像

### 更新docker-compose.yml

将您的 `docker-compose.yml` 中的镜像地址更新为GHCR：

```yaml
services:
  backend:
    image: ghcr.io/chuangyeshuo/mcprapi-backend:latest
    # 其他配置...
    
  frontend:
    image: ghcr.io/chuangyeshuo/mcprapi-frontend:latest
    # 其他配置...
```

### 拉取和运行镜像

```bash
# 拉取最新镜像
docker pull ghcr.io/chuangyeshuo/mcprapi-backend:latest
docker pull ghcr.io/chuangyeshuo/mcprapi-frontend:latest

# 拉取特定版本
docker pull ghcr.io/chuangyeshuo/mcprapi-backend:1.2.0
docker pull ghcr.io/chuangyeshuo/mcprapi-frontend:1.2.0

# 使用docker-compose启动
docker-compose -f docker-compose.ghcr.yml up -d
```

### 环境变量配置

创建 `.env` 文件来管理版本：

```bash
# .env
BACKEND_VERSION=1.2.0
FRONTEND_VERSION=1.2.0
MYSQL_ROOT_PASSWORD=your_secure_password
REDIS_PASSWORD=your_redis_password
JWT_SECRET=your_jwt_secret
```

然后在 `docker-compose.ghcr.yml` 中使用：

```yaml
services:
  backend:
    image: ghcr.io/chuangyeshuo/mcprapi-backend:${BACKEND_VERSION:-latest}
  frontend:
    image: ghcr.io/chuangyeshuo/mcprapi-frontend:${FRONTEND_VERSION:-latest}
```

## 镜像管理

### 查看镜像信息

```bash
# 列出本地镜像
docker images ghcr.io/chuangyeshuo/mcprapi*

# 查看镜像详细信息
docker inspect ghcr.io/chuangyeshuo/mcprapi-backend:latest
```

### 清理镜像

```bash
# 删除本地镜像
docker rmi ghcr.io/chuangyeshuo/mcprapi-backend:old-version

# 清理未使用的镜像
docker image prune -f
```

### 在GitHub上管理

1. 访问 `https://github.com/chuangyeshuo/mcprapi/packages`
2. 可以查看、删除、设置权限

## 最佳实践

### 1. 镜像标签策略

```bash
# 语义化版本
ghcr.io/chuangyeshuo/mcprapi-backend:1.2.0
ghcr.io/chuangyeshuo/mcprapi-backend:1.2
ghcr.io/chuangyeshuo/mcprapi-backend:1

# 分支标签
ghcr.io/chuangyeshuo/mcprapi-backend:main
ghcr.io/chuangyeshuo/mcprapi-backend:develop

# 特殊标签
ghcr.io/chuangyeshuo/mcprapi-backend:latest
ghcr.io/chuangyeshuo/mcprapi-backend:stable
```

### 2. 多阶段构建优化

您的Dockerfile已经使用了多阶段构建，这很好：
- `builder` 阶段：编译代码
- `development` 阶段：开发环境
- `production` 阶段：生产环境（最小化镜像）

### 3. 安全考虑

- ✅ 使用非root用户运行
- ✅ 设置健康检查
- ✅ 使用只读文件系统
- ✅ 禁用新权限

### 4. 镜像大小优化

```bash
# 查看镜像大小
docker images ghcr.io/chuangyeshuo/mcprapi*

# 查看镜像层
docker history ghcr.io/chuangyeshuo/mcprapi-backend:latest
```

## 常用命令速查

### 镜像管理

```bash
# 查看本地镜像
docker images | grep ghcr.io/chuangyeshuo

# 删除本地镜像
docker rmi ghcr.io/chuangyeshuo/mcprapi-backend:latest
docker rmi ghcr.io/chuangyeshuo/mcprapi-frontend:latest

# 清理未使用的镜像
docker image prune -f

# 查看镜像详细信息
docker inspect ghcr.io/chuangyeshuo/mcprapi-backend:latest
```

### 容器管理

```bash
# 查看运行中的容器
docker ps

# 查看容器日志
docker logs mcprapi-backend
docker logs mcprapi-frontend

# 进入容器
docker exec -it mcprapi-backend /bin/sh
docker exec -it mcprapi-frontend /bin/sh

# 重启服务
docker-compose -f docker-compose.ghcr.yml restart

# 停止并删除容器
docker-compose -f docker-compose.ghcr.yml down
```

### 版本管理

```bash
# 列出所有可用版本（需要GitHub CLI）
gh api repos/chuangyeshuo/mcprapi/packages/container/mcprapi-backend/versions

# 拉取特定版本
docker pull ghcr.io/chuangyeshuo/mcprapi-backend:1.2.0

# 更新到最新版本
docker pull ghcr.io/chuangyeshuo/mcprapi-backend:latest
docker pull ghcr.io/chuangyeshuo/mcprapi-frontend:latest
docker-compose -f docker-compose.ghcr.yml up -d
```

### 基本操作

```bash
# 登录
docker login ghcr.io

# 构建
docker build -t ghcr.io/chuangyeshuo/mcprapi-backend:latest ./backend

# 推送
docker push ghcr.io/chuangyeshuo/mcprapi-backend:latest

# 拉取
docker pull ghcr.io/chuangyeshuo/mcprapi-backend:latest

# 运行
docker run -d --name backend ghcr.io/chuangyeshuo/mcprapi-backend:latest

# 查看日志
docker logs backend

# 进入容器
docker exec -it backend sh
```

## 故障排除

### 常见问题

1. **推送失败 - 权限不足**
   ```bash
   # 检查token权限，确保包含write:packages
   ```

2. **拉取失败 - 镜像不存在**
   ```bash
   # 检查镜像名称和标签是否正确
   docker pull ghcr.io/chuangyeshuo/mcprapi-backend:latest
   ```

3. **构建失败**
   ```bash
   # 检查Dockerfile语法和依赖
   docker build --no-cache -t test ./backend
   ```

## 下一步

1. 设置GitHub Actions自动构建
2. 配置镜像扫描和安全检查
3. 实现多架构镜像支持（ARM64/AMD64）
4. 设置镜像签名验证