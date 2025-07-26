<template>
  <div class="dashboard-container">
    <div class="dashboard-text">欢迎使用 API 管理平台</div>
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="hover">
          <div slot="header" class="clearfix">
            <span>API 总数</span>
          </div>
          <div class="card-panel">
            <div class="card-panel-icon-wrapper">
              <svg-icon icon-class="dashboard" class-name="card-panel-icon" />
            </div>
            <div class="card-panel-description">
              <div class="card-panel-text">API 总数</div>
              <count-to :start-val="0" :end-val="apiCount" :duration="2000" class="card-panel-num" />
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div slot="header" class="clearfix">
            <span>业务线总数</span>
          </div>
          <div class="card-panel">
            <div class="card-panel-icon-wrapper" style="background: #40c9c6;">
              <svg-icon icon-class="business" class-name="card-panel-icon" />
            </div>
            <div class="card-panel-description">
              <div class="card-panel-text">业务线总数</div>
              <count-to :start-val="0" :end-val="businessCount" :duration="2000" class="card-panel-num" />
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div slot="header" class="clearfix">
            <span>部门总数</span>
          </div>
          <div class="card-panel">
            <div class="card-panel-icon-wrapper" style="background: #36a3f7;">
              <svg-icon icon-class="department" class-name="card-panel-icon" />
            </div>
            <div class="card-panel-description">
              <div class="card-panel-text">部门总数</div>
              <count-to :start-val="0" :end-val="departmentCount" :duration="2000" class="card-panel-num" />
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div slot="header" class="clearfix">
            <span>用户总数</span>
          </div>
          <div class="card-panel">
            <div class="card-panel-icon-wrapper" style="background: #f4516c;">
              <svg-icon icon-class="user" class-name="card-panel-icon" />
            </div>
            <div class="card-panel-description">
              <div class="card-panel-text">用户总数</div>
              <count-to :start-val="0" :end-val="userCount" :duration="2000" class="card-panel-num" />
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="12">
        <el-card shadow="hover">
          <div slot="header" class="clearfix">
            <span>API 分类统计</span>
          </div>
          <div style="height: 300px" ref="apiCategoryChart"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover">
          <div slot="header" class="clearfix">
            <span>业务线 API 数量统计</span>
          </div>
          <div style="height: 300px" ref="businessApiChart"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
// 导入 ECharts 图表库
import * as echarts from 'echarts';
import CountTo from 'vue-count-to';
import { mapGetters } from 'vuex';
import { getDashboardStats, getApiCategoryStats, getBusinessApiStats } from '@/api/dashboard'
import { getApiList } from '@/api/api'
import { getBusinessList } from '@/api/business'
import { getDepartmentList } from '@/api/department'
import { getUserList } from '@/api/user'

