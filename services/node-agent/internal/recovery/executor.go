package recovery

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/events"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/guardian"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
)

type Executor struct {
	guardian    *guardian.Service
	publish     func(context.Context, events.Event) error
	log         *logger.Logger
}

func NewExecutor(guardianService *guardian.Service, publish func(context.Context, events.Event) error, log *logger.Logger) *Executor {
	if log == nil {
		log = logger.New(logger.Info, nil)
	}
	return &Executor{guardian: guardianService, publish: publish, log: log}
}

func (e *Executor) ExecuteStep(ctx context.Context, step Step, target string) (bool, error) {
	switch step.Action {
	case RecoveryActionPing:
		return e.ping(ctx, target)
	case RecoveryActionPower:
		return false, e.power(ctx, guardian.PowerActionPress)
	case RecoveryActionReset:
		return false, e.reset(ctx, guardian.ResetActionPress)
	case RecoveryActionNotify:
		return false, e.notify(ctx, target)
	default:
		return false, fmt.Errorf("unsupported recovery action: %s", step.Action)
	}
}

func (e *Executor) ping(ctx context.Context, target string) (bool, error) {
	if target == "" {
		return false, fmt.Errorf("ping target is empty")
	}

	dialer := &net.Dialer{Timeout: 3 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", target)
	if err != nil {
		return false, nil
	}
	_ = conn.Close()
	return true, nil
}

func (e *Executor) power(ctx context.Context, action guardian.PowerAction) error {
	if e.guardian == nil {
		return fmt.Errorf("guardian service not configured")
	}
	return e.guardian.Power(ctx, action)
}

func (e *Executor) reset(ctx context.Context, action guardian.ResetAction) error {
	if e.guardian == nil {
		return fmt.Errorf("guardian service not configured")
	}
	return e.guardian.Reset(ctx, action)
}

func (e *Executor) notify(ctx context.Context, target string) error {
	if e.publish == nil {
		return fmt.Errorf("event publisher not configured")
	}
	_ = e.publish(ctx, events.SystemEvent("recovery_notify", "Recovery notification for target: "+target))
	return nil
}
