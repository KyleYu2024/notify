<template>
  <v-dialog v-model="isOpen" width="90%" max-width="800px">
    <v-card class="dialog-card">
      <v-card-title class="dialog-header">
        {{ plugin?.name }} 设置
      </v-card-title>
      <v-card-text class="dialog-body">
        <FormRender :model="localForm" :config="(plugin?.ui as UIConfig)" />
      </v-card-text>
      <v-card-actions class="dialog-actions">
        <v-spacer />
        <v-btn variant="text" @click="handleCancel">取消</v-btn>
        <v-btn variant="flat" :loading="saving" @click="handleSave">保存</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import FormRender from '@/components/FormRender.vue'
import type { PluginInfo, UIConfig } from '@/store/plugins'

interface Props {
  modelValue: boolean
  plugin: PluginInfo | null
  settings?: Record<string, any>
  saving?: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'save', value: Record<string, any>): void
}>()

const isOpen = computed({
  get: () => props.modelValue,
  set: (v: boolean) => emit('update:modelValue', v)
})

const localForm = ref<Record<string, any>>({ ...(props.settings || {}) })

watch(
  () => [props.settings, props.plugin] as const,
  () => {
    localForm.value = { ...(props.settings || {}) }
  }
)

function handleCancel() {
  emit('update:modelValue', false)
}

function handleSave() {
  emit('save', localForm.value)
}
</script>

<style scoped lang="less">
.dialog-card {
  display: flex;
  flex-direction: column;
  max-height: 80vh;
  overflow: hidden;
}

.dialog-header {
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  background: rgb(var(--v-theme-surface));
}

.dialog-body {
  flex: 1;
  overflow: auto;
}

.dialog-actions {
  border-top: 1px solid rgba(0, 0, 0, 0.06);
  background: rgb(var(--v-theme-surface));
}
</style>
