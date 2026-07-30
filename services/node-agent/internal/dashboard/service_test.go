package dashboard

import (
	"io"
	"testing"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/alert"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/config"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/scheduler"
)

func TestOverview_NoAlerts_ReturnsHealthy(t *testing.T) {
	s := scheduler.New(5*time.Second, logger.New(logger.Info, io.Discard), nil, nil, nil, nil)
	_ = s.Start()
	defer s.Stop()

	engine := alert.New(nil, logger.New(logger.Info, io.Discard))
	cfg := config.Default()

	service := NewService(s, engine, cfg, nil)
	overview := service.Overview()

	if overview.Status != "healthy" {
		t.Fatalf("expected status healthy, got %s", overview.Status)
	}
	if len(overview.ActiveAlerts) != 0 {
		t.Fatalf("expected 0 active alerts, got %d", len(overview.ActiveAlerts))
	}
	if overview.NodeName != cfg.Agent.Name {
		t.Fatalf("expected nodeName %s, got %s", cfg.Agent.Name, overview.NodeName)
	}
}

func TestOverview_WarningAlert_ReturnsWarning(t *testing.T) {
	s := scheduler.New(5*time.Second, logger.New(logger.Info, io.Discard), nil, nil, nil, nil)
	_ = s.Start()
	defer s.Stop()

	rules := []alert.Rule{
		{
			ID:        "cpu-warn",
			Name:      "CPU Usage",
			Metric:    "cpu.usage",
			Operator:  alert.GreaterThan,
			Threshold: 0.0,
			Severity:  alert.SeverityWarning,
			Enabled:   true,
		},
	}
	engine := alert.New(rules, logger.New(logger.Info, io.Discard))
	cfg := config.Default()

	service := NewService(s, engine, cfg, nil)
	overview := service.Overview()

	if overview.Status != "warning" {
		t.Fatalf("expected status warning, got %s", overview.Status)
	}
	if len(overview.ActiveAlerts) != 1 {
		t.Fatalf("expected 1 active alert, got %d", len(overview.ActiveAlerts))
	}
}

func TestOverview_CriticalAlert_ReturnsCritical(t *testing.T) {
	s := scheduler.New(5*time.Second, logger.New(logger.Info, io.Discard), nil, nil, nil, nil)
	_ = s.Start()
	defer s.Stop()

	rules := []alert.Rule{
		{
			ID:        "cpu-crit",
			Name:      "CPU Usage",
			Metric:    "cpu.usage",
			Operator:  alert.GreaterThanOrEqual,
			Threshold: 0.0,
			Severity:  alert.SeverityCritical,
			Enabled:   true,
		},
	}
	engine := alert.New(rules, logger.New(logger.Info, io.Discard))
	cfg := config.Default()

	service := NewService(s, engine, cfg, nil)
	overview := service.Overview()

	if overview.Status != "critical" {
		t.Fatalf("expected status critical, got %s", overview.Status)
	}
	if len(overview.ActiveAlerts) != 1 {
		t.Fatalf("expected 1 active alert, got %d", len(overview.ActiveAlerts))
	}
}

func TestOverview_MixedAlerts_PrefersCritical(t *testing.T) {
	s := scheduler.New(5*time.Second, logger.New(logger.Info, io.Discard), nil, nil, nil, nil)
	_ = s.Start()
	defer s.Stop()

	rules := []alert.Rule{
		{
			ID:        "cpu-warn",
			Name:      "CPU Usage",
			Metric:    "cpu.usage",
			Operator:  alert.GreaterThan,
			Threshold: 0.0,
			Severity:  alert.SeverityWarning,
			Enabled:   true,
		},
		{
			ID:        "mem-crit",
			Name:      "Memory Usage",
			Metric:    "memory.used_percent",
			Operator:  alert.GreaterThan,
			Threshold: 0.0,
			Severity:  alert.SeverityCritical,
			Enabled:   true,
		},
	}
	engine := alert.New(rules, logger.New(logger.Info, io.Discard))
	cfg := config.Default()

	service := NewService(s, engine, cfg, nil)
	overview := service.Overview()

	if overview.Status != "critical" {
		t.Fatalf("expected status critical when both warning and critical present, got %s", overview.Status)
	}
	if len(overview.ActiveAlerts) != 2 {
		t.Fatalf("expected 2 active alerts, got %d", len(overview.ActiveAlerts))
	}
}

func TestOverview_AlertsSortedByTriggerTime(t *testing.T) {
	s := scheduler.New(5*time.Second, logger.New(logger.Info, io.Discard), nil, nil, nil, nil)
	_ = s.Start()
	defer s.Stop()

	rules := []alert.Rule{
		{
			ID:        "cpu-warn",
			Name:      "CPU Usage",
			Metric:    "cpu.usage",
			Operator:  alert.GreaterThan,
			Threshold: 0.0,
			Severity:  alert.SeverityWarning,
			Enabled:   true,
		},
		{
			ID:        "mem-crit",
			Name:      "Memory Usage",
			Metric:    "memory.used_percent",
			Operator:  alert.GreaterThan,
			Threshold: 0.0,
			Severity:  alert.SeverityCritical,
			Enabled:   true,
		},
	}
	engine := alert.New(rules, logger.New(logger.Info, io.Discard))
	cfg := config.Default()

	service := NewService(s, engine, cfg, nil)
	overview := service.Overview()

	if len(overview.ActiveAlerts) != 2 {
		t.Fatalf("expected 2 active alerts, got %d", len(overview.ActiveAlerts))
	}
	for i := 1; i < len(overview.ActiveAlerts); i++ {
		if overview.ActiveAlerts[i].TriggeredAt.Before(overview.ActiveAlerts[i-1].TriggeredAt) {
			t.Fatalf("alerts not sorted: %v before %v", overview.ActiveAlerts[i].TriggeredAt, overview.ActiveAlerts[i-1].TriggeredAt)
		}
	}
}
