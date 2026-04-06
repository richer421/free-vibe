import type { ThemeConfig } from 'antd'

export const themeConfig: ThemeConfig = {
  token: {
    colorPrimary: '#2563eb',
    colorBgBase: '#edf2f7',
    colorBgLayout: '#edf2f7',
    colorBgContainer: '#ffffff',
    colorBorderSecondary: '#dbe4f0',
    colorText: '#0f172a',
    colorTextSecondary: '#64748b',
    borderRadius: 18,
    fontFamily:
      '"IBM Plex Sans", "Segoe UI", "PingFang SC", "Hiragino Sans GB", sans-serif',
    boxShadowTertiary: '0 20px 50px rgba(15, 23, 42, 0.08)',
  },
  components: {
    Layout: {
      bodyBg: '#edf2f7',
      headerBg: 'rgba(255, 255, 255, 0.82)',
      siderBg: '#0f172a',
      triggerBg: '#0f172a',
      triggerColor: '#f8fafc',
    },
    Menu: {
      darkItemBg: '#0f172a',
      darkItemSelectedBg: '#2563eb',
      darkItemSelectedColor: '#eff6ff',
      darkItemHoverBg: 'rgba(148, 163, 184, 0.12)',
      darkItemColor: '#cbd5e1',
    },
    Button: {
      fontWeight: 600,
      primaryShadow: '0 10px 24px rgba(37, 99, 235, 0.22)',
    },
    Card: {
      headerFontSize: 16,
    },
  },
}
