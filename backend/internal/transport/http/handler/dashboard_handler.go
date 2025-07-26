package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"mcprapi/backend/internal/domain/service"
	"mcprapi/backend/internal/transport/http/dto"
	"mcprapi/backend/internal/transport/middleware"
)

// DashboardHandler 仪表盘处理器
type DashboardHandler struct {
	dashboardService *service.DashboardService
	authService      *service.AuthService
}

// NewDashboardHandler 创建仪表盘处理器
func NewDashboardHandler(dashboardService *service.DashboardService, authService *service.AuthService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
		authService:      authService,
	}
}

// RegisterRoutes 注册路由
func (h *DashboardHandler) RegisterRoutes(r *gin.RouterGroup) {
	dashboard := r.Group("/dashboard")
	{
		dashboard.GET("/stats", h.GetDashboardStats)
		dashboard.GET("/api-category-stats", h.GetAPICategories)
		dashboard.GET("/business-api-stats", h.GetBusinessAPIStats)
		dashboard.GET("/department-api-stats", h.GetDepartmentAPIStats)
	}
}

// GetDashboardStats 获取仪表盘统计数据
// @Summary 获取仪表盘统计数据
// @Description 获取系统整体统计数据，包括用户、部门、角色、API等数量
// @Tags 仪表盘
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.Response{data=service.DashboardStats} "获取成功"
// @Failure 500 {object} dto.Response "服务器错误"
// @Router /dashboard/stats [get]
func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.dashboardService.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取统计数据成功",
		Data:    stats,
	})
}

// GetAPICategories 获取API分类统计
func (h *DashboardHandler) GetAPICategories(c *gin.Context) {
	categories, err := h.dashboardService.GetAPICategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取API分类统计成功",
		Data:    categories,
	})
}

// GetBusinessAPIStats 获取业务线API统计
// @Summary 获取业务线API统计
// @Description 获取业务线API统计数据，管理员可查看所有业务线，普通用户只能查看自己部门的业务线
// @Tags 仪表盘
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.Response{data=[]service.BusinessAPIStatsItem} "获取成功"
// @Failure 401 {object} dto.Response "未认证"
// @Failure 403 {object} dto.Response "无权限"
// @Failure 500 {object} dto.Response "服务器错误"
// @Router /dashboard/business-api-stats [get]
func (h *DashboardHandler) GetBusinessAPIStats(c *gin.Context) {
	// 获取当前用户
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Code:    dto.CodeUnauthorized,
			Message: "未认证的用户",
		})
		return
	}

	// 直接从token中获取用户部门ID
	userDeptID := middleware.GetCurrentDeptID(c)

	// 检查用户权限
	hasPermission, err := h.authService.CheckPermission(&service.CheckPermissionRequest{
		UserID:  fmt.Sprintf("%d", userID),
		APIPath: c.Request.URL.Path,
		Method:  c.Request.Method,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: "权限检查失败",
		})
		return
	}
	if !hasPermission.Allowed {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "无权访问该资源",
		})
		return
	}

	var stats []*service.BusinessAPIStatsItem
	if hasPermission.IsAdmin {
		// 管理员可以看到所有业务线统计
		stats, err = h.dashboardService.GetBusinessAPIStats()
	} else if userDeptID > 0 {
		// 非管理员用户只能看到自己部门的业务线统计
		stats, err = h.dashboardService.GetBusinessAPIStatsByDept(userDeptID)
	} else {
		// 没有部门的用户返回空统计
		stats = []*service.BusinessAPIStatsItem{}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取业务线API统计成功",
		Data:    stats,
	})
}

// GetDepartmentAPIStats 获取部门API统计
func (h *DashboardHandler) GetDepartmentAPIStats(c *gin.Context) {
	// 获取当前用户
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Code:    dto.CodeUnauthorized,
			Message: "未认证的用户",
		})
		return
	}

	// 直接从token中获取用户部门ID
	userDeptID := middleware.GetCurrentDeptID(c)

	// 检查用户权限
	hasPermission, err := h.authService.CheckPermission(&service.CheckPermissionRequest{
		UserID:  fmt.Sprintf("%d", userID),
		APIPath: c.Request.URL.Path,
		Method:  c.Request.Method,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: "权限检查失败",
		})
		return
	}
	if !hasPermission.Allowed {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "无权访问该资源",
		})
		return
	}

	var stats []*service.DepartmentAPIStatsItem
	if hasPermission.IsAdmin {
		// 管理员可以看到所有部门统计
		stats, err = h.dashboardService.GetDepartmentAPIStats()
	} else if userDeptID > 0 {
		// 非管理员用户只能看到自己部门的统计
		stats, err = h.dashboardService.GetDepartmentAPIStatsByDept(userDeptID)
	} else {
		// 没有部门的用户返回空统计
		stats = []*service.DepartmentAPIStatsItem{}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取部门API统计成功",
		Data:    stats,
	})
}
