<template>
  <div class="home-container">
    <!-- 背景效果 -->
    <div class="neon-background">
      <div class="gradient-orb orb-1"></div>
      <div class="gradient-orb orb-2"></div>
      <div class="gradient-orb orb-3"></div>
      <div class="particles" ref="particles"></div>
      <div class="grid-overlay"></div>
    </div>

    <div class="home-content">
      <!-- 顶部装饰线 -->
      <div class="top-decoration">
        <div class="deco-line line-left"></div>
        <div class="deco-dot"></div>
        <div class="deco-line line-right"></div>
      </div>

      <!-- 标题区域 -->
      <div class="title-section">
        <div class="title-wrapper">
          <h1 class="main-title font-display">
            <span class="title-line">NEON</span>
            <span class="title-line title-accent">LOTTERY</span>
          </h1>
          <div class="title-underline"></div>
        </div>
        <p class="subtitle font-body">🎰 霓虹嘉年华 · 幸运大抽奖 🎰</p>
        <p class="description font-body">选择您的公司，开启幸运之旅</p>
      </div>

      <!-- 加载状态 -->
      <div v-if="loading" class="loading-state">
        <div class="neon-spinner"></div>
        <p class="font-body">正在加载...</p>
      </div>

      <!-- 统计数据 -->
      <div v-if="!loading && stats" class="stats-section">
        <div class="stats-container">
          <div
            v-for="(stat, index) in statsCards"
            :key="index"
            class="stat-card glass"
            :style="{ animationDelay: `${index * 0.1}s` }"
          >
            <div class="stat-icon" :class="stat.iconClass">{{ stat.icon }}</div>
            <div class="stat-info">
              <div class="stat-value font-display">{{ stat.value }}</div>
              <div class="stat-label font-body">{{ stat.label }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 公司卡片 -->
      <div v-if="!loading && companies.length > 0" class="companies-section">
        <div class="section-header">
          <h2 class="section-title font-display">参与公司</h2>
          <div class="title-line-accent"></div>
        </div>
        <div class="companies-grid">
          <router-link
            v-for="(company, index) in companies"
            :key="company.id"
            :to="`/lottery?company=${company.code}`"
            class="company-card glass neon-border-flow"
            :style="{ animationDelay: `${index * 0.1}s` }"
          >
            <div class="card-inner">
              <div class="company-icon">{{ getIconEmoji(company.id) }}</div>
              <h3 class="company-name font-display">{{ company.name }}</h3>
              <p class="company-code font-mono">@{{ company.code }}</p>
              <div class="company-stats">
                <div class="stat-item">
                  <span class="stat-value font-display">{{ company.total_users || 0 }}</span>
                  <span class="stat-label">用户</span>
                </div>
                <div class="stat-divider"></div>
                <div class="stat-item">
                  <span class="stat-value font-display">{{ company.drawn_count || 0 }}</span>
                  <span class="stat-label">中奖</span>
                </div>
              </div>
              <div class="company-action">
                <span class="action-text font-display">进入抽奖</span>
                <span class="action-arrow">→</span>
              </div>
            </div>
            <div class="card-glow"></div>
          </router-link>
        </div>
      </div>

      <!-- 最近中奖 -->
      <div v-if="!loading && recentWinners.length > 0" class="recent-winners-section">
        <div class="section-header">
          <h2 class="section-title font-display">最近中奖</h2>
          <div class="title-line-accent"></div>
        </div>
        <div class="winners-list glass">
          <div
            v-for="(winner, index) in recentWinners"
            :key="winner.id"
            class="winner-item"
            :style="{ animationDelay: `${index * 0.08}s` }"
          >
            <div class="winner-avatar font-display">{{ (winner.user?.name || '未知')[0] }}</div>
            <div class="winner-details">
              <span class="winner-name font-display">{{ winner.user?.name || '未知' }}</span>
              <span class="winner-prize">{{ winner.level?.name || '' }}</span>
            </div>
            <div class="winner-company font-mono">{{ winner.company?.name || '' }}</div>
          </div>
        </div>
      </div>

      <!-- 统一登录入口 -->
      <div class="unified-access">
        <!-- 参与抽奖（统一登录入口） -->
        <router-link to="/lottery" class="access-link glass lottery-cta">
          <div class="access-icon">🎲</div>
          <div class="access-text">
            <h3 class="font-display">参与抽奖</h3>
            <p class="font-body">统一登录 · 用户和管理员都可参与</p>
          </div>
          <div class="access-arrow font-display">→</div>
          <div class="link-glow"></div>
          <div class="pulse-ring"></div>
        </router-link>

        <!-- 管理后台（仅已登录管理员可见） -->
        <router-link v-if="isLoggedIn" to="/admin/dashboard" class="access-link glass admin-cta">
          <div class="access-icon">🎛️</div>
          <div class="access-text">
            <h3 class="font-display">管理后台</h3>
            <p class="font-body">系统管理控制台</p>
          </div>
          <div class="access-arrow font-display">→</div>
          <div class="link-glow"></div>
        </router-link>
      </div>

      <!-- 底部装饰线 -->
      <div class="bottom-decoration">
        <div class="deco-line line-left"></div>
        <div class="deco-dot"></div>
        <div class="deco-line line-right"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import request from '../utils/request'

const router = useRouter()
const loading = ref(false)
const companies = ref([])
const stats = ref(null)
const recentWinners = ref([])
const particles = ref(null)
const isLoggedIn = ref(false)

// 检查是否已登录
const checkLoginStatus = () => {
  isLoggedIn.value = !!localStorage.getItem('admin_token')
}

// 监听路由变化，重新检查登录状态
router.afterEach(() => {
  checkLoginStatus()
})

// 组件挂载时检查登录状态
onMounted(() => {
  checkLoginStatus()
})

// 获取统计数据
const statsCards = computed(() => {
  if (!stats.value) return []
  return [
    {
      icon: '👥',
      iconClass: 'icon-cyan',
      value: stats.value.total_users || 0,
      label: '总用户数'
    },
    {
      icon: '🏆',
      iconClass: 'icon-yellow',
      value: stats.value.drawn_users || 0,
      label: '中奖人数'
    },
    {
      icon: '🏢',
      iconClass: 'icon-magenta',
      value: companies.value.length,
      label: '活跃公司'
    }
  ]
})

// 获取公司列表
const fetchCompanies = async () => {
  try {
    // 检查是否有登录的管理员
    const adminUserStr = localStorage.getItem('admin_user')
    const adminToken = localStorage.getItem('admin_token')

    if (adminUserStr && adminToken) {
      // 已登录，根据角色获取公司列表
      const adminUser = JSON.parse(adminUserStr)

      if (adminUser.is_super_admin) {
        // 超级管理员：获取所有激活的公司
        const data = await request.get('/admin/companies', {
          headers: {
            'Authorization': `Bearer ${adminToken}`
          }
        })
        // 只显示激活的公司，去重
        const uniqueCompanies = data.filter(c => c.is_active)
        companies.value = Array.from(new Map(uniqueCompanies.map(c => [c.id, c])).values())
      } else {
        // 普通管理员：只显示自己的公司
        if (adminUser.company) {
          companies.value = adminUser.company.is_active ? [adminUser.company] : []
        } else {
          companies.value = []
        }
      }
    } else {
      // 未登录：不显示公司列表
      companies.value = []
    }
  } catch (error) {
    // 出错时也不显示默认数据
    companies.value = []
  }
}

// 获取统计数据
const fetchStats = async () => {
  try {
    // 检查是否有登录的管理员
    const adminUserStr = localStorage.getItem('admin_user')
    const adminToken = localStorage.getItem('admin_token')

    if (adminUserStr && adminToken) {
      // 管理员已登录，使用管理员API
      const adminUser = JSON.parse(adminUserStr)

      const data = await request.get('/admin/stats', {
        params: { company_id: adminUser.company?.id }
      })
      stats.value = {
        total_users: data.total_users || 0,
        drawn_users: data.drawn_users || 0,
        total_records: data.total_records || 0
      }
    } else {
      // 未登录，使用默认值，不调用API
      stats.value = {
        total_users: 0,
        drawn_users: 0,
        total_records: 0
      }
    }
  } catch (error) {
    // 使用默认值，不影响页面显示
    stats.value = {
      total_users: 0,
      drawn_users: 0,
      total_records: 0
    }
  }
}

// 获取最近中奖记录
const fetchRecentWinners = async () => {
  try {
    // 检查是否有登录的管理员
    const adminUserStr = localStorage.getItem('admin_user')
    const adminToken = localStorage.getItem('admin_token')

    if (!adminUserStr || !adminToken) {
      // 未登录，不显示中奖记录
      recentWinners.value = []
      return
    }

    const adminUser = JSON.parse(adminUserStr)
    let companyID = null

    if (adminUser.company && !adminUser.is_super_admin) {
      // 普通管理员，使用自己公司的ID
      companyID = adminUser.company.id
    }

    // 管理员已登录，使用管理员API
    const data = await request.get('/admin/draw-records', {
      params: {
        company_id: companyID,
        page: 1,
        page_size: 5
      },
      headers: {
        'Authorization': `Bearer ${adminToken}`
      }
    })
    recentWinners.value = data.data || []
  } catch (error) {
    recentWinners.value = []
  }
}

// 获取公司图标
const getIconEmoji = (companyId) => {
  const icons = ['🎪', '🎭', '🎰', '🎲', '🎯', '🎱']
  return icons[companyId % icons.length] || '🎪'
}

// 创建粒子效果
const createParticles = () => {
  const container = particles.value
  if (!container) return

  for (let i = 0; i < 30; i++) {
    const particle = document.createElement('div')
    particle.className = 'particle'
    particle.style.left = `${Math.random() * 100}%`
    particle.style.top = `${Math.random() * 100}%`
    particle.style.animationDelay = `${Math.random() * 5}s`
    particle.style.animationDuration = `${5 + Math.random() * 10}s`
    container.appendChild(particle)
  }
}

onMounted(async () => {
  loading.value = true
  try {
    await Promise.all([
      fetchCompanies(),
      fetchStats(),
      fetchRecentWinners()
    ])
    createParticles()
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
/* ============================================
   主容器 - 深色霓虹背景
   ============================================ */

.home-container {
  min-height: 100vh;
  background: var(--bg-primary);
  position: relative;
  overflow-x: hidden;
}

.neon-background {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
  overflow: hidden;
}

/* 渐变光球 */
.gradient-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.3;
  animation: orbFloat 20s ease-in-out infinite;
}

.orb-1 {
  width: 600px;
  height: 600px;
  background: var(--neon-cyan);
  top: -200px;
  left: -200px;
  animation-delay: 0s;
}

.orb-2 {
  width: 500px;
  height: 500px;
  background: var(--neon-magenta);
  bottom: -150px;
  right: -150px;
  animation-delay: 7s;
}

.orb-3 {
  width: 400px;
  height: 400px;
  background: var(--neon-purple);
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation-delay: 14s;
}

@keyframes orbFloat {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  33% {
    transform: translate(50px, -50px) scale(1.1);
  }
  66% {
    transform: translate(-30px, 30px) scale(0.9);
  }
}

/* 粒子效果 */
.particles {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.particle {
  position: absolute;
  width: 4px;
  height: 4px;
  background: var(--neon-cyan);
  border-radius: 50%;
  opacity: 0.6;
  animation: particleFloat 10s ease-in-out infinite;
  box-shadow: 0 0 10px var(--neon-cyan);
}

@keyframes particleFloat {
  0%, 100% {
    transform: translateY(0) translateX(0);
    opacity: 0;
  }
  10% {
    opacity: 0.6;
  }
  90% {
    opacity: 0.6;
  }
  100% {
    transform: translateY(-100vh) translateX(50px);
    opacity: 0;
  }
}

/* 网格覆盖 */
.grid-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image:
    linear-gradient(rgba(0, 255, 245, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(0, 255, 245, 0.03) 1px, transparent 1px);
  background-size: 50px 50px;
  pointer-events: none;
}

/* ============================================
   主内容区域
   ============================================ */

.home-content {
  position: relative;
  z-index: 1;
  max-width: 1400px;
  margin: 0 auto;
  padding: var(--spacing-xl) var(--spacing-lg);
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* ============================================
   装饰元素
   ============================================ */

.top-decoration,
.bottom-decoration {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-md);
  margin: var(--spacing-xl) 0;
}

.deco-line {
  flex: 1;
  height: 1px;
  background: linear-gradient(90deg,
    transparent,
    var(--neon-cyan),
    var(--neon-magenta),
    transparent
  );
  position: relative;
  overflow: hidden;
}

.deco-line::after {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg,
    transparent,
    rgba(255, 255, 255, 0.8),
    transparent
  );
  animation: shimmer 3s ease-in-out infinite;
}

.line-left {
  animation-delay: 0s;
}

.line-right {
  animation-delay: 1.5s;
}

@keyframes shimmer {
  0% {
    left: -100%;
  }
  100% {
    left: 100%;
  }
}

.deco-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: var(--neon-cyan);
  box-shadow: 0 0 20px var(--neon-cyan),
              0 0 40px var(--neon-cyan);
  animation: dot-pulse 2s ease-in-out infinite;
}

@keyframes dot-pulse {
  0%, 100% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.3);
    opacity: 0.8;
  }
}

