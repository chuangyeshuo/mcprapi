package repository

import (
	"errors"

	"gorm.io/gorm"

	"mcprapi/backend/internal/domain/entity"
	"mcprapi/backend/internal/domain/repository"
)

// APIRepositoryImpl API仓库实现
type APIRepositoryImpl struct {
	db *gorm.DB
}

// NewAPIRepository 创建API仓库
func NewAPIRepository(db *gorm.DB) repository.APIRepository {
	return &APIRepositoryImpl{db: db}
}

// Create 创建API
func (r *APIRepositoryImpl) Create(api *entity.API) error {
	return r.db.Create(api).Error
}

// Update 更新API
func (r *APIRepositoryImpl) Update(api *entity.API) error {
	return r.db.Save(api).Error
}

// Delete 删除API
func (r *APIRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entity.API{}, id).Error
}

// GetByID 根据ID获取API
func (r *APIRepositoryImpl) GetByID(id uint) (*entity.API, error) {
	var api entity.API
	if err := r.db.First(&api, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &api, nil
}

// GetByPath 根据路径和方法获取API
func (r *APIRepositoryImpl) GetByPath(path, method string) (*entity.API, error) {
	var api entity.API
	if err := r.db.Where("path = ? AND method = ?", path, method).First(&api).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &api, nil
}

// GetByPathAndMethod 根据路径和方法获取API (保持向后兼容)
func (r *APIRepositoryImpl) GetByPathAndMethod(path, method string) (*entity.API, error) {
	return r.GetByPath(path, method)
}

// List 获取API列表
func (r *APIRepositoryImpl) List(page, pageSize int, query string, businessID uint, category string) ([]*entity.API, int64, error) {
	var apis []*entity.API
	var total int64

	offset := (page - 1) * pageSize

	db := r.db
	if query != "" {
		db = db.Where("name LIKE ? OR path LIKE ?", "%"+query+"%", "%"+query+"%")
	}

	if businessID > 0 {
		db = db.Where("business_id = ?", businessID)
	}

	if category != "" {
		db = db.Where("category = ?", category)
	}

	if err := db.Model(&entity.API{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset(offset).Limit(pageSize).Find(&apis).Error; err != nil {
		return nil, 0, err
	}

	return apis, total, nil
}

// ListAll 获取所有API列表
func (r *APIRepositoryImpl) ListAll() ([]*entity.API, error) {
	var apis []*entity.API
	err := r.db.Find(&apis).Error
	return apis, err
}

// GetTotalCount 获取API总数
func (r *APIRepositoryImpl) GetTotalCount() (int64, error) {
	var count int64
	if err := r.db.Model(&entity.API{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetCountByBusiness 获取业务线API数量
func (r *APIRepositoryImpl) GetCountByBusiness(businessID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&entity.API{}).Where("business_id = ?", businessID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByDepartment 获取部门API数量
func (r *APIRepositoryImpl) CountByDepartment(deptID uint) (int64, error) {
	var count int64
	// 通过业务线表关联查询部门下的API数量
	if err := r.db.Table("apis").
		Joins("JOIN businesses ON apis.business_id = businesses.id").
		Where("businesses.dept_id = ?", deptID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetCategoryStats 获取API分类统计
func (r *APIRepositoryImpl) GetCategoryStats() ([]*repository.CategoryStatsItem, error) {
	var stats []*repository.CategoryStatsItem
	
	// 按HTTP方法统计
	if err := r.db.Model(&entity.API{}).
		Select("method as name, count(*) as value").
		Group("method").
		Scan(&stats).Error; err != nil {
		return nil, err
	}
	
	return stats, nil
}

// ListByBusiness 获取业务线下的API列表
func (r *APIRepositoryImpl) ListByBusiness(businessID uint) ([]*entity.API, error) {
	var apis []*entity.API
	if err := r.db.Where("business_id = ?", businessID).Find(&apis).Error; err != nil {
		return nil, err
	}
	return apis, nil
}

// ListCategories 获取API分类列表
func (r *APIRepositoryImpl) ListCategories() ([]*entity.APICategory, error) {
	var categories []*entity.APICategory
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetByBusiness 根据业务线获取API列表
func (r *APIRepositoryImpl) GetByBusiness(businessID uint) ([]*entity.API, error) {
	var apis []*entity.API
	if err := r.db.Where("business_id = ?", businessID).Find(&apis).Error; err != nil {
		return nil, err
	}
	return apis, nil
}

// GetCategories 获取API分类列表
func (r *APIRepositoryImpl) GetCategories() ([]*entity.APICategory, error) {
	var categories []*entity.APICategory
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetCategoryByID 根据ID获取API分类
func (r *APIRepositoryImpl) GetCategoryByID(id uint) (*entity.APICategory, error) {
	var category entity.APICategory
	if err := r.db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// CreateCategory 创建API分类
func (r *APIRepositoryImpl) CreateCategory(category *entity.APICategory) error {
	return r.db.Create(category).Error
}

// UpdateCategory 更新API分类
func (r *APIRepositoryImpl) UpdateCategory(category *entity.APICategory) error {
	return r.db.Save(category).Error
}

// DeleteCategory 删除API分类
func (r *APIRepositoryImpl) DeleteCategory(id uint) error {
	return r.db.Delete(&entity.APICategory{}, id).Error
}

// GetCategoryByCode 根据代码获取API分类
func (r *APIRepositoryImpl) GetCategoryByCode(code string) (*entity.APICategory, error) {
	var category entity.APICategory
	if err := r.db.Where("code = ?", code).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// GetCategoriesByParentID 根据父ID获取子分类
func (r *APIRepositoryImpl) GetCategoriesByParentID(parentID uint) ([]*entity.APICategory, error) {
	var categories []*entity.APICategory
	if err := r.db.Where("parent_id = ?", parentID).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetAPIsByCategoryID 根据分类ID获取API列表
func (r *APIRepositoryImpl) GetAPIsByCategoryID(categoryID uint) ([]*entity.API, error) {
	var apis []*entity.API
	if err := r.db.Where("category_id = ?", categoryID).Find(&apis).Error; err != nil {
		return nil, err
	}
	return apis, nil
}

// ListWithDepartmentFilter 根据部门过滤获取API列表
func (r *APIRepositoryImpl) ListWithDepartmentFilter(page, pageSize int, query string, businessIDs []uint, category string) ([]*entity.API, int64, error) {
	var apis []*entity.API
	var total int64

	offset := (page - 1) * pageSize

	db := r.db
	if query != "" {
		db = db.Where("name LIKE ? OR path LIKE ?", "%"+query+"%", "%"+query+"%")
	}

	if len(businessIDs) > 0 {
		db = db.Where("business_id IN ?", businessIDs)
	}

	if category != "" {
		db = db.Where("category = ?", category)
	}

	if err := db.Model(&entity.API{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset(offset).Limit(pageSize).Find(&apis).Error; err != nil {
		return nil, 0, err
	}

	return apis, total, nil
}
