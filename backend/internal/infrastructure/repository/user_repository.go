package repository

import (
	"errors"

	"gorm.io/gorm"

	"mcprapi/backend/internal/domain/entity"
	"mcprapi/backend/internal/domain/repository"
)

// UserRepositoryImpl 用户仓库实现
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepositoryImpl{db: db}
}

// Create 创建用户
func (r *UserRepositoryImpl) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

// Update 更新用户
func (r *UserRepositoryImpl) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户
func (r *UserRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entity.User{}, id).Error
}

// GetByID 根据ID获取用户
func (r *UserRepositoryImpl) GetByID(id uint) (*entity.User, error) {
	var user entity.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepositoryImpl) GetByUsername(username string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// List 获取用户列表
func (r *UserRepositoryImpl) List(page, pageSize int, query string) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64

	offset := (page - 1) * pageSize

	db := r.db
	if query != "" {
		db = db.Where("username LIKE ? OR name LIKE ? OR email LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	}

	if err := db.Model(&entity.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ListAll 获取所有用户列表
func (r *UserRepositoryImpl) ListAll() ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.Find(&users).Error
	return users, err
}

// GetTotalCount 获取用户总数
func (r *UserRepositoryImpl) GetTotalCount() (int64, error) {
	var count int64
	err := r.db.Model(&entity.User{}).Count(&count).Error
	return count, err
}

// ListByDept 根据部门获取用户列表
func (r *UserRepositoryImpl) ListByDept(page, pageSize int, query string, deptID uint) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64

	offset := (page - 1) * pageSize

	db := r.db
	if query != "" {
		db = db.Where("username LIKE ? OR name LIKE ? OR email LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	}

	if deptID > 0 {
		db = db.Where("dept_id = ?", deptID)
	}

	if err := db.Model(&entity.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUserRoles 获取用户角色
func (r *UserRepositoryImpl) GetUserRoles(userID uint) ([]string, error) {
	var roleCodes []string
	if err := r.db.Model(&entity.Role{}).
		Select("roles.code").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Pluck("code", &roleCodes).Error; err != nil {
		return nil, err
	}
	return roleCodes, nil
}

// AssignRoles 为用户分配角色
func (r *UserRepositoryImpl) AssignRoles(userID uint, roleIDs []uint) error {
	// 开启事务
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除用户现有角色
		if err := tx.Where("user_id = ?", userID).Delete(&entity.UserRole{}).Error; err != nil {
			return err
		}

		// 分配新角色
		for _, roleID := range roleIDs {
			userRole := &entity.UserRole{
				UserID: userID,
				RoleID: roleID,
			}
			if err := tx.Create(userRole).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetUsersByRole 获取角色用户
func (r *UserRepositoryImpl) GetUsersByRole(roleID uint) ([]*entity.User, error) {
	var users []*entity.User
	if err := r.db.Model(&entity.User{}).Joins("JOIN user_roles ON user_roles.user_id = users.id").Where("user_roles.role_id = ?", roleID).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetUsersByDepartment 获取部门用户
func (r *UserRepositoryImpl) GetUsersByDepartment(departmentID uint) ([]*entity.User, error) {
	var users []*entity.User
	if err := r.db.Where("dept_id = ?", departmentID).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
