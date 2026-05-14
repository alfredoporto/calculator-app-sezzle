# Coverage Report

Generated locally with the commands below.

## Backend

Command:

```sh
cd backend
go test ./... -cover
```

Result:

```text
ok  	github.com/alfredos/sezzle-calculator/backend/cmd/server        coverage: 31.8% of statements
ok  	github.com/alfredos/sezzle-calculator/backend/internal/api       coverage: 90.0% of statements
ok  	github.com/alfredos/sezzle-calculator/backend/internal/calculator coverage: 88.1% of statements
```

## Frontend

Command:

```sh
cd frontend
npm test -- --coverage
```

Result:

```text
Test Files  2 passed (2)
Tests       7 passed (7)

Statements 90.66% (68/75)
Branches   77.77% (35/45)
Functions  100%   (18/18)
Lines      90.54% (67/74)
```

## Notes

Generated coverage artifacts are excluded from Git to keep the repository small
and avoid committing machine-generated HTML output. Re-run the commands above to
produce fresh local coverage artifacts.
