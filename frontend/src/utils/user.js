import { ref } from 'vue'

// 获取当前用户信息
export function useUser() {
  const userStr = ref(localStorage.getItem('lottery_user'))

  const isLoggedIn = ref(!!userStr.value)

  const currentUser = ref({})

  // 解析用户信息
  const parseUser = () => {
    if (!userStr.value) {
      currentUser.value = {}
      return
    }
    try {
      currentUser.value = JSON.parse(userStr.value)
    } catch (error) {
      currentUser.value = {}
    }
  }

  // 监听localStorage变化
  const checkLoginStatus = () => {
    userStr.value = localStorage.getItem('lottery_user')
    isLoggedIn.value = !!userStr.value
    parseUser()
  }

  // 初始化时检查一次
  checkLoginStatus()

  const setUser = (userData) => {
    // 如果有token，保存token
    if (userData.token) {
      localStorage.setItem('lottery_token', userData.token)
      // 从响应中提取纯用户数据
      const user = userData.user || userData
      localStorage.setItem('lottery_user', JSON.stringify(user))
    } else {
      // 兼容旧格式，直接保存用户数据
      localStorage.setItem('lottery_user', JSON.stringify(userData))
    }
    checkLoginStatus()
  }

  const clearUser = () => {
    localStorage.removeItem('lottery_user')
    localStorage.removeItem('lottery_token')
    checkLoginStatus()
  }

  // 监听storage事件（跨标签页同步）
  window.addEventListener('storage', (e) => {
    if (e.key === 'lottery_user' || e.key === null) {
      checkLoginStatus()
    }
  })

  return {
    userLoggedIn: isLoggedIn,
    currentUser,
    setUser,
    clearUser,
    checkLoginStatus
  }
}
