package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alfredos/sezzle-calculator/backend/internal/calculator"
)

func TestHandlerHealth(t *testing.T) {
	t.Parallel()

	handler := NewHandler(calculator.NewService())
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}

	var body map[string]string
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body["status"] != "ok" {
		t.Fatalf("status body = %q, want ok", body["status"])
	}
}

func TestHandlerCalculateSuccess(t *testing.T) {
	t.Parallel()

	handler := NewHandler(calculator.NewService())
	body := bytes.NewBufferString(`{"operation":"divide","operands":[10,2]}`)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/calculations", body)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
	}

	var got CalculationResponse
	if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if got.Operation != "divide" {
		t.Fatalf("operation = %q, want divide", got.Operation)
	}
	if got.Result != 5 {
		t.Fatalf("result = %v, want 5", got.Result)
	}
	if len(got.Operands) != 2 || got.Operands[0] != 10 || got.Operands[1] != 2 {
		t.Fatalf("operands = %v, want [10 2]", got.Operands)
	}
}

func TestHandlerCalculateValidationErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantCode   string
	}{
		{
			name:       "malformed json",
			body:       `{"operation":`,
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_JSON",
		},
		{
			name:       "trailing data",
			body:       `{"operation":"add","operands":[1,2]} garbage`,
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_JSON",
		},
		{
			name:       "invalid operation",
			body:       `{"operation":"mod","operands":[10,2]}`,
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_OPERATION",
		},
		{
			name:       "wrong operand count",
			body:       `{"operation":"add","operands":[10]}`,
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_OPERANDS",
		},
		{
			name:       "division by zero",
			body:       `{"operation":"divide","operands":[10,0]}`,
			wantStatus: http.StatusBadRequest,
			wantCode:   "DIVISION_BY_ZERO",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := NewHandler(calculator.NewService())
			request := httptest.NewRequest(http.MethodPost, "/api/v1/calculations", bytes.NewBufferString(tt.body))
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			if response.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body = %s", response.Code, tt.wantStatus, response.Body.String())
			}

			var got ErrorResponse
			if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
				t.Fatalf("decode response: %v", err)
			}

			if got.Error.Code != tt.wantCode {
				t.Fatalf("error code = %q, want %q", got.Error.Code, tt.wantCode)
			}
			if got.Error.Message == "" {
				t.Fatal("error message is empty")
			}
		})
	}
}

func TestHandlerMethodAndRouteHandling(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantAllow  string
	}{
		{name: "calculation get not allowed", method: http.MethodGet, path: "/api/v1/calculations", wantStatus: http.StatusMethodNotAllowed, wantAllow: http.MethodPost},
		{name: "health post not allowed", method: http.MethodPost, path: "/healthz", wantStatus: http.StatusMethodNotAllowed, wantAllow: http.MethodGet},
		{name: "unknown route", method: http.MethodGet, path: "/missing", wantStatus: http.StatusNotFound},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := NewHandler(calculator.NewService())
			request := httptest.NewRequest(tt.method, tt.path, nil)
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			if response.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, tt.wantStatus)
			}
			if got := response.Header().Get("Allow"); got != tt.wantAllow {
				t.Fatalf("Allow header = %q, want %q", got, tt.wantAllow)
			}
		})
	}
}
