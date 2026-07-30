package sqlite

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/rules"
)

func TestOpenRules_CreatesDatabase(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "rules.db")

	store, err := OpenRules(dbPath)
	if err != nil {
		t.Fatalf("OpenRules() failed: %v", err)
	}
	defer store.Close()

	if store == nil {
		t.Fatal("expected non-nil Store")
	}

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatal("expected database file to exist")
	}
}

func TestOpenRules_EmptyPath_ReturnsError(t *testing.T) {
	store, err := OpenRules("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
	if store != nil {
		t.Fatal("expected nil store on error")
	}
}

func TestStore_Save_InsertsRule(t *testing.T) {
	store := mustOpenRulesStore(t, "rules_save.db")
	defer store.Close()

	rule := rules.Rule{
		ID:      "test-rule",
		Name:    "Test Rule",
		Enabled: true,
		Trigger: rules.TriggerEvent,
		Conditions: []rules.Condition{
			{Field: "severity", Operator: rules.OpEquals, Value: "critical"},
		},
		Actions: []rules.Action{rules.ActionNotify},
	}

	if err := store.Save(context.Background(), rule); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	var count int
	if err := store.db.QueryRow(`SELECT COUNT(*) FROM rules`).Scan(&count); err != nil {
		t.Fatalf("count query failed: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 row, got %d", count)
	}
}

func TestStore_Enabled_ReturnsOnlyEnabled(t *testing.T) {
	store := mustOpenRulesStore(t, "rules_enabled.db")
	defer store.Close()

	rulesList := []rules.Rule{
		{ID: "enabled", Name: "Enabled Rule", Enabled: true, Trigger: rules.TriggerEvent, Conditions: nil, Actions: nil},
		{ID: "disabled", Name: "Disabled Rule", Enabled: false, Trigger: rules.TriggerEvent, Conditions: nil, Actions: nil},
	}

	ctx := context.Background()
	for _, r := range rulesList {
		if err := store.Save(ctx, r); err != nil {
			t.Fatalf("Save() failed: %v", err)
		}
	}

	enabled, err := store.Enabled(ctx)
	if err != nil {
		t.Fatalf("Enabled() failed: %v", err)
	}
	if len(enabled) != 1 {
		t.Fatalf("expected 1 enabled rule, got %d", len(enabled))
	}
	if enabled[0].ID != "enabled" {
		t.Fatalf("expected enabled rule, got %s", enabled[0].ID)
	}
}

func TestStore_Close_NilStore_Succeeds(t *testing.T) {
	var s *RulesStore
	if err := s.Close(); err != nil {
		t.Fatalf("Close() on nil store should succeed, got %v", err)
	}
}

func mustOpenRulesStore(t *testing.T, name string) *RulesStore {
	t.Helper()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, name)

	store, err := OpenRules(dbPath)
	if err != nil {
		t.Fatalf("OpenRules() failed: %v", err)
	}
	return store
}
