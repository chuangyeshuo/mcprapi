/*
- @Author: lidi10@staff.weibo.com

- @Date: 2025-07-26 22:16:45

- Copyright (c) 2023 by Weibo, All Rights Reserved.
*/
package service

import (
	"gorm.io/gorm"

	"mcprapi/backend/internal/infrastructure/database"
	"mcprapi/backend/pkg/casbinx"
)

// InitService 数据库初始化服务
type InitService struct {
	db       *gorm.DB
	enforcer *casbinx.Enforcer
}

// NewInitService 创建数据库初始化服务
func NewInitService(db *gorm.DB, enforcer *casbinx.Enforcer) *InitService {
	return &InitService{
		db:       db,
		enforcer: enforcer,
	}
}

// CheckInitialized 检查数据库是否已初始化
func (s *InitService) CheckInitialized() (bool, error) {
	var count int64
	err := s.db.Table("users").Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// InitializeDatabase 初始化数据库
func (s *InitService) InitializeDatabase() error {
	// 初始化数据库数据
	if err := database.InitData(s.db); err != nil {
		return err
	}

	// 重新加载Casbin策略到内存
	if s.enforcer != nil {
		if err := s.enforcer.LoadPolicy(); err != nil {
			return err
		}
	}

	return nil
}
