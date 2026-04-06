import type { ReactNode } from 'react'

export interface AppRouteModel {
  key: string
  path: string
  label: string
  icon: ReactNode
}
