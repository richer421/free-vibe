import type { MenuProps } from 'antd'

import { appRoutes } from '../router/routes'

export function useAppNavigation(): MenuProps['items'] {
  return appRoutes.map((route) => ({
    key: route.key,
    icon: route.icon,
    label: route.label,
  }))
}
