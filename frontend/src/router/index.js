import Vue from 'vue';

import VueRouter from 'vue-router';

import Layout from '@/layout';

Vue.use(VueRouter)

// 公共路由
export const constantRoutes = [
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },
  {
    path: '/404',
    component: () => import('@/views/error/404'),
    hidden: true
  },
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        component: () => import('@/views/dashboard/index'),
        name: 'Dashboard',
        meta: { title: '首页', icon: 'dashboard', affix: true }
      }
    ]
  }
]

// 动态路由，基于用户权限动态加载
export const asyncRoutes = [
  // 业务线管理 - 独立菜单
  {
    path: '/business',
    component: Layout,
    children: [
      {
        path: 'index',
        component: () => import('@/views/system/business/index'),
        name: 'Business',
        meta: { title: '业务创建', icon: 'component' }
      }
    ]
  },
  // 角色管理 - 独立菜单，与业务线管理同级
  {
    path: '/role',
    component: Layout,
    children: [
      {
        path: 'index',
        component: () => import('@/views/system/role/index'),
        name: 'Role',
        meta: { title: '角色创建', icon: 'peoples' }
      }
    ]
  },
  // 用户管理 - 独立菜单，移动到角色管理后面
  {
    path: '/user',
    component: Layout,
    children: [
      {
        path: 'index',
        component: () => import('@/views/system/user/index'),
        name: 'User',
        meta: { title: '用户管理', icon: 'user' }
      }
    ]
  },
  // API管理 - 移动到业务线管理后面
  {
    path: '/api',
    component: Layout,
    redirect: '/api/list',
    name: 'API',
    meta: { title: 'API创建', icon: 'api' },
    children: [
      {
        path: 'list',
        component: () => import('@/views/api/list/index'),
        name: 'APIList',
        meta: { title: 'API创建', icon: 'list' }
      }
    ]
  },
  // 角色权限管理 - 独立菜单，与业务线管理同级
  {
    path: '/role-permission',
    component: Layout,
    children: [
      {
        path: 'index',
        component: () => import('@/views/system/role/permission'),
        name: 'RolePermission',
        meta: { title: 'API授权', icon: 'lock' }
      }
    ]
  },
  // 系统管理 - 仅包含管理员功能
  {
    path: '/system',
    component: Layout,
    redirect: '/system/department',
    name: 'System',
    meta: { title: '系统管理', icon: 'setting' },
    children: [
      {
        path: 'department',
        component: () => import('@/views/system/department/index'),
        name: 'Department',
        meta: { title: '部门管理', icon: 'tree' }
      },
      {
        path: 'casbin',
        component: () => import('@/views/system/casbin/index'),
        name: 'Casbin',
        meta: { title: 'Casbin权限管理', icon: 'lock' }
      },
      {
        path: 'menu-permission',
        component: () => import('@/views/system/menu-permission/index'),
        name: 'MenuPermission',
        meta: { title: '菜单权限管理', icon: 'tree-table', roles: ['admin'] }
      }
    ]
  },
  // 404页面必须放在最后
  { path: '*', redirect: '/404', hidden: true }
]

const createRouter = () => new VueRouter({
  mode: 'history',
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRoutes
})

const router = createRouter()

// 解决Vue Router重复导航错误
const originalPush = VueRouter.prototype.push
VueRouter.prototype.push = function push(location) {
  return originalPush.call(this, location).catch(err => {
    // 忽略NavigationDuplicated错误，但记录其他错误
    if (err.name !== 'NavigationDuplicated') {
      console.error('Router push error:', err)
    }
    return Promise.resolve(false) // 返回resolved promise以避免未捕获的错误
  })
}

const originalReplace = VueRouter.prototype.replace
VueRouter.prototype.replace = function replace(location) {
  return originalReplace.call(this, location).catch(err => {
    // 忽略NavigationDuplicated错误，但记录其他错误
    if (err.name !== 'NavigationDuplicated') {
      console.error('Router replace error:', err)
    }
    return Promise.resolve(false) // 返回resolved promise以避免未捕获的错误
  })
}

// 重置路由
export function resetRouter () {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher
}

export default router