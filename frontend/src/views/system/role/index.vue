<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input
        v-model="listQuery.name"
        placeholder="角色名称"
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
      <el-table-column label="角色名称" width="150px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.name }}</span>
        </template>
      </el-table-column>
      <el-table-column label="角色编码" width="120px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.code }}</span>
        </template>
      </el-table-column>
      <el-table-column label="所属部门" width="150px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.dept_name || '未知部门' }}</span>
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
          <el-button type="primary" size="mini" @click="handleUpdate(row)">
            编辑
          </el-button>
          <el-button v-if="row.status!='deleted'" size="mini" type="danger" @click="handleDelete(row,$index)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />

    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible">
      <el-form ref="dataForm" :rules="rules" :model="temp" label-position="left" label-width="80px" style="width: 400px; margin-left:50px;">
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="temp.name" />
        </el-form-item>
        <el-form-item label="角色编码" prop="code">
          <el-input v-model="temp.code" placeholder="请输入角色编码，如：admin、user等" />
        </el-form-item>
        <el-form-item label="所属部门" prop="dept_id">
          <el-select 
            v-model="temp.dept_id" 
            placeholder="请选择部门" 
            style="width: 100%"
            :disabled="!roles.includes('admin')"
          >
            <el-option
              v-for="dept in departmentOptions"
              :key="dept.id"
              :label="dept.name"
              :value="dept.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="temp.description" type="textarea" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="temp.status" class="filter-item" placeholder="请选择">
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
import { mapGetters } from 'vuex';

import { getDepartmentList } from '@/api/department';
import {
  createRole,
  deleteRole,
  getRoleList,
  updateRole,
} from '@/api/role';
import Pagination
  from '@/components/Pagination'; // secondary package based on el-pagination
import waves from '@/directive/waves'; // waves directive
import { parseTime } from '@/utils';

export default {
  name: 'RoleTable',
  components: { Pagination },
  directives: { waves },
  computed: {
    ...mapGetters(['roles', 'deptId'])
  },
  filters: {
    statusFilter(status) {
      const statusMap = {
        1: 'success',
        0: 'info'
      }
      return statusMap[status]
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
      departmentOptions: [],
      statusOptions: [
        { key: 1, display_name: '启用' },
        { key: 0, display_name: '禁用' }
      ],
      temp: {
        id: undefined,
        name: '',
        code: '',
        description: '',
        dept_id: undefined,
        status: 1
      },
      dialogFormVisible: false,
      dialogStatus: '',
      textMap: {
        update: '编辑',
        create: '创建'
      },
      rules: {
        name: [{ required: true, message: '角色名称是必填项', trigger: 'blur' }],
        code: [{ required: true, message: '角色编码是必填项', trigger: 'blur' }],
        dept_id: [{ required: true, message: '所属部门是必填项', trigger: 'change' }]
      }
    }
  },
  created() {
    this.getList()
    this.getDepartmentOptions()
  },
  methods: {
    getList() {
      this.listLoading = true
      
      // 构建查询参数
      const queryParams = { ...this.listQuery }
      
      // 如果不是admin用户，只查询自己部门的角色
      if (!this.roles.includes('admin') && this.deptId) {
        queryParams.dept_id = this.deptId
      }
      
      getRoleList(queryParams).then(response => {
        this.list = response.data.items
        this.total = response.data.total
        this.listLoading = false
      })
    },
    getDepartmentOptions() {
      if (this.roles.includes('admin')) {
        // admin用户可以看到所有部门
        getDepartmentList({ page: 1, limit: 1000 }).then(response => {
          this.departmentOptions = response.data.items || []
        }).catch(() => {
          this.departmentOptions = []
        })
      } else {
        // 非admin用户只能看到自己的部门
        if (this.deptId) {
          getDepartmentList({ page: 1, limit: 1000 }).then(response => {
            const allDepts = response.data.items || []
            this.departmentOptions = allDepts.filter(dept => dept.id === this.deptId)
          }).catch(() => {
            this.departmentOptions = []
          })
        } else {
          this.departmentOptions = []
        }
      }
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
        description: '',
        dept_id: undefined,
        status: 1
      }
    },
    handleCreate() {
      this.resetTemp()
      
      // 如果不是admin用户，自动设置为当前用户的部门
      if (!this.roles.includes('admin') && this.deptId) {
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
          createRole(this.temp).then(() => {
            this.dialogFormVisible = false
            this.getList() // 重新获取列表以确保数据同步
            this.$notify({
              title: '成功',
              message: '创建成功',
              type: 'success',
              duration: 2000
            })
          }).catch(error => {
            // 显示友好的错误信息
            const message = error.response?.data?.message || '创建失败'
            this.$notify({
              title: '错误',
              message: message,
              type: 'error',
              duration: 3000
            })
          })
        }
      })
    },
    handleUpdate(row) {
      this.temp = Object.assign({}, row) // copy obj
      
      // // 编辑时，如果角色编码包含部门前缀，则去掉前缀显示
      // if (this.temp.code && this.temp.code.includes('_')) {
      //   // 找到对应的部门
      //   const dept = this.departmentOptions.find(d => d.id === this.temp.dept_id)
      //   if (dept && this.temp.code.startsWith(dept.code + '_')) {
      //     // 去掉部门前缀，只显示角色编码部分
      //     this.temp.code = this.temp.code.substring(dept.code.length + 1)
      //   }
      // }
      
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
          // 后端UpdateRoleRequest需要ID字段，所以不能从数据中移除ID
          updateRole(tempData.id, tempData).then(() => {
            const index = this.list.findIndex(v => v.id === this.temp.id)
            this.list.splice(index, 1, this.temp)
            this.dialogFormVisible = false
            this.$notify({
              title: '成功',
              message: '更新成功',
              type: 'success',
              duration: 2000
            })
            // 重新获取列表以确保数据同步
            this.getList()
          }).catch(error => {
            // 显示友好的错误信息
            const message = error.response?.data?.message || '更新失败'
            this.$notify({
              title: '错误',
              message: message,
              type: 'error',
              duration: 3000
            })
          })
        }
      })
    },
    handleDelete(row, index) {
      this.$confirm('此操作将永久删除该角色, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteRole(row.id).then(() => {
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