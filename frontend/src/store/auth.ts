import { defineStore } from 'pinia'
import http from '@/common/axiosConfig'
import axios from 'axios' // 保留axios用于特殊认证处理
import { ref } from 'vue'
import { useToast } from 'vue-toast-notification'

const toast = useToast()

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const isAuthenticated = ref(false)
  const username = ref('')
  const isAuthRequired = ref(true)
  const loading = ref(false)
  const apiBaseUrl = axios.defaults.baseURL
  const version = ref('')
  // 获取系统状态
  const checkSystemStatus = async () => {
    try {
      const response = await http.get('/health')

      const healthData = response.data
      // 从 adminAuthRequired 字段获取是否需要 admin 认证
      isAuthRequired.value = healthData.adminAuthRequired || false
      version.value = healthData.version || ''
      return healthData
    } catch (error) {
      console.error('检查系统状态失败:', error)
      throw error
    }
  }

  // 登录
  const login = async (usernameInput: string, password: string) => {
    loading.value = true
    try {
      // 设置基础认证头
      const credentials = btoa(`${usernameInput}:${password}`)

      // 测试认证 - 这里需要使用axios来设置特殊的认证头
      await axios.get('admin/apps', {
        headers: {
          Authorization: `Basic ${credentials}`,
        },
      })

      // 保存认证信息
      username.value = usernameInput
      isAuthenticated.value = true

      // 设置默认认证头
      axios.defaults.headers.common['Authorization'] = `Basic ${credentials}`

      // 保存到localStorage
      localStorage.setItem('auth_credentials', credentials)
      localStorage.setItem('username', usernameInput)

      toast.success('登录成功')
      return true
    } catch (error: any) {
      console.error('登录失败:', error)
      if (error.response?.status === 401) {
        toast.error('用户名或密码错误')
      } else {
        toast.error('登录失败，请检查网络连接')
      }
      throw error
    } finally {
      loading.value = false
    }
  }

  // 退出登录
  const logout = () => {
    isAuthenticated.value = false
    username.value = ''
    delete axios.defaults.headers.common['Authorization']
    localStorage.removeItem('auth_credentials')
    localStorage.removeItem('username')
    toast.info('已退出登录')
  }

  // 从localStorage恢复认证状态
  const restoreAuth = () => {
    const credentials = localStorage.getItem('auth_credentials')
    const savedUsername = localStorage.getItem('username')

    if (credentials && savedUsername) {
      username.value = savedUsername
      isAuthenticated.value = true
      axios.defaults.headers.common['Authorization'] = `Basic ${credentials}`
    }
  }

  return {
    // 状态
    isAuthenticated,
    username,
    isAuthRequired,
    loading,
    // 方法
    login,
    logout,
    restoreAuth,
    checkSystemStatus,
    apiBaseUrl,
    version,
  }
})
