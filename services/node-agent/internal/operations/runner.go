package operations

import (
	"context"
	"os/exec"
)

type OSRunner struct{}

func NewOSRunner() CommandRunner {
	return OSRunner{}
}

func (r OSRunner) Run(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	return cmd.Run()
}

type RecordingRunner struct {
	Calls []string
}

func NewRecordingRunner() *RecordingRunner {
	return &RecordingRunner{}
}

func (r *RecordingRunner) Run(ctx context.Context, name string, args ...string) error {
	r.Calls = append(r.Calls, name+" "+joinArgs(args))
	return nil
}

func joinArgs(args []string) string {
	result := ""
	for i, arg := range args {
		if i > 0 {
			result += " "
		}
		result += arg
	}
	return result
}
