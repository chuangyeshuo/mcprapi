/*
 * @Author: lidi10@staff.weibo.com
 * @Date: 2025-07-19 17:26:24
 * @LastEditTime: 2025-07-20 12:33:23
 * @LastEditors: lidi10@staff.weibo.com
 * @Description:
 * Copyright (c) 2023 by Weibo, All Rights Reserved.
 */
package repository

import (
	"mcprapi/backend/internal/domain/entity"
)

// APIRepository API仓库接口
type APIRepository interface {
	// Create 创建API
	Create(api *entity.API) error

	// Update 更新API
	Update(api *entity.API) error

	// Delete 删除API
	Delete(id uint) error

	// GetByID 根据ID获取API
	GetByID(id uint) (*entity.API, error)

	// GetByPath 根据路径和方法获取API
	GetByPath(path, method string) (*entity.API, error)

	// List 获取API列表
	List(page, pageSize int, query string, businessID uint, category string) ([]*entity.API, int64, error)

	// ListAll 获取所有API列表
	ListAll() ([]*entity.API, error)

	// ListByBusiness 获取业务线下的API列表
	ListByBusiness(businessID uint) ([]*entity.API, error)

	// ListCategories 获取API分类列表
	ListCategories() ([]*entity.APICategory, error)

	// CreateCategory 创建API分类
	CreateCategory(category *entity.APICategory) error

	// UpdateCategory 更新API分类
	UpdateCategory(category *entity.APICategory) error

	// DeleteCategory 删除API分类
	DeleteCategory(id uint) error

	// GetCategoryByID 根据ID获取API分类
	GetCategoryByID(id uint) (*entity.APICategory, error)

	// GetCategoryByCode 根据代码获取API分类
	GetCategoryByCode(code string) (*entity.APICategory, error)

	// GetCategoriesByParentID 根据父ID获取子分类
	GetCategoriesByParentID(parentID uint) ([]*entity.APICategory, error)

	// GetAPIsByCategoryID 根据分类ID获取API列表
	GetAPIsByCategoryID(categoryID uint) ([]*entity.API, error)

	// ListWithDepartmentFilter 根据部门过滤获取API列表
	ListWithDepartmentFilter(page, pageSize int, query string, businessIDs []uint, category string) ([]*entity.API, int64, error)

	// GetTotalCount 获取API总数
	GetTotalCount() (int64, error)

	// GetCountByBusiness 获取业务线API数量
	GetCountByBusiness(businessID uint) (int64, error)

	// CountByDepartment 获取部门API数量
	CountByDepartment(deptID uint) (int64, error)

	// GetCategoryStats 获取API分类统计
	GetCategoryStats() ([]*CategoryStatsItem, error)
}

// CategoryStatsItem API分类统计项
type CategoryStatsItem struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}
