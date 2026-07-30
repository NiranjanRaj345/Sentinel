package rules

import "testing"

func TestEvaluator_Equals(t *testing.T) {
	e := NewEvaluator()
	rule := Rule{
		Enabled: true,
		Conditions: []Condition{
			{Field: "severity", Operator: OpEquals, Value: "critical"},
		},
	}
	event := Event{Severity: "critical"}
	if !e.Evaluate(event, rule) {
		t.Fatal("expected match")
	}
}

func TestEvaluator_NotEquals(t *testing.T) {
	e := NewEvaluator()
	rule := Rule{
		Enabled: true,
		Conditions: []Condition{
			{Field: "severity", Operator: OpNotEquals, Value: "critical"},
		},
	}
	event := Event{Severity: "warning"}
	if !e.Evaluate(event, rule) {
		t.Fatal("expected match")
	}
}

func TestEvaluator_Contains(t *testing.T) {
	e := NewEvaluator()
	rule := Rule{
		Enabled: true,
		Conditions: []Condition{
			{Field: "message", Operator: OpContains, Value: "CPU"},
		},
	}
	event := Event{Message: "CPU usage high"}
	if !e.Evaluate(event, rule) {
		t.Fatal("expected match")
	}
}

func TestEvaluator_DisabledRule(t *testing.T) {
	e := NewEvaluator()
	rule := Rule{
		Enabled: false,
		Conditions: []Condition{
			{Field: "severity", Operator: OpEquals, Value: "critical"},
		},
	}
	event := Event{Severity: "critical"}
	if e.Evaluate(event, rule) {
		t.Fatal("expected no match for disabled rule")
	}
}

func TestEvaluator_MultipleConditions_AllMustMatch(t *testing.T) {
	e := NewEvaluator()
	rule := Rule{
		Enabled: true,
		Conditions: []Condition{
			{Field: "type", Operator: OpEquals, Value: "alert"},
			{Field: "severity", Operator: OpEquals, Value: "critical"},
		},
	}
	event := Event{Type: "alert", Severity: "critical"}
	if !e.Evaluate(event, rule) {
		t.Fatal("expected match when all conditions match")
	}
}

func TestEvaluator_MetadataField(t *testing.T) {
	e := NewEvaluator()
	rule := Rule{
		Enabled: true,
		Conditions: []Condition{
			{Field: "ruleId", Operator: OpEquals, Value: "cpu-warning"},
		},
	}
	event := Event{Metadata: map[string]interface{}{"ruleId": "cpu-warning"}}
	if !e.Evaluate(event, rule) {
		t.Fatal("expected match on metadata field")
	}
}