/* ============================================
   滚动提示
   ============================================ */

/* ============================================
   标题区域
   ============================================ */

.title-section {
  text-align: center;
  margin-bottom: var(--spacing-3xl);
  animation: fadeInDown 1s ease-out;
  position: relative;
}

.title-wrapper {
  position: relative;
  display: inline-block;
  perspective: 1000px;
}

.main-title {
  font-size: var(--font-size-6xl);
  font-weight: var(--font-weight-black);
  line-height: 1.1;
  margin-bottom: var(--spacing-md);
  letter-spacing: 8px;
  text-transform: uppercase;
  display: inline-block;
  animation: title-float 6s ease-in-out infinite;
}

@keyframes title-float {
  0%, 100% {
    transform: translateY(0) rotateX(0);
  }
  50% {
    transform: translateY(-10px) rotateX(5deg);
  }
}

.title-line {
  display: block;
  color: var(--text-primary);
  text-shadow: 0 0 30px var(--neon-cyan), 0 0 60px var(--neon-cyan);
  animation: textGlow 3s ease-in-out infinite;
}

.title-accent {
  color: var(--neon-magenta);
  text-shadow: 0 0 30px var(--neon-magenta), 0 0 60px var(--neon-magenta);
  animation-delay: 1.5s;
}

.title-underline {
  width: 0;
  height: 4px;
  background: linear-gradient(90deg, var(--neon-cyan), var(--neon-magenta));
  margin: var(--spacing-md) auto;
  border-radius: var(--radius-full);
  animation: expand-line 1s ease-out 0.5s forwards;
  box-shadow: 0 0 20px var(--neon-cyan);
}

