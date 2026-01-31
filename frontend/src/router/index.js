import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/Home.vue'),
    meta: { title: '首页 - 幸运抽奖' }
  },
  {
    path: '/lottery',
    name: 'Lottery',
    component: () => import('../views/Lottery.vue'),
    meta: { title: '抽奖活动 - 幸运抽奖' }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/Register.vue'),
    meta: { title: '扫码注册 - 幸运抽奖' }
  },
  {
    path: '/admin',
    name: 'AdminLogin',
    component: () => import('../views/admin/Login.vue'),
    meta: { title: '管理员登录 - 幸运抽奖' }
  },
  {
    path: '/admin/dashboard',
    name: 'AdminDashboard',
    component: () => import('../views/admin/Layout.vue'),
    meta: { requiresAuth: true, title: '管理后台 - 幸运抽奖' },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('../views/admin/Dashboard.vue'),
        meta: { title: '数据概览 - 管理后台 - 幸运抽奖' }
      },
      {
        path: 'companies',
        name: 'Companies',
        component: () => import('../views/admin/Companies.vue'),
        meta: { title: '企业管理 - 管理后台 - 幸运抽奖' }
      },
      {
        path: 'admins',
        name: 'AdminManagement',
        component: () => import('../views/admin/AdminManagement.vue'),
        meta: { title: '管理员管理 - 管理后台 - 幸运抽奖' }
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('../views/admin/Users.vue'),
        meta: { title: '用户管理 - 管理后台 - 幸运抽奖' }
      },
      {
        path: 'prizes',
        name: 'PrizeLevels',
        component: () => import('../views/admin/PrizeLevels.vue'),
        meta: { title: '奖品管理 - 管理后台 - 幸运抽奖' }
      },
      {
        path: 'records',
        name: 'DrawRecords',
        component: () => import('../views/admin/DrawRecords.vue'),
        meta: { title: '抽奖记录 - 管理后台 - 幸运抽奖' }
      },
      {
        path: 'operation-logs',
        name: 'OperationLogs',
        component: () => import('../views/admin/OperationLogs.vue'),
        meta: { title: '操作日志 - 管理后台 - 幸运抽奖' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

// 路由守卫 - 检查管理员认证和设置动态 title
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('admin_token')

  // 设置动态 title
  if (to.meta.title) {
    document.title = to.meta.title
  } else {
    document.title = '幸运抽奖'
  }

  if (to.meta.requiresAuth && !token) {
    next('/admin')
  } else if (to.path === '/admin' && token) {
    next('/admin/dashboard')
  } else {
    next()
  }
})

export default router
