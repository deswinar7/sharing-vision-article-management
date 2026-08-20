import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router'
import { afterEach, vi } from 'vitest'
import { PostsPage } from './PostsPage'

function renderPage() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  return render(<QueryClientProvider client={client}><MemoryRouter><PostsPage /></MemoryRouter></QueryClientProvider>)
}

afterEach(() => vi.unstubAllGlobals())

it('shows the loading state', () => {
  vi.stubGlobal('fetch', vi.fn(() => new Promise(() => undefined)))
  renderPage()
  expect(screen.getByText('Loading articles...')).toBeInTheDocument()
})

it('shows an error state and retry action', async () => {
  vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response(JSON.stringify({ error: { message: 'Unavailable' } }), { status: 500, headers: { 'Content-Type': 'application/json' } })))
  renderPage()
  expect(await screen.findByText('Failed to load articles.')).toBeInTheDocument()
  expect(screen.getByRole('button', { name: /retry/i })).toBeInTheDocument()
})

