# Sezzle Calculator

Full-stack calculator assessment built with React, TypeScript, Vite, and a Go REST API.

## Features

- Basic operations: addition, subtraction, multiplication, division.
- Advanced operations: exponentiation, square root, percentage.
- Frontend validation and backend validation.
- JSON API responses for success and errors.
- Backend service and handler tests.
- Frontend validation and happy-path tests.

## Project Structure

```text
backend/   Go REST API
frontend/  React + TypeScript + Vite app
```

## Prerequisites

- Go 1.25+
- Node.js 24+
- npm 11+
- Docker, optional

## Running Locally

Backend:

```sh
cd backend
go run ./cmd/server
```

Frontend:

```sh
cd frontend
npm install
npm run dev
```

The frontend dev server proxies `/api` and `/healthz` requests to the backend.

## Tests

Backend:

```sh
cd backend
go test ./... -cover
```

Frontend:

```sh
cd frontend
npm test -- --coverage
```

This project was implemented with a test-first workflow: calculator service tests,
handler tests, API client tests, and UI behavior tests were written before the
corresponding production code.

Current local coverage at completion:

- Backend calculator package: 88.1%
- Backend API package: 88.6%
- Frontend statements: 92.85%

## API

Health check:

```sh
curl http://localhost:8080/healthz
```

Calculate:

```sh
curl -X POST http://localhost:8080/api/v1/calculations \
  -H 'Content-Type: application/json' \
  -d '{"operation":"divide","operands":[10,2]}'
```

Success response:

```json
{
  "operation": "divide",
  "operands": [10, 2],
  "result": 5
}
```

Error response:

```json
{
  "error": {
    "code": "DIVISION_BY_ZERO",
    "message": "division by zero"
  }
}
```

## Design Decisions

- The backend uses Go's standard `net/http` package to keep dependencies small and behavior explicit.
- Calculator rules live in a service package so they can be tested independently from HTTP.
- HTTP models are separate from calculator logic to keep transport concerns out of business rules.
- The frontend calls the backend through a typed API client instead of embedding `fetch` directly in UI components.
- Numeric behavior uses `float64`, which is appropriate for a calculator but not for money math.
- The project is delivered as a monorepo to make review, local setup, and full-stack Docker execution straightforward.

## Docker

Build:

```sh
docker build -t sezzle-calculator .
```

Run:

```sh
docker run --rm -p 8080:8080 sezzle-calculator
```

Then open `http://localhost:8080`.

## Notes

- The backend has no database or external service dependencies.
- The API intentionally returns stable error codes instead of raw internal errors.
- This application uses `float64`; financial-grade decimal math is out of scope
  for a general calculator assessment.
