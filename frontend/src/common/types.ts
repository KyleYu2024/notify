// 通用类型定义

// 字段映射接口
export interface FieldMapping {
  title?: string
  content?: string
  image?: string
  url?: string
  level?: string
}

// 应用认证接口
export interface AppAuth {
  enabled: boolean
  token: string
}

// 通知应用接口
export interface NotificationApp {
  appId: string
  name: string
  description?: string
  enabled: boolean
  notifiers: string[]
  templateId: string
  pluginId?: string
  defaultImage?: string
  auth: AppAuth
  fieldMapping?: FieldMapping
}

// 通知服务实例接口
export interface NotifierInstance {
  type: string
  enabled: boolean
  [key: string]: any
}

// 消息模板接口
export interface MessageTemplate {
  id: string
  name: string
  template: string
}

// 通知请求接口
export interface NotificationRequest {
  title: string
  content: string
  level: string
  image?: string
  url?: string
  data?: Record<string, any>
  targets?: string[]
}

// API响应接口 - 适配BaseRes结构
export interface ApiResponse<T = any> {
  code: number // 状态码：0表示成功，其他表示错误
  msg: string // 消息描述
  data?: T // 实际数据
}

// 系统状态接口
export interface SystemStatus {
  auth_enabled: boolean
  admin_endpoints: Record<string, string>
  notification_endpoints: Record<string, string>
  supported_apps: string[]
  version?: string
}

// 通知服务类型
export type NotifierType =
  | 'wechatWorkAPPBot'
  | 'wechatWorkWebhookBot'
  | 'telegramAppBot'
  | 'dingTalkAppBot'
  | 'feishuAppBot'

export const NotifierTypeMap = {
  wechatWorkAPPBot: 'wechatWorkAPPBot',
  wechatWorkWebhookBot: 'wechatWorkWebhookBot',
  telegramAppBot: 'telegramAppBot',
  dingTalkAppBot: 'dingTalkAppBot',
  feishuAppBot: 'feishuAppBot',
} as const

// 通知服务类型选项
export const notifierTypeOptions = [
  { title: '企业微信应用', value: NotifierTypeMap.wechatWorkAPPBot },
  { title: '企业微信群机器人', value: NotifierTypeMap.wechatWorkWebhookBot },
  { title: 'Telegram', value: NotifierTypeMap.telegramAppBot },
  { title: '钉钉', value: NotifierTypeMap.dingTalkAppBot },
  { title: '飞书', value: NotifierTypeMap.feishuAppBot },
]

// 通知级别
export type NotificationLevel = 'info' | 'warning' | 'error' | 'success'

// 企业微信应用配置
export interface WechatWorkConfig {
  enabled: boolean
  corp_id: string
  agent_id: string
  secret: string
  targets?: string
  proxy?: string
}

// 企业微信群机器人配置
export interface WechatWorkWebhookConfig {
  enabled: boolean
  key: string
  proxy?: string
}

// Telegram配置
export interface TelegramConfig {
  enabled: boolean
  bot_token: string
  chat_id: string
  proxy?: string
}

// 钉钉配置
export interface DingTalkConfig {
  enabled: boolean
  access_token: string
  secret?: string
  proxy?: string
  targets?: string
}

// 飞书配置
export interface FeishuConfig {
  enabled: boolean
  app_id: string
  app_secret: string
  targets?: string
  proxy?: string
}

// 通知服务配置联合类型
export type NotifierConfig =
  | WechatWorkConfig
  | WechatWorkWebhookConfig
  | TelegramConfig
  | DingTalkConfig
  | FeishuConfig
