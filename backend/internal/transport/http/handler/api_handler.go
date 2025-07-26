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

// APIHandler API处理器
type APIHandler struct {
	apiService  *service.APIService
	authService *service.AuthService
}

// NewAPIHandler 创建API处理器
func NewAPIHandler(apiService *service.APIService, authService *service.AuthService) *APIHandler {
	return &APIHandler{
		apiService:  apiService,
		authService: authService,
	}
}

// Register 注册路由
func (h *APIHandler) Register(router *gin.RouterGroup) {
	// API路由
	apiRouter := router.Group("/api")
	{
		// 注意：以下路由已在基础认证路由组中注册，这里不再重复注册
		// - /list
		// - /:id
		// - /business/:id
		// - /category/list
		// - /category/:id
		// - POST "" (创建API)
		// - PUT "/:id" (更新API)
		// - POST "/check-permission" (权限校验)

		// 删除API
		apiRouter.DELETE("/:id", h.Delete)

		// API分类路由
		apiRouter.POST("/category", h.CreateCategory)
		apiRouter.PUT("/category/:id", h.UpdateCategory)
		apiRouter.DELETE("/category/:id", h.DeleteCategory)
	}
}

// CheckPermission 检查权限
func (h *APIHandler) CheckPermission(c *gin.Context) {
	// 从JWT token中获取当前用户ID
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Code:    dto.CodeUnauthorized,
			Message: "未认证的用户",
		})
		return
	}

	// 获取Token版本号并验证
	tokenVersion := middleware.GetCurrentTokenVersion(c)
	if err := h.authService.ValidateTokenVersion(userID, tokenVersion); err != nil {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Code:    dto.CodeUnauthorized,
			Message: err.Error(),
		})
		return
	}

	// 定义简化的请求结构体，只需要api_path和method
	var reqBody struct {
		APIPath string `json:"api_path" binding:"required"`
		Method  string `json:"method" binding:"required"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	// 构造完整的权限检查请求，自动填充用户ID
	req := service.CheckPermissionRequest{
		UserID:  fmt.Sprintf("%d", userID),
		APIPath: reqBody.APIPath,
		Method:  reqBody.Method,
	}

	resp, err := h.authService.CheckPermission(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "检查权限成功",
		Data:    resp,
	})
}

// Create 创建API
func (h *APIHandler) Create(c *gin.Context) {
	var req service.CreateAPIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
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
		Message: "创建API成功",
		Data:    api,
	})
}

// Update 更新API
func (h *APIHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的API ID",
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

	req.ID = uint(id)
	api, err := h.apiService.UpdateAPI(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "更新API成功",
		Data:    api,
	})
}

// Delete 删除API
func (h *APIHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的API ID",
		})
		return
	}

	err = h.apiService.DeleteAPI(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "删除API成功",
	})
}

// Get 获取API
func (h *APIHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的API ID",
		})
		return
	}

	api, err := h.apiService.GetAPI(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Code:    dto.CodeNotFound,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取API成功",
		Data:    api,
	})
}

// List 获取API列表
func (h *APIHandler) List(c *gin.Context) {
	// 获取当前用户ID
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Code:    dto.CodeUnauthorized,
			Message: "未认证的用户",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	query := c.DefaultQuery("query", "")
	businessID, _ := strconv.ParseUint(c.DefaultQuery("business_id", "0"), 10, 32)
	category := c.DefaultQuery("category", "")

	req := &service.ListAPIRequest{
		Page:       page,
		PageSize:   limit,
		Query:      query,
		BusinessID: uint(businessID),
		Category:   category,
		UserID:     userID,
	}

	resp, err := h.apiService.ListAPI(req)
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
		Message: "获取API列表成功",
		Data:    result,
	})
}

// ListByBusiness 获取业务线下的API列表
func (h *APIHandler) ListByBusiness(c *gin.Context) {
	businessID, err := strconv.ParseUint(c.Param("business_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的业务线ID",
		})
		return
	}

	apis, err := h.apiService.ListAPIByBusiness(uint(businessID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取业务线API列表成功",
		Data:    apis,
	})
}

// ListCategories 获取API分类列表
func (h *APIHandler) ListCategories(c *gin.Context) {
	categories, err := h.apiService.ListAPICategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取API分类列表成功",
		Data:    categories,
	})
}

// CreateCategory 创建API分类
func (h *APIHandler) CreateCategory(c *gin.Context) {
	var req service.CreateAPICategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	category, err := h.apiService.CreateAPICategory(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "创建API分类成功",
		Data:    category,
	})
}

// UpdateCategory 更新API分类
func (h *APIHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的分类ID",
		})
		return
	}

	var req service.UpdateAPICategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	req.ID = uint(id)
	category, err := h.apiService.UpdateAPICategory(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "更新API分类成功",
		Data:    category,
	})
}

// DeleteCategory 删除API分类
func (h *APIHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的分类ID",
		})
		return
	}

	err = h.apiService.DeleteAPICategory(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "删除API分类成功",
	})
}

// GetCategory 获取API分类
func (h *APIHandler) GetCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的分类ID",
		})
		return
	}

	category, err := h.apiService.GetAPICategory(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Code:    dto.CodeNotFound,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取API分类成功",
		Data:    category,
	})
}
