package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mcprapi/backend/internal/domain/service"
	"mcprapi/backend/internal/transport/http/dto"
	"mcprapi/backend/internal/transport/middleware"
)

// DeptAPIHandler 部门API管理处理器
type DeptAPIHandler struct {
	apiService            *service.APIService
	businessService       *service.BusinessService
	deptPermissionService *service.DeptPermissionService
}

// NewDeptAPIHandler 创建部门API管理处理器
func NewDeptAPIHandler(
	apiService *service.APIService,
	businessService *service.BusinessService,
	deptPermissionService *service.DeptPermissionService,
) *DeptAPIHandler {
	return &DeptAPIHandler{
		apiService:            apiService,
		businessService:       businessService,
		deptPermissionService: deptPermissionService,
	}
}

// Register 注册路由
func (h *DeptAPIHandler) Register(r *gin.RouterGroup) {
	// 部门级别的API管理路由 - 使用 dept-management 前缀避免与单个部门路由冲突
	r.GET("/dept-management/:dept_id/apis", h.ListDeptAPIs)
	r.POST("/dept-management/:dept_id/apis", h.CreateDeptAPI)
	r.PUT("/dept-management/:dept_id/apis/:id", h.UpdateDeptAPI)
	r.DELETE("/dept-management/:dept_id/apis/:id", h.DeleteDeptAPI)
	r.GET("/dept-management/:dept_id/apis/:id", h.GetDeptAPI)

	// 部门业务线下的API管理 - 使用不同的路由结构避免冲突
	r.GET("/dept-management/:dept_id/business-apis/:business_id", h.ListBusinessAPIs)
	r.POST("/dept-management/:dept_id/business-apis/:business_id", h.CreateBusinessAPI)
}

// checkDeptPermission 检查部门权限
func (h *DeptAPIHandler) checkDeptPermission(c *gin.Context, deptID uint, action string) bool {
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

// checkBusinessBelongsToDept 检查业务线是否属于部门
func (h *DeptAPIHandler) checkBusinessBelongsToDept(businessID, deptID uint) bool {
	business, err := h.businessService.GetBusiness(businessID)
	if err != nil {
		return false
	}
	return business.DeptID == deptID
}

// ListDeptAPIs 获取部门API列表
func (h *DeptAPIHandler) ListDeptAPIs(c *gin.Context) {
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
			Message: "无权限访问该部门API",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	name := c.DefaultQuery("name", "")
	businessID, _ := strconv.ParseUint(c.DefaultQuery("business_id", "0"), 10, 32)

	req := &service.ListAPIByDeptRequest{
		Page:       page,
		PageSize:   limit,
		Query:      name,
		DeptID:     uint(deptID),     // 添加部门ID过滤
		BusinessID: uint(businessID), // 可选的业务线ID过滤
	}

	resp, err := h.apiService.ListAPIByDept(req)
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
		Message: "获取部门API列表成功",
		Data:    result,
	})
}

// CreateDeptAPI 创建部门API
func (h *DeptAPIHandler) CreateDeptAPI(c *gin.Context) {
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
			Message: "无权限在该部门创建API",
		})
		return
	}

	var req service.CreateAPIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	// 检查业务线是否属于该部门
	if req.BusinessID > 0 && !h.checkBusinessBelongsToDept(req.BusinessID, uint(deptID)) {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "业务线不属于该部门",
		})
		return
	}

	api, err := h.apiService.CreateAPI(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "创建部门API成功",
		Data:    api,
	})
}

// UpdateDeptAPI 更新部门API
func (h *DeptAPIHandler) UpdateDeptAPI(c *gin.Context) {
	deptID, err := strconv.ParseUint(c.Param("dept_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的部门ID",
		})
		return
	}

	apiID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的API ID",
		})
		return
	}

	// 检查权限
	if !h.checkDeptPermission(c, uint(deptID), "write") {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "无权限修改该部门API",
		})
		return
	}

	// 检查API是否属于该部门
	api, err := h.apiService.GetAPI(uint(apiID))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Code:    dto.CodeNotFound,
			Message: "API不存在",
		})
		return
	}

	// 通过业务线检查API是否属于该部门
	if api.BusinessID > 0 && !h.checkBusinessBelongsToDept(api.BusinessID, uint(deptID)) {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "API不属于该部门",
		})
		return
	}

	var req service.UpdateAPIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	req.ID = uint(apiID)

	// 如果要修改业务线，检查新业务线是否属于该部门
	if req.BusinessID > 0 && !h.checkBusinessBelongsToDept(req.BusinessID, uint(deptID)) {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "目标业务线不属于该部门",
		})
		return
	}

	updatedAPI, err := h.apiService.UpdateAPI(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "更新部门API成功",
		Data:    updatedAPI,
	})
}

