import request from '@/utils/request'

// 获取角色菜单权限
export function getRoleMenuPermission(roleCode) {
  // 从localStorage获取权限配置
  const permissions = localStorage.getItem('roleMenuPermissions')
  if (permissions) {
    const parsed = JSON.parse(permissions)
    return Promise.resolve({
      code: 0,
      data: parsed[roleCode] || []
    })
  }
  return Promise.resolve({
    code: 0,
    data: []
  })
}

// 更新角色菜单权限
export function updateRoleMenuPermission(roleCode, menus) {
  // 保存到localStorage
  const permissions = localStorage.getItem('roleMenuPermissions')
  let parsed = {}
  if (permissions) {
    parsed = JSON.parse(permissions)
  }
  parsed[roleCode] = menus
  localStorage.setItem('roleMenuPermissions', JSON.stringify(parsed))
  
  return Promise.resolve({
    code: 0,
    message: '更新成功'
  })
}

// 获取所有角色的菜单权限配置
export function getAllRoleMenuPermissions() {
  const permissions = localStorage.getItem('roleMenuPermissions')
  if (permissions) {
    return Promise.resolve({
      code: 0,
      data: JSON.parse(permissions)
    })
  }
  return Promise.resolve({
    code: 0,
    data: {}
  })
}