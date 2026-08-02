package resources

import "context"

type Repository interface {
	List(ctx context.Context) ([]Resource, error)
	Save(ctx context.Context, resource Resource) error
	Close() error
}
