import { defineStore } from 'pinia'
import http from '@/common/axiosConfig'
import { useToast } from 'vue-toast-notification'

const toast = useToast()
export interface UIConfig {
  component: string
  text?: string
  html?: string
  content?: any
  slots?: any
  props?: any
  events?: any
}

export interface PluginInfo {
  id: string
  name: string
  version: string
  description?: string
  author?: string
  enabled: boolean
  loadedAt: number
  ui?: UIConfig
  settings?: Record<string, any>
  test_data?: any
}

interface State {
  loading: boolean
  list: PluginInfo[]
}

export const usePluginsStore = defineStore('plugins', {
  state: (): State => ({
    loading: false,
    list: [],
  }),
  getters: {
    getById: (state) => (id: string) => state.list.find((p) => p.id === id),
  },
  actions: {
    async fetchList() {
      this.loading = true
      try {
        const res = await http.get<PluginInfo[]>('/admin/plugins')
        if (res.code === 0 && Array.isArray(res.data)) {
          this.list = res.data
        }
      } finally {
        this.loading = false
      }
    },
    async fetchOne(id: string) {
      const res = await http.get<PluginInfo>(`/admin/plugins/${id}`)
      if (res.code === 0 && res.data) {
        const idx = this.list.findIndex((p) => p.id === id)
        if (idx >= 0) this.list[idx] = res.data
        else this.list.push(res.data)
      }
    },
    async setEnabled(id: string, enabled: boolean) {
      const url = enabled ? `/admin/plugins/${id}/enable` : `/admin/plugins/${id}/disable`
      const res = await http.request<{ message?: string }>({ method: 'PUT', url })
      if (res.code === 0) {
        toast.success('插件设置成功')
        const target = this.list.find((p) => p.id === id)
        if (target) target.enabled = enabled
      } else {
        console.error('插件设置失败:', res)
        toast.error(res.msg || '插件设置失败')
      }
      return res
    },
    async updateSettings(id: string, settings: Record<string, any>) {
      const res = await http.request<{ message?: string }>({
        method: 'PUT',
        url: `/admin/plugins/${id}/config`,
        data: { settings },
      })
      if (res.code === 0) {
        toast.success('插件配置更新成功')
        const target = this.list.find((p) => p.id === id)
        if (target) target.settings = settings
      }
      return res
    },
    async testPlugin(id: string, input: Record<string, any>) {
      const res = await http.post<any>(`/admin/plugins/${id}/test`, { input })
      return res
    },
  },
})
