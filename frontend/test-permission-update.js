// 测试权限配置更新脚本
console.log('=== 测试权限配置更新 ===');

// 模拟localStorage
const mockLocalStorage = {
  data: {},
  getItem: function(key) {
    return this.data[key] || null;
  },
  setItem: function(key, value) {
    this.data[key] = value;
  }
};

// 清空现有权限配置
mockLocalStorage.setItem('roleMenuPermissions', null);

// 模拟权限初始化逻辑
const savedPermissions = mockLocalStorage.getItem('roleMenuPermissions');
if (!savedPermissions) {
  // 定义所有非admin角色都可以访问的基础菜单
  const baseMenus = ['business', 'role', 'user', 'role-permission', 'api', 'api.list'];
  
  const defaultPermissions = {
    admin: ['business', 'role', 'user', 'role-permission', 'system', 'system.user', 'system.role', 'system.role-permission', 'system.department', 'system.casbin', 'system.menu-permission', 'api', 'api.list'],
    user: baseMenus, // 普通用户可以访问所有基础菜单
    test_role: baseMenus, // 测试角色可以访问所有基础菜单
    member: baseMenus, // member角色可以访问所有基础菜单
    member_mcp: baseMenus, // 会员角色可以访问所有基础菜单
    shop_mcp: baseMenus, // 电商角色可以访问所有基础菜单
    guest: baseMenus, // 访客角色可以访问所有基础菜单
    operator: baseMenus, // 操作员角色可以访问所有基础菜单
    manager: baseMenus // 管理员角色可以访问所有基础菜单
  };
  mockLocalStorage.setItem('roleMenuPermissions', JSON.stringify(defaultPermissions));
}

// 验证权限配置
const permissions = JSON.parse(mockLocalStorage.getItem('roleMenuPermissions'));
console.log('权限配置已更新：');
console.log(JSON.stringify(permissions, null, 2));

// 验证每个角色的权限
console.log('\n=== 权限验证结果 ===');
const baseMenus = ['business', 'role', 'user', 'role-permission', 'api', 'api.list'];
const roles = ['user', 'test_role', 'member', 'member_mcp', 'shop_mcp', 'guest', 'operator', 'manager'];

roles.forEach(role => {
  const rolePermissions = permissions[role] || [];
  const hasAllBaseMenus = baseMenus.every(menu => rolePermissions.includes(menu));
  console.log(`${role}: ${hasAllBaseMenus ? '✅ 拥有所有基础菜单权限' : '❌ 缺少基础菜单权限'}`);
  console.log(`  权限列表: [${rolePermissions.join(', ')}]`);
});

console.log('\n=== 基础菜单说明 ===');
console.log('业务创建: business');
console.log('角色创建: role');
console.log('用户管理: user');
console.log('API创建: api, api.list');
console.log('API授权: role-permission');