/*
 * @Author: lidi10@staff.weibo.com
 * @Date: 2025-07-19 17:26:54
 * @LastEditTime: 2025-07-26 13:30:45
 * @LastEditors: lidi10@staff.weibo.com
 * @Description:
 * Copyright (c) 2023 by Weibo, All Rights Reserved.
 */
package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"mcprapi/backend/internal/domain/entity"
	"mcprapi/backend/internal/domain/repository"
	"mcprapi/backend/internal/pkg/encrypt"
	"mcprapi/backend/pkg/casbinx"
)

// AuthService 认证服务
type AuthService struct {
	userRepo  repository.UserRepository
	roleRepo  repository.RoleRepository
	apiRepo   repository.APIRepository
	casbin    *casbinx.Enforcer
	jwtSecret string
}

// NewAuthService 创建认证服务
func NewAuthService(userRepo repository.UserRepository, roleRepo repository.RoleRepository, apiRepo repository.APIRepository, casbin *casbinx.Enforcer, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		roleRepo:  roleRepo,
		apiRepo:   apiRepo,
		casbin:    casbin,
		jwtSecret: jwtSecret,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string       `json:"token"`
	ExpiresAt int64        `json:"expires_at"`
	User      *entity.User `json:"user"`
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// 验证用户名和密码不能为空
	if req.Username == "" {
		return nil, errors.New("用户名不能为空")
	}
	if req.Password == "" {
		return nil, errors.New("密码不能为空")
	}

	// 根据用户名查找用户
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}
	if user == nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 验证密码
	if !encrypt.VerifyPassword(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("用户已禁用")
	}

	// 生成JWT令牌
	token, expiresAt, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user,
	}, nil
}

// Claims JWT声明
type Claims struct {
	UserID       uint   `json:"user_id"`
	Username     string `json:"username"`
	DeptID       uint   `json:"dept_id"`
	TokenVersion int    `json:"v,omitempty"` // Token版本号，用于安全控制
	jwt.StandardClaims
}

