import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'

import App from './App'

describe('App', () => {
  it('renders the admin shell dashboard', async () => {
    render(<App />)

    expect(
      await screen.findByRole('heading', { name: 'Operational Overview' }),
    ).toBeInTheDocument()
    expect(await screen.findByText('Workspace')).toBeInTheDocument()
    expect(await screen.findByText('Quick Actions')).toBeInTheDocument()
    expect(await screen.findByText('Recent Activity')).toBeInTheDocument()
  })
})
