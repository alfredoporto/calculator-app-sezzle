package api

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/alfredos/sezzle-calculator/backend/internal/calculator"
)

const maxRequestBytes = 1 << 20

type Handler struct {
	calculator calculator.Service
	logger     *slog.Logger
}

func NewHandler(calculator calculator.Service) Handler {
	return Handler{
		calculator: calculator,
		logger:     slog.Default(),
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/healthz":
		h.handleHealth(w, r)
	case "/api/v1/calculations":
		h.handleCalculation(w, r)
	default:
		writeError(w, http.StatusNotFound, "NOT_FOUND", "route not found")
	}
}

func (h Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h Handler) handleCalculation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBytes)

	var request CalculationRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "request body must be valid calculation JSON")
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "request body must contain exactly one JSON object")
		return
	}

	result, err := h.calculator.Calculate(calculator.Operation(request.Operation), request.Operands)
	if err != nil {
		h.writeCalculationError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, CalculationResponse{
		Operation: request.Operation,
		Operands:  request.Operands,
		Result:    result,
	})
}

func (h Handler) writeCalculationError(w http.ResponseWriter, err error) {
	calcErr, ok := calculator.AsError(err)
	if !ok {
		h.logger.Error("calculate", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	writeError(w, http.StatusBadRequest, string(calcErr.Code), calcErr.Message)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("write response", slog.String("error", err.Error()))
	}
}

func writeError(w http.ResponseWriter, status int, code string, message string) {
	writeJSON(w, status, ErrorResponse{
		Error: ErrorBody{
			Code:    code,
			Message: message,
		},
	})
}

var _ http.Handler = Handler{}