// 生成JWT令牌
func (s *AuthService) generateToken(user *entity.User) (string, int64, error) {
	expiresAt := time.Now().Add(time.Hour * 24).Unix() // 24小时过期

	claims := &Claims{
		UserID:       user.ID,
		Username:     user.Username,
		DeptID:       user.DeptID,
		TokenVersion: user.TokenVersion, // 包含用户的Token版本号
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
			IssuedAt:  time.Now().Unix(),
			Issuer:    "api-auth-system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt, nil
}

// GenerateUserToken 为用户生成Token（默认24小时过期）
func (s *AuthService) GenerateUserToken(user *entity.User) (string, int64, error) {
	return s.generateToken(user)
}

// GenerateUserTokenWithExpiry 为用户生成指定过期时间的Token
func (s *AuthService) GenerateUserTokenWithExpiry(user *entity.User, expireDays int) (string, int64, error) {
	expiresAt := time.Now().Add(time.Duration(expireDays) * 24 * time.Hour).Unix()

	claims := &Claims{
		UserID:       user.ID,
		Username:     user.Username,
		DeptID:       user.DeptID,
		TokenVersion: user.TokenVersion, // 包含用户的Token版本号
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
			IssuedAt:  time.Now().Unix(),
			Issuer:    "api-auth-system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt, nil
}

// GenerateUserTokenWithVersionIncrement 生成Token并递增版本号
func (s *AuthService) GenerateUserTokenWithVersionIncrement(user *entity.User, expireDays int) (string, int64, error) {
	// 递增用户的Token版本号
	user.TokenVersion++
	err := s.userRepo.Update(user)
	if err != nil {
		return "", 0, fmt.Errorf("更新用户Token版本号失败: %v", err)
	}

	// 生成新的Token
	return s.GenerateUserTokenWithExpiry(user, expireDays)
}

// ValidateTokenVersion 验证Token版本号
func (s *AuthService) ValidateTokenVersion(userID uint, tokenVersion int) error {
	// 如果Token中没有版本号，则跳过验证（兼容旧Token）
	if tokenVersion == 0 {
		return nil
	}

	// 获取用户当前的版本号
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("获取用户信息失败: %v", err)
	}

	// 比较版本号
	if tokenVersion < user.TokenVersion {
		return errors.New("Token已失效，请重新获取")
	}

	return nil
}

// CheckPermissionRequest 权限检查请求
type CheckPermissionRequest struct {
	UserID  string `json:"user_id"`
	APIPath string `json:"api_path"`
	Method  string `json:"method"`
}

// CheckPermissionResponse 权限检查响应
type CheckPermissionResponse struct {
	Allowed bool     `json:"allowed"`
	Roles   []string `json:"roles"`
	IsAdmin bool     `json:"is_admin"`
	DeptID  uint     `json:"dept_id"`
}

// CheckPermission 检查用户是否有权限访问API
func (s *AuthService) CheckPermission(req *CheckPermissionRequest) (*CheckPermissionResponse, error) {
	// 获取用户角色
	userID := uint(0)
	if req.UserID != "" {
		// 转换用户ID为uint
		_, err := fmt.Sscanf(req.UserID, "%d", &userID)
		if err != nil {
			return nil, errors.New("无效的用户ID")
		}
	}

	roles, err := s.userRepo.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	// 获取用户信息以获取部门ID
	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return &CheckPermissionResponse{
			Allowed: false,
			Roles:   []string{},
			IsAdmin: false,
			DeptID:  0,
		}, nil
	}

	// 检查是否为管理员
	isAdmin := false
	for _, role := range roles {
		if role == "admin" {
			isAdmin = true
			break
		}
	}

	// 特殊处理：对于dashboard、business、api这3个路径组，只要是本部门登录用户即可操作
	if s.isDepartmentLevelAPI(req.APIPath) {
		// 只要用户已登录且有有效的部门ID，就允许访问
		if user.DeptID > 0 {
			return &CheckPermissionResponse{
				Allowed: true,
				Roles:   roles,
				IsAdmin: isAdmin,
				DeptID:  user.DeptID,
			}, nil
		}
	}

	// 检查权限
	allowed := false
	for _, role := range roles {
		// 使用EnforceWithDept方法，包括部门信息
		if s.casbin.EnforceWithDept(role, req.APIPath, req.Method, user.DeptID) {
			allowed = true
			break
		}
		// 也尝试使用通配符部门（部门ID为0时会使用"*"）
		if s.casbin.EnforceWithDept(role, req.APIPath, req.Method, 0) {
			allowed = true
			break
		}
	}

	return &CheckPermissionResponse{
		Allowed: allowed,
		Roles:   roles,
		IsAdmin: isAdmin,
		DeptID:  user.DeptID,
	}, nil
}

// isDepartmentLevelAPI 检查是否为部门级别的API路径
func (s *AuthService) isDepartmentLevelAPI(apiPath string) bool {
	// 检查是否为dashboard、business、api这3个路径组
	if strings.HasPrefix(apiPath, "/api/v1/dashboard/") ||
		strings.HasPrefix(apiPath, "/api/v1/business/") ||
		strings.HasPrefix(apiPath, "/api/v1/api/") {
		return true
	}
	return false
}

// LogoutResponse 登出响应
type LogoutResponse struct {
	Message string `json:"message"`
}

// Logout 用户登出
func (s *AuthService) Logout() (*LogoutResponse, error) {
	// 在JWT无状态的情况下，服务端不需要做特殊处理
	// 客户端只需要删除本地存储的token即可
	// 如果需要实现token黑名单功能，可以在这里将token加入黑名单

	return &LogoutResponse{
		Message: "登出成功",
	}, nil
}

// GenerateQRCode 生成扫码登录二维码
func (s *AuthService) GenerateQRCode() (string, error) {
	// 生成唯一的二维码ID
	qrID := uuid.New().String()

	// 在实际应用中，这里应该将二维码ID存储到Redis中，设置过期时间
	// 并关联一个状态（等待扫码、已扫码、已确认等）

	// 返回二维码内容，前端可以使用这个内容生成二维码图片
	return qrID, nil
}
