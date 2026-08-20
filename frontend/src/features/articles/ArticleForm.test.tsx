import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter } from 'react-router'
import { vi } from 'vitest'
import { ArticleForm } from './ArticleForm'

const validContent = 'A'.repeat(200)

function renderForm(onSubmit = vi.fn().mockResolvedValue(undefined)) {
  render(<MemoryRouter><ArticleForm onSubmit={onSubmit} isSubmitting={false} /></MemoryRouter>)
  return onSubmit
}

function fillValidForm() {
  fireEvent.change(screen.getByLabelText('Title'), { target: { value: 'A descriptive article title' } })
  fireEvent.change(screen.getByLabelText('Category'), { target: { value: 'News' } })
  fireEvent.change(screen.getByLabelText('Content'), { target: { value: validContent } })
}

it('shows field validation for invalid values', async () => {
  renderForm()
  fireEvent.click(screen.getByRole('button', { name: 'Publish' }))
  expect(await screen.findByText('Title must contain at least 20 characters')).toBeInTheDocument()
  expect(screen.getByText('Content must contain at least 200 characters')).toBeInTheDocument()
  expect(screen.getByText('Category must contain at least 3 characters')).toBeInTheDocument()
})

it.each([['Publish', 'publish'], ['Save as Draft', 'draft']] as const)('submits %s with the correct status', async (button, status) => {
  const onSubmit = renderForm()
  fillValidForm()
  fireEvent.click(screen.getByRole('button', { name: button }))
  await waitFor(() => expect(onSubmit).toHaveBeenCalledWith(expect.objectContaining({ status })))
})

