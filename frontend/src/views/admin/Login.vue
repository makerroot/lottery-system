<template>
  <div class="admin-login-page">
    <!-- 背景 -->
    <div class="login-background">
      <div class="gradient-orb orb-1"></div>
      <div class="gradient-orb orb-2"></div>
      <div class="grid-overlay"></div>
      <div class="particles" ref="particlesRef"></div>
    </div>

    <!-- 登录卡片 -->
    <div class="login-container">
      <div class="login-card glass">
        <!-- Logo -->
        <div class="login-header">
          <div class="logo-icon">🎰</div>
          <h1 class="login-title font-display">NEON ADMIN</h1>
          <p class="login-subtitle font-body">管理后台登录</p>
        </div>

        <!-- 表单 -->
        <a-form
          :model="formData"
          layout="vertical"
          @submit="handleLogin"
          class="login-form"
        >
          <a-form-item>
            <label class="form-label font-body">
              <UserOutlined class="label-icon" />
              用户名
            </label>
            <a-input
              v-model:value="formData.username"
              placeholder="请输入用户名"
              size="large"
              class="neon-input"
              @keyup.enter="focusPassword"
            >
              <template #prefix>
                <UserOutlined />
              </template>
            </a-input>
          </a-form-item>

          <a-form-item>
            <label class="form-label font-body">
              <LockOutlined class="label-icon" />
              密码
            </label>
            <a-input-password
              ref="passwordInputRef"
              v-model:value="formData.password"
              placeholder="请输入密码"
              size="large"
              class="neon-input"
              @keyup.enter="handleLogin"
            >
              <template #prefix>
                <LockOutlined />
              </template>
            </a-input-password>
          </a-form-item>

          <a-form-item>
            <div class="form-options">
              <a-checkbox v-model:checked="formData.remember" class="neon-checkbox">
                记住我
              </a-checkbox>
            </div>
          </a-form-item>

          <a-form-item>
            <button
              type="submit"
              class="login-button font-display"
              :class="{ loading: loading }"
              :disabled="loading"
            >
              <span v-if="loading" class="button-spinner"></span>
              <span>{{ loading ? '登录中...' : '登录 →' }}</span>
            </button>
          </a-form-item>
        </a-form>

        <!-- 提示信息 -->
        <div class="login-tips">
          <div class="tip-item glass">
            <span class="tip-icon">💡</span>
            <div class="tip-content">
              <div class="tip-title font-body">默认账号</div>
              <div class="tip-desc font-mono">makerroot / 123456</div>
            </div>
          </div>
        </div>

        <!-- 返回首页 -->
        <div class="back-link">
          <router-link to="/" class="back-text font-body">
            ← 返回首页
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { UserOutlined, LockOutlined } from '@ant-design/icons-vue'
import request from '../../utils/request'

const router = useRouter()
const passwordInputRef = ref(null)
const particlesRef = ref(null)

const formData = reactive({
  username: '',
  password: '',
  remember: false
})

const loading = ref(false)

// 创建粒子效果
const createParticles = () => {
  const container = particlesRef.value
  if (!container) return

  for (let i = 0; i < 15; i++) {
    const particle = document.createElement('div')
    particle.className = 'particle'
    particle.style.left = `${Math.random() * 100}%`
    particle.style.top = `${Math.random() * 100}%`
    particle.style.animationDelay = `${Math.random() * 5}s`
    particle.style.animationDuration = `${5 + Math.random() * 10}s`
    container.appendChild(particle)
  }
}

onMounted(() => {
  const remembered = localStorage.getItem('remembered_username')
  if (remembered) {
    formData.username = remembered
    formData.remember = true
  }
  createParticles()
})

const focusPassword = () => {
  passwordInputRef.value?.focus()
}

