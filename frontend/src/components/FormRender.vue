<script setup lang="ts">
import type { UIConfig } from '@/store/plugins';
import { get, set, isString } from 'lodash-es'

// 定义 props
defineProps<{
  config: UIConfig // JSON 配置
  model: Record<string, any> // 数据模型
}>()

const EXPRESSION_CACHE = new Map<string, (model: any) => any>()
const EVENT_CACHE = new Map<string, (model: any, event: any) => any>()
const hasOwn = (obj: any, key: string) => Object.prototype.hasOwnProperty.call(obj, key)
const getByPath = (obj: any, path: string) => {
  if (!isString(path)) {
    return undefined
  }
  if (hasOwn(obj, path)) {
    return obj[path]
  }
  return get(obj, path)
}
const setByPath = (obj: any, path: string, value: any) => {
  if (!isString(path)) return
  if (hasOwn(obj, path)) {
    obj[path] = value
  } else {
    set(obj, path, value)
  }
}
const isExpression = (value: string) => value.startsWith('{{') && value.endsWith('}}')
const extractExpression = (value: string) => value.slice(2, -2).trim()

const compileExpression = (expression: string) => {
  const cached = EXPRESSION_CACHE.get(expression)
  if (cached) return cached
  const fn = new Function('model', `with(model) { return ${expression} }`) as (m: any) => any
  EXPRESSION_CACHE.set(expression, fn)
  return fn
}

const evalExpressionSafely = (expression: string, model: any, fallback?: any) => {
  try {
    return compileExpression(expression)(model)
  } catch (e) {
    console.error('表达式执行错误:', expression, e)
    return fallback
  }
}

const compileEventHandler = (code: string) => {
  const cached = EVENT_CACHE.get(code)
  if (cached) return cached
  const handler = new Function(
    'model',
    'event',
    `
      try {
        with(model) {
          return (${code})(event);
        }
      } catch(e) {
        console.error('事件处理函数执行错误:', e);
      }
    `,
  ) as (m: any, e: any) => any
  EVENT_CACHE.set(code, handler)
  return handler
}



const resolveComp = (name: string) => {
  try {
    // 先尝试按已注册组件解析
    return resolveComponent(name)
  } catch {
    // 回退为原始标签（如 div/span 等原生标签），避免抛错
    return name
  }
}

/**
 * 解析属性，支持 v-model 和动态绑定
 * @param rawProps 原始属性
 * @param model 数据模型
 * @returns 解析后的属性
 */
const parseProps = (rawProps: Record<string, any> = {}, model: Record<string, any>) => {
  if (!rawProps) {
    rawProps = {}
  }
  const parsedProps: Record<string, any> = {}

  for (const [key, value] of Object.entries(rawProps)) {
    if (key === 'modelvalue' || key === 'modelValue') {
      // 将 modelValue 转换为 value/onUpdate:value 的形式（兼容部分原生/三方组件）
      const current = isString(value) ? getByPath(model, value) : value
      parsedProps['value'] = current
      parsedProps['onUpdate:value'] = (newValue: any) => {
        if (isString(value)) setByPath(model, value, newValue)
      }
    } else if (['model', 'v-model'].includes(key)) {
      // 处理 v-model
      const current = isString(value) ? getByPath(model, value) : value
      parsedProps['modelValue'] = current
      parsedProps['onUpdate:modelValue'] = (newValue: any) => {
        if (isString(value)) setByPath(model, value, newValue)
      }
    } else if (['show', 'v-show'].includes(key)) {
      // 处理 v-show，实现显示隐藏
      let isVisible: any
      if (isString(value) && isExpression(value)) {
        const expression = extractExpression(value)
        isVisible = evalExpressionSafely(expression, model, true)
      } else if (isString(value)) {
        isVisible = getByPath(model, value)
      } else {
        isVisible = !!value
      }
      // 动态设置 style.display
      if (!parsedProps.style) {
        parsedProps.style = {}
      }
      parsedProps.style.display = isVisible ? '' : 'none'
    } else if (key.startsWith('model:') || key.startsWith('v-model:')) {
      // 处理 v-model:<prop>
      const propName = key.split(':')[1]
      const current = isString(value) ? getByPath(model, value) : value
      parsedProps[propName] = current
      parsedProps[`onUpdate:${propName}`] = (newValue: any) => {
        if (isString(value)) setByPath(model, value, newValue)
      }
    } else if (key.startsWith('on')) {
      // 处理事件监听，值是函数的代码 function xxx(e) { ... }
      if (isString(value)) {
        const handler = compileEventHandler(value)
        parsedProps[key] = (...args: any[]) => {
          const [event] = args
          return handler(model, event)
        }
      } else if (typeof value === 'function') {
        parsedProps[key] = value
      }
    } else {
      // 如果是表达式，需要绑定
      if (isString(value) && isExpression(value)) {
        const expression = extractExpression(value)
        parsedProps[key] = evalExpressionSafely(expression, model)
      } else if (isString(value)) {
        // 如果是数据模型路径/属性，直接绑定
        const v = getByPath(model, value)
        parsedProps[key] = v !== undefined ? v : value
      } else {
        // 其他情况直接赋值
        parsedProps[key] = value
      }
    }
  }

  return parsedProps
}

/**
 * 渲染插槽内容
 * @param slotContent 插槽配置
 * @param model 数据模型
 * @param slotScope 插槽作用域
 */
const renderSlotContent = (slotContent: any, model: any, slotScope: any) => {
  if (Array.isArray(slotContent)) {
    // 如果插槽内容是数组，递归渲染
    return slotContent.map(childConfig => renderComponent(childConfig, model, slotScope))
  }
  // 如果插槽内容是单个配置，递归渲染
  return renderComponent(slotContent, model, slotScope)
}

/**
 * 渲染组件函数（递归支持嵌套）
 * @param config JSON 配置
 * @param model 数据模型
 * @param slotScope 插槽作用域
 * @returns 渲染的组件 VNode
 */
const renderComponent = (config: any, model: any, slotScope: any = {}) => {

  let { component, props: componentProps = {}, content = [], slots = {}, html, text } = config
  if (!slots) {
    slots = {}
  }
  if (!content) {
    content = []
  }
  if (!componentProps) {
    componentProps = {}
  }
  // 动态解析组件
  const Component = resolveComp(component)

  // 解析属性
  const parsedProps = parseProps(componentProps, model)

  // 动态插槽解析
  const slotNodes: Record<string, any> = {}
  for (const [slotName, slotContent] of Object.entries(slots)) {
    slotNodes[slotName] = (slotScopeData: any) =>
      renderSlotContent(slotContent, model, { ...slotScope, ...slotScopeData })
  }

  // 渲染组件内容
  const renderContent = () => {
    // 如果配置了 `html`，直接渲染为 HTML 内容
    if (html) {
      return h(Component, { innerHTML: typeof html === 'string' ? html : getByPath(model, html) })
    }

    // 如果配置了 `text`，直接渲染为文本内容
    if (text) {
      return typeof text === 'string' ? text : getByPath(model, text)
    }

    // 如果配置了 `content`，递归渲染子组件
    if (Array.isArray(content)) {
      return content.map((childConfig: any) => renderComponent(childConfig, model, slotScope))
    }

    return null
  }

  // 渲染组件
  return h(Component, parsedProps, {
    ...slotNodes,
    default: renderContent,
  })
}
</script>

<template>
  <Component :is="renderComponent(config, model)" />
</template>
