<script setup lang="ts">
import { useAuthStore } from '@/store/auth'
import { onMounted, ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useConfirm } from 'vuetify-use-dialog'
import { useTheme } from 'vuetify'
import { useDisplay } from 'vuetify'
import { updateThemeColor } from '@/common/utils'

const router = useRouter()
const authStore = useAuthStore()
const createConfirm = useConfirm()
const theme = useTheme()
const display = useDisplay()

// 导航抽屉状态
const isDesktop = computed(() => display.mdAndUp.value)
const drawer = ref(isDesktop.value)

// 系统健康状态
const systemHealth = ref(true)
const healthChecking = ref(false)

// 主题相关
const isDark = computed(() => theme.global.current.value.dark)

// 切换主题 (通用方法)
const changeTheme = (themeName: string) => {
  theme.global.name.value = themeName
  localStorage.setItem('theme', themeName)
  updateThemeColor(themeName)
}

// 快捷切换 (Day <-> Night)
const toggleTheme = () => {
  const targetTheme = isDark.value ? 'day' : 'night'
  changeTheme(targetTheme)
}

// 菜单项配置
const menuItems = [
  {
    title: '仪表板',
    subtitle: '系统概览和快速操作',
    icon: 'mdi-view-dashboard',
    to: '/dashboard'
  },
  {
    title: '通知应用',
    subtitle: '管理通知应用配置',
    icon: 'mdi-application',
    to: '/apps'
  },
  {
    title: '通知服务配置',
    subtitle: '管理通知服务设置',
    icon: 'mdi-bell-cog',
    to: '/notifiers'
  },
  {
    title: '模板管理',
    subtitle: '管理消息模板',
    icon: 'mdi-file-document-multiple',
    to: '/templates'
  },
  {
    title: '插件管理',
    subtitle: '管理后端插件与设置',
    icon: 'mdi-puzzle',
    to: '/plugins'
  },
  {
    title: '日志',
    subtitle: '实时查看系统日志',
    icon: 'mdi-math-log',
    to: '/logs'
  }
]

// 检查系统健康状态
const checkSystemHealth = async () => {
  healthChecking.value = true
  try {
    await authStore.checkSystemStatus()
    systemHealth.value = true
  } catch (error) {
    systemHealth.value = false
  } finally {
    healthChecking.value = false
  }
}

// 打开GitHub页面
const openGithub = () => {
  window.open('https://github.com/jianxcao/notify', '_blank')
}

// 处理退出登录
const handleLogout = async () => {
  const isConfirmed = await createConfirm({
    title: '确认退出',
    content: '您确定要退出登录吗？',
    dialogProps: {
      width: '300px',
    },
    confirmationText: '退出',
    cancellationText: '取消'
  })

  if (isConfirmed) {
    authStore.logout()
    router.push('/login')
  }
}

// 监听屏幕大小变化，自动调整drawer状态
watch(isDesktop, (newValue) => {
  drawer.value = newValue
})

// 页面加载时检查系统健康状态
onMounted(() => {
  checkSystemHealth()
  setInterval(checkSystemHealth, 30000)
  updateThemeColor(theme.global.name.value)
})
</script>

