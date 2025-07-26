import request from '@/utils/request'

// 获取API列表
export function getApiList (params) {
  return request({
    url: '/v1/api/list',
    method: 'get',
    params
  })
}

// 获取API详情
export function getApiDetail (id) {
  return request({
    url: `/v1/api/${id}`,
    method: 'get'
  })
}

// 创建API
export function createApi (data) {
  return request({
    url: '/v1/api',
    method: 'post',
    data
  })
}

// 更新API
export function updateApi (id, data) {
  return request({
    url: `/v1/api/${id}`,
    method: 'put',
    data
  })
}

// 删除API
export function deleteApi (id) {
  return request({
    url: `/v1/api/${id}`,
    method: 'delete'
  })
}

// 根据业务线获取API列表
export function getApiListByBusiness (businessId) {
  return request({
    url: `/v1/api/business/${businessId}`,
    method: 'get'
  })
}