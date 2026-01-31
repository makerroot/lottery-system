<template>
  <div class="operation-logs-page">
    <!-- 权限不足提示 -->
    <div v-if="!isAdmin" class="permission-denied">
      <a-result
        status="403"
        title="请先登录"
        sub-title="您需要登录管理员账号才能访问此页面"
      >
        <template #extra>
          <router-link to="/admin">
            <a-button type="primary" size="large">前往登录</a-button>
          </router-link>
        </template>
      </a-result>
    </div>

    <!-- 正常内容 -->
    <template v-else>
      <!-- 页面头部 -->
      <div class="page-header">
        <div class="header-left">
          <h1 class="page-title font-display">📋 操作日志</h1>
          <p class="page-description font-body">
            {{ isSuperAdmin ? '查看所有系统操作记录' : '查看您所在公司的操作记录' }}
          </p>
        </div>
      </div>

      <!-- 统计卡片 -->
    <a-row :gutter="16" class="stats-row">
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <div class="stat-icon" style="background: rgba(24, 144, 255, 0.1); color: #1890ff;">📝</div>
          <div class="stat-info">
            <div class="stat-value font-display">{{ stats.create || 0 }}</div>
            <div class="stat-label font-body">新增操作</div>
          </div>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <div class="stat-icon" style="background: rgba(82, 196, 26, 0.1); color: #52c41a;">✏️</div>
          <div class="stat-info">
            <div class="stat-value font-display">{{ stats.update || 0 }}</div>
            <div class="stat-label font-body">更新操作</div>
          </div>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <div class="stat-icon" style="background: rgba(255, 77, 79, 0.1); color: #ff4d4f;">🗑️</div>
          <div class="stat-info">
            <div class="stat-value font-display">{{ stats.delete || 0 }}</div>
            <div class="stat-label font-body">删除操作</div>
          </div>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="12" :md="6">
        <a-card class="stat-card">
          <div class="stat-icon" style="background: rgba(115, 56, 240, 0.1); color: #7338ff;">🔐</div>
          <div class="stat-info">
            <div class="stat-value font-display">{{ stats.login || 0 }}</div>
            <div class="stat-label font-body">登录记录</div>
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 筛选条件 -->
    <a-card class="filter-card">
      <a-form layout="inline">
        <a-form-item label="操作类型">
          <a-select
            v-model:value="filters.action"
            placeholder="全部"
            style="width: 150px"
            allowClear
            @change="fetchLogs"
          >
            <a-select-option value="">全部</a-select-option>
            <a-select-option value="create">新增</a-select-option>
            <a-select-option value="update">更新</a-select-option>
            <a-select-option value="delete">删除</a-select-option>
            <a-select-option value="login">登录</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="资源类型">
          <a-select
            v-model:value="filters.resource"
            placeholder="全部"
            style="width: 150px"
            allowClear
            @change="fetchLogs"
          >
            <a-select-option value="">全部</a-select-option>
            <a-select-option value="admin">管理员</a-select-option>
            <a-select-option value="company">公司</a-select-option>
            <a-select-option value="user">用户</a-select-option>
            <a-select-option value="prize_level">奖项等级</a-select-option>
            <a-select-option value="prize">奖品</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="fetchLogs" :loading="loading">
            刷新
          </a-button>
        </a-form-item>
      </a-form>
    </a-card>

    <!-- 日志列表 -->
    <a-card class="logs-card">
      <a-table
        :columns="columns"
        :data-source="logs"
        :loading="loading"
        :pagination="pagination"
        @change="handleTableChange"
        rowKey="id"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'action'">
            <a-tag :color="getActionColor(record.action)">
              {{ getActionLabel(record.action) }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'resource'">
            <a-tag color="blue">{{ getResourceLabel(record.resource) }}</a-tag>
          </template>
          <template v-else-if="column.key === 'details'">
            <div class="details-text">{{ record.details }}</div>
          </template>
          <template v-else-if="column.key === 'created_at'">
            <span class="font-mono">{{ formatTime(record.created_at) }}</span>
          </template>
        </template>
      </a-table>
    </a-card>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import request from '../../utils/request'
import { useAdmin } from '../../utils/admin'

const router = useRouter()
const { isAdmin, currentUser, isSuperAdmin } = useAdmin()

// 页面加载时检查权限
onMounted(() => {
  // 检查是否登录
  if (!isAdmin.value) {
    message.error('请先登录')
    router.push('/admin')
    return
  }

  // 加载数据（普通管理员和超级管理员都可以查看）
  fetchLogs()
  fetchStats()
})

// 数据
const logs = ref([])
const stats = ref({})
const loading = ref(false)

// 筛选条件
const filters = ref({
  action: '',
  resource: ''
})

// 分页配置
const pagination = ref({
  current: 1,
  pageSize: 20,
  total: 0,
  showSizeChanger: true,
  showTotal: (total) => `共 ${total} 条`
})

// 表格列配置
const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
  { title: '操作人', dataIndex: 'admin_name', key: 'admin_name', width: 120 },
  { title: '操作类型', dataIndex: 'action', key: 'action', width: 100 },
  { title: '资源类型', dataIndex: 'resource', key: 'resource', width: 120 },
  { title: '操作详情', dataIndex: 'details', key: 'details', ellipsis: true },
  { title: 'IP地址', dataIndex: 'ip_address', key: 'ip_address', width: 150 },
  { title: '操作时间', dataIndex: 'created_at', key: 'created_at', width: 180 }
]

