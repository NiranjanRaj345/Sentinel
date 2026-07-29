package sqlite

import (
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
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
