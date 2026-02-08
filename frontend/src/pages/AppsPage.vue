<template>
  <!-- 页面加载状态 -->
  <LoadingView :loading="pageLoading" text="正在加载应用数据..." />

  <!-- 页面内容 -->
  <div v-if="!pageLoading">
    <!-- 页面标题 -->
    <v-row>
      <v-col>
        <h1>
          <v-icon icon="mdi-application" class="mr-2"></v-icon>
          通知应用管理
        </h1>
        <p class="text-body-2 text-medium-emphasis mt-1">
          管理通知应用配置，包括认证、通知服务、模板等设置
        </p>
      </v-col>
      <v-col cols="auto">
        <v-btn color="info" variant="outlined" size="small" @click="showHelpDialog = true" class="mr-2">
          <v-icon icon="mdi-help-circle" class="mr-2"></v-icon>
          使用帮助
        </v-btn>
        <v-btn color="primary" variant="flat" size="small" @click="showCreateDialog = true">
          <v-icon icon="mdi-plus" class="mr-2"></v-icon>
          创建应用
        </v-btn>
      </v-col>
    </v-row>

    <!-- 应用列表 -->
    <v-card class="mt-2 bg-transparent">
      <v-card-text>
        <v-row v-if="appsList.length > 0">
          <v-col v-for="app in appsList" :key="app.appId" cols="12" md="6" lg="4">
            <v-card class="app-card" :class="{ 'disabled-app': !app.enabled }" :elevation="app.enabled ? 3 : 0">
              <!-- 应用头部信息 -->
              <v-card-title class="d-flex align-center">
                <v-icon icon="mdi-application" class="mr-2"></v-icon>
                <div class="flex-grow-1">
                  <div class="text-h6">{{ app.name }}</div>
                  <div class="text-caption text-medium-emphasis">{{ app.appId }}</div>
                </div>
                <!-- 状态标识 -->
                <v-chip :color="app.enabled ? 'success' : 'error'" variant="flat" size="small">
                  {{ app.enabled ? '启用' : '禁用' }}
                </v-chip>
              </v-card-title>

              <!-- 应用描述 -->
              <v-card-text v-if="app.description" class="py-2">
                <p class="text-body-2 text-medium-emphasis mb-0">{{ app.description }}</p>
              </v-card-text>

              <!-- 应用详细信息 -->
              <v-card-text class="pt-2">
                <!-- 认证状态 -->
                <div class="mb-3">
                  <v-chip :color="app.auth?.enabled ? 'warning' : 'default'" variant="flat" size="small" class="mr-2">
                    <v-icon icon="mdi-shield-check" class="mr-1" size="small"></v-icon>
                    {{ app.auth?.enabled ? '需要认证' : '无需认证' }}
                  </v-chip>
                </div>

                <!-- 通知服务列表 -->
                <div class="mb-3">
                  <div class="text-body-2 text-medium-emphasis mb-1">通知服务:</div>
                  <v-chip-group v-if="app.notifiers && app.notifiers.length > 0">
                    <v-chip v-for="notifier in app.notifiers" :key="notifier" size="small" variant="outlined">
                      {{ notifier }}
                    </v-chip>
                  </v-chip-group>
                  <span v-else class="text-body-2 text-medium-emphasis">未配置通知服务</span>
                </div>

                <!-- 插件信息（优先展示） -->
                <div v-if="app.pluginId" class="mb-3">
                  <div class="text-body-2 text-medium-emphasis mb-1">插件:</div>
                  <v-chip size="small" variant="tonal" color="secondary">
                    <v-icon icon="mdi-puzzle" class="mr-1" size="small"></v-icon>
                    {{ app.pluginId }}
                  </v-chip>
                </div>

                <!-- 模板信息（仅在未配置插件时显示） -->
                <div v-else-if="app.templateId" class="mb-3">
                  <div class="text-body-2 text-medium-emphasis mb-1">消息模板:</div>
                  <v-chip size="small" variant="tonal" color="primary">
                    <v-icon icon="mdi-file-document" class="mr-1" size="small"></v-icon>
                    {{ app.templateId }}
                  </v-chip>
                </div>
              </v-card-text>

              <!-- 操作按钮 -->
              <v-card-actions>
                <div>
                  <v-btn icon="mdi-pencil" @click="editApp(app)" variant="text" size="small" color="primary">
                    <v-icon></v-icon>
                    <v-tooltip activator="parent" location="bottom">编辑应用</v-tooltip>
                  </v-btn>

                  <v-btn icon="mdi-send" @click="testApp(app)" variant="text" size="small" :disabled="!app.enabled"
                    color="success">
                    <v-icon></v-icon>
                    <v-tooltip activator="parent" location="bottom">测试通知</v-tooltip>
                  </v-btn>

                  <v-btn icon="mdi-content-copy" @click="copyNotifyUrl(app)" variant="text" size="small" color="info">
                    <v-icon></v-icon>
                    <v-tooltip activator="parent" location="bottom">复制通知地址</v-tooltip>
                  </v-btn>
                </div>

                <v-btn icon="mdi-delete" @click="deleteApp(app)" variant="text" size="small" color="error">
                  <v-icon></v-icon>
                  <v-tooltip activator="parent" location="bottom">删除应用</v-tooltip>
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-col>
        </v-row>

        <!-- 无数据状态 -->
        <v-row v-else-if="!appsStore.loading">
          <v-col cols="12" class="text-center py-8">
            <v-icon icon="mdi-application-outline" size="80" class="text-medium-emphasis mb-4"></v-icon>
            <h3 class="text-h6 text-medium-emphasis mb-2">暂无应用</h3>
            <p class="text-body-2 text-medium-emphasis mb-4">点击上方的"创建应用"按钮来创建您的第一个通知应用</p>
          </v-col>
        </v-row>

        <!-- 加载状态 -->
        <v-row v-else>
          <v-col cols="12" class="text-center py-8">
            <v-progress-circular indeterminate color="primary" size="40"></v-progress-circular>
            <p class="text-body-2 text-medium-emphasis mt-4">加载应用列表中...</p>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- 创建/编辑应用对话框 -->
    <AppEditDialog v-model="showCreateDialog" :editing-app="editingApp" :loading="appsStore.loading" @save="saveApp"
      @cancel="closeDialog" />

    <!-- 测试通知对话框 -->
    <TestNotificationDialog v-model="showTestDialog" :app="testingApp" :loading="appsStore.loading"
      @send="sendTestNotification" @cancel="showTestDialog = false" />

    <!-- 帮助说明对话框 -->
    <NotifyHelpDialog v-model="showHelpDialog" />
  </div>
