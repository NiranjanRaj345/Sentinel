package rules

import "context"

type Repository interface {
	Enabled(ctx context.Context) ([]Rule, error)
	Save(ctx context.Context, rule Rule) error
	Close() error
}
