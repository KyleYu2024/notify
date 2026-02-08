import { createApp } from 'vue'
import 'virtual:uno.css'
import 'vue-toast-notification/dist/theme-bootstrap.css'
import vuetify from '@/plugins/vuetify'
import 'dayjs/locale/zh-cn' // import locale
import App from './App.vue'
import VuetifyUseDialog from 'vuetify-use-dialog'
import ToastPlugin from 'vue-toast-notification'
import { createPinia } from 'pinia'
import router from '@/router'

const app = createApp(App)
app.use(ToastPlugin, {
  position: 'bottom-right',
})
app.use(vuetify)
app.use(VuetifyUseDialog)
app.use(createPinia())
app.use(router)

app.mount('#app')
