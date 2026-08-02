package resources

import "context"

type Provider interface {
	List(ctx context.Context) ([]Resource, error)
	Execute(ctx context.Context, action ResourceAction, name string) (Resource, error)
}
