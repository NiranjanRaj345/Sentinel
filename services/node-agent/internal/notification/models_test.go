package notification

import (
	"testing"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/events"
)

func TestNewNotification_SetsDefaults(t *testing.T) {
	before := time.Now().UTC()
	n := NewNotification("id-1", "Title", "Message", SeverityCritical, SourceAlert)
	after := time.Now().UTC()

	if n.ID != "id-1" {
		t.Fatalf("expected id id-1, got %s", n.ID)
	}
	if n.Title != "Title" {
		t.Fatalf("expected title Title, got %s", n.Title)
	}
	if n.Message != "Message" {
		t.Fatalf("expected message Message, got %s", n.Message)
	}
	if n.Severity != SeverityCritical {
		t.Fatalf("expected severity critical, got %s", n.Severity)
	}
	if n.Source != SourceAlert {
		t.Fatalf("expected source alert, got %s", n.Source)
	}
	if n.Status != StatusPending {
		t.Fatalf("expected status pending, got %s", n.Status)
	}
	if n.CreatedAt.Before(before) || n.CreatedAt.After(after) {
		t.Fatalf("expected createdAt between %v and %v, got %v", before, after, n.CreatedAt)
	}
	if n.SentAt != nil {
		t.Fatalf("expected nil sentAt, got %v", n.SentAt)
	}
	if n.Provider != "" {
		t.Fatalf("expected empty provider, got %s", n.Provider)
	}
	if n.Error != "" {
		t.Fatalf("expected empty error, got %s", n.Error)
	}
}

func TestRecentResponse_EmptyList(t *testing.T) {
	r := RecentResponse{Notifications: []Notification{}}
	if len(r.Notifications) != 0 {
		t.Fatalf("expected empty notifications list")
	}
}

func TestStatus_Constants(t *testing.T) {
	if StatusPending != "pending" {
		t.Fatalf("expected StatusPending=pending, got %s", StatusPending)
	}
	if StatusSent != "sent" {
		t.Fatalf("expected StatusSent=sent, got %s", StatusSent)
	}
	if StatusFailed != "failed" {
		t.Fatalf("expected StatusFailed=failed, got %s", StatusFailed)
	}
}

func TestSeverity_Constants(t *testing.T) {
	if SeverityInfo != events.SeverityInfo {
		t.Fatalf("expected SeverityInfo to match events.SeverityInfo")
	}
	if SeverityWarning != events.SeverityWarning {
		t.Fatalf("expected SeverityWarning to match events.SeverityWarning")
	}
	if SeverityCritical != events.SeverityCritical {
		t.Fatalf("expected SeverityCritical to match events.SeverityCritical")
	}
}
