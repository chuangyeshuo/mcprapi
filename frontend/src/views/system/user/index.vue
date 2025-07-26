<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input
        v-model="listQuery.username"
        placeholder="用户名"
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
      <el-table-column label="用户名" width="150px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.username }}</span>
        </template>
      </el-table-column>
      <el-table-column label="姓名" width="150px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.name }}</span>
        </template>
      </el-table-column>
      <el-table-column label="邮箱" width="200px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.email }}</span>
        </template>
      </el-table-column>
      <el-table-column label="角色" align="center" width="200">
        <template slot-scope="{row}">
          <span v-if="row.roles && row.roles.length > 0">
            {{ row.roles.join(', ') }}
          </span>
          <span v-else style="color: #999;">无角色</span>
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
      <el-table-column label="操作" align="center" width="350" class-name="small-padding fixed-width">
        <template slot-scope="{row,$index}">
          <el-button type="primary" size="mini" @click="handleUpdate(row)">
            编辑
          </el-button>
          <el-button type="success" size="mini" @click="handleRoles(row)">
            角色
          </el-button>
          <el-button type="warning" size="mini" style="background-color: #8B4513; border-color: #8B4513;" @click="handleToken(row)">
            Token
          </el-button>
          <el-button v-if="row.status!='deleted'" size="mini" type="danger" @click="handleDelete(row,$index)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />

    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible">
      <el-form ref="dataForm" :rules="rules" :model="temp" label-position="left" label-width="70px" style="width: 400px; margin-left:50px;">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="temp.username" placeholder="英文数字组合或邮箱前缀，如：john123 或 john.doe" />
        </el-form-item>
        <el-form-item label="姓名" prop="name">
          <el-input v-model="temp.name" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="temp.email" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="temp.password" type="password" />
        </el-form-item>
        <el-form-item label="头像" prop="avatar">
          <el-input v-model="temp.avatar" placeholder="请输入头像URL地址" />
        </el-form-item>
        <el-form-item label="部门" prop="dept_id">
          <el-select v-model="temp.dept_id" class="filter-item" placeholder="请选择部门" style="width: 100%;" :disabled="!roles.includes('admin')">
            <el-option v-for="item in departmentOptions" :key="item.id" :label="item.name" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="角色" prop="role_ids">
          <el-select v-model="temp.role_ids" multiple class="filter-item" placeholder="请选择角色" style="width: 100%;">
            <el-option v-for="item in roleOptions" :key="item.id" :label="item.name" :value="item.id" />
          </el-select>
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

    <!-- 角色管理对话框 -->
    <el-dialog title="用户角色管理" :visible.sync="roleDialogVisible" width="500px">
      <div style="margin-bottom: 20px;">
        <p><strong>用户：</strong>{{ currentUser.name }} ({{ currentUser.username }})</p>
        <p><strong>部门：</strong>{{ currentUser.dept_name }}</p>
      </div>
      
      <el-form>
        <el-form-item label="选择角色：">
          <el-checkbox-group v-model="roleForm.selectedRoleCodes">
            <el-checkbox 
              v-for="role in roleOptions" 
              :key="role.id" 
              :label="role.code"
            >
              {{ role.name }}
            </el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </el-form>
      
      <div slot="footer" class="dialog-footer">
        <el-button @click="roleDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveUserRoles">确定</el-button>
      </div>
    </el-dialog>

    <!-- Token管理对话框 -->
    <el-dialog title="用户Token管理" :visible.sync="tokenDialogVisible" width="600px">
      <div style="margin-bottom: 20px;">
        <p><strong>用户：</strong>{{ currentUser.name }} ({{ currentUser.username }})</p>
        <p><strong>邮箱：</strong>{{ currentUser.email }}</p>
      </div>
      
      <div style="margin-bottom: 20px;">
        <h4>JWT Token:</h4>
        <el-input
          v-model="userToken"
          type="textarea"
          :rows="8"
          readonly
          placeholder="加载中..."
          style="font-family: monospace; font-size: 12px;"
        />
      </div>
      
      <div slot="footer" class="dialog-footer" style="text-align: center;">
        <el-button type="danger" @click="showRefreshTokenDialog">刷新Token</el-button>
        <el-button type="warning" @click="showRefreshTokenWithVersionDialog">刷新Token+版本号</el-button>
        <el-button @click="tokenDialogVisible = false">关闭</el-button>
      </div>
    </el-dialog>

    <!-- 刷新Token对话框 -->
    <el-dialog :title="refreshTokenForm.withVersion ? '刷新Token+版本号' : '刷新Token'" :visible.sync="refreshTokenDialogVisible" width="500px">
      <div v-if="refreshTokenForm.withVersion" style="margin-bottom: 15px; padding: 10px; background-color: #fff6f7; border: 1px solid #fbc4c4; border-radius: 4px;">
        <i class="el-icon-warning" style="color: #f56c6c;"></i>
        <span style="color: #f56c6c; font-weight: bold;">警告：</span>
        <span style="color: #606266;">刷新Token+版本号将使所有旧Token立即失效，请确保其他系统已准备好使用新Token。</span>
      </div>
      
      <el-form ref="refreshTokenForm" :model="refreshTokenForm" :rules="refreshTokenRules" label-width="100px">
        <el-form-item label="过期时间" prop="expireDays">
          <el-input-number
            v-model="refreshTokenForm.expireDays"
            :min="1"
            :max="365"
            placeholder="请输入天数"
            style="width: 100%;"
          />
          <div style="color: #999; font-size: 12px; margin-top: 5px;">
            Token将在 {{ refreshTokenForm.expireDays }} 天后过期
          </div>
        </el-form-item>
      </el-form>
      
      <div slot="footer" class="dialog-footer">
        <el-button @click="refreshTokenDialogVisible = false">取消</el-button>
        <el-button 
          :type="refreshTokenForm.withVersion ? 'warning' : 'danger'" 
          @click="confirmRefreshToken"
        >
          {{ refreshTokenForm.withVersion ? '确认刷新+版本号' : '确认刷新' }}
        </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { getUserList, createUser, updateUser, deleteUser, getUserRoles, assignUserRoles, getUserToken, refreshUserToken, refreshUserTokenWithVersion } from '@/api/user'
