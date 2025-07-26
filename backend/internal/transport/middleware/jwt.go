package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"mcprapi/backend/internal/transport/http/dto"
)

// JWTClaims JWT声明
type JWTClaims struct {
	UserID       uint   `json:"user_id"`
	Username     string `json:"username"`
	DeptID       uint   `json:"dept_id"`
	TokenVersion int    `json:"v,omitempty"` // Token版本号，用于安全控制
	jwt.StandardClaims
}

// JWT JWT中间件
func JWT(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, dto.Response{
				Code:    dto.CodeUnauthorized,
				Message: "未提供认证令牌",
			})
			c.Abort()
			return
		}

		// 移除Bearer前缀
		token = strings.TrimPrefix(token, "Bearer ")

		// 解析JWT令牌
		claims, err := parseToken(token, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.Response{
				Code:    dto.CodeUnauthorized,
				Message: "无效的认证令牌: " + err.Error(),
			})
			c.Abort()
			return
		}

		// 检查令牌是否过期
		if claims.ExpiresAt < time.Now().Unix() {
			c.JSON(http.StatusUnauthorized, dto.Response{
				Code:    dto.CodeUnauthorized,
				Message: "认证令牌已过期",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("dept_id", claims.DeptID)
		c.Set("token_version", claims.TokenVersion) // 存储Token版本号

		c.Next()
	}
}

// 解析JWT令牌
func parseToken(tokenString, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的令牌")
}

// GetCurrentUser 获取当前用户ID
func GetCurrentUser(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}

	return userID.(uint)
}

// GetCurrentUsername 获取当前用户名
func GetCurrentUsername(c *gin.Context) string {
	username, exists := c.Get("username")
	if !exists {
		return ""
	}

	return username.(string)
}

// GetCurrentDeptID 获取当前用户部门ID
func GetCurrentDeptID(c *gin.Context) uint {
	deptID, exists := c.Get("dept_id")
	if !exists {
		return 0
	}

	return deptID.(uint)
}

// GetCurrentTokenVersion 获取当前Token版本号
func GetCurrentTokenVersion(c *gin.Context) int {
	tokenVersion, exists := c.Get("token_version")
	if !exists {
		return 0
	}

	return tokenVersion.(int)
}