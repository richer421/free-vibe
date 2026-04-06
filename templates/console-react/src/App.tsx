import { Spin } from 'antd'
import { Suspense, lazy } from 'react'

import { AdminLayout } from './layouts/AdminLayout'

const Dashboard = lazy(() => import('./pages/dashboard'))

function App() {
  return (
    <AdminLayout>
      <Suspense
        fallback={
          <div className="flex min-h-[320px] items-center justify-center">
            <Spin size="large" />
          </div>
        }
      >
        <Dashboard />
      </Suspense>
    </AdminLayout>
  )
}

export default App
