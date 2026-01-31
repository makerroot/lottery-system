import { ref, computed } from 'vue'
import api from './api'

// 公司配置状态
const companyConfig = ref(null)
const companyCode = ref(null)
const isLoading = ref(false)

// 默认配置
export const defaultConfig = {
  id: 1,
  code: 'DEFAULT',
  name: '默认公司',
  logo: '',
  theme_color: '#00fff5',
  bg_color: '#0a0f14',
  title: '🎉 幸运大抽奖',
  subtitle: 'Luck Lottery System',
  welcome_text: '欢迎参与抽奖活动！',
  rules_text: '每人只能抽一次，祝您好运！',
  draw_button_text: '点击抽奖',
  success_text: '恭喜您获得',
  contact_name: '管理员',
  contact_phone: '400-888-8888',
  contact_email: 'admin@example.com'
}

// 加载公司配置（强制刷新）
export async function loadCompanyConfig(code) {
  if (!code) {
    // 如果没有提供code，尝试从URL获取
    // 注意：参数可能在hash中，如：#/lottery?company=4000
    let urlCode = null

    // 先尝试从hash中获取
    if (window.location.hash.includes('?')) {
      const hashParts = window.location.hash.split('?')
      if (hashParts.length > 1) {
        const hashParams = new URLSearchParams(hashParts[1])
        urlCode = hashParams.get('company')
      }
    }

    // 如果hash中没有，尝试从search中获取
    if (!urlCode && window.location.search) {
      const urlParams = new URLSearchParams(window.location.search)
      urlCode = urlParams.get('company')
    }

    code = urlCode || 'DEFAULT'
  }

  companyCode.value = code
  isLoading.value = true

  try {
    // 每次都重新获取，不使用缓存（使用api实例，自动带token）
    const response = await api.get('/api/company', {
      params: { code: code }
    })
    companyConfig.value = response
    return response
  } catch (error) {
    // 如果加载失败，使用默认配置
    companyConfig.value = defaultConfig
    return defaultConfig
  } finally {
    isLoading.value = false
  }
}

// 获取当前公司配置
export function useCompany() {
  const setCompanyCode = (code) => {
    companyCode.value = code
  }

  return {
    companyConfig: computed(() => companyConfig.value || defaultConfig),
    companyCode: computed(() => companyCode.value || 'DEFAULT'),
    setCompanyCode,
    isLoading: computed(() => isLoading.value),
    themeColor: computed(() => companyConfig.value?.theme_color || defaultConfig.theme_color),
    bgColor: computed(() => companyConfig.value?.bg_color || defaultConfig.bg_color),
    title: computed(() => companyConfig.value?.title || defaultConfig.title),
    subtitle: computed(() => companyConfig.value?.subtitle || defaultConfig.subtitle),
    welcomeText: computed(() => companyConfig.value?.welcome_text || defaultConfig.welcome_text),
    rulesText: computed(() => companyConfig.value?.rules_text || defaultConfig.rules_text),
    drawButtonText: computed(() => companyConfig.value?.draw_button_text || defaultConfig.draw_button_text),
    successText: computed(() => companyConfig.value?.success_text || defaultConfig.successText)
  }
}

// 重置公司配置
export function resetCompanyConfig() {
  companyConfig.value = null
  companyCode.value = null
}
