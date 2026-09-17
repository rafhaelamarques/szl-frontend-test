import { describe, expect, it } from 'vitest'
import { appendDigit, formatResult } from './format'

describe('formatResult', () => {
  it('keeps small integers exact', () => {
    expect(formatResult(5)).toBe('5')
  })

  it('trims float noise', () => {
    expect(formatResult(0.1 + 0.2)).toBe('0.3')
  })
})

describe('appendDigit', () => {
  it('replaces the display', () => {
    expect(appendDigit('12', '7', true)).toBe('7')
  })

  it('avoids a second decimal', () => {
    expect(appendDigit('1.2', '.', false)).toBe('1.2')
  })

  it('replaces a lone zero', () => {
    expect(appendDigit('0', '4', false)).toBe('4')
  })
})
