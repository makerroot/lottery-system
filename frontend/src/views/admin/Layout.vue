<template>
  <div class="admin-layout-container">
    <!-- 侧边栏 -->
    <aside class="sidebar" :class="{ 'sidebar--collapsed': collapsed }">
      <!-- Logo -->
      <div class="sidebar-header">
        <div class="logo-icon">🎰</div>
        <h2 v-if="!collapsed" class="logo-text font-display">MAKERROOT</h2>
      </div>

      <!-- 导航菜单 -->
      <nav class="sidebar-nav">
        <router-link
          to="/admin/dashboard"
          class="nav-item"
          :class="{ active: isActive('/admin/dashboard') }"
        >
          <DashboardOutlined class="nav-icon" />
          <span class="nav-text font-body">数据概览</span>
          <div class="nav-indicator"></div>
        </router-link>

        <router-link
          to="/admin/dashboard/companies"
          class="nav-item"
          :class="{ active: isActive('/admin/dashboard/companies') }"
        >
          <TeamOutlined class="nav-icon" />
          <span class="nav-text font-body">公司管理</span>
          <div class="nav-indicator"></div>
        </router-link>

        <router-link
          to="/admin/dashboard/users"
          class="nav-item"
          :class="{ active: isActive('/admin/dashboard/users') }"
        >
          <UserOutlined class="nav-icon" />
          <span class="nav-text font-body">用户管理</span>
          <div class="nav-indicator"></div>
        </router-link>

        <router-link
          to="/admin/dashboard/prizes"
          class="nav-item"
          :class="{ active: isActive('/admin/dashboard/prizes') }"
        >
          <GiftOutlined class="nav-icon" />
          <span class="nav-text font-body">奖品管理</span>
          <div class="nav-indicator"></div>
        </router-link>

        <router-link
          to="/admin/dashboard/records"
          class="nav-item"
          :class="{ active: isActive('/admin/dashboard/records') }"
        >
          <FileTextOutlined class="nav-icon" />
          <span class="nav-text font-body">抽奖记录</span>
          <div class="nav-indicator"></div>
        </router-link>

        <router-link
          to="/admin/dashboard/operation-logs"
          class="nav-item"
          :class="{ active: isActive('/admin/dashboard/operation-logs') }"
        >
          <HistoryOutlined class="nav-icon" />
          <span class="nav-text font-body">操作日志</span>
          <div class="nav-indicator"></div>
        </router-link>

        <router-link
          to="/admin/dashboard/admins"
          class="nav-item"
          :class="{ active: isActive('/admin/dashboard/admins') }"
        >
          <SafetyOutlined class="nav-icon" />
          <span class="nav-text font-body">管理员管理</span>
          <div class="nav-indicator"></div>
        </router-link>
      </nav>

      <!-- 底部信息 -->
      <div v-if="!collapsed" class="sidebar-footer">
        <div class="footer-info">
          <p class="footer-version font-mono">v2.0.0</p>
          <p class="footer-copyright">© 2026 MakerRoot Admin</p>
        </div>
      </div>
    </aside>

    <!-- 主内容区 -->
    <div class="main-wrapper">
      <!-- 顶部栏 -->
      <header class="top-bar glass">
        <div class="top-left">
          <button
            @click="collapsed = !collapsed"
            class="toggle-btn"
            :title="collapsed ? '展开侧边栏' : '收起侧边栏'"
          >
            <MenuUnfoldOutlined v-if="collapsed" />
            <MenuFoldOutlined v-else />
          </button>

          <div class="breadcrumb">
            <span class="breadcrumb-item font-body">管理后台</span>
            <span class="breadcrumb-separator">/</span>
            <span class="breadcrumb-item current font-body">{{ currentPageName }}</span>
          </div>
        </div>

        <div class="top-right">
          <!-- 用户信息 -->
          <div class="user-section">
            <a-tag
              v-if="currentUser.is_super_admin"
              color="purple"
              class="super-tag"
            >
              超级管理员
            </a-tag>
            <span v-if="currentUser.company" class="company-name font-body">
              {{ currentUser.company.name }}
            </span>
          </div>

          <!-- 用户下拉菜单 -->
          <a-dropdown>
            <div class="user-trigger">
              <div class="user-avatar font-display">
                {{ (currentUser.username || 'A')[0].toUpperCase() }}
              </div>
              <span class="user-name font-body">{{ currentUser.username || 'Admin' }}</span>
              <DownOutlined class="dropdown-icon" />
            </div>
            <template #overlay>
              <a-menu class="user-dropdown-menu">
                <a-menu-item @click="showChangePasswordModal">
                  <LockOutlined />
                  修改密码
                </a-menu-item>
                <a-menu-divider />
                <a-menu-item @click="handleLogout" class="logout-item">
                  <LogoutOutlined />
                  退出登录
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </header>

      <!-- 内容区 -->
      <main class="main-content">
        <router-view />
      </main>
    </div>

    <!-- 修改密码弹窗 -->
    <a-modal
      v-model:open="changePasswordVisible"
      title="修改密码"
      width="500px"
      :maskClosable="false"
      @ok="handleChangePassword"
      @cancel="changePasswordVisible = false"
    >
      <a-form :model="passwordForm" layout="vertical">
        <a-form-item>
          <label class="form-label font-body">
            <span class="label-icon">🔑</span>
            当前密码
          </label>
          <a-input-password
            v-model:value="passwordForm.oldPassword"
            placeholder="请输入当前密码"
            class="neon-input"
          />
        </a-form-item>
        <a-form-item>
          <label class="form-label font-body">
            <span class="label-icon">🔒</span>
            新密码
          </label>
          <a-input-password
            v-model:value="passwordForm.newPassword"
            placeholder="请输入新密码（至少6位）"
            class="neon-input"
          />
        </a-form-item>
        <a-form-item>
          <label class="form-label font-body">
            <span class="label-icon">✓</span>
            确认新密码
          </label>
          <a-input-password
            v-model:value="passwordForm.confirmPassword"
            placeholder="请再次输入新密码"
            class="neon-input"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import request from '../../utils/request'
