<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input
        v-model="listQuery.query"
        placeholder="搜索角色"
        style="width: 200px;"
        class="filter-item"
        @keyup.enter.native="handleFilter"
      />
      <el-button v-waves class="filter-item" type="primary" icon="el-icon-search" @click="handleFilter">
        搜索
      </el-button>
      <el-button class="filter-item" style="margin-left: 10px;" type="primary" icon="el-icon-edit" @click="handleCreate">
        新增角色
      </el-button>
    </div>

    <el-table
      :key="tableKey"
      v-loading="listLoading"
      :data="list"
      border
      fit
      highlight-current-row
      style="width: 100%;"
    >
      <el-table-column label="ID" prop="id" sortable="custom" align="center" width="80">
        <template slot-scope="{row}">
          <span>{{ row.id }}</span>
        </template>
      </el-table-column>
      <el-table-column label="角色名称" width="150px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.name }}</span>
        </template>
      </el-table-column>
      <el-table-column label="角色代码" width="150px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.code }}</span>
        </template>
      </el-table-column>
      <el-table-column label="描述" min-width="200px">
        <template slot-scope="{row}">
          <span>{{ row.description }}</span>
        </template>
      </el-table-column>
      <el-table-column label="可访问菜单" min-width="300px">
        <template slot-scope="{row}">
          <el-tag
            v-for="menu in getAccessibleMenus(row.code)"
            :key="menu"
            type="success"
            size="mini"
            style="margin-right: 5px; margin-bottom: 5px;"
          >
            {{ getMenuName(menu) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="230" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <el-button type="primary" size="mini" @click="handleUpdate(row)">
            配置权限
          </el-button>
          <el-button v-if="row.code !== 'admin'" size="mini" type="danger" @click="handleDelete(row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />

    <!-- 权限配置对话框 -->
    <el-dialog :title="dialogTitle" :visible.sync="dialogFormVisible" width="60%">
      <el-form ref="dataForm" :model="temp" label-position="left" label-width="100px" style="width: 100%;">
        <el-form-item label="角色名称">
          <el-input v-model="temp.name" disabled />
        </el-form-item>
        <el-form-item label="角色代码">
          <el-input v-model="temp.code" disabled />
        </el-form-item>
        <el-form-item label="菜单权限">
          <el-tree
            ref="menuTree"
            :data="menuTreeData"
            :props="menuTreeProps"
            show-checkbox
            node-key="key"
            :default-checked-keys="temp.menus || []"
            :default-expand-all="true"
          />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">
          取消
        </el-button>
        <el-button type="primary" @click="updateMenuPermission">
          确定
        </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import {
  getRoleMenuPermission,
  updateRoleMenuPermission,
} from '@/api/menu-permission';
import { getAllRoles } from '@/api/user';
import Pagination from '@/components/Pagination';
import waves from '@/directive/waves';

export default {
  name: 'MenuPermission',
  components: { Pagination },
  directives: { waves },
  data() {
    return {
      tableKey: 0,
      list: [],
      total: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 20,
        query: ''
      },
      temp: {
        id: undefined,
        name: '',
        code: '',
        menus: []
      },
      dialogFormVisible: false,
      dialogTitle: '',
      roleMenuPermissions: {}, // 存储角色的菜单权限
      menuTreeData: [
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
      ],
      menuTreeProps: {
        children: 'children',
        label: 'label'
      },
      menuNameMap: {
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
      }
    }
  },
  created() {
    this.getList()
    this.loadRoleMenuPermissions()
  },
  methods: {
    getList() {
      this.listLoading = true
      getAllRoles().then(response => {
        this.list = response.data
        this.total = response.data.length
        this.listLoading = false
      })
    },
    loadRoleMenuPermissions() {
      // 从localStorage加载角色菜单权限配置
      const saved = localStorage.getItem('roleMenuPermissions')
      if (saved) {
        this.roleMenuPermissions = JSON.parse(saved)
      } else {
        // 定义所有非admin角色都可以访问的基础菜单
        const baseMenus = ['business', 'role', 'user', 'role-permission', 'api', 'api.list']
        
        // 默认配置：admin可以访问所有菜单，其他角色可以访问基础菜单
        this.roleMenuPermissions = {
          admin: ['business', 'role', 'user', 'role-permission', 'system', 'system.user', 'system.role', 'system.role-permission', 'system.department', 'system.casbin', 'system.menu-permission', 'api', 'api.list'],
          user: baseMenus, // 普通用户可以访问所有基础菜单
          test_role: baseMenus, // 测试角色可以访问所有基础菜单
          member: baseMenus, // member角色可以访问所有基础菜单
          member_mcp: baseMenus, // 会员角色可以访问所有基础菜单
          shop_mcp: baseMenus, // 电商角色可以访问所有基础菜单
          guest: baseMenus, // 访客角色可以访问所有基础菜单
          operator: baseMenus, // 操作员角色可以访问所有基础菜单
          manager: baseMenus // 管理员角色可以访问所有基础菜单
        }
        this.saveRoleMenuPermissions()
      }
    },
    saveRoleMenuPermissions() {
      localStorage.setItem('roleMenuPermissions', JSON.stringify(this.roleMenuPermissions))
    },
    getAccessibleMenus(roleCode) {
      return this.roleMenuPermissions[roleCode] || []
    },
    getMenuName(menuKey) {
      return this.menuNameMap[menuKey] || menuKey
    },
    handleFilter() {
      this.listQuery.page = 1
      this.getList()
    },
    handleCreate() {
      this.$message.info('请在角色管理页面创建新角色')
    },
    handleUpdate(row) {
      this.temp = Object.assign({}, row)
      this.temp.menus = this.roleMenuPermissions[row.code] || []
      this.dialogTitle = '配置菜单权限'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs.menuTree.setCheckedKeys(this.temp.menus)
      })
    },
    updateMenuPermission() {
      // 只获取完全选中的节点，不包含半选中的节点
      const checkedKeys = this.$refs.menuTree.getCheckedKeys()
      
      // 过滤掉可能的重复项并排序
      const uniqueCheckedKeys = [...new Set(checkedKeys)].sort()
      
      // 验证选中的权限是否都是有效的菜单项
      const validMenuKeys = this.getAllValidMenuKeys()
      const validCheckedKeys = uniqueCheckedKeys.filter(key => validMenuKeys.includes(key))
      
      // 如果有无效的权限，给出警告
      if (validCheckedKeys.length !== uniqueCheckedKeys.length) {
        const invalidKeys = uniqueCheckedKeys.filter(key => !validMenuKeys.includes(key))
        console.warn('发现无效的菜单权限:', invalidKeys)
      }
      
      // 调试信息
      console.log('角色:', this.temp.code)
      console.log('选中的菜单权限:', validCheckedKeys)
      console.log('权限数量:', validCheckedKeys.length)
      
      // 保存权限
      this.roleMenuPermissions[this.temp.code] = validCheckedKeys
      this.saveRoleMenuPermissions()
      
      this.dialogFormVisible = false
      
      // 显示更详细的成功信息
      const menuNames = validCheckedKeys.map(key => this.getMenuName(key)).join('、')
      if (validCheckedKeys.length > 0) {
        this.$message.success(`菜单权限配置成功！已为角色"${this.temp.name}"分配${validCheckedKeys.length}个权限：${menuNames}`)
      } else {
        this.$message.success(`菜单权限配置成功！已清空角色"${this.temp.name}"的所有菜单权限`)
      }
      
      // 如果当前用户的角色权限被修改，需要刷新页面
      const currentUserRoles = this.$store.getters.roles
      if (currentUserRoles.includes(this.temp.code)) {
        this.$message.info('您的权限已更新，页面将在3秒后刷新')
        setTimeout(() => {
          location.reload()
        }, 3000)
      }
    },
    // 获取所有有效的菜单键值
    getAllValidMenuKeys() {
      const keys = []
      const traverse = (nodes) => {
        nodes.forEach(node => {
          keys.push(node.key)
          if (node.children) {
            traverse(node.children)
          }
        })
      }
      traverse(this.menuTreeData)
      return keys
    },
    handleDelete(row) {
      this.$confirm('确定删除该角色吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.$message.info('请在角色管理页面删除角色')
      })
    }
  }
}
</script>

<style scoped>
.filter-container {
  padding-bottom: 10px;
}
.filter-item {
  display: inline-block;
  vertical-align: middle;
  margin-bottom: 10px;
}
</style>