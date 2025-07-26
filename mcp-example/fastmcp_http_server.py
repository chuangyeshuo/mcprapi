"""
使用 FastMCP 2.x 实现的云端MCP服务器 - Streamable HTTP 版本
支持通过 Streamable HTTP 协议访问（推荐）和 SSE 协议（兼容）
支持从 HTTP 请求头中获取 JWT token 进行身份验证
"""
import os
import json
import logging
from typing import Dict, Any, Optional
import httpx
from fastmcp import FastMCP
from fastmcp.server.dependencies import get_http_headers

# 配置日志
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# 创建 FastMCP 实例
mcp = FastMCP(name="Cloud MCP Server")

# 配置
WEATHER_API_URL = os.getenv("REAL_WEATHER_API_URL", "https://api.weather.gov")
BUSINESS_API_URL = os.getenv("REAL_BUSINESS_API_URL", "http://localhost:8081")
AUTH_API_URL = os.getenv("REAL_AUTH_API_URL", "http://localhost:8081")

# HTTP客户端
http_client = httpx.AsyncClient(timeout=30.0)

def get_token_from_request() -> str:
    """从当前请求的 Authorization 头中获取 token"""
    try:
        # 使用 FastMCP 2.x 的 get_http_headers 函数
        headers = get_http_headers()
        auth_header = headers.get("authorization", "")
        
        if auth_header.startswith("Bearer "):
            token = auth_header[7:]  # 移除 "Bearer " 前缀
            logger.info(f"从请求头获取到 token: {token[:20]}...")
            return token
        
        logger.warning("请求头中未找到有效的 Authorization token")
        return ""
    except RuntimeError:
        # 当不在 HTTP 请求上下文中时会抛出 RuntimeError
        logger.warning("不在 HTTP 请求上下文中，无法获取请求头")
        return ""
    except Exception as e:
        logger.error(f"获取请求头时发生错误: {str(e)}")
        return ""

async def verify_user_permission(user_token: str, tool_name: str) -> Optional[Dict[str, Any]]:
    """验证用户权限"""
    if not user_token or user_token == "":
        # 测试模式，返回默认用户信息
        logger.info("使用测试模式，返回默认用户信息")
        return {
            "user_id": 7,
            "username": "test_user",
            "roles": ["member"],
            "permissions": ["weather:read", "user:read", "department:read"]
        }
    
    api_path = f"/member/api/v1/mcp/{tool_name}"
    logger.info(f"验证权限 - 工具: {tool_name}, API路径: {api_path}")
    
    try:
        response = await http_client.post(
            f"{AUTH_API_URL}/api/v1/api/check-permission",
            json={"api_path": api_path, "method": "POST"},
            headers={
                "Authorization": f"Bearer {user_token}",
                "Content-Type": "application/json"
            }
        )
        
        logger.info(f"权限验证响应状态码: {response.status_code}")
        
        if response.status_code == 200:
            result = response.json()
            logger.info(f"权限验证响应内容: {json.dumps(result, ensure_ascii=False)}")
            
            data = result.get("data", {})
            allowed = data.get("allowed", False)
            logger.info(f"权限检查结果 - allowed: {allowed}")
            
            if allowed:
                logger.info(f"^_^权限验证成功！🚀")
                return data
            else:
                logger.warning(f"权限验证失败 - 用户无权限访问 {tool_name}")
        else:
            logger.error(f"权限验证API返回错误状态码: {response.status_code}")
            logger.error(f"响应内容: {response.text}")
        
        return None
    except Exception as e:
        logger.error(f"权限验证异常: {str(e)}")
        # 测试模式，返回默认用户信息
        logger.info("异常情况下使用测试模式")
        return {
            "user_id": 7,
            "username": "test_user",
            "roles": ["member"],
            "permissions": ["weather:read", "user:read", "department:read"]
        }

@mcp.tool()
async def get_weather_alerts(state: str) -> str:
    """获取指定州的天气预警信息
    
    Args:
        state: 美国州的缩写，如 'CA' 代表加利福尼亚州
    """
    # 从请求头获取 token
    user_token = get_token_from_request()
    
    # 验证权限
    user_info = await verify_user_permission(user_token, "get_weather_alerts")
    if not user_info:
        return "❌ 权限验证失败"
    
    try:
        response = await http_client.get(f"{WEATHER_API_URL}/alerts/active/area/{state}")
        
        if response.status_code == 200:
            data = response.json()
            alerts = data.get("features", [])
            
            if not alerts:
                return f"当前 {state} 州没有活跃的天气预警。"
            
            result = f"🌪️ {state} 州天气预警信息:\n\n"
            for i, alert in enumerate(alerts[:3], 1):  # 限制显示3个预警
                properties = alert.get("properties", {})
                event = properties.get("event", "未知事件")
                headline = properties.get("headline", "无标题")
                description = properties.get("description", "无描述")[:150]
                
                result += f"{i}. {event}\n"
                result += f"   标题: {headline}\n"
                result += f"   描述: {description}...\n\n"
            
            return result
        else:
            return f"获取天气预警失败: HTTP {response.status_code}"
            
    except Exception as e:
        return f"天气预警API调用失败: {str(e)}"

