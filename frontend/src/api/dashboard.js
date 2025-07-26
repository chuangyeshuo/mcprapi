import request from '@/utils/request'

// 获取仪表盘统计数据
export function getDashboardStats() {
  return request({
    url: '/v1/dashboard/stats',
    method: 'get'
  })
}

// 获取API分类统计
export function getApiCategoryStats() {
  return request({
    url: '/v1/dashboard/api-category-stats',
    method: 'get'
  })
}

// 获取业务线API统计
export function getBusinessApiStats() {
  return request({
    url: '/v1/dashboard/business-api-stats',
    method: 'get'
  })
}

// 获取部门API统计
export function getDepartmentApiStats() {
  return request({
    url: '/v1/dashboard/department-api-stats',
    method: 'get'
  })
}