// DeleteDeptAPI 删除部门API
func (h *DeptAPIHandler) DeleteDeptAPI(c *gin.Context) {
	deptID, err := strconv.ParseUint(c.Param("dept_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的部门ID",
		})
		return
	}

	apiID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的API ID",
		})
		return
	}

	// 检查权限
	if !h.checkDeptPermission(c, uint(deptID), "write") {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "无权限删除该部门API",
		})
		return
	}

	// 检查API是否属于该部门
	api, err := h.apiService.GetAPI(uint(apiID))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Code:    dto.CodeNotFound,
			Message: "API不存在",
		})
		return
	}

	// 通过业务线检查API是否属于该部门
	if api.BusinessID > 0 && !h.checkBusinessBelongsToDept(api.BusinessID, uint(deptID)) {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "API不属于该部门",
		})
		return
	}

	err = h.apiService.DeleteAPI(uint(apiID))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "删除部门API成功",
	})
}

// GetDeptAPI 获取部门API
func (h *DeptAPIHandler) GetDeptAPI(c *gin.Context) {
	deptID, err := strconv.ParseUint(c.Param("dept_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的部门ID",
		})
		return
	}

	apiID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的API ID",
		})
		return
	}

	// 检查权限
	if !h.checkDeptPermission(c, uint(deptID), "read") {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "无权限查看该部门API",
		})
		return
	}

	api, err := h.apiService.GetAPI(uint(apiID))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Code:    dto.CodeNotFound,
			Message: "API不存在",
		})
		return
	}

	// 通过业务线检查API是否属于该部门
	if api.BusinessID > 0 && !h.checkBusinessBelongsToDept(api.BusinessID, uint(deptID)) {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "API不属于该部门",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取部门API成功",
		Data:    api,
	})
}

// ListBusinessAPIs 获取业务线API列表
func (h *DeptAPIHandler) ListBusinessAPIs(c *gin.Context) {
	deptID, err := strconv.ParseUint(c.Param("dept_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的部门ID",
		})
		return
	}

	businessID, err := strconv.ParseUint(c.Param("business_id"), 10, 32)
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
			Message: "无权限访问该部门API",
		})
		return
	}

	// 检查业务线是否属于该部门
	if !h.checkBusinessBelongsToDept(uint(businessID), uint(deptID)) {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "业务线不属于该部门",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	name := c.DefaultQuery("name", "")

	req := &service.ListAPIByDeptRequest{
		Page:       page,
		PageSize:   limit,
		Query:      name,
		BusinessID: uint(businessID),
	}

	resp, err := h.apiService.ListAPIByDept(req)
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
		Message: "获取业务线API列表成功",
		Data:    result,
	})
}

// CreateBusinessAPI 在业务线下创建API
func (h *DeptAPIHandler) CreateBusinessAPI(c *gin.Context) {
	deptID, err := strconv.ParseUint(c.Param("dept_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的部门ID",
		})
		return
	}

	businessID, err := strconv.ParseUint(c.Param("business_id"), 10, 32)
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
			Message: "无权限在该部门创建API",
		})
		return
	}

	// 检查业务线是否属于该部门
	if !h.checkBusinessBelongsToDept(uint(businessID), uint(deptID)) {
		c.JSON(http.StatusForbidden, dto.Response{
			Code:    dto.CodeForbidden,
			Message: "业务线不属于该部门",
		})
		return
	}

	var req service.CreateAPIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	// 强制设置业务线ID
	req.BusinessID = uint(businessID)

	api, err := h.apiService.CreateAPI(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "在业务线下创建API成功",
		Data:    api,
	})
}
