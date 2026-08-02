package automation

import (
	"context"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/events"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/guardian"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/operations"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/rules"
)

type Engine struct {
	operationsService *operations.Service
	guardianService   *guardian.Service
	publish           func(context.Context, events.Event) error
	log               *logger.Logger
}

func NewEngine(operationsService *operations.Service, guardianService *guardian.Service, publish func(context.Context, events.Event) error, log *logger.Logger) *Engine {
	if log == nil {
		log = logger.New(logger.Info, nil)
	}
	return &Engine{
		operationsService: operationsService,
		guardianService:   guardianService,
		publish:           publish,
		log:               log,
	}
}

func (e *Engine) Dispatch(ctx context.Context, match rules.Match) error {
	if e == nil || e.operationsService == nil {
		return nil
	}

	for _, action := range match.Rule.Actions {
		switch action {
		case rules.ActionNotify:
			e.log.Info("automation notification: rule=%s event=%s", match.Rule.ID, match.Event.ID)
			if e.publish != nil {
				_ = e.publish(ctx, events.SystemEvent("rule_notification", "Rule matched: "+match.Rule.Name))
			}
		case rules.ActionExecute:
			e.executeOperation(ctx, match)
		case "guardian_power", "guardian_reset":
			e.executeGuardianAction(ctx, match, action)
		default:
			e.log.Warn("unsupported automation action: %s", action)
		}
	}

	return nil
}

func (e *Engine) executeGuardianAction(ctx context.Context, match rules.Match, action rules.Action) {
	if e.guardianService == nil {
		e.log.Error("automation guardian action skipped: guardian service not configured: rule=%s action=%s", match.Rule.ID, action)
		return
	}

	var guardianErr error
	switch action {
	case "guardian_power":
		guardianErr = e.guardianService.Power(ctx, guardian.PowerActionPress)
	case "guardian_reset":
		guardianErr = e.guardianService.Reset(ctx, guardian.ResetActionPress)
	}

	if guardianErr != nil {
		e.log.Error("automation guardian action failed: rule=%s action=%s err=%v", match.Rule.ID, action, guardianErr)
	}
}

func (e *Engine) executeOperation(ctx context.Context, match rules.Match) {
	result, err := e.operationsService.Execute(ctx, operations.Request{Action: operations.ActionRestart, Confirm: true})
	if err != nil {
		e.log.Error("automation execution failed: rule=%s err=%v", match.Rule.ID, err)
		return
	}

	if !result.Success {
		e.log.Error("automation execution unsuccessful: rule=%s message=%s", match.Rule.ID, result.Message)
	}
}
