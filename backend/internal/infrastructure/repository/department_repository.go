package repository

import (
	"errors"

	"gorm.io/gorm"

	"mcprapi/backend/internal/domain/entity"
	"mcprapi/backend/internal/domain/repository"
)

// DepartmentRepositoryImpl 部门仓库实现
type DepartmentRepositoryImpl struct {
	db *gorm.DB
}

// NewDepartmentRepository 创建部门仓库
func NewDepartmentRepository(db *gorm.DB) repository.DepartmentRepository {
	return &DepartmentRepositoryImpl{db: db}
}

// Create 创建部门
func (r *DepartmentRepositoryImpl) Create(department *entity.Department) error {
	return r.db.Create(department).Error
}

// Update 更新部门
func (r *DepartmentRepositoryImpl) Update(department *entity.Department) error {
	return r.db.Save(department).Error
}

// Delete 删除部门
func (r *DepartmentRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entity.Department{}, id).Error
}

// GetByID 根据ID获取部门
func (r *DepartmentRepositoryImpl) GetByID(id uint) (*entity.Department, error) {
	var department entity.Department
	if err := r.db.First(&department, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &department, nil
}

// GetByCode 根据编码获取部门
func (r *DepartmentRepositoryImpl) GetByCode(code string) (*entity.Department, error) {
	var department entity.Department
	if err := r.db.Where("code = ?", code).First(&department).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &department, nil
}

// List 获取部门列表
func (r *DepartmentRepositoryImpl) List(parentID uint) ([]*entity.Department, error) {
	var departments []*entity.Department

	db := r.db
	if parentID > 0 {
		db = db.Where("parent_id = ?", parentID)
	} else {
		// 如果parentID为0，获取顶级部门（parent_id为NULL或0）
		db = db.Where("parent_id IS NULL OR parent_id = 0")
	}

	if err := db.Find(&departments).Error; err != nil {
		return nil, err
	}

	return departments, nil
}

// ListAll 获取所有部门
func (r *DepartmentRepositoryImpl) ListAll() ([]*entity.Department, error) {
	var departments []*entity.Department
	if err := r.db.Find(&departments).Error; err != nil {
		return nil, err
	}
	return departments, nil
}

// GetTotalCount 获取部门总数
func (r *DepartmentRepositoryImpl) GetTotalCount() (int64, error) {
	var count int64
	err := r.db.Model(&entity.Department{}).Count(&count).Error
	return count, err
}

// GetChildren 获取子部门
func (r *DepartmentRepositoryImpl) GetChildren(parentID uint) ([]*entity.Department, error) {
	var departments []*entity.Department
	if err := r.db.Where("parent_id = ?", parentID).Find(&departments).Error; err != nil {
		return nil, err
	}
	return departments, nil
}

// BusinessRepositoryImpl 业务线仓库实现
type BusinessRepositoryImpl struct {
	db *gorm.DB
}

// NewBusinessRepository 创建业务线仓库
func NewBusinessRepository(db *gorm.DB) repository.BusinessRepository {
	return &BusinessRepositoryImpl{db: db}
}

// Create 创建业务线
func (r *BusinessRepositoryImpl) Create(business *entity.Business) error {
	return r.db.Create(business).Error
}

// Update 更新业务线
func (r *BusinessRepositoryImpl) Update(business *entity.Business) error {
	return r.db.Save(business).Error
}

// Delete 删除业务线
func (r *BusinessRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entity.Business{}, id).Error
}

// GetByID 根据ID获取业务线
func (r *BusinessRepositoryImpl) GetByID(id uint) (*entity.Business, error) {
	var business entity.Business
	if err := r.db.First(&business, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &business, nil
}

// GetByCode 根据编码获取业务线
func (r *BusinessRepositoryImpl) GetByCode(code string) (*entity.Business, error) {
	var business entity.Business
	if err := r.db.Where("code = ?", code).First(&business).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &business, nil
}

// List 获取业务线列表
func (r *BusinessRepositoryImpl) List(deptID uint) ([]*entity.Business, error) {
	var businesses []*entity.Business

	db := r.db
	if deptID > 0 {
		db = db.Where("dept_id = ?", deptID)
	}

	if err := db.Find(&businesses).Error; err != nil {
		return nil, err
	}

	return businesses, nil
}

// GetByDepartment 根据部门获取业务线
func (r *BusinessRepositoryImpl) GetByDepartment(departmentID uint) ([]*entity.Business, error) {
	var businesses []*entity.Business
	if err := r.db.Where("dept_id = ?", departmentID).Find(&businesses).Error; err != nil {
		return nil, err
	}
	return businesses, nil
}

// ListAll 获取所有业务线
func (r *BusinessRepositoryImpl) ListAll() ([]*entity.Business, error) {
	var businesses []*entity.Business
	if err := r.db.Find(&businesses).Error; err != nil {
		return nil, err
	}
	return businesses, nil
}

// GetTotalCount 获取业务线总数
func (r *BusinessRepositoryImpl) GetTotalCount() (int64, error) {
	var count int64
	err := r.db.Model(&entity.Business{}).Count(&count).Error
	return count, err
}
