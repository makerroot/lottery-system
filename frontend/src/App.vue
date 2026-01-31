<template>
  <a-config-provider :locale="zhCN">
    <router-view v-slot="{ Component }">
      <transition name="fade" mode="out-in">
        <component :is="Component" />
      </transition>
    </router-view>

    <!-- 全局技术支持信息 -->
    <footer class="global-footer">
      <div class="footer-content">
        <p class="footer-text">
          由 <span class="footer-highlight">makerroot</span> 提供技术支持
        </p>
      </div>
    </footer>
  </a-config-provider>
</template>

<script setup>
import zhCN from 'ant-design-vue/es/locale/zh_CN'
import { onMounted } from 'vue'

onMounted(() => {
  // 添加页面加载完成的class
  setTimeout(() => {
    document.body.classList.add('loaded')
  }, 100)
})
</script>

<style>
/* ============================================
   全局样式 - Neon Carnival Theme
   ============================================ */

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: var(--font-body);
  background: var(--bg-primary);
  color: var(--text-primary);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  overflow-x: hidden;
  line-height: 1.6;
}

/* 页面切换动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.4s cubic-bezier(0.4, 0, 0.2, 1),
              transform 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 全局默认字体 */
body {
  margin: 0;
  padding: 0;
  font-family: var(--font-body);
  font-size: var(--font-size-base);
  line-height: 1.6;
  color: var(--text-primary);
  background: var(--bg-primary);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

/* 所有文本元素默认使用 body 字体 */
h1, h2, h3, h4, h5, h6,
p, span, div, button, input, textarea, select {
  font-family: inherit;
}

/* Display 字体类 */
.font-display {
  font-family: var(--font-display);
  letter-spacing: 0.02em;
}

/* Body 字体类 */
.font-body {
  font-family: var(--font-body);
}

/* Mono 字体类 */
.font-mono {
  font-family: var(--font-mono);
  letter-spacing: -0.01em;
}

.fade-enter-from {
  opacity: 0;
  transform: translateY(20px) scale(0.98);
}

.fade-leave-to {
  opacity: 0;
  transform: translateY(-20px) scale(0.98);
}

/* 霓虹滚动条 */
::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}

::-webkit-scrollbar-track {
  background: var(--bg-secondary);
  border-radius: var(--radius-base);
}

::-webkit-scrollbar-thumb {
  background: linear-gradient(180deg, var(--neon-cyan), var(--neon-purple));
  border-radius: var(--radius-base);
  box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.3);
}

::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(180deg, var(--neon-magenta), var(--neon-yellow));
  box-shadow: var(--glow-cyan);
}