import { getUserAccessibleRoles } from '@/api/role'
import { getDepartmentList } from '@/api/department'
import waves from '@/directive/waves' // waves directive
import { parseTime } from '@/utils'
import Pagination from '@/components/Pagination' // secondary package based on el-pagination
import { mapGetters } from 'vuex'

export default {
  name: 'UserTable',
  components: { Pagination },
  directives: { waves },
  computed: {
    ...mapGetters([
      'roles',
      'deptId'
    ])
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
        username: undefined
      },
      statusOptions: [
        { key: 1, display_name: '启用' },
        { key: 0, display_name: '禁用' }
      ],
      departmentOptions: [],
      roleOptions: [],
      temp: {
        id: undefined,
        username: '',
        name: '',
        email: '',
        password: '',
        avatar: '', // 头像URL
        dept_id: undefined,
        role_ids: [],
        status: 1
      },
      dialogFormVisible: false,
      dialogStatus: '',
      textMap: {
        update: '编辑',
        create: '创建'
      },
      // 角色管理相关
      roleDialogVisible: false,
      currentUser: {},
      roleForm: {
        selectedRoleCodes: []
      },
      // Token管理相关
      tokenDialogVisible: false,
      refreshTokenDialogVisible: false,
      userToken: '',
      refreshTokenForm: {
        expireDays: 30,
        withVersion: false
      },
      refreshTokenRules: {
        expireDays: [
          { required: true, message: '过期天数是必填项', trigger: 'blur' },
          { type: 'number', min: 1, max: 365, message: '过期天数必须在1-365天之间', trigger: 'blur' }
        ]
      },
      rules: {
        username: [
          { required: true, message: '用户名是必填项', trigger: 'blur' },
          { 
            pattern: /^[a-zA-Z0-9._-]+$/, 
            message: '用户名只能包含英文字母、数字、点号、下划线和短横线', 
            trigger: 'blur' 
          },
          { 
            min: 3, 
            max: 50, 
            message: '用户名长度在 3 到 50 个字符', 
            trigger: 'blur' 
          }
        ],
        name: [{ required: true, message: '姓名是必填项', trigger: 'blur' }],
        email: [{ required: true, message: '邮箱是必填项', trigger: 'blur' }],
        password: [{ required: true, message: '密码是必填项', trigger: 'blur' }],
        dept_id: [{ required: true, message: '部门是必填项', trigger: 'change' }]
      }
    }
  },
  created() {
    this.getList()
    this.getDepartmentOptions()
    this.getRoleOptions()
  },
  methods: {
    getList() {
      this.listLoading = true
      
      // 构建查询参数
      const queryParams = { ...this.listQuery }
      
      // 如果不是admin用户，只获取自己部门的用户
      if (!this.roles.includes('admin') && this.deptId) {
        queryParams.dept_id = this.deptId
      }
      
      getUserList(queryParams).then(response => {
        this.list = response.data.items
        this.total = response.data.total
        this.listLoading = false
      })
    },
    getDepartmentOptions() {
      getDepartmentList({ page: 1, page_size: 1000 }).then(response => {
        const allDepartments = response.data.items || []
        
        // 如果不是admin用户，只显示自己的部门
        if (!this.roles.includes('admin') && this.deptId) {
          this.departmentOptions = allDepartments.filter(dept => dept.id === this.deptId)
        } else {
          this.departmentOptions = allDepartments
        }
      })
    },
    getRoleOptions() {
      getUserAccessibleRoles().then(response => {
        this.roleOptions = response.data || []
      })
    },
    handleFilter() {
      this.listQuery.page = 1
      this.getList()
    },
    resetTemp() {
      this.temp = {
        id: undefined,
        username: '',
        name: '',
        email: '',
        password: '',
        dept_id: undefined,
        role_ids: [],
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
          createUser(this.temp).then(() => {
            this.list.unshift(this.temp)
            this.dialogFormVisible = false
            this.$notify({
              title: '成功',
              message: '创建成功',
              type: 'success',
              duration: 2000
            })
            this.getList() // 重新获取列表以确保数据同步
          })
        }
      })
    },
    handleUpdate(row) {
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
          // 后端UpdateUserRequest需要ID字段，所以不能从数据中移除ID
          updateUser(tempData.id, tempData).then(() => {
            const index = this.list.findIndex(v => v.id === this.temp.id)
            this.list.splice(index, 1, this.temp)
            this.dialogFormVisible = false
            this.$notify({
              title: '成功',
              message: '更新成功',
              type: 'success',
              duration: 2000
            })
          })
        }
      })
    },
    // 处理角色按钮点击
    async handleRoles(row) {
      this.currentUser = { ...row }
      this.roleDialogVisible = true
      
      try {
        // 获取用户当前角色代码
        const response = await getUserRoles(row.id)
        // 后端返回的是角色代码数组
        this.roleForm.selectedRoleCodes = response.data || []
      } catch (error) {
        console.error('获取用户角色失败:', error)
        this.roleForm.selectedRoleCodes = []
      }
    },
    async saveUserRoles() {
      try {
        // 将角色代码转换为角色ID
        const roleIds = this.roleOptions
          .filter(role => this.roleForm.selectedRoleCodes.includes(role.code))
          .map(role => role.id)
        
        await assignUserRoles({
          user_id: this.currentUser.id,
          role_ids: roleIds
        })
        
        this.$message.success('角色分配成功')
        this.roleDialogVisible = false
        this.getList() // 刷新用户列表
      } catch (error) {
        console.error('分配角色失败:', error)
        this.$message.error('分配角色失败')
      }
    },
    handleDelete(row, index) {
      this.$confirm('此操作将永久删除该用户, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteUser(row.id).then(() => {
          this.$notify({
            title: '成功',
            message: '删除成功',
            type: 'success',
            duration: 2000
          })
          this.list.splice(index, 1)
        })
      })
    },
    // 处理Token按钮点击
    async handleToken(row) {
      this.currentUser = { ...row }
      this.tokenDialogVisible = true
      this.userToken = '加载中...'
      
      try {
        const response = await getUserToken(row.id)
        this.userToken = response.data.token || '无Token信息'
      } catch (error) {
        console.error('获取用户Token失败:', error)
        this.userToken = '获取Token失败，请稍后重试'
        this.$message.error('获取Token失败')
      }
    },
    // 显示刷新Token对话框
    showRefreshTokenDialog() {
      this.refreshTokenForm.expireDays = 30
      this.refreshTokenForm.withVersion = false
      this.refreshTokenDialogVisible = true
    },
    // 显示刷新Token+版本号对话框
    showRefreshTokenWithVersionDialog() {
      this.refreshTokenForm.expireDays = 30
      this.refreshTokenForm.withVersion = true
      this.refreshTokenDialogVisible = true
    },
    // 确认刷新Token
    async confirmRefreshToken() {
      this.$refs['refreshTokenForm'].validate(async (valid) => {
        if (valid) {
          try {
            let response
            if (this.refreshTokenForm.withVersion) {
              // 调用刷新Token+版本号API
              response = await refreshUserTokenWithVersion({
                user_id: this.currentUser.id,
                expire_days: this.refreshTokenForm.expireDays
              })
            } else {
              // 调用普通刷新Token API
              response = await refreshUserToken({
                user_id: this.currentUser.id,
                expire_days: this.refreshTokenForm.expireDays
              })
            }
            
            this.userToken = response.data.token
            this.refreshTokenDialogVisible = false
            
            if (this.refreshTokenForm.withVersion) {
              this.$message.success('Token+版本号刷新成功，所有旧Token已失效')
            } else {
              this.$message.success('Token刷新成功')
            }
          } catch (error) {
            console.error('刷新Token失败:', error)
            this.$message.error('刷新Token失败')
          }
        }
      })
    }
  }
}
</script>