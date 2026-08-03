package telegram

import (
	"fmt"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/notification"
)

func FormatNotification(n notification.Notification) string {
	switch n.Source {
	case notification.SourceRecovery:
		return FormatRecovery(n)
	case notification.SourceGuardian:
		return FormatGuardian(n)
	case notification.SourceAlert:
		return FormatAlert(n)
	case notification.SourceOperations:
		return FormatOperation(n)
	case notification.SourceObserver:
		return FormatObserver(n)
	default:
		return FormatGeneric(n)
	}
}

func FormatRecovery(n notification.Notification) string {
	emoji := "✅"
	if n.Severity == notification.SeverityCritical {
		emoji = "🚨"
	}
	return fmt.Sprintf(`🛰️ Sentinel

%s Recovery %s

%s

Time: %s`,
		emoji,
		statusText(n.Status),
		n.Message,
		n.CreatedAt.Format("2006-01-02 15:04"),
	)
}

func FormatGuardian(n notification.Notification) string {
	emoji := "⚡"
	if n.Severity == notification.SeverityWarning {
		emoji = "⚠️"
	}
	return fmt.Sprintf(`🛰️ Sentinel

%s Guardian

%s

Time: %s`,
		emoji,
		n.Message,
		n.CreatedAt.Format("2006-01-02 15:04"),
	)
}

func FormatAlert(n notification.Notification) string {
	emoji := "🔥"
	if n.Severity == notification.SeverityWarning {
		emoji = "⚠️"
	}
	return fmt.Sprintf(`🛰️ Sentinel

%s Critical Alert

%s

Time: %s`,
		emoji,
		n.Message,
		n.CreatedAt.Format("2006-01-02 15:04"),
	)
}

func FormatOperation(n notification.Notification) string {
	emoji := "🛠️"
	if n.Severity == notification.SeverityCritical {
		emoji = "❌"
	}
	return fmt.Sprintf(`🛰️ Sentinel

%s Operation

%s

Time: %s`,
		emoji,
		n.Message,
		n.CreatedAt.Format("2006-01-02 15:04"),
	)
}

func FormatObserver(n notification.Notification) string {
	emoji := "🌡️"
	if n.Severity == notification.SeverityCritical {
		emoji = "🚨"
	}
	return fmt.Sprintf(`🛰️ Sentinel

%s Observer

%s

Time: %s`,
		emoji,
		n.Message,
		n.CreatedAt.Format("2006-01-02 15:04"),
	)
}

func FormatGeneric(n notification.Notification) string {
	emoji := "ℹ️"
	if n.Severity == notification.SeverityWarning {
		emoji = "⚠️"
	}
	if n.Severity == notification.SeverityCritical {
		emoji = "🚨"
	}
	return fmt.Sprintf(`🛰️ Sentinel

%s %s

%s

Time: %s`,
		emoji,
		n.Title,
		n.Message,
		n.CreatedAt.Format("2006-01-02 15:04"),
	)
}

func statusText(status notification.Status) string {
	switch status {
	case notification.StatusSent:
		return "Successful"
	case notification.StatusFailed:
		return "Failed"
	case notification.StatusPending:
		return "Pending"
	default:
		return string(status)
	}
}

func FormatTestNotification(nodeName string) string {
	return fmt.Sprintf(`🛰️ Sentinel

✅ Test notification

Node: %s

Time: %s`,
		nodeName,
		time.Now().Format("2006-01-02 15:04"),
	)
}
