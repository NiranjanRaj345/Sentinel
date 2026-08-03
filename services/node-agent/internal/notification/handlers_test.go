package notification

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_Get_ReturnsRecent(t *testing.T) {
	repo := &stubRepo{recent: []Notification{{ID: "n1", Title: "T"}}}
	svc := NewService(repo, nil, nil)
	h := NewHandler(svc, nil)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/notifications/recent", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp RecentResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Notifications) != 1 || resp.Notifications[0].ID != "n1" {
		t.Fatalf("expected notification n1, got %v", resp.Notifications)
	}
}

func TestHandler_Get_EmptyList(t *testing.T) {
	repo := &stubRepo{recent: []Notification{}}
	svc := NewService(repo, nil, nil)
	h := NewHandler(svc, nil)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/notifications/recent", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp RecentResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Notifications) != 0 {
		t.Fatalf("expected 0 notifications, got %d", len(resp.Notifications))
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	repo := &stubRepo{}
	svc := NewService(repo, nil, nil)
	h := NewHandler(svc, nil)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/notifications/recent", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestTestHandler_Post_ReturnsQueued(t *testing.T) {
	repo := &stubRepo{}
	svc := NewService(repo, nil, nil)
	h := NewTestHandler(svc, nil)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/notifications/test", strings.NewReader(`{"provider":"telegram"}`))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp["status"] != "queued" {
		t.Fatalf("expected status queued, got %s", resp["status"])
	}
}

func TestTestHandler_MethodNotAllowed(t *testing.T) {
	repo := &stubRepo{}
	svc := NewService(repo, nil, nil)
	h := NewTestHandler(svc, nil)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/notifications/test", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d: %s", rr.Code, rr.Body.String())
	}
}
