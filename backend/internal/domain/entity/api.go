/*
 * @Author: lidi10@staff.weibo.com
 * @Date: 2025-07-19 17:26:05
 * @LastEditTime: 2025-07-19 17:27:22
 * @LastEditors: lidi10@staff.weibo.com
 * @Description:
 * Copyright (c) 2023 by Weibo, All Rights Reserved.
 */
package entity

import (
	"time"
)

// API API实体
type API struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100" json:"name"`
	Path        string    `gorm:"size:200;index" json:"path"`  // API路径
	Method      string    `gorm:"size:10" json:"method"`       // HTTP方法: GET, POST, PUT, DELETE
	Description string    `gorm:"size:200" json:"description"` // 描述
	BusinessID  uint      `gorm:"index" json:"business_id"`    // 所属业务线ID
	CategoryID  uint      `gorm:"index" json:"category_id"`    // 分类ID
	Status      int       `gorm:"default:1" json:"status"`     // 1: 启用, 0: 禁用
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 设置表名
func (API) TableName() string {
	return "apis"
}

// APICategory API分类
type APICategory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:50" json:"name"`
	Code      string    `gorm:"size:50;uniqueIndex" json:"code"`
	ParentID  uint      `gorm:"index" json:"parent_id"` // 父分类ID，0表示顶级分类
	Sort      int       `gorm:"default:0" json:"sort"`  // 排序
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 设置表名
func (APICategory) TableName() string {
	return "api_categories"
}
