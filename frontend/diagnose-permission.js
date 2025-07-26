// 全面诊断权限问题
console.log('=== 全面诊断权限问题 ===');

// 1. 检查当前用户信息
console.log('\n=== 当前用户信息 ===');
const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}');
console.log('用户信息:', userInfo);

// 2. 检查Vuex store中的用户角色
if (window.$store) {
  console.log('Vuex用户状态:', window.$store.state.user);
  console.log('用户角色:', window.$store.state.user.roles);
  console.log('当前路由:', window.$store.state.permission.routes);
} else {
  console.log('Vuex store未找到');
}

// 3. 检查当前权限配置
const currentPermissions = localStorage.getItem('roleMenuPermissions');
console.log('\n=== 当前权限配置 ===');
if (currentPermissions) {
  const permissions = JSON.parse(currentPermissions);
  console.log('权限配置:', permissions);
  console.log('shop_mcp权限:', permissions.shop_mcp);
  console.log('member_mcp权限:', permissions.member_mcp);
} else {
  console.log('未找到权限配置');
}

// 4. 强制重置权限配置
console.log('\n=== 强制重置权限配置 ===');
const correctPermissions = {
  admin: ['business', 'role-permission', 'system', 'system.user', 'system.role', 'system.role-permission', 'system.department', 'system.casbin', 'system.menu-permission', 'api', 'api.list'],
  user: ['api', 'api.list'],
  test_role: ['api', 'api.list'],
  member: ['business', 'role-permission', 'api', 'api.list'],
  shop: ['business', 'role-permission', 'api', 'api.list']
};

localStorage.setItem('roleMenuPermissions', JSON.stringify(correctPermissions));
console.log('已重置权限配置:', correctPermissions);

// 5. 模拟权限检查函数
function debugMenuPermission(roles, routeName, routePath) {
  console.log(`\n--- 检查路由权限 ---`);
  console.log(`路由名称: ${routeName}`);
  console.log(`路由路径: ${routePath}`);
  console.log(`用户角色: ${JSON.stringify(roles)}`);
  
  const roleMenuPermissions = JSON.parse(localStorage.getItem('roleMenuPermissions') || '{}');
  
  // 如果是admin角色，直接返回true
  if (roles.includes('admin')) {
    console.log('✅ Admin用户，直接通过');
    return true;
  }
  
  // 构建菜单key
  let menuKey = '';
  
  // 根据路由名称匹配
  if (routeName === 'Business') {
    menuKey = 'business';
  } else if (routeName === 'RolePermission') {
    menuKey = 'role-permission';
  } else if (routeName === 'API') {
    menuKey = 'api';
  } else if (routeName === 'APIList') {
    menuKey = 'api.list';
  } else if (routeName === 'System') {
    menuKey = 'system';
  }
  
  // 如果根据名称没有匹配到，尝试根据路径匹配
  if (!menuKey) {
    if (routePath === '/business' || routePath === '/business/index') {
      menuKey = 'business';
    } else if (routePath === '/role-permission' || routePath === '/role-permission/index') {
      menuKey = 'role-permission';
    } else if (routePath === '/api' || routePath === '/api/index') {
      menuKey = 'api';
    } else if (routePath === '/api/list') {
      menuKey = 'api.list';
    }
  }
  
  console.log(`菜单key: ${menuKey}`);
  
  // 检查用户角色是否有该菜单权限
  for (const role of roles) {
    const userMenus = roleMenuPermissions[role] || [];
    console.log(`角色 ${role} 的权限: ${JSON.stringify(userMenus)}`);
    
    if (userMenus.includes(menuKey)) {
      console.log(`✅ 角色 ${role} 有权限访问 ${menuKey}`);
      return true;
    }
    
    // 对于子菜单，也检查父菜单权限
    if (menuKey.includes('.')) {
      const parentKey = menuKey.split('.')[0];
      if (userMenus.includes(parentKey)) {
        console.log(`✅ 角色 ${role} 有父菜单权限 ${parentKey}`);
        return true;
      }
    }
  }
  
  console.log(`❌ 无权限访问 ${menuKey}`);
  return false;
}

// 6. 测试关键路由权限
console.log('\n=== 测试电商用户权限 ===');
const shopRoles = ['shop_mcp'];

debugMenuPermission(shopRoles, 'Business', '/business');
debugMenuPermission(shopRoles, 'RolePermission', '/role-permission');
debugMenuPermission(shopRoles, 'API', '/api');
debugMenuPermission(shopRoles, 'APIList', '/api/list');

console.log('\n=== 测试会员用户权限 ===');
const memberRoles = ['member_mcp'];

debugMenuPermission(memberRoles, 'Business', '/business');
debugMenuPermission(memberRoles, 'RolePermission', '/role-permission');
debugMenuPermission(memberRoles, 'API', '/api');
debugMenuPermission(memberRoles, 'APIList', '/api/list');

// 7. 强制重新生成路由
console.log('\n=== 强制重新生成路由 ===');
if (window.$store && window.$store.dispatch) {
  const currentRoles = window.$store.state.user.roles || [];
  console.log('当前用户角色:', currentRoles);
  
  window.$store.dispatch('permission/generateRoutes', currentRoles).then(() => {
    console.log('✅ 路由重新生成完成');
    console.log('新的路由:', window.$store.state.permission.routes);
    
    // 强制刷新页面
    setTimeout(() => {
      console.log('正在刷新页面...');
      window.location.reload();
    }, 1000);
  }).catch(err => {
    console.error('❌ 路由生成失败:', err);
  });
} else {
  console.log('❌ 无法访问Vuex store，请手动刷新页面');
  setTimeout(() => {
    window.location.reload();
  }, 2000);
}