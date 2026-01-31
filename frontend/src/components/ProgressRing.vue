<template>
  <div class="progress-ring" :style="{ width: size + 'px', height: size + 'px' }">
    <svg
      :width="size"
      :height="size"
      class="progress-ring__svg"
    >
      <circle
        :stroke="trackColor"
        :stroke-width="stroke"
        fill="transparent"
        :r="normalizedRadius"
        :cx="size / 2"
        :cy="size / 2"
      />
      <circle
        :stroke="strokeColor"
        :stroke-width="stroke"
        fill="transparent"
        :r="normalizedRadius"
        :cx="size / 2"
        :cy="size / 2"
        class="progress-ring__circle"
        :style="{ strokeDasharray: circumference + ' ' + circumference, strokeDashoffset: strokeDashoffset }"
      />
    </svg>
    <div class="progress-ring__content">
      <slot>
        <span class="progress-ring__text font-display">{{ formattedPercent }}%</span>
      </slot>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  percent: {
    type: Number,
    default: 0,
    validator: (value) => value >= 0 && value <= 100
  },
  size: {
    type: Number,
    default: 120
  },
  stroke: {
    type: Number,
    default: 8
  },
  strokeColor: {
    type: String,
    default: 'var(--neon-cyan)'
  },
  trackColor: {
    type: String,
    default: 'var(--bg-secondary)'
  }
})

const normalizedRadius = computed(() => props.size / 2 - props.stroke)
const circumference = computed(() => normalizedRadius.value * 2 * Math.PI)

const strokeDashoffset = computed(() => {
  const progress = props.percent / 100
  return circumference.value * (1 - progress)
})

const formattedPercent = computed(() => Math.round(props.percent))
</script>

<style scoped>
.progress-ring {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.progress-ring__svg {
  transform: rotate(-90deg);
}

.progress-ring__circle {
  transition: stroke-dashoffset 0.5s ease;
  transform-origin: 50% 50%;
  stroke-linecap: round;
}

.progress-ring__content {
  position: absolute;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
}

.progress-ring__text {
  font-size: 24px;
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  text-shadow: 0 0 10px var(--neon-cyan);
}

@keyframes progressAnimation {
  from {
    stroke-dashoffset: 100%;
  }
  to {
    stroke-dashoffset: var(--target-offset);
  }
}
</style>
