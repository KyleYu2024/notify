<template>
  <!-- 页面加载状态 -->
  <LoadingView :loading="pageLoading" text="正在加载模板数据..." v-if="pageLoading" />

  <!-- 页面内容 -->
  <div v-else>
    <!-- 页面标题 -->
    <v-row>
      <v-col>
        <h1 class="text-h6">
          <v-icon icon="mdi-file-document-multiple" class="mr-2"></v-icon>
          模板管理
        </h1>
        <p class="text-body-2 text-medium-emphasis mt-1">
          管理消息模板，支持模板复用和实时预览
        </p>
      </v-col>
      <v-col cols="auto">
        <v-btn @click="showImportDialog = true" variant="outlined" size="small" class="me-2">
          <v-icon icon="mdi-import" class="mr-1"></v-icon>
          导入
        </v-btn>
        <v-btn color="primary" @click="showCreateDialog = true" size="small">
          <v-icon icon="mdi-plus" class="mr-2"></v-icon>
          创建模板
        </v-btn>
      </v-col>
    </v-row>

    <!-- CSS Column 瀑布流模板列表 -->
    <div v-if="Object.keys(templatesStore.templates).length > 0" class="masonry-container">
      <div v-for="(template, templateId) in templatesStore.templates" :key="templateId" class="masonry-item">
        <v-card class="template-card" elevation="2">
          <v-card-title class="d-flex align-center">
            <v-icon icon="mdi-file-document" color="primary" class="mr-3"></v-icon>
            <div class="flex-grow-1">
              <div class="text-h6">{{ template.name }}</div>
              <div class="text-caption text-medium-emphasis">ID: {{ template.id }}</div>
            </div>
          </v-card-title>

          <v-card-text class="template-item">
            <div class="text-body-2 mb-2">标题:</div>
            <v-card variant="outlined" class="template-preview pa-2" style="max-height: 150px; overflow-y: auto;">
              <pre class="text-caption">{{ template.title || '无内容' }}</pre>
            </v-card>
          </v-card-text>
          <v-card-text class="template-item">
            <div class="text-body-2 mb-2">内容:</div>
            <v-card variant="outlined" class="template-preview pa-2" style="max-height: 150px; overflow-y: auto;">
              <pre class="text-caption">{{ template.content || '无内容' }}</pre>
            </v-card>
          </v-card-text>
          <v-card-text class="template-item">
            <div class="text-body-2 mb-2">链接:</div>
            <v-card variant="outlined" class="template-preview pa-2" style="max-height: 150px; overflow-y: auto;">
              <pre class="text-caption">{{ template.url || '无内容' }}</pre>
            </v-card>
          </v-card-text>
          <v-card-text>
            <div class="text-body-2 mb-2">图片:</div>
            <v-card variant="outlined" class="template-preview pa-2" style="max-height: 150px; overflow-y: auto;">
              <pre class="text-caption">{{ template.image || '无内容' }}</pre>
            </v-card>
          </v-card-text>
          <v-card-text>
            <div class="text-body-2 mb-2">目标:</div>
            <v-card variant="outlined" class="template-preview pa-2" style="max-height: 150px; overflow-y: auto;">
              <pre class="text-caption">{{ template.targets || '无内容' }}</pre>
            </v-card>
          </v-card-text>
          <v-card-actions>
            <v-btn variant="flat" size="small" @click="editTemplate(template)" :loading="templatesStore.loading">
              <v-icon icon="mdi-pencil" class="mr-1"></v-icon>
              编辑
            </v-btn>
            <v-btn variant="text" size="small" @click="exportTemplate(template)" :loading="templatesStore.loading">
              <v-icon icon="mdi-share-variant" class="mr-1"></v-icon>
              分享
            </v-btn>
            <v-spacer></v-spacer>
            <v-btn variant="text" size="small" color="error" @click="deleteTemplate(template.id)"
              :loading="templatesStore.loading">
              <v-icon icon="mdi-delete"></v-icon>
            </v-btn>
          </v-card-actions>
        </v-card>
      </div>
    </div>

    <!-- 无数据状态 -->
    <v-row v-else-if="!templatesStore.loading">
      <v-col cols="12" class="text-center py-8">
        <v-icon icon="mdi-file-document-multiple-outline" size="80" class="text-medium-emphasis mb-4"></v-icon>
        <h3 class="text-h6 text-medium-emphasis mb-2">暂无模板</h3>
        <p class="text-body-2 text-medium-emphasis mb-4">点击上方的"创建模板"按钮来创建您的第一个消息模板</p>
        <v-btn color="primary" variant="outlined" @click="showCreateDialog = true">
          <v-icon icon="mdi-plus" class="mr-2"></v-icon>
          创建模板
        </v-btn>
      </v-col>
    </v-row>

    <!-- 加载状态 -->
    <v-row v-else>
      <v-col cols="12" class="text-center py-8">
        <v-progress-circular indeterminate color="primary" size="40"></v-progress-circular>
        <p class="text-body-2 text-medium-emphasis mt-4">加载模板数据中...</p>
      </v-col>
    </v-row>

    <!-- 创建/编辑模板对话框 -->
    <TemplateEditDialog v-model="showCreateDialog" :editing-template="editingTemplate" :loading="templatesStore.loading"
      @save="handleSaveTemplate" @cancel="handleCancelEdit" />

    <!-- 导入模板对话框 -->
    <TemplateImportDialog v-model="showImportDialog" :existing-templates="templatesStore.templates"
      :loading="templatesStore.loading" @import="handleImportTemplates" @cancel="handleCancelImport" />

  </div>
