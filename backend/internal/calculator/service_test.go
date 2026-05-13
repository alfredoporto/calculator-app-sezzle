package calculator

import "testing"

func TestServiceCalculateCoreOperations(t *testing.T) {
	t.Parallel()

	service := NewService()

	tests := []struct {
		name      string
		operation Operation
		operands  []float64
		want      float64
	}{
		{name: "addition", operation: OperationAdd, operands: []float64{10, 2}, want: 12},
		{name: "subtraction", operation: OperationSubtract, operands: []float64{10, 2}, want: 8},
		{name: "multiplication", operation: OperationMultiply, operands: []float64{10, 2}, want: 20},
		{name: "division", operation: OperationDivide, operands: []float64{10, 2}, want: 5},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := service.Calculate(tt.operation, tt.operands)
			if err != nil {
				t.Fatalf("Calculate() error = %v", err)
			}

			if got != tt.want {
				t.Fatalf("Calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceCalculateAdvancedOperations(t *testing.T) {
	t.Parallel()

	service := NewService()

	tests := []struct {
		name      string
		operation Operation
		operands  []float64
		want      float64
	}{
		{name: "power", operation: OperationPower, operands: []float64{2, 3}, want: 8},
		{name: "square root", operation: OperationSquareRoot, operands: []float64{81}, want: 9},
		{name: "percentage", operation: OperationPercentage, operands: []float64{25}, want: 0.25},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := service.Calculate(tt.operation, tt.operands)
			if err != nil {
				t.Fatalf("Calculate() error = %v", err)
			}

			if got != tt.want {
				t.Fatalf("Calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceCalculateValidationErrors(t *testing.T) {
	t.Parallel()

	service := NewService()

	tests := []struct {
		name      string
		operation Operation
		operands  []float64
		wantCode  ErrorCode
	}{
		{name: "unsupported operation", operation: Operation("mod"), operands: []float64{10, 2}, wantCode: ErrorCodeInvalidOperation},
		{name: "missing binary operand", operation: OperationAdd, operands: []float64{10}, wantCode: ErrorCodeInvalidOperands},
		{name: "extra unary operand", operation: OperationSquareRoot, operands: []float64{9, 3}, wantCode: ErrorCodeInvalidOperands},
		{name: "division by zero", operation: OperationDivide, operands: []float64{10, 0}, wantCode: ErrorCodeDivisionByZero},
		{name: "negative square root", operation: OperationSquareRoot, operands: []float64{-1}, wantCode: ErrorCodeInvalidOperands},
		{name: "non finite result", operation: OperationPower, operands: []float64{10, 400}, wantCode: ErrorCodeNonFiniteResult},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := service.Calculate(tt.operation, tt.operands)
			if err == nil {
				t.Fatal("Calculate() error = nil, want error")
			}

			calcErr, ok := AsError(err)
			if !ok {
				t.Fatalf("Calculate() error type = %T, want CalcError", err)
			}

			if calcErr.Code != tt.wantCode {
				t.Fatalf("Calculate() error code = %s, want %s", calcErr.Code, tt.wantCode)
			}
		})
	}
}
