package handlers

import (
	// "encoding/json"
	// "net/http"
	"time"
	// "github.com/NiranjanRaj345/sentinel/services/node-agent/internal/version"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Agent     AgentInfo `json:"agent"`
}

type AgentInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
