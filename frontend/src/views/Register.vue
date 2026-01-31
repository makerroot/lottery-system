<template>
  <div class="register-page">
    <!-- 背景效果 -->
    <div class="background-effects">
      <div class="gradient-orb orb-1"></div>
      <div class="gradient-orb orb-2"></div>
    </div>

    <!-- 注册容器 -->
    <div class="register-container">
      <div class="register-card glass">
        <!-- 公司信息 -->
        <div v-if="companyInfo" class="company-header">
          <h1 class="company-title">{{ companyInfo.company.name }}</h1>
          <p class="company-subtitle">扫码参与抽奖活动</p>
        </div>

        <!-- 注册表单 -->
        <a-form
          v-if="!registeredSuccess"
          :model="registerForm"
          layout="vertical"
          @finish="handleRegister"
          class="register-form"
        >
          <a-alert
            v-if="errorMessage"
            :message="errorMessage"
            type="error"
            show-icon
            closable
            @close="errorMessage = ''"
            style="margin-bottom: 16px"
          />

          <a-form-item label="姓名">
            <a-input
              v-model:value="registerForm.name"
              placeholder="请输入真实姓名"
              size="large"
              :disabled="registering"
            >
              <template #prefix>
                <IdcardOutlined />
              </template>
            </a-input>
          </a-form-item>

          <a-form-item label="手机号（选填）">
            <a-input
              v-model:value="registerForm.phone"
              placeholder="请输入手机号"
              size="large"
              :disabled="registering"
            >
              <template #prefix>
                <PhoneOutlined />
              </template>
            </a-input>
          </a-form-item>

          <a-form-item>
            <a-button
              type="primary"
              html-type="submit"
              size="large"
              block
              :loading="registering"
              class="register-button"
            >
              {{ registering ? '注册中...' : '🎉 立即参与抽奖' }}
            </a-button>
          </a-form-item>

          <div class="register-tips">
            <p class="tip-text">
              <InfoCircleOutlined /> 注册成功后将自动加入抽奖池
            </p>
            <p class="tip-text">
              <SafetyOutlined /> 您的信息将被安全保护
            </p>
          </div>
        </a-form>

        <!-- 注册成功 -->
        <div v-else class="success-container">
          <div class="success-icon">✅</div>
          <h2 class="success-title">参与成功！</h2>
          <p class="success-message">您已成功加入抽奖池</p>
          <div class="success-info">
            <p><strong>姓名：</strong>{{ registerForm.name }}</p>
            <p v-if="registerForm.phone"><strong>手机号：</strong>{{ registerForm.phone }}</p>
          </div>
          <a-button type="primary" size="large" @click="goToLottery" class="view-lottery-btn">
            查看抽奖页面
          </a-button>
        </div>
      </div>

      <!-- 参与统计 -->
      <div v-if="companyInfo && !registeredSuccess" class="stats-card glass">
        <div class="stat-item">
          <div class="stat-icon">👥</div>
          <div class="stat-content">
            <div class="stat-value">{{ companyInfo.stats.total_users }}</div>
            <div class="stat-label">已参与</div>
          </div>
        </div>
        <div class="stat-divider"></div>
        <div class="stat-item">
          <div class="stat-icon">⏳</div>
          <div class="stat-content">
            <div class="stat-value">{{ companyInfo.stats.undrawn_users }}</div>
            <div class="stat-label">待抽奖</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import {
  IdcardOutlined,
  PhoneOutlined,
  InfoCircleOutlined,
  SafetyOutlined
} from '@ant-design/icons-vue'
import api from '../utils/api'

const route = useRoute()
const router = useRouter()

const companyCode = ref('')
const companyInfo = ref(null)
const registering = ref(false)
const registeredSuccess = ref(false)
const errorMessage = ref('')

const registerForm = ref({
  name: '',
  phone: ''
})

// 获取公司信息
const fetchCompanyInfo = async () => {
  try {
    const code = route.query.company_code
    if (!code) {
      errorMessage.value = '缺少公司代码参数'
      return
    }

    companyCode.value = code
    const data = await api.get(`/api/company-info?company_code=${code}`)
    companyInfo.value = data
  } catch (error) {
    errorMessage.value = error.response?.data?.error || '获取公司信息失败'
  }
}

