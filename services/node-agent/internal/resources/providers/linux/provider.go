package linux

import (
	"context"
	"os/exec"
	"strings"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/resources"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/services"
)

type linuxProvider struct {
	log            *logger.Logger
	knownResources []resources.Resource
	runner         commandRunner
}

type commandRunner interface {
	run(ctx context.Context, name string, args ...string) (string, error)
}

type osRunner struct{}

func (osRunner) run(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func NewLinuxProvider(log *logger.Logger) resources.Provider {
	if log == nil {
		log = logger.New(logger.Info, nil)
	}
	known := []resources.Resource{
		{Name: "Sunshine", Type: resources.ResourceTypeRemoteDesktop, Description: "Remote desktop host", Provider: "linux"},
		{Name: "Docker", Type: resources.ResourceTypeContainerRuntime, Description: "Container runtime", Provider: "linux"},
		{Name: "Tailscale", Type: resources.ResourceTypeVPN, Description: "Mesh VPN", Provider: "linux"},
		{Name: "Jellyfin", Type: resources.ResourceTypeMediaServer, Description: "Media server", Provider: "linux"},
		{Name: "Prometheus", Type: resources.ResourceTypeMonitoring, Description: "Metrics collector", Provider: "linux"},
		{Name: "Grafana", Type: resources.ResourceTypeMonitoring, Description: "Metrics dashboard", Provider: "linux"},
		{Name: "Sentinel", Type: resources.ResourceTypeApplication, Description: "Sentinel node agent", Provider: "linux"},
	}
	return &linuxProvider{log: log, knownResources: known, runner: osRunner{}}
}

func (p *linuxProvider) Name() string {
	return "linux"
}

func (p *linuxProvider) List(ctx context.Context) ([]resources.Resource, error) {
	serviceItems, err := p.listServices(ctx)
	if err != nil {
		return nil, err
	}

	statusMap := make(map[string]services.ServiceItem)
	for _, item := range serviceItems {
		statusMap[item.Name] = item
	}

	var result []resources.Resource
	for _, known := range p.knownResources {
		item, ok := statusMap[known.Name]
		if !ok {
			item = services.ServiceItem{Name: known.Name, Status: services.ServiceStatusUnknown}
		}
		resource := p.toResource(known, item)
		result = append(result, resource)
	}
	return result, nil
}

func (p *linuxProvider) Execute(ctx context.Context, action resources.ResourceAction, name string) (resources.Resource, error) {
	unit := name + ".service"

	var cmd string
	var args []string
	switch action {
	case resources.ResourceActionStart:
		cmd = "systemctl"
		args = []string{"start", unit}
	case resources.ResourceActionStop:
		cmd = "systemctl"
		args = []string{"stop", unit}
	case resources.ResourceActionRestart:
		cmd = "systemctl"
		args = []string{"restart", unit}
	default:
		return resources.Resource{Name: name, Health: resources.HealthUnknown, Message: "unsupported action"}, nil
	}

	out, err := p.runner.run(ctx, cmd, args...)
	status, statusOut := p.activeStatus(ctx, unit)
	message := strings.TrimSpace(out)
	if message == "" {
		message = strings.TrimSpace(statusOut)
	}
	if err != nil {
		message = err.Error()
	}

	resource := resources.Resource{
		Name:    name,
		Status:  string(status),
		Message: message,
		Health:  p.mapHealth(status),
	}
	p.log.Info("resource operation %s %s: %v", action, unit, resource.Health)
	return resource, nil
}

func (p *linuxProvider) listServices(ctx context.Context) ([]services.ServiceItem, error) {
	out, err := p.runner.run(ctx, "systemctl", "list-unit-files", "--type=service", "--no-pager", "--no-legend")
	if err != nil {
		return nil, err
	}

	var items []services.ServiceItem
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		name := strings.TrimSuffix(parts[0], ".service")
		status, _ := parseStatus(parts[1])
		items = append(items, services.ServiceItem{Name: name, Status: status})
	}
	return items, nil
}

func (p *linuxProvider) activeStatus(ctx context.Context, unit string) (services.ServiceStatus, string) {
	out, err := p.runner.run(ctx, "systemctl", "is-active", unit)
	value := strings.TrimSpace(strings.ToLower(out))
	switch value {
	case "active":
		return services.ServiceStatusRunning, out
	case "inactive":
		return services.ServiceStatusStopped, out
	case "failed":
		return services.ServiceStatusFailed, out
	default:
		if err != nil {
			return services.ServiceStatusUnknown, err.Error()
		}
		return services.ServiceStatusUnknown, out
	}
}

func (p *linuxProvider) toResource(known resources.Resource, item services.ServiceItem) resources.Resource {
	resource := resources.Resource{
		Name:     known.Name,
		Type:     known.Type,
		Provider: known.Provider,
		Status:   string(item.Status),
		Message:  item.Message,
		Health:   p.mapHealth(item.Status),
	}
	return resource
}

func (p *linuxProvider) mapHealth(status services.ServiceStatus) resources.Health {
	switch status {
	case services.ServiceStatusRunning:
		return resources.HealthHealthy
	case services.ServiceStatusStopped:
		return resources.HealthUnavailable
	case services.ServiceStatusFailed:
		return resources.HealthUnavailable
	default:
		return resources.HealthUnknown
	}
}

func parseStatus(raw string) (services.ServiceStatus, string) {
	switch strings.ToLower(raw) {
	case "enabled":
		return services.ServiceStatusRunning, raw
	case "disabled":
		return services.ServiceStatusStopped, raw
	case "static":
		return services.ServiceStatusStopped, raw
	case "masked":
		return services.ServiceStatusStopped, raw
	default:
		return services.ServiceStatusUnknown, raw
	}
}
