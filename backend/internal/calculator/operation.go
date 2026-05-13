package calculator

type Operation string

const (
	OperationAdd        Operation = "add"
	OperationSubtract   Operation = "subtract"
	OperationMultiply   Operation = "multiply"
	OperationDivide     Operation = "divide"
	OperationPower      Operation = "power"
	OperationSquareRoot Operation = "sqrt"
	OperationPercentage Operation = "percentage"
)

func (o Operation) IsSupported() bool {
	switch o {
	case OperationAdd,
		OperationSubtract,
		OperationMultiply,
		OperationDivide,
		OperationPower,
		OperationSquareRoot,
		OperationPercentage:
		return true
	default:
		return false
	}
}

func (o Operation) ExpectedOperandCount() (int, bool) {
	switch o {
	case OperationAdd,
		OperationSubtract,
		OperationMultiply,
		OperationDivide,
		OperationPower:
		return 2, true
	case OperationSquareRoot,
		OperationPercentage:
		return 1, true
	default:
		return 0, false
	}
}