// 处理注册
const handleRegister = async () => {
  // 验证表单
  if (!registerForm.value.name) {
    errorMessage.value = '请输入姓名'
    return
  }

  registering.value = true
  errorMessage.value = ''

  try {
    const response = await api.post('/api/self-register', {
      name: registerForm.value.name,
      phone: registerForm.value.phone
    }, {
      params: { company_code: companyCode.value }
    })

    registeredSuccess.value = true
    message.success('参与成功！')
  } catch (error) {
    if (error.response?.status === 409) {
      errorMessage.value = error.response.data?.error || '用户已存在'
    } else {
      errorMessage.value = error.response?.data?.error || '注册失败，请稍后重试'
    }
  } finally {
    registering.value = false
  }
}

// 跳转到抽奖页面
const goToLottery = () => {
  router.push(`/lottery?company=${companyCode.value}`)
}

onMounted(() => {
  fetchCompanyInfo()
})
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  position: relative;
  overflow: hidden;
}

.background-effects {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
  pointer-events: none;
}

.gradient-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.3;
  animation: float 20s ease-in-out infinite;
}

.orb-1 {
  width: 400px;
  height: 400px;
  background: #667eea;
  top: -100px;
  left: -100px;
}

.orb-2 {
  width: 300px;
  height: 300px;
  background: #764ba2;
  bottom: -50px;
  right: -50px;
  animation-delay: 7s;
}

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  50% {
    transform: translate(30px, -30px) scale(1.1);
  }
}

.register-container {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 480px;
}

.register-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  padding: 40px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.company-header {
  text-align: center;
  margin-bottom: 32px;
}

.company-title {
  font-size: 28px;
  font-weight: bold;
  color: #1a1a1a;
  margin-bottom: 8px;
}

.company-subtitle {
  font-size: 14px;
  color: #666;
  margin: 0;
}

.register-form {
  margin-top: 24px;
}

.register-button {
  height: 50px;
  font-size: 16px;
  font-weight: bold;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 12px;
  margin-top: 16px;
}

.register-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.4);
}

.register-tips {
  margin-top: 24px;
  padding: 16px;
  background: rgba(102, 126, 234, 0.05);
  border-radius: 12px;
  border: 1px solid rgba(102, 126, 234, 0.2);
}

.tip-text {
  font-size: 13px;
  color: #666;
  margin: 8px 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.success-container {
  text-align: center;
  padding: 20px 0;
}

.success-icon {
  font-size: 64px;
  margin-bottom: 16px;
  animation: successBounce 0.6s ease-out;
}

@keyframes successBounce {
  0% {
    transform: scale(0);
  }
  50% {
    transform: scale(1.2);
  }
  100% {
    transform: scale(1);
  }
}

.success-title {
  font-size: 24px;
  font-weight: bold;
  color: #1a1a1a;
  margin-bottom: 8px;
}

.success-message {
  font-size: 14px;
  color: #666;
  margin-bottom: 24px;
}

.success-info {
  background: rgba(82, 196, 26, 0.05);
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 24px;
  text-align: left;
}

.success-info p {
  margin: 8px 0;
  font-size: 14px;
  color: #1a1a1a;
}

.view-lottery-btn {
  height: 50px;
  font-size: 16px;
  font-weight: bold;
  background: linear-gradient(135deg, #52c41a 0%, #73d13d 100%);
  border: none;
  border-radius: 12px;
}

.stats-card {
  margin-top: 20px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  padding: 20px;
  display: flex;
  align-items: center;
  justify-content: space-around;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.stat-icon {
  font-size: 32px;
}

.stat-content {
  text-align: left;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #1a1a1a;
  line-height: 1;
}

.stat-label {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.stat-divider {
  width: 1px;
  height: 40px;
  background: rgba(0, 0, 0, 0.1);
}

@media (max-width: 576px) {
  .register-card {
    padding: 24px;
  }

  .company-title {
    font-size: 24px;
  }
}
</style>
