// 强制修复权限问题
console.log('=== 强制修复权限问题 ===');

// 1. 清除所有相关缓存
console.log('清除缓存...');
localStorage.removeItem('roleMenuPermissions');
localStorage.removeItem('roles');
localStorage.removeItem('permissions');
sessionStorage.clear();

// 2. 设置正确的权限配置
const correctPermissions = {
  admin: ['business', 'role-permission', 'system', 'system.user', 'system.role', 'system.role-permission', 'system.department', 'system.casbin', 'system.menu-permission', 'api', 'api.list'],
  user: ['api', 'api.list'],
  test_role: ['api', 'api.list'],
  member: ['business', 'role-permission', 'api', 'api.list'],
  shop: ['business', 'role-permission', 'api', 'api.list']
};

localStorage.setItem('roleMenuPermissions', JSON.stringify(correctPermissions));
console.log('✅ 权限配置已重置');

// 3. 检查当前用户角色
let currentRoles = [];
if (window.$store && window.$store.state.user.roles) {
  currentRoles = window.$store.state.user.roles;
  console.log('当前用户角色:', currentRoles);
} else {
  console.log('⚠️ 无法获取用户角色，将刷新页面');
}

// 4. 强制重新生成路由
if (window.$store && currentRoles.length > 0) {
  console.log('重新生成路由...');
  
  // 重置路由状态
  window.$store.commit('permission/SET_ROUTES', []);
  
  // 重新生成路由
  window.$store.dispatch('permission/generateRoutes', currentRoles).then((accessRoutes) => {
    console.log('✅ 路由重新生成成功');
    console.log('可访问路由数量:', accessRoutes.length);
    
    // 重新添加路由
    window.$router.addRoutes(accessRoutes);
    
    console.log('✅ 路由已重新添加');
    
    // 延迟刷新页面
    setTimeout(() => {
      console.log('刷新页面以应用更改...');
      window.location.reload();
    }, 1000);
  }).catch(err => {
    console.error('❌ 路由生成失败:', err);
    setTimeout(() => {
      window.location.reload();
    }, 1000);
  });
} else {
  console.log('直接刷新页面...');
  setTimeout(() => {
    window.location.reload();
  }, 1000);
}

console.log('=== 修复脚本执行完成 ===');