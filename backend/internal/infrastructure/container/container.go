package container

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"

	repository "mcprapi/backend/internal/domain/repository"
	"mcprapi/backend/internal/domain/service"
	"mcprapi/backend/internal/infrastructure/cache"
	"mcprapi/backend/internal/infrastructure/database"
	repo "mcprapi/backend/internal/infrastructure/repository"
	"mcprapi/backend/internal/pkg/logger"
	"mcprapi/backend/pkg/casbinx"
)

// Container 依赖注入容器
type Container struct {
	// 配置
	Config *viper.Viper

	// 数据库
	DB *gorm.DB

	// 缓存
	Redis *cache.Redis

	// 日志
	Logger *logger.Logger

	// Casbin
	Enforcer *casbinx.Enforcer

	// 仓库
	UserRepository       repository.UserRepository
	RoleRepository       repository.RoleRepository
	APIRepository        repository.APIRepository
	DepartmentRepository repository.DepartmentRepository
	BusinessRepository   repository.BusinessRepository

	// 服务
	AuthService           *service.AuthService
	UserService           *service.UserService
	RoleService           *service.RoleService
	APIService            *service.APIService
	DepartmentService     *service.DepartmentService
	BusinessService       *service.BusinessService
	CasbinService         *service.CasbinService
	DeptPermissionService *service.DeptPermissionService
	DashboardService      *service.DashboardService
}

// New 创建依赖注入容器
func New(config *viper.Viper) (*Container, error) {
	// 创建容器
	container := &Container{
		Config: config,
	}

	// 初始化日志
	if err := container.initLogger(); err != nil {
		return nil, err
	}

	// 初始化数据库
	if err := container.initDB(); err != nil {
		return nil, err
	}

	// 初始化缓存
	if err := container.initRedis(); err != nil {
		return nil, err
	}

	// 初始化Casbin
	if err := container.initCasbin(); err != nil {
		return nil, err
	}

	// 初始化仓库
	container.initRepository()

	// 初始化服务
	container.initService()

	return container, nil
}

// initLogger 初始化日志
func (c *Container) initLogger() error {
	// 获取日志配置
	logConfig := logger.Config{
		Level:     c.Config.GetString("log.level"),
		FilePath:  c.Config.GetString("log.file_path"),
		FileMode:  c.Config.GetBool("log.file_mode"),
		FormatStr: c.Config.GetString("log.format_str"),
	}

	// 创建日志记录器
	logger, err := logger.New(logConfig)
	if err != nil {
		return err
	}

	c.Logger = logger
	return nil
}

// initDB 初始化数据库
func (c *Container) initDB() error {
	// 获取数据库配置
	dbConfig := database.Config{
		Host:     c.Config.GetString("mysql.host"),
		Port:     c.Config.GetInt("mysql.port"),
		Username: c.Config.GetString("mysql.username"),
		Password: c.Config.GetString("mysql.password"),
		Database: c.Config.GetString("mysql.database"),
		Charset:  c.Config.GetString("mysql.charset"),
	}

	// 创建数据库连接
	db, err := database.NewMySQL(dbConfig)
	if err != nil {
		return err
	}

	// 自动迁移表结构
	if err := database.AutoMigrate(db); err != nil {
		return err
	}

	// 初始化数据
	if err := database.InitData(db); err != nil {
		return err
	}

	c.DB = db
	return nil
}

// initRedis 初始化缓存
func (c *Container) initRedis() error {
	// 获取缓存配置
	redisConfig := cache.Config{
		Host:     c.Config.GetString("redis.host"),
		Port:     c.Config.GetInt("redis.port"),
		Password: c.Config.GetString("redis.password"),
		DB:       c.Config.GetInt("redis.db"),
	}

	// 创建缓存连接
	redis, err := cache.NewRedis(redisConfig)
	if err != nil {
		return err
	}

	c.Redis = redis
	return nil
}

// initCasbin 初始化Casbin
func (c *Container) initCasbin() error {
	// 获取Casbin配置
	modelPath := c.Config.GetString("casbin.model")

	// 创建Casbin执行器
	enforcer, err := casbinx.NewEnforcerWithDB(modelPath, c.DB)
	if err != nil {
		return err
	}

	c.Enforcer = enforcer
	return nil
}

// initRepository 初始化仓库
func (c *Container) initRepository() {
	c.UserRepository = repo.NewUserRepository(c.DB)
	c.RoleRepository = repo.NewRoleRepository(c.DB)
	c.APIRepository = repo.NewAPIRepository(c.DB)
	c.DepartmentRepository = repo.NewDepartmentRepository(c.DB)
	c.BusinessRepository = repo.NewBusinessRepository(c.DB)
}

// initService 初始化服务
func (c *Container) initService() {
	jwtSecret := c.Config.GetString("jwt.secret")
	c.AuthService = service.NewAuthService(c.UserRepository, c.RoleRepository, c.APIRepository, c.Enforcer, jwtSecret)
	c.UserService = service.NewUserService(c.UserRepository, c.DepartmentRepository, c.RoleRepository)
	c.RoleService = service.NewRoleService(c.RoleRepository, c.DepartmentRepository, c.APIRepository, c.UserRepository, c.Enforcer)
	c.APIService = service.NewAPIService(c.APIRepository, c.BusinessRepository, c.UserRepository, c.DepartmentRepository)
	c.DepartmentService = service.NewDepartmentService(c.DepartmentRepository)
	c.BusinessService = service.NewBusinessService(c.BusinessRepository, c.DepartmentRepository)
	c.CasbinService = service.NewCasbinService(c.Enforcer, c.DB)
	c.DeptPermissionService = service.NewDeptPermissionService(c.UserRepository, c.RoleRepository, c.DepartmentRepository, c.Enforcer)
	c.DashboardService = service.NewDashboardService(c.APIService, c.BusinessService, c.DepartmentService, c.UserService, c.APIRepository)
}

// Close 关闭容器
func (c *Container) Close() error {
	// 关闭日志
	if c.Logger != nil {
		c.Logger.Close()
	}

	// 关闭缓存
	if c.Redis != nil {
		c.Redis.Close()
	}

	return nil
}
