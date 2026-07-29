package handlers

import (
	"encoding/json"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/metrics"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/version"
	"net/http"
	"time"
)

func Metrics(w http.ResponseWriter, r *http.Request) {
	info, err := metrics.GetInfo()
	if err != nil {
		http.Error(w, "failed to collect metrics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func Health(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Agent: AgentInfo{
			Name:    version.Build.Name,
			Version: version.Build.Version,
		},
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)
	}
}
