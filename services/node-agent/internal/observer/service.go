package observer

import (
	"context"
	"fmt"
	"sync"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/events"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
)

type Service struct {
	client    *Client
	publish   func(context.Context, events.Event) error
	status    StatusResponse
	statusMu  sync.RWMutex
	env       EnvironmentResponse
	envMu     sync.RWMutex
	log       *logger.Logger
}

func NewService(client *Client, publish func(context.Context, events.Event) error, log *logger.Logger) *Service {
	if client == nil {
		client = NewClient("http://localhost:8082", log)
	}
	if log == nil {
		log = logger.New(logger.Info, nil)
	}
	return &Service{client: client, publish: publish, log: log}
}

func (s *Service) Status(ctx context.Context) (StatusResponse, error) {
	if s == nil || s.client == nil {
		return StatusResponse{}, fmt.Errorf("observer client not configured")
	}
	status, err := s.client.Status(ctx)
	if err != nil {
		return status, err
	}
	s.setStatus(status)
	return status, nil
}

func (s *Service) Environment(ctx context.Context) (EnvironmentResponse, error) {
	if s == nil || s.client == nil {
		return EnvironmentResponse{}, fmt.Errorf("observer client not configured")
	}
	env, err := s.client.Environment(ctx)
	if err != nil {
		return env, err
	}
	s.setEnvironment(env)
	return env, nil
}

func (s *Service) setStatus(status StatusResponse) {
	s.statusMu.Lock()
	s.status = status
	s.statusMu.Unlock()
}

func (s *Service) cachedStatus() StatusResponse {
	s.statusMu.RLock()
	defer s.statusMu.RUnlock()
	return s.status
}

func (s *Service) setEnvironment(env EnvironmentResponse) {
	s.envMu.Lock()
	s.env = env
	s.envMu.Unlock()
}

func (s *Service) cachedEnvironment() EnvironmentResponse {
	s.envMu.RLock()
	defer s.envMu.RUnlock()
	return s.env
}