</template>

<script setup lang="ts">
import LoadingView from '@/components/LoadingView.vue'
import TemplateEditDialog from '@/components/dialog/TemplateEditDialog.vue'
import TemplateImportDialog from '@/components/dialog/TemplateImportDialog.vue'
import { useTemplatesStore, type IMessageTemplate } from '@/store/templates'
import { onMounted, ref } from 'vue'
import { useConfirm } from 'vuetify-use-dialog'

const templatesStore = useTemplatesStore()
const createConfirm = useConfirm()
// 页面加载状态
const pageLoading = ref(true)

// 对话框状态
const showCreateDialog = ref(false)
const showImportDialog = ref(false)

// 编辑状态
const editingTemplate = ref<IMessageTemplate | null>(null)

// 加载模板
const loadTemplates = async () => {
  await templatesStore.fetchTemplates()
}

// 编辑模板
const editTemplate = (template: IMessageTemplate) => {
  editingTemplate.value = template
  showCreateDialog.value = true
}


// 删除模板
const deleteTemplate = async (templateId: string) => {
  const template = templatesStore.templates[templateId]

  const isConfirmed = await createConfirm({
    title: '确认删除',
    content: `您确定要删除模板 "${template?.name}" 吗？此操作不可撤销。`,
    dialogProps: {
      width: '500px',
    },
    confirmationText: '删除',
    confirmationButtonProps: {
      color: 'error',
      variant: 'outlined'
    },
    cancellationText: '取消'
  })

  if (isConfirmed) {
    try {
      await templatesStore.deleteTemplate(templateId)
    } catch (error) {
      console.error('删除模板失败:', error)
    }
  }
}

// 处理保存模板
const handleSaveTemplate = async (formData: any) => {
  try {
    if (editingTemplate.value) {
      await templatesStore.updateTemplate(editingTemplate.value.id, formData)
    } else {
      await templatesStore.createTemplate(formData)
    }
    handleCancelEdit()
  } catch (error) {
    console.error('保存模板失败:', error)
  }
}

// 处理取消编辑
const handleCancelEdit = () => {
  showCreateDialog.value = false
  editingTemplate.value = null
}

