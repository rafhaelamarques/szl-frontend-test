import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { Calculator } from './Calculator'

afterEach(() => {
  vi.unstubAllGlobals()
})

function mockFetchOk(result: number) {
  vi.stubGlobal(
    'fetch',
    vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({ result }),
    }),
  )
}

function mockFetchError(message: string, status = 422) {
  vi.stubGlobal(
    'fetch',
    vi.fn().mockResolvedValue({
      ok: false,
      status,
      json: async () => ({ error: message }),
    }),
  )
}

describe('Calculator', () => {
  it('calls the API on equals and shows the result', async () => {
    mockFetchOk(5)
    const user = userEvent.setup()
    render(<Calculator />)

    await user.click(screen.getByRole('button', { name: '1' }))
    await user.click(screen.getByRole('button', { name: '0' }))
    await user.click(screen.getByRole('button', { name: 'divide' }))
    await user.click(screen.getByRole('button', { name: '2' }))
    await user.click(screen.getByRole('button', { name: 'equals' }))

    expect(await screen.findByRole('status')).toHaveTextContent('5')
    expect(fetch).toHaveBeenCalledWith(
      '/api/v1/calculate',
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ operation: 'divide', a: 10, b: 2 }),
      }),
    )
  })

  it('surfaces API errors', async () => {
    mockFetchError('division by zero')
    const user = userEvent.setup()
    render(<Calculator />)

    await user.click(screen.getByRole('button', { name: '8' }))
    await user.click(screen.getByRole('button', { name: 'divide' }))
    await user.click(screen.getByRole('button', { name: '0' }))
    await user.click(screen.getByRole('button', { name: 'equals' }))

    expect(await screen.findByRole('alert')).toHaveTextContent('division by zero')
  })

  it('does not call the API without a pending operation', async () => {
    mockFetchOk(0)
    const user = userEvent.setup()
    render(<Calculator />)

    await user.click(screen.getByRole('button', { name: 'equals' }))

    expect(fetch).not.toHaveBeenCalled()
  })

  it('sends sqrt with only operand a', async () => {
    mockFetchOk(3)
    const user = userEvent.setup()
    render(<Calculator />)

    await user.click(screen.getByRole('button', { name: '9' }))
    await user.click(screen.getByRole('button', { name: 'square root' }))

    expect(await screen.findByRole('status')).toHaveTextContent('3')
    expect(fetch).toHaveBeenCalledWith(
      '/api/v1/calculate',
      expect.objectContaining({
        body: JSON.stringify({ operation: 'sqrt', a: 9 }),
      }),
    )
  })
})
