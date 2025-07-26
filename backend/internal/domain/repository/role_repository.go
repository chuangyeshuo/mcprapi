package repository

import (
	"mcprapi/backend/internal/domain/entity"
)

// RoleRepository 角色仓库接口
type RoleRepository interface {
	// Create 创建角色
	Create(role *entity.Role) error

	// Update 更新角色
	Update(role *entity.Role) error

	// Delete 删除角色
	Delete(id uint) error

	// GetByID 根据ID获取角色
	GetByID(id uint) (*entity.Role, error)

	// GetByCode 根据编码获取角色
	GetByCode(code string) (*entity.Role, error)

	// GetByName 根据名称获取角色
	GetByName(name string) (*entity.Role, error)

	// List 获取角色列表
	List(page, pageSize int, query string, deptID uint) ([]*entity.Role, int64, error)

	// ListAll 获取所有角色
	ListAll() ([]*entity.Role, error)

	// ListByDept 根据部门ID获取角色列表
	ListByDept(deptID uint) ([]*entity.Role, error)

	// GetRolePermissions 获取角色权限
	GetRolePermissions(roleID uint) ([]map[string]string, error)

	// UpdateRolePermissions 更新角色权限
	UpdateRolePermissions(roleID uint, permissions []map[string]string) error
}
