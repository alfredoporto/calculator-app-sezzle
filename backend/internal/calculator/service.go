package calculator

import (
	"errors"
	"fmt"
	"math"
)

type ErrorCode string

const (
	ErrorCodeInvalidOperation ErrorCode = "INVALID_OPERATION"
	ErrorCodeInvalidOperands  ErrorCode = "INVALID_OPERANDS"
	ErrorCodeDivisionByZero   ErrorCode = "DIVISION_BY_ZERO"
	ErrorCodeNonFiniteResult  ErrorCode = "NON_FINITE_RESULT"
)

type CalcError struct {
	Code    ErrorCode
	Message string
}

func (e CalcError) Error() string {
	return e.Message
}

func AsError(err error) (CalcError, bool) {
	var calcErr CalcError
	if errors.As(err, &calcErr) {
		return calcErr, true
	}

	return CalcError{}, false
}

type Service struct{}

func NewService() Service {
	return Service{}
}

func (s Service) Calculate(operation Operation, operands []float64) (float64, error) {
	if !operation.IsSupported() {
		return 0, CalcError{
			Code:    ErrorCodeInvalidOperation,
			Message: fmt.Sprintf("unsupported operation %q", operation),
		}
	}

	if err := validateOperandCount(operation, operands); err != nil {
		return 0, err
	}

	result, err := calculate(operation, operands)
	if err != nil {
		return 0, err
	}

	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, CalcError{
			Code:    ErrorCodeNonFiniteResult,
			Message: "calculation result is not finite",
		}
	}

	return result, nil
}

func validateOperandCount(operation Operation, operands []float64) error {
	want, ok := operation.ExpectedOperandCount()
	if !ok {
		return CalcError{
			Code:    ErrorCodeInvalidOperation,
			Message: fmt.Sprintf("unsupported operation %q", operation),
		}
	}

	if len(operands) != want {
		return CalcError{
			Code:    ErrorCodeInvalidOperands,
			Message: fmt.Sprintf("%s expects %d operand(s)", operation, want),
		}
	}

	return nil
}

func calculate(operation Operation, operands []float64) (float64, error) {
	switch operation {
	case OperationAdd:
		return operands[0] + operands[1], nil
	case OperationSubtract:
		return operands[0] - operands[1], nil
	case OperationMultiply:
		return operands[0] * operands[1], nil
	case OperationDivide:
		if operands[1] == 0 {
			return 0, CalcError{
				Code:    ErrorCodeDivisionByZero,
				Message: "division by zero",
			}
		}
		return operands[0] / operands[1], nil
	case OperationPower:
		return math.Pow(operands[0], operands[1]), nil
	case OperationSquareRoot:
		if operands[0] < 0 {
			return 0, CalcError{
				Code:    ErrorCodeInvalidOperands,
				Message: "square root requires a non-negative operand",
			}
		}
		return math.Sqrt(operands[0]), nil
	case OperationPercentage:
		return operands[0] / 100, nil
	default:
		return 0, CalcError{
			Code:    ErrorCodeInvalidOperation,
			Message: fmt.Sprintf("unsupported operation %q", operation),
		}
	}
}
