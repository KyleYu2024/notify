<template>
  <LoadingView :loading="pageLoading" text="正在加载数据..." v-if="pageLoading" />

  <div v-else>
    <v-row>
      <v-col>
        <h1 class="text-h6">
          <v-icon icon="mdi-view-dashboard" class="mr-2"></v-icon>
          系统仪表板
        </h1>
        <p class="text-body-2 text-medium-emphasis mt-1">
          通知管理系统概览，查看系统状态和统计信息
        </p>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12" sm="6" md="4">
        <v-card class="stat-card text-white" color="success" variant="flat" @click="navigateToNotifiers">
          <v-card-text class="d-flex align-center">
            <div class="flex-grow-1">
              <div class="text-h4 font-weight-bold">{{ notifiersCount }}</div>
              <div class="text-body-2">通知服务</div>
            </div>
            <v-icon icon="mdi-bell-cog" size="48" class="opacity-80"></v-icon>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" sm="6" md="4">
        <v-card class="stat-card text-white" color="primary" variant="flat" @click="navigateToApps">
          <v-card-text class="d-flex align-center">
            <div class="flex-grow-1">
              <div class="text-h4 font-weight-bold">{{ appsCount }}<span class="text-body-2 ml-1">已启用({{
                enabledAppsCount }})</span></div>
              <div class="text-body-2">通知应用</div>
            </div>
            <v-icon icon="mdi-application" size="48" class="opacity-80"></v-icon>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" sm="6" md="4">
        <v-card class="stat-card text-white" color="warning" variant="flat" @click="navigateToTemplates">
          <v-card-text class="d-flex align-center">
            <div class="flex-grow-1">
              <div class="text-h4 font-weight-bold">{{ templatesCount }}</div>
              <div class="text-body-2">消息模板</div>
            </div>
            <v-icon icon="mdi-file-document" size="48" class="opacity-80"></v-icon>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12" md="6">
        <v-card class="mb-4 h-100" elevation="0" border rounded="xl">
          <v-card-title class="d-flex align-center">
            <v-icon icon="mdi-heart-pulse" class="mr-2" color="error"></v-icon>
            系统状态
          </v-card-title>
          <v-card-text>
            <v-list density="comfortable">
              <v-list-item>
                <template v-slot:prepend>
                  <v-icon :icon="systemStatus.backend ? 'mdi-check-circle' : 'mdi-alert-circle'"
                    :color="systemStatus.backend ? 'success' : 'error'"></v-icon>
                </template>
                <v-list-item-title>后端服务</v-list-item-title>
                <v-list-item-subtitle>
                  {{ systemStatus.backend ? '运行正常' : '连接失败' }}
                </v-list-item-subtitle>
              </v-list-item>

              <v-list-item>
                <template v-slot:prepend>
                  <v-icon :icon="systemStatus.auth ? 'mdi-shield-check' : 'mdi-shield-off'"
                    :color="systemStatus.auth ? 'warning' : 'success'"></v-icon>
                </template>
                <v-list-item-title>认证状态</v-list-item-title>
                <v-list-item-subtitle>
                  {{ systemStatus.auth ? '需要认证' : '无需认证' }}
                </v-list-item-subtitle>
              </v-list-item>

              <v-list-item>
                <template v-slot:prepend>
                  <v-icon :icon="systemStatus.config ? 'mdi-cog' : 'mdi-cog-off'"
                    :color="systemStatus.config ? 'success' : 'error'"></v-icon>
                </template>
                <v-list-item-title>配置状态</v-list-item-title>
                <v-list-item-subtitle>
                  {{ systemStatus.config ? '配置正常' : '配置异常' }}
                </v-list-item-subtitle>
              </v-list-item>
            </v-list>

            <v-btn variant="tonal" size="small" @click="checkSystemStatus" :loading="statusLoading" class="mt-4" color="primary">
              <v-icon icon="mdi-refresh" class="mr-2"></v-icon>
              刷新状态
            </v-btn>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="6" v-if="!display.mobile.value">
        <v-card class="mb-4 h-100" elevation="0" border rounded="xl">
          <v-card-title class="d-flex align-center">
            <v-icon icon="mdi-information" class="mr-2" color="info"></v-icon>
            系统信息
          </v-card-title>
          <v-card-text>
            <v-list density="comfortable">
              <v-list-item>
                <template v-slot:prepend>
                  <v-icon icon="mdi-clock" color="primary"></v-icon>
                </template>
                <v-list-item-title>当前时间</v-list-item-title>
                <v-list-item-subtitle>{{ currentTime }}</v-list-item-subtitle>
              </v-list-item>

              <v-list-item>
                <template v-slot:prepend>
                  <v-icon icon="mdi-server" color="info"></v-icon>
                </template>
                <v-list-item-title>当前版本</v-list-item-title>
                <v-list-item-subtitle>{{ authStore.version || 'v0.2.3' }}</v-list-item-subtitle>
              </v-list-item>
              
              <v-list-item>
                <template v-slot:prepend>
                  <v-icon icon="mdi-github" color="grey"></v-icon>
                </template>
                <v-list-item-title>项目地址</v-list-item-title>
                <v-list-item-subtitle>
                  <a href="https://github.com/kyleyu2024/notify" target="_blank" class="text-decoration-none text-primary">GitHub</a>
                </v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useDisplay } from 'vuetify'