const handleLogin = async (e) => {
  e?.preventDefault()

  if (!formData.username.trim()) {
    message.warning('请输入用户名')
    return
  }

  if (!formData.password.trim()) {
    message.warning('请输入密码')
    return
  }

  loading.value = true
  try {
    const data = await request.post('/admin/login', {
      username: formData.username.trim(),
      password: formData.password.trim()
    })

    localStorage.setItem('admin_token', data.token)
    localStorage.setItem('admin_user', JSON.stringify(data.user))

    if (formData.remember) {
      localStorage.setItem('remembered_username', formData.username)
    } else {
      localStorage.removeItem('remembered_username')
    }

    message.success('登录成功！欢迎回来 👋')

    setTimeout(() => {
      router.push('/admin/dashboard')
    }, 500)
  } catch (error) {
    const errorMsg = error.response?.data?.error || '登录失败，请检查用户名和密码'
    message.error(errorMsg)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.admin-login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  background: var(--bg-primary);
  overflow: hidden;
}

.login-background {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
  overflow: hidden;
  pointer-events: none;
}

.gradient-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(120px);
  opacity: 0.15;
  animation: orbFloat 20s ease-in-out infinite;
}

.orb-1 {
  width: 600px;
  height: 600px;
  background: var(--neon-cyan);
  top: -200px;
  left: -200px;
}

.orb-2 {
  width: 500px;
  height: 500px;
  background: var(--neon-magenta);
  bottom: -150px;
  right: -150px;
  animation-delay: 10s;
}

@keyframes orbFloat {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  50% {
    transform: translate(50px, -50px) scale(1.1);
  }
}

.grid-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image:
    linear-gradient(rgba(0, 255, 245, 0.05) 1px, transparent 1px),
    linear-gradient(90deg, rgba(0, 255, 245, 0.05) 1px, transparent 1px);
  background-size: 40px 40px;
  animation: gridMove 20s linear infinite;
}

@keyframes gridMove {
  0% {
    transform: translate(0, 0);
  }
  100% {
    transform: translate(40px, 40px);
  }
}

.particles {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.particle {
  position: absolute;
  width: 3px;
  height: 3px;
  background: var(--neon-cyan);
  border-radius: 50%;
  opacity: 0.6;
  animation: particleFloat 10s ease-in-out infinite;
  box-shadow: 0 0 8px var(--neon-cyan);
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
    transform: translateY(-100vh) translateX(30px);
    opacity: 0;
  }
}

.login-container {
  position: relative;
  z-index: 100;
  padding: var(--spacing-lg);
}

.login-card {
  width: 100%;
  max-width: 440px;
  padding: var(--spacing-3xl);
  border-radius: var(--radius-3xl);
  border: 2px solid var(--neon-cyan);
  box-shadow: var(--shadow-4), var(--glow-cyan);
  animation: scaleIn 0.6s ease-out;
  overflow: visible;
  position: relative;
  z-index: 10;
}

.login-header {
  text-align: center;
  margin-bottom: var(--spacing-3xl);
}

.logo-icon {
  font-size: 80px;
  margin-bottom: var(--spacing-lg);
  animation: float 3s ease-in-out infinite;
}

.login-title {
  font-size: var(--font-size-4xl);
  color: var(--neon-cyan);
  margin-bottom: var(--spacing-sm);
  letter-spacing: 6px;
  text-shadow: 0 0 30px var(--neon-cyan);
  animation: textGlow 3s ease-in-out infinite;
}

.login-subtitle {
  font-size: var(--font-size-base);
  color: var(--text-primary);
  letter-spacing: 2px;
}

.login-form {
  margin-bottom: var(--spacing-2xl);
  overflow: visible;
  position: relative;
  z-index: 5;
}

.form-label {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  color: var(--text-primary);
  font-size: var(--font-size-sm);
  text-transform: uppercase;
  letter-spacing: 1px;
  margin-bottom: var(--spacing-sm);
  font-weight: var(--font-weight-semibold);
}

.label-icon {
  color: var(--neon-cyan);
  font-size: var(--font-size-lg);
}

