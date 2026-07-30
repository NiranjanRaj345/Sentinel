package alert

import (
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/metrics"
)

type Engine struct {
	rules []Rule
	log   *logger.Logger
}

func New(rules []Rule, log *logger.Logger) *Engine {
	return &Engine{
		rules: rules,
		log:   log,
	}
}

func (e *Engine) Evaluate(info metrics.Info) []Event {
	var events []Event

	for _, rule := range e.rules {
		if !rule.Enabled {
			continue
		}

		value, ok := metricValue(info, rule.Metric)
		if !ok {
			continue
		}

		if !compare(value, rule.Operator, rule.Threshold) {
			continue
		}

		event := Event{
			RuleID:      rule.ID,
			RuleName:    rule.Name,
			Metric:      rule.Metric,
			Value:       value,
			Threshold:   rule.Threshold,
			Severity:    rule.Severity,
			TriggeredAt: info.Metadata.Timestamp,
		}

		events = append(events, event)

		e.log.Info(
			"%s %s (%.1f %s %.1f)",
			severityLabel(rule.Severity),
			rule.Name,
			value,
			rule.Operator,
			rule.Threshold,
		)
	}

	return events
}

func metricValue(info metrics.Info, metric string) (float64, bool) {
	switch metric {
	case "cpu.usage":
		return info.CPU.UsagePercent, true

	case "memory.used_percent":
		return info.Memory.UsagePercent, true

	case "disk.used_percent":
		if len(info.Disks) == 0 {
			return 0, false
		}
		return info.Disks[0].UsagePercent, true

	default:
		return 0, false
	}
}

func compare(value float64, op Operator, threshold float64) bool {
	switch op {
	case GreaterThan:
		return value > threshold
	case GreaterThanOrEqual:
		return value >= threshold
	case LessThan:
		return value < threshold
	case LessThanOrEqual:
		return value <= threshold
	default:
		return false
	}
}

func severityLabel(severity Severity) string {
	switch severity {
	case SeverityWarning:
		return "WARNING"
	case SeverityCritical:
		return "CRITICAL"
	default:
		return "INFO"
	}
}
