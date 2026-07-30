package dashboard

import (
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/alert"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/metrics"
)

type Overview struct {
	NodeName       string        `json:"nodeName"`
	Version        string        `json:"version"`
	Status         string        `json:"status"`
	UptimeMs       int64         `json:"uptimeMs"`
	Snapshot       metrics.Info  `json:"snapshot"`
	ActiveAlerts   []alert.Event `json:"activeAlerts"`
	LastCollection time.Time     `json:"lastCollection"`
}
