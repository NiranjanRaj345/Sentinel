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
	Uptime         time.Duration `json:"uptime"`
	Snapshot       metrics.Info  `json:"snapshot"`
	ActiveAlerts   []alert.Event `json:"activeAlerts"`
	LastCollection time.Time     `json:"lastCollection"`
}
