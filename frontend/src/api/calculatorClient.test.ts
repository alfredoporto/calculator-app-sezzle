import { describe, expect, it, vi } from 'vitest';
import {
  calculate,
  CalculatorApiError,
  type CalculationRequest
} from './calculatorClient';

describe('calculate', () => {
  it('returns a calculation response', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      new Response(
        JSON.stringify({
          operation: 'add',
          operands: [10, 2],
          result: 12
        }),
        { status: 200, headers: { 'Content-Type': 'application/json' } }
      )
    );

    const request: CalculationRequest = { operation: 'add', operands: [10, 2] };
    const result = await calculate(request, { fetcher: fetchMock });

    expect(fetchMock).toHaveBeenCalledWith('/api/v1/calculations', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(request),
      signal: undefined
    });
    expect(result).toEqual({ operation: 'add', operands: [10, 2], result: 12 });
  });

  it('throws a typed API error for backend error responses', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      new Response(
        JSON.stringify({
          error: {
            code: 'DIVISION_BY_ZERO',
            message: 'division by zero'
          }
        }),
        { status: 400, headers: { 'Content-Type': 'application/json' } }
      )
    );

    await expect(
      calculate({ operation: 'divide', operands: [10, 0] }, { fetcher: fetchMock })
    ).rejects.toMatchObject({
      code: 'DIVISION_BY_ZERO',
      message: 'division by zero',
      status: 400
    });
  });

  it('throws a generic API error for unexpected error payloads', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      new Response('not found', {
        status: 404,
        headers: { 'Content-Type': 'text/plain' }
      })
    );

    await expect(
      calculate({ operation: 'add', operands: [10, 2] }, { fetcher: fetchMock })
    ).rejects.toBeInstanceOf(CalculatorApiError);
  });
});

