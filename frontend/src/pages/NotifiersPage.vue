<template>
  <!-- 页面加载状态 -->
  <LoadingView :loading="pageLoading" text="正在加载通知服务数据..." />

  <!-- 页面内容 -->
  <div v-if="!notifiersStore.loading">
    <!-- 页面标题 -->
    <v-row>
      <v-col>
        <h1 class="text-h6">
          <v-icon icon=" mdi-bell-cog" class="mr-2"></v-icon>
          通知服务
        </h1>
        <p class="text-body-2 text-medium-emphasis mt-1">
          管理通知服务实例，包括企业微信、Telegram、钉钉等通知渠道
        </p>
      </v-col>
      <v-col cols="auto">
        <v-btn color="primary" size="small" @click="addNotifier" :loading="notifiersStore.loading" variant="flat">
          <v-icon icon="mdi-plus" class="mr-1"></v-icon>
          新增通知服务
        </v-btn>
      </v-col>

    </v-row>

    <!-- 通知服务列表 -->
    <v-card class="mt-4 bg-transparent">
      <v-card-text>
        <v-row v-if="Object.keys(notifiersStore.notifiers).length > 0">
          <v-col v-for="(notifier, key) in notifiersStore.notifiers" :key="key" cols="12" md="6" lg="4">
            <v-card class="notifier-card" :class="{ 'notifier-disabled': !notifier.enabled }" elevation="4">
              <v-card-title class="d-flex align-center">
                <v-icon :icon="getNotifierIcon(notifier.type)" :color="notifier.enabled ? 'primary' : 'grey'"
                  class="mr-3" size="32"></v-icon>
                <div class="flex-grow-1">
                  <div class="text-h6">{{ key }}</div>
                  <div class="text-caption text-medium-emphasis">
                    {{ getNotifierTypeName(notifier.type) }}
                  </div>
                </div>
                <v-chip :color="notifier.enabled ? 'success' : 'error'" :text="notifier.enabled ? '启用' : '禁用'"
                  size="small" variant="flat"></v-chip>
              </v-card-title>
              <v-card-actions>
                <v-btn variant="flat" size="small" @click="editNotifier(key)" :loading="notifiersStore.loading">
                  <v-icon icon="mdi-pencil" class="mr-1"></v-icon>
                  编辑
                </v-btn>
                <v-spacer></v-spacer>
                <v-btn variant="text" size="small" color="error" @click="deleteNotifier(key)"
                  :loading="notifiersStore.loading">
                  <v-icon icon="mdi-delete"></v-icon>
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-col>
        </v-row>

        <!-- 无数据状态 -->
        <v-row v-else-if="!notifiersStore.loading">
          <v-col cols="12" class="text-center py-8">
            <v-icon icon="mdi-bell-cog-outline" size="80" class="text-medium-emphasis mb-4"></v-icon>
            <h3 class="text-h6 text-medium-emphasis mb-2">暂无通知服务</h3>
            <p class="text-body-2 text-medium-emphasis mb-4">点击上方的"新增通知服务"按钮来添加您的第一个通知服务</p>
            <v-btn color="primary" variant="outlined" @click="addNotifier">
              <v-icon icon="mdi-plus" class="mr-2"></v-icon>
              新增通知服务
            </v-btn>
          </v-col>
        </v-row>

        <!-- 加载状态 -->
        <v-row v-else>
          <v-col cols="12" class="text-center py-8">
            <v-progress-circular indeterminate color="primary" size="40"></v-progress-circular>
            <p class="text-body-2 text-medium-emphasis mt-4">加载通知服务中...</p>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- 编辑通知服务对话框 -->
    <NotifierEditDialog v-model="showEditDialog" :notifier-key="editNotifierKey" :loading="notifiersStore.loading"
      @save="handleSaveNotifier" @cancel="handleCancelEdit" />

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useNotifiersStore } from '@/store/notifiers'
import { useConfirm } from 'vuetify-use-dialog'
import LoadingView from '@/components/LoadingView.vue'
import NotifierEditDialog from '@/components/dialog/NotifierEditDialog.vue'
import { getNotifierTypeName, getNotifierIcon } from '@/common/utils'

const notifiersStore = useNotifiersStore()
const createConfirm = useConfirm()

// 页面加载状态
const pageLoading = ref(true)

// 对话框状态
const showEditDialog = ref(false)

// 编辑状态
const editNotifierKey = ref('')



// 加载通知服务
const loadNotifiers = async () => {
  await notifiersStore.fetchNotifiers()
}

// 新增通知服务
const addNotifier = () => {
  editNotifierKey.value = '' // 新增时，key为空
  showEditDialog.value = true
}

// 编辑通知服务
const editNotifier = async (key: string) => {
  editNotifierKey.value = key
  showEditDialog.value = true
}

// 删除通知服务
const deleteNotifier = async (key: string) => {
  const isConfirmed = await createConfirm({
    title: '确认删除',
    content: `您确定要删除通知服务 "${key}" 吗？此操作不可撤销。`,
    confirmationText: '删除',
    dialogProps: {
      width: '500px',
    },
    confirmationButtonProps: {
      color: 'error',
      variant: 'flat'
    },
    cancellationText: '取消'
  })

  if (isConfirmed) {
    try {
      await notifiersStore.deleteNotifier(key)
    } catch (error) {
      console.error('删除通知服务失败:', error)
    }
  }
}

// 处理保存通知服务
const handleSaveNotifier = async (key: string, config: any) => {
  try {
    await notifiersStore.updateNotifier(key, config)
    handleCancelEdit()
  } catch (error) {
    console.error('保存通知服务失败:', error)
  }
}

// 处理取消编辑
const handleCancelEdit = () => {
  showEditDialog.value = false
  editNotifierKey.value = ''
}

// 页面加载时获取数据
onMounted(async () => {
  try {
    await loadNotifiers()
  } catch (error) {
    console.error('页面初始化失败:', error)
  } finally {
    pageLoading.value = false
  }
})
</script>

<style scoped>
.notifier-card {
  transition: all 0.3s ease;
}

.notifier-card:hover {
  transform: translateY(-2px);
}

.notifier-disabled {
  opacity: 0.6;
}
</style>
