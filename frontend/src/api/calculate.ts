import type { Operation } from '../calculator/operations'

export type CalculateSuccess = { result: number }
export type CalculateFailure = { error: string }

export class ApiError extends Error {
  readonly status: number

  constructor(message: string, status: number) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

export async function calculate(
  operation: Operation,
  a: number,
  b?: number,
): Promise<number> {
  const body: { operation: Operation; a: number; b?: number } = { operation, a }
  if (b !== undefined) {
    body.b = b
  }

  const res = await fetch('/api/v1/calculate', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })

  const data = (await res.json()) as CalculateSuccess | CalculateFailure
  if (!res.ok) {
    const message = 'error' in data ? data.error : 'request failed'
    throw new ApiError(message, res.status)
  }
  if (!('result' in data) || typeof data.result !== 'number') {
    throw new ApiError('malformed response', res.status)
  }
  return data.result
}
