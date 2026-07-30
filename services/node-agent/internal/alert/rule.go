package alert

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
