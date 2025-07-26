<template>
  <div class="role-permission-container">
    <div class="header">
      <h2>角色权限管理</h2>
      <p class="description">为角色分配可访问的API权限</p>
    </div>

    <div class="content">
      <!-- 角色选择 -->
      <div class="role-selector">
        <label>选择角色：</label>
        <el-select
          v-model="selectedRoleId"
          placeholder="请选择角色"
          @change="handleRoleChange"
          style="width: 300px"
        >
          <el-option
            v-for="role in roles"
            :key="role.id"
            :label="`${role.name} (${role.code})`"
            :value="role.id"
          />
        </el-select>
      </div>

      <!-- 部门选择（仅admin用户可见） -->
      <div v-if="isAdmin" class="dept-selector">
        <label>选择部门：</label>
        <el-select
          v-model="selectedDeptId"
          placeholder="请选择部门（不选择表示所有部门）"
          clearable
          style="width: 300px"
        >
          <el-option
            v-for="dept in departmentList"
            :key="dept.id"
            :label="dept.name"
            :value="dept.id"
          />
        </el-select>
        <span class="dept-tip">提示：不选择部门表示为所有部门分配权限</span>
      </div>

      <!-- API权限配置 -->
      <div v-if="selectedRoleId" class="permission-config">
        <div class="section-header">
          <h3>API权限配置</h3>
          <div class="actions">
            <el-button @click="expandAll">展开全部</el-button>
            <el-button @click="collapseAll">收起全部</el-button>
            <el-button type="primary" @click="savePermissions" :loading="saving">
              保存权限
            </el-button>
          </div>
        </div>

        <!-- 业务线分组 -->
        <div class="business-groups">
          <div
            v-for="business in businessList"
            :key="business.id"
            class="business-group"
          >
            <div class="business-header" @click="toggleBusiness(business.id)">
              <i
                :class="[
                  'el-icon-arrow-right',
                  { 'expanded': expandedBusiness.includes(business.id) }
                ]"
              ></i>
              <el-checkbox
                :value="isBusinessAllSelected(business.id)"
                :indeterminate="isBusinessIndeterminate(business.id)"
                @change="handleBusinessCheckChange(business.id, $event)"
              >
                {{ business.name }}
              </el-checkbox>
              <span class="api-count">({{ getBusinessAPICount(business.id) }} 个API)</span>
            </div>

            <div
              v-show="expandedBusiness.includes(business.id)"
              class="api-list"
            >
              <div
                v-for="api in getBusinessAPIs(business.id)"
                :key="api.id"
                class="api-item"
              >
                <el-checkbox
                  :value="selectedAPIIds.includes(api.id)"
                  @change="handleAPICheckChange(api.id, $event)"
                >
                  <div class="api-info">
                    <span class="api-name">{{ api.name }}</span>
                    <span class="api-method" :class="api.method.toLowerCase()">
                      {{ api.method }}
                    </span>
                    <span class="api-path">{{ api.path }}</span>
                  </div>
                  <div v-if="api.description" class="api-description">
                    {{ api.description }}
                  </div>
                </el-checkbox>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-else class="empty-state">
        <i class="el-icon-user"></i>
        <p>请先选择一个角色</p>
      </div>
    </div>
  </div>
</template>

<script>
import { getAllRoles, getUserAccessibleRoles, getRoleAPIPermissions, updateRoleAPIPermissions } from '@/api/role'
import { getApiList } from '@/api/api'
import { getBusinessList } from '@/api/business'
import { getDepartmentList } from '@/api/department'

