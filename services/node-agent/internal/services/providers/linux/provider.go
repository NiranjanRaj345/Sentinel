package linux

import (
	"context"
	"os/exec"
	"strings"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/services"
)

type linuxProvider struct {
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

func NewLinuxProvider(log *logger.Logger) services.Provider {
	if log == nil {
		log = logger.New(logger.Info, nil)
	}
	return &linuxProvider{log: log, runner: osRunner{}}
}

func (p *linuxProvider) Name() string {
	return "linux"
}

func (p *linuxProvider) List(ctx context.Context) ([]services.ServiceItem, error) {
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

func (p *linuxProvider) Execute(ctx context.Context, item services.ServiceItem) (services.ServiceItem, error) {
	unit := item.Name + ".service"

	var cmd string
	var args []string
	switch item.Action {
	case services.ActionStart:
		cmd = "systemctl"
		args = []string{"start", unit}
	case services.ActionStop:
		cmd = "systemctl"
		args = []string{"stop", unit}
	case services.ActionRestart:
		cmd = "systemctl"
		args = []string{"restart", unit}
	default:
		return services.ServiceItem{Name: item.Name, Status: services.ServiceStatusUnknown, Message: "unsupported action"}, nil
	}

	out, err := p.runner.run(ctx, cmd, args...)
	status, statusOut := p.activeStatus(ctx, unit)
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
	p.log.Info("service operation %s %s: %v", item.Action, unit, result.Status)
	return result, nil
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
