<template>
  <div>
    <v-text-field v-model="config.access_token" label="Access Token *" :rules="[rules.required]" class="mb-4"
      @input="handleConfigChange"></v-text-field>

    <v-text-field v-model="config.secret" label="签名密钥 *" class="mb-4" @input="handleConfigChange"></v-text-field>
    <v-text-field v-model="config.targets" label="目标 *" hint="用户手机号或用户id，多个用逗号分隔" class="mb-4" persistent-hint
      @input="handleConfigChange"></v-text-field>
    <v-text-field v-model="config.proxy" label="代理服务器" hint="可选，格式: http://proxy.example.com:8080" persistent-hint
      class="mb-4" @input="handleConfigChange"></v-text-field>

    <v-alert type="info" variant="tonal" class="mb-4">
      <div class="text-body-2">
        <strong>如何获取配置信息：</strong><br>
        1. 打开钉钉PC客户端，进入要发送消息的群聊<br>
        2. 点击群设置 → 智能群助手 → 添加机器人 → 自定义<br>
        3. 配置机器人名称，安全设置选择"加签"<br>
        4. 记录生成的 <strong>Webhook地址</strong> 中的 access_token 参数<br>
        5. 记录生成的 <strong>签名密钥</strong><br>
        6. 目标填写要@的用户手机号，多个用逗号分隔（可选）
      </div>
    </v-alert>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { DingTalkConfig } from '@/common/types'

interface Props {
  modelValue: Partial<DingTalkConfig>
}

interface Emits {
  (e: 'update:modelValue', value: Partial<DingTalkConfig>): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 内部配置状态
const config = ref<Partial<DingTalkConfig>>({
  access_token: '',
  secret: '',
  proxy: '',
  targets: '',
  ...props.modelValue
})

// 验证规则
const rules = {
  required: (value: any) => !!value || '此字段为必填项',
  url: (value: string) => {
    if (!value) return true
    try {
      new URL(value)
      return true
    } catch {
      return '请输入有效的URL'
    }
  }
}

// 监听 props 变化
watch(() => props.modelValue, (newValue) => {
  config.value = {
    access_token: '',
    secret: '',
    proxy: '',
    targets: '',
    ...newValue
  }
}, { deep: true })

// 配置变化处理
const handleConfigChange = () => {
  emit('update:modelValue', { ...config.value })
}
</script>
