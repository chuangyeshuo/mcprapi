package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mcprapi/backend/internal/domain/service"
	"mcprapi/backend/internal/transport/http/dto"
	"mcprapi/backend/internal/transport/middleware"
)

// DeptPermissionHandler 部门权限管理处理器
type DeptPermissionHandler struct {
	deptPermissionService *service.DeptPermissionService
}

// NewDeptPermissionHandler 创建部门权限管理处理器
func NewDeptPermissionHandler(deptPermissionService *service.DeptPermissionService) *DeptPermissionHandler {
	return &DeptPermissionHandler{
		deptPermissionService: deptPermissionService,
	}
}

// Register 注册路由
func (h *DeptPermissionHandler) Register(r *gin.RouterGroup) {
	// 部门权限管理路由
	r.POST("/dept-permission/grant-admin", h.GrantDeptAdmin)
	r.POST("/dept-permission/revoke-admin", h.RevokeDeptAdmin)
	r.POST("/dept-permission/check", h.CheckDeptPermission)
	r.GET("/departments/:id/admins", h.GetDeptAdmins)
	r.POST("/system/init-admin", h.InitSystemAdmin)
}

// GrantDeptAdmin 授予部门管理员权限
func (h *DeptPermissionHandler) GrantDeptAdmin(c *gin.Context) {
	var req service.GrantDeptAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	// 检查当前用户是否有权限执行此操作
	currentUserID := middleware.GetCurrentUser(c)
	if currentUserID == 0 {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Code:    dto.CodeUnauthorized,
			Message: "未认证的用户",
		})
		return
	}

	// 只有系统管理员可以授予部门管理员权限
	checkReq := &service.CheckDeptPermissionRequest{
		UserID: currentUserID,
		DeptID: req.DeptID,
		Action: "admin",
	}
	hasPermission, err := h.deptPermissionService.CheckDeptPermission(checkReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: "检查权限失败: " + err.Error(),
		})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "无权限执行此操作",
		})
		return
	}

	if err := h.deptPermissionService.GrantDeptAdmin(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "授予部门管理员权限成功",
	})
}

// RevokeDeptAdmin 撤销部门管理员权限
func (h *DeptPermissionHandler) RevokeDeptAdmin(c *gin.Context) {
	var req service.RevokeDeptAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	// 检查当前用户是否有权限执行此操作
	currentUserID := middleware.GetCurrentUser(c)
	if currentUserID == 0 {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Code:    dto.CodeUnauthorized,
			Message: "未认证的用户",
		})
		return
	}

	// 只有系统管理员可以撤销部门管理员权限
	checkReq := &service.CheckDeptPermissionRequest{
		UserID: currentUserID,
		DeptID: req.DeptID,
		Action: "admin",
	}
	hasPermission, err := h.deptPermissionService.CheckDeptPermission(checkReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: "检查权限失败: " + err.Error(),
		})
		return
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "无权限执行此操作",
		})
		return
	}

	if err := h.deptPermissionService.RevokeDeptAdmin(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "撤销部门管理员权限成功",
	})
}

// CheckDeptPermission 检查部门权限
func (h *DeptPermissionHandler) CheckDeptPermission(c *gin.Context) {
	var req service.CheckDeptPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	hasPermission, err := h.deptPermissionService.CheckDeptPermission(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "检查部门权限成功",
		Data: map[string]interface{}{
			"has_permission": hasPermission,
		},
	})
}

// GetDeptAdmins 获取部门管理员列表
func (h *DeptPermissionHandler) GetDeptAdmins(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的部门ID",
		})
		return
	}

	admins, err := h.deptPermissionService.GetDeptAdmins(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取部门管理员列表成功",
		Data:    admins,
	})
}

// InitSystemAdmin 初始化系统管理员
func (h *DeptPermissionHandler) InitSystemAdmin(c *gin.Context) {
	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	if err := h.deptPermissionService.InitSystemAdmin(req.UserID); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "初始化系统管理员成功",
	})
}