package rules

import (
	"context"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
)

type Service struct {
	repo      Repository
	dispatcher Dispatcher
	evaluator *Evaluator
	log       *logger.Logger
}

func NewService(repo Repository, dispatcher Dispatcher, evaluator *Evaluator, log *logger.Logger) *Service {
	if evaluator == nil {
		evaluator = NewEvaluator()
	}
	if log == nil {
		log = logger.New(logger.Info, nil)
	}
	return &Service{
		repo:      repo,
		dispatcher: dispatcher,
		evaluator: evaluator,
		log:       log,
	}
}

func (s *Service) Evaluate(ctx context.Context, event Event) []Match {
	if s == nil || s.repo == nil {
		return nil
	}
	rules, err := s.repo.Enabled(ctx)
	if err != nil {
		s.log.Error("failed to load rules: %v", err)
		return nil
	}
	var matches []Match
	for _, rule := range rules {
		if s.evaluator.Evaluate(event, rule) {
			matches = append(matches, Match{Rule: rule, Event: event})
			if s.dispatcher != nil {
				if err := s.dispatcher.Dispatch(ctx, Match{Rule: rule, Event: event}); err != nil {
					s.log.Error("failed to dispatch rule %s: %v", rule.ID, err)
				}
			}
		}
	}
	return matches
}

func (s *Service) Close() error {
	if s == nil || s.repo == nil {
		return nil
	}
	return s.repo.Close()
}

func SeedRules() []Rule {
	return []Rule{
		{
			ID:       "critical-alert-notify",
			Name:     "Critical Alert Notification",
			Enabled:  true,
			Trigger:  TriggerEvent,
			Actions:  []Action{ActionNotify},
			Conditions: []Condition{
				{Field: "severity", Operator: OpEquals, Value: "critical"},
			},
		},
		{
			ID:       "operation-failure-notify",
			Name:     "Operation Failure Notification",
			Enabled:  true,
			Trigger:  TriggerEvent,
			Actions:  []Action{ActionNotify},
			Conditions: []Condition{
				{Field: "type", Operator: OpEquals, Value: "operation"},
				{Field: "severity", Operator: OpEquals, Value: "warning"},
			},
		},
	}
}
