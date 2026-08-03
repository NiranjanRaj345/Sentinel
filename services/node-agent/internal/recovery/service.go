package recovery

import (
	"context"
	"fmt"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/events"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/notification"
)

type Service struct {
	engine       *Engine
	repo         Repository
	publish      func(context.Context, events.Event) error
	notify       func(context.Context, notification.Notification)
	activePolicy *Policy
	log          *logger.Logger
}

func NewService(engine *Engine, repo Repository, publish func(context.Context, events.Event) error, notify func(context.Context, notification.Notification), log *logger.Logger) *Service {
	if engine == nil {
		engine = NewEngine(nil, repo, publish, notify, log)
	}
	if log == nil {
		log = logger.New(logger.Info, nil)
	}
	return &Service{engine: engine, repo: repo, publish: publish, notify: notify, log: log}
}

func (s *Service) Execute(ctx context.Context, policyID, target string) (Execution, error) {
	if s == nil || s.engine == nil {
		return Execution{}, fmt.Errorf("recovery service not configured")
	}

	policy, err := s.getPolicy(policyID)
	if err != nil {
		return Execution{}, err
	}

	execution, err := s.engine.Execute(ctx, policy, target)
	if err != nil {
		return execution, err
	}

	if s.repo != nil {
		_ = s.repo.Save(ctx, execution)
	}

	return execution, nil
}

func (s *Service) Recent(limit int) ([]Execution, error) {
	if s == nil || s.repo == nil {
		return nil, fmt.Errorf("recovery repository not configured")
	}
	return s.repo.Recent(limit)
}

func (s *Service) Close() error {
	if s == nil || s.repo == nil {
		return nil
	}
	return s.repo.Close()
}

func (s *Service) getPolicy(id string) (Policy, error) {
	if s.activePolicy != nil && s.activePolicy.ID == id {
		return *s.activePolicy, nil
	}

	switch id {
	case DesktopRecoveryPolicy.ID:
		return DesktopRecoveryPolicy, nil
	default:
		return Policy{}, fmt.Errorf("unknown recovery policy: %s", id)
	}
}
