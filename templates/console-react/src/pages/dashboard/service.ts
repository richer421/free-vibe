import { appConfig } from '../../config/app'
import { formatDeltaLabel } from '../../utils/formatDeltaLabel'

export interface DashboardQuickActionModel {
  title: string
  description: string
}

export interface DashboardActivityModel {
  title: string
  meta: string
  tag: string
}

export interface DashboardQueueItemModel {
  label: string
  value: number
}

export interface DashboardPageModel {
  stats: Array<{
    label: string
    value: string
    helper: string
    tone?: 'default' | 'accent'
  }>
  quickActions: DashboardQuickActionModel[]
  activities: DashboardActivityModel[]
  deliveryQueue: DashboardQueueItemModel[]
  operatorNotes: string[]
  recommendedFiles: string[]
}

export function getDashboardPageModel(): DashboardPageModel {
  return {
    stats: [
      {
        label: 'Revenue',
        value: '¥128K',
        helper: formatDeltaLabel(12.4),
      },
      {
        label: 'Active Users',
        value: '2,431',
        helper: formatDeltaLabel(8.1),
      },
      {
        label: 'Orders',
        value: '312',
        helper: formatDeltaLabel(5.7),
      },
      {
        label: 'Alerts',
        value: '07',
        helper: 'Needs review',
        tone: 'accent',
      },
    ],
    quickActions: [
      {
        title: 'Create project',
        description: 'Start a new business module with standard ownership fields.',
      },
      {
        title: 'Invite member',
        description: 'Grant workspace access for product, ops, and engineering roles.',
      },
      {
        title: 'Export report',
        description: 'Generate a weekly operating snapshot for stakeholders.',
      },
    ],
    activities: [
      {
        title: 'Membership Center release approved',
        meta: '12 minutes ago',
        tag: 'Release',
      },
      {
        title: 'Growth campaign warning acknowledged',
        meta: '27 minutes ago',
        tag: 'Alert',
      },
      {
        title: 'Order platform backlog trimmed',
        meta: '43 minutes ago',
        tag: 'Ops',
      },
    ],
    deliveryQueue: [
      {
        label: 'Release readiness',
        value: 84,
      },
      {
        label: 'QA completion',
        value: 71,
      },
      {
        label: 'Monitoring coverage',
        value: 92,
      },
    ],
    operatorNotes: [
      'Replace static metrics with your own domain data source.',
      'Split this page into route-level modules when dashboard scope grows.',
      'Keep permissions at the shell boundary instead of inside leaf widgets.',
    ],
    recommendedFiles: [
      'src/pages/dashboard/index.tsx',
      'src/pages/dashboard/service.ts',
      'src/router/routes.tsx',
      appConfig.recommendedFile,
    ],
  }
}
