package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mcprapi/backend/internal/domain/service"
	"mcprapi/backend/internal/transport/http/dto"
	"mcprapi/backend/internal/transport/middleware"
)

// DeptRoleHandler 部门角色管理处理器
type DeptRoleHandler struct {
	roleService           *service.RoleService
	deptPermissionService *service.DeptPermissionService
}

// NewDeptRoleHandler 创建部门角色管理处理器
func NewDeptRoleHandler(
	roleService *service.RoleService,
	deptPermissionService *service.DeptPermissionService,
) *DeptRoleHandler {
	return &DeptRoleHandler{
		roleService:           roleService,
		deptPermissionService: deptPermissionService,
	}
}

// Register 注册路由
func (h *DeptRoleHandler) Register(r *gin.RouterGroup) {
	// 部门角色管理路由 - 使用 dept-management 前缀避免与单个部门路由冲突
	r.GET("/dept-management/:dept_id/roles", h.ListDeptRoles)
	r.POST("/dept-management/:dept_id/roles", h.CreateDeptRole)
	r.PUT("/dept-management/:dept_id/roles/:id", h.UpdateDeptRole)
	r.DELETE("/dept-management/:dept_id/roles/:id", h.DeleteDeptRole)
	r.GET("/dept-management/:dept_id/roles/:id", h.GetDeptRole)
	r.POST("/dept-management/:dept_id/roles/:id/permissions", h.AssignDeptRolePermissions)
	r.GET("/dept-management/:dept_id/roles/:id/permissions", h.GetDeptRolePermissions)
}

// checkDeptPermission 检查部门权限
func (h *DeptRoleHandler) checkDeptPermission(c *gin.Context, deptID uint, action string) bool {
	currentUserID := middleware.GetCurrentUser(c)
	if currentUserID == 0 {
		return false
	}

	checkReq := &service.CheckDeptPermissionRequest{
		UserID: currentUserID,
		DeptID: deptID,
		Action: action,
	}

	hasPermission, err := h.deptPermissionService.CheckDeptPermission(checkReq)
	if err != nil {
		return false
	}

	return hasPermission
}

// ListDeptRoles 获取部门角色列表
func (h *DeptRoleHandler) ListDeptRoles(c *gin.Context) {
	deptID, err := strconv.ParseUint(c.Param("dept_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的部门ID",
		})
		return
	}

	// 检查权限
	if !h.checkDeptPermission(c, uint(deptID), "read") {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "无权限访问该部门角色",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	name := c.DefaultQuery("name", "")

	req := &service.ListRoleRequest{
		Page:     page,
		PageSize: limit,
		Query:    name,
		DeptID:   uint(deptID), // 添加部门ID过滤
	}

	resp, err := h.roleService.ListRoleByDept(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: err.Error(),
		})
		return
	}

	result := map[string]interface{}{
		"items": resp.Items,
		"total": resp.Total,
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取部门角色列表成功",
		Data:    result,
	})
}

// CreateDeptRole 创建部门角色
func (h *DeptRoleHandler) CreateDeptRole(c *gin.Context) {
	deptID, err := strconv.ParseUint(c.Param("dept_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的部门ID",
		})
		return
	}

	// 检查权限
	if !h.checkDeptPermission(c, uint(deptID), "write") {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "无权限在该部门创建角色",
		})
		return
	}

	var req service.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	// 强制设置部门ID
	req.DeptID = uint(deptID)

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
		Message: "创建部门角色成功",
		Data:    role,
	})
}

