<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input
        v-model="listQuery.name"
        placeholder="业务名称"
        style="width: 200px;"
        class="filter-item"
        @keyup.enter.native="handleFilter"
      />
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
      <el-table-column label="业务名称" width="150px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.name }}</span>
        </template>
      </el-table-column>
      <el-table-column label="业务编码" width="120px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.code }}</span>
        </template>
      </el-table-column>
      <el-table-column label="所属部门" width="120px" align="center">
        <template slot-scope="{row}">
          <span>{{ getDepartmentName(row.dept_id) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="负责人" width="100px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.owner || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="联系邮箱" width="150px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.email || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="描述" width="150px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.description || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="状态" class-name="status-col" width="80">
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
      <el-table-column label="操作" align="center" width="180" class-name="small-padding fixed-width">
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
        <el-form-item label="业务名称" prop="name">
          <el-input v-model="temp.name" placeholder="请输入业务名称" />
        </el-form-item>
        <el-form-item label="业务编码" prop="code">
          <el-input v-model="temp.code" placeholder="请输入业务编码，如：user-center" />
        </el-form-item>
        <el-form-item label="所属部门" prop="dept_id">
          <el-select 
            v-model="temp.dept_id" 
            class="filter-item" 
            placeholder="请选择所属部门" 
            style="width: 100%;"
            :disabled="!isAdmin"
          >
            <el-option v-for="item in departmentOptions" :key="item.id" :label="item.name" :value="item.id" />
          </el-select>
          <div v-if="!isAdmin" class="form-tip">
            <small class="text-muted">非管理员用户只能在自己的部门下创建业务线</small>
          </div>
        </el-form-item>
        <el-form-item label="负责人" prop="owner">
          <el-input v-model="temp.owner" placeholder="请输入负责人姓名" />
        </el-form-item>
        <el-form-item label="联系邮箱" prop="email">
          <el-input v-model="temp.email" placeholder="请输入联系邮箱" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="temp.description" type="textarea" placeholder="请输入业务描述" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="temp.status" class="filter-item" placeholder="请选择状态" style="width: 100%;">
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
import {
  createBusiness,
  deleteBusiness,
  getBusinessList,
  updateBusiness,
} from '@/api/business';
import { getDepartmentList } from '@/api/department';
import Pagination
  from '@/components/Pagination'; // secondary package based on el-pagination
import waves from '@/directive/waves'; // waves directive
import { parseTime } from '@/utils';

export default {
  name: 'BusinessTable',
  components: { Pagination },
  directives: { waves },
  filters: {
    statusFilter(status) {
      const statusMap = {
        1: 'success',
        0: 'info'
      }
      return statusMap[status]
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
        name: undefined
      },
      statusOptions: [
        { key: 1, display_name: '启用' },
        { key: 0, display_name: '禁用' }
      ],
      departmentOptions: [],
      temp: {
        id: undefined,
        name: '',
        code: '',
        dept_id: undefined,
        description: '',
        owner: '',
        email: '',
        status: 1
      },
      dialogFormVisible: false,
      dialogStatus: '',
      textMap: {
        update: '编辑',
        create: '创建'
      },
      rules: {
        name: [{ required: true, message: '业务名称是必填项', trigger: 'blur' }],
        code: [
          { required: true, message: '业务编码是必填项', trigger: 'blur' },
          { pattern: /^[a-zA-Z0-9_-]+$/, message: '业务编码只能包含字母、数字、下划线和横线', trigger: 'blur' }
        ],
        dept_id: [{ required: true, message: '所属部门是必填项', trigger: 'change' }],
        email: [
          { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
        ]
      }
    }
  },
  created() {
    this.getList()
    this.getDepartmentOptions()
  },
  methods: {
    async getList() {
      this.listLoading = true
      try {
        let response
        
        // 如果是管理员，获取所有业务线；否则只获取当前用户部门的业务线
        if (this.isAdmin) {
          response = await getBusinessList(this.listQuery)
        } else {
          // 非管理员用户，添加部门ID过滤
          const params = {
            ...this.listQuery,
            dept_id: this.deptId
          }
          response = await getBusinessList(params)
        }
        
        this.list = response.data.items
        this.total = response.data.total
      } catch (error) {
        this.$message.error('获取业务线列表失败')
        console.error(error)
      } finally {
        this.listLoading = false
      }
    },
    getDepartmentOptions() {
      getDepartmentList({ page: 1, limit: 1000 }).then(response => {
        const allDepartments = response.data.items || []
        
        // 如果是管理员，显示所有部门；否则只显示当前用户的部门
        if (this.isAdmin) {
          this.departmentOptions = allDepartments
        } else {
          // 非管理员用户只能看到自己的部门
          this.departmentOptions = allDepartments.filter(dept => dept.id === this.deptId)
        }
      }).catch(() => {
        this.departmentOptions = []
      })
    },
    getDepartmentName(deptId) {
      if (!deptId || !this.departmentOptions.length) return '-'
      const dept = this.departmentOptions.find(item => item.id === deptId)
      return dept ? dept.name : '-'
    },
    handleFilter() {
      this.listQuery.page = 1
      this.getList()
    },
    resetTemp() {
      this.temp = {
        id: undefined,
        name: '',
        code: '',
        dept_id: undefined,
        description: '',
        owner: '',
        email: '',
        status: 1
      }
    },
    handleCreate() {
      this.resetTemp()
      // 如果不是管理员，自动设置为当前用户的部门
      if (!this.isAdmin && this.deptId) {
        this.temp.dept_id = this.deptId
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
          createBusiness(this.temp).then(() => {
            this.getList() // 重新获取列表而不是直接添加到前端
            this.dialogFormVisible = false
            this.$notify({
              title: '成功',
              message: '创建成功',
              type: 'success',
              duration: 2000
            })
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
      // 检查权限：非管理员用户只能编辑自己部门的业务线
      if (!this.isAdmin && row.dept_id !== this.deptId) {
        this.$message.warning('您只能编辑自己部门的业务线')
        return
      }
      
      this.temp = Object.assign({}, row) // copy obj
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
          // 后端UpdateBusinessRequest需要ID字段，所以不能从数据中移除ID
          updateBusiness(tempData.id, tempData).then(() => {
            this.getList() // 重新获取列表以确保数据同步
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
              duration: 2000
            })
          })
        }
      })
    },
    handleDelete(row, index) {
      // 检查权限：非管理员用户只能删除自己部门的业务线
      if (!this.isAdmin && row.dept_id !== this.deptId) {
        this.$message.warning('您只能删除自己部门的业务线')
        return
      }
      
      this.$confirm('此操作将永久删除该业务, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteBusiness(row.id).then(() => {
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