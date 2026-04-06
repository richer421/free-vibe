import PlusOutlined from '@ant-design/icons/PlusOutlined'
import type { ComponentType } from 'react'

export interface AppConfigModel {
  shellLabel: string
  moduleName: string
  moduleDescription: string
  recommendedActionTitle: string
  recommendedActionDescription: string
  recommendedFile: string
  primaryActionIcon: ComponentType
}

export const appConfig: AppConfigModel = {
  shellLabel: 'Console React',
  moduleName: '__MODULE_NAME__',
  moduleDescription: 'A stable React + Ant Design admin shell for internal products.',
  recommendedActionTitle: 'Recommended first extension',
  recommendedActionDescription:
    'Add routing and permission boundaries before wiring business pages.',
  recommendedFile: 'src/config/app.ts',
  primaryActionIcon: PlusOutlined,
}
