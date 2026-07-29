package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/storage/sqlite"
)

func HistoryLatest(store *sqlite.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info, err := store.Latest()
		if err != nil {
			if err.Error() == "no historical metrics available" {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, fmt.Sprintf("failed to load latest metrics: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(info); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}

func HistoryRange(store *sqlite.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fromStr := r.URL.Query().Get("from")
		toStr := r.URL.Query().Get("to")

		if fromStr == "" || toStr == "" {
			http.Error(w, "missing 'from' or 'to' parameter", http.StatusBadRequest)
			return
		}

		from, err := time.Parse(time.RFC3339, fromStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid 'from' parameter: %v", err), http.StatusBadRequest)
			return
		}

		to, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid 'to' parameter: %v", err), http.StatusBadRequest)
			return
		}

		if from.After(to) {
			http.Error(w, "invalid range: 'from' must not be after 'to'", http.StatusBadRequest)
			return
		}

		results, err := store.Range(from, to)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to load history: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(results); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}
