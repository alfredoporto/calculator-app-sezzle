export type Operation =
  | 'add'
  | 'subtract'
  | 'multiply'
  | 'divide'
  | 'power'
  | 'sqrt'
  | 'percentage';

export type CalculationRequest = {
  operation: Operation;
  operands: number[];
};

export type CalculationResponse = {
  operation: Operation;
  operands: number[];
  result: number;
};

type ErrorResponse = {
  error?: {
    code?: string;
    message?: string;
  };
};

type CalculateOptions = {
  fetcher?: typeof fetch;
  signal?: AbortSignal;
};

export class CalculatorApiError extends Error {
  readonly code: string;
  readonly status: number;

  constructor(code: string, message: string, status: number) {
    super(message);
    this.name = 'CalculatorApiError';
    this.code = code;
    this.status = status;
  }
}

export async function calculate(
  request: CalculationRequest,
  options: CalculateOptions = {}
): Promise<CalculationResponse> {
  const fetcher = options.fetcher ?? fetch;
  const response = await fetcher('/api/v1/calculations', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(request),
    signal: options.signal
  });

  if (!response.ok) {
    throw await toApiError(response);
  }

  return (await response.json()) as CalculationResponse;
}

async function toApiError(response: Response): Promise<CalculatorApiError> {
  try {
    const payload = (await response.json()) as ErrorResponse;
    const code = payload.error?.code ?? 'REQUEST_FAILED';
    const message = payload.error?.message ?? 'Calculation request failed';
    return new CalculatorApiError(code, message, response.status);
  } catch {
    return new CalculatorApiError(
      'REQUEST_FAILED',
      'Calculation request failed',
      response.status
    );
  }
}