// 分享模板到剪贴板
const exportTemplate = (template: IMessageTemplate) => {
  templatesStore.shareTemplateToClipboard(template)
}

// 处理导入模板
const handleImportTemplates = async (data: { templates: IMessageTemplate[], conflictStrategy: 'skip' | 'overwrite' | 'generate' }) => {
  try {
    const { templates, conflictStrategy } = data
    await templatesStore.importTemplates(
      templates,
      conflictStrategy === 'overwrite',
      conflictStrategy === 'generate'
    )
    showImportDialog.value = false
  } catch (error) {
    console.error('导入模板失败:', error)
  }
}

// 处理取消导入
const handleCancelImport = () => {
  showImportDialog.value = false
}

// 页面加载时获取数据
onMounted(async () => {
  try {
    await loadTemplates()
  } catch (error) {
    console.error('页面初始化失败:', error)
  } finally {
    pageLoading.value = false
  }
})
</script>

<style scoped>
pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
}

.template-item {
  padding-block: 8px;
}

/* CSS Column 瀑布流容器 */
.masonry-container {
  /* 基础多列布局 */
  column-count: 3;
  column-gap: 16px;
  column-fill: balance;
  /* 平衡各列高度 */
  margin-top: 16px;

  /* 避免列间断页 */
  orphans: 1;
  widows: 1;
}

/* 瀑布流项目 */
.masonry-item {
  /* 防止元素被分割到两列之间 */
  break-inside: avoid;
  page-break-inside: avoid;
  /* 兼容旧版浏览器 */

  /* 项目间距 */
  margin-bottom: 16px;

  /* 确保元素是块级的 */
  display: block;
  width: 100%;
}

/* 模板卡片样式 */
.template-card {
  width: 100%;
  transition: all 0.3s ease;
  border-radius: 8px;
  overflow: hidden;

  /* 防止卡片被分割 */
  break-inside: avoid;
  page-break-inside: avoid;
}

/* 卡片悬停效果 */
.template-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15) !important;
}

/* 模板预览区域 */
.template-preview {
  background-color: rgba(0, 0, 0, 0.02);
  border-radius: 4px;
}

/* 响应式断点 */
@media (max-width: 1200px) {
  .masonry-container {
    column-count: 2;
  }
}

@media (max-width: 768px) {
  .masonry-container {
    column-count: 1;
    column-gap: 0;
    margin-top: 12px;
  }

  .masonry-item {
    margin-bottom: 12px;
  }
}

/* 超小屏幕优化 */
@media (max-width: 480px) {
  .template-card {
    margin: 0;
  }

  .template-card .v-card-text {
    padding: 8px 12px;
  }

  .template-card .v-card-actions {
    padding: 6px 12px 12px;
    flex-wrap: wrap;
    gap: 4px;
  }
}

/* 打印样式优化 */
@media print {
  .masonry-container {
    column-count: 2;
    column-gap: 20px;
  }

  .template-card {
    break-inside: avoid;
    page-break-inside: avoid;
  }
}

/* 加载状态 */
.masonry-container.loading {
  opacity: 0.6;
  pointer-events: none;
}

/* 动画优化 */
.masonry-item {
  animation: fadeInUp 0.4s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 确保在列布局中的良好展示 */
.template-card .v-card-title,
.template-card .v-card-text,
.template-card .v-card-actions {
  break-inside: avoid;
}

/* 高对比度模式支持 */
@media (prefers-contrast: high) {
  .template-card {
    border: 2px solid #000;
  }

  .template-preview {
    background-color: rgba(0, 0, 0, 0.1);
    border: 1px solid #666;
  }
}

/* 减少动画模式支持 */
@media (prefers-reduced-motion: reduce) {

  .template-card,
  .masonry-item {
    transition: none;
    animation: none;
  }

  .template-card:hover {
    transform: none;
  }
}
</style>
