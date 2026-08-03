package operations

import (
	"context"
	"errors"
	"testing"
)

func TestService_Execute_ValidationFailure(t *testing.T) {
	svc := NewService(nil, nil, NewValidator(nil), nil, nil, nil)
	_, err := svc.Execute(context.Background(), Request{Action: ActionRestart, Confirm: false})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestService_Execute_Success(t *testing.T) {
	provider := &stubProvider{supported: true}
	auditor := &stubAuditor{}
	svc := NewService(provider, auditor, NewValidator(provider), nil, nil, nil)

	result, err := svc.Execute(context.Background(), Request{Action: ActionRestart, Confirm: true})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !result.Success {
		t.Fatal("expected success")
	}
	if len(auditor.records) != 1 {
		t.Fatalf("expected 1 audit record, got %d", len(auditor.records))
	}
}

func TestService_Execute_ProviderFailure(t *testing.T) {
	provider := &stubProvider{supported: true, execErr: errors.New("boom")}
	svc := NewService(provider, NewAuditor(nil), NewValidator(provider), nil, nil, nil)

	_, err := svc.Execute(context.Background(), Request{Action: ActionRestart, Confirm: true})
	if err == nil {
		t.Fatal("expected provider error")
	}
}

type stubAuditor struct {
	records []Result
}

func (s *stubAuditor) Record(result Result) {
	s.records = append(s.records, result)
}
