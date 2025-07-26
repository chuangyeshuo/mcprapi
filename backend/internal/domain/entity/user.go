package entity

import (
	"time"
)

// User 用户实体
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:50;uniqueIndex" json:"username"`
	Name      string    `gorm:"size:100" json:"name"`
	Email     string    `gorm:"size:100;uniqueIndex" json:"email"`
	Password  string    `gorm:"size:255" json:"-"`
	Avatar    string    `gorm:"size:500" json:"avatar"` // 头像URL
	DeptID    uint      `gorm:"index" json:"dept_id"`
	Status    int       `gorm:"default:1" json:"status"` // 1: 启用, 0: 禁用
	TokenVersion int    `gorm:"default:1" json:"token_version"` // Token版本号，用于安全控制
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 设置表名
func (User) TableName() string {
	return "users"
}