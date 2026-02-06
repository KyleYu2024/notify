<template>
  <v-container class="fill-height">
    <!-- 页面加载状态 -->
    <LoadingView :loading="pageLoading" text="正在检查系统状态..." v-if="pageLoading" />

    <!-- 登录页面内容 -->
    <v-row align="center" justify="center" v-else>
      <v-col cols="12" sm="8" md="4">
        <v-card class="pa-8" elevation="8">
          <v-card-title class="text-center mb-6">
            <v-icon icon="mdi-bell-ring" size="48" class="mb-4" color="primary"></v-icon>
            <h2 class="text-h4">通知管理系统</h2>
          </v-card-title>

          <!-- 系统状态检查 -->
          <v-alert v-if="!authStore.isAuthRequired" type="info" class="mb-4" variant="tonal">
            <v-icon icon="mdi-information"></v-icon>
            系统当前无需认证，点击下方按钮直接进入
          </v-alert>

          <v-card-text v-if="!authStore.isAuthRequired">
            <v-btn @click="handleLogin" color="primary" size="large" block :loading="authStore.loading">
              <v-icon icon="mdi-login" class="mr-2"></v-icon>
              进入系统
            </v-btn>
          </v-card-text>

          <v-card-text v-if="authStore.isAuthRequired">
            <v-form @submit.prevent="handleLogin">
              <v-text-field v-model="form.username" label="用户名" prepend-inner-icon="mdi-account" variant="outlined"
                :disabled="!authStore.isAuthRequired" :rules="authStore.isAuthRequired ? [rules.required] : []"
                class="mb-4"></v-text-field>

              <v-text-field v-model="form.password" label="密码" type="password" prepend-inner-icon="mdi-lock"
                variant="outlined" :disabled="!authStore.isAuthRequired"
                :rules="authStore.isAuthRequired ? [rules.required] : []" class="mb-6"></v-text-field>

              <v-btn type="submit" color="primary" size="large" block :loading="authStore.loading"
                :disabled="authStore.isAuthRequired && (!form.username || !form.password)">
                <v-icon icon="mdi-login" class="mr-2"></v-icon>
                登录
              </v-btn>
            </v-form>
          </v-card-text>

          <!-- 系统信息 -->
          <v-card-actions class="justify-center">
            <v-btn variant="text" size="small" @click="checkSystemStatus" :loading="checking">
              <v-icon icon="mdi-refresh" class="mr-1"></v-icon>
              检查系统状态
            </v-btn>
          </v-card-actions>

          <!-- 连接状态指示器 -->
          <v-card-text class="text-center">
            <v-chip :color="systemStatus === 'online' ? 'success' : 'error'" variant="flat" size="small">
              <v-icon :icon="systemStatus === 'online' ? 'mdi-wifi' : 'mdi-wifi-off'" class="mr-1"></v-icon>
              {{ systemStatus === 'online' ? '系统在线' : '连接失败' }}
            </v-chip>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import LoadingView from '@/components/LoadingView.vue'

const router = useRouter()
const authStore = useAuthStore()

// 页面加载状态
const pageLoading = ref(true)

// 表单数据
const form = ref({
  username: '',
  password: ''
})

// 系统状态
const systemStatus = ref<'online' | 'offline'>('offline')
const checking = ref(false)

// 表单验证规则
const rules = {
  required: (value: string) => !!value || '此字段为必填项'
}

// 检查系统状态
const checkSystemStatus = async () => {
  checking.value = true
  try {
    await authStore.checkSystemStatus()
    systemStatus.value = 'online'
  } catch (error) {
    systemStatus.value = 'offline'
  } finally {
    checking.value = false
  }
}

// 处理登录
const handleLogin = async () => {
  try {
    if (authStore.isAuthRequired) {
      await authStore.login(form.value.username, form.value.password)
    } else {
      // 无需认证的情况下直接设置为已认证
      authStore.isAuthenticated = true
    }
    // 登录成功后跳转到主页
    router.push('/')
  } catch (error) {
    console.error('登录失败:', error)
  }
}

// 页面加载时检查系统状态
onMounted(async () => {
  try {
    await checkSystemStatus()
    // 检查是否需要认证
    if (!authStore.isAuthRequired) {
      // 如果不需要认证，直接设置为已认证状态
      authStore.isAuthenticated = true
      router.replace('/')
      authStore.username = '未认证'
      setTimeout(() => {
        pageLoading.value = false
      }, 320)
    } else {
      pageLoading.value = false
    }
  } catch (error) {
    console.error('初始化失败:', error)
    pageLoading.value = false
  }
})
</script>

<style scoped>
.fill-height {
  min-height: 100vh;
  min-width: 100vw;
}

.v-card {
  @supports (backdrop-filter: blur(1px)) or (-webkit-backdrop-filter: blur(1px)) {
    background: transparent;
    backdrop-filter: blur(30px);
  }
}
</style>
