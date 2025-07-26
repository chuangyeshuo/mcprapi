import {
  asyncRoutes,
  constantRoutes,
} from '@/router';

/**
 * 使用meta.role判断当前用户是否有权限
 * @param roles
 * @param route
 */
function hasPermission (roles, route) {
  if (route.meta && route.meta.roles) {
    return roles.some(role => route.meta.roles.includes(role))
  } else {
    return true
  }
}

/**
 * 检查用户是否有菜单权限（基于localStorage配置）
 * @param roles
 * @param route
 */
function hasMenuPermission (roles, route) {
  // 获取菜单权限配置
  const roleMenuPermissions = JSON.parse(localStorage.getItem('roleMenuPermissions') || '{}')
  
  // 如果是admin角色，直接返回true
  if (roles.includes('admin')) {
    return true
  }
  
  // 定义所有非admin角色都可以访问的基础菜单
  const baseMenus = ['business', 'role', 'user', 'role-permission', 'api', 'api.list']
  
  // 构建菜单key
  let menuKey = ''
  
  // 根据路由名称匹配
  if (route.name === 'Business') {
    menuKey = 'business'
  } else if (route.name === 'Role') {
    menuKey = 'role'
  } else if (route.name === 'User') {
    menuKey = 'user'
  } else if (route.name === 'RolePermission') {
    menuKey = 'role-permission'
  } else if (route.name === 'API') {
    menuKey = 'api'
  } else if (route.name === 'APIList') {
    menuKey = 'api.list'
  } else if (route.name === 'System') {
    menuKey = 'system'
  } else if (route.name === 'User') {
    menuKey = 'system.user'
  } else if (route.name === 'Department') {
    menuKey = 'system.department'
  } else if (route.name === 'Casbin') {
    menuKey = 'system.casbin'
  } else if (route.name === 'MenuPermission') {
    menuKey = 'system.menu-permission'
  }
  
  // 如果根据名称没有匹配到，尝试根据路径匹配
  if (!menuKey) {
    if (route.path === '/business' || route.path === '/business/index') {
      menuKey = 'business'
    } else if (route.path === '/role' || route.path === '/role/index') {
      menuKey = 'role'
    } else if (route.path === '/user' || route.path === '/user/index') {
      menuKey = 'user'
    } else if (route.path === '/role-permission' || route.path === '/role-permission/index') {
      menuKey = 'role-permission'
    } else if (route.path === '/api' || route.path === '/api/index') {
      menuKey = 'api'
    } else if (route.path === '/api/list') {
      menuKey = 'api.list'
    } else if (route.path === '/system') {
      menuKey = 'system'
    } else if (route.path && route.path.includes('/system/')) {
      const subPath = route.path.replace('/system/', '').replace('/', '.')
      menuKey = 'system.' + subPath
    }
  }
  
  // 如果是基础菜单，所有非admin角色都有权限
  if (baseMenus.includes(menuKey)) {
    return true
  }
  
  // 如果没有配置菜单key，使用原有的角色权限检查
  if (!menuKey) {
    return hasPermission(roles, route)
  }
  
  // 检查用户角色是否有该菜单权限
  for (const role of roles) {
    const userMenus = roleMenuPermissions[role] || []
    if (userMenus.includes(menuKey)) {
      return true
    }
    // 对于子菜单，也检查父菜单权限
    if (menuKey.includes('.')) {
      const parentKey = menuKey.split('.')[0]
      if (userMenus.includes(parentKey)) {
        return true
      }
    }
  }
  
  return false
}

/**
 * 递归过滤异步路由表，返回符合用户角色权限的路由表
 * @param routes asyncRoutes
 * @param roles
 */
export function filterAsyncRoutes (routes, roles) {
  const res = []

  routes.forEach(route => {
    const tmp = { ...route }
    
    // 检查当前路由权限
    const hasCurrentPermission = hasMenuPermission(roles, tmp)
    
    if (hasCurrentPermission) {
      // 如果有子路由，递归过滤子路由
      if (tmp.children) {
        tmp.children = filterAsyncRoutes(tmp.children, roles)
        // 即使子路由被过滤完了，父路由仍然保留（可能有直接访问的情况）
      }
      res.push(tmp)
    } else {
      // 即使父路由没有权限，也要检查子路由
      if (tmp.children) {
        const filteredChildren = filterAsyncRoutes(tmp.children, roles)
        if (filteredChildren.length > 0) {
          tmp.children = filteredChildren
          res.push(tmp)
        }
      }
    }
  })

  return res
}

const state = {
  routes: [],
  addRoutes: []
}

const mutations = {
  SET_ROUTES: (state, routes) => {
    state.addRoutes = routes
    state.routes = constantRoutes.concat(routes)
  }
}

const actions = {
  generateRoutes ({ commit }, roles) {
    return new Promise(resolve => {
      // 初始化默认菜单权限配置
      const savedPermissions = localStorage.getItem('roleMenuPermissions')
      if (!savedPermissions) {
        // 定义所有非admin角色都可以访问的基础菜单
        const baseMenus = ['business', 'role', 'user', 'role-permission', 'api', 'api.list']
        
        const defaultPermissions = {
          admin: ['business', 'role', 'user', 'role-permission', 'system', 'system.user', 'system.role', 'system.role-permission', 'system.department', 'system.casbin', 'system.menu-permission', 'api', 'api.list'],
          user: baseMenus, // 普通用户可以访问所有基础菜单
          test_role: baseMenus, // 测试角色可以访问所有基础菜单
          member: baseMenus, // member角色可以访问所有基础菜单
          member_mcp: baseMenus, // 会员角色可以访问所有基础菜单
          shop_mcp: baseMenus, // 电商角色可以访问所有基础菜单
          guest: baseMenus, // 访客角色可以访问所有基础菜单
          operator: baseMenus, // 操作员角色可以访问所有基础菜单
          manager: baseMenus // 管理员角色可以访问所有基础菜单
        }
        localStorage.setItem('roleMenuPermissions', JSON.stringify(defaultPermissions))
      }
      
      let accessedRoutes
      if (roles.includes('admin')) {
        accessedRoutes = asyncRoutes || []
      } else {
        accessedRoutes = filterAsyncRoutes(asyncRoutes, roles)
      }
      commit('SET_ROUTES', accessedRoutes)
      resolve(accessedRoutes)
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}