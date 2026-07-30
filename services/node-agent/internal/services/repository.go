package services

import "context"

type Repository interface {
	List(ctx context.Context) ([]ServiceItem, error)
	Save(ctx context.Context, item ServiceItem) error
	Close() error
}