@keyframes expand-line {
  to {
    width: 200px;
  }
}

.subtitle {
  font-size: var(--font-size-xl);
  color: var(--neon-cyan);
  margin-bottom: var(--spacing-sm);
  letter-spacing: 4px;
  text-shadow: 0 0 20px var(--neon-cyan);
  animation: fadeIn 1.5s ease-out 0.3s both;
}

.description {
  font-size: var(--font-size-base);
  color: var(--text-secondary);
  letter-spacing: 2px;
  animation: fadeIn 1.5s ease-out 0.5s both;
}

/* 快速操作按钮 */
.quick-actions {
  display: flex;
  gap: var(--spacing-lg);
  justify-content: center;
  margin-top: var(--spacing-xl);
  animation: fadeInUp 1s ease-out 0.7s both;
}

.quick-btn {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-md) var(--spacing-xl);
  border-radius: var(--radius-lg);
  text-decoration: none;
  color: inherit;
  border: 2px solid var(--neon-cyan);
  transition: all var(--transition-base);
  position: relative;
  overflow: hidden;
  background: rgba(0, 255, 245, 0.05);
}

.quick-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg,
    transparent,
    rgba(0, 255, 245, 0.3),
    transparent
  );
  transition: left var(--transition-base);
}

.quick-btn:hover::before {
  left: 100%;
}

