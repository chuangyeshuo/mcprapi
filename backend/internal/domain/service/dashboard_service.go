package service

import (
	"mcprapi/backend/internal/domain/repository"
)

// DashboardService 仪表盘服务
type DashboardService struct {
	apiService        *APIService
	businessService   *BusinessService
	departmentService *DepartmentService
	userService       *UserService
	apiRepo           repository.APIRepository
}

// NewDashboardService 创建仪表盘服务
func NewDashboardService(
	apiService *APIService,
	businessService *BusinessService,
	departmentService *DepartmentService,
	userService *UserService,
	apiRepo repository.APIRepository,
) *DashboardService {
	return &DashboardService{
		apiService:        apiService,
		businessService:   businessService,
		departmentService: departmentService,
		userService:       userService,
		apiRepo:           apiRepo,
	}
}

// DashboardStats 仪表盘统计数据
type DashboardStats struct {
	APICount        int64 `json:"api_count"`
	BusinessCount   int64 `json:"business_count"`
	DepartmentCount int64 `json:"department_count"`
	UserCount       int64 `json:"user_count"`
}

// GetDashboardStats 获取仪表盘统计数据
func (s *DashboardService) GetDashboardStats() (*DashboardStats, error) {
	// 获取API总数
	apiCount, err := s.apiService.GetTotalCount()
	if err != nil {
		return nil, err
	}

	// 获取业务线总数
	businessCount, err := s.businessService.GetTotalCount()
	if err != nil {
		return nil, err
	}

	// 获取部门总数
	departmentCount, err := s.departmentService.GetTotalCount()
	if err != nil {
		return nil, err
	}

	// 获取用户总数
	userCount, err := s.userService.GetTotalCount()
	if err != nil {
		return nil, err
	}

	return &DashboardStats{
		APICount:        apiCount,
		BusinessCount:   businessCount,
		DepartmentCount: departmentCount,
		UserCount:       userCount,
	}, nil
}

// GetAPICategories 获取API分类统计
func (s *DashboardService) GetAPICategories() ([]*repository.CategoryStatsItem, error) {
	return s.apiService.GetCategoryStats()
}

// BusinessAPIStatsItem 业务线API统计项
type BusinessAPIStatsItem struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

// GetBusinessAPIStats 获取业务线API统计
func (s *DashboardService) GetBusinessAPIStats() ([]*BusinessAPIStatsItem, error) {
	// 获取所有业务线
	businesses, err := s.businessService.ListAllBusiness()
	if err != nil {
		return nil, err
	}

	var stats []*BusinessAPIStatsItem
	for _, business := range businesses {
		// 获取该业务线下的API数量
		count, err := s.apiRepo.GetCountByBusiness(business.ID)
		if err != nil {
			continue // 跳过错误的业务线
		}

		stats = append(stats, &BusinessAPIStatsItem{
			Name:  business.Name,
			Value: count,
		})
	}

	return stats, nil
}

// GetBusinessAPIStatsByDept 获取指定部门的业务线API统计
func (s *DashboardService) GetBusinessAPIStatsByDept(deptID uint) ([]*BusinessAPIStatsItem, error) {
	// 获取指定部门的业务线
	businesses, err := s.businessService.ListBusiness(deptID)
	if err != nil {
		return nil, err
	}

	var stats []*BusinessAPIStatsItem
	for _, business := range businesses {
		// 获取该业务线下的API数量
		count, err := s.apiRepo.GetCountByBusiness(business.ID)
		if err != nil {
			continue // 跳过错误的业务线
		}

		stats = append(stats, &BusinessAPIStatsItem{
			Name:  business.Name,
			Value: count,
		})
	}

	return stats, nil
}

// GetDepartmentAPIStatsByDept 获取指定部门的API统计
func (s *DashboardService) GetDepartmentAPIStatsByDept(deptID uint) ([]*DepartmentAPIStatsItem, error) {
	// 获取指定部门信息
	dept, err := s.departmentService.GetDepartment(deptID)
	if err != nil {
		return nil, err
	}

	// 获取该部门的API数量
	apiCount, err := s.apiRepo.CountByDepartment(deptID)
	if err != nil {
		return nil, err
	}

	// 返回单个部门的统计
	stats := []*DepartmentAPIStatsItem{
		{
			Name:  dept.Name,
			Value: apiCount,
		},
	}

	return stats, nil
}

// DepartmentAPIStatsItem 部门API统计项
type DepartmentAPIStatsItem struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

// GetDepartmentAPIStats 获取部门API统计
func (s *DashboardService) GetDepartmentAPIStats() ([]*DepartmentAPIStatsItem, error) {
	// 获取所有部门
	departments, err := s.departmentService.ListAllDepartment()
	if err != nil {
		return nil, err
	}

	var stats []*DepartmentAPIStatsItem
	for _, dept := range departments {
		// 获取该部门下的业务线
		businesses, err := s.businessService.ListBusiness(dept.ID)
		if err != nil {
			continue // 跳过错误的部门
		}

		var apiCount int64
		for _, business := range businesses {
			// 获取该业务线下的API数量
			count, err := s.apiRepo.GetCountByBusiness(business.ID)
			if err != nil {
				continue // 跳过错误的业务线
			}
			apiCount += count
		}

		stats = append(stats, &DepartmentAPIStatsItem{
			Name:  dept.Name,
			Value: apiCount,
		})
	}

	return stats, nil
}