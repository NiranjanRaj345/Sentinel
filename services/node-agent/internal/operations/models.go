package operations

import "context"

type Action string

const (
	ActionSleep    Action = "sleep"
	ActionRestart  Action = "restart"
	ActionShutdown Action = "shutdown"
)

type Request struct {
	Action   Action `json:"action"`
	Confirm  bool   `json:"confirm"`
}

type Result struct {
	Action      Action        `json:"action"`
	Success     bool          `json:"success"`
	StartedAt   interface{}   `json:"startedAt"`
	FinishedAt  interface{}   `json:"finishedAt"`
	Message     string        `json:"message"`
}

type ValidationError struct {
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return e.Message
}

type CommandRunner interface {
	Run(ctx context.Context, name string, args ...string) error
}