.quick-btn:hover {
  transform: translateY(-4px);
  box-shadow: 0 10px 30px rgba(0, 255, 245, 0.3),
              0 0 30px rgba(0, 255, 245, 0.2);
  border-color: var(--neon-magenta);
}

.btn-icon {
  font-size: var(--font-size-2xl);
  filter: drop-shadow(0 0 10px var(--neon-cyan));
  transition: transform var(--transition-base);
}

.quick-btn:hover .btn-icon {
  transform: scale(1.2) rotate(10deg);
  filter: drop-shadow(0 0 10px var(--neon-magenta));
}

.btn-text {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  letter-spacing: 2px;
  text-transform: uppercase;
  color: var(--text-primary);
}

.btn-arrow {
  font-size: var(--font-size-xl);
  color: var(--neon-cyan);
  transition: transform var(--transition-base);
}

.quick-btn:hover .btn-arrow {
  transform: translateX(8px);
  color: var(--neon-magenta);
}

.login-btn {
  border-color: var(--neon-magenta);
  background: rgba(255, 0, 255, 0.05);
}

.login-btn .btn-icon {
  filter: drop-shadow(0 0 10px var(--neon-magenta));
}

.login-btn:hover {
  border-color: var(--neon-cyan);
}

/* ============================================
   统计卡片
   ============================================ */

