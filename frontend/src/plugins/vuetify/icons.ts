import type { IconAliases } from 'vuetify'
import { Icon } from '@iconify/vue'
import { aliases } from 'vuetify/iconsets/mdi'
import { h } from 'vue'

const alertTypeIcon = {
  success: 'mdi-check-circle-outline',
  info: 'mdi-information-outline',
  warning: 'mdi-alert-outline',
  error: 'mdi-alert-circle-outline',
}
const modifiedAliases: IconAliases = Object.assign(aliases, alertTypeIcon)
export const iconify = {
  component: (props: any) => h(Icon, props),
}

export const icons: any = {
  defaultSet: 'iconify',
  mergedAliases: modifiedAliases,
  sets: {
    iconify,
  },
}
