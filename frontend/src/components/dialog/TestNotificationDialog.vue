<template>
  <v-dialog v-model="dialogVisible" max-width="600px" @update:model-value="handleDialogClose">
    <v-card>
      <v-card-title>测试通知 - {{ app?.name }}</v-card-title>
      <v-card-text>
        <v-textarea v-if="Object.keys(templateTest).length > 0" v-model="templateTest" label="数据" rows="3"
          class="mb-4"></v-textarea>
        <v-form v-else>
          <!-- 固定字段 -->
          <v-text-field v-model="form.title" label="标题" class="mb-4"></v-text-field>
          <v-textarea v-model="form.content" label="内容" rows="3" class="mb-4"></v-textarea>
          <v-text-field v-model="form.image" label="图片" class="mb-4"></v-text-field>
          <v-text-field v-model="form.url" label="链接URL" class="mb-4"></v-text-field>
          <v-text-field v-model="form.targets" label="目标" class="mb-4"></v-text-field>

          <!-- 动态字段 -->
          <v-divider class="my-4"></v-divider>
          <div class="d-flex align-center justify-space-between mb-3">
            <h4>动态字段</h4>
            <v-btn size="small" color="primary" variant="outlined" @click="addCustomField">
              <v-icon left>mdi-plus</v-icon>
              添加字段
            </v-btn>
          </div>

          <div v-for="(field, index) in customFields" :key="index" class="mb-3">
            <v-row>
              <v-col cols="4">
                <v-text-field v-model="field.name" label="字段名" density="compact" variant="outlined"></v-text-field>
              </v-col>
              <v-col cols="7">
                <v-text-field v-model="field.value" label="字段值" density="compact" variant="outlined"></v-text-field>
              </v-col>
              <v-col cols="1" class="d-flex align-center">
                <v-btn size="small" color="error" variant="text" icon="mdi-delete"
                  @click="removeCustomField(index)"></v-btn>
              </v-col>
            </v-row>
          </div>

          <div v-if="customFields.length === 0" class="text-center text-grey mb-3">
            <p>暂无动态字段，点击上方按钮添加</p>
          </div>
        </v-form>


      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn text="取消" @click="handleCancel"></v-btn>
        <v-btn variant="flat" color="primary" text="发送" @click="handleSend" :loading="loading"></v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import type { NotificationApp } from '@/common/types'
import { usePluginsStore } from '@/store/plugins'
import { ref, watch, defineProps, defineEmits } from 'vue'
const pluginsStore = usePluginsStore()

interface CustomField {
  name: string
  value: string
}

interface TestNotificationData {
  title: string
  content: string
  url: string
  image: string
  targets: string
  [key: string]: any
}

interface Props {
  app?: NotificationApp | null
  loading?: boolean
}

interface Emits {
  (e: 'send', data: TestNotificationData | Record<string, any>): void
  (e: 'cancel'): void
}


const props = withDefaults(defineProps<Props>(), {
  loading: false
})

const plugin = computed(() => {
  if (props.app?.pluginId) {
    return pluginsStore.getById(props.app.pluginId)
  }
  return null
})

const dialogVisible = defineModel<boolean>()

const emit = defineEmits<Emits>()

// 表单数据
const form = ref<TestNotificationData>({
  title: '测试通知',
  content: '这是一条测试通知消息',
  url: '',
  image: '',
  targets: ''
})

const templateTest = ref<string>('')

// 动态字段
const customFields = ref<CustomField[]>([])

// 添加自定义字段
const addCustomField = () => {
  customFields.value.push({
    name: '',
    value: ''
  })
}

// 删除自定义字段
const removeCustomField = (index: number) => {
  customFields.value.splice(index, 1)
}

// 监听对话框打开时重置表单
watch(dialogVisible, (visible) => {
  if (visible) {
    form.value = {
      title: '测试通知',
      content: '这是一条测试通知消息',
      url: '',
      image: '',
      targets: ''
    }
    console.debug(plugin.value)
    if (plugin.value) {
      templateTest.value = JSON.stringify(plugin.value?.test_data || {}, null, 2)
    } else {
      templateTest.value = ''
    }
    customFields.value = []
  }
})

// 处理发送
const handleSend = () => {
  if (plugin.value && templateTest.value) {
    const data = JSON.parse(templateTest.value)
    emit('send', data)
    return
  } else if (form.value) {
    const data = {
      ...form.value,
      ...customFields.value.filter(field => field.name.trim() !== '' || field.value.trim() !== '').reduce((acc, field) => {
        acc[field.name] = field.value
        return acc
      }, {} as Record<string, string>)
    }
    emit('send', data)
  }
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
