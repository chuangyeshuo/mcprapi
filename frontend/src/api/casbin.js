import request from '@/utils/request';

// 获取Casbin策略列表
export function getCasbinPolicyList(params) {
  return request({
    url: '/v1/casbin/policy/list',
    method: 'get',
    params
  })
}

// 添加Casbin策略
export function addCasbinPolicy(data) {
  return request({
    url: '/v1/casbin/policy',
    method: 'post',
    data
  })
}

// 更新Casbin策略
export function updateCasbinPolicy(data) {
  return request({
    url: '/v1/casbin/policy',
    method: 'put',
    data
  })
}

// 删除Casbin策略
export function deleteCasbinPolicy(id) {
  return request({
    url: `/v1/casbin/policy/${id}`,
    method: 'delete'
  })
}

// 批量删除Casbin策略
export function batchDeleteCasbinPolicy(ids) {
  return request({
    url: '/v1/casbin/policy/batch',
    method: 'delete',
    data: { ids }
  })
}

// 重新加载Casbin策略
export function reloadCasbinPolicy() {
  return request({
    url: '/v1/casbin/policy/reload',
    method: 'post'
  })
}