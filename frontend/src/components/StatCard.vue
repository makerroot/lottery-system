<template>
  <div class="stat-card" :class="cardClass">
    <div class="stat-icon-wrapper" :style="{ background: iconBg }">
      <div class="stat-icon">{{ icon }}</div>
    </div>
    <div class="stat-content">
      <div class="stat-value font-display">{{ formattedValue }}</div>
      <div class="stat-label font-body">{{ title }}</div>
      <div v-if="showTrend" class="stat-trend font-body" :class="trendClass">
        <span class="trend-icon">{{ trendIcon }}</span>
        <span class="trend-value">{{ trend }}</span>
      </div>
    </div>
    <div v-if="loading" class="stat-skeleton"></div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  title: {
    type: String,
    required: true
  },
  value: {
    type: [Number, String],
    default: 0
  },
  icon: {
    type: String,
    default: '📊'
  },
  iconColor: {
    type: String,
    default: '#667eea'
  },
  trend: {
    type: String,
    default: ''
  },
  trendUp: {
    type: Boolean,
    default: true
  },
  showTrend: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  },
  variant: {
    type: String,
    default: 'default', // default, gradient, glass
    validator: (value) => ['default', 'gradient', 'glass'].includes(value)
  }
})

const cardClass = computed(() => `stat-card--${props.variant}`)

const iconBg = computed(() => {
  // 如果是CSS变量，直接使用
  if (props.iconColor.startsWith('var(')) {
    return props.variant === 'gradient' ? props.iconColor : props.iconColor.replace(')', ', 0.1)').replace('var(', 'var(')
  }
  // 如果是普通颜色，添加透明度
  const opacity = props.variant === 'gradient' ? '1' : '0.1'
  return props.iconColor + opacity
})

const formattedValue = computed(() => {
  if (typeof props.value === 'number') {
    return props.value.toLocaleString()
  }
  return props.value
})

const trendClass = computed(() => ({
  'stat-trend--up': props.trendUp,
  'stat-trend--down': !props.trendUp
}))

const trendIcon = computed(() => props.trendUp ? '↗' : '↘')
</script>

<style scoped>
.stat-card {
  position: relative;
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
  padding: var(--spacing-xl);
  background: rgba(26, 26, 36, 0.6);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-2);
  transition: all var(--transition-base);
  overflow: hidden;
  min-height: 120px;
}

.stat-card:hover {
  border-color: var(--neon-cyan);
  box-shadow: var(--glow-cyan), var(--shadow-3);
  transform: translateY(-4px);
}

.stat-card--gradient {
  background: var(--primary-gradient);
  color: var(--text-inverse);
}

.stat-card--gradient .stat-label {
  color: var(--text-inverse);
  opacity: 0.9;
}

.stat-card--glass {
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-color);
}

.stat-icon-wrapper {
  width: 64px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-lg);
  font-size: 32px;
  flex-shrink: 0;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

.stat-content {
  flex: 1;
  min-width: 0;
}

.stat-value {
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-bold);
  line-height: 1;
  margin-bottom: var(--spacing-sm);
  color: var(--text-primary);
}

.stat-label {
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
  margin-bottom: var(--spacing-xs);
}

.stat-trend {
  display: inline-flex;
  align-items: center;
  gap: var(--spacing-xs);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
  padding: 2px 8px;
  border-radius: var(--radius-sm);
}

.stat-trend--up {
  color: var(--success-color);
  background: rgba(0, 255, 136, 0.1);
}

.stat-trend--down {
  color: var(--error-color);
  background: rgba(255, 51, 102, 0.1);
}

.stat-skeleton {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(
    90deg,
    var(--bg-secondary) 25%,
    var(--bg-tertiary) 50%,
    var(--bg-secondary) 75%
  );
  background-size: 1000px 100%;
  animation: shimmer 2s infinite;
  border-radius: var(--radius-xl);
  z-index: 1;
}

@media (max-width: 640px) {
  .stat-card {
    padding: var(--spacing-md);
  }

  .stat-icon-wrapper {
    width: 48px;
    height: 48px;
    font-size: 24px;
  }

  .stat-value {
    font-size: var(--font-size-2xl);
  }
}
</style>
