package service

import (
	"errors"
	"fmt"
	"strconv"

	"mcprapi/backend/internal/domain/entity"
	"mcprapi/backend/internal/domain/repository"
	"mcprapi/backend/pkg/casbinx"
)

// RoleService 角色服务
type RoleService struct {
	roleRepo       repository.RoleRepository
	deptRepo       repository.DepartmentRepository
	apiRepo        repository.APIRepository
	userRepo       repository.UserRepository
	casbinEnforcer *casbinx.Enforcer
}

// NewRoleService 创建角色服务
func NewRoleService(roleRepo repository.RoleRepository, deptRepo repository.DepartmentRepository, apiRepo repository.APIRepository, userRepo repository.UserRepository, casbinEnforcer *casbinx.Enforcer) *RoleService {
	return &RoleService{
		roleRepo:       roleRepo,
		deptRepo:       deptRepo,
		apiRepo:        apiRepo,
		userRepo:       userRepo,
		casbinEnforcer: casbinEnforcer,
	}
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name           string `json:"name" binding:"required"`
	Code           string `json:"code" binding:"required"`
	Description    string `json:"description"`
	DeptID         uint   `json:"dept_id" binding:"required"`
	CurrentUserID  uint   `json:"-"` // 当前用户ID，用于权限验证，不从JSON绑定
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	ID          uint   `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
	DeptID      uint   `json:"dept_id" binding:"required"`
	Status      int    `json:"status" binding:"oneof=0 1"`
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(req *CreateRoleRequest) (*entity.Role, error) {
	// 检查部门是否存在
	dept, err := s.deptRepo.GetByID(req.DeptID)
	if err != nil {
		return nil, errors.New("部门不存在")
	}
	if dept == nil {
		return nil, errors.New("部门不存在")
	}

	// 权限验证：非admin用户只能在自己部门创建角色
	if req.CurrentUserID > 0 {
		currentUser, err := s.userRepo.GetByID(req.CurrentUserID)
		if err != nil {
			return nil, errors.New("当前用户不存在")
		}

		// 获取用户角色
		userRoles, err := s.userRepo.GetUserRoles(req.CurrentUserID)
		if err != nil {
			return nil, errors.New("获取用户角色失败")
		}

		// 检查是否是admin
		isAdmin := false
		for _, roleCode := range userRoles {
			if roleCode == "admin" || roleCode == "system_admin" {
				isAdmin = true
				break
			}
		}

		// 非admin用户只能在自己部门创建角色
		if !isAdmin && currentUser.DeptID != req.DeptID {
			return nil, errors.New("无权限在该部门创建角色")
		}
	}

	// 生成带部门编码的角色编码：部门编码_角色编码
	roleCodeWithDept := fmt.Sprintf("%s_%s", dept.Code, req.Code)

	// 检查角色名称在该部门下是否已存在
	existingRoles, err := s.roleRepo.ListByDept(req.DeptID)
	if err != nil {
		return nil, err
	}
	for _, existingRole := range existingRoles {
		if existingRole.Name == req.Name {
			return nil, errors.New("该部门下角色名称已存在，请使用其他名称")
		}
	}

	// 检查生成的角色编码是否已存在（全局检查）
	existingRoleByCode, _ := s.roleRepo.GetByCode(roleCodeWithDept)
	if existingRoleByCode != nil {
		return nil, errors.New("角色编码已存在，请使用其他编码")
	}

	// 创建角色
	role := &entity.Role{
		Name:        req.Name,
		Code:        roleCodeWithDept, // 使用拼接了部门编码的角色编码
		Description: req.Description,
		DeptID:      dept.ID,
		Status:      1, // 默认启用
	}

	err = s.roleRepo.Create(role)
	if err != nil {
		return nil, err
	}

	// 为新创建的角色添加默认权限（不使用SavePolicy避免清空策略表）
	err = s.addDefaultPermissionsToRoleWithoutSave(role.Code, role.DeptID)
	if err != nil {
		// 如果添加权限失败，记录错误但不回滚角色创建
		// 可以考虑添加日志记录
		fmt.Printf("Warning: Failed to add default permissions to role %s: %v\n", role.Code, err)
	}

	return role, nil
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(req *UpdateRoleRequest) (*entity.Role, error) {
	// 检查角色是否存在
	role, err := s.roleRepo.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("角色不存在")
	}

	// 检查部门是否存在
	dept, err := s.deptRepo.GetByID(req.DeptID)
	if err != nil {
		return nil, errors.New("部门不存在")
	}
	if dept == nil {
		return nil, errors.New("部门不存在")
	}

	// 编辑时直接使用用户输入的角色编码，不再自动拼接部门编码
	roleCode := req.Code

	// 检查角色编码是否已存在（排除自身）
	existingRole, _ := s.roleRepo.GetByCode(roleCode)
	if existingRole != nil && existingRole.ID != req.ID {
		return nil, errors.New("角色编码已存在，请使用其他编码")
	}

	// 检查角色名称在该部门下是否已存在（排除自身）
	existingRoles, err := s.roleRepo.ListByDept(req.DeptID)
	if err != nil {
		return nil, err
	}
	for _, existingRole := range existingRoles {
		if existingRole.Name == req.Name && existingRole.ID != req.ID {
			return nil, errors.New("该部门下角色名称已存在，请使用其他名称")
		}
	}

	// 更新角色
	role.Name = req.Name
	role.Code = roleCode // 直接使用用户输入的角色编码
	role.Description = req.Description
	role.DeptID = req.DeptID
	role.Status = req.Status

	err = s.roleRepo.Update(role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(id uint) error {
	// 检查角色是否存在
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("角色不存在")
	}

	// 删除角色
	return s.roleRepo.Delete(id)
}

// GetRole 获取角色
func (s *RoleService) GetRole(id uint) (*entity.Role, error) {
	return s.roleRepo.GetByID(id)
}

// ListRoleRequest 获取角色列表请求
type ListRoleRequest struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
	Query    string `json:"query" form:"query"`
	DeptID   uint   `json:"dept_id" form:"dept_id"`
}

// RoleWithDepartment 包含部门信息的角色
type RoleWithDepartment struct {
	*entity.Role
	DeptName string `json:"dept_name"`
}

// ListRoleResponse 获取角色列表响应
type ListRoleResponse struct {
	Total int64                 `json:"total"`
	Items []*RoleWithDepartment `json:"items"`
}

// ListRole 获取角色列表
func (s *RoleService) ListRole(req *ListRoleRequest) (*ListRoleResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	roles, total, err := s.roleRepo.List(req.Page, req.PageSize, req.Query, req.DeptID)
	if err != nil {
		return nil, err
	}

	// 为每个角色获取部门名称
	rolesWithDept := make([]*RoleWithDepartment, len(roles))
	for i, role := range roles {
		deptName := "未知部门"
		if role.DeptID > 0 {
			dept, err := s.deptRepo.GetByID(role.DeptID)
			if err == nil && dept != nil {
				deptName = dept.Name
			}
		}

		rolesWithDept[i] = &RoleWithDepartment{
			Role:     role,
			DeptName: deptName,
		}
	}

	return &ListRoleResponse{
		Total: total,
		Items: rolesWithDept,
	}, nil
}

// Permission 权限
type Permission struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

// UpdateRolePermissionsRequest 更新角色权限请求
type UpdateRolePermissionsRequest struct {
	RoleID      uint         `json:"role_id" binding:"required"`
	Permissions []Permission `json:"permissions" binding:"required"`
}

// UpdateRolePermissionsByAPIIDsRequest 通过API ID更新角色权限请求
type UpdateRolePermissionsByAPIIDsRequest struct {
	RoleID uint   `json:"role_id" binding:"required"`
	APIIDs []uint `json:"api_ids" binding:"required"`
	DeptID uint   `json:"dept_id,omitempty"` // 部门ID，用于权限范围控制
}

// UpdateRolePermissions 更新角色权限
func (s *RoleService) UpdateRolePermissions(req *UpdateRolePermissionsRequest) error {
	// 检查角色是否存在
	role, err := s.roleRepo.GetByID(req.RoleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("角色不存在")
	}

	// 选择性清除角色权限，保留指定的API权限
	s.deleteRolePermissionsExcept(role.Code, []string{
		"/api/v1/apis*",
		"/api/v1/apis/*",
		"/api/v1/business/*",
		"/api/v1/business*",
		"/api/v1/departments*",
		"/api/v1/departments/*",
		"/api/v1/role",
		"/api/v1/role/*",
		"/api/v1/user",
		"/api/v1/user/*",
	})

	// 添加新权限
	for _, perm := range req.Permissions {
		_, err = s.casbinEnforcer.AddPolicyWithDept(role.Code, perm.Path, perm.Method, "*", "allow")
		if err != nil {
			return err
		}
	}

	// 同步策略到内存（安全方式）
	return s.casbinEnforcer.SyncPolicyToMemory()
}

// addDefaultPermissionsToRoleWithoutSave 为角色添加默认权限（不保存策略，使用安全同步）
func (s *RoleService) addDefaultPermissionsToRoleWithoutSave(roleCode string, deptID uint) error {
	// 定义默认权限列表
	defaultPermissions := []struct {
		path   string
		method string
	}{
		{"/api/v1/apis*", "*"},
		{"/api/v1/apis/*", "*"},
		{"/api/v1/business/*", "*"},
		{"/api/v1/business*", "*"},
		{"/api/v1/departments*", "GET"},
		{"/api/v1/departments/*", "GET"},
		{"/api/v1/role", "*"},
		{"/api/v1/role/*", "*"},
		{"/api/v1/user", "*"},
		{"/api/v1/user/*", "*"},
	}

	// 将部门ID转换为字符串
	deptIDStr := fmt.Sprintf("%d", deptID)

	// 为角色添加每个默认权限
	for _, perm := range defaultPermissions {
		// 使用选择的部门ID而不是"*"
		_, err := s.casbinEnforcer.AddPolicyWithDeptAndSync(roleCode, perm.path, perm.method, deptIDStr, "allow")
		if err != nil {
			return fmt.Errorf("failed to add permission %s %s: %v", perm.path, perm.method, err)
		}
	}

	// 同步策略到内存（安全方式）
	return s.casbinEnforcer.SyncPolicyToMemory()
}

// GetUserAccessibleRoles 获取用户可访问的角色
func (s *RoleService) GetUserAccessibleRoles(userID uint) ([]*entity.Role, error) {
	// 检查用户是否存在
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 获取用户角色
	userRoles, err := s.userRepo.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	// 检查用户是否是admin
	isAdmin := false
	for _, roleCode := range userRoles {
		if roleCode == "admin" {
			isAdmin = true
			break
		}
	}

	// 如果是admin，返回所有角色
	if isAdmin {
		return s.roleRepo.ListAll()
	}

	// 非admin用户，只返回自己部门的角色
	return s.roleRepo.ListByDept(user.DeptID)
}

// deleteRolePermissionsExcept 删除角色权限，但保留指定的路径
func (s *RoleService) deleteRolePermissionsExcept(roleCode string, preservePaths []string) {
	// 获取角色的所有权限
	policies := s.casbinEnforcer.GetFilteredPolicy(0, roleCode)

	// 删除不在保留列表中的权限
	for _, policy := range policies {
		if len(policy) >= 3 {
			path := policy[1]
			method := policy[2]

			// 检查是否在保留列表中
			shouldPreserve := false
			for _, preservePath := range preservePaths {
				if path == preservePath {
					shouldPreserve = true
					break
				}
			}

			// 如果不在保留列表中，则删除
			if !shouldPreserve {
				if len(policy) >= 5 {
					// 包含部门信息的权限
					s.casbinEnforcer.RemovePolicyWithDept(roleCode, path, method, policy[3], policy[4])
				} else if len(policy) >= 4 {
					// 不包含部门信息的权限
					s.casbinEnforcer.RemovePolicy(roleCode, path, method)
				}
			}
		}
	}
}

// GetRolePermissions 获取角色权限
func (s *RoleService) GetRolePermissions(roleID uint) ([]Permission, error) {
	// 检查角色是否存在
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("角色不存在")
	}

	// 获取角色权限
	policies := s.casbinEnforcer.GetFilteredPolicy(0, role.Code)

	permissions := make([]Permission, 0, len(policies))
	for _, policy := range policies {
		if len(policy) >= 3 {
			permissions = append(permissions, Permission{
				Path:   policy[1],
				Method: policy[2],
			})
		}
	}

	return permissions, nil
}

// GetAllRoles 获取所有角色
func (s *RoleService) GetAllRoles() ([]*entity.Role, error) {
	return s.roleRepo.ListAll()
}

// UpdateRolePermissionsByAPIIDs 通过API ID更新角色权限
func (s *RoleService) UpdateRolePermissionsByAPIIDs(req *UpdateRolePermissionsByAPIIDsRequest) error {
	// 检查角色是否存在
	role, err := s.roleRepo.GetByID(req.RoleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("角色不存在")
	}

	// 选择性清除角色权限，保留指定的API权限
	s.deleteRolePermissionsExcept(role.Code, []string{
		"/api/v1/apis*",
		"/api/v1/apis/*",
		"/api/v1/business/*",
		"/api/v1/business*",
		"/api/v1/departments*",
		"/api/v1/departments/*",
		"/api/v1/role",
		"/api/v1/role/*",
		"/api/v1/user",
		"/api/v1/user/*",
	})

	// 确定部门参数
	deptParam := "*" // 默认为全部门
	if req.DeptID > 0 {
		deptParam = strconv.FormatUint(uint64(req.DeptID), 10)
	}

	// 根据API ID获取API信息并添加权限
	for _, apiID := range req.APIIDs {
		api, err := s.apiRepo.GetByID(apiID)
		if err != nil || api == nil {
			continue // 跳过不存在的API
		}

		// 添加权限（包含部门维度）
		_, err = s.casbinEnforcer.AddPolicyWithDept(role.Code, api.Path, api.Method, deptParam, "allow")
		if err != nil {
			return err
		}
	}

	// 同步策略到内存（安全方式）
	return s.casbinEnforcer.SyncPolicyToMemory()
}

// GetRoleAPIIDs 获取角色拥有的API ID列表
func (s *RoleService) GetRoleAPIIDs(roleID uint) ([]uint, error) {
	// 检查角色是否存在
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("角色不存在")
	}

	// 获取角色权限
	policies := s.casbinEnforcer.GetFilteredPolicy(0, role.Code)

	// 获取所有API，用于匹配
	// 这里需要一个获取所有API的方法，我们先假设有这个方法
	// 实际实现中可能需要在API仓库中添加这个方法
	var apiIDs []uint
	for _, policy := range policies {
		if len(policy) >= 3 {
			path := policy[1]
			method := policy[2]

			// 根据path和method查找API
			api, err := s.apiRepo.GetByPath(path, method)
			if err == nil && api != nil {
				apiIDs = append(apiIDs, api.ID)
			}
		}
	}

	return apiIDs, nil
}

// ListRoleByDept 按部门获取角色列表
func (s *RoleService) ListRoleByDept(req *ListRoleRequest) (*ListRoleResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 检查部门是否存在
	if req.DeptID > 0 {
		_, err := s.deptRepo.GetByID(req.DeptID)
		if err != nil {
			return nil, errors.New("部门不存在")
		}
	}

	roles, total, err := s.roleRepo.List(req.Page, req.PageSize, req.Query, req.DeptID)
	if err != nil {
		return nil, err
	}

	// 为每个角色获取部门名称
	rolesWithDept := make([]*RoleWithDepartment, len(roles))
	for i, role := range roles {
		deptName := "未知部门"
		if role.DeptID > 0 {
			dept, err := s.deptRepo.GetByID(role.DeptID)
			if err == nil && dept != nil {
				deptName = dept.Name
			}
		}

		rolesWithDept[i] = &RoleWithDepartment{
			Role:     role,
			DeptName: deptName,
		}
	}

	return &ListRoleResponse{
		Total: total,
		Items: rolesWithDept,
	}, nil
}

// AssignRolePermissionsRequest 分配角色权限请求
type AssignRolePermissionsRequest struct {
	RoleID      uint         `json:"role_id" binding:"required"`
	Permissions []Permission `json:"permissions" binding:"required"`
	DeptID      uint         `json:"dept_id"` // 部门ID，用于权限范围控制
}

// AssignRolePermissions 分配角色权限（部门级别）
func (s *RoleService) AssignRolePermissions(req *AssignRolePermissionsRequest) error {
	// 检查角色是否存在
	role, err := s.roleRepo.GetByID(req.RoleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("角色不存在")
	}

	// 检查角色是否属于指定部门
	if req.DeptID > 0 && role.DeptID != req.DeptID {
		return errors.New("角色不属于指定部门")
	}

	// 选择性清除角色权限，保留指定的API权限
	s.deleteRolePermissionsExcept(role.Code, []string{
		"/api/v1/apis*",
		"/api/v1/apis/*",
		"/api/v1/business/*",
		"/api/v1/business*",
		"/api/v1/departments*",
		"/api/v1/departments/*",
		"/api/v1/role",
		"/api/v1/role/*",
		"/api/v1/user",
		"/api/v1/user/*",
	})

	// 添加新权限（带部门信息）
	for _, perm := range req.Permissions {
		if req.DeptID > 0 {
			// 使用带部门的权限添加方法
			_, err = s.casbinEnforcer.AddPolicyWithDept(role.Code, perm.Path, perm.Method, strconv.FormatUint(uint64(req.DeptID), 10), "allow")
		} else {
			_, err = s.casbinEnforcer.AddPolicy(role.Code, perm.Path, perm.Method, "allow")
		}
		if err != nil {
			return err
		}
	}

	// 不保存策略到存储，避免清空策略表
	return nil
}