export default {
  name: 'RolePermission',
  data() {
    return {
      roles: [],
      selectedRoleId: null,
      selectedDeptId: null,
      selectedAPIIds: [],
      businessList: [],
      apiList: [],
      departmentList: [],
      expandedBusiness: [],
      saving: false
    }
  },
  computed: {
    // 根据业务线分组的API
    apisByBusiness() {
      const groups = {}
      this.apiList.forEach(api => {
        if (!groups[api.business_id]) {
          groups[api.business_id] = []
        }
        groups[api.business_id].push(api)
      })
      return groups
    },
    // 获取当前用户的部门ID
    currentUserDeptId() {
      return this.$store.getters.deptId
    },
    // 判断当前用户是否为admin
    isAdmin() {
      const roles = this.$store.getters.roles || []
      return roles.includes('admin')
    }
  },
  async mounted() {
    await this.loadRoles()
    await this.loadBusinessList()
    await this.loadAPIList()
    if (this.isAdmin) {
      await this.loadDepartmentList()
    }
  },
  methods: {
    // 加载角色列表
    async loadRoles() {
      try {
        const response = await getUserAccessibleRoles()
        this.roles = response.data || []
      } catch (error) {
        console.error('加载角色失败:', error)
        this.$message.error('加载角色失败')
      }
    },

    // 加载业务线列表
    async loadBusinessList() {
      try {
        const response = await getBusinessList({ page: 1, limit: 1000 })
        this.businessList = response.data?.items || []
      } catch (error) {
        this.$message.error('加载业务线列表失败')
        console.error(error)
      }
    },

    // 加载API列表
    async loadAPIList() {
      try {
        const response = await getApiList({ page: 1, limit: 1000 })
        this.apiList = response.data?.items || []
      } catch (error) {
        this.$message.error('加载API列表失败')
        console.error(error)
      }
    },

    // 加载部门列表
    async loadDepartmentList() {
      try {
        const response = await getDepartmentList({ page: 1, limit: 1000 })
        this.departmentList = response.data?.items || []
      } catch (error) {
        this.$message.error('加载部门列表失败')
        console.error(error)
      }
    },

    // 角色选择变化
    async handleRoleChange(roleId) {
      if (roleId) {
        await this.loadRolePermissions(roleId)
      } else {
        this.selectedAPIIds = []
      }
    },

    // 加载角色权限
    async loadRolePermissions(roleId) {
      try {
        const response = await getRoleAPIPermissions(roleId)
        this.selectedAPIIds = response.data || []
      } catch (error) {
        this.$message.error('加载角色权限失败')
        console.error(error)
        this.selectedAPIIds = []
      }
    },

    // 获取业务线的API
    getBusinessAPIs(businessId) {
      return this.apisByBusiness[businessId] || []
    },

    // 获取业务线API数量
    getBusinessAPICount(businessId) {
      return this.getBusinessAPIs(businessId).length
    },

    // 检查业务线是否全选
    isBusinessAllSelected(businessId) {
      const apis = this.getBusinessAPIs(businessId)
      return apis.length > 0 && apis.every(api => this.selectedAPIIds.includes(api.id))
    },

    // 检查业务线是否半选
    isBusinessIndeterminate(businessId) {
      const apis = this.getBusinessAPIs(businessId)
      const selectedCount = apis.filter(api => this.selectedAPIIds.includes(api.id)).length
      return selectedCount > 0 && selectedCount < apis.length
    },

    // 业务线展开/收起
    toggleBusiness(businessId) {
      const index = this.expandedBusiness.indexOf(businessId)
      if (index > -1) {
        this.expandedBusiness.splice(index, 1)
      } else {
        this.expandedBusiness.push(businessId)
      }
    },

    // 业务线全选/取消全选
    handleBusinessCheckChange(businessId, checked) {
      const apis = this.getBusinessAPIs(businessId)
      if (checked) {
        // 全选
        apis.forEach(api => {
          if (!this.selectedAPIIds.includes(api.id)) {
            this.selectedAPIIds.push(api.id)
          }
        })
      } else {
        // 取消全选
        apis.forEach(api => {
          const index = this.selectedAPIIds.indexOf(api.id)
          if (index > -1) {
            this.selectedAPIIds.splice(index, 1)
          }
        })
      }
    },

    // API选择变化
    handleAPICheckChange(apiId, checked) {
      if (checked) {
        if (!this.selectedAPIIds.includes(apiId)) {
          this.selectedAPIIds.push(apiId)
        }
      } else {
        const index = this.selectedAPIIds.indexOf(apiId)
        if (index > -1) {
          this.selectedAPIIds.splice(index, 1)
        }
      }
    },

    // 展开全部
    expandAll() {
      this.expandedBusiness = this.businessList.map(b => b.id)
    },

    // 收起全部
    collapseAll() {
      this.expandedBusiness = []
    },

    // 保存权限
    async savePermissions() {
      if (!this.selectedRoleId) {
        this.$message.warning('请先选择角色')
        return
      }

      this.saving = true
      try {
        // 确定要使用的部门ID
        let deptId = null
        if (this.isAdmin) {
          // admin用户可以选择部门，如果没有选择则使用null（表示所有部门）
          deptId = this.selectedDeptId || null
        } else {
          // 非admin用户使用自己的部门ID
          deptId = this.currentUserDeptId
        }
        
        await updateRoleAPIPermissions(this.selectedRoleId, this.selectedAPIIds, deptId)
        this.$message.success('权限保存成功')
      } catch (error) {
        this.$message.error('权限保存失败')
        console.error(error)
      } finally {
        this.saving = false
      }
    }
  }
}
</script>

