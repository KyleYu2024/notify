<template>
  <v-dialog v-model="dialogVisible" max-width="800px">
    <v-card class="notify-help-dialog">
      <v-card-title class="d-flex align-center">
        <v-icon icon="mdi-help-circle" class="mr-2"></v-icon>
        通知API使用说明
      </v-card-title>

      <v-card-text>
        <v-tabs v-model="helpTab" color="primary" class="mb-4">
          <v-tab value="endpoint">调用地址</v-tab>
          <v-tab value="auth">认证方式</v-tab>
          <v-tab value="params">请求参数</v-tab>
          <v-tab value="examples">示例代码</v-tab>
        </v-tabs>

        <v-tabs-window v-model="helpTab">
          <!-- 调用地址 -->
          <v-tabs-window-item value="endpoint">
            <h4 class="mb-2">API调用地址</h4>
            <p class="text-body-2 mb-3">每个应用都有唯一的通知API地址：</p>
            <CodeBlock :code="`${getCurrentBaseUrl()}/api/v1/notify/应用ID`" />
            <p class="text-body-2 mt-3 mb-2">请求方式：<strong>POST</strong></p>
            <p class="text-body-2">Content-Type：<strong>application/json</strong></p>
          </v-tabs-window-item>

          <!-- 认证方式 -->
          <v-tabs-window-item value="auth">
            <h4 class="mb-2">认证方式</h4>
            <div class="mb-3">
              <h5 class="mb-2">启用认证的应用</h5>
              <p class="text-body-2">需要在请求头中携带认证令牌：</p>
              <CodeBlock code="Authorization: Bearer 您的访问令牌" />
              <CodeBlock code="http://localhost:5174/api/v1/notify/appid?token=您的访问令牌" />
            </div>
            <div>
              <h5 class="mb-2">未启用认证的应用</h5>
              <p class="text-body-2">可以直接发送请求，无需认证头。</p>
            </div>
          </v-tabs-window-item>

          <!-- 请求参数 -->
          <v-tabs-window-item value="params">
            <h4 class="mb-2">请求参数</h4>
            <p class="text-body-2 mb-3">支持以下参数（具体支持的参数取决于您选择的消息模板）：</p>

            <v-table density="compact">
              <thead>
                <tr>
                  <th>参数名</th>
                  <th>类型</th>
                  <th>说明</th>
                  <th>示例</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td><code>title</code></td>
                  <td>string</td>
                  <td>通知标题</td>
                  <td>"系统告警"</td>
                </tr>
                <tr>
                  <td><code>content</code></td>
                  <td>string</td>
                  <td>通知内容</td>
                  <td>"服务器CPU使用率过高"</td>
                </tr>
                <tr>
                  <td><code>image</code></td>
                  <td>string</td>
                  <td>图片URL</td>
                  <td>"https://example.com/image.png"</td>
                </tr>
                <tr>
                  <td><code>url</code></td>
                  <td>string</td>
                  <td>跳转链接</td>
                  <td>"https://example.com/dashboard"</td>
                </tr>
              </tbody>
            </v-table>
          </v-tabs-window-item>

          <!-- 示例代码 -->
          <v-tabs-window-item value="examples">
            <h4 class="mb-2">示例代码</h4>

            <h5 class="mt-4 mb-2">cURL 示例</h5>
            <CodeBlock :code="curlExample" language="bash" />

            <h5 class="mt-4 mb-2">JavaScript 示例</h5>
            <CodeBlock :code="javascriptExample" language="javascript" />

            <h5 class="mt-4 mb-2">Python 示例</h5>
            <CodeBlock :code="pythonExample" language="python" />
          </v-tabs-window-item>
        </v-tabs-window>
      </v-card-text>

      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn text="关闭" @click="handleClose"></v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { getCurrentBaseUrl } from '@/common/utils'

interface Props {
  modelValue: boolean
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 对话框显示状态
const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

// 帮助对话框Tab状态
const helpTab = ref('endpoint')

// 示例代码
const curlExample = computed(() => `curl -X POST ${getCurrentBaseUrl()}/api/v1/notify/your-app-id \\
  -H "Content-Type: application/json" \\
  -H "Authorization: Bearer your-token" \\
  -d '{
    "title": "系统告警",
    "content": "服务器CPU使用率过高",
    "level": "warning"
  }'`)

const javascriptExample = computed(() => `fetch('${getCurrentBaseUrl()}/api/v1/notify/your-app-id', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer your-token'
  },
  body: JSON.stringify({
    title: '系统告警',
    content: '服务器CPU使用率过高',
    level: 'warning'
  })
})`)

const pythonExample = computed(() => `import requests

url = '${getCurrentBaseUrl()}/api/v1/notify/your-app-id'
headers = {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer your-token'
}
data = {
    'title': '系统告警',
    'content': '服务器CPU使用率过高',
    'level': 'warning'
}

response = requests.post(url, headers=headers, json=data)`)

// 处理关闭
const handleClose = () => {
  dialogVisible.value = false
}
</script>



<style lang="less" scoped>
@import '@/styles/mix.less';

.notify-help-dialog {
  max-height: 80vh;
  overflow-y: auto;
  .scrollbar();
  background-color: rgba(var(--v-theme-surface), 0.5);
  backdrop-filter: blur(10px);
}
</style>
