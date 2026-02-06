import { defineStore } from 'pinia'
import http from '@/common/axiosConfig'
import { ref } from 'vue'
import { useToast } from 'vue-toast-notification'
import type { NotificationApp } from '@/common/types'

const toast = useToast()

export const useAppsStore = defineStore('apps', () => {
  // 状态
  const apps = ref<Record<string, NotificationApp>>({})
  const loading = ref(false)

  // 获取所有应用
  const fetchApps = async () => {
    loading.value = true
    try {
      const response = await http.get('/admin/apps')
      apps.value = response.data
    } catch (error: any) {
      console.error('获取应用列表失败:', error)
      toast.error('获取应用列表失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 获取单个应用配置
  const getApp = async (appId: string) => {
    try {
      const response = await http.get(`/admin/apps/${appId}`)
      return response.data
    } catch (error: any) {
      console.error('获取应用配置失败:', error)
      toast.error('获取应用配置失败')
      throw error
    }
  }

  // 创建应用
  const createApp = async (appConfig: NotificationApp) => {
    loading.value = true
    try {
      const response = await http.post('/admin/apps', appConfig)

      if (response.code === 0) {
        // 更新本地状态
        apps.value[appConfig.appId] = appConfig
        fetchApps()
        toast.success(`应用 ${appConfig.name} 创建成功`)
      } else {
        toast.error(response.msg || '创建应用失败')
      }
      return response.data
    } catch (error: any) {
      console.error('创建应用失败:', error)
      if (error.response?.status === 409) {
        toast.error('应用ID已存在')
      } else {
        toast.error('创建应用失败')
      }
      throw error
    } finally {
      loading.value = false
    }
  }

  // 更新应用配置
  const updateApp = async (appId: string, appConfig: NotificationApp) => {
    loading.value = true
    try {
      const response = await http.request({
        method: 'PUT',
        url: `/admin/apps/${appId}`,
        data: appConfig,
      })

      if (response.code === 0) {
        // 更新本地状态
        apps.value[appId] = appConfig
        fetchApps()
        toast.success(`应用 ${appConfig.name} 配置更新成功`)
      } else {
        toast.error(response.msg || '更新应用配置失败')
      }
      return response.data
    } catch (error: any) {
      console.error('更新应用配置失败:', error)
      toast.error('更新应用配置失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 删除应用
  const deleteApp = async (appId: string) => {
    loading.value = true
    try {
      const response = await http.delete(`/admin/apps/${appId}`)

      if (response.code === 0) {
        // 从本地状态中移除
        delete apps.value[appId]
        fetchApps()
        toast.success(`应用 ${appId} 删除成功`)
      } else {
        toast.error(response.msg || '删除应用失败')
      }
    } catch (error: any) {
      console.error('删除应用失败:', error)
      toast.error('删除应用失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 发送测试通知
  const sendTestNotification = async (appId: string, testData: any) => {
    loading.value = true
    try {
      const response = await http.post(`/notify/${appId}`, testData)
      if (response.code === 0) {
        toast.success('测试通知发送成功')
      } else {
        if (response.code === 5001) {
          const msg = (response.msg || '发送测试通知失败').split('\n')
          msg.forEach((m) => toast.error(m))
        } else {
          toast.error(response.msg || '发送测试通知失败')
        }
      }
      return response.data
    } catch (error: any) {
      console.error('发送测试通知失败:', error)
      toast.error('发送测试通知失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 获取应用列表（仅name和enabled状态）
  const getAppsList = () => {
    return Object.entries(apps.value).map(([appId, app]) => ({
      appId,
      name: app.name,
      enabled: app.enabled,
      description: app.description,
    }))
  }

  return {
    // 状态
    apps,
    loading,

    // 方法
    fetchApps,
    getApp,
    createApp,
    updateApp,
    deleteApp,
    sendTestNotification,
    getAppsList,
  }
})
