package windows

import (
	"context"
	"os/exec"
	"strings"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/resources"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/services"
)

type windowsProvider struct {
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

func NewWindowsProvider(log *logger.Logger) resources.Provider {
	if log == nil {
		log = logger.New(logger.Info, nil)
	}
	known := []resources.Resource{
		{Name: "Sunshine", Type: resources.ResourceTypeRemoteDesktop, Description: "Remote desktop host", Provider: "windows"},
		{Name: "Docker", Type: resources.ResourceTypeContainerRuntime, Description: "Container runtime", Provider: "windows"},
		{Name: "Tailscale", Type: resources.ResourceTypeVPN, Description: "Mesh VPN", Provider: "windows"},
		{Name: "Jellyfin", Type: resources.ResourceTypeMediaServer, Description: "Media server", Provider: "windows"},
		{Name: "Prometheus", Type: resources.ResourceTypeMonitoring, Description: "Metrics collector", Provider: "windows"},
		{Name: "Grafana", Type: resources.ResourceTypeMonitoring, Description: "Metrics dashboard", Provider: "windows"},
		{Name: "Sentinel", Type: resources.ResourceTypeApplication, Description: "Sentinel node agent", Provider: "windows"},
	}
	return &windowsProvider{log: log, knownResources: known, runner: osRunner{}}
}

func (p *windowsProvider) Name() string {
	return "windows"
}

func (p *windowsProvider) List(ctx context.Context) ([]resources.Resource, error) {
	serviceItems, err := p.listServices(ctx)
	if err != nil {
		return nil, err
	}

	statusMap := make(map[string]services.ServiceItem)
	for _, item := range serviceItems {
		statusMap[strings.ToLower(item.Name)] = item
	}

	var result []resources.Resource
	for _, known := range p.knownResources {
		item, ok := statusMap[strings.ToLower(known.Name)]
		if !ok {
			item = services.ServiceItem{Name: known.Name, Status: services.ServiceStatusUnknown}
		}
		result = append(result, resources.Resource{
			Name:     known.Name,
			Type:     known.Type,
			Provider: known.Provider,
			Status:   string(item.Status),
			Message:  item.Message,
			Health:   p.mapHealth(item.Status),
		})
	}
	return result, nil
}

func (p *windowsProvider) Execute(ctx context.Context, action resources.ResourceAction, name string) (resources.Resource, error) {
	var cmd string
	var args []string
	switch action {
	case resources.ResourceActionStart:
		cmd = "powershell"
		args = []string{"-NoProfile", "-Command", "Start-Service -Name '" + name + "'"}
	case resources.ResourceActionStop:
		cmd = "powershell"
		args = []string{"-NoProfile", "-Command", "Stop-Service -Name '" + name + "'"}
	case resources.ResourceActionRestart:
		cmd = "powershell"
		args = []string{"-NoProfile", "-Command", "Restart-Service -Name '" + name + "'"}
	default:
		return resources.Resource{Name: name, Health: resources.HealthUnknown, Message: "unsupported action"}, nil
	}

	out, err := p.runner.run(ctx, cmd, args...)
	status, statusOut := p.activeStatus(ctx, name)
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
	p.log.Info("resource operation %s %s: %v", action, name, resource.Health)
	return resource, nil
}

func (p *windowsProvider) listServices(ctx context.Context) ([]services.ServiceItem, error) {
	out, err := p.runner.run(ctx, "powershell", "-NoProfile", "-Command", "Get-Service | Select-Object Name,Status | ConvertTo-Csv -NoTypeInformation")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(out, "\n")
	var items []services.ServiceItem
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.EqualFold(line, "Name,Status") {
			continue
		}
		parts := strings.SplitN(line, ",", 2)
		if len(parts) < 2 {
			continue
		}
		name := strings.Trim(parts[0], "\"")
		status := windowsServiceStatus(strings.Trim(parts[1], "\""))
		items = append(items, services.ServiceItem{Name: name, Status: status})
	}
	return items, nil
}

func (p *windowsProvider) activeStatus(ctx context.Context, name string) (services.ServiceStatus, string) {
	out, err := p.runner.run(ctx, "powershell", "-NoProfile", "-Command", "Get-Service -Name '"+name+"' | Select-Object -ExpandProperty Status")
	value := strings.TrimSpace(strings.ToLower(out))
	switch value {
	case "running":
		return services.ServiceStatusRunning, out
	case "stopped":
		return services.ServiceStatusStopped, out
	case "startpending", "stoppending":
		return services.ServiceStatusRunning, out
	default:
		if err != nil {
			return services.ServiceStatusUnknown, err.Error()
		}
		return services.ServiceStatusUnknown, out
	}
}

func (p *windowsProvider) mapHealth(status services.ServiceStatus) resources.Health {
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

func windowsServiceStatus(raw string) services.ServiceStatus {
	switch strings.ToLower(raw) {
	case "running":
		return services.ServiceStatusRunning
	case "stopped":
		return services.ServiceStatusStopped
	case "startpending", "stoppending":
		return services.ServiceStatusRunning
	default:
		return services.ServiceStatusUnknown
	}
}
