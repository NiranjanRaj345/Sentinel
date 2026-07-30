package services

import "context"

type Provider interface {
	List(ctx context.Context) ([]ServiceItem, error)
	Execute(ctx context.Context, item ServiceItem) (ServiceItem, error)
}
