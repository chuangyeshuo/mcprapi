// 清理并重置权限配置
console.log('=== 清理并重置电商角色权限配置 ===');

// 清除现有的权限配置
localStorage.removeItem('roleMenuPermissions');
console.log('已清除现有权限配置');

// 设置新的权限配置
const newPermissions = {
  admin: ['system', 'system.user', 'system.role', 'system.role-permission', 'system.department', 'system.business', 'system.casbin', 'system.menu-permission', 'api', 'api.list'],
  user: ['api', 'api.list'],
  test_role: ['api', 'api.list'],
  member_mcp: ['system', 'system.business', 'api', 'api.list'], // 会员角色
  shop_mcp: ['system', 'system.business', 'api', 'api.list'] // 电商角色
};

localStorage.setItem('roleMenuPermissions', JSON.stringify(newPermissions));
console.log('已设置新的权限配置');

// 验证电商角色权限
const shopPermissions = newPermissions.shop_mcp;
console.log('电商角色(shop_mcp)的权限:', shopPermissions);

// 检查关键权限
console.log('检查权限:');
console.log('- 系统管理权限:', shopPermissions.includes('system'));
console.log('- 业务线管理权限:', shopPermissions.includes('system.business'));
console.log('- API管理权限:', shopPermissions.includes('api'));
console.log('- API列表权限:', shopPermissions.includes('api.list'));

console.log('\n=== 权限配置完成 ===');
console.log('电商用户现在应该能够访问:');
console.log('1. 系统管理菜单');
console.log('2. 业务线管理页面');
console.log('3. API管理');
console.log('4. API列表');

console.log('\n请刷新页面以应用新的权限配置！');