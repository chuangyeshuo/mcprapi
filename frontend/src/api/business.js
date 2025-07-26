import request from '@/utils/request'

// 获取业务线列表
export function getBusinessList (params) {
  return request({
    url: '/v1/business/list',
    method: 'get',
    params
  })
}

// 获取业务线详情
export function getBusinessDetail (id) {
  return request({
    url: `/v1/business/${id}`,
    method: 'get'
  })
}

// 创建业务线
export function createBusiness (data) {
  return request({
    url: '/v1/business',
    method: 'post',
    data
  })
}

// 更新业务线
export function updateBusiness (id, data) {
  return request({
    url: `/v1/business/${id}`,
    method: 'put',
    data
  })
}

// 删除业务线
export function deleteBusiness (id) {
  return request({
    url: `/v1/business/${id}`,
    method: 'delete'
  })
}

// 根据部门ID获取业务线列表
export function getBusinessListByDepartment (departmentId) {
  return request({
    url: `/v1/business/department/${departmentId}`,
    method: 'get'
  })
}

// 获取所有业务线（不分页）
export function getAllBusiness () {
  return request({
    url: '/v1/business/all',
    method: 'get'
  })
}