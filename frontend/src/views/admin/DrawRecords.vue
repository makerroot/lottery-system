<template>
  <div class="drawrecords-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title font-display">📋 抽奖记录</h1>
        <p class="page-subtitle font-body">查看所有中奖记录和详细信息</p>
      </div>
      <div class="header-right">
        <a-space>
          <a-input-search
            v-model:value="search"
            placeholder="搜索手机号或姓名"
            style="width: 300px"
            @search="handleSearch"
            allow-clear
            class="neon-input"
          />
          <a-button @click="exportData" class="export-btn neon-button">
            <ExportOutlined /> 导出数据
          </a-button>
        </a-space>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-cards">
      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
          🎊
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.total_records }}</div>
          <div class="stat-label font-body">总记录数</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #52c41a 0%, #73d13d 100%);">
          🎁
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.today_records }}</div>
          <div class="stat-label font-body">今日中奖</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #faad14 0%, #ffc53d 100%);">
          🏢
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.companies_count }}</div>
          <div class="stat-label font-body">参与公司</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #1890ff 0%, #40a9ff 100%);">
          📈
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.draw_rate }}%</div>
          <div class="stat-label font-body">中奖率</div>
        </div>
      </div>
    </div>

    <!-- 表格 -->
    <a-table
      :columns="columns"
      :data-source="records"
      :pagination="pagination"
      @change="handleTableChange"
      row-key="id"
      class="records-table"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'company'">
          <a-tag v-if="record.company" :color="record.company.theme_color || 'blue'">
            {{ record.company.name }}
          </a-tag>
          <span v-else class="no-data">-</span>
        </template>
        <template v-else-if="column.key === 'user'">
          <div class="user-cell">
            <div class="user-avatar">
              {{ (record.user?.name || '未')[0] }}
            </div>
            <div class="user-info">
              <div class="user-name font-body">{{ record.user?.name || '未设置' }}</div>
              <div class="user-phone">{{ record.user?.phone || '-' }}</div>
            </div>
          </div>
        </template>
        <template v-else-if="column.key === 'prize'">
          <div class="prize-cell">
            <a-tag :color="getPrizeColor(record.level?.name)" class="prize-level-tag">
              {{ getPrizeIcon(record.level?.name) }} {{ record.level?.name }}
            </a-tag>
            <div class="prize-name">{{ record.prize?.name || '-' }}</div>
          </div>
        </template>
        <template v-else-if="column.key === 'time'">
          <div class="time-cell">
            <div class="time-date">{{ formatDate(record.created_at) }}</div>
          </div>
        </template>
      </template>
    </a-table>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { message } from 'ant-design-vue'
import { ExportOutlined } from '@ant-design/icons-vue'
import request from '../../utils/request'

const records = ref([])
const search = ref('')
const stats = ref({
  total_records: 0,
  today_records: 0,
  companies_count: 0,
  draw_rate: 0
})

const pagination = ref({
  current: 1,
  pageSize: 20,
  total: 0,
  showSizeChanger: true,
  showTotal: (total) => `共 ${total} 条`
})

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
  { title: '所属公司', key: 'company', width: 150 },
  { title: '中奖用户', key: 'user', width: 250 },
  { title: '中奖奖品', key: 'prize', width: 200 },
  { title: '抽奖时间', key: 'time', width: 200 }
]

const fetchRecords = async () => {
  try {
    const data = await request.get('/admin/draw-records', {
      params: {
        page: pagination.value.current,
        page_size: pagination.value.pageSize,
        search: search.value
      }
    })
    records.value = data.data || data
    pagination.value.total = data.total || data.length

    // 计算统计数据
    stats.value = {
      total_records: data.total || data.length || 0,
      today_records: records.value.filter(r => {
        const recordDate = new Date(r.created_at).toDateString()
        const today = new Date().toDateString()
        return recordDate === today
      }).length,
      companies_count: [...new Set(records.value.map(r => r.company_id))].length,
      draw_rate: 0 // 需要从后端获取总用户数来计算
    }
  } catch (error) {
    message.error('获取记录失败')
  }
}

const handleSearch = () => {
  pagination.value.current = 1
  fetchRecords()
}

