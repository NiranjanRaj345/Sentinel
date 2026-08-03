package recovery

import "time"

type RecoveryAction string

const (
	RecoveryActionPing   RecoveryAction = "ping"
	RecoveryActionPower  RecoveryAction = "power"
	RecoveryActionReset  RecoveryAction = "reset"
	RecoveryActionNotify RecoveryAction = "notify"
)

type Step struct {
	Action  RecoveryAction `json:"action"`
	Delay   time.Duration  `json:"delay"`
	Retries int            `json:"retries"`
}

type Policy struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Steps   []Step `json:"steps"`
}

type ExecutionStatus string

const (
	ExecutionStatusRunning  ExecutionStatus = "running"
	ExecutionStatusSucceeded ExecutionStatus = "succeeded"
	ExecutionStatusFailed   ExecutionStatus = "failed"
)

type Execution struct {
	ID        string            `json:"id"`
	PolicyID  string            `json:"policyId"`
	Status    ExecutionStatus   `json:"status"`
	Current   int               `json:"current"`
	Attempts  int               `json:"attempts"`
	Message   string            `json:"message"`
	StartedAt time.Time         `json:"startedAt"`
	FinishedAt *time.Time       `json:"finishedAt"`
}
