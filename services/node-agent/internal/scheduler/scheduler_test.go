package scheduler

import (
	"io"
	"sync"
	"testing"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
)

func TestScheduler_Snapshot_BeforeStart_ReturnsError(t *testing.T) {
	s := New(5*time.Second, logger.New(logger.Info, io.Discard))

	_, err := s.Snapshot()
	if err == nil {
		t.Fatal("expected error before Start()")
	}
}

func TestScheduler_Start_InitialCollectionSucceeds(t *testing.T) {
	s := New(100*time.Millisecond, logger.New(logger.Info, io.Discard))

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
	s := New(50*time.Millisecond, logger.New(logger.Info, io.Discard))

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
	s := New(100*time.Millisecond, logger.New(logger.Info, io.Discard))

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
	s := New(100*time.Millisecond, logger.New(logger.Info, io.Discard))

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
