package service

import (
	"errors"
	"strings"

	"mcprapi/backend/internal/domain/entity"
	"mcprapi/backend/internal/domain/repository"
)

// APIService API服务
type APIService struct {
	apiRepo      repository.APIRepository
	businessRepo repository.BusinessRepository
	userRepo     repository.UserRepository
	deptRepo     repository.DepartmentRepository
}

// NewAPIService 创建API服务
func NewAPIService(apiRepo repository.APIRepository, businessRepo repository.BusinessRepository, userRepo repository.UserRepository, deptRepo repository.DepartmentRepository) *APIService {
	return &APIService{
		apiRepo:      apiRepo,
		businessRepo: businessRepo,
		userRepo:     userRepo,
		deptRepo:     deptRepo,
	}
}

// CreateAPIRequest 创建API请求
type CreateAPIRequest struct {
	Name        string `json:"name" binding:"required"`
	Path        string `json:"path" binding:"required"`
	Method      string `json:"method" binding:"required,oneof=GET POST PUT DELETE"`
	Description string `json:"description"`
	DeptID      uint   `json:"dept_id" binding:"required"`
	BusinessID  uint   `json:"business_id" binding:"required"`
	CategoryID  uint   `json:"category_id"`
}

// UpdateAPIRequest 更新API请求
type UpdateAPIRequest struct {
	ID          uint   `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Path        string `json:"path" binding:"required"`
	Method      string `json:"method" binding:"required,oneof=GET POST PUT DELETE"`
	Description string `json:"description"`
	DeptID      uint   `json:"dept_id" binding:"required"`
	BusinessID  uint   `json:"business_id" binding:"required"`
	CategoryID  uint   `json:"category_id"`
	Status      int    `json:"status" binding:"oneof=0 1"`
}

// CreateAPI 创建API
func (s *APIService) CreateAPI(req *CreateAPIRequest) (*entity.API, error) {
	// 检查部门是否存在
	dept, err := s.deptRepo.GetByID(req.DeptID)
	if err != nil {
		return nil, errors.New("部门不存在")
	}

	// 检查业务线是否存在
	business, err := s.businessRepo.GetByID(req.BusinessID)
	if err != nil {
		return nil, errors.New("业务线不存在")
	}

	// 检查业务线是否属于指定部门
	if business.DeptID != req.DeptID {
		return nil, errors.New("业务线不属于指定部门")
	}

	// 自动在路径前拼接部门代码
	finalPath := req.Path
	if dept.Code != "" && !strings.HasPrefix(req.Path, "/"+dept.Code+"/") {
		// 如果路径不是以部门代码开头，则自动添加
		if strings.HasPrefix(req.Path, "/") {
			finalPath = "/" + dept.Code + req.Path
		} else {
			finalPath = "/" + dept.Code + "/" + req.Path
		}
	}

	// 检查API路径是否已存在
	existingAPI, _ := s.apiRepo.GetByPath(finalPath, req.Method)
	if existingAPI != nil {
		return nil, errors.New("API路径和方法已存在")
	}

	// 创建API
	api := &entity.API{
		Name:        req.Name,
		Path:        finalPath,
		Method:      req.Method,
		Description: req.Description,
		BusinessID:  business.ID,
		CategoryID:  req.CategoryID,
		Status:      1, // 默认启用
	}

	err = s.apiRepo.Create(api)
	if err != nil {
		return nil, err
	}

	return api, nil
}

// UpdateAPI 更新API
func (s *APIService) UpdateAPI(req *UpdateAPIRequest) (*entity.API, error) {
	// 检查API是否存在
	api, err := s.apiRepo.GetByID(req.ID)
	if err != nil {
		return nil, errors.New("API不存在")
	}

	// 检查部门是否存在
	dept, err := s.deptRepo.GetByID(req.DeptID)
	if err != nil {
		return nil, errors.New("部门不存在")
	}

	// 检查业务线是否存在
	business, err := s.businessRepo.GetByID(req.BusinessID)
	if err != nil {
		return nil, errors.New("业务线不存在")
	}

	// 检查业务线是否属于指定部门
	if business.DeptID != req.DeptID {
		return nil, errors.New("业务线不属于指定部门")
	}

	// 自动在路径前拼接部门代码
	finalPath := req.Path
	if dept.Code != "" && !strings.HasPrefix(req.Path, "/"+dept.Code+"/") {
		// 如果路径不是以部门代码开头，则自动添加
		if strings.HasPrefix(req.Path, "/") {
			finalPath = "/" + dept.Code + req.Path
		} else {
			finalPath = "/" + dept.Code + "/" + req.Path
		}
	}

	// 检查API路径是否已存在（排除自身）
	existingAPI, _ := s.apiRepo.GetByPath(finalPath, req.Method)
	if existingAPI != nil && existingAPI.ID != req.ID {
		return nil, errors.New("API路径和方法已存在")
	}

	// 更新API
	api.Name = req.Name
	api.Path = finalPath
	api.Method = req.Method
	api.Description = req.Description
	api.BusinessID = req.BusinessID
	api.CategoryID = req.CategoryID
	api.Status = req.Status

	err = s.apiRepo.Update(api)
	if err != nil {
		return nil, err
	}

	return api, nil
}

// DeleteAPI 删除API
func (s *APIService) DeleteAPI(id uint) error {
	// 检查API是否存在
	_, err := s.apiRepo.GetByID(id)
	if err != nil {
		return errors.New("API不存在")
	}

	// 删除API
	return s.apiRepo.Delete(id)
}

// GetAPI 获取API
func (s *APIService) GetAPI(id uint) (*entity.API, error) {
	return s.apiRepo.GetByID(id)
}

// ListAPIRequest 获取API列表请求
type ListAPIRequest struct {
	Page       int    `json:"page" form:"page"`
	PageSize   int    `json:"page_size" form:"page_size"`
	Query      string `json:"query" form:"query"`
	BusinessID uint   `json:"business_id" form:"business_id"`
	Category   string `json:"category" form:"category"`
	UserID     uint   `json:"user_id" form:"user_id"` // 当前用户ID，用于权限控制
}

// APIWithDepartment API带部门信息
type APIWithDepartment struct {
	*entity.API
	BusinessName   string `json:"business_name"`   // 业务线名称
	DepartmentID   uint   `json:"dept_id"`         // 部门ID
	DepartmentName string `json:"department_name"` // 部门名称
}

// ListAPIResponse 获取API列表响应
type ListAPIResponse struct {
	Total int64                `json:"total"`
	Items []*APIWithDepartment `json:"items"`
}

// ListAPI 获取API列表
func (s *APIService) ListAPI(req *ListAPIRequest) (*ListAPIResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 基于用户部门的权限控制
	var allowedBusinessIDs []uint
	if req.UserID > 0 {
		// 获取当前用户信息
		user, err := s.userRepo.GetByID(req.UserID)
		if err != nil {
			return nil, errors.New("用户不存在")
		}

		// 检查是否为admin用户（假设admin用户名为"admin"）
		if user.Username != "admin" {
			// 普通用户只能看到自己部门的业务线API
			if user.DeptID > 0 {
				// 获取用户部门下的所有业务线
				businesses, err := s.businessRepo.List(user.DeptID)
				if err != nil {
					return nil, err
				}

				// 提取业务线ID列表
				for _, business := range businesses {
					allowedBusinessIDs = append(allowedBusinessIDs, business.ID)
				}
			}
		}
		// admin用户不设置allowedBusinessIDs，可以看到所有API
	}

	apis, total, err := s.apiRepo.ListWithDepartmentFilter(req.Page, req.PageSize, req.Query, allowedBusinessIDs, req.Category)
	if err != nil {
		return nil, err
	}

	// 转换为包含部门信息的结构
	apiWithDepts := make([]*APIWithDepartment, len(apis))
	for i, api := range apis {
		apiWithDept := &APIWithDepartment{
			API: api,
		}

		// 获取业务线信息
		if api.BusinessID > 0 {
			business, err := s.businessRepo.GetByID(api.BusinessID)
			if err == nil && business != nil {
				apiWithDept.BusinessName = business.Name
				apiWithDept.DepartmentID = business.DeptID

				// 获取部门信息
				if business.DeptID > 0 {
					dept, err := s.deptRepo.GetByID(business.DeptID)
					if err == nil && dept != nil {
						apiWithDept.DepartmentName = dept.Name
					}
				}
			}
		}

		apiWithDepts[i] = apiWithDept
	}

	return &ListAPIResponse{
		Total: total,
		Items: apiWithDepts,
	}, nil
}

// GetTotalCount 获取API总数
func (s *APIService) GetTotalCount() (int64, error) {
	return s.apiRepo.GetTotalCount()
}

// GetCategoryStats 获取API分类统计
func (s *APIService) GetCategoryStats() ([]*repository.CategoryStatsItem, error) {
	return s.apiRepo.GetCategoryStats()
}

// ListAllAPIs 获取所有API列表
func (s *APIService) ListAllAPIs() ([]*entity.API, error) {
	return s.apiRepo.ListAll()
}

// ListAPIByBusiness 获取业务线下的API列表
func (s *APIService) ListAPIByBusiness(businessID uint) ([]*entity.API, error) {
	return s.apiRepo.ListByBusiness(businessID)
}

// ListAPICategories 获取API分类列表
func (s *APIService) ListAPICategories() ([]*entity.APICategory, error) {
	return s.apiRepo.ListCategories()
}

// CreateAPICategoryRequest 创建API分类请求
type CreateAPICategoryRequest struct {
	Name     string `json:"name" binding:"required"`
	Code     string `json:"code" binding:"required"`
	ParentID uint   `json:"parent_id"`
	Sort     int    `json:"sort"`
}

// UpdateAPICategoryRequest 更新API分类请求
type UpdateAPICategoryRequest struct {
	ID       uint   `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Code     string `json:"code" binding:"required"`
	ParentID uint   `json:"parent_id"`
	Sort     int    `json:"sort"`
}