<style scoped>
.role-permission-container {
  padding: 20px;
  background: #fff;
  border-radius: 8px;
}

.header {
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #ebeef5;
}

.header h2 {
  margin: 0 0 8px 0;
  color: #303133;
  font-size: 24px;
  font-weight: 600;
}

.description {
  margin: 0;
  color: #909399;
  font-size: 14px;
}

.role-selector {
  margin-bottom: 30px;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 6px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.role-selector label {
  font-weight: 500;
  color: #606266;
  white-space: nowrap;
}

.dept-selector {
  margin-bottom: 30px;
  padding: 20px;
  background: #f0f9ff;
  border-radius: 6px;
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.dept-selector label {
  font-weight: 500;
  color: #606266;
  white-space: nowrap;
}

.dept-tip {
  color: #909399;
  font-size: 12px;
  margin-left: 12px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid #ebeef5;
}

.section-header h3 {
  margin: 0;
  color: #303133;
  font-size: 18px;
  font-weight: 600;
}

.actions {
  display: flex;
  gap: 8px;
}

.business-groups {
  border: 1px solid #ebeef5;
  border-radius: 6px;
  overflow: hidden;
}

.business-group {
  border-bottom: 1px solid #ebeef5;
}

.business-group:last-child {
  border-bottom: none;
}

.business-header {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  background: #fafafa;
  cursor: pointer;
  transition: background-color 0.2s;
  gap: 8px;
}

.business-header:hover {
  background: #f0f0f0;
}

.business-header .el-icon-arrow-right {
  transition: transform 0.2s;
  color: #909399;
}

.business-header .el-icon-arrow-right.expanded {
  transform: rotate(90deg);
}

.api-count {
  color: #909399;
  font-size: 12px;
  margin-left: auto;
}

.api-list {
  background: #fff;
}

.api-item {
  padding: 12px 20px 12px 48px;
  border-bottom: 1px solid #f5f7fa;
}

.api-item:last-child {
  border-bottom: none;
}

.api-info {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 4px;
}

.api-name {
  font-weight: 500;
  color: #303133;
}

.api-method {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  text-transform: uppercase;
}

.api-method.get {
  background: #e7f7ff;
  color: #1890ff;
}

.api-method.post {
  background: #f6ffed;
  color: #52c41a;
}

.api-method.put {
  background: #fff7e6;
  color: #fa8c16;
}

.api-method.delete {
  background: #fff2f0;
  color: #ff4d4f;
}

.api-path {
  color: #606266;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
}

.api-description {
  color: #909399;
  font-size: 12px;
  margin-left: 24px;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #909399;
}

.empty-state i {
  font-size: 48px;
  margin-bottom: 16px;
  display: block;
}

.empty-state p {
  margin: 0;
  font-size: 16px;
}
</style>