</template>

<script setup lang="ts">
import LoadingView from '@/components/LoadingView.vue'
import AppEditDialog from '@/components/dialog/AppEditDialog.vue'
import TestNotificationDialog from '@/components/dialog/TestNotificationDialog.vue'
import { useAppsStore } from '@/store/apps'
import { useNotifiersStore } from '@/store/notifiers'
import { useTemplatesStore } from '@/store/templates'
import { usePluginsStore } from '@/store/plugins'
import { computed, onMounted, ref } from 'vue'
import { useConfirm } from 'vuetify-use-dialog'
import { copyToClipboard, getCurrentBaseUrl } from '@/common/utils'
import { useToast } from 'vue-toast-notification'
import type { NotificationApp } from '@/common/types'

const appsStore = useAppsStore()
const notifiersStore = useNotifiersStore()
const templatesStore = useTemplatesStore()
const pluginsStore = usePluginsStore()
const createConfirm = useConfirm()
const toast = useToast()

// 页面加载状态
const pageLoading = ref(true)

// 移除了表格配置，改为卡片展示

// 对话框状态
const showCreateDialog = ref(false)
const showTestDialog = ref(false)
const showHelpDialog = ref(false)

// 编辑状态
const editingApp = ref<NotificationApp | null>(null)

const testingApp = ref<NotificationApp | null>(null)

// 计算属性
const appsList = computed(() => {
  return Object.entries(appsStore.apps).map(([_, app]) => app)
})



// 加载数据
const loadApps = async () => {
  await Promise.all([
    appsStore.fetchApps(),
    notifiersStore.fetchNotifiers(),
    templatesStore.fetchTemplates(),
    pluginsStore.fetchList(),
  ])
}

// 编辑应用
const editApp = (app: NotificationApp) => {
  editingApp.value = app
  showCreateDialog.value = true
}

// 测试应用
const testApp = (app: NotificationApp) => {
  testingApp.value = app
  showTestDialog.value = true
}

// 复制通知地址
const copyNotifyUrl = async (app: NotificationApp) => {
  const baseUrl = getCurrentBaseUrl()
  const notifyUrl = `${baseUrl}/api/v1/notify/${app.appId}`

  const success = await copyToClipboard(notifyUrl)
  if (success) {
    toast.success(`已复制通知地址: ${notifyUrl}`)
  } else {
    toast.error('复制失败，请手动复制')
  }
}

// 删除应用
const deleteApp = async (app: NotificationApp) => {
  const isConfirmed = await createConfirm({
    title: '确认删除',
    content: `您确定要删除应用 "${app.name}" 吗？此操作不可撤销。`,
    dialogProps: {
      width: '500px',
    },
    confirmationText: '删除',
    cancellationText: '取消'
  })

  if (isConfirmed) {
    try {
      await appsStore.deleteApp(app.appId)
    } catch (error) {
      console.error('删除应用失败:', error)
    }
  }
}

// 保存应用
const saveApp = async (appData: any) => {
  try {
    if (editingApp.value) {
      await appsStore.updateApp(editingApp.value.appId, appData)
    } else {
      await appsStore.createApp(appData)
    }
    closeDialog()
  } catch (error) {
    console.error('保存应用失败:', error)
  }
}

// 发送测试通知
const sendTestNotification = async (testData: Record<string, any>) => {
  if (!testingApp.value) return

  try {
    await appsStore.sendTestNotification(testingApp.value.appId, testData)
    showTestDialog.value = false
  } catch (error) {
    console.error('发送测试通知失败:', error)
  }
}

// 关闭对话框
const closeDialog = () => {
  showCreateDialog.value = false
  editingApp.value = null
}

// 页面加载时获取数据
onMounted(async () => {
  try {
    await loadApps()
  } catch (error) {
    console.error('页面初始化失败:', error)
  } finally {
    pageLoading.value = false
  }
})
</script>

<style scoped lang="less">
.app-card {
  height: 100%;
  position: relative;
  padding-bottom: 52px;

  .v-card-title {
    padding-bottom: 8px;
  }

  .v-card-text {
    padding-top: 8px;
    padding-bottom: 8px;
  }

  .v-card-actions {
    border-top: 1px solid rgba(var(--v-border-color), 0.2);
    position: absolute;
    bottom: 0;
    z-index: 3;
    padding: 12px 16px;
    justify-content: space-between;
    padding-top: 0;
    width: 100%;
  }

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.12) !important;
  }

  &.disabled-app {
    opacity: 0.7;
    filter: grayscale(0.3);

    &:hover {
      transform: none;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.12) !important;
    }
  }
}

.flex-grow-1 {
  flex: 1 1 auto;
  min-width: 0;
  /* 防止文本溢出 */
}

.text-h6 {
  line-height: 1.2;
  word-break: break-word;
}

.text-caption {
  line-height: 1.2;
  word-break: break-word;
}
</style>
