package events

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
)

type Service struct {
	repo      Repository
	publisher Publisher
	log       *logger.Logger
}

func NewService(repo Repository, publisher Publisher, log *logger.Logger) *Service {
	if log == nil {
		log = logger.New(logger.Info, os.Stderr)
	}
	return &Service{
		repo:      repo,
		publisher: publisher,
		log:       log,
	}
}

func (s *Service) Publish(ctx context.Context, event Event) error {
	if s == nil {
		return nil
	}

	if event.Metadata == nil {
		event.Metadata = make(map[string]interface{})
	}

	if err := validate(event); err != nil {
		return fmt.Errorf("validate event: %w", err)
	}

	event.ID = generateID()
	event.CreatedAt = time.Now().UTC()

	if s.repo != nil {
		if err := s.repo.Save(ctx, event); err != nil {
			s.log.Error("failed to persist event: %v", err)
		}
	}

	if s.publisher != nil {
		if err := s.publisher.Publish(ctx, event); err != nil {
			s.log.Error("failed to publish event: %v", err)
		}
	}

	return nil
}

func (s *Service) Recent(limit int) ([]Event, error) {
	if s == nil || s.repo == nil {
		return nil, fmt.Errorf("events repository not configured")
	}
	return s.repo.Recent(limit)
}

func generateID() string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

func validate(event Event) error {
	if event.ID != "" {
		return fmt.Errorf("event id must not be set")
	}
	if event.Type == "" {
		return fmt.Errorf("event type is required")
	}
	if event.Severity == "" {
		return fmt.Errorf("event severity is required")
	}
	if event.Source == "" {
		return fmt.Errorf("event source is required")
	}
	if event.Title == "" {
		return fmt.Errorf("event title is required")
	}
	if !event.CreatedAt.IsZero() {
		return fmt.Errorf("event timestamp must not be set")
	}
	return nil
}