import { useAuthStore } from '@/store/auth'
import { useAppsStore } from '@/store/apps'
import { useNotifiersStore } from '@/store/notifiers'
import { useTemplatesStore } from '@/store/templates'
import LoadingView from '@/components/LoadingView.vue'

const router = useRouter()
const authStore = useAuthStore()
const appsStore = useAppsStore()
const notifiersStore = useNotifiersStore()
const templatesStore = useTemplatesStore()
const display = useDisplay()

// 页面加载状态
const pageLoading = ref(true)

// 系统状态
const systemStatus = ref({
  backend: false,
  auth: false,
  config: false
})

const statusLoading = ref(false)

// 当前时间
const currentTime = ref('')
let timeInterval: any = null

// 计算属性
const appsCount = computed(() => Object.keys(appsStore.apps).length)
const notifiersCount = computed(() => Object.keys(notifiersStore.notifiers).length)
const templatesCount = computed(() => Object.keys(templatesStore.templates).length)
const enabledAppsCount = computed(() => {
  return Object.values(appsStore.apps).filter(app => app.enabled).length
})

// 更新当前时间
const updateCurrentTime = () => {
  currentTime.value = new Date().toLocaleString('zh-CN', {
    hour12: false,
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// 检查系统状态
const checkSystemStatus = async () => {
  statusLoading.value = true
  try {
    const status = await authStore.checkSystemStatus()
    systemStatus.value = {
      backend: true,
      auth: status.adminAuthRequired || false,
      config: true
    }
  } catch (error) {
    systemStatus.value = {
      backend: false,
      auth: false,
      config: false
    }
  } finally {
    statusLoading.value = false
  }
}

// 导航到不同页面
const navigateToApps = () => {
  router.push('/apps')
}

const navigateToNotifiers = () => {
  router.push('/notifiers')
}

const navigateToTemplates = () => {
  router.push('/templates')
}

// 加载数据
const loadData = async () => {
  try {
    await Promise.all([
      appsStore.fetchApps(),
      notifiersStore.fetchNotifiers(),
      templatesStore.fetchTemplates()
    ])
  } catch (e) {
    console.warn('部分数据加载失败', e)
  }
}

// 页面加载
onMounted(async () => {
  try {
    updateCurrentTime()
    checkSystemStatus()
    // 每秒更新一次时间，体验更好
    timeInterval = setInterval(updateCurrentTime, 1000)
    await loadData()
  } catch (error) {
    console.error('页面初始化失败:', error)
  } finally {
    pageLoading.value = false
  }
})

onUnmounted(() => {
  if (timeInterval) clearInterval(timeInterval)
})
</script>

<style scoped lang="less">
.stat-card {
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
  /* 调整圆角为 16px，更现代的大圆角风格 */
  border-radius: 16px !important;
  min-height: 120px;
  display: flex;
  align-items: center;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 16px rgba(0,0,0,0.1) !important; /* 悬停时阴影也稍微扩散一点，配合大圆角 */
  }
}
</style>