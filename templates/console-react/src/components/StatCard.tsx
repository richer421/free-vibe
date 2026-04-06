import { Card, Tag, Typography } from 'antd'

export interface StatCardModel {
  label: string
  value: string
  helper: string
  tone?: 'default' | 'accent'
}

export function StatCard({
  label,
  value,
  helper,
  tone = 'default',
}: StatCardModel) {
  const isAccent = tone === 'accent'

  return (
    <Card
      style={{
        background: isAccent ? '#0f172a' : '#ffffff',
        borderColor: isAccent ? '#0f172a' : '#e2e8f0',
        boxShadow: isAccent
          ? '0 18px 40px rgba(15, 23, 42, 0.12)'
          : '0 10px 28px rgba(148, 163, 184, 0.12)',
      }}
      styles={{ body: { padding: 20 } }}
    >
      <Typography.Text
        style={{ color: isAccent ? '#cbd5e1' : '#64748b', fontSize: 13 }}
      >
        {label}
      </Typography.Text>
      <div className="mt-4 flex items-end justify-between gap-3">
        <Typography.Title
          level={2}
          style={{ margin: 0, color: isAccent ? '#f8fafc' : '#0f172a' }}
        >
          {value}
        </Typography.Title>
        <Tag color={isAccent ? 'processing' : 'blue'}>{helper}</Tag>
      </div>
    </Card>
  )
}
