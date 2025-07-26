package casbinx

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// Enforcer Casbin执行器
type Enforcer struct {
	enforcer *casbin.Enforcer
}

// NewEnforcer 创建Casbin执行器
func NewEnforcer(modelPath, policyType string) (*Enforcer, error) {
	var adapter persist.Adapter
	var err error

	// 根据策略类型创建适配器
	switch policyType {
	case "file":
		// 使用文件存储策略
		adapter = fileadapter.NewAdapter("configs/policy.csv")
	case "mysql":
		// 使用MySQL存储策略
		// 注意：这里需要传入数据库连接，实际应用中应该从配置中获取
		adapter, err = gormadapter.NewAdapter("mysql", "root:password@tcp(127.0.0.1:3306)/api_auth_dev", true)
		if err != nil {
			return nil, err
		}
	default:
		// 默认使用文件存储策略
		adapter = fileadapter.NewAdapter("")
	}

	// 加载模型
	m, err := model.NewModelFromFile(modelPath)
	if err != nil {
		return nil, err
	}

	// 创建执行器
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	// 加载策略
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	return &Enforcer{enforcer: enforcer}, nil
}

// NewEnforcerWithDB 使用现有的数据库连接创建Casbin执行器
func NewEnforcerWithDB(modelPath string, db *gorm.DB) (*Enforcer, error) {
	// 创建适配器
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	// 加载模型
	m, err := model.NewModelFromFile(modelPath)
	if err != nil {
		return nil, err
	}

	// 创建执行器
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	// 加载策略
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	return &Enforcer{enforcer: enforcer}, nil
}

// Enforce 检查权限
func (e *Enforcer) Enforce(sub, obj, act string) bool {
	result, _ := e.enforcer.Enforce(sub, obj, act)
	return result
}

// EnforceWithDept 检查权限（包含部门维度）
func (e *Enforcer) EnforceWithDept(sub, obj, act string, deptID uint) bool {
	// 将部门ID转换为字符串
	deptStr := "*" // 默认为通配符，表示不限制部门
	if deptID > 0 {
		deptStr = fmt.Sprintf("%d", deptID)
	}

	result, _ := e.enforcer.Enforce(sub, obj, act, deptStr)
	return result
}

// AddPolicy 添加策略
func (e *Enforcer) AddPolicy(sub, obj, act, eft string) (bool, error) {
	return e.enforcer.AddPolicy(sub, obj, act, eft)
}

// AddPolicyWithDept 添加策略（包含部门维度）
func (e *Enforcer) AddPolicyWithDept(sub, obj, act, dept, eft string) (bool, error) {
	return e.enforcer.AddPolicy(sub, obj, act, dept, eft)
}

// RemovePolicy 删除策略
func (e *Enforcer) RemovePolicy(sub, obj, act string) (bool, error) {
	return e.enforcer.RemovePolicy(sub, obj, act)
}

// RemovePolicyWithDept 删除策略（包含部门维度）
func (e *Enforcer) RemovePolicyWithDept(sub, obj, act, dept, eft string) (bool, error) {
	return e.enforcer.RemovePolicy(sub, obj, act, dept, eft)
}

// AddRoleForUser 为用户添加角色
func (e *Enforcer) AddRoleForUser(user, role string) (bool, error) {
	return e.enforcer.AddRoleForUser(user, role)
}

// DeleteRoleForUser 删除用户角色
func (e *Enforcer) DeleteRoleForUser(user, role string) (bool, error) {
	return e.enforcer.DeleteRoleForUser(user, role)
}

// DeleteRolesForUser 删除用户所有角色
func (e *Enforcer) DeleteRolesForUser(user string) (bool, error) {
	return e.enforcer.DeleteRolesForUser(user)
}

// GetRolesForUser 获取用户角色
func (e *Enforcer) GetRolesForUser(user string) ([]string, error) {
	return e.enforcer.GetRolesForUser(user)
}

// GetUsersForRole 获取角色用户
func (e *Enforcer) GetUsersForRole(role string) ([]string, error) {
	return e.enforcer.GetUsersForRole(role)
}

// HasRoleForUser 检查用户是否有角色
func (e *Enforcer) HasRoleForUser(user, role string) (bool, error) {
	return e.enforcer.HasRoleForUser(user, role)
}

// GetPolicy 获取所有策略
func (e *Enforcer) GetPolicy() [][]string {
	return e.enforcer.GetPolicy()
}

// GetFilteredPolicy 获取过滤后的策略
func (e *Enforcer) GetFilteredPolicy(fieldIndex int, fieldValues ...string) [][]string {
	return e.enforcer.GetFilteredPolicy(fieldIndex, fieldValues...)
}

// SavePolicy 保存策略
func (e *Enforcer) SavePolicy() error {
	return e.enforcer.SavePolicy()
}

// SafeSavePolicy 安全保存策略 - 只在确保内存策略完整时才保存
func (e *Enforcer) SafeSavePolicy() error {
	// 获取当前内存中的策略数量
	policies := e.enforcer.GetPolicy()
	groupingPolicies := e.enforcer.GetGroupingPolicy()
	
	// 检查策略是否为空（防止意外清空数据库）
	if len(policies) == 0 && len(groupingPolicies) == 0 {
		return fmt.Errorf("策略为空，拒绝保存以防止清空数据库")
	}
	
	// 如果策略数量合理，则保存
	return e.enforcer.SavePolicy()
}

// GetGroupingPolicy 获取分组策略
func (e *Enforcer) GetGroupingPolicy() [][]string {
	return e.enforcer.GetGroupingPolicy()
}

// SyncPolicyToMemory 将数据库中的策略同步到内存（推荐使用）
func (e *Enforcer) SyncPolicyToMemory() error {
	// 重新从数据库加载策略到内存
	return e.enforcer.LoadPolicy()
}

// AddPolicyAndSync 添加策略并同步到内存
func (e *Enforcer) AddPolicyAndSync(sub, obj, act, eft string) (bool, error) {
	// 添加策略到内存
	added, err := e.enforcer.AddPolicy(sub, obj, act, eft)
	if err != nil {
		return false, err
	}
	
	// 如果成功添加，则同步到内存（从数据库重新加载）
	if added {
		if err := e.SyncPolicyToMemory(); err != nil {
			return false, err
		}
	}
	
	return added, nil
}

// AddPolicyWithDeptAndSync 添加策略（包含部门维度）并同步到内存
func (e *Enforcer) AddPolicyWithDeptAndSync(sub, obj, act, dept, eft string) (bool, error) {
	// 添加策略到内存
	added, err := e.enforcer.AddPolicy(sub, obj, act, dept, eft)
	if err != nil {
		return false, err
	}
	
	// 如果成功添加，则同步到内存（从数据库重新加载）
	if added {
		if err := e.SyncPolicyToMemory(); err != nil {
			return false, err
		}
	}
	
	return added, nil
}

// LoadPolicy 加载策略
func (e *Enforcer) LoadPolicy() error {
	return e.enforcer.LoadPolicy()
}

// DeletePermissionsForUser 删除用户所有权限
func (e *Enforcer) DeletePermissionsForUser(user string) (bool, error) {
	return e.enforcer.DeletePermissionsForUser(user)
}
