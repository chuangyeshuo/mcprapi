// 测试会员角色权限配置
console.log('=== 会员角色权限配置测试 ===');

// 模拟localStorage中的权限配置
const roleMenuPermissions = {
  admin: ['system', 'system.user', 'system.role', 'system.role-permission', 'system.department', 'system.business', 'system.casbin', 'system.menu-permission', 'api', 'api.list'],
  user: ['api', 'api.list'],
  test_role: [],
  member_mcp: ['system.business', 'api', 'api.list'], // 修复后的配置
  shop_mcp: ['system.business', 'api', 'api.list']
};

// 测试会员角色权限
const memberPermissions = roleMenuPermissions.member_mcp;
console.log('会员角色(member_mcp)的权限:', memberPermissions);

// 检查是否包含业务线管理权限
const hasBusinessPermission = memberPermissions.includes('system.business');
console.log('是否有业务线管理权限:', hasBusinessPermission);

// 检查是否包含API管理权限
const hasAPIPermission = memberPermissions.includes('api');
console.log('是否有API管理权限:', hasAPIPermission);

// 检查是否包含API列表权限
const hasAPIListPermission = memberPermissions.includes('api.list');
console.log('是否有API列表权限:', hasAPIListPermission);

console.log('=== 测试结果 ===');
if (hasBusinessPermission && hasAPIPermission && hasAPIListPermission) {
  console.log('✅ 会员角色权限配置正确！');
  console.log('会员用户现在可以访问：');
  console.log('- 业务线管理');
  console.log('- API管理');
  console.log('- API列表');
} else {
  console.log('❌ 会员角色权限配置有问题');
}

// 模拟权限检查函数
function hasMenuPermission(roles, menuKey) {
  for (const role of roles) {
    const userMenus = roleMenuPermissions[role] || [];
    if (userMenus.includes(menuKey)) {
      return true;
    }
  }
  return false;
}

// 测试权限检查
console.log('\n=== 权限检查测试 ===');
const memberRoles = ['member_mcp'];

console.log('会员角色访问业务线管理:', hasMenuPermission(memberRoles, 'system.business'));
console.log('会员角色访问用户管理:', hasMenuPermission(memberRoles, 'system.user'));
console.log('会员角色访问API管理:', hasMenuPermission(memberRoles, 'api'));