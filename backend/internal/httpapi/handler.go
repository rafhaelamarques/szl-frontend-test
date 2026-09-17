package httpapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"calculator/internal/calc"
)

type calculateRequest struct {
	Operation calc.Operation `json:"operation"`
	A         *float64       `json:"a"`
	B         *float64       `json:"b"`
}

type calculateResponse struct {
	Result float64 `json:"result"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("POST /api/v1/calculate", calculate)
	return withCORS(withRecover(mux))
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func calculate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req calculateRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.A == nil {
		writeError(w, http.StatusBadRequest, "field a is required")
		return
	}
	if calc.RequiresB(req.Operation) && req.B == nil {
		writeError(w, http.StatusUnprocessableEntity, calc.ErrMissingOperand.Error())
		return
	}

	b := 0.0
	if req.B != nil {
		b = *req.B
	}

	result, err := calc.Compute(req.Operation, *req.A, b)
	if err != nil {
		writeError(w, statusFor(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, calculateResponse{Result: result})
}

func statusFor(err error) int {
	switch {
	case errors.Is(err, calc.ErrUnknownOperation):
		return http.StatusBadRequest
	case errors.Is(err, calc.ErrDivisionByZero),
		errors.Is(err, calc.ErrNegativeSqrt),
		errors.Is(err, calc.ErrMissingOperand),
		errors.Is(err, calc.ErrNonFiniteResult):
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func withRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic: %v", rec)
				writeError(w, http.StatusInternalServerError, "internal error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
