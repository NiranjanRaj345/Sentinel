package alert

import "time"

type Event struct {
	RuleID   string
	RuleName string

	Metric string

	Value float64

	Threshold float64

	Severity Severity

	TriggeredAt time.Time
}
