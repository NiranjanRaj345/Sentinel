package notification

import "context"

type Dispatcher struct {
	service *Service
}

func NewDispatcher(service *Service) *Dispatcher {
	return &Dispatcher{service: service}
}

func (d *Dispatcher) Notify(ctx context.Context, notification Notification) {
	if d == nil || d.service == nil {
		return
	}
	d.service.Send(ctx, notification)
}
