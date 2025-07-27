# Docker 部署指南

## 端口配置

本项目已优化Docker部署配置，默认端口如下：

- **后端API服务**: 8081
- **前端Web服务**: 8082
- **MySQL数据库**: 3306
- **Redis缓存**: 6379

## 快速开始

### 1. 环境准备

确保您的系统已安装：
- Docker (>= 20.0)
- Docker Compose (>= 2.0)

检查环境：
```bash
make env-check
```

### 2. 初始化项目

```bash
# 初始化项目（会自动创建 .env 文件）
make init

# 编辑环境变量文件
vim .env
```

### 3. 生产环境部署

```bash
# 构建并启动生产环境
make docker-build
make docker-up

# 查看服务状态
make docker-ps

# 查看日志
make docker-logs

# 查看服务日志
docker-compose logs
```

**注意：** 
- 首次启动时会通过 MySQL 初始化脚本自动创建数据库结构
- 默认管理员账户会在后端服务启动时自动创建：`admin/admin`
- 后端服务会等待数据库完全启动后再启动

访问地址：
- 前端应用: http://localhost:8082
- 后端API: http://localhost:8081
- API文档: http://localhost:8081/swagger/index.html
- 健康检查: http://localhost:8081/health

### 4. 开发环境部署

```bash
# 启动开发环境（自动初始化数据库）
docker-compose -f docker-compose.dev.yml up -d

# 查看服务状态
docker-compose -f docker-compose.dev.yml ps

# 查看服务日志
docker-compose -f docker-compose.dev.yml logs
```

**注意：** 
- 首次启动时会通过 MySQL 初始化脚本自动创建数据库结构
- 默认管理员账户会在后端服务启动时自动创建：`admin/admin`
- 后端服务会等待数据库完全启动后再启动
- 开发环境支持热重载功能

开发环境额外服务：
- 健康检查: http://localhost:8081/health
- 数据库管理: http://localhost:8083 (Adminer)
- Redis管理: http://localhost:8084 (Redis Commander)

## 环境变量配置

主要环境变量说明：

```bash
# 端口配置
BACKEND_PORT=8081          # 后端服务端口
FRONTEND_PORT=8082         # 前端服务端口
MYSQL_PORT=3306           # MySQL端口
REDIS_PORT=6379           # Redis端口

# 数据库配置
MYSQL_ROOT_PASSWORD=your_strong_password
MYSQL_DATABASE=api_auth
MYSQL_USER=mcprapi
MYSQL_PASSWORD=your_mysql_password

# Redis配置
REDIS_PASSWORD=your_redis_password

# 安全配置
JWT_SECRET=your_very_long_jwt_secret_key
CORS_ORIGINS=http://localhost:8082,https://yourdomain.com

# 应用配置
APP_ENV=production
LOG_LEVEL=info
RATE_LIMIT=100
```

## Docker Compose 文件说明

### 生产环境 (docker-compose.yml)

- **多阶段构建**: 优化镜像大小和构建效率
- **健康检查**: 确保服务正常启动
- **安全配置**: 非root用户运行，只读文件系统
- **性能优化**: 内存限制，连接池配置
- **依赖管理**: 服务启动顺序控制

### 开发环境 (docker-compose.dev.yml)

- **热重载**: 支持代码修改实时生效
- **调试工具**: 集成数据库和Redis管理工具
- **开发配置**: 详细日志，调试模式
- **卷挂载**: 源代码实时同步

## 常用命令

### 生产环境管理

```bash
# 启动服务
make docker-up

# 停止服务
make docker-down

# 重启服务
make docker-restart

# 查看日志
make docker-logs

# 重新构建
make docker-build
```

### 开发环境管理

```bash
# 启动开发环境
make docker-dev-up

# 停止开发环境
make docker-dev-down

# 查看开发日志
make docker-dev-logs

# 重新构建开发环境
make docker-dev-build
```

### 容器管理

```bash
# 查看容器状态
make docker-ps

# 进入后端容器
make docker-exec-backend

# 进入前端容器
make docker-exec-frontend

# 进入MySQL容器
make docker-exec-mysql

# 进入Redis容器
make docker-exec-redis

# 清理所有资源
make docker-clean
```

## 安全特性

### 容器安全
- 非root用户运行
- 只读文件系统
- 最小权限原则
- 安全选项配置

### 网络安全
- 内部网络隔离
- CORS配置
- 安全头设置
- 限流保护

### 数据安全
- 数据卷持久化
- 密码加密存储
- JWT令牌保护
- 数据库连接加密

## 性能优化

### 镜像优化
- 多阶段构建
- Alpine Linux基础镜像
- 依赖缓存优化
- 镜像层最小化

### 运行时优化
- 内存限制配置
- 连接池优化
- 缓存策略
- Gzip压缩

### 监控配置
- 健康检查
- 日志收集
- 性能指标
- 错误追踪

## 故障排除

### 常见问题

1. **端口冲突**
   ```bash
   # 检查端口占用
   lsof -i :8081
   lsof -i :8082
   
   # 修改 .env 文件中的端口配置
   ```

2. **容器启动失败**
   ```bash
   # 查看详细日志
   docker-compose logs [service_name]
   
   # 检查健康状态
   docker-compose ps
   ```

3. **数据库连接失败**
   ```bash
   # 检查MySQL容器状态
   docker exec -it mcprapi-mysql mysqladmin ping -h localhost -u root -p
   
   # 查看数据库日志
   docker logs mcprapi-mysql
   ```

4. **前端无法访问后端**
   ```bash
   # 检查网络连接
   docker network ls
   docker network inspect mcprapi-network
   
   # 检查CORS配置
   ```

### 日志查看

```bash
# 查看所有服务日志
make docker-logs

# 查看特定服务日志
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f mysql
docker-compose logs -f redis
```

### 数据备份

```bash
# 备份MySQL数据
docker exec mcprapi-mysql mysqldump -u root -p api_auth > backup.sql

# 备份Redis数据
docker exec mcprapi-redis redis-cli --rdb /data/backup.rdb

# 备份Docker卷
docker run --rm -v mcprapi_mysql-data:/data -v $(pwd):/backup alpine tar czf /backup/mysql-backup.tar.gz /data
```

## 部署清单

部署前检查清单：

- [ ] 环境变量配置完成
- [ ] 端口无冲突
- [ ] Docker环境正常
- [ ] 防火墙配置
- [ ] SSL证书配置（生产环境）
- [ ] 域名解析配置（生产环境）
- [ ] 监控配置
- [ ] 备份策略

## 更新升级

```bash
# 停止服务
make docker-down

# 拉取最新代码
git pull

# 重新构建镜像
make docker-build

# 启动服务
make docker-up

# 验证服务
make docker-ps
```

## 技术支持

如遇到问题，请：

1. 查看本文档的故障排除部分
2. 检查项目日志文件
3. 查看GitHub Issues
4. 联系技术支持团队