import {
  DashboardOutlined,
  TeamOutlined,
  UserOutlined,
  GiftOutlined,
  FileTextOutlined,
  HistoryOutlined,
  SafetyOutlined,
  LogoutOutlined,
  LockOutlined,
  DownOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined
} from '@ant-design/icons-vue'

const router = useRouter()
const route = useRoute()
const collapsed = ref(false)

// 修改密码相关
const changePasswordVisible = ref(false)
const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const isActive = (path) => {
  // 数据概览页面需要精确匹配
  if (path === '/admin/dashboard') {
    return route.path === '/admin/dashboard' || route.path === '/admin/dashboard/'
  }
  // 其他页面精确匹配
  return route.path === path
}

const currentPageName = computed(() => {
  const path = route.path
  if (path === '/admin/dashboard' || path === '/admin/dashboard/') {
    return '数据概览'
  } else if (path.includes('/admin/dashboard/companies')) {
    return '公司管理'
  } else if (path.includes('/admin/dashboard/users')) {
    return '用户管理'
  } else if (path.includes('/admin/dashboard/prizes')) {
    return '奖品管理'
  } else if (path.includes('/admin/dashboard/records')) {
    return '抽奖记录'
  } else if (path.includes('/admin/dashboard/operation-logs')) {
    return '操作日志'
  } else if (path.includes('/admin/dashboard/admins')) {
    return '管理员管理'
  }
  return '管理后台'
})

const currentUser = computed(() => {
  const user = localStorage.getItem('admin_user')
  return user ? JSON.parse(user) : {}
})

const handleLogout = () => {
  localStorage.removeItem('admin_token')
  localStorage.removeItem('admin_user')
  message.success('已退出登录')
  router.push('/admin')
}

const showChangePasswordModal = () => {
  passwordForm.value = {
    oldPassword: '',
    newPassword: '',
    confirmPassword: ''
  }
  changePasswordVisible.value = true
}

const handleChangePassword = async () => {
  // 验证表单
  if (!passwordForm.value.oldPassword) {
    message.warning('请输入当前密码')
    return
  }
  if (!passwordForm.value.newPassword) {
    message.warning('请输入新密码')
    return
  }
  if (passwordForm.value.newPassword.length < 6) {
    message.warning('新密码长度至少为6位')
    return
  }
  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    message.warning('两次输入的新密码不一致')
    return
  }
  if (passwordForm.value.oldPassword === passwordForm.value.newPassword) {
    message.warning('新密码不能与当前密码相同')
    return
  }

  try {
    await request.post('/admin/change-password', {
      old_password: passwordForm.value.oldPassword,
      new_password: passwordForm.value.newPassword
    })
    message.success('密码修改成功，请重新登录')
    changePasswordVisible.value = false
    // 延迟退出登录，让用户看到成功提示
    setTimeout(() => {
      handleLogout()
    }, 1500)
  } catch (error) {
    message.error(error.response?.data?.error || '密码修改失败')
  }
}
</script>

