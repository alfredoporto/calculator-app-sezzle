import { useState } from 'react';
import { type CalculationResponse } from './api/calculatorClient';
import { CalculatorForm } from './components/CalculatorForm';
import { ResultCard } from './components/ResultCard';
import './styles.css';

export default function App() {
  const [result, setResult] = useState<number | null>(null);
  const [error, setError] = useState<string | null>(null);

  function handleResult(response: CalculationResponse) {
    setResult(response.result);
    setError(null);
  }

  function handleError(message: string) {
    setError(message || null);
    if (message) {
      setResult(null);
    }
  }

  return (
    <main className="app-shell">
      <section className="calculator-layout" aria-labelledby="app-title">
        <div className="intro">
          <p className="eyebrow">Sezzle assessment</p>
          <h1 id="app-title">Calculator</h1>
          <p>
            Run calculations through the Go API with validation on both sides of
            the request.
          </p>
        </div>

        <div className="workspace">
          <CalculatorForm onResult={handleResult} onError={handleError} />
          <ResultCard result={result} error={error} />
        </div>
      </section>
    </main>
  );
}

