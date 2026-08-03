package notification

import (
	"context"
	"time"
)

type Repository interface {
	Save(ctx context.Context, notification Notification) error
	Recent(limit int) ([]Notification, error)
	UpdateStatus(ctx context.Context, id string, status Status, sentAt *time.Time, errorMsg string) error
	Close() error
}
