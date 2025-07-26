package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mcprapi/backend/internal/domain/service"
	"mcprapi/backend/internal/transport/http/dto"
	"mcprapi/backend/internal/transport/middleware"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *service.UserService
	authService *service.AuthService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService *service.UserService, authService *service.AuthService) *UserHandler {
	return &UserHandler{
		userService: userService,
		authService: authService,
	}
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Code:    dto.CodeUnauthorized,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "登录成功",
		Data:    resp,
	})
}

// Logout 用户登出
func (h *UserHandler) Logout(c *gin.Context) {
	resp, err := h.authService.Logout()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: resp.Message,
		Data:    resp,
	})
}

// GenerateQRCode 生成扫码登录二维码
func (h *UserHandler) GenerateQRCode(c *gin.Context) {
	qrCode, err := h.authService.GenerateQRCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: "生成二维码失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "生成二维码成功",
		Data:    gin.H{"qr_code": qrCode},
	})
}

// CheckQRCode 检查二维码状态
func (h *UserHandler) CheckQRCode(c *gin.Context) {
	qrID := c.Param("id")
	if qrID == "" {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "二维码ID不能为空",
		})
		return
	}

	// 这里应该检查二维码状态，简化处理
	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "二维码状态检查成功",
		Data:    gin.H{"status": "pending"},
	})
}

// GetInfo 获取当前用户信息
func (h *UserHandler) GetInfo(c *gin.Context) {
	// 从JWT中间件获取当前用户ID
	userID := middleware.GetCurrentUser(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Code:    dto.CodeUnauthorized,
			Message: "未认证的用户",
		})
		return
	}

	user, err := h.userService.GetUser(userID)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Code:    dto.CodeUnauthorized,
			Message: "用户不存在或已被删除，请重新登录",
		})
		return
	}

	// 获取用户角色
	roles, err := h.userService.GetUserRoles(userID)
	if err != nil {
		// 如果获取角色失败，设置默认角色
		roles = []string{"user"}
	}

	// 构建返回数据
	userInfo := map[string]interface{}{
		"id":          user.ID,
		"username":    user.Username,
		"name":        user.Name,
		"email":       user.Email,
		"dept_id":     user.DeptID,
		"status":      user.Status,
		"created_at":  user.CreatedAt,
		"updated_at":  user.UpdatedAt,
		"roles":       roles,
		"permissions": []string{},  // 暂时返回空权限数组
		"avatar":      user.Avatar, // 返回用户头像URL
		"userId":      user.ID,     // 前端期望的字段名
		"deptId":      user.DeptID, // 前端期望的字段名
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取用户信息成功",
		Data:    userInfo,
	})
}

// Register 注册路由
func (h *UserHandler) Register(router *gin.RouterGroup) {
	userRouter := router.Group("/user")
	{
		// 获取用户列表
		userRouter.GET("/list", h.List)
		// 创建用户
		userRouter.POST("", h.Create)
		// 更新用户
		userRouter.PUT("/:id", h.Update)
		// 删除用户
		userRouter.DELETE("/:id", h.Delete)
		// 获取用户
		userRouter.GET("/:id", h.Get)
		// 分配角色
		userRouter.POST("/:id/roles", h.AssignRoles)
		// 分配角色（新路由，匹配前端调用）
		userRouter.POST("/assign-roles", h.AssignRoles)
		// 获取用户角色
		userRouter.GET("/:id/roles", h.GetUserRoles)
		// 注意：Token相关路由已在main.go的basicAuthAPI中注册
	}
}

// Create 创建用户
func (h *UserHandler) Create(c *gin.Context) {
	var req service.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

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
		Message: "创建用户成功",
		Data:    user,
	})
}

// Update 更新用户
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的用户ID",
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

	req.ID = uint(id)
	user, err := h.userService.UpdateUser(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "更新用户成功",
		Data:    user,
	})
}

// Delete 删除用户
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的用户ID",
		})
		return
	}

	err = h.userService.DeleteUser(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "删除用户成功",
	})
}

// Get 获取用户
func (h *UserHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的用户ID",
		})
		return
	}

	user, err := h.userService.GetUser(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Code:    dto.CodeNotFound,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取用户成功",
		Data:    user,
	})
}

// List 获取用户列表
func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	username := c.DefaultQuery("username", "")
	deptID, _ := strconv.ParseUint(c.DefaultQuery("dept_id", "0"), 10, 32)

	req := &service.ListUserRequest{
		Page:     page,
		PageSize: limit,
		Query:    username,
		DeptID:   uint(deptID),
	}

	var resp *service.ListUserResponse
	var err error

	// 如果指定了部门ID，使用按部门过滤的方法
	if req.DeptID > 0 {
		resp, err = h.userService.ListUserByDept(req)
	} else {
		resp, err = h.userService.ListUser(req)
	}

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
		Message: "获取用户列表成功",
		Data:    result,
	})
}

// AssignRoles 分配角色
func (h *UserHandler) AssignRoles(c *gin.Context) {
	var req service.AssignRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	err := h.userService.AssignRoles(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeBusinessError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "分配角色成功",
	})
}

// GetUserRoles 获取用户角色
func (h *UserHandler) GetUserRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的用户ID",
		})
		return
	}

	roles, err := h.userService.GetUserRoles(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取用户角色成功",
		Data:    roles,
	})
}

// GetUserToken 获取用户Token
func (h *UserHandler) GetUserToken(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "无效的用户ID",
		})
		return
	}

	// 获取用户信息
	user, err := h.userService.GetUser(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Code:    dto.CodeNotFound,
			Message: "用户不存在",
		})
		return
	}

	// 生成Token
	token, _, err := h.authService.GenerateUserToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: "生成Token失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "获取用户Token成功",
		Data:    gin.H{"token": token},
	})
}

// RefreshUserToken 刷新用户Token（不递增版本号）
func (h *UserHandler) RefreshUserToken(c *gin.Context) {
	var req struct {
		UserID     uint `json:"user_id" binding:"required"`
		ExpireDays int  `json:"expire_days" binding:"required,min=1,max=365"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	// 获取用户信息
	user, err := h.userService.GetUser(req.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Code:    dto.CodeNotFound,
			Message: "用户不存在",
		})
		return
	}

	// 生成新的Token
	token, expiresAt, err := h.authService.GenerateUserTokenWithExpiry(user, req.ExpireDays)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: "刷新Token失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "刷新Token成功",
		Data: map[string]interface{}{
			"token":      token,
			"expires_at": expiresAt,
			"user":       user,
		},
	})
}

// RefreshUserTokenWithVersion 刷新用户Token并递增版本号
func (h *UserHandler) RefreshUserTokenWithVersion(c *gin.Context) {
	var req struct {
		UserID     uint `json:"user_id" binding:"required"`
		ExpireDays int  `json:"expire_days" binding:"required,min=1,max=365"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    dto.CodeInvalidParams,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	// 获取用户信息
	user, err := h.userService.GetUser(req.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Code:    dto.CodeNotFound,
			Message: "用户不存在",
		})
		return
	}

	// 生成新的Token并递增版本号
	token, expiresAt, err := h.authService.GenerateUserTokenWithVersionIncrement(user, req.ExpireDays)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Code:    dto.CodeInternalError,
			Message: "生成Token失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    dto.CodeSuccess,
		Message: "刷新Token+版本号成功，所有旧Token已失效",
		Data: map[string]interface{}{
			"token":      token,
			"expires_at": expiresAt,
			"user":       user,
		},
	})
}