.stats-section {
  margin-bottom: var(--spacing-3xl);
}

.stats-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: var(--spacing-lg);
}

.stat-card {
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
  padding: var(--spacing-lg);
  border-radius: var(--radius-xl);
  border: 1px solid var(--border-color);
  animation: fadeInUp 0.6s ease-out backwards;
  transition: all var(--transition-base);
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-3);
  border-color: var(--neon-cyan);
}

.stat-icon {
  font-size: var(--font-size-4xl);
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-lg);
  background: rgba(255, 255, 255, 0.05);
}

.icon-cyan {
  color: var(--neon-cyan);
  text-shadow: 0 0 20px var(--neon-cyan);
}

.icon-yellow {
  color: var(--neon-yellow);
  text-shadow: 0 0 20px var(--neon-yellow);
}

.icon-magenta {
  color: var(--neon-magenta);
  text-shadow: 0 0 20px var(--neon-magenta);
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: var(--font-size-3xl);
  color: var(--text-primary);
  line-height: 1;
  margin-bottom: var(--spacing-xs);
}

.stat-label {
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 1px;
}

/* ============================================
   公司展示区
   ============================================ */

.companies-section {
  margin-bottom: var(--spacing-3xl);
}

.section-header {
  text-align: center;
  margin-bottom: var(--spacing-xl);
  animation: fadeIn 0.8s ease-out;
}

.section-title {
  font-size: var(--font-size-3xl);
  color: var(--text-primary);
  margin-bottom: var(--spacing-sm);
  letter-spacing: 4px;
  text-transform: uppercase;
  text-shadow: 0 0 20px var(--neon-cyan);
}

.title-line-accent {
  width: 100px;
  height: 3px;
  background: var(--primary-gradient);
  margin: 0 auto;
  border-radius: var(--radius-full);
  box-shadow: 0 0 10px var(--neon-cyan);
}

.companies-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: var(--spacing-lg);
}

.company-card {
  position: relative;
  border-radius: var(--radius-xl);
  padding: var(--spacing-lg);
  text-decoration: none;
  color: inherit;
  display: block;
  overflow: hidden;
  animation: fadeInUp 0.6s ease-out backwards;
  transition: all var(--transition-base);
  border: 2px solid transparent;
}

.company-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg,
    transparent,
    var(--neon-cyan),
    var(--neon-magenta),
    transparent
  );
  opacity: 0;
  transition: opacity var(--transition-base);
}

.company-card:hover {
  transform: translateY(-8px);
  box-shadow: var(--shadow-3), var(--glow-cyan);
  border-color: var(--neon-cyan);
}

.company-card:hover::before {
  opacity: 1;
}

.card-inner {
  position: relative;
  z-index: 1;
}

.company-icon {
  font-size: var(--font-size-5xl);
  text-align: center;
  margin-bottom: var(--spacing-md);
  filter: drop-shadow(0 0 20px var(--neon-magenta));
  animation: float 3s ease-in-out infinite;
}

.company-name {
  font-size: var(--font-size-2xl);
  text-align: center;
  color: var(--text-primary);
  margin-bottom: var(--spacing-xs);
  text-transform: uppercase;
  letter-spacing: 2px;
}

.company-code {
  font-size: var(--font-size-sm);
  text-align: center;
  color: var(--neon-cyan);
  margin-bottom: var(--spacing-md);
  text-shadow: 0 0 10px var(--neon-cyan);
}

.company-stats {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: var(--spacing-lg);
  margin-bottom: var(--spacing-md);
}

.stat-item {
  text-align: center;
}

.stat-item .stat-value {
  display: block;
  font-size: var(--font-size-xl);
  color: var(--neon-yellow);
  text-shadow: 0 0 10px var(--neon-yellow);
}

.stat-item .stat-label {
  font-size: var(--font-size-xs);
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 1px;
}

.stat-divider {
  width: 1px;
  height: 30px;
  background: var(--border-color);
}

