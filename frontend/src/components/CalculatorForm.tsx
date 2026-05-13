import { FormEvent, useMemo, useState } from 'react';
import {
  calculate,
  CalculatorApiError,
  type CalculationResponse,
  type Operation
} from '../api/calculatorClient';

type OperationOption = {
  value: Operation;
  label: string;
  operandCount: 1 | 2;
};

const operations: OperationOption[] = [
  { value: 'add', label: 'Addition', operandCount: 2 },
  { value: 'subtract', label: 'Subtraction', operandCount: 2 },
  { value: 'multiply', label: 'Multiplication', operandCount: 2 },
  { value: 'divide', label: 'Division', operandCount: 2 },
  { value: 'power', label: 'Exponentiation', operandCount: 2 },
  { value: 'sqrt', label: 'Square Root', operandCount: 1 },
  { value: 'percentage', label: 'Percentage', operandCount: 1 }
];

type CalculatorFormProps = {
  onResult: (response: CalculationResponse) => void;
  onError: (message: string) => void;
};

export function CalculatorForm({ onResult, onError }: CalculatorFormProps) {
  const [operation, setOperation] = useState<Operation>('add');
  const [firstOperand, setFirstOperand] = useState('0');
  const [secondOperand, setSecondOperand] = useState('0');
  const [validationError, setValidationError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const selectedOperation = useMemo(
    () => operations.find((item) => item.value === operation) ?? operations[0],
    [operation]
  );
  const needsSecondOperand = selectedOperation.operandCount === 2;

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setValidationError(null);
    onError('');

    const parsedFirst = parseOperand(firstOperand);
    if (parsedFirst === null) {
      const message = 'Enter a valid number for the first operand.';
      setValidationError(message);
      onError(message);
      return;
    }

    const operands = [parsedFirst];
    if (needsSecondOperand) {
      const parsedSecond = parseOperand(secondOperand);
      if (parsedSecond === null) {
        const message = 'Enter a valid number for the second operand.';
        setValidationError(message);
        onError(message);
        return;
      }
      operands.push(parsedSecond);
    }

    setIsSubmitting(true);
    try {
      const response = await calculate({ operation, operands });
      onResult(response);
    } catch (error) {
      onError(toUserMessage(error));
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <form className="calculator-form" onSubmit={handleSubmit}>
      <label htmlFor="operation">Operation</label>
      <select
        id="operation"
        value={operation}
        onChange={(event) => setOperation(event.target.value as Operation)}
      >
        {operations.map((item) => (
          <option key={item.value} value={item.value}>
            {item.label}
          </option>
        ))}
      </select>

      <label htmlFor="first-operand">First operand</label>
      <input
        id="first-operand"
        inputMode="decimal"
        value={firstOperand}
        onChange={(event) => setFirstOperand(event.target.value)}
      />

      {needsSecondOperand ? (
        <>
          <label htmlFor="second-operand">Second operand</label>
          <input
            id="second-operand"
            inputMode="decimal"
            value={secondOperand}
            onChange={(event) => setSecondOperand(event.target.value)}
          />
        </>
      ) : null}

      {validationError ? (
        <p className="error-message" role="alert">
          {validationError}
        </p>
      ) : null}

      <button type="submit" disabled={isSubmitting}>
        {isSubmitting ? 'Calculating...' : 'Calculate'}
      </button>
    </form>
  );
}

function parseOperand(value: string): number | null {
  if (value.trim() === '') {
    return null;
  }

  const parsed = Number(value);
  if (!Number.isFinite(parsed)) {
    return null;
  }

  return parsed;
}

function toUserMessage(error: unknown): string {
  if (error instanceof CalculatorApiError) {
    return error.message;
  }

  if (error instanceof Error) {
    return error.message;
  }

  return 'Calculation failed.';
}
