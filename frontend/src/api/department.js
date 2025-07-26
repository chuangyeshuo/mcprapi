import request from '@/utils/request'

// 获取部门列表
export function getDepartmentList (params) {
  return request({
    url: '/v1/department/list',
    method: 'get',
    params
  })
}

// 获取部门详情
export function getDepartmentDetail (id) {
  return request({
    url: `/v1/department/${id}`,
    method: 'get'
  })
}

// 创建部门
export function createDepartment (data) {
  return request({
    url: '/v1/department',
    method: 'post',
    data
  })
}

// 更新部门
export function updateDepartment (id, data) {
  return request({
    url: `/v1/department/${id}`,
    method: 'put',
    data
  })
}

// 删除部门
export function deleteDepartment (id) {
  return request({
    url: `/v1/department/${id}`,
    method: 'delete'
  })
}

// 获取子部门
export function getChildDepartments (parentId) {
  return request({
    url: `/v1/department/${parentId}/children`,
    method: 'get'
  })
}

// 获取所有部门（树形结构）
export function getDepartmentTree () {
  return request({
    url: '/v1/department/tree',
    method: 'get'
  })
}