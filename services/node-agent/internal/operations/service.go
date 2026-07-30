package operations

import (
	"context"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/events"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
)

type Service struct {
	provider  Provider
	auditor   Auditor
	validator Validator
	log       *logger.Logger
	publish   func(context.Context, events.Event) error
}

func NewService(provider Provider, auditor Auditor, validator Validator, publish func(context.Context, events.Event) error, log *logger.Logger) *Service {
	if provider == nil {
		provider = noopProvider{}
	}
	if auditor == nil {
		auditor = NewAuditor(log)
	}
	if validator == nil {
		validator = NewValidator(provider)
	}
	return &Service{
		provider:  provider,
		auditor:   auditor,
		validator: validator,
		publish:   publish,
		log:       log,
	}
}

func (s *Service) Execute(ctx context.Context, req Request) (Result, error) {
	if err := s.validator.Validate(req); err != nil {
		return Result{}, err
	}

	result, err := s.provider.Execute(ctx, req.Action)
	if err != nil {
		return Result{}, err
	}

	s.auditor.Record(result)

	if s.publish != nil {
		event := events.OperationFailure(string(req.Action), result.Message)
		if result.Success {
			event = events.OperationSuccess(string(req.Action), result.Message)
		}
		_ = s.publish(ctx, event)
	}

	return result, nil
}

type noopProvider struct{}

func (noopProvider) Name() string {
	return "noop"
}

func (noopProvider) Supports(action Action) bool {
	return false
}

func (noopProvider) Execute(ctx context.Context, action Action) (Result, error) {
	startedAt := time.Now().UTC()
	return Result{
		Action:     action,
		Success:    false,
		StartedAt:  startedAt,
		FinishedAt: startedAt,
		Message:    "no provider configured",
	}, nil
}
