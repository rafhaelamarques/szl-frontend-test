export const OPERATIONS = [
  'add',
  'subtract',
  'multiply',
  'divide',
  'power',
  'sqrt',
  'percentage',
] as const

export type Operation = (typeof OPERATIONS)[number]

export const BINARY_OPERATIONS = [
  'add',
  'subtract',
  'multiply',
  'divide',
  'power',
  'percentage',
] as const satisfies readonly Operation[]

export type BinaryOperation = (typeof BINARY_OPERATIONS)[number]

export const OP_LABEL: Record<Operation, string> = {
  add: '+',
  subtract: '−',
  multiply: '×',
  divide: '÷',
  power: 'xʸ',
  sqrt: '√',
  percentage: '%',
}
