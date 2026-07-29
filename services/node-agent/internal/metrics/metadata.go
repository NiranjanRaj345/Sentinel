package metrics

import "time"

type AgentInfo struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Platform     string `json:"platform"`
	Architecture string `json:"architecture"`
	GoVersion    string `json:"go_version"`
}

type Metadata struct {
	Timestamp            time.Time `json:"timestamp"`
	CollectionDurationMS int64     `json:"collection_duration_ms"`
	Agent                AgentInfo `json:"agent"`
}
