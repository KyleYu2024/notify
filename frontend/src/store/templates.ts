import { defineStore } from 'pinia'
import http from '@/common/axiosConfig'
import { ref } from 'vue'
import { useToast } from 'vue-toast-notification'
import { copyToClipboard } from '@/common/utils'

const toast = useToast()

export interface IMessageTemplate {
  id: string
  name: string
  content: string
  title: string
  image: string
  url: string
  targets: string
}

export const useTemplatesStore = defineStore('templates', () => {
  // 状态
  const templates = ref<Record<string, IMessageTemplate>>({})
  const loading = ref(false)

  // 获取所有模板
  const fetchTemplates = async () => {
    loading.value = true
    try {
      const response = await http.get('/admin/templates')
      templates.value = response.data
    } catch (error: any) {
      console.error('获取模板列表失败:', error)
      toast.error('获取模板列表失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 获取单个模板
  const getTemplate = async (templateId: string) => {
    try {
      const response = await http.get(`/admin/templates/${templateId}`)
      return response.data
    } catch (error: any) {
      console.error('获取模板失败:', error)
      toast.error('获取模板失败')
      throw error
    }
  }

  // 创建模板
  const createTemplate = async (templateData: IMessageTemplate) => {
    loading.value = true
    try {
      const response = await http.post('/admin/templates', templateData)
      if (response.code === 0) {
        // 更新本地状态
        templates.value[templateData.id] = templateData

        toast.success(`模板 ${templateData.name} 创建成功`)
      } else {
        toast.error(response.msg || '创建模板失败')
      }
      return response.data
    } catch (error: any) {
      console.error('创建模板失败:', error)
      if (error.response?.status === 409) {
        toast.error('模板ID已存在')
      } else {
        toast.error('创建模板失败')
      }
      throw error
    } finally {
      loading.value = false
    }
  }

  // 更新模板
  const updateTemplate = async (templateId: string, templateData: IMessageTemplate) => {
    loading.value = true
    try {
      const response = await http.request({
        method: 'PUT',
        url: `/admin/templates/${templateId}`,
        data: templateData,
      })
      if (response.code === 0) {
        // 更新本地状态
        templates.value[templateId] = templateData
        fetchTemplates()
        toast.success(`模板 ${templateData.name} 更新成功`)
      } else {
        toast.error(response.msg || '更新模板失败')
      }
      return response.data
    } catch (error: any) {
      console.error('更新模板失败:', error)
      toast.error('更新模板失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 删除模板
  const deleteTemplate = async (templateId: string) => {
    loading.value = true
    try {
      const response = await http.delete(`/admin/templates/${templateId}`)

      if (response.code === 0) {
        // 从本地状态中移除
        delete templates.value[templateId]
        fetchTemplates()
        toast.success(`模板 ${templateId} 删除成功`)
      } else {
        toast.error(response.msg || '删除模板失败')
      }
    } catch (error: any) {
      toast.error('删除模板失败')
    } finally {
      loading.value = false
    }
  }

  // 获取模板选项列表（用于下拉选择）
  const getTemplateOptions = () => {
    return Object.entries(templates.value).map(([id, template]) => ({
      value: id,
      title: template.name,
      subtitle: `ID: ${id}`,
    }))
  }

  // 预览模板渲染结果
  const previewTemplate = (templateContent: string, data: Record<string, any>) => {
    try {
      // 简单的模板预览，将变量替换为示例数据
      let preview = templateContent

      // 替换常见变量
      const variables: Record<string, any> = {
        title: data.title || '示例标题',
        content: data.content || '示例内容',
        level: data.level || 'info',
        message: data.message || '示例消息',
        timestamp: new Date().toLocaleString('zh-CN'),
        image: data.image || '',
        url: data.url || '',
        ...data,
      }

      Object.entries(variables).forEach(([key, value]) => {
        // 替换 {{.Key}} 格式的变量
        const regex = new RegExp(`{{\\s*\\.${key}\\s*}}`, 'g')
        preview = preview.replace(regex, String(value))

        // 替换带条件的变量 {{if .Key}}content{{end}}
        const ifRegex = new RegExp(`{{\\s*if\\s+\\.${key}\\s*}}([^{]*){{\\s*end\\s*}}`, 'g')
        preview = preview.replace(ifRegex, value ? '$1' : '')
      })

      // 处理 upper 过滤器
      preview = preview.replace(/{{\\s*\\.(\w+)\\s*\|\\s*upper\\s*}}/g, (match, key) => {
        const value = variables[key]
        return value ? String(value).toUpperCase() : match
      })

      return preview
    } catch (error) {
      console.error('模板预览失败:', error)
      return '模板预览失败'
    }
  }

  // 分享模板到剪贴板
  const shareTemplateToClipboard = async (template: IMessageTemplate) => {
    try {
      const exportData = {
        version: '1.0',
        exportTime: new Date().toISOString(),
        exportType: 'single',
        templates: [template],
      }

      const jsonString = JSON.stringify(exportData, null, 2)

      const success = await copyToClipboard(jsonString)

      if (success) {
        toast.success(`模板 "${template.name}" 已分享到剪贴板`)
      } else {
        toast.error('分享到剪贴板失败')
      }
    } catch (error: any) {
      console.error('分享模板失败:', error)
      toast.error('分享模板失败')
    }
  }

  // 生成随机ID
  const generateRandomId = (baseName: string = 'template'): string => {
    const timestamp = Date.now().toString(36)
    const randomStr = Math.random().toString(36).substring(2, 8)
    let newId = `${baseName}_${timestamp}_${randomStr}`

    // 确保ID唯一
    let counter = 1
    const originalId = newId
    while (templates.value[newId]) {
      newId = `${originalId}_${counter}`
      counter++
    }

    return newId
  }

  // 批量导入模板
  const importTemplates = async (
    templatesToImport: IMessageTemplate[],
    overwrite: boolean = false,
    generateNewIds: boolean = false
  ) => {
    loading.value = true
    let successCount = 0
    let skipCount = 0
    let errorCount = 0
    let renamedCount = 0

    try {
      for (const template of templatesToImport) {
        try {
          const exists = templates.value[template.id]
          let templateToImport = { ...template }

          if (exists) {
            if (overwrite) {
              await updateTemplate(template.id, templateToImport)
              successCount++
            } else if (generateNewIds) {
              // 生成新的随机ID
              const newId = generateRandomId(template.id.replace(/[^a-zA-Z0-9_-]/g, ''))
              templateToImport.id = newId
              await createTemplate(templateToImport)
              successCount++
              renamedCount++
            } else {
              skipCount++
              continue
            }
          } else {
            await createTemplate(templateToImport)
            successCount++
          }
        } catch (error) {
          console.error(`导入模板 ${template.id} 失败:`, error)
          errorCount++
        }
      }

      // 显示导入结果
      if (successCount > 0) {
        let message = `成功导入 ${successCount} 个模板`
        if (renamedCount > 0) {
          message += `，其中 ${renamedCount} 个模板使用了新的ID`
        }
        if (skipCount > 0) {
          message += `，跳过 ${skipCount} 个已存在的模板`
        }
        toast.success(message)
      }

      if (errorCount > 0) {
        toast.error(`${errorCount} 个模板导入失败`)
      }

      if (successCount === 0 && skipCount > 0) {
        toast.info(`跳过了 ${skipCount} 个已存在的模板`)
      }
    } catch (error: any) {
      console.error('批量导入模板失败:', error)
      toast.error('批量导入模板失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  return {
    // 状态
    templates,
    loading,

    // 方法
    fetchTemplates,
    getTemplate,
    createTemplate,
    updateTemplate,
    deleteTemplate,
    getTemplateOptions,
    previewTemplate,

    // 导入导出方法
    shareTemplateToClipboard,
    importTemplates,
  }
})
