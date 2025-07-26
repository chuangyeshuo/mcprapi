# 🚀 MCP RAPI 完整部署指南

> **MCP RAPI** - 现代化API权限管理系统，支持多租户架构和MCP协议集成

## 📋 目录

- [🎯 部署方式对比](#-部署方式对比)
- [🐳 Docker Compose 一键部署（推荐）](#-docker-compose-一键部署推荐)
- [💻 手动部署](#-手动部署)
- [⚙️ 配置文件详解](#️-配置文件详解)
- [🔧 环境变量说明](#-环境变量说明)
- [🛠️ 故障排除](#️-故障排除)
- [📊 性能优化](#-性能优化)
- [🔒 安全配置](#-安全配置)
- [📈 监控与日志](#-监控与日志)

## 🎯 部署方式对比

| 特性 | Docker Compose 部署 | 手动部署 |
|------|-------------------|----------|
| **部署难度** | ⭐ 极简单 | ⭐⭐⭐⭐ 复杂 |
| **环境一致性** | ✅ 完全一致 | ❌ 依赖本地环境 |
| **依赖管理** | ✅ 自动处理 | ❌ 手动安装 |
| **服务隔离** | ✅ 容器隔离 | ❌ 共享系统资源 |
| **扩展性** | ✅ 易于扩展 | ❌ 需要手动配置 |
| **回滚能力** | ✅ 版本控制 | ❌ 手动备份 |
| **开发调试** | ✅ 热重载支持 | ✅ 直接调试 |
| **资源占用** | 📊 中等 | 📊 较低 |
| **学习成本** | 📚 低 | 📚 高 |

### 🏆 推荐选择

- **🐳 Docker Compose**: 适合生产环境、团队协作、快速部署
- **💻 手动部署**: 适合学习研究、资源受限、深度定制

---

## 🐳 Docker Compose 一键部署（推荐）

### 🌟 特性优势

- **🚀 一键启动**: 单条命令完成所有服务部署
- **🔄 自动初始化**: 首次启动自动创建数据库和管理员账户
- **🔥 热重载**: 开发环境支持代码修改实时生效
- **🛠️ 管理工具**: 集成数据库和Redis管理界面
- **🛡️ 安全配置**: 生产级安全设置
- **📊 健康检查**: 自动监控服务状态

### 📦 环境要求

```bash
# 检查Docker版本
docker --version          # >= 20.0
docker-compose --version  # >= 2.0

# 检查系统资源
free -h                   # 内存 >= 4GB
df -h                     # 磁盘 >= 2GB
```

### 🚀 开发环境部署

```bash
# 1. 克隆项目
git clone https://github.com/chuangyeshuo/mcprapi.git
cd mcprapi

# 2. 一键启动开发环境
docker-compose -f docker-compose.dev.yml up -d

# 3. 查看启动状态
docker-compose -f docker-compose.dev.yml ps

# 4. 等待初始化完成（约2-3分钟）
docker-compose -f docker-compose.dev.yml logs -f backend-dev
```

**🌐 开发环境访问地址：**

| 服务 | 地址 | 说明 |
|------|------|------|
| 🎨 **前端应用** | http://localhost:8082 | Vue.js 管理界面 |
| 🔧 **后端API** | http://localhost:8081 | Go API 服务 |
| 📚 **API文档** | http://localhost:8081/swagger/index.html | Swagger 文档 |
| ❤️ **健康检查** | http://localhost:8081/health | 服务状态 |
| 🗄️ **数据库管理** | http://localhost:8083 | Adminer 管理界面 |
| 📊 **Redis管理** | http://localhost:8084 | Redis Commander |

**🔑 默认登录信息：**
```
用户名: admin
密码: admin
```

### 🏭 生产环境部署

```bash
# 1. 配置环境变量
cp .env.example .env
vim .env  # 修改生产环境配置

# 2. 一键启动生产环境
docker-compose up -d

# 3. 查看启动状态
docker-compose ps

# 4. 检查服务健康状态
curl http://localhost:8081/health
```

**🌐 生产环境访问地址：**

| 服务 | 地址 | 说明 |
|------|------|------|
| 🎨 **前端应用** | http://localhost:8082 | 生产级前端 |
| 🔧 **后端API** | http://localhost:8081 | 生产级API |
| 📚 **API文档** | http://localhost:8081/swagger/index.html | API文档 |
| ❤️ **健康检查** | http://localhost:8081/health | 服务监控 |

### 🔧 Docker Compose 管理命令

```bash
# 服务管理
docker-compose -f docker-compose.dev.yml up -d      # 启动开发环境
docker-compose -f docker-compose.dev.yml down       # 停止开发环境
docker-compose -f docker-compose.dev.yml restart    # 重启开发环境
docker-compose -f docker-compose.dev.yml ps         # 查看服务状态

docker-compose up -d                                 # 启动生产环境
docker-compose down                                  # 停止生产环境
docker-compose restart                               # 重启生产环境
docker-compose ps                                    # 查看服务状态

# 日志查看
docker-compose -f docker-compose.dev.yml logs -f backend-dev    # 查看后端日志
docker-compose -f docker-compose.dev.yml logs -f frontend-dev   # 查看前端日志
docker-compose -f docker-compose.dev.yml logs -f mysql-dev      # 查看数据库日志

# 容器操作
docker exec -it mcprapi-backend-dev sh              # 进入后端容器
docker exec -it mcprapi-frontend-dev sh             # 进入前端容器
docker exec -it mcprapi-mysql-dev mysql -u root -p  # 进入数据库

# 数据管理
docker-compose -f docker-compose.dev.yml down -v    # 删除所有数据
docker volume ls                                     # 查看数据卷
docker volume rm mcprapi_mysql-dev-data             # 删除数据库数据
```

---

## 💻 手动部署

### 📋 环境要求

```bash
# 后端要求
Go >= 1.21
MySQL >= 8.0
Redis >= 6.0

# 前端要求
Node.js >= 18.0
npm >= 8.0

# 系统要求
Linux/macOS/Windows
内存 >= 2GB
磁盘 >= 1GB
```

### 🗄️ 数据库准备

```bash
# 1. 安装MySQL
# Ubuntu/Debian
sudo apt update
sudo apt install mysql-server

# CentOS/RHEL
sudo yum install mysql-server

# macOS
brew install mysql

# 2. 启动MySQL服务
sudo systemctl start mysql
sudo systemctl enable mysql

# 3. 创建数据库和用户
mysql -u root -p
```

```sql
-- 创建数据库
CREATE DATABASE api_auth CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建用户
CREATE USER 'mcprapi'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON api_auth.* TO 'mcprapi'@'localhost';
FLUSH PRIVILEGES;
```

### 📊 Redis 准备

```bash
# 1. 安装Redis
# Ubuntu/Debian
sudo apt install redis-server

# CentOS/RHEL
sudo yum install redis

# macOS
brew install redis

# 2. 启动Redis服务
sudo systemctl start redis
sudo systemctl enable redis

# 3. 配置Redis密码（可选）
sudo vim /etc/redis/redis.conf
# 取消注释并设置: requirepass your_redis_password

# 重启Redis
sudo systemctl restart redis
```

### 🔧 后端部署

```bash
# 1. 进入后端目录
cd backend

# 2. 安装Go依赖
go mod download

# 3. 复制配置文件
cp configs/dev.yaml.example configs/dev.yaml

# 4. 编辑配置文件
vim configs/dev.yaml
```

**配置文件示例 (configs/dev.yaml):**
```yaml
# 服务器配置
server:
  port: 8081
  mode: debug

# 数据库配置
mysql:
  host: localhost
  port: 3306
  username: mcprapi
  password: your_password
  database: api_auth
  charset: utf8mb4

# Redis配置
redis:
  host: localhost
  port: 6379
  password: your_redis_password
  db: 0

# JWT配置
jwt:
  secret: your_very_long_jwt_secret_key
  expire: 24h

# 日志配置
log:
  level: debug
  file_path: ./logs/app.log
  file_mode: true
  format_str: "[%time%] [%level%] %msg%"

# Casbin配置
casbin:
  model: configs/casbin_model.conf
```

```bash
# 5. 初始化数据库
go run scripts/init_admin.go

# 6. 启动后端服务
go run cmd/main.go --config configs/dev.yaml

# 🚀 后端服务运行在 http://localhost:8081
```

### 🎨 前端部署

```bash
# 1. 进入前端目录
cd frontend

# 2. 安装Node.js依赖
npm install

# 3. 配置API地址
vim .env.development
```

**环境配置文件 (.env.development):**
```bash
# API配置
VUE_APP_API_BASE_URL=http://localhost:8081

# 应用配置
VUE_APP_TITLE=MCP RAPI 管理系统
VUE_APP_VERSION=1.0.0

# 开发配置
NODE_ENV=development
```

```bash
# 4. 启动开发服务器
npm run serve

# 🎨 前端服务运行在 http://localhost:8080
```

### 🏭 生产环境构建

```bash
# 后端生产构建
cd backend
go build -o mcprapi cmd/main.go

# 启动生产服务
./mcprapi --config configs/prod.yaml

# 前端生产构建
cd frontend
npm run build

# 使用Nginx部署
sudo cp -r dist/* /var/www/html/
```

---

## ⚙️ 配置文件详解

### 🐳 Docker Compose 配置差异

#### 开发环境 (docker-compose.dev.yml)

```yaml
# 特性：热重载、调试工具、详细日志
services:
  backend-dev:
    build:
      target: development  # 开发构建目标
    environment:
      GIN_MODE: debug     # 调试模式
      LOG_LEVEL: debug    # 详细日志
      HOT_RELOAD: true    # 热重载
    volumes:
      - ./backend:/app    # 源码挂载
    ports:
      - "8081:8080"       # 端口映射

  # 额外的管理工具
  adminer:              # 数据库管理
    image: adminer:4.8.1
    ports:
      - "8083:8080"

  redis-commander:      # Redis管理
    image: rediscommander/redis-commander
    ports:
      - "8084:8081"
```

#### 生产环境 (docker-compose.yml)

```yaml
# 特性：安全优化、性能优化、最小权限
services:
  backend:
    build:
      target: production  # 生产构建目标
    environment:
      APP_ENV: production # 生产模式
      LOG_LEVEL: info     # 简化日志
    security_opt:
      - no-new-privileges:true  # 安全配置
    read_only: true       # 只读文件系统
    tmpfs:
      - /tmp:noexec,nosuid,size=100m  # 临时文件系统
    
    # 无管理工具，专注性能和安全
```

### 🔧 主要差异对比

| 配置项 | 开发环境 | 生产环境 |
|--------|----------|----------|
| **构建目标** | development | production |
| **日志级别** | debug | info |
| **热重载** | ✅ 启用 | ❌ 禁用 |
| **源码挂载** | ✅ 实时同步 | ❌ 构建时复制 |
| **管理工具** | ✅ Adminer + Redis Commander | ❌ 无 |
| **安全配置** | 🔓 宽松 | 🔒 严格 |
| **文件系统** | 📝 可写 | 📖 只读 |
| **端口配置** | 🔓 多端口暴露 | 🔒 最小端口 |
| **资源限制** | 🔓 无限制 | 🔒 严格限制 |

---

## 🔧 环境变量说明

### 📝 .env 文件配置

```bash
# 复制环境变量模板
cp .env.example .env
```

**完整环境变量说明：**

```bash
# ===========================================
# 🌐 服务端口配置
# ===========================================
BACKEND_PORT=8081          # 后端API服务端口
FRONTEND_PORT=8082         # 前端Web服务端口
MYSQL_PORT=3306           # MySQL数据库端口（生产）
MYSQL_DEV_PORT=3307       # MySQL数据库端口（开发）
REDIS_PORT=6379           # Redis缓存端口（生产）
REDIS_DEV_PORT=6380       # Redis缓存端口（开发）

# ===========================================
# 🗄️ 数据库配置
# ===========================================
MYSQL_ROOT_PASSWORD=your_strong_root_password
MYSQL_DATABASE=api_auth                    # 生产数据库名
MYSQL_DEV_DATABASE=api_auth_dev           # 开发数据库名
MYSQL_USER=mcprapi
MYSQL_PASSWORD=your_mysql_password

# ===========================================
# 📊 Redis配置
# ===========================================
REDIS_PASSWORD=your_redis_password

# ===========================================
# 🔐 安全配置
# ===========================================
JWT_SECRET=your_very_long_jwt_secret_key_at_least_32_characters
CORS_ORIGINS=http://localhost:8082,https://yourdomain.com

# ===========================================
# 🚀 应用配置
# ===========================================
APP_ENV=production                         # 应用环境: development/production
LOG_LEVEL=info                            # 日志级别: debug/info/warn/error
RATE_LIMIT=100                            # API限流: 每分钟请求数

# ===========================================
# 🌍 外部访问配置
# ===========================================
API_BASE_URL=http://localhost:8081       # API基础URL
FRONTEND_BASE_URL=http://localhost:8082   # 前端基础URL

# ===========================================
# 🔧 构建配置
# ===========================================
BUILD_ENV=production                      # 构建环境
NODE_ENV=production                       # Node.js环境
```

### 🔒 安全建议

```bash
# 生成强密码
openssl rand -base64 32  # 生成32字符随机密码

# JWT密钥生成
openssl rand -hex 64     # 生成64字符十六进制密钥

# 检查密码强度
echo "your_password" | pwscore  # 密码评分（需安装libpwquality-tools）
```

---

## 🛠️ 故障排除

### 🔍 常见问题诊断

#### 1. 🐳 Docker 相关问题

**容器启动失败：**
```bash
# 检查容器状态
docker-compose -f docker-compose.dev.yml ps

# 查看容器日志
docker-compose -f docker-compose.dev.yml logs backend-dev

# 检查端口占用
netstat -tulpn | grep :8081
lsof -i :8081

# 清理并重新启动
docker-compose -f docker-compose.dev.yml down -v
docker-compose -f docker-compose.dev.yml up -d
```

**数据库连接失败：**
```bash
# 检查数据库容器状态
docker-compose -f docker-compose.dev.yml logs mysql-dev

# 进入数据库容器测试连接
docker exec -it mcprapi-mysql-dev mysql -u root -p

# 检查网络连接
docker network ls
docker network inspect mcprapi_mcprapi-dev-network
```

**健康检查失败：**
```bash
# 检查健康检查状态
docker inspect mcprapi-backend-dev | grep -A 10 "Health"

# 手动测试健康检查
curl -f http://localhost:8081/health

# 进入容器内部测试
docker exec -it mcprapi-backend-dev wget -O- http://localhost:8080/health
```

#### 2. 💻 手动部署问题

**Go 编译错误：**
```bash
# 清理模块缓存
go clean -modcache
go mod download

# 检查Go版本
go version  # 需要 >= 1.21

# 更新依赖
go mod tidy
```

**数据库连接错误：**
```bash
# 测试数据库连接
mysql -h localhost -u mcprapi -p api_auth

# 检查MySQL服务状态
sudo systemctl status mysql

# 查看MySQL错误日志
sudo tail -f /var/log/mysql/error.log
```

**Redis连接错误：**
```bash
# 测试Redis连接
redis-cli -h localhost -p 6379 -a your_password ping

# 检查Redis服务状态
sudo systemctl status redis

# 查看Redis日志
sudo tail -f /var/log/redis/redis-server.log
```

**前端构建错误：**
```bash
# 清理node_modules
rm -rf node_modules package-lock.json
npm install

# 检查Node.js版本
node --version  # 需要 >= 18.0
npm --version   # 需要 >= 8.0

# 使用yarn替代npm
npm install -g yarn
yarn install
yarn serve
```

### 📊 性能监控

```bash
# 查看容器资源使用
docker stats

# 查看系统资源
htop
free -h
df -h

# 查看网络连接
netstat -tulpn
ss -tulpn
```

### 🔧 维护命令

```bash
# Docker清理
docker system prune -a              # 清理未使用的镜像
docker volume prune                 # 清理未使用的数据卷
docker network prune                # 清理未使用的网络

# 日志清理
docker-compose -f docker-compose.dev.yml logs --tail=100 > logs_backup.txt
sudo truncate -s 0 /var/lib/docker/containers/*/*-json.log

# 数据备份
docker exec mcprapi-mysql-dev mysqldump -u root -p api_auth_dev > backup.sql

# 数据恢复
docker exec -i mcprapi-mysql-dev mysql -u root -p api_auth_dev < backup.sql
```

---

## 📊 性能优化

### 🚀 Docker 性能优化

#### 容器资源限制
```yaml
# docker-compose.yml 中添加资源限制
services:
  backend:
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 2G
        reservations:
          cpus: '0.5'
          memory: 512M
    
  mysql:
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 1G
        reservations:
          cpus: '0.25'
          memory: 256M
```

#### 数据库性能优化
```sql
-- MySQL 配置优化 (my.cnf)
[mysqld]
# 连接池配置
max_connections = 200
max_connect_errors = 10000

# 缓存配置
innodb_buffer_pool_size = 1G
query_cache_size = 256M
query_cache_type = 1

# 日志配置
slow_query_log = 1
slow_query_log_file = /var/log/mysql/slow.log
long_query_time = 2

# 索引优化
innodb_flush_log_at_trx_commit = 2
innodb_log_file_size = 256M
```

#### Redis 性能优化
```bash
# Redis 配置优化 (redis.conf)
# 内存配置
maxmemory 512mb
maxmemory-policy allkeys-lru

# 持久化配置
save 900 1
save 300 10
save 60 10000

# 网络配置
tcp-keepalive 300
timeout 0
```

### 🔧 应用性能优化

#### Go 后端优化
```yaml
# configs/prod.yaml
server:
  read_timeout: 30s
  write_timeout: 30s
  idle_timeout: 120s
  max_header_bytes: 1048576

# 数据库连接池
mysql:
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: 3600s

# Redis 连接池
redis:
  pool_size: 10
  min_idle_conns: 5
  dial_timeout: 5s
  read_timeout: 3s
  write_timeout: 3s
```

#### 前端优化
```javascript
// vue.config.js
module.exports = {
  productionSourceMap: false,
  configureWebpack: {
    optimization: {
      splitChunks: {
        chunks: 'all',
        cacheGroups: {
          vendor: {
            name: 'chunk-vendors',
            test: /[\\/]node_modules[\\/]/,
            priority: 10,
            chunks: 'initial'
          }
        }
      }
    }
  }
}
```

---

## 🔒 安全配置

### 🛡️ Docker 安全

#### 容器安全配置
```yaml
# docker-compose.yml 安全配置
services:
  backend:
    security_opt:
      - no-new-privileges:true
      - apparmor:docker-default
    read_only: true
    tmpfs:
      - /tmp:noexec,nosuid,size=100m
    user: "1000:1000"
    cap_drop:
      - ALL
    cap_add:
      - NET_BIND_SERVICE
```

#### 网络安全
```yaml
# 自定义网络配置
networks:
  mcprapi-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16
    driver_opts:
      com.docker.network.bridge.name: mcprapi-br
      com.docker.network.bridge.enable_icc: "false"
```

### 🔐 应用安全

#### JWT 安全配置
```yaml
# configs/prod.yaml
jwt:
  secret: ${JWT_SECRET}  # 至少64字符
  expire: 2h             # 短期过期时间
  refresh_expire: 168h   # 7天刷新令牌
  issuer: "mcprapi"
  algorithm: "HS256"
```

#### CORS 安全配置
```yaml
cors:
  allowed_origins:
    - "https://yourdomain.com"
    - "https://admin.yourdomain.com"
  allowed_methods:
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
  allowed_headers:
    - "Authorization"
    - "Content-Type"
  max_age: 86400
```

#### 数据库安全
```sql
-- 创建只读用户
CREATE USER 'mcprapi_readonly'@'%' IDENTIFIED BY 'strong_password';
GRANT SELECT ON api_auth.* TO 'mcprapi_readonly'@'%';

-- 创建备份用户
CREATE USER 'mcprapi_backup'@'localhost' IDENTIFIED BY 'backup_password';
GRANT SELECT, LOCK TABLES ON api_auth.* TO 'mcprapi_backup'@'localhost';

-- 删除默认用户
DROP USER IF EXISTS ''@'localhost';
DROP USER IF EXISTS ''@'%';
```

### 🔒 SSL/TLS 配置

#### Nginx SSL 配置
```nginx
server {
    listen 443 ssl http2;
    server_name yourdomain.com;
    
    ssl_certificate /etc/ssl/certs/yourdomain.crt;
    ssl_certificate_key /etc/ssl/private/yourdomain.key;
    
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;
    ssl_prefer_server_ciphers off;
    
    add_header Strict-Transport-Security "max-age=63072000" always;
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    
    location / {
        proxy_pass http://localhost:8082;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    location /api {
        proxy_pass http://localhost:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

---

## 📈 监控与日志

### 📊 系统监控

#### Prometheus + Grafana 监控
```yaml
# docker-compose.monitoring.yml
version: '3.8'
services:
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus-data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/etc/prometheus/console_libraries'
      - '--web.console.templates=/etc/prometheus/consoles'

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    volumes:
      - grafana-data:/var/lib/grafana
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin

volumes:
  prometheus-data:
  grafana-data:
```

#### 应用指标收集
```go
// 在 Go 应用中添加 Prometheus 指标
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
        },
        []string{"method", "endpoint"},
    )
)
```

### 📝 日志管理

#### ELK Stack 日志收集
```yaml
# docker-compose.logging.yml
version: '3.8'
services:
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:7.15.0
    environment:
      - discovery.type=single-node
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
    ports:
      - "9200:9200"
    volumes:
      - elasticsearch-data:/usr/share/elasticsearch/data

  logstash:
    image: docker.elastic.co/logstash/logstash:7.15.0
    volumes:
      - ./logging/logstash.conf:/usr/share/logstash/pipeline/logstash.conf
    ports:
      - "5044:5044"
    depends_on:
      - elasticsearch

  kibana:
    image: docker.elastic.co/kibana/kibana:7.15.0
    ports:
      - "5601:5601"
    environment:
      - ELASTICSEARCH_HOSTS=http://elasticsearch:9200
    depends_on:
      - elasticsearch

volumes:
  elasticsearch-data:
```

#### 结构化日志配置
```yaml
# configs/prod.yaml
log:
  level: info
  format: json
  output: stdout
  fields:
    service: mcprapi
    version: 1.0.0
  file:
    enabled: true
    path: /var/log/mcprapi/app.log
    max_size: 100MB
    max_backups: 10
    max_age: 30
    compress: true
```

### 🚨 告警配置

#### Prometheus 告警规则
```yaml
# monitoring/alert-rules.yml
groups:
  - name: mcprapi-alerts
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value }} errors per second"

      - alert: DatabaseConnectionFailed
        expr: mysql_up == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Database connection failed"
          description: "MySQL database is not responding"

      - alert: HighMemoryUsage
        expr: container_memory_usage_bytes / container_spec_memory_limit_bytes > 0.9
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High memory usage"
          description: "Container memory usage is above 90%"
```

#### 健康检查端点
```go
// 健康检查实现
func (h *HealthHandler) Check(c *gin.Context) {
    status := map[string]interface{}{
        "status": "healthy",
        "timestamp": time.Now().Unix(),
        "version": "1.0.0",
        "checks": map[string]interface{}{
            "database": h.checkDatabase(),
            "redis": h.checkRedis(),
            "casbin": h.checkCasbin(),
        },
    }
    
    c.JSON(http.StatusOK, status)
}
```

---

## 📚 相关文档

- [🚀 快速启动指南](QUICK_START.md)
- [🏗️ 系统架构文档](API多租户授权管理系统架构文档.md)
- [🗃️ 数据库初始化指南](DATABASE_INIT.md)
- [🏢 部门管理流程](新增部门流程文档.md)
- [🤖 MCP集成指南](mcp-example/FASTMCP_HTTP_CONFIG.md)

---

## 🤝 技术支持

如果您在部署过程中遇到问题，请：

1. 📖 查看相关文档
2. 🔍 搜索已知问题
3. 🐛 提交Issue
4. 💬 参与讨论

**联系方式：**
- 📧 Email: support@mcprapi.com
- 💬 GitHub Issues: [提交问题](https://github.com/chuangyeshuo/mcprapi/issues)
- 📱 Discord: [加入讨论](https://discord.gg/mcprapi)

---

## 📄 许可证

本项目采用 [MIT License](LICENSE) 开源协议。

## 🙏 致谢

感谢所有为 MCP RAPI 项目做出贡献的开发者和社区成员！

---

**🎯 快速链接：**
- [⚡ 一键启动开发环境](#-开发环境部署)
- [🏭 生产环境部署](#-生产环境部署)
- [🛠️ 故障排除指南](#️-故障排除)
- [📊 性能优化建议](#-性能优化)