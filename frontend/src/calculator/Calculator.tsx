import { useState } from 'react'
import { calculate } from '../api/calculate'
import { appendDigit, formatResult } from './format'
import {
  OP_LABEL,
  type BinaryOperation,
  type Operation,
} from './operations'
import './Calculator.css'

type Pending = { a: number; operation: BinaryOperation }

export function Calculator() {
  const [display, setDisplay] = useState('0')
  const [pending, setPending] = useState<Pending | null>(null)
  const [replace, setReplace] = useState(true)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const current = () => Number(display)

  async function run(operation: Operation, a: number, b?: number) {
    setLoading(true)
    setError(null)
    try {
      const result = await calculate(operation, a, b)
      setDisplay(formatResult(result))
      setPending(null)
      setReplace(true)
    } catch (err) {
      const message = err instanceof Error ? err.message : 'request failed'
      setError(message)
    } finally {
      setLoading(false)
    }
  }

  function onDigit(digit: string) {
    if (loading) return
    setError(null)
    setDisplay((d) => appendDigit(d, digit, replace))
    setReplace(false)
  }

  function onClear() {
    setDisplay('0')
    setPending(null)
    setReplace(true)
    setError(null)
  }

  function onBinary(operation: BinaryOperation) {
    if (loading) return
    setError(null)
    if (pending && !replace) {
      void (async () => {
        setLoading(true)
        try {
          const result = await calculate(pending.operation, pending.a, current())
          setDisplay(formatResult(result))
          setPending({ a: result, operation })
          setReplace(true)
        } catch (err) {
          const message = err instanceof Error ? err.message : 'request failed'
          setError(message)
        } finally {
          setLoading(false)
        }
      })()
      return
    }
    setPending({ a: current(), operation })
    setReplace(true)
  }

  function onEquals() {
    if (loading || !pending) return
    void run(pending.operation, pending.a, current())
  }

  function onSqrt() {
    if (loading) return
    void run('sqrt', current())
  }

  return (
    <div className="calc">
      <div className="calc__display" role="status" aria-live="polite">
        {display}
      </div>
      {error ? (
        <p className="calc__error" role="alert">
          {error}
        </p>
      ) : (
        <p className="calc__hint">
          {pending
            ? `${formatResult(pending.a)} ${OP_LABEL[pending.operation]} …`
            : loading
              ? 'Calculating…'
              : 'Ready'}
        </p>
      )}
      <div className="calc__keys">
        <button type="button" className="calc__btn calc__btn--util" onClick={onClear} aria-label="clear">
          AC
        </button>
        <button type="button" className="calc__btn calc__btn--op" onClick={onSqrt} aria-label="square root">
          {OP_LABEL.sqrt}
        </button>
        <button
          type="button"
          className="calc__btn calc__btn--op"
          onClick={() => onBinary('percentage')}
          aria-label="percent of"
        >
          {OP_LABEL.percentage}
        </button>
        <button
          type="button"
          className="calc__btn calc__btn--op"
          onClick={() => onBinary('divide')}
          aria-label="divide"
        >
          {OP_LABEL.divide}
        </button>

        <button type="button" className="calc__btn" onClick={() => onDigit('7')}>
          7
        </button>
        <button type="button" className="calc__btn" onClick={() => onDigit('8')}>
          8
        </button>
        <button type="button" className="calc__btn" onClick={() => onDigit('9')}>
          9
        </button>
        <button
          type="button"
          className="calc__btn calc__btn--op"
          onClick={() => onBinary('multiply')}
          aria-label="multiply"
        >
          {OP_LABEL.multiply}
        </button>

        <button type="button" className="calc__btn" onClick={() => onDigit('4')}>
          4
        </button>
        <button type="button" className="calc__btn" onClick={() => onDigit('5')}>
          5
        </button>
        <button type="button" className="calc__btn" onClick={() => onDigit('6')}>
          6
        </button>
        <button
          type="button"
          className="calc__btn calc__btn--op"
          onClick={() => onBinary('subtract')}
          aria-label="subtract"
        >
          {OP_LABEL.subtract}
        </button>

        <button type="button" className="calc__btn" onClick={() => onDigit('1')}>
          1
        </button>
        <button type="button" className="calc__btn" onClick={() => onDigit('2')}>
          2
        </button>
        <button type="button" className="calc__btn" onClick={() => onDigit('3')}>
          3
        </button>
        <button type="button" className="calc__btn calc__btn--op" onClick={() => onBinary('add')} aria-label="add">
          {OP_LABEL.add}
        </button>

        <button type="button" className="calc__btn" onClick={() => onDigit('0')}>
          0
        </button>
        <button type="button" className="calc__btn" onClick={() => onDigit('.')} aria-label="decimal point">
          .
        </button>
        <button
          type="button"
          className="calc__btn calc__btn--op"
          onClick={() => onBinary('power')}
          aria-label="exponentiation"
        >
          {OP_LABEL.power}
        </button>
        <button type="button" className="calc__btn calc__btn--eq" onClick={onEquals} aria-label="equals" disabled={loading}>
          =
        </button>
      </div>
    </div>
  )
}
