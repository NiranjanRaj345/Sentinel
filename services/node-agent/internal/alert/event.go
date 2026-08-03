package alert

import "time"

type Event struct {
	RuleID   	string    `json:"ruleId"`
	RuleName 	string    `json:"ruleName"`
	Metric   	string    `json:"metric"`
	Value    	float64   `json:"value"`
	Threshold 	float64   `json:"threshold"`
	Severity 	Severity  `json:"severity"`
	TriggeredAt time.Time `json:"triggeredAt"`
}
