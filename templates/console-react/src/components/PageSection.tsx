import { Card, Typography } from 'antd'
import type { CSSProperties, ReactNode } from 'react'

export interface PageSectionProps {
  title: string
  description?: string
  extra?: ReactNode
  children: ReactNode
  style?: CSSProperties
}

export function PageSection({
  title,
  description,
  extra,
  children,
  style,
}: PageSectionProps) {
  return (
    <Card
      style={style}
      styles={{
        body: { padding: 24 },
        header: {
          minHeight: 'auto',
          padding: '20px 24px 0',
          borderBottom: 'none',
        },
      }}
      title={
        <div className="flex items-start justify-between gap-4">
          <div>
            <Typography.Title level={4} style={{ margin: 0, color: '#0f172a' }}>
              {title}
            </Typography.Title>
            {description ? (
              <Typography.Paragraph
                style={{ margin: '6px 0 0', color: '#64748b' }}
              >
                {description}
              </Typography.Paragraph>
            ) : null}
          </div>
          {extra}
        </div>
      }
    >
      {children}
    </Card>
  )
}
