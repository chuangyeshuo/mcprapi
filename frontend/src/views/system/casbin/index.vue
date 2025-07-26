<template>
  <div class="app-container">
    <!-- 页面标题和说明 -->
    <div class="page-header">
      <h2>Casbin权限管理</h2>
      <p class="description">管理系统的访问控制策略，支持基于角色的权限控制(RBAC)和部门级权限隔离</p>
    </div>

    <!-- 策略统计卡片 -->
    <el-row :gutter="20" class="stats-container">
      <el-col :span="6">
        <el-card shadow="hover" class="stats-card">
          <div class="stats-content">
            <div class="stats-icon policy-icon">
              <i class="el-icon-document"></i>
            </div>
            <div class="stats-info">
              <div class="stats-number">{{ policyStats.total }}</div>
              <div class="stats-label">总策略数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stats-card">
          <div class="stats-content">
            <div class="stats-icon role-icon">
              <i class="el-icon-user"></i>
            </div>
            <div class="stats-info">
              <div class="stats-number">{{ policyStats.roleCount }}</div>
              <div class="stats-label">策略(p)</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stats-card">
          <div class="stats-content">
            <div class="stats-icon inherit-icon">
              <i class="el-icon-share"></i>
            </div>
            <div class="stats-info">
              <div class="stats-number">{{ policyStats.inheritCount }}</div>
              <div class="stats-label">角色继承(g)</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stats-card">
          <div class="stats-content">
            <div class="stats-icon dept-icon">
              <i class="el-icon-office-building"></i>
            </div>
            <div class="stats-info">
              <div class="stats-number">{{ policyStats.deptCount }}</div>
              <div class="stats-label">部门策略</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <div class="filter-container">
      <el-input
        v-model="listQuery.subject"
        placeholder="搜索部门/角色"
        style="width: 200px;"
        class="filter-item"
        prefix-icon="el-icon-search"
        @keyup.enter.native="handleFilter"
        clearable
      />
      <el-select
        v-model="listQuery.ptype"
        placeholder="策略类型"
        clearable
        style="width: 120px"
        class="filter-item"
      >
        <el-option value="p" label="策略(p)" />
        <el-option value="g" label="角色继承(g)" />
      </el-select>
      <el-select
        v-model="listQuery.dept_id"
        placeholder="部门"
        clearable
        style="width: 150px"
        class="filter-item"
      >
        <el-option
          v-for="dept in departments"
          :key="dept.id"
          :label="dept.name"
          :value="dept.id"
        />
      </el-select>
      <el-button
        v-waves
        class="filter-item"
        type="primary"
        icon="el-icon-search"
        @click="handleFilter"
      >
        搜索
      </el-button>
      <el-button
        class="filter-item"
        style="margin-left: 10px;"
        type="primary"
        icon="el-icon-plus"
        @click="handleCreate"
      >
        添加策略
      </el-button>
      <el-button
        class="filter-item"
        style="margin-left: 10px;"
        type="danger"
        icon="el-icon-delete"
        :disabled="multipleSelection.length === 0"
        @click="handleBatchDelete"
      >
        批量删除 ({{ multipleSelection.length }})
      </el-button>
      <el-button
        class="filter-item"
        style="margin-left: 10px;"
        type="success"
        icon="el-icon-refresh"
        @click="handleReload"
        :loading="reloadLoading"
      >
        重新加载
      </el-button>
      <el-button
        class="filter-item"
        style="margin-left: 10px;"
        type="warning"
        icon="el-icon-download"
        @click="handleExport"
      >
        导出策略
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
      @sort-change="sortChange"
      @selection-change="handleSelectionChange"
    >
      <el-table-column
        type="selection"
        width="55"
        align="center"
      />
      <el-table-column label="ID" prop="id" align="center" width="80" sortable="custom">
        <template slot-scope="{row}">
          <span class="id-badge">{{ row.id }}</span>
        </template>
      </el-table-column>
      <el-table-column label="策略类型" width="100px" align="center">
        <template slot-scope="{row}">
          <el-tag :type="row.ptype === 'p' ? 'primary' : 'success'" size="small">
            {{ row.ptype === 'p' ? '策略(p)' : '角色继承(g)' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="主体(部门/角色)" min-width="150px">
        <template slot-scope="{row}">
          <div class="subject-cell">
            <i :class="getSubjectIcon(row.v0)" class="subject-icon"></i>
            <span class="subject-name">{{ row.v0 }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="资源(路径)" min-width="200px">
        <template slot-scope="{row}">
          <div class="resource-cell">
            <el-tooltip :content="row.v1" placement="top">
              <code class="resource-path">{{ row.v1 }}</code>
            </el-tooltip>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="动作(方法)" width="100px" align="center">
        <template slot-scope="{row}">
          <el-tag :type="getMethodTagType(row.v2)" size="small">
            {{ row.v2 }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="部门" width="120px" align="center">
        <template slot-scope="{row}">
          <span v-if="row.v3 === '*'" class="dept-all">
            <i class="el-icon-office-building"></i> 全部
          </span>
          <span v-else class="dept-specific">
            {{ getDepartmentName(row.v3) || row.v3 }}
          </span>
        </template>
      </el-table-column>
      <el-table-column label="效果" width="80px" align="center">
        <template slot-scope="{row}">
          <el-tag :type="row.v4 === 'allow' ? 'success' : 'danger'" size="small">
            {{ row.v4 || 'allow' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="200px" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <el-button type="primary" size="mini" icon="el-icon-edit" @click="handleUpdate(row)">
            编辑
          </el-button>
          <el-button type="danger" size="mini" icon="el-icon-delete" @click="handleDelete(row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination
      v-show="total>0"
      :total="total"
      :page.sync="listQuery.page"
      :limit.sync="listQuery.page_size"
      @pagination="getList"
    />

    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible" width="600px">
      <el-form
        ref="dataForm"
        :rules="rules"
        :model="temp"
        label-position="left"
        label-width="120px"
        style="width: 500px; margin-left:20px;"
      >
        <el-form-item label="策略类型" prop="ptype">
          <el-select v-model="temp.ptype" class="filter-item" placeholder="请选择策略类型">
            <el-option value="p" label="策略(p)" />
            <el-option value="g" label="角色继承(g)" />
          </el-select>
        </el-form-item>
        <el-form-item label="主体(部门/角色)" prop="v0">
          <el-input v-model="temp.v0" placeholder="如: admin, user, role_name" />
        </el-form-item>
        <el-form-item label="资源(路径)" prop="v1">
          <el-input v-model="temp.v1" placeholder="如: /api/v1/*, /api/v1/users" />
        </el-form-item>
        <el-form-item label="动作(方法)" prop="v2">
          <el-select v-model="temp.v2" class="filter-item" placeholder="请选择HTTP方法" clearable>
            <el-option label="所有方法 (*)" value="*" />
            <el-option label="GET" value="GET" />
            <el-option label="POST" value="POST" />
            <el-option label="PUT" value="PUT" />
            <el-option label="DELETE" value="DELETE" />
            <el-option label="PATCH" value="PATCH" />
          </el-select>
        </el-form-item>
        <el-form-item label="部门" prop="v3">
          <el-select v-model="temp.v3" class="filter-item" placeholder="请选择部门" clearable>
            <el-option label="所有部门 (*)" value="*" />
            <el-option
              v-for="dept in departments"
              :key="dept.id"
              :label="dept.name"
              :value="dept.id.toString()"
            />
          </el-select>
          <div class="form-tip">
            <i class="el-icon-info"></i>
            <span>默认使用当前用户所属部门，可手动修改为其他部门或全部部门</span>
          </div>
        </el-form-item>
        <el-form-item label="效果" prop="v4">
          <el-select v-model="temp.v4" class="filter-item" placeholder="请选择效果">
            <el-option label="允许 (allow)" value="allow" />
            <el-option label="拒绝 (deny)" value="deny" />
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">
          取消
        </el-button>
        <el-button v-if="dialogStatus === 'create'" type="primary" @click="createData()">
          确认
        </el-button>
        <el-button v-else type="primary" @click="updateData()">
          确认
        </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import {
  addCasbinPolicy,
  batchDeleteCasbinPolicy,
  deleteCasbinPolicy,
  getCasbinPolicyList,
  reloadCasbinPolicy,
  updateCasbinPolicy,
} from '@/api/casbin';
import { getDepartmentList } from '@/api/department';
import Pagination
  from '@/components/Pagination'; // secondary package based on el-pagination
import waves from '@/directive/waves'; // waves directive
import { mapGetters } from 'vuex';

export default {
  name: 'CasbinPolicy',
  components: { Pagination },
  directives: { waves },
  computed: {
    ...mapGetters(['deptId'])
  },
  data() {
    return {
      tableKey: 0,
      list: null,
      total: 0,
      listLoading: true,
      reloadLoading: false,
      multipleSelection: [], // 多选数据
      listQuery: {
        page: 1,
        page_size: 20,
        subject: undefined,
        ptype: undefined,
        dept_id: undefined
      },
      temp: {
        ptype: 'p',
        v0: '',
        v1: '',
        v2: '*',
        v3: '*',
        v4: 'allow'
      },
      dialogFormVisible: false,
      dialogStatus: '',
      textMap: {
        create: '添加策略',
        update: '编辑策略'
      },
      rules: {
        ptype: [{ required: true, message: '策略类型是必填项', trigger: 'change' }],
        v0: [{ required: true, message: '主体是必填项', trigger: 'blur' }],
        v1: [{ required: true, message: '资源是必填项', trigger: 'blur' }],
        v2: [{ required: true, message: '动作是必填项', trigger: 'change' }]
      },
      departments: [],
      policyStats: {
        total: 0,
        roleCount: 0,
        inheritCount: 0,
        deptCount: 0
      }
    }
  },
  created() {
    this.getList()
    this.loadDepartments()
  },
  methods: {
    getList() {
      this.listLoading = true
      getCasbinPolicyList(this.listQuery).then(response => {
        this.list = response.data.list
        this.total = response.data.total
        this.updateStatistics()
        this.listLoading = false
      }).catch(() => {
        this.listLoading = false
      })
    },
    handleFilter() {
      this.listQuery.page = 1
      this.getList()
    },
    sortChange(data) {
      const { prop, order } = data
      if (prop === 'id') {
        this.sortByID(order)
      }
    },
    sortByID(order) {
      if (order === 'ascending') {
        this.listQuery.sort = '+id'
      } else {
        this.listQuery.sort = '-id'
      }
      this.handleFilter()
    },
    resetTemp() {
      this.temp = {
        ptype: 'p',
        v0: '',
        v1: '',
        v2: '*',
        v3: this.deptId ? this.deptId.toString() : '*', // 自动设置当前用户的部门ID
        v4: 'allow'
      }
    },
    handleCreate() {
      this.resetTemp()
      this.dialogStatus = 'create'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    createData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          // 确保部门信息存在
          if (!this.temp.v3) {
            this.temp.v3 = this.deptId ? this.deptId.toString() : '*'
          }
          
          addCasbinPolicy(this.temp).then(() => {
            this.dialogFormVisible = false
            this.$notify({
              title: '成功',
              message: `策略添加成功${this.temp.v3 !== '*' ? '，已关联到当前部门' : ''}`,
              type: 'success',
              duration: 2000
            })
            this.getList()
          })
        }
      })
    },
    handleUpdate(row) {
      this.temp = Object.assign({}, row) // 复制行数据到temp
      this.dialogStatus = 'update'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    updateData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          const tempData = Object.assign({}, this.temp)
          // 确保部门信息存在
          if (!tempData.v3) {
            tempData.v3 = this.deptId ? this.deptId.toString() : '*'
          }
          
          updateCasbinPolicy(tempData).then(() => {
            this.dialogFormVisible = false
            this.$notify({
              title: '成功',
              message: `策略更新成功${tempData.v3 !== '*' ? '，已关联到当前部门' : ''}`,
              type: 'success',
              duration: 2000
            })
            this.getList()
          })
        }
      })
    },
    handleDelete(row) {
      this.$confirm('确认删除该策略?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteCasbinPolicy(row.id).then(() => {
          this.$notify({
            title: '成功',
            message: '策略删除成功',
            type: 'success',
            duration: 2000
          })
          this.getList()
        })
      })
    },
    handleSelectionChange(val) {
      this.multipleSelection = val
    },
    handleBatchDelete() {
      if (this.multipleSelection.length === 0) {
        this.$message.warning('请选择要删除的策略')
        return
      }
      
      this.$confirm(`确定要删除选中的 ${this.multipleSelection.length} 条策略吗?`, '批量删除确认', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        const ids = this.multipleSelection.map(item => item.id)
        batchDeleteCasbinPolicy(ids).then(() => {
          this.$notify({
            title: '成功',
            message: '批量删除成功',
            type: 'success',
            duration: 2000
          })
          this.getList()
          this.multipleSelection = []
        }).catch(() => {
          this.$notify({
            title: '错误',
            message: '批量删除失败',
            type: 'error',
            duration: 2000
          })
        })
      })
    },
    handleReload() {
      this.$confirm('确认重新加载所有策略?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.reloadLoading = true
        reloadCasbinPolicy().then(() => {
          this.$notify({
            title: '成功',
            message: '策略重新加载成功',
            type: 'success',
            duration: 2000
          })
          this.reloadLoading = false
          this.getList()
        }).catch(() => {
          this.reloadLoading = false
        })
      })
    },
    handleExport() {
      // 简单的CSV导出功能
      const headers = ['ID', '策略类型', '主体', '资源', '动作', '部门', '效果']
      const csvContent = [
        headers.join(','),
        ...this.list.map(row => [
          row.id,
          row.ptype === 'p' ? '策略' : '角色继承',
          row.v0,
          row.v1,
          row.v2,
          row.v3 === '*' ? '全部部门' : this.getDepartmentName(row.v3),
          row.v4 || 'allow'
        ].join(','))
      ].join('\n')
      
      const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
      const link = document.createElement('a')
      const url = URL.createObjectURL(blob)
      link.setAttribute('href', url)
      link.setAttribute('download', 'casbin-policies.csv')
      link.style.visibility = 'hidden'
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      
      this.$notify({
        title: '成功',
        message: '策略导出成功',
        type: 'success',
        duration: 2000
      })
    },
    formatJson(filterVal) {
      return this.list.map(v => filterVal.map(j => {
        if (j === 'ptype') {
          return v[j] === 'p' ? '策略' : '角色继承'
        } else if (j === 'v3' && v[j] === '*') {
          return '全部部门'
        } else if (j === 'v4') {
          return v[j] || 'allow'
        } else {
          return v[j]
        }
      }))
    },
    getSubjectIcon(subject) {
      if (subject === 'admin' || subject.includes('admin')) {
        return 'el-icon-user-solid'
      } else if (subject.startsWith('role_') || subject === 'user') {
        return 'el-icon-user'
      } else {
        return 'el-icon-s-custom'
      }
    },
    getMethodTagType(method) {
      const typeMap = {
        'GET': 'success',
        'POST': 'primary',
        'PUT': 'warning',
        'DELETE': 'danger',
        'PATCH': 'info',
        '*': ''
      }
      return typeMap[method] || 'info'
    },
    getDepartmentName(deptId) {
      if (deptId === '*') return '全部'
      const dept = this.departments.find(d => d.id.toString() === deptId)
      return dept ? dept.name : deptId
    },
    updateStatistics() {
      if (!this.list) return
      this.policyStats.total = this.list.length
      this.policyStats.roleCount = this.list.filter(item => item.ptype === 'p').length
      this.policyStats.inheritCount = this.list.filter(item => item.ptype === 'g').length
      this.policyStats.deptCount = new Set(this.list.map(item => item.v3).filter(dept => dept && dept !== '*')).size
    },
    loadDepartments() {
      getDepartmentList().then(response => {
        // 后端返回格式：{ code: 200, data: { items: [...], total: ... } }
        this.departments = response.data.items || response.data.list || response.data || []
      }).catch(error => {
        console.error('获取部门列表失败:', error)
        // 如果API调用失败，使用默认数据
        this.departments = [
          { id: 1, name: '技术部' },
          { id: 2, name: '市场部' },
          { id: 3, name: '人事部' },
          { id: 4, name: '财务部' }
        ]
      })
    }
  }
}
</script>

