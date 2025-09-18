<template>
  <!-- 应用栏 -->
  <v-app-bar density="comfortable" class="header">
    <template v-slot:prepend>
      <v-app-bar-nav-icon v-if="!isDesktop || !drawer" @click="drawer = !drawer">
      </v-app-bar-nav-icon>
    </template>
    <v-toolbar-title>
      通知管理系统
    </v-toolbar-title>
    <v-spacer></v-spacer>
    <!-- 用户菜单 -->
    <v-menu offset-y>
      <template v-slot:activator="{ props }">
        <v-btn v-bind="props" icon>
          <v-icon icon="mdi-account-circle"></v-icon>
        </v-btn>
      </template>

      <v-card class="user-menu">
        <!-- 用户信息 -->
        <v-card-item>
          <v-card-title class="d-flex align-center">
            <v-icon icon="mdi-account" class="mr-2"></v-icon>
            {{ authStore.username }}
          </v-card-title>
        </v-card-item>

        <v-divider></v-divider>

        <!-- 系统状态 -->
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

          <!-- 主题切换 -->
          <v-list-item @click.stop="themeMenu = !themeMenu">
            <template v-slot:prepend>
              <v-icon :icon="getThemeIcon(currentTheme)"></v-icon>
            </template>

            <v-menu offset-x v-model="themeMenu">
              <template v-slot:activator="{ props }">
                <v-list-item-title v-bind="props">主题切换</v-list-item-title>
              </template>
              <v-list density="compact" class="user-menu">
                <v-list-item v-for="theme in availableThemes" :key="theme.key" @click="changeTheme(theme.key)"
                  :active="currentTheme === theme.key">
                  <template v-slot:prepend>
                    <v-icon :icon="theme.icon" :color="theme.color"></v-icon>
                  </template>
                  <v-list-item-title>{{ theme.name }}</v-list-item-title>
                </v-list-item>
              </v-list>
            </v-menu>
          </v-list-item>

          <v-divider></v-divider>

          <!-- GitHub 链接 -->
          <v-list-item @click="openGithub">
            <template v-slot:prepend>
              <v-icon icon="mdi-github" color="grey-darken-1"></v-icon>
            </template>
            <v-list-item-title>GitHub</v-list-item-title>
            <template v-slot:append>
              <v-icon icon="mdi-open-in-new" size="small" color="grey"></v-icon>
            </template>
          </v-list-item>

          <v-divider v-if="authStore.isAuthRequired"></v-divider>

          <!-- 退出登录 -->
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

  <!-- 导航抽屉 -->
  <v-navigation-drawer v-model="drawer" :temporary="!isDesktop" :permanent="isDesktop" class="navigation-drawer">
    <v-list>
      <v-list-item prepend-icon="mdi-account" :title="authStore.username"></v-list-item>
    </v-list>
    <v-divider></v-divider>

    <v-list nav>
      <v-list-item v-for="item in menuItems" :key="item.title" :to="item.to" :prepend-icon="item.icon"
        :title="item.title" :subtitle="item.subtitle" color="primary"></v-list-item>
    </v-list>
    <template v-slot:append>
      <div class="pa-2">
        <!-- GitHub 链接 -->
        <v-btn color="grey-darken-1" variant="outlined" block @click="openGithub" class="mb-2">
          <v-icon icon="mdi-github" class="mr-2"></v-icon>
          GitHub
          <v-icon icon="mdi-open-in-new" size="small" class="ml-1"></v-icon>
        </v-btn>

        <!-- 退出登录 -->
        <v-btn v-if="authStore.isAuthRequired" color="error" variant="outlined" block @click="handleLogout">
          <v-icon icon="mdi-logout" class="mr-2"></v-icon>
          退出登录
        </v-btn>
      </div>
    </template>
  </v-navigation-drawer>

  <!-- 主内容区域 -->
  <v-main>
    <v-container fluid>
      <router-view></router-view>
    </v-container>
  </v-main>

</template>

<script setup lang="ts">
import { useAuthStore } from '@/store/auth'
import { onMounted, ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useConfirm } from 'vuetify-use-dialog'
import { useTheme } from 'vuetify'
import { useDisplay } from 'vuetify'

const router = useRouter()
const authStore = useAuthStore()
const createConfirm = useConfirm()
const theme = useTheme()
const display = useDisplay()
const themeMenu = ref(false)

// 导航抽屉状态
const isDesktop = computed(() => display.mdAndUp.value)
const drawer = ref(isDesktop.value)

// 系统健康状态
const systemHealth = ref(true)
const healthChecking = ref(false)

// 主题相关
const currentTheme = computed(() => theme.global.name.value)

const availableThemes = [
  {
    key: 'light',
    name: '浅色主题',
    icon: 'mdi-white-balance-sunny',
    color: 'orange'
  },
  {
    key: 'dark',
    name: '深色主题',
    icon: 'mdi-moon-waning-crescent',
    color: 'blue'
  },
  {
    key: 'purple',
    name: '紫色主题',
    icon: 'mdi-palette',
    color: 'purple'
  }
]

// 获取主题图标
const getThemeIcon = (themeName: string) => {
  const themeConfig = availableThemes.find(t => t.key === themeName)
  return themeConfig?.icon || 'mdi-palette'
}

// 切换主题
const changeTheme = (themeName: string) => {
  theme.global.name.value = themeName
  localStorage.setItem('theme', themeName)
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

  // 定期检查系统状态
  setInterval(checkSystemHealth, 30000) // 每30秒检查一次
})
</script>

<style scoped lang="less">
.header {
  padding-top: env(safe-area-inset-top);
}

.navigation-drawer {
  padding-top: env(safe-area-inset-top);
  padding-bottom: env(safe-area-inset-bottom);
}

@supports (backdrop-filter: blur(30px)) {
  .header {
    background-color: rgba(var(--v-theme-background), 0.5);
    backdrop-filter: blur(10px);
  }

  .navigation-drawer {
    background-color: rgba(var(--v-theme-background), 0.8);
    backdrop-filter: blur(30px);
  }

  .user-menu {
    background-color: rgba(var(--v-theme-surface), 0.7) !important;
    backdrop-filter: blur(30px);
  }

  .user-menu-list {
    background-color: transparent
  }
}
</style>
