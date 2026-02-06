<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from '@/store/auth'

const authStore = useAuthStore()

// 应用启动时恢复认证状态
onMounted(() => {
  authStore.restoreAuth()
})
</script>

<template>
  <v-app class="app">
    <router-view />
  </v-app>
</template>

<style lang="less" scoped>
.app {
  padding-top: env(safe-area-inset-top);
  padding-bottom: env(safe-area-inset-bottom);
}
</style>

<style>
html {
  scrollbar-width: thin !important;
  scrollbar-color: rgba(var(--v-theme-perfect-scrollbar-thumb), 1) transparent;
}
</style>
