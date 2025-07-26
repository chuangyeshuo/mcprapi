import 'element-ui/lib/theme-chalk/index.css';
import './assets/styles/index.scss';
import './icons'; // svg图标
import './permission'; // 权限控制

import Vue from 'vue';

import ElementUI from 'element-ui';

import { parseTime } from '@/utils';

import App from './App.vue';
import router from './router';
import store from './store';

Vue.use(ElementUI, { size: 'medium' })

// 注册全局过滤器
Vue.filter('parseTime', (time, cFormat) => {
  return parseTime(time, cFormat)
})

// 专门用于格式化创建时间和更新时间的过滤器
Vue.filter('formatDateTime', (time) => {
  if (!time) return '-'
  
  try {
    // 直接使用Date构造函数解析ISO 8601格式
    const date = new Date(time)
    
    // 检查日期是否有效
    if (isNaN(date.getTime())) {
      return '-'
    }
    
    // 格式化为 YYYY-MM-DD HH:mm
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')
    
    return `${year}-${month}-${day} ${hours}:${minutes}`
  } catch (error) {
    console.warn('时间格式化失败:', time, error)
    return '-'
  }
})

Vue.config.productionTip = false

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')