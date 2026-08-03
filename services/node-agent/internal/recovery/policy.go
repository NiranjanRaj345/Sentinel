package recovery

import "time"

var DesktopRecoveryPolicy = Policy{
	ID:      "desktop_recovery",
	Name:    "Desktop Recovery",
	Enabled: true,
	Steps: []Step{
		{Action: RecoveryActionPing, Delay: 30 * time.Second, Retries: 3},
		{Action: RecoveryActionPower, Delay: 2 * time.Minute, Retries: 1},
		{Action: RecoveryActionPing, Delay: 30 * time.Second, Retries: 3},
		{Action: RecoveryActionReset, Delay: 30 * time.Second, Retries: 1},
		{Action: RecoveryActionNotify, Delay: 0, Retries: 1},
	},
}

func SeedPolicies() []Policy {
	return []Policy{DesktopRecoveryPolicy}
}
