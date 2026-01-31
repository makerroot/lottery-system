<template>
  <div class="dashboard-container">
    <!-- 页面头部 -->
    <div class="dashboard-header">
      <div class="header-content">
        <h1 class="dashboard-title font-display">📊 数据仪表盘</h1>
        <p class="dashboard-subtitle font-body">实时监控系统运营数据</p>
      </div>
      <div class="dashboard-actions">
        <a-button @click="refreshData" :loading="loading" class="neon-button" size="large">
          <template #icon><ReloadOutlined /></template>
          刷新数据
        </a-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-grid">
      <StatCard
        title="总用户数"
        :value="stats.total_users"
        icon="👥"
        icon-color="var(--neon-purple)"
        :trend="calculateTrend('total_users')"
        :show-trend="true"
        :trend-up="true"
        :loading="loading"
        variant="glass"
      />
      <StatCard
        title="已抽奖人数"
        :value="stats.drawn_users"
        icon="🏆"
        icon-color="var(--success-color)"
        :trend="calculateTrend('drawn_users')"
        :show-trend="true"
        :trend-up="true"
        :loading="loading"
        variant="glass"
      />
      <StatCard
        title="抽奖记录"
        :value="stats.total_records"
        icon="📊"
        icon-color="var(--warning-color)"
        :show-trend="false"
        :loading="loading"
        variant="glass"
      />
      <StatCard
        title="活跃公司"
        :value="companyCount"
        icon="🏢"
        icon-color="var(--info-color)"
        :show-trend="false"
        :loading="loading"
        variant="glass"
      />
    </div>

    <!-- 图表区域 -->
    <a-row :gutter="[24, 24]" class="charts-section">
      <!-- 抽奖进度环形图 -->
      <a-col :xs="24" :sm="24" :md="12" :lg="8">
        <div class="chart-card glass-card">
          <div class="card-header">
            <h3 class="card-title font-display">🎯 抽奖参与率</h3>
          </div>
          <div class="ring-chart-container">
            <ProgressRing
              :percent="participationRate"
              :size="180"
              :stroke="12"
              stroke-color="var(--neon-purple)"
            />
            <div class="ring-legend">
              <div class="legend-item">
                <span class="legend-dot"></span>
                <span class="legend-label font-body">已参与: {{ stats.drawn_users }}</span>
              </div>
              <div class="legend-item">
                <span class="legend-dot secondary"></span>
                <span class="legend-label font-body">未参与: {{ stats.total_users - stats.drawn_users }}</span>
              </div>
            </div>
          </div>
        </div>
      </a-col>

      <!-- 奖项库存进度条 -->
      <a-col :xs="24" :sm="24" :md="12" :lg="16">
        <div class="chart-card glass-card">
          <div class="card-header">
            <h3 class="card-title font-display">🏆 奖项库存状态</h3>
          </div>
          <div class="prize-levels-list">
            <div
              v-for="level in stats.levels"
              :key="level.id"
              class="level-item"
            >
              <div class="level-info">
                <div class="level-name font-body">
                  <span class="level-badge" :style="{ background: getLevelColor(level.name) }">
                    {{ level.name.charAt(0) }}
                  </span>
                  <span class="level-text">{{ level.name }}</span>
                </div>
                <div class="level-stock font-mono">
                  {{ level.used_stock }} / {{ level.total_stock }}
                </div>
              </div>
              <div class="level-progress">
                <a-progress
                  :percent="getStockPercent(level)"
                  :stroke-color="getProgressColor(level)"
                  :show-info="false"
                  :stroke-width="8"
                />
              </div>
              <div class="level-meta">
                <span class="level-probability font-mono">概率: {{ (level.probability * 100).toFixed(1) }}%</span>
                <a-tag :color="getStockPercent(level) >= 100 ? 'error' : 'success'" size="small">
                  {{ getStockPercent(level) >= 100 ? '已抽完' : '库存充足' }}
                </a-tag>
              </div>
            </div>
            <a-empty v-if="!stats.levels || stats.levels.length === 0" description="暂无奖项数据" />
          </div>
        </div>
      </a-col>
    </a-row>

    <!-- 公司排行榜和最近中奖 -->
    <a-row :gutter="[24, 24]" class="ranking-section">
      <!-- 公司排行榜 -->
      <a-col :xs="24" :lg="12">
        <div class="chart-card glass-card">
          <div class="card-header">
            <h3 class="card-title font-display">🏢 公司排行榜</h3>
          </div>
          <div class="ranking-list">
            <div
              v-for="(company, index) in topCompanies"
              :key="company.id"
              class="ranking-item"
              :class="{ 'ranking-top': index < 3 }"
            >
              <div class="ranking-rank">
                <span v-if="index === 0" class="rank-badge rank-gold">🥇</span>
                <span v-else-if="index === 1" class="rank-badge rank-silver">🥈</span>
                <span v-else-if="index === 2" class="rank-badge rank-bronze">🥉</span>
                <span v-else class="rank-number font-display">{{ index + 1 }}</span>
              </div>
              <div class="ranking-content">
                <div class="company-name font-body">{{ company.name }}</div>
                <div class="company-stats">
                  <span class="stat-item">👥 {{ company.total_users || 0 }}</span>
                  <span class="stat-item">🏆 {{ company.drawn_count || 0 }}</span>
                </div>
              </div>
              <div class="ranking-rate">
                <div class="rate-value font-display">{{ getDrawRate(company) }}%</div>
                <div class="rate-label font-body">中奖率</div>
              </div>
            </div>
            <a-empty v-if="topCompanies.length === 0" description="暂无公司数据" />
          </div>
        </div>
      </a-col>

      <!-- 最近中奖记录 -->
      <a-col :xs="24" :lg="12">
        <div class="chart-card glass-card">
          <div class="card-header">
            <h3 class="card-title font-display">🏆 最近中奖</h3>
          </div>
          <div class="winners-list">
            <div
              v-for="record in recentWinners"
              :key="record.id"
              class="winner-item"
            >
              <div class="winner-avatar font-display">
                {{ getAvatarText(record.user?.name) }}
              </div>
              <div class="winner-info">
                <div class="winner-name font-body">{{ record.user?.name || '未知' }}</div>
                <div class="winner-prize font-body">{{ record.level?.name }}</div>
              </div>
              <div class="winner-company">
                <a-tag :color="record.company?.theme_color || 'blue'" size="small">
                  {{ record.company?.name }}
                </a-tag>
              </div>
              <div class="winner-time font-mono">
                {{ formatTime(record.created_at) }}
              </div>
            </div>
            <a-empty v-if="recentWinners.length === 0" description="暂无中奖记录" />
          </div>
          <div class="winners-footer">
            <router-link to="/admin/dashboard/records" class="view-all-link">
              查看全部 →
            </router-link>
          </div>
        </div>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { message } from 'ant-design-vue'
