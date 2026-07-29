package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/metrics"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/storage/sqlite"
)

func openTempStore(t *testing.T) *sqlite.Store {
	t.Helper()
	tmpDir := t.TempDir()
	store, err := sqlite.Open(tmpDir + "/test.db")
	if err != nil {
		t.Fatalf("open temp store: %v", err)
	}
	return store
}

func mustSave(t *testing.T, store *sqlite.Store, info metrics.Info) {
	t.Helper()
	if err := store.Save(info); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}
}

func TestHistoryLatest_ReturnsSnapshot(t *testing.T) {
	store := openTempStore(t)
	defer store.Close()

	mustSave(t, store, metrics.Info{Metadata: metrics.Metadata{Timestamp: time.Now().UTC()}})

	h := HistoryLatest(store)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history/latest", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	var info metrics.Info
	if err := json.NewDecoder(rr.Body).Decode(&info); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if info.Metadata.Timestamp.IsZero() {
		t.Fatal("expected non-zero timestamp")
	}
}

func TestHistoryLatest_NoData_Returns404(t *testing.T) {
	store := openTempStore(t)
	defer store.Close()

	h := HistoryLatest(store)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history/latest", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d: %s", http.StatusNotFound, rr.Code, rr.Body.String())
	}
}

func TestHistoryRange_ReturnsSnapshots(t *testing.T) {
	store := openTempStore(t)
	defer store.Close()

	base := time.Date(2026, 7, 29, 15, 0, 0, 0, time.UTC)
	mustSave(t, store, metrics.Info{Metadata: metrics.Metadata{Timestamp: base}})
	mustSave(t, store, metrics.Info{Metadata: metrics.Metadata{Timestamp: base.Add(30 * time.Minute)}})

	h := HistoryRange(store)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history?from=2026-07-29T15:00:00Z&to=2026-07-29T16:00:00Z", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	var results []metrics.Info
	if err := json.NewDecoder(rr.Body).Decode(&results); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 snapshots, got %d", len(results))
	}
}

func TestHistoryRange_InvalidFrom_Returns400(t *testing.T) {
	store := openTempStore(t)
	defer store.Close()

	h := HistoryRange(store)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history?from=invalid&to=2026-07-29T16:00:00Z", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestHistoryRange_FromAfterTo_Returns400(t *testing.T) {
	store := openTempStore(t)
	defer store.Close()

	h := HistoryRange(store)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history?from=2026-07-29T17:00:00Z&to=2026-07-29T16:00:00Z", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestHistoryRange_MissingParams_Returns400(t *testing.T) {
	store := openTempStore(t)
	defer store.Close()

	h := HistoryRange(store)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/history?from=2026-07-29T15:00:00Z", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}