/* 输入框基础样式 - 确保文字可见 */
.neon-input :deep(.ant-input),
.neon-input :deep(.ant-input-password input) {
  background: rgba(255, 255, 255, 0.95) !important;
  border: 1px solid rgba(217, 217, 217, 0.8);
  border-radius: var(--radius-lg);
  color: #1a1a1a !important;
  transition: all var(--transition-base);
  position: relative;
  z-index: 10;
}

.neon-input :deep(.ant-input::placeholder),
.neon-input :deep(.ant-input-password input::placeholder) {
  color: #8c8c8c !important;
}

.neon-input :deep(.ant-input:focus),
.neon-input :deep(.ant-input-focused),
.neon-input :deep(.ant-input-password:focus),
.neon-input :deep(.ant-input-password-focused) {
  border-color: var(--neon-cyan) !important;
  box-shadow: 0 0 0 2px rgba(0, 255, 245, 0.2);
  background: rgba(255, 255, 255, 1) !important;
}

/* 密码输入框容器样式 */
.neon-input :deep(.ant-input-password) {
  background: rgba(255, 255, 255, 0.95) !important;
  border: 1px solid rgba(217, 217, 217, 0.8);
  border-radius: var(--radius-lg);
  position: relative;
  z-index: 10;
}

.neon-input :deep(.ant-input-password:hover) {
  border-color: var(--neon-cyan);
  background: rgba(255, 255, 255, 1) !important;
}

/* 确保密码输入框内的input元素样式正确 */
.neon-input :deep(.ant-input-password .ant-input) {
  background: transparent !important;
  color: #1a1a1a !important;
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.neon-checkbox :deep(.ant-checkbox-wrapper) {
  color: var(--text-primary);
}

.neon-checkbox :deep(.ant-checkbox-checked .ant-checkbox-inner) {
  background-color: var(--neon-cyan);
  border-color: var(--neon-cyan);
}

.login-button {
  width: 100%;
  padding: var(--spacing-lg);
  background: var(--primary-gradient);
  color: var(--text-inverse);
  font-size: var(--font-size-xl);
  border-radius: var(--radius-full);
  border: none;
  cursor: pointer;
  box-shadow: 0 0 40px var(--neon-cyan);
  transition: all var(--transition-base);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-md);
}

.login-button:hover:not(:disabled) {
  transform: translateY(-4px);
  box-shadow: 0 0 60px var(--neon-cyan), 0 0 100px var(--neon-magenta);
}

.login-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.login-button.loading {
  animation: pulse 1s ease-in-out infinite;
}

.button-spinner {
  width: 20px;
  height: 20px;
  border: 3px solid transparent;
  border-top-color: var(--text-inverse);
  border-right-color: var(--text-inverse);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.login-tips {
  margin-bottom: var(--spacing-2xl);
}

.tip-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-md);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-color);
  background: rgba(255, 255, 255, 0.02);
}

.tip-icon {
  font-size: var(--font-size-2xl);
  flex-shrink: 0;
}

.tip-content {
  flex: 1;
}

.tip-title {
  font-size: var(--font-size-sm);
  color: var(--text-primary);
  margin-bottom: var(--spacing-xs);
  font-weight: var(--font-weight-semibold);
}

.tip-desc {
  font-size: var(--font-size-xs);
  color: var(--neon-cyan);
  font-family: var(--font-mono);
}

.back-link {
  text-align: center;
}

.back-text {
  color: var(--text-primary);
  font-size: var(--font-size-sm);
  padding: var(--spacing-sm) var(--spacing-lg);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-full);
  display: inline-block;
  transition: all var(--transition-base);
}

.back-text:hover {
  color: var(--neon-cyan);
  border-color: var(--neon-cyan);
  box-shadow: var(--glow-cyan);
}

@media (max-width: 768px) {
  .login-card {
    padding: var(--spacing-2xl);
  }

  .login-title {
    font-size: var(--font-size-3xl);
  }

  .logo-icon {
    font-size: 60px;
  }

  .particle {
    display: none;
  }
}

@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
  }

  .particle {
    display: none;
  }
}
</style>
