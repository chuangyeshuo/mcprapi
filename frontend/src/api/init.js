import request from '@/utils/request'

// Check database initialization status
export function checkInitStatus () {
  return request({
    url: '/v1/init/status',
    method: 'get'
  })
}

// Initialize database
export function initDatabase () {
  return request({
    url: '/v1/init/database',
    method: 'post'
  })
}