/* 全局加载动画 */
@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 页面淡入效果 */
body.loaded {
  animation: fadeIn 0.6s ease-in;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

/* 霓虹选中文字样式 */
::selection {
  background: var(--neon-cyan);
  color: var(--bg-primary);
  text-shadow: none;
}

::-moz-selection {
  background: var(--neon-cyan);
  color: var(--bg-primary);
  text-shadow: none;
}

/* 全局链接样式 */
a {
  color: var(--neon-cyan);
  text-decoration: none;
  transition: all var(--transition-fast);
}

a:hover {
  color: var(--neon-magenta);
  text-shadow: 0 0 10px rgba(255, 0, 255, 0.5);
}

/* 全局按钮样式重置 */
button {
  font-family: inherit;
  cursor: pointer;
  border: none;
  outline: none;
}

/* 全局输入框样式重置 */
input,
textarea,
select {
  font-family: inherit;
  outline: none;
}

/* 聚焦轮廓 */
:focus-visible {
  outline: 2px solid var(--neon-cyan);
  outline-offset: 2px;
  box-shadow: var(--glow-cyan);
}

/* 图片优化 */
img {
  max-width: 100%;
  height: auto;
  display: block;
}

/* 代码块样式 */
code,
pre {
  font-family: var(--font-mono);
  background: rgba(0, 255, 245, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
  color: var(--neon-cyan);
}

pre {
  padding: var(--spacing-md);
  overflow-x: auto;
  border: 1px solid var(--border-color);
}

/* 表格全局样式 */
table {
  width: 100%;
  border-collapse: collapse;
}

/* Ant Design 组件样式覆盖 */
.ant-card {
  background: rgba(26, 26, 36, 0.8) !important;
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-color) !important;
}

.ant-card-head {
  background: transparent !important;
  border-bottom: 1px solid var(--border-color) !important;
}

.ant-card-head-title {
  color: var(--text-primary) !important;
}

/* Ant Design 标签样式 - 确保良好的对比度 */
.ant-tag {
  font-weight: var(--font-weight-medium);
  border: 1px solid currentColor;
}

/* 白色标签特别处理 - 避免在深色背景上看不清 */
.ant-tag-white,
.ant-tag[style*="white"],
.ant-tag[color="white"] {
  background: var(--neon-cyan) !important;
  color: var(--bg-primary) !important;
  border-color: var(--neon-cyan) !important;
}

/* 确保所有标签文字颜色正确 */
.ant-tag {
  color: var(--bg-primary) !important;
}

.ant-tag-blue {
  background: var(--neon-cyan) !important;
  border-color: var(--neon-cyan) !important;
  color: var(--bg-primary) !important;
}

.ant-tag-green {
  background: var(--success-color) !important;
  border-color: var(--success-color) !important;
  color: var(--bg-primary) !important;
}

.ant-tag-red {
  background: var(--error-color) !important;
  border-color: var(--error-color) !important;
  color: var(--bg-primary) !important;
}

.ant-tag-cyan,
.ant-tag-purple,
.ant-tag-orange {
  color: var(--bg-primary) !important;
}

.ant-tag-default {
  background: rgba(255, 255, 255, 0.1) !important;
  border-color: var(--border-color) !important;
  color: var(--text-secondary) !important;
}

/* 确保小文本有足够的对比度 */
.level-probability,
.level-meta,
.rate-label,
.winner-time,
.footer-copyright,
.card-stat-label,
.stat-label {
  color: var(--text-secondary) !important;
}

/* 超小文本对比度增强 */
.font-size-xs,
.tip-desc,
.footer-version {
  color: var(--text-secondary) !important;
}

/* 确保描述文本可读 */
.level-description,
.tip-content,
.user-phone,
.company-code,
.time-date,
.prize-name {
  color: var(--text-secondary) !important;
}

/* 辅助信息对比度 */
.hint-text,
.helper-text,
.meta-text,
.sub-text {
  color: var(--text-secondary) !important;
}

/* 标签内的文字颜色 */
.ant-tag,
.level-badge {
  color: var(--bg-primary) !important;
  font-weight: var(--font-weight-medium);
}

/* 图例文字增强对比度 */
.legend-item {
  color: var(--text-secondary) !important;
}

/* 状态标签 */
.status-tag,
.state-tag,
.badge {
  background: rgba(255, 255, 255, 0.15) !important;
  color: var(--text-primary) !important;
  border: 1px solid var(--border-color) !important;
}

.ant-card-body {
  color: var(--text-primary) !important;
}

/* 按钮组件覆盖 */
.ant-btn-primary {
  background: var(--primary-gradient) !important;
  border: none !important;
  box-shadow: 0 4px 16px rgba(0, 255, 245, 0.3) !important;
}

.ant-btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 24px rgba(0, 255, 245, 0.5) !important;
}

/* 输入框组件覆盖 */
.ant-input {
  background: rgba(255, 255, 255, 0.05) !important;
  border: 1px solid var(--border-color) !important;
  color: var(--text-primary) !important;
}

.ant-input:hover,
.ant-input:focus {
  border-color: var(--neon-cyan) !important;
  box-shadow: 0 0 10px rgba(0, 255, 245, 0.3) !important;
}

/* 选择器组件覆盖 */
.ant-select-selector {
  background: rgba(255, 255, 255, 0.05) !important;
  border: 1px solid var(--border-color) !important;
  color: var(--text-primary) !important;
}

.ant-select-focused .ant-select-selector {
  border-color: var(--neon-cyan) !important;
  box-shadow: 0 0 10px rgba(0, 255, 245, 0.3) !important;
}

/* 模态框覆盖 */
.ant-modal-content {
  background: rgba(26, 26, 36, 0.95) !important;
  backdrop-filter: blur(30px);
  border: 2px solid var(--neon-cyan) !important;
  box-shadow: var(--shadow-4), var(--glow-cyan) !important;
}

