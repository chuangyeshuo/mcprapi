package service

import (
	"errors"
	"strings"

	"mcprapi/backend/internal/domain/entity"
	"mcprapi/backend/internal/domain/repository"
)

// BusinessService 业务线服务
type BusinessService struct {
	businessRepo repository.BusinessRepository
	deptRepo     repository.DepartmentRepository
}

// NewBusinessService 创建业务线服务
func NewBusinessService(businessRepo repository.BusinessRepository, deptRepo repository.DepartmentRepository) *BusinessService {
	return &BusinessService{
		businessRepo: businessRepo,
		deptRepo:     deptRepo,
	}
}

// CreateBusinessRequest 创建业务线请求
type CreateBusinessRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	DeptID      uint   `json:"dept_id" binding:"required"`
	Description string `json:"description"`
	Owner       string `json:"owner"`
	Email       string `json:"email" binding:"omitempty,email"`
}

// UpdateBusinessRequest 更新业务线请求
type UpdateBusinessRequest struct {
	ID          uint   `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	DeptID      uint   `json:"dept_id" binding:"required"`
	Description string `json:"description"`
	Owner       string `json:"owner"`
	Email       string `json:"email" binding:"omitempty,email"`
	Status      int    `json:"status" binding:"oneof=0 1"`
}

// CreateBusiness 创建业务线
func (s *BusinessService) CreateBusiness(req *CreateBusinessRequest) (*entity.Business, error) {
	// 检查部门是否存在
	_, err := s.deptRepo.GetByID(req.DeptID)
	if err != nil {
		return nil, errors.New("部门不存在")
	}

	// 检查业务线编码是否已存在
	existingBusiness, _ := s.businessRepo.GetByCode(req.Code)
	if existingBusiness != nil {
		return nil, errors.New("业务线编码已存在")
	}

	// 创建业务线
	business := &entity.Business{
		Name:        req.Name,
		Code:        req.Code,
		DeptID:      req.DeptID,
		Description: req.Description,
		Owner:       req.Owner,
		Email:       req.Email,
		Status:      1, // 默认启用
	}

	err = s.businessRepo.Create(business)
	if err != nil {
		return nil, err
	}

	return business, nil
}

// UpdateBusiness 更新业务线
func (s *BusinessService) UpdateBusiness(req *UpdateBusinessRequest) (*entity.Business, error) {
	// 检查业务线是否存在
	business, err := s.businessRepo.GetByID(req.ID)
	if err != nil {
		return nil, errors.New("业务线不存在")
	}

	// 检查部门是否存在
	_, err = s.deptRepo.GetByID(req.DeptID)
	if err != nil {
		return nil, errors.New("部门不存在")
	}

	// 检查业务线编码是否已存在（排除自身）
	existingBusiness, _ := s.businessRepo.GetByCode(req.Code)
	if existingBusiness != nil && existingBusiness.ID != req.ID {
		return nil, errors.New("业务线编码已存在")
	}

	// 更新业务线
	business.Name = req.Name
	business.Code = req.Code
	business.DeptID = req.DeptID
	business.Description = req.Description
	business.Owner = req.Owner
	business.Email = req.Email
	business.Status = req.Status

	err = s.businessRepo.Update(business)
	if err != nil {
		return nil, err
	}

	return business, nil
}

// DeleteBusiness 删除业务线
func (s *BusinessService) DeleteBusiness(id uint) error {
	// 检查业务线是否存在
	_, err := s.businessRepo.GetByID(id)
	if err != nil {
		return errors.New("业务线不存在")
	}

	// 删除业务线
	return s.businessRepo.Delete(id)
}

// GetBusiness 获取业务线
func (s *BusinessService) GetBusiness(id uint) (*entity.Business, error) {
	return s.businessRepo.GetByID(id)
}

// ListBusiness 获取业务线列表
func (s *BusinessService) ListBusiness(deptID uint) ([]*entity.Business, error) {
	return s.businessRepo.List(deptID)
}

// ListAllBusiness 获取所有业务线
func (s *BusinessService) ListAllBusiness() ([]*entity.Business, error) {
	return s.businessRepo.ListAll()
}

// Create 创建业务线（别名）
func (s *BusinessService) Create(req *CreateBusinessRequest) (*entity.Business, error) {
	return s.CreateBusiness(req)
}

// Update 更新业务线（别名）
func (s *BusinessService) Update(req *UpdateBusinessRequest) (*entity.Business, error) {
	return s.UpdateBusiness(req)
}

// Delete 删除业务线（别名）
func (s *BusinessService) Delete(id uint) error {
	return s.DeleteBusiness(id)
}

// Get 获取业务线（别名）
func (s *BusinessService) Get(id uint) (*entity.Business, error) {
	return s.GetBusiness(id)
}

// List 获取业务线列表（分页）
func (s *BusinessService) List(page, pageSize int, query string, deptID uint) ([]*entity.Business, int64, error) {
	// 这里需要实现分页查询逻辑
	// 暂时返回所有业务线
	businesses, err := s.businessRepo.ListAll()
	if err != nil {
		return nil, 0, err
	}
	
	// 简单的过滤逻辑
	var filtered []*entity.Business
	for _, business := range businesses {
		if deptID > 0 && business.DeptID != deptID {
			continue
		}
		if query != "" && !strings.Contains(business.Name, query) && !strings.Contains(business.Code, query) {
			continue
		}
		filtered = append(filtered, business)
	}
	
	return filtered, int64(len(filtered)), nil
}

// GetByDepartment 获取部门业务线
func (s *BusinessService) GetByDepartment(deptID uint) ([]*entity.Business, error) {
	return s.businessRepo.List(deptID)
}

// GetAll 获取所有业务线（别名）
func (s *BusinessService) GetAll() ([]*entity.Business, error) {
	return s.ListAllBusiness()
}

// ListBusinessRequest 获取业务线列表请求
type ListBusinessRequest struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
	Query    string `json:"query" form:"query"`
	DeptID   uint   `json:"dept_id" form:"dept_id"` // 部门ID过滤
}

// BusinessWithDepartment 包含部门信息的业务线
type BusinessWithDepartment struct {
	*entity.Business
	DeptName string `json:"dept_name"`
}

// ListBusinessResponse 获取业务线列表响应
type ListBusinessResponse struct {
	Total int64                     `json:"total"`
	Items []*BusinessWithDepartment `json:"items"`
}

// ListBusinessByDept 按部门获取业务线列表
func (s *BusinessService) ListBusinessByDept(req *ListBusinessRequest) (*ListBusinessResponse, error) {
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

	businesses, total, err := s.List(req.Page, req.PageSize, req.Query, req.DeptID)
	if err != nil {
		return nil, err
	}

	// 为每个业务线获取部门名称
	businessesWithDept := make([]*BusinessWithDepartment, len(businesses))
	for i, business := range businesses {
		deptName := "未知部门"
		if business.DeptID > 0 {
			dept, err := s.deptRepo.GetByID(business.DeptID)
			if err == nil && dept != nil {
				deptName = dept.Name
			}
		}
		
		businessesWithDept[i] = &BusinessWithDepartment{
			Business: business,
			DeptName: deptName,
		}
	}

	return &ListBusinessResponse{
		Total: total,
		Items: businessesWithDept,
	}, nil
}

// GetTotalCount 获取业务线总数
func (s *BusinessService) GetTotalCount() (int64, error) {
	return s.businessRepo.GetTotalCount()
}