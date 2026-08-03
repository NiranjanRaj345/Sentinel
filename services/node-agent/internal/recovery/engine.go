package recovery

import (
	"context"
	"fmt"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/events"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
)

type Engine struct {
	executor *Executor
	repo     Repository
	publish  func(context.Context, events.Event) error
	log      *logger.Logger
}

func NewEngine(executor *Executor, repo Repository, publish func(context.Context, events.Event) error, log *logger.Logger) *Engine {
	if executor == nil {
		executor = NewExecutor(nil, nil, log)
	}
	if log == nil {
		log = logger.New(logger.Info, nil)
	}
	return &Engine{executor: executor, repo: repo, publish: publish, log: log}
}

func (e *Engine) Execute(ctx context.Context, policy Policy, target string) (Execution, error) {
	execution := Execution{
		ID:        generateID(),
		PolicyID:  policy.ID,
		Status:    ExecutionStatusRunning,
		StartedAt: time.Now().UTC(),
	}

	for i, step := range policy.Steps {
		execution.Current = i + 1
		if err := e.saveExecution(ctx, execution); err != nil {
			e.log.Error("failed to save recovery execution: %v", err)
		}

		if step.Delay > 0 {
			select {
			case <-ctx.Done():
				execution.Status = ExecutionStatusFailed
				execution.Message = "recovery cancelled"
				finishExecution(&execution)
				return execution, ctx.Err()
			case <-time.After(step.Delay):
			}
		}

		success, stepErr := e.executeWithRetries(ctx, step, target)
		execution.Attempts++

		if stepErr != nil {
			execution.Status = ExecutionStatusFailed
			execution.Message = fmt.Sprintf("step %d failed: %v", i+1, stepErr)
			finishExecution(&execution)
			e.emitEvent(ctx, execution)
			return execution, stepErr
		}

		if !success {
			continue
		}

		execution.Status = ExecutionStatusSucceeded
		execution.Message = fmt.Sprintf("succeeded at step %d", i+1)
		finishExecution(&execution)
		e.emitEvent(ctx, execution)
		return execution, nil
	}

	execution.Status = ExecutionStatusFailed
	execution.Message = "recovery exhausted all steps without success"
	finishExecution(&execution)
	e.emitEvent(ctx, execution)
	return execution, fmt.Errorf("recovery policy exhausted: %s", policy.ID)
}

func (e *Engine) executeWithRetries(ctx context.Context, step Step, target string) (bool, error) {
	for attempt := 0; attempt <= step.Retries; attempt++ {
		if ctx.Err() != nil {
			return false, ctx.Err()
		}

		success, err := e.executor.ExecuteStep(ctx, step, target)
		if err != nil {
			return false, err
		}
		if success {
			return true, nil
		}
	}

	return false, nil
}

func finishExecution(execution *Execution) {
	now := time.Now().UTC()
	execution.FinishedAt = &now
}

func (e *Engine) saveExecution(ctx context.Context, execution Execution) error {
	if e.repo == nil {
		return nil
	}
	return e.repo.Save(ctx, execution)
}

func (e *Engine) emitEvent(ctx context.Context, execution Execution) {
	if e.publish == nil {
		return
	}

	severity := events.SeverityInfo
	if execution.Status == ExecutionStatusFailed {
		severity = events.SeverityCritical
	}

	_ = e.publish(ctx, events.Event{
		Type:     events.EventTypeSystem,
		Severity: severity,
		Source:   events.SourceResources,
		Title:    fmt.Sprintf("Recovery %s", execution.Status),
		Message:  execution.Message,
		Metadata: map[string]interface{}{"executionId": execution.ID, "policyId": execution.PolicyID},
	})
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UTC().UnixNano())
}
