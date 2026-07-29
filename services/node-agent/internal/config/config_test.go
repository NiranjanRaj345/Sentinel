package config

import (
	"testing"
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
