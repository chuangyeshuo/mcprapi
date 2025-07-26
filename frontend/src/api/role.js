import request from '@/utils/request'

// 获取角色列表
export function getRoleList (params) {
  return request({
    url: '/v1/role/list',
    method: 'get',
    params
  })
}

// 获取角色详情
export function getRoleDetail (id) {
  return request({
    url: `/v1/role/${id}`,
    method: 'get'
  })
}

// 创建角色
export function createRole (data) {
  return request({
    url: '/v1/role',
    method: 'post',
    data
  })
}

// 更新角色
export function updateRole (id, data) {
  return request({
    url: `/v1/role/${id}`,
    method: 'put',
    data
  })
}

// 删除角色
export function deleteRole (id) {
  return request({
    url: `/v1/role/${id}`,
    method: 'delete'
  })
}

// 获取角色权限
export function getRolePermissions (roleId) {
  return request({
    url: `/v1/role/${roleId}/permissions`,
    method: 'get'
  })
}

// 更新角色权限
export function updateRolePermissions (roleId, permissions) {
  return request({
    url: `/v1/role/${roleId}/permissions`,
    method: 'put',
    data: { permissions }
  })
}

// 获取角色API权限
export function getRoleAPIPermissions (roleId) {
  return request({
    url: `/v1/role/${roleId}/api-permissions`,
    method: 'get'
  })
}

// 更新角色API权限
export function updateRoleAPIPermissions (roleId, apiIds, deptId = null) {
  const data = { api_ids: apiIds }
  if (deptId) {
    data.dept_id = deptId
  }
  return request({
    url: `/v1/role/${roleId}/api-permissions`,
    method: 'put',
    data
  })
}

// 获取所有角色（不分页）
export function getAllRoles () {
  return request({
    url: '/v1/role/all',
    method: 'get'
  })
}

// 获取用户可访问的角色
export function getUserAccessibleRoles () {
  return request({
    url: '/v1/role/user-accessible',
    method: 'get'
  })
}