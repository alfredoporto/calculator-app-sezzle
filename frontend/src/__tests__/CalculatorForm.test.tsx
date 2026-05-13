import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import App from '../App';
import { calculate } from '../api/calculatorClient';

vi.mock('../api/calculatorClient', async () => {
  const actual = await vi.importActual<typeof import('../api/calculatorClient')>(
    '../api/calculatorClient'
  );

  return {
    ...actual,
    calculate: vi.fn()
  };
});

const calculateMock = vi.mocked(calculate);

describe('CalculatorForm', () => {
  beforeEach(() => {
    calculateMock.mockReset();
  });

  it('validates inputs before submitting', async () => {
    const user = userEvent.setup();
    render(<App />);

    await user.clear(screen.getByLabelText(/first operand/i));
    await user.click(screen.getByRole('button', { name: /calculate/i }));

    expect(
      screen.getAllByText(/enter a valid number for the first operand/i)
    ).toHaveLength(2);
    expect(calculateMock).not.toHaveBeenCalled();
  });

  it('clears a previous result after client-side validation fails', async () => {
    const user = userEvent.setup();
    calculateMock.mockResolvedValue({
      operation: 'add',
      operands: [10, 2],
      result: 12
    });

    render(<App />);

    await user.clear(screen.getByLabelText(/first operand/i));
    await user.type(screen.getByLabelText(/first operand/i), '10');
    await user.clear(screen.getByLabelText(/second operand/i));
    await user.type(screen.getByLabelText(/second operand/i), '2');
    await user.click(screen.getByRole('button', { name: /calculate/i }));

    expect(await screen.findByText('12')).toBeInTheDocument();

    await user.clear(screen.getByLabelText(/first operand/i));
    await user.click(screen.getByRole('button', { name: /calculate/i }));

    expect(screen.queryByText('12')).not.toBeInTheDocument();
    expect(
      screen.getAllByText(/enter a valid number for the first operand/i)
    ).toHaveLength(2);
  });

  it('submits a calculation and displays the result', async () => {
    const user = userEvent.setup();
    calculateMock.mockResolvedValue({
      operation: 'add',
      operands: [10, 2],
      result: 12
    });

    render(<App />);

    await user.clear(screen.getByLabelText(/first operand/i));
    await user.type(screen.getByLabelText(/first operand/i), '10');
    await user.clear(screen.getByLabelText(/second operand/i));
    await user.type(screen.getByLabelText(/second operand/i), '2');
    await user.click(screen.getByRole('button', { name: /calculate/i }));

    expect(calculateMock).toHaveBeenCalledWith({
      operation: 'add',
      operands: [10, 2]
    });
    expect(await screen.findByText('12')).toBeInTheDocument();
  });

  it('displays backend errors', async () => {
    const user = userEvent.setup();
    calculateMock.mockRejectedValue(new Error('division by zero'));

    render(<App />);

    await user.selectOptions(screen.getByLabelText(/operation/i), 'divide');
    await user.clear(screen.getByLabelText(/first operand/i));
    await user.type(screen.getByLabelText(/first operand/i), '10');
    await user.clear(screen.getByLabelText(/second operand/i));
    await user.type(screen.getByLabelText(/second operand/i), '0');
    await user.click(screen.getByRole('button', { name: /calculate/i }));

    await waitFor(() => {
      expect(screen.getByText(/division by zero/i)).toBeInTheDocument();
    });
  });
});
