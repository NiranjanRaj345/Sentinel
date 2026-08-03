package nodes

import (
	"context"
	"time"
)

type Repository interface {
	Save(ctx context.Context, node Node) error
	UpdateLastSeen(ctx context.Context, id string, lastSeen time.Time) error
	UpdateStatus(ctx context.Context, id string, status Status) error
	Get(ctx context.Context, id string) (Node, error)
	List(ctx context.Context) ([]Node, error)
	Remove(ctx context.Context, id string) error
	Close() error
}
