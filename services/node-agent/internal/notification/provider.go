package notification

import "context"

type Provider interface {
	Name() string

	Send(ctx context.Context, notification Notification) error
}