.company-action {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: var(--spacing-sm);
  padding-top: var(--spacing-md);
  border-top: 1px solid var(--border-color);
}

.action-text {
  font-size: var(--font-size-base);
  color: var(--neon-cyan);
  text-transform: uppercase;
  letter-spacing: 2px;
  transition: all var(--transition-fast);
}

.action-arrow {
  font-size: var(--font-size-xl);
  color: var(--neon-cyan);
  transition: transform var(--transition-base);
}

.company-card:hover .action-text {
  color: var(--neon-magenta);
  text-shadow: 0 0 10px var(--neon-magenta);
}

.company-card:hover .action-arrow {
  transform: translateX(8px);
  color: var(--neon-magenta);
}

.card-glow {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 0;
  height: 0;
  background: radial-gradient(circle, var(--neon-cyan) 0%, transparent 70%);
  opacity: 0;
  transition: all var(--transition-slow);
  pointer-events: none;
}

.company-card:hover .card-glow {
  width: 300px;
  height: 300px;
  opacity: 0.1;
}

/* ============================================
   最近中奖
   ============================================ */

.recent-winners-section {
  margin-bottom: var(--spacing-3xl);
}

.winners-list {
  border-radius: var(--radius-xl);
  padding: var(--spacing-lg);
  border: 1px solid var(--border-color);
}

.winner-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-md) 0;
  border-bottom: 1px solid var(--border-color);
  animation: slideInLeft 0.5s ease-out backwards;
  transition: all var(--transition-fast);
}

.winner-item:last-child {
  border-bottom: none;
}

.winner-item:hover {
  background: rgba(0, 255, 245, 0.05);
  padding-left: var(--spacing-md);
  padding-right: var(--spacing-md);
  margin-left: calc(-1 * var(--spacing-md));
  margin-right: calc(-1 * var(--spacing-md));
  border-radius: var(--radius-base);
}

.winner-avatar {
  width: 50px;
  height: 50px;
  border-radius: var(--radius-full);
  background: var(--primary-gradient);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-xl);
  color: var(--text-inverse);
  flex-shrink: 0;
  box-shadow: 0 0 20px rgba(0, 255, 245, 0.3);
}

.winner-details {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.winner-name {
  font-size: var(--font-size-lg);
  color: var(--text-primary);
  font-weight: var(--font-weight-semibold);
}

.winner-prize {
  font-size: var(--font-size-sm);
  color: var(--neon-yellow);
  padding: 4px 12px;
  background: rgba(255, 204, 0, 0.1);
  border-radius: var(--radius-full);
  align-self: flex-start;
  border: 1px solid rgba(255, 204, 0, 0.3);
  text-shadow: 0 0 10px var(--neon-yellow);
}

.winner-company {
  font-size: var(--font-size-xs);
  color: var(--text-tertiary);
  padding: 4px 8px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: var(--radius-base);
  flex-shrink: 0;
  border: 1px solid var(--border-color);
}

/* ============================================
   统一登录入口
   ============================================ */

.unified-access {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-lg);
  animation: fadeIn 1s ease-out;
  margin-bottom: var(--spacing-xl);
}

.access-link {
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
  padding: var(--spacing-lg) var(--spacing-xl);
  border-radius: var(--radius-xl);
  text-decoration: none;
  color: inherit;
  transition: all var(--transition-base);
  width: 100%;
  max-width: 600px;
  position: relative;
  overflow: hidden;
}

/* 抽奖入口样式 */
.lottery-cta {
  border: 2px solid var(--neon-cyan);
  background: linear-gradient(135deg,
    rgba(0, 255, 245, 0.15) 0%,
    rgba(0, 255, 245, 0.05) 100%
  );
  animation: pulse-border 2s ease-in-out infinite;
}

.lottery-cta .access-icon {
  filter: drop-shadow(0 0 20px var(--neon-cyan));
}

.lottery-cta:hover {
  border-color: var(--neon-cyan);
  box-shadow: 0 0 40px rgba(0, 255, 245, 0.4);
  transform: translateY(-4px);
}

/* 管理后台入口样式 */
.admin-cta {
  border: 2px solid var(--neon-magenta);
  background: linear-gradient(135deg,
    rgba(255, 0, 255, 0.15) 0%,
    rgba(255, 0, 255, 0.05) 100%
  );
}