.ant-modal-header {
  background: transparent !important;
  border-bottom: 1px solid var(--border-color) !important;
}

.ant-modal-body {
  background: rgba(26, 26, 36, 0.5) !important;
}

.ant-modal-title {
  color: var(--neon-cyan) !important;
}

.ant-modal-close-x {
  color: var(--text-secondary) !important;
}

.ant-modal-close-x:hover {
  color: var(--error-color) !important;
}

/* 表格组件覆盖 */
.ant-table {
  background: transparent !important;
}

.ant-table-thead > tr > th {
  background: rgba(0, 255, 245, 0.05) !important;
  border-bottom: 2px solid var(--neon-cyan) !important;
  color: var(--neon-cyan) !important;
}

.ant-table-tbody > tr > td {
  border-bottom: 1px solid var(--border-color) !important;
  color: var(--text-primary) !important;
}

.ant-table-tbody > tr:hover > td {
  background: rgba(0, 255, 245, 0.05) !important;
}

/* 标签覆盖 */
.ant-tag {
  background: rgba(255, 255, 255, 0.05) !important;
  border: 1px solid var(--border-color) !important;
  color: var(--text-primary) !important;
}

/* 消息提示覆盖 */
.ant-message {
  top: 20px !important;
}

.ant-message-notice-content {
  background: rgba(26, 26, 36, 0.95) !important;
  backdrop-filter: blur(20px);
  border: 1px solid var(--neon-cyan) !important;
  box-shadow: var(--shadow-3), var(--glow-cyan) !important;
  color: var(--text-primary) !important;
}

/* 通知覆盖 */
.ant-notification {
  background: rgba(26, 26, 36, 0.95) !important;
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-color) !important;
}

/* 加载指示器覆盖 */
.ant-spin-dot-item {
  background-color: var(--neon-cyan) !important;
}

/* 分页器覆盖 */
.ant-pagination-item {
  background: rgba(255, 255, 255, 0.05) !important;
  border: 1px solid var(--border-color) !important;
}

.ant-pagination-item:hover {
  border-color: var(--neon-cyan) !important;
}

.ant-pagination-item-active {
  background: var(--primary-gradient) !important;
  border-color: transparent !important;
}

.ant-pagination-item-active a {
  color: var(--text-inverse) !important;
}

/* 表单样式 */
.ant-form-item-label > label {
  color: #ffffff !important;
}

.ant-form-item-extra {
  color: #ffffff !important;
}

.ant-form-item-explain-error {
  color: var(--error-color) !important;
}

/* 减少动画偏好支持 */
@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}

/* 移动端优化 */
@media (max-width: 768px) {
  body {
    font-size: var(--font-size-sm);
  }

  .ant-card {
    margin: 0 -10px !important;
    border-radius: 0 !important;
  }
}

/* ============================================
   全局Footer
   ============================================ */

.global-footer {
  margin-top: var(--spacing-2xl);
  padding: var(--spacing-xl) var(--spacing-md);
  border-top: 1px solid var(--border-color);
  background: rgba(26, 26, 36, 0.5);
  backdrop-filter: blur(10px);
}

.footer-content {
  max-width: var(--container-max-width);
  margin: 0 auto;
  text-align: center;
}

.footer-text {
  color: var(--text-secondary);
  font-size: var(--font-size-sm);
  letter-spacing: 1px;
  margin: 0;
  transition: all var(--transition-fast);
}

.footer-highlight {
  color: var(--neon-cyan);
  font-weight: var(--font-weight-semibold);
  font-family: var(--font-display);
  text-shadow: 0 0 10px rgba(0, 255, 245, 0.3);
  transition: all var(--transition-fast);
  display: inline-block;
}

.global-footer:hover .footer-highlight {
  color: var(--neon-magenta);
  text-shadow: 0 0 15px rgba(255, 0, 255, 0.5);
  transform: scale(1.05);
}

/* 移动端Footer优化 */
@media (max-width: 768px) {
  .global-footer {
    padding: var(--spacing-md);
    margin-top: var(--spacing-xl);
  }

  .footer-text {
    font-size: var(--font-size-xs);
  }
}
</style>
