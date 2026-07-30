package operations

import (
	"context"
	"testing"
	"time"
)

func TestValidator_UnknownAction(t *testing.T) {
	v := NewValidator(nil)
	err := v.Validate(Request{Action: "unknown", Confirm: true})
	if err == nil {
		t.Fatal("expected validation error for unknown action")
	}
}

func TestValidator_MissingConfirmation(t *testing.T) {
	v := NewValidator(nil)
	err := v.Validate(Request{Action: ActionRestart, Confirm: false})
	if err == nil {
		t.Fatal("expected validation error for missing confirmation")
	}
}

func TestValidator_UnsupportedAction(t *testing.T) {
	provider := &stubProvider{supported: false}
	v := NewValidator(provider)
	err := v.Validate(Request{Action: ActionRestart, Confirm: true})
	if err == nil {
		t.Fatal("expected validation error for unsupported action")
	}
}

func TestValidator_ValidRequest(t *testing.T) {
	provider := &stubProvider{supported: true}
	v := NewValidator(provider)
	err := v.Validate(Request{Action: ActionRestart, Confirm: true})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

type stubProvider struct {
	supported bool
	execErr   error
}

func (s *stubProvider) Name() string {
	return "stub"
}

func (s *stubProvider) Supports(action Action) bool {
	return s.supported
}

func (s *stubProvider) Execute(ctx context.Context, action Action) (Result, error) {
	if s.execErr != nil {
		return Result{}, s.execErr
	}
	return Result{Action: action, Success: true, StartedAt: time.Now().UTC(), FinishedAt: time.Now().UTC(), Message: "ok"}, nil
}
