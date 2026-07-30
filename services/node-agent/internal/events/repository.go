package events

import "context"

type Repository interface {
	Save(ctx context.Context, event Event) error
	Recent(limit int) ([]Event, error)
	Close() error
}