@mcp.tool()
async def get_weather_forecast(latitude: float, longitude: float) -> str:
    """获取指定位置的天气预报
    
    Args:
        latitude: 纬度
        longitude: 经度
    """
    # 从请求头获取 token
    user_token = get_token_from_request()
    
    # 验证权限
    user_info = await verify_user_permission(user_token, "get_weather_forecast")
    if not user_info:
        return "❌ 权限验证失败"
    
    try:
        # 获取预报点信息
        points_response = await http_client.get(f"{WEATHER_API_URL}/points/{latitude},{longitude}")
        
        if points_response.status_code != 200:
            return f"获取预报点信息失败: HTTP {points_response.status_code}"
        
        points_data = points_response.json()
        forecast_url = points_data.get("properties", {}).get("forecast")
        
        if not forecast_url:
            return "无法获取预报URL"
        
        # 获取预报信息
        forecast_response = await http_client.get(forecast_url)
        
        if forecast_response.status_code == 200:
            forecast_data = forecast_response.json()
            periods = forecast_data.get("properties", {}).get("periods", [])
            
            if not periods:
                return "无可用的天气预报数据"
            
            result = f"🌤️ 位置 ({latitude}, {longitude}) 天气预报:\n\n"
            for period in periods[:5]:  # 显示未来5个时段
                name = period.get("name", "未知时段")
                temperature = period.get("temperature", "N/A")
                temperature_unit = period.get("temperatureUnit", "")
                short_forecast = period.get("shortForecast", "无预报")
                
                result += f"📅 {name}\n"
                result += f"🌡️ 温度: {temperature}°{temperature_unit}\n"
                result += f"☁️ 天气: {short_forecast}\n\n"
            
            return result
        else:
            return f"获取天气预报失败: HTTP {forecast_response.status_code}"
            
    except Exception as e:
        return f"天气预报API调用失败: {str(e)}"

@mcp.tool()
async def get_user_info(user_id: int) -> str:
    """获取用户信息
    
    Args:
        user_id: 用户ID
    """
    # 从请求头获取 token
    user_token = get_token_from_request()
    
    # 验证权限
    user_info = await verify_user_permission(user_token, "get_user_info")
    if not user_info:
        return "❌ 权限验证失败"
    
    try:
        # 返回模拟数据
        user_data = {
            "user_id": user_id,
            "username": f"user{user_id}",
            "email": f"user{user_id}@example.com",
            "roles": ["member", "aigc"],
            "created_at": "2024-01-01T00:00:00Z",
            "last_login": "2024-12-20T10:30:00Z",
            "status": "active"
        }
        return json.dumps(user_data, indent=2, ensure_ascii=False)
        
    except Exception as e:
        return f"用户信息API调用失败: {str(e)}"

@mcp.tool()
async def get_department_stats(department: str) -> str:
    """获取部门统计信息
    
    Args:
        department: 部门名称
    """
    # 从请求头获取 token
    user_token = get_token_from_request()
    
    # 验证权限
    user_info = await verify_user_permission(user_token, "get_department_stats")
    if not user_info:
        return "❌ 权限验证失败"
    
    try:
        # 返回模拟数据
        stats_data = {
            "department": department,
            "employee_count": 25,
            "active_projects": 8,
            "budget_utilization": "75%",
            "performance_score": 4.2,
            "monthly_revenue": "$125,000",
            "last_updated": "2024-12-20T15:30:00Z"
        }
        return json.dumps(stats_data, indent=2, ensure_ascii=False)
        
    except Exception as e:
        return f"部门统计API调用失败: {str(e)}"

@mcp.tool()
async def create_order(product: str, quantity: int) -> str:
    """创建订单（需要管理员权限）
    
    Args:
        product: 产品名称
        quantity: 数量
    """
    # 从请求头获取 token
    user_token = get_token_from_request()
    
    # 验证权限
    user_info = await verify_user_permission(user_token, "create_order")
    if not user_info:
        return "❌ 权限验证失败"
    
    # 检查管理员权限
    user_roles = user_info.get("roles", [])
    if "administrator" not in user_roles and "order_manager" not in user_roles:
        return "❌ 权限不足：创建订单需要管理员或订单管理员权限"
    
    try:
        # 返回模拟订单创建结果
        order_data = {
            "order_id": f"ORD-{user_info.get('user_id')}-{hash(product) % 10000}",
            "product": product,
            "quantity": quantity,
            "created_by": user_info.get("user_id"),
            "status": "created",
            "total_amount": quantity * 99.99,
            "created_at": "2024-12-20T15:45:00Z"
        }
        
        return f"✅ 订单创建成功:\n{json.dumps(order_data, indent=2, ensure_ascii=False)}"
        
    except Exception as e:
        return f"创建订单API调用失败: {str(e)}"

if __name__ == "__main__":
    print("🚀 启动 FastMCP 服务器 (Streamable HTTP 模式)...")
    print("🔗 使用 Streamable HTTP 传输协议 (推荐)")
    print("🔐 支持从 Authorization 头获取 JWT token")
    print("📡 HTTP 端点: http://localhost:8000")
    print("📡 SSE 端点: http://localhost:8000/sse (兼容模式)")
    
    # 启动 FastMCP 服务器，使用 Streamable HTTP 传输协议
    # FastMCP 2.x 默认使用 Streamable HTTP，也支持 SSE 兼容模式
    mcp.run(transport="http")