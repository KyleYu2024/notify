<template>
  <v-dialog v-model="dialogVisible" max-width="800px" @update:model-value="handleDialogClose">
    <v-card class="dialog-card">
      <v-card-title>
        {{ editingApp ? '编辑应用' : '创建应用' }}
      </v-card-title>

      <v-card-text>
        <v-form ref="formRef">
          <v-row>
            <v-col cols="12" md="6">
              <v-text-field v-model="form.appId" label="应用ID *" :rules="[rules.required, rules.appId]"
                :disabled="!!editingApp" hint="只能包含字母、数字、下划线和连字符" persistent-hint></v-text-field>
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field v-model="form.name" label="应用名称 *" :rules="[rules.required]"></v-text-field>
            </v-col>
          </v-row>

          <v-row>
            <v-col cols="12">
              <v-textarea v-model="form.description" label="应用描述" rows="2" class="mb-4"></v-textarea>
            </v-col>
          </v-row>

          <v-row>
            <v-col cols="12" md="6">
              <v-switch v-model="form.enabled" label="启用应用" color="primary"></v-switch>
            </v-col>
            <v-col cols="12" md="6">
              <v-switch v-model="form.auth.enabled" label="启用认证" color="warning"></v-switch>
            </v-col>
          </v-row>

          <!-- 认证配置 -->
          <v-expand-transition>
            <v-card v-if="form.auth.enabled" elevation="2" class="mb-4">
              <v-card-title class="text-body-1">认证配置</v-card-title>
              <v-card-text>
                <v-text-field v-model="form.auth.token" label="访问令牌 *"
                  :rules="form.auth.enabled ? [rules.required] : []" hint="客户端调用时需要在Authorization头中携带此令牌"
                  persistent-hint></v-text-field>
              </v-card-text>
            </v-card>
          </v-expand-transition>

          <v-row>
            <v-col cols="12" md="6">
              <v-select v-model="form.notifiers" label="通知服务 *" :items="notifiersList" item-title="name"
                item-value="key" multiple chips :rules="[rules.required]"></v-select>
            </v-col>
            <v-col cols="12" md="6">
              <v-select v-model="form.pluginId" label="插件（可选）" :items="pluginOptions" item-title="name" item-value="id"
                clearable hint="配置了插件后，模板将不生效" persistent-hint></v-select>
            </v-col>
          </v-row>

          <v-row>
            <v-col cols="12" md="6">
              <v-select v-model="form.templateId" label="消息模板 *" :items="templatesStore.getTemplateOptions()"
                item-title="title" item-value="value" :rules="form.pluginId ? [] : [rules.required]"
                :disabled="!!form.pluginId" hint="当选择了插件时，模板将被禁用" persistent-hint></v-select>
            </v-col>
          </v-row>

          <v-row>
            <v-col cols="12">
              <v-text-field v-model="form.defaultImage" label="默认图片URL" :rules="[rules.url]" hint="通知中未指定图片时使用的默认图片"
                persistent-hint></v-text-field>
            </v-col>
          </v-row>
        </v-form>
      </v-card-text>

      <v-card-actions class="operation-actions" elevation="12">
        <v-spacer></v-spacer>
        <v-btn text="取消" @click="handleCancel"></v-btn>
        <v-btn variant="flat" color="primary" text="保存" @click="handleSave" :loading="loading"></v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, defineProps, defineEmits } from 'vue'
import { useNotifiersStore } from '@/store/notifiers'
import { useTemplatesStore } from '@/store/templates'
import { usePluginsStore } from '@/store/plugins'
import type { NotificationApp } from '@/common/types'

interface Props {
  modelValue: boolean
  editingApp?: any
  loading?: boolean
}

interface Emits {
  (e: 'save', appData: NotificationApp): void
  (e: 'cancel'): void
}
const dialogVisible = defineModel<boolean>()
const props = withDefaults(defineProps<Props>(), {
  loading: false
})

const emit = defineEmits<Emits>()

const notifiersStore = useNotifiersStore()
const templatesStore = useTemplatesStore()
const pluginsStore = usePluginsStore()

// 表单引用和验证状态
const formRef = ref()


// 表单数据
const form = ref<NotificationApp>({
  appId: '',
  name: '',
  description: '',
  enabled: true,
  notifiers: [],
  templateId: '',
  pluginId: '',
  defaultImage: '',
  auth: {
    enabled: false,
    token: ''
  }
})



// 重置表单
const resetForm = () => {
  form.value = {
    appId: '',
    name: '',
    description: '',
    enabled: true,
    notifiers: [],
    templateId: '',
    pluginId: '',
    defaultImage: '',
    auth: {
      enabled: false,
      token: ''
    },
  }
}

// 验证规则
const rules = {
  required: (value: any) => !!value || '此字段为必填项',
  appId: (value: string) => /^[a-zA-Z0-9_-]+$/.test(value) || 'ID只能包含字母、数字、下划线和连字符',
  url: (value: string) => {
    if (!value) return true // 允许空值
    try {
      new URL(value)
      return true
    } catch {
      return '请输入有效的URL格式'
    }
  }
}


const notifiersList = computed(() => {
  return Object.entries(notifiersStore.notifiers).map(([key, notifier]) => ({
    key,
    name: `${key} (${notifier.type})`,
    enabled: notifier.enabled
  })).filter(item => item.enabled)
})

// 插件选项（仅显示已启用插件）
const pluginOptions = computed(() => {
  return pluginsStore.list.filter(p => p.enabled).map(p => ({ id: p.id, name: `${p.name} (${p.id})` }))
})


// 监听编辑应用变化
watch(dialogVisible, (visible) => {
  const app = props.editingApp
  if (visible && app) {
    form.value = {
      appId: app.appId,
      name: app.name,
      description: app.description || '',
      enabled: app.enabled,
      notifiers: app.notifiers || [],
      templateId: app.templateId,
      pluginId: app.pluginId || '',
      defaultImage: app.defaultImage || '',
      auth: app.auth ? { ...app.auth } : { enabled: false, token: '' },
    }
  } else {
    // 重置表单
    resetForm()
  }
}, { immediate: true })



// 处理保存
const handleSave = async () => {
  const result = await formRef.value.validate()
  if (!result.valid) {
    return
  }
  const appData = {
    ...form.value,
  }
  await emit('save', appData)
}

// 处理取消
const handleCancel = () => {
  emit('cancel')
}

// 处理对话框关闭
const handleDialogClose = (value: boolean) => {
  if (!value) {
    emit('cancel')
  }
}
</script>

<style lang="less" scoped>
@import '@/styles/mix.less';

.dialog-card {
  max-height: 90vh;
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
