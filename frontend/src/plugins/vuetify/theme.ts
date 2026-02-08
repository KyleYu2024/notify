import type { VuetifyOptions } from 'vuetify'

// 默认逻辑：读取本地存储或系统偏好，只在 day 和 night 之间选择
let defaultTheme = 'day'

try {
  const t = localStorage.getItem('theme')
  if (t === 'macosDark' || t === 'dark' || t === 'purple') {
    defaultTheme = 'night'
  } else if (t === 'google' || t === 'light') {
    defaultTheme = 'day'
  } else if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
    defaultTheme = 'night'
  }
} catch (e) {}

const theme: VuetifyOptions['theme'] = {
  defaultTheme,
  themes: {
    // === 白天模式 (基于 Google 风格) ===
    day: {
      dark: false,
      colors: {
        primary: '#1A73E8',      // Google Blue
        secondary: '#5F6368',    // Grey text
        'on-secondary': '#ffffff',
        success: '#1E8E3E',      // Google Green
        info: '#1A73E8',         
        
        // [!code focus] 修改为深橙色
        warning: '#F57C00',      // Deep Orange (原 #F9AB00 是黄色)
        
        error: '#D93025',        // Google Red
        'on-primary': '#FFFFFF',
        'on-success': '#FFFFFF',
        'on-warning': '#FFFFFF', // 深橙色背景配白字很清晰
        
        background: '#F0F2F5',   
        'on-background': '#202124', 
        surface: '#FFFFFF',
        'on-surface': '#202124',
        
        'grey-50': '#F8F9FA',
        'grey-100': '#F1F3F4',
        'grey-200': '#E8EAED',
        'grey-300': '#DADCE0',
        'grey-400': '#BDC1C6',
        'grey-500': '#9AA0A6',
        'grey-600': '#80868B',
        'grey-700': '#5F6368',
        'grey-800': '#3C4043',
        'grey-900': '#202124',
        
        'perfect-scrollbar-thumb': '#DADCE0',
        'skin-bordered-background': '#fff',
        'skin-bordered-surface': '#fff',
      },
      variables: {
        'code-color': '#d400ff',
        'overlay-scrim-background': '#202124',
        'overlay-scrim-opacity': 0.5,
        'border-color': '#DADCE0',
        'table-header-background': '#F8F9FA',
        'custom-background': '#F8F9FA',
        'shadow-key-umbra-opacity': 'rgba(0, 0, 0, 0.04)',
        'shadow-key-penumbra-opacity': 'rgba(0, 0, 0, 0.08)',
        'shadow-key-ambient-opacity': 'rgba(0, 0, 0, 0.04)',
      },
    },

    // === 夜晚模式 (基于 macOS 深色风格) ===
    night: {
      dark: true,
      colors: {
        primary: '#0A84FF',      // macOS Blue
        secondary: '#8E8E93',    // macOS Gray
        'on-secondary': '#fff',
        success: '#30D158',      // macOS Green
        info: '#0A84FF',
        
        // [!code focus] 修改为 macOS 系统橙色
        warning: '#FF9F0A',      // macOS Orange (原 #FFD60A 是黄色)
        
        error: '#FF453A',        // macOS Red
        'on-primary': '#FFFFFF',
        'on-success': '#FFFFFF',
        'on-warning': '#FFFFFF', // 橙色配白字
        
        background: '#0F0F0F',   // 稍微提升背景色 (原 #000000)
        'on-background': '#FFFFFF',
        surface: '#1C1C1E',      
        'on-surface': '#FFFFFF',

        'grey-50': '#1C1C1E',
        'grey-100': '#2C2C2E',
        'grey-200': '#3A3A3C',
        'grey-300': '#48484A',
        'grey-400': '#636366',
        'grey-500': '#8E8E93',
        'grey-600': '#AEAEB2',
        'grey-700': '#C7C7CC',
        'grey-800': '#D1D1D6',
        'grey-900': '#E5E5EA',

        'perfect-scrollbar-thumb': '#48484A',
        'skin-bordered-background': '#1C1C1E',
        'skin-bordered-surface': '#1C1C1E',
      },
      variables: {
        'code-color': '#d400ff',
        'overlay-scrim-background': '#000000',
        'overlay-scrim-opacity': 0.7,
        'border-color': '#2C2C2E', 
        'table-header-background': '#1C1C1E',
        'custom-background': '#000000',
        'shadow-key-umbra-opacity': 'rgba(0, 0, 0, 0.6)',
        'shadow-key-penumbra-opacity': 'rgba(0, 0, 0, 0.4)',
        'shadow-key-ambient-opacity': 'rgba(0, 0, 0, 0.3)',
      },
    },
  },
}

export default theme