package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/alfredos/sezzle-calculator/backend/internal/api"
	"github.com/alfredos/sezzle-calculator/backend/internal/calculator"
)

func TestWithFrontendDelegatesAPIRoutes(t *testing.T) {
	t.Parallel()

	handler := testFrontendHandler()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/calculations", strings.NewReader(`{"operation":"multiply","operands":[6,7]}`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
	}
	if got := response.Body.String(); !strings.Contains(got, `"result":42`) {
		t.Fatalf("body = %s, want result 42", got)
	}
}

func TestWithFrontendDelegatesHealthRoute(t *testing.T) {
	t.Parallel()

	handler := testFrontendHandler()
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if got := response.Body.String(); !strings.Contains(got, `"status":"ok"`) {
		t.Fatalf("body = %s, want health status", got)
	}
}

func TestWithFrontendServesStaticAssets(t *testing.T) {
	t.Parallel()

	handler := testFrontendHandler()
	request := httptest.NewRequest(http.MethodGet, "/assets/app.js", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if got := response.Body.String(); got != "console.log('calculator');" {
		t.Fatalf("body = %q, want static asset content", got)
	}
}

func TestWithFrontendFallsBackToSPAIndex(t *testing.T) {
	t.Parallel()

	handler := testFrontendHandler()
	request := httptest.NewRequest(http.MethodGet, "/calculator/history", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if got := response.Body.String(); !strings.Contains(got, "<div id=\"root\"></div>") {
		t.Fatalf("body = %s, want SPA index", got)
	}
}

func testFrontendHandler() http.Handler {
	frontend := fstest.MapFS{
		"index.html": {
			Data: []byte(`<!doctype html><html><body><div id="root"></div></body></html>`),
		},
		"assets/app.js": {
			Data: []byte("console.log('calculator');"),
		},
	}

	return withFrontend(api.NewHandler(calculator.NewService()), frontend)
}
