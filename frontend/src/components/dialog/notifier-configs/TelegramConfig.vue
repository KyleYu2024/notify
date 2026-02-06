<template>
  <div>
    <v-text-field v-model="config.bot_token" label="Bot Token *" :rules="[rules.required]" hint="从 @BotFather 获取"
      persistent-hint class="mb-4" @input="handleConfigChange"></v-text-field>

    <v-text-field v-model="config.chat_id" label="Chat ID *" :rules="[rules.required]" hint="群组或频道ID，可以是负数, 多个用逗号分隔"
      persistent-hint class="mb-4" @input="handleConfigChange"></v-text-field>

    <v-text-field v-model="config.proxy" label="代理服务器" hint="可选，格式: http://proxy.example.com:8080" persistent-hint
      class="mb-4" @input="handleConfigChange"></v-text-field>

    <v-alert type="info" variant="tonal" class="mb-4">
      <div class="text-body-2">
        <strong>如何获取配置信息：</strong><br>
        <strong>获取 Bot Token：</strong><br>
        1. 在 Telegram 中搜索 @BotFather<br>
        2. 发送 /start 开始对话<br>
        3. 发送 /newbot 创建新机器人<br>
        4. 按提示设置机器人名称和用户名<br>
        5. 获得形如 <strong>1234567890:ABCDEFghijklmnopQRSTUVwxyz</strong> 的 token<br><br>
        <strong>获取 Chat ID：</strong><br>
        1. 将机器人添加到目标群组或频道<br>
        2. 在浏览器访问：https://api.telegram.org/bot&lt;TOKEN&gt;/getUpdates<br>
        3. 从返回结果中找到 chat.id 字段值
      </div>
    </v-alert>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { TelegramConfig } from '@/common/types'

interface Props {
  modelValue: Partial<TelegramConfig>
}

interface Emits {
  (e: 'update:modelValue', value: Partial<TelegramConfig>): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 内部配置状态
const config = ref<Partial<TelegramConfig>>({
  bot_token: '',
  chat_id: '',
  proxy: '',
  ...props.modelValue
})

// 验证规则
const rules = {
  required: (value: any) => !!value || '此字段为必填项'
}

// 监听 props 变化
watch(() => props.modelValue, (newValue) => {
  config.value = {
    bot_token: '',
    chat_id: '',
    proxy: '',
    ...newValue
  }
}, { deep: true })

// 配置变化处理
const handleConfigChange = () => {
  emit('update:modelValue', { ...config.value })
}
</script>
