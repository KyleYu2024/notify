<template>
  <v-dialog v-model="dialog" max-width="800px">
    <v-card class="template-import-dialog">
      <v-card-title>
        <v-icon icon="mdi-import" class="mr-2"></v-icon>
        导入模板
      </v-card-title>

      <v-card-text>
        <v-textarea v-model="importText" label="粘贴模板JSON" placeholder="请粘贴分享的模板JSON内容..." variant="outlined" rows="8"
          auto-grow hint="粘贴从其他地方分享的模板JSON内容" persistent-hint class="mb-4" @input="handleTextChange" />

        <!-- 预览解析结果 -->
        <v-card v-if="parsedTemplates.length > 0" variant="outlined" class="mt-4">
          <v-card-title class="text-h6">
            解析结果
            <v-chip class="ml-2" size="small" color="primary">
              {{ parsedTemplates.length }} 个模板
            </v-chip>
          </v-card-title>

          <v-card-text>
            <v-list density="compact">
              <v-list-item v-for="template in parsedTemplates" :key="template.id" class="px-0">
                <template v-slot:prepend>
                  <v-icon :icon="getTemplateStatusIcon(template)" :color="getTemplateStatusColor(template)"
                    class="mr-3" />
                </template>

                <v-list-item-title>
                  {{ template.name }}
                  <v-chip v-if="getTemplateStatus(template) === 'exists'" size="x-small" color="warning" class="ml-2">
                    已存在
                  </v-chip>
                </v-list-item-title>

                <v-list-item-subtitle>
                  ID: {{ template.id }}
                </v-list-item-subtitle>

                <template v-slot:append>
                  <v-checkbox v-model="selectedTemplates" :value="template.id" hide-details density="compact" />
                </template>
              </v-list-item>
            </v-list>

            <!-- 导入选项 -->
            <v-divider class="my-4" />

            <v-row>
              <v-col cols="12">
                <div class="text-subtitle-2 mb-3">ID冲突处理策略</div>
                <v-radio-group v-model="conflictStrategy" hide-details>
                  <v-radio label="跳过已存在的模板" value="skip">
                    <template v-slot:label>
                      <div>
                        <div class="text-body-2">跳过已存在的模板</div>
                        <div class="text-caption text-medium-emphasis">保持现有模板不变，忽略重复的模板</div>
                      </div>
                    </template>
                  </v-radio>

                  <v-radio label="覆盖已存在的模板" value="overwrite">
                    <template v-slot:label>
                      <div>
                        <div class="text-body-2">覆盖已存在的模板</div>
                        <div class="text-caption text-medium-emphasis">用导入的模板替换现有的同ID模板</div>
                      </div>
                    </template>
                  </v-radio>

                  <v-radio label="自动生成新ID" value="generate">
                    <template v-slot:label>
                      <div>
                        <div class="text-body-2">自动生成新ID</div>
                        <div class="text-caption text-medium-emphasis">为冲突的模板自动生成新的唯一ID</div>
                      </div>
                    </template>
                  </v-radio>
                </v-radio-group>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>

        <!-- 错误信息 -->
        <v-alert v-if="errorMessage" type="error" variant="outlined" class="mt-4" closable
          @click:close="errorMessage = ''">
          {{ errorMessage }}
        </v-alert>
      </v-card-text>

      <v-card-actions class="operation-actions">
        <v-btn text="取消" @click="handleCancel" />
        <v-spacer />
        <v-btn color="primary" text="导入" @click="handleImport" :disabled="!canImport" :loading="loading"
          variant="flat" />
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import type { IMessageTemplate } from '@/store/templates'

interface Props {
  existingTemplates?: Record<string, IMessageTemplate>
  loading?: boolean
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'import', data: { templates: IMessageTemplate[], conflictStrategy: 'skip' | 'overwrite' | 'generate' }): void
  (e: 'cancel'): void
}

const props = withDefaults(defineProps<Props>(), {
  existingTemplates: () => ({}),
  loading: false
})

const emit = defineEmits<Emits>()

const dialog = defineModel<boolean>()

// 导入状态
const importText = ref('')

// 解析状态
const parsedTemplates = ref<IMessageTemplate[]>([])
const selectedTemplates = ref<string[]>([])
const conflictStrategy = ref<'skip' | 'overwrite' | 'generate'>('generate')
const errorMessage = ref('')

// 计算属性
const canImport = computed(() => {
  return parsedTemplates.value.length > 0 && selectedTemplates.value.length > 0
})

// 获取模板状态
const getTemplateStatus = (template: IMessageTemplate) => {
  return props.existingTemplates[template.id] ? 'exists' : 'new'
}

const getTemplateStatusIcon = (template: IMessageTemplate) => {
  return getTemplateStatus(template) === 'exists' ? 'mdi-alert-circle' : 'mdi-check-circle'
}

const getTemplateStatusColor = (template: IMessageTemplate) => {
  return getTemplateStatus(template) === 'exists' ? 'warning' : 'success'
}

// 处理文本变化
const handleTextChange = () => {
  if (importText.value.trim()) {
    parseTemplateData(importText.value)
  } else {
    clearParsedData()
  }
}

// 解析模板数据
const parseTemplateData = (jsonString: string) => {
  try {
    errorMessage.value = ''
    const data = JSON.parse(jsonString)

    // 支持多种格式
    let templates: IMessageTemplate[] = []

    if (Array.isArray(data)) {
      // 数组格式：[template1, template2, ...]
      templates = data
    } else if (data.templates && Array.isArray(data.templates)) {
      // 对象格式：{ templates: [template1, template2, ...], version: "1.0" }
      templates = data.templates
    } else if (data.id && data.name) {
      // 单个模板格式
      templates = [data]
    } else if (typeof data === 'object') {
      // Record格式：{ "template1": {...}, "template2": {...} }
      templates = Object.values(data)
    } else {
      throw new Error('不支持的数据格式')
    }

    // 验证模板数据
    const validTemplates = templates.filter(template => {
      return template.id && template.name &&
        typeof template.title === 'string' &&
        typeof template.content === 'string'
    })

    if (validTemplates.length === 0) {
      throw new Error('未找到有效的模板数据')
    }

    parsedTemplates.value = validTemplates
    selectedTemplates.value = validTemplates.map(t => t.id)

  } catch (error: any) {
    console.error('解析模板数据失败:', error)
    errorMessage.value = `解析失败: ${error.message || '无效的JSON格式'}`
    clearParsedData()
  }
}

// 清空解析数据
const clearParsedData = () => {
  parsedTemplates.value = []
  selectedTemplates.value = []
}

// 处理导入
const handleImport = () => {
  const templatesToImport = parsedTemplates.value.filter(
    template => selectedTemplates.value.includes(template.id)
  )

  emit('import', {
    templates: templatesToImport,
    conflictStrategy: conflictStrategy.value
  })
}

// 处理取消
const handleCancel = () => {
  // 重置状态
  importText.value = ''
  clearParsedData()
  conflictStrategy.value = 'generate'
  errorMessage.value = ''

  emit('cancel')
}

// 监听对话框关闭
watch(dialog, (value) => {
  if (!value) {
    handleCancel()
  }
})
</script>

<style lang="less" scoped>
@import '@/styles/mix.less';

.template-import-dialog {
  max-height: 80vh;
  overflow-y: auto;
  .scrollbar();
  background-color: rgba(var(--v-theme-surface), 0.6);
  backdrop-filter: blur(30px);
}

.operation-actions {
  position: sticky;
  bottom: 0;
  z-index: 3;
  background-color: rgba(var(--v-theme-surface), 0.8);
  backdrop-filter: blur(10px);
}
</style>
