package automation

import (
	"context"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/rules"
)

type Service struct {
	engine    *Engine
	repo      Repository
	log       *logger.Logger
}

func NewService(engine *Engine, repo Repository, log *logger.Logger) *Service {
	if engine == nil {
		engine = NewEngine(nil, nil, log)
	}
	if log == nil {
		log = logger.New(logger.Info, nil)
	}
	return &Service{engine: engine, repo: repo, log: log}
}

func (s *Service) Dispatch(ctx context.Context, match rules.Match) error {
	if s == nil || s.engine == nil {
		return nil
	}
	if err := s.engine.Dispatch(ctx, match); err != nil {
		return err
	}
	if s.repo != nil {
		_ = s.repo.Save(ctx, ExecutionRecord{
			RuleID:    match.Rule.ID,
			RuleName:  match.Rule.Name,
			Action:    "evaluate",
			Success:   true,
			Message:   "dispatched",
			CreatedAt: time.Now().UTC(),
		})
	}
	return nil
}
