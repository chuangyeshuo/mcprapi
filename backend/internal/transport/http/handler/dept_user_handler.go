package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mcprapi/backend/internal/domain/service"
	"mcprapi/backend/internal/transport/http/dto"
	"mcprapi/backend/internal/transport/middleware"
)

// DeptUserHandler 部门用户管理处理器
type DeptUserHandler struct {
	userService           *service.UserService
	deptPermissionService *service.DeptPermissionService
}

// NewDeptUserHandler 创建部门用户管理处理器
func NewDeptUserHandler(
	userService *service.UserService,
	deptPermissionService *service.DeptPermissionService,
) *DeptUserHandler {
	return &DeptUserHandler{
		userService:           userService,
		deptPermissionService: deptPermissionService,
	}
}

// Register 注册路由
func (h *DeptUserHandler) Register(r *gin.RouterGroup) {
	// 部门用户管理路由 - 使用 dept-management 前缀避免与单个部门路由冲突
	r.GET("/dept-management/:dept_id/users", h.ListDeptUsers)
	r.POST("/dept-management/:dept_id/users", h.CreateDeptUser)
	r.PUT("/dept-management/:dept_id/users/:id", h.UpdateDeptUser)
	r.DELETE("/dept-management/:dept_id/users/:id", h.DeleteDeptUser)
	r.GET("/dept-management/:dept_id/users/:id", h.GetDeptUser)
	r.POST("/dept-management/:dept_id/users/:id/roles", h.AssignDeptUserRoles)
}

// checkDeptPermission 检查部门权限
func (h *DeptUserHandler) checkDeptPermission(c *gin.Context, deptID uint, action string) bool {
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

// ListDeptUsers 获取部门用户列表
func (h *DeptUserHandler) ListDeptUsers(c *gin.Context) {
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
			Message: "无权限访问该部门用户",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	username := c.DefaultQuery("username", "")

	req := &service.ListUserRequest{
		Page:     page,
		PageSize: limit,
		Query:    username,
		DeptID:   uint(deptID), // 添加部门ID过滤
	}

	resp, err := h.userService.ListUserByDept(req)
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
		Message: "获取部门用户列表成功",
		Data:    result,
	})
}

// CreateDeptUser 创建部门用户
func (h *DeptUserHandler) CreateDeptUser(c *gin.Context) {
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
			Message: "无权限在该部门创建用户",
		})
		return
	}

	var req service.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	// 强制设置部门ID
	req.DeptID = uint(deptID)

	user, err := h.userService.CreateUser(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "创建部门用户成功",
		Data:    user,
	})
}

// UpdateDeptUser 更新部门用户
func (h *DeptUserHandler) UpdateDeptUser(c *gin.Context) {
	deptID, err := strconv.ParseUint(c.Param("dept_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的部门ID",
		})
		return
	}

	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的用户ID",
		})
		return
	}

	// 检查权限
	if !h.checkDeptPermission(c, uint(deptID), "write") {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "无权限修改该部门用户",
		})
		return
	}

	// 检查用户是否属于该部门
	user, err := h.userService.GetUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Code:    dto.CodeNotFound,
			Message: "用户不存在",
		})
		return
	}

	if user.DeptID != uint(deptID) {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "用户不属于该部门",
		})
		return
	}

	var req service.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	req.ID = uint(userID)
	// 确保不能修改部门ID
	req.DeptID = uint(deptID)

	updatedUser, err := h.userService.UpdateUser(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "更新部门用户成功",
		Data:    updatedUser,
	})
}

// DeleteDeptUser 删除部门用户
func (h *DeptUserHandler) DeleteDeptUser(c *gin.Context) {
	deptID, err := strconv.ParseUint(c.Param("dept_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的部门ID",
		})
		return
	}

	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的用户ID",
		})
		return
	}

	// 检查权限
	if !h.checkDeptPermission(c, uint(deptID), "write") {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "无权限删除该部门用户",
		})
		return
	}

	// 检查用户是否属于该部门
	user, err := h.userService.GetUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Code:    dto.CodeNotFound,
			Message: "用户不存在",
		})
		return
	}

	if user.DeptID != uint(deptID) {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "用户不属于该部门",
		})
		return
	}

	err = h.userService.DeleteUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "删除部门用户成功",
	})
}

// GetDeptUser 获取部门用户
func (h *DeptUserHandler) GetDeptUser(c *gin.Context) {
	deptID, err := strconv.ParseUint(c.Param("dept_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的部门ID",
		})
		return
	}

	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的用户ID",
		})
		return
	}

	// 检查权限
	if !h.checkDeptPermission(c, uint(deptID), "read") {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "无权限查看该部门用户",
		})
		return
	}

	user, err := h.userService.GetUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Code:    dto.CodeNotFound,
			Message: "用户不存在",
		})
		return
	}

	if user.DeptID != uint(deptID) {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "用户不属于该部门",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取部门用户成功",
		Data:    user,
	})
}

// AssignDeptUserRoles 分配部门用户角色
func (h *DeptUserHandler) AssignDeptUserRoles(c *gin.Context) {
	deptID, err := strconv.ParseUint(c.Param("dept_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的部门ID",
		})
		return
	}

	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的用户ID",
		})
		return
	}

	// 检查权限
	if !h.checkDeptPermission(c, uint(deptID), "write") {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "无权限分配该部门用户角色",
		})
		return
	}

	// 检查用户是否属于该部门
	user, err := h.userService.GetUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Code:    dto.CodeNotFound,
			Message: "用户不存在",
		})
		return
	}

	if user.DeptID != uint(deptID) {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "用户不属于该部门",
		})
		return
	}

	var req service.AssignRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	req.UserID = uint(userID)

	// TODO: 需要验证角色是否属于该部门或者是系统角色
	err = h.userService.AssignRoles(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "分配部门用户角色成功",
	})
}