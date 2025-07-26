package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"mcprapi/backend/internal/domain/entity"
	"mcprapi/backend/internal/pkg/encrypt"
)

// Config MySQL配置
type Config struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Charset  string `mapstructure:"charset"`
}

// NewMySQL 创建MySQL连接
func NewMySQL(config Config) (*gorm.DB, error) {
	// 首先尝试创建数据库（如果不存在）
	if err := createDatabaseIfNotExists(config); err != nil {
		return nil, fmt.Errorf("创建数据库失败: %v", err)
	}

	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.Charset,
	)

	// 创建GORM配置
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: logger.Default.LogMode(logger.Info),
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	// 获取底层连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生命周期

	return db, nil
}

// createDatabaseIfNotExists 创建数据库（如果不存在）
func createDatabaseIfNotExists(config Config) error {
	// 构建不包含数据库名的DSN，用于连接MySQL服务器
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=%s&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Charset,
	)

	// 连接MySQL服务器
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接MySQL服务器失败: %v", err)
	}

	// 获取底层SQL连接
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取SQL连接失败: %v", err)
	}
	defer sqlDB.Close()

	// 创建数据库（如果不存在）
	createSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET %s COLLATE %s_unicode_ci",
		config.Database, config.Charset, config.Charset)

	if err := db.Exec(createSQL).Error; err != nil {
		return fmt.Errorf("创建数据库失败: %v", err)
	}

	return nil
}

// AutoMigrate 自动迁移表结构
func AutoMigrate(db *gorm.DB) error {
	// 自动迁移表结构
	return db.AutoMigrate(
		&entity.User{},
		&entity.Role{},
		&entity.UserRole{},
		&entity.Department{},
		&entity.Business{},
		&entity.API{},
		&entity.APICategory{},
	)
}

// InitData 初始化数据
func InitData(db *gorm.DB) error {
	// 检查是否已有数据
	var count int64
	db.Model(&entity.User{}).Count(&count)
	if count > 0 {
		return nil // 已有数据，不需要初始化
	}

	// 开启事务
	return db.Transaction(func(tx *gorm.DB) error {
		// 创建默认部门
		department := &entity.Department{
			Name:   "默认部门",
			Code:   "default",
			Level:  1,
			Status: 1,
		}
		if err := tx.Create(department).Error; err != nil {
			return err
		}

		// 创建默认角色
		adminRole := &entity.Role{
			Name:        "管理员",
			Code:        "admin",
			Description: "系统管理员",
			DeptID:      department.ID,
			Status:      1,
		}
		if err := tx.Create(adminRole).Error; err != nil {
			return err
		}

		userRole := &entity.Role{
			Name:        "普通用户",
			Code:        "user",
			Description: "普通用户",
			DeptID:      department.ID,
			Status:      1,
		}
		if err := tx.Create(userRole).Error; err != nil {
			return err
		}

		// 创建默认用户
		admin := &entity.User{
			Username: "admin",
			Name:     "管理员",
			Email:    "admin@example.com",
			Password: encrypt.GenerateHash("123456"), // 使用123456作为密码
			DeptID:   department.ID,
			Status:   1,
		}
		if err := tx.Create(admin).Error; err != nil {
			return err
		}

		// 创建member用户
		member := &entity.User{
			Username: "member",
			Name:     "普通用户",
			Email:    "member@example.com",
			Password: encrypt.GenerateHash("123456"), // 使用123456作为密码
			DeptID:   department.ID,
			Status:   1,
		}
		if err := tx.Create(member).Error; err != nil {
			return err
		}

		// 分配角色
		adminUserRole := &entity.UserRole{
			UserID: admin.ID,
			RoleID: adminRole.ID,
		}
		if err := tx.Create(adminUserRole).Error; err != nil {
			return err
		}

		// 为member用户分配普通用户角色
		memberUserRole := &entity.UserRole{
			UserID: member.ID,
			RoleID: userRole.ID,
		}
		if err := tx.Create(memberUserRole).Error; err != nil {
			return err
		}

		// 创建API分类
		apiCategory := &entity.APICategory{
			Name: "系统管理",
			Code: "system",
			Sort: 1,
		}
		if err := tx.Create(apiCategory).Error; err != nil {
			return err
		}

		// 创建默认业务线
		business := &entity.Business{
			Name:        "默认业务线",
			Code:        "default",
			DeptID:      department.ID,
			Description: "默认业务线",
			Owner:       "管理员",
			Email:       "admin@example.com",
			Status:      1,
		}
		if err := tx.Create(business).Error; err != nil {
			return err
		}

		// 创建默认API
		api := &entity.API{
			Name:        "用户登录",
			Path:        "/api/v1/auth/login",
			Method:      "POST",
			Description: "用户登录接口",
			BusinessID:  business.ID,
			CategoryID:  apiCategory.ID,
			Status:      1,
		}
		if err := tx.Create(api).Error; err != nil {
			return err
		}

		// 创建默认Casbin策略
		// 为admin用户添加全部权限
		adminPolicy := map[string]interface{}{
			"ptype": "p",
			"v0":    "admin",
			"v1":    "/api/v1/*",
			"v2":    "*",
			"v3":    "*",
			"v4":    "allow",
		}
		if err := tx.Table("casbin_rule").Create(adminPolicy).Error; err != nil {
			return err
		}

		return nil
	})
}
