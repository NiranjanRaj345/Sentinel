package guardian

import (
	"context"
	"fmt"
	"sync"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/events"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/notification"
)

type Service struct {
	client    *Client
	publish   func(context.Context, events.Event) error
	notify    func(context.Context, notification.Notification)
	status    StatusResponse
	statusMu  sync.RWMutex
	log       *logger.Logger
}

func NewService(client *Client, publish func(context.Context, events.Event) error, notify func(context.Context, notification.Notification), log *logger.Logger) *Service {
	if client == nil {
		client = NewClient("http://localhost:8081", log)
	}
	if log == nil {
		log = logger.New(logger.Info, nil)
	}
	return &Service{client: client, publish: publish, notify: notify, log: log}
}

func (s *Service) Status(ctx context.Context) (StatusResponse, error) {
	if s == nil || s.client == nil {
		return StatusResponse{}, fmt.Errorf("guardian client not configured")
	}
	status, err := s.client.Status(ctx)
	if err != nil {
		return status, err
	}
	s.setStatus(status)
	return status, nil
}

func (s *Service) Power(ctx context.Context, action PowerAction) error {
	if s == nil || s.client == nil {
		return fmt.Errorf("guardian client not configured")
	}
	if err := s.client.Power(ctx, action); err != nil {
		return err
	}
	if s.publish != nil {
		_ = s.publish(ctx, events.SystemEvent("guardian_power", "Guardian power action: "+string(action)))
	}
	if s.notify != nil {
		severity := notification.SeverityInfo
		if action == PowerActionPress {
			severity = notification.SeverityWarning
		}
		s.notify(ctx, notification.NewNotification(
			"guardian-power-"+string(action),
			"Guardian power action",
			"Guardian power action: "+string(action),
			severity,
			notification.SourceGuardian,
		))
	}
	return nil
}

func (s *Service) Reset(ctx context.Context, action ResetAction) error {
	if s == nil || s.client == nil {
		return fmt.Errorf("guardian client not configured")
	}
	if err := s.client.Reset(ctx, action); err != nil {
		return err
	}
	if s.publish != nil {
		_ = s.publish(ctx, events.SystemEvent("guardian_reset", "Guardian reset action: "+string(action)))
	}
	if s.notify != nil {
		severity := notification.SeverityInfo
		if action == ResetActionPress {
			severity = notification.SeverityWarning
		}
		s.notify(ctx, notification.NewNotification(
			"guardian-reset-"+string(action),
			"Guardian reset action",
			"Guardian reset action: "+string(action),
			severity,
			notification.SourceGuardian,
		))
	}
	return nil
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
