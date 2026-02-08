<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from '@/store/auth'
import { useTheme } from 'vuetify'
import { updateThemeColor } from '@/common/utils'

const authStore = useAuthStore()
const theme = useTheme()

// 应用启动时恢复认证状态
onMounted(() => {
  authStore.restoreAuth()
  updateThemeColor(theme.global.name.value)
})
</script>

<template>
  <v-app class="app">
    <router-view />
  </v-app>
</template>

<style lang="less" scoped>
.app {
  /* 移除全局 Safe Area Padding，由具体组件（如 v-app-bar）自行处理 */
}
</style>

<style>
html {
  scrollbar-width: thin !important;
  scrollbar-color: rgba(var(--v-theme-perfect-scrollbar-thumb), 1) transparent;
}
</style>