const handleTableChange = (pag) => {
  pagination.value.current = pag.current
  pagination.value.pageSize = pag.pageSize
  fetchRecords()
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getPrizeIcon = (levelName) => {
  if (!levelName) return '🎁'
  if (levelName.includes('一等')) return '🥇'
  if (levelName.includes('二等')) return '🥈'
  if (levelName.includes('三等')) return '🥉'
  if (levelName.includes('参与')) return '🎁'
  return '🏆'
}

const getPrizeColor = (levelName) => {
  if (!levelName) return 'default'
  if (levelName.includes('一等')) return 'red'
  if (levelName.includes('二等')) return 'orange'
  if (levelName.includes('三等')) return 'green'
  if (levelName.includes('参与')) return 'blue'
  return 'purple'
}

const exportData = () => {
  message.info('导出功能开发中...')
}

onMounted(() => {
  fetchRecords()
})
</script>

<style scoped>
.drawrecords-page {
  padding: var(--spacing-lg);
}

/* 页面头部 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--spacing-xl);
  flex-wrap: wrap;
  gap: var(--spacing-md);
}

.header-left {
  flex: 1;
}

.page-title {
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-bold);
  margin: 0 0 var(--spacing-xs) 0;
  color: var(--text-primary);
}

.page-subtitle {
  font-size: var(--font-size-base);
  color: var(--text-primary);
  margin: 0;
}

.header-right {
  flex-shrink: 0;
}

.export-btn {
  height: 40px;
  font-weight: var(--font-weight-semibold);
}

/* 统计卡片 */
.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: var(--spacing-lg);
  margin-bottom: var(--spacing-xl);
}

.stat-card {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-lg);
  background: rgba(26, 26, 36, 0.6);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-1);
  transition: all var(--transition-bounce);
  animation: fadeInUp 0.6s ease-out;
}

.stat-card:hover {
  border-color: var(--neon-cyan);
  box-shadow: var(--glow-cyan), var(--shadow-2);
  transform: translateY(-2px);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  color: white;
  flex-shrink: 0;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-bold);
  line-height: 1;
  margin-bottom: var(--spacing-xs);
}

.stat-label {
  font-size: var(--font-size-sm);
  color: var(--text-primary);
}

/* 表格样式 */
.records-table {
  background: rgba(26, 26, 36, 0.6);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-xl);
  overflow: hidden;
  animation: fadeInUp 0.6s ease-out 0.1s both;
  transition: all var(--transition-base);
}

.records-table:hover {
  border-color: var(--neon-cyan);
  box-shadow: var(--glow-cyan);
}

.user-cell {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-full);
  background: var(--primary-gradient);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-bold);
  flex-shrink: 0;
}

.user-info {
  flex: 1;
}

.user-name {
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin-bottom: 2px;
}

.user-phone {
  font-size: var(--font-size-sm);
  color: var(--text-primary);
  font-family: var(--font-mono);
}

.prize-cell {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.prize-level-tag {
  align-self: flex-start;
  font-weight: var(--font-weight-semibold);
}

.prize-name {
  font-size: var(--font-size-sm);
  color: var(--text-primary);
}

.time-cell {
  display: flex;
  flex-direction: column;
}

.time-date {
  font-size: var(--font-size-sm);
  color: var(--text-primary);
}

.no-data {
  color: var(--text-tertiary);
}

/* 响应式 */
@media (max-width: 768px) {
  .drawrecords-page {
    padding: var(--spacing-md);
  }

  .page-title {
    font-size: var(--font-size-2xl);
  }

  .stats-cards {
    grid-template-columns: repeat(2, 1fr);
  }

  .page-header {
    flex-direction: column;
    align-items: stretch;
  }

  .header-right {
    width: 100%;
  }

  .header-right .ant-space {
    width: 100%;
    flex-direction: column;
  }

  .header-right .ant-input-search,
  .header-right .ant-btn {
    width: 100% !important;
  }
}

/* 输入框统一样式 */
.neon-input :deep(.ant-input) {
  background: rgba(255, 255, 255, 0.95) !important;
  border: 1px solid rgba(217, 217, 217, 0.8);
  border-radius: var(--radius-lg);
  color: #1a1a1a !important;
  transition: all var(--transition-base);
}

.neon-input :deep(.ant-input::placeholder) {
  color: #8c8c8c !important;
}

.neon-input :deep(.ant-input:focus),
.neon-input :deep(.ant-input-focused) {
  border-color: var(--neon-cyan) !important;
  box-shadow: 0 0 0 2px rgba(0, 255, 245, 0.2);
  background: rgba(255, 255, 255, 1) !important;
}
</style>
