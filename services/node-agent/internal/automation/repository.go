package automation

import "context"

type Repository interface {
	Save(ctx context.Context, record ExecutionRecord) error
	List(ctx context.Context, limit int) ([]ExecutionRecord, error)
	Close() error
}

type ExecutionRecord struct {
	ID        string `json:"id"`
	RuleID    string `json:"ruleId"`
	RuleName  string `json:"ruleName"`
	Action    string `json:"action"`
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	CreatedAt interface{} `json:"createdAt"`
}
