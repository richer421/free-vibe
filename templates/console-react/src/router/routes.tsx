import BellOutlined from '@ant-design/icons/BellOutlined'
import DashboardOutlined from '@ant-design/icons/DashboardOutlined'
import SettingOutlined from '@ant-design/icons/SettingOutlined'
import TeamOutlined from '@ant-design/icons/TeamOutlined'

import type { AppRouteModel } from '../types/navigation'

export const appRoutes: AppRouteModel[] = [
  {
    key: 'overview',
    path: '/',
    label: 'Overview',
    icon: <DashboardOutlined />,
  },
  {
    key: 'members',
    path: '/members',
    label: 'Members',
    icon: <TeamOutlined />,
  },
  {
    key: 'alerts',
    path: '/alerts',
    label: 'Alerts',
    icon: <BellOutlined />,
  },
  {
    key: 'settings',
    path: '/settings',
    label: 'Settings',
    icon: <SettingOutlined />,
  },
]
