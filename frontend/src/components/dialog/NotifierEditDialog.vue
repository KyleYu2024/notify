<template>
  <v-dialog v-model="showDialog" max-width="500px">
    <v-card class="edit-dialog">
      <v-card-title>{{ title }}</v-card-title>
      <v-card-text>
        <v-form ref="formRef" v-model="formValid">
          <!-- 新增通知服务时显示名称输入框 -->
          <v-text-field v-if="!props.notifierKey" v-model="form.name" label="通知服务名称" placeholder="请输入通知服务名称"
            :rules="nameRules" class="mb-4" required></v-text-field>

          <!-- 新增通知服务时显示类型选择器 -->
          <v-select v-if="!props.notifierKey" v-model="form.type" :items="notifierTypeOptions" label="通知服务类型"
            :rules="typeRules" class="mb-4" required></v-select>

          <v-switch v-model="form.enabled" label="启用通知服务" color="primary" class="mb-4"></v-switch>

          <!-- 动态组件渲染不同的通知服务配置 -->
          <div v-if="currentNotifierType && configComponent">
            <component :is="configComponent" v-model="form.config" @update:modelValue="handleConfigUpdate" />
          </div>
        </v-form>
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
import { useNotifiersStore, type INotifierInstance } from '@/store/notifiers'
import { ref, computed, watch } from 'vue'
import { notifierConfigComponents, type NotifierConfigType } from './notifier-configs'
import { NotifierTypeMap, notifierTypeOptions } from '@/common/types'

const notifierStore = useNotifiersStore()

interface Props {
  notifierKey?: string
  loading?: boolean
}
const showDialog = defineModel<boolean>()

interface Emits {
  (e: 'save', key: string, config: any): void
  (e: 'cancel'): void
}

const title = computed(() => {
  return props.notifierKey ? '编辑通知服务 - ' + props.notifierKey : '新增通知服务'
})

const props = withDefaults(defineProps<Props>(), {
  notifierKey: '',
  loading: false
})

// 表单数据
const form = ref<INotifierInstance & { name: string }>({
  name: '',
  enabled: true,
  type: NotifierTypeMap.wechatWorkAPPBot,
  config: {}
})

const notifier = computed(() => {
  if (!props.notifierKey) return null
  return notifierStore.notifiers[props.notifierKey] || null
})

const notifierType = computed(() => {
  return notifier.value?.type || ''
})

// 当前使用的通知服务类型（编辑时使用existing，新增时使用form中的type）
const currentNotifierType = computed(() => {
  return props.notifierKey ? notifierType.value : form.value.type
})

// 当前配置组件
const configComponent = computed(() => {
  return notifierConfigComponents[currentNotifierType.value as NotifierConfigType] || null
})



// 验证规则
const nameRules = [
  (v: string) => !!v || '通知服务名称不能为空',
  (v: string) => (v && v.length >= 2) || '通知服务名称至少2个字符',
  (v: string) => (v && v.length <= 50) || '通知服务名称最多50个字符',
  (v: string) => {
    if (props.notifierKey) return true // 编辑时不检查重复
    return !notifierStore.notifiers[v] || '通知服务名称已存在'
  }
]

const typeRules = [
  (v: string) => !!v || '请选择通知服务类型'
]

const emit = defineEmits<Emits>()

watch(showDialog, (val) => {
  if (val) {
    if (notifier.value) {
      // 编辑模式
      form.value = { ...notifier.value, name: props.notifierKey }
    } else {
      // 新增模式，默认选择微信
      form.value = {
        name: '',
        enabled: true,
        type: NotifierTypeMap.wechatWorkAPPBot,
        config: {}
      }
    }
  }
}, { immediate: true })

// 内部状态
const formValid = ref(false)
const formRef = ref()


// 处理配置更新
const handleConfigUpdate = (newConfig: any) => {
  form.value.config = { ...newConfig }
}

// 处理保存
const handleSave = async () => {
  const res = await formRef.value.validate()
  if (!res.valid) return

  const saveKey = props.notifierKey || form.value.name
  const saveData = {
    enabled: form.value.enabled,
    type: form.value.type,
    config: form.value.config
  }

  await emit('save', saveKey, saveData)
  showDialog.value = false
}

// 处理取消
const handleCancel = () => {
  emit('cancel')
}
</script>

<style lang="less" scoped>
@import '@/styles/mix.less';

.edit-dialog {
  max-height: 80vh;
  background-color: rgba(var(--v-theme-surface), 0.5);
  backdrop-filter: blur(20px);
}

.operation-actions {
  position: sticky;
  bottom: 0;
  z-index: 3;
  background-color: transparent;
}
</style>