import { ReloadOutlined } from '@ant-design/icons-vue'
import request from '../../utils/request'
import StatCard from '../../components/StatCard.vue'
import ProgressRing from '../../components/ProgressRing.vue'

const loading = ref(false)
const stats = ref({
  total_users: 0,
  drawn_users: 0,
  total_records: 0,
  levels: []
})

const companies = ref([])
const recentWinners = ref([])

// 上次数据（用于计算趋势）
const previousStats = ref({
  total_users: 0,
  drawn_users: 0
})

// 自动刷新定时器
let refreshTimer = null

// 计算公司数量
const companyCount = computed(() => companies.value.length)

// 计算参与率
const participationRate = computed(() => {
  if (stats.value.total_users === 0) return 0
  return Math.round((stats.value.drawn_users / stats.value.total_users) * 100)
})

// Top 5 公司
const topCompanies = computed(() => {
  return companies.value
    .sort((a, b) => (b.drawn_count || 0) - (a.drawn_count || 0))
    .slice(0, 5)
})

// 获取统计数据
const fetchStats = async () => {
  try {
    const data = await request.get('/admin/stats')
    // 保存上次数据
    previousStats.value = {
      total_users: stats.value.total_users,
      drawn_users: stats.value.drawn_users
    }
    stats.value = data
  } catch (error) {
  }
}

// 获取公司列表
const fetchCompanies = async () => {
  try {
    const data = await request.get('/admin/companies')
    companies.value = data.filter(c => c.is_active)
  } catch (error) {
  }
}

// 获取最近中奖记录
const fetchRecentWinners = async () => {
  try {
    const data = await request.get('/admin/draw-records', {
      params: { page: 1, page_size: 5 }
    })
    recentWinners.value = data.data || []
  } catch (error) {
  }
}

// 刷新所有数据
const refreshData = async () => {
  loading.value = true
  try {
    await Promise.all([
      fetchStats(),
      fetchCompanies(),
      fetchRecentWinners()
    ])
    message.success('数据刷新成功')
  } catch (error) {
    message.error('数据刷新失败')
  } finally {
    loading.value = false
  }
}

