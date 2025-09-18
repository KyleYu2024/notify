<template>
  <div>
    <LoadingView :loading="pageLoading" text="正在加载插件数据..." />
    <div v-if="!pageLoading">
      <div class="flex items-center justify-between mb-4">
        <div class="text-lg font-semibold">插件管理</div>
        <v-btn color="primary" variant="flat" @click="refresh" :loading="store.loading">
          刷新
        </v-btn>
      </div>
      <v-card class="mt-4 bg-transparent">
        <v-card-text>
          <v-row v-if="store.list.length > 0">
            <v-col v-for="p in store.list" :key="p.id" cols="12" md="6" lg="4" class="p-1!">
              <v-card class="plugin-card" :class="{ 'plugin-disabled': !p.enabled }" elevation="4">
                <v-card-item>
                  <div>
                    <div class="text-base font-medium">{{ p.name }} ({{ p.id }})</div>
                    <div class="text-sm text-gray-500">v{{ p.version }} · {{ p.description }}</div>
                  </div>
                </v-card-item>
                <v-card-actions class="mt-auto justify-between gap-3">
                  <v-switch color="primary" :model-value="p.enabled"
                    @update:model-value="(v) => toggleEnabled(p.id, !!v)" />
                  <v-btn variant="text" icon="mdi-cog" aria-label="设置" @click="openSettings(p)" v-if="p.ui" />
                </v-card-actions>
              </v-card>
            </v-col>
          </v-row>
          <v-row v-else-if="!store.loading">
            <v-col cols="12" class="text-center py-8">
              <v-icon icon="mdi-puzzle-outline" size="80" class="text-medium-emphasis mb-4"></v-icon>
              <h3 class="text-h6 text-medium-emphasis mb-2">暂无插件</h3>
              <p class="text-body-2 text-medium-emphasis mb-4">点击下方“刷新”，或在服务器部署插件后再试</p>
              <v-btn color="primary" variant="outlined" @click="refresh" :loading="store.loading">
                <v-icon icon="mdi-refresh" class="mr-2"></v-icon>
                刷新
              </v-btn>
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
      <PluginSettingsDialog v-model="show" :plugin="current" :settings="form" :saving="saving" @save="onSaveSettings" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { usePluginsStore } from '@/store/plugins'
import type { PluginInfo } from '@/store/plugins'
import PluginSettingsDialog from '@/components/dialog/PluginSettingsDialog.vue'

const store = usePluginsStore()

const show = ref(false)
const current = ref<PluginInfo | null>(null)
const form = ref<Record<string, any>>({})
const saving = ref(false)
const pageLoading = ref(true)

async function refresh() {
  await store.fetchList()
}

function openSettings(p: PluginInfo) {
  console.debug('openSettings', p)
  current.value = p
  form.value = { ...(p.settings || {}) }
  show.value = true
}

async function onSaveSettings(payload: Record<string, any>) {
  if (!current.value) return
  saving.value = true
  try {
    const res = await store.updateSettings(current.value.id, payload)
    if (res.code === 0) {
      show.value = false
      form.value = { ...payload }
    }
  } finally {
    saving.value = false
  }
}

async function toggleEnabled(id: string, v: boolean) {
  await store.setEnabled(id, v)
}


onMounted(async () => {
  try {
    await refresh()
  } finally {
    pageLoading.value = false
  }
})
</script>

<style scoped lang="less">
.plugin-card {
  height: 150px;
  display: flex;
  flex-direction: column;
}
</style>