<style scoped>
.admin-layout-container {
  min-height: 100vh;
  display: flex;
  background: var(--bg-primary);
}

/* ============================================
   侧边栏
   ============================================ */

.sidebar {
  width: 260px;
  background: rgba(20, 20, 32, 0.95);
  backdrop-filter: blur(20px);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  transition: width var(--transition-base);
  position: fixed;
  height: 100vh;
  z-index: 100;
}

.sidebar--collapsed {
  width: 80px;
}

.sidebar-header {
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-md);
  padding: 0 var(--spacing-lg);
  border-bottom: 1px solid var(--border-color);
}

.logo-icon {
  font-size: 40px;
  animation: float 3s ease-in-out infinite;
}

.logo-text {
  font-size: var(--font-size-2xl);
  color: var(--neon-cyan);
  letter-spacing: 4px;
  text-shadow: 0 0 20px var(--neon-cyan);
  animation: textGlow 3s ease-in-out infinite;
}

.sidebar-nav {
  flex: 1;
  padding: var(--spacing-xl) var(--spacing-md);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.nav-item {
  position: relative;
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-md) var(--spacing-lg);
  border-radius: var(--radius-lg);
  color: var(--text-primary);
  text-decoration: none;
  transition: all var(--transition-base);
  margin-bottom: var(--spacing-sm);
}

.nav-item::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 4px;
  height: 0;
  background: var(--neon-cyan);
  border-radius: 0 var(--radius-full);
  transition: height var(--transition-base);
}

.nav-indicator {
  position: absolute;
  right: var(--spacing-md);
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--neon-cyan);
  opacity: 0;
  box-shadow: 0 0 10px var(--neon-cyan);
  transition: all var(--transition-base);
}

.nav-icon {
  font-size: var(--font-size-lg);
  flex-shrink: 0;
}

.nav-text {
  font-size: var(--font-size-base);
  flex: 1;
  white-space: nowrap;
}

.sidebar--collapsed .nav-text {
  display: none;
}

.sidebar--collapsed .nav-indicator {
  display: none;
}

.nav-item:hover {
  background: rgba(0, 255, 245, 0.1);
  color: var(--text-primary);
  transform: translateX(4px);
}

.nav-item.active {
  background: rgba(0, 255, 245, 0.15);
  color: var(--neon-cyan);
}

.nav-item.active::before {
  height: 24px;
}

.nav-item.active .nav-indicator {
  opacity: 1;
  animation: pulse 2s ease-in-out infinite;
}

.sidebar-footer {
  padding: var(--spacing-xl);
  border-top: 1px solid var(--border-color);
}

.footer-info {
  text-align: center;
}

.footer-version {
  font-size: var(--font-size-xs);
  color: var(--neon-cyan);
  font-family: var(--font-mono);
  margin-bottom: var(--spacing-xs);
}

.footer-copyright {
  font-size: var(--font-size-xs);
  color: var(--text-tertiary);
}

/* ============================================
   主内容区
   ============================================ */

.main-wrapper {
  flex: 1;
  margin-left: 260px;
  transition: margin-left var(--transition-base);
}

.sidebar--collapsed + .main-wrapper {
  margin-left: 80px;
}

.top-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-lg) var(--spacing-2xl);
  margin-bottom: var(--spacing-xl);
  border-radius: var(--radius-xl);
  border: 1px solid var(--border-color);
  min-height: 72px;
}

.top-left {
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
}

.toggle-btn {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-base);
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border-color);
  color: var(--text-primary);
  cursor: pointer;
  transition: all var(--transition-base);
}

.toggle-btn:hover {
  background: rgba(0, 255, 245, 0.1);
  border-color: var(--neon-cyan);
  color: var(--neon-cyan);
  box-shadow: 0 0 15px rgba(0, 255, 245, 0.3);
}

.breadcrumb {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.breadcrumb-item {
  font-size: var(--font-size-sm);
  color: var(--text-primary);
}

.breadcrumb-item.current {
  color: var(--neon-cyan);
  font-weight: var(--font-weight-semibold);
}

.breadcrumb-separator {
  color: var(--text-tertiary);
}

.top-right {
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
}

.user-section {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.super-tag {
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-semibold);
  padding: 4px 12px;
  border-radius: var(--radius-full);
  border: 1px solid var(--neon-purple);
  background: rgba(185, 79, 255, 0.1);
  color: var(--neon-purple);
}

.company-name {
  font-size: var(--font-size-sm);
  color: var(--text-tertiary);
  padding: 4px 12px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: var(--radius-full);
}

.user-trigger {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-sm) var(--spacing-md);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--transition-base);
}