.admin-cta .access-icon {
  filter: drop-shadow(0 0 20px var(--neon-magenta));
}

.admin-cta:hover {
  border-color: var(--neon-magenta);
  box-shadow: 0 0 40px rgba(255, 0, 255, 0.4);
  transform: translateY(-4px);
}

.access-icon {
  font-size: var(--font-size-4xl);
  filter: drop-shadow(0 0 20px var(--neon-cyan));
  transition: transform var(--transition-base);
  animation: float 3s ease-in-out infinite;
}

.access-link:hover .access-icon {
  transform: scale(1.1) rotate(5deg);
}

.access-text {
  flex: 1;
}

.access-text h3 {
  font-size: var(--font-size-xl);
  color: var(--text-primary);
  margin-bottom: var(--spacing-xs);
  text-transform: uppercase;
  letter-spacing: 2px;
}

.access-text p {
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
}

.access-arrow {
  font-size: var(--font-size-2xl);
  color: var(--neon-cyan);
  transition: transform var(--transition-base);
}

.lottery-cta .access-arrow {
  color: var(--neon-cyan);
}

.admin-cta .access-arrow {
  color: var(--neon-magenta);
}

.access-link:hover .access-arrow {
  transform: translateX(10px);
}

.link-glow {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 0;
  height: 0;
  background: radial-gradient(circle, var(--neon-cyan) 0%, transparent 70%);
  opacity: 0;
  transition: all var(--transition-slow);
  pointer-events: none;
}

.lottery-cta:hover .link-glow {
  width: 400px;
  height: 400px;
  opacity: 0.15;
  background: radial-gradient(circle, var(--neon-cyan) 0%, transparent 70%);
}

.admin-cta:hover .link-glow {
  width: 400px;
  height: 400px;
  opacity: 0.15;
  background: radial-gradient(circle, var(--neon-magenta) 0%, transparent 70%);
}

/* 脉冲环效果 */
.pulse-ring {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 100%;
  height: 100%;
  border: 2px solid var(--neon-cyan);
  border-radius: var(--radius-xl);
  opacity: 0;
  animation: pulse-ring 2s ease-out infinite;
  pointer-events: none;
}

@keyframes pulse-border {
  0%, 100% {
    box-shadow: 0 0 20px rgba(0, 255, 245, 0.3),
                0 0 40px rgba(0, 255, 245, 0.2);
  }
  50% {
    box-shadow: 0 0 30px rgba(0, 255, 245, 0.5),
                0 0 60px rgba(0, 255, 245, 0.3);
  }
}

@keyframes pulse-ring {
  0% {
    transform: translate(-50%, -50%) scale(1);
    opacity: 0.6;
  }
  100% {
    transform: translate(-50%, -50%) scale(1.1);
    opacity: 0;
  }
}

/* ============================================
   加载状态
   ============================================ */

.loading-state {
  text-align: center;
  padding: var(--spacing-3xl) var(--spacing-lg);
}

.neon-spinner {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  border: 3px solid transparent;
  border-top-color: var(--neon-cyan);
  border-right-color: var(--neon-magenta);
  animation: spin 1s linear infinite;
  margin: 0 auto var(--spacing-md);
  box-shadow: 0 0 30px var(--neon-cyan), 0 0 60px var(--neon-magenta);
}

.loading-state p {
  color: var(--text-secondary);
  font-size: var(--font-size-base);
}

/* ============================================
   响应式设计
   ============================================ */

@media (max-width: 768px) {
  .home-content {
    padding: var(--spacing-lg) var(--spacing-md);
  }

  .main-title {
    font-size: var(--font-size-4xl);
    letter-spacing: 4px;
  }

  .subtitle {
    font-size: var(--font-size-lg);
    letter-spacing: 2px;
  }

  .description {
    font-size: var(--font-size-sm);
  }

  .stats-container {
    grid-template-columns: 1fr;
  }

  .companies-grid {
    grid-template-columns: 1fr;
  }

  .winner-item {
    flex-wrap: wrap;
  }

  .winner-company {
    margin-left: auto;
  }

  .section-title {
    font-size: var(--font-size-2xl);
    letter-spacing: 2px;
  }

  .orb-1, .orb-2, .orb-3 {
    width: 300px;
    height: 300px;
  }

  .access-link {
    max-width: 100%;
  }
}

/* 减少动画偏好 */
@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }

  .particle {
    display: none;
  }
}
</style>
