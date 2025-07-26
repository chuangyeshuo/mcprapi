package entity

import (
	"time"
)

// Role 角色实体
type Role struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:50;uniqueIndex" json:"name"`
	Code        string    `gorm:"size:50;uniqueIndex" json:"code"`
	Description string    `gorm:"size:200" json:"description"`
	DeptID      uint      `gorm:"index" json:"dept_id"` // 所属部门ID
	Status      int       `gorm:"default:1" json:"status"` // 1: 启用, 0: 禁用
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 设置表名
func (Role) TableName() string {
	return "roles"
}

// UserRole 用户角色关联
type UserRole struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index:idx_user_role,unique" json:"user_id"`
	RoleID    uint      `gorm:"index:idx_user_role,unique" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 设置表名
func (UserRole) TableName() string {
	return "user_roles"
}