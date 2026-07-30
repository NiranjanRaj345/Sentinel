package operations

import "context"

type Provider interface {
	Name() string
	Supports(action Action) bool
	Execute(ctx context.Context, action Action) (Result, error)
}
