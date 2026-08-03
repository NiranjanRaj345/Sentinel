package notification

import (
	"context"
	"fmt"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/events"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
)

type Service struct {
	repo      Repository
	publish   func(context.Context, events.Event) error
	log       *logger.Logger
	providers []Provider
}

func NewService(repo Repository, publish func(context.Context, events.Event) error, log *logger.Logger, providers ...Provider) *Service {
	if log == nil {
		log = logger.New(logger.Info, nil)
	}
	return &Service{
		repo:      repo,
		publish:   publish,
		log:       log,
		providers: providers,
	}
}

func (s *Service) AddProvider(provider Provider) {
	if s == nil {
		return
	}
	s.providers = append(s.providers, provider)
}

func (s *Service) SetProviders(providers []Provider) {
	if s == nil {
		return
	}
	s.providers = providers
}

func (s *Service) Send(ctx context.Context, notification Notification) {
	if s == nil {
		return
	}

	if s.repo != nil {
		if err := s.repo.Save(ctx, notification); err != nil {
			s.log.Error("failed to persist notification: %v", err)
		}
	}

	if len(s.providers) == 0 {
		now := time.Now().UTC()
		_ = s.repo.UpdateStatus(ctx, notification.ID, StatusFailed, &now, "no providers configured")

		if s.publish != nil {
			_ = s.publish(ctx, events.Event{
				Type:     events.EventTypeSystem,
				Severity: events.SeverityWarning,
				Source:   SourceSystem,
				Title:    "Notification failed",
				Message:  fmt.Sprintf("No notification providers configured for: %s", notification.Title),
				Metadata: map[string]interface{}{"notificationId": notification.ID},
			})
		}
		return
	}

	var lastErr error
	var sentAt *time.Time

	for _, provider := range s.providers {
		if err := provider.Send(ctx, notification); err != nil {
			s.log.Error("provider %s failed to send notification: %v", provider.Name(), err)
			lastErr = err
			continue
		}

		now := time.Now().UTC()
		sentAt = &now
		_ = s.repo.UpdateStatus(ctx, notification.ID, StatusSent, sentAt, "")

		if s.publish != nil {
			_ = s.publish(ctx, events.Event{
				Type:     events.EventTypeSystem,
				Severity: events.SeverityInfo,
				Source:   SourceSystem,
				Title:    "Notification sent",
				Message:  fmt.Sprintf("Notification '%s' sent via %s", notification.Title, provider.Name()),
				Metadata: map[string]interface{}{"notificationId": notification.ID, "provider": provider.Name()},
			})
		}
		return
	}

	if lastErr != nil {
		now := time.Now().UTC()
		errMsg := ""
		if lastErr != nil {
			errMsg = lastErr.Error()
		}
		_ = s.repo.UpdateStatus(ctx, notification.ID, StatusFailed, &now, errMsg)

		if s.publish != nil {
			_ = s.publish(ctx, events.Event{
				Type:     events.EventTypeSystem,
				Severity: events.SeverityWarning,
				Source:   SourceSystem,
				Title:    "Notification failed",
				Message:  fmt.Sprintf("Failed to send notification '%s': %v", notification.Title, lastErr),
				Metadata: map[string]interface{}{"notificationId": notification.ID},
			})
		}
	}
}

func (s *Service) Recent(limit int) ([]Notification, error) {
	if s == nil || s.repo == nil {
		return nil, fmt.Errorf("notification service is not initialized")
	}
	return s.repo.Recent(limit)
}

func (s *Service) Close() error {
	if s == nil || s.repo == nil {
		return nil
	}
	return s.repo.Close()
}
