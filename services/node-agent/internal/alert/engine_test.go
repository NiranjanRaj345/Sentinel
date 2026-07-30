package alert

import (
	"io"
	"testing"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/metrics"
)

func TestEvaluate_NoRules(t *testing.T) {
	e := New(nil, logger.New(logger.Info, io.Discard))

	events := e.Evaluate(metrics.Info{})

	if len(events) != 0 {
		t.Fatalf("expected 0 events, got %d", len(events))
	}
}

func TestEvaluate_CPUWarning(t *testing.T) {
	e := New([]Rule{
		{
			ID:        "cpu-warn",
			Name:      "CPU Usage",
			Metric:    "cpu.usage",
			Operator:  GreaterThan,
			Threshold: 80.0,
			Severity:  SeverityWarning,
			Enabled:   true,
		},
	}, logger.New(logger.Info, io.Discard))

	events := e.Evaluate(metrics.Info{
		CPU: metrics.CPUInfo{UsagePercent: 92.4},
	})

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	event := events[0]
	if event.RuleID != "cpu-warn" {
		t.Fatalf("expected rule ID cpu-warn, got %s", event.RuleID)
	}
	if event.Value != 92.4 {
		t.Fatalf("expected value 92.4, got %f", event.Value)
	}
	if event.Threshold != 80.0 {
		t.Fatalf("expected threshold 80.0, got %f", event.Threshold)
	}
	if event.Severity != SeverityWarning {
		t.Fatalf("expected severity warning, got %s", event.Severity)
	}
}

func TestEvaluate_CPUCritical(t *testing.T) {
	e := New([]Rule{
		{
			ID:        "cpu-crit",
			Name:      "CPU Usage",
			Metric:    "cpu.usage",
			Operator:  GreaterThanOrEqual,
			Threshold: 90.0,
			Severity:  SeverityCritical,
			Enabled:   true,
		},
	}, logger.New(logger.Info, io.Discard))

	events := e.Evaluate(metrics.Info{
		CPU: metrics.CPUInfo{UsagePercent: 90.0},
	})

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	if events[0].Severity != SeverityCritical {
		t.Fatalf("expected severity critical, got %s", events[0].Severity)
	}
}

func TestEvaluate_DisabledRule(t *testing.T) {
	e := New([]Rule{
		{
			ID:        "cpu-disabled",
			Name:      "CPU Usage",
			Metric:    "cpu.usage",
			Operator:  GreaterThan,
			Threshold: 80.0,
			Severity:  SeverityWarning,
			Enabled:   false,
		},
	}, logger.New(logger.Info, io.Discard))

	events := e.Evaluate(metrics.Info{
		CPU: metrics.CPUInfo{UsagePercent: 95.0},
	})

	if len(events) != 0 {
		t.Fatalf("expected 0 events for disabled rule, got %d", len(events))
	}
}

func TestEvaluate_UnknownMetric(t *testing.T) {
	e := New([]Rule{
		{
			ID:        "unknown",
			Name:      "Unknown Metric",
			Metric:    "foo.bar",
			Operator:  GreaterThan,
			Threshold: 50.0,
			Severity:  SeverityWarning,
			Enabled:   true,
		},
	}, logger.New(logger.Info, io.Discard))

	events := e.Evaluate(metrics.Info{})

	if len(events) != 0 {
		t.Fatalf("expected 0 events for unknown metric, got %d", len(events))
	}
}

func TestEvaluate_MultipleAlerts(t *testing.T) {
	e := New([]Rule{
		{
			ID:        "cpu-warn",
			Name:      "CPU Usage",
			Metric:    "cpu.usage",
			Operator:  GreaterThan,
			Threshold: 80.0,
			Severity:  SeverityWarning,
			Enabled:   true,
		},
		{
			ID:        "mem-warn",
			Name:      "Memory Usage",
			Metric:    "memory.used_percent",
			Operator:  GreaterThan,
			Threshold: 85.0,
			Severity:  SeverityWarning,
			Enabled:   true,
		},
		{
			ID:        "disk-crit",
			Name:      "Disk Usage",
			Metric:    "disk.used_percent",
			Operator:  GreaterThan,
			Threshold: 95.0,
			Severity:  SeverityCritical,
			Enabled:   true,
		},
	}, logger.New(logger.Info, io.Discard))

	events := e.Evaluate(metrics.Info{
		CPU:    metrics.CPUInfo{UsagePercent: 92.0},
		Memory: metrics.MemoryInfo{UsagePercent: 88.0},
		Disks: []metrics.DiskInfo{
			{UsagePercent: 96.0},
		},
	})

	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}

	ruleIDs := make(map[string]bool)
	for _, event := range events {
		ruleIDs[event.RuleID] = true
	}

	if !ruleIDs["cpu-warn"] {
		t.Error("expected cpu-warn event")
	}
	if !ruleIDs["mem-warn"] {
		t.Error("expected mem-warn event")
	}
	if !ruleIDs["disk-crit"] {
		t.Error("expected disk-crit event")
	}
}

func TestEvaluate_LessThanOperator(t *testing.T) {
	e := New([]Rule{
		{
			ID:        "disk-free",
			Name:      "Free Disk",
			Metric:    "disk.used_percent",
			Operator:  LessThan,
			Threshold: 10.0,
			Severity:  SeverityInfo,
			Enabled:   true,
		},
	}, logger.New(logger.Info, io.Discard))

	events := e.Evaluate(metrics.Info{
		Disks: []metrics.DiskInfo{
			{UsagePercent: 5.0},
		},
	})

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
}

func TestEvaluate_GreaterThanOrEqual(t *testing.T) {
	e := New([]Rule{
		{
			ID:        "cpu-crit",
			Name:      "CPU Usage",
			Metric:    "cpu.usage",
			Operator:  GreaterThanOrEqual,
			Threshold: 90.0,
			Severity:  SeverityCritical,
			Enabled:   true,
		},
	}, logger.New(logger.Info, io.Discard))

	events := e.Evaluate(metrics.Info{
		CPU: metrics.CPUInfo{UsagePercent: 90.0},
	})

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
}
