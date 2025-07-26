package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mcprapi/backend/internal/domain/entity"
	"mcprapi/backend/internal/domain/service"
	"mcprapi/backend/internal/transport/http/dto"
	"mcprapi/backend/internal/transport/middleware"
)

// BusinessHandler 业务线处理器
type BusinessHandler struct {
	businessService *service.BusinessService
	authService     *service.AuthService
}

// NewBusinessHandler 创建业务线处理器
func NewBusinessHandler(businessService *service.BusinessService, authService *service.AuthService) *BusinessHandler {
	return &BusinessHandler{
		businessService: businessService,
		authService:     authService,
	}
}

// Register 注册路由
func (h *BusinessHandler) Register(router *gin.RouterGroup) {
	businessRouter := router.Group("/business")
	{
		// 创建业务线
		businessRouter.POST("", h.Create)
		// 更新业务线
		businessRouter.PUT("/:id", h.Update)
		// 删除业务线
		businessRouter.DELETE("/:id", h.Delete)
		// 获取业务线
		businessRouter.GET("/:id", h.Get)
		// 注意：/list 路由已在基础认证路由组中注册，这里不再重复注册
		// 注意：/department/:id 路由已在基础认证路由组中注册，这里不再重复注册
		// 注意：/all 路由已在基础认证路由组中注册，这里不再重复注册
	}
}

// Create 创建业务线
func (h *BusinessHandler) Create(c *gin.Context) {
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
	var req service.CreateBusinessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的请求参数",
			Data:    err.Error(),
		})
		return
	}

	// 创建业务线
	business, err := h.businessService.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: "创建业务线失败",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "创建业务线成功",
		Data:    business,
	})
}

// Update 更新业务线
func (h *BusinessHandler) Update(c *gin.Context) {
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

	// 获取业务线ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的业务线ID",
			Data:    err.Error(),
		})
		return
	}

	// 解析请求
	var req service.UpdateBusinessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的请求参数",
			Data:    err.Error(),
		})
		return
	}

	// 设置业务线ID
	req.ID = uint(id)

	// 更新业务线
	business, err := h.businessService.Update(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: "更新业务线失败",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "更新业务线成功",
		Data:    business,
	})
}

// Delete 删除业务线
func (h *BusinessHandler) Delete(c *gin.Context) {
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

	// 获取业务线ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的业务线ID",
			Data:    err.Error(),
		})
		return
	}

	// 删除业务线
	if err := h.businessService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: "删除业务线失败",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "删除业务线成功",
	})
}

// Get 获取业务线
func (h *BusinessHandler) Get(c *gin.Context) {
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

	// 获取业务线ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的业务线ID",
			Data:    err.Error(),
		})
		return
	}

	// 获取业务线
	business, err := h.businessService.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: "获取业务线失败",
			Data:    err.Error(),
		})
		return
	}

	if business == nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Code:    dto.CodeNotFound,
			Message: "业务线不存在",
		})
		return
	}

	// 非管理员用户只能访问自己部门的业务线
	if !hasPermission.IsAdmin && business.DeptID != userDeptID {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "无权访问其他部门的业务线",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取业务线成功",
		Data:    business,
	})
}

// List 获取业务线列表
func (h *BusinessHandler) List(c *gin.Context) {
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

	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	query := c.DefaultQuery("query", "")
	departmentIDStr := c.DefaultQuery("dept_id", "0")
	departmentID, _ := strconv.ParseUint(departmentIDStr, 10, 64)

	// 检查用户权限
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

	// 根据用户角色和部门限制数据访问
	finalDepartmentID := uint(departmentID)
	if !hasPermission.IsAdmin {
		// 非管理员用户只能查看自己部门的业务线，直接使用token中的部门ID
		finalDepartmentID = userDeptID
		// 如果用户没有部门，返回空列表
		if userDeptID == 0 {
			result := map[string]interface{}{
				"items": []*entity.Business{},
				"total": int64(0),
			}
			c.JSON(http.StatusOK, dto.Response{
				Code:    dto.CodeSuccess,
				Message: "获取业务线列表成功",
				Data:    result,
			})
			return
		}
	}

	// 获取业务线列表
	businesses, total, err := h.businessService.List(page, limit, query, finalDepartmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: "获取业务线列表失败",
			Data:    err.Error(),
		})
		return
	}

	// 构造前端期望的响应格式
	result := map[string]interface{}{
		"items": businesses,
		"total": total,
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取业务线列表成功",
		Data:    result,
	})
}

// GetByDepartment 获取部门业务线
func (h *BusinessHandler) GetByDepartment(c *gin.Context) {
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

	// 检查用户权限
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

	// 非管理员用户只能访问自己部门的业务线
	if !hasPermission.IsAdmin && uint(id) != userDeptID {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "无权访问其他部门的业务线",
		})
		return
	}

	// 获取部门业务线
	businesses, err := h.businessService.GetByDepartment(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: "获取部门业务线失败",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取部门业务线成功",
		Data:    businesses,
	})
}

// GetAll 获取所有业务线
func (h *BusinessHandler) GetAll(c *gin.Context) {
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

	var businesses []*entity.Business
	if hasPermission.IsAdmin {
		// 管理员可以获取所有业务线
		businesses, err = h.businessService.GetAll()
	} else if userDeptID > 0 {
		// 非管理员用户只能获取自己部门的业务线，直接使用token中的部门ID
		businesses, err = h.businessService.GetByDepartment(userDeptID)
	} else {
		// 没有部门的用户返回空列表
		businesses = []*entity.Business{}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: "获取所有业务线失败",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取所有业务线成功",
		Data:    businesses,
	})
}
