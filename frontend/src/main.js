import { createApp } from 'vue'
import router from './router'
import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/reset.css'
import './styles/design-system.css'
import './styles/global-enhancements.css'
import './assets/input-fix.css' // 全局输入框样式修复
import App from './App.vue'

const app = createApp(App)
app.use(Antd)
app.use(router)
app.mount('#app')