// 计算趋势
const calculateTrend = (key) => {
  const current = stats.value[key]
  const previous = previousStats.value[key]
  if (previous === 0) return '+0%'
  const change = ((current - previous) / previous * 100).toFixed(1)
  return (change > 0 ? '+' : '') + change + '%'
}

// 获取库存百分比
const getStockPercent = (level) => {
  if (level.total_stock === 0) return 0
  return Math.round((level.used_stock / level.total_stock) * 100)
}

// 获取进度条颜色
const getProgressColor = (level) => {
  const percent = getStockPercent(level)
  if (percent >= 100) return 'var(--error-color)'
  if (percent >= 80) return 'var(--warning-color)'
  return 'var(--success-color)'
}

// 获取等级颜色
const getLevelColor = (name) => {
  const colorMap = {
    '一等奖': 'var(--neon-magenta)',
    '二等奖': 'var(--neon-yellow)',
    '三等奖': 'var(--neon-cyan)',
    '参与奖': 'var(--neon-purple)'
  }
  return colorMap[name] || 'var(--neon-purple)'
}

// 获取头像文字
const getAvatarText = (name) => {
  if (!name) return '?'
  return name.charAt(0)
}

// 计算中奖率
const getDrawRate = (company) => {
  if (!company.total_users || company.total_users === 0) return 0
  return ((company.drawn_count || 0) / company.total_users * 100).toFixed(1)
}

// 格式化时间
const formatTime = (dateStr) => {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now - date
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  return `${days}天前`
}

onMounted(async () => {
  await refreshData()
  // 自动刷新（每30秒）
  refreshTimer = setInterval(refreshData, 30000)
})

onUnmounted(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
})
</script>

<style scoped>
/* ============================================
   Dashboard Container
   ============================================ */

.dashboard-container {
  padding: var(--spacing-xl);
  max-width: 1600px;
  margin: 0 auto;
}

/* ============================================
   Header Section
   ============================================ */

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-2xl);
  padding: var(--spacing-xl) var(--spacing-2xl);
  background: rgba(26, 26, 36, 0.4);
  backdrop-filter: blur(20px);
  border-radius: var(--radius-xl);
  border: 1px solid var(--border-color);
  flex-wrap: wrap;
  gap: var(--spacing-md);
}

.header-content {
  flex: 1;
  min-width: 0;
}

.dashboard-title {
  font-size: var(--font-size-4xl);
  font-weight: var(--font-weight-bold);
  margin: 0 0 var(--spacing-xs) 0;
  color: var(--text-primary);
  line-height: 1.2;
}

.dashboard-subtitle {
  font-size: var(--font-size-base);
  color: var(--text-primary);
  margin: 0;
}

.dashboard-actions {
  flex-shrink: 0;
}

/* ============================================
   Stats Grid
   ============================================ */

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: var(--spacing-lg);
  margin-bottom: var(--spacing-2xl);
}

/* ============================================
   Chart Cards
   ============================================ */

.charts-section,
.ranking-section {
  margin-bottom: var(--spacing-2xl);
}

.glass-card {
  height: 100%;
  background: rgba(26, 26, 36, 0.6);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-xl);
  padding: var(--spacing-xl);
  transition: all var(--transition-base);
  box-shadow: var(--shadow-2);
}

.glass-card:hover {
  border-color: var(--neon-cyan);
  box-shadow: var(--glow-cyan), var(--shadow-3);
  transform: translateY(-2px);
}

.card-header {
  margin-bottom: var(--spacing-lg);
  padding-bottom: var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
}

.card-title {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-bold);
  margin: 0;
  color: var(--text-primary);
}

/* ============================================
   Ring Chart
   ============================================ */

.ring-chart-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: var(--spacing-lg) 0;
}

.ring-legend {
  margin-top: var(--spacing-lg);
  display: flex;
  gap: var(--spacing-xl);
  flex-wrap: wrap;
  justify-content: center;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  font-size: var(--font-size-sm);
}

.legend-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: var(--neon-cyan);
  box-shadow: 0 0 10px var(--neon-cyan);
}

.legend-dot.secondary {
  background: var(--text-tertiary);
  box-shadow: none;
}

.legend-label {
  color: var(--text-primary);
}

/* ============================================
   Prize Levels List
   ============================================ */

.prize-levels-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.level-item {
  padding: var(--spacing-md);
  background: rgba(255, 255, 255, 0.03);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-color);
  transition: all var(--transition-base);
}

.level-item:hover {
  background: rgba(0, 255, 245, 0.05);
  border-color: var(--neon-cyan);
  box-shadow: 0 0 15px rgba(0, 255, 245, 0.2);
  transform: translateX(4px);
}

.level-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-md);
}

.level-name {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
}

