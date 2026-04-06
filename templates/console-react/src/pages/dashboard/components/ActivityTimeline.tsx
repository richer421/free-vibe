import ClockCircleOutlined from '@ant-design/icons/ClockCircleOutlined'
import { List, Tag } from 'antd'

import type { DashboardActivityModel } from '../service'

export interface ActivityTimelineProps {
  activities: DashboardActivityModel[]
}

export function ActivityTimeline({ activities }: ActivityTimelineProps) {
  return (
    <List
      itemLayout="horizontal"
      dataSource={activities}
      renderItem={(activity) => (
        <List.Item style={{ paddingInline: 0 }}>
          <List.Item.Meta
            avatar={
              <div className="flex h-10 w-10 items-center justify-center rounded-2xl bg-blue-100 text-blue-600">
                <ClockCircleOutlined />
              </div>
            }
            title={<span className="font-semibold text-slate-900">{activity.title}</span>}
            description={
              <div className="flex items-center gap-2 text-sm text-slate-500">
                <span>{activity.meta}</span>
                <Tag>{activity.tag}</Tag>
              </div>
            }
          />
        </List.Item>
      )}
    />
  )
}
