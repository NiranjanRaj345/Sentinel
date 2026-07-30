package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/metrics"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/storage/sqlite"
)

type HistoryPoint struct {
	Timestamp    time.Time `json:"timestamp"`
	CPUUsage     float64   `json:"cpuUsage"`
	MemoryUsage  float64   `json:"memoryUsage"`
	DiskUsage    float64   `json:"diskUsage"`
}

type HistoryResponse struct {
	Period string         `json:"period"`
	Points []HistoryPoint `json:"points"`
}

type HistoryHandler struct {
	store *sqlite.Store
	log   *logger.Logger
}

func NewHistoryHandler(store *sqlite.Store, log *logger.Logger) *HistoryHandler {
	return &HistoryHandler{store: store, log: log}
}

func (h *HistoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		period = "1h"
	}

	to := time.Now().UTC()
	from, err := parsePeriod(period, to)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid period: %v", err), http.StatusBadRequest)
		return
	}

	snapshots, err := h.store.Range(from, to)
	if err != nil {
		h.log.Error("failed to load history: %v", err)
		http.Error(w, "failed to load history", http.StatusInternalServerError)
		return
	}

	points := make([]HistoryPoint, 0, len(snapshots))
	for _, snap := range snapshots {
		points = append(points, HistoryPoint{
			Timestamp:    snap.Metadata.Timestamp,
			CPUUsage:     snap.CPU.UsagePercent,
			MemoryUsage:  snap.Memory.UsagePercent,
			DiskUsage:    diskUsage(snap.Disks),
		})
	}

	response := HistoryResponse{
		Period: period,
		Points: points,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func parsePeriod(period string, to time.Time) (time.Time, error) {
	switch period {
	case "1h":
		return to.Add(-1 * time.Hour), nil
	case "24h":
		return to.Add(-24 * time.Hour), nil
	case "7d":
		return to.Add(-7 * 24 * time.Hour), nil
	default:
		d, err := time.ParseDuration(period)
		if err != nil {
			return time.Time{}, fmt.Errorf("unsupported period: %s", period)
		}
		return to.Add(-d), nil
	}
}

func diskUsage(disks []metrics.DiskInfo) float64 {
	if len(disks) == 0 {
		return 0
	}
	return disks[0].UsagePercent
}