.level-badge {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-base);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: var(--font-weight-bold);
  font-size: var(--font-size-base);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.level-text {
  font-size: var(--font-size-base);
}

.level-stock {
  font-size: var(--font-size-sm);
  color: var(--text-primary);
  font-weight: var(--font-weight-medium);
}

.level-progress {
  margin-bottom: var(--spacing-md);
}

.level-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.level-probability {
  font-size: var(--font-size-xs);
  color: var(--text-tertiary);
}

/* ============================================
   Ranking & Winners
   ============================================ */

.ranking-list,
.winners-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.ranking-item,
.winner-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-md);
  background: rgba(255, 255, 255, 0.03);
  border-radius: var(--radius-lg);
  border: 1px solid transparent;
  transition: all var(--transition-base);
}

.ranking-item:hover,
.winner-item:hover {
  background: rgba(0, 255, 245, 0.05);
  border-color: var(--neon-cyan);
  transform: translateX(4px);
}

.ranking-item.ranking-top {
  background: linear-gradient(135deg, rgba(255, 215, 0, 0.1) 0%, rgba(255, 215, 0, 0.05) 100%);
  border: 1px solid rgba(255, 215, 0, 0.2);
}

/* Ranking */
.ranking-rank {
  width: 48px;
  display: flex;
  justify-content: center;
  flex-shrink: 0;
}

.rank-badge {
  font-size: 28px;
  line-height: 1;
}

.rank-number {
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  color: var(--text-tertiary);
}

.ranking-content {
  flex: 1;
  min-width: 0;
}

.company-name {
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--spacing-xs);
  color: var(--text-primary);
  font-size: var(--font-size-base);
}

.company-stats {
  display: flex;
  gap: var(--spacing-lg);
  font-size: var(--font-size-sm);
  color: var(--text-primary);
}

.stat-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.ranking-rate {
  text-align: center;
  flex-shrink: 0;
}

.rate-value {
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  color: var(--neon-cyan);
  line-height: 1;
}

.rate-label {
  font-size: var(--font-size-xs);
  color: var(--text-tertiary);
  margin-top: 2px;
}

/* Winners */
.winner-avatar {
  width: 44px;
  height: 44px;
  border-radius: var(--radius-full);
  background: var(--primary-gradient);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: var(--font-weight-bold);
  flex-shrink: 0;
  font-size: var(--font-size-lg);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.winner-info {
  flex: 1;
  min-width: 0;
}

.winner-name {
  font-weight: var(--font-weight-semibold);
  margin-bottom: 2px;
  color: var(--text-primary);
  font-size: var(--font-size-sm);
}

.winner-prize {
  font-size: var(--font-size-sm);
  color: var(--text-primary);
}

.winner-company {
  flex-shrink: 0;
}

.winner-time {
  font-size: var(--font-size-xs);
  color: var(--text-tertiary);
  flex-shrink: 0;
}

.winners-footer {
  margin-top: var(--spacing-lg);
  padding-top: var(--spacing-md);
  border-top: 1px solid var(--border-color);
  text-align: center;
}

.view-all-link {
  color: var(--neon-cyan);
  text-decoration: none;
  font-weight: var(--font-weight-medium);
  font-size: var(--font-size-sm);
  transition: all var(--transition-base);
  display: inline-flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.view-all-link:hover {
  color: var(--neon-magenta);
  text-shadow: 0 0 10px var(--neon-magenta);
}

/* ============================================
   Responsive
   ============================================ */

@media (max-width: 1200px) {
  .dashboard-container {
    padding: var(--spacing-lg);
  }

  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .dashboard-container {
    padding: var(--spacing-md);
  }

  .dashboard-header {
    flex-direction: column;
    align-items: stretch;
    padding: var(--spacing-lg);
  }

  .dashboard-title {
    font-size: var(--font-size-3xl);
  }

  .stats-grid {
    grid-template-columns: 1fr;
  }

  .ranking-item {
    flex-wrap: wrap;
  }

  .ranking-rate {
    width: 100%;
    margin-top: var(--spacing-sm);
    text-align: left;
    display: flex;
    align-items: center;
    gap: var(--spacing-md);
  }

  .winner-item {
    flex-wrap: wrap;
  }

  .winner-time {
    width: 100%;
    margin-top: var(--spacing-xs);
  }
}

@media (max-width: 480px) {
  .dashboard-title {
    font-size: var(--font-size-2xl);
  }

  .glass-card {
    padding: var(--spacing-md);
  }

  .ring-legend {
    flex-direction: column;
    gap: var(--spacing-sm);
  }
}
</style>