// 获取操作标签
const getActionLabel = (action) => {
  const labels = {
    create: '新增',
    update: '更新',
    delete: '删除',
    login: '登录',
    logout: '退出'
  }
  return labels[action] || action
}

// 获取操作颜色
const getActionColor = (action) => {
  const colors = {
    create: 'success',
    update: 'processing',
    delete: 'error',
    login: 'blue',
    logout: 'default'
  }
  return colors[action] || 'default'
}

// 获取资源标签
const getResourceLabel = (resource) => {
  const labels = {
    admin: '管理员',
    company: '公司',
    user: '用户',
    prize_level: '奖项等级',
    prize: '奖品'
  }
  return labels[resource] || resource
}

// 格式化时间
const formatTime = (timeStr) => {
  const date = new Date(timeStr)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

// 获取日志列表
const fetchLogs = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.value.current,
      page_size: pagination.value.pageSize
    }

    if (filters.value.action) {
      params.action = filters.value.action
    }
    if (filters.value.resource) {
      params.resource = filters.value.resource
    }

    const data = await request.get('/admin/operation-logs', { params })
    logs.value = data.data || []
    pagination.value.total = data.total || 0
  } catch (error) {
    message.error('获取操作日志失败')
  } finally {
    loading.value = false
  }
}

// 获取统计数据
const fetchStats = async () => {
  try {
    const data = await request.get('/admin/operation-stats')
    const statsMap = {}
    data.forEach(item => {
      statsMap[item.action] = item.count
    })
    stats.value = statsMap
  } catch (error) {
    // 静默失败，不影响主要功能
  }
}

// 表格变化处理
const handleTableChange = (pag, filters, sorter) => {
  pagination.value.current = pag.current
  pagination.value.pageSize = pag.pageSize
  fetchLogs()
}
</script>

<style scoped>
.operation-logs-page {
  padding: var(--spacing-xl);
}

.page-header {
  margin-bottom: var(--spacing-xl);
}

.page-title {
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  margin: 0 0 var(--spacing-xs);
}

.page-description {
  font-size: var(--font-size-base);
  color: var(--text-secondary);
  margin: 0;
}

.stats-row {
  margin-bottom: var(--spacing-lg);
}

.stat-card {
  background: var(--company-color-bg);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.stat-card :deep(.ant-card-body) {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-lg);
}

.stat-icon {
  width: 50px;
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-2xl);
  border-radius: var(--radius-lg);
}

.stat-value {
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  line-height: 1.2;
}

.stat-label {
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
  margin-top: var(--spacing-xs);
}

.filter-card {
  margin-bottom: var(--spacing-lg);
}

.logs-card :deep(.ant-table) {
  background: transparent;
}

.logs-card :deep(.ant-table-thead > tr > th) {
  background: var(--bg-elevated);
  color: var(--text-primary);
  border-bottom: 1px solid var(--border-color);
}

.logs-card :deep(.ant-table-tbody > tr > td) {
  background: transparent;
  border-bottom: 1px solid var(--border-color);
  color: var(--text-secondary);
}

.logs-card :deep(.ant-table-tbody > tr:hover > td) {
  background: var(--company-color-bg);
}

.details-text {
  max-width: 400px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
