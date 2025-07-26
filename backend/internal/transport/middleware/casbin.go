/*
 * @Author: lidi10@staff.weibo.com
 * @Date: 2025-07-19 17:31:36
 * @LastEditTime: 2025-07-20 08:48:41
 * @LastEditors: lidi10@staff.weibo.com
 * @Description:
 * Copyright (c) 2023 by Weibo, All Rights Reserved.
 */
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mcprapi/backend/internal/domain/repository"
	"mcprapi/backend/internal/transport/http/dto"
	"mcprapi/backend/pkg/casbinx"
)

// Casbin Casbin中间件
func Casbin(enforcer *casbinx.Enforcer, userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户角色
		userID := GetCurrentUser(c)
		if userID == 0 {
			c.JSON(http.StatusUnauthorized, dto.Response{
				Code:    dto.CodeUnauthorized,
				Message: "未认证的用户",
			})
			c.Abort()
			return
		}

		// 获取当前用户部门ID
		userDeptID := GetCurrentDeptID(c)

		// 获取请求路径和方法
		path := c.Request.URL.Path
		method := c.Request.Method

		// 获取用户角色
		roles, err := getUserRoles(userID, userRepo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.Response{
				Code:    dto.CodeInternalError,
				Message: "获取用户角色失败",
			})
			c.Abort()
			return
		}

		// 如果用户没有角色，拒绝访问
		if len(roles) == 0 {
			c.JSON(http.StatusForbidden, dto.Response{
				Code:    dto.CodeForbidden,
				Message: "无权访问该资源",
			})
			c.Abort()
			return
		}

		// 检查权限 - 增加部门维度检查
		allowed := false
		for _, role := range roles {
			// 使用新的4参数Enforce方法：role, path, method, deptID
			if enforcer.EnforceWithDept(role, path, method, userDeptID) {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, dto.Response{
				Code:    dto.CodeForbidden,
				Message: "无权访问该资源",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// 获取用户角色
func getUserRoles(userID uint, userRepo repository.UserRepository) ([]string, error) {
	// 从数据库中获取用户角色
	roles, err := userRepo.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	// 如果用户没有分配角色，返回空数组
	if len(roles) == 0 {
		return []string{}, nil
	}

	return roles, nil
}
