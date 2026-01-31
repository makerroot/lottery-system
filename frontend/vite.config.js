import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      },
      '/admin': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  },
  build: {
    // 增加 chunk 大小警告限制到 1500 KB（Ant Design Vue 较大）
    chunkSizeWarningLimit: 1500,

    // 配置代码分割
    rollupOptions: {
      output: {
        // 手动分包策略
        manualChunks(id) {
          // Vue 核心库
          if (id.includes('node_modules/vue/') || id.includes('node_modules/@vue/')) {
            return 'vue-vendor'
          }

          // Ant Design Vue 图标
          if (id.includes('@ant-design/icons-vue')) {
            return 'antd-icons'
          }

          // Ant Design Vue 组件库
          if (id.includes('node_modules/ant-design-vue/')) {
            return 'antd-vue'
          }

          // Axios
          if (id.includes('node_modules/axios/')) {
            return 'axios'
          }

          // Dayjs
          if (id.includes('node_modules/dayjs/')) {
            return 'dayjs'
          }

          // 其他 node_modules
          if (id.includes('node_modules/')) {
            return 'vendor'
          }
        }
      }
    }
  }
})
