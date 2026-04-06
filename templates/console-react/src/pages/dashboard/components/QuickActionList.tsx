import ArrowRightOutlined from '@ant-design/icons/ArrowRightOutlined'

import type { DashboardQuickActionModel } from '../service'

export interface QuickActionListProps {
  actions: DashboardQuickActionModel[]
}

export function QuickActionList({ actions }: QuickActionListProps) {
  return (
    <div className="space-y-3">
      {actions.map((action) => (
        <button
          key={action.title}
          type="button"
          className="flex w-full items-start justify-between rounded-3xl border border-slate-200 bg-slate-50 px-4 py-4 text-left transition hover:border-blue-300 hover:bg-blue-50"
        >
          <div>
            <div className="font-semibold text-slate-900">{action.title}</div>
            <div className="mt-1 text-sm leading-6 text-slate-500">
              {action.description}
            </div>
          </div>
          <ArrowRightOutlined className="mt-1 text-slate-400" />
        </button>
      ))}
    </div>
  )
}
