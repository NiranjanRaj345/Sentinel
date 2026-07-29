package sqlite

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/metrics"
)

func TestOpen_CreatesDatabase(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() failed: %v", err)
	}

	if store == nil {
		t.Fatal("expected non-nil Store")
	}

	if err := store.Close(); err != nil {
		t.Fatalf("Close() failed: %v", err)
	}

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatal("expected database file to exist")
	}
}

func TestOpen_CreatesSchema(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "schema_test.db")

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() failed: %v", err)
	}
	defer store.Close()

	var tableName string
	err = store.db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='metrics_snapshots'`).Scan(&tableName)
	if err != nil {
		t.Fatalf("metrics_snapshots table not found: %v", err)
	}
	if tableName != "metrics_snapshots" {
		t.Fatalf("expected metrics_snapshots, got %s", tableName)
	}

	var version int
	err = store.db.QueryRow(`SELECT COALESCE(MAX(version), 0) FROM schema_version`).Scan(&version)
	if err != nil {
		t.Fatalf("schema_version table not found: %v", err)
	}
	if version != SchemaVersion {
		t.Fatalf("expected schema version %d, got %d", SchemaVersion, version)
	}
}

func TestOpen_TwiceSucceeds(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "double_open.db")

	store1, err := Open(dbPath)
	if err != nil {
		t.Fatalf("first Open() failed: %v", err)
	}
	defer store1.Close()

	store2, err := Open(dbPath)
	if err != nil {
		t.Fatalf("second Open() failed: %v", err)
	}
	defer store2.Close()

	if store1 == nil || store2 == nil {
		t.Fatal("expected non-nil stores")
	}
}

func TestOpen_EmptyPath_ReturnsError(t *testing.T) {
	store, err := Open("")
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

func TestStore_SchemaVersion_Exists(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "version_test.db")

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() failed: %v", err)
	}
	defer store.Close()

	var currentVersion int
	if err := store.db.QueryRow(`SELECT COALESCE(MAX(version), 0) FROM schema_version`).Scan(&currentVersion); err != nil {
		t.Fatalf("failed to read schema_version: %v", err)
	}
	if currentVersion != SchemaVersion {
		t.Fatalf("expected schema version %d, got %d", SchemaVersion, currentVersion)
	}
}

func TestStore_Save_InsertsRow(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "save_test.db")

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() failed: %v", err)
	}
	defer store.Close()

	info := metrics.Info{
		Metadata: metrics.Metadata{
			Timestamp: time.Now().UTC(),
		},
		CPU: metrics.CPUInfo{
			UsagePercent: 42.0,
		},
		Memory: metrics.MemoryInfo{
			TotalBytes:   1024,
			UsedBytes:    512,
			UsagePercent: 50.0,
		},
		Disks: []metrics.DiskInfo{},
		Network: metrics.NetworkInfo{
			Hostname:   "test-host",
			Interfaces: []metrics.NetworkInterface{},
			IO:         metrics.NetworkIO{},
		},
		Processes: []metrics.ProcessInfo{},
	}

	if err := store.Save(info); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	var count int
	if err := store.db.QueryRow(`SELECT COUNT(*) FROM metrics_snapshots`).Scan(&count); err != nil {
		t.Fatalf("count query failed: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 row, got %d", count)
	}
}

func TestStore_Save_NilStore_ReturnsError(t *testing.T) {
	var s *Store
	info := metrics.Info{}
	if err := s.Save(info); err == nil {
		t.Fatal("expected error when saving to nil store")
	}
}

func TestStore_Latest_EmptyDatabase(t *testing.T) {
	store := mustOpenStore(t, "latest_empty.db")
	defer store.Close()

	_, err := store.Latest()
	if err == nil {
		t.Fatal("expected error when database is empty")
	}
}

func TestStore_Latest_ReturnsMostRecent(t *testing.T) {
	store := mustOpenStore(t, "latest_recent.db")
	defer store.Close()

	older := metrics.Info{Metadata: metrics.Metadata{Timestamp: time.Now().UTC().Add(-2 * time.Hour)}}
	newer := metrics.Info{Metadata: metrics.Metadata{Timestamp: time.Now().UTC()}}

	if err := store.Save(older); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}
	if err := store.Save(newer); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	info, err := store.Latest()
	if err != nil {
		t.Fatalf("Latest() failed: %v", err)
	}

	if !info.Metadata.Timestamp.Equal(newer.Metadata.Timestamp) {
		t.Fatalf("expected latest timestamp %v, got %v", newer.Metadata.Timestamp, info.Metadata.Timestamp)
	}
}

func TestStore_Range_ReturnsSnapshotsInOrder(t *testing.T) {
	store := mustOpenStore(t, "range_ordered.db")
	defer store.Close()

	base := time.Date(2026, 7, 29, 15, 0, 0, 0, time.UTC)
	snapshots := []metrics.Info{
		{Metadata: metrics.Metadata{Timestamp: base}},
		{Metadata: metrics.Metadata{Timestamp: base.Add(30 * time.Minute)}},
		{Metadata: metrics.Metadata{Timestamp: base.Add(60 * time.Minute)}},
	}

	for _, snap := range snapshots {
		if err := store.Save(snap); err != nil {
			t.Fatalf("Save() failed: %v", err)
		}
	}

	results, err := store.Range(base, base.Add(60*time.Minute))
	if err != nil {
		t.Fatalf("Range() failed: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 snapshots, got %d", len(results))
	}

	for i, snap := range results {
		if !snap.Metadata.Timestamp.Equal(snapshots[i].Metadata.Timestamp) {
			t.Fatalf("snapshot %d: expected %v, got %v", i, snapshots[i].Metadata.Timestamp, snap.Metadata.Timestamp)
		}
	}
}

func TestStore_Range_EmptyResult(t *testing.T) {
	store := mustOpenStore(t, "range_empty.db")
	defer store.Close()

	from := time.Date(2026, 7, 29, 15, 0, 0, 0, time.UTC)
	to := from.Add(1 * time.Hour)

	results, err := store.Range(from, to)
	if err != nil {
		t.Fatalf("Range() failed: %v", err)
	}

	if len(results) != 0 {
		t.Fatalf("expected 0 snapshots, got %d", len(results))
	}
}

func mustOpenStore(t *testing.T, name string) *Store {
	t.Helper()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, name)

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() failed: %v", err)
	}
	return store
}