.user-trigger:hover {
  background: rgba(255, 255, 255, 0.05);
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--primary-gradient);
  color: var(--text-inverse);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-bold);
}

.user-name {
  font-size: var(--font-size-sm);
  color: var(--text-primary);
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dropdown-icon {
  font-size: var(--font-size-xs);
  color: var(--text-tertiary);
  transition: transform var(--transition-base);
}

.user-trigger:hover .dropdown-icon {
  transform: rotate(180deg);
}

.user-dropdown-menu {
  min-width: 160px;
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-3);
}

.logout-item {
  color: var(--error-color);
}

.logout-item:hover {
  background: rgba(255, 51, 102, 0.1);
}

.main-content {
  padding: 0 var(--spacing-2xl) var(--spacing-2xl);
  animation: fadeIn 0.3s ease-out;
}

/* ============================================
   响应式
   ============================================ */

@media (max-width: 768px) {
  .sidebar {
    transform: translateX(-100%);
    z-index: 1000;
  }

  .main-wrapper {
    margin-left: 0;
  }

  .sidebar--collapsed + .main-wrapper {
    margin-left: 0;
  }

  .breadcrumb {
    display: none;
  }

  .company-name {
    display: none;
  }

  .user-name {
    display: none;
  }

  .top-bar {
    padding: var(--spacing-md) var(--spacing-lg);
  }
}

/* ============================================
   修改密码弹窗输入框样式
   ============================================ */

/* 修复密码输入框文字颜色问题 */
:deep(.ant-modal-content) {
  background: rgba(30, 30, 40, 0.98) !important;
  border: 1px solid var(--border-color);
}

:deep(.ant-modal-header) {
  background: transparent !important;
  border-bottom: 1px solid var(--border-color);
}

:deep(.ant-modal-title) {
  color: var(--text-primary) !important;
}

:deep(.ant-modal-body) {
  background: transparent !important;
}

/* 修改密码弹窗的输入框样式 */
:deep(.ant-modal .ant-input),
:deep(.ant-modal .ant-input-password input) {
  background: rgba(255, 255, 255, 0.95) !important;
  border: 1px solid rgba(217, 217, 217, 0.8);
  border-radius: var(--radius-lg);
  color: #1a1a1a !important;
  font-size: 14px;
}

:deep(.ant-modal .ant-input::placeholder),
:deep(.ant-modal .ant-input-password input::placeholder) {
  color: #8c8c8c !important;
}

:deep(.ant-modal .ant-input:focus),
:deep(.ant-modal .ant-input-password:focus) {
  border-color: var(--neon-cyan) !important;
  box-shadow: 0 0 0 2px rgba(0, 255, 245, 0.2);
  background: rgba(255, 255, 255, 1) !important;
}

:deep(.ant-modal .ant-input-password) {
  background: rgba(255, 255, 255, 0.95) !important;
  border: 1px solid rgba(217, 217, 217, 0.8);
  border-radius: var(--radius-lg);
}

:deep(.ant-modal .ant-input-password:hover) {
  border-color: var(--neon-cyan);
  background: rgba(255, 255, 255, 1) !important;
}

/* 确保密码输入框内的input元素样式正确 */
:deep(.ant-modal .ant-input-password .ant-input) {
  background: transparent !important;
  color: #1a1a1a !important;
}

/* 表单标签颜色 */
:deep(.ant-modal .ant-form-item-label > label) {
  color: var(--text-primary) !important;
}

/* 表单标签样式 */
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

/* 输入框统一样式 */
.neon-input :deep(.ant-input),
.neon-input :deep(.ant-input-password input) {
  background: rgba(255, 255, 255, 0.95) !important;
  border: 1px solid rgba(217, 217, 217, 0.8);
  border-radius: var(--radius-lg);
  color: #1a1a1a !important;
  transition: all var(--transition-base);
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

.neon-input :deep(.ant-input-password) {
  background: rgba(255, 255, 255, 0.95) !important;
  border: 1px solid rgba(217, 217, 217, 0.8);
  border-radius: var(--radius-lg);
}

.neon-input :deep(.ant-input-password:hover) {
  border-color: var(--neon-cyan);
  background: rgba(255, 255, 255, 1) !important;
}

.neon-input :deep(.ant-input-password .ant-input) {
  background: transparent !important;
  color: #1a1a1a !important;
}
</style>
