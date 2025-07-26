// 调试会员用户权限配置
console.log('=== 调试会员用户权限配置 ===');

// 1. 检查当前localStorage中的权限配置
const currentPermissions = localStorage.getItem('roleMenuPermissions');
console.log('当前权限配置:', currentPermissions);

if (currentPermissions) {
  const permissions = JSON.parse(currentPermissions);
  console.log('解析后的权限配置:', permissions);
  console.log('member权限:', permissions.member);
  console.log('shop权限:', permissions.shop);
}

// 2. 清除所有相关的localStorage数据
console.log('\n=== 清除旧的权限配置 ===');
localStorage.removeItem('roleMenuPermissions');
localStorage.removeItem('roles');
localStorage.removeItem('permissions');

// 3. 重置为正确的权限配置
const correctPermissions = {
  admin: ['business', 'role-permission', 'system', 'system.user', 'system.role', 'system.role-permission', 'system.department', 'system.casbin', 'system.menu-permission', 'api', 'api.list'],
  user: ['api', 'api.list'],
  test_role: ['api', 'api.list'],
  member: ['business', 'role-permission', 'api', 'api.list'],
  shop: ['business', 'role-permission', 'api', 'api.list']
};

console.log('设置正确的权限配置:', correctPermissions);
localStorage.setItem('roleMenuPermissions', JSON.stringify(correctPermissions));

// 4. 验证设置是否成功
const verifyPermissions = JSON.parse(localStorage.getItem('roleMenuPermissions'));
console.log('验证设置后的权限配置:', verifyPermissions);

// 5. 模拟权限检查函数（与前端代码保持一致）
function testMenuPermission(roles, routeName, routePath) {
  const roleMenuPermissions = JSON.parse(localStorage.getItem('roleMenuPermissions') || '{}');
  
  // 如果是admin角色，直接返回true
  if (roles.includes('admin')) {
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
  
  console.log(`路由 ${routeName}(${routePath}) -> 菜单key: ${menuKey}`);
  
  // 检查用户角色是否有该菜单权限
  for (const role of roles) {
    const userMenus = roleMenuPermissions[role] || [];
    if (userMenus.includes(menuKey)) {
      return true;
    }
    // 对于子菜单，也检查父菜单权限
    if (menuKey.includes('.')) {
      const parentKey = menuKey.split('.')[0];
      if (userMenus.includes(parentKey)) {
        return true;
      }
    }
  }
  
  return false;
}

// 6. 测试会员用户的菜单权限
console.log('\n=== 测试会员用户菜单权限 ===');
const memberRoles = ['member_mcp'];

console.log('业务线管理权限:', testMenuPermission(memberRoles, 'Business', '/business/index'));
console.log('角色权限管理权限:', testMenuPermission(memberRoles, 'RolePermission', '/role-permission/index'));
console.log('API管理权限:', testMenuPermission(memberRoles, 'API', '/api'));
console.log('API列表权限:', testMenuPermission(memberRoles, 'APIList', '/api/list'));

// 7. 测试电商用户的菜单权限
console.log('\n=== 测试电商用户菜单权限 ===');
const shopRoles = ['shop_mcp'];

console.log('业务线管理权限:', testMenuPermission(shopRoles, 'Business', '/business/index'));
console.log('角色权限管理权限:', testMenuPermission(shopRoles, 'RolePermission', '/role-permission/index'));
console.log('API管理权限:', testMenuPermission(shopRoles, 'API', '/api'));
console.log('API列表权限:', testMenuPermission(shopRoles, 'APIList', '/api/list'));

// 8. 强制刷新页面以应用新配置
console.log('\n=== 完成权限配置重置 ===');
console.log('正在刷新页面以应用新的权限配置...');

// 延迟刷新，让用户看到日志
setTimeout(() => {
  window.location.reload();
}, 2000);