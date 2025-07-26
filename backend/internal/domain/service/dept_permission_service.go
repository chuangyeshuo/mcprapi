package service

import (
	"errors"
	"fmt"

	"mcprapi/backend/internal/domain/entity"
	"mcprapi/backend/internal/domain/repository"
	"mcprapi/backend/pkg/casbinx"
)

// DeptPermissionService 部门权限管理服务
type DeptPermissionService struct {
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	deptRepo       repository.DepartmentRepository
	casbinEnforcer *casbinx.Enforcer
}

// NewDeptPermissionService 创建部门权限管理服务
func NewDeptPermissionService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	deptRepo repository.DepartmentRepository,
	casbinEnforcer *casbinx.Enforcer,
) *DeptPermissionService {
	return &DeptPermissionService{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		deptRepo:       deptRepo,
		casbinEnforcer: casbinEnforcer,
	}
}

// GrantDeptAdminRequest 授予部门管理员权限请求
type GrantDeptAdminRequest struct {
	UserID uint `json:"user_id" binding:"required"`
	DeptID uint `json:"dept_id" binding:"required"`
}

// RevokeDeptAdminRequest 撤销部门管理员权限请求
type RevokeDeptAdminRequest struct {
	UserID uint `json:"user_id" binding:"required"`
	DeptID uint `json:"dept_id" binding:"required"`
}

// CheckDeptPermissionRequest 检查部门权限请求
type CheckDeptPermissionRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	DeptID uint   `json:"dept_id" binding:"required"`
	Action string `json:"action" binding:"required"` // read, write, admin
}

