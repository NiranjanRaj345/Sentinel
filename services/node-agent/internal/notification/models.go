package notification

import (
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/events"
)

type Status string

const (
	StatusPending Status = "pending"
	StatusSent    Status = "sent"
	StatusFailed  Status = "failed"
)

type Severity = events.Severity

const (
	SeverityInfo     Severity = events.SeverityInfo
	SeverityWarning  Severity = events.SeverityWarning
	SeverityCritical Severity = events.SeverityCritical
)

type Source = events.Source

const (
	SourceOperations Source = events.SourceOperations
	SourceAlert      Source = events.SourceAlert
	SourceScheduler  Source = events.SourceScheduler
	SourceResources  Source = events.SourceResources
	SourceGuardian   Source = events.SourceGuardian
	SourceRecovery   Source = "recovery"
	SourceAutomation Source = "automation"
	SourceObserver   Source = "observer"
	SourceSystem     Source = events.SourceSystem
)

type Notification struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Message   string     `json:"message"`
	Severity  Severity   `json:"severity"`
	Source    Source     `json:"source"`
	Provider  string     `json:"provider,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	SentAt    *time.Time `json:"sentAt,omitempty"`
	Status    Status     `json:"status"`
	Error     string     `json:"error,omitempty"`
}

type RecentResponse struct {
	Notifications []Notification `json:"notifications"`
}

func NewNotification(
	id string,
	title string,
	message string,
	severity Severity,
	source Source,
) Notification {
	return Notification{
		ID:        id,
		Title:     title,
		Message:   message,
		Severity:  severity,
		Source:    source,
		CreatedAt: time.Now().UTC(),
		Status:    StatusPending,
	}
}
