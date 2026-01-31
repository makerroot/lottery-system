import axios from 'axios'

// 创建用于抽奖页面的axios实例
// 使用相对路径，通过 Caddy 反向代理访问后端
export const api = axios.create({
  baseURL: '/',  // 使用相对路径，自动适配当前域名和端口
  timeout: 10000
})

// 请求拦截器：自动添加token
api.interceptors.request.use(
  config => {
    // 优先使用用户token，如果没有则使用管理员token
    const token = localStorage.getItem('lottery_token') || localStorage.getItem('admin_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器：处理401错误
api.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    if (error.response?.status === 401) {
      // 清除用户和管理员信息
      localStorage.removeItem('lottery_token')
      localStorage.removeItem('lottery_user')
      localStorage.removeItem('admin_token')
      localStorage.removeItem('admin_user')
      // 可以选择跳转到首页或显示登录提示
    }
    return Promise.reject(error)
  }
)

export default api
