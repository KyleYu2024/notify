// 通用工具函数

import { NotifierTypeMap } from './types'

/**
 * 格式化时间戳
 */
export function formatTimestamp(timestamp: string | number | Date): string {
  const date = new Date(timestamp)
  if (isNaN(date.getTime())) {
    return '无效时间'
  }
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  })
}

/**
 * 格式化字节大小
 */
export function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B'

  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

/**
 * 深拷贝对象
 */
export function deepClone<T>(obj: T): T {
  if (obj === null || typeof obj !== 'object') {
    return obj
  }

  if (obj instanceof Date) {
    return new Date(obj.getTime()) as unknown as T
  }

  if (obj instanceof Array) {
    return obj.map((item) => deepClone(item)) as unknown as T
  }

  if (typeof obj === 'object') {
    const clonedObj = {} as T
    for (const key in obj) {
      if (obj.hasOwnProperty(key)) {
        clonedObj[key] = deepClone(obj[key])
      }
    }
    return clonedObj
  }

  return obj
}

/**
 * 防抖函数
 */
export function debounce<T extends (...args: any[]) => void>(
  func: T,
  delay: number
): (...args: Parameters<T>) => void {
  let timeoutId: ReturnType<typeof setTimeout>

  return (...args: Parameters<T>) => {
    clearTimeout(timeoutId)
    timeoutId = setTimeout(() => func(...args), delay)
  }
}

/**
 * 节流函数
 */
export function throttle<T extends (...args: any[]) => void>(
  func: T,
  delay: number
): (...args: Parameters<T>) => void {
  let lastCall = 0

  return (...args: Parameters<T>) => {
    const now = Date.now()
    if (now - lastCall >= delay) {
      lastCall = now
      func(...args)
    }
  }
}

/**
 * 生成随机字符串
 */
export function generateRandomString(length: number = 8): string {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  let result = ''
  for (let i = 0; i < length; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  return result
}

/**
 * 验证邮箱格式
 */
export function isValidEmail(email: string): boolean {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

/**
 * 验证URL格式
 */
export function isValidUrl(url: string): boolean {
  try {
    new URL(url)
    return true
  } catch {
    return false
  }
}

/**
 * 获取通知服务类型显示名称
 */
export function getNotifierTypeName(type: string): string {
  const names: Record<string, string> = {
    [NotifierTypeMap.wechatWorkAPPBot]: '企业微信',
    [NotifierTypeMap.telegramAppBot]: 'Telegram',
    [NotifierTypeMap.dingTalkAppBot]: '钉钉',
  }
  return names[type] || type
}

/**
 * 获取通知服务图标
 */
export function getNotifierIcon(type: string): string {
  const icons: Record<string, string> = {
    [NotifierTypeMap.wechatWorkAPPBot]: 'mdi-wechat',
    [NotifierTypeMap.telegramAppBot]: 'mdi-telegram',
    [NotifierTypeMap.dingTalkAppBot]: 'mdi-message-processing',
  }
  return icons[type] || 'mdi-bell'
}

/**
 * 获取通知级别颜色
 */
export function getNotificationLevelColor(level: string): string {
  const colors: Record<string, string> = {
    info: 'primary',
    warning: 'warning',
    error: 'error',
    success: 'success',
  }
  return colors[level] || 'primary'
}

/**
 * 格式化敏感信息
 */
export function maskSensitiveInfo(value: string, visibleChars: number = 4): string {
  if (!value || value.length <= visibleChars) {
    return '••••••••'
  }

  const visible = value.substring(0, visibleChars)
  const masked = '•'.repeat(Math.max(value.length - visibleChars, 4))
  return visible + masked
}

/**
 * 检查对象是否为空
 */
export function isEmpty(obj: any): boolean {
  if (obj == null) return true
  if (Array.isArray(obj) || typeof obj === 'string') return obj.length === 0
  if (typeof obj === 'object') return Object.keys(obj).length === 0
  return false
}

/**
 * 延迟函数
 */
export function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

/**
 * 安全的JSON解析
 */
export function safeJsonParse<T>(json: string, defaultValue: T): T {
  try {
    return JSON.parse(json)
  } catch {
    return defaultValue
  }
}

/**
 * 格式化错误消息
 */
export function formatErrorMessage(error: any): string {
  if (typeof error === 'string') return error
  // 适配新的BaseRes错误格式
  if (error?.response?.data?.msg) return error.response.data.msg
  if (error?.response?.data?.error) return error.response.data.error
  if (error?.response?.data?.message) return error.response.data.message
  if (error?.message) return error.message
  return '发生未知错误'
}

/**
 * 复制文本到剪贴板
 */
export async function copyToClipboard(text: string): Promise<boolean> {
  try {
    if (navigator.clipboard && window.isSecureContext) {
      // 现代浏览器的 Clipboard API
      await navigator.clipboard.writeText(text)
      return true
    } else {
      // 降级方案：使用传统的 execCommand
      const textArea = document.createElement('textarea')
      textArea.value = text
      textArea.style.position = 'fixed'
      textArea.style.left = '-999999px'
      textArea.style.top = '-999999px'
      document.body.appendChild(textArea)
      textArea.focus()
      textArea.select()

      const success = document.execCommand('copy')
      document.body.removeChild(textArea)
      return success
    }
  } catch (error) {
    console.error('复制到剪贴板失败:', error)
    return false
  }
}

/**
 * 获取当前域名和端口
 */
export function getCurrentBaseUrl(): string {
  return window.location.origin
}
