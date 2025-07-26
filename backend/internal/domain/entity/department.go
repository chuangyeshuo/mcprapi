package entity

import (
	"time"
)

// Department 部门实体
type Department struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100" json:"name"`
	Code      string    `gorm:"size:50;uniqueIndex" json:"code"`
	ParentID  uint      `gorm:"index" json:"parent_id"` // 父部门ID，0表示顶级部门
	Level     int       `gorm:"default:1" json:"level"`  // 层级，1表示集团，2表示部门，3表示子部门
	Sort      int       `gorm:"default:0" json:"sort"`   // 排序
	Status    int       `gorm:"default:1" json:"status"` // 1: 启用, 0: 禁用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 设置表名
func (Department) TableName() string {
	return "departments"
}

// Business 业务线实体
type Business struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:100" json:"name"`
	Code         string    `gorm:"size:50;uniqueIndex" json:"code"`
	DeptID       uint      `gorm:"index" json:"dept_id"`       // 所属部门ID
	Description  string    `gorm:"size:200" json:"description"` // 描述
	Owner        string    `gorm:"size:50" json:"owner"`       // 负责人
	Email        string    `gorm:"size:100" json:"email"`      // 联系邮箱
	Status       int       `gorm:"default:1" json:"status"`    // 1: 启用, 0: 禁用
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 设置表名
func (Business) TableName() string {
	return "businesses"
}