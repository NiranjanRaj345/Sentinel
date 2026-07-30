package linux

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/operations"
)

func TestLinuxProvider_SupportsKnownActions(t *testing.T) {
	provider := NewLinuxProvider(logger.New(0, io.Discard), operations.NewOSRunner())

	if !provider.Supports(operations.ActionSleep) {
		t.Error("expected sleep to be supported")
	}
	if !provider.Supports(operations.ActionRestart) {
		t.Error("expected restart to be supported")
	}
	if !provider.Supports(operations.ActionShutdown) {
		t.Error("expected shutdown to be supported")
	}
	if provider.Supports("unknown") {
		t.Error("expected unknown action to be unsupported")
	}
}

func TestLinuxProvider_Execute_UsesRunner(t *testing.T) {
	called := false
	var gotCmd string
	var gotArgs []string
	runner := &stubRunner{
		runFunc: func(ctx context.Context, name string, args ...string) error {
			called = true
			gotCmd = name
			gotArgs = args
			return nil
		},
	}

	provider := NewLinuxProvider(logger.New(0, io.Discard), runner)
	result, err := provider.Execute(context.Background(), operations.ActionRestart)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected runner to be called")
	}
	if gotCmd != "systemctl" {
		t.Fatalf("expected systemctl, got %s", gotCmd)
	}
	if len(gotArgs) != 1 || gotArgs[0] != "reboot" {
		t.Fatalf("expected [reboot], got %v", gotArgs)
	}
	if !result.Success {
		t.Fatal("expected success")
	}
}

func TestLinuxProvider_Execute_RunnerFailure(t *testing.T) {
	runner := &stubRunner{
		runFunc: func(ctx context.Context, name string, args ...string) error {
			return errors.New("runner failed")
		},
	}

	provider := NewLinuxProvider(logger.New(0, io.Discard), runner)
	result, err := provider.Execute(context.Background(), operations.ActionSleep)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Fatal("expected failure")
	}
	if result.Message != "runner failed" {
		t.Fatalf("expected error message, got %s", result.Message)
	}
}

func TestLinuxProvider_Execute_UnsupportedAction(t *testing.T) {
	provider := NewLinuxProvider(logger.New(0, io.Discard), operations.NewOSRunner())
	_, err := provider.Execute(context.Background(), "unknown")

	if err == nil {
		t.Fatal("expected error for unsupported action")
	}
}

type stubRunner struct {
	runFunc func(ctx context.Context, name string, args ...string) error
}

func (r *stubRunner) Run(ctx context.Context, name string, args ...string) error {
	return r.runFunc(ctx, name, args...)
}
