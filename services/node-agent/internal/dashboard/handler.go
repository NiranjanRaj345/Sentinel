package dashboard

import (
	"encoding/json"
	"net/http"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
)

func OverviewHandler(service *Service, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		overview := service.Overview()

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(overview); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}
