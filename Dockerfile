FROM node:24-alpine AS frontend-build
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM golang:1.25-alpine AS backend-build
WORKDIR /app/backend
COPY backend/go.mod ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

FROM alpine:3.22
WORKDIR /app
RUN addgroup -S app && adduser -S app -G app
COPY --from=backend-build /out/server /app/server
COPY --from=frontend-build /app/frontend/dist /app/public
ENV HTTP_ADDR=:8080
ENV FRONTEND_DIST=/app/public
EXPOSE 8080
USER app
ENTRYPOINT ["/app/server"]

