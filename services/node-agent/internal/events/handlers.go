package events

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func EventsRecent(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		eventsList, err := service.Recent(100)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to load events: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(eventsList); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}
