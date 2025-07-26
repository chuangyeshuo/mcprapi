package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mcprapi/backend/internal/domain/service"
	"mcprapi/backend/internal/transport/http/dto"
	"mcprapi/backend/internal/transport/middleware"
)

// DeptBusinessHandler 部门业务线管理处理器
type DeptBusinessHandler struct {
	businessService       *service.BusinessService
	deptPermissionService *service.DeptPermissionService
}

// NewDeptBusinessHandler 创建部门业务线管理处理器
func NewDeptBusinessHandler(
	businessService *service.BusinessService,
	deptPermissionService *service.DeptPermissionService,
) *DeptBusinessHandler {
	return &DeptBusinessHandler{
		businessService:       businessService,
		deptPermissionService: deptPermissionService,
	}
}

// Register 注册路由
func (h *DeptBusinessHandler) Register(r *gin.RouterGroup) {
	// 部门业务线管理路由 - 使用 dept-management 前缀避免与单个部门路由冲突
	r.GET("/dept-management/:dept_id/businesses", h.ListDeptBusinesses)
	r.POST("/dept-management/:dept_id/businesses", h.CreateDeptBusiness)
	r.PUT("/dept-management/:dept_id/businesses/:id", h.UpdateDeptBusiness)
	r.DELETE("/dept-management/:dept_id/businesses/:id", h.DeleteDeptBusiness)
	r.GET("/dept-management/:dept_id/businesses/:id", h.GetDeptBusiness)
}

// checkDeptPermission 检查部门权限
func (h *DeptBusinessHandler) checkDeptPermission(c *gin.Context, deptID uint, action string) bool {
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

// ListDeptBusinesses 获取部门业务线列表
func (h *DeptBusinessHandler) ListDeptBusinesses(c *gin.Context) {
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
			Message: "无权限访问该部门业务线",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	name := c.DefaultQuery("name", "")

	req := &service.ListBusinessRequest{
		Page:     page,
		PageSize: limit,
		Query:    name,
		DeptID:   uint(deptID), // 添加部门ID过滤
	}

	resp, err := h.businessService.ListBusinessByDept(req)
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
		Message: "获取部门业务线列表成功",
		Data:    result,
	})
}

// CreateDeptBusiness 创建部门业务线
func (h *DeptBusinessHandler) CreateDeptBusiness(c *gin.Context) {
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
			Message: "无权限在该部门创建业务线",
		})
		return
	}

	var req service.CreateBusinessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	// 强制设置部门ID
	req.DeptID = uint(deptID)

	business, err := h.businessService.CreateBusiness(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "创建部门业务线成功",
		Data:    business,
	})
}

// UpdateDeptBusiness 更新部门业务线
func (h *DeptBusinessHandler) UpdateDeptBusiness(c *gin.Context) {
	deptID, err := strconv.ParseUint(c.Param("dept_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的部门ID",
		})
		return
	}

	businessID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的业务线ID",
		})
		return
	}

	// 检查权限
	if !h.checkDeptPermission(c, uint(deptID), "write") {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "无权限修改该部门业务线",
		})
		return
	}

	// 检查业务线是否属于该部门
	business, err := h.businessService.GetBusiness(uint(businessID))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Code:    dto.CodeNotFound,
			Message: "业务线不存在",
		})
		return
	}

	if business.DeptID != uint(deptID) {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "业务线不属于该部门",
		})
		return
	}

	var req service.UpdateBusinessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	req.ID = uint(businessID)
	// 确保不能修改部门ID
	req.DeptID = uint(deptID)

	updatedBusiness, err := h.businessService.UpdateBusiness(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "更新部门业务线成功",
		Data:    updatedBusiness,
	})
}

// DeleteDeptBusiness 删除部门业务线
func (h *DeptBusinessHandler) DeleteDeptBusiness(c *gin.Context) {
	deptID, err := strconv.ParseUint(c.Param("dept_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的部门ID",
		})
		return
	}

	businessID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的业务线ID",
		})
		return
	}

	// 检查权限
	if !h.checkDeptPermission(c, uint(deptID), "write") {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "无权限删除该部门业务线",
		})
		return
	}

	// 检查业务线是否属于该部门
	business, err := h.businessService.GetBusiness(uint(businessID))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Code:    dto.CodeNotFound,
			Message: "业务线不存在",
		})
		return
	}

	if business.DeptID != uint(deptID) {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "业务线不属于该部门",
		})
		return
	}

	err = h.businessService.DeleteBusiness(uint(businessID))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "删除部门业务线成功",
	})
}

// GetDeptBusiness 获取部门业务线
func (h *DeptBusinessHandler) GetDeptBusiness(c *gin.Context) {
	deptID, err := strconv.ParseUint(c.Param("dept_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的部门ID",
		})
		return
	}

	businessID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的业务线ID",
		})
		return
	}

	// 检查权限
	if !h.checkDeptPermission(c, uint(deptID), "read") {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "无权限查看该部门业务线",
		})
		return
	}

	business, err := h.businessService.GetBusiness(uint(businessID))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Code:    dto.CodeNotFound,
			Message: "业务线不存在",
		})
		return
	}

	if business.DeptID != uint(deptID) {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "业务线不属于该部门",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取部门业务线成功",
		Data:    business,
	})
}