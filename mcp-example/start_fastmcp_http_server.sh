#!/bin/bash

# 启动 FastMCP HTTP 服务器
echo "🚀 启动 FastMCP HTTP/SSE 服务器..."

# 设置环境变量
export REAL_WEATHER_API_URL="https://api.weather.gov"
export REAL_BUSINESS_API_URL="http://localhost:8081"
export REAL_AUTH_API_URL="http://localhost:8081"
export SERVER_SECRET="your_server_secret_key"

# 启动 MCP 服务器 (SSE 模式，支持 HTTP)
echo "📡 服务器将在 http://localhost:3000 上运行"
echo "🔗 可以通过 Claude Desktop 配置文件连接"
python fastmcp_http_server.py