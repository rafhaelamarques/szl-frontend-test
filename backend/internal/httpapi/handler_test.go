package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealth(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	NewMux().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}

func TestCalculate(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantResult *float64
		wantError  string
	}{
		{
			name:       "add",
			body:       `{"operation":"add","a":2,"b":3}`,
			wantStatus: http.StatusOK,
			wantResult: ptr(5),
		},
		{
			name:       "divide by zero",
			body:       `{"operation":"divide","a":1,"b":0}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantError:  "division by zero",
		},
		{
			name:       "sqrt negative",
			body:       `{"operation":"sqrt","a":-4}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantError:  "square root of negative number",
		},
		{
			name:       "missing a",
			body:       `{"operation":"add","b":1}`,
			wantStatus: http.StatusBadRequest,
			wantError:  "field a is required",
		},
		{
			name:       "missing b",
			body:       `{"operation":"add","a":1}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantError:  "missing operand",
		},
		{
			name:       "unknown op",
			body:       `{"operation":"mod","a":1,"b":2}`,
			wantStatus: http.StatusBadRequest,
			wantError:  "unknown operation",
		},
		{
			name:       "invalid json",
			body:       `{`,
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid JSON body",
		},
		{
			name:       "sqrt unary",
			body:       `{"operation":"sqrt","a":16}`,
			wantStatus: http.StatusOK,
			wantResult: ptr(4),
		},
	}

	mux := NewMux()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			mux.ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d, body=%s", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if tt.wantResult != nil {
				var res struct {
					Result float64 `json:"result"`
				}
				if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
					t.Fatalf("decode: %v", err)
				}
				if res.Result != *tt.wantResult {
					t.Fatalf("result = %v, want %v", res.Result, *tt.wantResult)
				}
				return
			}
			var res struct {
				Error string `json:"error"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if res.Error != tt.wantError {
				t.Fatalf("error = %q, want %q", res.Error, tt.wantError)
			}
		})
	}
}

func TestCORSPreflight(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/api/v1/calculate", nil)
	NewMux().ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("missing CORS origin")
	}
}

func ptr(v float64) *float64 { return &v }
