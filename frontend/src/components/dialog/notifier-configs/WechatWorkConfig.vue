<template>
  <div>
    <v-text-field v-model="config.corp_id" label="企业ID *" :rules="[rules.required]" class="mb-4"
      @input="handleConfigChange"></v-text-field>

    <v-text-field v-model="config.agent_id" label="应用ID *" :rules="[rules.required]" class="mb-4"
      @input="handleConfigChange"></v-text-field>

    <v-text-field v-model="config.secret" label="应用密钥 *" :rules="[rules.required]" class="mb-4"
      @input="handleConfigChange"></v-text-field>

    <v-text-field v-model="config.targets" label="目标 *" hint="用户id，多个用逗号分隔" class="mb-4" persistent-hint
      @input="handleConfigChange"></v-text-field>

    <v-text-field v-model="config.proxy" label="代理服务器" hint="可选，格式: http://proxy.example.com:8080" persistent-hint
      class="mb-4" @input="handleConfigChange"></v-text-field>

    <v-alert type="info" variant="tonal" class="mb-4">
      <div class="text-body-2">
        <strong>如何获取配置信息：</strong><br>
        1. 登录企业微信管理后台：https://work.weixin.qq.com/<br>
        2. 进入"应用管理" → "应用" → "创建应用"<br>
        3. 创建自建应用后，可以获得：<br>
        &nbsp;&nbsp;&nbsp;• <strong>企业ID</strong>：在"我的企业" → "企业信息"中查看<br>
        &nbsp;&nbsp;&nbsp;• <strong>应用ID</strong>：在应用详情页的 AgentId<br>
        &nbsp;&nbsp;&nbsp;• <strong>应用密钥</strong>：在应用详情页的 Secret<br>
        4. 目标可以是用户ID、部门ID或标签ID，多个用逗号分隔<br>
        &nbsp;&nbsp;&nbsp;用户ID格式：@userId，部门ID格式：@departmentId
      </div>
    </v-alert>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { WechatWorkConfig } from '@/common/types'

interface Props {
  modelValue: Partial<WechatWorkConfig>
}

interface Emits {
  (e: 'update:modelValue', value: Partial<WechatWorkConfig>): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 内部配置状态
const config = ref<Partial<WechatWorkConfig>>({
  corp_id: '',
  agent_id: '',
  secret: '',
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
    corp_id: '',
    agent_id: '',
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
