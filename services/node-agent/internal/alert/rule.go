package alert

import "fmt"

type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"
)

type Operator string

const (
	GreaterThan        Operator = ">"
	GreaterThanOrEqual Operator = ">="
	LessThan           Operator = "<"
	LessThanOrEqual    Operator = "<="
)

type Rule struct {
	ID          string
	Name        string
	Description string

	Metric    string
	Operator  Operator
	Threshold float64

	Severity Severity

	Enabled bool
}

func (r Rule) Validate() error {
	if r.ID == "" {
		return fmt.Errorf("alert rule id cannot be empty")
	}

	if r.Name == "" {
		return fmt.Errorf("alert rule name cannot be empty")
	}

	if !isValidMetric(r.Metric) {
		return fmt.Errorf("alert rule metric must be one of: cpu.usage, memory.used_percent, disk.used_percent")
	}

	if !isValidOperator(r.Operator) {
		return fmt.Errorf("alert rule operator must be one of: >, >=, <, <=")
	}

	if !isValidSeverity(r.Severity) {
		return fmt.Errorf("alert rule severity must be one of: info, warning, critical")
	}

	if r.Threshold < 0 || r.Threshold > 100 {
		return fmt.Errorf("alert rule threshold must be between 0 and 100")
	}

	return nil
}

func isValidMetric(metric string) bool {
	switch metric {
	case "cpu.usage", "memory.used_percent", "disk.used_percent":
		return true
	default:
		return false
	}
}

func isValidOperator(op Operator) bool {
	switch op {
	case GreaterThan, GreaterThanOrEqual, LessThan, LessThanOrEqual:
		return true
	default:
		return false
	}
}

func isValidSeverity(severity Severity) bool {
	switch severity {
	case SeverityInfo, SeverityWarning, SeverityCritical:
		return true
	default:
		return false
	}
}
