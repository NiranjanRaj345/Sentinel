package sqlite

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/events"
)

func TestOpenEvents_CreatesDatabase(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	store, err := OpenEvents(dbPath)
	if err != nil {
		t.Fatalf("OpenEvents() failed: %v", err)
	}
	defer store.Close()

	if store == nil {
		t.Fatal("expected non-nil Store")
	}

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatal("expected database file to exist")
	}
}

func TestOpenEvents_EmptyPath_ReturnsError(t *testing.T) {
	store, err := OpenEvents("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
	if store != nil {
		t.Fatal("expected nil store on error")
	}
}

func TestStore_Close_NilStore_Succeeds(t *testing.T) {
	var s *Store
	if err := s.Close(); err != nil {
		t.Fatalf("Close() on nil store should succeed, got %v", err)
	}
}

func TestStore_Save_InsertsRow(t *testing.T) {
	store := mustOpenStore(t, "save_test.db")
	defer store.Close()

	event := events.Event{
		Type:     events.EventTypeSystem,
		Severity: events.SeverityInfo,
		Source:   events.SourceScheduler,
		Title:    "test",
		Message:  "message",
		Metadata: map[string]interface{}{"key": "value"},
	}

	if err := store.Save(context.Background(), event); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	var count int
	if err := store.db.QueryRow(`SELECT COUNT(*) FROM events`).Scan(&count); err != nil {
		t.Fatalf("count query failed: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 row, got %d", count)
	}
}

func TestStore_Recent_ReturnsMostRecentFirst(t *testing.T) {
	store := mustOpenStore(t, "recent_test.db")
	defer store.Close()

	eventsList := []events.Event{
		{
			ID:        "a",
			Type:      events.EventTypeSystem,
			Severity:  events.SeverityInfo,
			Source:    events.SourceScheduler,
			Title:     "first",
			Message:   "first message",
			CreatedAt: time.Date(2026, 7, 29, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        "b",
			Type:      events.EventTypeSystem,
			Severity:  events.SeverityWarning,
			Source:    events.SourceScheduler,
			Title:     "second",
			Message:   "second message",
			CreatedAt: time.Date(2026, 7, 29, 11, 0, 0, 0, time.UTC),
		},
	}

	ctx := context.Background()
	for _, e := range eventsList {
		if err := store.Save(ctx, e); err != nil {
			t.Fatalf("Save() failed: %v", err)
		}
	}

	results, err := store.Recent(10)
	if err != nil {
		t.Fatalf("Recent() failed: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 events, got %d", len(results))
	}
	if results[0].Title != "second" {
		t.Fatalf("expected most recent first, got %s", results[0].Title)
	}
}

func TestStore_Recent_Limit(t *testing.T) {
	store := mustOpenStore(t, "limit_test.db")
	defer store.Close()

	ctx := context.Background()
	for i := 0; i < 5; i++ {
		event := events.Event{
			ID:        string(rune('a' + i)),
			Type:      events.EventTypeSystem,
			Severity:  events.SeverityInfo,
			Source:    events.SourceScheduler,
			Title:     "event",
			Message:   "message",
			CreatedAt: time.Date(2026, 7, 29, 10, 0, 0, 0, time.UTC).Add(time.Duration(i) * time.Minute),
		}
		if err := store.Save(ctx, event); err != nil {
			t.Fatalf("Save() failed: %v", err)
		}
	}

	results, err := store.Recent(2)
	if err != nil {
		t.Fatalf("Recent() failed: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 events with limit 2, got %d", len(results))
	}
}

func TestStore_Recent_NilStore_ReturnsError(t *testing.T) {
	var s *Store
	_, err := s.Recent(10)
	if err == nil {
		t.Fatal("expected error when store is nil")
	}
}

func TestStore_Save_NilStore_ReturnsError(t *testing.T) {
	var s *Store
	err := s.Save(context.Background(), events.Event{})
	if err == nil {
		t.Fatal("expected error when store is nil")
	}
}

func mustOpenStore(t *testing.T, name string) *Store {
	t.Helper()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, name)

	store, err := OpenEvents(dbPath)
	if err != nil {
		t.Fatalf("OpenEvents() failed: %v", err)
	}
	return store
}
