<template>
  <div class="page-container">
    <div class="header">
      日志
    </div>
    <v-virtual-scroll :items="reversedLogs" :height="scrollHeight" class="log-virtual-scroll" item-resizable
      :item-size="estimatedItemSize">
      <template v-slot:default="{ item: log }">
        <div class="log-line">
          <div class="log-line-left">
            <v-chip class="ma-1 level-chip" size="x-small" :color="levelColor(log.level)" label>{{ log.level }}</v-chip>
            <span class="timestamp">{{ log.time }}</span>
          </div>
          <span class="message">
            <span>{{ log.msg }}</span>
            <span v-if="hasExtraFields(log)">{{ _.omit(log, ['level', 'time', 'msg', '_k']) }}</span>
          </span>
        </div>
      </template>
    </v-virtual-scroll>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import * as _ from 'lodash-es'

interface LogEntry {
  level: string
  time: string
  msg: string
  [key: string]: any
}
// 估算的每个日志项高度（像素）
const estimatedItemSize = 85

const logs = ref<LogEntry[]>([])
let eventSource: EventSource | null = null
let nextId = 1

// 计算滚动容器高度 (视口高度 - 头部高度)
const scrollHeight = computed(() => {
  return 'calc(100vh - 88px)'
})

// 反转日志列表，让最新的日志显示在顶部
const reversedLogs = computed(() => {
  return [...logs.value].reverse()
})

function levelColor(level: string) {
  switch (level) {
    case 'DEBUG':
      return 'indigo'
    case 'INFO':
      return 'green'
    case 'WARN':
      return 'orange'
    case 'ERROR':
      return 'red'
    default:
      return 'grey'
  }
}

// 检查是否有额外字段需要显示
function hasExtraFields(log: LogEntry) {
  const extraFields = _.omit(log, ['level', 'time', 'msg', '_k'])
  return Object.keys(extraFields).length > 0
}

onMounted(() => {
  eventSource = new EventSource('/api/v1/logs/stream')
  eventSource.onmessage = (e) => {
    const logArray = JSON.parse(e.data)

    // 处理数组中的每条日志
    logArray.forEach((logEntry: any) => {
      logEntry.time = new Date(logEntry.time).toLocaleString()
      logEntry._k = nextId++
    })
    // 直接添加所有日志
    logs.value.push(...logArray)
    // 控制日志数量上限
    if (logs.value.length > 500) {
      const removeCount = logs.value.length - 500
      logs.value.splice(0, removeCount)
    }
  }
})

onBeforeUnmount(() => {
  if (eventSource) {
    eventSource.close()
  }
})
</script>

<style scoped lang="less">
@import '@/styles/mix.less';

.page-container {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.header {
  font-size: 22px;
  font-weight: 500;
  flex-shrink: 0;
  height: 32px;
  line-height: 32px;
}

.log-virtual-scroll {
  flex: 1;
  overflow-x: hidden;

  :deep(.v-virtual-scroll__container) {
    .scrollbar();
  }
}

.log-line {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
  font-size: 0.9rem;
  line-height: 1.4;
  padding-block: 8px;
  border-bottom: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));

  .log-line-left {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .level-chip {
    width: 60px;
    justify-content: center;
    padding: 0 2px;
    width: 70px;
    flex: 0 0 70px;
  }

  .timestamp {
    color: gray;
    margin: 0 4px;
    flex: 0 0 150px;
  }

  .message {
    word-break: break-all;
    display: flex;
    flex-direction: column;
    gap: 4px;
    flex: 1;
    padding-left: 8px;
  }
}
</style>
