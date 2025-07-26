package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mcprapi/backend/internal/domain/service"
	"mcprapi/backend/internal/transport/http/dto"
	"mcprapi/backend/internal/transport/middleware"
)

// RoleHandler 角色处理器
type RoleHandler struct {
	roleService *service.RoleService
}

// NewRoleHandler 创建角色处理器
func NewRoleHandler(roleService *service.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// Register 注册路由
func (h *RoleHandler) Register(router *gin.RouterGroup) {
	roleRouter := router.Group("/role")
	{
		// 创建角色
		roleRouter.POST("", h.Create)
		// 更新角色
		roleRouter.PUT("/:id", h.Update)
		// 删除角色
		roleRouter.DELETE("/:id", h.Delete)
		// 获取角色
		roleRouter.GET("/:id", h.Get)
		// 获取角色列表
		roleRouter.GET("/list", h.List)
		// 更新角色权限
		roleRouter.PUT("/:id/permissions", h.UpdatePermissions)
		// 获取角色权限
		roleRouter.GET("/:id/permissions", h.GetPermissions)
		// 通过API ID更新角色权限
		roleRouter.PUT("/:id/api-permissions", h.UpdateAPIPermissions)
		// 获取角色的API权限ID列表
		roleRouter.GET("/:id/api-permissions", h.GetAPIPermissions)
		// 获取所有角色
		roleRouter.GET("/all", h.GetAll)
		// 获取用户可访问的角色
		roleRouter.GET("/user-accessible", h.GetUserAccessible)
	}
}

// Create 创建角色
func (h *RoleHandler) Create(c *gin.Context) {
	var req service.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	// 设置当前用户ID用于权限验证
	req.CurrentUserID = middleware.GetCurrentUser(c)

	role, err := h.roleService.CreateRole(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "创建角色成功",
		Data:    role,
	})
}

// Update 更新角色
func (h *RoleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的角色ID",
		})
		return
	}

	var req service.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	req.ID = uint(id)
	role, err := h.roleService.UpdateRole(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "更新角色成功",
		Data:    role,
	})
}

// Delete 删除角色
func (h *RoleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的角色ID",
		})
		return
	}

	err = h.roleService.DeleteRole(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "删除角色成功",
	})
}

// Get 获取角色
func (h *RoleHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的角色ID",
		})
		return
	}

	role, err := h.roleService.GetRole(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Code:    dto.CodeNotFound,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取角色成功",
		Data:    role,
	})
}

// List 获取角色列表
func (h *RoleHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	query := c.DefaultQuery("query", "")
	deptID, _ := strconv.ParseUint(c.DefaultQuery("dept_id", "0"), 10, 32)

	req := &service.ListRoleRequest{
		Page:     page,
		PageSize: limit,
		Query:    query,
		DeptID:   uint(deptID),
	}

	resp, err := h.roleService.ListRole(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: err.Error(),
		})
		return
	}

	// 构造前端期望的响应格式
	result := map[string]interface{}{
		"items": resp.Items,
		"total": resp.Total,
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取角色列表成功",
		Data:    result,
	})
}

// UpdatePermissions 更新角色权限
func (h *RoleHandler) UpdatePermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的角色ID",
		})
		return
	}

	var req service.UpdateRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	req.RoleID = uint(id)
	err = h.roleService.UpdateRolePermissions(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "更新角色权限成功",
	})
}

// GetPermissions 获取角色权限
func (h *RoleHandler) GetPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的角色ID",
		})
		return
	}

	permissions, err := h.roleService.GetRolePermissions(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取角色权限成功",
		Data:    permissions,
	})
}

// GetAll 获取所有角色
func (h *RoleHandler) GetAll(c *gin.Context) {
	roles, err := h.roleService.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取所有角色成功",
		Data:    roles,
	})
}

// GetUserAccessible 获取用户可访问的角色
func (h *RoleHandler) GetUserAccessible(c *gin.Context) {
	// 从JWT中获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Code:    dto.CodeUnauthorized,
			Message: "未授权",
		})
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "用户ID格式错误",
		})
		return
	}

	roles, err := h.roleService.GetUserAccessibleRoles(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取用户可访问角色成功",
		Data:    roles,
	})
}

// UpdateAPIPermissions 通过API ID更新角色权限
func (h *RoleHandler) UpdateAPIPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的角色ID",
		})
		return
	}

	// 创建临时结构体来接收请求体，避免RoleID验证问题
	var reqBody struct {
		APIIDs []uint `json:"api_ids" binding:"required"`
		DeptID *uint  `json:"dept_id,omitempty"` // 可选的部门ID
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	// 获取当前用户的部门ID
	currentDeptID := middleware.GetCurrentDeptID(c)
	
	// 构造完整的请求结构体
	req := service.UpdateRolePermissionsByAPIIDsRequest{
		RoleID: uint(id),
		APIIDs: reqBody.APIIDs,
	}

	// 确定使用的部门ID：优先使用请求中的部门ID，否则使用当前用户的部门ID
	if reqBody.DeptID != nil && *reqBody.DeptID > 0 {
		req.DeptID = *reqBody.DeptID
	} else if currentDeptID > 0 {
		req.DeptID = currentDeptID
	}
	// 如果都为0或nil，则DeptID保持为0，service层会使用默认值"*"

	err = h.roleService.UpdateRolePermissionsByAPIIDs(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "更新角色API权限成功",
	})
}

// GetAPIPermissions 获取角色的API权限ID列表
func (h *RoleHandler) GetAPIPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的角色ID",
		})
		return
	}

	apiIDs, err := h.roleService.GetRoleAPIIDs(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取角色API权限成功",
		Data:    apiIDs,
	})
}
