package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mcprapi/backend/internal/domain/service"
	"mcprapi/backend/internal/transport/http/dto"
	"mcprapi/backend/internal/transport/middleware"
)

// CasbinHandler Casbin权限管理处理器
type CasbinHandler struct {
	casbinService *service.CasbinService
	authService   *service.AuthService
}

// NewCasbinHandler 创建Casbin处理器
func NewCasbinHandler(casbinService *service.CasbinService, authService *service.AuthService) *CasbinHandler {
	return &CasbinHandler{
		casbinService: casbinService,
		authService:   authService,
	}
}

// Register 注册路由
func (h *CasbinHandler) Register(router *gin.RouterGroup) {
	casbinGroup := router.Group("/casbin")
	{
		casbinGroup.GET("/policy/list", h.ListPolicy)
		casbinGroup.POST("/policy", h.AddPolicy)
		casbinGroup.PUT("/policy", h.UpdatePolicy)
		casbinGroup.DELETE("/policy/:id", h.DeletePolicy)
		casbinGroup.DELETE("/policy/batch", h.BatchDeletePolicy)
		casbinGroup.POST("/policy/reload", h.ReloadPolicy)
	}
}

// ListPolicy 获取策略列表
func (h *CasbinHandler) ListPolicy(c *gin.Context) {
	var req service.ListPolicyRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误",
		})
		return
	}

	result, err := h.casbinService.ListPolicy(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取策略列表成功",
		Data:    result,
	})
}

// AddPolicy 添加策略
func (h *CasbinHandler) AddPolicy(c *gin.Context) {
	var req service.AddPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误",
		})
		return
	}

	// 如果前端没有指定部门或指定为空，则使用当前用户的部门
	if req.V3 == "" || req.V3 == "*" {
		currentDeptID := middleware.GetCurrentDeptID(c)
		if currentDeptID > 0 {
			req.V3 = strconv.FormatUint(uint64(currentDeptID), 10)
		} else {
			req.V3 = "*" // 如果获取不到部门信息，则设置为全部部门
		}
	}

	if err := h.casbinService.AddPolicy(&req); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "添加策略成功",
	})
}

// DeletePolicy 删除策略
func (h *CasbinHandler) DeletePolicy(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的策略ID",
		})
		return
	}

	req := service.DeletePolicyRequest{ID: id}
	if err := h.casbinService.DeletePolicy(&req); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "删除策略成功",
	})
}

// BatchDeletePolicy 批量删除策略
func (h *CasbinHandler) BatchDeletePolicy(c *gin.Context) {
	var req service.BatchDeletePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	if err := h.casbinService.BatchDeletePolicy(&req); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "批量删除策略成功",
	})
}

// UpdatePolicy 更新策略
func (h *CasbinHandler) UpdatePolicy(c *gin.Context) {
	var req service.UpdatePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	// 如果前端没有指定部门或指定为空，则使用当前用户的部门
	if req.V3 == "" {
		currentDeptID := middleware.GetCurrentDeptID(c)
		if currentDeptID > 0 {
			req.V3 = strconv.FormatUint(uint64(currentDeptID), 10)
		} else {
			req.V3 = "*" // 如果获取不到部门信息，则设置为全部部门
		}
	}

	if err := h.casbinService.UpdatePolicy(&req); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "更新策略成功",
	})
}

// ReloadPolicy 重新加载策略
func (h *CasbinHandler) ReloadPolicy(c *gin.Context) {
	if err := h.casbinService.ReloadPolicy(); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "重新加载策略成功",
	})
}