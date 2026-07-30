package events

import (
	"testing"
)

func TestOperationSuccess_Fields(t *testing.T) {
	event := OperationSuccess("restart", "operation completed")

	if event.Type != EventTypeOperation {
		t.Fatalf("expected type %s, got %s", EventTypeOperation, event.Type)
	}
	if event.Severity != SeverityInfo {
		t.Fatalf("expected severity %s, got %s", SeverityInfo, event.Severity)
	}
	if event.Source != SourceOperations {
		t.Fatalf("expected source %s, got %s", SourceOperations, event.Source)
	}
	if event.Metadata["success"] != true {
		t.Fatalf("expected success=true in metadata")
	}
}

func TestOperationFailure_Fields(t *testing.T) {
	event := OperationFailure("shutdown", "permission denied")

	if event.Type != EventTypeOperation {
		t.Fatalf("expected type %s, got %s", EventTypeOperation, event.Type)
	}
	if event.Severity != SeverityWarning {
		t.Fatalf("expected severity %s, got %s", SeverityWarning, event.Severity)
	}
	if event.Metadata["success"] != false {
		t.Fatalf("expected success=false in metadata")
	}
}

func TestAlertRaised_Fields(t *testing.T) {
	event := AlertRaised("cpu-high", "CPU Usage", SeverityCritical, 95.5, 90.0)

	if event.Type != EventTypeAlert {
		t.Fatalf("expected type %s, got %s", EventTypeAlert, event.Type)
	}
	if event.Severity != SeverityCritical {
		t.Fatalf("expected severity %s, got %s", SeverityCritical, event.Severity)
	}
	if event.Source != SourceAlert {
		t.Fatalf("expected source %s, got %s", SourceAlert, event.Source)
	}
	if event.Metadata["raised"] != true {
		t.Fatalf("expected raised=true in metadata")
	}
}

func TestAlertCleared_Fields(t *testing.T) {
	event := AlertCleared("cpu-high", "CPU Usage")

	if event.Type != EventTypeAlert {
		t.Fatalf("expected type %s, got %s", EventTypeAlert, event.Type)
	}
	if event.Source != SourceAlert {
		t.Fatalf("expected source %s, got %s", SourceAlert, event.Source)
	}
	if event.Metadata["raised"] != false {
		t.Fatalf("expected raised=false in metadata")
	}
}

func TestSystemEvent_Fields(t *testing.T) {
	event := SystemEvent("scheduler_started", "background scheduler started")

	if event.Type != EventTypeSystem {
		t.Fatalf("expected type %s, got %s", EventTypeSystem, event.Type)
	}
	if event.Severity != SeverityInfo {
		t.Fatalf("expected severity %s, got %s", SeverityInfo, event.Severity)
	}
	if event.Source != SourceScheduler {
		t.Fatalf("expected source %s, got %s", SourceScheduler, event.Source)
	}
}
