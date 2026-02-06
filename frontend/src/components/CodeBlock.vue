<template>
  <div class="code-block-container">
    <div class="code-block" :class="`language-${language}`">
      <div class="code-header">
        <span class="language-label" v-if="language">{{ languageLabel }}</span>
        <v-btn icon="mdi-content-copy" variant="text" size="small" @click="copyCode" class="copy-btn"
          :loading="copying">
          <v-icon size="16"></v-icon>
          <v-tooltip activator="parent" location="bottom">复制代码</v-tooltip>
        </v-btn>
      </div>
      <pre class="code-content"><code>{{ code }}</code></pre>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { copyToClipboard } from '@/common/utils'
import { useToast } from 'vue-toast-notification'

interface Props {
  code: string
  language?: string
}

const props = withDefaults(defineProps<Props>(), {
  language: 'text'
})

const toast = useToast()
const copying = ref(false)

// 语言标签映射
const languageLabels: Record<string, string> = {
  bash: 'Bash',
  javascript: 'JavaScript',
  python: 'Python',
  json: 'JSON',
  text: 'Text'
}

const languageLabel = computed(() => {
  return languageLabels[props.language] || props.language.toUpperCase()
})

// 复制代码
const copyCode = async () => {
  if (copying.value) return

  copying.value = true
  try {
    const success = await copyToClipboard(props.code)
    if (success) {
      toast.success('代码已复制到剪贴板')
    } else {
      toast.error('复制失败，请手动复制')
    }
  } finally {
    copying.value = false
  }
}
</script>

<style scoped lang="less">
@import '@/styles/mix.less';

.code-block-container {
  margin: 8px 0;
}

.code-block {
  background-color: rgba(var(--v-theme-code), 1);
  border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  border-radius: 8px;
  overflow: hidden;
  position: relative;

  .code-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 12px;
    background-color: rgba(var(--v-theme-kbd), 1);
    border-bottom: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
    min-height: 40px;

    .language-label {
      font-size: 0.75rem;
      font-weight: 500;
      color: rgba(var(--v-theme-on-surface), 0.6);
      text-transform: uppercase;
      letter-spacing: 0.5px;
    }

    .copy-btn {
      opacity: 0.7;
      transition: opacity 0.2s;

      &:hover {
        opacity: 1;
      }
    }
  }

  .code-content {
    margin: 0;
    padding: 16px;
    font-family: 'Courier New', Consolas, 'Liberation Mono', Menlo, Courier, monospace;
    font-size: 0.875rem;
    line-height: 1.5;
    white-space: pre;
    overflow-x: auto;
    .scrollbar();
    color: rgba(var(--v-theme-on-code), 1);

    code {
      background: none;
      padding: 0;
      border-radius: 0;
      font-family: inherit;
      font-size: inherit;
      color: inherit;
    }
  }


}
</style>
