export function formatResult(n: number): string {
  if (!Number.isFinite(n)) {
    return 'Error'
  }
  const asInt = Number.isInteger(n)
  if (asInt && Math.abs(n) < 1e12) {
    return String(n)
  }
  // display rounding only; JSON result stays full float64
  const trimmed = Number(n.toPrecision(12))
  return String(trimmed)
}

export function appendDigit(current: string, digit: string, replace: boolean): string {
  if (replace) {
    return digit === '.' ? '0.' : digit
  }
  if (digit === '.' && current.includes('.')) {
    return current
  }
  if (current === '0' && digit !== '.') {
    return digit
  }
  return current + digit
}
