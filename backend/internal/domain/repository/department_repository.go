package repository

import (
	"mcprapi/backend/internal/domain/entity"
)

// DepartmentRepository 部门仓库接口
type DepartmentRepository interface {
	// Create 创建部门
	Create(dept *entity.Department) error

	// Update 更新部门
	Update(dept *entity.Department) error

	// Delete 删除部门
	Delete(id uint) error

	// GetByID 根据ID获取部门
	GetByID(id uint) (*entity.Department, error)

	// GetByCode 根据编码获取部门
	GetByCode(code string) (*entity.Department, error)

	// List 获取部门列表
	List(parentID uint) ([]*entity.Department, error)

	// ListAll 获取所有部门
	ListAll() ([]*entity.Department, error)

	// GetTotalCount 获取部门总数
	GetTotalCount() (int64, error)
}

// BusinessRepository 业务线仓库接口
type BusinessRepository interface {
	// Create 创建业务线
	Create(business *entity.Business) error

	// Update 更新业务线
	Update(business *entity.Business) error

	// Delete 删除业务线
	Delete(id uint) error

	// GetByID 根据ID获取业务线
	GetByID(id uint) (*entity.Business, error)

	// GetByCode 根据编码获取业务线
	GetByCode(code string) (*entity.Business, error)

	// List 获取业务线列表
	List(deptID uint) ([]*entity.Business, error)

	// ListAll 获取所有业务线
	ListAll() ([]*entity.Business, error)

	// GetTotalCount 获取业务线总数
	GetTotalCount() (int64, error)
}