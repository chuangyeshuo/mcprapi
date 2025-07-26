// 测试会员用户API列表权限修复脚本
console.log('=== 会员用户API列表权限修复测试 ===');

// 1. 清理现有权限配置
localStorage.removeItem('roleMenuPermissions');
console.log('✓ 已清理现有权限配置');

// 2. 设置新的权限配置
const newPermissions = {
  admin: ['business', 'role-permission', 'system', 'system.user', 'system.role', 'system.role-permission', 'system.department', 'system.casbin', 'system.menu-permission', 'api', 'api.list'],
  user: ['api', 'api.list'],
  test_role: ['api', 'api.list'],
  member: ['business', 'role-permission', 'api', 'api.list'], // 会员角色权限
  shop: ['business', 'role-permission', 'api', 'api.list']
};

localStorage.setItem('roleMenuPermissions', JSON.stringify(newPermissions));
console.log('✓ 已设置新的权限配置');

// 3. 验证会员角色权限
const memberPermissions = newPermissions.member_mcp;
console.log('会员角色(member_mcp)的权限:', memberPermissions);

// 4. 检查API相关权限
const hasAPIAccess = memberPermissions.includes('api');
const hasAPIListAccess = memberPermissions.includes('api.list');

console.log('API管理权限:', hasAPIAccess ? '✓ 有权限' : '✗ 无权限');
console.log('API列表权限:', hasAPIListAccess ? '✓ 有权限' : '✗ 无权限');

// 5. 检查业务线管理权限
const hasBusinessAccess = memberPermissions.includes('business');
console.log('业务线管理权限:', hasBusinessAccess ? '✓ 有权限' : '✗ 无权限');

// 6. 检查角色权限管理权限
const hasRolePermissionAccess = memberPermissions.includes('role-permission');
console.log('角色权限管理权限:', hasRolePermissionAccess ? '✓ 有权限' : '✗ 无权限');

// 7. 模拟菜单权限检查函数
function hasMenuPermission(roles, menuKey) {
  const roleMenuPermissions = JSON.parse(localStorage.getItem('roleMenuPermissions') || '{}');
  
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

// 8. 测试菜单权限检查
const memberRoles = ['member_mcp'];
console.log('\n=== 菜单权限检查测试 ===');
console.log('API管理菜单:', hasMenuPermission(memberRoles, 'api') ? '✓ 可访问' : '✗ 不可访问');
console.log('API列表菜单:', hasMenuPermission(memberRoles, 'api.list') ? '✓ 可访问' : '✗ 不可访问');
console.log('业务线管理菜单:', hasMenuPermission(memberRoles, 'business') ? '✓ 可访问' : '✗ 不可访问');
console.log('角色权限管理菜单:', hasMenuPermission(memberRoles, 'role-permission') ? '✓ 可访问' : '✗ 不可访问');

console.log('\n=== 修复完成 ===');
console.log('请刷新页面以应用新的权限配置');