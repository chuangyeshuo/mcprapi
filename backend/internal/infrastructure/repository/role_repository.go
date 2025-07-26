package repository

import (
	"errors"

	"gorm.io/gorm"

	"mcprapi/backend/internal/domain/entity"
	"mcprapi/backend/internal/domain/repository"
)

// RoleRepositoryImpl 角色仓库实现
type RoleRepositoryImpl struct {
	db *gorm.DB
}

// NewRoleRepository 创建角色仓库
func NewRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &RoleRepositoryImpl{db: db}
}

// Create 创建角色
func (r *RoleRepositoryImpl) Create(role *entity.Role) error {
	return r.db.Create(role).Error
}

// Update 更新角色
func (r *RoleRepositoryImpl) Update(role *entity.Role) error {
	return r.db.Save(role).Error
}

// Delete 删除角色
func (r *RoleRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entity.Role{}, id).Error
}

// GetByID 根据ID获取角色
func (r *RoleRepositoryImpl) GetByID(id uint) (*entity.Role, error) {
	var role entity.Role
	if err := r.db.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

// GetByCode 根据编码获取角色
func (r *RoleRepositoryImpl) GetByCode(code string) (*entity.Role, error) {
	var role entity.Role
	if err := r.db.Where("code = ?", code).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

// GetByName 根据名称获取角色
func (r *RoleRepositoryImpl) GetByName(name string) (*entity.Role, error) {
	var role entity.Role
	if err := r.db.Where("name = ?", name).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

// List 获取角色列表
func (r *RoleRepositoryImpl) List(page, pageSize int, query string, deptID uint) ([]*entity.Role, int64, error) {
	var roles []*entity.Role
	var total int64

	offset := (page - 1) * pageSize

	db := r.db
	if query != "" {
		db = db.Where("name LIKE ? OR code LIKE ?", "%"+query+"%", "%"+query+"%")
	}

	if deptID > 0 {
		db = db.Where("dept_id = ?", deptID)
	}

	if err := db.Model(&entity.Role{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset(offset).Limit(pageSize).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// ListAll 获取所有角色
func (r *RoleRepositoryImpl) ListAll() ([]*entity.Role, error) {
	var roles []*entity.Role
	if err := r.db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// ListByDept 根据部门ID获取角色列表
func (r *RoleRepositoryImpl) ListByDept(deptID uint) ([]*entity.Role, error) {
	var roles []*entity.Role
	if err := r.db.Where("dept_id = ?", deptID).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// GetRolePermissions 获取角色权限
func (r *RoleRepositoryImpl) GetRolePermissions(roleID uint) ([]map[string]string, error) {
	// 这里应该从Casbin中获取角色权限
	// 示例：
	// 1. 获取Casbin执行器
	// 2. 获取角色权限
	// 3. 返回权限列表
	return []map[string]string{}, nil
}

// UpdateRolePermissions 更新角色权限
func (r *RoleRepositoryImpl) UpdateRolePermissions(roleID uint, permissions []map[string]string) error {
	// 这里应该更新Casbin中的角色权限
	// 示例：
	// 1. 获取Casbin执行器
	// 2. 删除角色现有权限
	// 3. 添加新权限
	// 4. 保存策略
	return nil
}
