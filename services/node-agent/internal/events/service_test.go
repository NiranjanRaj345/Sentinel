package events

import (
	"context"
	"errors"
	"testing"
)

type fakeRepo struct {
	events []Event
	err    error
}

func (f *fakeRepo) Save(ctx context.Context, event Event) error {
	if f.err != nil {
		return f.err
	}
	f.events = append(f.events, event)
	return nil
}

func (f *fakeRepo) Recent(limit int) ([]Event, error) {
	if f.err != nil {
		return nil, f.err
	}
	if limit <= 0 {
		limit = len(f.events)
	}
	start := 0
	if len(f.events) > limit {
		start = len(f.events) - limit
	}
	return f.events[start:], nil
}

func (f *fakeRepo) Close() error {
	return nil
}

type fakePublisher struct {
	events []Event
	err    error
}

func (f *fakePublisher) Publish(ctx context.Context, event Event) error {
	if f.err != nil {
		return f.err
	}
	f.events = append(f.events, event)
	return nil
}

func TestService_Publish_AssignsIDAndTimestamp(t *testing.T) {
	repo := &fakeRepo{}
	pub := &fakePublisher{}
	svc := NewService(repo, pub, nil)

	if err := svc.Publish(context.Background(), SystemEvent("test", "message")); err != nil {
		t.Fatalf("Publish() failed: %v", err)
	}

	if len(repo.events) != 1 {
		t.Fatalf("expected 1 persisted event, got %d", len(repo.events))
	}

	persisted := repo.events[0]
	if persisted.ID == "" {
		t.Fatal("expected event ID to be assigned")
	}
	if persisted.CreatedAt.IsZero() {
		t.Fatal("expected event timestamp to be assigned")
	}
}

func TestService_Publish_PersistsAndPublishes(t *testing.T) {
	repo := &fakeRepo{}
	pub := &fakePublisher{}
	svc := NewService(repo, pub, nil)

	event := SystemEvent("test", "message")
	if err := svc.Publish(context.Background(), event); err != nil {
		t.Fatalf("Publish() failed: %v", err)
	}

	if len(repo.events) != 1 {
		t.Fatalf("expected 1 persisted event, got %d", len(repo.events))
	}
	if len(pub.events) != 1 {
		t.Fatalf("expected 1 published event, got %d", len(pub.events))
	}
}

func TestService_Publish_NilService_NoOp(t *testing.T) {
	var svc *Service
	if err := svc.Publish(context.Background(), Event{}); err != nil {
		t.Fatalf("expected no error from nil service, got %v", err)
	}
}

func TestService_Publish_NilRepo_StillPublishes(t *testing.T) {
	pub := &fakePublisher{}
	svc := NewService(nil, pub, nil)

	event := SystemEvent("test", "message")
	if err := svc.Publish(context.Background(), event); err != nil {
		t.Fatalf("Publish() failed: %v", err)
	}

	if len(pub.events) != 1 {
		t.Fatalf("expected 1 published event, got %d", len(pub.events))
	}
}

func TestService_Recent_DelegatesToRepo(t *testing.T) {
	repo := &fakeRepo{}
	svc := NewService(repo, nil, nil)

	for i := 0; i < 3; i++ {
		svc.Publish(context.Background(), SystemEvent("t", "m"))
	}

	events, err := svc.Recent(2)
	if err != nil {
		t.Fatalf("Recent() failed: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
}

func TestService_Recent_NoRepo_ReturnsError(t *testing.T) {
	svc := NewService(nil, nil, nil)
	_, err := svc.Recent(10)
	if err == nil {
		t.Fatal("expected error when repo is nil")
	}
}

func TestService_Publish_RepoFailure_ContinuesToPublish(t *testing.T) {
	repo := &fakeRepo{err: errors.New("db full")}
	pub := &fakePublisher{}
	svc := NewService(repo, pub, nil)

	event := SystemEvent("test", "message")
	if err := svc.Publish(context.Background(), event); err != nil {
		t.Fatalf("Publish() should not fail: %v", err)
	}

	if len(pub.events) != 1 {
		t.Fatalf("expected 1 published event even when repo fails, got %d", len(pub.events))
	}
}

func TestService_Publish_PublisherFailure_ContinuesPersist(t *testing.T) {
	repo := &fakeRepo{}
	pub := &fakePublisher{err: errors.New("broker down")}
	svc := NewService(repo, pub, nil)

	event := SystemEvent("test", "message")
	if err := svc.Publish(context.Background(), event); err != nil {
		t.Fatalf("Publish() should not fail: %v", err)
	}

	if len(repo.events) != 1 {
		t.Fatalf("expected 1 persisted event even when publisher fails, got %d", len(repo.events))
	}
}
