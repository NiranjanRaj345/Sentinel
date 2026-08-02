package resources

import "github.com/NiranjanRaj345/sentinel/services/node-agent/internal/services"

type HealthEvaluator struct{}

func NewHealthEvaluator() *HealthEvaluator {
	return &HealthEvaluator{}
}

func (e *HealthEvaluator) Evaluate(serviceStatus services.ServiceStatus, message string) Health {
	switch serviceStatus {
	case services.ServiceStatusRunning:
		return HealthHealthy
	case services.ServiceStatusStopped:
		return HealthUnavailable
	case services.ServiceStatusFailed:
		return HealthUnavailable
	default:
		return HealthUnknown
	}
}
