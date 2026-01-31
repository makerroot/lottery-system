import axios from 'axios'
import { message } from 'ant-design-vue'

// 创建axios实例
// 使用相对路径，通过 Caddy 反向代理访问后端
const request = axios.create({
  baseURL: '/',  // 使用相对路径，自动适配当前域名和端口
  timeout: 10000
})

// 请求拦截器
request.interceptors.request.use(
  config => {
    // 从localStorage获取token
    const token = localStorage.getItem('admin_token')
    if (token) {
      // 设置Authorization header
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  response => {
    // 直接返回response.data
    return response.data
  },
  error => {

    // 处理不同的错误状态码
    if (error.response) {
      const { status, data } = error.response

      switch (status) {
        case 401:
          // 未授权，清除token并跳转到登录页
          message.error('登录已过期，请重新登录')
          localStorage.removeItem('admin_token')
          localStorage.removeItem('admin_user')
          window.location.href = '/admin'
          break
        case 403:
          message.error('没有权限访问')
          break
        case 404:
          message.error('请求的资源不存在')
          break
        case 500:
          message.error('服务器错误，请稍后重试')
          break
        default:
          message.error(data.error || '请求失败')
      }
    } else if (error.request) {
      // 请求已发送但没有收到响应
      message.error('网络错误，请检查网络连接')
    } else {
      // 请求配置出错
      message.error('请求配置错误')
    }

    return Promise.reject(error)
  }
)

export default request
