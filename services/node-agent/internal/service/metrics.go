package service

import (
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/metrics"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/scheduler"
)

type MetricsService struct {
	scheduler *scheduler.Scheduler
}

func NewMetricsService(scheduler *scheduler.Scheduler) *MetricsService {
	return &MetricsService{
		scheduler: scheduler,
	}
}

func (s *MetricsService) GetInfo() (metrics.Info, error) {
	return s.scheduler.Snapshot()
}