// CreateAPICategory 创建API分类
func (s *APIService) CreateAPICategory(req *CreateAPICategoryRequest) (*entity.APICategory, error) {
	// 检查分类代码是否已存在
	existingCategory, _ := s.apiRepo.GetCategoryByCode(req.Code)
	if existingCategory != nil {
		return nil, errors.New("分类代码已存在")
	}

	// 创建分类
	category := &entity.APICategory{
		Name:     req.Name,
		Code:     req.Code,
		ParentID: req.ParentID,
		Sort:     req.Sort,
	}

	err := s.apiRepo.CreateCategory(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// UpdateAPICategory 更新API分类
func (s *APIService) UpdateAPICategory(req *UpdateAPICategoryRequest) (*entity.APICategory, error) {
	// 检查分类是否存在
	category, err := s.apiRepo.GetCategoryByID(req.ID)
	if err != nil {
		return nil, errors.New("分类不存在")
	}

	// 检查分类代码是否已存在（排除自身）
	existingCategory, _ := s.apiRepo.GetCategoryByCode(req.Code)
	if existingCategory != nil && existingCategory.ID != req.ID {
		return nil, errors.New("分类代码已存在")
	}

	// 更新分类
	category.Name = req.Name
	category.Code = req.Code
	category.ParentID = req.ParentID
	category.Sort = req.Sort

	err = s.apiRepo.UpdateCategory(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteAPICategory 删除API分类
func (s *APIService) DeleteAPICategory(id uint) error {
	// 检查分类是否存在
	_, err := s.apiRepo.GetCategoryByID(id)
	if err != nil {
		return errors.New("分类不存在")
	}

	// 检查是否有子分类
	children, err := s.apiRepo.GetCategoriesByParentID(id)
	if err == nil && len(children) > 0 {
		return errors.New("存在子分类，无法删除")
	}

	// 检查是否有API使用该分类
	apis, err := s.apiRepo.GetAPIsByCategoryID(id)
	if err == nil && len(apis) > 0 {
		return errors.New("存在使用该分类的API，无法删除")
	}

	// 删除分类
	return s.apiRepo.DeleteCategory(id)
}

// GetAPICategory 获取API分类
func (s *APIService) GetAPICategory(id uint) (*entity.APICategory, error) {
	return s.apiRepo.GetCategoryByID(id)
}

// ListAPIByDeptRequest 按部门获取API列表请求
type ListAPIByDeptRequest struct {
	Page       int    `json:"page" form:"page"`
	PageSize   int    `json:"page_size" form:"page_size"`
	Query      string `json:"query" form:"query"`
	DeptID     uint   `json:"dept_id" form:"dept_id" binding:"required"`
	BusinessID uint   `json:"business_id" form:"business_id"`
	Category   string `json:"category" form:"category"`
}

// ListAPIByDept 按部门获取API列表
func (s *APIService) ListAPIByDept(req *ListAPIByDeptRequest) (*ListAPIResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 检查部门是否存在
	_, err := s.deptRepo.GetByID(req.DeptID)
	if err != nil {
		return nil, errors.New("部门不存在")
	}

	// 获取部门下的所有业务线
	businesses, err := s.businessRepo.List(req.DeptID)
	if err != nil {
		return nil, err
	}

	// 提取业务线ID列表
	var allowedBusinessIDs []uint
	for _, business := range businesses {
		// 如果指定了业务线ID，只包含该业务线
		if req.BusinessID > 0 {
			if business.ID == req.BusinessID {
				allowedBusinessIDs = append(allowedBusinessIDs, business.ID)
			}
		} else {
			allowedBusinessIDs = append(allowedBusinessIDs, business.ID)
		}
	}

	// 如果指定了业务线但不属于该部门，返回空结果
	if req.BusinessID > 0 && len(allowedBusinessIDs) == 0 {
		return &ListAPIResponse{
			Total: 0,
			Items: []*APIWithDepartment{},
		}, nil
	}

	apis, total, err := s.apiRepo.ListWithDepartmentFilter(req.Page, req.PageSize, req.Query, allowedBusinessIDs, req.Category)
	if err != nil {
		return nil, err
	}

	// 转换为包含部门信息的结构
	apiWithDepts := make([]*APIWithDepartment, len(apis))
	for i, api := range apis {
		apiWithDept := &APIWithDepartment{
			API: api,
		}

		// 获取业务线信息
		if api.BusinessID > 0 {
			business, err := s.businessRepo.GetByID(api.BusinessID)
			if err == nil && business != nil {
				apiWithDept.BusinessName = business.Name
				apiWithDept.DepartmentID = business.DeptID

				// 获取部门信息
				if business.DeptID > 0 {
					dept, err := s.deptRepo.GetByID(business.DeptID)
					if err == nil && dept != nil {
						apiWithDept.DepartmentName = dept.Name
					}
				}
			}
		}

		apiWithDepts[i] = apiWithDept
	}

	return &ListAPIResponse{
		Total: total,
		Items: apiWithDepts,
	}, nil
}
