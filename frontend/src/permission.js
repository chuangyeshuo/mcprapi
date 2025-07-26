import 'nprogress/nprogress.css';

import { Message } from 'element-ui';
import NProgress from 'nprogress';

import { getToken } from '@/utils/auth';

import router from './router';
import store from './store';

NProgress.configure({ showSpinner: false })

const whiteList = ['/login'] // 不重定向白名单

router.beforeEach(async (to, from, next) => {
  // 开始进度条
  NProgress.start()

  // 确定用户是否已登录
  const hasToken = getToken()

  if (hasToken) {
    if (to.path === '/login') {
      // 如果已登录，重定向到首页
      next({ path: '/' })
      NProgress.done()
      return
    }
    
    // 确定用户是否已获取其权限角色
    const hasRoles = store.getters.roles && store.getters.roles.length > 0
    if (hasRoles) {
      next()
      return
    }
    
    try {
      // 获取用户信息
      const { roles } = await store.dispatch('user/getInfo')

      // 根据角色生成可访问路由图
      const accessRoutes = await store.dispatch('permission/generateRoutes', roles)

      // 动态添加可访问路由
      router.addRoutes(accessRoutes)

      // 确保addRoutes已完成，使用replace避免历史记录堆积
      // hack方法 确保addRoutes已完成
      next({ ...to, replace: true })
      return
    } catch (error) {
      // 移除token并转到登录页重新登录
      await store.dispatch('user/resetToken')
      Message.error(error || '验证失败，请重新登录')
      // 重定向到登录页面
      next(`/login?redirect=${to.path}`)
      NProgress.done()
      return
    }
  }
  
  /* 没有token */
  if (whiteList.indexOf(to.path) !== -1) {
    // 在免登录白名单中，直接进入
    next()
    return
  }
  
  // 其他无权访问的页面将重定向到登录页面
  // 避免重复重定向
  if (to.path !== '/login') {
    next(`/login?redirect=${to.path}`)
  } else {
    next()
  }
  NProgress.done()
})

router.afterEach(() => {
  // 完成进度条
  NProgress.done()
})