package windows

import (
	"context"
	"os/exec"
	"strings"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/services"
)

type windowsProvider struct {
	log    *logger.Logger
	runner commandRunner
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

func NewWindowsProvider(log *logger.Logger) services.Provider {
	if log == nil {
		log = logger.New(logger.Info, nil)
	}
	return &windowsProvider{log: log, runner: osRunner{}}
}

func (p *windowsProvider) Name() string {
	return "windows"
}

func (p *windowsProvider) List(ctx context.Context) ([]services.ServiceItem, error) {
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

func (p *windowsProvider) Execute(ctx context.Context, item services.ServiceItem) (services.ServiceItem, error) {
	var cmd string
	var args []string
	switch item.Action {
	case services.ActionStart:
		cmd = "powershell"
		args = []string{"-NoProfile", "-Command", "Start-Service -Name '" + item.Name + "'"}
	case services.ActionStop:
		cmd = "powershell"
		args = []string{"-NoProfile", "-Command", "Stop-Service -Name '" + item.Name + "'"}
	case services.ActionRestart:
		cmd = "powershell"
		args = []string{"-NoProfile", "-Command", "Restart-Service -Name '" + item.Name + "'"}
	default:
		return services.ServiceItem{Name: item.Name, Status: services.ServiceStatusUnknown, Message: "unsupported action"}, nil
	}

	out, err := p.runner.run(ctx, cmd, args...)
	status, statusOut := p.activeStatus(ctx, item.Name)
	result := services.ServiceItem{Name: item.Name, Status: status, Action: item.Action, Message: strings.TrimSpace(out)}
	if err != nil {
		result.Status = services.ServiceStatusFailed
		result.Message = strings.TrimSpace(out)
		if result.Message == "" {
			result.Message = err.Error()
		}
	}
	if result.Message == "" {
		result.Message = strings.TrimSpace(statusOut)
	}
	p.log.Info("service operation %s %s: %v", item.Action, item.Name, result.Status)
	return result, nil
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
