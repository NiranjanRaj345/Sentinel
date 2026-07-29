package service

import "github.com/NiranjanRaj345/sentinel/services/node-agent/internal/metrics"

// MetricsService provides access to system metrics.
type MetricsService struct{}

// NewMetricsService creates a new MetricsService.
func NewMetricsService() *MetricsService {
	return &MetricsService{}
}

// GetInfo returns the current system metrics.
func (s *MetricsService) GetInfo() (metrics.Info, error) {
	return metrics.GetInfo()
}
