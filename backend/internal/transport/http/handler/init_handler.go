package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mcprapi/backend/internal/domain/service"
	"mcprapi/backend/internal/transport/http/dto"
)

// InitHandler 数据库初始化处理器
type InitHandler struct {
	initService *service.InitService
}

// NewInitHandler 创建数据库初始化处理器
func NewInitHandler(initService *service.InitService) *InitHandler {
	return &InitHandler{
		initService: initService,
	}
}

// CheckStatus 检查数据库初始化状态
// @Summary 检查数据库初始化状态
// @Description 检查数据库是否已经初始化
// @Tags 系统初始化
// @Accept json
// @Produce json
// @Success 200 {object} dto.Response{data=map[string]bool} "成功"
// @Failure 500 {object} dto.Response "服务器错误"
// @Router /api/v1/init/status [get]
func (h *InitHandler) CheckStatus(c *gin.Context) {
	initialized, err := h.initService.CheckInitialized()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: "检查初始化状态失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取初始化状态成功",
		Data: map[string]bool{
			"initialized": initialized,
		},
	})
}

// Initialize 初始化数据库
// @Summary 初始化数据库
// @Description 执行数据库初始化，创建默认用户、角色、部门等基础数据
// @Tags 系统初始化
// @Accept json
// @Produce json
// @Success 200 {object} dto.Response "初始化成功"
// @Failure 400 {object} dto.Response "数据库已初始化"
// @Failure 500 {object} dto.Response "服务器错误"
// @Router /api/v1/init/database [post]
func (h *InitHandler) Initialize(c *gin.Context) {
	// 检查是否已经初始化
	initialized, err := h.initService.CheckInitialized()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: "检查初始化状态失败: " + err.Error(),
		})
		return
	}

	if initialized {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "数据库已经初始化，无需重复初始化",
		})
		return
	}

	// 执行初始化
	if err := h.initService.InitializeDatabase(); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: "数据库初始化失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "数据库初始化成功",
		Data: map[string]interface{}{
			"admin_username": "admin",
			"admin_password": "123456",
			"member_username": "member", 
			"member_password": "123456",
			"note": "请及时修改默认密码",
		},
	})
}

// Register 注册路由
func (h *InitHandler) Register(r *gin.RouterGroup) {
	initGroup := r.Group("/init")
	{
		initGroup.GET("/status", h.CheckStatus)
		initGroup.POST("/database", h.Initialize)
	}
}