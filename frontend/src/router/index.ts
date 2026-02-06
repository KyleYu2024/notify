import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/store/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/pages/LoginPage.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/',
      component: () => import('@/components/Layout.vue'),
      meta: { requiresAuth: true },
      redirect: '/dashboard',
      children: [
        {
          path: '/dashboard',
          name: 'Dashboard',
          component: () => import('@/pages/DashboardPage.vue'),
          meta: {
            requiresAuth: true,
            title: '仪表板',
          },
        },
        {
          path: '/dashboard',
          redirect: '/',
        },
        {
          path: '/apps',
          name: 'Apps',
          component: () => import('@/pages/AppsPage.vue'),
          meta: {
            requiresAuth: true,
            title: '通知应用',
          },
        },
        {
          path: '/notifiers',
          name: 'Notifiers',
          component: () => import('@/pages/NotifiersPage.vue'),
          meta: {
            requiresAuth: true,
            title: '通知服务配置',
          },
        },
        {
          path: '/templates',
          name: 'Templates',
          component: () => import('@/pages/TemplatesPage.vue'),
          meta: {
            requiresAuth: true,
            title: '模板管理',
          },
        },
        {
          path: '/logs',
          name: 'Logs',
          component: () => import('@/pages/LogsPage.vue'),
          meta: {
            requiresAuth: true,
            title: '日志',
          },
        },
        {
          path: '/plugins',
          name: 'Plugins',
          component: () => import('@/pages/PluginsPage.vue'),
          meta: {
            requiresAuth: true,
            title: '插件管理',
          },
        },
      ],
    },
    {
      // 404 页面
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
})

// 路由守卫
router.beforeEach(async (to, _, next) => {
  const authStore = useAuthStore()

  // 如果是登录页面，直接通过
  if (to.path === '/login') {
    // 如果已经登录，重定向到主页
    if (authStore.isAuthenticated) {
      next('/')
    } else {
      next()
    }
    return
  }

  // 检查是否需要认证
  if (to.meta.requiresAuth) {
    // 如果未登录，重定向到登录页
    if (!authStore.isAuthenticated) {
      // 尝试从localStorage恢复认证状态
      authStore.restoreAuth()

      // 如果仍未认证，重定向到登录页
      if (!authStore.isAuthenticated) {
        next('/login')
        return
      }
    }
  }

  next()
})

export default router

export function navigateTo(path: string) {
  router.push(path)
}
