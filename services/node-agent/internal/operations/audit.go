package operations

import "github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"

type Auditor interface {
	Record(result Result)
}

type auditLogger struct {
	log *logger.Logger
}

func NewAuditor(log *logger.Logger) Auditor {
	if log == nil {
		return noopAuditor{}
	}
	return &auditLogger{log: log}
}

func (a *auditLogger) Record(result Result) {
	if a.log == nil {
		return
	}

	status := "success"
	if !result.Success {
		status = "failure"
	}

	a.log.Info("operation audit: action=%s status=%s message=%s", result.Action, status, result.Message)
}

type noopAuditor struct{}

func (noopAuditor) Record(result Result) {}
