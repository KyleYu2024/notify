import { defineStore } from 'pinia'
import http from '@/common/axiosConfig'
import { ref } from 'vue'
import { useToast } from 'vue-toast-notification'
import { NotifierTypeMap } from '@/common/types'

const toast = useToast()

export interface INotifierInstance {
  type: string
  enabled: boolean
  config: Record<string, any>
}

export const useNotifiersStore = defineStore('notifiers', () => {
  // 状态
  const notifiers = ref<Record<string, INotifierInstance>>({})
  const loading = ref(false)

  // 获取所有通知服务
  const fetchNotifiers = async () => {
    loading.value = true
    try {
      const response = await http.get('/admin/notifiers')
      notifiers.value = response.data
    } catch (error: any) {
      console.error('获取通知服务列表失败:', error)
      toast.error('获取通知服务列表失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 获取单个通知服务配置
  const getNotifier = async (notifierName: string) => {
    try {
      const response = await http.get(`/admin/notifiers/${notifierName}`)
      return response.data
    } catch (error: any) {
      console.error('获取通知服务配置失败:', error)
      toast.error('获取通知服务配置失败')
      throw error
    }
  }

  // 更新通知服务配置
  const updateNotifier = async (notifierName: string, config: INotifierInstance) => {
    loading.value = true
    try {
      const response = await http.request({
        method: 'PUT',
        url: `/admin/notifiers/${notifierName}`,
        data: config,
      })

      if (response.code === 0) {
        await fetchNotifiers()
        toast.success(`通知服务 ${notifierName} 配置更新成功`)
      } else {
        toast.error(response.msg || '更新通知服务配置失败')
      }
      return response.data
    } catch (error: any) {
      console.error('更新通知服务配置失败:', error)
      toast.error('更新通知服务配置失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 删除通知服务
  const deleteNotifier = async (notifierName: string) => {
    loading.value = true
    try {
      const response = await http.delete(`/admin/notifiers/${notifierName}`)

      if (response.code === 0) {
        // 从本地状态中移除
        delete notifiers.value[notifierName]
        fetchNotifiers()
        toast.success(`通知服务 ${notifierName} 删除成功`)
      } else {
        toast.error(response.msg || '删除通知服务失败')
      }
    } catch (error: any) {
      console.error('删除通知服务失败:', error)
      toast.error('删除通知服务失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 获取通知服务类型列表
  const getNotifierTypes = () => {
    return [
      {
        type: NotifierTypeMap.wechatWorkAPPBot,
        name: '企业微信',
        fields: ['agent_id', 'corp_id', 'secret'],
      },
      {
        type: NotifierTypeMap.telegramAppBot,
        name: 'Telegram',
        fields: ['bot_token', 'chat_id'],
      },
      {
        type: NotifierTypeMap.dingTalkAppBot,
        name: '钉钉',
        fields: ['webhook_url', 'secret'],
      },
    ]
  }

  return {
    // 状态
    notifiers,
    loading,

    // 方法
    fetchNotifiers,
    getNotifier,
    updateNotifier,
    deleteNotifier,
    getNotifierTypes,
  }
})
