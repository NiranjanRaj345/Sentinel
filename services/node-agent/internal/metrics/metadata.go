package metrics

import "time"

type AgentInfo struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Platform     string `json:"platform"`
	Architecture string `json:"architecture"`
	GoVersion    string `json:"goVersion"`
}

type Metadata struct {
	Timestamp            time.Time `json:"timestamp"`
	CollectionDurationMS int64     `json:"collectionDurationMs"`
	Agent                AgentInfo `json:"agent"`
}
