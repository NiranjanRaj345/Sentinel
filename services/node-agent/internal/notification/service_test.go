package notification

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/events"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
)

type stubRepo struct {
	saved    []Notification
	recent   []Notification
	updated  []string
	closeErr error
}

func (s *stubRepo) Save(ctx context.Context, n Notification) error {
	s.saved = append(s.saved, n)
	return nil
}

func (s *stubRepo) Recent(limit int) ([]Notification, error) {
	return s.recent, nil
}

func (s *stubRepo) UpdateStatus(ctx context.Context, id string, status Status, sentAt *time.Time, errorMsg string) error {
	s.updated = append(s.updated, id)
	return nil
}

func (s *stubRepo) Close() error {
	return s.closeErr
}

type stubProvider struct {
	name    string
	sendErr error
}

func (p *stubProvider) Name() string {
	return p.name
}

func (p *stubProvider) Send(ctx context.Context, n Notification) error {
	return p.sendErr
}

func TestService_Send_NoProviders_MarksFailed(t *testing.T) {
	var published []events.Event
	repo := &stubRepo{}
	log := logger.New(logger.Info, io.Discard)
	svc := NewService(repo, func(ctx context.Context, e events.Event) error {
		published = append(published, e)
		return nil
	}, log)

	svc.Send(context.Background(), NewNotification("n1", "T", "M", SeverityInfo, SourceSystem))

	if len(repo.saved) != 1 {
		t.Fatalf("expected 1 saved notification, got %d", len(repo.saved))
	}
	if len(repo.updated) != 1 || repo.updated[0] != "n1" {
		t.Fatalf("expected status update for n1, got %v", repo.updated)
	}
	if len(published) != 1 {
		t.Fatalf("expected 1 published event, got %d", len(published))
	}
	if published[0].Title != "Notification failed" {
		t.Fatalf("expected failed event title, got %s", published[0].Title)
	}
}

func TestService_Send_ProviderSucceeds_MarksSent(t *testing.T) {
	var published []events.Event
	repo := &stubRepo{}
	log := logger.New(logger.Info, io.Discard)
	svc := NewService(repo, func(ctx context.Context, e events.Event) error {
		published = append(published, e)
		return nil
	}, log, &stubProvider{name: "test"})

	svc.Send(context.Background(), NewNotification("n2", "T", "M", SeverityInfo, SourceSystem))

	if len(repo.updated) != 1 || repo.updated[0] != "n2" {
		t.Fatalf("expected status update for n2, got %v", repo.updated)
	}
	if len(published) != 1 {
		t.Fatalf("expected 1 published event, got %d", len(published))
	}
	if published[0].Title != "Notification sent" {
		t.Fatalf("expected sent event title, got %s", published[0].Title)
	}
}

func TestService_Send_ProviderFails_MarksFailed(t *testing.T) {
	var published []events.Event
	repo := &stubRepo{}
	log := logger.New(logger.Info, io.Discard)
	svc := NewService(repo, func(ctx context.Context, e events.Event) error {
		published = append(published, e)
		return nil
	}, log, &stubProvider{name: "test", sendErr: errors.New("boom")})

	svc.Send(context.Background(), NewNotification("n3", "T", "M", SeverityInfo, SourceSystem))

	if len(repo.updated) != 1 || repo.updated[0] != "n3" {
		t.Fatalf("expected status update for n3, got %v", repo.updated)
	}
	if len(published) != 1 {
		t.Fatalf("expected 1 published event, got %d", len(published))
	}
	if published[0].Title != "Notification failed" {
		t.Fatalf("expected failed event title, got %s", published[0].Title)
	}
}

func TestService_Recent_ReturnsRepoResults(t *testing.T) {
	repo := &stubRepo{recent: []Notification{{ID: "r1"}}}
	svc := NewService(repo, nil, nil)

	recent, err := svc.Recent(5)
	if err != nil {
		t.Fatalf("recent: %v", err)
	}
	if len(recent) != 1 || recent[0].ID != "r1" {
		t.Fatalf("expected recent with r1, got %v", recent)
	}
}

func TestService_Close_NilRepo_Noop(t *testing.T) {
	svc := NewService(nil, nil, nil)
	if err := svc.Close(); err != nil {
		t.Fatalf("expected nil error for nil repo, got %v", err)
	}
}

func TestService_NilService_Noop(t *testing.T) {
	var svc *Service
	svc.Send(context.Background(), Notification{})
	_, _ = svc.Recent(1)
	_ = svc.Close()
}
