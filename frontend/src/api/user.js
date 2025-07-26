import request from '@/utils/request'

// 用户登录
export function login (data) {
  return request({
    url: '/v1/auth/login',
    method: 'post',
    data
  })
}

// 获取用户信息
export function getInfo () {
  return request({
    url: '/v1/user/info',
    method: 'get'
  })
}

// 用户登出
export function logout () {
  return request({
    url: '/v1/auth/logout',
    method: 'post'
  })
}

// 获取用户列表
export function getUserList (params) {
  return request({
    url: '/v1/user/list',
    method: 'get',
    params
  })
}

// 获取用户详情
export function getUserDetail (id) {
  return request({
    url: `/v1/user/${id}`,
    method: 'get'
  })
}

// 创建用户
export function createUser (data) {
  return request({
    url: '/v1/user',
    method: 'post',
    data
  })
}

// 更新用户
export function updateUser (id, data) {
  return request({
    url: `/v1/user/${id}`,
    method: 'put',
    data
  })
}

// 删除用户
export function deleteUser (id) {
  return request({
    url: `/v1/user/${id}`,
    method: 'delete'
  })
}

// 获取用户角色
export function getUserRoles (userId) {
  return request({
    url: `/v1/user/${userId}/roles`,
    method: 'get'
  })
}

// 分配用户角色
export function assignUserRoles (data) {
  return request({
    url: '/v1/user/assign-roles',
    method: 'post',
    data
  })
}

// 获取所有角色
export function getAllRoles () {
  return request({
    url: '/v1/role/all',
    method: 'get'
  })
}

// 生成扫码登录二维码
export function generateQrCode () {
  return request({
    url: '/v1/auth/qrcode',
    method: 'get'
  })
}

// 检查扫码登录状态
export function checkQrCodeStatus (qrCodeId) {
  return request({
    url: `/v1/auth/qrcode/${qrCodeId}/status`,
    method: 'get'
  })
}

// 获取用户Token
export function getUserToken (userId) {
  return request({
    url: `/v1/user/${userId}/token`,
    method: 'get'
  })
}

// 刷新用户Token
export function refreshUserToken (data) {
  return request({
    url: '/v1/user/refresh-token',
    method: 'post',
    data
  })
}

// 刷新用户Token并递增版本号
export function refreshUserTokenWithVersion (data) {
  return request({
    url: '/v1/user/refresh-token-with-version',
    method: 'post',
    data
  })
}