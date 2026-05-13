type ResultCardProps = {
  result: number | null;
  error: string | null;
};

export function ResultCard({ result, error }: ResultCardProps) {
  return (
    <section className="result-panel" aria-live="polite">
      <h2>Result</h2>
      {error ? <p className="error-message">{error}</p> : null}
      {result !== null && !error ? <p className="result-value">{result}</p> : null}
      {result === null && !error ? (
        <p className="empty-result">Enter values and calculate.</p>
      ) : null}
    </section>
  );
}

