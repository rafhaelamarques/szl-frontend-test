package calc

import (
	"errors"
	"math"
)

type Operation string

const (
	OpAdd        Operation = "add"
	OpSubtract   Operation = "subtract"
	OpMultiply   Operation = "multiply"
	OpDivide     Operation = "divide"
	OpPower      Operation = "power"
	OpSqrt       Operation = "sqrt"
	OpPercentage Operation = "percentage"
)

var (
	ErrUnknownOperation = errors.New("unknown operation")
	ErrDivisionByZero   = errors.New("division by zero")
	ErrNegativeSqrt     = errors.New("square root of negative number")
	ErrMissingOperand   = errors.New("missing operand")
	ErrNonFiniteResult  = errors.New("result is not a finite number")
)

// RequiresB reports whether the operation needs operand b.
func RequiresB(op Operation) bool {
	return op != OpSqrt
}

// Compute evaluates a binary (or unary, for sqrt) arithmetic operation.
// percentage is defined as (a / 100) * b — "a percent of b".
func Compute(op Operation, a, b float64) (float64, error) {
	if !isFinite(a) || (RequiresB(op) && !isFinite(b)) {
		return 0, ErrNonFiniteResult
	}

	var result float64
	switch op {
	case OpAdd:
		result = a + b
	case OpSubtract:
		result = a - b
	case OpMultiply:
		result = a * b
	case OpDivide:
		if b == 0 {
			return 0, ErrDivisionByZero
		}
		result = a / b
	case OpPower:
		result = math.Pow(a, b)
	case OpSqrt:
		if a < 0 {
			return 0, ErrNegativeSqrt
		}
		result = math.Sqrt(a)
	case OpPercentage:
		result = (a / 100) * b
	default:
		return 0, ErrUnknownOperation
	}

	if !isFinite(result) {
		return 0, ErrNonFiniteResult
	}
	return result, nil
}

func isFinite(n float64) bool {
	return !math.IsNaN(n) && !math.IsInf(n, 0)
}
