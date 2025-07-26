# 🚀 MCP RAPI 快速启动指南

## 📋 系统要求

- **Docker** >= 20.0
- **Docker Compose** >= 2.0
- **内存** >= 4GB
- **磁盘空间** >= 2GB

## ⚡ 一键启动

### 🐳 开发环境（推荐）

```bash
# 1. 克隆项目
git clone https://github.com/chuangyeshuo/mcprapi.git
cd mcprapi

# 2. 一键启动（自动初始化数据库）
docker-compose -f docker-compose.dev.yml up -d

# 3. 查看启动状态
docker-compose -f docker-compose.dev.yml ps

# 4. 查看初始化日志（可选）
docker-compose -f docker-compose.dev.yml logs db-init-dev
```

**✨ 自动化特性：**
- 🔄 **自动数据库初始化**：首次启动时自动创建管理员账户
- 🔥 **热重载开发**：代码修改后自动重启
- 🛠️ **开发工具集成**：包含数据库管理和Redis管理工具

**开发环境服务地址：**
- 🌐 **前端应用**: http://localhost:8082
- 🔧 **后端API**: http://localhost:8081
- 📚 **API文档**: http://localhost:8081/swagger/index.html
- ❤️ **健康检查**: http://localhost:8081/health
- 🗄️ **数据库管理**: http://localhost:8083 (Adminer)
- 📊 **Redis管理**: http://localhost:8084 (Redis Commander)

### 🏭 生产环境

```bash
# 1. 克隆项目
git clone https://github.com/chuangyeshuo/mcprapi.git
cd mcprapi

# 2. 配置环境变量（可选）
cp .env.example .env
# 编辑 .env 文件设置生产环境参数

# 3. 一键启动（自动初始化数据库）
docker-compose up -d

# 4. 查看启动状态
docker-compose ps

# 5. 查看初始化日志（可选）
docker-compose logs db-init
```

**🔒 生产环境特性：**
- 🛡️ **安全优化**：只读文件系统、安全配置
- 🚀 **性能优化**：生产级别的资源配置
- 🔄 **自动初始化**：首次启动自动创建管理员账户

**生产环境服务地址：**
- 🌐 **前端应用**: http://localhost:8082
- 🔧 **后端API**: http://localhost:8081
- 📚 **API文档**: http://localhost:8081/swagger/index.html
- ❤️ **健康检查**: http://localhost:8081/health

## 🔑 默认登录信息

```
用户名: admin
密码: admin
```

> ⚠️ **安全提醒**: 首次登录后请立即修改默认密码！

## 📊 服务状态检查

### 查看所有容器状态
```bash
# 开发环境
docker-compose -f docker-compose.dev.yml ps

# 生产环境
docker-compose ps
```

### 查看服务日志
```bash
# 查看所有服务日志
docker-compose -f docker-compose.dev.yml logs -f

# 查看特定服务日志
docker-compose -f docker-compose.dev.yml logs -f backend-dev
docker-compose -f docker-compose.dev.yml logs -f frontend-dev
```

### 健康检查状态
```bash
# 检查后端健康状态
curl http://localhost:8081/health

# 预期响应
{
  "status": "ok",
  "timestamp": 1703123456,
  "service": "mcprapi-backend",
  "version": "1.0.0"
}
```

## 🛠️ 常见问题解决

### 1. 容器启动失败

**检查端口占用：**
```bash
# 检查端口是否被占用
lsof -i :8081  # 后端端口
lsof -i :8082  # 前端端口
lsof -i :3307  # MySQL端口（开发环境）
lsof -i :6380  # Redis端口（开发环境）
```

**解决方案：**
- 停止占用端口的进程
- 或修改 `.env` 文件中的端口配置

### 2. 数据库连接失败

**检查数据库状态：**
```bash
docker-compose -f docker-compose.dev.yml logs mysql-dev
```

**解决方案：**
- 等待数据库完全启动（约30-60秒）
- 检查数据库密码配置
- 重启数据库容器

### 3. 健康检查失败

**检查健康检查状态：**
```bash
docker inspect mcprapi-backend-dev | grep -A 10 "Health"
```

**解决方案：**
- 确保 `/health` 路由已注册（已修复）
- 检查容器内网络连接
- 查看后端服务日志

### 4. 前端无法访问后端

**检查网络连接：**
```bash
# 进入前端容器测试连接
docker exec -it mcprapi-frontend-dev sh
wget -O- http://backend-dev:8080/health
```

**解决方案：**
- 检查 Docker 网络配置
- 确认后端服务已启动
- 检查防火墙设置

## 🔧 开发工具

### 进入容器调试
```bash
# 进入后端容器
docker exec -it mcprapi-backend-dev sh

# 进入前端容器
docker exec -it mcprapi-frontend-dev sh

# 进入数据库容器
docker exec -it mcprapi-mysql-dev mysql -u root -p
```

### 重启特定服务
```bash
# 重启后端服务
docker-compose -f docker-compose.dev.yml restart backend-dev

# 重启前端服务
docker-compose -f docker-compose.dev.yml restart frontend-dev
```

### 查看实时日志
```bash
# 实时查看后端日志
docker-compose -f docker-compose.dev.yml logs -f backend-dev

# 实时查看前端日志
docker-compose -f docker-compose.dev.yml logs -f frontend-dev
```

## 🧹 清理和重置

### 停止所有服务
```bash
# 开发环境
docker-compose -f docker-compose.dev.yml down

# 生产环境
docker-compose down
```

### 完全清理（包括数据）
```bash
# ⚠️ 警告：这将删除所有数据！
docker-compose -f docker-compose.dev.yml down -v
docker system prune -f
```

### 重新构建镜像
```bash
# 重新构建并启动
docker-compose -f docker-compose.dev.yml up -d --build
```

## 📈 性能监控

### 资源使用情况
```bash
# 查看容器资源使用
docker stats

# 查看磁盘使用
docker system df
```

### 数据库性能
```bash
# 连接到数据库查看状态
docker exec -it mcprapi-mysql-dev mysql -u root -p -e "SHOW PROCESSLIST;"
```

## 🔗 相关文档

- [📖 完整文档](README.md)
- [🏗️ 系统架构](API多租户授权管理系统架构文档.md)
- [🐳 Docker部署](DOCKER_DEPLOYMENT.md)
- [🗄️ 数据库初始化](DATABASE_INIT.md)
- [🤖 MCP集成](mcp-example/FASTMCP_HTTP_CONFIG.md)

## 🆘 获取帮助

如果遇到问题，请：

1. 查看本文档的常见问题部分
2. 检查 [GitHub Issues](https://github.com/chuangyeshuo/mcprapi/issues)
3. 提交新的 Issue 并提供详细信息：
   - 操作系统版本
   - Docker 版本
   - 错误日志
   - 复现步骤

---

**🎉 恭喜！您的 MCP RAPI 系统已成功启动！**