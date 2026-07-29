package metrics

import "time"

// AgentInfo describes the running Sentinel agent.
type AgentInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Metadata describes the metrics collection.
type Metadata struct {
	Timestamp            time.Time `json:"timestamp"`
	CollectionDurationMS int64     `json:"collection_duration_ms"`
	Agent                AgentInfo `json:"agent"`
}
