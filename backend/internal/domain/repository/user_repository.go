package repository

import (
	"mcprapi/backend/internal/domain/entity"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	// Create 创建用户
	Create(user *entity.User) error

	// Update 更新用户
	Update(user *entity.User) error

	// Delete 删除用户
	Delete(id uint) error

	// GetByID 根据ID获取用户
	GetByID(id uint) (*entity.User, error)

	// GetByUsername 根据用户名获取用户
	GetByUsername(username string) (*entity.User, error)

	// List 获取用户列表
	List(page, pageSize int, query string) ([]*entity.User, int64, error)

	// ListAll 获取所有用户列表
	ListAll() ([]*entity.User, error)

	// GetTotalCount 获取用户总数
	GetTotalCount() (int64, error)

	// ListByDept 根据部门获取用户列表
	ListByDept(page, pageSize int, query string, deptID uint) ([]*entity.User, int64, error)

	// GetUserRoles 获取用户角色
	GetUserRoles(userID uint) ([]string, error)

	// AssignRoles 分配角色给用户
	AssignRoles(userID uint, roleIDs []uint) error
}