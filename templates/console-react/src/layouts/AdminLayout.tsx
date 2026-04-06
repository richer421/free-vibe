import { Avatar, Button, Layout, Menu, Space, Typography } from 'antd'
import type { ReactNode } from 'react'

import { appConfig } from '../config/app'
import { useAppNavigation } from '../hooks/useAppNavigation'

export interface AdminLayoutProps {
  children: ReactNode
}

export function AdminLayout({ children }: AdminLayoutProps) {
  const navigationItems = useAppNavigation()
  const PrimaryActionIcon = appConfig.primaryActionIcon

  return (
    <Layout style={{ minHeight: '100vh', background: '#edf2f7' }}>
      <Layout.Sider width={264} theme="dark">
        <div className="flex h-full flex-col px-5 py-6">
          <div className="space-y-4">
            <span className="inline-flex rounded-full border border-blue-400/20 bg-blue-400/10 px-3 py-1 text-[11px] font-semibold uppercase tracking-[0.28em] text-blue-100">
              {appConfig.shellLabel}
            </span>
            <div>
              <Typography.Title level={4} style={{ margin: 0, color: '#f8fafc' }}>
                {appConfig.moduleName}
              </Typography.Title>
              <Typography.Paragraph
                style={{ margin: '8px 0 0', color: '#94a3b8' }}
              >
                {appConfig.moduleDescription}
              </Typography.Paragraph>
            </div>
          </div>

          <Menu
            mode="inline"
            theme="dark"
            defaultSelectedKeys={['overview']}
            items={navigationItems}
            style={{
              marginTop: 32,
              background: 'transparent',
              borderInlineEnd: 'none',
            }}
          />

          <div className="mt-auto rounded-3xl border border-blue-400/20 bg-blue-400/10 p-4">
            <Typography.Text style={{ color: '#dbeafe', fontWeight: 600 }}>
              {appConfig.recommendedActionTitle}
            </Typography.Text>
            <Typography.Paragraph
              style={{ margin: '8px 0 0', color: '#bfdbfe' }}
            >
              {appConfig.recommendedActionDescription}
            </Typography.Paragraph>
          </div>
        </div>
      </Layout.Sider>

      <Layout style={{ background: 'transparent' }}>
        <Layout.Header
          style={{
            height: 'auto',
            padding: '24px 32px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            borderBottom: '1px solid rgba(226, 232, 240, 0.9)',
            background: 'rgba(255, 255, 255, 0.82)',
            backdropFilter: 'blur(14px)',
          }}
        >
          <div>
            <Typography.Text
              style={{
                fontSize: 12,
                fontWeight: 600,
                letterSpacing: '0.24em',
                textTransform: 'uppercase',
                color: '#64748b',
              }}
            >
              Workspace
            </Typography.Text>
            <Typography.Title level={2} style={{ margin: '6px 0 0', color: '#0f172a' }}>
              Operational Overview
            </Typography.Title>
            <Typography.Paragraph
              style={{ margin: '8px 0 0', color: '#64748b' }}
            >
              A default console shell with room for metrics, actions, and activity.
            </Typography.Paragraph>
          </div>

          <Space size="middle">
            <Button>Export Report</Button>
            <Button type="primary" icon={<PrimaryActionIcon />}>
              New Task
            </Button>
            <Avatar
              style={{
                backgroundColor: '#dbeafe',
                color: '#1d4ed8',
                fontWeight: 700,
              }}
            >
              FV
            </Avatar>
          </Space>
        </Layout.Header>

        <Layout.Content style={{ padding: '32px' }}>{children}</Layout.Content>
      </Layout>
    </Layout>
  )
}
