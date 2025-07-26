/**
 * 业务线管理部门权限控制测试脚本
 * 
 * 测试场景：
 * 1. admin角色：可以查看和操作所有部门的业务线
 * 2. 非admin角色：只能查看和操作自己所属部门的业务线
 * 
 * 使用方法：
 * 1. 在浏览器控制台中运行此脚本
 * 2. 脚本会模拟不同角色和部门的权限控制
 */

console.log('=== 业务线管理部门权限控制测试 ===');

// 清理现有的用户信息和权限配置
localStorage.removeItem('vue_admin_template_token');
localStorage.removeItem('roleMenuPermissions');

// 模拟用户信息
const mockUsers = {
  admin: {
    userId: 'admin001',
    name: '系统管理员',
    roles: ['admin'],
    deptId: 1,
    permissions: ['business', 'api', 'api.list']
  },
  shop_user: {
    userId: 'shop001', 
    name: '电商部门用户',
    roles: ['shop_mcp'],
    deptId: 2,
    permissions: ['business', 'api', 'api.list']
  },
  member_user: {
    userId: 'member001',
    name: '会员部门用户', 
    roles: ['member_mcp'],
    deptId: 3,
    permissions: ['business', 'api', 'api.list']
  }
};

// 模拟业务线数据
const mockBusinessList = [
  { id: 1, name: '用户中心', code: 'user-center', dept_id: 1, owner: '张三', status: 1 },
  { id: 2, name: '电商平台', code: 'ecommerce', dept_id: 2, owner: '李四', status: 1 },
  { id: 3, name: '会员系统', code: 'member-system', dept_id: 3, owner: '王五', status: 1 },
  { id: 4, name: '支付系统', code: 'payment', dept_id: 1, owner: '赵六', status: 1 }
];

// 模拟部门数据
const mockDepartments = [
  { id: 1, name: '技术部' },
  { id: 2, name: '电商部' },
  { id: 3, name: '会员部' }
];

// 权限检查函数
function checkBusinessPermission(userRole, userDeptId, businessDeptId) {
  // admin角色可以访问所有业务线
  if (userRole.includes('admin')) {
    return true;
  }
  
  // 非admin角色只能访问自己部门的业务线
  return userDeptId === businessDeptId;
}

// 过滤业务线列表
function filterBusinessList(userRole, userDeptId, businessList) {
  if (userRole.includes('admin')) {
    return businessList; // admin可以看到所有业务线
  }
  
  return businessList.filter(business => business.dept_id === userDeptId);
}

// 测试不同角色的权限
function testRolePermissions() {
  console.log('\n--- 测试不同角色的业务线访问权限 ---');
  
  Object.entries(mockUsers).forEach(([roleKey, user]) => {
    console.log(`\n${user.name} (${user.roles.join(', ')}) - 部门ID: ${user.deptId}`);
    
    // 获取该用户可以看到的业务线
    const visibleBusinessList = filterBusinessList(user.roles, user.deptId, mockBusinessList);
    console.log('可见的业务线：');
    visibleBusinessList.forEach(business => {
      const deptName = mockDepartments.find(d => d.id === business.dept_id)?.name || '未知部门';
      console.log(`  - ${business.name} (${business.code}) - ${deptName}`);
    });
    
    // 测试操作权限
    console.log('操作权限测试：');
    mockBusinessList.forEach(business => {
      const canOperate = checkBusinessPermission(user.roles, user.deptId, business.dept_id);
      const deptName = mockDepartments.find(d => d.id === business.dept_id)?.name || '未知部门';
      console.log(`  - ${business.name} (${deptName}): ${canOperate ? '✅ 可操作' : '❌ 无权限'}`);
    });
  });
}

// 模拟API请求参数
function simulateApiParams() {
  console.log('\n--- 模拟API请求参数 ---');
  
  Object.entries(mockUsers).forEach(([roleKey, user]) => {
    console.log(`\n${user.name} 的API请求参数：`);
    
    const baseParams = {
      page: 1,
      limit: 20,
      name: undefined
    };
    
    let apiParams;
    if (user.roles.includes('admin')) {
      apiParams = baseParams; // admin不需要额外过滤
      console.log('  GET /v1/business/list', JSON.stringify(apiParams));
    } else {
      apiParams = {
        ...baseParams,
        dept_id: user.deptId // 非admin需要添加部门过滤
      };
      console.log('  GET /v1/business/list', JSON.stringify(apiParams));
    }
  });
}

// 设置测试用户信息到localStorage
function setTestUser(userKey) {
  const user = mockUsers[userKey];
  if (!user) {
    console.error('用户不存在:', userKey);
    return;
  }
  
  // 设置用户token（模拟登录）
  localStorage.setItem('vue_admin_template_token', `mock_token_${user.userId}`);
  
  // 设置权限配置
  const rolePermissions = {};
  Object.keys(mockUsers).forEach(key => {
    const u = mockUsers[key];
    rolePermissions[u.roles[0]] = u.permissions;
  });
  localStorage.setItem('roleMenuPermissions', JSON.stringify(rolePermissions));
  
  console.log(`\n✅ 已设置测试用户: ${user.name}`);
  console.log('请刷新页面查看效果，或手动调用以下方法测试：');
  console.log('- this.$store.dispatch("user/getInfo") // 获取用户信息');
  console.log('- this.getList() // 重新获取业务线列表');
}

// 运行测试
testRolePermissions();
simulateApiParams();

console.log('\n=== 测试完成 ===');
console.log('\n可用的测试方法：');
console.log('- setTestUser("admin") // 设置为管理员用户');
console.log('- setTestUser("shop_user") // 设置为电商部门用户'); 
console.log('- setTestUser("member_user") // 设置为会员部门用户');

// 导出测试方法到全局
window.setTestUser = setTestUser;
window.testBusinessPermission = testRolePermissions;

console.log('\n💡 提示：');
console.log('1. 运行 setTestUser("shop_user") 切换到电商部门用户');
console.log('2. 刷新页面查看业务线管理页面的权限控制效果');
console.log('3. 尝试编辑/删除不同部门的业务线，观察权限限制');