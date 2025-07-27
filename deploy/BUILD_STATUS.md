# 构建状态

## 最新构建信息

**构建时间**: 2024年12月

**构建状态**: ✅ 成功

## 镜像详情

### 后端镜像
- **镜像地址**: `ghcr.io/chuangyeshuo/mcprapi-backend:latest`
- **镜像摘要**: `sha256:fdc3ce25d58fe4ce1da2db004d4b4c38b3c188d35de8e58087727d1cb11727f0`
- **构建状态**: ✅ 构建成功
- **推送状态**: ✅ 推送成功

### 前端镜像
- **镜像地址**: `ghcr.io/chuangyeshuo/mcprapi-frontend:latest`
- **镜像摘要**: `sha256:e1024a94bb9f0491b47324c4561783c67773119f4a1205585fdc80205fbbd4ec`
- **构建状态**: ✅ 构建成功
- **推送状态**: ✅ 推送成功

## 构建日志摘要

### 后端构建
- 基础镜像: `golang:1.23-alpine`
- 运行时镜像: `alpine:3.18`
- 构建时间: 约 3-5 分钟
- 镜像大小: 优化后的多阶段构建

### 前端构建
- 基础镜像: `node:18-alpine`
- 运行时镜像: `nginx:alpine`
- 构建时间: 约 2-4 分钟
- 镜像大小: 优化后的多阶段构建

## 使用说明

### 快速部署
```bash
# 使用部署脚本
cd deploy
./deploy-ghcr.sh

# 或直接拉取镜像
docker pull ghcr.io/chuangyeshuo/mcprapi-backend:latest
docker pull ghcr.io/chuangyeshuo/mcprapi-frontend:latest
```

### 验证镜像
```bash
# 验证镜像摘要
docker inspect ghcr.io/chuangyeshuo/mcprapi-backend:latest | grep -A 1 "RepoDigests"
docker inspect ghcr.io/chuangyeshuo/mcprapi-frontend:latest | grep -A 1 "RepoDigests"
```

## 故障排除

如果遇到镜像拉取问题，请参考：
- [GitHub Container Registry 指南](./GITHUB_CONTAINER_REGISTRY_GUIDE.md)
- [部署文档](./README.md)
- [权限检查工具](./check-token-permissions.sh)

## 下次构建

要重新构建镜像，请运行：
```bash
cd deploy
./build-and-push.sh
```

---
*此文件由构建脚本自动更新*