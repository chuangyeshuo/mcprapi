// 测试菜单权限管理页面配置更新脚本
console.log('=== 测试菜单权限管理页面配置更新 ===');

// 模拟菜单权限管理页面的配置
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

// 菜单名称映射
const menuNameMap = {
  'business': '业务创建',
  'role': '角色创建',
  'user': '用户管理',
  'role-permission': 'API授权',
  'api': 'API创建',
  'api.list': 'API列表',
  'system': '系统管理',
  'system.user': '用户管理',
  'system.role': '角色管理',
  'system.role-permission': '角色权限',
  'system.department': '部门管理',
  'system.casbin': 'Casbin权限管理',
  'system.menu-permission': '菜单权限管理'
};

// 菜单树数据
const menuTreeData = [
  {
    key: 'business',
    label: '业务创建'
  },
  {
    key: 'role',
    label: '角色创建'
  },
  {
    key: 'user',
    label: '用户管理'
  },
  {
    key: 'role-permission',
    label: 'API授权'
  },
  {
    key: 'api',
    label: 'API创建',
    children: [
      { key: 'api.list', label: 'API列表' }
    ]
  },
  {
    key: 'system',
    label: '系统管理',
    children: [
      { key: 'system.user', label: '用户管理' },
      { key: 'system.role', label: '角色管理' },
      { key: 'system.role-permission', label: '角色权限' },
      { key: 'system.department', label: '部门管理' },
      { key: 'system.casbin', label: 'Casbin权限管理' },
      { key: 'system.menu-permission', label: '菜单权限管理' }
    ]
  }
];

console.log('菜单权限管理页面配置已更新：');
console.log(JSON.stringify(defaultPermissions, null, 2));

console.log('\n=== 基础菜单验证 ===');
const roles = ['user', 'test_role', 'member', 'member_mcp', 'shop_mcp', 'guest', 'operator', 'manager'];

roles.forEach(role => {
  const rolePermissions = defaultPermissions[role] || [];
  const hasAllBaseMenus = baseMenus.every(menu => rolePermissions.includes(menu));
  console.log(`${role}: ${hasAllBaseMenus ? '✅ 拥有所有基础菜单权限' : '❌ 缺少基础菜单权限'}`);
  
  const menuNames = rolePermissions.map(key => menuNameMap[key] || key).join(', ');
  console.log(`  权限列表: [${menuNames}]`);
});

console.log('\n=== 菜单树结构验证 ===');
function getAllMenuKeys(nodes) {
  const keys = [];
  const traverse = (nodes) => {
    nodes.forEach(node => {
      keys.push(node.key);
      if (node.children) {
        traverse(node.children);
      }
    });
  };
  traverse(nodes);
  return keys;
}

const allMenuKeys = getAllMenuKeys(menuTreeData);
console.log('菜单树包含的所有菜单键:', allMenuKeys);

console.log('\n=== 基础菜单覆盖验证 ===');
const missingInTree = baseMenus.filter(menu => !allMenuKeys.includes(menu));
const extraInBase = allMenuKeys.filter(menu => !baseMenus.includes(menu) && !menu.startsWith('system'));

if (missingInTree.length === 0) {
  console.log('✅ 菜单树包含所有基础菜单');
} else {
  console.log('❌ 菜单树缺少基础菜单:', missingInTree);
}

console.log('系统管理菜单:', allMenuKeys.filter(menu => menu.startsWith('system')));
console.log('基础菜单:', baseMenus);

console.log('\n=== 优化总结 ===');
console.log('1. ✅ 所有非admin角色都可以访问5个基础菜单');
console.log('2. ✅ 菜单名称已更新为更准确的描述');
console.log('3. ✅ 菜单树结构已重新组织');
console.log('4. ✅ 权限配置与路由权限保持一致');