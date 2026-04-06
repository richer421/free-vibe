import FireOutlined from '@ant-design/icons/FireOutlined'
import RiseOutlined from '@ant-design/icons/RiseOutlined'
import { Button, Progress, Space, Tag, Typography } from 'antd'

import { PageSection } from '../../components/PageSection'
import { StatCard } from '../../components/StatCard'
import { ActivityTimeline } from './components/ActivityTimeline'
import { QuickActionList } from './components/QuickActionList'
import { useDashboardModel } from './hooks/useDashboardModel'

export function Dashboard() {
  const dashboardModel = useDashboardModel()

  return (
    <div className="space-y-6">
      <section className="grid gap-4 xl:grid-cols-4">
        {dashboardModel.stats.map((item) => (
          <StatCard key={item.label} {...item} />
        ))}
      </section>

      <section className="grid gap-6 2xl:grid-cols-[1.35fr_0.65fr]">
        <PageSection
          title="Performance Trend"
          description="Use this block as the default place for trend charts or KPI timelines."
          extra={<Tag color="blue">This week</Tag>}
          style={{
            borderColor: '#e2e8f0',
            boxShadow: '0 10px 28px rgba(148, 163, 184, 0.12)',
          }}
        >
          <div className="grid gap-4 lg:grid-cols-[1.15fr_0.85fr]">
            <div className="rounded-3xl border border-slate-200 bg-slate-50 p-5">
              <div className="flex items-end justify-between gap-4">
                <div>
                  <Typography.Text className="text-sm text-slate-500">
                    Gross merchandise volume
                  </Typography.Text>
                  <Typography.Title level={3} style={{ margin: '8px 0 0' }}>
                    ¥423,000
                  </Typography.Title>
                </div>
                <Space size="small" className="text-slate-500">
                  <RiseOutlined />
                  <span>Upward trend</span>
                </Space>
              </div>

              <div className="mt-6 flex h-40 items-end gap-3">
                {[38, 52, 47, 64, 72, 68, 86].map((height, index) => (
                  <div key={height} className="flex flex-1 flex-col items-center gap-2">
                    <div
                      className="w-full rounded-t-2xl bg-gradient-to-b from-blue-500 to-sky-300"
                      style={{ height: `${height}%` }}
                    />
                    <span className="text-xs text-slate-400">D{index + 1}</span>
                  </div>
                ))}
              </div>
            </div>

            <div className="space-y-3 rounded-3xl border border-slate-200 bg-white p-5">
              <Typography.Text className="text-sm font-semibold text-slate-600">
                Delivery queue
              </Typography.Text>
              {dashboardModel.deliveryQueue.map((item) => (
                <div key={item.label}>
                  <div className="mb-2 flex items-center justify-between text-sm text-slate-600">
                    <span>{item.label}</span>
                    <span>{item.value}%</span>
                  </div>
                  <Progress percent={item.value} showInfo={false} strokeColor="#2563eb" />
                </div>
              ))}
            </div>
          </div>
        </PageSection>

        <PageSection
          title="Quick Actions"
          description="Seed common operations here before you introduce routing."
          style={{
            borderColor: '#e2e8f0',
            boxShadow: '0 10px 28px rgba(148, 163, 184, 0.12)',
          }}
        >
          <QuickActionList actions={dashboardModel.quickActions} />
        </PageSection>
      </section>

      <section className="grid gap-6 2xl:grid-cols-[0.92fr_1.08fr]">
        <PageSection
          title="Recent Activity"
          description="Recent events give the template a default operational rhythm."
          style={{
            borderColor: '#e2e8f0',
            boxShadow: '0 10px 28px rgba(148, 163, 184, 0.12)',
          }}
        >
          <ActivityTimeline activities={dashboardModel.activities} />
        </PageSection>

        <PageSection
          title="Operator Notes"
          description="Keep a lightweight secondary block for guidance, handoff notes, or SLA hints."
          extra={<Button type="text">View archive</Button>}
          style={{
            borderColor: '#e2e8f0',
            boxShadow: '0 10px 28px rgba(148, 163, 184, 0.12)',
          }}
        >
          <div className="grid gap-4 lg:grid-cols-2">
            <div className="rounded-3xl border border-slate-200 bg-slate-50 p-5">
              <div className="flex items-center gap-2 text-sm font-semibold text-slate-700">
                <FireOutlined className="text-amber-500" />
                Attention points
              </div>
              <ul className="mt-4 space-y-3 text-sm leading-6 text-slate-600">
                {dashboardModel.operatorNotes.map((note) => (
                  <li key={note}>{note}</li>
                ))}
              </ul>
            </div>

            <div className="rounded-3xl border border-dashed border-slate-300 bg-white p-5">
              <Typography.Text className="text-sm font-semibold text-slate-700">
                Recommended next files
              </Typography.Text>
              <div className="mt-4 space-y-3">
                {dashboardModel.recommendedFiles.map((file) => (
                  <div
                    key={file}
                    className="rounded-2xl bg-slate-950 px-4 py-3 text-sm text-slate-100"
                  >
                    {file}
                  </div>
                ))}
              </div>
            </div>
          </div>
        </PageSection>
      </section>
    </div>
  )
}

export default Dashboard
