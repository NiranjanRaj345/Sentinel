package resources

import (
	"context"
	"fmt"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/events"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
)

type Service struct {
	provider       Provider
	repo           Repository
	publish        func(context.Context, events.Event) error
	evaluator      *HealthEvaluator
	log            *logger.Logger
	previousHealth map[string]Health
}

func NewService(provider Provider, repo Repository, publish func(context.Context, events.Event) error, log *logger.Logger) *Service {
	if provider == nil {
		provider = noopProvider{}
	}
	if log == nil {
		log = logger.New(logger.Info, nil)
	}
	return &Service{
		provider:       provider,
		repo:           repo,
		publish:        publish,
		evaluator:      NewHealthEvaluator(),
		log:            log,
		previousHealth: make(map[string]Health),
	}
}

func (s *Service) List(ctx context.Context) ([]Resource, error) {
	if s == nil || s.provider == nil {
		return nil, nil
	}
	resources, err := s.provider.List(ctx)
	if err != nil {
		return nil, err
	}

	var result []Resource
	for _, resource := range resources {
		s.emitHealthChangeIfNeeded(resource)
		result = append(result, resource)
	}
	return result, nil
}

func (s *Service) Execute(ctx context.Context, action ResourceAction, name string) (Resource, error) {
	if s == nil || s.provider == nil {
		return Resource{}, fmt.Errorf("resource provider not configured")
	}
	result, err := s.provider.Execute(ctx, action, name)
	if err != nil {
		return result, err
	}
	if s.repo != nil {
		_ = s.repo.Save(ctx, result)
	}
	return result, nil
}

func (s *Service) Close() error {
	if s == nil || s.repo == nil {
		return nil
	}
	return s.repo.Close()
}

func (s *Service) emitHealthChangeIfNeeded(resource Resource) {
	if s.publish == nil {
		return
	}
	prev, ok := s.previousHealth[resource.Name]
	if !ok || prev != resource.Health {
		_ = s.publish(context.Background(), events.ResourceHealthChanged(resource.Name, string(resource.Health), resource.Message))
		s.previousHealth[resource.Name] = resource.Health
	}
}

type noopProvider struct{}

func (noopProvider) List(ctx context.Context) ([]Resource, error) { return nil, nil }
func (noopProvider) Execute(ctx context.Context, action ResourceAction, name string) (Resource, error) { return Resource{}, nil }
