package service

import (
	"errors"
	"strings"

	"mcprapi/backend/internal/domain/entity"
	"mcprapi/backend/internal/domain/repository"
)

// DepartmentService 部门服务
type DepartmentService struct {
	deptRepo repository.DepartmentRepository
}

// NewDepartmentService 创建部门服务
func NewDepartmentService(deptRepo repository.DepartmentRepository) *DepartmentService {
	return &DepartmentService{
		deptRepo: deptRepo,
	}
}

// CreateDepartmentRequest 创建部门请求
type CreateDepartmentRequest struct {
	Name     string `json:"name" binding:"required"`
	Code     string `json:"code" binding:"required"`
	ParentID uint   `json:"parent_id"`
	Level    int    `json:"level" binding:"required,oneof=1 2 3"`
	Sort     int    `json:"sort"`
}

// UpdateDepartmentRequest 更新部门请求
type UpdateDepartmentRequest struct {
	ID       uint   `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Code     string `json:"code" binding:"required"`
	ParentID uint   `json:"parent_id"`
	Level    int    `json:"level" binding:"required,oneof=1 2 3"`
	Sort     int    `json:"sort"`
	Status   int    `json:"status" binding:"oneof=0 1"`
}

// CreateDepartment 创建部门
func (s *DepartmentService) CreateDepartment(req *CreateDepartmentRequest) (*entity.Department, error) {
	// 检查部门编码是否已存在
	existingDept, _ := s.deptRepo.GetByCode(req.Code)
	if existingDept != nil {
		return nil, errors.New("部门编码已存在")
	}

	// 检查父部门是否存在
	if req.ParentID > 0 {
		parentDept, err := s.deptRepo.GetByID(req.ParentID)
		if err != nil {
			return nil, errors.New("父部门不存在")
		}

		// 检查层级关系
		if parentDept.Level >= req.Level {
			return nil, errors.New("子部门层级必须大于父部门层级")
		}
	}

	// 创建部门
	dept := &entity.Department{
		Name:     req.Name,
		Code:     req.Code,
		ParentID: req.ParentID,
		Level:    req.Level,
		Sort:     req.Sort,
		Status:   1, // 默认启用
	}

	err := s.deptRepo.Create(dept)
	if err != nil {
		return nil, err
	}

	return dept, nil
}

// UpdateDepartment 更新部门
func (s *DepartmentService) UpdateDepartment(req *UpdateDepartmentRequest) (*entity.Department, error) {
	// 检查部门是否存在
	dept, err := s.deptRepo.GetByID(req.ID)
	if err != nil {
		return nil, errors.New("部门不存在")
	}

	// 检查部门编码是否已存在（排除自身）
	existingDept, _ := s.deptRepo.GetByCode(req.Code)
	if existingDept != nil && existingDept.ID != req.ID {
		return nil, errors.New("部门编码已存在")
	}

	// 检查父部门是否存在
	if req.ParentID > 0 {
		parentDept, err := s.deptRepo.GetByID(req.ParentID)
		if err != nil {
			return nil, errors.New("父部门不存在")
		}

		// 检查层级关系
		if parentDept.Level >= req.Level {
			return nil, errors.New("子部门层级必须大于父部门层级")
		}

		// 检查是否形成循环依赖
		if req.ParentID == req.ID {
			return nil, errors.New("不能将自己设为父部门")
		}
	}

	// 更新部门
	dept.Name = req.Name
	dept.Code = req.Code
	dept.ParentID = req.ParentID
	dept.Level = req.Level
	dept.Sort = req.Sort
	dept.Status = req.Status

	err = s.deptRepo.Update(dept)
	if err != nil {
		return nil, err
	}

	return dept, nil
}

// DeleteDepartment 删除部门
func (s *DepartmentService) DeleteDepartment(id uint) error {
	// 检查部门是否存在
	_, err := s.deptRepo.GetByID(id)
	if err != nil {
		return errors.New("部门不存在")
	}

	// 检查是否有子部门
	depts, err := s.deptRepo.List(id)
	if err != nil {
		return err
	}
	if len(depts) > 0 {
		return errors.New("请先删除子部门")
	}

	// 删除部门
	return s.deptRepo.Delete(id)
}

// GetDepartment 获取部门
func (s *DepartmentService) GetDepartment(id uint) (*entity.Department, error) {
	return s.deptRepo.GetByID(id)
}

// ListDepartment 获取部门列表
func (s *DepartmentService) ListDepartment(parentID uint) ([]*entity.Department, error) {
	return s.deptRepo.List(parentID)
}

// ListAllDepartment 获取所有部门
func (s *DepartmentService) ListAllDepartment() ([]*entity.Department, error) {
	return s.deptRepo.ListAll()
}

// Create 创建部门（别名）
func (s *DepartmentService) Create(req *CreateDepartmentRequest) (*entity.Department, error) {
	return s.CreateDepartment(req)
}

// Update 更新部门（别名）
func (s *DepartmentService) Update(req *UpdateDepartmentRequest) (*entity.Department, error) {
	return s.UpdateDepartment(req)
}

// Delete 删除部门（别名）
func (s *DepartmentService) Delete(id uint) error {
	return s.DeleteDepartment(id)
}

// Get 获取部门（别名）
func (s *DepartmentService) Get(id uint) (*entity.Department, error) {
	return s.GetDepartment(id)
}

// List 获取部门列表（支持分页和查询）
func (s *DepartmentService) List(page, size int, query string) ([]*entity.Department, int64, error) {
	// 这里需要根据实际的repository接口来实现
	// 暂时返回所有部门
	departments, err := s.deptRepo.ListAll()
	if err != nil {
		return nil, 0, err
	}
	
	// 简单的过滤逻辑
	var filtered []*entity.Department
	for _, dept := range departments {
		if query == "" || strings.Contains(dept.Name, query) || strings.Contains(dept.Code, query) {
			filtered = append(filtered, dept)
		}
	}
	
	total := int64(len(filtered))
	
	// 简单的分页逻辑
	start := (page - 1) * size
	end := start + size
	if start > len(filtered) {
		return []*entity.Department{}, total, nil
	}
	if end > len(filtered) {
		end = len(filtered)
	}
	
	return filtered[start:end], total, nil
}

// GetChildren 获取子部门
func (s *DepartmentService) GetChildren(parentID uint) ([]*entity.Department, error) {
	return s.deptRepo.List(parentID)
}

// GetTree 获取部门树
func (s *DepartmentService) GetTree() ([]*entity.Department, error) {
	// 获取所有部门，按层级和排序排列
	allDepts, err := s.deptRepo.ListAll()
	if err != nil {
		return nil, err
	}
	
	// 简单返回所有部门，前端可以根据ParentID构建树形结构
	return allDepts, nil
}

// GetTotalCount 获取部门总数
func (s *DepartmentService) GetTotalCount() (int64, error) {
	return s.deptRepo.GetTotalCount()
}