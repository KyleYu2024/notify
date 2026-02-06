<template>
  <div>
    <!-- 飞书应用配置 -->
    <v-text-field v-model="config.app_id" label="应用 ID *" hint="飞书应用的 App ID" persistent-hint :rules="[rules.required]"
      class="mb-4" @input="handleConfigChange"></v-text-field>

    <v-text-field v-model="config.app_secret" label="应用密钥 *" hint="飞书应用的 App Secret" persistent-hint
      :rules="[rules.required]" class="mb-4" @input="handleConfigChange"></v-text-field>

    <v-text-field v-model="config.targets" label="目标用户" hint="接收者ID，支持多种类型，多个用逗号分隔（可选）" persistent-hint class="mb-4"
      @input="handleConfigChange"></v-text-field>

    <v-text-field v-model="config.proxy" label="代理服务器" hint="可选，格式: http://proxy.example.com:8080" persistent-hint
      class="mb-4" @input="handleConfigChange"></v-text-field>

    <v-alert type="info" variant="tonal" class="mb-4">
      <div class="text-body-2">
        <strong>配置说明：</strong><br><br>
        <strong>1. 创建飞书应用</strong><br>
        • 访问飞书开放平台：<a href="https://open.feishu.cn/" target="_blank">https://open.feishu.cn/</a><br>
        • 创建企业自建应用，获取 <strong>App ID</strong> 和 <strong>App Secret</strong><br>
        • 配置应用权限：im:message（发送消息权限）<br>
        • 发布应用并获得管理员同意<br><br>

        <strong>2. 配置接收者</strong><br>
        目标用户支持多种ID格式（系统会自动识别类型）：<br>
        • <strong>open_id</strong>: 以 "ou_" 开头（如：ou_7d8a6e6...）- 用户在应用中的身份<br>
        • <strong>union_id</strong>: 以 "on_" 开头（如：on_94648c...）- 用户在开发商下的统一身份<br>
        • <strong>user_id</strong>: 用户ID（如：4d7a3c6g）- 用户在租户内的身份<br>
        • <strong>chat_id</strong>: 以 "oc_" 开头（如：oc_234jkl...）- 群组ID<br>
        • <strong>email</strong>: 用户邮箱地址（如：user@example.com）<br>
        <strong><a
            href="https://open.feishu.cn/document/faq/trouble-shooting/how-to-obtain-openid">获取openid</a></strong>
        <br>
        <strong><a
            href="https://open.feishu.cn/document/server-docs/group/chat/chat-id-description">获取chat_id</a></strong>
        <br><br>
        <strong>3. 测试发送</strong><br>
        配置完成后可通过API测试接口发送测试消息验证配置是否正确。
      </div>
    </v-alert>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { FeishuConfig } from '@/common/types'

interface Props {
  modelValue: Partial<FeishuConfig>
}

interface Emits {
  (e: 'update:modelValue', value: Partial<FeishuConfig>): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 内部配置状态
const config = ref<Partial<FeishuConfig>>({
  app_id: '',
  app_secret: '',
  proxy: '',
  targets: '',
  ...props.modelValue
})

// 验证规则
const rules = {
  required: (value: any) => !!value || '此字段为必填项'
}

// 监听 props 变化
watch(() => props.modelValue, (newValue) => {
  config.value = {
    app_id: '',
    app_secret: '',
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
