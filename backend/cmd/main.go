package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"mcprapi/backend/internal/infrastructure/container"
	"mcprapi/backend/internal/transport/http/handler"
	"mcprapi/backend/internal/transport/middleware"
)

func main() {
	// 解析命令行参数
	var configFile string
	flag.StringVar(&configFile, "config", "configs/dev.yaml", "配置文件路径")
	flag.Parse()

	// 加载配置文件
	config := viper.New()
	config.SetConfigFile(configFile)
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 初始化依赖注入容器
	c, err := container.New(config)
	if err != nil {
		log.Fatalf("初始化依赖注入容器失败: %v", err)
	}
	defer c.Close()

	// 初始化处理器
	userHandler := handler.NewUserHandler(c.UserService, c.AuthService)
	roleHandler := handler.NewRoleHandler(c.RoleService)
	apiHandler := handler.NewAPIHandler(c.APIService, c.AuthService)
	departmentHandler := handler.NewDepartmentHandler(c.DepartmentService, c.AuthService)
	businessHandler := handler.NewBusinessHandler(c.BusinessService, c.AuthService)
	casbinHandler := handler.NewCasbinHandler(c.CasbinService, c.AuthService)
	dashboardHandler := handler.NewDashboardHandler(c.DashboardService, c.AuthService)

	// 初始化部门级别处理器
	deptPermissionHandler := handler.NewDeptPermissionHandler(c.DeptPermissionService)
	deptUserHandler := handler.NewDeptUserHandler(c.UserService, c.DeptPermissionService)
	deptRoleHandler := handler.NewDeptRoleHandler(c.RoleService, c.DeptPermissionService)
	deptBusinessHandler := handler.NewDeptBusinessHandler(c.BusinessService, c.DeptPermissionService)
	deptAPIHandler := handler.NewDeptAPIHandler(c.APIService, c.BusinessService, c.DeptPermissionService)

	// 设置Gin模式
	gin.SetMode(config.GetString("gin.mode"))

	// 创建Gin引擎
	r := gin.Default()

	// 注册中间件
	r.Use(middleware.CORS())

	// 注册路由
	// 公共API
	publicAPI := r.Group("/api/v1")
	{
		// 用户登录
		publicAPI.POST("/auth/login", userHandler.Login)
		// 用户登出
		publicAPI.POST("/auth/logout", userHandler.Logout)
		// 扫码登录
		publicAPI.POST("/auth/qrcode", userHandler.GenerateQRCode)
		publicAPI.GET("/auth/qrcode/:id", userHandler.CheckQRCode)
	}

	// 基础认证API（只需要JWT认证，不需要权限控制）
	basicAuthAPI := r.Group("/api/v1")
	basicAuthAPI.Use(middleware.JWT(config.GetString("jwt.secret")))
	{
		// 获取当前用户信息 - 所有登录用户都可以访问
		basicAuthAPI.GET("/user/info", userHandler.GetInfo)
		// Token相关路由 - 所有登录用户都可以访问
		basicAuthAPI.GET("/user/:id/token", userHandler.GetUserToken)
		basicAuthAPI.POST("/user/refresh-token", userHandler.RefreshUserToken)
		basicAuthAPI.POST("/user/refresh-token-with-version", userHandler.RefreshUserTokenWithVersion)
		// 获取所有业务线 - 所有登录用户都可以访问
		basicAuthAPI.GET("/business/all", businessHandler.GetAll)
		// API列表 - 所有登录用户都可以访问（非admin只能看到自己部门的）
		basicAuthAPI.GET("/api/list", apiHandler.List)
		// 部门列表 - 所有登录用户都可以访问
		basicAuthAPI.GET("/department/list", departmentHandler.List)
		// 业务线列表 - 所有登录用户都可以访问
		basicAuthAPI.GET("/business/list", businessHandler.List)
		// 根据部门获取业务线 - 所有登录用户都可以访问
		basicAuthAPI.GET("/business/department/:id", businessHandler.GetByDepartment)
		// API相关路由 - 所有登录用户都可以访问
		basicAuthAPI.GET("/api/:id", apiHandler.Get)
		basicAuthAPI.PUT("/api/:id", apiHandler.Update)
		basicAuthAPI.GET("/api/business/:id", apiHandler.ListByBusiness)
		basicAuthAPI.GET("/api/category/list", apiHandler.ListCategories)
		basicAuthAPI.GET("/api/category/:id", apiHandler.GetCategory)
		// 创建API - 所有登录用户都可以访问
		basicAuthAPI.POST("/api", apiHandler.Create)
		// 权限校验 - 所有登录用户都可以访问
		basicAuthAPI.POST("/api/check-permission", apiHandler.CheckPermission)

		// Dashboard路由 - 所有登录用户都可以访问
		dashboardHandler.RegisterRoutes(basicAuthAPI)
	}

	// 认证API（需要JWT认证和权限控制）
	authAPI := r.Group("/api/v1")
	authAPI.Use(middleware.JWT(config.GetString("jwt.secret")))
	authAPI.Use(middleware.Casbin(c.Enforcer, c.UserRepository)) // 添加权限控制中间件
	{
		// 用户管理
		userHandler.Register(authAPI)

		// 角色管理
		roleHandler.Register(authAPI)

		// API管理
		apiHandler.Register(authAPI)

		// 部门管理
		departmentHandler.Register(authAPI)

		// 业务线管理
		businessHandler.Register(authAPI)

		// 部门级别管理路由（需要部门权限检查）
		deptPermissionHandler.Register(authAPI)
		deptUserHandler.Register(authAPI)
		deptRoleHandler.Register(authAPI)
		deptBusinessHandler.Register(authAPI)
		deptAPIHandler.Register(authAPI)
	}

	// Casbin权限管理API（只需要JWT认证，不需要Casbin权限控制）
	casbinAPI := r.Group("/api/v1")
	casbinAPI.Use(middleware.JWT(config.GetString("jwt.secret")))
	{
		// Casbin权限管理
		casbinHandler.Register(casbinAPI)
	}

	// 启动HTTP服务器
	port := config.GetInt("server.port")
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	// 优雅关闭
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("启动HTTP服务器失败: %v", err)
		}
	}()

	log.Printf("服务器已启动，监听端口: %d", port)

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("关闭服务器...")

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭HTTP服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("关闭服务器失败: %v", err)
	}

	log.Println("服务器已关闭")
}