<template>
  <v-app-bar density="comfortable" class="header" :elevation="0" border="b">
    <template v-slot:prepend>
      <v-app-bar-nav-icon v-if="!isDesktop || !drawer" @click="drawer = !drawer">
      </v-app-bar-nav-icon>
    </template>
    <v-toolbar-title class="text-body-1 font-weight-medium">
      通知管理系统
    </v-toolbar-title>
    <v-spacer></v-spacer>

    <v-btn icon variant="text" @click="toggleTheme" class="mr-2">
      <v-icon>
        {{ isDark ? 'mdi-weather-night' : 'mdi-white-balance-sunny' }}
      </v-icon>
      <v-tooltip activator="parent" location="bottom">
        {{ isDark ? '切换到白天模式' : '切换到夜晚模式' }}
      </v-tooltip>
    </v-btn>

    <v-menu offset-y>
      <template v-slot:activator="{ props }">
        <v-btn v-bind="props" icon>
          <v-icon icon="mdi-account-circle"></v-icon>
        </v-btn>
      </template>

      <v-card class="user-menu">
        <v-card-item>
          <v-card-title class="d-flex align-center">
            <v-icon icon="mdi-account" class="mr-2"></v-icon>
            {{ authStore.username || 'Admin' }}
          </v-card-title>
        </v-card-item>

        <v-divider></v-divider>

        <v-list class="user-menu-list">
          <v-list-item @click="checkSystemHealth" :disabled="healthChecking">
            <template v-slot:prepend>
              <v-icon :icon="systemHealth ? 'mdi-heart' : 'mdi-heart-broken'"
                :color="systemHealth ? 'success' : 'error'"></v-icon>
            </template>
            <v-list-item-title>系统状态</v-list-item-title>
            <template v-slot:append>
              <v-progress-circular v-if="healthChecking" indeterminate size="20"></v-progress-circular>
            </template>
          </v-list-item>

          <v-divider></v-divider>

          <v-list-item @click="openGithub">
            <template v-slot:prepend>
              <v-icon icon="mdi-github" color="grey-darken-1"></v-icon>
            </template>
            <v-list-item-title>GitHub</v-list-item-title>
          </v-list-item>

          <v-divider v-if="authStore.isAuthRequired"></v-divider>

          <v-list-item @click="handleLogout" color="error" v-if="authStore.isAuthRequired">
            <template v-slot:prepend>
              <v-icon icon="mdi-logout" color="error"></v-icon>
            </template>
            <v-list-item-title>退出登录</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-card>
    </v-menu>
  </v-app-bar>

  <v-navigation-drawer v-model="drawer" :temporary="!isDesktop" :permanent="isDesktop" class="navigation-drawer"
    :elevation="0" border="e">
    <div class="d-flex align-center justify-center pa-4 mb-2">
       <v-img src="/logo.png" max-width="40" class="mr-3"></v-img>
       <span class="text-h6 font-weight-bold" style="letter-spacing: 0.5px;">Notify</span>
    </div>
    
    <v-divider></v-divider>

    <v-list nav density="comfortable" class="mt-2">
      <v-list-item v-for="item in menuItems" :key="item.title" :to="item.to" :prepend-icon="item.icon"
        :title="item.title" :subtitle="item.subtitle" color="primary" rounded="xl" class="mb-1"></v-list-item>
    </v-list>
  </v-navigation-drawer>

  <v-main class="bg-background">
    <v-container fluid class="pa-6">
      <router-view></router-view>
    </v-container>
  </v-main>

</template>

<style scoped lang="less">
.header {
  padding-top: env(safe-area-inset-top);
  background-color: rgb(var(--v-theme-surface));
}

.navigation-drawer {
  padding-top: env(safe-area-inset-top);
  padding-bottom: env(safe-area-inset-bottom);
}

/* 激活菜单项样式增强 */
:deep(.v-list-item--active) {
  background-color: rgb(var(--v-theme-primary), 0.1);
  color: rgb(var(--v-theme-primary));
}

/* 磨砂玻璃效果 */
@supports (backdrop-filter: blur(30px)) {
  .header {
    background-color: rgba(var(--v-theme-background), 0.85);
    backdrop-filter: blur(10px);
  }

  .navigation-drawer {
    background-color: rgba(var(--v-theme-background), 0.85);
    backdrop-filter: blur(30px);
  }

  .user-menu {
    background-color: rgba(var(--v-theme-surface), 0.85) !important;
    backdrop-filter: blur(30px);
  }

  /* 针对 Night 模式的特殊处理 */
  :deep(.v-theme--night) .v-app-bar {
    background-color: rgba(15, 15, 15, 0.75) !important;
    backdrop-filter: blur(20px);
    -webkit-backdrop-filter: blur(20px);
    border-bottom: 0.5px solid rgba(255, 255, 255, 0.1);
  }
}
</style>