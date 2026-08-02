package events

import "time"

type EventType string

const (
	EventTypeOperation EventType = "operation"
	EventTypeAlert     EventType = "alert"
	EventTypeSystem    EventType = "system"
)

type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"
)

type Source string

const (
	SourceOperations Source = "operations"
	SourceAlert      Source = "alert"
	SourceScheduler  Source = "scheduler"
	SourceResources  Source = "resources"
)

type Event struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	Severity  Severity               `json:"severity"`
	Source    Source                 `json:"source"`
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt time.Time              `json:"createdAt"`
}

func OperationSuccess(action, message string) Event {
	return Event{
		Type:     EventTypeOperation,
		Severity: SeverityInfo,
		Source:   SourceOperations,
		Title:    "Operation succeeded: " + action,
		Message:  message,
		Metadata: map[string]interface{}{"action": action, "success": true},
	}
}

func OperationFailure(action, message string) Event {
	return Event{
		Type:     EventTypeOperation,
		Severity: SeverityWarning,
		Source:   SourceOperations,
		Title:    "Operation failed: " + action,
		Message:  message,
		Metadata: map[string]interface{}{"action": action, "success": false},
	}
}

func AlertRaised(ruleID, ruleName string, severity Severity, value, threshold float64) Event {
	return Event{
		Type:     EventTypeAlert,
		Severity: severity,
		Source:   SourceAlert,
		Title:    "Alert raised: " + ruleName,
		Message:  "Rule " + ruleID + " triggered",
		Metadata: map[string]interface{}{
			"ruleId":    ruleID,
			"ruleName":  ruleName,
			"raised":    true,
			"value":     value,
			"threshold": threshold,
		},
	}
}

func AlertCleared(ruleID, ruleName string) Event {
	return Event{
		Type:     EventTypeAlert,
		Severity: SeverityInfo,
		Source:   SourceAlert,
		Title:    "Alert cleared: " + ruleName,
		Message:  "Rule " + ruleID + " no longer active",
		Metadata: map[string]interface{}{"ruleId": ruleID, "ruleName": ruleName, "raised": false},
	}
}

func SystemEvent(title, message string) Event {
	return Event{
		Type:     EventTypeSystem,
		Severity: SeverityInfo,
		Source:   SourceScheduler,
		Title:    title,
		Message:  message,
		Metadata: map[string]interface{}{},
	}
}

func ResourceHealthChanged(name string, health string, message string) Event {
	severity := SeverityInfo
	if health == "unavailable" {
		severity = SeverityCritical
	} else if health == "degraded" {
		severity = SeverityWarning
	}

	return Event{
		Type:     EventTypeSystem,
		Severity: severity,
		Source:   SourceResources,
		Title:    "Resource health changed: " + name,
		Message:  message,
		Metadata: map[string]interface{}{"resource": name, "health": health},
	}
}
