package service

import (
	"errors"

	"mcprapi/backend/internal/domain/entity"
	"mcprapi/backend/internal/domain/repository"
	"mcprapi/backend/internal/pkg/encrypt"
)

// UserService 用户服务
type UserService struct {
	userRepo repository.UserRepository
	deptRepo repository.DepartmentRepository
	roleRepo repository.RoleRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo repository.UserRepository, deptRepo repository.DepartmentRepository, roleRepo repository.RoleRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		deptRepo: deptRepo,
		roleRepo: roleRepo,
	}
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Avatar   string `json:"avatar"` // 头像URL，可选
	DeptID   uint   `json:"dept_id" binding:"required"`
	RoleIDs  []uint `json:"role_ids"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	ID       uint   `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Avatar   string `json:"avatar"` // 头像URL，可选
	DeptID   uint   `json:"dept_id" binding:"required"`
	Status   int    `json:"status" binding:"oneof=0 1"`
}

// CreateUser 创建用户
func (s *UserService) CreateUser(req *CreateUserRequest) (*entity.User, error) {
	// 检查部门是否存在
	_, err := s.deptRepo.GetByID(req.DeptID)
	if err != nil {
		return nil, errors.New("部门不存在")
	}

	// 检查用户名是否已存在
	existingUser, _ := s.userRepo.GetByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 创建用户
	user := &entity.User{
		Username: req.Username,
		Name:     req.Name,
		Email:    req.Email,
		Password: encrypt.GenerateHash(req.Password), // 加密密码
		Avatar:   req.Avatar, // 头像URL
		DeptID:   req.DeptID,
		Status:   1, // 默认启用
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	// 分配角色
	if len(req.RoleIDs) > 0 {
		err = s.userRepo.AssignRoles(user.ID, req.RoleIDs)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(req *UpdateUserRequest) (*entity.User, error) {
	// 检查用户是否存在
	user, err := s.userRepo.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	// 检查部门是否存在
	_, err = s.deptRepo.GetByID(req.DeptID)
	if err != nil {
		return nil, errors.New("部门不存在")
	}

	// 检查用户名是否已存在（排除自身）
	existingUser, _ := s.userRepo.GetByUsername(req.Username)
	if existingUser != nil && existingUser.ID != req.ID {
		return nil, errors.New("用户名已存在")
	}

	// 更新用户
	user.Username = req.Username
	user.Name = req.Name
	user.Email = req.Email
	user.Avatar = req.Avatar // 更新头像URL
	user.DeptID = req.DeptID
	user.Status = req.Status

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id uint) error {
	// 检查用户是否存在
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	// 删除用户
	return s.userRepo.Delete(id)
}

// GetUser 获取用户
func (s *UserService) GetUser(id uint) (*entity.User, error) {
	return s.userRepo.GetByID(id)
}

// ListUserRequest 获取用户列表请求
type ListUserRequest struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
	Query    string `json:"query" form:"query"`
	DeptID   uint   `json:"dept_id" form:"dept_id"` // 部门ID过滤
}

// UserWithRoles 包含角色信息的用户
type UserWithRoles struct {
	*entity.User
	Roles []string `json:"roles"`
}

// ListUserResponse 获取用户列表响应
type ListUserResponse struct {
	Total int64             `json:"total"`
	Items []*UserWithRoles  `json:"items"`
}

// ListUser 获取用户列表
func (s *UserService) ListUser(req *ListUserRequest) (*ListUserResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	users, total, err := s.userRepo.List(req.Page, req.PageSize, req.Query)
	if err != nil {
		return nil, err
	}

	// 为每个用户获取角色信息
	usersWithRoles := make([]*UserWithRoles, len(users))
	for i, user := range users {
		roles, err := s.userRepo.GetUserRoles(user.ID)
		if err != nil {
			// 如果获取角色失败，设置为空数组
			roles = []string{}
		}
		
		usersWithRoles[i] = &UserWithRoles{
			User:  user,
			Roles: roles,
		}
	}

	return &ListUserResponse{
		Total: total,
		Items: usersWithRoles,
	}, nil
}

// GetTotalCount 获取用户总数
func (s *UserService) GetTotalCount() (int64, error) {
	return s.userRepo.GetTotalCount()
}

// ListAllUsers 获取所有用户列表
func (s *UserService) ListAllUsers() ([]*entity.User, error) {
	return s.userRepo.ListAll()
}

// AssignRolesRequest 分配角色请求
type AssignRolesRequest struct {
	UserID  uint   `json:"user_id" binding:"required"`
	RoleIDs []uint `json:"role_ids" binding:"required"`
}

// AssignRoles 分配角色
func (s *UserService) AssignRoles(req *AssignRolesRequest) error {
	// 检查用户是否存在
	user, err := s.userRepo.GetByID(req.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	// 检查角色是否存在
	for _, roleID := range req.RoleIDs {
		role, err := s.roleRepo.GetByID(roleID)
		if err != nil {
			return err
		}
		if role == nil {
			return errors.New("角色不存在")
		}
	}

	// 分配角色
	return s.userRepo.AssignRoles(req.UserID, req.RoleIDs)
}

// GetUserRoles 获取用户角色
func (s *UserService) GetUserRoles(userID uint) ([]string, error) {
	// 检查用户是否存在
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	// 获取用户角色
	return s.userRepo.GetUserRoles(userID)
}

// ListUserByDept 按部门获取用户列表
func (s *UserService) ListUserByDept(req *ListUserRequest) (*ListUserResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 检查部门是否存在
	if req.DeptID > 0 {
		_, err := s.deptRepo.GetByID(req.DeptID)
		if err != nil {
			return nil, errors.New("部门不存在")
		}
	}

	users, total, err := s.userRepo.ListByDept(req.Page, req.PageSize, req.Query, req.DeptID)
	if err != nil {
		return nil, err
	}

	// 为每个用户获取角色信息
	usersWithRoles := make([]*UserWithRoles, len(users))
	for i, user := range users {
		roles, err := s.userRepo.GetUserRoles(user.ID)
		if err != nil {
			// 如果获取角色失败，设置为空数组
			roles = []string{}
		}
		
		usersWithRoles[i] = &UserWithRoles{
			User:  user,
			Roles: roles,
		}
	}

	return &ListUserResponse{
		Total: total,
		Items: usersWithRoles,
	}, nil
}