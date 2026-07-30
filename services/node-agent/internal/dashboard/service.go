package dashboard

import (
	"sort"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/alert"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/config"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/scheduler"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/version"
)

type Service struct {
	scheduler *scheduler.Scheduler
	engine    *alert.Engine
	cfg       config.Config
	startedAt time.Time
}

func NewService(
	scheduler *scheduler.Scheduler,
	engine *alert.Engine,
	cfg config.Config,
) *Service {
	return &Service{
		scheduler: scheduler,
		engine:    engine,
		cfg:       cfg,
		startedAt: time.Now(),
	}
}

func (s *Service) Overview() Overview {
	snapshot, _ := s.scheduler.Snapshot()
	events := s.engine.ActiveEvents(snapshot)
	if events == nil {
		events = []alert.Event{}
	}
	stats := s.scheduler.Stats()

	status := "healthy"
	for _, event := range events {
		if event.Severity == alert.SeverityCritical {
			status = "critical"
			break
		}
		if event.Severity == alert.SeverityWarning && status == "healthy" {
			status = "warning"
		}
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].TriggeredAt.Before(events[j].TriggeredAt)
	})

	return Overview{
		NodeName:       s.cfg.Agent.Name,
		Version:        version.Build.Version,
		Status:         status,
		Uptime:         time.Since(s.startedAt),
		Snapshot:       snapshot,
		ActiveAlerts:   events,
		LastCollection: stats.LastCollectionAt,
	}
}
