<template>
  <div>
    <v-text-field v-model="config.key" label="群机器人 Key *" :rules="[rules.required]" class="mb-4"
      hint="例如: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" persistent-hint @input="handleConfigChange">
    </v-text-field>

    <v-text-field v-model="config.proxy" label="代理服务器" hint="可选，格式: http://proxy.example.com:8080" persistent-hint
      class="mb-4" @input="handleConfigChange">
    </v-text-field>

    <v-alert type="info" variant="tonal" class="mb-4">
      <div class="text-body-2">
        <strong>如何获取 Key：</strong><br>
        1. 在企业微信群聊中，点击右上角"..." → "群机器人"<br>
        2. 点击"添加机器人" → "自定义机器人"<br>
        3. 配置机器人名称和头像，点击"添加"<br>
        4. 从生成的 Webhook URL 中复制 key 参数值<br>
        &nbsp;&nbsp;&nbsp;URL 格式：https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=<strong>xxxxxxxx</strong>
      </div>
    </v-alert>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { WechatWorkWebhookConfig } from '@/common/types'

interface Props {
  modelValue: Partial<WechatWorkWebhookConfig>
}

interface Emits {
  (e: 'update:modelValue', value: Partial<WechatWorkWebhookConfig>): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 内部配置状态
const config = ref<Partial<WechatWorkWebhookConfig>>({
  key: '',
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
    key: '',
    proxy: '',
    ...newValue
  }
}, { deep: true })

// 配置变化处理
const handleConfigChange = () => {
  emit('update:modelValue', { ...config.value })
}
</script>
