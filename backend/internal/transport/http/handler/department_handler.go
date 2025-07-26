package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mcprapi/backend/internal/domain/service"
	"mcprapi/backend/internal/transport/http/dto"
	"mcprapi/backend/internal/transport/middleware"
)

// DepartmentHandler 部门处理器
type DepartmentHandler struct {
	departmentService *service.DepartmentService
	authService       *service.AuthService
}

// NewDepartmentHandler 创建部门处理器
func NewDepartmentHandler(departmentService *service.DepartmentService, authService *service.AuthService) *DepartmentHandler {
	return &DepartmentHandler{
		departmentService: departmentService,
		authService:       authService,
	}
}

// Register 注册路由
func (h *DepartmentHandler) Register(router *gin.RouterGroup) {
	departmentRouter := router.Group("/department")
	{
		// 创建部门
		departmentRouter.POST("", h.Create)
		// 更新部门
		departmentRouter.PUT("/:id", h.Update)
		// 删除部门
		departmentRouter.DELETE("/:id", h.Delete)
		// 获取部门
		departmentRouter.GET("/:id", h.Get)
		// 注意：/list 路由已在基础认证路由组中注册，这里不再重复注册
		// 获取子部门
		departmentRouter.GET("/:id/children", h.GetChildren)
		// 获取部门树
		departmentRouter.GET("/tree", h.GetTree)
	}
}

// Create 创建部门
func (h *DepartmentHandler) Create(c *gin.Context) {
	// 获取当前用户
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Code:    dto.CodeUnauthorized,
			Message: "未认证的用户",
		})
		return
	}

	// 检查权限
	hasPermission, err := h.authService.CheckPermission(&service.CheckPermissionRequest{
		UserID:  fmt.Sprintf("%d", userID),
		APIPath: c.Request.URL.Path,
		Method:  c.Request.Method,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: "权限检查失败",
			Data:    err.Error(),
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

	// 解析请求
	var req service.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的请求参数",
			Data:    err.Error(),
		})
		return
	}

	// 创建部门
	department, err := h.departmentService.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: "创建部门失败",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "创建部门成功",
		Data:    department,
	})
}

// Update 更新部门
func (h *DepartmentHandler) Update(c *gin.Context) {
	// 获取当前用户
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Code:    dto.CodeUnauthorized,
			Message: "未认证的用户",
		})
		return
	}

	// 检查权限
	hasPermission, err := h.authService.CheckPermission(&service.CheckPermissionRequest{
		UserID:  fmt.Sprintf("%d", userID),
		APIPath: c.Request.URL.Path,
		Method:  c.Request.Method,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: "权限检查失败",
			Data:    err.Error(),
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

	// 获取部门ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的部门ID",
			Data:    err.Error(),
		})
		return
	}

	// 解析请求
	var req service.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的请求参数",
			Data:    err.Error(),
		})
		return
	}

	// 设置部门ID
	req.ID = uint(id)

	// 更新部门
	department, err := h.departmentService.Update(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: "更新部门失败",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "更新部门成功",
		Data:    department,
	})
}

// Delete 删除部门
func (h *DepartmentHandler) Delete(c *gin.Context) {
	// 获取当前用户
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Code:    dto.CodeUnauthorized,
			Message: "未认证的用户",
		})
		return
	}

	// 检查权限
	hasPermission, err := h.authService.CheckPermission(&service.CheckPermissionRequest{
		UserID:  fmt.Sprintf("%d", userID),
		APIPath: c.Request.URL.Path,
		Method:  c.Request.Method,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: "权限检查失败",
			Data:    err.Error(),
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

	// 获取部门ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的部门ID",
			Data:    err.Error(),
		})
		return
	}

	// 删除部门
	if err := h.departmentService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: "删除部门失败",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "删除部门成功",
	})
}

// Get 获取部门
func (h *DepartmentHandler) Get(c *gin.Context) {
	// 获取当前用户
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Code:    dto.CodeUnauthorized,
			Message: "未认证的用户",
		})
		return
	}

	// 检查权限
	hasPermission, err := h.authService.CheckPermission(&service.CheckPermissionRequest{
		UserID:  fmt.Sprintf("%d", userID),
		APIPath: c.Request.URL.Path,
		Method:  c.Request.Method,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: "权限检查失败",
			Data:    err.Error(),
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

	// 获取部门ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的部门ID",
			Data:    err.Error(),
		})
		return
	}

	// 获取部门
	department, err := h.departmentService.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: "获取部门失败",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取部门成功",
		Data:    department,
	})
}

// List 获取部门列表
func (h *DepartmentHandler) List(c *gin.Context) {
	// 获取当前用户
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Code:    dto.CodeUnauthorized,
			Message: "未认证的用户",
		})
		return
	}

	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	query := c.DefaultQuery("query", "")

	// 获取部门列表
	departments, total, err := h.departmentService.List(page, limit, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: "获取部门列表失败",
			Data:    err.Error(),
		})
		return
	}

	// 构造前端期望的响应格式
	result := map[string]interface{}{
		"items": departments,
		"total": total,
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取部门列表成功",
		Data:    result,
	})
}

// GetChildren 获取子部门
func (h *DepartmentHandler) GetChildren(c *gin.Context) {
	// 获取当前用户
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Code:    dto.CodeUnauthorized,
			Message: "未认证的用户",
		})
		return
	}

	// 检查权限
	hasPermission, err := h.authService.CheckPermission(&service.CheckPermissionRequest{
		UserID:  fmt.Sprintf("%d", userID),
		APIPath: c.Request.URL.Path,
		Method:  c.Request.Method,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: "权限检查失败",
			Data:    err.Error(),
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

	// 获取部门ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的部门ID",
			Data:    err.Error(),
		})
		return
	}

	// 获取子部门
	children, err := h.departmentService.GetChildren(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: "获取子部门失败",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取子部门成功",
		Data:    children,
	})
}

// GetTree 获取部门树
func (h *DepartmentHandler) GetTree(c *gin.Context) {
	// 获取当前用户
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Code:    dto.CodeUnauthorized,
			Message: "未认证的用户",
		})
		return
	}

	// 检查权限
	hasPermission, err := h.authService.CheckPermission(&service.CheckPermissionRequest{
		UserID:  fmt.Sprintf("%d", userID),
		APIPath: c.Request.URL.Path,
		Method:  c.Request.Method,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: "权限检查失败",
			Data:    err.Error(),
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

	// 获取部门树
	tree, err := h.departmentService.GetTree()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: "获取部门树失败",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取部门树成功",
		Data:    tree,
	})
}
