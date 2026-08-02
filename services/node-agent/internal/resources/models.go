package resources

type ResourceType string

const (
	ResourceTypeRemoteDesktop ResourceType = "remote_desktop"
	ResourceTypeVPN           ResourceType = "vpn"
	ResourceTypeContainerRuntime ResourceType = "container_runtime"
	ResourceTypeMediaServer   ResourceType = "media_server"
	ResourceTypeDatabase      ResourceType = "database"
	ResourceTypeMonitoring    ResourceType = "monitoring"
	ResourceTypeApplication   ResourceType = "application"
)

type Health string

const (
	HealthHealthy     Health = "healthy"
	HealthDegraded    Health = "degraded"
	HealthUnavailable Health = "unavailable"
	HealthUnknown      Health = "unknown"
)

type Resource struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Type        ResourceType `json:"type"`
	Health      Health     `json:"health"`
	Version     string     `json:"version,omitempty"`
	Description string     `json:"description,omitempty"`
	Provider    string     `json:"provider,omitempty"`
	Status      string     `json:"status,omitempty"`
	Message     string     `json:"message,omitempty"`
}

type ResourceAction string

const (
	ResourceActionStart   ResourceAction = "start"
	ResourceActionStop    ResourceAction = "stop"
	ResourceActionRestart ResourceAction = "restart"
)