export default {
  name: 'DashboardView',
  components: {
    CountTo
  },
  data() {
    return {
      apiCount: 0,
      businessCount: 0,
      departmentCount: 0,
      userCount: 0,
      apiCategoryChart: null,
      businessApiChart: null,
      loading: false
    }
  },
  computed: {
    ...mapGetters([
      'name'
    ])
  },
  mounted() {
    this.fetchDashboardData()
    this.initCharts()
    window.addEventListener('resize', this.resizeCharts)
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.resizeCharts)
    if (this.apiCategoryChart) {
      this.apiCategoryChart.dispose()
    }
    if (this.businessApiChart) {
      this.businessApiChart.dispose()
    }
  },
  methods: {
    async fetchDashboardData() {
      this.loading = true
      try {
        // 尝试使用专门的dashboard API
        try {
          const response = await getDashboardStats()
          if (response.code === 0) {
            this.apiCount = response.data.api_count
            this.businessCount = response.data.business_count
            this.departmentCount = response.data.department_count
            this.userCount = response.data.user_count
            return
          }
        } catch (error) {
          console.log('Dashboard API not available, using fallback method')
        }

        // 如果dashboard API不可用，使用现有API获取统计数据
        const [apiRes, businessRes, departmentRes, userRes] = await Promise.all([
          getApiList({ page: 1, size: 1 }).catch(() => ({ data: { total: 0 } })),
          getBusinessList({ page: 1, size: 1 }).catch(() => ({ data: { total: 0 } })),
          getDepartmentList({ page: 1, size: 1 }).catch(() => ({ data: { total: 0 } })),
          getUserList({ page: 1, size: 1 }).catch(() => ({ data: { total: 0 } }))
        ])

        this.apiCount = apiRes.data?.total || 0
        this.businessCount = businessRes.data?.total || 0
        this.departmentCount = departmentRes.data?.total || 0
        this.userCount = userRes.data?.total || 0
      } catch (error) {
        console.error('Failed to load dashboard stats:', error)
        this.$message.error('加载统计数据失败')
      } finally {
        this.loading = false
      }
    },
    async initCharts() {
      this.$nextTick(async () => {
        await this.initApiCategoryChart()
        await this.initBusinessApiChart()
      })
    },
    async initApiCategoryChart() {
      this.apiCategoryChart = echarts.init(this.$refs.apiCategoryChart)
      
      try {
        // 尝试使用专门的API分类统计接口
        let categoryData = []
        try {
          const response = await getApiCategoryStats()
          if (response.code === 0) {
            categoryData = response.data
          }
        } catch (error) {
          console.log('API category stats not available, using fallback method')
          // 如果专门的接口不可用，获取所有API并统计
          const apiResponse = await getApiList({ page: 1, size: 1000 })
          if (apiResponse.code === 0 && apiResponse.data.list) {
            const methodCounts = {}
            apiResponse.data.list.forEach(api => {
              const method = api.method || 'GET'
              methodCounts[method] = (methodCounts[method] || 0) + 1
            })
            categoryData = Object.entries(methodCounts).map(([name, value]) => ({ name, value }))
          }
        }

        // 如果没有数据，使用默认数据
        if (categoryData.length === 0) {
          categoryData = [
            { value: 0, name: 'GET' },
            { value: 0, name: 'POST' },
            { value: 0, name: 'PUT' },
            { value: 0, name: 'DELETE' },
            { value: 0, name: 'PATCH' }
          ]
        }

        const option = {
          tooltip: {
            trigger: 'item',
            formatter: '{a} <br/>{b}: {c} ({d}%)'
          },
          legend: {
            orient: 'vertical',
            left: 10,
            data: categoryData.map(item => item.name)
          },
          series: [
            {
              name: 'API分类',
              type: 'pie',
              radius: ['50%', '70%'],
              avoidLabelOverlap: false,
              label: {
                show: false,
                position: 'center'
              },
              emphasis: {
                label: {
                  show: true,
                  fontSize: '18',
                  fontWeight: 'bold'
                }
              },
              labelLine: {
                show: false
              },
              data: categoryData
            }
          ]
        }
        
        this.apiCategoryChart.setOption(option)
      } catch (error) {
        console.error('Failed to load API category chart:', error)
      }
    },
    async initBusinessApiChart() {
      // 使用 echarts 库初始化业务 API 图表
      this.businessApiChart = echarts.init(this.$refs.businessApiChart)
      
      try {
        // 尝试使用专门的业务线API统计接口
        let businessApiData = []
        try {
          const response = await getBusinessApiStats()
          if (response.code === 0) {
            businessApiData = response.data
          }
        } catch (error) {
          console.log('Business API stats not available, using fallback method')
          // 如果专门的接口不可用，获取业务线列表并统计每个业务线的API数量
          const businessResponse = await getBusinessList({ page: 1, size: 100 })
          if (businessResponse.code === 0 && businessResponse.data.list) {
            const businesses = businessResponse.data.list
            const businessApiCounts = await Promise.all(
              businesses.map(async (business) => {
                try {
                  const apiResponse = await getApiList({ businessId: business.id, page: 1, size: 1 })
                  return {
                    name: business.name,
                    value: apiResponse.data?.total || 0
                  }
                } catch (error) {
                  return {
                    name: business.name,
                    value: 0
                  }
                }
              })
            )
            businessApiData = businessApiCounts
          }
        }

        // 如果没有数据，使用默认数据
        if (businessApiData.length === 0) {
          businessApiData = [
            { name: '暂无数据', value: 0 }
          ]
        }

        const option = {
          tooltip: {
            trigger: 'axis',
            axisPointer: {
              type: 'shadow'
            }
          },
          grid: {
            left: '3%',
            right: '4%',
            bottom: '3%',
            containLabel: true
          },
          xAxis: {
            type: 'value',
            boundaryGap: [0, 0.01]
          },
          yAxis: {
            type: 'category',
            data: businessApiData.map(item => item.name)
          },
          series: [
            {
              name: 'API数量',
              type: 'bar',
              data: businessApiData.map(item => item.value)
            }
          ]
        }
        
        this.businessApiChart.setOption(option)
      } catch (error) {
        console.error('Failed to load business API chart:', error)
      }
    },
    resizeCharts() {
      if (this.apiCategoryChart) {
        this.apiCategoryChart.resize()
      }
      if (this.businessApiChart) {
        this.businessApiChart.resize()
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.dashboard {
  &-container {
    margin: 30px;
  }
  &-text {
    font-size: 30px;
    line-height: 46px;
    margin-bottom: 20px;
  }
}

.card-panel {
  height: 108px;
  display: flex;
  font-size: 12px;
  position: relative;
  overflow: hidden;
  color: #666;
  background: #fff;
  box-shadow: 4px 4px 40px rgba(0, 0, 0, .05);
  border-color: rgba(0, 0, 0, .05);
  
  &-icon-wrapper {
    float: left;
    margin: 14px 0 0 14px;
    padding: 16px;
    transition: all 0.38s ease-out;
    border-radius: 6px;
    background: #34bfa3;
  }
  
  &-icon {
    float: left;
    font-size: 48px;
    color: #fff;
  }
  
  &-description {
    float: right;
    flex: 1;
    padding: 26px 15px 0 15px;
    margin-left: 10px;
    text-align: right;
  }
  
  &-text {
    line-height: 18px;
    color: rgba(0, 0, 0, 0.45);
    font-size: 16px;
    margin-bottom: 12px;
  }
  
  &-num {
    font-size: 20px;
    font-weight: bold;
  }
}
</style>