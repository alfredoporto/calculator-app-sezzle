.PHONY: test build backend-test frontend-test frontend-build docker-build run-backend run-frontend

test: backend-test frontend-test

build: frontend-build

backend-test:
	cd backend && go test ./... -cover

frontend-test:
	cd frontend && npm test -- --coverage

frontend-build:
	cd frontend && npm run build

docker-build:
	docker build -t sezzle-calculator .

run-backend:
	cd backend && go run ./cmd/server

run-frontend:
	cd frontend && npm run dev

