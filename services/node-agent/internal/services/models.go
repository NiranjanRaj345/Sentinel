package services

type ServiceStatus string

const (
	ServiceStatusRunning ServiceStatus = "running"
	ServiceStatusStopped ServiceStatus = "stopped"
	ServiceStatusFailed  ServiceStatus = "failed"
	ServiceStatusUnknown ServiceStatus = "unknown"
)

type Action string

const (
	ActionStart   Action = "start"
	ActionStop    Action = "stop"
	ActionRestart Action = "restart"
)

type ServiceItem struct {
	Name       string        `json:"name"`
	Status     ServiceStatus `json:"status"`
	Action     Action        `json:"action,omitempty"`
	Message    string        `json:"message,omitempty"`
}