// UpdateDeptRole 更新部门角色
func (h *DeptRoleHandler) UpdateDeptRole(c *gin.Context) {
	deptID, err := strconv.ParseUint(c.Param("dept_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的部门ID",
		})
		return
	}

	roleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的角色ID",
		})
		return
	}

	// 检查权限
	if !h.checkDeptPermission(c, uint(deptID), "write") {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "无权限修改该部门角色",
		})
		return
	}

	// 检查角色是否属于该部门
	role, err := h.roleService.GetRole(uint(roleID))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Code:    dto.CodeNotFound,
			Message: "角色不存在",
		})
		return
	}

	if role.DeptID != uint(deptID) {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "角色不属于该部门",
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

	req.ID = uint(roleID)
	// 确保不能修改部门ID
	req.DeptID = uint(deptID)

	updatedRole, err := h.roleService.UpdateRole(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "更新部门角色成功",
		Data:    updatedRole,
	})
}

// DeleteDeptRole 删除部门角色
func (h *DeptRoleHandler) DeleteDeptRole(c *gin.Context) {
	deptID, err := strconv.ParseUint(c.Param("dept_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的部门ID",
		})
		return
	}

	roleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的角色ID",
		})
		return
	}

	// 检查权限
	if !h.checkDeptPermission(c, uint(deptID), "write") {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "无权限删除该部门角色",
		})
		return
	}

	// 检查角色是否属于该部门
	role, err := h.roleService.GetRole(uint(roleID))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Code:    dto.CodeNotFound,
			Message: "角色不存在",
		})
		return
	}

	if role.DeptID != uint(deptID) {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "角色不属于该部门",
		})
		return
	}

	err = h.roleService.DeleteRole(uint(roleID))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "删除部门角色成功",
	})
}

// GetDeptRole 获取部门角色
func (h *DeptRoleHandler) GetDeptRole(c *gin.Context) {
	deptID, err := strconv.ParseUint(c.Param("dept_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的部门ID",
		})
		return
	}

	roleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的角色ID",
		})
		return
	}

	// 检查权限
	if !h.checkDeptPermission(c, uint(deptID), "read") {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "无权限查看该部门角色",
		})
		return
	}

	role, err := h.roleService.GetRole(uint(roleID))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Code:    dto.CodeNotFound,
			Message: "角色不存在",
		})
		return
	}

	if role.DeptID != uint(deptID) {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "角色不属于该部门",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取部门角色成功",
		Data:    role,
	})
}

// AssignDeptRolePermissions 分配部门角色权限
func (h *DeptRoleHandler) AssignDeptRolePermissions(c *gin.Context) {
	deptID, err := strconv.ParseUint(c.Param("dept_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的部门ID",
		})
		return
	}

	roleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的角色ID",
		})
		return
	}

	// 检查权限
	if !h.checkDeptPermission(c, uint(deptID), "write") {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "无权限分配该部门角色权限",
		})
		return
	}

	// 检查角色是否属于该部门
	role, err := h.roleService.GetRole(uint(roleID))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Code:    dto.CodeNotFound,
			Message: "角色不存在",
		})
		return
	}

	if role.DeptID != uint(deptID) {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "角色不属于该部门",
		})
		return
	}

	var req service.AssignRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	req.RoleID = uint(roleID)
	req.DeptID = uint(deptID) // 确保权限分配在部门范围内

	err = h.roleService.AssignRolePermissions(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "分配部门角色权限成功",
	})
}

// GetDeptRolePermissions 获取部门角色权限
func (h *DeptRoleHandler) GetDeptRolePermissions(c *gin.Context) {
	deptID, err := strconv.ParseUint(c.Param("dept_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的部门ID",
		})
		return
	}

	roleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的角色ID",
		})
		return
	}

	// 检查权限
	if !h.checkDeptPermission(c, uint(deptID), "read") {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "无权限查看该部门角色权限",
		})
		return
	}

	// 检查角色是否属于该部门
	role, err := h.roleService.GetRole(uint(roleID))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Code:    dto.CodeNotFound,
			Message: "角色不存在",
		})
		return
	}

	if role.DeptID != uint(deptID) {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "角色不属于该部门",
		})
		return
	}

	permissions, err := h.roleService.GetRolePermissions(uint(roleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取部门角色权限成功",
		Data:    permissions,
	})
}