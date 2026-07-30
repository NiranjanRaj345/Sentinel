package config

import (
	"testing"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/alert"
)

func TestConfigValidate_Valid(t *testing.T) {
	cfg := Default()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestServerValidate_InvalidPort(t *testing.T) {
	cfg := Default()
	cfg.Server.Port = 70000
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLoggingValidate_InvalidLevel(t *testing.T) {
	cfg := Default()
	cfg.Logging.Level = "verbose"
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMetricsValidate_InvalidInterval(t *testing.T) {
	cfg := Default()
	cfg.Metrics.Interval = "abc"
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestStorageValidate_EmptyPath(t *testing.T) {
	cfg := Default()
	cfg.Storage.Path = ""
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty storage path, got nil")
	}
}

func TestAlertsValidate_EmptyRuleID(t *testing.T) {
	cfg := Default()
	cfg.Alerts.Rules = []alert.Rule{
		{ID: ""},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty rule ID, got nil")
	}
}

func TestAlertsValidate_InvalidMetric(t *testing.T) {
	cfg := Default()
	cfg.Alerts.Rules = []alert.Rule{
		{
			ID:        "test",
			Name:      "Test",
			Metric:    "invalid.metric",
			Operator:  alert.GreaterThan,
			Threshold: 50,
			Severity:  alert.SeverityWarning,
			Enabled:   true,
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for invalid metric, got nil")
	}
}

func TestAlertsValidate_InvalidOperator(t *testing.T) {
	cfg := Default()
	cfg.Alerts.Rules = []alert.Rule{
		{
			ID:        "test",
			Name:      "Test",
			Metric:    "cpu.usage",
			Operator:  "==",
			Threshold: 50,
			Severity:  alert.SeverityWarning,
			Enabled:   true,
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for invalid operator, got nil")
	}
}

func TestAlertsValidate_InvalidSeverity(t *testing.T) {
	cfg := Default()
	cfg.Alerts.Rules = []alert.Rule{
		{
			ID:        "test",
			Name:      "Test",
			Metric:    "cpu.usage",
			Operator:  alert.GreaterThan,
			Threshold: 50,
			Severity:  "invalid",
			Enabled:   true,
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for invalid severity, got nil")
	}
}

func TestAlertsValidate_ThresholdOutOfRange(t *testing.T) {
	cfg := Default()
	cfg.Alerts.Rules = []alert.Rule{
		{
			ID:        "test",
			Name:      "Test",
			Metric:    "cpu.usage",
			Operator:  alert.GreaterThan,
			Threshold: 150,
			Severity:  alert.SeverityWarning,
			Enabled:   true,
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for out-of-range threshold, got nil")
	}
}
