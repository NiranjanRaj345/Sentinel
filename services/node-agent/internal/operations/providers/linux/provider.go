package linux

import (
	"context"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/operations"
)

type linuxProvider struct {
	log  *logger.Logger
	runner operations.CommandRunner
}

func NewLinuxProvider(log *logger.Logger, runner operations.CommandRunner) operations.Provider {
	if runner == nil {
		runner = operations.NewOSRunner()
	}
	return &linuxProvider{log: log, runner: runner}
}

func (p *linuxProvider) Name() string {
	return "linux"
}

func (p *linuxProvider) Supports(action operations.Action) bool {
	switch action {
	case operations.ActionSleep, operations.ActionRestart, operations.ActionShutdown:
		return true
	default:
		return false
	}
}

func (p *linuxProvider) Execute(ctx context.Context, action operations.Action) (operations.Result, error) {
	startedAt := time.Now().UTC()

	var cmd string
	switch action {
	case operations.ActionSleep:
		cmd = "systemctl"
	case operations.ActionRestart:
		cmd = "systemctl"
	case operations.ActionShutdown:
		cmd = "systemctl"
	default:
		return operations.Result{}, operations.ValidationError{Message: "unsupported action"}
	}

	var args []string
	switch action {
	case operations.ActionSleep:
		args = []string{"suspend"}
	case operations.ActionRestart:
		args = []string{"reboot"}
	case operations.ActionShutdown:
		args = []string{"poweroff"}
	}

	p.log.Info("executing operation: %s %v", cmd, args)

	err := p.runner.Run(ctx, cmd, args...)
	finishedAt := time.Now().UTC()

	result := operations.Result{
		Action:     action,
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
	}

	if err != nil {
		p.log.Error("operation failed: %v", err)
		result.Success = false
		result.Message = err.Error()
		return result, nil
	}

	result.Success = true
	result.Message = "operation completed"
	return result, nil
}
