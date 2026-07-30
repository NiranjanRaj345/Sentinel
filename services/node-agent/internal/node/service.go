package node

import "github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"

type Service struct {
	provider Provider
	log      *logger.Logger
}

func NewService(provider Provider, log *logger.Logger) *Service {
	if provider == nil {
		provider = noopProvider{}
	}
	return &Service{provider: provider, log: log}
}

func (s *Service) Capabilities() []CapabilityStatus {
	return s.provider.Capabilities()
}

type noopProvider struct{}

func (noopProvider) Name() string {
	return "unknown"
}

func (noopProvider) Capabilities() []CapabilityStatus {
	return nil
}
