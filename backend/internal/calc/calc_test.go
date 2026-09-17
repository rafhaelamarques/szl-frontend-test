package calc

import (
	"errors"
	"math"
	"testing"
)

func TestCompute(t *testing.T) {
	tests := []struct {
		name    string
		op      Operation
		a, b    float64
		want    float64
		wantErr error
	}{
		{name: "add", op: OpAdd, a: 2, b: 3, want: 5},
		{name: "subtract", op: OpSubtract, a: 5, b: 8, want: -3},
		{name: "multiply", op: OpMultiply, a: 4, b: 2.5, want: 10},
		{name: "divide", op: OpDivide, a: 10, b: 4, want: 2.5},
		{name: "divide by zero", op: OpDivide, a: 1, b: 0, wantErr: ErrDivisionByZero},
		{name: "power", op: OpPower, a: 2, b: 10, want: 1024},
		{name: "sqrt", op: OpSqrt, a: 9, want: 3},
		{name: "sqrt of zero", op: OpSqrt, a: 0, want: 0},
		{name: "sqrt negative", op: OpSqrt, a: -1, wantErr: ErrNegativeSqrt},
		{name: "percentage of", op: OpPercentage, a: 25, b: 200, want: 50},
		{name: "unknown op", op: Operation("nope"), a: 1, b: 1, wantErr: ErrUnknownOperation},
		{name: "power overflow", op: OpPower, a: 10, b: 1000, wantErr: ErrNonFiniteResult},
		{name: "nan operand", op: OpAdd, a: math.NaN(), b: 1, wantErr: ErrNonFiniteResult},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Compute(tt.op, tt.a, tt.b)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Compute() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Compute() unexpected error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("Compute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequiresB(t *testing.T) {
	if RequiresB(OpSqrt) {
		t.Fatal("sqrt must not require b")
	}
	if !RequiresB(OpAdd) {
		t.Fatal("add must require b")
	}
}
