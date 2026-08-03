package operations

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_Post_InvalidJSON_Returns400(t *testing.T) {
	svc := NewService(nil, NewAuditor(nil), NewValidator(nil), nil, nil, nil)
	h := NewHandler(svc, nil)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/operations", strings.NewReader("not-json"))
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestHandler_Post_ValidationFailure_Returns400(t *testing.T) {
	svc := NewService(nil, NewAuditor(nil), NewValidator(nil), nil, nil, nil)
	h := NewHandler(svc, nil)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/operations", strings.NewReader(`{"action":"restart","confirm":false}`))
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}

	var body map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if body["error"] == "" {
		t.Fatal("expected error message in response")
	}
}

func TestHandler_Post_Success_ReturnsResult(t *testing.T) {
	provider := &stubProvider{supported: true}
	svc := NewService(provider, NewAuditor(nil), NewValidator(provider), nil, nil, nil)
	h := NewHandler(svc, nil)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/operations", strings.NewReader(`{"action":"restart","confirm":true}`))
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	var result Result
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	if !result.Success {
		t.Fatal("expected success")
	}
}

func TestHandler_MethodNotAllowed_Returns405(t *testing.T) {
	svc := NewService(nil, NewAuditor(nil), NewValidator(nil), nil, nil, nil)
	h := NewHandler(svc, nil)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/operations", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected %d, got %d: %s", http.StatusMethodNotAllowed, rr.Code, rr.Body.String())
	}
}
