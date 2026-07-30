package operations

import (
	"context"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
)

type Service struct {
	provider  Provider
	auditor   Auditor
	validator Validator
	log       *logger.Logger
}

func NewService(provider Provider, auditor Auditor, validator Validator, log *logger.Logger) *Service {
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
		log:      log,
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
