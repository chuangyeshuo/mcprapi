package service

import (
	"errors"

	"gorm.io/gorm"

	"mcprapi/backend/pkg/casbinx"
)

// CasbinService Casbin权限管理服务
type CasbinService struct {
	enforcer *casbinx.Enforcer
	db       *gorm.DB
}

// NewCasbinService 创建Casbin服务
func NewCasbinService(enforcer *casbinx.Enforcer, db *gorm.DB) *CasbinService {
	return &CasbinService{
		enforcer: enforcer,
		db:       db,
	}
}

// PolicyRule 策略规则
type PolicyRule struct {
	ID    int    `json:"id" gorm:"column:id"`
	PType string `json:"ptype" gorm:"column:ptype"`
	V0    string `json:"v0" gorm:"column:v0"`
	V1    string `json:"v1" gorm:"column:v1"`
	V2    string `json:"v2" gorm:"column:v2"`
	V3    string `json:"v3" gorm:"column:v3"`
	V4    string `json:"v4" gorm:"column:v4"`
	V5    string `json:"v5" gorm:"column:v5"`
}

// ListPolicyRequest 获取策略列表请求
type ListPolicyRequest struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
	PType    string `json:"ptype" form:"ptype"`
	Subject  string `json:"subject" form:"subject"`
}

// ListPolicyResponse 获取策略列表响应
type ListPolicyResponse struct {
	List  []*PolicyRule `json:"list"`
	Total int64         `json:"total"`
}

// AddPolicyRequest 添加策略请求
type AddPolicyRequest struct {
	PType string `json:"ptype" binding:"required"`
	V0    string `json:"v0" binding:"required"`
	V1    string `json:"v1" binding:"required"`
	V2    string `json:"v2" binding:"required"`
	V3    string `json:"v3"`
	V4    string `json:"v4"`
	V5    string `json:"v5"`
}

// UpdatePolicyRequest 更新策略请求
type UpdatePolicyRequest struct {
	ID    int    `json:"id" binding:"required"`
	PType string `json:"ptype" binding:"required"`
	V0    string `json:"v0" binding:"required"`
	V1    string `json:"v1" binding:"required"`
	V2    string `json:"v2" binding:"required"`
	V3    string `json:"v3"`
	V4    string `json:"v4"`
	V5    string `json:"v5"`
}

// DeletePolicyRequest 删除策略请求
type DeletePolicyRequest struct {
	ID int `json:"id" binding:"required"`
}

// BatchDeletePolicyRequest 批量删除策略请求
type BatchDeletePolicyRequest struct {
	IDs []int `json:"ids" binding:"required,min=1"`
}

// ListPolicy 获取策略列表
func (s *CasbinService) ListPolicy(req *ListPolicyRequest) (*ListPolicyResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 构建查询条件
	query := s.db.Table("casbin_rule")

	// 按策略类型过滤
	if req.PType != "" {
		query = query.Where("ptype = ?", req.PType)
	}

	// 按主体过滤
	if req.Subject != "" {
		query = query.Where("v0 LIKE ?", "%"+req.Subject+"%")
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	var rules []*PolicyRule
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Find(&rules).Error; err != nil {
		return nil, err
	}

	return &ListPolicyResponse{
		List:  rules,
		Total: total,
	}, nil
}

// AddPolicy 添加策略
func (s *CasbinService) AddPolicy(req *AddPolicyRequest) error {
	// 设置默认值
	dept := req.V3
	if dept == "" {
		dept = "*"
	}
	effect := req.V4
	if effect == "" {
		effect = "allow"
	}

	// 直接在数据库中插入记录
	insertData := map[string]interface{}{
		"ptype": req.PType,
		"v0":    req.V0,
		"v1":    req.V1,
		"v2":    req.V2,
		"v3":    dept,
		"v4":    effect,
		"v5":    req.V5,
	}

	if err := s.db.Table("casbin_rule").Create(insertData).Error; err != nil {
		return errors.New("策略添加失败")
	}

	// 重新加载策略到内存
	return s.enforcer.LoadPolicy()
}

// UpdatePolicy 更新策略
func (s *CasbinService) UpdatePolicy(req *UpdatePolicyRequest) error {
	// 从数据库中查找策略
	var rule PolicyRule
	if err := s.db.Table("casbin_rule").Where("id = ?", req.ID).First(&rule).Error; err != nil {
		return errors.New("策略不存在")
	}

	// 设置默认值
	dept := req.V3
	if dept == "" {
		dept = "*"
	}
	effect := req.V4
	if effect == "" {
		effect = "allow"
	}

	// 直接更新数据库记录
	updateData := map[string]interface{}{
		"ptype": req.PType,
		"v0":    req.V0,
		"v1":    req.V1,
		"v2":    req.V2,
		"v3":    dept,
		"v4":    effect,
		"v5":    req.V5,
	}

	if err := s.db.Table("casbin_rule").Where("id = ?", req.ID).Updates(updateData).Error; err != nil {
		return errors.New("策略更新失败")
	}

	// 重新加载策略到内存
	return s.enforcer.LoadPolicy()
}

// DeletePolicy 删除策略
func (s *CasbinService) DeletePolicy(req *DeletePolicyRequest) error {
	// 从数据库中查找策略
	var rule PolicyRule
	if err := s.db.Table("casbin_rule").Where("id = ?", req.ID).First(&rule).Error; err != nil {
		return errors.New("策略不存在")
	}

	// 直接从数据库删除记录
	if err := s.db.Table("casbin_rule").Where("id = ?", req.ID).Delete(&rule).Error; err != nil {
		return errors.New("策略删除失败")
	}

	// 重新加载策略到内存
	return s.enforcer.LoadPolicy()
}

// BatchDeletePolicy 批量删除策略
func (s *CasbinService) BatchDeletePolicy(req *BatchDeletePolicyRequest) error {
	// 检查策略是否存在
	var count int64
	if err := s.db.Table("casbin_rule").Where("id IN ?", req.IDs).Count(&count).Error; err != nil {
		return errors.New("查询策略失败")
	}

	if count != int64(len(req.IDs)) {
		return errors.New("部分策略不存在")
	}

	// 批量删除策略
	if err := s.db.Table("casbin_rule").Where("id IN ?", req.IDs).Delete(&PolicyRule{}).Error; err != nil {
		return errors.New("批量删除策略失败")
	}

	// 重新加载策略到内存
	return s.enforcer.LoadPolicy()
}

// ReloadPolicy 重新加载策略
func (s *CasbinService) ReloadPolicy() error {
	return s.enforcer.LoadPolicy()
}
