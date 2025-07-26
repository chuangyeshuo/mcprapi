<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input
        v-model="listQuery.path"
        placeholder="API路径"
        style="width: 200px;"
        class="filter-item"
        @keyup.enter.native="handleFilter"
      />
      <el-select v-model="listQuery.method" placeholder="请求方法" clearable style="width: 120px" class="filter-item">
        <el-option v-for="item in methodOptions" :key="item" :label="item" :value="item" />
      </el-select>
      <el-select 
        v-if="isAdmin"
        v-model="listQuery.dept_id" 
        placeholder="请选择部门" 
        clearable 
        class="filter-item" 
        @change="onDepartmentChange"
      >
        <el-option
          v-for="item in departmentOptions"
          :key="item.id"
          :label="item.name"
          :value="item.id"
        />
      </el-select>
      <el-select v-model="listQuery.business_id" placeholder="请选择业务线" clearable class="filter-item">
        <el-option
          v-for="item in filteredBusinessOptions"
          :key="item.id"
          :label="item.name"
          :value="item.id"
        />
      </el-select>
      <el-button v-waves class="filter-item" type="primary" icon="el-icon-search" @click="handleFilter">
        搜索
      </el-button>
      <el-button class="filter-item" style="margin-left: 10px;" type="primary" icon="el-icon-edit" @click="handleCreate">
        添加
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
      <el-table-column label="API名称" width="150px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.name }}</span>
        </template>
      </el-table-column>
      <el-table-column label="API路径" width="200px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.path }}</span>
        </template>
      </el-table-column>
      <el-table-column label="请求方法" width="100px" align="center">
        <template slot-scope="{row}">
          <el-tag :type="row.method | methodFilter">{{ row.method }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="业务线" width="120px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.business_name || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="部门" width="120px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.department_name || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="描述" width="200px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.description }}</span>
        </template>
      </el-table-column>
      <el-table-column label="状态" class-name="status-col" width="100">
        <template slot-scope="{row}">
          <el-tag :type="row.status | statusFilter">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="150px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.created_at | formatDateTime }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="230" class-name="small-padding fixed-width">
        <template slot-scope="{row,$index}">
          <el-button 
            v-if="isAdmin || row.dept_id === deptId"
            type="primary" 
            size="mini" 
            @click="handleUpdate(row)"
          >
            编辑
          </el-button>
          <el-button 
            v-if="(isAdmin || row.dept_id === deptId) && row.status!='deleted'" 
            size="mini" 
            type="danger" 
            @click="handleDelete(row,$index)"
          >
            删除
          </el-button>
          <span v-if="!isAdmin && row.dept_id !== deptId" class="text-muted">
            无权限操作
          </span>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />

    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible">
      <el-form ref="dataForm" :rules="rules" :model="temp" label-position="left" label-width="80px" style="width: 450px; margin-left:50px;">
        <el-form-item label="API名称" prop="name">
          <el-input v-model="temp.name" placeholder="请输入API名称" />
        </el-form-item>
        <el-form-item label="API路径" prop="path">
          <el-input v-model="temp.path" placeholder="请输入API路径，如：/api/v1/users" />
        </el-form-item>
        <el-form-item label="请求方法" prop="method">
          <el-select v-model="temp.method" class="filter-item" placeholder="请选择请求方法">
            <el-option v-for="item in methodOptions" :key="item" :label="item" :value="item" />
          </el-select>
        </el-form-item>
        <el-form-item label="部门" prop="dept_id">
          <el-select 
            v-model="temp.dept_id" 
            class="filter-item" 
            placeholder="请选择部门" 
            style="width: 100%;" 
            @change="onTempDepartmentChange"
            :disabled="!isAdmin"
          >
            <el-option v-for="item in departmentOptions" :key="item.id" :label="item.name" :value="item.id" />
          </el-select>
          <div v-if="!isAdmin" class="form-tip">
            <small class="text-muted">非管理员用户只能在自己的部门下创建API</small>
          </div>
        </el-form-item>
        <el-form-item label="业务线" prop="business_id">
          <el-select v-model="temp.business_id" class="filter-item" placeholder="请选择业务线" style="width: 100%;">
            <el-option v-for="item in tempFilteredBusinessOptions" :key="item.id" :label="item.name" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="temp.description" type="textarea" placeholder="请输入API描述" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="temp.status" class="filter-item" placeholder="请选择状态">
            <el-option v-for="item in statusOptions" :key="item.key" :label="item.display_name" :value="item.key" />
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">
          取消
        </el-button>
        <el-button type="primary" @click="dialogStatus==='create'?createData():updateData()">
          确认
        </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { getApiList, createApi, updateApi, deleteApi } from '@/api/api'
import { getAllBusiness } from '@/api/business'
import { getDepartmentList } from '@/api/department'
import waves from '@/directive/waves' // waves directive
import { parseTime } from '@/utils'
import Pagination from '@/components/Pagination' // secondary package based on el-pagination

export default {
  name: 'APITable',
  components: { Pagination },
  directives: { waves },
  filters: {
    statusFilter(status) {
      const statusMap = {
        1: 'success',
        0: 'info'
      }
      return statusMap[status]
    },
    methodFilter(method) {
      const methodMap = {
        'GET': 'success',
        'POST': 'warning',
        'PUT': 'info',
        'DELETE': 'danger'
      }
      return methodMap[method]
    }
  },
  computed: {
    ...mapGetters(['roles', 'deptId']),
    // 判断是否为管理员
    isAdmin() {
      return this.roles.includes('admin')
    }
  },
  data() {
    return {
      tableKey: 0,
      list: null,
      total: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 20,
        path: undefined,
        method: undefined,
        dept_id: undefined,
        business_id: undefined
      },
      methodOptions: ['GET', 'POST', 'PUT', 'DELETE'],
      statusOptions: [
        { key: 1, display_name: '启用' },
        { key: 0, display_name: '禁用' }
      ],
      businessOptions: [],
      departmentOptions: [],
      filteredBusinessOptions: [],
      tempFilteredBusinessOptions: [],
      temp: {
        id: undefined,
        name: '',
        path: '',
        method: 'GET',
        description: '',
        dept_id: undefined,
        business_id: undefined,
        status: 1
      },
      dialogFormVisible: false,
      dialogStatus: '',
      textMap: {
        update: '编辑',
        create: '创建'
      },
      rules: {
        name: [{ required: true, message: 'API名称是必填项', trigger: 'blur' }],
        path: [{ required: true, message: 'API路径是必填项', trigger: 'blur' }],
        method: [{ required: true, message: '请求方法是必填项', trigger: 'change' }],
        dept_id: [{ required: true, message: '部门是必填项', trigger: 'change' }],
        business_id: [{ required: true, message: '业务线是必填项', trigger: 'change' }]
      }
    }
  },
  created() {
    this.getList()
    this.getBusinessOptions()
    this.getDepartmentOptions()
  },
  methods: {
    async getList() {
      this.listLoading = true
      try {
        let queryParams = { ...this.listQuery }
        
        // 如果不是管理员，自动添加当前用户的部门ID过滤
        if (!this.isAdmin && this.deptId) {
          queryParams.dept_id = this.deptId
        }
        
        const response = await getApiList(queryParams)
        this.list = response.data.items
        this.total = response.data.total
      } catch (error) {
        this.$message.error('获取API列表失败')
        console.error(error)
      } finally {
        this.listLoading = false
      }
    },
    async getBusinessOptions() {
      try {
        const response = await getAllBusiness()
        let businessList = response.data || []
        
        // 如果不是管理员，只显示当前用户部门的业务线
        if (!this.isAdmin && this.deptId) {
          businessList = businessList.filter(business => business.dept_id === this.deptId)
        }
        
        this.businessOptions = businessList
        this.updateFilteredBusinessOptions()
      } catch (error) {
        this.businessOptions = []
        this.filteredBusinessOptions = []
        console.error('获取业务线选项失败:', error)
      }
    },
    getDepartmentOptions() {
      getDepartmentList().then(response => {
        this.departmentOptions = response.data.items || []
      }).catch(() => {
        this.departmentOptions = []
      })
    },
    updateFilteredBusinessOptions() {
      if (this.listQuery.dept_id) {
        // 根据选择的部门过滤业务线
        this.filteredBusinessOptions = this.businessOptions.filter(business => 
          business.dept_id === this.listQuery.dept_id
        )
      } else {
        // 显示所有业务线
        this.filteredBusinessOptions = [...this.businessOptions]
      }
    },
    onDepartmentChange() {
      // 部门变化时，清空业务线选择并更新过滤后的业务线选项
      this.listQuery.business_id = undefined
      this.updateFilteredBusinessOptions()
    },
    onTempDepartmentChange() {
      // 表单中部门变化时，清空业务线选择并更新过滤后的业务线选项
      this.temp.business_id = undefined
      this.updateTempFilteredBusinessOptions()
    },
    updateTempFilteredBusinessOptions() {
      if (this.temp.dept_id) {
        // 根据选择的部门过滤业务线
        this.tempFilteredBusinessOptions = this.businessOptions.filter(business => 
          business.dept_id === this.temp.dept_id
        )
      } else {
        // 显示所有业务线
        this.tempFilteredBusinessOptions = [...this.businessOptions]
      }
    },
    getBusinessName(businessId) {
      if (!businessId || !this.businessOptions.length) return '-'
      const business = this.businessOptions.find(item => item.id === businessId)
      return business ? business.name : '-'
    },
    getDepartmentName(deptId) {
      if (!deptId || !this.departmentOptions.length) return '-'
      const department = this.departmentOptions.find(item => item.id === deptId)
      return department ? department.name : '-'
    },
    handleFilter() {
      this.listQuery.page = 1
      this.getList()
    },
    resetTemp() {
      this.temp = {
        id: undefined,
        name: '',
        path: '',
        method: 'GET',
        description: '',
        dept_id: undefined,
        business_id: undefined,
        status: 1
      }
      this.tempFilteredBusinessOptions = [...this.businessOptions]
    },
    handleCreate() {
      this.resetTemp()
      // 如果不是管理员，自动设置为当前用户的部门
      if (!this.isAdmin && this.deptId) {
        this.temp.dept_id = this.deptId
        this.updateTempFilteredBusinessOptions()
      }
      this.dialogStatus = 'create'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    createData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          createApi(this.temp).then(() => {
            this.list.unshift(this.temp)
            this.dialogFormVisible = false
            this.$notify({
              title: '成功',
              message: '创建成功',
              type: 'success',
              duration: 2000
            })
            this.getList() // 重新获取列表以显示最新数据
          }).catch(error => {
            this.$notify({
              title: '错误',
              message: error.response?.data?.message || '创建失败',
              type: 'error',
              duration: 3000
            })
          })
        }
      })
    },
    handleUpdate(row) {
      // 检查权限：非管理员用户只能编辑自己部门的API
      if (!this.isAdmin && row.dept_id !== this.deptId) {
        this.$message.warning('您只能编辑自己部门的API')
        return
      }
      
      this.temp = Object.assign({}, row) // copy obj
      this.updateTempFilteredBusinessOptions() // 初始化过滤后的业务线选项
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
          // 后端UpdateAPIRequest需要ID字段，所以不能从数据中移除ID
          updateApi(tempData.id, tempData).then(() => {
            const index = this.list.findIndex(v => v.id === this.temp.id)
            this.list.splice(index, 1, this.temp)
            this.dialogFormVisible = false
            this.$notify({
              title: '成功',
              message: '更新成功',
              type: 'success',
              duration: 2000
            })
          }).catch(error => {
            this.$notify({
              title: '错误',
              message: error.response?.data?.message || '更新失败',
              type: 'error',
              duration: 3000
            })
          })
        }
      })
    },
    handleDelete(row, index) {
      // 检查权限：非管理员用户只能删除自己部门的API
      if (!this.isAdmin && row.dept_id !== this.deptId) {
        this.$message.warning('您只能删除自己部门的API')
        return
      }
      
      this.$confirm('此操作将永久删除该API, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteApi(row.id).then(() => {
          this.$notify({
            title: '成功',
            message: '删除成功',
            type: 'success',
            duration: 2000
          })
          this.list.splice(index, 1)
        })
      })
    }
  }
}
</script>

<style scoped>
.text-muted {
  color: #999;
  font-size: 12px;
}

.form-tip {
  margin-top: 5px;
}

.form-tip small {
  color: #909399;
}
</style>