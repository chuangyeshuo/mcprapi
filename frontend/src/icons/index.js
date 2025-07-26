import Vue from 'vue'
import SvgIcon from '@/components/SvgIcon' // svg组件

// 全局注册
Vue.component('svg-icon', SvgIcon)

// 自动导入所有svg图标
const req = require.context('./svg', false, /\.svg$/)
const requireAll = requireContext => requireContext.keys().map(requireContext)
requireAll(req)