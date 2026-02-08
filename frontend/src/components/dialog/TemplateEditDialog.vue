<template>
  <v-dialog v-model="dialog" max-width="1000px">
    <v-card class="edit-template-dialog">
      <v-card-title>
        {{ editingTemplate ? '编辑模板' : '创建模板' }}
      </v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12" md="6">
            <v-form ref="formRef" v-model="formValid">
              <v-text-field v-model="form.id" label="模板ID *" hint="唯一标识符，建议使用下划线分隔" persistent-hint
                :rules="[rules.required, rules.templateId]" :disabled="!!editingTemplate" class="mb-4"></v-text-field>
              <v-text-field v-model="form.name" label="模板名称 *" :rules="[rules.required]" class="mb-4"></v-text-field>
              <v-textarea v-model="form.title" label="标题" hint="支持Go模板语法，如 {{.Title}}" persistent-hint
                :rules="[rules.required]" rows="3" auto-grow class="mb-4"></v-textarea>
              <v-textarea v-model="form.content" label="内容" hint="支持Go模板语法，如 {{.Content}}" persistent-hint
                :rules="[rules.required]" rows="3" auto-grow class="mb-4"></v-textarea>
              <v-textarea v-model="form.url" label="链接" hint="支持Go模板语法，如  {{.url}}" persistent-hint
                :rules="[rules.required]" rows="2" auto-grow class="mb-4"></v-textarea>
              <v-textarea v-model="form.image" label="图片" hint="支持Go模板语法，如{{.image}}" persistent-hint
                :rules="[rules.required]" rows="2" auto-grow class="mb-4"></v-textarea>
              <v-textarea v-model="form.targets" label="目标" hint="支持Go模板语法，如 {{.targets}}" persistent-hint
                :rules="[rules.required]" rows="2" auto-grow class="mb-4"></v-textarea>
            </v-form>
          </v-col>
          <v-col cols="12" md="6">
            <v-expansion-panels class="mb-4" v-model="expanded">
              <v-expansion-panel title="模板变量说明">
                <v-expansion-panel-text>
                  <v-list density="compact">
                    <v-list-item v-for="variable in templateVariables" :key="variable.name" class="px-0">
                      <template v-slot:prepend>
                        <v-code>{{ variable.name }}</v-code>
                      </template>
                      <v-list-item-title class="text-caption">
                        {{ variable.description }}
                      </v-list-item-title>
                    </v-list-item>
                  </v-list>
                  <v-alert type="info" variant="tonal" class="mt-4"
                    text="请求body中的所有参数都可以在这里配置，如 {{.title}} 等"></v-alert>
                </v-expansion-panel-text>
              </v-expansion-panel>
            </v-expansion-panels>
          </v-col>

        </v-row>
      </v-card-text>

      <v-card-actions class="operation-actions">
        <v-btn text="取消" @click="handleCancel"></v-btn>
        <v-btn color="primary" text="保存" @click="handleSave" :disabled="!formValid" :loading="loading"
          variant="flat"></v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import type { IMessageTemplate } from '@/store/templates'
import { useDisplay } from 'vuetify'

interface Props {
  editingTemplate?: IMessageTemplate | null
  loading?: boolean
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'save', data: any): void
  (e: 'cancel'): void
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

const emit = defineEmits<Emits>()

const dialog = defineModel<boolean>()
const display = useDisplay()

// 表单状态
const formValid = ref(false)
const formRef = ref()

// 表单数据
const form = ref({
  id: '',
  name: '',
  title: '{{.title}}',
  content: '{{.content}}',
  url: '{{.url}}',
  image: '{{.image}}',
  targets: '{{.targets}}',
})

const expanded = ref<number | undefined>()

watchEffect(() => {
  if (display.mdAndUp.value) {
    console.debug('mdAndUp', display.mdAndUp)
    expanded.value = 0
  }
})


// 验证规则
const rules = {
  required: (value: any) => !!value || '此字段为必填项',
  templateId: (value: string) => /^[a-zA-Z0-9_-]+$/.test(value) || 'ID只能包含字母、数字、下划线和连字符'
}

// 模板变量列表
const templateVariables = [
  { name: '{{.title}}', description: '通知标题' },
  { name: '{{.content}}', description: '通知内容' },
  { name: '{{.message}}', description: '消息内容（兼容字段）' },
  { name: '{{.image}}', description: '图片URL' },
  { name: '{{.url}}', description: '链接URL' },
  { name: '{{.timestamp}}', description: '时间戳' },
  { name: '{{.targets}}', description: '目标' },
  { name: '{{if .url}}...{{end}}', description: '条件渲染' }
]

// 监听编辑模板变化
watch(() => props.editingTemplate, (template) => {
  if (template) {
    form.value = {
      id: template.id,
      name: template.name,
      title: template.title,
      content: template.content,
      url: template.url,
      image: template.image,
      targets: template.targets,
    }
  } else {
    form.value = {
      id: '',
      name: '',
      title: '{{.title}}',
      content: '{{.content}}',
      url: '{{.url}}',
      image: '{{.image}}',
      targets: '{{.targets}}',
    }
  }
}, { immediate: true })

// 处理保存
const handleSave = async () => {
  const res = await formRef.value.validate()
  if (!res) return
  await emit('save', form.value)
  dialog.value = false
}

// 处理取消
const handleCancel = () => {
  emit('cancel')
}
</script>

<style lang="less" scoped>
@import '@/styles/mix.less';

.edit-template-dialog {
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
