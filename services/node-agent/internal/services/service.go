package services

import (
	"context"
	"fmt"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
)

type Service struct {
	provider Provider
	repo     Repository
	log      *logger.Logger
}

func NewService(provider Provider, repo Repository, log *logger.Logger) *Service {
	if provider == nil {
		provider = noopProvider{}
	}
	if log == nil {
		log = logger.New(logger.Info, nil)
	}
	return &Service{provider: provider, repo: repo, log: log}
}

func (s *Service) List(ctx context.Context) ([]ServiceItem, error) {
	if s == nil || s.provider == nil {
		return nil, nil
	}
	return s.provider.List(ctx)
}

func (s *Service) Execute(ctx context.Context, action Action, name string) (ServiceItem, error) {
	if s == nil || s.provider == nil {
		return ServiceItem{}, fmt.Errorf("service provider not configured")
	}
	result, err := s.provider.Execute(ctx, ServiceItem{Name: name, Action: action, Status: ServiceStatusUnknown})
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

type noopProvider struct{}

func (noopProvider) List(ctx context.Context) ([]ServiceItem, error) { return nil, nil }
func (noopProvider) Execute(ctx context.Context, item ServiceItem) (ServiceItem, error) { return item, nil }
