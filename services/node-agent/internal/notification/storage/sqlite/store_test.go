package sqlite

import (
	"context"
	"testing"
	"time"

	notifrepo "github.com/NiranjanRaj345/sentinel/services/node-agent/internal/notification"
)

func openTempNotificationStore(t *testing.T) *Store {
	t.Helper()
	store, err := OpenNotifications(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("open temp notification store: %v", err)
	}
	return store
}

func TestStore_Save_And_Recent(t *testing.T) {
	store := openTempNotificationStore(t)
	defer store.Close()

	ctx := context.Background()
	now := time.Now().UTC()
	n := notifrepo.Notification{
		ID:        "notif-1",
		Title:     "Test",
		Message:   "Message",
		Severity:  notifrepo.SeverityInfo,
		Source:    notifrepo.SourceAlert,
		Status:    notifrepo.StatusPending,
		CreatedAt: now,
	}

	if err := store.Save(ctx, n); err != nil {
		t.Fatalf("save notification: %v", err)
	}

	recent, err := store.Recent(10)
	if err != nil {
		t.Fatalf("recent notifications: %v", err)
	}

	if len(recent) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(recent))
	}

	got := recent[0]
	if got.ID != n.ID {
		t.Fatalf("expected id %s, got %s", n.ID, got.ID)
	}
	if got.Title != n.Title {
		t.Fatalf("expected title %s, got %s", n.Title, got.Title)
	}
	if got.Status != n.Status {
		t.Fatalf("expected status %s, got %s", n.Status, got.Status)
	}
	if got.SentAt != nil {
		t.Fatalf("expected nil sentAt, got %v", got.SentAt)
	}
}

func TestStore_UpdateStatus(t *testing.T) {
	store := openTempNotificationStore(t)
	defer store.Close()

	ctx := context.Background()
	now := time.Now().UTC()
	n := notifrepo.Notification{
		ID:        "notif-2",
		Title:     "Test",
		Message:   "Message",
		Severity:  notifrepo.SeverityWarning,
		Source:    notifrepo.SourceSystem,
		Status:    notifrepo.StatusPending,
		CreatedAt: now,
	}

	if err := store.Save(ctx, n); err != nil {
		t.Fatalf("save notification: %v", err)
	}

	sentAt := time.Now().UTC()
	if err := store.UpdateStatus(ctx, n.ID, notifrepo.StatusSent, &sentAt, ""); err != nil {
		t.Fatalf("update status: %v", err)
	}

	recent, err := store.Recent(10)
	if err != nil {
		t.Fatalf("recent notifications: %v", err)
	}

	if len(recent) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(recent))
	}

	got := recent[0]
	if got.Status != notifrepo.StatusSent {
		t.Fatalf("expected status sent, got %s", got.Status)
	}
	if got.SentAt == nil {
		t.Fatalf("expected non-nil sentAt")
	}
	if got.Error != "" {
		t.Fatalf("expected empty error, got %s", got.Error)
	}
}

func TestStore_Recent_Empty(t *testing.T) {
	store := openTempNotificationStore(t)
	defer store.Close()

	recent, err := store.Recent(10)
	if err != nil {
		t.Fatalf("recent notifications: %v", err)
	}

	if len(recent) != 0 {
		t.Fatalf("expected 0 notifications, got %d", len(recent))
	}
}

func TestStore_Save_Multiple_OrderedByCreatedAtDesc(t *testing.T) {
	store := openTempNotificationStore(t)
	defer store.Close()

	ctx := context.Background()
	for i := 0; i < 3; i++ {
		n := notifrepo.Notification{
			ID:        "notif-" + string(rune('A'+i)),
			Title:     "Title " + string(rune('A'+i)),
			Message:   "Message",
			Severity:  notifrepo.SeverityInfo,
			Source:    notifrepo.SourceSystem,
			Status:    notifrepo.StatusPending,
			CreatedAt: time.Now().UTC().Add(time.Duration(i) * time.Second),
		}
		if err := store.Save(ctx, n); err != nil {
			t.Fatalf("save notification: %v", err)
		}
	}

	recent, err := store.Recent(10)
	if err != nil {
		t.Fatalf("recent notifications: %v", err)
	}

	if len(recent) != 3 {
		t.Fatalf("expected 3 notifications, got %d", len(recent))
	}

	for i := 1; i < len(recent); i++ {
		if !recent[i].CreatedAt.Before(recent[i-1].CreatedAt) {
			t.Fatalf("expected notifications ordered by created_at desc")
		}
	}
}
