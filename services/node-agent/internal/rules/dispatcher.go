package rules

import "context"

type Dispatcher interface {
	Dispatch(ctx context.Context, match Match) error
}
