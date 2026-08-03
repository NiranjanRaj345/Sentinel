package scheduler

import (
	"io"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/storage/sqlite"
)

func openTempStore(t *testing.T) *sqlite.Store {
	t.Helper()
	store, err := sqlite.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open temp store: %v", err)
	}
	return store
}

func TestScheduler_Snapshot_BeforeStart_ReturnsError(t *testing.T) {
	s := New(5*time.Second, logger.New(logger.Info, io.Discard), nil, nil, nil, nil, nil)

	_, err := s.Snapshot()
	if err == nil {
		t.Fatal("expected error before Start()")
	}
}

func TestScheduler_Start_InitialCollectionSucceeds(t *testing.T) {
	store := openTempStore(t)
	defer store.Close()

	s := New(100*time.Millisecond, logger.New(logger.Info, io.Discard), store, nil, nil, nil, nil)

	err := s.Start()
	if err != nil {
		t.Fatalf("expected Start() to succeed, got %v", err)
	}
	defer s.Stop()

	info, err := s.Snapshot()
	if err != nil {
		t.Fatalf("expected Snapshot() to succeed after Start(), got %v", err)
	}
	if info.Metadata.Timestamp.IsZero() {
		t.Fatal("expected non-zero timestamp in snapshot")
	}
}

func TestScheduler_Snapshot_UpdatesOverTime(t *testing.T) {
	store := openTempStore(t)
	defer store.Close()

	s := New(50*time.Millisecond, logger.New(logger.Info, io.Discard), store, nil, nil, nil, nil)

	if err := s.Start(); err != nil {
		t.Fatalf("Start() failed: %v", err)
	}
	defer s.Stop()

	info1, err := s.Snapshot()
	if err != nil {
		t.Fatalf("first Snapshot() failed: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	info2, err := s.Snapshot()
	if err != nil {
		t.Fatalf("second Snapshot() failed: %v", err)
	}

	if !info2.Metadata.Timestamp.After(info1.Metadata.Timestamp) {
		t.Fatalf("expected snapshot timestamp to update: %v vs %v", info2.Metadata.Timestamp, info1.Metadata.Timestamp)
	}
}

func TestScheduler_Snapshot_Concurrent(t *testing.T) {
	store := openTempStore(t)
	defer store.Close()

	s := New(100*time.Millisecond, logger.New(logger.Info, io.Discard), store, nil, nil, nil, nil)

	if err := s.Start(); err != nil {
		t.Fatalf("Start() failed: %v", err)
	}
	defer s.Stop()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				info, err := s.Snapshot()
				if err != nil {
					t.Errorf("Snapshot() failed: %v", err)
					return
				}
				if info.Metadata.Timestamp.IsZero() {
					t.Error("expected non-zero timestamp")
					return
				}
			}
		}()
	}

	wg.Wait()
}

func TestScheduler_Stop_TerminatesGoroutine(t *testing.T) {
	store := openTempStore(t)
	defer store.Close()

	s := New(100*time.Millisecond, logger.New(logger.Info, io.Discard), store, nil, nil, nil, nil)

	if err := s.Start(); err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	s.Stop()

	_, err := s.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot() failed after Stop(): %v", err)
	}
}

func TestScheduler_Stats_PopulatedAfterStart(t *testing.T) {
	store := openTempStore(t)
	defer store.Close()

	s := New(100*time.Millisecond, logger.New(logger.Info, io.Discard), store, nil, nil, nil, nil)

	if err := s.Start(); err != nil {
		t.Fatalf("Start() failed: %v", err)
	}
	defer s.Stop()

	stats := s.Stats()
	if stats.StartedAt.IsZero() {
		t.Fatal("expected StartedAt to be set after Start()")
	}
	if stats.LastCollectionAt.IsZero() {
		t.Fatal("expected LastCollectionAt to be set after initial collection")
	}
	if stats.SuccessfulCollections != 1 {
		t.Fatalf("expected 1 successful collection, got %d", stats.SuccessfulCollections)
	}
	if stats.FailedCollections != 0 {
		t.Fatalf("expected 0 failed collections, got %d", stats.FailedCollections)
	}
	if stats.LastError != "" {
		t.Fatalf("expected LastError to be empty, got %q", stats.LastError)
	}
}

func TestScheduler_Stats_SuccessfulCollection_IncrementsCounter(t *testing.T) {
	store := openTempStore(t)
	defer store.Close()

	s := New(50*time.Millisecond, logger.New(logger.Info, io.Discard), store, nil, nil, nil, nil)

	if err := s.Start(); err != nil {
		t.Fatalf("Start() failed: %v", err)
	}
	defer s.Stop()

	time.Sleep(150 * time.Millisecond)

	stats := s.Stats()
	if stats.SuccessfulCollections < 2 {
		t.Fatalf("expected at least 2 successful collections, got %d", stats.SuccessfulCollections)
	}
	if stats.LastCollectionDuration <= 0 {
		t.Fatalf("expected LastCollectionDuration > 0, got %v", stats.LastCollectionDuration)
	}
}

func TestScheduler_Stats_ConcurrentAccess(t *testing.T) {
	store := openTempStore(t)
	defer store.Close()

	s := New(100*time.Millisecond, logger.New(logger.Info, io.Discard), store, nil, nil, nil, nil)

	if err := s.Start(); err != nil {
		t.Fatalf("Start() failed: %v", err)
	}
	defer s.Stop()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				stats := s.Stats()
				if stats.StartedAt.IsZero() {
					t.Error("expected StartedAt to be set")
					return
				}
			}
		}()
	}

	wg.Wait()
}
