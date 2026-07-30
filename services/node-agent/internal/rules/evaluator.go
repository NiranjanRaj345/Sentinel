package rules

import "strings"

type Evaluator struct{}

func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

func (e *Evaluator) Evaluate(event Event, rule Rule) bool {
	if !rule.Enabled {
		return false
	}
	if len(rule.Conditions) == 0 {
		return true
	}
	for _, condition := range rule.Conditions {
		if !e.evaluateCondition(event, condition) {
			return false
		}
	}
	return true
}

func (e *Evaluator) evaluateCondition(event Event, condition Condition) bool {
	value, ok := eventField(event, condition.Field)
	if !ok {
		return false
	}
	switch condition.Operator {
	case OpEquals:
		return value == condition.Value
	case OpNotEquals:
		return value != condition.Value
	case OpContains:
		return contains(value, condition.Value)
	default:
		return false
	}
}

func eventField(event Event, field string) (string, bool) {
	switch field {
	case "type":
		return event.Type, true
	case "severity":
		return event.Severity, true
	case "source":
		return event.Source, true
	case "title":
		return event.Title, true
	case "message":
		return event.Message, true
	default:
		if event.Metadata != nil {
			v, ok := event.Metadata[field]
			if !ok {
				return "", false
			}
			switch typed := v.(type) {
			case string:
				return typed, true
			case int:
				return string(rune(typed)), true
			case float64:
				return string(rune(typed)), true
			case bool:
				if typed {
					return "true", true
				}
				return "false", true
			default:
				return "", false
			}
		}
		return "", false
	}
}

func contains(haystack, needle string) bool {
	return strings.Contains(haystack, needle)
}