// GrantDeptAdmin 授予用户部门管理员权限
func (s *DeptPermissionService) GrantDeptAdmin(req *GrantDeptAdminRequest) error {
	// 检查用户是否存在
	user, err := s.userRepo.GetByID(req.UserID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 检查部门是否存在
	dept, err := s.deptRepo.GetByID(req.DeptID)
	if err != nil {
		return errors.New("部门不存在")
	}

	// 检查用户是否属于该部门
	if user.DeptID != req.DeptID {
		return errors.New("用户不属于该部门")
	}

	// 创建部门管理员角色（如果不存在）
	deptAdminRoleCode := fmt.Sprintf("dept_admin_%d", req.DeptID)
	deptAdminRole, err := s.roleRepo.GetByCode(deptAdminRoleCode)
	if err != nil {
		// 角色不存在，创建新角色
		deptAdminRole = &entity.Role{
			Name:        fmt.Sprintf("%s部门管理员", dept.Name),
			Code:        deptAdminRoleCode,
			Description: fmt.Sprintf("%s部门的管理员角色", dept.Name),
			DeptID:      req.DeptID,
			Status:      1,
		}
		if err := s.roleRepo.Create(deptAdminRole); err != nil {
			return fmt.Errorf("创建部门管理员角色失败: %v", err)
		}
	}

	// 为用户分配部门管理员角色
	userRoleCode := fmt.Sprintf("user_%d", req.UserID)
	if _, err := s.casbinEnforcer.AddRoleForUser(userRoleCode, deptAdminRoleCode); err != nil {
		return fmt.Errorf("分配角色失败: %v", err)
	}

	// 添加部门管理员权限策略
	policies := [][]string{
		// 用户管理权限
		{deptAdminRoleCode, fmt.Sprintf("/api/v1/dept/%d/users", req.DeptID), "GET", fmt.Sprintf("%d", req.DeptID), "allow"},
		{deptAdminRoleCode, fmt.Sprintf("/api/v1/dept/%d/users", req.DeptID), "POST", fmt.Sprintf("%d", req.DeptID), "allow"},
		{deptAdminRoleCode, fmt.Sprintf("/api/v1/dept/%d/users/*", req.DeptID), "PUT", fmt.Sprintf("%d", req.DeptID), "allow"},
		{deptAdminRoleCode, fmt.Sprintf("/api/v1/dept/%d/users/*", req.DeptID), "DELETE", fmt.Sprintf("%d", req.DeptID), "allow"},

		// 角色管理权限
		{deptAdminRoleCode, fmt.Sprintf("/api/v1/dept/%d/roles", req.DeptID), "GET", fmt.Sprintf("%d", req.DeptID), "allow"},
		{deptAdminRoleCode, fmt.Sprintf("/api/v1/dept/%d/roles", req.DeptID), "POST", fmt.Sprintf("%d", req.DeptID), "allow"},
		{deptAdminRoleCode, fmt.Sprintf("/api/v1/dept/%d/roles/*", req.DeptID), "PUT", fmt.Sprintf("%d", req.DeptID), "allow"},
		{deptAdminRoleCode, fmt.Sprintf("/api/v1/dept/%d/roles/*", req.DeptID), "DELETE", fmt.Sprintf("%d", req.DeptID), "allow"},

		// 业务线管理权限
		{deptAdminRoleCode, fmt.Sprintf("/api/v1/dept/%d/businesses", req.DeptID), "GET", fmt.Sprintf("%d", req.DeptID), "allow"},
		{deptAdminRoleCode, fmt.Sprintf("/api/v1/dept/%d/businesses", req.DeptID), "POST", fmt.Sprintf("%d", req.DeptID), "allow"},
		{deptAdminRoleCode, fmt.Sprintf("/api/v1/dept/%d/businesses/*", req.DeptID), "PUT", fmt.Sprintf("%d", req.DeptID), "allow"},
		{deptAdminRoleCode, fmt.Sprintf("/api/v1/dept/%d/businesses/*", req.DeptID), "DELETE", fmt.Sprintf("%d", req.DeptID), "allow"},

		// API管理权限（通过业务线）
		{deptAdminRoleCode, "/api/v1/business/*/apis", "GET", fmt.Sprintf("%d", req.DeptID), "allow"},
		{deptAdminRoleCode, "/api/v1/business/*/apis", "POST", fmt.Sprintf("%d", req.DeptID), "allow"},
		{deptAdminRoleCode, "/api/v1/business/*/apis/*", "PUT", fmt.Sprintf("%d", req.DeptID), "allow"},
		{deptAdminRoleCode, "/api/v1/business/*/apis/*", "DELETE", fmt.Sprintf("%d", req.DeptID), "allow"},
	}

	for _, policy := range policies {
		if _, err := s.casbinEnforcer.AddPolicyWithDept(policy[0], policy[1], policy[2], policy[3], policy[4]); err != nil {
			return fmt.Errorf("添加权限策略失败: %v", err)
		}
	}

	// 同步策略到内存（安全方式）
	return s.casbinEnforcer.SyncPolicyToMemory()
}

// RevokeDeptAdmin 撤销用户部门管理员权限
func (s *DeptPermissionService) RevokeDeptAdmin(req *RevokeDeptAdminRequest) error {
	// 检查用户是否存在
	if _, err := s.userRepo.GetByID(req.UserID); err != nil {
		return errors.New("用户不存在")
	}

	// 检查部门是否存在
	if _, err := s.deptRepo.GetByID(req.DeptID); err != nil {
		return errors.New("部门不存在")
	}

	// 删除用户的部门管理员角色
	userRoleCode := fmt.Sprintf("user_%d", req.UserID)
	deptAdminRoleCode := fmt.Sprintf("dept_admin_%d", req.DeptID)

	if _, err := s.casbinEnforcer.DeleteRoleForUser(userRoleCode, deptAdminRoleCode); err != nil {
		return fmt.Errorf("撤销角色失败: %v", err)
	}

	// 同步策略到内存（安全方式）
	return s.casbinEnforcer.SyncPolicyToMemory()
}

// CheckDeptPermission 检查用户是否有部门权限
func (s *DeptPermissionService) CheckDeptPermission(req *CheckDeptPermissionRequest) (bool, error) {
	// 获取用户角色
	userRoleCode := fmt.Sprintf("user_%d", req.UserID)
	roles, err := s.casbinEnforcer.GetRolesForUser(userRoleCode)
	if err != nil {
		return false, err
	}

	// 检查是否有系统管理员权限
	for _, role := range roles {
		if role == "admin" {
			return true, nil // 系统管理员有所有权限
		}
	}

	// 检查是否有部门管理员权限
	deptAdminRoleCode := fmt.Sprintf("dept_admin_%d", req.DeptID)
	for _, role := range roles {
		if role == deptAdminRoleCode {
			return true, nil // 部门管理员有本部门权限
		}
	}

	// 检查具体的权限
	path := fmt.Sprintf("/api/v1/dept/%d/*", req.DeptID)
	for _, role := range roles {
		if s.casbinEnforcer.EnforceWithDept(role, path, req.Action, req.DeptID) {
			return true, nil
		}
	}

	return false, nil
}

// GetDeptAdmins 获取部门管理员列表
func (s *DeptPermissionService) GetDeptAdmins(deptID uint) ([]entity.User, error) {
	// 检查部门是否存在
	if _, err := s.deptRepo.GetByID(deptID); err != nil {
		return nil, errors.New("部门不存在")
	}

	// 获取部门管理员角色的用户
	deptAdminRoleCode := fmt.Sprintf("dept_admin_%d", deptID)
	userCodes, err := s.casbinEnforcer.GetUsersForRole(deptAdminRoleCode)
	if err != nil {
		return nil, err
	}

	var users []entity.User
	for _, userCode := range userCodes {
		// 从用户代码中提取用户ID
		var userID uint
		if _, err := fmt.Sscanf(userCode, "user_%d", &userID); err != nil {
			continue
		}

		user, err := s.userRepo.GetByID(userID)
		if err != nil {
			continue
		}

		users = append(users, *user)
	}

	return users, nil
}

// InitSystemAdmin 初始化系统管理员权限
func (s *DeptPermissionService) InitSystemAdmin(userID uint) error {
	// 创建系统管理员角色（如果不存在）
	systemAdminRole, err := s.roleRepo.GetByCode("admin")
	if err != nil {
		// 角色不存在，创建新角色
		systemAdminRole = &entity.Role{
			Name:        "系统管理员",
			Code:        "admin",
			Description: "系统超级管理员，拥有所有权限",
			DeptID:      0, // 0表示不属于任何部门
			Status:      1,
		}
		if err := s.roleRepo.Create(systemAdminRole); err != nil {
			return fmt.Errorf("创建系统管理员角色失败: %v", err)
		}
	}

	// 为用户分配系统管理员角色
	userRoleCode := fmt.Sprintf("user_%d", userID)
	if _, err := s.casbinEnforcer.AddRoleForUser(userRoleCode, "admin"); err != nil {
		return fmt.Errorf("分配系统管理员角色失败: %v", err)
	}

	// 添加系统管理员权限策略（拥有所有权限）
	if _, err := s.casbinEnforcer.AddPolicyWithDept("admin", "/api/v1/*", "*", "*", "allow"); err != nil {
		return fmt.Errorf("添加系统管理员权限策略失败: %v", err)
	}

	// 同步策略到内存（安全方式）
	return s.casbinEnforcer.SyncPolicyToMemory()
}
