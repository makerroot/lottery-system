import { computed } from 'vue'

// 获取当前管理员信息
export function useAdmin() {
  const userStr = localStorage.getItem('admin_user')

  const isAdmin = computed(() => {
    return !!userStr
  })

  const currentUser = computed(() => {
    if (!userStr) return {}
    try {
      return JSON.parse(userStr)
    } catch (error) {
      return {}
    }
  })

  const isSuperAdmin = computed(() => {
    return currentUser.value?.is_super_admin || false
  })

  const canManageAllCompanies = computed(() => {
    return isSuperAdmin.value
  })

  return {
    isAdmin,
    currentUser,
    isSuperAdmin,
    canManageAllCompanies
  }
}
