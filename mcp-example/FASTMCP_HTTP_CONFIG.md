# FastMCP 2.x HTTP 配置指南

## 概述

本指南介绍如何配置和使用 FastMCP 2.x 版本的 HTTP 服务器，并将其与 Claude Desktop 集成。

## 服务器配置

### 1. 启动 FastMCP 服务器

运行以下命令启动服务器：

python fastmcp_http_server.py

服务器将启动并显示：

```
🚀 启动 FastMCP 服务器 (Streamable HTTP 模式)...
🔗 使用 Streamable HTTP 传输协议 (推荐)
📡 HTTP 端点: http://localhost:8000
📡 SSE 端点: http://localhost:8000/sse (兼容模式)

🏎️  FastMCP version: 2.10.6
🤝 MCP version:     1.12.2

Starting MCP server 'Cloud MCP Server' with transport 'http' on http://127.0.0.1:8000/mcp/
```

### 2. Claude Desktop 配置

将以下配置添加到您的 Claude Desktop 配置文件中：

```json
{
  "mcpServers": {
    "fastmcp-streamable-http": {
      "type": "streamableHttp",
      "url": "http://localhost:8000/mcp/",
      "headers": {
        "Content-Type": "application/json",
        "Authorization": "Bearer your-jwt-token"
      }
    }
  }
}
```

## 可用工具

FastMCP 服务器提供以下工具：

1. **get_weather_alerts** - 获取指定州的天气预警信息
2. **get_weather_forecast** - 获取指定位置的天气预报
3. **get_user_info** - 获取用户信息
4. **get_department_stats** - 获取部门统计信息
5. **create_order** - 创建订单（需要管理员权限）

## 使用步骤

1. 启动 FastMCP HTTP Stream 服务器
2. 确认服务器在 `http://localhost:8000` 上运行
3. 使用提供的配置文件配置 Claude Desktop
4. 重启 Claude Desktop 以加载新配置
5. 开始使用 FastMCP 工具

## FastMCP 2.x 的优势

### Streamable HTTP 传输

- ✅ **更高效**: 比 SSE 更好的性能和资源利用
- ✅ **更稳定**: 更可靠的连接管理
- ✅ **更标准**: 基于标准 HTTP 协议
- ✅ **更兼容**: 更好的防火墙和代理兼容性

### 版本升级优势

- 🚀 **FastMCP 2.10.6**: 最新版本，包含所有最新功能和修复
- 🤝 **MCP 1.12.2**: 支持最新的 MCP 协议规范
- 🔧 **更好的错误处理**: 改进的错误报告和调试信息
- 📊 **性能优化**: 更快的启动时间和更低的内存使用

## 认证配置

服务器支持通过 HTTP Headers 进行 JWT 认证：

```json
{
  "mcpServers": {
    "fastmcp-streamable-http": {
      "type": "streamableHttp",
      "url": "http://localhost:8000/mcp/",
      "headers": {
        "Content-Type": "application/json",
        "Authorization": "Bearer YOUR_JWT_TOKEN"
      }
    }
  }
}
```

JWT Token 包含用户信息：

- 用户ID
- 用户名
- 部门ID
- 权限信息

```curl
curl -s http://localhost:8081/api/v1/auth/login -H "Content-Type: application/json" -d '{"username":"lidi10","password":"123456"}'
```

```json
{"code":0,"message":"登录成功","data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo3LCJ1c2VybmFtZSI6ImxpZGkxMCIsImRlcHRfaWQiOjIsImV4cCI6MTc1MzU4ODcyNSwiaWF0IjoxNzUzNTAyMzI1LCJpc3MiOiJhcGktYXV0aC1zeXN0ZW0ifQ.edRZQJXPvq43ppcs6GdnG800voWh3UJAPY6lU4Ovprw","expires_at":1753588725,"user":{"id":7,"username":"lidi10","name":"李迪","email":"lidi10@126.com","avatar":"https://avatars.githubusercontent.com/u/1?v=4","dept_id":2,"status":1,"created_at":"2025-07-20T21:24:39.842+08:00","updated_at":"2025-07-20T21:24:39.842+08:00"}}}
```

## 故障排除

### 常见问题

1. **连接失败**

   - 确认服务器正在运行：`lsof -i :8000`
   - 检查端点路径是否正确：`/mcp/`
   - 验证 URL 格式：`http://localhost:8000/mcp/`
2. **工具不可用**

   - 检查服务器日志中的错误信息
   - 确认 FastMCP 版本为 2.x
   - 重启 Claude Desktop
3. **权限错误**

   - 检查认证令牌配置
   - 确认用户权限设置

### 调试技巧

- 查看服务器控制台输出
- 使用 curl 测试端点响应
- 检查 Claude Desktop 的 MCP 设置格式

## 兼容性说明

- **FastMCP 2.x**: 推荐使用 Streamable HTTP
- **向后兼容**: 仍支持 SSE 传输（已弃用）
- **Claude Desktop**: 需要支持 `streamableHttp` 类型
