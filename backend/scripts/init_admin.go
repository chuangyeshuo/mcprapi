package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"

	"mcprapi/backend/internal/infrastructure/container"
	"mcprapi/backend/internal/infrastructure/database"
)

// createDatabaseIfNotExists 检查并创建数据库
func createDatabaseIfNotExists(config database.Config) error {
	// 构建不包含数据库名的DSN，用于连接MySQL服务器
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=%s&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Charset,
	)

	// 连接MySQL服务器
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("连接MySQL服务器失败: %v", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		return fmt.Errorf("无法连接到MySQL服务器: %v", err)
	}

	// 检查数据库是否存在
	var count int
	query := "SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?"
	err = db.QueryRow(query, config.Database).Scan(&count)
	if err != nil {
		return fmt.Errorf("检查数据库是否存在失败: %v", err)
	}

	// 如果数据库不存在，则创建
	if count == 0 {
		createQuery := fmt.Sprintf("CREATE DATABASE `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", config.Database)
		_, err = db.Exec(createQuery)
		if err != nil {
			return fmt.Errorf("创建数据库失败: %v", err)
		}
		fmt.Printf("数据库 '%s' 创建成功！\n", config.Database)
	} else {
		fmt.Printf("数据库 '%s' 已存在\n", config.Database)
	}

	return nil
}

// connectToDatabase 连接到指定数据库
func connectToDatabase(config database.Config) (*gorm.DB, error) {
	// 构建完整的DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.Charset,
	)

	// 创建GORM配置
	gormConfig := &gorm.Config{}

	// 连接数据库
	db, err := gorm.Open(mysqlDriver.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	return db, nil
}

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

	// 获取数据库配置
	dbConfig := database.Config{
		Host:     config.GetString("mysql.host"),
		Port:     config.GetInt("mysql.port"),
		Username: config.GetString("mysql.username"),
		Password: config.GetString("mysql.password"),
		Database: config.GetString("mysql.database"),
		Charset:  config.GetString("mysql.charset"),
	}

	// 设置默认字符集
	if dbConfig.Charset == "" {
		dbConfig.Charset = "utf8mb4"
	}

	fmt.Printf("正在检查数据库连接配置...\n")
	fmt.Printf("主机: %s:%d\n", dbConfig.Host, dbConfig.Port)
	fmt.Printf("用户: %s\n", dbConfig.Username)
	fmt.Printf("数据库: %s\n", dbConfig.Database)

	// 检查并创建数据库
	if err := createDatabaseIfNotExists(dbConfig); err != nil {
		log.Fatalf("创建数据库失败: %v", err)
	}

	// 连接到数据库
	db, err := connectToDatabase(dbConfig)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	fmt.Println("数据库连接成功！")

	// 自动迁移表结构
	fmt.Println("正在执行数据库迁移...")
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 初始化数据
	fmt.Println("正在初始化基础数据...")
	if err := database.InitData(db); err != nil {
		log.Fatalf("初始化数据失败: %v", err)
	}

	// 初始化依赖注入容器（用于获取服务）
	fmt.Println("正在初始化服务容器...")
	c, err := container.New(config)
	if err != nil {
		log.Fatalf("初始化依赖注入容器失败: %v", err)
	}
	defer c.Close()

	// 查找admin用户
	fmt.Println("正在配置admin用户权限...")
	adminUser, err := c.UserRepository.GetByUsername("admin")
	if err != nil {
		log.Fatalf("查找admin用户失败: %v", err)
	}

	// 为admin用户配置系统管理员权限
	if err := c.DeptPermissionService.InitSystemAdmin(adminUser.ID); err != nil {
		log.Fatalf("配置admin用户系统管理员权限失败: %v", err)
	}

	fmt.Println("\n🎉 数据库初始化完成！")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("默认管理员账号:")
	fmt.Println("  用户名: admin")
	fmt.Println("  密码: 123456")
	fmt.Println("  部门: 默认部门")
	fmt.Println("  角色: 管理员")
	fmt.Println("  权限: 系统管理员 (拥有 /api/v1/* 的所有操作权限)")
	fmt.Println("")
	fmt.Println("默认普通用户账号:")
	fmt.Println("  用户名: member")
	fmt.Println("  密码: 123456")
	fmt.Println("  部门: 默认部门")
	fmt.Println("  角色: 普通用户")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("⚠️  请在生产环境中立即修改默认密码！")
}