<style scoped>
.app-container {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin: 0;
}

.page-description {
  color: #909399;
  margin-top: 8px;
  font-size: 14px;
}

.stats-cards {
  display: flex;
  gap: 20px;
  margin-bottom: 20px;
}

.stat-card {
  flex: 1;
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  border-left: 4px solid #409eff;
}

.stat-card.success {
  border-left-color: #67c23a;
}

.stat-card.warning {
  border-left-color: #e6a23c;
}

.stat-card.info {
  border-left-color: #909399;
}

.stat-number {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 8px;
}

.stat-label {
  color: #909399;
  font-size: 14px;
}

.filter-container {
  padding: 20px;
  background: #fff;
  border-radius: 8px;
  margin-bottom: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.filter-item {
  margin-right: 10px;
  margin-bottom: 10px;
}

.id-badge {
  display: inline-block;
  padding: 2px 8px;
  background: #f0f9ff;
  color: #1890ff;
  border-radius: 4px;
  font-weight: 500;
}

.subject-cell {
  display: flex;
  align-items: center;
}

.subject-icon {
  margin-right: 8px;
  color: #409eff;
}

.subject-name {
  font-weight: 500;
}

.resource-cell {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

.resource-path {
  background: #f5f5f5;
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 12px;
  color: #e74c3c;
  border: 1px solid #e0e0e0;
}

.dept-all {
  color: #67c23a;
  font-weight: 500;
}

.dept-specific {
  color: #409eff;
  font-weight: 500;
}

.table-container {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.pagination-container {
  padding: 20px;
  text-align: center;
  background: #fff;
  border-radius: 8px;
  margin-top: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.dialog-footer {
  text-align: right;
  padding-top: 20px;
}

.form-tip {
  margin-top: 5px;
  font-size: 12px;
  color: #909399;
  display: flex;
  align-items: center;
}

.form-tip i {
  margin-right: 4px;
  color: #409eff;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .stats-cards {
    flex-direction: column;
  }
  
  .filter-container {
    padding: 15px;
  }
  
  .filter-item {
    width: 100%;
    margin-bottom: 15px;
  }
}
</style>